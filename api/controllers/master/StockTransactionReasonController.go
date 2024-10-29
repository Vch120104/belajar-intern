package mastercontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterservice "after-sales/api/services/master"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type StockTransactionReasonController interface {
	GetStockTransactionReasonByCode(writer http.ResponseWriter, request *http.Request)
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
