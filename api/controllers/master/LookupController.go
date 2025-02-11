package mastercontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type LookupController interface {
	GetOprItemPrice(writer http.ResponseWriter, request *http.Request)
	ItemOprCode(writer http.ResponseWriter, request *http.Request)
	ItemOprCodeByCode(writer http.ResponseWriter, request *http.Request)
	ItemOprCodeByID(writer http.ResponseWriter, request *http.Request)
	GetLineTypeByItemCode(writer http.ResponseWriter, request *http.Request)
	GetLineTypeByReferenceType(writer http.ResponseWriter, request *http.Request)
	GetCampaignMaster(writer http.ResponseWriter, request *http.Request)
	ItemOprCodeWithPrice(writer http.ResponseWriter, request *http.Request)
	ItemOprCodeWithPriceByID(writer http.ResponseWriter, request *http.Request)
	ItemOprCodeWithPriceByCode(writer http.ResponseWriter, request *http.Request)
	VehicleUnitMaster(writer http.ResponseWriter, request *http.Request)
	GetVehicleUnitByID(writer http.ResponseWriter, request *http.Request)
	GetVehicleUnitByChassisNumber(writer http.ResponseWriter, request *http.Request)
	CustomerByTypeAndAddress(writer http.ResponseWriter, request *http.Request)
	CustomerByTypeAndAddressByID(writer http.ResponseWriter, request *http.Request)
	CustomerByTypeAndAddressByCode(writer http.ResponseWriter, request *http.Request)
	WorkOrderService(writer http.ResponseWriter, request *http.Request)
	WorkOrderAtpmRegistration(writer http.ResponseWriter, request *http.Request)
	ListItemLocation(writer http.ResponseWriter, request *http.Request)
	WarehouseGroupByCompany(writer http.ResponseWriter, request *http.Request)
	ItemListTrans(writer http.ResponseWriter, request *http.Request)
	ItemListTransPL(writer http.ResponseWriter, request *http.Request)
	ReferenceTypeWorkOrder(writer http.ResponseWriter, request *http.Request)
	ReferenceTypeWorkOrderByID(writer http.ResponseWriter, request *http.Request)
	ReferenceTypeSalesOrder(writer http.ResponseWriter, request *http.Request)
	ReferenceTypeSalesOrderByID(writer http.ResponseWriter, request *http.Request)
	LocationAvailable(writer http.ResponseWriter, request *http.Request)
	ItemDetailForItemInquiry(writer http.ResponseWriter, request *http.Request)
	ItemSubstituteDetailForItemInquiry(writer http.ResponseWriter, request *http.Request)
	GetPartNumberItemImport(writer http.ResponseWriter, request *http.Request)
	LocationItem(writer http.ResponseWriter, request *http.Request)
	ItemLocUOM(writer http.ResponseWriter, request *http.Request)
	ItemLocUOMById(writer http.ResponseWriter, request *http.Request)
	ItemLocUOMByCode(writer http.ResponseWriter, request *http.Request)
	ItemMasterForFreeAccs(writer http.ResponseWriter, request *http.Request)
	ItemMasterForFreeAccsById(writer http.ResponseWriter, request *http.Request)
	ItemMasterForFreeAccsByCode(writer http.ResponseWriter, request *http.Request)
}

type LookupControllerImpl struct {
	LookupService masterservice.LookupService
}

func NewLookupController(LookupService masterservice.LookupService) LookupController {
	return &LookupControllerImpl{
		LookupService: LookupService,
	}
}

