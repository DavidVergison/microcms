package savearticle

import (
	"github.com/DavidVergison/microcms/articles/entities"
	customError "github.com/DavidVergison/microcms/articles/errors"
	"github.com/DavidVergison/microcms/articles/save-article/adapters"
)

var (
	errForbidden = customError.CustomErrorForbidden{}
	errDatabase  = customError.CustomErrorWithTheDatabase{}
	errInvalid   = customError.CustomErrorBadRequest{}
	errNotFound  = customError.CustomErrorArticleNotFound{}
)

func isTheArticleValid(article entities.Article) bool {
	if article.ID == "" ||
		article.Location.Category == "" ||
		article.Location.SubCategory == "" ||
		article.Location.Site == "" {
		return false
	}
	return true
}

func SaveArticle(articlesDbAdapter adapters.ArticleDbInterface, userDbAdapter adapters.SiteDbInterface, article entities.Article) error {

	if !isTheArticleValid(article) {
		return errInvalid
	}

	site, err := getSiteIfAllowed(userDbAdapter, article)
	if err != nil {
		return err
	}

	article.ReadingKey = site.ReadingKey

	updated, err := articlesDbAdapter.SaveArticle(article)
	if err != nil {
		errDatabase.Detail = err.Error()
		return errDatabase
	}
	if !updated {
		return errNotFound
	}

	return nil
}

func getSiteIfAllowed(userDbAdapter adapters.SiteDbInterface, article entities.Article) (entities.Site, error) {
	authorizedSites, err := userDbAdapter.GetSitesByUser(article.User)
	if err != nil {
		errDatabase.Detail = err.Error()
		return entities.Site{}, errDatabase
	}

	site := selectSite(article.Location.Site, authorizedSites)
	if site.Site == "" {
		return entities.Site{}, errForbidden
	}

	return site, nil
}

func selectSite(a string, list []entities.Site) entities.Site {
	for _, b := range list {
		if b.Site == a {
			return b
		}
	}
	return entities.Site{}
}
