package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// 🔐 Key type ต้องเป็นชนิดเดียวกัน
type key string

const userIDKey key = "userID"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Room  string `json:"room"`
}

var mockTokens = map[string]string{ // username → token
	"nueng": "user123",
	"admin": "admin456",
}

// === Rate Limit ===
var loginAttempts = make(map[string]int)
var lastAttempt = make(map[string]int64)
var mu sync.Mutex

const maxAttempts = 5
const blockDuration = 60 // seconds

func isRateLimited(ip string) bool {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now().Unix()
	if last, ok := lastAttempt[ip]; ok {
		if now-last > blockDuration {
			loginAttempts[ip] = 0
		}
	}
	lastAttempt[ip] = now

	if loginAttempts[ip] >= maxAttempts {
		return true
	}
	loginAttempts[ip]++
	return false
}

// POST /auth/login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	/*fmt.Println("Login..")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	fmt.Println("body : ", string(body))*/

	ip := strings.Split(r.RemoteAddr, ":")[0]
	if isRateLimited(ip) {
		http.Error(w, "Too many login attempts. Please try again later.", http.StatusTooManyRequests)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, ok := mockTokens[req.Username]
	if !ok || req.Password != "password" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	res := LoginResponse{
		Token: token,
		Room:  "room_" + req.Username,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// POST /auth/logout (mock)
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logged out (mock)"))
}

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username := reverseLookupUsername(userID)

	profile := map[string]interface{}{
		"userID":   userID,
		"username": username,
		"role":     "agent",
		"room":     "room_" + username,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// ใช้ lookup กลับจาก token เป็น username (mock logic)
func reverseLookupUsername(token string) string {
	for user, tk := range mockTokens {
		if tk == token {
			return user
		}
	}
	return "unknown"
}

func GetUserID(r *http.Request) string {
	if val := r.Context().Value(userIDKey); val != nil {
		return val.(string)
	}
	return ""
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		fmt.Println("token1 : ", token)

		if token == "" {
			token = r.URL.Query().Get("authorization")
		}

		fmt.Println("token2 : ", token)

		// ✅ รองรับทั้ง Bearer token และ token ธรรมดา
		token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))

		fmt.Println("token3 : ", token)

		if token == "" {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		userID, err := ValidateToken(token)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func ValidateToken(token string) (string, error) {
	// สมมุติว่า token = "user123" ถึงจะผ่าน
	if token == "user123" {
		return "user123", nil
	}
	return "", errors.New("invalid token")
}
