package transactionjpcbcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbservice "after-sales/api/services/transaction/JPCB"
	"after-sales/api/utils"
	"errors"
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

func (r *BayMasterControllerImpl) UpdateBayMaster(writer http.ResponseWriter, request *http.Request) {
	valueRequest := transactionjpcbpayloads.BayMasterUpdateRequest{}
	helper.ReadFromRequestBody(request, &valueRequest)

	update, err := r.bayMasterService.UpdateBayMaster(valueRequest)
	if err != nil {
		if err.Err.Error() == "already start" {
			exceptions.NewAppException(writer, request, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusOK,
				Message:    "Already Start",
				Data:       nil,
				Err:        errors.New("already start"),
			})
			return
		}
		if err.Err.Error() == "bay not found" {
			exceptions.NewAppException(writer, request, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusOK,
				Message:    "Bay Not Found",
				Data:       nil,
				Err:        errors.New("bay not found"),
			})
			return
		}
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, update, "Bay updated successfully", http.StatusOK)
}
