package main

import (
	"log"
	"net/http"

	"github.com/TaylorOno/bookmarker/cmd/config"
	"github.com/TaylorOno/bookmarker/cmd/metrics"
	"github.com/TaylorOno/bookmarker/cmd/routes"
	"github.com/TaylorOno/bookmarker/service"
	"github.com/TaylorOno/bookmarker/service/repository"
	"github.com/go-playground/validator/v10"
)

func main() {
	session, err := config.NewAWSSessions("id", "secret", "us-west-2", "http://localhost:8000")
	if err != nil {
		log.Fatal(err.Error())
	}

	reporter:= metrics.NewConsoleReporter()
	dynamoClient := config.NewDynamoClient(session)
	repository := repository.NewDynamoRepository(dynamoClient, "bookmarks").AddReporter(reporter)

	bookmarkerService := &service.Service{
		Repo: repository,
	}

	server := routes.Server{
		BookmarkService: bookmarkerService,
		Validate:        validator.New(),
	}
	router := server.SetRoutes(reporter)

	http.ListenAndServe(":8080", router)
}
