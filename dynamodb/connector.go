package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type DynamoDbConnector struct {
	Dynamo dynamodbiface.DynamoDBAPI
}

func NewDynamoDbConnector(region string) DynamoDbConnector {
	return DynamoDbConnector{
		Dynamo: getDBClient(region),
	}
}

func getDBClient(region string) dynamodbiface.DynamoDBAPI {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String(region),
		},
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess, aws.NewConfig().WithRegion(region))

	return dynamodbiface.DynamoDBAPI(svc)
}

type MyDynamo struct {
	Db dynamodbiface.DynamoDBAPI
}
