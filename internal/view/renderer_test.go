package view

import (
	"path/filepath"
	"testing"
)

func TestTemplatesParse(t *testing.T) {
	root := filepath.Join("..", "..", "web", "templates")
	if _, err := New(root); err != nil {
		t.Fatalf("parse templates: %v", err)
	}
}

func TestDisplayFormatters(t *testing.T) {
	tests := map[string]struct {
		got  string
		want string
	}{
		"date":           {formatDate("2026-01-01"), "01 Jan 2026"},
		"empty date":     {formatDate(""), ""},
		"date time":      {formatDateTime("2026-08-12 07:05"), "12 Agu 2026, 07:05"},
		"rupiah":         {formatRupiah("10000000.00"), "Rp 10.000.000"},
		"rupiah decimal": {formatRupiah("1250.50"), "Rp 1.250,5"},
		"empty rupiah":   {formatRupiah(""), ""},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.got != test.want {
				t.Fatalf("got %q, want %q", test.got, test.want)
			}
		})
	}
}
