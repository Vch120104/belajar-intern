package mastercontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type LookupController interface {
	ItemOprCode(writer http.ResponseWriter, request *http.Request)
	CampaignMaster(writer http.ResponseWriter, request *http.Request)
	ItemOprCodeWithPrice(writer http.ResponseWriter, request *http.Request)
	VehicleUnitMaster(writer http.ResponseWriter, request *http.Request)
	GetVehicleUnitByID(writer http.ResponseWriter, request *http.Request)
	GetVehicleUnitByChassisNumber(writer http.ResponseWriter, request *http.Request)
	WorkOrderService(writer http.ResponseWriter, request *http.Request)
}

type LookupControllerImpl struct {
	LookupService masterservice.LookupService
}

func NewLookupController(LookupService masterservice.LookupService) LookupController {
	return &LookupControllerImpl{
		LookupService: LookupService,
	}
}

func (r *LookupControllerImpl) ItemOprCode(writer http.ResponseWriter, request *http.Request) {
	linetypeStrId := chi.URLParam(request, "linetype_id")
	linetypeId, err := strconv.Atoi(linetypeStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Line Type ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()
	queryParams := map[string]string{}
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	criteria := utils.BuildFilterCondition(queryParams)
	lookup, totalPages, totalRows, baseErr := r.LookupService.ItemOprCode(linetypeId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(writer, lookup, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *LookupControllerImpl) ItemOprCodeWithPrice(writer http.ResponseWriter, request *http.Request) {
	linetypeStrId := chi.URLParam(request, "linetype_id")
	linetypeId, err := strconv.Atoi(linetypeStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Line Type ID", http.StatusBadRequest)
		return
	}

	companyStrId := chi.URLParam(request, "company_id")
	companyId, err := strconv.Atoi(companyStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Company ID", http.StatusBadRequest)
		return
	}

	operationItemCodeStrId := chi.URLParam(request, "operation_item_id")
	operationItemId, err := strconv.Atoi(operationItemCodeStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Operation Item ID", http.StatusBadRequest)
		return
	}

	brandStrId := chi.URLParam(request, "brand_id")
	brandId, err := strconv.Atoi(brandStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Brand ID", http.StatusBadRequest)
		return
	}

	modelStrId := chi.URLParam(request, "model_id")
	modelId, err := strconv.Atoi(modelStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Model ID", http.StatusBadRequest)
		return
	}

	jobTypeStrId := chi.URLParam(request, "job_type_id")
	jobTypeId, err := strconv.Atoi(jobTypeStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Job Type ID", http.StatusBadRequest)
		return
	}

	variantStrId := chi.URLParam(request, "variant_id")
	variantId, err := strconv.Atoi(variantStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Variant ID", http.StatusBadRequest)
		return
	}

	currencyStrId := chi.URLParam(request, "currency_id")
	currencyId, err := strconv.Atoi(currencyStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Currency ID", http.StatusBadRequest)
		return
	}

	billCodeStrId := chi.URLParam(request, "bill_code")
	if billCodeStrId == "" {
		payloads.NewHandleError(writer, "Invalid Billcode", http.StatusBadRequest)
		return
	}

	whsGroupStrId := chi.URLParam(request, "warehouse_group")
	if whsGroupStrId == "" {
		payloads.NewHandleError(writer, "Invalid Warehouse", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()
	queryParams := map[string]string{}
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	lookup, totalPages, totalRows, baseErr := r.LookupService.ItemOprCodeWithPrice(linetypeId, companyId, operationItemId, brandId, modelId, jobTypeId, variantId, currencyId, billCodeStrId, whsGroupStrId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(writer, lookup, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *LookupControllerImpl) VehicleUnitMaster(writer http.ResponseWriter, request *http.Request) {
	brandStrId := chi.URLParam(request, "brand_id")
	brandId, err := strconv.Atoi(brandStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Brand ID", http.StatusBadRequest)
		return
	}

	modelStrId := chi.URLParam(request, "model_id")
	modelId, err := strconv.Atoi(modelStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Model ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()
	queryParams := map[string]string{}
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	criteria := utils.BuildFilterCondition(queryParams)
	lookup, totalPages, totalRows, baseErr := r.LookupService.VehicleUnitMaster(brandId, modelId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(writer, lookup, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *LookupControllerImpl) GetVehicleUnitByID(writer http.ResponseWriter, request *http.Request) {
	vehicleStrId := chi.URLParam(request, "vehicle_id")
	vehicleId, err := strconv.Atoi(vehicleStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Vehicle ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()
	queryParams := map[string]string{}
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	criteria := utils.BuildFilterCondition(queryParams)
	lookup, totalPages, totalRows, baseErr := r.LookupService.GetVehicleUnitByID(vehicleId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(writer, lookup, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *LookupControllerImpl) GetVehicleUnitByChassisNumber(writer http.ResponseWriter, request *http.Request) {
	chassisNumber := chi.URLParam(request, "vehicle_chassis_number")

	queryValues := request.URL.Query()
	queryParams := map[string]string{}
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	lookup, totalPages, totalRows, baseErr := r.LookupService.GetVehicleUnitByChassisNumber(chassisNumber, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(writer, lookup, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *LookupControllerImpl) CampaignMaster(writer http.ResponseWriter, request *http.Request) {
	companyStrId := chi.URLParam(request, "company_id")
	companyId, err := strconv.Atoi(companyStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Campaign ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()
	queryParams := map[string]string{}
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	criteria := utils.BuildFilterCondition(queryParams)
	lookup, totalPages, totalRows, baseErr := r.LookupService.CampaignMaster(companyId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(writer, lookup, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *LookupControllerImpl) WorkOrderService(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()
	queryParams := map[string]string{}
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	lookup, totalPages, totalRows, baseErr := r.LookupService.WorkOrderService(paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(writer, lookup, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}
