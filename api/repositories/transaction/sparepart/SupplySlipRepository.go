package transactionsparepartrepository

import (
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
)

type SupplySlipRepository interface {
	GetSupplySlipById(int32) (transactionsparepartpayloads.SupplySlipResponse, error)
	GetSupplySlipDetailById(int32) (transactionsparepartpayloads.SupplySlipDetailResponse, error)
}
