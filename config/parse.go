package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Disassemble bool `yaml:"disassemble"`
}

var CONFIG Config

func ParseConfig() {
	file, err := os.Open("config.yml")
	if err != nil {
		fmt.Printf("Warning: error opening config.yml: %v", err)
		return
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&CONFIG)
	if err != nil {
		fmt.Printf("error parsing config.yml: %v", err)
		return
	}
	fmt.Printf("Config loaded: %+v\n", CONFIG)
}
