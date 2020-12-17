package routes_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/TaylorOno/bookmarker/cmd/routes"
	"github.com/TaylorOno/bookmarker/internal/repository"
	"github.com/TaylorOno/bookmarker/internal/service"
	"github.com/TaylorOno/bookmarker/mocks"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BookmarkGet", func() {

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

	Context("Get Bookmark", func() {
		It("Returns 200 if bookmark is returned", func() {
			req, _ := http.NewRequest(http.MethodGet, "/test/book1", nil)
			bookmarker.EXPECT().GetBookmark(gomock.Any(), gomock.Any()).Return(service.Bookmark{
				Book:                 "book1",
				Series:               "",
				Status:               "IN_PROGRESS",
				Page:                 36,
				AdditionalProperties: nil,
			}, nil)
			result := httptest.NewRecorder()
			server.GetBookmark(result, req)
			Expect(result.Result().StatusCode).To(Equal(200))
		})

		It("Returns 404 if no bookmark is found", func() {
			req, _ := http.NewRequest(http.MethodGet, "/test/book0", nil)
			result := httptest.NewRecorder()
			bookmarker.EXPECT().GetBookmark(gomock.Any(), gomock.Any()).Return(service.Bookmark{}, repository.NotFoundException)
			server.GetBookmark(result, req)
			Expect(result.Result().StatusCode).To(Equal(404))
		})

		It("Returns 500 if request fails", func() {
			req, _ := http.NewRequest(http.MethodGet, "/test/book999", nil)
			bookmarker.EXPECT().GetBookmark(gomock.Any(), gomock.Any()).Return(service.Bookmark{}, errors.New("get error"))
			result := httptest.NewRecorder()
			server.GetBookmark(result, req)
			Expect(result.Result().StatusCode).To(Equal(500))
		})
	})
})
