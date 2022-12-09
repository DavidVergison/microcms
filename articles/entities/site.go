package entities

type SiteStruct struct {
	Site       string
	Username   string
	ReadingKey string
}

type Site SiteStruct

func NewSite(newSite SiteStruct) Site {
	return Site(newSite)
}
