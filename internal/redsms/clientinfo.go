package redsms

import (
	"log"
	"net/http"
)

func ClientInfo(args []string) {
	if len(args) != 0 {
		log.Fatal("Использование: go run ./cmd/client-info")
	}
	request(http.MethodGet, "client/info", nil)
}
