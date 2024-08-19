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

func NewPurchaseOrderControllerImpl(PurchaseOrderService transactionsparepartservice.PurchaseOrderService) PurchaseOrderController {
	return &PurchaseOrderControllerImpl{service: PurchaseOrderService}
}

// GetAllPurchaserOrderWithPagination
//
//	@Summary		Get All Purchase Order
//	@Description	REST API Purchase Order
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Purchase Order
//	@Param			page								query		string	true	"page"
//	@Param			limit								query		string	true	"limit"
//	@Param			purchase_order_document_number		query		string	false	"purchase_order_document_number"
//	@Param			item_group_id						query		string	false	"item_group_id"
//	@Param			order_type_id						query		string	false	"order_type_id"
//	@Param			purchase_order_status_id			query		string	false	"purchase_order_status_id"
//	@Param			warehouse_group_id					query		string	false	"warehouse_group_id"
//	@Param			warehouse_id						query		string	false	"warehouse_id"
//	@Param			supplier_id							query		string	false	"supplier_id"
//	@Param			cost_center_id						query		string	false	"cost_center_id"
//	@Param			purchase_order_date_from			query		string	false	"purchase_request_date_from"
//	@Param			purchase_order_date_to				query		string	false	"purchase_request_date_to"
//	@Param			purchase_request_document_number	query		string	false	"purchase_request_document_number"
//	@Param			created_by_user_id					query		string	false	"created_by_user_id"
//	@Param			sort_by								query		string	false	"sort_by"
//	@Param			sort_of								query		string	false	"sort_of"
//	@Success		200									{object}	[]transactionsparepartpayloads.PurchaseRequestGetAllListResponses
//	@Failure		500,400,401,404,403,422				{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-order/ [get]
func (controller *PurchaseOrderControllerImpl) GetAllPurchaserOrderWithPagination(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"purchase_order_document_number": queryValues.Get("purchase_order_document_number"),
		"A.item_group_id":                queryValues.Get("item_group_id"),
		"order_type_id":                  queryValues.Get("order_type_id"),
		"purchase_order_status_id":       queryValues.Get("purchase_order_status_id"),
		"warehouse_id":                   queryValues.Get("warehouse_id"),
		"warehouse_group_id":             queryValues.Get("warehouse_group_id"),
		"supplier_id":                    queryValues.Get("supplier_id"),
		"cost_center_id":                 queryValues.Get("cost_center_id"),
		"created_by_user_id":             queryValues.Get("created_by_user_id"),
	}
	DateParams := map[string]string{
		"purchase_request_date_from": queryValues.Get("purchase_request_date_from"),
		"purchase_request_date_to":   queryValues.Get("purchase_request_date_to"),
		"PurchaseRequestDocNo":       queryValues.Get("purchase_request_document_number"),
	}
	paginations := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	filterCondition := utils.BuildFilterCondition(queryParams)
	res, err := controller.service.GetAllPurchaseOrder(filterCondition, paginations, DateParams)
	if err != nil {
		helper.ReturnError(writer, request, err)
	}
	payloads.NewHandleSuccessPagination(writer, res, "Success Get All Data", 200, res.Limit, res.Page, res.TotalRows, res.TotalPages)
}
