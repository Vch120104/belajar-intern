package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

// @Summary Insert Item Claim
// @Description Insert Item Claim
// @Tags Transaction Sparepart: Claim Supplier
// @Accept json
// @Produce json
// @Param InsertItemClaim body transactionsparepartpayloads.ClaimSupplierInsertPayload true "Insert Item Claim"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/cliam-supplier [post]
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

// @Summary Insert Item Claim Detail
// @Description Insert Item Claim Detail
// @Tags Transaction Sparepart: Claim Supplier
// @Accept json
// @Produce json
// @Param InsertItemClaimDetail body transactionsparepartpayloads.ClaimSupplierInsertDetailPayload true "Insert Item Claim Detail"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/cliam-supplier/detail [post]
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

// @Summary Get Item Claim By ID
// @Description Get Item Claim By ID
// @Tags Transaction Sparepart: Claim Supplier
// @Accept json
// @Produce json
// @Param claim_system_number path string true "Claim System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/cliam-supplier/{claim_system_number} [get]
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

// @Summary Submit Item Claim
// @Description Submit Item Claim
// @Tags Transaction Sparepart: Claim Supplier
// @Accept json
// @Produce json
// @Param claim_system_number path string true "Claim System Number"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/cliam-supplier/submit/{claim_system_number} [post]
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

// @Summary Get All Item Claim
// @Description Get All Item Claim
// @Tags Transaction Sparepart: Claim Supplier
// @Accept json
// @Produce json
// @Param claim_system_number query string false "Claim System Number"
// @Param goods_receive_document_number query string false "Goods Receive Document Number"
// @Param vehicle_brand_id query string false "Vehicle Brand ID"
// @Param profit_center_id query string false "Profit Center ID"
// @Param claim_document_number query string false "Claim Document Number"
// @Param reference_document_number query string false "Reference Document Number"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.ResponsePagination
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/cliam-supplier [get]
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

// @Summary Get Item Claim Detail By Header ID
// @Description Get Item Claim Detail By Header ID
// @Tags Transaction Sparepart: Claim Supplier
// @Accept json
// @Produce json
// @Param claim_system_number query string false "Claim System Number"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.ResponsePagination
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/cliam-supplier/detail [get]
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
