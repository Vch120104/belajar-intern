package masteritemrepositoryimpl

import (
	exceptions "after-sales/api/exceptions"
	masteritemlevelrepo "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strings"

	// masteritemlevelservice "after-sales/api/services/master/item_level"
	masteritementities "after-sales/api/entities/master/item"
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type ItemLevelImpl struct {
}

func StartItemLevelRepositoryImpl() masteritemlevelrepo.ItemLevelRepository {
	return &ItemLevelImpl{}
}

func (r *ItemLevelImpl) GetItemLevelLookUpbyId(tx *gorm.DB, filter []utils.FilterCondition, itemLevelId int) (masteritemlevelpayloads.GetItemLevelLookUp, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemLevel1{}
	response := masteritemlevelpayloads.GetItemLevelLookUp{}

	query := tx.Model(&entities).
		Select(`
			mtr_item_level_1.item_level_1_id,
			mtr_item_level_1.item_level_1_code,
			mtr_item_level_1.item_level_1_name,
			mil2.item_level_2_id,
			mil2.item_level_2_code,
			mil2.item_level_2_name,
			mil3.item_level_3_id,
			mil3.item_level_3_code,
			mil3.item_level_3_name,
			mil4.item_level_4_id,
			mil4.item_level_4_code,
			mil4.item_level_4_name,
			mtr_item_level_1.is_active
		`).
		Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_1_id = mtr_item_level_1.item_level_1_id").
		Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_2_id = mil2.item_level_2_id").
		Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_3_id = mil3.item_level_3_id").
		Where("mtr_item_level_1.item_level_1_id = ?", itemLevelId)

	whereQuery := utils.ApplyFilter(query, filter)
	err := whereQuery.First(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

func (r *ItemLevelImpl) GetItemLevelLookUp(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination, itemClassId int) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemLevel1{}
	responses := []masteritemlevelpayloads.GetItemLevelLookUp{}

	query := tx.Model(&entities).
		Select(`
			mtr_item_level_1.item_level_1_id,
			mtr_item_level_1.item_level_1_code,
			mtr_item_level_1.item_level_1_name,
			mil2.item_level_2_id,
			mil2.item_level_2_code,
			mil2.item_level_2_name,
			mil3.item_level_3_id,
			mil3.item_level_3_code,
			mil3.item_level_3_name,
			mil4.item_level_4_id,
			mil4.item_level_4_code,
			mil4.item_level_4_name,
			mtr_item_level_1.is_active
		`).
		Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_1_id = mtr_item_level_1.item_level_1_id").
		Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_2_id = mil2.item_level_2_id").
		Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_3_id = mil3.item_level_3_id").
		Where("mtr_item_level_1.item_class_id = ?", itemClassId)

	queryFilter := utils.ApplyFilter(query, filter)
	err := queryFilter.Scopes(pagination.Paginate(&entities, &pages, queryFilter)).Order("mtr_item_level_1.item_level_1_id").Scan(&responses).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.Rows = responses

	return pages, nil

}

func (r *ItemLevelImpl) GetItemLevelDropDown(tx *gorm.DB, itemLevel int) ([]masteritemlevelpayloads.GetItemLevelDropdownResponse, *exceptions.BaseErrorResponse) {
	response := []masteritemlevelpayloads.GetItemLevelDropdownResponse{}
	var query *gorm.DB

	switch itemLevel {
	case 1:
		entities := masteritementities.ItemLevel1{}
		query = tx.Model(&entities).
			Select(`
				item_level_1_id AS item_level_id,
				1 AS item_level,
				item_level_1_code AS item_level_code,
				item_level_1_name AS item_level_name
			`)
	case 2:
		entities := masteritementities.ItemLevel2{}
		query = tx.Model(&entities).
			Select(`
				item_level_2_id AS item_level_id,
				2 AS item_level,
				item_level_2_code AS item_level_code,
				item_level_2_name AS item_level_name
			`)
	case 3:
		entities := masteritementities.ItemLevel3{}
		query = tx.Model(&entities).
			Select(`
				item_level_3_id AS item_level_id,
				3 AS item_level,
				item_level_3_code AS item_level_code,
				item_level_3_name AS item_level_name
			`)
	case 4:
		entities := masteritementities.ItemLevel4{}
		query = tx.Model(&entities).
			Select(`
				item_level_4_id AS item_level_id,
				4 AS item_level,
				item_level_4_code AS item_level_code,
				item_level_4_name AS item_level_name
			`)
	default:
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("item_level is unavailable"),
		}
	}

	err := query.Scan(&response).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	for i := 0; i < len(response); i++ {
		response[i].ItemLevelCodeName = response[i].ItemLevelCode + " - " + response[i].ItemLevelName
	}

	return response, nil
}

func (r *ItemLevelImpl) GetAll(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	response := []masteritemlevelpayloads.GetAllItemLevelResponse{}

	entityLevel1 := masteritementities.ItemLevel1{}
	queryLevel1 := tx.Model(&entityLevel1).
		Select(`
			mtr_item_level_1.is_active,
			mtr_item_level_1.item_level_1_id AS item_level_id,
			'1' AS item_level,
			mtr_item_level_1.item_level_1_code AS item_level_code,
			mtr_item_level_1.item_level_1_name AS item_level_name,
			mic.item_class_id,
			mic.item_class_code,
			'' AS item_level_parent
		`).
		Joins("INNER JOIN mtr_item_class mic ON mic.item_class_id = mtr_item_level_1.item_class_id")

	entityLevel2 := masteritementities.ItemLevel2{}
	queryLevel2 := tx.Model(&entityLevel2).
		Select(`
			mtr_item_level_2.is_active,
			item_level_2_id AS item_level_id,
			'2' AS item_level,
			item_level_2_code AS item_level_code,
			item_level_2_name AS item_level_name,
			0 AS item_class_id,
			'' AS item_class_code,
			mil1.item_level_1_code AS item_level_parent
		`).
		Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = mtr_item_level_2.item_level_1_id")

	entityLevel3 := masteritementities.ItemLevel3{}
	queryLevel3 := tx.Model(&entityLevel3).
		Select(`
			mtr_item_level_3.is_active,
			item_level_3_id AS item_level_id,
			'3' AS item_level,
			item_level_3_code AS item_level_code,
			item_level_3_name AS item_level_name,
			0 AS item_class_id,
			'' AS item_class_code,
			mil2.item_level_2_code AS item_level_parent
		`).
		Joins("INNER JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = mtr_item_level_3.item_level_2_id")

	entityLevel4 := masteritementities.ItemLevel4{}
	queryLevel4 := tx.Model(&entityLevel4).
		Select(`
			mtr_item_level_4.is_active,
			item_level_4_id AS item_level_id,
			'4' AS item_level,
			item_level_4_code AS item_level_code,
			item_level_4_name AS item_level_name,
			0 AS item_class_id,
			'' AS item_class_code,
			mil3.item_level_3_code AS item_level_parent
		`).
		Joins("INNER JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = mtr_item_level_4.item_level_3_id")

	unionQuery := tx.Table("(? UNION ALL ? UNION ALL ? UNION ALL ?) A", queryLevel1, queryLevel2, queryLevel3, queryLevel4)
	whereQuery := utils.ApplyFilter(unionQuery, filter)
	err := whereQuery.Scan(&response).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	mapResponse := []map[string]interface{}{}
	for _, data := range response {
		temp := map[string]interface{}{
			"is_active":         data.IsActive,
			"item_level_id":     data.ItemLevelId,
			"item_level":        data.ItemLevel,
			"item_level_code":   data.ItemLevelCode,
			"item_level_name":   data.ItemLevelName,
			"item_class_id":     data.ItemClassId,
			"item_class_code":   data.ItemClassCode,
			"item_level_parent": data.ItemLevelParent,
		}
		mapResponse = append(mapResponse, temp)
	}

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponse, &pages)

	pages.Rows = dataPaginate
	pages.TotalPages = totalPages
	pages.TotalRows = int64(totalRows)

	return pages, nil
}

