package config

import "os"

type Config struct {
	DataLocation string `json:"data_location"`
}

func Parse() (Config, error) {
	obsidianPrefix := "/Library/Mobile Documents/iCloud~md~obsidian/Documents"
	homedir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	return Config{
		DataLocation: homedir + obsidianPrefix + "/notes/daily",
	}, nil
}
