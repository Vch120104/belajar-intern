package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PurchaseOrderControllerImpl struct {
	service transactionsparepartservice.PurchaseOrderService
}

type PurchaseOrderController interface {
	GetAllPurchaserOrderWithPagination(writer http.ResponseWriter, request *http.Request)
	GetByIdPurchaseOrder(writer http.ResponseWriter, request *http.Request)
	GetPurchaseOrderDetailByHeaderId(writer http.ResponseWriter, request *http.Request)
	NewPurchaseOrderHeader(writer http.ResponseWriter, request *http.Request)
	UpdatePurchaseOrderHeader(writer http.ResponseWriter, request *http.Request)
	GetPurchaseOrderDetailById(writer http.ResponseWriter, request *http.Request)
	NewPurchaseOrderDetail(writer http.ResponseWriter, request *http.Request)
	DeletePurchaseOrderDetailMultiId(writer http.ResponseWriter, request *http.Request)
	SavePurchaseOrderDetail(writer http.ResponseWriter, request *http.Request)
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
//
// @Tags Transaction Sparepart: Purchase Order
//
//	@Param			page								query		string	true	"page"
//	@Param			limit								query		string	true	"limit"
//	@Param			purchase_order_document_number		query		string	false	"purchase_order_document_number"
//	@Param			item_group_id						query		string	false	"item_group_id"
//	@Param			order_type_id						query		string	false	"order_type_id"
//	@Param			purchase_order_date_from			query		string	false	"purchase_order_date_from"
//	@Param			purchase_order_date_to				query		string	false	"purchase_order_date_to"
//	@Param			purchase_order_status_id			query		string	false	"purchase_order_status_id"
//	@Param			warehouse_group_id					query		string	false	"warehouse_group_id"
//	@Param			warehouse_id						query		string	false	"warehouse_id"
//	@Param			supplier_id							query		string	false	"supplier_id"
//	@Param			cost_center_id						query		string	false	"cost_center_id"
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
		"trx_item_purchase_order.purchase_order_document_number": queryValues.Get("purchase_order_document_number"),
		"trx_item_purchase_order.item_group_id":                  queryValues.Get("item_group_id"),
		"trx_item_purchase_order.order_type_id":                  queryValues.Get("order_type_id"),
		"trx_item_purchase_order.purchase_order_status_id":       queryValues.Get("purchase_order_status_id"),
		"trx_item_purchase_order.warehouse_id":                   queryValues.Get("warehouse_id"),
		"trx_item_purchase_order.warehouse_group_id":             queryValues.Get("warehouse_group_id"),
		"trx_item_purchase_order.supplier_id":                    queryValues.Get("supplier_id"),
		"trx_item_purchase_order.cost_center_id":                 queryValues.Get("cost_center_id"),
		"trx_item_purchase_order.created_by_user_id":             queryValues.Get("created_by_user_id"),
	}
	DateParams := map[string]string{
		"purchase_order_date_from": queryValues.Get("purchase_order_date_from"),
		"purchase_order_date_to":   queryValues.Get("purchase_order_date_to"),
		"PurchaseRequestDocNo":     queryValues.Get("purchase_request_document_number"),
	}
	paginations := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	fmt.Println("dfgdfgdfgdfg")
	filterCondition := utils.BuildFilterCondition(queryParams)
	res, err := controller.service.GetAllPurchaseOrder(filterCondition, paginations, DateParams)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, res.Rows, "Success Get All Data", 200, res.Limit, res.Page, res.TotalRows, res.TotalPages)
}

// GetByIdPurchaseOrder
//
//	@Summary		Get By Id Purchase Order
//	@Description	REST API Get By Id Purchase Order
//	@Accept			json
//	@Produce		json
//
// @Tags Transaction Sparepart: Purchase Order
//
//	@Param			purchase_order_system_number		path		string	false	"purchase_order_system_number"
//	@Success		200									{object}	transactionsparepartpayloads.PurchaseOrderGetByIdResponses
//	@Failure		500,400,401,404,403,422				{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-order/by-id/{purchase_order_system_number} [get]
func (controller *PurchaseOrderControllerImpl) GetByIdPurchaseOrder(writer http.ResponseWriter, request *http.Request) {
	purchaseOrderSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "purchase_order_system_number"))
	result, err := controller.service.GetByIdPurchaseOrder(purchaseOrderSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)

}

