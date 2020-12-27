//go:generate mockgen -destination=../../tests/mocks/mock_bookmarker.go -package=mocks -source routes.go

package routes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/TaylorOno/bookmarker/cmd/middleware"
	"github.com/TaylorOno/bookmarker/service"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	MissingBodyException = errors.New("no request body")
)

type Bookmarker interface {
	SaveBookmark(ctx context.Context, b service.NewBookmarkRequest) (service.Bookmark, error)
	DeleteBookmark(ctx context.Context, b service.DeleteBookmarkRequest) error
	GetBookmark(ctx context.Context, b service.BookmarkRequest) (service.Bookmark, error)
	GetBookmarkList(ctx context.Context, b service.BookmarkListRequest) ([]service.Bookmark, error)
}

type Server struct {
	BookmarkService Bookmarker
	Validate        *validator.Validate
}

type Reporter interface {
	ObserverHistogram(name string, value float64, labels ...string)
	ObserverSummary(name string, value float64, labels ...string)
}

//SetRoutes registers the http routes and handlers
func (s *Server) SetRoutes(reporter Reporter) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.NewInboundObserver(reporter))
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/{user}", s.SaveBookmark).Methods(http.MethodPost)
	r.HandleFunc("/{user}/{Book}", s.DeleteBookmark).Methods(http.MethodDelete)
	r.HandleFunc("/{user}", s.GetBookmarks).Methods(http.MethodGet)
	r.HandleFunc("/{user}/{Book}", s.GetBookmark).Methods(http.MethodGet)
	return r
}

func getMethodBody(req *http.Request, request *service.NewBookmarkRequest) error {
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
