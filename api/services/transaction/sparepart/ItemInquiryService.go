package transactionsparepartservice

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
)

type ItemInquiryService interface {
	GetAllItemInquiry(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdItemInquiry(filter transactionsparepartpayloads.ItemInquiryGetByIdFilter) (transactionsparepartpayloads.ItemInquiryGetByIdResponse, *exceptions.BaseErrorResponse)
}
