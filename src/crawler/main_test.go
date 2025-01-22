package main

import (
	"fmt"
	"testing"
)

// TODO: Finalize Tests
func TestURLStatus(t *testing.T) {
	c := NewCrawler(baseUrl, 10)
	doc := c.FetchUrl(2, 1)
	fmt.Println(doc)
	// if _, ok := doc.(*goquery.Document); !ok {
	// 	t.Errorf("wrong type: got %T, want *goquery.Document", doc)
	// }
}
