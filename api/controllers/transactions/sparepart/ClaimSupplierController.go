package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type ClaimSupplierController interface {
	InsertItemClaim(writer http.ResponseWriter, request *http.Request)
	InsertItemClaimDetail(writer http.ResponseWriter, request *http.Request)
	GetItemClaimById(writer http.ResponseWriter, request *http.Request)
	GetItemClaimDetailByHeaderId(writer http.ResponseWriter, request *http.Request)
	SubmitItemClaim(writer http.ResponseWriter, request *http.Request)
	GetAllItemClaim(writer http.ResponseWriter, request *http.Request)
}
type ClaimSupplierControllerImpl struct {
	service transactionsparepartservice.ClaimSupplierService
}

func NewClaimSupplierControllerImpl(service transactionsparepartservice.ClaimSupplierService) ClaimSupplierController {
	return &ClaimSupplierControllerImpl{service: service}
}
func (controller *ClaimSupplierControllerImpl) InsertItemClaim(writer http.ResponseWriter, request *http.Request) {
	var claimSupplierInsertPayload transactionsparepartpayloads.ClaimSupplierInsertPayload
	helper.ReadFromRequestBody(request, &claimSupplierInsertPayload)
	res, err := controller.service.InsertItemClaim(claimSupplierInsertPayload)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Inserted item claim Header", http.StatusCreated)
}
func (controller *ClaimSupplierControllerImpl) InsertItemClaimDetail(writer http.ResponseWriter, request *http.Request) {
	var claimSupplierDetailInsertPayload transactionsparepartpayloads.ClaimSupplierInsertDetailPayload
	helper.ReadFromRequestBody(request, &claimSupplierDetailInsertPayload)
	res, err := controller.service.InsertItemClaimDetail(claimSupplierDetailInsertPayload)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	//payloads.NewHandleSuccess(writer, res, "Successfully Inserted item claim Detail", http.StatusInternalServerError)
	payloads.NewHandleSuccess(writer, res, "Successfully Inserted item claim Header", http.StatusCreated)

}
func (controller *ClaimSupplierControllerImpl) GetItemClaimById(writer http.ResponseWriter, request *http.Request) {
	claimIdStr := chi.URLParam(request, "claim_system_number")
	claimId, errs := strconv.Atoi(claimIdStr)
	if errs != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: errs, StatusCode: http.StatusInternalServerError})
		return
	}
	res, err := controller.service.GetItemClaimById(claimId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Inserted item claim Header", http.StatusCreated)
}
func (controller *ClaimSupplierControllerImpl) SubmitItemClaim(writer http.ResponseWriter, request *http.Request) {
	claimId := chi.URLParam(request, "claim_system_number")
	claimIds, err := strconv.Atoi(claimId)
	if err != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Err:        err,
			Message:    "failed to convert claim system number to int",
		})
	}
	res, serErr := controller.service.SubmitItemClaim(claimIds)
	if serErr != nil {
		helper.ReturnError(writer, request, serErr)
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Submitted item claim Header", http.StatusCreated)
}

func (controller *ClaimSupplierControllerImpl) GetAllItemClaim(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"claim_system_number":           queryValues.Get("claim_system_number"),
		"goods_receive_document_number": queryValues.Get("goods_receive_document_number"),
		"vehicle_brand_id":              queryValues.Get("vehicle_brand_id"),
		"profit_center_id":              queryValues.Get("profit_center_id"),
		"claim_document_number":         queryValues.Get("claim_document_number"),
		"reference_document_number":     queryValues.Get("reference_document_number"),
	}
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	filterCondition := utils.BuildFilterCondition(queryParams)
	result, err := controller.service.GetAllItemClaim(pagination, filterCondition)
	if err != nil {
		//helper.ReturnError(writer, request, err)
		exceptions.NewAppException(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)

}

func (controller *ClaimSupplierControllerImpl) GetItemClaimDetailByHeaderId(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"claim_system_number": queryValues.Get("claim_system_number"),
	}
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	filterCondition := utils.BuildFilterCondition(queryParams)
	result, err := controller.service.GetItemClaimDetailByHeaderId(pagination, filterCondition)
	if err != nil {
		//helper.ReturnError(writer, request, err)
		exceptions.NewAppException(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)

}
