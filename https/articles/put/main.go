package main

import (
	"net/http"
	"reflect"

	db "github.com/DavidVergison/microcms/dynamodb"
	dbArticle "github.com/DavidVergison/microcms/dynamodb/article-connector"
	dbSite "github.com/DavidVergison/microcms/dynamodb/site-connector"
	tools "github.com/DavidVergison/microcms/https"
	articleFormatter "github.com/DavidVergison/microcms/https/article-formatter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	customErr "github.com/DavidVergison/microcms/articles/errors"
	savearticle "github.com/DavidVergison/microcms/articles/save-article"
)

var cnx db.DynamoDbConnector

func init() {
	cnx = db.NewDynamoDbConnector("eu-west-3")
}

func putArticle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	articleDB := dbArticle.NewArticleConnector(cnx)
	userDB := dbSite.NewSiteConnector(cnx)

	article, err := articleFormatter.ArticleFormatter(request)
	if err != nil {
		return formatterErrorHandler(err)
	}

	// usecase
	err = savearticle.SaveArticle(articleDB, userDB, article)

	if err != nil {
		return usecaseErrorHandler(err)
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
}

func formatterErrorHandler(err error) (events.APIGatewayProxyResponse, error) {
	errMap := map[reflect.Type]int{
		reflect.TypeOf(articleFormatter.CustomErrUnableToUnmarshal{}): http.StatusBadRequest,
		reflect.TypeOf(customErr.CustomErrorBadRequest{}):             http.StatusBadRequest,
	}
	return events.APIGatewayProxyResponse{StatusCode: errMap[reflect.TypeOf(err)]}, err
}
func usecaseErrorHandler(err error) (events.APIGatewayProxyResponse, error) {
	errMap := map[reflect.Type]int{
		reflect.TypeOf(customErr.CustomErrorBadRequest{}):      http.StatusBadRequest,
		reflect.TypeOf(customErr.CustomErrorArticleNotFound{}): http.StatusNotFound,
		reflect.TypeOf(customErr.CustomErrorForbidden{}):       http.StatusForbidden,
		reflect.TypeOf(customErr.CustomErrorWithTheDatabase{}): http.StatusBadGateway,
	}
	return events.APIGatewayProxyResponse{StatusCode: errMap[reflect.TypeOf(err)]}, err
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tools.LogAsJson(request)

	response, err := tools.HandleCors(request, putArticle)
	if err != nil {
		tools.LogAsJson(err)
	}

	tools.LogAsJson(response)
	return response, nil
}

func main() {
	lambda.Start(handler)
}
