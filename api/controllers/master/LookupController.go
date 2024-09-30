package mastercontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type LookupController interface {
	ItemOprCode(writer http.ResponseWriter, request *http.Request)
	ItemOprCodeByCode(writer http.ResponseWriter, request *http.Request)
	ItemOprCodeByID(writer http.ResponseWriter, request *http.Request)
	GetLineTypeByItemCode(writer http.ResponseWriter, request *http.Request)
	CampaignMaster(writer http.ResponseWriter, request *http.Request)
	ItemOprCodeWithPrice(writer http.ResponseWriter, request *http.Request)
	VehicleUnitMaster(writer http.ResponseWriter, request *http.Request)
	GetVehicleUnitByID(writer http.ResponseWriter, request *http.Request)
	GetVehicleUnitByChassisNumber(writer http.ResponseWriter, request *http.Request)
	CustomerByTypeAndAddress(writer http.ResponseWriter, request *http.Request)
	CustomerByTypeAndAddressByID(writer http.ResponseWriter, request *http.Request)
	CustomerByTypeAndAddressByCode(writer http.ResponseWriter, request *http.Request)
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

func (r *LookupControllerImpl) ItemOprCodeByCode(writer http.ResponseWriter, request *http.Request) {
	linetypeStrId := chi.URLParam(request, "linetype_id")
	linetypeId, err := strconv.Atoi(linetypeStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Line Type ID", http.StatusBadRequest)
		return
	}
	fmt.Println("linetypeId", linetypeId)

	itemCode := chi.URLParam(request, "item_code")
	if itemCode == "" {
		payloads.NewHandleError(writer, "Invalid Item Code", http.StatusBadRequest)
		return
	}
	fmt.Println("itemCode", itemCode)

	queryValues := request.URL.Query()
	queryParams := map[string]string{}
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	lookup, totalPages, totalRows, baseErr := r.LookupService.ItemOprCodeByCode(linetypeId, itemCode, paginate, criteria)
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

func (r *LookupControllerImpl) ItemOprCodeByID(writer http.ResponseWriter, request *http.Request) {
	linetypeStrId := chi.URLParam(request, "linetype_id")
	linetypeId, err := strconv.Atoi(linetypeStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Line Type ID", http.StatusBadRequest)
		return
	}

	itemStrId := chi.URLParam(request, "item_id")
	itemId, err := strconv.Atoi(itemStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Item ID", http.StatusBadRequest)
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
	lookup, totalPages, totalRows, baseErr := r.LookupService.ItemOprCodeByID(linetypeId, itemId, paginate, criteria)
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

func (r *LookupControllerImpl) CustomerByTypeAndAddressByID(writer http.ResponseWriter, request *http.Request) {
	customerStrId := chi.URLParam(request, "customer_id")
	customerId, err := strconv.Atoi(customerStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Customer ID", http.StatusBadRequest)
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
	lookup, totalPages, totalRows, baseErr := r.LookupService.CustomerByTypeAndAddressByID(customerId, paginate, criteria)
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

func (r *LookupControllerImpl) CustomerByTypeAndAddress(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{}
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	criteria := utils.BuildFilterCondition(queryParams)
	lookup, totalPages, totalRows, baseErr := r.LookupService.CustomerByTypeAndAddress(paginate, criteria)
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

func (r *LookupControllerImpl) CustomerByTypeAndAddressByCode(writer http.ResponseWriter, request *http.Request) {

	customerCodeStrId := chi.URLParam(request, "customer_code")

	queryValues := request.URL.Query()
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	queryParams := map[string]string{}
	criteria := utils.BuildFilterCondition(queryParams)

	lookup, totalPages, totalRows, baseErr := r.LookupService.CustomerByTypeAndAddressByCode(customerCodeStrId, paginate, criteria)

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

// GetLineTypeByItemCode godoc
// @Summary Get Line Type By Item Code
// @Description Get Line Type By Item Code
// @Tags Master
// @Accept json
// @Produce json
// @Param item_code path string true "Item Code"
// @Success 200 {object} ItemOprCodeResponse
// @Router /master/lookup/line-type/{item_code} [get]
func (r *LookupControllerImpl) GetLineTypeByItemCode(writer http.ResponseWriter, request *http.Request) {
	itemCode := chi.URLParam(request, "item_code")
	if itemCode == "" {
		payloads.NewHandleError(writer, "Invalid Item Code", http.StatusBadRequest)
		return
	}

	lookup, baseErr := r.LookupService.GetLineTypeByItemCode(itemCode)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, lookup, "Get Data Successfully", http.StatusOK)
}
