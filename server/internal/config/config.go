package config

import "os"

type Config struct {
	Addr          string
	DatabaseURL   string
	EncryptionKey string
}

func Load(addr string) Config {
	return Config{
		Addr:          addr,
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		EncryptionKey: os.Getenv("ENCRYPTION_KEY"),
	}
}
