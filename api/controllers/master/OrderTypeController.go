package mastercontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterservice "after-sales/api/services/master"
	"net/http"
)

type OrderTypeController interface {
	GetAllOrderType(writer http.ResponseWriter, request *http.Request)
}

type OrderTypeControllerImpl struct {
	OrderTypeService masterservice.OrderTypeService
}

func NewOrderTypeControllerImpl(orderTypeService masterservice.OrderTypeService) OrderTypeController {
	return &OrderTypeControllerImpl{
		OrderTypeService: orderTypeService,
	}
}

func (r *OrderTypeControllerImpl) GetAllOrderType(writer http.ResponseWriter, request *http.Request) {
	result, err := r.OrderTypeService.GetAllOrderType()
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}
