package config

type Config struct{
	Env string `yaml:"env" env-default:"local"`
} 