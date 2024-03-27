package config

import (
	"encoding/json"
	"log"
	"os"
)

type AppConfig struct {
	App    *App    `json:"app"`
	Redis  *Redis  `json:"redis"`
	Pulsar *Pulsar `json:"pulsar"`
}

type App struct {
	Port string `json:"port"`
}

type Redis struct {
	HostName   *string `json:"hostName"`
	Password   *string `json:"password"`
	DB         *int    `json:"db"`
	MaxRetries *int    `json:"maxRetries"`
}

type Pulsar struct {
	HostName         *string `json:"hostName"`
	Topic            *Topic  `json:"topic"`
	SubscriptionName *string `json:"subscriptionName"`
}

type Topic struct {
	CameraEvent *string `json:"cameraEvent"`
}

var config *AppConfig

func InitConfig() {
	configFile, err := os.ReadFile("config/config.json")
	if err != nil {
		log.Fatalf("Error reading config file. Error: %s", err.Error())
	}
	if err = json.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("Error unmarshalling config file. Error: %s", err.Error())
	}
	log.Println("Initializing config")
}

func GetConfig() *AppConfig {
	if config != nil {
		return config
	}
	return nil
}
