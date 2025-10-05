package helpers

import (
	"math"
	"strconv"
	"testcase/internal/utils"

	"github.com/gin-gonic/gin"
)

type PaginationParams struct {
	Page   int    `json:"page" form:"page"`
	Limit  int    `json:"limit" form:"limit"`
	Sort   string `json:"sort" form:"sort"`
	Order  string `json:"order" form:"order"`
	Search string `json:"search" form:"search"`
	Filter string `json:"filter" form:"filter"`
}

func ParsePaginationParams(c *gin.Context) *PaginationParams {
	params := &PaginationParams{
		Page:   1,
		Limit:  10,
		Sort:   "created_at",
		Order:  "desc",
		Search: "",
		Filter: "",
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			if limit > 100 {
				limit = 100
			}
			params.Limit = limit
		}
	}

	if sort := c.Query("sort"); sort != "" {
		params.Sort = sort
	}

	if order := c.Query("order"); order != "" {
		if order == "asc" || order == "desc" {
			params.Order = order
		}
	}

	params.Search = c.Query("search")

	params.Filter = c.Query("filter")

	return params
}

func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

func (p *PaginationParams) GetOrderBy() string {
	return p.Sort + " " + p.Order
}

func CreatePaginationResult(data interface{}, total int64, params *PaginationParams) *utils.PaginationResult {
	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	result := &utils.PaginationResult{
		List: data,
		Metadata: utils.PaginationMeta{
			Total:       total,
			Page:        params.Page,
			Limit:       params.Limit,
			TotalPages:  totalPages,
			HasNextPage: params.Page < totalPages,
			HasPrevPage: params.Page > 1,
		},
	}

	if result.Metadata.HasNextPage {
		nextPage := params.Page + 1
		result.Metadata.NextPage = &nextPage
	}

	if result.Metadata.HasPrevPage {
		prevPage := params.Page - 1
		result.Metadata.PreviousPage = &prevPage
	}

	return result
}

func CreatePaginationMeta(total int64, params *PaginationParams) *utils.PaginationMeta {
	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	meta := &utils.PaginationMeta{
		Total:       total,
		Page:        params.Page,
		Limit:       params.Limit,
		TotalPages:  totalPages,
		HasNextPage: params.Page < totalPages,
		HasPrevPage: params.Page > 1,
	}

	if meta.HasNextPage {
		nextPage := params.Page + 1
		meta.NextPage = &nextPage
	}

	if meta.HasPrevPage {
		prevPage := params.Page - 1
		meta.PreviousPage = &prevPage
	}

	return meta
}

func ValidatePaginationParams(params *PaginationParams) *PaginationParams {
	if params.Page < 1 {
		params.Page = 1
	}

	if params.Limit < 1 {
		params.Limit = 10
	}

	if params.Limit > 100 {
		params.Limit = 100
	}

	if params.Order != "asc" && params.Order != "desc" {
		params.Order = "desc"
	}

	return params
}

func (p *PaginationParams) PaginateQuery(query interface{}) interface{} {
	return query
}
