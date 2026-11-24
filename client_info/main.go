package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"os"

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
	if len(os.Args) != 1 {
		log.Fatal("Запустите без аргументов: go run client_info/main.go")
	}
	path := "client/info"
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
	req.Header.Set("sig", fmt.Sprintf("%x", md5.Sum([]byte(ts+key))))
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
