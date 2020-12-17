// Code generated by MockGen. DO NOT EDIT.
// Source: routes.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	service "github.com/TaylorOno/bookmarker/internal/service"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockBookmarker is a mock of Bookmarker interface
type MockBookmarker struct {
	ctrl     *gomock.Controller
	recorder *MockBookmarkerMockRecorder
}

// MockBookmarkerMockRecorder is the mock recorder for MockBookmarker
type MockBookmarkerMockRecorder struct {
	mock *MockBookmarker
}

// NewMockBookmarker creates a new mock instance
func NewMockBookmarker(ctrl *gomock.Controller) *MockBookmarker {
	mock := &MockBookmarker{ctrl: ctrl}
	mock.recorder = &MockBookmarkerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBookmarker) EXPECT() *MockBookmarkerMockRecorder {
	return m.recorder
}

// SaveBookmark mocks base method
func (m *MockBookmarker) SaveBookmark(ctx context.Context, b service.NewBookmarkRequest) (service.Bookmark, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveBookmark", ctx, b)
	ret0, _ := ret[0].(service.Bookmark)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveBookmark indicates an expected call of SaveBookmark
func (mr *MockBookmarkerMockRecorder) SaveBookmark(ctx, b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveBookmark", reflect.TypeOf((*MockBookmarker)(nil).SaveBookmark), ctx, b)
}

// DeleteBookmark mocks base method
func (m *MockBookmarker) DeleteBookmark(ctx context.Context, b service.DeleteBookmarkRequest) (service.Bookmark, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBookmark", ctx, b)
	ret0, _ := ret[0].(service.Bookmark)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteBookmark indicates an expected call of DeleteBookmark
func (mr *MockBookmarkerMockRecorder) DeleteBookmark(ctx, b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBookmark", reflect.TypeOf((*MockBookmarker)(nil).DeleteBookmark), ctx, b)
}

// GetBookmark mocks base method
func (m *MockBookmarker) GetBookmark(ctx context.Context, b service.BookmarkRequest) (service.Bookmark, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBookmark", ctx, b)
	ret0, _ := ret[0].(service.Bookmark)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBookmark indicates an expected call of GetBookmark
func (mr *MockBookmarkerMockRecorder) GetBookmark(ctx, b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBookmark", reflect.TypeOf((*MockBookmarker)(nil).GetBookmark), ctx, b)
}

// GetBookmarkList mocks base method
func (m *MockBookmarker) GetBookmarkList(ctx context.Context, b service.BookmarkListRequest) ([]service.Bookmark, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBookmarkList", ctx, b)
	ret0, _ := ret[0].([]service.Bookmark)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBookmarkList indicates an expected call of GetBookmarkList
func (mr *MockBookmarkerMockRecorder) GetBookmarkList(ctx, b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBookmarkList", reflect.TypeOf((*MockBookmarker)(nil).GetBookmarkList), ctx, b)
}
