package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/TaylorOno/bookmarker/internal/repository"
	"github.com/TaylorOno/bookmarker/internal/service"
)

//GetBookmark given a user and a book returns a matching bookmark if any.
func (s *Server) GetBookmark(w http.ResponseWriter, req *http.Request) {
	BookmarkRequest, err := s.CreateBookmarkRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := s.BookmarkService.GetBookmark(req.Context(), BookmarkRequest)
	if err != nil {
		if errors.Is(err, repository.NotFoundException) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(result)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

//CreateBookmarkRequest extracts service.BookmarkRequest from a http.Request.
func (s *Server) CreateBookmarkRequest(req *http.Request) (service.BookmarkRequest, error) {
	return service.BookmarkRequest{
		UserId: getUserID(req.URL.Path),
		Book:   getBook(req.URL.Path),
	}, nil
}
