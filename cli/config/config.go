package config

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	RegisterURL     string `json:"register_url"`
	AuthURL         string `json:"auth_url"`
	SendLogoPassURL string `json:"send_logo_pass_url"`
	GetLogoPassURL  string `json:"get_logo_pass_url"`
	SendBin         string `json:"send_bin_url"`
	GetBin          string `json:"get_bin_url"`
	SyncInterval    int    `json:"sync_interval"`
}

func New(data []byte) (*Config, error) {
	if data == nil || len(data) == 0 {
		return nil, fmt.Errorf("config data is empty")
	}
	var cfg *Config
	err := json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}
	return cfg, nil
}
