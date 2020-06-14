package coursesbot

import (
	"fmt"
	"os"
)

type Config struct {
	APIKey string
	Debug  bool
}

func NewConfig() (*Config, error) {
	key := os.Getenv("API_KEY")
	debug := os.Getenv("DEBUG") == "1"
	if key == "" {
		return nil, fmt.Errorf("API_KEY is missiong")
	}
	return &Config{key, debug}, nil
}
