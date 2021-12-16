package ctrl

import (
	"fmt"
	"strconv"

	"github.com/haorenfsa/milvus-ops/errs"
	models "github.com/haorenfsa/milvus-ops/model"
)

type Queryer interface {
	DefaultQuery(key string, defaultValue string) string
	Query(key string) string
}

// parsePaginatorFromQuery returns a pagination by querying given Queryer.
func parsePaginatorFromQuery(q Queryer) (models.Pagination, error) {
	var err error
	p := models.Pagination{}

	val := q.DefaultQuery("page", "1")
	if p.Page, err = strconv.Atoi(val); err != nil {
		return p, fmt.Errorf("%w invalid 'page' value '%s': %s", errs.ErrBadRequest, val, err)
	}

	val = q.DefaultQuery("per_page", "10")
	if p.PerPage, err = strconv.Atoi(val); err != nil {
		return p, fmt.Errorf("%w invalid 'per_page' value '%s': %s", errs.ErrBadRequest, val, err)
	}

	return p, nil
}
