package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const obsidianPrefix = "/Library/Mobile Documents/iCloud~md~obsidian/Documents"

func TestParse(t *testing.T) {
	home := "/Users/tester"
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	cfg, err := Parse()
	require.NoError(t, err)

	require.Equal(t, home+obsidianPrefix+"/notes/daily", cfg.DataLocation)
	require.Equal(t, home+obsidianPrefix+"/notes/daily/archive", cfg.ArchiveLocation)
}

func TestGetDataFilePath(t *testing.T) {
	cfg := Config{DataLocation: "/tmp/data"}
	date := "2025-01-02"

	path := cfg.GetDataFilePath(date)

	require.Equal(t, "/tmp/data/2025-01-02.md", path)
}

func TestGetArchiveFilePath(t *testing.T) {
	cfg := Config{ArchiveLocation: "/tmp/archive"}
	date := "2025-01-02"

	path := cfg.GetArchiveFilePath(date)

	require.Equal(t, "/tmp/archive/2025-01-02.md", path)
}
