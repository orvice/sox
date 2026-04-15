package storage

import "github.com/aws/aws-sdk-go/service/dynamodb"

type DynamoDBStorage struct {
	client *dynamodb.DynamoDB
}
