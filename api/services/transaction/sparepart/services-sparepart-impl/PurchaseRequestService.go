package transactionsparepartserviceimpl

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
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

type PurchaseRequestServiceImpl struct {
	PurchaseRequestRepo transactionsparepartrepository.PurchaseRequestRepository
	DB                  *gorm.DB
	RedisClient         *redis.Client
}

func NewPurchaseRequestImpl(PurchaseRequestRepo transactionsparepartrepository.PurchaseRequestRepository, db *gorm.DB, redis *redis.Client) transactionsparepartservice.PurchaseRequestService {
	return &PurchaseRequestServiceImpl{
		PurchaseRequestRepo: PurchaseRequestRepo,
		DB:                  db,
		RedisClient:         redis,
	}
}

func (p *PurchaseRequestServiceImpl) GetAllPurchaseRequest(filterCondition []utils.FilterCondition, pages pagination.Pagination, Dateparams map[string]string) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	//TODO implement me
	tx := p.DB.Begin()
	result, err := p.PurchaseRequestRepo.GetAllPurchaseRequest(tx, filterCondition, pages, Dateparams)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (p *PurchaseRequestServiceImpl) GetByIdPurchaseRequest(id int) (transactionsparepartpayloads.PurchaseRequestGetByIdNormalizeResponses, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	result, err := p.PurchaseRequestRepo.GetByIdPurchaseRequest(tx, id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (p *PurchaseRequestServiceImpl) GetAllPurchaseRequestDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	//TODO implement me
	tx := p.DB.Begin()
	result, err := p.PurchaseRequestRepo.GetAllPurchaseRequestDetail(tx, filterCondition, pages)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (p *PurchaseRequestServiceImpl) GetByIdPurchaseRequestDetail(id int) (transactionsparepartpayloads.PurchaseRequestDetailResponsesPayloads, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	result, err := p.PurchaseRequestRepo.GetByIdPurchaseRequestDetail(tx, id)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (p *PurchaseRequestServiceImpl) PurchaseRequestSaveNewHeader(request transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest) (transactionsparepartentities.PurchaseRequestEntities, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	result, err := p.PurchaseRequestRepo.PurchaseRequestSaveHeader(tx, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (p *PurchaseRequestServiceImpl) PurchaseRequestSaveNewDetail(payloads transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads) (transactionsparepartentities.PurchaseRequestDetail, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	result, err := p.PurchaseRequestRepo.PurchaseRequestSaveDetail(tx, payloads)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}
