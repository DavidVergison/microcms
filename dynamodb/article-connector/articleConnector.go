package articleconnector

import (
	"github.com/DavidVergison/microcms/articles/entities"
	db "github.com/DavidVergison/microcms/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type ArticleConnector struct {
	db        dynamodbiface.DynamoDBAPI
	tableName string
}

func NewArticleConnector(db db.DynamoDbConnector) ArticleConnector {
	return ArticleConnector{
		tableName: "Articles",
		db:        db.Dynamo,
	}
}

func (adapter ArticleConnector) SaveArticle(article entities.Article) (bool, error) {

	av, err := dynamodbattribute.MarshalMap(article)
	input := &dynamodb.PutItemInput{
		TableName:           aws.String(adapter.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_exists(ID)"),
		ReturnValues:        aws.String("ALL_OLD"),
	}

	result, err := adapter.db.PutItem(input)
	if err != nil {
		_, ok := err.(*dynamodb.ConditionalCheckFailedException)
		if ok {
			return false, nil
		}
		return false, err
	}

	return len(result.Attributes) > 0, nil
}

func buildSaveArticleExpression(article entities.Article) (expression.Expression, error) {
	expr, err := expression.NewBuilder().
		WithCondition(
			expression.Name("ID").Equal(expression.Value(article.ID)),
		).
		Build()
	return expr, err
}
