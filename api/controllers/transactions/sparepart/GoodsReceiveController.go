package transactionsparepartcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type GoodsReceiveController interface {
	GetAllGoodsReceive(writer http.ResponseWriter, request *http.Request)
	GetGoodsReceiveById(writer http.ResponseWriter, request *http.Request)
	InsertGoodsReceive(writer http.ResponseWriter, request *http.Request)
	UpdateGoodsReceive(writer http.ResponseWriter, request *http.Request)
	InsertGoodsReceiveDetail(writer http.ResponseWriter, request *http.Request)
	UpdateGoodsReceiveDetail(writer http.ResponseWriter, request *http.Request)
	LocationItemGoodsReceive(writer http.ResponseWriter, request *http.Request)
	DeleteGoodsReceive(writer http.ResponseWriter, request *http.Request)
	DeleteGoodsReceiveDetail(writer http.ResponseWriter, request *http.Request)
}

type GoodsReceiveControllerImpl struct {
	service transactionsparepartservice.GoodsReceiveService
}

func NewGoodsReceiveController(service transactionsparepartservice.GoodsReceiveService) GoodsReceiveController {
	return &GoodsReceiveControllerImpl{service: service}
}

// @Summary Get All Goods Receive
// @Description Get All Goods Receive
// @Tags Transaction : Sparepart Goods Receive
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param company_id query string false "Company ID"
// @Param goods_receive_document_number query string false "Goods Receive Document Number"
// @Param item_group_id query string false "Item Group ID"
// @Param reference_type_good_receive_id query string false "Reference Type Good Receive ID"
// @Param reference_document_number query string false "Reference Document Number"
// @Param supplier_id query string false "Supplier ID"
// @Param goods_receive_status_id query string false "Goods Receive Status ID"
// @Param supplier_delivery_order_number query string false "Supplier Delivery Order Number"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.ResponsePagination
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/goods-receive [get]
func (controller *GoodsReceiveControllerImpl) GetAllGoodsReceive(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"ig.company_id":                          queryValues.Get("company_id"),
		"ig.goods_receive_document_number":       queryValues.Get("goods_receive_document_number"),
		"itemgroup.item_group_id":                queryValues.Get("item_group_id"),
		"reftype.reference_type_good_receive_id": queryValues.Get("reference_type_good_receive_id"),
		"ig.reference_document_number":           queryValues.Get("reference_document_number"),
		"ig.supplier_id":                         queryValues.Get("supplier_id"),
		"ig.goods_receive_status_id":             queryValues.Get("goods_receive_status_id"),
		"ig.supplier_delivery_order_number":      queryValues.Get("supplier_delivery_order_number"),
	}
	paginations := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	filterCondition := utils.BuildFilterCondition(queryParams)
	res, err := controller.service.GetAllGoodsReceive(filterCondition, paginations)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, res.Rows, "Success Get All Data", http.StatusOK, res.Limit, res.Page, res.TotalRows, res.TotalPages)

}

// @Summary Get Goods Receive By ID
// @Description Get Goods Receive By ID
// @Tags Transaction : Sparepart Goods Receive
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param goods_receive_id path string true "Goods Receive ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/goods-receive/{goods_receive_id} [get]
func (controller *GoodsReceiveControllerImpl) GetGoodsReceiveById(writer http.ResponseWriter, request *http.Request) {
	GoodsReceiveSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "goods_receive_id"))
	res, err := controller.service.GetGoodsReceiveById(GoodsReceiveSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Get Data Successfully!", http.StatusOK)

}

// @Summary Insert Goods Receive
// @Description Insert Goods Receive
// @Tags Transaction : Sparepart Goods Receive
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param GoodsReceiveInsertPayloads body transactionsparepartpayloads.GoodsReceiveInsertPayloads true "Goods Receive Insert Payloads"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/goods-receive [post]
func (controller *GoodsReceiveControllerImpl) InsertGoodsReceive(writer http.ResponseWriter, request *http.Request) {
	var GoodsReceiveHeaderPayloads transactionsparepartpayloads.GoodsReceiveInsertPayloads
	helper.ReadFromRequestBody(request, &GoodsReceiveHeaderPayloads)
	res, err := controller.service.InsertGoodsReceive(GoodsReceiveHeaderPayloads)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Inserted Goods Receive Header", http.StatusCreated)
}

// @Summary Update Goods Receive
// @Description Update Goods Receive
// @Tags Transaction : Sparepart Goods Receive
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param goods_receive_id path string true "Goods Receive ID"
// @Param GoodsReceiveUpdatePayloads body transactionsparepartpayloads.GoodsReceiveUpdatePayloads true "Goods Receive Update Payloads"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/goods-receive/{goods_receive_id} [put]
func (controller *GoodsReceiveControllerImpl) UpdateGoodsReceive(writer http.ResponseWriter, request *http.Request) {
	var GoodsReceiveHeaderPayloads transactionsparepartpayloads.GoodsReceiveUpdatePayloads
	GoodsReceiveSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "goods_receive_id"))
	helper.ReadFromRequestBody(request, &GoodsReceiveHeaderPayloads)
	res, err := controller.service.UpdateGoodsReceive(GoodsReceiveHeaderPayloads, GoodsReceiveSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Update Goods Receive Header", http.StatusOK)
}

