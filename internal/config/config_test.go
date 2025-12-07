package config

import (
	"testing"

	"gotest.tools/v3/assert"
)

const obsidianPrefix = "/Library/Mobile Documents/iCloud~md~obsidian/Documents"

func TestParse(t *testing.T) {
	home := "/Users/tester"
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	cfg, err := Parse()
	assert.NilError(t, err)

	assert.Equal(t, home+obsidianPrefix+"/notes/daily", cfg.DataLocation)
	assert.Equal(t, home+obsidianPrefix+"/notes/daily/archive", cfg.ArchiveLocation)
}

func TestGetDataFilePath(t *testing.T) {
	cfg := Config{DataLocation: "/tmp/data"}
	date := "2025-01-02"

	path := cfg.GetDataFilePath(date)

	assert.Equal(t, "/tmp/data/2025-01-02.md", path)
}

func TestGetArchiveFilePath(t *testing.T) {
	cfg := Config{ArchiveLocation: "/tmp/archive"}
	date := "2025-01-02"

	path := cfg.GetArchiveFilePath(date)

	assert.Equal(t, "/tmp/archive/2025-01-02.md", path)
}
