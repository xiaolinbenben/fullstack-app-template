package config

import (
	"errors"
	"net"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	Addr          string
	DatabaseURL   string
	EncryptionKey string
}

func Load(addr string) (Config, error) {
	databaseURL, err := DatabaseURL()
	if err != nil {
		return Config{}, err
	}
	return Config{
		Addr:          addr,
		DatabaseURL:   databaseURL,
		EncryptionKey: os.Getenv("ENCRYPTION_KEY"),
	}, nil
}

// DatabaseURL 用官方 PostgreSQL 变量拼连接串。
// 主机名默认是 postgres。本机 make server-dev 和 make test 传入 127.0.0.1。
func DatabaseURL() (string, error) {
	user := strings.TrimSpace(os.Getenv("POSTGRES_USER"))
	password := strings.TrimSpace(os.Getenv("POSTGRES_PASSWORD"))
	name := strings.TrimSpace(os.Getenv("POSTGRES_DB"))
	if user == "" || password == "" || name == "" {
		return "", errors.New("缺少 POSTGRES_USER、POSTGRES_PASSWORD 或 POSTGRES_DB")
	}
	if !postgresToken(user) || !postgresToken(password) || !postgresToken(name) {
		return "", errors.New("POSTGRES_USER、POSTGRES_PASSWORD 和 POSTGRES_DB 只允许字母、数字和连字符")
	}

	host := strings.TrimSpace(os.Getenv("POSTGRES_HOST"))
	if host == "" {
		host = "postgres"
	}
	if host != "postgres" && host != "127.0.0.1" {
		return "", errors.New("数据库主机名只允许 postgres 或 127.0.0.1")
	}

	return (&url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(user, password),
		Host:     net.JoinHostPort(host, "5432"),
		Path:     "/" + name,
		RawQuery: "sslmode=disable",
	}).String(), nil
}

func postgresToken(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' {
			continue
		}
		return false
	}
	return true
}
