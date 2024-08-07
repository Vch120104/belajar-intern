package transactionsparepartcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"net/http"
)

type PurchaseOrderControllerImpl struct {
	service transactionsparepartservice.PurchaseOrderService
}

type PurchaseOrderController interface {
	GetAllPurchaserOrderWithPagination(writer http.ResponseWriter, request *http.Request)
}

func NewPurchaseOrderImpl(PurchaseOrderService transactionsparepartservice.PurchaseOrderService) PurchaseOrderController {
	return &PurchaseOrderControllerImpl{service: PurchaseOrderService}
}
func (controller *PurchaseOrderControllerImpl) GetAllPurchaserOrderWithPagination(writer http.ResponseWriter, request *http.Request) {
	QueryValues := request.URL.Query()
	queryParams := map[string]string{
		"das": "adsa",
	}
	paginations := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(QueryValues, "limit"),
		Page:   utils.NewGetQueryInt(QueryValues, "page"),
		SortOf: QueryValues.Get("sort_of"),
		SortBy: QueryValues.Get("sort_by"),
	}
	filterCondition := utils.BuildFilterCondition(queryParams)
	res, err := controller.service.GetAllPurchaseOrder(filterCondition, paginations)
	if err != nil {
		helper.ReturnError(writer, request, err)
	}
	payloads.NewHandleSuccessPagination(writer, res, "Success Get All Data", 200, res.Limit, res.Page, res.TotalRows, res.TotalPages)
}