// @Summary Get Opr Item Price
// @Description Get Opr Item Price
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param line_type_id path int true "Line Type ID"
// @Param package_id query string false "Package ID"
// @Param package_code query string false "Package Code"
// @Param package_name query string false "Package Name"
// @Param profit_center_name query string false "Profit Center Name"
// @Param model_code query string false "Model Code"
// @Param model_description query string false "Model Description"
// @Param package_price query string false "Package Price"
// @Param brand_id query string false "Brand ID"
// @Param model_id query string false "Model ID"
// @Param variant_id query string false "Variant ID"
// @Param operation_id query string false "Operation ID"
// @Param operation_code query string false "Operation Code"
// @Param operation_name query string false "Operation Name"
// @Param operation_entries_code query string false "Operation Entries Code"
// @Param operation_entries_description query string false "Operation Entries Description"
// @Param operation_key_code query string false "Operation Key Code"
// @Param operation_key_description query string false "Operation Key Description"
// @Param item_id query string false "Item ID"
// @Param item_code query string false "Item Code"
// @Param item_name query string false "Item Name"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-opr-code/{line_type_id} [get]
func (r *LookupControllerImpl) ItemOprCode(writer http.ResponseWriter, request *http.Request) {
	linetypeIdStr := chi.URLParam(request, "line_type_id")
	linetypeId, err := strconv.Atoi(linetypeIdStr)
	if err != nil || linetypeId == 0 {
		payloads.NewHandleError(writer, "Invalid Line Type Id", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"package_id":                    queryValues.Get("package_id"),
		"package_code":                  queryValues.Get("package_code"),
		"package_name":                  queryValues.Get("package_name"),
		"profit_center_name":            queryValues.Get("profit_center_name"),
		"model_code":                    queryValues.Get("model_code"),
		"model_description":             queryValues.Get("model_description"),
		"package_price":                 queryValues.Get("package_price"),
		"brand_id":                      queryValues.Get("brand_id"),
		"model_id":                      queryValues.Get("model_id"),
		"variant_id":                    queryValues.Get("variant_id"),
		"operation_id":                  queryValues.Get("operation_id"),
		"operation_code":                queryValues.Get("operation_code"),
		"operation_name":                queryValues.Get("operation_name"),
		"operation_entries_code":        queryValues.Get("operation_entries_code"),
		"operation_entries_description": queryValues.Get("operation_entries_description"),
		"operation_key_code":            queryValues.Get("operation_key_code"),
		"operation_key_description":     queryValues.Get("operation_key_description"),
		"item_id":                       queryValues.Get("item_id"),
		"item_code":                     queryValues.Get("item_code"),
		"item_name":                     queryValues.Get("item_name"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	criteria := utils.BuildFilterCondition(queryParams)

	lookup, baseErr := r.LookupService.ItemOprCode(linetypeId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		lookup.Limit,
		lookup.Page,
		int64(lookup.TotalRows),
		lookup.TotalPages,
	)
}

// @Summary Get Opr Item Price By Code
// @Description Get Opr Item Price By Code
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param line_type_id path int true "Line Type ID"
// @Param opr_item_code path string true "Item Code"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-opr-code/{line_type_id}/{opr_item_code} [get]
func (r *LookupControllerImpl) ItemOprCodeByCode(writer http.ResponseWriter, request *http.Request) {
	linetypeIdStr := chi.URLParam(request, "line_type_id")
	linetypeId, err := strconv.Atoi(linetypeIdStr)
	if err != nil || linetypeId == 0 {
		payloads.NewHandleError(writer, "Invalid Line Type Id", http.StatusBadRequest)
		return
	}

	encodedCampaignCode := chi.URLParam(request, "*")
	if len(encodedCampaignCode) > 0 && encodedCampaignCode[0] == '/' {
		encodedCampaignCode = encodedCampaignCode[1:]
	}

	itemCodeUnescaped, err := url.PathUnescape(encodedCampaignCode)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Campaign Code",
		})
		return
	}

	itemCodeUnescaped, err = url.PathUnescape(itemCodeUnescaped)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Campaign Code after second decoding",
		})
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

	lookup, baseErr := r.LookupService.ItemOprCodeByCode(linetypeId, itemCodeUnescaped, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
	)
}

