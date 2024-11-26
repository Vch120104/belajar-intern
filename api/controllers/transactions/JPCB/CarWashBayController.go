package transactionjpcbcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbservice "after-sales/api/services/transaction/JPCB"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type BayMasterController interface {
	GetAllCarWashBay(writer http.ResponseWriter, request *http.Request)
	GetAllActiveCarWashBay(writer http.ResponseWriter, request *http.Request)
	GetAllDeactiveCarWashBay(writer http.ResponseWriter, request *http.Request)
	ChangeStatusCarWashBay(writer http.ResponseWriter, request *http.Request)
	GetAllCarWashBayDropDown(writer http.ResponseWriter, request *http.Request)
	PostCarWashBay(writer http.ResponseWriter, request *http.Request)
	PutCarWashBay(writer http.ResponseWriter, request *http.Request)
	GetCarWashBayById(writer http.ResponseWriter, request *http.Request)
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
	valueRequest := transactionjpcbpayloads.CarWashBayUpdateRequest{}
	helper.ReadFromRequestBody(request, &valueRequest)
	if validationErr := validation.ValidationForm(writer, request, &valueRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

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

func (r *BayMasterControllerImpl) PostCarWashBay(writer http.ResponseWriter, request *http.Request) {
	var CarWashBayPostPayloads transactionjpcbpayloads.CarWashBayPostRequest
	helper.ReadFromRequestBody(request, &CarWashBayPostPayloads)

	response, err := r.bayMasterService.PostCarWashBay(CarWashBayPostPayloads)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Successfully Inserted Car Wash Bay", http.StatusCreated)
}

func (r *BayMasterControllerImpl) PutCarWashBay(writer http.ResponseWriter, request *http.Request) {
	var CarWashBayPutPayloads transactionjpcbpayloads.CarWashBayPutRequest
	helper.ReadFromRequestBody(request, &CarWashBayPutPayloads)

	response, err := r.bayMasterService.PutCarWashBay(CarWashBayPutPayloads)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Successfully Updated Car Wash Bay", http.StatusOK)
}

func (r *BayMasterControllerImpl) GetCarWashBayById(writer http.ResponseWriter, request *http.Request) {
	carWashBayId, errA := strconv.Atoi(chi.URLParam(request, "car_wash_bay_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.bayMasterService.GetCarWashBayById(carWashBayId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}
