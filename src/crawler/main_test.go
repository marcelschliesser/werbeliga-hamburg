package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestHost(t *testing.T) {
	c := NewCrawler(os.Getenv("URL"), 10)

	u, err := url.Parse(c.baseUrl)
	if err != nil {
		panic(err)
	}

	host := u.Host
	scheme := u.Scheme
	req, err := http.NewRequest("GET", fmt.Sprintf("%v://%v", scheme, host), nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	req.Header.Set("User-Agent", "MyApp/1.0")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	if statusCode != 200 {
		t.Errorf("wrong status code: got %v, want 200", statusCode)
	}
}
