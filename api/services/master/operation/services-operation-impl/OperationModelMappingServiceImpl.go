package masteroperationserviceimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OperationModelMappingServiceImpl struct {
	operationModelMappingRepo masteroperationrepository.OperationModelMappingRepository
	DB                        *gorm.DB
	RedisClient               *redis.Client // Redis client
}

func StartOperationModelMappingService(operationModelMappingRepo masteroperationrepository.OperationModelMappingRepository, db *gorm.DB, redisClient *redis.Client) masteroperationservice.OperationModelMappingService {
	return &OperationModelMappingServiceImpl{
		operationModelMappingRepo: operationModelMappingRepo,
		DB:                        db,
		RedisClient:               redisClient,
	}
}

func (s *OperationModelMappingServiceImpl) GetOperationModelMappingById(id int) (masteroperationpayloads.OperationModelMappingResponse, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.GetOperationModelMappingById(tx, id)
	if err != nil {
		return masteroperationpayloads.OperationModelMappingResponse{}, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetOperationModelMappingLookup(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.GetOperationModelMappingLookup(tx, filterCondition, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetOperationModelMappingByBrandModelOperationCode(request masteroperationpayloads.OperationModelModelBrandOperationCodeRequest) (masteroperationpayloads.OperationModelMappingResponse, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.GetOperationModelMappingByBrandModelOperationCode(tx, request)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) SaveOperationModelMapping(req masteroperationpayloads.OperationModelMappingResponse) (bool, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.SaveOperationModelMapping(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) ChangeStatusOperationModelMapping(Id int) (bool, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.ChangeStatusOperationModelMapping(tx, Id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) SaveOperationModelMappingFrt(request masteroperationpayloads.OperationModelMappingFrtRequest) (bool, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.SaveOperationModelMappingFrt(tx, request)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) DeleteOperationLevel(ids []int) (bool, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	result, err := s.operationModelMappingRepo.DeleteOperationLevel(tx, ids)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *OperationModelMappingServiceImpl) DeactivateOperationFrt(id string) (bool, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.DeactivateOperationFrt(tx, id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) ActivateOperationFrt(id string) (bool, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.ActivateOperationFrt(tx, id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetAllOperationDocumentRequirement(id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.GetAllOperationDocumentRequirement(tx, id, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetAllOperationFrt(id int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, totalPages, totalRows, err := s.operationModelMappingRepo.GetAllOperationFrt(tx, id, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *OperationModelMappingServiceImpl) GetOperationDocumentRequirementById(id int) (masteroperationpayloads.OperationModelMappingDocumentRequirementRequest, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.GetOperationDocumentRequirementById(tx, id)
	if err != nil {
		return masteroperationpayloads.OperationModelMappingDocumentRequirementRequest{}, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetOperationFrtById(id int) (masteroperationpayloads.OperationModelMappingFrtRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.GetOperationFrtById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteroperationpayloads.OperationModelMappingFrtRequest{}, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) SaveOperationModelMappingDocumentRequirement(request masteroperationpayloads.OperationModelMappingDocumentRequirementRequest) (bool, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.SaveOperationModelMappingDocumentRequirement(tx, request)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) DeactivateOperationDocumentRequirement(id string) (bool, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.DeactivateOperationDocumentRequirement(tx, id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) ActivateOperationDocumentRequirement(id string) (bool, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.ActivateOperationDocumentRequirement(tx, id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) SaveOperationLevel(request masteroperationpayloads.OperationLevelRequest) (bool, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.SaveOperationLevel(tx, request)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetAllOperationLevel(id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.GetAllOperationLevel(tx, id, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetOperationLevelById(id int) (masteroperationpayloads.OperationLevelByIdResponse, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.GetOperationLevelById(tx, id)
	if err != nil {
		return masteroperationpayloads.OperationLevelByIdResponse{}, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) DeactivateOperationLevel(id string) (bool, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.DeactivateOperationLevel(tx, id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) ActivateOperationLevel(id string) (bool, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.operationModelMappingRepo.ActivateOperationLevel(tx, id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) UpdateOperationModelMapping(operationModelMappingId int, request masteroperationpayloads.OperationModelMappingUpdate) (masteroperationentities.OperationModelMapping, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()

	update, err := s.operationModelMappingRepo.UpdateOperationModelMapping(tx, operationModelMappingId, request)

	if err != nil {
		return update, err
	}
	return update, nil
}

func (s *OperationModelMappingServiceImpl) UpdateOperationFrt(operationFrtId int, request masteroperationpayloads.OperationFrtUpdate) (masteroperationentities.OperationFrt, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()

	update, err := s.operationModelMappingRepo.UpdateOperationFrt(tx, operationFrtId, request)

	if err != nil {
		return update, err
	}
	return update, nil
}

func (s *OperationModelMappingServiceImpl) SaveOperationModelMappingAndFRT(requestHeader masteroperationpayloads.OperationModelMappingResponse, requestDetail masteroperationpayloads.OperationModelMappingFrtRequest) (bool, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	_, err = s.operationModelMappingRepo.SaveOperationModelMapping(tx, requestHeader)

	if err != nil {
		return false, err
	}

	latestID, errGet := s.operationModelMappingRepo.GetOperationModelMappingLatestId(tx)
	if errGet != nil {
		return false, errGet
	}

	requestDetail.OperationModelMappingId = latestID

	results, errDetail := s.operationModelMappingRepo.SaveOperationModelMappingFrt(tx, requestDetail)
	defer helper.CommitOrRollback(tx, err)
	if errDetail != nil {
		return false, errDetail
	}

	return results, nil
}

func (s *OperationModelMappingServiceImpl) CopyOperationModelMappingToOtherModel(headerId int, request masteroperationpayloads.OperationModelMappingCopyRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.CopyOperationModelMappingToOtherModel(tx, headerId, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}
