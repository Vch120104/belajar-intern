package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	aftersalesserviceapiutils "after-sales/api/utils/aftersales-service"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"net/http"

	"gorm.io/gorm"
)

type ForecastMasterRepositoryImpl struct {
}

func StartForecastMasterRepositoryImpl() masterrepository.ForecastMasterRepository {
	return &ForecastMasterRepositoryImpl{}
}

func (r *ForecastMasterRepositoryImpl) GetForecastMasterById(tx *gorm.DB, forecastMasterId int) (masterpayloads.ForecastMasterResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.ForecastMaster{}
	response := masterpayloads.ForecastMasterResponse{}

	err := tx.Model(&entities).
		Where(masterentities.ForecastMaster{
			ForecastMasterId: forecastMasterId,
		}).
		First(&response).
		Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

func (r *ForecastMasterRepositoryImpl) SaveForecastMaster(tx *gorm.DB, req masterpayloads.ForecastMasterResponse) (masterentities.ForecastMaster, *exceptions.BaseErrorResponse) {
	entities := masterentities.ForecastMaster{
		SupplierId:                 req.SupplierId,
		CompanyId:                  req.CompanyId,
		MovingCodeId:               req.MovingCodeId,
		OrderTypeId:                req.OrderTypeId,
		ForecastMasterLeadTime:     req.ForecastMasterLeadTime,
		ForecastMasterSafetyFactor: req.ForecastMasterSafetyFactor,
		ForecastMasterOrderCycle:   req.ForecastMasterOrderCycle,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return masterentities.ForecastMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *ForecastMasterRepositoryImpl) ChangeStatusForecastMaster(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masterentities.ForecastMaster

	result := tx.Model(&entities).
		Where("forecast_master_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	if entities.IsActive {
		entities.IsActive = false
	} else {
		entities.IsActive = true
	}

	result = tx.Save(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}

func (r *ForecastMasterRepositoryImpl) GetAllForecastMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var entities []masterentities.ForecastMaster

	baseModelQuery := tx.Model(&masterentities.ForecastMaster{})

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&entities).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(entities) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, entity := range entities {
		// Fetch supplier data
		getSupplierResponse, supplierErr := generalserviceapiutils.GetSupplierMasterByID(entity.SupplierId)
		if supplierErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        supplierErr.Err,
			}
		}

		// Fetch order type
		getOrderTypeResponse, orderTypeErr := aftersalesserviceapiutils.GetOrderTypeById(entity.OrderTypeId)
		if orderTypeErr != nil {
			return pages, orderTypeErr
		}

		// Fetch moving code
		getMovingCodeResponse, movingCodeErr := aftersalesserviceapiutils.GetMovingCodeById(entity.MovingCodeId)
		if movingCodeErr != nil {
			return pages, movingCodeErr
		}

		result := map[string]interface{}{
			"forecast_master_id":            entity.ForecastMasterId,
			"is_active":                     entity.IsActive,
			"company_id":                    entity.CompanyId,
			"supplier_id":                   entity.SupplierId,
			"supplier_name":                 getSupplierResponse.SupplierName,
			"moving_code_id":                entity.MovingCodeId,
			"moving_code_description":       getMovingCodeResponse.MovingCodeName,
			"order_type_id":                 entity.OrderTypeId,
			"order_type_name":               getOrderTypeResponse.OrderTypeName,
			"forecast_master_lead_time":     entity.ForecastMasterLeadTime,
			"forecast_master_safety_factor": entity.ForecastMasterSafetyFactor,
			"forecast_master_order_cycle":   entity.ForecastMasterOrderCycle,
		}

		results = append(results, result)
	}

	pages.Rows = results
	return pages, nil
}

func (r *ForecastMasterRepositoryImpl) UpdateForecastMaster(tx *gorm.DB, req masterpayloads.ForecastMasterResponse, id int) (masterentities.ForecastMaster, *exceptions.BaseErrorResponse) {
	var entity masterentities.ForecastMaster

	err := tx.Model(&entity).Where("forecast_master_id = ?", id).First(&entity).Updates(req)
	if err.Error != nil {
		return masterentities.ForecastMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err.Error,
		}
	}
	return entity, nil
}
