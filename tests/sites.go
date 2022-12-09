package tests

import (
	"github.com/DavidVergison/microcms/articles/entities"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/brianvoe/gofakeit/v6"
)

func TestGetSite() entities.Site {
	return entities.Site{
		Site:       gofakeit.Adverb(),
		Username:   gofakeit.Name(),
		ReadingKey: gofakeit.Animal(),
	}
}

func TestGetAttributesValuesForSite(site entities.Site) []map[string]*dynamodb.AttributeValue {
	return []map[string]*dynamodb.AttributeValue{
		{
			"username":    {S: aws.String(site.Username)},
			"site":        {S: aws.String(site.Site)},
			"reading_key": {S: aws.String(site.ReadingKey)},
		},
	}
}
