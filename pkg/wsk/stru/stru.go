package stru

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// โครงสร้างสำหรับเก็บ client
type Client struct {
	Conn *websocket.Conn
	ID   string
}

type Request struct {
	TranType string          `json:"tranType"`           // request, response
	Room     string          `json:"room,omitempty"`     // ชื่อห้อง
	Type     string          `json:"type"`               // command, text, image, confirm, etc.
	Action   string          `json:"action"`             // คำสั่งที่ต้องการให้ client ทำ
	Message  string          `json:"message"`            // กรณี tranType=request  มันคือข้อความ (ใช้กับ text/image), กรณี tranType=response มันคือข้อความประกอบสถานะ เช่น success หรือ error...(ข้อความเออเร่อ)...
	Selector string          `json:"selector,omitempty"` // กรณี single
	TabID    int             `json:"tab_id,omitempty"`   // id ของแทบที่ user เปิดใช้งาน side panel อยู่
	Data     json.RawMessage `json:"data,omitempty"`     // รองรับ dynamic structure กรณี tranType=request คือมีค่าส่งเข้ามาเพิ่มเติมในนี้, tranType=response คือค่าที่ส่งคืนกลับเซิฟเวอร์
}

type FormField struct {
	Selector string `json:"selector"`
	Value    string `json:"value"`
}

// ฟังก์ชันสำหรับส่งคำสั่งไปยัง client ใน room ที่ต้องการ
func (req *Request) SendToRoom(userID string, client *Client) (err error) {
	msg, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = client.Conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Printf("Error sending to %s: %v\n", userID, err)
		return err
	}

	log.Printf("✅ Sent command to %s \n", userID)
	return nil
}
