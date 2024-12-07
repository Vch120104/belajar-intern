package masterrepository

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OrderTypeRepository interface {
	GetAllOrderType(tx *gorm.DB, filterConditions []utils.FilterCondition) ([]masterpayloads.GetOrderTypeResponse, *exceptions.BaseErrorResponse)
	GetOrderTypeById(tx *gorm.DB, id int) (masterpayloads.GetOrderTypeResponse, *exceptions.BaseErrorResponse)
	GetOrderTypeByName(tx *gorm.DB, name string) ([]masterpayloads.GetOrderTypeResponse, *exceptions.BaseErrorResponse)
	SaveOrderType(tx *gorm.DB, req masterpayloads.OrderTypeSaveRequest) (masterentities.OrderType, *exceptions.BaseErrorResponse)
	UpdateOrderType(tx *gorm.DB, id int, req masterpayloads.OrderTypeUpdateRequest) (masterentities.OrderType, *exceptions.BaseErrorResponse)
	ChangeStatusOrderType(tx *gorm.DB, id int) (masterentities.OrderType, *exceptions.BaseErrorResponse)
	DeleteOrderType(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse)
}
