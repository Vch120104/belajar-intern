package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"log"

	"gorm.io/gorm"
)

type MarkupMasterRepositoryImpl struct {
	myDB *gorm.DB
}

func StartMarkupMasterRepositoryImpl(db *gorm.DB) masteritemrepository.MarkupMasterRepository {
	return &MarkupMasterRepositoryImpl{myDB: db}
}

func (r *MarkupMasterRepositoryImpl) WithTrx(trxHandle *gorm.DB) masteritemrepository.MarkupMasterRepository {
	if trxHandle == nil {
		log.Println("Transaction Database Not Found!")
		return r
	}
	r.myDB = trxHandle
	return r
}

func (r *MarkupMasterRepositoryImpl) GetMarkupMasterList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	var responses []masteritementities.MarkupMaster

	baseModelQuery := r.myDB.Model(&responses)
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

func (r *MarkupMasterRepositoryImpl) GetMarkupMasterById(Id int) (masteritempayloads.MarkupMasterResponse, error) {
	entities := masteritementities.MarkupMaster{}
	response := masteritempayloads.MarkupMasterResponse{}

	rows, err := r.myDB.Model(&entities).
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

func (r *MarkupMasterRepositoryImpl) SaveMarkupMaster(req masteritempayloads.MarkupMasterResponse) (bool, error) {
	entities := masteritementities.MarkupMaster{
		IsActive:                req.IsActive,
		MarkupMasterId:          req.MarkupMasterId,
		MarkupMasterCode:        req.MarkupMasterCode,
		MarkupMasterDescription: req.MarkupMasterDescription,
	}

	err := r.myDB.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}
func (r *MarkupMasterRepositoryImpl) ChangeStatusMasterMarkupMaster(Id int) (bool, error) {
	var entities masteritementities.MarkupMaster

	result := r.myDB.Model(&entities).
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

	result = r.myDB.Save(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}
func (r *MarkupMasterRepositoryImpl) GetMarkupMasterByCode(markupCode string) (masteritempayloads.MarkupMasterResponse, error) {
	response := masteritempayloads.MarkupMasterResponse{}
	var entities masteritementities.MarkupMaster
	rows, err := r.myDB.Model(&entities).
		Where("markup_master_code = ?", markupCode).
		First(&response).Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}
