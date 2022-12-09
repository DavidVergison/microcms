package tests

import (
	"encoding/base64"
	"encoding/json"

	"github.com/DavidVergison/microcms/articles/entities"
	articleFormatter "github.com/DavidVergison/microcms/https/article-formatter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/brianvoe/gofakeit/v6"
)

func TestGetArticle() entities.Article {
	return entities.Article{
		Location: entities.LocationStruct{
			Site:        gofakeit.AchAccount(),
			Category:    gofakeit.MinecraftArmorTier(),
			SubCategory: gofakeit.MinecraftArmorPart(),
		},
		ReadingKey: gofakeit.Animal(),
		Resume:     gofakeit.Vegetable(),
		Content:    `{"name":"` + gofakeit.Breakfast() + `"}`,
		Meta: map[string]string{
			"A": gofakeit.Adjective(),
		},
		ID: gofakeit.BitcoinAddress(),
	}
}

func testGetRequestBodyForArticle(article entities.Article) string {
	var c complex
	json.Unmarshal([]byte(article.Content), &c)
	dto := writeArticleRequestTestDto{
		Category:    article.Location.Category,
		SubCategory: article.Location.SubCategory,
		ID:          article.ID,
		Site:        article.Location.Site,
		Meta:        article.Meta,
		Content: complex{
			Name: c.Name,
		},
		Resume: article.Resume,
	}
	json, _ := json.Marshal(dto)
	return string(json)
}

func testGetToken(username string) string {
	claim := articleFormatter.TokenClaims{
		Username: username,
	}
	jsonClaim, _ := json.Marshal(claim)
	return "xxx." + base64.StdEncoding.EncodeToString(jsonClaim) + ".xxx"
}

func TestGetRequestForArticle(article entities.Article, username string) events.APIGatewayProxyRequest {
	body := testGetRequestBodyForArticle(article)

	return events.APIGatewayProxyRequest{
		Body: body,
		Headers: map[string]string{
			"Authorization": "Bearer " + testGetToken(username),
		},
	}
}

func TestGetExpressionBuilderForArticle(article entities.Article) (expression.Expression, error) {
	expr := expression.
		Set(expression.Name("content"), expression.Value(article.Content)).
		Set(expression.Name("meta"), expression.Value(article.Meta)).
		Set(expression.Name("reading_key"), expression.Value(article.ReadingKey)).
		Set(expression.Name("resume"), expression.Value(article.Resume)).
		Set(expression.Name("category"), expression.Value(article.Location.Category)).
		Set(expression.Name("site"), expression.Value(article.Location.Site)).
		Set(expression.Name("subCategory"), expression.Value(article.Location.SubCategory))

	return expression.NewBuilder().
		WithUpdate(expr).
		Build()
}

type writeArticleRequestTestDto struct {
	Category    string            `json:"category"`
	SubCategory string            `json:"subcategory"`
	ID          string            `json:"id"`
	Site        string            `json:"website"`
	Meta        map[string]string `json:"meta"`
	Content     complex           `json:"content"`
	Resume      string            `json:"resume"`
}
type complex struct {
	Name string `json:"name"`
}
