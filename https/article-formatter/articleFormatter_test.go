package articleFormatter

import (
	"encoding/json"
	"testing"

	"github.com/DavidVergison/microcms/https/dto"
	"github.com/aws/aws-lambda-go/events"
)

func TestArticleFormatter(t *testing.T) {

	t.Run("happyFlow", func(t *testing.T) {

		// GIVEN
		request := getRequest()
		wantedUsername := "janedoe@example.com"

		// WHEN
		article, err := ArticleFormatter(request)

		// THEN
		if err != nil {
			t.Fatalf("Non expected error : %s", err.Error())
		}
		if article.User != wantedUsername {
			t.Fatalf("Non expected user (waiting for %s, get %s)", wantedUsername, article.User)
		}

	})

	t.Run("extractUserFromBearer", func(t *testing.T) {

		// GIVEN
		wantedUsername := "janedoe@example.com"
		token := `xxx.eyJzdWIiOiJhYWFhYWFhYS1iYmJiLWNjY2MtZGRkZC1lZWVlZWVlZWVlZWUiLCJkZXZpY2Vfa2V5IjoiYWFhYWFhYWEtYmJiYi1jY2NjLWRkZGQtZWVlZWVlZWVlZWVlIiwiY29nbml0bzpncm91cHMiOlsiYWRtaW4iXSwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU2MjE5MDUyNCwiaXNzIjoiaHR0cHM6Ly9jb2duaXRvLWlkcC51cy13ZXN0LTIuYW1hem9uYXdzLmNvbS91cy13ZXN0LTJfZXhhbXBsZSIsImV4cCI6MTU2MjE5NDEyNCwiaWF0IjoxNTYyMTkwNTI0LCJvcmlnaW5fanRpIjoiYWFhYWFhYWEtYmJiYi1jY2NjLWRkZGQtZWVlZWVlZWVlZWVlIiwianRpIjoiYWFhYWFhYWEtYmJiYi1jY2NjLWRkZGQtZWVlZWVlZWVlZWVlIiwiY2xpZW50X2lkIjoiNTdjYmlzaGs0ajI0cGFiYzEyMzQ1Njc4OTAiLCJ1c2VybmFtZSI6ImphbmVkb2VAZXhhbXBsZS5jb20ifQ==.xxx`

		// WHEN
		username, err := extractUserFromBearer(token)

		// THEN
		if err != nil {
			t.Fatalf("Non expected error : %s", err.Error())
		}
		if username != wantedUsername {
			t.Fatalf("Non expected result (waiting for %s, get %s)", wantedUsername, username)
		}

	})
}

func getRequest() events.APIGatewayProxyRequest {
	payload := dto.WriteArticleRequestDto{
		Category:    "cat1",
		SubCategory: "subcat1",
		ID:          "XYZ",
		Meta: map[string]string{
			"A": "a",
			"B": "b",
		},
		Content: dto.JsonPayload{
			Value: `"a":"b"`,
		},
	}

	jsonByte, _ := json.Marshal(payload)

	return events.APIGatewayProxyRequest{
		Body: string(jsonByte),
		Headers: map[string]string{
			"Authorization": "Bearer xxx.eyJzdWIiOiJhYWFhYWFhYS1iYmJiLWNjY2MtZGRkZC1lZWVlZWVlZWVlZWUiLCJkZXZpY2Vfa2V5IjoiYWFhYWFhYWEtYmJiYi1jY2NjLWRkZGQtZWVlZWVlZWVlZWVlIiwiY29nbml0bzpncm91cHMiOlsiYWRtaW4iXSwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU2MjE5MDUyNCwiaXNzIjoiaHR0cHM6Ly9jb2duaXRvLWlkcC51cy13ZXN0LTIuYW1hem9uYXdzLmNvbS91cy13ZXN0LTJfZXhhbXBsZSIsImV4cCI6MTU2MjE5NDEyNCwiaWF0IjoxNTYyMTkwNTI0LCJvcmlnaW5fanRpIjoiYWFhYWFhYWEtYmJiYi1jY2NjLWRkZGQtZWVlZWVlZWVlZWVlIiwianRpIjoiYWFhYWFhYWEtYmJiYi1jY2NjLWRkZGQtZWVlZWVlZWVlZWVlIiwiY2xpZW50X2lkIjoiNTdjYmlzaGs0ajI0cGFiYzEyMzQ1Njc4OTAiLCJ1c2VybmFtZSI6ImphbmVkb2VAZXhhbXBsZS5jb20ifQ==.xxx",
		},
	}
}
