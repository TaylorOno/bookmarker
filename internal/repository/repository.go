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
)

var (
	//NotFoundException special exception for not found
	NotFoundException = errors.New("not found")
)

type BookmarkRepository interface {
	CreateBookmark(ctx context.Context, bookmark UserBookmark) (UserBookmark, error)

	GetBookmark(ctx context.Context, user string, book string) (UserBookmark, error)
	GetBookmarks(ctx context.Context, user string, filter string, limit int64) ([]UserBookmark, error)

	UpdateBookmark(ctx context.Context, bookmark UserBookmark) (UserBookmark, error)

	DeleteBookmark(ctx context.Context, user string, book string) error
}
