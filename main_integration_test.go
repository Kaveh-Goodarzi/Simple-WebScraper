//go:build integration
// +build integration

package main

import (
	"testing"
)
func TestScraperWithInternet(t *testing.T) {
	s := NewScraper()

	err := s.Start("https://www.scrapingcourse.com/ecommerce")
	if err != nil {
		t.Fatalf("scraper failed: %v", err)
	}

	if len(s.products) == 0 {
		t.Fatalf("expected products, got zero")
	}
}
