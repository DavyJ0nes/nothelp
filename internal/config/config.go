package config

import "os"

type Config struct {
	DataLocation    string
	WeeklyLocation  string
	ArchiveLocation string
}

func Parse() (Config, error) {
	obsidianPrefix := "/Library/Mobile Documents/iCloud~md~obsidian/Documents"
	homedir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	return Config{
		DataLocation:    homedir + obsidianPrefix + "/notes/daily",
		WeeklyLocation:  homedir + obsidianPrefix + "/notes/weekly",
		ArchiveLocation: homedir + obsidianPrefix + "/notes/daily/archive",
	}, nil
}

func (c Config) GetDataFilePath(date string) string {
	return c.DataLocation + "/" + date + ".md"
}

func (c Config) GetWeeklyFilePath(week string) string {
	return c.WeeklyLocation + "/" + week + ".md"
}

func (c Config) GetArchiveFilePath(date string) string {
	return c.ArchiveLocation + "/" + date + ".md"
}
