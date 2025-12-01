// internal/utils/pagination.go
package utils

import (
	"net/http"
	"strconv"
	"student-portal/internal/commons/constants"
)

// PaginationQuery holds the parsed page and limit values from the request.
type PaginationQuery struct {
	Page   int
	Limit  int
	Offset int
}

// PaginationResponse holds the metadata and data for a paginated response.
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalCount int64       `json:"total_count"`
	TotalPages int         `json:"total_pages"`
}

// NewPaginationQuery parses the 'page' and 'limit' query parameters from the request.
func NewPaginationQuery(r *http.Request) PaginationQuery {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = constants.DefaultPage
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = constants.DefaultLimit
	}

	if limit > constants.MaxLimit {
		limit = constants.MaxLimit
	}

	offset := (page - 1) * limit

	return PaginationQuery{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

// NewPaginationResponse creates a PaginationResponse object.
func NewPaginationResponse(data interface{}, query PaginationQuery, totalCount int64) PaginationResponse {
	totalPages := 0
	if query.Limit > 0 {
		totalPages = int((totalCount + int64(query.Limit) - 1) / int64(query.Limit))
	}

	return PaginationResponse{
		Data:       data,
		Page:       query.Page,
		Limit:      query.Limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}
}
