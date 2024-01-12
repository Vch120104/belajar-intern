package payloads

import "github.com/gin-gonic/gin"

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

func HandleSuccess(c *gin.Context, data interface{}, message string, status int) {
	res := Response{
		StatusCode: status,
		Message:    message,
		Data:       data,
	}

	c.JSON(status, res)
}

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
