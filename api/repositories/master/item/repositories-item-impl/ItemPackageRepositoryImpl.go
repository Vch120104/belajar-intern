package masteritemrepositoryimpl

import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemlevelrepo "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type ItemPackageRepositoryImpl struct {
}

// ChangeStatusItemPackage implements masteritemrepository.ItemPackageRepository.

func StartItemPackageRepositoryImpl() masteritemlevelrepo.ItemPackageRepository {
	return &ItemPackageRepositoryImpl{}
}

func (r *ItemPackageRepositoryImpl) GetAllItemPackage(tx *gorm.DB, internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	var responses []masteritempayloads.GetAllItemPackageResponse
	var getItemGroupResponses []masteritempayloads.GetItemGroupResponse

	var itemGroupCode string

	//apply external services filter
	for i := 0; i < len(externalFilterCondition); i++ {
		itemGroupCode = externalFilterCondition[i].ColumnValue
	}

	tableStruct := masteritementities.ItemPackage{}

	baseModelQuery := tx.Model(&tableStruct)

	whereQuery := utils.ApplyFilter(baseModelQuery, internalFilterCondition)

	rows, err := whereQuery.Scan(&responses).Rows()

	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	if len(responses) == 0 {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	itemGroupUrl := config.EnvConfigs.GeneralServiceUrl + "/filter-item-group?item_group_code=" + itemGroupCode

	errUrlItemPackage := utils.Get(itemGroupUrl, &getItemGroupResponses, nil)

	if errUrlItemPackage != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	joinedData := utils.DataFrameInnerJoin(responses, getItemGroupResponses, "ItemGroupId")

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (*ItemPackageRepositoryImpl) GetItemPackageById(tx *gorm.DB, Id int) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse) {

	tableStruct := masteritementities.ItemPackage{}
	response := []masteritempayloads.GetAllItemPackageResponse{}

	getItemGroupResponses := []masteritempayloads.ItemGroupResponse{}

	baseModelQuery := tx.Model(&tableStruct)

	rows, err := baseModelQuery.Where(masteritementities.ItemPackage{
		ItemPackageId: Id,
	}).First(&response).Rows()

	if err != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(response) == 0 {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	defer rows.Close()

	itemGroupUrl := config.EnvConfigs.GeneralServiceUrl + "api/general/filter-item-group"
	errUrlItemPackage := utils.Get(itemGroupUrl, &getItemGroupResponses, nil)

	if errUrlItemPackage != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	joinedData := utils.DataFrameInnerJoin(response, getItemGroupResponses, "ItemGroupId")

	return joinedData, nil
}

func (r *ItemPackageRepositoryImpl) SaveItemPackage(tx *gorm.DB, request masteritempayloads.SaveItemPackageRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masteritementities.ItemPackage{
		IsActive:        request.IsActive,
		ItemGroupId:     request.ItemGroupId,
		ItemPackageId:   request.ItemPackageId,
		ItemPackageCode: request.ItemPackageCode,
		ItemPackageName: request.ItemPackageName,
		ItemPackageSet:  request.ItemPackageSet,
		Description:     request.Description,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *ItemPackageRepositoryImpl) ChangeStatusItemPackage(tx *gorm.DB, id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities masteritementities.ItemPackage

	result := tx.Model(&entities).
		Where("item_package_id = ?", id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
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
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}
