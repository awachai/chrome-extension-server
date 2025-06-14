package wsk

import (
	"ai-web-agent-server/pkg/wsk/stru"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestStru(t *testing.T) {
	req := stru.Request{}
	req.TranType = "request"
	req.Type = "commard"
	req.Action = "get_dom"

	jsonSTR, err := json.Marshal(req)
	if err != nil {
		log.Println("❌ Error:", err)
	}

	fmt.Println(string(jsonSTR))

}

func TestClick1(t *testing.T) {
	err := SendClickCommandTo("nueng", "#fQuote > div:nth-child(4) > button")
	if err != nil {
		log.Println("❌ Error:", err)
	}
	fmt.Println("✅ Click command sent to #yearcar for user nueng")
}

func TestClick(w http.ResponseWriter, r *http.Request) {
	err := SendClickCommandTo("user123", "#fQuote > div:nth-child(4) > button")
	if err != nil {
		log.Println("❌ Error:", err)
	}
	w.Write([]byte("คลิกที่ #yearcar ถูกส่งไปแล้ว"))
}
