package masterservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type FieldActionService interface {
	GetAllFieldAction(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	SaveFieldAction(req masterpayloads.FieldActionResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	GetFieldActionHeaderById(Id int) (masterpayloads.FieldActionResponse, *exceptionsss_test.BaseErrorResponse)
	GetAllFieldActionVehicleDetailById(Id int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetFieldActionVehicleDetailById(Id int) (masterpayloads.FieldActionDetailResponse, *exceptionsss_test.BaseErrorResponse)
	GetAllFieldActionVehicleItemDetailById(Id int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetFieldActionVehicleItemDetailById(Id int) (masterpayloads.FieldActionItemDetailResponse, *exceptionsss_test.BaseErrorResponse)
	PostFieldActionVehicleItemDetail(Id int, req masterpayloads.FieldActionItemDetailResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	PostFieldActionVehicleDetail(Id int, req masterpayloads.FieldActionDetailResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	PostMultipleVehicleDetail(headerId int, id string) (bool, *exceptionsss_test.BaseErrorResponse)
	PostVehicleItemIntoAllVehicleDetail(headerId int, req masterpayloads.FieldActionItemDetailResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusFieldAction(id int) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusFieldActionVehicle(id int) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusFieldActionVehicleItem(id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
