package config

import "os"

type Config struct {
	DataLocation string `json:"data_location"`
}

func Parse() (Config, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	return Config{
		DataLocation: homedir + "/notes/daily",
	}, nil
}
