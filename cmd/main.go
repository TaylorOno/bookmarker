package main

import (
	"log"
	"net/http"

	"github.com/TaylorOno/bookmarker/cmd/config"
	"github.com/TaylorOno/bookmarker/cmd/routes"
	"github.com/TaylorOno/bookmarker/internal/repository"
	"github.com/TaylorOno/bookmarker/internal/service"
	"github.com/go-playground/validator/v10"
)

func main() {
	session, err := config.CreateAWSSessions("id", "secret", "us-west-1", "http://localhost:8000")
	if err != nil {
		log.Fatal(err.Error())
	}

	repo := repository.CreateDynamoRepository(session, "bookmarks")
	config.CreateTableIfNotExist(repo)

	bookmarkerService := &service.Service{
		Repo: repo,
	}

	server := routes.Server{
		BookmarkService: bookmarkerService,
		Validate:        validator.New(),
	}
	router := server.SetRoutes()

	http.ListenAndServe(":8080", router)
}
