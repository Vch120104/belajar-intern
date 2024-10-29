package mastercontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
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
func (controller *StockTransactionTypeControllerImpl) GetStockTransactionTypeByCode(writer http.ResponseWriter, request *http.Request) {
	StockTransactionTypeCode := chi.URLParam(request, "stock_transaction_type_code")
	res, err := controller.service.GetStockTransactionTypeByCode(StockTransactionTypeCode)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Get Stock Trasaction By Code", http.StatusOK)
}

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