// @Summary Get Opr Item Price By ID
// @Description Get Opr Item Price By ID
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param line_type_id path int true "Line Type ID"
// @Param id path int true "Item ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-opr-code/{line_type_id}/{id} [get]
func (r *LookupControllerImpl) ItemOprCodeByID(writer http.ResponseWriter, request *http.Request) {
	linetypeIdStr := chi.URLParam(request, "line_type_id")
	linetypeId, err := strconv.Atoi(linetypeIdStr)
	if err != nil || linetypeId == 0 {
		payloads.NewHandleError(writer, "Invalid Line Type Id", http.StatusBadRequest)
		return
	}

	itemStrId := chi.URLParam(request, "id")
	itemId, err := strconv.Atoi(itemStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Operation Item ID", http.StatusBadRequest)
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
	lookup, baseErr := r.LookupService.ItemOprCodeByID(linetypeId, itemId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
	)
}

// @Summary Get Opr Item Price With Price
// @Description Get Opr Item Price With Price
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param line_type_id path int true "Line Type ID"
// @Param company_id query int true "Company ID"
// @Param opr_item_id query int true "Opr Item ID"
// @Param brand_id query int true "Brand ID"
// @Param model_id query int true "Model ID"
// @Param trx_type_id query int true "Trx Type ID"
// @Param job_type_id query int true "Job Type ID"
// @Param variant_id query int true "Variant ID"
// @Param currency_id query int true "Currency ID"
// @Param whs_group query string true "Warehouse Group"
// @Param package_id query string false "Package ID"
// @Param package_code query string false "Package Code"
// @Param package_name query string false "Package Name"
// @Param profit_center_id query string false "Profit Center ID"
// @Param brand_id query string false "Brand ID"
// @Param model_id query string false "Model ID"
// @Param package_price query string false "Package Price"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-opr-code-with-price [get]
func (r *LookupControllerImpl) ItemOprCodeWithPrice(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	linetypeId, err := strconv.Atoi(queryValues.Get("line_type_id"))
	if err != nil || linetypeId == 0 {
		payloads.NewHandleError(writer, "Invalid Line Type Id", http.StatusBadRequest)
		return
	}

	companyId, err := strconv.Atoi(queryValues.Get("company_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Company Id", http.StatusBadRequest)
		return
	}

	oprItemCode, err := strconv.Atoi(queryValues.Get("opr_item_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Opr Item Id", http.StatusBadRequest)
		return
	}

	brandId, err := strconv.Atoi(queryValues.Get("brand_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Brand Id", http.StatusBadRequest)
		return
	}

	modelId, err := strconv.Atoi(queryValues.Get("model_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Model Id", http.StatusBadRequest)
		return
	}

	trxTypeId, err := strconv.Atoi(queryValues.Get("trx_type_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Trx Type Id", http.StatusBadRequest)
		return
	}

	jobTypeId, err := strconv.Atoi(queryValues.Get("job_type_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Job Type Id", http.StatusBadRequest)
		return
	}

	variantId, err := strconv.Atoi(queryValues.Get("variant_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Variant Id", http.StatusBadRequest)
		return
	}

	currencyId, err := strconv.Atoi(queryValues.Get("currency_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Currency Id", http.StatusBadRequest)
		return
	}

	whsGroup := queryValues.Get("whs_group")
	if whsGroup == "" {
		payloads.NewHandleError(writer, "Invalid Warehouse Group", http.StatusBadRequest)
		return
	}

	var criteria []utils.FilterCondition
	if linetypeId == 1 {
		queryParams := map[string]string{
			"mtr_package.package_id":       queryValues.Get("package_id"),
			"mtr_package.package_code":     queryValues.Get("package_code"),
			"mtr_package.package_name":     queryValues.Get("package_name"),
			"mtr_package.profit_center_id": queryValues.Get("profit_center_id"),
			"mtr_package.brand_id":         queryValues.Get("brand_id"),
			"mtr_package.model_id":         queryValues.Get("model_id"),
			"mtr_package.package_price":    queryValues.Get("package_price"),
		}

		criteria = utils.BuildFilterCondition(queryParams)
	} else if linetypeId == 2 {
		queryParams := map[string]string{
			"mtr_operation_model_mapping.operation_id": queryValues.Get("operation_id"),
		}
		criteria = utils.BuildFilterCondition(queryParams)
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	lookup, baseErr := r.LookupService.ItemOprCodeWithPrice(linetypeId, companyId, oprItemCode, brandId, modelId, trxTypeId, jobTypeId, variantId, currencyId, whsGroup, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		lookup.Limit,
		lookup.Page,
		int64(lookup.TotalRows),
		lookup.TotalPages,
	)
}

// @Summary Get Vehicle Unit Master
// @Description Get Vehicle Unit Master
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param brand_id path int true "Brand ID"
// @Param model_id path int true "Model ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/vehicle-unit-master/{brand_id}/{model_id} [get]
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
	lookup, baseErr := r.LookupService.GetVehicleUnitMaster(brandId, modelId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		lookup.Limit,
		lookup.Page,
		int64(lookup.TotalRows),
		lookup.TotalPages,
	)
}

// @Summary Get Vehicle Unit By ID
// @Description Get Vehicle Unit By ID
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param vehicle_id path int true "Vehicle ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/vehicle-unit-master/{vehicle_id} [get]
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
	lookup, baseErr := r.LookupService.GetVehicleUnitByID(vehicleId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		lookup.Limit,
		lookup.Page,
		int64(lookup.TotalRows),
		lookup.TotalPages,
	)
}

// @Summary Get Vehicle Unit By Chassis Number
// @Description Get Vehicle Unit By Chassis Number
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param vehicle_chassis_number path string true "Vehicle Chassis Number"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/vehicle-unit-master/by-code/{vehicle_chassis_number} [get]
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
	lookup, baseErr := r.LookupService.GetVehicleUnitByChassisNumber(chassisNumber, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		lookup.Limit,
		lookup.Page,
		int64(lookup.TotalRows),
		lookup.TotalPages,
	)
}

// @Summary Get Campaign Master
// @Description Get Campaign Master
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param company_id path int true "Company ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/campaign-master/{company_id} [get]
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
	lookup, baseErr := r.LookupService.GetCampaignMaster(companyId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		lookup.Limit,
		lookup.Page,
		int64(lookup.TotalRows),
		lookup.TotalPages,
	)
}

// @Summary Work Order Service
// @Description Work Order Service
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/work-order-service [get]
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
	lookup, baseErr := r.LookupService.WorkOrderService(paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		lookup.Limit,
		lookup.Page,
		int64(lookup.TotalRows),
		lookup.TotalPages,
	)
}

// @Summary Work Order ATPM Registration
// @Description Work Order ATPM Registration
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/work-order-atpm-registration [get]
func (r *LookupControllerImpl) WorkOrderAtpmRegistration(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()
	queryParams := map[string]string{}
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	lookup, baseErr := r.LookupService.WorkOrderAtpmRegistration(paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}

		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		lookup.Rows,
		"Get Data Successfully",
		http.StatusOK,
		lookup.Limit,
		lookup.Page,
		int64(lookup.TotalRows),
		lookup.TotalPages,
	)
}