// GetPurchaseOrderDetailByHeaderId
//
//	@Summary		Get Purchase Reqeust Detail
//	@Description	REST API Get Purchase Request Detail
//	@Accept			json
//	@Produce		json
//
// @Tags Transaction Sparepart: Purchase Order
//
//	@Param			page								query		string	true	"page"
//	@Param			limit								query		string	true	"limit"
//	@Param			sort_by								query		string	false	"sort_by"
//	@Param			sort_of								query		string	false	"sort_of"
//	@Param			purchase_order_system_number		path		string	false	"purchase_order_system_number"
//	@Success		200									{object}	transactionsparepartpayloads.PurchaseOrderGetDetail
//	@Failure		500,400,401,404,403,422				{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-order/detail [get]
func (controller *PurchaseOrderControllerImpl) GetPurchaseOrderDetailByHeaderId(writer http.ResponseWriter, request *http.Request) {
	//purchaseOrderSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "purchase_order_system_number"))
	queryValues := request.URL.Query()
	purchaseOrderSystemNumbers, _ := strconv.Atoi(queryValues.Get("purchase_order_system_number"))
	paginations := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	result, err := controller.service.GetByIdPurchaseOrderDetail(purchaseOrderSystemNumbers, paginations)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// NewPurchaseOrderHeader
//
//	@Summary		Create New Purchase Order Header2
//	@Description	Create New Purchase Order Header
//	@Accept			json
//	@Produce		json
//
// @Tags Transaction Sparepart: Purchase Order
//
//	@Param			reqBody					body		transactionsparepartpayloads.PurchaseOrderNewPurchaseOrderResponses	true	"Purchase Request Header Data"
//	@Success		201						{object}	payloads.Response
//	@Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-order [post]
func (controller *PurchaseOrderControllerImpl) NewPurchaseOrderHeader(writer http.ResponseWriter, request *http.Request) {

	var purchaseRequest transactionsparepartpayloads.PurchaseOrderNewPurchaseOrderResponses

	helper.ReadFromRequestBody(request, &purchaseRequest)
	if validationErr := validation.ValidationForm(writer, request, &purchaseRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := controller.service.NewPurchaseOrderHeader(purchaseRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "save success", http.StatusCreated)
}

// UpdatePurchaseOrderHeader
//
//	@Summary		Update Purchase Request order
//	@Description	Update Purchase Request order
//	@Accept			json
//	@Produce		json
//
// @Tags Transaction Sparepart: Purchase Order
//
//	@Param			purchase_order_system_number	path		int	true	"purchase_order_system_number"
//	@Param			reqBody					body		transactionsparepartpayloads.PurchaseOrderNewPurchaseOrderPayloads	true	"Purchase Request Header Data"
//	@Success		201						{object}	payloads.Response
//	@Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-order/{purchase_order_system_number} [put]
func (controller *PurchaseOrderControllerImpl) UpdatePurchaseOrderHeader(writer http.ResponseWriter, request *http.Request) {
	var puchaseRequestHeader transactionsparepartpayloads.PurchaseOrderNewPurchaseOrderPayloads
	PurchaseOrderSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "purchase_order_system_number"))

	helper.ReadFromRequestBody(request, &puchaseRequestHeader)
	if validationErr := validation.ValidationForm(writer, request, &puchaseRequestHeader); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	success, err := controller.service.UpdatePurchaseOrderHeader(PurchaseOrderSystemNumber, puchaseRequestHeader)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "save success", http.StatusOK)
}

