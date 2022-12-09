package errors

// CustomErrorForbidden This user can not write in this site
type CustomErrorForbidden struct{}

func (m CustomErrorForbidden) Error() string {
	return "This user can not write in this site"
}
