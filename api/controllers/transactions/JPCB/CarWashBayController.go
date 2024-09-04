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
	GetAllCarWashBay(writer http.ResponseWriter, request *http.Request)
	GetAllActiveCarWashBay(writer http.ResponseWriter, request *http.Request)
	GetAllDeactiveCarWashBay(writer http.ResponseWriter, request *http.Request)
	ChangeStatusCarWashBay(writer http.ResponseWriter, request *http.Request)
	GetAllCarWashBayDropDown(writer http.ResponseWriter, request *http.Request)
}

type BayMasterControllerImpl struct {
	bayMasterService transactionjpcbservice.BayMasterService
}

func NewCarWashBayController(bayMasterService transactionjpcbservice.BayMasterService) BayMasterController {
	return &BayMasterControllerImpl{
		bayMasterService: bayMasterService,
	}
}

func (r *BayMasterControllerImpl) GetAllCarWashBay(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

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
	paginatedData, totalPages, totalRows, err := r.bayMasterService.GetAllCarWashBay(criteria, paginate)

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

func (r *BayMasterControllerImpl) GetAllActiveCarWashBay(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

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
	paginatedData, totalPages, totalRows, err := r.bayMasterService.GetAllActiveCarWashBay(criteria, paginate)

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

func (r *BayMasterControllerImpl) GetAllDeactiveCarWashBay(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"company_id": queryValues.Get("company_id"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	responseData, err := r.bayMasterService.GetAllDeactiveCarWashBay(criteria)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(responseData), "Get Data Successfully", http.StatusOK)
}

func (r *BayMasterControllerImpl) ChangeStatusCarWashBay(writer http.ResponseWriter, request *http.Request) {
	valueRequest := transactionjpcbpayloads.BayMasterUpdateRequest{}
	helper.ReadFromRequestBody(request, &valueRequest)

	update, err := r.bayMasterService.ChangeStatusCarWashBay(valueRequest)
	if err != nil {
		if err.Err.Error() == "already start" {
			exceptions.NewAppException(writer, request, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusOK,
				Message:    "Can't remove Bay, Car Wash status for this bay is already Started",
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

func (r *BayMasterControllerImpl) GetAllCarWashBayDropDown(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"company_id": queryValues.Get("company_id"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	responseData, err := r.bayMasterService.GetAllCarWashBayDropDown(criteria)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(responseData), "Get Data Successfully", http.StatusOK)
}
