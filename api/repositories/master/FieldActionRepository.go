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
	GetAllFieldAction(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveFieldAction(tx *gorm.DB, req masterpayloads.FieldActionRequest) (bool, *exceptions.BaseErrorResponse)

	GetFieldActionHeaderById(tx *gorm.DB, Id int) (masterpayloads.FieldActionResponse, *exceptions.BaseErrorResponse)
	GetAllFieldActionVehicleDetailById(tx *gorm.DB, Id int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetFieldActionVehicleDetailById(tx *gorm.DB, Id int) (masterpayloads.FieldActionDetailResponse, *exceptions.BaseErrorResponse)
	GetFieldActionVehicleItemDetailById(tx *gorm.DB, Id int) (masterpayloads.FieldActionEligibleVehicleItemOperationResp, *exceptions.BaseErrorResponse)
	GetAllFieldActionVehicleItemOperationDetailById(tx *gorm.DB, Id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	PostFieldActionVehicleItemDetail(tx *gorm.DB, req masterpayloads.FieldActionEligibleVehicleItemOperationRequest, id int) (bool, *exceptions.BaseErrorResponse)
	PostFieldActionVehicleDetail(tx *gorm.DB, req masterpayloads.FieldActionDetailResponse, id int) (bool, *exceptions.BaseErrorResponse)
	PostMultipleVehicleDetail(tx *gorm.DB, headerId int, id string) (bool, *exceptions.BaseErrorResponse)
	PostVehicleItemIntoAllVehicleDetail(tx *gorm.DB, headerId int, req masterpayloads.FieldActionEligibleVehicleItemOperationRequest) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusFieldAction(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusFieldActionVehicle(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusFieldActionVehicleItem(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse)
}
