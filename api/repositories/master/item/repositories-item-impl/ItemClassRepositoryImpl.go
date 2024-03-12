package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	masteritempayloads "after-sales/api/payloads/master/item"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ItemClassRepositoryImpl struct {
}

func StartItemClassRepositoryImpl() masteritemrepository.ItemClassRepository {
	return &ItemClassRepositoryImpl{}
}

func (r *ItemClassRepositoryImpl) GetAllItemClass(tx *gorm.DB, filterCondition []utils.FilterCondition) ([]map[string]interface{}, error) {
	entities := []masteritementities.ItemClass{}
	var responses []masteritempayloads.ItemClassResponse
	var getLineTypeResponse []masteritempayloads.LineTypeResponse
	var getItemGroupResponse []masteritempayloads.ItemGroupResponse
	var c *gin.Context
	var internalServiceFilter, externalServiceFilter []utils.FilterCondition
	var groupName, lineTypeCode string
	responseStruct := reflect.TypeOf(masteritempayloads.ItemClassResponse{})

	for i := 0; i < len(filterCondition); i++ {
		flag := false
		for j := 0; j < responseStruct.NumField(); j++ {
			if filterCondition[i].ColumnField == responseStruct.Field(j).Tag.Get("parent_entity")+"."+responseStruct.Field(j).Tag.Get("json") {
				internalServiceFilter = append(internalServiceFilter, filterCondition[i])
				flag = true
				break
			}
		}
		if !flag {
			externalServiceFilter = append(externalServiceFilter, filterCondition[i])
		}
	}

	//apply external services filter
	for i := 0; i < len(externalServiceFilter); i++ {
		if strings.Contains(externalServiceFilter[i].ColumnField, "line_type_code") {
			lineTypeCode = externalServiceFilter[i].ColumnValue
		} else {
			groupName = externalServiceFilter[i].ColumnValue
		}
	}

	//define base model
	baseModelQuery := tx.Model(&entities)
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, internalServiceFilter)
	//apply pagination and execute
	rows, err := whereQuery.Scan(&responses).Rows()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if len(responses) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	groupServiceUrl := "http://10.1.32.26:8000/general-service/api/general/filter-item-group?item_group_name=" + groupName

	errUrlItemGroup := utils.Get(c, groupServiceUrl, &getItemGroupResponse, nil)

	if errUrlItemGroup != nil {
		return nil, errUrlItemGroup
	}

	joinedData := utils.DataFrameInnerJoin(responses, getItemGroupResponse, "ItemGroupId")

	lineTypeUrl := "http://10.1.32.26:8000/general-service/api/general/line-type?line_type_code=" + lineTypeCode

	errUrlLineType := utils.Get(c, lineTypeUrl, &getLineTypeResponse, nil)

	if errUrlLineType != nil {
		return nil, errUrlLineType
	}

	joinedDataSecond := utils.DataFrameInnerJoin(joinedData, getLineTypeResponse, "LineTypeId")

	return joinedDataSecond, nil
}

func (r *ItemClassRepositoryImpl) GetItemClassById(tx *gorm.DB, Id int) (masteritempayloads.ItemClassResponse, error) {
	entities := masteritementities.ItemClass{}
	response := masteritempayloads.ItemClassResponse{}

	rows, err := tx.Model(&entities).
		Where(masteritementities.ItemClass{
			ItemClassId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *ItemClassRepositoryImpl) SaveItemClass(tx *gorm.DB, request masteritempayloads.ItemClassResponse) (bool, error) {
	entities := masteritementities.ItemClass{
		IsActive:      request.IsActive,
		ItemClassId:   request.ItemClassId,
		ItemClassCode: request.ItemClassCode,
		ItemGroupID:   request.ItemGroupId,
		LineTypeID:    request.LineTypeId,
		ItemClassName: request.ItemClassName,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *ItemClassRepositoryImpl) ChangeStatusItemClass(tx *gorm.DB, Id int) (bool, error) {
	var entities masteritementities.ItemClass

	result := tx.Model(&entities).
		Where("item_class_id = ?", Id).
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
