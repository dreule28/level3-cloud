package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Addr      string
	Namespace string
	LogsPath  string

	AuthUser string
	AuthPass string

	JWTSecret   string
	JWTIssuer   string
	JWTAudience string
	JWTTL       time.Duration
}

func getenv(name, defaultVal string) string {
	if env, ok := os.LookupEnv(name); ok {
		return env
	}
	return defaultVal
}

func MustLoad() Config {
	cfg := Config{
		Addr:      getenv("ADDR", ":8080"),
		Namespace: getenv("NAMESPACE", "paas-postgres"),
		LogsPath:  getenv("LOG_STORE_PATH", "data/logs.jsonl"),

		AuthUser: getenv("AUTH_USER", "admin"),
		AuthPass: getenv("AUTH_PASS", "password"),

		JWTSecret:   getenv("JWT_SECRET", "secret"),
		JWTIssuer:   getenv("JWT_ISSUER", "paas-api"),
		JWTAudience: getenv("JWT_AUDIENCE", "paas-ui"),
		JWTTL:       15 * time.Minute,
	}
	return cfg
}

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
