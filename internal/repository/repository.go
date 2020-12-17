//go:generate mockgen -destination=../../mocks/mock_repository.go -package=mocks -source repository.go
package repository

import "context"

type BookmarkRepository interface {
	CreateBookmark(ctx context.Context, bookmark UserBookmark) (UserBookmark, error)

	GetBookmark(ctx context.Context, user string, book string) (UserBookmark, error)
	GetRecentBookmarks(ctx context.Context, user string, limit int64) ([]UserBookmark, error)

	UpdateBookmark(ctx context.Context, bookmark UserBookmark) (UserBookmark, error)

	DeleteBookmark(ctx context.Context, bookmark UserBookmark) (UserBookmark, error)
}
