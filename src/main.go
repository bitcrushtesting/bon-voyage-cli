// Copyright 2024 Bitcrush Testing

package main

import (
	"bon-voyage-cli/cmd"
	"bon-voyage-cli/connection"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the configuration structure
type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func main() {

	config, err := loadConfig("./config.yaml")
	if err != nil {
		os.Exit(1)
	}

	connection.Init(config.Host, config.Port)

	info, err := connection.Information()
	if err != nil {
		fmt.Println("Could not gather server information", err)
		os.Exit(1)
	}

	fmt.Println("------------------------------------")
	fmt.Println("Server name:", info.Name)
	fmt.Println("Server id:", info.Id)
	fmt.Println("------------------------------------")

	cmd.Execute()

	fmt.Println("------------------------------------")
}

func loadConfig(configFile string) (*Config, error) {

	configData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
		return nil, err
	}
	return &config, nil
}
