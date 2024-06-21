package masterservice

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type FieldActionService interface {
	GetAllFieldAction(filterCondition []utils.FilterCondition, pages pagination.Pagination)  (pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveFieldAction(req masterpayloads.FieldActionRequest) (bool, *exceptions.BaseErrorResponse)
	GetFieldActionHeaderById(Id int) (masterpayloads.FieldActionResponse, *exceptions.BaseErrorResponse)
	GetAllFieldActionVehicleDetailById(Id int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetFieldActionVehicleDetailById(Id int) (masterpayloads.FieldActionDetailResponse, *exceptions.BaseErrorResponse)
	GetAllFieldActionVehicleItemDetailById(Id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetFieldActionVehicleItemDetailById(Id int) (masterpayloads.FieldActionItemDetailResponse, *exceptions.BaseErrorResponse)
	PostFieldActionVehicleItemDetail(Id int, req masterpayloads.FieldActionItemDetailResponse) (bool, *exceptions.BaseErrorResponse)
	PostFieldActionVehicleDetail(Id int, req masterpayloads.FieldActionDetailResponse) (bool, *exceptions.BaseErrorResponse)
	PostMultipleVehicleDetail(headerId int, id string) (bool, *exceptions.BaseErrorResponse)
	PostVehicleItemIntoAllVehicleDetail(headerId int, req masterpayloads.FieldActionItemDetailResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusFieldAction(id int) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusFieldActionVehicle(id int) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusFieldActionVehicleItem(id int) (bool, *exceptions.BaseErrorResponse)
}