// @Summary Insert Goods Receive Detail
// @Description Insert Goods Receive Detail
// @Tags Transaction : Sparepart Goods Receive
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param GoodsReceiveDetailInsertPayloads body transactionsparepartpayloads.GoodsReceiveDetailInsertPayloads true "Goods Receive Detail Insert Payloads"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/goods-receive/detail [post]
func (controller *GoodsReceiveControllerImpl) InsertGoodsReceiveDetail(writer http.ResponseWriter, request *http.Request) {
	var GoodsReceiveDetailPayloads transactionsparepartpayloads.GoodsReceiveDetailInsertPayloads
	helper.ReadFromRequestBody(request, &GoodsReceiveDetailPayloads)
	res, err := controller.service.InsertGoodsReceiveDetail(GoodsReceiveDetailPayloads)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Insert Goods Receive Detail", http.StatusCreated)
}

// @Summary Update Goods Receive Detail
// @Description Update Goods Receive Detail
// @Tags Transaction : Sparepart Goods Receive
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param goods_receive_detail_system_number path string true "Goods Receive Detail System Number"
// @Param GoodsReceiveDetailUpdatePayloads body transactionsparepartpayloads.GoodsReceiveDetailUpdatePayloads true "Goods Receive Detail Update Payloads"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/goods-receive/detail/{goods_receive_detail_system_number} [put]
func (controller *GoodsReceiveControllerImpl) UpdateGoodsReceiveDetail(writer http.ResponseWriter, request *http.Request) {
	GoodsReceiveDetailSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "goods_receive_detail_system_number"))
	fmt.Println(GoodsReceiveDetailSystemNumber)
	var GoodsReceiveDetailPayloads transactionsparepartpayloads.GoodsReceiveDetailUpdatePayloads
	helper.ReadFromRequestBody(request, &GoodsReceiveDetailPayloads)
	res, err := controller.service.UpdateGoodsReceiveDetail(GoodsReceiveDetailPayloads, GoodsReceiveDetailSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	if res {
		payloads.NewHandleSuccess(writer, res, "Status updated successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// @Summary Location Item Goods Receive
// @Description Location Item Goods Receive
// @Tags Transaction : Sparepart Goods Receive
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param warehouse_location_name query string false "Warehouse Location Name"
// @Param item_code query string false "Item Code"
// @Param company_id query string false "Company ID"
// @Param warehouse_code query string false "Warehouse Code"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.ResponsePagination
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/goods-receive/location-item [get]
func (controller *GoodsReceiveControllerImpl) LocationItemGoodsReceive(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"B.warehouse_location_name": queryValues.Get("warehouse_location_name"),
		"item.item_code":            queryValues.Get("item_code"),
		"whs.company_id":            queryValues.Get("company_id"),
		"whs.warehouse_code":        queryValues.Get("warehouse_code"),
	}
	paginations := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	filterCondition := utils.BuildFilterCondition(queryParams)
	res, err := controller.service.LocationItemGoodsReceive(filterCondition, paginations)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, res.Rows, "Success Get All Data", http.StatusOK, res.Limit, res.Page, res.TotalRows, res.TotalPages)

	panic("implement me")
}

// @Summary Delete Goods Receive
// @Description Delete Goods Receive
// @Tags Transaction : Sparepart Goods Receive
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param goods_receive_id path string true "Goods Receive ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/goods-receive/{goods_receive_id} [delete]
func (controller *GoodsReceiveControllerImpl) DeleteGoodsReceive(writer http.ResponseWriter, request *http.Request) {
	GoodsReceiveSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "goods_receive_id"))
	res, err := controller.service.DeleteGoodsReceive(GoodsReceiveSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "goods receive deleted successfull", http.StatusOK)
}

// @Summary Delete Goods Receive Detail
// @Description Delete Goods Receive Detail
// @Tags Transaction : Sparepart Goods Receive
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param goods_receive_detail_id path string true "Goods Receive Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/goods-receive/detail/{goods_receive_detail_id} [delete]
func (controller *GoodsReceiveControllerImpl) DeleteGoodsReceiveDetail(writer http.ResponseWriter, request *http.Request) {
	GoodsReceiveDetailSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "goods_receive_detail_id"))
	res, err := controller.service.DeleteGoodsReceiveDetail(GoodsReceiveDetailSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "goods receive detail deleted successfull", http.StatusOK)
}
