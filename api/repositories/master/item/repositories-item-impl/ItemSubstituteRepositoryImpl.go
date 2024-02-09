package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type ItemSubstituteRepositoryImpl struct {
}

func StartItemSubstituteRepositoryImpl(db *gorm.DB) masteritemrepository.ItemSubstituteRepository {
	return &ItemSubstituteRepositoryImpl{}
}

func (r *ItemSubstituteRepositoryImpl) GetAllItemSubstitute(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := []masteritementities.ItemSubstitute{}

	baseModelQuery := tx.Model(&entities)

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&entities).Rows()

	if len(entities) == 0 {
		return pages, gorm.ErrRecordNotFound
	}

	if err != nil {
		return pages, err
	}
	defer rows.Close()

	pages.Rows = entities

	return pages, nil
}

func (r *ItemSubstituteRepositoryImpl) GetByIdItemSubstitute(tx *gorm.DB,id int) (masteritempayloads.ItemSubstitutePayloads, error) {
	entities := masteritementities.ItemSubstitute{}
	response := masteritempayloads.ItemSubstitutePayloads{}

	rows, err := tx.Model(&entities).Where(masteritementities.ItemSubstitute{ItemSubstituteId: id}).First(&response).Rows()

	if err != nil {
		return response, err
	}
	defer rows.Close()
	return response, nil
}

func (r *ItemSubstituteRepositoryImpl) GetAllItemSubstituteDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := []masteritementities.ItemSubstituteDetail{}

	baseModelQuery := tx.Model(&entities)

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&entities).Rows()

	if len(entities) == 0 {
		return pages, gorm.ErrRecordNotFound
	}

	if err != nil {
		return pages, err
	}
	defer rows.Close()

	pages.Rows = entities

	return pages, nil
}

func (r *ItemSubstituteRepositoryImpl) GetByIdItemSubstituteDetail(tx *gorm.DB, id int) (masteritempayloads.ItemSubstituteDetailPayloads, error) {
	entities := masteritementities.ItemSubstituteDetail{}
	response := masteritempayloads.ItemSubstituteDetailPayloads{}

	rows, err := tx.Model(&entities).Where(masteritementities.ItemSubstituteDetail{ItemSubstituteId: id}).First(&response).Rows()

	if err != nil {
		return response, err
	}
	defer rows.Close()
	return response, nil
}

func (r *ItemSubstituteRepositoryImpl) SaveItemSubstitute(tx *gorm.DB, req masteritempayloads.ItemSubstitutePayloads) (bool, error) {
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

func (r *ItemSubstituteRepositoryImpl) SaveItemSubstituteDetail(tx *gorm.DB, req masteritempayloads.ItemSubstituteDetailPayloads) (bool, error) {
	entities := masteritementities.ItemSubstituteDetail{
		ItemSubstituteDetailId: req.ItemSubstituteDetailId,
		ItemSubstituteId:       req.ItemSubstituteId,
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
		Where("operation_group_id = ?", id).
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
	var entities masteritementities.ItemSubstituteDetail
	strid := strings.Split(id, ",")

	var strids []int

	for _, numid := range strid {
		num, err := strconv.Atoi(numid)
		if err != nil {
			return false, err
		}
		strids = append(strids, num)
	}
	for _, value := range strids {
		var entityToUpdate masteritementities.ItemSubstituteDetail
		err := tx.Model(&entities).Where(masteritementities.ItemSubstituteDetail{
			ItemSubstituteDetailId: int(value),
		}).First(&entityToUpdate).Error
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
	var entities masteritementities.ItemSubstituteDetail

	strid := strings.Split(id, ",")

	var strids []int

	for _, numid := range strid {
		num, err := strconv.Atoi(numid)
		if err != nil {
			return false, err
		}
		strids = append(strids, num)
	}
	for _, value := range strids {
		var entityToUpdate masteritementities.ItemSubstituteDetail
		err := tx.Model(&entities).Where(masteritementities.ItemSubstituteDetail{
			ItemSubstituteDetailId: int(value),
		}).First(&entityToUpdate).Error
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
