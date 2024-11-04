package transactionsparepartcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"net/http"
)

type StockTransactionController interface {
	StockTransactionInsert(writer http.ResponseWriter, request *http.Request)
}
type StockTransactionControllerImpl struct {
	service transactionsparepartservice.StockTransactionService
}

func (s *StockTransactionControllerImpl) StockTransactionInsert(writer http.ResponseWriter, request *http.Request) {
	var stocktransaction transactionsparepartpayloads.StockTransactionInsertPayloads

	helper.ReadFromRequestBody(request, &stocktransaction)
	res, err := s.service.StockTransactionInsert(stocktransaction)
	if err != nil {
		//helper.ReturnError(writer, request, err)+-
		return
	}
	payloads.NewHandleSuccess(writer, res, "Insert SuccessFull!", http.StatusCreated)
}

func StartStockTransactionControllerImpl(service transactionsparepartservice.StockTransactionService) StockTransactionController {
	return &StockTransactionControllerImpl{service: service}
}