// @Summary Customer By Type And Address By ID
// @Description Customer By Type And Address By ID
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param customer_id path int true "Customer ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/new-bill-to/{customer_id} [get]
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
	lookup, baseErr := r.LookupService.CustomerByTypeAndAddressByID(customerId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		lookup.Limit,
		lookup.Page,
		int64(lookup.TotalRows),
		lookup.TotalPages,
	)
}

// @Summary Customer By Type And Address
// @Description Customer By Type And Address
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/new-bill-to [get]
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
	lookup, baseErr := r.LookupService.CustomerByTypeAndAddress(paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		lookup.Limit,
		lookup.Page,
		int64(lookup.TotalRows),
		lookup.TotalPages,
	)
}

// @Summary Customer By Type And Address By Code
// @Description Customer By Type And Address By Code
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param customer_code path string true "Customer Code"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/new-bill-to/{customer_code} [get]
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

	lookup, baseErr := r.LookupService.CustomerByTypeAndAddressByCode(customerCodeStrId, paginate, criteria)

	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		lookup.Limit,
		lookup.Page,
		int64(lookup.TotalRows),
		lookup.TotalPages,
	)
}

// GetLineTypeByItemCode godoc
// @Summary Get Line Type By Item Code
// @Description Get Line Type By Item Code
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param item_code path string true "Item Code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/line-type/{item_code} [get]
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

// @Summary List Item Location
// @Description List Item Location
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param company_id query int true "Company ID"
// @Param warehouse_code query string false "Warehouse Code"
// @Param warehouse_name query string false "Warehouse Name"
// @Param warehouse_group_code query string false "Warehouse Group Code"
// @Param warehouse_group_name query string false "Warehouse Group Name"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-location-warehouse [get]
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

// @Summary Warehouse Group By Company
// @Description Warehouse Group By Company
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param company_id path int true "Company ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/warehouse-group/{company_id [get]
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

// @Summary Item List Trans
// @Description Item List Trans
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param item_code query string false "Item Code"
// @Param item_name query string false "Item Name"
// @Param item_class_id query string false "Item Class ID"
// @Param item_class_code query string false "Item Class Code"
// @Param item_class_name query string false "Item Class Name"
// @Param item_type_code query string false "Item Type Code"
// @Param item_level_1_code query string false "Item Level 1 Code"
// @Param item_level_2_code query string false "Item Level 2 Code"
// @Param item_level_3_code query string false "Item Level 3 Code"
// @Param item_level_4_code query string false "Item Level 4 Code"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-list-trans [get]
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

// @Summary Item List Trans PL
// @Description Item List Trans PL
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param company_id query int true "Company ID"
// @Param brand_id query string false "Brand ID"
// @Param item_group_id query string false "Item Group ID"
// @Param item_code query string false "Item Code"
// @Param item_name query string false "Item Name"
// @Param item_class_code query string false "Item Class Code"
// @Param item_type_code query string false "Item Type Code"
// @Param item_level_1_code query string false "Item Level 1 Code"
// @Param item_level_2_code query string false "Item Level 2 Code"
// @Param item_level_3_code query string false "Item Level 3 Code"
// @Param item_level_4_code query string false "Item Level 4 Code"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-list-trans-pl [get]
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

// @Summary Reference Type Work Order
// @Description Reference Type Work Order
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param work_order_system_number query string false "Work Order System Number"
// @Param work_order_document_number query string false "Work Order Document Number"
// @Param work_order_date query string false "Work Order Date"
// @Param work_order_status_description query string false "Work Order Status Description"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/reference-type-work-order [get]
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

	lookup, baseErr := r.LookupService.ReferenceTypeWorkOrder(paginate, filters)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		lookup.Limit,
		lookup.Page,
		int64(lookup.TotalRows),
		lookup.TotalPages,
	)
}

// @Summary Reference Type Work Order By ID
// @Description Reference Type Work Order By ID
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param work_order_system_number path int true "Work Order System Number"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/reference-type-work-order/{work_order_system_number} [get]
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
	lookup, baseErr := r.LookupService.ReferenceTypeWorkOrderByID(referenceTypeId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		lookup.Limit,
		lookup.Page,
		int64(lookup.TotalRows),
		lookup.TotalPages,
	)
}

// @Summary Reference Type Sales Order
// @Description Reference Type Sales Order
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param work_order_system_number query string false "Work Order System Number"
// @Param work_order_document_number query string false "Work Order Document Number"
// @Param work_order_date query string false "Work Order Date"
// @Param work_order_status_description query string false "Work Order Status Description"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/reference-type-sales-order [get]
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

	lookup, baseErr := r.LookupService.ReferenceTypeSalesOrder(paginate, filters)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		lookup.Limit,
		lookup.Page,
		int64(lookup.TotalRows),
		lookup.TotalPages,
	)
}

