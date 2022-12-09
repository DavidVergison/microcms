package errors

// CustomErrorArticleNotFound article not found
type CustomErrorArticleNotFound struct{}

func (m CustomErrorArticleNotFound) Error() string {
	return "article not found"
}
