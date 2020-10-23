package config

import (
	"bufio"
	"flag"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DbURL string `envconfig:"GRABBER_DB_URL"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DbURL: "postgresql://postgres:postgres@postgresdb:5432/grabberdb?sslmode=disable",
	}
	err := envconfig.Process("GRABBER", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// SetEnvFromFile sets environment variables from file if file specified as CLI argument
func SetEnvFromFile() error {
	envFile := flag.String("env-file", "", "Path to file with env values")
	flag.Parse()

	if *envFile == "" {
		return nil
	}

	file, err := os.Open(*envFile)
	if file != nil {
		defer func() {
			_ = file.Close()
		}()
	}
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.SplitN(scanner.Text(), "=", 2)
		if len(s) < 2 {
			continue
		}

		if err := os.Setenv(s[0], s[1]); err != nil {
			return err
		}
	}

	return scanner.Err()
}
