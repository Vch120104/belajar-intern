package masteroperationrepositoryimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	"after-sales/api/utils"
	"log"

	"gorm.io/gorm"
)

type OperationGroupRepositoryImpl struct {
	myDB *gorm.DB
}

func StartOperationGroupRepositoryImpl(db *gorm.DB) masteroperationrepository.OperationGroupRepository {
	return &OperationGroupRepositoryImpl{myDB: db}
}

func (r *OperationGroupRepositoryImpl) WithTrx(trxHandle *gorm.DB) masteroperationrepository.OperationGroupRepository {
	if trxHandle == nil {
		log.Println("Transaction Database Not Found!")
		return r
	}
	r.myDB = trxHandle
	return r
}

func (r *OperationGroupRepositoryImpl) GetAllOperationGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := []masteroperationentities.OperationGroup{}
	//define base model
	baseModelQuery := r.myDB.Model(&entities)
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	//apply pagination and execute
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

func (r *OperationGroupRepositoryImpl) GetAllOperationGroupIsActive() ([]masteroperationpayloads.OperationGroupResponse, error) {
	var OperationGroups []masteroperationentities.OperationGroup
	response := []masteroperationpayloads.OperationGroupResponse{}

	err := r.myDB.Model(&OperationGroups).Where("is_active = 'true'").Scan(&response).Error

	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *OperationGroupRepositoryImpl) GetOperationGroupById(Id int) (masteroperationpayloads.OperationGroupResponse, error) {
	entities := masteroperationentities.OperationGroup{}
	response := masteroperationpayloads.OperationGroupResponse{}

	rows, err := r.myDB.Model(&entities).
		Where(masteroperationentities.OperationGroup{
			OperationGroupId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *OperationGroupRepositoryImpl) GetOperationGroupByCode(Code string) (masteroperationpayloads.OperationGroupResponse, error) {
	entities := masteroperationentities.OperationGroup{}
	response := masteroperationpayloads.OperationGroupResponse{}

	rows, err := r.myDB.Model(&entities).
		Where(masteroperationentities.OperationGroup{
			OperationGroupCode: Code,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *OperationGroupRepositoryImpl) SaveOperationGroup(req masteroperationpayloads.OperationGroupResponse) (bool, error) {
	entities := masteroperationentities.OperationGroup{
		IsActive:                  req.IsActive,
		OperationGroupId:          req.OperationGroupId,
		OperationGroupCode:        req.OperationGroupCode,
		OperationGroupDescription: req.OperationGroupDescription,
	}

	err := r.myDB.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *OperationGroupRepositoryImpl) ChangeStatusOperationGroup(Id int) (bool, error) {
	var entities masteroperationentities.OperationGroup

	result := r.myDB.Model(&entities).
		Where("operation_group_id = ?", Id).
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
