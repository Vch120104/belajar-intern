package payloads

import (
	"after-sales/api/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type ResponsePagination struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Page       int         `json:"page"`
	Limit      int         `json:"page_limit"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

type ResponsePaginationHeader struct {
	Header interface{} `json:"data"`
	Data   ResponsePagination
}

// Deprecated: please change to the latest one without *gin.Context
func HandleSuccess(c *gin.Context, data interface{}, message string, status int) {
	res := Response{
		StatusCode: status,
		Message:    message,
		Data:       data,
	}

	c.JSON(status, res)
}

// Deprecated: please change to the latest one without *gin.Context
func HandleSuccessPagination(c *gin.Context, data interface{}, message string, status int, limit int, page int, totalRows int64, totalPages int) {
	res := ResponsePagination{
		StatusCode: status,
		Message:    message,
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Data:       data,
	}
	c.JSON(status, res)
}

func NewHandleSuccess(writer http.ResponseWriter, data interface{}, message string, status int) {
	res := Response{
		StatusCode: status,
		Message:    message,
		Data:       data,
	}

	helper.WriteToResponseBody(writer, res)
}

func NewHandleSuccessPagination(writer http.ResponseWriter, data interface{}, message string, status int, limit int, page int, totalRows int64, totalPages int) {
	res := ResponsePagination{
		StatusCode: status,
		Message:    message,
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Data:       data,
	}

	helper.WriteToResponseBody(writer, res)
}
