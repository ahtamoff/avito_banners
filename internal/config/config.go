package config

import (
	"avito_banners/pkg/logging"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string     `yaml:"env" env-default:"local"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
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

var instanse *Config
var once sync.Once

func GetConfig() *Config{
	once.Do(func(){
		logger := logging.GetLogger()
		logger.Info("read configuration")
		instanse := &Config{}

		if err := cleanenv.ReadConfig("config.yaml", instanse); err != nil {
			help, _ := cleanenv.GetDescription(instanse, nil)
			logger.Info(help)
			logger.Fatal(err)
		}


	})
	return instanse
}