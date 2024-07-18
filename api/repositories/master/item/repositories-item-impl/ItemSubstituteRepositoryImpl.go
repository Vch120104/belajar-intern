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
	"time"

	"gorm.io/gorm"
)

type ItemSubstituteRepositoryImpl struct {
}

func StartItemSubstituteRepositoryImpl() masteritemrepository.ItemSubstituteRepository {
	return &ItemSubstituteRepositoryImpl{}
}

func (r *ItemSubstituteRepositoryImpl) GetAllItemSubstitute(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination, from time.Time, to time.Time) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemSubstitute{}
	payloads := []masteritempayloads.ItemSubstitutePayloads{}

	query := tx.Model(entities).Select("mtr_item_substitute.*, Item.item_code, Item.item_name").
		Joins("Item", tx.Select(""))

	whereQuery := utils.ApplyFilter(query, filterCondition)

	if !from.IsZero() && !to.IsZero() {
		whereQuery.Where("effective_date BETWEEN ? AND ?", from, to)
	} else if !from.IsZero() {
		whereQuery.Where("effective_date >= ?", from)
	}

	err := whereQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&payloads).Error

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

	pages.Rows = payloads

	return pages, nil
}

func (r *ItemSubstituteRepositoryImpl) GetByIdItemSubstitute(tx *gorm.DB, id int) (masteritempayloads.ItemSubstitutePayloads, *exceptions.BaseErrorResponse) {
	entity := masteritementities.ItemSubstitute{}
	response := masteritempayloads.ItemSubstitutePayloads{}

	err := tx.Model(entity).Select("mtr_item_substitute.*, Item.item_code, Item.item_name, Item.item_class_id, mi").
		Where(masteritementities.ItemSubstitute{ItemSubstituteId: id}).
		Joins("Item", tx.Select("")).
		Joins("JOIN mtr_item_class ON Item.item_class_id = mtr_item_class.item_class_id", tx.Select("")).
		First(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	return response, nil
}

func (r *ItemSubstituteRepositoryImpl) GetAllItemSubstituteDetail(tx *gorm.DB, pages pagination.Pagination, id int) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := []masteritementities.ItemSubstituteDetail{}
	payloads := []masteritempayloads.ItemSubstituteDetailPayloads{}

	query := tx.Model(entities).Select("mtr_item_substitute_detail.*, Item.item_code, Item.item_name").
		Joins("Item", tx.Select("")).Where("mtr_item_substitute_detail.item_substitute_id = ?", id)

	err := query.Scopes(pagination.Paginate(&entities, &pages, query)).Scan(&payloads).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.Rows = payloads

	return pages, nil
}

func (r *ItemSubstituteRepositoryImpl) GetByIdItemSubstituteDetail(tx *gorm.DB, id int) (masteritempayloads.ItemSubstituteDetailGetPayloads, *exceptions.BaseErrorResponse) {
	response := masteritempayloads.ItemSubstituteDetailGetPayloads{}
	tableStruct := masteritempayloads.ItemSubstituteDetailPayloads{}
	baseModelQuery := utils.CreateJoinSelectStatement(tx, tableStruct).Where(masteritementities.ItemSubstituteDetail{ItemSubstituteDetailId: id})

	rows, err := baseModelQuery.First(&response).Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()
	return response, nil
}

func (r *ItemSubstituteRepositoryImpl) SaveItemSubstitute(tx *gorm.DB, req masteritempayloads.ItemSubstitutePostPayloads) (bool, *exceptions.BaseErrorResponse) {

	// parseEffectiveDate, _ := time.Parse("2006-01-02T15:04:05.000Z", req.EffectiveDate)

	entities := masteritementities.ItemSubstitute{
		SubstituteTypeCode: req.SubstituteTypeCode,
		ItemSubstituteId:   req.ItemSubstituteId,
		EffectiveDate:      req.EffectiveDate,
		ItemId:             req.ItemId,
	}
	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
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

func (r *ItemSubstituteRepositoryImpl) ChangeStatusItemSubstitute(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
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
