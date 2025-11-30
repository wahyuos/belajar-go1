package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config struct untuk menampung konfigurasi aplikasi
type Config struct {
    // ... (field DB yang sudah ada) ...
    GinMode    string // <-- NEW: Field untuk menampung GIN_MODE
}

// LoadConfig memuat variabel lingkungan dari file .env
func GinMode() (*Config, error) {
    if err := godotenv.Load(".env"); err != nil {
        fmt.Println("No .env file found, loading from environment variables.")
    }

    return &Config{
        // ... (pemuatan field DB yang sudah ada) ...
        GinMode:    os.Getenv("GIN_MODE"), // <-- Pemuatan GIN_MODE
    }, nil
}