package models

import "github.com/pkg/errors"

const paginationMaxLimit = 50

// Page indicates requested page of data.
type Page struct {
	Offset int64
	Limit  int64
}

// NewPage returns a new Page instance give the offset and limit.
func NewPage(offset, limit int64) Page {
	return Page{Offset: offset, Limit: limit}
}

// Validate the offset and limit.
func (p Page) Validate() error {
	if p.Offset < 0 {
		return errors.New("offset value can not be negative")
	}
	if p.Limit < 1 || paginationMaxLimit < p.Limit {
		return errors.Errorf("limit should be between 1 and %d", paginationMaxLimit)
	}
	return nil
}
