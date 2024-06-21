package masteritemrepositoryimpl

import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemlevelrepo "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type ItemPackageRepositoryImpl struct {
}

// ChangeStatusItemPackage implements masteritemrepository.ItemPackageRepository.

func StartItemPackageRepositoryImpl() masteritemlevelrepo.ItemPackageRepository {
	return &ItemPackageRepositoryImpl{}
}

func (r *ItemPackageRepositoryImpl) GetAllItemPackage(tx *gorm.DB, internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
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
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	if len(responses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	itemGroupUrl := config.EnvConfigs.GeneralServiceUrl + "filter-item-group?item_group_code=" + itemGroupCode

	errUrlItemPackage := utils.Get(itemGroupUrl, &getItemGroupResponses, nil)

	if errUrlItemPackage != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	joinedData := utils.DataFrameInnerJoin(responses, getItemGroupResponses, "ItemGroupId")

	fmt.Print(joinedData)

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (*ItemPackageRepositoryImpl) GetItemPackageById(tx *gorm.DB, Id int) (masteritempayloads.GetItemPackageResponse, *exceptions.BaseErrorResponse) {

	tableStruct := masteritementities.ItemPackage{}
	response := masteritempayloads.GetItemPackageResponse{}

	getItemGroupResponses := masteritempayloads.ItemGroupResponse{}

	baseModelQuery := tx.Model(&tableStruct).Select("mtr_item_package.*")

	err := baseModelQuery.Where(masteritementities.ItemPackage{
		ItemPackageId: Id,
	}).First(&response).Error

	fmt.Println(response)

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	itemGroupUrl := config.EnvConfigs.GeneralServiceUrl + "item-group/" + strconv.Itoa(response.ItemGroupId)
	errUrlItemPackage := utils.Get(itemGroupUrl, &getItemGroupResponses, nil)

	response.ItemGroupName = &getItemGroupResponses.ItemGroupName
	response.ItemGroupCode = &getItemGroupResponses.ItemGroupCode

	if errUrlItemPackage != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	return response, nil
}

func (r *ItemPackageRepositoryImpl) SaveItemPackage(tx *gorm.DB, request masteritempayloads.SaveItemPackageRequest) (bool, *exceptions.BaseErrorResponse) {
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
		if strings.Contains(err.Error(), "duplicate") {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return true, nil
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
