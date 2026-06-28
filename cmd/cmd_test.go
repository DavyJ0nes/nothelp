package cmd

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/davyj0nes/nothelp/internal/config"
)

func TestExists(t *testing.T) {
	// Test case where file exists
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	if !exists(file.Name()) {
		t.Errorf("expected file to exist")
	}

	// Test case where file does not exist
	if exists("nonexistentfile") {
		t.Errorf("expected file to not exist")
	}
}

func TestThisWeek(t *testing.T) {
	got := thisWeek()

	matched, err := regexp.MatchString(`^\d{4}-W\d{2}$`, got)
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Errorf("expected ISO week format YYYY-Www, got %q", got)
	}
}

func TestGetWeeklyFilePathCreatesFile(t *testing.T) {
	dir := t.TempDir()
	conf := config.Config{WeeklyLocation: dir}
	week := "2026-W27"

	filePath, err := getWeeklyFilePath(conf, week)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := filepath.Join(dir, week+".md")
	if filePath != want {
		t.Errorf("expected path %q, got %q", want, filePath)
	}

	body, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("expected weekly note to be created, got %v", err)
	}
	if !strings.Contains(string(body), "Weekly Review — "+week) {
		t.Errorf("expected rendered week %q in note, got:\n%s", week, body)
	}
}
