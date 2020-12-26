package routes_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/TaylorOno/bookmarker/tests/mocks"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/TaylorOno/bookmarker/cmd/routes"
)

var _ = Describe("BookmarkDelete", func() {
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

	Context("Delete Bookmark", func() {
		It("Returns 200 if bookmark deleted returned", func() {
			req, _ := http.NewRequest(http.MethodGet, "/test/book1", nil)
			bookmarker.EXPECT().DeleteBookmark(gomock.Any(), gomock.Any()).Return(nil)
			result := httptest.NewRecorder()
			server.DeleteBookmark(result, req)
			Expect(result.Result().StatusCode).To(Equal(200))
		})

		It("Returns 500 if request fails", func() {
			req, _ := http.NewRequest(http.MethodGet, "/test/book999", nil)
			bookmarker.EXPECT().DeleteBookmark(gomock.Any(), gomock.Any()).Return(errors.New("get error"))
			result := httptest.NewRecorder()
			server.DeleteBookmark(result, req)
			Expect(result.Result().StatusCode).To(Equal(500))
		})
	})
})
