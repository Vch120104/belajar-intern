package transactionjpcbrepository

import (
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type BayMasterRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllActive(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllDeactive(tx *gorm.DB, filterCondition []utils.FilterCondition) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	CarWashBayDropDown(tx *gorm.DB, filterCondition []utils.FilterCondition) ([]transactionjpcbpayloads.CarWashBayDropDownResponse, *exceptions.BaseErrorResponse)
	ChangeStatus(tx *gorm.DB, request transactionjpcbpayloads.CarWashBayUpdateRequest) (transactionjpcbentities.BayMaster, *exceptions.BaseErrorResponse)
	PostCarWashBay(tx *gorm.DB, request transactionjpcbpayloads.CarWashBayPostRequest) (transactionjpcbentities.BayMaster, *exceptions.BaseErrorResponse)
	UpdateCarWashBay(tx *gorm.DB, request transactionjpcbpayloads.CarWashBayPutRequest) (transactionjpcbentities.BayMaster, *exceptions.BaseErrorResponse)
	GetCarWashBayById(tx *gorm.DB, carWashBayId int) (transactionjpcbentities.BayMaster, *exceptions.BaseErrorResponse)
}
