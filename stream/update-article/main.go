package main

import (
	"context"

	tools "github.com/DavidVergison/microcms/https"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.DynamoDBEvent) {
	tools.LogAsJson(event)
}

func main() {
	lambda.Start(handler)
}
