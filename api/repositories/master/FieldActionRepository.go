package masterrepository

import (
	// masterpayloads "after-sales/api/payloads/master"
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type FieldActionRepository interface {
	GetAllFieldAction(*gorm.DB, []utils.FilterCondition, pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	SaveFieldAction(tx *gorm.DB, req masterpayloads.FieldActionRequest) (bool, *exceptionsss_test.BaseErrorResponse)

	GetFieldActionHeaderById(tx *gorm.DB, Id int) (masterpayloads.FieldActionResponse, *exceptionsss_test.BaseErrorResponse)
	GetAllFieldActionVehicleDetailById(tx *gorm.DB, Id int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetFieldActionVehicleDetailById(tx *gorm.DB, Id int) (masterpayloads.FieldActionDetailResponse, *exceptionsss_test.BaseErrorResponse)
	GetAllFieldActionVehicleItemDetailById(tx *gorm.DB, Id int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetFieldActionVehicleItemDetailById(tx *gorm.DB, Id int) (masterpayloads.FieldActionItemDetailResponse, *exceptionsss_test.BaseErrorResponse)
	PostFieldActionVehicleItemDetail(tx *gorm.DB, req masterpayloads.FieldActionItemDetailResponse, id int) (bool, *exceptionsss_test.BaseErrorResponse)
	PostFieldActionVehicleDetail(tx *gorm.DB, req masterpayloads.FieldActionDetailResponse, id int) (bool, *exceptionsss_test.BaseErrorResponse)
	PostMultipleVehicleDetail(tx *gorm.DB, headerId int, id string) (bool, *exceptionsss_test.BaseErrorResponse)
	PostVehicleItemIntoAllVehicleDetail(tx *gorm.DB, headerId int, req masterpayloads.FieldActionItemDetailResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusFieldAction(tx *gorm.DB, id int) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusFieldActionVehicle(tx *gorm.DB, id int) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusFieldActionVehicleItem(tx *gorm.DB, id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
