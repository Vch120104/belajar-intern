package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"strings"

	"gorm.io/gorm"
)

type ItemSubstituteRepositoryImpl struct {
}

func StartItemSubstituteRepositoryImpl() masteritemrepository.ItemSubstituteRepository {
	return &ItemSubstituteRepositoryImpl{}
}

func (r *ItemSubstituteRepositoryImpl) GetAllItemSubstitute(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := []masteritementities.ItemSubstitute{}
	payloads := []masteritempayloads.ItemSubstitutePayloads{}
	tableStruct := masteritempayloads.ItemSubstitutePayloads{}
	baseModelQuery := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&payloads).Rows()

	if len(payloads) == 0 {
		return pages, gorm.ErrRecordNotFound
	}

	if err != nil {
		return pages, err
	}
	defer rows.Close()

	pages.Rows = payloads

	return pages, nil
}

func (r *ItemSubstituteRepositoryImpl) GetByIdItemSubstitute(tx *gorm.DB, id int) (masteritempayloads.ItemSubstitutePayloads, error) {
	entities := masteritementities.ItemSubstitute{}
	response := masteritempayloads.ItemSubstitutePayloads{}

	rows, err := tx.Model(&entities).Where(masteritementities.ItemSubstitute{ItemSubstituteId: id}).First(&response).Rows()

	if err != nil {
		return response, err
	}
	defer rows.Close()
	return response, nil
}

func (r *ItemSubstituteRepositoryImpl) GetAllItemSubstituteDetail(tx *gorm.DB, pages pagination.Pagination, id int) (pagination.Pagination, error) {
	entities := []masteritementities.ItemSubstituteDetail{}
	payloads := []masteritempayloads.ItemSubstituteDetailPayloads{}
	tableStruct := masteritempayloads.ItemSubstituteDetailPayloads{}

	baseModelQuery := utils.CreateJoinSelectStatement(tx, tableStruct).Where(masteritementities.ItemSubstituteDetail{ItemSubstituteId: id})

	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, baseModelQuery)).Scan(&payloads).Rows()

	if len(payloads) == 0 {
		return pages, gorm.ErrRecordNotFound
	}

	if err != nil {
		return pages, err
	}
	defer rows.Close()

	pages.Rows = payloads

	return pages, nil
}

func (r *ItemSubstituteRepositoryImpl) GetByIdItemSubstituteDetail(tx *gorm.DB, id int) (masteritempayloads.ItemSubstituteDetailGetPayloads, error) {
	entities := masteritementities.ItemSubstituteDetail{}
	response := masteritempayloads.ItemSubstituteDetailGetPayloads{}

	rows, err := tx.Model(&entities).Where(masteritementities.ItemSubstituteDetail{ItemSubstituteDetailId: id}).First(&response).Rows()

	if err != nil {
		return response, err
	}
	defer rows.Close()
	return response, nil
}

func (r *ItemSubstituteRepositoryImpl) SaveItemSubstitute(tx *gorm.DB, req masteritempayloads.ItemSubstitutePostPayloads) (bool, error) {
	entities := masteritementities.ItemSubstitute{
		SubstituteTypeCode: req.SubstituteTypeCode,
		ItemSubstituteId:   req.ItemSubstituteId,
		EffectiveDate:      req.EffectiveDate,
		ItemId:             req.ItemId,
	}
	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *ItemSubstituteRepositoryImpl) SaveItemSubstituteDetail(tx *gorm.DB, req masteritempayloads.ItemSubstituteDetailPostPayloads, id int) (bool, error) {
	entities := masteritementities.ItemSubstituteDetail{
		ItemSubstituteDetailId: req.ItemSubstituteDetailId,
		ItemId:                 req.ItemId,
		ItemSubstituteId:       id,
		Quantity:               req.Quantity,
		Sequence:               req.Sequence,
	}
	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *ItemSubstituteRepositoryImpl) ChangeStatusItemOperation(tx *gorm.DB, id int) (bool, error) {
	var entities masteritementities.ItemSubstitute

	result := tx.Model(&entities).
		Where("item_substitute_id = ?", id).
		First(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	if entities.IsActive {
		entities.IsActive = false
	} else {
		entities.IsActive = true
	}

	result = tx.Save(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

func (r *ItemSubstituteRepositoryImpl) DeactivateItemSubstituteDetail(tx *gorm.DB, id string) (bool, error) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteritementities.ItemSubstituteDetail
		err := tx.Model(&entityToUpdate).Where("item_substitute_detail_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, err
		}

		entityToUpdate.IsActive = false
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, result.Error
		}
	}

	return true, nil
}

func (r *ItemSubstituteRepositoryImpl) ActivateItemSubstituteDetail(tx *gorm.DB, id string) (bool, error) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteritementities.ItemSubstituteDetail
		err := tx.Model(&entityToUpdate).Where("item_substitute_detail_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, err
		}

		entityToUpdate.IsActive = true
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, result.Error
		}
	}

	return true, nil
}
