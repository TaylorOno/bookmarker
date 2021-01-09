package main

import (
	"log"
	"net/http"

	"github.com/TaylorOno/bookmarker/cmd/metrics"

	"github.com/TaylorOno/bookmarker/cmd/config"
	"github.com/TaylorOno/bookmarker/cmd/routes"
	"github.com/TaylorOno/bookmarker/service"
	"github.com/TaylorOno/bookmarker/service/repository"
	"github.com/go-playground/validator/v10"
)

func main() {
	reporter := metrics.NewConsoleReporter()

	session, err := config.NewAWSSessions("id", "secret", "us-west-2", "http://192.168.3.144:8000")
	if err != nil {
		log.Fatal(err.Error())
	}

	dynamoClient := config.NewDynamoClient(session)
	dynamoRepo := repository.NewDynamoRepository(dynamoClient, "bookmarks").AddReporter(reporter)
	bookmarkerService := service.NewBookmarker(dynamoRepo)

	server := routes.Server{BookmarkService: bookmarkerService, Validate: validator.New()}

	router := server.SetRoutes(reporter)

	log.Println("starting server on 8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err.Error())
	}
}
