package config

import "testing"

func TestDatabaseURL(t *testing.T) {
	t.Setenv("POSTGRES_USER", "app")
	t.Setenv("POSTGRES_PASSWORD", "abc-123")
	t.Setenv("POSTGRES_DB", "app")
	t.Setenv("POSTGRES_HOST", "")

	got, err := DatabaseURL()
	if err != nil {
		t.Fatal(err)
	}
	want := "postgres://app:abc-123@postgres:5432/app?sslmode=disable"
	if got != want {
		t.Fatalf("DatabaseURL() = %q, want %q", got, want)
	}
}

func TestDatabaseURLUsesLocalHost(t *testing.T) {
	t.Setenv("POSTGRES_USER", "app")
	t.Setenv("POSTGRES_PASSWORD", "abc-123")
	t.Setenv("POSTGRES_DB", "app")
	t.Setenv("POSTGRES_HOST", "127.0.0.1")

	got, err := DatabaseURL()
	if err != nil {
		t.Fatal(err)
	}
	want := "postgres://app:abc-123@127.0.0.1:5432/app?sslmode=disable"
	if got != want {
		t.Fatalf("DatabaseURL() = %q, want %q", got, want)
	}
}

func TestDatabaseURLRejectsMissingOrInvalidValues(t *testing.T) {
	t.Setenv("POSTGRES_USER", "app")
	t.Setenv("POSTGRES_PASSWORD", "")
	t.Setenv("POSTGRES_DB", "app")
	t.Setenv("POSTGRES_HOST", "")
	if _, err := DatabaseURL(); err == nil {
		t.Fatal("missing password returned a URL")
	}

	t.Setenv("POSTGRES_PASSWORD", "bad password")
	if _, err := DatabaseURL(); err == nil {
		t.Fatal("password with a space returned a URL")
	}

	t.Setenv("POSTGRES_PASSWORD", "abc-123")
	t.Setenv("POSTGRES_HOST", "db.example.com")
	if _, err := DatabaseURL(); err == nil {
		t.Fatal("unexpected host returned a URL")
	}
}

func TestLoad(t *testing.T) {
	t.Setenv("POSTGRES_USER", "app")
	t.Setenv("POSTGRES_PASSWORD", "abc-123")
	t.Setenv("POSTGRES_DB", "app")
	t.Setenv("POSTGRES_HOST", "postgres")
	t.Setenv("ENCRYPTION_KEY", "test-key")

	cfg, err := Load(":8000")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Addr != ":8000" {
		t.Fatalf("Addr = %q, want :8000", cfg.Addr)
	}
	if cfg.EncryptionKey != "test-key" {
		t.Fatalf("EncryptionKey = %q, want test-key", cfg.EncryptionKey)
	}
	want := "postgres://app:abc-123@postgres:5432/app?sslmode=disable"
	if cfg.DatabaseURL != want {
		t.Fatalf("DatabaseURL = %q, want %q", cfg.DatabaseURL, want)
	}
}
