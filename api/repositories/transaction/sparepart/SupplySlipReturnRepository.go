package transactionsparepartrepository

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type SupplySlipReturnRepository interface {
	SaveSupplySlipReturn(tx *gorm.DB, request transactionsparepartentities.SupplySlipReturn) (transactionsparepartentities.SupplySlipReturn, *exceptions.BaseErrorResponse)
	SaveSupplySlipReturnDetail(tx *gorm.DB, request transactionsparepartentities.SupplySlipReturnDetail) (transactionsparepartentities.SupplySlipReturnDetail, *exceptions.BaseErrorResponse)
	GetAllSupplySlipReturn(tx *gorm.DB, internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetSupplySlipReturnById(tx *gorm.DB, Id int, pagination pagination.Pagination, supplySlip map[string]interface{}) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetSupplySlipReturnDetailById(tx *gorm.DB, Id int) (transactionsparepartpayloads.SupplySlipReturnDetailResponse, *exceptions.BaseErrorResponse)
	UpdateSupplySlipReturn(tx *gorm.DB, req transactionsparepartentities.SupplySlipReturn, id int) (transactionsparepartentities.SupplySlipReturn, *exceptions.BaseErrorResponse)
	UpdateSupplySlipReturnDetail(tx *gorm.DB, req transactionsparepartentities.SupplySlipReturnDetail, id int) (transactionsparepartentities.SupplySlipReturnDetail, *exceptions.BaseErrorResponse)
	GetSupplySlipId(tx *gorm.DB, Id int) (int, *exceptions.BaseErrorResponse)
}