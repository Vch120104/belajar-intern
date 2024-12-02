package transactionsparepartserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ItemInquiryServiceImpl struct {
	ItemInquiryRepo transactionsparepartrepository.ItemInquiryRepository
	DB              *gorm.DB
	redis           *redis.Client
}

func StartItemInquiryService(itemInquiryRepo transactionsparepartrepository.ItemInquiryRepository, db *gorm.DB, rdb *redis.Client) transactionsparepartservice.ItemInquiryService {
	return &ItemInquiryServiceImpl{
		ItemInquiryRepo: itemInquiryRepo,
		DB:              db,
		redis:           rdb,
	}
}

func (i *ItemInquiryServiceImpl) GetAllItemInquiry(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
	result, err := i.ItemInquiryRepo.GetAllItemInquiry(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return pages, err
	}
	return result, nil
}

func (i *ItemInquiryServiceImpl) GetByIdItemInquiry(filter transactionsparepartpayloads.ItemInquiryGetByIdFilter) (transactionsparepartpayloads.ItemInquiryGetByIdResponse, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
	result, err := i.ItemInquiryRepo.GetByIdItemInquiry(tx, filter)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}
