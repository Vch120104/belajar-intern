package transactionsparepartcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
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
}

type GoodsReceiveControllerImpl struct {
	service transactionsparepartservice.GoodsReceiveService
}

func NewGoodsReceiveController(service transactionsparepartservice.GoodsReceiveService) GoodsReceiveController {
	return &GoodsReceiveControllerImpl{service: service}
}
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

// get by id good receive
// goods-receive/{goods_receive_id}
func (controller *GoodsReceiveControllerImpl) GetGoodsReceiveById(writer http.ResponseWriter, request *http.Request) {
	GoodsReceiveSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "goods_receive_id"))
	res, err := controller.service.GetGoodsReceiveById(GoodsReceiveSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Get Data Successfully!", http.StatusOK)

}
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

// goods-receive/detail post
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

// /detail put
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

// goods-receive/location-item get
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

// delete goods receive
func (controller *GoodsReceiveControllerImpl) DeleteGoodsReceive(writer http.ResponseWriter, request *http.Request) {
	GoodsReceiveSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "goods_receive_id"))
	res, err := controller.service.DeleteGoodsReceive(GoodsReceiveSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "goods receive deleted successfull", http.StatusOK)
}
