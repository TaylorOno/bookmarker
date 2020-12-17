package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/TaylorOno/bookmarker/internal/repository"
	"github.com/TaylorOno/bookmarker/internal/service"
)

//GetBookmarks returns a lists of bookmarks for a given user.
//filter query param for IN_PROGRESS of FINISHED.
//limit query param can be used to limit the items that are returned (default 30).
func (s *Server) GetBookmarks(w http.ResponseWriter, req *http.Request) {
	BookmarkListRequest, err := s.CreateBookmarkListRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := s.BookmarkService.GetBookmarkList(req.Context(), BookmarkListRequest)
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

//CreateBookmarkListRequest extracts service.BookmarkListRequest from a http.Request.
func (s *Server) CreateBookmarkListRequest(req *http.Request) (service.BookmarkListRequest, error) {
	request := service.BookmarkListRequest{
		UserId: getUserID(req.URL.Path),
		Limit:  getLimit(req.URL.Query()),
		Filter: getFilter(req.URL.Query()),
	}

	err := s.Validate.Struct(request)
	if err != nil {
		return request, err
	}

	return request, nil
}
