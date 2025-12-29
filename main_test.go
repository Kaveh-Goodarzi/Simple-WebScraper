package main

import "testing"

func TestExtractProductFromHTML(t *testing.T) {
	p := extractProductFromHTML(
		"Test Product",
		"$99",
		"/img.png",
		"/product/1",
	)

	if p.Name != "Test Product" {
		t.Fatalf("expected name 'Test Product', got '%s'", p.Name)
	}

	if p.Price != "$99" {
		t.Fatalf("expected price '$99', got '%s'", p.Price)
	}

	if p.Image != "/img.png" {
		t.Fatalf("unexpected image value")
	}

	if p.URL != "/product/1" {
		t.Fatalf("unexpected url value")
	}
}