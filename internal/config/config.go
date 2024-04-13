package config

import (
	"log"
	"os"
	"time"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string     `yaml:"env" env-default:"local"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
	DB 			DB 		   `yaml`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost: 8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"2s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"300s"`
}

type DB struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
	Name string `yaml:"name"`
}

func MustLoad() *Config{
	if err := godotenv.Load("../../.env"); err != nil {
        log.Print("No .env file found")
    }
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")	
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}