package config

import (
	"os"
	"strconv"
)

type Config struct {
	ReplicationFactor int
}

func LoadConfig() Config {
	repFactor, err := strconv.Atoi(os.Getenv("REPLICATION_FACTOR"))
	if err != nil {
		repFactor = 3
	}

	return Config{
		ReplicationFactor: repFactor,
	}
}
