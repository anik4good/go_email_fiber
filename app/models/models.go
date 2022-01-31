package models

type Config struct {
	DB_USERNAME string `yaml:"DB_USERNAME"`
	DB_PASSWORD string `yaml:"DB_PASSWORD"`
	DB_IP       string `yaml:"DB_IP"`
	DB_PORT     string `yaml:"DB_PORT"`
	DB_NAME     string `yaml:"DB_NAME"`
}

type QueuedEmail struct {
	ID     uint64
	Name   string
	Email  string
	Status uint64
}