// @Summary Reference Type Sales Order By ID
// @Description Reference Type Sales Order By ID
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param sales_order_system_number path int true "Sales Order System Number"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/reference-type-sales-order/{sales_order_system_number} [get]
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
	lookup, baseErr := r.LookupService.ReferenceTypeSalesOrderByID(referenceTypeId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		lookup.Limit,
		lookup.Page,
		int64(lookup.TotalRows),
		lookup.TotalPages,
	)
}

// @Summary Get Line Type By Reference Type
// @Description Get Line Type By Reference Type
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param reference_type_id path int true "Reference Type ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/line-type-reference/{reference_type_id} [get]
func (r *LookupControllerImpl) GetLineTypeByReferenceType(writer http.ResponseWriter, request *http.Request) {
	referenceTypeStr := chi.URLParam(request, "reference_type_id")
	referenceType, err := strconv.Atoi(referenceTypeStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Reference Type", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()
	pages := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	lookup, baseErr := r.LookupService.GetLineTypeByReferenceType(referenceType, pages)
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

// @Summary Location Available
// @Description Location Available
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param company_id query int true "Company ID"
// @Param warehouse_id query int true "Warehouse ID"
// @Param is_active query string false "Is Active"
// @Param warehouse_location_code query string false "Warehouse Location Code"
// @Param warehouse_location_name query string false "Warehouse Location Name"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/location-available [get]
func (r *LookupControllerImpl) LocationAvailable(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"company_id":                       queryValues.Get("company_id"),
		"warehouse_id":                     queryValues.Get("warehouse_id"),
		"mtr_warehouse_location.is_active": queryValues.Get("is_active"),
		"mtr_warehouse_location.warehouse_location_code": queryValues.Get("warehouse_location_code"),
		"mtr_warehouse_location.warehouse_location_name": queryValues.Get("warehouse_location_name"),
	}

	if queryParams["company_id"] == "" {
		payloads.NewHandleError(writer, "company_id cannot be empty", http.StatusBadRequest)
		return
	}
	if queryParams["warehouse_id"] == "" {
		payloads.NewHandleError(writer, "warehouse_id cannot be empty", http.StatusBadRequest)
		return
	}

	criteria := utils.BuildFilterCondition(queryParams)

	pages := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	location, baseErr := r.LookupService.LocationAvailable(criteria, pages)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(writer, location.Rows, "Get Data Successfully!", http.StatusOK, location.Limit, location.Page, location.TotalRows, location.TotalPages)
}

// @Summary Item Detail For Item Inquiry
// @Description Item Detail For Item Inquiry
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param brand_id query string true "Brand ID"
// @Param item_id query string true "Item ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-detail/item-inquiry [get]
func (r *LookupControllerImpl) ItemDetailForItemInquiry(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"brand_id": queryValues.Get("brand_id"),
		"item_id":  queryValues.Get("item_id"),
	}

	if queryParams["brand_id"] == "" {
		payloads.NewHandleError(writer, "brand_id cannot be empty", http.StatusBadRequest)
		return
	}
	if queryParams["item_id"] == "" {
		payloads.NewHandleError(writer, "item_id cannot be empty", http.StatusBadRequest)
		return
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	item, baseErr := r.LookupService.ItemDetailForItemInquiry(criteria, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}
	payloads.NewHandleSuccessPagination(writer, item.Rows, "Get Data Successfully!", http.StatusOK, item.Limit, item.Page, item.TotalRows, item.TotalPages)
}

// @Summary Item Substitute Detail For Item Inquiry
// @Description Item Substitute Detail For Item Inquiry
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param item_id query string true "Item ID"
// @Param company_id query string true "Company ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-substitute/detail/item-inquiry [get]
func (r *LookupControllerImpl) ItemSubstituteDetailForItemInquiry(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"item_id":    queryValues.Get("item_id"),
		"company_id": queryValues.Get("company_id"),
	}

	if queryParams["item_id"] == "" {
		payloads.NewHandleError(writer, "item_id cannot be empty", http.StatusBadRequest)
		return
	}
	if queryParams["company_id"] == "" {
		payloads.NewHandleError(writer, "company_id cannot be empty", http.StatusBadRequest)
		return
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	item, baseErr := r.LookupService.ItemSubstituteDetailForItemInquiry(criteria, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}
	payloads.NewHandleSuccessPagination(writer, item.Rows, "Get Data Successfully!", http.StatusOK, item.Limit, item.Page, item.TotalRows, item.TotalPages)
}

// @Summary Get PartNumber Item Import
// @Description Get PartNumber Item Import
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param item_code query string false "Item Code"
// @Param item_name query string false "Item Name"
// @Param supplier_code query string false "Supplier Code"
// @Param supplier_name query string false "Supplier Name"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-import/part-number [get]
func (r *LookupControllerImpl) GetPartNumberItemImport(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	internalFilterCondition := map[string]string{
		"mtr_item.item_code": queryValues.Get("item_code"),
		"mtr_item.item_name": queryValues.Get("item_name"),
	}
	externalFilterCondition := map[string]string{
		"supplier_code": queryValues.Get("supplier_code"),
		"supplier_name": queryValues.Get("supplier_name"),
	}

	internalCriteria := utils.BuildFilterCondition(internalFilterCondition)
	externalCriteria := utils.BuildFilterCondition(externalFilterCondition)

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	item, baseErr := r.LookupService.GetPartNumberItemImport(internalCriteria, externalCriteria, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}
	payloads.NewHandleSuccessPagination(writer, item.Rows, "Get Data Successfully!", http.StatusOK, item.Limit, item.Page, item.TotalRows, item.TotalPages)
}

// @Summary Location Item
// @Description Location Item
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param warehouse_location_id query string false "Warehouse Location ID"
// @Param warehouse_location_code query string false "Warehouse Location Code"
// @Param warehouse_location_name query string false "Warehouse Location Name"
// @Param company_id query string false "Company ID"
// @Param warehouse_id query string false "Warehouse ID"
// @Param item_id query string false "Item ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/location-item [get]
func (r *LookupControllerImpl) LocationItem(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	filterCondition := map[string]string{
		"mwl.warehouse_location_id":   queryValues.Get("warehouse_location_id"),
		"mwl.warehouse_location_code": queryValues.Get("warehouse_location_code"),
		"mwl.warehouse_location_name": queryValues.Get("warehouse_location_name"),
		"mwm.company_id":              queryValues.Get("company_id"),
		"mwm.warehouse_id":            queryValues.Get("warehouse_id"),
		"mtr_location_item.item_id":   queryValues.Get("item_id"),
	}

	criteria := utils.BuildFilterCondition(filterCondition)

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	if paginate.GetSortBy() == "" {
		paginate.SortBy = "mwl.warehouse_location_id"
	}

	item, baseErr := r.LookupService.LocationItem(criteria, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}
	payloads.NewHandleSuccessPagination(writer, item.Rows, "Get Data Successfully!", http.StatusOK, item.Limit, item.Page, item.TotalRows, item.TotalPages)

}

// @Summary Item Loc UOM
// @Description Item Loc UOM
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param company_id query string false "Company ID"
// @Param item_id query string false "Item ID"
// @Param item_code query string false "Item Code"
// @Param item_name query string false "Item Name"
// @Param uom_code query string false "UOM Code"
// @Param quantity_available query string false "Quantity Available"
// @Param is_active query string false "Is Active"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-loc-uom [get]
func (r *LookupControllerImpl) ItemLocUOM(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	filterCondition := map[string]string{
		"company_id":         queryValues.Get("company_id"),
		"mi.item_id":         queryValues.Get("item_id"),
		"mi.item_code":       queryValues.Get("item_code"),
		"mi.item_name":       queryValues.Get("item_name"),
		"mu.uom_code":        queryValues.Get("uom_code"),
		"quantity_available": queryValues.Get("quantity_available"),
		"mi.is_active":       queryValues.Get("is_active"),
	}

	if filterCondition["company_id"] == "" {
		payloads.NewHandleError(writer, "company_id cannot be empty", http.StatusBadRequest)
		return
	}

	criteria := utils.BuildFilterCondition(filterCondition)

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	if paginate.GetSortBy() == "" {
		paginate.SortBy = "item_id"
	}

	item, baseErr := r.LookupService.ItemLocUOM(criteria, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}
	payloads.NewHandleSuccessPagination(writer, item.Rows, "Get Data Successfully!", http.StatusOK, item.Limit, item.Page, item.TotalRows, item.TotalPages)
}

// @Summary Item Loc UOM By ID
// @Description Item Loc UOM By ID
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param company_id path int true "Company ID"
// @Param item_id path int true "Item ID"
// @Param item_code query string false "Item Code"
// @Param item_name query string false "Item Name"
// @Param uom_code query string false "UOM Code"
// @Param quantity_available query string false "Quantity Available"
// @Param is_active query string false "Is Active"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-loc-uom/by-id/{company_id}/{item_id} [get]
func (r *LookupControllerImpl) ItemLocUOMById(writer http.ResponseWriter, request *http.Request) {
	companyId, _ := strconv.Atoi(chi.URLParam(request, "company_id"))
	itemId, _ := strconv.Atoi(chi.URLParam(request, "item_id"))

	item, baseErr := r.LookupService.ItemLocUOMById(companyId, itemId)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}
	payloads.NewHandleSuccess(writer, item, "Get Data Successfully!", http.StatusOK)
}

