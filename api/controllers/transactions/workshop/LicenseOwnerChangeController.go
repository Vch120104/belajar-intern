package transactionworkshopcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type LicenseOwnerChangeControllerImpl struct {
	LicenseOwnerChangeService transactionworkshopservice.LicenseOwnerChangeService
}
type LicenseOwnerChangeController interface {
	GetAll(writer http.ResponseWriter, request *http.Request)
	GetHistoryByChassisNumber(writer http.ResponseWriter, request *http.Request)
}

func NewLicenseOwnerChangeController(LicenseOwnerChangeService transactionworkshopservice.LicenseOwnerChangeService) LicenseOwnerChangeController {
	return &LicenseOwnerChangeControllerImpl{
		LicenseOwnerChangeService: LicenseOwnerChangeService,
	}
}

// @Summary Get All License Owner Change
// @Description Retrieve all license owner change with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Transaction : Workshop License Owner Change
// @Param vehicle_brand query string false "Vehicle Brand"
// @Param model_code query string false "Model Code"
// @Param vehicle_id query string false "Vehicle ID"
// @Param tnkb_old query string false "TNKB Old"
// @Param tnkb_new query string false "TNKB New"
// @Param owner_name_old query string false "Owner Name Old"
// @Param owner_name_new query string false "Owner Name New"
// @Param change_date_from query string false "Change Date From"
// @Param change_date_to query string false "Change Date To"
// @Param limit query string true "Items per page"
// @Param page query string true "Page number"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/license-owner-change [get]
func (l *LicenseOwnerChangeControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_license_owner_change.brand_id":         queryValues.Get("vehicle_brand"),
		"trx_license_owner_change.model_id":         queryValues.Get("model_code"),
		"trx_license_owner_change.vehicle_id":       queryValues.Get("vehicle_id"),
		"trx_license_owner_change.tnkb_old":         queryValues.Get("tnkb_old"),
		"trx_license_owner_change.tnkb_new":         queryValues.Get("tnkb_new"),
		"trx_license_owner_change.owner_name_old":   queryValues.Get("owner_name_old"),
		"trx_license_owner_change.owner_name_new":   queryValues.Get("owner_name_new"),
		"trx_license_owner_change.change_date_from": queryValues.Get("change_date_from"),
		"trx_license_owner_change.change_date_to":   queryValues.Get("change_date_to"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := l.LicenseOwnerChangeService.GetAll(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	// Respond with the paginated data
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

// @Summary Get License Owner Change History By Chassis Number
// @Description Retrieve license owner change history by chassis number with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Transaction : Workshop License Owner Change
// @Param vehicle_chassis_number path string true "Vehicle Chassis Number"
// @Param change_type query string false "Change Type"
// @Param limit query string true "Items per page"
// @Param page query string true "Page number"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/license-owner-change/history/{vehicle_chassis_number} [get]
func (l *LicenseOwnerChangeControllerImpl) GetHistoryByChassisNumber(writer http.ResponseWriter, request *http.Request) {
	chassisNumber := chi.URLParam(request, "vehicle_chassis_number")
	if chassisNumber == "" {
		payloads.NewHandleError(writer, "Vehicle Chassis Number is required", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterParams := map[string]string{
		"change_type": queryValues.Get("change_type"),
	}

	criteria := utils.BuildFilterCondition(filterParams)

	result, err := l.LicenseOwnerChangeService.GetHistoryByChassisNumber(chassisNumber, criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		result.Rows,
		"Get History Successfully!",
		http.StatusOK,
		result.Limit,
		result.Page,
		int64(result.TotalRows),
		result.TotalPages,
	)
}
