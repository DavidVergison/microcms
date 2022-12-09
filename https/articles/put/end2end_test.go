package main

import (
	"errors"
	"net/http"
	"testing"

	dynamock "github.com/DavidVergison/go-dynamock"
	"github.com/DavidVergison/microcms/articles/entities"
	"github.com/DavidVergison/microcms/tests"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var mock *dynamock.DynaMock

func TestGetName(t *testing.T) {

	// GIVEN an which that does not pass pass the sanity check
	// WHEN you try to update it
	// THEN you have an error "bad request"
	t.Run("an which that does not pass pass the sanity check", func(t *testing.T) {
		article := tests.TestGetArticle()
		site := tests.TestGetSite()
		articleWithoutID := article
		articleWithoutID.ID = ""
		articleWithoutSite := article
		articleWithoutSite.Location.Site = ""
		articleWithoutCat := article
		articleWithoutCat.Location.Category = ""
		articleWithoutSubCat := article
		articleWithoutSubCat.Location.SubCategory = ""
		articlesList := []entities.Article{
			articleWithoutID,
			articleWithoutSite,
			articleWithoutCat,
			articleWithoutSubCat,
		}

		for _, article := range articlesList {
			// set the mock
			request := tests.TestGetRequestForArticle(article, site.Username)
			cnx.Dynamo, mock = dynamock.New()
			initArticleMock(mock, &article, &site, true, true, nil)
			initSiteMock(mock, site,nil)

			// do the thing
			response, _ := handler(request)

			// test
			if response.StatusCode != http.StatusBadRequest {
				t.Fatalf("waiting for %v, got %v", http.StatusBadRequest, response.StatusCode)
			}
		}
	})

	// GIVEN an article the user did not own
	// WHEN you try to update it
	// THEN you have an error "forbidden"
	t.Run("The user must owned the article to update it", func(t *testing.T) {
		article := tests.TestGetArticle()
		site := tests.TestGetSite()
		//article.Location.Site = site.Site

		// set the mock
		request := tests.TestGetRequestForArticle(article, "xyz")
		cnx.Dynamo, mock = dynamock.New()
		initArticleMock(mock, &article, &site, false, false, nil)
		initSiteMock(mock, site,nil)

		// do the thing
		response, _ := handler(request)

		// test
		if response.StatusCode != http.StatusForbidden {
			t.Fatalf("waiting for %v, got %v", http.StatusForbidden, response.StatusCode)
		}
	})

	// GIVEN an article
	// WHEN everything is ok
	// THEN you have an "OK" status"
	t.Run("If everything is ok, the response is 200", func(t *testing.T) {
		article := tests.TestGetArticle()
		site := tests.TestGetSite()

		// set the mock
		cnx.Dynamo, mock = dynamock.New()
		initArticleMock(mock, &article, &site, true, true, nil)
		initSiteMock(mock, site,nil)
		request := tests.TestGetRequestForArticle(article, site.Username)

		// do the thing
		response, _ := handler(request)

		// test
		if response.StatusCode != http.StatusOK {
			t.Fatalf("waiting for %v, got %v", http.StatusOK, response.StatusCode)
		}
	})

	// GIVEN an article
	// WHEN there is a problem with the DB (Article)
	// THEN you have an error "bad gateway"
	t.Run("If there is a pb with the DB (Article), the error is 502", func(t *testing.T) {
		article := tests.TestGetArticle()
		site := tests.TestGetSite()
		//article.Location.Site = site.Site

		// set the mock
		cnx.Dynamo, mock = dynamock.New()
		initArticleMock(mock, &article, &site, true, false, errors.New("random"))
		initSiteMock(mock, site,nil)
		request := tests.TestGetRequestForArticle(article, site.Username)

		// do the thing
		response, _ := handler(request)

		// test
		if response.StatusCode != http.StatusBadGateway {
			t.Fatalf("waiting for %v, got %v", http.StatusBadGateway, response.StatusCode)
		}
	})

	// GIVEN an article
	// WHEN there is a problem with the DB (Site)
	// THEN you have an error "bad gateway"
	t.Run("If there is a pb with the DB (Site), the error is 502", func(t *testing.T) {
		article := tests.TestGetArticle()
		site := tests.TestGetSite()
		//article.Location.Site = site.Site

		// set the mock
		cnx.Dynamo, mock = dynamock.New()
		initArticleMock(mock, &article, &site, true, true, nil)
		initSiteMock(mock, site, errors.New("random"))
		request := tests.TestGetRequestForArticle(article, site.Username)

		// do the thing
		response, _ := handler(request)

		// test
		if response.StatusCode != http.StatusBadGateway {
			t.Fatalf("waiting for %v, got %v", http.StatusBadGateway, response.StatusCode)
		}
	})

	// GIVEN a new article
	// WHEN you try to update it
	// THEN you have an error "not found"
	t.Run("An Article must exist to be updated", func(t *testing.T) {
		article := tests.TestGetArticle()
		article.ID = "xyz"
		site := tests.TestGetSite()
		//article.Location.Site = site.Site

		// set the mock
		cnx.Dynamo, mock = dynamock.New()
		initArticleMock(mock, &article, &site, true, false, &dynamodb.ConditionalCheckFailedException{})
		initSiteMock(mock, site, nil)
		request := tests.TestGetRequestForArticle(article, site.Username)

		// do the thing
		response, _ := handler(request)

		// test
		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("waiting for %v, got %v", http.StatusNotFound, response.StatusCode)
		}
	})

}

type testData struct {
	article entities.Article
	site    entities.Site
}

func initSiteMock(mock *dynamock.DynaMock, site entities.Site, err error) {
	if err != nil {
		mock.ExpectQuery().Table("Sites").
			WillReturnError(err)
	} else {
		mock.ExpectQuery().Table("Sites").
			WillReturns(dynamodb.QueryOutput{
				Items: tests.TestGetAttributesValuesForSite(site),
			})
	}

}

func initArticleMock(mock *dynamock.DynaMock, article *entities.Article, site *entities.Site, own bool, updated bool, err error) {
	if own {
		article.Location.Site = site.Site
		article.ReadingKey = site.ReadingKey
	}

	output := dynamodb.PutItemOutput{}
	if updated {
		output.Attributes = map[string]*dynamodb.AttributeValue{
			"X": {
				S: aws.String("Y"),
			},
		}
	}

	article.User = site.Username
	av, _ := dynamodbattribute.MarshalMap(article)

	if err != nil {
		mock.ExpectPutItem().ToTable("Articles").
			WithItems(av).
			WillReturnError(err)
	} else {
		mock.ExpectPutItem().ToTable("Articles").
			WithItems(av).
			WillReturns(output)
	}

}
