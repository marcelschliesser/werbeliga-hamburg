package main

import (
	"fmt"
	"os"
	"testing"
)

// TODO: Finalize Tests
func TestURLStatus(t *testing.T) {
	c := NewCrawler(os.Getenv("URL"), 10)
	doc := c.FetchUrl(2, 1)
	fmt.Println(doc)
	// if _, ok := doc.(*goquery.Document); !ok {
	// 	t.Errorf("wrong type: got %T, want *goquery.Document", doc)
	// }
}
