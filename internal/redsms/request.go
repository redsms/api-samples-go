package redsms

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
	"sort"
	"strings"
	"time"
)

func required(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatal("Задайте ", name)
	}
	return value
}

func request(method, path string, fields url.Values) {
	login, key := required("REDSMS_LOGIN"), required("REDSMS_API_KEY")
	if method == http.MethodGet && len(fields) > 0 {
		path += "?" + fields.Encode()
	}
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var values strings.Builder
	for _, k := range keys {
		values.WriteString(fields.Get(k))
	}
	var nonce [16]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		log.Fatal(err)
	}
	ts := fmt.Sprintf("%x", nonce)
	var requestBody io.Reader
	if method == http.MethodPost {
		requestBody = strings.NewReader(fields.Encode())
	}
	req, err := http.NewRequest(method, "https://cp.redsms.ru/api/"+path, requestBody)
	if err != nil {
		log.Fatal(err)
	}
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("login", login)
	req.Header.Set("ts", ts)
	req.Header.Set("sig", fmt.Sprintf("%x", md5.Sum([]byte(values.String()+ts+key))))
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
