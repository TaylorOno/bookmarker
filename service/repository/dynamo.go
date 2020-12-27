//go:generate mockgen -destination=../../tests/mocks/mock_dynamo.go -package=mocks -source dynamo.go

package repository

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	_dynamoLatencyHistogram = "dynamo_latency_histogram"
	_dynamoLatencySummary   = "dynamo_latency_summary"
	_dynamoCapacityUsed     = "dynamo_capacity_used"
)

var (
	//maxPageSize page size for queries
	maxPageSize = 100
)

type DynamoClient interface {
	PutItemWithContext(aws.Context, *dynamodb.PutItemInput, ...request.Option) (*dynamodb.PutItemOutput, error)
	QueryWithContext(aws.Context, *dynamodb.QueryInput, ...request.Option) (*dynamodb.QueryOutput, error)
	GetItemWithContext(aws.Context, *dynamodb.GetItemInput, ...request.Option) (*dynamodb.GetItemOutput, error)
	DeleteItemWithContext(aws.Context, *dynamodb.DeleteItemInput, ...request.Option) (*dynamodb.DeleteItemOutput, error)
	CreateTable(*dynamodb.CreateTableInput) (*dynamodb.CreateTableOutput, error)
	WaitUntilTableExistsWithContext(aws.Context, *dynamodb.DescribeTableInput, ...request.WaiterOption) error
}

type DynamoReporter interface {
	ObserverHistogram(name string, value float64, labels ...string)
	ObserverCount(name string, value float64, labels ...string)
	ObserverSummary(name string, value float64, labels ...string)
}

//Dynamo holds dynamo client and table name
type Dynamo struct {
	Client    DynamoClient
	Reporter  DynamoReporter
	TableName string
}

//CreateBookmark adds a bookmark for a user
func (d *Dynamo) CreateBookmark(ctx context.Context, bookmark UserBookmark) (UserBookmark, error) {
	start := time.Now()
	createItemInput, err := d.createItemInput(bookmark)
	if err != nil {
		return bookmark, err
	}

	putResult, err := d.Client.PutItemWithContext(ctx, createItemInput)
	if err != nil {
		return bookmark, err
	}

	d.observer(start, *putResult.ConsumedCapacity.CapacityUnits, _add)
	return bookmark, nil
}

func (d *Dynamo) createItemInput(bookmark UserBookmark) (*dynamodb.PutItemInput, error) {
	var itemInput *dynamodb.PutItemInput

	av, err := dynamodbattribute.MarshalMap(bookmark)
	if err != nil {
		return itemInput, err
	}

	return &dynamodb.PutItemInput{
		Item:                   av,
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(d.TableName),
	}, nil
}

//UpdateBookmark updates a users bookmark with new values this is done by deleting the old bookmark and inserting the new one
func (d *Dynamo) UpdateBookmark(ctx context.Context, bookmark UserBookmark) (UserBookmark, error) {
	return bookmark, errors.New("not implemented")
}

//DeleteBookmark deletes a users bookmark
func (d *Dynamo) DeleteBookmark(ctx context.Context, user string, book string) error {
	start := time.Now()
	deleteItemInput := d.createDeleteItemInput(user, book)

	deleteResult, err := d.Client.DeleteItemWithContext(ctx, deleteItemInput)
	if err != nil {
		return err
	}

	d.observer(start, *deleteResult.ConsumedCapacity.CapacityUnits, _delete)
	return nil
}

func (d *Dynamo) createDeleteItemInput(user string, book string) *dynamodb.DeleteItemInput {
	key := bookmarkKey{
		UserId: user,
		Book:   book,
	}

	keyMap, _ := dynamodbattribute.MarshalMap(key)

	return &dynamodb.DeleteItemInput{
		Key:                    keyMap,
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(d.TableName),
	}
}

//GetBookmark returns a users bookmark item for a specific book
func (d *Dynamo) GetBookmark(ctx context.Context, user string, book string) (UserBookmark, error) {
	start := time.Now()
	var bookmark UserBookmark
	getItemInput := d.createGetItemInput(user, book)

	getResult, err := d.Client.GetItemWithContext(ctx, getItemInput)
	if err != nil {
		return bookmark, err
	}

	if len(getResult.Item) <= 0 {
		return bookmark, NotFoundException
	}

	err = dynamodbattribute.UnmarshalMap(getResult.Item, &bookmark)
	if err != nil {
		return bookmark, err
	}

	d.observer(start, *getResult.ConsumedCapacity.CapacityUnits, _get)
	return bookmark, nil
}

func (d *Dynamo) createGetItemInput(user string, book string) *dynamodb.GetItemInput {
	key := bookmarkKey{
		UserId: user,
		Book:   book,
	}

	keyMap, _ := dynamodbattribute.MarshalMap(key)

	return &dynamodb.GetItemInput{
		Key:                    keyMap,
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(d.TableName),
	}
}

//GetBookmarks returns a list of a users bookmarks from newest to oldest
func (d *Dynamo) GetBookmarks(ctx context.Context, user string, filter string, limit int) ([]UserBookmark, error) {
	start := time.Now()
	var bookmarks []UserBookmark

	itemQueryInput, err := d.createFilterQuery(user, filter)
	if err != nil {
		return bookmarks, err
	}

	queryResult, err := d.Client.QueryWithContext(ctx, itemQueryInput)
	if err != nil {
		return bookmarks, err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(queryResult.Items, &bookmarks)
	if err != nil {
		return bookmarks, err
	}

	limit = min(len(bookmarks), limit)

	d.observer(start, *queryResult.ConsumedCapacity.CapacityUnits, _query)
	return bookmarks[:limit], nil
}

func min(i int, j int) int {
	if i < j {
		return i
	}
	return j
}

func (d *Dynamo) createFilterQuery(user string, statusFilter string) (*dynamodb.QueryInput, error) {
	var expr expression.Expression
	var err error

	key := expression.Key("UserId").Equal(expression.Value(user))
	filter := expression.Name("Status").Equal(expression.Value(statusFilter))
	expr, err = expression.NewBuilder().WithKeyCondition(key).WithFilter(filter).Build()

	if statusFilter == _nonefilter {
		expr, err = expression.NewBuilder().WithKeyCondition(key).Build()
	}

	return &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		KeyConditionExpression:    expr.KeyCondition(),
		IndexName:                 aws.String("History"),
		Select:                    aws.String("ALL_PROJECTED_ATTRIBUTES"),
		ReturnConsumedCapacity:    aws.String("TOTAL"),
		Limit:                     aws.Int64(int64(maxPageSize)),
		ScanIndexForward:          aws.Bool(false),
		TableName:                 aws.String(d.TableName),
	}, err
}

func (d *Dynamo) observer(start time.Time, capacity float64, operation string) {
	if d.Reporter != nil {
		d.Reporter.ObserverHistogram(_dynamoLatencyHistogram, float64(time.Since(start).Milliseconds()), operation, d.TableName)
		d.Reporter.ObserverSummary(_dynamoLatencySummary, float64(time.Since(start).Milliseconds()), operation, d.TableName)
		d.Reporter.ObserverCount(_dynamoCapacityUsed, capacity, operation, d.TableName)
	}
}
