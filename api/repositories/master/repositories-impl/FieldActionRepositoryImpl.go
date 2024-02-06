package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	// masterpayloads "after-sales/api/payloads/master"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type FieldActionRepositoryImpl struct {
}

func StartFieldActionRepositoryImpl() masterrepository.FieldActionRepository {
	return &FieldActionRepositoryImpl{}
}

func (r *FieldActionRepositoryImpl) GetAllFieldAction(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := []masterentities.FieldAction{}
	//define base model
	baseModelQuery := tx.Model(&entities)
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

func (r *FieldActionRepositoryImpl) SaveFieldAction(tx *gorm.DB, req masterpayloads.FieldActionResponse) (bool, error) {
	entities := masterentities.FieldAction{
		IsActive:                 req.IsActive,
		FieldActionSystemNumber: req.FieldActionSystemNumber,
		FieldActionDocumentNumber:          req.FieldActionDocumentNo,
		ApprovalValue:   req.ApprovalValue,
		BrandId:   req.BrandId,
		FieldActionName:   req.FieldActionName,
		FieldActionPeriodFrom:   req.FieldActionPeriodFrom,
		FieldActionPeriodTo:   req.FieldActionPeriodTo,
		IsNeverExpired:   req.IsNeverExpired,
		RemarkPopup:   req.RemarkPopup,
		IsCritical:   req.IsCritical,
		RemarkInvoice:   req.RemarkInvoice,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *FieldActionRepositoryImpl) GetFieldActionHeaderById(tx *gorm.DB, Id int) (masterpayloads.FieldActionResponse, error) {
	entities := masterentities.FieldAction{}
	response := masterpayloads.FieldActionResponse{}

	rows, err := tx.Model(&entities).
		Where(masterentities.FieldAction{
			FieldActionSystemNumber: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *FieldActionRepositoryImpl) GetFieldActionDetailById(tx *gorm.DB, Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := []masterentities.FieldActionEligibleVehicle{}

	//define base model
	baseModelQuery := tx.Model(&entities)
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
