package mastercontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type StockTransactionReasonController interface {
	GetStockTransactionReasonByCode(writer http.ResponseWriter, request *http.Request)
	InsertStockTransactionReason(writer http.ResponseWriter, request *http.Request)
	GetStockTransactionReasonById(writer http.ResponseWriter, request *http.Request)
	GetAllStockTransactionReason(writer http.ResponseWriter, request *http.Request)
}

type StockTransactionReasonControllerImpl struct {
	service masterservice.StockTransactionReasonService
}

func StartStockTransactionReasonController(service masterservice.StockTransactionReasonService) StockTransactionReasonController {
	return &StockTransactionReasonControllerImpl{service: service}
}

// @Summary Get Stock Transaction Reason By Code
// @Description REST API Stock Transaction Reason
// @Accept json
// @Produce json
// @Tags Master : Stock Transaction Reason
// @Security AuthorizationKeyAuth
// @Param stock_transaction_reason_code path string true "stock_transaction_reason_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/stock-transaction-reason/{stock_transaction_reason_code} [get]
func (controller *StockTransactionReasonControllerImpl) GetStockTransactionReasonByCode(writer http.ResponseWriter, request *http.Request) {
	StockTransactionReasonCode := chi.URLParam(request, "stock_transaction_reason_code")
	res, err := controller.service.GetStockTransactionReasonByCode(StockTransactionReasonCode)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Get Stock Transaction Reason By Code", http.StatusOK)
}

// @Summary Insert Stock Transaction Reason
// @Description REST API Stock Transaction Reason
// @Accept json
// @Produce json
// @Tags Master : Stock Transaction Reason
// @Security AuthorizationKeyAuth
// @Param request body masterpayloads.StockTransactionReasonInsertPayloads true "Stock Transaction Reason Insert Payloads"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/stock-transaction-reason [post]
func (controller *StockTransactionReasonControllerImpl) InsertStockTransactionReason(writer http.ResponseWriter, request *http.Request) {
	formRequest := masterpayloads.StockTransactionReasonInsertPayloads{}
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)

	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	res, err := controller.service.InsertStockTransactionReason(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Create Transaction Reason", http.StatusCreated)

}

// @Summary Get Stock Transaction Reason By Id
// @Description REST API Stock Transaction Reason
// @Accept json
// @Produce json
// @Tags Master : Stock Transaction Reason
// @Security AuthorizationKeyAuth
// @Param stock_transaction_reason_id path string true "stock_transaction_reason_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/stock-transaction-reason/{stock_transaction_reason_id} [get]
func (controller *StockTransactionReasonControllerImpl) GetStockTransactionReasonById(writer http.ResponseWriter, request *http.Request) {
	//temp :=
	StockTransactionReasonId, _ := strconv.Atoi(chi.URLParam(request, "stock_transaction_reason_id"))
	res, err := controller.service.GetStockTransactionReasonById(StockTransactionReasonId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Get Stock Transaction Reason By Id", http.StatusOK)
}

// @Summary Get All Stock Transaction Reason
// @Description REST API Stock Transaction Reason
// @Accept json
// @Produce json
// @Tags Master : Stock Transaction Reason
// @Security AuthorizationKeyAuth
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param sort_of query string false "sort_of"
// @Param sort_by query string false "sort_by"
// @Success 200 {object} payloads.ResponsePagination
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/stock-transaction-reason [get]
func (controller *StockTransactionReasonControllerImpl) GetAllStockTransactionReason(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{}
	paginations := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	filter := utils.BuildFilterCondition(queryParams)
	res, err := controller.service.GetAllStockTransactionReason(filter, paginations)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, res.Rows, "Success Get All Stock Transaction Reason", http.StatusOK, res.Limit, res.Page, res.TotalRows, res.TotalPages)

}
