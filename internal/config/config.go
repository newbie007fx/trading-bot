package config

import "os"

type Config struct {
	ProjectID  string
	Collection string
	DocumentID string
}

func Load() Config {
	return Config{
		ProjectID:  os.Getenv("GCP_PROJECT"),
		Collection: getEnv("STATE_COLLECTION", "bot_state"),
		DocumentID: getEnv("STATE_DOCUMENT", "eth_main"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
