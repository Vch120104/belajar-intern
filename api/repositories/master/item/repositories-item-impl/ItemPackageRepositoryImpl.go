package masteritemrepositoryimpl

import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemlevelrepo "after-sales/api/repositories/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemPackageRepositoryImpl struct {
}

func StartItemPackageRepositoryImpl() masteritemlevelrepo.ItemPackageRepository {
	return &ItemPackageRepositoryImpl{}
}

func (r *ItemPackageRepositoryImpl) GetAllItemPackage(tx *gorm.DB, internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, error) {
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
		return nil, 0, 0, err
	}

	defer rows.Close()

	if len(responses) == 0 {
		return nil, 0, 0, gorm.ErrRecordNotFound
	}

	itemGroupUrl := config.EnvConfigs.GeneralServiceUrl + "api/general/filter-item-group?item_group_code=" + itemGroupCode

	errUrlItemPackage := utils.Get(itemGroupUrl, &getItemGroupResponses, nil)

	if errUrlItemPackage != nil {
		return nil, 0, 0, errUrlItemPackage
	}

	joinedData := utils.DataFrameInnerJoin(responses, getItemGroupResponses, "ItemGroupId")
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (*ItemPackageRepositoryImpl) GetItemPackageById(tx *gorm.DB, Id int) ([]map[string]interface{}, error) {

	tableStruct := masteritementities.ItemPackage{}
	response := []masteritempayloads.GetAllItemPackageResponse{}

	getItemGroupResponses := []masteritempayloads.ItemGroupResponse{}

	baseModelQuery := tx.Model(&tableStruct)

	rows, err := baseModelQuery.Where(masteritementities.ItemPackage{
		ItemPackageId: Id,
	}).First(&response).Rows()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	itemGroupUrl := config.EnvConfigs.GeneralServiceUrl + "api/general/filter-item-group"
	errUrlItemPackage := utils.Get(itemGroupUrl, &getItemGroupResponses, nil)

	if errUrlItemPackage != nil {
		return nil, err
	}

	joinedData := utils.DataFrameInnerJoin(response, getItemGroupResponses, "ItemGroupId")

	return joinedData, nil
}

func (r *ItemPackageRepositoryImpl) SaveItemPackage(tx *gorm.DB, request masteritempayloads.SaveItemPackageRequest) (bool, error) {
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
		return false, err
	}

	return true, nil
}
