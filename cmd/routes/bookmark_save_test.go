package routes_test

import (
	"errors"
	"fmt"
	service2 "github.com/TaylorOno/bookmarker/service"
	"net/http"
	"net/http/httptest"

	"github.com/TaylorOno/bookmarker/cmd/routes"
	"github.com/TaylorOno/bookmarker/tests/mocks"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BookmarkSave", func() {
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

	Context("SaveBookmark", func() {
		It("Returns 200 if successful", func() {
			req, _ := http.NewRequest(http.MethodPost, "/test", bodyFromFile("valid_bookmark.json"))
			bookmarker.EXPECT().SaveBookmark(gomock.Any(), gomock.Any())
			result := httptest.NewRecorder()
			server.SaveBookmark(result, req)
			Expect(result.Result().StatusCode).To(Equal(200))
		})

		It("Returns 400 if missing required fields", func() {
			req, _ := http.NewRequest(http.MethodPost, "/test", bodyFromFile("invalid_bookmark.json"))
			result := httptest.NewRecorder()
			server.SaveBookmark(result, req)
			Expect(result.Result().StatusCode).To(Equal(400))

		})

		It("Returns 500 if fails to save", func() {
			req, _ := http.NewRequest(http.MethodPost, "/test", bodyFromFile("valid_bookmark.json"))
			bookmarker.EXPECT().SaveBookmark(gomock.Any(), gomock.Any()).Return(service2.Bookmark{}, errors.New("save error"))
			result := httptest.NewRecorder()
			server.SaveBookmark(result, req)
			Expect(result.Result().StatusCode).To(Equal(500))
		})
	})

	Context("CreateSaveRequest", func() {
		It("returns a NewBookmarkRequest", func() {
			req, _ := http.NewRequest(http.MethodPost, "/test", bodyFromFile("valid_bookmark.json"))
			result, err := server.CreateSaveRequest(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(fmt.Sprintf("%T", result)).To(Equal("service.NewBookmarkRequest"))
		})

		It("returns an error if missing body", func() {
			req, _ := http.NewRequest(http.MethodPost, "/test", nil)
			_, err := server.CreateSaveRequest(req)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if decode fails body", func() {
			req, _ := http.NewRequest(http.MethodPost, "/test", bodyFromFile("invalid_bookmark.json"))
			_, err := server.CreateSaveRequest(req)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if missing required field", func() {
			req, _ := http.NewRequest(http.MethodPost, "/test", bodyFromFile("missing_field_bookmark.json"))
			_, err := server.CreateSaveRequest(req)
			Expect(err).To(HaveOccurred())
		})
	})
})
