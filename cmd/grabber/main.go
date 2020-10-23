package main

import (
	"encoding/json"
	"grabber/internal/backend"
	"grabber/internal/config"
	"grabber/internal/grabber"
	"grabber/internal/repository"
	"grabber/internal/rest/restapi"
	"grabber/internal/rest/restapi/operations"
	"log"
	"time"

	"github.com/go-openapi/loads"
)

func main() {
	if err := config.SetEnvFromFile(); err != nil {
		log.Fatal("failed to set environment variables from provided file: ", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load configuration: ", err)
	}

	repo, err := repository.NewPostgresRepository(cfg.DbURL)
	if err != nil {
		log.Fatal("repository creation error: ", err)
	}

	err = repo.Migrate()
	if err != nil {
		log.Fatal("database migration application error: ", err)
	}

	// Init and start server
	api := operations.NewGrabberAPI(validateSpec(restapi.SwaggerJSON, restapi.FlatSwaggerJSON))
	server := restapi.NewServer(api)
	defer func() {
		if err := server.Shutdown(); err != nil {
			log.Print("Server shutdown failed with error: ", err)
		}
	}()
	handler := backend.InitAndBindToAPI(repo, grabber.NewHtmlGrabber(), api)
	server.GracefulTimeout = time.Duration(15) * time.Second
	server.SetHandler(handler)
	server.Port = 8701

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

func validateSpec(orig, flat json.RawMessage) *loads.Document {
	swaggerSpec, err := loads.Embedded(orig, flat)
	if err != nil {
		log.Fatalln(err)
	}
	return swaggerSpec
}
