package service_test

import (
	"context"
	"github.com/TaylorOno/bookmarker/internal/service"
	"github.com/TaylorOno/bookmarker/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Service", func() {
	var (
		mockCtrl          *gomock.Controller
		repository        *mocks.MockBookmarkRepository
		bookermarkService *service.Service
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		repository = mocks.NewMockBookmarkRepository(mockCtrl)
		bookermarkService = &service.Service{
			Repo: repository,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("SaveBookmark", func() {
		It("Calls repository CreateBookmark", func() {
			var b service.NewBookmarkRequest
			repository.EXPECT().CreateBookmark(gomock.Any(), gomock.Any())
			bookermarkService.SaveBookmark(context.Background(), b)
		})
	})

	Context("GetBookmark", func() {
		It("Calls repository GetBookmark", func() {
			var b service.BookmarkRequest
			repository.EXPECT().GetBookmark(gomock.Any(), gomock.Any(), gomock.Any())
			bookermarkService.GetBookmark(context.Background(), b)
		})
	})

	Context("GetBookmarkList", func() {
		It("Calls repository GetBookmarks", func() {
			var b service.BookmarkListRequest
			repository.EXPECT().GetRecentBookmarks(gomock.Any(), gomock.Any(), gomock.Any())
			bookermarkService.GetBookmarkList(context.Background(), b)
		})
	})
})
