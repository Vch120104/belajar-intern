package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type LandedCostMasterRepositoryImpl struct {
}

func StartLandedCostMasterRepositoryImpl() masteritemrepository.LandedCostMasterRepository {
	return &LandedCostMasterRepositoryImpl{}
}

func (r *LandedCostMasterRepositoryImpl) GetAllLandedCost(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities masteritementities.LandedCost
	var payloads []masteritempayloads.LandedCostMasterPayloads

	query := tx.Model(&entities).
		Select("mtr_landed_cost.*, mtr_shipping_method.shipping_method_description, mtr_landed_cost_type.landed_cost_type_name").
		Joins("LEFT JOIN dms_microservices_general_dev.dbo.mtr_shipping_method ON mtr_landed_cost.shipping_method_id = mtr_shipping_method.shipping_method_id").
		Joins("LEFT JOIN dms_microservices_general_dev.dbo.mtr_landed_cost_type ON mtr_landed_cost.landed_cost_type_id = mtr_landed_cost_type.landed_cost_type_id")

	query = utils.ApplyFilter(query, filterCondition)
	query = query.Scopes(pagination.Paginate(&pages, query))

	err := query.Scan(&payloads).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(payloads) == 0 {
		pages.Rows = []map[string]interface{}{}
		pages.TotalRows = 0
		pages.TotalPages = 0
		return pages, nil
	}

	result := []map[string]interface{}{}
	for _, payload := range payloads {
		// Fetch Shipping Method by ID
		shippingMethod, errSM := generalserviceapiutils.GetShippingMethodById(payload.ShippingMethodId)
		if errSM != nil {
			return pages, errSM
		}

		// Fetch Landed Cost Type by ID
		landedCostType, errLC := generalserviceapiutils.GetLandedCostTypeById(payload.LandedCostTypeId)
		if errLC != nil {
			return pages, errLC
		}

		// Fetch Supplier by ID
		supplier, errSup := generalserviceapiutils.GetSupplierMasterByID(payload.SupplierId)
		if errSup != nil {
			return pages, errSup
		}

		// Fetch Company by ID
		company, errComp := generalserviceapiutils.GetCompanyDataById(payload.CompanyId)
		if errComp != nil {
			return pages, errComp
		}

		// Combine data into final result
		data := map[string]interface{}{
			"landed_cost_id":        payload.LandedCostId,
			"shipping_method_name":  shippingMethod.ShippingMethodName,
			"landed_cost_type_name": landedCostType.LandedCostTypeName,
			"landed_cost_type_code": landedCostType.LandedCostTypeCode,
			"is_active":             payload.IsActive,
			"landed_cost_factor":    payload.LandedCostFactor,
			"company_id":            payload.CompanyId,
			"company_name":          company.CompanyName,
			"supplier_id":           payload.SupplierId,
			"supplier_code":         supplier.SupplierCode,
			"supplier_name":         supplier.SupplierName,
			"shipping_method_id":    payload.ShippingMethodId,
			"landed_cost_type_id":   payload.LandedCostTypeId,
		}
		result = append(result, data)
	}

	pages.Rows = result
	return pages, nil
}

func (r *LandedCostMasterRepositoryImpl) GetByIdLandedCost(tx *gorm.DB, id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tableStruct := masteritementities.LandedCost{}
	response := masteritempayloads.LandedCostMasterPayloads{}

	baseModelQuery := tx.Model(tableStruct)

	err := baseModelQuery.Where("landed_cost_id = ?", id).First(&response).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	// Fetch LandedCostType
	GetLandedCostType, errLandedCost := generalserviceapiutils.GetLandedCostTypeById(response.LandedCostTypeId)
	if errLandedCost != nil {
		return nil, errLandedCost
	}

	// Fetch ShippingMethod
	GetShippingMethodType, errShippingMethod := generalserviceapiutils.GetShippingMethodById(response.ShippingMethodId)
	if errShippingMethod != nil {
		return nil, errShippingMethod
	}

	result := map[string]interface{}{
		"company_id":                   response.CompanyId,
		"is_active":                    response.IsActive,
		"landed_cost_factor":           response.LandedCostFactor,
		"landed_cost_id":               response.LandedCostId,
		"landed_cost_type_code":        GetLandedCostType.LandedCostTypeCode,
		"landed_cost_type_description": GetLandedCostType.LandedCostTypeDescription,
		"landed_cost_type_id":          GetLandedCostType.LandedCostTypeId,
		"shipping_method_code":         GetShippingMethodType.ShippingMethodCode,
		"shipping_method_id":           GetShippingMethodType.ShippingMethodId,
		"supplier_id":                  response.SupplierId,
	}

	return result, nil
}

