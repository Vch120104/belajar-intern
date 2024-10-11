package transactionsparepartcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type BinningListController interface {
	GetBinningListById(writer http.ResponseWriter, request *http.Request)
}
type BinningListControllerImpl struct {
	service transactionsparepartservice.BinningListService
}

func NewBinningListControllerImpl(service transactionsparepartservice.BinningListService) BinningListController {
	return &BinningListControllerImpl{service: service}
}
func (controller *BinningListControllerImpl) GetBinningListById(writer http.ResponseWriter, request *http.Request) {
	BinningStockSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "binning_stock_system_number"))
	result, err := controller.service.GetBinningListById(BinningStockSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Successfully Get Data Binning Stock !", http.StatusOK)
}
