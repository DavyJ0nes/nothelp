package cmd

import (
	"os"
	"testing"
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
