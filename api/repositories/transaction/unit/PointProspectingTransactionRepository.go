package repositories

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionunitpayloads "after-sales/api/payloads/transaction/unit"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type PointProspectingRepository interface {
	GetAllCompanyData(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllSalesRepresentative(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetSalesByCompanyCode(*gorm.DB, float64, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	Process(*gorm.DB, transactionunitpayloads.ProcessRequest) (bool, *exceptions.BaseErrorResponse)
}
