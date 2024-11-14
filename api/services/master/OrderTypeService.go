package masterservice

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
)

type OrderTypeService interface {
	GetAllOrderType() ([]masterpayloads.GetOrderTypeResponse, *exceptions.BaseErrorResponse)
	GetOrderTypeById(id int) (masterpayloads.GetOrderTypeResponse, *exceptions.BaseErrorResponse)
	GetOrderTypeByName(name string) ([]masterpayloads.GetOrderTypeResponse, *exceptions.BaseErrorResponse)
	SaveOrderType(req masterpayloads.OrderTypeSaveRequest) (masterentities.OrderType, *exceptions.BaseErrorResponse)
	UpdateOrderType(id int, req masterpayloads.OrderTypeUpdateRequest) (masterentities.OrderType, *exceptions.BaseErrorResponse)
	ChangeStatusOrderType(id int) (masterentities.OrderType, *exceptions.BaseErrorResponse)
	DeleteOrderType(id int) (bool, *exceptions.BaseErrorResponse)
}
