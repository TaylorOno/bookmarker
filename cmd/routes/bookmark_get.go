package routes

import (
	"encoding/json"
	"errors"
	service2 "github.com/TaylorOno/bookmarker/service"
	"log"
	"net/http"

	"github.com/TaylorOno/bookmarker/service/repository"
)

//GetBookmark given a user and a book returns a matching bookmark if any.
func (s *Server) GetBookmark(w http.ResponseWriter, req *http.Request) {
	BookmarkRequest := s.CreateBookmarkRequest(req)

	result, err := s.BookmarkService.GetBookmark(req.Context(), BookmarkRequest)
	if err != nil {
		if errors.Is(err, repository.NotFoundException) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(result)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

//CreateBookmarkRequest extracts service.BookmarkRequest from a http.Request.
func (s *Server) CreateBookmarkRequest(req *http.Request) service2.BookmarkRequest {
	return service2.BookmarkRequest{
		UserId: getUserID(req.URL.Path),
		Book:   getBook(req.URL.Path),
	}
}
