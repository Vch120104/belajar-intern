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

func (r *LandedCostMasterRepositoryImpl) GetAllLandedCost(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var shippingmethodpayloads []masteritempayloads.ShippingMethodResponse
	var landedcostpayloads []masteritempayloads.LandedCostTypeResponse
	var entities []masteritementities.LandedCost
	var payloads []masteritempayloads.LandedCostMasterPayloads

	baseModelQuery := tx.Model(&entities)
	rows := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, baseModelQuery))
	Where := utils.ApplyFilterExact(rows, filterCondition)
	final, err := Where.Scan(&payloads).Rows()

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer final.Close()
	if len(payloads) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	ShippingMethodUrl := config.EnvConfigs.GeneralServiceUrl + "shipping-method"

	errshippingmethod := utils.Get(ShippingMethodUrl, &shippingmethodpayloads, nil)

	if errshippingmethod != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	JoinedData1 := utils.DataFrameInnerJoin(payloads, shippingmethodpayloads, "ShippingMethodId")

	landedcosturl := config.EnvConfigs.GeneralServiceUrl + "landed-cost-type"

	errLandedCost := utils.Get(landedcosturl, &landedcostpayloads, nil)

	if errLandedCost != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	JoinedData2 := utils.DataFrameInnerJoin(JoinedData1, landedcostpayloads, "LandedCostTypeId")

	paginatedata, totalPages, totalrows := pagination.NewDataFramePaginate(JoinedData2, &pages)
	return paginatedata, totalPages, totalrows, nil
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
	joinedData1 := utils.DataFrameInnerJoin([]masteritempayloads.LandedCostMasterPayloads{response}, []masteritempayloads.LandedCostTypeResponse{GetLandedCostType}, "LandedCostTypeId")
	ShippingMethodUrl := config.EnvConfigs.GeneralServiceUrl + "shipping-method/" + strconv.Itoa(response.ShippingMethodId)

	errshippingmethod := utils.Get(ShippingMethodUrl, &GetShippinMethodType, nil)

	if errshippingmethod != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	JoinedData2 := utils.DataFrameInnerJoin(joinedData1, []masteritempayloads.ShippingMethodResponse{GetShippinMethodType}, "ShippingMethodId")

	return JoinedData2[0], nil

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

func (r *LandedCostMasterRepositoryImpl) SaveLandedCost(tx *gorm.DB, req masteritempayloads.LandedCostMasterRequest) (bool, *exceptions.BaseErrorResponse) {
	entitiesGet := masteritementities.LandedCost{}
	rows, _ := tx.Model(&entitiesGet).Where(masteritementities.LandedCost{CompanyId: req.CompanyId, SupplierId: req.SupplierId, ShippingMethodId: req.ShippingMethodId, LandedCostTypeId: req.LandedCostTypeId}).Rows()
	defer rows.Close()

	if !rows.Next() {
		entities := masteritementities.LandedCost{
			IsActive:         req.IsActive,
			CompanyId:        req.CompanyId,
			SupplierId:       req.SupplierId,
			ShippingMethodId: req.ShippingMethodId,
			LandedCostTypeId: req.LandedCostTypeId,
			LandedCostId:     req.LandedCostId,
			LandedCostfactor: req.LandedCostFactor,
		}
		err := tx.Save(&entities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		return true, nil
	}

	return true, nil
}

func (r *LandedCostMasterRepositoryImpl) DeactivateLandedCostmaster(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteritementities.LandedCost
		err := tx.Model(&entityToUpdate).Where("landed_cost_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = false
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}

func (r *LandedCostMasterRepositoryImpl) ActivateLandedCostMaster(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteritementities.LandedCost
		err := tx.Model(&entityToUpdate).Where("landed_cost_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = true
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}

func (r *LandedCostMasterRepositoryImpl) UpdateLandedCostMaster(tx *gorm.DB, id int, req masteritempayloads.LandedCostMasterUpdateRequest) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.LandedCost

	result := tx.Model(&entities).Where("landed_cost_id = ?", id).Update("landed_cost_factor", req.LandedCostfactor)
	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        result.Error,
		}
	}
	return true, nil
}
