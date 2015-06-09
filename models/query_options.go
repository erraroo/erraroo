package models

const DefaultPerPage = 10

type QueryOptions struct {
	PerPage int
	Page    int
}

func (o QueryOptions) PageOrDefault() int {
	if o.Page <= 0 {
		return 1
	}

	return o.Page
}

func (o QueryOptions) Offset() int {
	return (o.PageOrDefault() - 1) * o.PerPageOrDefault()
}

func (o QueryOptions) PerPageOrDefault() int {
	if o.PerPage <= 0 {
		return DefaultPerPage
	}

	return o.PerPage
}
