package transactionjpcbcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbservice "after-sales/api/services/transaction/JPCB"
	"after-sales/api/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type BayMasterController interface {
	GetAllBayMaster(writer http.ResponseWriter, request *http.Request)
	GetAllActiveBayCarWashScreen(writer http.ResponseWriter, request *http.Request)
	GetAllDeactiveBayCarWashScreen(writer http.ResponseWriter, request *http.Request)
	UpdateBayMaster(writer http.ResponseWriter, request *http.Request)
}

type BayMasterControllerImpl struct {
	bayMasterService transactionjpcbservice.BayMasterService
}

func BayController(bayMasterService transactionjpcbservice.BayMasterService) BayMasterController {
	return &BayMasterControllerImpl{
		bayMasterService: bayMasterService,
	}
}

// getAllBayMaster implements BayMasterController.
func (r *BayMasterControllerImpl) GetAllBayMaster(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query() // Retrieve query parameters

	queryParams := map[string]string{
		"company_id": queryValues.Get("company_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: chi.URLParam(request, "sort_of"),
		SortBy: chi.URLParam(request, "sort_by"),
	}
	print(queryParams)

	criteria := utils.BuildFilterCondition(queryParams)
	paginatedData, totalPages, totalRows, err := r.bayMasterService.GetAllBayMaster(criteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// GetAllActiveBayCarWashScreen implements BayMasterController.
func (r *BayMasterControllerImpl) GetAllActiveBayCarWashScreen(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query() // Retrieve query parameters

	queryParams := map[string]string{
		"company_id": queryValues.Get("company_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: chi.URLParam(request, "sort_of"),
		SortBy: chi.URLParam(request, "sort_by"),
	}
	print(queryParams)

	criteria := utils.BuildFilterCondition(queryParams)
	paginatedData, totalPages, totalRows, err := r.bayMasterService.GetAllActiveBayCarWashScreen(criteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// GetAllDeactiveBayCarWashScreen implements BayMasterController.
func (r *BayMasterControllerImpl) GetAllDeactiveBayCarWashScreen(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query() // Retrieve query parameters

	queryParams := map[string]string{
		"company_id": queryValues.Get("company_id"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	responseData, err := r.bayMasterService.GetAllDeactiveBayCarWashScreen(criteria)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(responseData), "Get Data Successfully", http.StatusOK)
}

// UpdateBayMaster implements BayMasterController.
func (r *BayMasterControllerImpl) UpdateBayMaster(writer http.ResponseWriter, request *http.Request) {
	// companyId, _ := strconv.Atoi(chi.URLParam(request, "company_id"))
	// car_wash_bay_id, _ := strconv.Atoi(chi.URLParam(request, "car_wash_bay_id"))

	var valueRequest transactionjpcbpayloads.BayMasterUpdateRequest
	helper.ReadFromRequestBody(request, &valueRequest)

	update, err := r.bayMasterService.UpdateBayMaster(valueRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
	}

	payloads.NewHandleSuccess(writer, update, "Bay updated successfully", http.StatusOK)
}

//  agreementID, _ := strconv.Atoi(chi.URLParam(request, "agreement_id"))
// 	valueID, _ := strconv.Atoi(chi.URLParam(request, "agreement_discount_id"))

// 	var valueRequest masterpayloads.DiscountValueRequest
// 	helper.ReadFromRequestBody(request, &valueRequest)

// 	update, err := r.AgreementService.UpdateDiscountValue(int(agreementID), int(valueID), valueRequest)
// 	if err != nil {
// 		exceptions.NewAppException(writer, request, err)
// 		return
// 	}

// 	payloads.NewHandleSuccess(writer, update, "Discount value updated successfully", http.StatusOK)
