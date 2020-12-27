package middleware_test

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/TaylorOno/bookmarker/cmd/middleware"
	"github.com/TaylorOno/bookmarker/tests/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metrics", func() {
	var (
		mockCtrl     *gomock.Controller
		mockReporter *mocks.MockReporter
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockReporter = mocks.NewMockReporter(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("NewInboundObserver", func() {
		It("Creates a MiddlewareFunc", func() {
			result := middleware.NewInboundObserver(mockReporter)
			Expect(reflect.TypeOf(result).Name()).To(Equal("MiddlewareFunc"))
		})
	})

	Context("Middleware", func() {
		It("Returns a HandlerFunc", func() {
			result := middleware.NewInboundObserver(mockReporter)
			handler := result(http.NotFoundHandler())
			Expect(reflect.TypeOf(handler).Name()).To(Equal("HandlerFunc"))
		})
	})

	Context("Handler", func() {
		It("Calls observer methods", func() {
			mockReporter.EXPECT().ObserverHistogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
			mockReporter.EXPECT().ObserverSummary(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
			result := middleware.NewInboundObserver(mockReporter)
			handler := result(http.NotFoundHandler())
			req, _ := http.NewRequest("get", "/test", nil)
			r := mux.NewRouter()
			r.Handle("/test", handler)
			r.ServeHTTP(httptest.NewRecorder(), req)
		})
	})
})
