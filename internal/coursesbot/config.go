package coursesbot

import (
	"fmt"
	"os"
)

type Config struct {
	APIKey string
	Debug  bool
	Proxy  *ProxyConfig
}

type ProxyConfig struct {
	Type     string
	URL      string
	User     string
	Password string
}

func NewConfig() (*Config, error) {
	key := os.Getenv("API_KEY")
	debug := os.Getenv("DEBUG") == "1"
	if key == "" {
		return nil, fmt.Errorf("API_KEY is missiong")
	}

	var proxyConf *ProxyConfig

	if os.Getenv("SOCKS5_URL") != "" {
		proxyConf = &ProxyConfig{
			Type:     "socks5",
			URL:      os.Getenv("SOCKS5_URL"),
			User:     os.Getenv("SOCKS5_USER"),
			Password: os.Getenv("SOCKS5_PASSWORD"),
		}
	}

	return &Config{key, debug, proxyConf}, nil
}
