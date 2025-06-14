package wsk

import (
	"ai-web-agent-server/pkg/auth"
	"ai-web-agent-server/pkg/wsk/stru"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// ใช้สำหรับอัปเกรด HTTP เป็น WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // เปิดให้ทุก origin (ปรับตามความปลอดภัยที่ต้องการ)
	},
}

// เก็บ client ทั้งหมด
var clients = make(map[string]*stru.Client)
var mutex = &sync.Mutex{}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r) // ดึงจาก context

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade error:", err)
		return
	}

	// ดึง userID จาก query เช่น ws://localhost:8080/ws?user=nueng
	if userID == "" {
		userID = conn.RemoteAddr().String()
	}

	client := &stru.Client{Conn: conn, ID: userID}
	mutex.Lock()
	clients[userID] = client
	mutex.Unlock()

	log.Printf("✅ Client connected: %s\n", userID)

	for {
		var msg interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("❌ Error from %s: %v\n", userID, err)
			break
		}
		log.Printf("📨 Message from %s: %v\n", userID, msg)

		// ตัวอย่าง: ส่งกลับข้อความเดิม
		/*response := map[string]interface{}{
			"type":    "echo",
			"from":    userID,
			"message": msg,
		}
		conn.WriteJSON(response)*/
	}

	mutex.Lock()
	delete(clients, userID)
	mutex.Unlock()
	conn.Close()
	log.Printf("❌ Client disconnected: %s\n", userID)
}

// ส่งคำสั่ง click ไปยัง client ตาม userID
func SendClickCommandTo(userID string, selector string) error {
	client, ok := clients[userID]
	if !ok {
		log.Printf("Client %s not found\n", userID)
		return nil
	}

	req := stru.Request{
		Type:     "command",
		Action:   "click",
		Selector: selector,
	}

	err := req.SendToRoom(userID, client)
	if err != nil {
		log.Printf("❌ Error sending click command to %s: %v\n", userID, err)
		return err
	}

	log.Printf("✅ Sent click command to %s (selector: %s)\n", userID, selector)
	return nil
}

func SendOpenURLCommand(userID string, url string) error {
	client, ok := clients[userID]
	if !ok {
		log.Printf("Client %s not found\n", userID)
		return nil
	}
	// เตรียม payload สำหรับ data field
	payload := map[string]string{
		"url": url,
	}
	rawData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal data: %v", err)
		return err
	}

	req := stru.Request{
		TranType: "request",
		Room:     userID,
		Type:     "command",
		Action:   "open_url",
		Message:  "เปิดหน้า URL ให้ผู้ใช้งาน",
		Data:     rawData,
	}

	return req.SendToRoom(userID, client)
}
