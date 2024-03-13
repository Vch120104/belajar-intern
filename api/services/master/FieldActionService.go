package masterservice

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type FieldActionService interface {
	GetAllFieldAction(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination
	SaveFieldAction(req masterpayloads.FieldActionResponse) bool
	GetFieldActionHeaderById(Id int) masterpayloads.FieldActionResponse
	GetAllFieldActionVehicleDetailById(Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination
	GetFieldActionVehicleDetailById(Id int) masterpayloads.FieldActionDetailResponse
	GetAllFieldActionVehicleItemDetailById(Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination
	GetFieldActionVehicleItemDetailById(Id int) masterpayloads.FieldActionItemDetailResponse
	PostFieldActionVehicleItemDetail(Id int, req masterpayloads.FieldActionItemDetailResponse) bool
	PostFieldActionVehicleDetail(Id int, req masterpayloads.FieldActionDetailResponse) bool
	PostMultipleVehicleDetail(headerId int, companyId int, id string) bool
	PostVehicleItemIntoAllVehicleDetail(headerId int, req masterpayloads.FieldActionItemDetailResponse) bool
}
