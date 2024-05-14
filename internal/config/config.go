package config

import (
	"os"
)

type Config struct {
	EtherscanApiKey string
	CovalentApiKey  string
}

func NewConfig() *Config {
	return &Config{
		EtherscanApiKey: os.Getenv("ETHERSCAN_API_KEY"),
		CovalentApiKey:  os.Getenv("COVALENT_API_KEY"),
	}
}