func (r *ItemLevelImpl) GetById(tx *gorm.DB, itemLevel int, itemLevelId int) (masteritemlevelpayloads.GetItemLevelResponseById, *exceptions.BaseErrorResponse) {
	itemLevelResponse := masteritemlevelpayloads.GetItemLevelResponseById{}
	var query *gorm.DB

	switch itemLevel {
	case 1:
		entities := masteritementities.ItemLevel1{}
		query = tx.Model(&entities).
			Select(`
				is_active,
				item_level_1_id AS item_level_id,
				1 AS item_level,
				item_class_id,
				'' AS item_level_parent,
				item_level_1_code AS item_level_code,
				item_level_1_name AS item_level_name
			`).
			Where("item_level_1_id = ?", itemLevelId)
	case 2:
		entities := masteritementities.ItemLevel2{}
		query = tx.Model(&entities).
			Select(`
				mtr_item_level_2.is_active,
				item_level_2_id as item_level_id,
				2 AS item_level,
				0 AS item_class_id,
				mil1.item_level_1_code AS item_level_parent,
				item_level_2_code AS item_level_code,
				item_level_2_name AS item_level_name
			`).
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = mtr_item_level_2.item_level_1_id").
			Where("item_level_2_id = ?", itemLevelId)
	case 3:
		entities := masteritementities.ItemLevel3{}
		query = tx.Model(&entities).
			Select(`
				mtr_item_level_3.is_active,
				item_level_3_id as item_level_id,
				3 AS item_level,
				0 AS item_class_id,
				mil2.item_level_2_code AS item_level_parent,
				item_level_3_code AS item_level_code,
				item_level_3_name AS item_level_name
			`).
			Joins("INNER JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = mtr_item_level_3.item_level_2_id").
			Where("item_level_3_id = ?", itemLevelId)
	case 4:
		entities := masteritementities.ItemLevel4{}
		query = tx.Model(&entities).
			Select(`
				mtr_item_level_4.is_active,
				item_level_4_id as item_level_id,
				4 AS item_level,
				0 AS item_class_id,
				mil3.item_level_3_code AS item_level_parent,
				item_level_4_code AS item_level_code,
				item_level_4_name AS item_level_name
			`).
			Joins("INNER JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = mtr_item_level_4.item_level_3_id").
			Where("item_level_4_id = ?", itemLevelId)
	default:
		return itemLevelResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("item_level is unavailable"),
		}
	}

	err := query.First(&itemLevelResponse).Error
	if err != nil {
		return itemLevelResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return itemLevelResponse, nil
}

