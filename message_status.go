package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

func required(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatal("Задайте ", name)
	}
	return value
}

func main() {
	uuidPattern := regexp.MustCompile("(?i)^[0-9a-f]{8}(-[0-9a-f]{4}){3}-[0-9a-f]{12}$")
	if len(os.Args) != 2 || !uuidPattern.MatchString(os.Args[1]) {
		log.Fatal("Укажите UUID: go run message_status.go UUID")
	}
	fields := "uuid,status,errCode"
	query := url.Values{"fields": {fields}}
	path := "message/" + os.Args[1] + "?" + query.Encode()
	login, key := required("REDSMS_LOGIN"), required("REDSMS_API_KEY")
	var nonce [16]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		log.Fatal(err)
	}
	ts := fmt.Sprintf("%x", nonce)
	req, err := http.NewRequest(http.MethodGet, "https://cp.redsms.ru/api/"+path, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("login", login)
	req.Header.Set("ts", ts)
	req.Header.Set("sig", fmt.Sprintf("%x", md5.Sum([]byte(fields+ts+key))))
	client := http.Client{Timeout: 15 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("HTTP %s: %s", res.Status, body)
	}
	var result struct {
		Success bool              `json:"success"`
		Errors  []json.RawMessage `json:"errors"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
	if !result.Success || len(result.Errors) > 0 {
		os.Exit(1)
	}
}
