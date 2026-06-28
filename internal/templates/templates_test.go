package templates_test

import (
	"strings"
	"testing"

	"github.com/davyj0nes/nothelp/internal/templates"
)

func TestParse(t *testing.T) {
	date := "2025-01-19"
	result, err := templates.Parse(date)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !strings.Contains(string(result), date) {
		t.Errorf("expected rendered template to contain date %q, got:\n%s", date, result)
	}
}

func TestParseWeekly(t *testing.T) {
	week := "2026-W27"
	result, err := templates.ParseWeekly(week)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !strings.Contains(string(result), "Weekly Review — "+week) {
		t.Errorf("expected rendered template to contain week %q, got:\n%s", week, result)
	}
}
