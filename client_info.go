package main

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	login, key := os.Getenv("REDSMS_LOGIN"), os.Getenv("REDSMS_API_KEY")
	if login == "" || key == "" {
		log.Fatal("Задайте REDSMS_LOGIN и REDSMS_API_KEY")
	}
	var nonce [16]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		log.Fatal(err)
	}
	ts := fmt.Sprintf("%x", nonce)
	req, err := http.NewRequest(http.MethodGet, "https://cp.redsms.ru/api/client/info", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("login", login)
	req.Header.Set("ts", ts)
	req.Header.Set("sig", fmt.Sprintf("%x", md5.Sum([]byte(ts+key))))
	client := http.Client{Timeout: 15 * time.Second}
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
	fmt.Println(string(body))
}
