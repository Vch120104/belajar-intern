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
	var responses []masterpayloads.ItemOperationPost
	var entities masterentities.MappingItemOperation

	query := tx.Model(&entities).Select("mtr_item_operation.*").Table("mtr_item_operation")
	WhereQuery := utils.ApplyFilter(query, filterCondition)

	err := WhereQuery.Scopes(pagination.Paginate(&pages, query)).Scan(&responses).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	pages.Rows = responses
	return pages, nil
}

func (r *ItemOperationRepositoryImpl) GetAllItemOperationLineType(tx *gorm.DB, lineTypeId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []masterpayloads.ItemOperationPost
	var entities masterentities.MappingItemOperation

	query := tx.Model(&entities).Select("mtr_item_operation.*").Table("mtr_item_operation")
	WhereQuery := utils.ApplyFilter(query, filterCondition)

	err := WhereQuery.Scopes(pagination.Paginate(&pages, query)).Scan(&responses).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	pages.Rows = responses
	return pages, nil
}

func (r *ItemOperationRepositoryImpl) GetByIdItemOperation(tx *gorm.DB, id int) (masterpayloads.ItemOperationPost, *exceptions.BaseErrorResponse) {
	var responses masterpayloads.ItemOperationPost
	err := tx.Select("mtr_item_operation.*").Table("mtr_item_operation").
		Where("item_operation_id=?", id).Scan(&responses).Error
	if err != nil {
		return masterpayloads.ItemOperationPost{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	return responses, nil
}

func (r *ItemOperationRepositoryImpl) PostItemOperation(tx *gorm.DB, req masterpayloads.ItemOperationPost) (masterentities.MappingItemOperation, *exceptions.BaseErrorResponse) {
	entities := masterentities.MappingItemOperation{
		ItemOperationId: req.ItemOperationId,
		LineTypeId:      req.LineTypeId,
		ItemId:          req.ItemId,
		OperationId:     req.OperationId,
		PackageId:       req.PackageId,
	}
	err := tx.Save(&entities).Error
	if err != nil {
		return masterentities.MappingItemOperation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return entities, nil
}

func (r *ItemOperationRepositoryImpl) DeleteItemOperation(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity masterentities.MappingItemOperation
	err := tx.Model(masterentities.MappingItemOperation{}).Where("item_operation_id = ?", id).Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return true, nil
}

func (r *ItemOperationRepositoryImpl) UpdateItemOperation(tx *gorm.DB, id int, req masterpayloads.ItemOperationPost) (masterentities.MappingItemOperation, *exceptions.BaseErrorResponse) {
	entities := masterentities.MappingItemOperation{
		ItemOperationId: req.ItemOperationId,
		LineTypeId:      req.LineTypeId,
		ItemId:          req.ItemId,
		OperationId:     req.OperationId,
		PackageId:       req.PackageId,
	}
	err := tx.Model(masterentities.MappingItemOperation{}).Where("item_operation_id=?", id).Updates(&entities).Error
	if err != nil {
		return masterentities.MappingItemOperation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return entities, nil
}
