package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/validation"
	"fmt"
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
	if validationErr := validation.ValidationForm(writer, request, &stocktransaction); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	res, err := s.service.StockTransactionInsert(stocktransaction)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	fmt.Println(res)
	payloads.NewHandleSuccess(writer, res, "Insert SuccessFull!", http.StatusCreated)
}

func StartStockTransactionControllerImpl(service transactionsparepartservice.StockTransactionService) StockTransactionController {
	return &StockTransactionControllerImpl{service: service}
}
