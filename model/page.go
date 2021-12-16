package model

// Pagination ..
type Pagination struct {
	Page    int
	PerPage int
}

// Offset return pagination's offset, which used by database
func (p Pagination) Offset() int {
	page := p.Page - 1
	if page < 0 {
		page = 0
	}
	return page * p.PerPage
}

// Limit just return the perPage
func (p Pagination) Limit() int {
	return p.PerPage
}
