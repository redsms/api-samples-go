package main

import (
	"os"

	"github.com/redsms/api-samples-go/internal/redsms"
)

func main() {
	redsms.SendSMS(os.Args[1:])
}
