package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App        `yaml:"app"`
	Postgresql `yaml:"postgresql"`
	Nats       `yaml:"nats-subscriber"`
}

type App struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port string `yaml:"port" env-default:"80"`
}

type Postgresql struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type Nats struct {
	ClusterID string `yaml:"cluster_id"`
	ClientID  string `yaml:"client_id"`
	Subject   string `yaml:"subject"`
	NatsUrl   string `yaml:"nats_url"`
}

var (
	instance *Config
	once     sync.Once
)

func GetCongfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Println(help)
			log.Fatal(err)
		}
	})
	return instance
}