// @Summary Item Loc UOM By Code
// @Description Item Loc UOM By Code
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param company_id path int true "Company ID"
// @Param item_code query string true "Item Code"
// @Param item_name query string false "Item Name"
// @Param uom_code query string false "UOM Code"
// @Param quantity_available query string false "Quantity Available"
// @Param is_active query string false "Is Active"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-loc-uom/by-code/{company_id} [get]
func (r *LookupControllerImpl) ItemLocUOMByCode(writer http.ResponseWriter, request *http.Request) {
	companyId, _ := strconv.Atoi(chi.URLParam(request, "company_id"))
	itemCode := request.URL.Query().Get("item_code")

	item, baseErr := r.LookupService.ItemLocUOMByCode(companyId, itemCode)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}
	payloads.NewHandleSuccess(writer, item, "Get Data Successfully!", http.StatusOK)
}

// @Summary Item Opr Code With Price By ID
// @Description Item Opr Code With Price By ID
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param line_type_id path int true "Line Type ID"
// @Param opr_item_id path int true "Opr Item ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-opr-code-with-price/{line_type_id}/by-id/{opr_item_id} [get]
func (r *LookupControllerImpl) ItemOprCodeWithPriceByID(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	linetypeIdStr := chi.URLParam(request, "line_type_id")
	oprItemIdStr := chi.URLParam(request, "opr_item_id")

	if linetypeIdStr == "" {
		payloads.NewHandleError(writer, "Missing Line Type Id in URL", http.StatusBadRequest)
		return
	}
	if oprItemIdStr == "" {
		payloads.NewHandleError(writer, "Missing Opr Item Id in URL", http.StatusBadRequest)
		return
	}

	linetypeId, err := strconv.Atoi(linetypeIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Line Type Id", http.StatusBadRequest)
		return
	}

	oprItemId, err := strconv.Atoi(oprItemIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Opr Item Id", http.StatusBadRequest)
		return
	}

	queryParams := map[string]string{}
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	criteria := utils.BuildFilterCondition(queryParams)
	lookup, baseErr := r.LookupService.ItemOprCodeWithPriceByID(linetypeId, oprItemId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
	)
}

