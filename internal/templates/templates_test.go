package templates_test

import (
	"testing"

	"github.com/davyj0nes/nothelp/internal/templates"
)

func TestParse(t *testing.T) {
	date := "2025-01-19"
	result, err := templates.Parse(date)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Logf("Result: %s", result)
}
