package config

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

type ProviderConfig struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type CloudflareConfig struct {
	Key        string `json:"key"`
	EMail      string `json:"email"`
	RecordType string `json:"record_type"`
	RecordName string `json:"record_name"`
	ZoneName   string `json:"zone_name"`
}

type RefreshCongfig struct {
	Name             string           `json:"name"`
	Providers        []string         `json:"providers"`
	CloudflareConfig CloudflareConfig `json:"cloudflare_config"`
	RefreshInterval  time.Duration    `json:"refresh_interval"`
}

type Config struct {
	ProviderList   []ProviderConfig `json:"provider_list"`
	RefreshCongfig []RefreshCongfig `json:"refresh_config"`
}

func Parse(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	raw, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	if err := json.Unmarshal(raw, c); err != nil {
		return nil, err
	}
	return c, nil
}
