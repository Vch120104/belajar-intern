package transactionsparepartrepository

import (
	exceptions "after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"

	"gorm.io/gorm"
)

type SupplySlipRepository interface {
	GetSupplySlipById(tx *gorm.DB, Id int) (transactionsparepartpayloads.SupplySlipResponse, *exceptions.BaseErrorResponse)
	GetSupplySlipDetailById(tx *gorm.DB, Id int) (transactionsparepartpayloads.SupplySlipDetailResponse, *exceptions.BaseErrorResponse)
}
