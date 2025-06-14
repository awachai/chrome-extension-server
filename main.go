package main

import (
	"ai-web-agent-server/pkg/auth"
	"ai-web-agent-server/pkg/mdware"
	"ai-web-agent-server/pkg/wsk"
	httphandle "ai-web-agent-server/pkg/wsk/http_handle"
	"log"
	"net/http"
)

func main() {
	// ✅ Public routes
	http.HandleFunc("/auth/login", mdware.CorsMiddleware(auth.LoginHandler))
	http.HandleFunc("/auth/logout", auth.LogoutHandler)
	http.HandleFunc("/auth/getprofile", auth.AuthMiddleware(auth.GetProfileHandler))

	http.HandleFunc("/ws", auth.AuthMiddleware(wsk.HandleWebSocket)) //auth.AuthMiddleware(

	// Test endpoint for HTTP requests
	http.HandleFunc("/test/click", httphandle.TestClick)
	http.HandleFunc("/test/url", httphandle.TestURL)

	log.Println("🚀 WebSocket server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
