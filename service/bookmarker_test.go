package service_test

import (
	"context"
	"errors"

	"github.com/TaylorOno/bookmarker/service"
	"github.com/TaylorOno/bookmarker/service/repository"
	"github.com/TaylorOno/bookmarker/tests/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	var (
		mockCtrl          *gomock.Controller
		mockRepo          *mocks.MockBookmarkRepository
		bookermarkService *service.Service
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockBookmarkRepository(mockCtrl)
		bookermarkService = &service.Service{
			Repo: mockRepo,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("SaveBookmark", func() {
		It("Calls mockRepo CreateBookmark", func() {
			var b service.NewBookmarkRequest
			mockRepo.EXPECT().CreateBookmark(gomock.Any(), gomock.Any())
			bookermarkService.SaveBookmark(context.Background(), b)
		})

		It("Returns error if repository fails", func() {
			var b service.NewBookmarkRequest
			mockRepo.EXPECT().CreateBookmark(gomock.Any(), gomock.Any()).Return(repository.UserBookmark{}, errors.New("error"))
			_, err := bookermarkService.SaveBookmark(context.Background(), b)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("DeleteBookmark", func() {
		It("Calls mockRepo DeleteBookmark", func() {
			var b service.DeleteBookmarkRequest
			mockRepo.EXPECT().DeleteBookmark(gomock.Any(), gomock.Any(), gomock.Any())
			bookermarkService.DeleteBookmark(context.Background(), b)
		})

		It("Returns error if repository fails", func() {
			var b service.DeleteBookmarkRequest
			mockRepo.EXPECT().DeleteBookmark(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
			err := bookermarkService.DeleteBookmark(context.Background(), b)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("GetBookmark", func() {
		It("Calls mockRepo GetBookmark", func() {
			var b service.BookmarkRequest
			mockRepo.EXPECT().GetBookmark(gomock.Any(), gomock.Any(), gomock.Any())
			bookermarkService.GetBookmark(context.Background(), b)
		})

		It("Returns error if repository fails", func() {
			var b service.BookmarkRequest
			mockRepo.EXPECT().GetBookmark(gomock.Any(), gomock.Any(), gomock.Any()).Return(repository.UserBookmark{}, errors.New("error"))
			_, err := bookermarkService.GetBookmark(context.Background(), b)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("GetBookmarkList", func() {
		It("Calls mockRepo GetBookmarks", func() {
			var b service.BookmarkListRequest
			mockRepo.EXPECT().GetBookmarks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]repository.UserBookmark{{UserId: "test"}}, nil)
			bookermarkService.GetBookmarkList(context.Background(), b)
		})

		It("Returns error if repository fails", func() {
			var b service.BookmarkListRequest
			mockRepo.EXPECT().GetBookmarks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]repository.UserBookmark{}, errors.New("error"))
			_, err := bookermarkService.GetBookmarkList(context.Background(), b)
			Expect(err).To(HaveOccurred())
		})
	})
})
