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
)

type SupplySlipControllerImpl struct {
	supplyslipservice transactionsparepartservice.SupplySlipService
}

type SupplySlipController interface {
	GetSupplySlipByID(writer http.ResponseWriter, request *http.Request)
	GetAllSupplySlip(writer http.ResponseWriter, request *http.Request)
	SaveSupplySlip(writer http.ResponseWriter, request *http.Request)
	SaveSupplySlipDetail(writer http.ResponseWriter, request *http.Request)
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
	// Get ID from URL
	// id := mux.Vars(request)["id"]

	// Get data from service
	// data, err := r.supplyslipservice.GetSupplySlipByID(id)
	// if err != nil {
	// 	// Return error
	// 	exceptions.NewNotFoundException(writer, request, err)
	// 	return
	// }

	// Return success
	// payloads.NewHandleSuccess(writer, data, "Get Data Successfully", http.StatusOK)
}

func (r *SupplySlipControllerImpl) GetAllSupplySlip(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_supply_slip.supply_system_number": queryValues.Get("supply_system_number"),
		"supply_type_id":                       queryValues.Get("supply_type_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.supplyslipservice.GetAllSupplySlip(criteria, paginate)

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
