package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	REST_PORT string `envconfig:"REST_PORT" default:"8080"`

	DBHost     string `envconfig:"MYSQL_DB_HOST" default:"svc-mix-id-1"`
	DBPort     string `envconfig:"MYSQL_DB_PORT" default:"3306"`
	DBUsername string `envconfig:"MYSQL_DB_USERNAME" default:"quest"`
	DBPassword string `envconfig:"MYSQL_DB_PASSWORD" default:""`
	DBName     string `envconfig:"MYSQL_DB_NAME" default:"quest"`
}

func New() Config {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	return c
}
