package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemlevelrepo "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	aftersalesserviceapiutils "after-sales/api/utils/aftersales-service"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type ItemPackageRepositoryImpl struct {
}

func StartItemPackageRepositoryImpl() masteritemlevelrepo.ItemPackageRepository {
	return &ItemPackageRepositoryImpl{}
}

// GetItemPackageByCode implements masteritemrepository.ItemPackageRepository.
func (r *ItemPackageRepositoryImpl) GetItemPackageByCode(tx *gorm.DB, itemPackageCode string) (masteritempayloads.GetItemPackageResponse, *exceptions.BaseErrorResponse) {
	tableStruct := masteritementities.ItemPackage{}
	response := masteritempayloads.GetItemPackageResponse{}

	baseModelQuery := tx.Model(&tableStruct).Select("mtr_item_package.*")
	err := baseModelQuery.Where(masteritementities.ItemPackage{
		ItemPackageCode: itemPackageCode,
	}).First(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	itemGroupResponse, itemGroupErr := aftersalesserviceapiutils.GetItemGroupById(response.ItemGroupId)
	if itemGroupErr != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: itemGroupErr.StatusCode,
			Err:        itemGroupErr.Err,
		}
	}

	response.ItemGroupName = itemGroupResponse.ItemGroupName
	response.ItemGroupCode = itemGroupResponse.ItemGroupCode

	return response, nil
}

func (r *ItemPackageRepositoryImpl) GetAllItemPackage(tx *gorm.DB, internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	responses := []masteritempayloads.GetAllItemPackageResponse{}
	entities := masteritementities.ItemPackage{}

	query := tx.Model(&entities).
		Select(`
			mtr_item_package.is_active,
			mtr_item_package.item_package_id,
			mtr_item_package.item_package_code,
			mtr_item_package.item_package_name,
			mtr_item_group.item_group_id AS item_group_id,
			mtr_item_group.item_group_code AS item_group_code,
			mtr_item_group.item_group_name AS item_group_name,
			mtr_item_package.item_package_set
		`).
		Joins("LEFT JOIN mtr_item_group ON mtr_item_group.item_group_id = mtr_item_package.item_group_id")
	queryFilter := utils.ApplyFilter(query, internalFilterCondition)

	for _, filter := range externalFilterCondition {
		if filter.ColumnField == "item_group_code" {
			queryFilter = queryFilter.Where("mtr_item_group.item_group_code = ?", filter.ColumnValue)
		}
	}

	err := queryFilter.Scopes(pagination.Paginate(&pages, queryFilter)).
		Order("mtr_item_package.item_package_id").
		Scan(&responses).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []masteritempayloads.GetAllItemPackageResponse{}
		return pages, nil
	}

	pages.Rows = responses

	return pages, nil
}

func (*ItemPackageRepositoryImpl) GetItemPackageById(tx *gorm.DB, Id int) (masteritempayloads.GetItemPackageResponse, *exceptions.BaseErrorResponse) {
	tableStruct := masteritementities.ItemPackage{}
	response := masteritempayloads.GetItemPackageResponse{}

	baseModelQuery := tx.Model(&tableStruct).Select("mtr_item_package.*")
	err := baseModelQuery.Where(masteritementities.ItemPackage{
		ItemPackageId: Id,
	}).First(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	itemGroupResponse, itemGroupErr := aftersalesserviceapiutils.GetItemGroupById(response.ItemGroupId)
	if itemGroupErr != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: itemGroupErr.StatusCode,
			Err:        itemGroupErr.Err,
		}
	}

	response.ItemGroupName = itemGroupResponse.ItemGroupName
	response.ItemGroupCode = itemGroupResponse.ItemGroupCode

	return response, nil
}

func (r *ItemPackageRepositoryImpl) SaveItemPackage(tx *gorm.DB, request masteritempayloads.SaveItemPackageRequest) (masteritementities.ItemPackage, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemPackage{
		IsActive:        request.IsActive,
		ItemGroupId:     request.ItemGroupId,
		ItemPackageId:   request.ItemPackageId,
		ItemPackageCode: request.ItemPackageCode,
		ItemPackageName: request.ItemPackageName,
		ItemPackageSet:  request.ItemPackageSet,
		Description:     request.Description,
	}

	result := masteritementities.ItemPackage{}

	err := tx.Save(&entities).Where(masteritementities.ItemPackage{ItemPackageCode: request.ItemPackageCode}).First(&result).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return result, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {
			return result, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return result, nil
}

func (r *ItemPackageRepositoryImpl) ChangeStatusItemPackage(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.ItemPackage

	result := tx.Model(&entities).
		Where("item_package_id = ?", id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
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
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}
