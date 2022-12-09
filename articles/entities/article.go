package entities

type LocationStruct struct {
	Site        string
	Category    string
	SubCategory string
}

type ArticleStruct struct {
	Location   LocationStruct
	User       string
	ReadingKey string
	ID         string
	Meta       map[string]string
	Resume     string
	Content    string
}

type Article ArticleStruct


