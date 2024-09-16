package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=100"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"omitempty,oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
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

	tags := qs.Get("tags")
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	since := qs.Get("since")
	if since != "" {
		fq.Since = parseTime(since)
	}

	until := qs.Get("until")
	if since != "" {
		fq.Until = parseTime(until)
	}

	return fq, nil
}

func parseTime(s string) string {
	t, err := time.Parse(time.DateTime, s)

	if err != nil {
		return ""
	}

	return t.Format(time.DateTime)
}
