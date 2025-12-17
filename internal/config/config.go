package config

import "os"

type Config struct {
	ProjectID     string
	ProjectNumber string
	Collection    string
	DatabaseID    string
	DocumentID    string
	Location      string
	Mode          string
}

func Load() Config {
	return Config{
		ProjectID:     os.Getenv("GCP_PROJECT"),
		ProjectNumber: getEnv("GCP_PROJECT_NUMBER", "651545901471"),
		Collection:    getEnv("STATE_COLLECTION", "bot_state"),
		DocumentID:    getEnv("STATE_DOCUMENT", "eth_main"),
		DatabaseID:    getEnv("FIRESTORE_DB", "basicdata"),
		Location:      getEnv("LOCATION", "asia-southeast2"),
		Mode:          getEnv("MODE", "SIMULATION"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
