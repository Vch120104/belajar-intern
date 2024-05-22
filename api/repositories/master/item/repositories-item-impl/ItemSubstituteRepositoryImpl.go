package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type ItemSubstituteRepositoryImpl struct {
}

func StartItemSubstituteRepositoryImpl() masteritemrepository.ItemSubstituteRepository {
	return &ItemSubstituteRepositoryImpl{}
}

func (r *ItemSubstituteRepositoryImpl) GetAllItemSubstitute(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := []masteritementities.ItemSubstitute{}
	payloads := []masteritempayloads.ItemSubstitutePayloads{}
	tableStruct := masteritempayloads.ItemSubstitutePayloads{}
	baseModelQuery := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&payloads).Rows()
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(payloads) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	defer rows.Close()

	pages.Rows = payloads

	return pages, nil
}

func (r *ItemSubstituteRepositoryImpl) GetByIdItemSubstitute(tx *gorm.DB, id int) (masteritempayloads.ItemSubstitutePayloads, *exceptions.BaseErrorResponse) {
	response := masteritempayloads.ItemSubstitutePayloads{}
	tableStruct := masteritempayloads.ItemSubstitutePayloads{}
	baseModelQuery := utils.CreateJoinSelectStatement(tx, tableStruct).Where(masteritementities.ItemSubstitute{ItemSubstituteId: id})

	rows, err := baseModelQuery.First(&response).Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if response != (masteritempayloads.ItemSubstitutePayloads{}) {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	defer rows.Close()
	return response, nil
}

func (r *ItemSubstituteRepositoryImpl) GetAllItemSubstituteDetail(tx *gorm.DB, pages pagination.Pagination, id int) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := []masteritementities.ItemSubstituteDetail{}
	payloads := []masteritempayloads.ItemSubstituteDetailPayloads{}
	tableStruct := masteritempayloads.ItemSubstituteDetailPayloads{}

	baseModelQuery := utils.CreateJoinSelectStatement(tx, tableStruct).Where(masteritementities.ItemSubstituteDetail{ItemSubstituteId: id})

	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, baseModelQuery)).Scan(&payloads).Rows()
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	if len(payloads) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	defer rows.Close()

	pages.Rows = payloads

	return pages, nil
}

func (r *ItemSubstituteRepositoryImpl) GetByIdItemSubstituteDetail(tx *gorm.DB, id int) (masteritempayloads.ItemSubstituteDetailGetPayloads, *exceptions.BaseErrorResponse) {
	response := masteritempayloads.ItemSubstituteDetailGetPayloads{}
	tableStruct := masteritempayloads.ItemSubstituteDetailPayloads{}
	baseModelQuery := utils.CreateJoinSelectStatement(tx, tableStruct).Where(masteritementities.ItemSubstituteDetail{ItemSubstituteId: id})

	rows, err := baseModelQuery.First(&response).Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	if response != (masteritempayloads.ItemSubstituteDetailGetPayloads{}) {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	defer rows.Close()
	return response, nil
}

func (r *ItemSubstituteRepositoryImpl) SaveItemSubstitute(tx *gorm.DB, req masteritempayloads.ItemSubstitutePostPayloads) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemSubstitute{
		SubstituteTypeCode: req.SubstituteTypeCode,
		ItemSubstituteId:   req.ItemSubstituteId,
		EffectiveDate:      req.EffectiveDate,
		ItemId:             req.ItemId,
	}
	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	return true, nil
}

func (r *ItemSubstituteRepositoryImpl) SaveItemSubstituteDetail(tx *gorm.DB, req masteritempayloads.ItemSubstituteDetailPostPayloads, id int) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemSubstituteDetail{
		ItemSubstituteDetailId: req.ItemSubstituteDetailId,
		ItemId:                 req.ItemId,
		ItemSubstituteId:       id,
		Quantity:               req.Quantity,
		Sequence:               req.Sequence,
	}
	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	return true, nil
}

func (r *ItemSubstituteRepositoryImpl) ChangeStatusItemOperation(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.ItemSubstitute

	result := tx.Model(&entities).
		Where("item_substitute_id = ?", id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
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
			StatusCode: http.StatusNotFound,
			Err:        result.Error,
		}
	}

	return true, nil
}

func (r *ItemSubstituteRepositoryImpl) DeactivateItemSubstituteDetail(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteritementities.ItemSubstituteDetail
		err := tx.Model(&entityToUpdate).Where("item_substitute_detail_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = false
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}

func (r *ItemSubstituteRepositoryImpl) ActivateItemSubstituteDetail(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteritementities.ItemSubstituteDetail
		err := tx.Model(&entityToUpdate).Where("item_substitute_detail_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = true
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}
