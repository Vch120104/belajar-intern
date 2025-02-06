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

// @Summary Get All Car Wash
// @Description Get All Car Wash
// @Tags Transaction JPCB: Car Wash
// @Accept json
// @Produce json
// @Param company_id query string false "Company ID"
// @Param work_order_document_number query string false "Work Order Document Number"
// @Param promise_time query string false "Promise Time"
// @Param promise_date query string false "Promise Date"
// @Param car_wash_bay_id query string false "Car Wash Bay ID"
// @Param car_wash_status_id query string false "Car Wash Status ID"
// @Param car_wash_priority_id query string false "Car Wash Priority ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of path string false "Sort Of"
// @Param sort_by path string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/car-wash [get]
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
	result, err := r.CarWashService.GetAll(criteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

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

// @Summary Update Car Wash Priority
// @Description Update Car Wash Priority
// @Tags Transaction JPCB: Car Wash
// @Accept json
// @Produce json
// @Param work_order_system_number path string true "Work Order System Number"
// @Param car_wash_priority_id path string true "Car Wash Priority ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/car-wash/update-priority [put]
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

// @Summary Get All Car Wash Priority Drop Down
// @Description Get All Car Wash Priority Drop Down
// @Tags Transaction JPCB: Car Wash
// @Accept json
// @Produce json
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/car-wash/priority/dropdown [get]
func (r *CarWashControllerImpl) GetAllCarWashPriorityDropDown(writer http.ResponseWriter, request *http.Request) {
	response, err := r.CarWashService.GetAllCarWashPriorityDropDown()
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Get Data Successfully", http.StatusOK)
}

// @Summary Delete Car Wash
// @Description Delete Car Wash
// @Tags Transaction JPCB: Car Wash
// @Accept json
// @Produce json
// @Param work_order_system_number path int true "Work Order System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/car-wash/{work_order_system_number} [delete]
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

// @Summary Post Car Wash
// @Description Post Car Wash
// @Tags Transaction JPCB: Car Wash
// @Accept json
// @Produce json
// @Param request body transactionjpcbpayloads.CarWashPostRequestProps true "Car Wash Post Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/car-wash [post]
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

// @Summary Car Wash Screen
// @Description Car Wash Screen
// @Tags Transaction JPCB: Car Wash
// @Accept json
// @Produce json
// @Param company_id query string false "Company ID"
// @Param car_wash_status_id query string false "Car Wash Status ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/car-wash/screen [get]
func (r *CarWashControllerImpl) CarWashScreen(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	companyId, strConvError := strconv.Atoi(queryValues.Get("company_id"))
	carWashStatusId, strConvErrorCarWashStatus := strconv.Atoi(queryValues.Get("car_wash_status_id"))
	if strConvError != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        strConvError,
		})
	}
	if strConvErrorCarWashStatus != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        strConvError,
		})
	}

	data, err := r.CarWashService.GetAllCarWashScreen(companyId, carWashStatusId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, data, "Successfully get data", http.StatusOK)
}

// @Summary Update Bay Number Car Wash Screen
// @Description Update Bay Number Car Wash Screen
// @Tags Transaction JPCB: Car Wash
// @Accept json
// @Produce json
// @Param request body transactionjpcbpayloads.CarWashScreenUpdateBayNumberRequest true "Car Wash Screen Update Bay Number Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/car-wash/screen/update-bay [put]
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

// @Summary Start Car Wash
// @Description Start Car Wash
// @Tags Transaction JPCB: Car Wash
// @Accept json
// @Produce json
// @Param request body transactionjpcbpayloads.CarWashScreenUpdateBayNumberRequest true "Car Wash Screen Update Bay Number Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/car-wash/start [put]
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

// @Summary Stop Car Wash
// @Description Stop Car Wash
// @Tags Transaction JPCB: Car Wash
// @Accept json
// @Produce json
// @Param request body transactionjpcbpayloads.StopCarWashScreenRequest true "Stop Car Wash Screen Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/car-wash/stop [put]
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

// @Summary Cancel Car Wash
// @Description Cancel Car Wash
// @Tags Transaction JPCB: Car Wash
// @Accept json
// @Produce json
// @Param request body transactionjpcbpayloads.StopCarWashScreenRequest true "Stop Car Wash Screen Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/car-wash/cancel [put]
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

// @Summary Get Car Wash By Work Order System Number
// @Description Get Car Wash By Work Order System Number
// @Tags Transaction JPCB: Car Wash
// @Accept json
// @Produce json
// @Param work_order_system_number path int true "Work Order System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/car-wash/{work_order_system_number} [get]
func (r *CarWashControllerImpl) GetCarWashByWorkOrderSystemNumber(writer http.ResponseWriter, request *http.Request) {
	workOrderSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))

	data, err := r.CarWashService.GetCarWashByWorkOrderSystemNumber(workOrderSystemNumber)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, data, "Successfully retrieve data", http.StatusOK)
}
