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
	"after-sales/api/validation"
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
// @Tags Transaction Spare Part : Supply Slip
// @Param supply_slip_id path int true "Supply Slip ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
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

// @Summary Get All Supply Slip
// @Description Get All Supply Slip
// @Tags Transaction Spare Part : Supply Slip
// @Accept json
// @Produce json
// @Param supply_document_number query string false "Supply Document Number"
// @Param supply_date_from query string false "Supply Date From"
// @Param supply_date_to query string false "Supply Date To"
// @Param work_order_document_number query string false "Work Order Document Number"
// @Param supply_type_id query string false "Supply Type ID"
// @Param approval_status_id query string false "Approval Status ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/supply-slip [get]
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

	result, err := r.supplyslipservice.GetAllSupplySlip(internalCriteria, externalCriteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		result.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		result.Limit,
		result.Page,
		int64(result.TotalRows),
		result.TotalPages,
	)
}

// @Summary Save Supply Slip
// @Description Save Supply Slip
// @Tags Transaction Spare Part : Supply Slip
// @Accept json
// @Produce json
// @Param SupplySlip body transactionsparepartentities.SupplySlip true "Supply Slip"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/supply-slip [post]
func (r *SupplySlipControllerImpl) SaveSupplySlip(writer http.ResponseWriter, request *http.Request) {

	var formRequest transactionsparepartentities.SupplySlip
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	var message string

	create, err := r.supplyslipservice.SaveSupplySlip(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	message = "Create Data Successfully!"

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Save Supply Slip Detail
// @Description Save Supply Slip Detail
// @Tags Transaction Spare Part : Supply Slip
// @Accept json
// @Produce json
// @Param SupplySlipDetail body transactionsparepartentities.SupplySlipDetail true "Supply Slip Detail"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/supply-slip/detail [post]
func (r *SupplySlipControllerImpl) SaveSupplySlipDetail(writer http.ResponseWriter, request *http.Request) {

	var formRequest transactionsparepartentities.SupplySlipDetail
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	var message string

	create, err := r.supplyslipservice.SaveSupplySlipDetail(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	message = "Create Data Successfully!"

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Get Supply Slip Detail By ID
// @Description Get Supply Slip Detail By ID
// @Tags Transaction Spare Part : Supply Slip
// @Accept json
// @Produce json
// @Param supply_detail_system_number path int true "Supply Detail System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/supply-slip/detail/{supply_detail_system_number} [get]
func (r *SupplySlipControllerImpl) GetSupplySlipDetailByID(writer http.ResponseWriter, request *http.Request) {

	supplyDetailId, _ := strconv.Atoi(chi.URLParam(request, "supply_detail_system_number"))

	result, err := r.supplyslipservice.GetSupplySlipDetailById(supplyDetailId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Update Supply Slip
// @Description Update Supply Slip
// @Tags Transaction Spare Part : Supply Slip
// @Accept json
// @Produce json
// @Param supply_system_number path int true "Supply System Number"
// @Param SupplySlip body transactionsparepartentities.SupplySlip true "Supply Slip"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/supply-slip/{supply_system_number} [put]
func (r *SupplySlipControllerImpl) UpdateSupplySlip(writer http.ResponseWriter, request *http.Request) {
	supplyId, _ := strconv.Atoi(chi.URLParam(request, "supply_system_number"))
	var formRequest transactionsparepartentities.SupplySlip
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	result, err := r.supplyslipservice.UpdateSupplySlip(formRequest, supplyId)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Update Data Successfully!", http.StatusOK)
}

// @Summary Update Supply Slip Detail
// @Description Update Supply Slip Detail
// @Tags Transaction Spare Part : Supply Slip
// @Accept json
// @Produce json
// @Param supply_detail_system_number path int true "Supply Detail System Number"
// @Param SupplySlipDetail body transactionsparepartentities.SupplySlipDetail true "Supply Slip Detail"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/supply-slip/detail/{supply_detail_system_number} [put]
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

// @Summary Submit Supply Slip
// @Description Submit Supply Slip
// @Tags Transaction Spare Part : Supply Slip
// @Accept json
// @Produce json
// @Param supply_system_number path int true "Supply System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/supply-slip/submit/{supply_system_number} [put]
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
