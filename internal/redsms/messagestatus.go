package redsms

import (
	"log"
	"net/http"
	"net/url"
	"regexp"
)

var uuidPattern = regexp.MustCompile("(?i)^[0-9a-f]{8}(-[0-9a-f]{4}){3}-[0-9a-f]{12}$")

func MessageStatus(args []string) {
	if len(args) != 1 || !uuidPattern.MatchString(args[0]) {
		log.Fatal("Использование: go run ./cmd/message-status UUID")
	}
	fields := url.Values{"fields": {"uuid,status,errCode"}}
	request(http.MethodGet, "message/"+args[0], fields)
}
