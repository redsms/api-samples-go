package main

import (
	"os"

	"github.com/redsms/api-samples-go/internal/redsms"
)

func main() {
	redsms.ClientInfo(os.Args[1:])
}
