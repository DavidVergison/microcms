package articleFormatter

import (
	"encoding/json"
	"strings"

	"encoding/base64"

	"github.com/DavidVergison/microcms/articles/entities"
	"github.com/DavidVergison/microcms/https/dto"
	"github.com/aws/aws-lambda-go/events"
)

var (
	errUnableToUnmarshal = CustomErrUnableToUnmarshal{}
	errUnableToUseToken  = CustomErrUnableToUseToken{}
)

func ArticleFormatter(request events.APIGatewayProxyRequest) (entities.Article, error) {
	var dto dto.WriteArticleRequestDto
	err := json.Unmarshal([]byte(request.Body), &dto)
	if err != nil {
		return entities.Article{}, errUnableToUnmarshal
	}

	bearer := strings.Split(request.Headers["Authorization"], "Bearer ")[1]
	user, err := extractUserFromBearer(bearer)
	if err != nil {
		return entities.Article{}, err
	}

	return entities.Article{
		Location: entities.LocationStruct{
			Site:        dto.Site,
			Category:    dto.Category,
			SubCategory: dto.SubCategory,
		},
		ID:      dto.ID,
		Meta:    dto.Meta,
		Content: dto.Content.Value,
		User:    user,
		Resume:  dto.Resume,
	}, nil
}

func extractUserFromBearer(bearer string) (string, error) {
	splittedBearer := strings.Split(bearer, ".")
	if len(splittedBearer) != 3 {
		return "", errUnableToUseToken
	}

	jsonPayload, err := base64.StdEncoding.DecodeString(splittedBearer[1])
	if err != nil {
		return "", errUnableToUseToken
	}

	var claim TokenClaims

	err = json.Unmarshal(jsonPayload, &claim)
	if err != nil {
		return "", errUnableToUseToken
	}

	return claim.Username, nil
}

type TokenClaims struct {
	Username string `json:"username"`
}

// CustomErrUnableToUnmarshal Unable To Unmarshal the request body
type CustomErrUnableToUnmarshal struct{}

func (m CustomErrUnableToUnmarshal) Error() string {
	return "Unable To Unmarshal the request body"
}

// CustomErrUnableToUseToken Unable To use the bearer token
type CustomErrUnableToUseToken struct{}

func (m CustomErrUnableToUseToken) Error() string {
	return "Unable To use the bearer token"
}
