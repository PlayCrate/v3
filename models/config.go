package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	ListenAddress string `json:"listenAddress"`
	Auth          string `json:"Auth"`
	DBConnString  string `json:"connString"`
	CutOffTime    int64  `json:"cutoffTime"`
	V1Auth        string `json:"v1-auth"`
	Prod          bool   `json:"PROD"`
	Cron          string `json:"Cron"`
}

// NewConfig creates a configuration from file
func NewConfig(filePath string) *Config {
	cfg := loadConfiguration(filePath)
	log.Printf("Loaded configuration from %s\n", filePath)
	return cfg
}

func loadConfiguration(filePath string) *Config {
	cfgFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Unable to read config file:", err)
	}

	// Create an instance of the Config struct
	config := Config{}

	// Parse the JSON content into the Config struct
	err = json.Unmarshal(cfgFile, &config)
	if err != nil {
		log.Fatal("Unable to parse config file:", err)
	}
	return &config
}
