package main

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=100"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"omitempty,oneof=asc desc"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")

	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return PaginatedFeedQuery{}, err
		}

		fq.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		f, err := strconv.Atoi(offset)
		if err != nil {
			return PaginatedFeedQuery{}, err
		}

		fq.Offset = f
	}

	fq.Sort = qs.Get("sort")

	return fq, nil
}
