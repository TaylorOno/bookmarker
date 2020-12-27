package repository_test

import (
	"github.com/TaylorOno/bookmarker/service/repository"
	"github.com/TaylorOno/bookmarker/tests/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DynamoConfig", func() {
	var (
		mockCtrl     *gomock.Controller
		dynamoClient *mocks.MockDynamoClient
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		dynamoClient = mocks.NewMockDynamoClient(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("NewDynamoRepository", func() {
		It("Checks or Creates Table and returns a Dynamo Client", func() {
			var client *repository.Dynamo
			dynamoClient.EXPECT().CreateTable(gomock.Any())
			dynamoClient.EXPECT().WaitUntilTableExistsWithContext(gomock.Any(), gomock.Any(), gomock.Any())
			client = repository.NewDynamoRepository(dynamoClient, "testTable")
			Expect(client).ToNot(BeNil())
		})
	})
})
