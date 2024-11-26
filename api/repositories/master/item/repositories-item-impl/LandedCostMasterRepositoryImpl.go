package masteritemrepositoryimpl

import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"net/http"
	"strconv"
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
	var shippingmethodpayloads []masteritempayloads.ShippingMethodResponse
	var landedcostpayloads []masteritempayloads.LandedCostTypeResponse

	query := tx.Model(&entities).
		Select("mtr_landed_cost.*, shipping_method.shipping_method_name, landed_cost_type.landed_cost_type_name").
		Joins("LEFT JOIN shipping_method ON mtr_landed_cost.shipping_method_id = shipping_method.shipping_method_id").
		Joins("LEFT JOIN landed_cost_type ON mtr_landed_cost.landed_cost_type_id = landed_cost_type.landed_cost_type_id")

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

	errShippingMethod := utils.Get(config.EnvConfigs.GeneralServiceUrl+"shipping-method", &shippingmethodpayloads, nil)
	if errShippingMethod != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errShippingMethod,
		}
	}

	errLandedCostType := utils.Get(config.EnvConfigs.GeneralServiceUrl+"landed-cost-type", &landedcostpayloads, nil)
	if errLandedCostType != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errLandedCostType,
		}
	}

	joinedData := utils.DataFrameLeftJoin(payloads, shippingmethodpayloads, "ShippingMethodId")
	joinedData = utils.DataFrameLeftJoin(joinedData, landedcostpayloads, "LandedCostTypeId")

	result := []map[string]interface{}{}
	for _, res := range joinedData {
		data := map[string]interface{}{
			"landed_cost_id":        res["LandedCostId"],
			"shipping_method_name":  res["ShippingMethodName"],
			"landed_cost_type_name": res["LandedCostTypeName"],
			"is_active":             res["IsActive"],
			"cost_value":            res["CostValue"],
			"effective_date":        res["EffectiveDate"],
			"created_at":            res["CreatedAt"],
			"updated_at":            res["UpdatedAt"],
		}
		result = append(result, data)
	}

	pages.Rows = result
	return pages, nil
}

func (r *LandedCostMasterRepositoryImpl) GetByIdLandedCost(tx *gorm.DB, id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tableStruct := masteritementities.LandedCost{}
	response := masteritempayloads.LandedCostMasterPayloads{}
	var GetLandedCostType masteritempayloads.LandedCostTypeResponse
	var GetShippinMethodType masteritempayloads.ShippingMethodResponse

	baseModelQuery := tx.Model(tableStruct)

	err := baseModelQuery.Where("landed_cost_id = ?", id).First(&response).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	LandedCostTypeUrl := config.EnvConfigs.GeneralServiceUrl + "landed-cost-type/" + strconv.Itoa(response.LandedCostTypeId)
	errLandedCost := utils.Get(LandedCostTypeUrl, &GetLandedCostType, nil)
	if errLandedCost != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	joinedData1, errdf := utils.DataFrameInnerJoin([]masteritempayloads.LandedCostMasterPayloads{response}, []masteritempayloads.LandedCostTypeResponse{GetLandedCostType}, "LandedCostTypeId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	ShippingMethodUrl := config.EnvConfigs.GeneralServiceUrl + "shipping-method/" + strconv.Itoa(response.ShippingMethodId)

	errshippingmethod := utils.Get(ShippingMethodUrl, &GetShippinMethodType, nil)

	if errshippingmethod != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	JoinedData2, errdf := utils.DataFrameInnerJoin(joinedData1, []masteritempayloads.ShippingMethodResponse{GetShippinMethodType}, "ShippingMethodId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	result := map[string]interface{}{
		"company_id":                   JoinedData2[0]["CompanyId"],
		"is_active":                    JoinedData2[0]["IsActive"],
		"landed_cost_factor":           JoinedData2[0]["LandedCostFactor"],
		"landed_cost_id":               JoinedData2[0]["LandedCostTypeId"],
		"landed_cost_type_code":        JoinedData2[0]["LandedCostTypeCode"],
		"landed_cost_type_description": JoinedData2[0]["LandedCostTypeDescription"],
		"landed_cost_type_id":          JoinedData2[0]["LandedCostTypeId"],
		"shipping_method_code":         JoinedData2[0]["ShippingMethodCode"],
		"shipping_method_id":           JoinedData2[0]["ShippingMethodId"],
		"supplier_id":                  JoinedData2[0]["SupplierId"],
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
