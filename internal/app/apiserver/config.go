package apiserver

import "os"

type Config struct {
	BindAddr    string
	LogLevel    string
	DatabaseURL string
	SessionKey  string
}

// NewConfig Возвращаем конфиг.
func NewConfig() *Config {
	return &Config{
		BindAddr: getEnv("BANK_BIND_ADDR", ""),
		LogLevel: getEnv("BANK_LOG_LEVEL", ""),
		DatabaseURL: getEnv("BANK_DATABASE_URL", ""),
		SessionKey: getEnv("BANK_SESSION_KEY", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}

	return defaultVal
}
