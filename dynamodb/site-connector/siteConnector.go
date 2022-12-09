package siteconnector

import (
	"log"

	"github.com/DavidVergison/microcms/articles/entities"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	db "github.com/DavidVergison/microcms/dynamodb"
)

type SiteConnector struct {
	db        dynamodbiface.DynamoDBAPI
	tableName string
}

func NewSiteConnector(db db.DynamoDbConnector) SiteConnector {
	return SiteConnector{
		tableName: "Sites",
		db:        db.Dynamo,
	}
}

func (adapter SiteConnector) GetSitesByUser(username string) ([]entities.Site, error) {

	keyCond := expression.Key("username").
		Equal(expression.Value(username))

	expr, err := expression.NewBuilder().
		WithKeyCondition(keyCond).
		Build()
	if err != nil {
		log.Fatalf(err.Error())
	}

	input := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 aws.String(adapter.tableName),
	}

	response := []entities.Site{}

	result, err := adapter.db.Query(input)
	if err != nil {
		return response, err
	}

	for _, it := range result.Items {
		record := entities.SiteStruct{
			Username:   *it["username"].S,
			Site:       *it["site"].S,
			ReadingKey: *it["reading_key"].S,
		}

		response = append(response, entities.NewSite(record))
	}

	return response, nil
}
