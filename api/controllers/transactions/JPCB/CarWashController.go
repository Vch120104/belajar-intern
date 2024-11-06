package transactionjpcbcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbservice "after-sales/api/services/transaction/JPCB"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CarWashController interface {
	GetAllCarWash(writer http.ResponseWriter, request *http.Request)
	UpdatePriority(writer http.ResponseWriter, request *http.Request)
	GetAllCarWashPriorityDropDown(writer http.ResponseWriter, request *http.Request)
	DeleteCarWash(writer http.ResponseWriter, request *http.Request)
	PostCarWash(writer http.ResponseWriter, request *http.Request)
	CarWashScreen(writer http.ResponseWriter, request *http.Request)
	UpdateBayNumberCarWashScreenn(writer http.ResponseWriter, request *http.Request)
	StartCarWash(writer http.ResponseWriter, request *http.Request)
	StopCarWash(writer http.ResponseWriter, request *http.Request)
	CancelCarWash(writer http.ResponseWriter, request *http.Request)
	GetCarWashByWorkOrderSystemNumber(writer http.ResponseWriter, request *http.Request)
}

type CarWashControllerImpl struct {
	CarWashService transactionjpcbservice.CarWashService
}

func NewCarWashController(carWashService transactionjpcbservice.CarWashService) CarWashController {
	return &CarWashControllerImpl{
		CarWashService: carWashService,
	}
}

func (r *CarWashControllerImpl) GetAllCarWash(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		//TODO
		"trx_car_wash.company_id":                   queryValues.Get("company_id"),
		"trx_work_order.work_order_document_number": queryValues.Get("work_order_document_number"),
		"trx_work_order.promise_time":               queryValues.Get("promise_time"), //TODO delete this
		"trx_work_order.promise_date":               queryValues.Get("promise_date"),
		//filter by tnkb, tnkb is from another service
		"trx_car_wash.car_wash_bay_id":      queryValues.Get("car_wash_bay_id"),
		"trx_car_wash.car_wash_status_id":   queryValues.Get("car_wash_status_id"), //TODO delete this
		"trx_car_wash.car_wash_priority_id": queryValues.Get("car_wash_priority_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: chi.URLParam(request, "sort_of"),
		SortBy: chi.URLParam(request, "sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	paginatedData, totalPages, totalRows, err := r.CarWashService.GetAll(criteria, paginate)

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

func (r *CarWashControllerImpl) UpdatePriority(writer http.ResponseWriter, request *http.Request) {
	var formRequest transactionjpcbpayloads.CarWashUpdatePriorityRequest
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	response, err := r.CarWashService.UpdatePriority(formRequest.WorkOrderSystemNumber, formRequest.CarWashPriorityId)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

func (r *CarWashControllerImpl) GetAllCarWashPriorityDropDown(writer http.ResponseWriter, request *http.Request) {
	response, err := r.CarWashService.GetAllCarWashPriorityDropDown()
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Get Data Successfully", http.StatusOK)
}

func (r *CarWashControllerImpl) DeleteCarWash(writer http.ResponseWriter, request *http.Request) {
	workOrderSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))

	delete, err := r.CarWashService.DeleteCarWash(workOrderSystemNumber)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if delete {
		payloads.NewHandleSuccess(writer, true, "Data deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

func (r *CarWashControllerImpl) PostCarWash(writer http.ResponseWriter, request *http.Request) {
	var formRequest transactionjpcbpayloads.CarWashPostRequestProps
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	insert, err := r.CarWashService.PostCarWash(formRequest.WorkOrderSystemNumber)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, insert, "Data created successfully", http.StatusOK)
}

func (r *CarWashControllerImpl) CarWashScreen(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	companyId, strConvError := strconv.Atoi(queryValues.Get("company_id"))
	if strConvError != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        strConvError,
		})
	}

	data, err := r.CarWashService.GetAllCarWashScreen(companyId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, data, "Successfully get data", http.StatusOK)
}

func (r *CarWashControllerImpl) UpdateBayNumberCarWashScreenn(writer http.ResponseWriter, request *http.Request) {
	var formRequest transactionjpcbpayloads.CarWashScreenUpdateBayNumberRequest
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	data, err := r.CarWashService.UpdateBayNumberCarWashScreen(formRequest.CarWashBayId, formRequest.WorkOrderSystemNumber)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, data, "Successfully update data", http.StatusOK)
}

func (r *CarWashControllerImpl) StartCarWash(writer http.ResponseWriter, request *http.Request) {
	var formRequest transactionjpcbpayloads.CarWashScreenUpdateBayNumberRequest

	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	data, err := r.CarWashService.StartCarWash(formRequest.WorkOrderSystemNumber, formRequest.CarWashBayId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, data, "Successfully start carwash", http.StatusOK)
}

// StopCarWash implements CarWashController.
func (r *CarWashControllerImpl) StopCarWash(writer http.ResponseWriter, request *http.Request) {
	var formRequest transactionjpcbpayloads.StopCarWashScreenRequest

	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	data, err := r.CarWashService.StopCarWash(formRequest.WorkOrderSystemNumber)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, data, "Successfully start carwash", http.StatusOK)
}

// CancelCarWash implements CarWashController.
func (r *CarWashControllerImpl) CancelCarWash(writer http.ResponseWriter, request *http.Request) {
	var formRequest transactionjpcbpayloads.StopCarWashScreenRequest

	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	data, err := r.CarWashService.CancelCarWash(formRequest.WorkOrderSystemNumber)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, data, "Successfully start carwash", http.StatusOK)
}

// GetCarWashByWorkOrderSystemNumber implements CarWashController.
func (r *CarWashControllerImpl) GetCarWashByWorkOrderSystemNumber(writer http.ResponseWriter, request *http.Request) {
	workOrderSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))

	data, err := r.CarWashService.GetCarWashByWorkOrderSystemNumber(workOrderSystemNumber)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, data, "Successfully retrieve data", http.StatusOK)
}