// @Summary Item Opr Code With Price By Code
// @Description Item Opr Code With Price By Code
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param line_type_id path int true "Line Type ID"
// @Param opr_item_code query string true "Opr Item Code"
// @Param opr_item_name query string false "Opr Item Name"
// @Param brand_id query int false "Brand ID"
// @Param model_id query int false "Model ID"
// @Param job_type_id query int false "Job Type ID"
// @Param variant_id query int false "Variant ID"
// @Param currency_id query int false "Currency ID"
// @Param trx_type_id query int false "Bill Code"
// @Param whs_group query string false "Warehouse Group"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-opr-code-with-price/{line_type_id}/by-code/{opr_item_code} [get]
func (r *LookupControllerImpl) ItemOprCodeWithPriceByCode(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	linetypeIdStr := chi.URLParam(request, "line_type_id")
	encodedoprItemCodeCode := chi.URLParam(request, "*")
	if len(encodedoprItemCodeCode) > 0 && encodedoprItemCodeCode[0] == '/' {
		encodedoprItemCodeCode = encodedoprItemCodeCode[1:]
	}

	itemCodeUnescaped, err := url.PathUnescape(encodedoprItemCodeCode)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid operation item Code",
		})
		return
	}

	itemCodeUnescaped, err = url.PathUnescape(itemCodeUnescaped)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid operation item Code after second decoding",
		})
		return
	}

	if linetypeIdStr == "" {
		payloads.NewHandleError(writer, "Missing Line Type Id in URL", http.StatusBadRequest)
		return
	}

	linetypeId, err := strconv.Atoi(linetypeIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Line Type Id", http.StatusBadRequest)
		return
	}

	queryParams := map[string]string{}
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	criteria := utils.BuildFilterCondition(queryParams)
	lookup, baseErr := r.LookupService.ItemOprCodeWithPriceByCode(linetypeId, itemCodeUnescaped, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(
		writer,
		lookup.Rows,
		"Get Data Successfully!",
		http.StatusOK,
	)
}

// @Summary Get Opr Item Price
// @Description Get Opr Item Price
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param line_type_id query int true "Line Type ID"
// @Param company_id query int true "Company ID"
// @Param opr_item_id query int true "Opr Item ID"
// @Param brand_id query int true "Brand ID"
// @Param model_id query int true "Model ID"
// @Param job_type_id query int true "Job Type ID"
// @Param variant_id query int true "Variant ID"
// @Param currency_id query int true "Currency ID"
// @Param trx_type_id query int true "Bill Code"
// @Param whs_group query string true "Warehouse Group"
// @Param limit query int false "Limit"
func (r *LookupControllerImpl) GetOprItemPrice(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	linetypeId, err := strconv.Atoi(queryValues.Get("line_type_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Line Type Id", http.StatusBadRequest)
		return
	}

	companyId, err := strconv.Atoi(queryValues.Get("company_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Company Id", http.StatusBadRequest)
		return
	}

	oprItemCode, err := strconv.Atoi(queryValues.Get("opr_item_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Opr Item Id", http.StatusBadRequest)
		return
	}

	brandId, err := strconv.Atoi(queryValues.Get("brand_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Brand Id", http.StatusBadRequest)
		return
	}

	modelId, err := strconv.Atoi(queryValues.Get("model_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Model Id", http.StatusBadRequest)
		return
	}

	jobTypeId, err := strconv.Atoi(queryValues.Get("job_type_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Job Type Id", http.StatusBadRequest)
		return
	}

	variantId, err := strconv.Atoi(queryValues.Get("variant_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Variant Id", http.StatusBadRequest)
		return
	}

	currencyId, err := strconv.Atoi(queryValues.Get("currency_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Currency Id", http.StatusBadRequest)
		return
	}

	billCode, err := strconv.Atoi(queryValues.Get("trx_type_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Bill Code", http.StatusBadRequest)
		return
	}

	whsGroup := queryValues.Get("whs_group")
	if whsGroup == "" {
		payloads.NewHandleError(writer, "Warehouse Group is required", http.StatusBadRequest)
		return
	}

	price, baseErr := r.LookupService.GetOprItemPrice(
		linetypeId, companyId, oprItemCode, brandId, modelId, jobTypeId, variantId, currencyId, billCode, whsGroup,
	)

	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(
		writer,
		map[string]float64{"price": price},
		"Get Data Successfully!",
		http.StatusOK,
	)
}

