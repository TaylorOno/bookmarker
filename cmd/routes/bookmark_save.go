package routes

import (
	"encoding/json"
	service2 "github.com/TaylorOno/bookmarker/service"
	"log"
	"net/http"
)

//SaveBookmark saves a users bookmark object.
func (s *Server) SaveBookmark(w http.ResponseWriter, req *http.Request) {
	var newBookmarkRequest service2.NewBookmarkRequest
	newBookmarkRequest, err := s.CreateSaveRequest(req)
	if err != nil {
		errorBody, _ := json.Marshal(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorBody)
		return
	}

	_, err = s.BookmarkService.SaveBookmark(req.Context(), newBookmarkRequest)
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

//CreateSaveRequest extracts service.NewBookmarkRequest from a http.Request.
func (s *Server) CreateSaveRequest(req *http.Request) (service2.NewBookmarkRequest, error) {
	var request service2.NewBookmarkRequest
	userId := getUserID(req.URL.Path)
	err := getMethodBody(req, &request)
	if err != nil {
		return request, err
	}

	request.UserId = userId
	err = s.Validate.Struct(request)
	if err != nil {
		return request, err
	}

	return request, nil
}
