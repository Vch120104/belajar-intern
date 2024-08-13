package transactionsparepartrepository

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type SupplySlipRepository interface {
	GetSupplySlipById(tx *gorm.DB, Id int, pagination pagination.Pagination) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetSupplySlipDetailById(tx *gorm.DB, Id int) (transactionsparepartpayloads.SupplySlipDetailResponse, *exceptions.BaseErrorResponse)
	SaveSupplySlip(tx *gorm.DB, request transactionsparepartentities.SupplySlip) (transactionsparepartentities.SupplySlip, *exceptions.BaseErrorResponse)
	SaveSupplySlipDetail(tx *gorm.DB, request transactionsparepartentities.SupplySlipDetail) (transactionsparepartentities.SupplySlipDetail, *exceptions.BaseErrorResponse)
	GetAllSupplySlip(tx *gorm.DB, internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	UpdateSupplySlip(tx *gorm.DB, req transactionsparepartentities.SupplySlip, id int) (transactionsparepartentities.SupplySlip, *exceptions.BaseErrorResponse)
	UpdateSupplySlipDetail(tx *gorm.DB, req transactionsparepartentities.SupplySlipDetail, id int) (transactionsparepartentities.SupplySlipDetail, *exceptions.BaseErrorResponse)
}
