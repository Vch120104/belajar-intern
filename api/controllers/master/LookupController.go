package mastercontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type LookupController interface {
	ItemOprCode(writer http.ResponseWriter, request *http.Request)
	ItemOprCodeByCode(writer http.ResponseWriter, request *http.Request)
	ItemOprCodeByID(writer http.ResponseWriter, request *http.Request)
	GetLineTypeByItemCode(writer http.ResponseWriter, request *http.Request)
	GetCampaignMaster(writer http.ResponseWriter, request *http.Request)
	ItemOprCodeWithPrice(writer http.ResponseWriter, request *http.Request)
	VehicleUnitMaster(writer http.ResponseWriter, request *http.Request)
	GetVehicleUnitByID(writer http.ResponseWriter, request *http.Request)
	GetVehicleUnitByChassisNumber(writer http.ResponseWriter, request *http.Request)
	CustomerByTypeAndAddress(writer http.ResponseWriter, request *http.Request)
	CustomerByTypeAndAddressByID(writer http.ResponseWriter, request *http.Request)
	CustomerByTypeAndAddressByCode(writer http.ResponseWriter, request *http.Request)
	WorkOrderService(writer http.ResponseWriter, request *http.Request)
	ListItemLocation(writer http.ResponseWriter, request *http.Request)
	WarehouseGroupByCompany(writer http.ResponseWriter, request *http.Request)
	ItemListTrans(writer http.ResponseWriter, request *http.Request)
	ItemListTransPL(writer http.ResponseWriter, request *http.Request)
	ReferenceTypeWorkOrder(writer http.ResponseWriter, request *http.Request)
	ReferenceTypeWorkOrderByID(writer http.ResponseWriter, request *http.Request)
	ReferenceTypeSalesOrder(writer http.ResponseWriter, request *http.Request)
	ReferenceTypeSalesOrderByID(writer http.ResponseWriter, request *http.Request)
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
	queryParams := map[string]string{
		"opr_item_code": queryValues.Get("opr_item_code"),
		"opr_item_name": queryValues.Get("opr_item_name"),
	}
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
	itemCodeUnescaped, err := url.PathUnescape(itemCode)
	if err != nil {
		payloads.NewHandleError(writer, "Failed to decode Item Code", http.StatusBadRequest)
		return
	}

	if itemCodeUnescaped == "" {
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
	lookup, totalPages, totalRows, baseErr := r.LookupService.ItemOprCodeByCode(linetypeId, itemCodeUnescaped, paginate, criteria)
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

	billCodeStr := chi.URLParam(request, "transaction_type_id")
	billCodeStrId, err := strconv.Atoi(billCodeStr)
	if err != nil {
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
	lookup, totalPages, totalRows, baseErr := r.LookupService.GetVehicleUnitMaster(brandId, modelId, paginate, criteria)
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

func (r *LookupControllerImpl) GetCampaignMaster(writer http.ResponseWriter, request *http.Request) {
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
	lookup, totalPages, totalRows, baseErr := r.LookupService.GetCampaignMaster(companyId, paginate, criteria)
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

func (r *LookupControllerImpl) ListItemLocation(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	companyId, convErr := strconv.Atoi(queryValues.Get("company_id"))
	if convErr != nil {
		payloads.NewHandleError(writer, "company_id cannot be empty", http.StatusInternalServerError)
		return
	}

	queryParams := map[string]string{
		"warehouse_code":       queryValues.Get("warehouse_code"),
		"warehouse_name":       queryValues.Get("warehouse_name"),
		"warehouse_group_code": queryValues.Get("warehouse_group_code"),
		"warehouse_group_name": queryValues.Get("warehouse_group_name"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	warehouse, baseErr := r.LookupService.ListItemLocation(companyId, criteria, paginate)
	if baseErr != nil {
		exceptions.NewNotFoundException(writer, request, baseErr)
		return
	}
	payloads.NewHandleSuccessPagination(writer, warehouse.Rows, "Get Data Successfully!", http.StatusOK, warehouse.Limit, warehouse.Page, warehouse.TotalRows, warehouse.TotalPages)
}

func (r *LookupControllerImpl) WarehouseGroupByCompany(writer http.ResponseWriter, request *http.Request) {
	companyIdstr := chi.URLParam(request, "company_id")
	if companyIdstr == "" {
		payloads.NewHandleError(writer, "Invalid Company Id", http.StatusBadRequest)
	}

	companyId, _ := strconv.Atoi(companyIdstr)

	warehouse, baseErr := r.LookupService.WarehouseGroupByCompany(companyId)
	if baseErr != nil {
		exceptions.NewNotFoundException(writer, request, baseErr)
		return
	}
	payloads.NewHandleSuccess(writer, warehouse, "Get Data Successfully", http.StatusOK)
}

func (r *LookupControllerImpl) ItemListTrans(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"mtr_item.item_code":     queryValues.Get("item_code"),
		"mtr_item.item_name":     queryValues.Get("item_name"),
		"mtr_item.item_class_id": queryValues.Get("item_class_id"),
		"mic.item_class_code":    queryValues.Get("item_class_code"),
		"mic.item_class_name":    queryValues.Get("item_class_name"),
		"mit.item_type_code":     queryValues.Get("item_type_code"),
		"mil1.item_level_1_code": queryValues.Get("item_level_1_code"),
		"mil2.item_level_2_code": queryValues.Get("item_level_2_code"),
		"mil3.item_level_3_code": queryValues.Get("item_level_3_code"),
		"mil4.item_level_4_code": queryValues.Get("item_level_4_code"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	item, baseErr := r.LookupService.ItemListTrans(criteria, paginate)
	if baseErr != nil {
		helper.ReturnError(writer, request, baseErr)
		return
	}
	payloads.NewHandleSuccessPagination(writer, item.Rows, "Get Data Successfully!", http.StatusOK, item.Limit, item.Page, item.TotalRows, item.TotalPages)
}

func (r *LookupControllerImpl) ItemListTransPL(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	companyIdstr := queryValues.Get("company_id")
	if companyIdstr == "" {
		companyIdstr = "0"
	}

	companyId, _ := strconv.Atoi(companyIdstr)

	queryParams := map[string]string{
		"mid.brand_id":           queryValues.Get("brand_id"),
		"mtr_item.item_group_id": queryValues.Get("item_group_id"),
		"mtr_item.item_code":     queryValues.Get("item_code"),
		"mtr_item.item_name":     queryValues.Get("item_name"),
		"mic.item_class_code":    queryValues.Get("item_class_code"),
		"mit.item_type_code":     queryValues.Get("item_type_code"),
		"mil1.item_level_1_code": queryValues.Get("item_level_1_code"),
		"mil2.item_level_2_code": queryValues.Get("item_level_2_code"),
		"mil3.item_level_3_code": queryValues.Get("item_level_3_code"),
		"mil4.item_level_4_code": queryValues.Get("item_level_4_code"),
	}

	if queryParams["mid.brand_id"] == "" {
		payloads.NewHandleError(writer, "brand_id is required", http.StatusBadRequest)
		return
	}

	if queryParams["mtr_item.item_group_id"] == "" {
		payloads.NewHandleError(writer, "item_group_id is required", http.StatusBadRequest)
		return
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	item, baseErr := r.LookupService.ItemListTransPL(companyId, criteria, paginate)
	if baseErr != nil {
		item.Rows = []interface{}{}
		item.TotalRows = 0
		item.TotalPages = 0
	}
	payloads.NewHandleSuccessPagination(writer, item.Rows, "Get Data Successfully!", http.StatusOK, item.Limit, item.Page, item.TotalRows, item.TotalPages)
}

func (r *LookupControllerImpl) ReferenceTypeWorkOrder(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"work_order_system_number":      queryValues.Get("work_order_system_number"),
		"work_order_document_number":    queryValues.Get("work_order_document_number"),
		"work_order_date":               queryValues.Get("work_order_date"),
		"work_order_status_description": queryValues.Get("work_order_status_description"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	var filters []utils.FilterCondition
	for key, value := range queryParams {
		if value != "" {
			filters = append(filters, utils.FilterCondition{
				ColumnField: key,
				ColumnValue: value,
			})
		}
	}

	referenceType, totalPages, totalRows, baseErr := r.LookupService.ReferenceTypeWorkOrder(paginate, filters)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(writer, referenceType, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *LookupControllerImpl) ReferenceTypeWorkOrderByID(writer http.ResponseWriter, request *http.Request) {

	referenceTypeIdStr := chi.URLParam(request, "work_order_system_number")
	referenceTypeId, err := strconv.Atoi(referenceTypeIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Reference Type ID", http.StatusBadRequest)
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
	referenceType, totalPages, totalRows, baseErr := r.LookupService.ReferenceTypeWorkOrderByID(referenceTypeId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(writer, referenceType, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *LookupControllerImpl) ReferenceTypeSalesOrder(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"work_order_system_number":      queryValues.Get("work_order_system_number"),
		"work_order_document_number":    queryValues.Get("work_order_document_number"),
		"work_order_date":               queryValues.Get("work_order_date"),
		"work_order_status_description": queryValues.Get("work_order_status_description"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	var filters []utils.FilterCondition
	for key, value := range queryParams {
		if value != "" {
			filters = append(filters, utils.FilterCondition{
				ColumnField: key,
				ColumnValue: value,
			})
		}
	}

	referenceType, totalPages, totalRows, baseErr := r.LookupService.ReferenceTypeSalesOrder(paginate, filters)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(writer, referenceType, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *LookupControllerImpl) ReferenceTypeSalesOrderByID(writer http.ResponseWriter, request *http.Request) {

	referenceTypeIdStr := chi.URLParam(request, "sales_order_system_number")
	referenceTypeId, err := strconv.Atoi(referenceTypeIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Reference Type ID", http.StatusBadRequest)
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
	referenceType, totalPages, totalRows, baseErr := r.LookupService.ReferenceTypeSalesOrderByID(referenceTypeId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(writer, referenceType, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}
