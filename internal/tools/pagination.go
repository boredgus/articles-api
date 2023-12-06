package tools

import (
	"errors"
	"strconv"
)

var PageNotANumberErr = errors.New("page should be number")
var LimitNotANumberErr = errors.New("limit should be number")
var PageOutOfRangeErr = errors.New("page should more or equal than 0 (default 0)")
var LimitOutOfRangeErr = errors.New("limit should be from 1 to 50 (default 10)")

type Pagination struct{}

func NewPagination() Pagination {
	return Pagination{}
}

func (n Pagination) Parse(page, limit string) (p, l int, e error) {
	if len(page) == 0 {
		page = "0"
	}
	p, err := strconv.Atoi(page)
	if err != nil {
		return 0, 0, PageNotANumberErr
	}
	if p < 0 {
		return 0, 0, PageOutOfRangeErr
	}
	if len(limit) == 0 {
		limit = "10"
	}
	l, err = strconv.Atoi(limit)
	if err != nil {
		return 0, 0, LimitNotANumberErr
	}
	if l < 0 || l > 50 {
		return 0, 0, LimitOutOfRangeErr
	}
	return
}
