package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type URL string

type URLS []URL

func (urls URLS) IsValid() bool {
	validate := validator.New(validator.WithRequiredStructEnabled())

	for _, u := range urls {
		err := validate.Var(u, "url")
		// _, err := url.ParseRequestURI(string(u))
		if err != nil {
			panic(fmt.Sprintf("invalid urls: %s", string(u)))
		}
	}
	return true
}

func NewURLS(urls ...string) URLS {
	var result URLS
	for _, url := range urls {
		result = append(result, URL(url))
	}
	return result
}

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
	Staging     Environment = "staging"
)

type Config struct {
	PG_HOST                   string
	PG_PORT                   int
	PG_NAME                   string
	PG_USER                   string
	PG_PASS                   string
	PG_SSLMODE                string
	PASSWORD_HASH_SALT        string
	PASETO_PRIVATE_KEY_SECRET string
	PASETO_PUBLIC_KEY_SECRET  string
	PASETO_REFRESH_TOKEN_TTL  int
	PASETO_ACCESS_TOKEN_TTL   int
	DATABASE_URL              string
	CERBOS_URL                string
	RABBIT_MQ_URL             string
	PORT                      int
	RUN_SEEDS                 bool
	ENV                       Environment
	NATS_URL                  string
	NATS_JWT                  string
	NATS_SEED                 string
	JAEGER_ENDPOINT           string
	SERVICE_NAME              string
}

func Get(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("%s: %s", key, err)
			return fallback
		}
		return i
	}
	return fallback
}

func GetEnvironment() Environment {
	if env := Get("ENV", ""); env == "" {
		return Development
	} else {
		return Environment(env)
	}
}

func LoadConfig() *Config {
	if env := os.Getenv("ENV"); env != "" {
		var config Config

		if err := json.Unmarshal([]byte(env), &config); err != nil {
			fmt.Println("error parsing secret json:", err)
			return nil
		}

		return &config
	}

	if os.Getenv("AWS_EXECUTION_ENV") == "" {
		if err := godotenv.Load(); err != nil {
			log.Println("error loading env file:", err)
		}
	}

	natsServers := NewURLS(strings.Split(Get("NATS_URL", "nats://localhost:4222"), ",")...)
	natsServers.IsValid()

	cerbosServers := NewURLS(strings.Split(Get("CERBOS_URL", "http://localhost:3592"), ",")...)
	cerbosServers.IsValid()

	return &Config{
		PG_HOST:                   Get("PG_HOST", "localhost"),
		PG_PORT:                   GetInt("PG_PORT", 5432),
		PG_NAME:                   Get("PG_NAME", "kale"),
		PG_USER:                   Get("PG_USER", "postgres"),
		PG_PASS:                   Get("PG_PASS", "postgres"),
		PG_SSLMODE:                Get("PG_SSLMODE", "disable"),
		PASSWORD_HASH_SALT:        Get("PASSWORD_HASH_SALT", ""),
		PASETO_PRIVATE_KEY_SECRET: Get("PASETO_PRIVATE_KEY_SECRET", "love"),
		PASETO_PUBLIC_KEY_SECRET:  Get("PASETO_PUBLIC_KEY_SECRET", "opus"),
		PASETO_ACCESS_TOKEN_TTL:   GetInt("PASETO_ACCESS_TOKEN_TTL", 1),      // TTL in minutes
		PASETO_REFRESH_TOKEN_TTL:  GetInt("PASETO_REFRESH_TOKEN_TTL", 10080), // TTL in minues
		DATABASE_URL:              Get("DATABASE_URL", ""),
		CERBOS_URL:                Get("CERBOS_URL", "http://localhost:3592"),
		RABBIT_MQ_URL:             Get("RABBIT_MQ_URL", "amqp://guest:guest@localhost:5672"),
		PORT:                      GetInt("PORT", 3000),
		ENV:                       GetEnvironment(),
		RUN_SEEDS:                 true,
		NATS_URL:                  Get("NATS_URL", "nats://localhost:4222"),
		NATS_JWT:                  Get("NATS_JWT", ""),
		NATS_SEED:                 Get("NATS_SEED", ""),
		JAEGER_ENDPOINT:           Get("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
		SERVICE_NAME:              Get("SERVICE_NAME", "auth"),
	}
}

var GetConfig = sync.OnceValue(LoadConfig)
