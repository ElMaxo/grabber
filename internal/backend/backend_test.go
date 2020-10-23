package backend

import (
	"grabber/internal/config"
	"grabber/internal/repository"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	repo, err := repository.NewPostgresRepository(cfg.DbURL)
	if err != nil {
		log.Fatal("failed to create repository: ", err)
	}

	err = repo.Migrate()
	if err != nil {
		log.Fatal("migrations application error")
	}

	os.Exit(m.Run())
}
