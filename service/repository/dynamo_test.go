package repository_test

import (
	"context"
	"errors"
	"time"

	"github.com/TaylorOno/bookmarker/service/repository"
	"github.com/TaylorOno/bookmarker/tests/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dynamo", func() {
	repository.SetTimeOut(100 * time.Millisecond)

	var (
		mockCtrl      *gomock.Controller
		dynamoClient  *mocks.MockDynamoClient
		dynamo        *repository.Dynamo
		putItemOutput = &dynamodb.PutItemOutput{
			ConsumedCapacity: &dynamodb.ConsumedCapacity{
				CapacityUnits: aws.Float64(1),
			},
		}
		deleteItemOutput = &dynamodb.DeleteItemOutput{
			ConsumedCapacity: &dynamodb.ConsumedCapacity{
				CapacityUnits: aws.Float64(1),
			},
		}
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		dynamoClient = mocks.NewMockDynamoClient(mockCtrl)
		dynamo = &repository.Dynamo{
			Client:    dynamoClient,
			TableName: "test",
		}
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

	Context("CreateBookmark", func() {
		It("Calls dynamo save with user bookmark", func() {
			var argumentCapture *dynamodb.PutItemInput
			bookmark := repository.UserBookmark{UserId: "user", Book: "book"}
			dynamoClient.EXPECT().
				PutItemWithContext(gomock.Any(), gomock.Any()).
				DoAndReturn(func(a aws.Context, b *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
					argumentCapture = b
					return putItemOutput, nil
				})
			dynamo.CreateBookmark(context.Background(), bookmark)
			Expect(*argumentCapture.TableName).To(Equal("test"))
		})

		It("Returns error if save fails", func() {
			bookmark := repository.UserBookmark{UserId: "user", Book: "book"}
			dynamoClient.EXPECT().PutItemWithContext(gomock.Any(), gomock.Any()).Return(putItemOutput, errors.New("exception"))
			_, err := dynamo.CreateBookmark(context.Background(), bookmark)
			Expect(err).To(HaveOccurred())
		})

		It("Returns error if save times out", func() {
			bookmark := repository.UserBookmark{UserId: "user", Book: "book"}
			dynamoClient.EXPECT().
				PutItemWithContext(gomock.Any(), gomock.Any()).
				DoAndReturn(func(a context.Context, b *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
					time.Sleep(150 * time.Millisecond)
					select {
					case <-a.Done():
						return nil, a.Err()
					default:
						return nil, nil
					}
				})
			_, err := dynamo.CreateBookmark(context.Background(), bookmark)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("UpdateBookmark", func() {
		It("Calls dynamo put with user bookmark", func() {

		})
	})

	Context("DeleteBookmark", func() {
		It("Calls dynamo delete with user and book", func() {
			var argumentCapture *dynamodb.DeleteItemInput
			dynamoClient.EXPECT().
				DeleteItemWithContext(gomock.Any(), gomock.Any()).
				DoAndReturn(func(a aws.Context, b *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
					argumentCapture = b
					return deleteItemOutput, nil
				})

			err := dynamo.DeleteBookmark(context.Background(), "test", "book")
			Expect(err).ToNot(HaveOccurred())
			Expect(*argumentCapture.TableName).To(Equal("test"))
			Expect(*(*argumentCapture.Key["UserId"]).S).To(Equal("test"))
			Expect(*(*argumentCapture.Key["Book"]).S).To(Equal("book"))
		})

		It("Returns error if delete fails", func() {
			dynamoClient.EXPECT().DeleteItemWithContext(gomock.Any(), gomock.Any()).Return(deleteItemOutput, errors.New("exception"))
			err := dynamo.DeleteBookmark(context.Background(), "test", "user")
			Expect(err).To(HaveOccurred())
		})

		It("Returns error if get times out", func() {
			dynamoClient.EXPECT().
				DeleteItemWithContext(gomock.Any(), gomock.Any()).
				DoAndReturn(func(a context.Context, b *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
					time.Sleep(150 * time.Millisecond)
					select {
					case <-a.Done():
						return nil, a.Err()
					default:
						return nil, nil
					}
				})
			err := dynamo.DeleteBookmark(context.Background(), "test", "book")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("GetBookmark", func() {
		It("Calls dynamo get with user and book", func() {
			var argumentCapture *dynamodb.GetItemInput
			getItemOutput := testGetResponse(&repository.UserBookmark{UserId: "test", Book: "book"})

			dynamoClient.EXPECT().
				GetItemWithContext(gomock.Any(), gomock.Any()).
				DoAndReturn(func(a aws.Context, b *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
					argumentCapture = b
					return getItemOutput, nil
				})
			result, err := dynamo.GetBookmark(context.Background(), "test", "book")
			Expect(err).ToNot(HaveOccurred())
			Expect(*argumentCapture.TableName).To(Equal("test"))
			Expect(*(*argumentCapture.Key["UserId"]).S).To(Equal("test"))
			Expect(*(*argumentCapture.Key["Book"]).S).To(Equal("book"))
			Expect(result.UserId).To(Equal("test"))
		})

		It("Returns error if get fails", func() {
			getItemOutput := testGetResponse(nil)
			dynamoClient.EXPECT().GetItemWithContext(gomock.Any(), gomock.Any()).Return(getItemOutput, errors.New("exception"))
			_, err := dynamo.GetBookmark(context.Background(), "test", "book")
			Expect(err).To(HaveOccurred())
		})

		It("Returns error if item is not found", func() {
			getItemOutput := testGetResponse(nil)
			dynamoClient.EXPECT().GetItemWithContext(gomock.Any(), gomock.Any()).Return(getItemOutput, nil)
			_, err := dynamo.GetBookmark(context.Background(), "test", "book")
			Expect(err).To(HaveOccurred())
		})

		It("Returns error if get times out", func() {
			dynamoClient.EXPECT().
				GetItemWithContext(gomock.Any(), gomock.Any()).
				DoAndReturn(func(a context.Context, b *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
					time.Sleep(150 * time.Millisecond)
					select {
					case <-a.Done():
						return nil, a.Err()
					default:
						return nil, nil
					}
				})
			_, err := dynamo.GetBookmark(context.Background(), "test", "book")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("GetBookmarks", func() {
		It("Calls dynamo query with user and filter", func() {
			var argumentCapture *dynamodb.QueryInput
			getQueryOutput := testQueryResponse()

			dynamoClient.EXPECT().
				QueryWithContext(gomock.Any(), gomock.Any()).
				DoAndReturn(func(a aws.Context, b *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
					argumentCapture = b
					return getQueryOutput, nil
				})
			_, err := dynamo.GetBookmarks(context.Background(), "test", "NONE", 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(*argumentCapture.TableName).To(Equal("test"))
			Expect(*argumentCapture.IndexName).To(Equal("History"))
			Expect(argumentCapture.FilterExpression).To(BeNil())
		})

		It("Calls dynamo query with user and filter FINISHED", func() {
			var argumentCapture *dynamodb.QueryInput
			getQueryOutput := testQueryResponse()

			dynamoClient.EXPECT().
				QueryWithContext(gomock.Any(), gomock.Any()).
				DoAndReturn(func(a aws.Context, b *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
					argumentCapture = b
					return getQueryOutput, nil
				})
			_, err := dynamo.GetBookmarks(context.Background(), "test", "FINISHED", 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(*argumentCapture.TableName).To(Equal("test"))
			Expect(*argumentCapture.IndexName).To(Equal("History"))
			Expect(*argumentCapture.FilterExpression).To(Equal("#0 = :0"))
			Expect(*argumentCapture.ExpressionAttributeValues[":0"].S).To(Equal("FINISHED"))
		})

		It("Calls dynamo query with user and filter FINISHED", func() {
			var argumentCapture *dynamodb.QueryInput
			getQueryOutput := testQueryResponse()

			dynamoClient.EXPECT().
				QueryWithContext(gomock.Any(), gomock.Any()).
				DoAndReturn(func(a aws.Context, b *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
					argumentCapture = b
					return getQueryOutput, nil
				})
			_, err := dynamo.GetBookmarks(context.Background(), "test", "IN_PROGRESS", 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(*argumentCapture.TableName).To(Equal("test"))
			Expect(*argumentCapture.IndexName).To(Equal("History"))
			Expect(*argumentCapture.FilterExpression).To(Equal("#0 = :0"))
			Expect(*argumentCapture.ExpressionAttributeValues[":0"].S).To(Equal("IN_PROGRESS"))
		})

		It("Returns error if query fails", func() {
			getQueryOutput := testQueryResponse()
			dynamoClient.EXPECT().QueryWithContext(gomock.Any(), gomock.Any()).Return(getQueryOutput, errors.New("exception"))
			_, err := dynamo.GetBookmarks(context.Background(), "test", "NONE", 1)
			Expect(err).To(HaveOccurred())
		})
	})
})

func testGetResponse(item *repository.UserBookmark) *dynamodb.GetItemOutput {
	av, _ := dynamodbattribute.MarshalMap(item)
	return &dynamodb.GetItemOutput{
		Item: av,
		ConsumedCapacity: &dynamodb.ConsumedCapacity{
			CapacityUnits: aws.Float64(1),
		},
	}
}

func testQueryResponse() *dynamodb.QueryOutput {
	bookmark, _ := dynamodbattribute.MarshalMap(repository.UserBookmark{UserId: "test", Book: "book"})
	av := []map[string]*dynamodb.AttributeValue{bookmark}
	return &dynamodb.QueryOutput{
		Items: av,
		ConsumedCapacity: &dynamodb.ConsumedCapacity{
			CapacityUnits: aws.Float64(1),
		},
	}
}
