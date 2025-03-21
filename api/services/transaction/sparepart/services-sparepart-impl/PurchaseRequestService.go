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
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PurchaseRequestServiceImpl struct {
	PurchaseRequestRepo transactionsparepartrepository.PurchaseRequestRepository
	DB                  *gorm.DB
	RedisClient         *redis.Client
}

func (p *PurchaseRequestServiceImpl) GetAllItemTypePurchaseRequest(filterCondition []utils.FilterCondition, pages pagination.Pagination, companyId int) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	result, err := p.PurchaseRequestRepo.GetAllItemTypePrRequest(tx, filterCondition, pages, companyId)
	if err != nil {
		return result, err
	}
	return result, nil
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
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	result, err := p.PurchaseRequestRepo.GetAllPurchaseRequest(tx, filterCondition, pages, Dateparams)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (p *PurchaseRequestServiceImpl) GetByIdPurchaseRequest(id int) (transactionsparepartpayloads.PurchaseRequestGetByIdResponses, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
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
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	result, err := p.PurchaseRequestRepo.GetAllPurchaseRequestDetail(tx, filterCondition, pages)
	defer helper.CommitOrRollbackTrx(tx)

	if err != nil {
		return result, err
	}
	return result, nil
}

func (p *PurchaseRequestServiceImpl) GetByIdPurchaseRequestDetail(id int) (transactionsparepartpayloads.PurchaseRequestDetailResponsesPayloads, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	result, err := p.PurchaseRequestRepo.GetByIdPurchaseRequestDetail(tx, id)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (p *PurchaseRequestServiceImpl) NewPurchaseRequestHeader(request transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest) (transactionsparepartentities.PurchaseRequestEntities, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	result, err := p.PurchaseRequestRepo.NewPurchaseRequestHeader(tx, request)

	if err != nil {
		return result, err
	}
	return result, nil
}

func (p *PurchaseRequestServiceImpl) NewPurchaseRequestDetail(payloads transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads) (transactionsparepartentities.PurchaseRequestDetail, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	result, err := p.PurchaseRequestRepo.NewPurchaseRequestDetail(tx, payloads)

	if err != nil {
		return result, err
	}
	return result, nil
}

func (p *PurchaseRequestServiceImpl) SavePurchaseRequestUpdateHeader(request transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, id int) (transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, *exceptions.BaseErrorResponse) {
	//TODO implement me
	tx := p.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	res, err := p.PurchaseRequestRepo.SavePurchaseRequestHeader(tx, request, id)

	if err != nil {
		return res, err
	}
	return res, nil
}

func (p *PurchaseRequestServiceImpl) SavePurchaseRequestUpdateDetail(payloads transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, id int) (transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	res, err := p.PurchaseRequestRepo.SavePurchaseRequestDetail(tx, payloads, id)

	if err != nil {
		return res, err
	}
	return res, nil
}

func (p *PurchaseRequestServiceImpl) VoidPurchaseRequest(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	defer helper.CommitOrRollbackTrx(tx)
	res, err := p.PurchaseRequestRepo.VoidPurchaseRequest(tx, id)

	if err != nil {
		return res, err
	}
	return res, nil
}
func (p *PurchaseRequestServiceImpl) SubmitPurchaseRequest(request transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, id int) (transactionsparepartpayloads.PurchaseRequestGetByIdResponses, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	defer helper.CommitOrRollbackTrx(tx)
	res, err := p.PurchaseRequestRepo.SubmitPurchaseRequest(tx, request, id)

	if err != nil {
		return res, err
	}
	return res, nil
}

func (p *PurchaseRequestServiceImpl) InsertPurchaseRequestUpdateDetail(payloads transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, id int) (transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	defer helper.CommitOrRollbackTrx(tx)
	res, err := p.PurchaseRequestRepo.SavePurchaseRequestDetail(tx, payloads, id)

	if err != nil {
		return res, err
	}
	return res, nil
}
func (p *PurchaseRequestServiceImpl) GetByIdItemTypePurchaseRequest(companyId int, id int) (transactionsparepartpayloads.PurchaseRequestItemGetAll, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	defer helper.CommitOrRollbackTrx(tx)
	res, err := p.PurchaseRequestRepo.GetByIdPurchaseRequestItemPr(tx, companyId, id)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (p *PurchaseRequestServiceImpl) GetByCodeItemTypePurchaseRequest(companyId int, itemCode string, brandId int) (transactionsparepartpayloads.PurchaseRequestItemGetAll, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	defer helper.CommitOrRollbackTrx(tx)
	res, err := p.PurchaseRequestRepo.GetByCodePurchaseRequestItemPr(tx, companyId, itemCode, brandId)
	if err != nil {
		return res, err
	}
	return res, nil
}
func (p *PurchaseRequestServiceImpl) VoidPurchaseRequestDetail(stringId string) (bool, *exceptions.BaseErrorResponse) {
	tx := p.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	defer helper.CommitOrRollbackTrx(tx)
	res, err := p.PurchaseRequestRepo.VoidPurchaseRequestDetailMultiId(tx, stringId)

	if err != nil {
		return res, err
	}
	return res, nil
}
