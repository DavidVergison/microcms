package errors

// CustomErrorBadRequest the article is invalid
type CustomErrorBadRequest struct{}

func (m CustomErrorBadRequest) Error() string {
	return "This article is invalid"
}
