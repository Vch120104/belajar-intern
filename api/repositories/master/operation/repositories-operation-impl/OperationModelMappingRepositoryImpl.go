package masteroperationrepositoryimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	"time"

	masteroperationpayloads "after-sales/api/payloads/master/operation"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	"after-sales/api/utils"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type OperationModelMappingRepositoryImpl struct {
	myDB *gorm.DB
}

func StartOperationModelMappingRepositoryImpl(db *gorm.DB) masteroperationrepository.OperationModelMappingRepository {
	return &OperationModelMappingRepositoryImpl{myDB: db}
}

func (r *OperationModelMappingRepositoryImpl) WithTrx(trxHandle *gorm.DB) masteroperationrepository.OperationModelMappingRepository {
	if trxHandle == nil {
		log.Println("Transaction Database Not Found!")
		return r
	}
	r.myDB = trxHandle
	return r
}

func (r *OperationModelMappingRepositoryImpl) GetOperationModelMappingById(Id int) (masteroperationpayloads.OperationModelMappingResponse, error) {
	entities := masteroperationentities.OperationModelMapping{}
	response := masteroperationpayloads.OperationModelMappingResponse{}

	rows, err := r.myDB.Model(&entities).
		Where(masteroperationentities.OperationModelMapping{
			OperationModelMappingId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *OperationModelMappingRepositoryImpl) GetOperationModelMappingByBrandModelOperationCode(request masteroperationpayloads.OperationModelModelBrandOperationCodeRequest) (masteroperationpayloads.OperationModelMappingResponse, error) {
	entities := masteroperationentities.OperationModelMapping{}
	response := masteroperationpayloads.OperationModelMappingResponse{}

	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	r.myDB.Logger = newLogger

	rows, err := r.myDB.Model(&entities).
		Where("brand_id = ? AND model_id = ? AND operation_id = ?", request.BrandId, request.ModelId, request.OperationId).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *OperationModelMappingRepositoryImpl) GetOperationModelMappingLookup(filterCondition []utils.FilterCondition) ([]map[string]interface{}, error) {
	var responses []masteroperationpayloads.OperationModelMappingResponse
	var getBrandResponse []masteroperationpayloads.BrandResponse
	var getModelResponse []masteroperationpayloads.ModelResponse
	var c *gin.Context

	// define table struct
	tableStruct := masteroperationpayloads.OperationModelMappingLookup{}

	//join table
	joinTable := utils.CreateJoinSelectStatement(r.myDB, tableStruct)

	//apply filter
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	//execute
	rows, err := whereQuery.Scan(&responses).Rows()

	brandUrl := "http://10.1.32.26:8000/sales-service/api/sales/unit-brand?page=0&limit=10"
	errUrlBrand := utils.Get(c, brandUrl, &getBrandResponse, nil)
	if errUrlBrand != nil {
		return nil, errUrlBrand
	}

	modelUrl := "http://10.1.32.26:8000/sales-service/api/sales/unit-model?page=0&limit=10"
	errUrlModel := utils.Get(c, modelUrl, &getModelResponse, nil)
	if errUrlModel != nil {
		return nil, errUrlBrand
	}

	joinedData := utils.DataFrameInnerJoin(responses, getBrandResponse, "BrandId")

	joinedDataSecond := utils.DataFrameInnerJoin(joinedData, getModelResponse, "ModelId")

	if err != nil {
		return joinedDataSecond, err
	}

	defer rows.Close()

	return joinedDataSecond, nil
}

func (r *OperationModelMappingRepositoryImpl) SaveOperationModelMapping(request masteroperationpayloads.OperationModelMappingResponse) (bool, error) {
	entities := masteroperationentities.OperationModelMapping{
		IsActive:                request.IsActive,
		OperationModelMappingId: request.OperationModelMappingId,
		BrandId:                 request.BrandId,
		ModelId:                 request.ModelId,
		OperationId:             request.OperationId,
		OperationUsingIncentive: request.OperationUsingIncentive,
		OperationUsingActual:    request.OperationUsingActual,
		OperationPdi:            request.OperationPdi,
	}

	err := r.myDB.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *OperationModelMappingRepositoryImpl) ChangeStatusOperationModelMapping(Id int) (bool, error) {
	var entities masteroperationentities.OperationModelMapping

	result := r.myDB.Model(&entities).
		Where("operation_model_mapping_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	// Toggle the IsActive value
	if entities.IsActive {
		entities.IsActive = false
	} else {
		entities.IsActive = true
	}

	result = r.myDB.Save(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}
