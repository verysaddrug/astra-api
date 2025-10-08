package config

import (
	"os"
	"testing"
)

func TestLoadConfig_WithEnvFile(t *testing.T) {
	// Create a temporary .env file
	envContent := `DB_HOST=localhost
DB_PORT=5432
DB_USER=testuser
DB_PASSWORD=testpass
DB_NAME=testdb
ADMIN_TOKEN=admintoken123
AUTO_MIGRATE=true`

	envFile := "test.env"
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("failed to create test env file: %v", err)
	}
	defer os.Remove(envFile)

	config := LoadConfig(envFile)

	if config.DBHost != "localhost" {
		t.Fatalf("expected DBHost 'localhost', got %s", config.DBHost)
	}
	if config.DBPort != "5432" {
		t.Fatalf("expected DBPort '5432', got %s", config.DBPort)
	}
	if config.DBUser != "testuser" {
		t.Fatalf("expected DBUser 'testuser', got %s", config.DBUser)
	}
	if config.DBPassword != "testpass" {
		t.Fatalf("expected DBPassword 'testpass', got %s", config.DBPassword)
	}
	if config.DBName != "testdb" {
		t.Fatalf("expected DBName 'testdb', got %s", config.DBName)
	}
	if config.AdminToken != "admintoken123" {
		t.Fatalf("expected AdminToken 'admintoken123', got %s", config.AdminToken)
	}
	if !config.AutoMigrate {
		t.Fatal("expected AutoMigrate to be true")
	}
}

func TestLoadConfig_WithEnvFile_AutoMigrateFalse(t *testing.T) {
	// Clear any existing environment variables first
	os.Unsetenv("AUTO_MIGRATE")

	// Create a temporary .env file
	envContent := `DB_HOST=localhost
DB_PORT=5432
DB_USER=testuser
DB_PASSWORD=testpass
DB_NAME=testdb
ADMIN_TOKEN=admintoken123
AUTO_MIGRATE=false`

	envFile := "test.env"
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("failed to create test env file: %v", err)
	}
	defer os.Remove(envFile)

	config := LoadConfig(envFile)

	if config.AutoMigrate {
		t.Fatal("expected AutoMigrate to be false")
	}
}

func TestLoadConfig_WithEnvFile_AutoMigrateNotSet(t *testing.T) {
	// Clear any existing environment variables first
	os.Unsetenv("AUTO_MIGRATE")

	// Create a temporary .env file
	envContent := `DB_HOST=localhost
DB_PORT=5432
DB_USER=testuser
DB_PASSWORD=testpass
DB_NAME=testdb
ADMIN_TOKEN=admintoken123`

	envFile := "test.env"
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("failed to create test env file: %v", err)
	}
	defer os.Remove(envFile)

	config := LoadConfig(envFile)

	if config.AutoMigrate {
		t.Fatal("expected AutoMigrate to be false when not set")
	}
}

func TestLoadConfig_WithNonExistentEnvFile(t *testing.T) {
	// Clear any existing environment variables
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("ADMIN_TOKEN")
	os.Unsetenv("AUTO_MIGRATE")

	// Test with non-existent env file
	config := LoadConfig("nonexistent.env")

	// Should still return a config with empty values
	if config == nil {
		t.Fatal("expected config to be non-nil")
	}
	if config.DBHost != "" {
		t.Fatalf("expected empty DBHost, got %s", config.DBHost)
	}
	if config.AutoMigrate {
		t.Fatal("expected AutoMigrate to be false")
	}
}

func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("DB_HOST", "envhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "envuser")
	os.Setenv("DB_PASSWORD", "envpass")
	os.Setenv("DB_NAME", "envdb")
	os.Setenv("ADMIN_TOKEN", "envtoken")
	os.Setenv("AUTO_MIGRATE", "true")
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("ADMIN_TOKEN")
		os.Unsetenv("AUTO_MIGRATE")
	}()

	config := LoadConfig("nonexistent.env")

	if config.DBHost != "envhost" {
		t.Fatalf("expected DBHost 'envhost', got %s", config.DBHost)
	}
	if config.DBPort != "3306" {
		t.Fatalf("expected DBPort '3306', got %s", config.DBPort)
	}
	if config.DBUser != "envuser" {
		t.Fatalf("expected DBUser 'envuser', got %s", config.DBUser)
	}
	if config.DBPassword != "envpass" {
		t.Fatalf("expected DBPassword 'envpass', got %s", config.DBPassword)
	}
	if config.DBName != "envdb" {
		t.Fatalf("expected DBName 'envdb', got %s", config.DBName)
	}
	if config.AdminToken != "envtoken" {
		t.Fatalf("expected AdminToken 'envtoken', got %s", config.AdminToken)
	}
	if !config.AutoMigrate {
		t.Fatal("expected AutoMigrate to be true")
	}
}

func TestLoadConfig_EmptyValues(t *testing.T) {
	// Clear any existing environment variables
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("ADMIN_TOKEN")
	os.Unsetenv("AUTO_MIGRATE")

	config := LoadConfig("nonexistent.env")

	if config.DBHost != "" {
		t.Fatalf("expected empty DBHost, got %s", config.DBHost)
	}
	if config.DBPort != "" {
		t.Fatalf("expected empty DBPort, got %s", config.DBPort)
	}
	if config.DBUser != "" {
		t.Fatalf("expected empty DBUser, got %s", config.DBUser)
	}
	if config.DBPassword != "" {
		t.Fatalf("expected empty DBPassword, got %s", config.DBPassword)
	}
	if config.DBName != "" {
		t.Fatalf("expected empty DBName, got %s", config.DBName)
	}
	if config.AdminToken != "" {
		t.Fatalf("expected empty AdminToken, got %s", config.AdminToken)
	}
	if config.AutoMigrate {
		t.Fatal("expected AutoMigrate to be false")
	}
}
