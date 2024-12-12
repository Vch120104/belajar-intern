package transactionsparepartrepository

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	// transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemLocationTransferRepository interface {
	GetAllItemLocationTransfer(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	// InsertItemLocationTransfer(tx *gorm.DB, request transactionsparepartpayloads.InsertItemLocationTransferRequest) (transactionjpcbpayloads.SettingTechnicianGetByIdResponse, *exceptions.BaseErrorResponse)
}