// GetPurchaseOrderDetailById
//
//	@Summary		Get Purchase Order Detail Per Id
//	@Description	Get Purchase Order Detail Per Id
//	@Accept			json
//	@Produce		json
//
// @Tags Transaction Sparepart: Purchase Order
//
//	@Param			purchase_order_detail_system_number	path		int	true	"purchase_order_detail_system_number"	true	"Purchase Request Header Data"
//	@Success		200						{object}	transactionsparepartpayloads.PurchaseOrderGetDetail
//	@Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-order/detail/by-id/{purchase_order_detail_system_number} [get]
func (controller *PurchaseOrderControllerImpl) GetPurchaseOrderDetailById(writer http.ResponseWriter, request *http.Request) {

	PurchaseOrderSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "purchase_order_detail_system_number"))

	//helper.ReadFromRequestBody(request, &puchaseRequestHeader)
	success, err := controller.service.GetPurchaseOrderDetailById(PurchaseOrderSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "Get Data Success", http.StatusOK)
}

// NewPurchaseOrderDetail
//
//	@Summary		Create New Purchase Order Detail
//	@Description	Create New Purchase Order Detail
//	@Accept			json
//	@Produce		json
//
// @Tags Transaction Sparepart: Purchase Order
//
//	@Param			reqBody					body		transactionsparepartpayloads.PurchaseOrderDetailPayloads	true	"Purchase Request Header Data"
//	@Success		201						{object}	transactionsparepartentities.PurchaseOrderDetailEntities
//	@Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-order/detail [post]
func (controller *PurchaseOrderControllerImpl) NewPurchaseOrderDetail(writer http.ResponseWriter, request *http.Request) {
	var payload transactionsparepartpayloads.PurchaseOrderDetailPayloads

	helper.ReadFromRequestBody(request, &payload)
	if validationErr := validation.ValidationForm(writer, request, &payload); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := controller.service.NewPurchaseOrderDetail(payload)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "Create Data Success", http.StatusCreated)
}

// DeletePurchaseOrderDetailMultiId
//
// @Summary			Void Detail Purchase Order Detail Multi Id
// @Description		Void Detail Purchase Order Detail Multi Id
// @Accept			json
// @Produce			json
// @Tags			Transaction : Purchase Order
// @Param			purchase_order_detail_system_number	path		string true	"purchase_order_detail_system_number"
// @Success		201						{object}	payloads.Response
// @Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
// @Router			/v1/purchase-order/detail/{purchase_order_detail_system_number} [delete]
func (controller *PurchaseOrderControllerImpl) DeletePurchaseOrderDetailMultiId(writer http.ResponseWriter, request *http.Request) {
	PurchaseOrderDetailMultiId := chi.URLParam(request, "purchase_order_detail_system_number")
	success, baseErr := controller.service.DeletePurchaseOrderDetailMultiId(PurchaseOrderDetailMultiId)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			helper.ReturnError(writer, request, baseErr)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}
	if success {
		payloads.NewHandleSuccess(writer, nil, "Purchase Order Detail voided successfully", http.StatusOK)
	} else {
		helper.ReturnError(writer, request, baseErr)
	}
}

// SavePurchaseOrderDetail
//
//	@Summary		Save Purchase Order Detail
//	@Description	Save Purchase Order Detail
//	@Accept			json
//	@Produce		json
//
// @Tags Transaction Sparepart: Purchase Order
//
//	@Param			reqBody					body		transactionsparepartpayloads.PurchaseOrderSaveDetailPayloads	true	"Purchase Request Detail Data"
//	@Success		201						{object}	transactionsparepartentities.PurchaseOrderDetailEntities
//	@Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-order/detail [patch]
func (controller *PurchaseOrderControllerImpl) SavePurchaseOrderDetail(writer http.ResponseWriter, request *http.Request) {
	var payload transactionsparepartpayloads.PurchaseOrderSaveDetailPayloads
	helper.ReadFromRequestBody(request, &payload)
	if validationErr := validation.ValidationForm(writer, request, &payload); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	success, err := controller.service.SavePurchaseOrderDetail(payload)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "save Detail Success", http.StatusOK)
}
