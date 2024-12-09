package transactionworkshopcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"fmt"
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

// GetAll implements LicenseOwnerChangeController.
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

	fmt.Println("Query Params:", queryParams)

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	// Build the filter condition based on queryParams
	criteria := utils.BuildFilterCondition(queryParams)
	fmt.Println("Filter Conditions:", criteria)

	// Call the service to get data with pagination and filter
	paginatedData, totalPages, totalRows, err := l.LicenseOwnerChangeService.GetAll(criteria, paginate)
	if err != nil {
		// Handle error if service returns an error
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	// If data exists, return successful response with paginated data
	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		// If no data found, return 404
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// GetHistoryByChassisNumber implements LicenseOwnerChangeController.
func (l *LicenseOwnerChangeControllerImpl) GetHistoryByChassisNumber(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	chassisNumber := chi.URLParam(request, "vehicle_chassis_number")
	if chassisNumber == "" {
		payloads.NewHandleError(writer, "Vehicle Chassis Number is requried", http.StatusBadRequest)
		return
	}

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

	results, totalPages, totalRows, err := l.LicenseOwnerChangeService.GetHistoryByChassisNumber(chassisNumber, criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	if len(results) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(results), "Get History Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "No history data found for the given chassis number", http.StatusNotFound)
	}
}
