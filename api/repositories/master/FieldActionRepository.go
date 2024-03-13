package masterrepository

import (
	// masterpayloads "after-sales/api/payloads/master"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type FieldActionRepository interface {
	GetAllFieldAction(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
	SaveFieldAction(tx *gorm.DB, req masterpayloads.FieldActionResponse) (bool, error)

	GetFieldActionHeaderById(tx *gorm.DB, Id int) (masterpayloads.FieldActionResponse, error)
	GetAllFieldActionVehicleDetailById(tx *gorm.DB, Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error)
	GetFieldActionVehicleDetailById(tx *gorm.DB, Id int) (masterpayloads.FieldActionDetailResponse, error)
	GetAllFieldActionVehicleItemDetailById(tx *gorm.DB, Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error)
	GetFieldActionVehicleItemDetailById(tx *gorm.DB, Id int) (masterpayloads.FieldActionItemDetailResponse, error)
	PostFieldActionVehicleItemDetail(tx *gorm.DB, req masterpayloads.FieldActionItemDetailResponse, id int) (bool, error)
	PostFieldActionVehicleDetail(tx *gorm.DB, req masterpayloads.FieldActionDetailResponse, id int) (bool, error)
	PostMultipleVehicleDetail(tx *gorm.DB, headerId int, companyId int, id string) (bool, error)
	PostVehicleItemIntoAllVehicleDetail(tx *gorm.DB, headerId int, req masterpayloads.FieldActionItemDetailResponse) (bool, error)
}
