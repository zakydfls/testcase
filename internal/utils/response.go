package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	RequestID string      `json:"request_id,omitempty"`
	Code      int         `json:"code"`
	Success   bool        `json:"success"`
	Key       string      `json:"key,omitempty"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
}
type PaginationResult struct {
	Data         interface{} `json:"data"`
	Total        int64       `json:"total"`
	Page         int         `json:"page"`
	Limit        int         `json:"limit"`
	TotalPages   int         `json:"total_pages"`
	HasNextPage  bool        `json:"has_next_page"`
	HasPrevPage  bool        `json:"has_previous_page"`
	NextPage     *int        `json:"next_page"`
	PreviousPage *int        `json:"previous_page"`
}

func PaginatedResponse(ctx *gin.Context, result interface{}, totalItems int64, page, limit int, message string) {
	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))

	var nextPage, prevPage *int
	if page < totalPages {
		n := page + 1
		nextPage = &n
	}
	if page > 1 {
		p := page - 1
		prevPage = &p
	}

	hasNextPage := page < totalPages
	hasPrevPage := page > 1

	requestID := ctx.GetString("request_id")
	response := gin.H{
		"code":    http.StatusOK,
		"message": message,
		"result":  result,
		"metadata": PaginationResult{
			Data:         result,
			Total:        totalItems,
			Page:         page,
			Limit:        limit,
			TotalPages:   totalPages,
			HasNextPage:  hasNextPage,
			HasPrevPage:  hasPrevPage,
			NextPage:     nextPage,
			PreviousPage: prevPage,
		},
	}

	if requestID != "" {
		response["request_id"] = requestID
	}

	ctx.JSON(http.StatusOK, response)
}

func SuccessResponse(ctx *gin.Context, data interface{}, message string, statusCode int) {
	requestID := ctx.GetString("request_id")
	response := Response{
		Code:    statusCode,
		Message: message,
		Data:    data,
		Key:     "success",
	}

	if requestID != "" {
		response.RequestID = requestID
	}

	ctx.JSON(statusCode, response)
}
func ErrorResponse(c *gin.Context, errCode ErrorCode, err interface{}) {
	var message string
	var errorKey string

	switch e := err.(type) {
	case *AppError:
		message = e.GetDisplayMessage()
		errorKey = e.ErrorCode.Key
	case error:
		message = e.Error()
		errorKey = errCode.Key
	case string:
		message = e
		errorKey = errCode.Key
	default:
		message = errCode.Message
		errorKey = errCode.Key
	}

	response := Response{
		Success: false,
		Message: message,
		Error:   errorKey,
		Code:    errCode.Code,
		Data:    nil,
	}

	c.JSON(errCode.HttpStatus, response)
}

func HandleRouteNotFound(ctx *gin.Context) {
	requestID := ctx.GetString("request_id")
	errorInfo := GetError(ErrRouteNotFound)

	response := Response{
		Success: true,
		Code:    errorInfo.Code,
		Key:     errorInfo.Key,
		Message: errorInfo.Message,
	}

	if requestID != "" {
		response.RequestID = requestID
	}

	ctx.JSON(http.StatusNotFound, response)
}
