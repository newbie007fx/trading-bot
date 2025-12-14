package config

import "os"

type Config struct {
	ProjectID  string
	Collection string
	DatabaseID string
	DocumentID string
}

func Load() Config {
	return Config{
		ProjectID:  os.Getenv("GCP_PROJECT"),
		Collection: getEnv("STATE_COLLECTION", "bot_state"),
		DocumentID: getEnv("STATE_DOCUMENT", "eth_main"),
		DatabaseID: getEnv("FIRESTORE_DB", "basicdata"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
