package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type SalesOrderController interface {
	InsertSalesOrderHeader(writer http.ResponseWriter, request *http.Request)
	GetSalesOrderByID(writer http.ResponseWriter, request *http.Request)
	GetAllSalesOrder(writer http.ResponseWriter, request *http.Request)
	VoidSalesOrder(writer http.ResponseWriter, request *http.Request)
	InsertSalesOrderDetail(writer http.ResponseWriter, request *http.Request)
	DeleteSalesOrderDetail(writer http.ResponseWriter, request *http.Request)
	SalesOrderProposedDiscountMultiId(writer http.ResponseWriter, request *http.Request)
	UpdateSalesOrderHeader(writer http.ResponseWriter, request *http.Request)
	SubmitSalesOrderHeader(writer http.ResponseWriter, request *http.Request)
}

type SalesOrderControllerImpl struct {
	SalesOrderService transactionsparepartservice.SalesOrderServiceInterface
}

func StartSalesOrderControllerImpl(SalesOrderService transactionsparepartservice.SalesOrderServiceInterface) SalesOrderController {
	return &SalesOrderControllerImpl{SalesOrderService: SalesOrderService}
}

func (s *SalesOrderControllerImpl) InsertSalesOrderHeader(writer http.ResponseWriter, request *http.Request) {
	SalesOrderBody := transactionsparepartpayloads.SalesOrderInsertHeaderPayload{}
	errReadRequestBody := jsonchecker.ReadFromRequestBody(request, &SalesOrderBody)
	if errReadRequestBody != nil {
		helper.ReturnError(writer, request, errReadRequestBody)
	}
	errValidation := validation.ValidationForm(writer, request, &SalesOrderBody)
	if errValidation != nil {
		helper.ReturnError(writer, request, errValidation)
		return
	}
	res, err := s.SalesOrderService.InsertSalesOrderHeader(SalesOrderBody)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "success to insert sales order estimation", http.StatusCreated)
}

func (s *SalesOrderControllerImpl) GetSalesOrderByID(writer http.ResponseWriter, request *http.Request) {
	salesOrderIdStr := chi.URLParam(request, "sales_order_system_number")
	salesOrderId, errId := strconv.Atoi(salesOrderIdStr)
	if errId != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get sales order id please check input",
			Err:        errId,
		})
		return
	}
	res, err := s.SalesOrderService.GetSalesOrderByID(salesOrderId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "success to get sales order estimation", http.StatusOK)
}
func (s *SalesOrderControllerImpl) GetAllSalesOrder(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	queryParams := map[string]string{
		// tambahin item group, stock keep, brand_id, stock_keeping = 1
		// item detail tidak boleh duplicate
		//
		"mtr_item.item_group_id": queryValues.Get("item_group_id"),
		"DT.brand_id":            queryValues.Get("brand_id"),
		"item_code":              queryValues.Get("item_code"),
		"item_name":              queryValues.Get("item_name"),
		"item_class_name":        queryValues.Get("item_class_name"),
		"IT.item_type_code":      queryValues.Get("item_type_code"),
		"L1.item_level_1_code":   queryValues.Get("item_level_1"),
		"L2.item_level_2_code":   queryValues.Get("item_level_2"),
		"L3.item_level_3_code":   queryValues.Get("item_level_3"),
		"L4.item_level_4_code":   queryValues.Get("item_level_4"),
		"quantity":               queryValues.Get("quantity"),
	}
	filter := utils.BuildFilterCondition(queryParams)
	result, err := s.SalesOrderService.GetAllSalesOrder(pagination, filter)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data successfully", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)

}
func (s *SalesOrderControllerImpl) VoidSalesOrder(writer http.ResponseWriter, request *http.Request) {
	salesOrderIdStr := chi.URLParam(request, "sales_order_system_number")
	salesOrderId, errId := strconv.Atoi(salesOrderIdStr)
	if errId != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errId,
			Message:    "sales order system number is not found from path",
		})
	}
	res, err := s.SalesOrderService.VoidSalesOrder(salesOrderId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "success to void sales order", http.StatusOK)
}
func (s *SalesOrderControllerImpl) InsertSalesOrderDetail(writer http.ResponseWriter, request *http.Request) {
	var formRequest transactionsparepartpayloads.SalesOrderDetailInsertPayload
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	res, err := s.SalesOrderService.InsertSalesOrderDetail(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "succesfull insert sales order detail", http.StatusCreated)
	return
}
func (s *SalesOrderControllerImpl) DeleteSalesOrderDetail(writer http.ResponseWriter, request *http.Request) {
	salesOrderDetailId, errConvert := strconv.Atoi(chi.URLParam(request, "sales_order_detail_system_number"))
	if errConvert != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errConvert,
			Message:    "failed to convert id parameters",
		})
	}
	res, err := s.SalesOrderService.DeleteSalesOrderDetail(salesOrderDetailId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	var ResponseDeleteSalesOrder transactionsparepartpayloads.SalesOrderDeleteDetailResponse
	if res == true {
		ResponseDeleteSalesOrder = transactionsparepartpayloads.SalesOrderDeleteDetailResponse{
			DeleteMessage: "success to delete sales order detail",
			DeleteStatus:  true,
		}
	} else {
		helper.ReturnError(writer, request, err)
	}
	payloads.NewHandleSuccess(writer, ResponseDeleteSalesOrder, "success to delete sales order detail", http.StatusOK)
}
func (s *SalesOrderControllerImpl) SalesOrderProposedDiscountMultiId(writer http.ResponseWriter, request *http.Request) {
	queryValue := request.URL.Query()
	salesOrderDetailMultiId := chi.URLParam(request, "sales_order_detail_multi_id")
	proposedDiscount := utils.NewGetQueryfloat(queryValue, "proposed_discount")

	result, errResult := s.SalesOrderService.SalesOrderProposedDiscountMultiId(salesOrderDetailMultiId, proposedDiscount)
	if errResult != nil {
		helper.ReturnError(writer, request, errResult)
		return
	}
	payloads.NewHandleSuccess(writer, result, "success to insert proposed discount", http.StatusOK)

}
func (s *SalesOrderControllerImpl) UpdateSalesOrderHeader(writer http.ResponseWriter, request *http.Request) {
	salesOrderId := chi.URLParam(request, "sales_order_system_number")
	salesOrderSystemNumber, errConvertInt := strconv.Atoi(salesOrderId)
	if errConvertInt != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errConvertInt,
			Message:    "failed to convert id parameters",
		})
	}
	requestBody := transactionsparepartpayloads.SalesOrderUpdatePayload{}
	err := jsonchecker.ReadFromRequestBody(request, &requestBody)
	if err != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to read request body",
		})
		return
	}
	res, errResult := s.SalesOrderService.UpdateSalesOrderHeader(requestBody, salesOrderSystemNumber)
	if errResult != nil {
		helper.ReturnError(writer, request, errResult)
		return
	}
	payloads.NewHandleSuccess(writer, res, "success to update sales order header", http.StatusOK)
}
func (s *SalesOrderControllerImpl) SubmitSalesOrderHeader(writer http.ResponseWriter, request *http.Request) {
	//get sales order
	salesOrderId, errConvertId := strconv.Atoi(chi.URLParam(request, "sales_order_system_number"))
	if errConvertId != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errConvertId,
			Message:    "failed to convert id parameters",
		})
		return
	}
	res, err := s.SalesOrderService.SubmitSalesOrderHeader(salesOrderId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "success to submit sales order header", http.StatusOK)
}
