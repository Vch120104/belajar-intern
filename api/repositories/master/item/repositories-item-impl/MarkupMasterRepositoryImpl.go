package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type MarkupMasterRepositoryImpl struct {
}

func StartMarkupMasterRepositoryImpl() masteritemrepository.MarkupMasterRepository {
	return &MarkupMasterRepositoryImpl{}
}

func (r *MarkupMasterRepositoryImpl) GetMarkupMasterList(tx *gorm.DB,filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	var responses []masteritementities.MarkupMaster

	baseModelQuery := tx.Model(&responses)
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	//apply pagination and execute
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&responses, &pages, whereQuery)).Scan(&responses).Rows()

	if len(responses) == 0 {
		return pages, gorm.ErrRecordNotFound
	}

	if err != nil {
		return pages, err
	}

	defer rows.Close()

	pages.Rows = responses

	return pages, nil
}

func (r *MarkupMasterRepositoryImpl) GetMarkupMasterById(tx *gorm.DB,Id int) (masteritempayloads.MarkupMasterResponse, error) {
	entities := masteritementities.MarkupMaster{}
	response := masteritempayloads.MarkupMasterResponse{}

	rows, err := tx.Model(&entities).
		Where(masteritementities.MarkupMaster{
			MarkupMasterId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *MarkupMasterRepositoryImpl) SaveMarkupMaster(tx *gorm.DB,req masteritempayloads.MarkupMasterResponse) (bool, error) {
	entities := masteritementities.MarkupMaster{
		IsActive:                req.IsActive,
		MarkupMasterId:          req.MarkupMasterId,
		MarkupMasterCode:        req.MarkupMasterCode,
		MarkupMasterDescription: req.MarkupMasterDescription,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}
func (r *MarkupMasterRepositoryImpl) ChangeStatusMasterMarkupMaster(tx *gorm.DB,Id int) (bool, error) {
	var entities masteritementities.MarkupMaster

	result := tx.Model(&entities).
		Where("markup_master_id = ?", Id).
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
func (r *MarkupMasterRepositoryImpl) GetMarkupMasterByCode(tx *gorm.DB,markupCode string) (masteritempayloads.MarkupMasterResponse, error) {
	response := masteritempayloads.MarkupMasterResponse{}
	var entities masteritementities.MarkupMaster
	rows, err := tx.Model(&entities).
		Where("markup_master_code = ?", markupCode).
		First(&response).Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}
