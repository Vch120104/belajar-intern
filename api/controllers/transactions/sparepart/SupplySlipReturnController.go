package transactionsparepartcontroller

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
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

func (r *SupplySlipReturnControllerImpl) SaveSupplySlipReturn(writer http.ResponseWriter, request *http.Request) {

	var formRequest transactionsparepartentities.SupplySlipReturn
	helper.ReadFromRequestBody(request, &formRequest)
	var message string

	create, err := r.supplySlipReturnService.SaveSupplySlipReturn(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	message = "Create Data Successfully!"

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *SupplySlipReturnControllerImpl) SaveSupplySlipReturnDetail(writer http.ResponseWriter, request *http.Request) {

	var formRequest transactionsparepartentities.SupplySlipReturnDetail
	helper.ReadFromRequestBody(request, &formRequest)
	var message string

	create, err := r.supplySlipReturnService.SaveSupplySlipReturnDetail(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	message = "Create Data Successfully!"

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

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

	paginatedData, totalPages, totalRows, err := r.supplySlipReturnService.GetAllSupplySlipReturn(internalCriteria, externalCriteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

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

func (r *SupplySlipReturnControllerImpl) GetSupplySlipReturnDetailById(writer http.ResponseWriter, request *http.Request) {

	supplyReturnDetailId, _ := strconv.Atoi(chi.URLParam(request, "supply_return_detail_system_number"))

	result, err := r.supplySlipReturnService.GetSupplySlipReturnDetailById(supplyReturnDetailId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *SupplySlipReturnControllerImpl) UpdateSupplySlipReturn(writer http.ResponseWriter, request *http.Request) {
	supplyReturnId, _ := strconv.Atoi(chi.URLParam(request, "supply_return_system_number"))
	var formRequest transactionsparepartentities.SupplySlipReturn
	helper.ReadFromRequestBody(request, &formRequest)
	result, err := r.supplySlipReturnService.UpdateSupplySlipReturn(formRequest, supplyReturnId)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Update Data Successfully!", http.StatusOK)
}

func (r *SupplySlipReturnControllerImpl) UpdateSupplySlipReturnDetail(writer http.ResponseWriter, request *http.Request) {
	supplyReturnDetailId, _ := strconv.Atoi(chi.URLParam(request, "supply_return_detail_system_number"))
	var formRequest transactionsparepartentities.SupplySlipReturnDetail
	helper.ReadFromRequestBody(request, &formRequest)
	result, err := r.supplySlipReturnService.UpdateSupplySlipReturnDetail(formRequest, supplyReturnDetailId)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Update Data Successfully!", http.StatusOK)
}