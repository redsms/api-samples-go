package main

import (
	"os"

	"github.com/redsms/api-samples-go/internal/redsms"
)

func main() {
	redsms.MessageStatus(os.Args[1:])
}
