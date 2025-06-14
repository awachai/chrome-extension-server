package httphandle

import (
	"ai-web-agent-server/pkg/wsk"
	"log"
	"net/http"
)

func TestClick(w http.ResponseWriter, r *http.Request) {
	err := wsk.SendClickCommandTo("user123", "#fQuote > div:nth-child(4) > button")
	if err != nil {
		log.Println("❌ Error:", err)
	}
	w.Write([]byte("คลิกที่ #yearcar ถูกส่งไปแล้ว"))
}

func TestURL(w http.ResponseWriter, r *http.Request) {
	err := wsk.SendOpenURLCommand("user123", "https://www.google.com")
	if err != nil {
		log.Println("❌ Error:", err)
	}
	w.Write([]byte("url google ถูกส่งไปแล้ว"))
}
