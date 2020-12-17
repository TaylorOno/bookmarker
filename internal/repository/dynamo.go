package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	//Timeout for request to dynamo defaults to 1 second
	Timeout = 1 * time.Second

	//NotFoundException special exception for not found
	NotFoundException = errors.New("not found")
)

//Dynamo holds dynamo client and table name
type Dynamo struct {
	Client    dynamodbiface.DynamoDBAPI
	TableName string
}

//CreateDynamoRepository returns the DynamoDB repository implementation
func CreateDynamoRepository(awsSession *session.Session, tableName string) *Dynamo {
	return &Dynamo{
		Client:    dynamodb.New(awsSession),
		TableName: tableName,
	}
}

//CreateBookmark adds a bookmark for a user
func (d *Dynamo) CreateBookmark(ctx context.Context, bookmark UserBookmark) (UserBookmark, error) {
	av, err := dynamodbattribute.MarshalMap(bookmark)
	if err != nil {
		return bookmark, err
	}

	ctx, cancelFn := context.WithTimeout(ctx, Timeout)
	defer cancelFn()

	createItemInput := &dynamodb.PutItemInput{
		Item:                   av,
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(d.TableName),
	}

	result, err := d.Client.PutItemWithContext(ctx, createItemInput)
	if err != nil {
		return bookmark, err
	}

	log.Print(result.ConsumedCapacity.String())
	return bookmark, nil
}

//UpdateBookmark updates a users bookmark with new values this is done by deleting the old bookmark and inserting the new one
func (d *Dynamo) UpdateBookmark(ctx context.Context, bookmark UserBookmark) (UserBookmark, error) {
	return bookmark, errors.New("not implemented")
}

//DeleteBookmark deletes a users bookmark
func (d *Dynamo) DeleteBookmark(ctx context.Context, bookmark UserBookmark) (UserBookmark, error) {
	return bookmark, errors.New("not implemented")
}

//ReadBookmark returns a users bookmark item for a specific book
func (d *Dynamo) GetBookmark(ctx context.Context, user string, book string) (UserBookmark, error) {
	var bookmark UserBookmark
	getItemInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"UserId": {
				S: aws.String(user),
			},
			"Book": {
				S: aws.String(book),
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(d.TableName),
	}

	ctx, cancelFn := context.WithTimeout(ctx, Timeout)
	defer cancelFn()

	result, err := d.Client.GetItemWithContext(ctx, getItemInput)
	log.Print(result.ConsumedCapacity.String())
	if err != nil {
		return bookmark, err
	}

	if len(result.Item) <= 0 {
		log.Print()
		return bookmark, NotFoundException
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &bookmark)
	if err != nil {
		return bookmark, err
	}

	return bookmark, nil
}

//GetRecentBookmarks returns a list of a users bookmarks from newest to oldest
func (d *Dynamo) GetRecentBookmarks(ctx context.Context, user string, limit int64) ([]UserBookmark, error) {
	var bookmarks []UserBookmark

	key := expression.Key("UserId").Equal(expression.Value(user))
	filter := expression.Name("Status").NotEqual(expression.Value("FINISHED"))
	expr, err := expression.NewBuilder().WithKeyCondition(key).WithFilter(filter).Build()
	if err != nil {
		fmt.Println(err)
	}

	itemQueryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		IndexName:                 aws.String("History"),
		KeyConditionExpression:    expr.KeyCondition(),
		Limit:                     aws.Int64(limit),
		Select:                    aws.String("ALL_PROJECTED_ATTRIBUTES"),
		ReturnConsumedCapacity:    aws.String("TOTAL"),
		ScanIndexForward:          aws.Bool(false),
		TableName:                 aws.String(d.TableName),
	}

	//itemQueryInput := &dynamodb.QueryInput{
	//	IndexName:              aws.String("History"),
	//	KeyConditionExpression: aws.String("UserId=:UserId"),
	//	ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
	//		":UserId": {
	//			S: aws.String(user),
	//		},
	//	},
	//	Select:                 aws.String("ALL_PROJECTED_ATTRIBUTES"),
	//	Limit:                  aws.Int64(limit),
	//	ReturnConsumedCapacity: aws.String("TOTAL"),
	//	ScanIndexForward:       aws.Bool(false),
	//	TableName:              aws.String(d.TableName),
	//}

	ctx, cancelFn := context.WithTimeout(ctx, Timeout)
	defer cancelFn()

	result, err := d.Client.QueryWithContext(ctx, itemQueryInput)
	if err != nil {
		return bookmarks, err
	}
	log.Print(fmt.Sprintf("GetBookmarks WCU: %v", *result.ConsumedCapacity.CapacityUnits))

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &bookmarks)
	if err != nil {
		return bookmarks, err
	}

	return bookmarks, nil
}
