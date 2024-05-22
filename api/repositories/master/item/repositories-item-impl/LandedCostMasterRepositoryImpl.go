package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
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
	entities := []masteritementities.LandedCost{}
	payloads := []masteritempayloads.LandedCostMasterPayloads{}

	baseModelQuery := tx.Model(&entities)
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, baseModelQuery)).Scan(&payloads).Rows()

	if len(payloads) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	pages.Rows = payloads

	return pages, nil
}

func (r *LandedCostMasterRepositoryImpl) GetByIdLandedCost(tx *gorm.DB, id int) (masteritempayloads.LandedCostMasterPayloads, *exceptions.BaseErrorResponse) {
	entities := masteritementities.LandedCost{}
	payloads := masteritempayloads.LandedCostMasterPayloads{}
	rows, err := tx.Model(&entities).Where(masteritementities.LandedCost{LandedCostId: id}).First(&payloads).Rows()
	if err != nil {
		return payloads, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()
	return payloads, nil
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

func (r *LandedCostMasterRepositoryImpl) SaveLandedCost(tx *gorm.DB, req masteritempayloads.LandedCostMasterPayloads) (bool, *exceptions.BaseErrorResponse) {
	entitiesGet := masteritementities.LandedCost{}
	rows, _ := tx.Model(&entitiesGet).Where(masteritementities.LandedCost{CompanyId: req.CompanyId, SupplierId: req.SupplierId, ShippingMethodId: req.ShippingMethodId, LandedCostTypeId: req.LandedCostTypeId, LandedCostDescription: req.LandedCostDescription}).Rows()
	defer rows.Close()

	if !rows.Next() {
		entities := masteritementities.LandedCost{
			IsActive:               req.IsActive,
			CompanyId:              req.CompanyId,
			SupplierId:             req.SupplierId,
			ShippingMethodId:       req.ShippingMethodId,
			LandedCostTypeId:       req.LandedCostTypeId,
			LandedCostId:           req.LandedCostId,
			LandedCostDescription:  req.LandedCostDescription,
			LandedCostMasterFactor: req.LandedCostMasterFactor,
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
