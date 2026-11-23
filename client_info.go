package main

import (
	"log"
	"os"

	"github.com/redsms/api-samples-go/internal/redsms"
)

// Keep the original commands available for existing users.
func main() {
	args := os.Args[1:]
	switch {
	case len(args) == 0:
		redsms.ClientInfo(args)
	case args[0] == "send":
		redsms.SendSMS(args[1:])
	case args[0] == "status":
		redsms.MessageStatus(args[1:])
	default:
		log.Fatal("Использование: go run client_info.go [send | status UUID]")
	}
}
