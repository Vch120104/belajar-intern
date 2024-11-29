package transactionworkshopserviceimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ServiceRequestServiceImpl struct {
	ServiceRequestRepository transactionworkshoprepository.ServiceRequestRepository
	DB                       *gorm.DB
	RedisClient              *redis.Client // Redis client
}

func OpenServiceRequestServiceImpl(ServiceRequestRepo transactionworkshoprepository.ServiceRequestRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.ServiceRequestService {
	return &ServiceRequestServiceImpl{
		ServiceRequestRepository: ServiceRequestRepo,
		DB:                       db,
		RedisClient:              redisClient,
	}
}

// Function to generate document service number
func (s *ServiceRequestServiceImpl) GenerateDocumentNumberServiceRequest(ServiceRequestId int) (string, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	documentNumber, err := s.ServiceRequestRepository.GenerateDocumentNumberServiceRequest(tx, ServiceRequestId)
	if err != nil {
		return "", err
	}
	return documentNumber, nil
}

func (s *ServiceRequestServiceImpl) NewStatus(filter []utils.FilterCondition) ([]transactionworkshopentities.ServiceRequestMasterStatus, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	statuses, err := s.ServiceRequestRepository.NewStatus(tx, filter)
	if err != nil {
		return []transactionworkshopentities.ServiceRequestMasterStatus{}, err
	}
	return statuses, nil
}

func (s *ServiceRequestServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	results, totalPages, totalRows, repoErr := s.ServiceRequestRepository.GetAll(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	return results, totalPages, totalRows, nil
}

func (s *ServiceRequestServiceImpl) GetById(id int, pages pagination.Pagination) (transactionworkshoppayloads.ServiceRequestResponse, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	result, repoErr := s.ServiceRequestRepository.GetById(tx, id, pages)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}

func (s *ServiceRequestServiceImpl) New(request transactionworkshoppayloads.ServiceRequestSaveRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	save, err := s.ServiceRequestRepository.New(tx, request)
	if err != nil {
		return transactionworkshopentities.ServiceRequest{}, err
	}

	return save, nil
}

func (s *ServiceRequestServiceImpl) Save(id int, request transactionworkshoppayloads.ServiceRequestSaveDataRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	save, err := s.ServiceRequestRepository.Save(tx, id, request)
	if err != nil {
		return transactionworkshopentities.ServiceRequest{}, err
	}

	return save, nil
}

func (s *ServiceRequestServiceImpl) Submit(id int) (bool, string, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	submit, newDocumentNumber, err := s.ServiceRequestRepository.Submit(tx, id)
	if err != nil {
		return false, "", err
	}

	return submit, newDocumentNumber, nil
}

func (s *ServiceRequestServiceImpl) Void(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	void, err := s.ServiceRequestRepository.Void(tx, id)
	if err != nil {
		return false, err
	}
	return void, nil
}

func (s *ServiceRequestServiceImpl) CloseOrder(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	closeOrder, err := s.ServiceRequestRepository.CloseOrder(tx, id)
	if err != nil {
		return false, err
	}
	return closeOrder, nil
}

func (s *ServiceRequestServiceImpl) GetAllServiceDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	results, totalPages, totalRows, repoErr := s.ServiceRequestRepository.GetAllServiceDetail(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	return results, totalPages, totalRows, nil

}

func (s *ServiceRequestServiceImpl) GetServiceDetailById(idsys int) (transactionworkshoppayloads.ServiceDetailResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	results, repoErr := s.ServiceRequestRepository.GetServiceDetailById(tx, idsys)
	if repoErr != nil {
		if repoErr.StatusCode == http.StatusNotFound {
			return transactionworkshoppayloads.ServiceDetailResponse{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Message: "Data not found"}
		}
		return transactionworkshoppayloads.ServiceDetailResponse{}, repoErr
	}

	return results, nil

}

func (s *ServiceRequestServiceImpl) AddServiceDetail(idsys int, request transactionworkshoppayloads.ServiceDetailSaveRequest) (transactionworkshopentities.ServiceRequestDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	success, err := s.ServiceRequestRepository.AddServiceDetail(tx, idsys, request)
	if err != nil {
		return transactionworkshopentities.ServiceRequestDetail{}, err
	}

	return success, nil
}

func (s *ServiceRequestServiceImpl) UpdateServiceDetail(idsys int, idservice int, request transactionworkshoppayloads.ServiceDetailUpdateRequest) (transactionworkshopentities.ServiceRequestDetail, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	save, err := s.ServiceRequestRepository.UpdateServiceDetail(tx, idsys, idservice, request)
	if err != nil {
		return transactionworkshopentities.ServiceRequestDetail{}, err
	}

	return save, nil
}

func (s *ServiceRequestServiceImpl) DeleteServiceDetail(idsys int, idservice int) (bool, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	delete, err := s.ServiceRequestRepository.DeleteServiceDetail(tx, idsys, idservice)
	if err != nil {
		return false, err
	}

	return delete, nil
}

func (s *ServiceRequestServiceImpl) DeleteServiceDetailMultiId(idsys int, idservice []int) (bool, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	deletemultiid, err := s.ServiceRequestRepository.DeleteServiceDetailMultiId(tx, idsys, idservice)
	if err != nil {
		return false, err
	}

	return deletemultiid, nil
}

func (s *ServiceRequestServiceImpl) NewServiceType(filter []utils.FilterCondition) ([]transactionworkshopentities.ServiceRequestMasterServiceType, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	types, err := s.ServiceRequestRepository.NewServiceType(tx, filter)
	if err != nil {
		return []transactionworkshopentities.ServiceRequestMasterServiceType{}, err
	}
	return types, nil
}
