package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//NewDynamoClient initializes a DynamoDB client
func NewDynamoClient(session *session.Session) *dynamodb.DynamoDB {
	return dynamodb.New(session)
}

//newAWSSessions creates a aws session
func NewAWSSessions(id string, secret string, region string, endpoint string) (*session.Session, error) {
	awsCredentials := credentials.NewStaticCredentials(id, secret, "")

	awsConfig := aws.NewConfig().
		WithCredentials(awsCredentials).
		WithRegion(region).
		WithEndpoint(endpoint)

	return session.NewSession(awsConfig)
}