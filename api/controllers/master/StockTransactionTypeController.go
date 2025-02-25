package mastercontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type StockTransactionTypeController interface {
	GetStockTransactionTypeByCode(writer http.ResponseWriter, request *http.Request)
	GetAllStockTransactionType(writer http.ResponseWriter, request *http.Request)
}
type StockTransactionTypeControllerImpl struct {
	service masterservice.StockTransactionTypeService
}

func NewStockTransactionTypeController(service masterservice.StockTransactionTypeService) StockTransactionTypeController {

	return &StockTransactionTypeControllerImpl{service: service}
}

// @Summary Get Stock Transaction Type By Code
// @Description REST API Stock Transaction Type
// @Accept json
// @Produce json
// @Tags Master : Stock Transaction Type
// @Security BearerAuth
// @Param stock_transaction_type_code path string true "stock_transaction_type_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/stock-transaction-type/{stock_transaction_type_code} [get]
func (controller *StockTransactionTypeControllerImpl) GetStockTransactionTypeByCode(writer http.ResponseWriter, request *http.Request) {
	StockTransactionTypeCode := chi.URLParam(request, "stock_transaction_type_code")
	res, err := controller.service.GetStockTransactionTypeByCode(StockTransactionTypeCode)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Get Stock Trasaction By Code", http.StatusOK)
}

// @Summary Get All Stock Transaction Type
// @Description REST API Stock Transaction Type
// @Accept json
// @Produce json
// @Tags Master : Stock Transaction Type
// @Security BearerAuth
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param sort_of query string false "sort_of"
// @Param sort_by query string false "sort_by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/stock-transaction-type [get]
func (controller *StockTransactionTypeControllerImpl) GetAllStockTransactionType(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{}
	paginations := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	filter := utils.BuildFilterCondition(queryParams)
	res, err := controller.service.GetAllStockTransactionType(filter, paginations)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, res.Rows, "Get Stock Transaction Type Controller", http.StatusOK, res.Limit, res.Page, res.TotalRows, res.TotalPages)
}
