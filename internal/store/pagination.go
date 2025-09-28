package store

import (
	"net/http"
	"strconv"
	"strings"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lt=101"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Search string   `json:"search" validate:"max=100"`
	Tags   []string `json:"tags" validate:"max=100"`
}

func (fd PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	query := r.URL.Query()
	if l := query.Get("limit"); l != "" {
		limit, err := strconv.Atoi(l)
		if err != nil {
			return fd, err
		}
		fd.Limit = limit
	}
	if of := query.Get("offset"); of != "" {
		offset, err := strconv.Atoi(of)
		if err != nil {
			return fd, err
		}
		fd.Offset = offset
	}
	if sort := query.Get("sort"); sort != "" {
		fd.Sort = sort
	}
	if search := query.Get("search"); search != "" {
		fd.Search = search
	}
	if tags := query.Get("tags"); tags != "" {
		fd.Tags = strings.Split(tags, ",")
	}

	return fd, nil
}
