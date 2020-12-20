package service

import (
	"context"
	"strings"
	"time"

	"github.com/TaylorOno/bookmarker/internal/repository"
)

type Service struct {
	Repo repository.BookmarkRepository
}

func (s *Service) SaveBookmark(ctx context.Context, b NewBookmarkRequest) (Bookmark, error) {
	result, err := s.Repo.CreateBookmark(ctx, createUserBookmark(b))
	if err != nil {
		return Bookmark{}, err
	}

	return newBookmark(result), nil
}

func createUserBookmark(b NewBookmarkRequest) repository.UserBookmark {
	return repository.UserBookmark{
		UserId:               b.UserId,
		LastUpdated:          time.Now().UTC().Format("2006-01-02T15:04:05Z07:00.000"),
		Book:                 strings.ToLower(b.Book),
		Series:               b.Series,
		Status:               b.Status,
		Page:                 b.Page,
		AdditionalProperties: b.AdditionalProperties,
	}
}

func (s *Service) DeleteBookmark(ctx context.Context, b DeleteBookmarkRequest) error {
	err := s.Repo.DeleteBookmark(ctx, b.UserId, strings.ToLower(b.Book))
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetBookmark(ctx context.Context, b BookmarkRequest) (Bookmark, error) {
	userBookmark, err := s.Repo.GetBookmark(ctx, b.UserId, strings.ToLower(b.Book))
	if err != nil {
		return Bookmark{}, err
	}

	return newBookmark(userBookmark), nil
}

func (s *Service) GetBookmarkList(ctx context.Context, b BookmarkListRequest) ([]Bookmark, error) {
	userBookmarks, err := s.Repo.GetBookmarks(ctx, b.UserId, b.Filter, b.Limit)
	if err != nil {
		return []Bookmark{}, err
	}

	return toBookMarks(userBookmarks), nil
}

func toBookMarks(userBookmarks []repository.UserBookmark) []Bookmark {
	bookmarks := make([]Bookmark, 0, len(userBookmarks))
	for _, b := range userBookmarks {
		bookmarks = append(bookmarks, newBookmark(b))
	}
	return bookmarks
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
