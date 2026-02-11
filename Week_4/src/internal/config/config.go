package config

import (
	"log"
	"os"

)

type Config struct {
	Addr		string
	Namespace	string
}

func getenv(name, key string) string {
	env := os.Getenv(name)
	if env == "" {
		return key
	}
	return env
}

func MustLoad() Config {
	cfg := Config{
		Addr: 		getenv("ADDR", ":8080"),
		Namespace:	getenv("NAMESPACE", "paas-postgres"),
	}
	return cfg
}

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}