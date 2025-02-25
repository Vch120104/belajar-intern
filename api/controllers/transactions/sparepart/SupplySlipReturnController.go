package transactionsparepartcontroller

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SupplySlipReturnControllerImpl struct {
	supplySlipReturnService transactionsparepartservice.SupplySlipReturnService
}

type SupplySlipReturnController interface {
	SaveSupplySlipReturn(writer http.ResponseWriter, request *http.Request)
	SaveSupplySlipReturnDetail(writer http.ResponseWriter, request *http.Request)
	GetAllSupplySlipDetail(writer http.ResponseWriter, request *http.Request)
	GetSupplySlipReturnById(writer http.ResponseWriter, request *http.Request)
	GetSupplySlipReturnDetailById(writer http.ResponseWriter, request *http.Request)
	UpdateSupplySlipReturn(writer http.ResponseWriter, request *http.Request)
	UpdateSupplySlipReturnDetail(writer http.ResponseWriter, request *http.Request)
}

func NewSupplySlipReturnController(supplySlipReturnService transactionsparepartservice.SupplySlipReturnService) SupplySlipReturnController {
	return &SupplySlipReturnControllerImpl{
		supplySlipReturnService: supplySlipReturnService,
	}
}

// @Summary Save Supply Slip Return
// @Description Save Supply Slip Return
// @Tags Transaction : Sparepart Supply Slip Return
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param body body transactionsparepartentities.SupplySlipReturn true "Supply Slip Return Object"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/supply-slip-return [post]
func (r *SupplySlipReturnControllerImpl) SaveSupplySlipReturn(writer http.ResponseWriter, request *http.Request) {

	var formRequest transactionsparepartentities.SupplySlipReturn
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	var message string

	create, err := r.supplySlipReturnService.SaveSupplySlipReturn(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	message = "Create Data Successfully!"

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Save Supply Slip Return Detail
// @Description Save Supply Slip Return Detail
// @Tags Transaction : Sparepart Supply Slip Return
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param body body transactionsparepartentities.SupplySlipReturnDetail true "Supply Slip Return Detail Object"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/supply-slip-return/detail [post]
func (r *SupplySlipReturnControllerImpl) SaveSupplySlipReturnDetail(writer http.ResponseWriter, request *http.Request) {

	var formRequest transactionsparepartentities.SupplySlipReturnDetail
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	var message string

	create, err := r.supplySlipReturnService.SaveSupplySlipReturnDetail(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	message = "Create Data Successfully!"

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Get All Supply Slip Detail
// @Description Get All Supply Slip Detail
// @Tags Transaction : Sparepart Supply Slip Return
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param supply_return_document_number query string false "Supply Return Document Number"
// @Param supply_return_date_from query string false "Supply Return Date From"
// @Param supply_return_date_to query string false "Supply Return Date To"
// @Param supply_document_number query string false "Supply Document Number"
// @Param work_order_document_number query string false "Work Order Document Number"
// @Param customer_name query string false "Customer Name"
// @Param approval_status_id query string false "Approval Status Id"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/supply-slip-return [get]
func (r *SupplySlipReturnControllerImpl) GetAllSupplySlipDetail(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()

	internalFilterCondition := map[string]string{
		"trx_supply_slip_return.supply_return_document_number": queryValues.Get("supply_return_document_number"),
		"trx_supply_slip_return.supply_return_date_from":       queryValues.Get("supply_return_date_from"),
		"trx_supply_slip_return.supply_return_date_to":         queryValues.Get("supply_return_date_to"),
		"trx_supply_slip.supply_document_number":               queryValues.Get("supply_document_number"),
		"trx_work_order.work_order_document_number":            queryValues.Get("work_order_document_number"),
	}

	externalFilterCondition := map[string]string{
		"customer_name":      queryValues.Get("customer_name"),
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

	result, err := r.supplySlipReturnService.GetAllSupplySlipReturn(internalCriteria, externalCriteria, paginate)

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

// @Summary Get Supply Slip Return By Id
// @Description Get Supply Slip Return By Id
// @Tags Transaction : Sparepart Supply Slip Return
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param supply_return_system_number path int true "Supply Return System Number"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/supply-slip-return/{supply_return_system_number} [get]
func (r *SupplySlipReturnControllerImpl) GetSupplySlipReturnById(writer http.ResponseWriter, request *http.Request) {

	supplyReturnId, _ := strconv.Atoi(chi.URLParam(request, "supply_return_system_number"))

	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result, err := r.supplySlipReturnService.GetSupplySlipReturnById(supplyReturnId, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Supply Slip Return Detail By Id
// @Description Get Supply Slip Return Detail By Id
// @Tags Transaction : Sparepart Supply Slip Return
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param supply_return_detail_system_number path int true "Supply Return Detail System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/supply-slip-return/detail/{supply_return_detail_system_number} [get]
func (r *SupplySlipReturnControllerImpl) GetSupplySlipReturnDetailById(writer http.ResponseWriter, request *http.Request) {

	supplyReturnDetailId, _ := strconv.Atoi(chi.URLParam(request, "supply_return_detail_system_number"))

	result, err := r.supplySlipReturnService.GetSupplySlipReturnDetailById(supplyReturnDetailId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Update Supply Slip Return
// @Description Update Supply Slip Return
// @Tags Transaction : Sparepart Supply Slip Return
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param supply_return_system_number path int true "Supply Return System Number"
// @Param body body transactionsparepartentities.SupplySlipReturn true "Supply Slip Return Object"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/supply-slip-return/{supply_return_system_number} [put]
func (r *SupplySlipReturnControllerImpl) UpdateSupplySlipReturn(writer http.ResponseWriter, request *http.Request) {
	supplyReturnId, _ := strconv.Atoi(chi.URLParam(request, "supply_return_system_number"))
	var formRequest transactionsparepartentities.SupplySlipReturn
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	result, err := r.supplySlipReturnService.UpdateSupplySlipReturn(formRequest, supplyReturnId)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Update Data Successfully!", http.StatusOK)
}

// @Summary Update Supply Slip Return Detail
// @Description Update Supply Slip Return Detail
// @Tags Transaction : Sparepart Supply Slip Return
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param supply_return_detail_system_number path int true "Supply Return Detail System Number"
// @Param body body transactionsparepartentities.SupplySlipReturnDetail true "Supply Slip Return Detail Object"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/supply-slip-return/detail/{supply_return_detail_system_number} [put]
func (r *SupplySlipReturnControllerImpl) UpdateSupplySlipReturnDetail(writer http.ResponseWriter, request *http.Request) {
	supplyReturnDetailId, _ := strconv.Atoi(chi.URLParam(request, "supply_return_detail_system_number"))
	var formRequest transactionsparepartentities.SupplySlipReturnDetail
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	result, err := r.supplySlipReturnService.UpdateSupplySlipReturnDetail(formRequest, supplyReturnDetailId)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Update Data Successfully!", http.StatusOK)
}