// @Summary Item Master For Free Accs
// @Description Item Master For Free Accs
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param company_id query int true "Company ID"
// @Param item_id query int false "Item ID"
// @Param item_code query string false "Item Code"
// @Param item_name query string false "Item Name"
// @Param uom_code query string false "UOM Code"
// @Param is_active query string false "Is Active"
// @Param brand_id query int false "Brand ID"
// @Param model_id query int false "Model ID"
// @Param variant_id query int false "Variant ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-freeaccs [get]
func (r *LookupControllerImpl) ItemMasterForFreeAccs(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	filterCondition := map[string]string{
		"company_id":                 queryValues.Get("company_id"),
		"mtr_item.item_id":           queryValues.Get("item_id"),
		"mtr_item.item_code":         queryValues.Get("item_code"),
		"mtr_item.item_name":         queryValues.Get("item_name"),
		"mtr_uom.uom_code":           queryValues.Get("uom_code"),
		"mtr_item.is_active":         queryValues.Get("is_active"),
		"mtr_item_detail.brand_id":   queryValues.Get("brand_id"),
		"mtr_item_detail.model_id":   queryValues.Get("model_id"),
		"mtr_item_detail.variant_id": queryValues.Get("variant_id"),
	}

	criteria := utils.BuildFilterCondition(filterCondition)

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	if paginate.GetSortBy() == "" {
		paginate.SortBy = "item_id"
	}

	item, baseErr := r.LookupService.ItemMasterForFreeAccs(criteria, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}
	payloads.NewHandleSuccessPagination(writer, item.Rows, "Get Data Successfully!", http.StatusOK, item.Limit, item.Page, item.TotalRows, item.TotalPages)
}

// @Summary Item Master For Free Accs By ID
// @Description Item Master For Free Accs By ID
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param company_id path int true "Company ID"
// @Param item_id path int true "Item ID"
// @Param item_code query string false "Item Code"
// @Param item_name query string false "Item Name"
// @Param uom_code query string false "UOM Code"
// @Param is_active query string false "Is Active"
// @Param brand_id query int false "Brand ID"
// @Param model_id query int false "Model ID"
// @Param variant_id query int false "Variant ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-freeaccs/by-id/{company_id}/{item_id} [get]
func (r *LookupControllerImpl) ItemMasterForFreeAccsById(writer http.ResponseWriter, request *http.Request) {
	companyId, _ := strconv.Atoi(chi.URLParam(request, "company_id"))
	itemId, _ := strconv.Atoi(chi.URLParam(request, "item_id"))

	item, baseErr := r.LookupService.ItemMasterForFreeAccsById(companyId, itemId)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}
	payloads.NewHandleSuccess(writer, item, "Get Data Successfully!", http.StatusOK)
}

// @Summary Item Master For Free Accs By Code
// @Description Item Master For Free Accs By Code
// @Tags Master Lookup :
// @Accept json
// @Produce json
// @Param company_id path int true "Company ID"
// @Param item_code query string true "Item Code"
// @Param item_name query string false "Item Name"
// @Param uom_code query string false "UOM Code"
// @Param is_active query string false "Is Active"
// @Param brand_id query int false "Brand ID"
// @Param model_id query int false "Model ID"
// @Param variant_id query int false "Variant ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/lookup/item-freeaccs/by-code/{company_id} [get]
func (r *LookupControllerImpl) ItemMasterForFreeAccsByCode(writer http.ResponseWriter, request *http.Request) {
	companyId, _ := strconv.Atoi(chi.URLParam(request, "company_id"))
	itemCode := request.URL.Query().Get("item_code")

	item, baseErr := r.LookupService.ItemMasterForFreeAccsByCode(companyId, itemCode)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}
	payloads.NewHandleSuccess(writer, item, "Get Data Successfully!", http.StatusOK)
}
