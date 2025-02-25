package mastercontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type OrderTypeController interface {
	GetAllOrderType(writer http.ResponseWriter, request *http.Request)
	GetOrderTypeById(writer http.ResponseWriter, request *http.Request)
	GetOrderTypeByName(writer http.ResponseWriter, request *http.Request)
	SaveOrderType(writer http.ResponseWriter, request *http.Request)
	UpdateOrderType(writer http.ResponseWriter, request *http.Request)
	ChangeStatusOrderType(writer http.ResponseWriter, request *http.Request)
	DeleteOrderType(writer http.ResponseWriter, request *http.Request)
}

type OrderTypeControllerImpl struct {
	OrderTypeService masterservice.OrderTypeService
}

func NewOrderTypeControllerImpl(orderTypeService masterservice.OrderTypeService) OrderTypeController {
	return &OrderTypeControllerImpl{
		OrderTypeService: orderTypeService,
	}
}

// @Summary Get All Order Type
// @Description Get All Order Type
// @Tags Master : Order Type
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param is_active query string false "is_active"
// @Param order_type_id query string false "order_type_id"
// @Param order_type_code query string false "order_type_code"
// @Param order_type_name query string false "order_type_name"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/order-type [get]
func (r *OrderTypeControllerImpl) GetAllOrderType(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"is_active":       queryValues.Get("is_active"),
		"order_type_id":   queryValues.Get("order_type_id"),
		"order_type_code": queryValues.Get("order_type_code"),
		"order_type_name": queryValues.Get("order_type_name"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result, err := r.OrderTypeService.GetAllOrderType(filterCondition)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Order Type By ID
// @Description Get Order Type By ID
// @Tags Master : Order Type
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param order_type_id path string true "order_type_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/order-type/{order_type_id} [get]
func (r *OrderTypeControllerImpl) GetOrderTypeById(writer http.ResponseWriter, request *http.Request) {
	orderTypeId, errA := strconv.Atoi(chi.URLParam(request, "order_type_id"))
	if errA != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("failed to read url param, please check your param input"),
		})
		return
	}
	if orderTypeId == 0 {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("ID cannot be 0"),
		})
		return
	}

	result, err := r.OrderTypeService.GetOrderTypeById(orderTypeId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Order Type By Name
// @Description Get Order Type By Name
// @Tags Master : Order Type
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param order_type_name query string true "order_type_name"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/order-type/by-name [get]
func (r *OrderTypeControllerImpl) GetOrderTypeByName(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	orderTypeName := queryValues.Get("order_type_name")
	if orderTypeName == "" {
		payloads.NewHandleError(writer, "order_type_name is required", http.StatusBadRequest)
		return
	}

	result, err := r.OrderTypeService.GetOrderTypeByName(orderTypeName)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Order Type
// @Description Save Order Type
// @Tags Master : Order Type
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/order-type [post]
func (r *OrderTypeControllerImpl) SaveOrderType(writer http.ResponseWriter, request *http.Request) {
	formRequest := masterpayloads.OrderTypeSaveRequest{}
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.OrderTypeService.SaveOrderType(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Create Data Successfully!", http.StatusCreated)
}

// @Summary Update Order Type
// @Description Update Order Type
// @Tags Master : Order Type
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param order_type_id path string true "order_type_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/order-type/{order_type_id} [put]
func (r *OrderTypeControllerImpl) UpdateOrderType(writer http.ResponseWriter, request *http.Request) {
	orderTypeId, errA := strconv.Atoi(chi.URLParam(request, "order_type_id"))
	if errA != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("failed to read url param, please check your param input"),
		})
		return
	}
	if orderTypeId == 0 {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("ID cannot be 0"),
		})
		return
	}

	formRequest := masterpayloads.OrderTypeUpdateRequest{}
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	update, err := r.OrderTypeService.UpdateOrderType(orderTypeId, formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, update, "Update Data Successfully!", http.StatusOK)
}

// @Summary Change Status Order Type
// @Description Change Status Order Type
// @Tags Master : Order Type
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param order_type_id path string true "order_type_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/order-type/{order_type_id} [patch]
func (r *OrderTypeControllerImpl) ChangeStatusOrderType(writer http.ResponseWriter, request *http.Request) {
	orderTypeId, errA := strconv.Atoi(chi.URLParam(request, "order_type_id"))
	if errA != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("failed to read url param, please check your param input"),
		})
		return
	}
	if orderTypeId == 0 {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("ID cannot be 0"),
		})
		return
	}

	update, err := r.OrderTypeService.ChangeStatusOrderType(orderTypeId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, update, "Update Data Successfully!", http.StatusOK)
}

// @Summary Delete Order Type
// @Description Delete Order Type
// @Tags Master : Order Type
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param order_type_id path string true "order_type_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/order-type/{order_type_id} [delete]
func (r *OrderTypeControllerImpl) DeleteOrderType(writer http.ResponseWriter, request *http.Request) {
	orderTypeId, errA := strconv.Atoi(chi.URLParam(request, "order_type_id"))
	if errA != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("failed to read url param, please check your param input"),
		})
		return
	}
	if orderTypeId == 0 {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("ID cannot be 0"),
		})
		return
	}

	delete, err := r.OrderTypeService.DeleteOrderType(orderTypeId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, delete, "Delete Data Successfully!", http.StatusOK)
}
