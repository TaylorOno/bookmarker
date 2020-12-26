package routes

import (
	service2 "github.com/TaylorOno/bookmarker/service"
	"log"
	"net/http"
)

//DeleteBookmark removes a bookmark from the users list if it exists.
func (s *Server) DeleteBookmark(w http.ResponseWriter, req *http.Request) {
	DeleteBookmarkRequest := s.DeleteBookmarkRequest(req)

	err := s.BookmarkService.DeleteBookmark(req.Context(), DeleteBookmarkRequest)
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

//DeleteBookmarkRequest extracts service.DeleteBookmarkRequest from a http.Request.
func (s *Server) DeleteBookmarkRequest(req *http.Request) service2.DeleteBookmarkRequest {
	return service2.DeleteBookmarkRequest{
		UserId: getUserID(req.URL.Path),
		Book:   getBook(req.URL.Path),
	}
}
