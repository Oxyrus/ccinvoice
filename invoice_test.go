package main

import (
	"testing"
	"time"
)

func TestGetInvoiceNumber(t *testing.T) {
	dog := Dog{
		ID:   1,
		Name: "Fido",
	}

	date := time.Now()
	want := "FID" + date.Format("20060102")
	if got := getInvoiceNumber(dog); got != want {
		t.Errorf("getInvoiceNumber() = %q, want %q", got, want)
	}
}

func TestGetNextMonday(t *testing.T) {
	date := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	want := time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC)
	if got := nextMonday(date); got != want {
		t.Errorf("getNextMonday() = %q, want %q", got, want)
	}

	// Check that if the date is already a Monday, it returns the following TestGetNextMonday
	date = time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC)
	want = time.Date(2020, 1, 13, 0, 0, 0, 0, time.UTC)
	if got := nextMonday(date); got != want {
		t.Errorf("getNextMonday() = %q, want %q", got, want)
	}
}

// Test currently fails because it relies on the web server running
func TestGeneratePdf(t *testing.T) {
	dog := Dog{
		ID:        1,
		Name:      "Fido",
		OwnerName: "John Doe",
		Address:   "123 Fake Street",
		City:      "Fakeville",
		Email:     "johndoe@example.com",
		Service:   "walk",
		Quantity:  1,
		Price:     25,
	}

	want := "./public/FID" + time.Now().Format("20060102") + ".pdf"
	got, err := generatePdf(dog)
	if err != nil {
		t.Errorf("generatePdf() error = %q", err)
	}
	if got != want {
		t.Errorf("generateInvoice() = %q, want %q", got, want)
	}
}
