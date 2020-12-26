package routes_test

import (
	"errors"
	"fmt"
	service2 "github.com/TaylorOno/bookmarker/service"
	"net/http"
	"net/http/httptest"

	"github.com/TaylorOno/bookmarker/cmd/routes"
	"github.com/TaylorOno/bookmarker/service/repository"
	"github.com/TaylorOno/bookmarker/tests/mocks"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BookmarkList", func() {
	var (
		mockCtrl   *gomock.Controller
		bookmarker *mocks.MockBookmarker
		server     *routes.Server
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		bookmarker = mocks.NewMockBookmarker(mockCtrl)
		server = &routes.Server{
			BookmarkService: bookmarker,
			Validate:        validator.New(),
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("GetBookmarks", func() {
		It("Returns 200 if bookmarks are returned", func() {
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			bookmarker.EXPECT().GetBookmarkList(gomock.Any(), gomock.Any()).Return([]service2.Bookmark{{
				Book:                 "book1",
				Series:               "",
				Status:               "IN_PROGRESS",
				Page:                 36,
				AdditionalProperties: nil,
			}}, nil)
			result := httptest.NewRecorder()
			server.GetBookmarks(result, req)
			Expect(result.Result().StatusCode).To(Equal(200))
		})

		It("Returns 400 if filter is invalid", func() {
			req, _ := http.NewRequest(http.MethodGet, "/test?filter=ADVENTURE", nil)
			result := httptest.NewRecorder()
			server.GetBookmarks(result, req)
			Expect(result.Result().StatusCode).To(Equal(400))
		})

		It("Returns 404 if no bookmarks are found", func() {
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			result := httptest.NewRecorder()
			bookmarker.EXPECT().GetBookmarkList(gomock.Any(), gomock.Any()).Return([]service2.Bookmark{}, repository.NotFoundException)
			server.GetBookmarks(result, req)
			Expect(result.Result().StatusCode).To(Equal(404))
		})

		It("Returns 500 if request fails", func() {
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			bookmarker.EXPECT().GetBookmarkList(gomock.Any(), gomock.Any()).Return([]service2.Bookmark{}, errors.New("get error"))
			result := httptest.NewRecorder()
			server.GetBookmarks(result, req)
			Expect(result.Result().StatusCode).To(Equal(500))
		})
	})

	Context("CreateBookmarkListRequest", func() {
		It("returns a BookmarkListRequest", func() {
			req, _ := http.NewRequest(http.MethodGet, "/test-user?limit=7&filter=FINISHED", nil)
			result, err := server.CreateBookmarkListRequest(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(fmt.Sprintf("%T", result)).To(Equal("service.BookmarkListRequest"))
			Expect(result.UserId).To(Equal("test-user"))
			Expect(result.Limit).To(Equal(7))
			Expect(result.Filter).To(Equal("FINISHED"))
		})

		It("sets default limit to 30 if not provided", func() {
			req, _ := http.NewRequest(http.MethodGet, "/test-user?filter=FINISHED", nil)
			result, err := server.CreateBookmarkListRequest(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(fmt.Sprintf("%T", result)).To(Equal("service.BookmarkListRequest"))
			Expect(result.UserId).To(Equal("test-user"))
			Expect(result.Limit).To(Equal(30))
			Expect(result.Filter).To(Equal("FINISHED"))
		})

		It("sets default filter to NONE if not provided", func() {
			req, _ := http.NewRequest(http.MethodGet, "/test-user", nil)
			result, err := server.CreateBookmarkListRequest(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(fmt.Sprintf("%T", result)).To(Equal("service.BookmarkListRequest"))
			Expect(result.UserId).To(Equal("test-user"))
			Expect(result.Limit).To(Equal(30))
			Expect(result.Filter).To(Equal("NONE"))
		})

		It("returns an error if invalid filter", func() {
			req, _ := http.NewRequest(http.MethodGet, "/test-user?limit=7&filter=SPAGHETTI", nil)
			_, err := server.CreateBookmarkListRequest(req)
			Expect(err).To(HaveOccurred())
		})
	})
})
