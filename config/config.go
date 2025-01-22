package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Env struct {
	EMAIL_AWS_ACCESS_KEY_ID     string `mapstructure:"EMAIL_AWS_ACCESS_KEY_ID"`
	EMAIL_AWS_SECRET_ACCESS_KEY string `mapstructure:"EMAIL_AWS_SECRET_ACCESS_KEY"`
	EMAIL_AWS_REGION            string `mapstructure:"EMAIL_AWS_REGION"`
	EMAIL_AWS_ENDPOINT          string `mapstructure:"EMAIL_AWS_ENDPOINT"`
	EMAIL_TEMPLATES_BUCKET      string `mapstructure:"EMAIL_TEMPLATES_BUCKET"`
}

var (
	env  *Env
	once sync.Once
)

func loadEnv(path string) (*Env, error) {
	err := godotenv.Load(path)
	if err != nil {
		if os.Getenv("APP_ENV") != "production" {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	env := &Env{
		EMAIL_AWS_ACCESS_KEY_ID:     loadEnvVar("AWS_ACCESS_KEY_ID", ""),
		EMAIL_AWS_SECRET_ACCESS_KEY: loadEnvVar("AWS_SECRET_ACCESS_KEY", ""),
		EMAIL_AWS_REGION:            loadEnvVar("AWS_REGION", ""),
		EMAIL_AWS_ENDPOINT:          loadEnvVar("AWS_ENDPOINT", ""),
		EMAIL_TEMPLATES_BUCKET:      loadEnvVar("EMAIL_TEMPLATES_BUCKET", ""),
	}

	v := validator.New()
	err = v.Struct(env)
	return env, err
}

func loadEnvVar[T any](key string, defaultValue T) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	var result any
	switch any(defaultValue).(type) {
	case string:
		result = value
	case int:
		parsed, err := strconv.Atoi(value)
		if err != nil {
			panic(fmt.Sprintf("Erro ao converter variável %s para int: %v", key, err))
		}
		result = parsed
	case bool:
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			panic(fmt.Sprintf("Erro ao converter variável %s para bool: %v", key, err))
		}
		result = parsed
	case float64:
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic(fmt.Sprintf("Erro ao converter variável %s para float64: %v", key, err))
		}
		result = parsed
	default:
		panic(fmt.Sprintf("Tipo não suportado para variável %s", key))
	}

	return result.(T)
}

func newEnv() *Env {
	env, err := loadEnv("../.env")
	if err != nil {
		panic(err)
	}
	return env
}
func GetEnv() *Env {
	once.Do(func() {
		env = newEnv()
	})
	return env
}
