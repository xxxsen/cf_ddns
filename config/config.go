package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/xxxsen/common/logger"
)

type ProviderConfig struct {
	Name string      `json:"name"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type NotifierConfig struct {
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
	Disable          bool             `json:"disable"`
	Providers        []string         `json:"providers"`
	CloudflareConfig CloudflareConfig `json:"cloudflare_config"`
	RefreshInterval  int              `json:"refresh_interval"`
	Notifier         string           `json:"notifier"`
}

type Config struct {
	ProviderList  []ProviderConfig `json:"provider_list"`
	NotifierList  []NotifierConfig `json:"notifier_list"`
	RefresherList []RefreshCongfig `json:"refresher_list"`
	LogConfig     logger.LogConfig `json:"log_config"`
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
