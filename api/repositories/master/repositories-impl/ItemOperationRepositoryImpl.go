package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"net/http"

	"gorm.io/gorm"
)

type ItemOperationRepositoryImpl struct {
}

func StartItemOperationRepositoryImpl() masterrepository.ItemOperationRepository {
	return &ItemOperationRepositoryImpl{}
}

func (r *ItemOperationRepositoryImpl) GetAllItemOperation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []masterpayloads.ItemOperationGet
	var entities masterentities.ItemOperation

	query := tx.Select("mtr_item_operation.*,mtr_item.item_name,mtr_operation_code.operation_name").
		Joins("Join mtr_item on mtr_item.Item_id=mtr_item_operation.item_id").
		Joins("Join mtr_operation_model_mapping on mtr_operation_model_mapping.operation_model_mapping_id=mtr_item_operation.operation_model_mapping_id").
		Joins("Join mtr_operation_code on mtr_operation_code.operation_id=mtr_operation_model_mapping.operation_id")

	WhereQuery := utils.ApplyFilter(query, filterCondition)

	err := WhereQuery.Scopes(pagination.Paginate(&entities, &pages, query)).Scan(&responses).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	pages.Rows = responses
	return pages, nil
}

func (r *ItemOperationRepositoryImpl) GetByIdItemOperation(tx *gorm.DB, id int) (masterpayloads.ItemOperationGet, *exceptions.BaseErrorResponse) {
	var responses masterpayloads.ItemOperationGet
	err := tx.Select("mtr_item_operation.*,mtr_item.item_name,mtr_operation_code.operation_name").Table("mtr_item_operation").
		Joins("Join mtr_item on mtr_item.Item_id=mtr_item_operation.item_id").
		Joins("Join mtr_operation_model_mapping on mtr_operation_model_mapping.operation_model_mapping_id=mtr_item_operation.operation_model_mapping_id").
		Joins("Join mtr_operation_code on mtr_operation_code.operation_id=mtr_operation_model_mapping.operation_id").Where("item_operation_id=?", id).Scan(&responses).Error
	if err != nil {
		return masterpayloads.ItemOperationGet{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	return responses, nil
}

func (r *ItemOperationRepositoryImpl) PostItemOperation(tx *gorm.DB, req masterpayloads.ItemOperationPost) (masterentities.ItemOperation, *exceptions.BaseErrorResponse) {
	entities := masterentities.ItemOperation{
		ItemId:                  req.ItemId,
		OperationModelMappingId: req.OperationModelMappingId,
		LineTypeId:              req.LineTypeId,
	}
	err := tx.Save(&entities).Error
	if err != nil {
		return masterentities.ItemOperation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return entities, nil
}

func (r *ItemOperationRepositoryImpl) DeleteItemOperation(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity masterentities.ItemOperation
	err := tx.Model(masterentities.ItemOperation{}).Where("item_operation_id = ?", id).Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return true, nil
}

func (r *ItemOperationRepositoryImpl) UpdateItemOperation(tx *gorm.DB, id int, req masterpayloads.ItemOperationPost) (masterentities.ItemOperation, *exceptions.BaseErrorResponse) {
	entities := masterentities.ItemOperation{
		ItemId:                  req.ItemId,
		OperationModelMappingId: req.OperationModelMappingId,
		LineTypeId:              req.LineTypeId,
	}
	err := tx.Model(masterentities.ItemOperation{}).Where("item_operation_id=?", id).Updates(&entities).Error
	if err != nil {
		return masterentities.ItemOperation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return entities, nil
}
