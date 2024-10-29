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
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
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
func (controller *StockTransactionReasonControllerImpl) GetStockTransactionReasonByCode(writer http.ResponseWriter, request *http.Request) {
	StockTransactionReasonCode := chi.URLParam(request, "stock_transaction_reason_code")
	res, err := controller.service.GetStockTransactionReasonByCode(StockTransactionReasonCode)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Get Stock Transaction Reason By Code", http.StatusOK)
}

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

// stock-transaction-reasonby-id/{stock_transaction_reason_id}
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

// stock-transaction-reason get
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
