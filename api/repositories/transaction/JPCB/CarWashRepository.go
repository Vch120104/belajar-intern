package transactionjpcbrepository

import (
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type CarWashRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	UpdatePriority(tx *gorm.DB, workOrderSystemNumber, carWashPriorityId int) (transactionjpcbentities.CarWash, *exceptions.BaseErrorResponse)
	GetAllCarWashPriorityDropDown(tx *gorm.DB) ([]transactionjpcbpayloads.CarWashPriorityDropDownResponse, *exceptions.BaseErrorResponse)
	PostCarWash(tx *gorm.DB, workOrderSystemNumber int) (transactionjpcbpayloads.CarWashPostResponse, *exceptions.BaseErrorResponse)
	DeleteCarWash(tx *gorm.DB, workOrderSystemNumber int) (bool, *exceptions.BaseErrorResponse)
	GetCarWashByWorkOrderSystemNumber(tx *gorm.DB, workOrderSystemNumber int) (transactionjpcbpayloads.CarWashGetAllResponse, *exceptions.BaseErrorResponse)

	GetAllCarWashScreen(tx *gorm.DB, companyId int, carWashStatusId int) ([]transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse)
	UpdateBayNumberCarWashScreen(tx *gorm.DB, bayNumber, workOrderSystemNumber int) (transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse)
	StartCarWash(tx *gorm.DB, workOrderSystemNumber, carWashBayId int) (transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse)
	StopCarWash(tx *gorm.DB, workOrderSystemNumber int) (transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse)
	CancelCarWash(tx *gorm.DB, workOrderSystemNumber int) (transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse)
}
