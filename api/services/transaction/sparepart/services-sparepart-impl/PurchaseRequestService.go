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
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (p *PurchaseRequestServiceImpl) GetByIdPurchaseRequest(id int) (transactionsparepartpayloads.PurchaseRequestGetByIdResponses, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	result, err := p.PurchaseRequestRepo.GetByIdPurchaseRequest(tx, id)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (p *PurchaseRequestServiceImpl) GetAllPurchaseRequestDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	//TODO implement me
	tx := p.DB.Begin()
	result, err := p.PurchaseRequestRepo.GetAllPurchaseRequestDetail(tx, filterCondition, pages)
	defer helper.CommitOrRollbackTrx(tx)

	if err != nil {
		return result, err
	}
	return result, nil
}

func (p *PurchaseRequestServiceImpl) GetByIdPurchaseRequestDetail(id int) (transactionsparepartpayloads.PurchaseRequestDetailResponsesPayloads, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	result, err := p.PurchaseRequestRepo.GetByIdPurchaseRequestDetail(tx, id)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (p *PurchaseRequestServiceImpl) NewPurchaseRequestHeader(request transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest) (transactionsparepartentities.PurchaseRequestEntities, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	result, err := p.PurchaseRequestRepo.NewPurchaseRequestHeader(tx, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (p *PurchaseRequestServiceImpl) NewPurchaseRequestDetail(payloads transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads) (transactionsparepartentities.PurchaseRequestDetail, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	result, err := p.PurchaseRequestRepo.NewPurchaseRequestDetail(tx, payloads)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (p *PurchaseRequestServiceImpl) SavePurchaseRequestUpdateHeader(request transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, id int) (transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, *exceptions.BaseErrorResponse) {
	//TODO implement me
	tx := p.DB.Begin()
	res, err := p.PurchaseRequestRepo.SavePurchaseRequestHeader(tx, request, id)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return res, err
	}
	return res, nil
}

func (p *PurchaseRequestServiceImpl) SavePurchaseRequestUpdateDetail(payloads transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, id int) (transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	res, err := p.PurchaseRequestRepo.SavePurchaseRequestDetail(tx, payloads, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (p *PurchaseRequestServiceImpl) VoidPurchaseRequest(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	res, err := p.PurchaseRequestRepo.VoidPurchaseRequest(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return res, err
	}
	return res, nil
}
func (p *PurchaseRequestServiceImpl) InsertPurchaseRequestUpdateHeader(request transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, id int) (transactionsparepartpayloads.PurchaseRequestGetByIdResponses, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	res, err := p.PurchaseRequestRepo.InsertPurchaseRequestHeader(tx, request, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (p *PurchaseRequestServiceImpl) InsertPurchaseRequestUpdateDetail(payloads transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, id int) (transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	res, err := p.PurchaseRequestRepo.SavePurchaseRequestDetail(tx, payloads, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return res, err
	}
	return res, nil
}
