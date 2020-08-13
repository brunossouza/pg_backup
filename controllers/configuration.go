package controllers

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config struct
type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
	} `yaml:"database"`
	Backup struct {
		Directory    string   `yaml:"path_directory"`
		BackupFormat string   `yaml:"format"`
		Databases    []string `yaml:"databases"`
		Encode       string   `yaml:"encode"`
	} `yaml:"backup"`
}

// ReadFile read configuration file
func ReadFile(cfg *Config) {
	f, err := os.Open("config.yml")
	CheckError(err)
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	CheckError(err)
}
