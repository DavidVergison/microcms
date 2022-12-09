package adapters

import "github.com/DavidVergison/microcms/articles/entities"

type SiteDbInterface interface {
	GetSitesByUser(username string) ([]entities.Site, error)
}
