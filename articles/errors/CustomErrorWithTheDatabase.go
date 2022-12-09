package errors

// CustomErrForbidden This user can not write in this site
type CustomErrorWithTheDatabase struct {
	Detail string
}

func (m CustomErrorWithTheDatabase) Error() string {
	return "There is an error with the database : " + m.Detail
}
