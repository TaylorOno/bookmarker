package config

import (
	"fmt"
	"log"

	"github.com/TaylorOno/bookmarker/internal/repository"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//CreateAWSSessions creates a aws session
func CreateAWSSessions(id string, secret string, region string, endpoint string) (*session.Session, error) {
	awsCredentials := credentials.NewStaticCredentials(id, secret, "")

	awsConfig := aws.NewConfig().
		WithCredentials(awsCredentials).
		WithRegion(region).
		WithEndpoint(endpoint)

	return session.NewSession(awsConfig)
}

//CreateTableIfNotExist if user has permission creates the DynamoDB bookmark table if it does not exist.
func CreateTableIfNotExist(dynamo *repository.Dynamo) {
	name := "bookmarks"
	bookmarkTableDescription := createTableRequest(name)

	_, err := dynamo.Client.CreateTable(&bookmarkTableDescription)
	if err != nil {
		fmt.Println(err.Error())
	}

	tableDescription := dynamodb.DescribeTableInput{TableName: &name}
	err = dynamo.Client.WaitUntilTableExistsWithContext(aws.BackgroundContext(), &tableDescription, request.WithWaiterMaxAttempts(5))
	if err != nil {
		log.Fatalf("Table: %v does not exist. Error %v", name, err.Error())
	}
}

func createTableRequest(tableName string) dynamodb.CreateTableInput {
	attributeDefinitions := tableAttributes()
	primaryKeySchema := primaryKey()
	localSecondaryIndex := secondaryIndex()
	provisionedThroughput := provisionThroughPut(5, 5)

	return dynamodb.CreateTableInput{
		AttributeDefinitions:  attributeDefinitions,
		KeySchema:             primaryKeySchema,
		LocalSecondaryIndexes: localSecondaryIndex,
		ProvisionedThroughput: provisionedThroughput,
		TableName:             aws.String(tableName),
	}
}

func provisionThroughPut(read int64, write int64) *dynamodb.ProvisionedThroughput {
	return &dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(read),
		WriteCapacityUnits: aws.Int64(write),
	}
}

func secondaryIndex() []*dynamodb.LocalSecondaryIndex {
	return []*dynamodb.LocalSecondaryIndex{
		{
			IndexName: aws.String("History"),
			KeySchema: secondaryKey(),
			Projection: &dynamodb.Projection{
				ProjectionType: aws.String("INCLUDE"),
				NonKeyAttributes: []*string{
					aws.String("Page"),
					aws.String("Status"),
				},
			},
		},
	}
}

func secondaryKey() []*dynamodb.KeySchemaElement {
	return []*dynamodb.KeySchemaElement{
		{
			AttributeName: aws.String("UserId"),
			KeyType:       aws.String("HASH"),
		},
		{
			AttributeName: aws.String("LastUpdated"),
			KeyType:       aws.String("RANGE"),
		},
	}
}

func primaryKey() []*dynamodb.KeySchemaElement {
	return []*dynamodb.KeySchemaElement{
		{
			AttributeName: aws.String("UserId"),
			KeyType:       aws.String("HASH"),
		},
		{
			AttributeName: aws.String("Book"),
			KeyType:       aws.String("RANGE"),
		},
	}
}

func tableAttributes() []*dynamodb.AttributeDefinition {
	return []*dynamodb.AttributeDefinition{
		{
			AttributeName: aws.String("UserId"),
			AttributeType: aws.String("S"),
		},
		{
			AttributeName: aws.String("Book"),
			AttributeType: aws.String("S"),
		},
		{
			AttributeName: aws.String("LastUpdated"),
			AttributeType: aws.String("S"),
		},
	}
}
