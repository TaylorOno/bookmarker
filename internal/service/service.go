package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/TaylorOno/bookmarker/internal/repository"
)

type Service struct {
	Repo repository.BookmarkRepository
}

func (s *Service) SaveBookmark(ctx context.Context, b NewBookmarkRequest) (Bookmark, error) {
	var bookmark Bookmark
	newBookmark := repository.UserBookmark{
		UserId:               b.UserId,
		LastUpdated:          time.Now().UTC().Format("2006-01-02T15:04:05Z07:00.000"),
		Book:                 strings.ToLower(b.Book),
		Series:               b.Series,
		Status:               b.Status,
		Page:                 b.Page,
		AdditionalProperties: b.AdditionalProperties,
	}
	result, err := s.Repo.CreateBookmark(ctx, newBookmark)
	if err != nil {
		return bookmark, err
	}

	bookmark = Bookmark{
		Book:                 result.Book,
		Series:               result.Series,
		Status:               result.Status,
		Page:                 result.Page,
		AdditionalProperties: result.AdditionalProperties,
	}
	return bookmark, nil
}
func (s *Service) DeleteBookmark(ctx context.Context, b DeleteBookmarkRequest) (Bookmark, error) {
	var bookmark Bookmark
	return bookmark, errors.New("not implemented")
}
func (s *Service) GetBookmark(ctx context.Context, b BookmarkRequest) (Bookmark, error) {
	var bookmark Bookmark
	userBookmark, err := s.Repo.GetBookmark(ctx, b.UserId, strings.ToLower(b.Book))
	if err != nil {
		return bookmark, err
	}

	return newBookmark(userBookmark), nil
}

func newBookmark(userBookmark repository.UserBookmark) Bookmark {
	return Bookmark{
		Book:                 userBookmark.Book,
		LastUpdated:          userBookmark.LastUpdated,
		Series:               userBookmark.Series,
		Status:               userBookmark.Status,
		Page:                 userBookmark.Page,
		AdditionalProperties: userBookmark.AdditionalProperties,
	}
}
func (s *Service) GetBookmarkList(ctx context.Context, b BookmarkListRequest) ([]Bookmark, error) {
	userBookmarks, err := s.Repo.GetRecentBookmarks(ctx, b.UserId, b.Limit)
	if err != nil {
		return []Bookmark{}, err
	}

	bookmarks := make([]Bookmark, 0, len(userBookmarks))
	for _, b := range userBookmarks {
		bookmarks = append(bookmarks, newBookmark(b))
	}

	return bookmarks, nil
}
