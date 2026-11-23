package redsms

import (
	"log"
	"net/http"
	"net/url"
)

func SendSMS(args []string) {
	if len(args) != 0 {
		log.Fatal("Использование: go run ./cmd/send-sms")
	}
	fields := url.Values{}
	fields.Set("route", "sms")
	fields.Set("from", required("REDSMS_FROM"))
	fields.Set("to", required("REDSMS_TO"))
	fields.Set("text", "Тестовое сообщение от redsms.ru")
	request(http.MethodPost, "message", fields)
}
