package transactionjpcbservice

import (
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	"after-sales/api/utils"
)

type CarWashService interface {
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	UpdatePriority(workOrderSystemNumber, carWashPriorityId int) (transactionjpcbentities.CarWash, *exceptions.BaseErrorResponse)
	GetAllCarWashPriorityDropDown() ([]transactionjpcbpayloads.CarWashPriorityDropDownResponse, *exceptions.BaseErrorResponse)
	DeleteCarWash(workOrderSystemNumber int) (bool, *exceptions.BaseErrorResponse)
	PostCarWash(workOrderSystemNumber int) (transactionjpcbpayloads.CarWashPostResponse, *exceptions.BaseErrorResponse)
	GetCarWashByWorkOrderSystemNumber(workOrderSystemNumber int) (transactionjpcbpayloads.CarWashGetAllResponse, *exceptions.BaseErrorResponse)

	GetAllCarWashScreen(companyId int, carWashStatusId int) ([]transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse)
	UpdateBayNumberCarWashScreen(bayNumber, workOrderSystemNumber int) (transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse)
	StartCarWash(workOrderSystemNumber, carWashBayId int) (transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse)
	StopCarWash(workOrderSystemNumber int) (transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse)
	CancelCarWash(workOrderSystemNumber int) (transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse)
}
