package transactionsparepartcontroller

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
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

type SupplySlipControllerImpl struct {
	supplyslipservice transactionsparepartservice.SupplySlipService
}

type SupplySlipController interface {
	GetSupplySlipByID(writer http.ResponseWriter, request *http.Request)
	GetAllSupplySlip(writer http.ResponseWriter, request *http.Request)
	SaveSupplySlip(writer http.ResponseWriter, request *http.Request)
	SaveSupplySlipDetail(writer http.ResponseWriter, request *http.Request)
	GetSupplySlipDetailByID(writer http.ResponseWriter, request *http.Request)
	UpdateSupplySlip(writer http.ResponseWriter, request *http.Request)
	UpdateSupplySlipDetail(writer http.ResponseWriter, request *http.Request)
	SubmitSupplySlip(writer http.ResponseWriter, request *http.Request)
}

func NewSupplySlipController(supplyslipservice transactionsparepartservice.SupplySlipService) SupplySlipController {
	return &SupplySlipControllerImpl{
		supplyslipservice: supplyslipservice,
	}
}

// GetSupplySlipByID retrieves a supply slip by ID
// @Summary Get Supply Slip By ID
// @Description Retrieve a supply slip by its ID
// @Accept json
// @Produce json
// @Tags Transaction : Spare Part Supply Slip
// @Param supply_slip_id path int true "Supply Slip ID"
// @Success 200 {object} payloads.Response
// @Failure 500,404 {object} exceptions.BaseErrorResponse
// @Router /v1/supply-slip/{supply_slip_id} [get]
func (r *SupplySlipControllerImpl) GetSupplySlipByID(writer http.ResponseWriter, request *http.Request) {

	supplyId, _ := strconv.Atoi(chi.URLParam(request, "supply_system_number"))

	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result, err := r.supplyslipservice.GetSupplySliptById(supplyId, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *SupplySlipControllerImpl) GetAllSupplySlip(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()

	internalFilterCondition := map[string]string{
		"trx_supply_slip.supply_document_number":    queryValues.Get("supply_document_number"),
		"trx_supply_slip.supply_date_from":          queryValues.Get("supply_date_from"),
		"trx_supply_slip.supply_date_to":            queryValues.Get("supply_date_to"),
		"trx_work_order.work_order_document_number": queryValues.Get("work_order_document_number"),
	}

	externalFilterCondition := map[string]string{
		"supply_type_id":     queryValues.Get("supply_type_id"),
		"approval_status_id": queryValues.Get("approval_status_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	internalCriteria := utils.BuildFilterCondition(internalFilterCondition)
	externalCriteria := utils.BuildFilterCondition(externalFilterCondition)

	paginatedData, totalPages, totalRows, err := r.supplyslipservice.GetAllSupplySlip(internalCriteria, externalCriteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *SupplySlipControllerImpl) SaveSupplySlip(writer http.ResponseWriter, request *http.Request) {

	var formRequest transactionsparepartentities.SupplySlip
	helper.ReadFromRequestBody(request, &formRequest)
	var message string

	create, err := r.supplyslipservice.SaveSupplySlip(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	message = "Create Data Successfully!"

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *SupplySlipControllerImpl) SaveSupplySlipDetail(writer http.ResponseWriter, request *http.Request) {

	var formRequest transactionsparepartentities.SupplySlipDetail
	helper.ReadFromRequestBody(request, &formRequest)
	var message string

	create, err := r.supplyslipservice.SaveSupplySlipDetail(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	message = "Create Data Successfully!"

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *SupplySlipControllerImpl) GetSupplySlipDetailByID(writer http.ResponseWriter, request *http.Request) {

	supplyDetailId, _ := strconv.Atoi(chi.URLParam(request, "supply_detail_system_number"))

	result, err := r.supplyslipservice.GetSupplySlipDetailById(supplyDetailId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *SupplySlipControllerImpl) UpdateSupplySlip(writer http.ResponseWriter, request *http.Request) {
	supplyId, _ := strconv.Atoi(chi.URLParam(request, "supply_system_number"))
	var formRequest transactionsparepartentities.SupplySlip
	helper.ReadFromRequestBody(request, &formRequest)
	result, err := r.supplyslipservice.UpdateSupplySlip(formRequest, supplyId)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Update Data Successfully!", http.StatusOK)
}

func (r *SupplySlipControllerImpl) UpdateSupplySlipDetail(writer http.ResponseWriter, request *http.Request) {
	supplyDetailId, _ := strconv.Atoi(chi.URLParam(request, "supply_detail_system_number"))
	var formRequest transactionsparepartentities.SupplySlipDetail
	helper.ReadFromRequestBody(request, &formRequest)
	result, err := r.supplyslipservice.UpdateSupplySlipDetail(formRequest, supplyDetailId)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Update Data Successfully!", http.StatusOK)
}

func (r *SupplySlipControllerImpl) SubmitSupplySlip(writer http.ResponseWriter, request *http.Request) {
	supplySlipId := chi.URLParam(request, "supply_system_number")
	supplySlipInt, err := strconv.Atoi(supplySlipId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid supply slip ID", http.StatusBadRequest)
		return
	}

	success, newDocumentNumber, baseErr := r.supplyslipservice.SubmitSupplySlip(supplySlipInt)
	if baseErr != nil {
		if baseErr.Message == "Document number has already been generated" {
			payloads.NewHandleError(writer, baseErr.Message, http.StatusConflict)
		} else if baseErr.Message == "No supply slip data found" {
			payloads.NewHandleError(writer, baseErr.Message, http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		responseData := transactionsparepartpayloads.SubmitSupplySlipResponse{
			DocumentNumber:     newDocumentNumber,
			SupplySystemNumber: supplySlipInt,
		}
		payloads.NewHandleSuccess(writer, responseData, "supply slip submitted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to submit supply slip", http.StatusInternalServerError)
	}
}
