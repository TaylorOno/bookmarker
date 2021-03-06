//go:generate mockgen -destination=../../mocks/mock_repository.go -package=mocks -source repository.go
package repository

import (
	"context"
	"errors"
)

const (
	_inProgressfilter = "IN_PROGRESS"
	_finishedfilter   = "FINISHED"
	_nonefilter       = "NONE"

	_add    = "add"
	_get    = "get"
	_update = "update"
	_delete = "delete"
	_query  = "query"
)

var (
	//NotFoundException special exception for not found
	NotFoundException = errors.New("not found")
)

type BookmarkRepository interface {
	CreateBookmark(ctx context.Context, bookmark UserBookmark) (UserBookmark, error)
	GetBookmark(ctx context.Context, user string, book string) (UserBookmark, error)
	GetBookmarks(ctx context.Context, user string, filter string, limit int) ([]UserBookmark, error)
	UpdateBookmark(ctx context.Context, bookmark UserBookmark) (UserBookmark, error)
	DeleteBookmark(ctx context.Context, user string, book string) error
}

type bookmarkKey struct {
	UserId string
	Book   string
}

type UserBookmark struct {
	UserId      string
	LastUpdated string
	Book        string
	Series      string
	Status      string
	Page        int

	//AdditionalProperties provide for extendable data model there are no guarantees on any fields provided.  Data will
	//be projected into secondary indexes so be cautious of field size.
	AdditionalProperties map[string]interface{}
}
