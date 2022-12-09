package adapters

import (
	"github.com/DavidVergison/microcms/articles/entities"
)

type ArticleDbInterface interface {
	SaveArticle(article entities.Article) (bool, error)
}