func (r *ItemLevelImpl) Save(tx *gorm.DB, request masteritemlevelpayloads.SaveItemLevelRequest) (bool, *exceptions.BaseErrorResponse) {
	var err error

	if request.ItemLevel > 1 && request.ItemLevelParent == 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("item_level_parent is required"),
		}
	}

	switch request.ItemLevel {
	case 1:
		entities := masteritementities.ItemLevel1{
			IsActive:       request.IsActive,
			ItemLevel1Id:   request.ItemLevelId,
			ItemLevel1Name: request.ItemLevelName,
			ItemLevel1Code: request.ItemLevelCode,
			ItemClassId:    request.ItemClassId,
		}
		err = tx.Save(&entities).Error
	case 2:
		entityLevel1 := masteritementities.ItemLevel1{}
		err = tx.Model(&entityLevel1).Where(masteritementities.ItemLevel1{ItemLevel1Id: request.ItemLevelParent}).First(&entityLevel1).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		entities := masteritementities.ItemLevel2{
			IsActive:       request.IsActive,
			ItemLevel2Id:   request.ItemLevelId,
			ItemLevel1Id:   entityLevel1.ItemLevel1Id,
			ItemLevel2Code: request.ItemLevelCode,
			ItemLevel2Name: request.ItemLevelName,
		}
		err = tx.Save(&entities).Error
	case 3:
		entityLevel2 := masteritementities.ItemLevel2{}
		err = tx.Model(&entityLevel2).Where(masteritementities.ItemLevel2{ItemLevel2Id: request.ItemLevelParent}).First(&entityLevel2).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		entities := masteritementities.ItemLevel3{
			IsActive:       request.IsActive,
			ItemLevel3Id:   request.ItemLevelId,
			ItemLevel2Id:   entityLevel2.ItemLevel2Id,
			ItemLevel3Code: request.ItemLevelCode,
			ItemLevel3Name: request.ItemLevelName,
		}
		err = tx.Save(&entities).Error
	case 4:
		entityLevel3 := masteritementities.ItemLevel3{}
		err = tx.Model(&entityLevel3).Where(masteritementities.ItemLevel3{ItemLevel3Id: request.ItemLevelParent}).First(&entityLevel3).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		entities := masteritementities.ItemLevel4{
			IsActive:       request.IsActive,
			ItemLevel4Id:   request.ItemLevelId,
			ItemLevel3Id:   entityLevel3.ItemLevel3Id,
			ItemLevel4Code: request.ItemLevelCode,
			ItemLevel4Name: request.ItemLevelName,
		}
		err = tx.Save(&entities).Error
	default:
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("item_level is unavailable"),
		}
	}

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

func (r *ItemLevelImpl) ChangeStatus(tx *gorm.DB, itemLevel int, itemLevelId int) (bool, *exceptions.BaseErrorResponse) {
	var err error

	switch itemLevel {
	case 1:
		entities := masteritementities.ItemLevel1{}
		err = tx.Model(&entities).Where(masteritementities.ItemLevel1{ItemLevel1Id: itemLevelId}).First(&entities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		entities.IsActive = !entities.IsActive
		err = tx.Save(&entities).Error
	case 2:
		entities := masteritementities.ItemLevel2{}
		err = tx.Model(&entities).Where(masteritementities.ItemLevel2{ItemLevel2Id: itemLevelId}).First(&entities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		entities.IsActive = !entities.IsActive
		err = tx.Save(&entities).Error
	case 3:
		entities := masteritementities.ItemLevel3{}
		err = tx.Model(&entities).Where(masteritementities.ItemLevel3{ItemLevel3Id: itemLevelId}).First(&entities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		entities.IsActive = !entities.IsActive
		err = tx.Save(&entities).Error
	case 4:
		entities := masteritementities.ItemLevel4{}
		err = tx.Model(&entities).Where(masteritementities.ItemLevel4{ItemLevel4Id: itemLevelId}).First(&entities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		entities.IsActive = !entities.IsActive
		err = tx.Save(&entities).Error
	default:
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("item_level is unavailable"),
		}
	}

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}
