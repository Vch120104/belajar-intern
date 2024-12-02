package transactionsparepartrepository

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemInquiryRepository interface {
	GetAllItemInquiry(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdItemInquiry(tx *gorm.DB, filter transactionsparepartpayloads.ItemInquiryGetByIdFilter) (transactionsparepartpayloads.ItemInquiryGetByIdResponse, *exceptions.BaseErrorResponse)
}
