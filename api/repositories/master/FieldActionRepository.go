package masterrepository

import (
	// masterpayloads "after-sales/api/payloads/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type FieldActionRepository interface {
	GetAllFieldAction(*gorm.DB, []utils.FilterCondition, pagination.Pagination)(pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveFieldAction(tx *gorm.DB, req masterpayloads.FieldActionRequest) (bool, *exceptions.BaseErrorResponse)

	GetFieldActionHeaderById(tx *gorm.DB, Id int) (masterpayloads.FieldActionResponse, *exceptions.BaseErrorResponse)
	GetAllFieldActionVehicleDetailById(tx *gorm.DB, Id int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetFieldActionVehicleDetailById(tx *gorm.DB, Id int) (masterpayloads.FieldActionDetailResponse, *exceptions.BaseErrorResponse)
	GetFieldActionVehicleItemDetailById(*gorm.DB, int, int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetAllFieldActionVehicleItemDetailById(tx *gorm.DB, Id int, pages pagination.Pagination) ([]map[string]interface{},int,int, *exceptions.BaseErrorResponse)
	PostFieldActionVehicleItemDetail(tx *gorm.DB, req masterpayloads.FieldActionItemDetailResponse, id int) (bool, *exceptions.BaseErrorResponse)
	PostFieldActionVehicleDetail(tx *gorm.DB, req masterpayloads.FieldActionDetailResponse, id int) (bool, *exceptions.BaseErrorResponse)
	PostMultipleVehicleDetail(tx *gorm.DB, headerId int, id string) (bool, *exceptions.BaseErrorResponse)
	PostVehicleItemIntoAllVehicleDetail(tx *gorm.DB, headerId int, req masterpayloads.FieldActionItemDetailResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusFieldAction(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusFieldActionVehicle(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusFieldActionVehicleItem(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse)
}