// func (r *LandedCostMasterRepositoryImpl) SaveLandedCost(tx *gorm.DB, req masteritempayloads.LandedCostMasterPayloads) (bool, error) {
// 	entitiesGet := masteritementities.LandedCost{}
// 	rows, _:= tx.Model(&entitiesGet).Where(masteritementities.LandedCost{CompanyId: req.CompanyId,SupplierId: req.SupplierId,ShippingMethodId: req.ShippingMethodId,LandedCostTypeId: req.LandedCostTypeId,LandedCostDescription: req.LandedCostDescription}).Rows()
// 	if rows==nil{
// 		entities := masteritementities.LandedCost{
// 			IsActive:                req.IsActive,
// 			CompanyId:               req.CompanyId,
// 			SupplierId:              req.SupplierId,
// 			ShippingMethodId:        req.ShippingMethodId,
// 			LandedCostTypeId:        req.LandedCostTypeId,
// 			LandedCostId:            req.LandedCostId,
// 			LandedCostDescription:   req.LandedCostDescription,
// 			LandedCostMasterFactor: req.LandedCostMasterFactor,
// 		}
// 		err := tx.Save(&entities).Error
// 		if err != nil {
// 			return false, err
// 		}
// 		return true, nil
// 	}
// 	return true,nil
// }

func (r *LandedCostMasterRepositoryImpl) SaveLandedCost(tx *gorm.DB, req masteritempayloads.LandedCostMasterRequest) (masteritementities.LandedCost, *exceptions.BaseErrorResponse) {
	var existingLandedCost masteritementities.LandedCost

	err := tx.Model(masteritementities.LandedCost{}).
		Where(map[string]interface{}{
			"company_id":          req.CompanyId,
			"supplier_id":         req.SupplierId,
			"shipping_method_id":  req.ShippingMethodId,
			"landed_cost_type_id": req.LandedCostTypeId,
		}).
		First(&existingLandedCost).Error

	if err == nil {
		return existingLandedCost, nil
	}

	if err != gorm.ErrRecordNotFound {
		return masteritementities.LandedCost{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}

	newLandedCost := masteritementities.LandedCost{
		IsActive:         req.IsActive,
		CompanyId:        req.CompanyId,
		SupplierId:       req.SupplierId,
		ShippingMethodId: req.ShippingMethodId,
		LandedCostTypeId: req.LandedCostTypeId,
		LandedCostId:     req.LandedCostId,
		LandedCostfactor: req.LandedCostFactor,
	}

	err = tx.Save(&newLandedCost).Error
	if err != nil {
		return masteritementities.LandedCost{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Return the newly saved entity with no error
	return newLandedCost, nil
}

func (r *LandedCostMasterRepositoryImpl) DeactivateLandedCostmaster(tx *gorm.DB, id string) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")
	var results []map[string]interface{}

	for _, idStr := range idSlice {
		var entityToUpdate masteritementities.LandedCost // Create a new instance for each loop iteration
		err := tx.Model(&entityToUpdate).Where("landed_cost_id = ?", idStr).First(&entityToUpdate).Error
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = false
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        result.Error,
			}
		}

		// Append the updated entity details to the results slice
		results = append(results, map[string]interface{}{
			"is_active":      entityToUpdate.IsActive,
			"landed_cost_id": entityToUpdate.LandedCostId,
		})
	}

	return results, nil
}

func (r *LandedCostMasterRepositoryImpl) ActivateLandedCostMaster(tx *gorm.DB, id string) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")
	var results []map[string]interface{}

	for _, idStr := range idSlice {
		var entityToUpdate masteritementities.LandedCost // Create a new instance for each loop iteration
		err := tx.Model(&entityToUpdate).Where("landed_cost_id = ?", idStr).First(&entityToUpdate).Error
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = true
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        result.Error,
			}
		}

		// Append the updated entity details to the results slice
		results = append(results, map[string]interface{}{
			"is_active":      entityToUpdate.IsActive,
			"landed_cost_id": entityToUpdate.LandedCostId,
		})
	}

	return results, nil
}

func (r *LandedCostMasterRepositoryImpl) UpdateLandedCostMaster(tx *gorm.DB, id int, req masteritempayloads.LandedCostMasterUpdateRequest) (masteritementities.LandedCost, *exceptions.BaseErrorResponse) {
	var entities masteritementities.LandedCost

	// Update the landed_cost_factor
	result := tx.Model(&entities).Where("landed_cost_id = ?", id).Update("landed_cost_factor", req.LandedCostfactor)
	if result.Error != nil {
		return masteritementities.LandedCost{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        result.Error,
		}
	}

	// Fetch the updated entity
	err := tx.Model(&entities).Where("landed_cost_id = ?", id).First(&entities).Error
	if err != nil {
		return masteritementities.LandedCost{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	return entities, nil
}
