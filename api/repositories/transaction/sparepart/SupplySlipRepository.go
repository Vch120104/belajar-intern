package transactionsparepartrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"

	"gorm.io/gorm"
)

type SupplySlipRepository interface {
	GetSupplySlipById(tx *gorm.DB, Id int) (transactionsparepartpayloads.SupplySlipResponse, *exceptionsss_test.BaseErrorResponse)
	GetSupplySlipDetailById(tx *gorm.DB, Id int) (transactionsparepartpayloads.SupplySlipDetailResponse, *exceptionsss_test.BaseErrorResponse)
}
