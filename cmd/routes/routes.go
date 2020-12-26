//go:generate mockgen -destination=../../mocks/mock_bookmarker.go -package=mocks -source routes.go

package routes

import (
	"context"
	"encoding/json"
	"errors"
	service2 "github.com/TaylorOno/bookmarker/service"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var (
	MissingBodyException = errors.New("no request body")
)

type Bookmarker interface {
	SaveBookmark(ctx context.Context, b service2.NewBookmarkRequest) (service2.Bookmark, error)
	DeleteBookmark(ctx context.Context, b service2.DeleteBookmarkRequest) error
	GetBookmark(ctx context.Context, b service2.BookmarkRequest) (service2.Bookmark, error)
	GetBookmarkList(ctx context.Context, b service2.BookmarkListRequest) ([]service2.Bookmark, error)
}

type Server struct {
	BookmarkService Bookmarker
	Validate        *validator.Validate
}

//SetRoutes registers the http routes and handlers
func (s *Server) SetRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/{user}", s.SaveBookmark).Methods(http.MethodPost)
	r.HandleFunc("/{user}/{Book}", s.DeleteBookmark).Methods(http.MethodDelete)
	r.HandleFunc("/{user}", s.GetBookmarks).Methods(http.MethodGet)
	r.HandleFunc("/{user}/{Book}", s.GetBookmark).Methods(http.MethodGet)
	return r
}

func getMethodBody(req *http.Request, request *service2.NewBookmarkRequest) error {
	if req.Body == nil {
		return MissingBodyException
	}
	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(request)
	if err != nil {
		return err
	}
	return nil
}

func getUserID(p string) string {
	paths := strings.Split(p, "/")
	if len(paths) < 1 {
		return ""
	}
	return paths[1]
}

func getBook(p string) string {
	paths := strings.Split(p, "/")
	if len(paths) < 2 {
		return ""
	}
	return paths[2]
}

func getFilter(query url.Values) string {
	filter := query.Get("filter")
	if len(filter) == 0 {
		return "NONE"
	}
	return filter
}

func getLimit(query url.Values) int {
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		return 30
	}
	return limit
}
