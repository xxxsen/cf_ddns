package config

import (
	"encoding/json"
	"io"
	"os"
)

type ProviderConfig struct {
	Name string      `json:"name"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type CloudflareConfig struct {
	Key        string `json:"key"`
	EMail      string `json:"email"`
	RecordType string `json:"record_type"`
	RecordName string `json:"record_name"`
	ZoneName   string `json:"zone_name"`
	TTL        int    `json:"ttl"`
	Proxied    bool   `json:"proxied"`
}

type RefreshCongfig struct {
	Name             string           `json:"name"`
	Providers        []string         `json:"providers"`
	CloudflareConfig CloudflareConfig `json:"cloudflare_config"`
	RefreshInterval  int              `json:"refresh_interval"`
}

type Config struct {
	ProviderList  []ProviderConfig `json:"provider_list"`
	RefresherList []RefreshCongfig `json:"refresher_list"`
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
