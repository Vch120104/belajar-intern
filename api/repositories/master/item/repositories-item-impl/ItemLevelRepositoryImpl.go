package masteritemrepositoryimpl

import (
	exceptions "after-sales/api/exceptions"
	masteritemlevelrepo "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

// GetItemLevelLookUpbyId implements masteritemrepository.ItemLevelRepository.
func (r *ItemLevelImpl) GetItemLevelLookUpbyId(tx *gorm.DB, itemLevelId int) (masteritemlevelpayloads.GetItemLevelLookUp, *exceptions.BaseErrorResponse) {
	model := masteritementities.ItemLevel{}
	responses := masteritemlevelpayloads.GetItemLevelLookUp{}
	err := tx.Model(model).
		Select(`mtr_item_level.item_level_code AS [item_level_1],
	mtr_item_level.item_level_name AS [item_level_1_name],
	B.item_level_code AS [item_level_2],
	B.item_level_name AS [item_level_2_name],
	C.item_level_code AS [item_level_3],
	C.item_level_name AS [item_level_3_name],
	D.item_level_code AS [item_level_4],
	D.item_level_name AS [item_level_4_name],
	mtr_item_level.item_level_id AS [item_level_id],
	mtr_item_level.is_active AS [is_active]`).Joins("LEFT OUTER JOIN mtr_item_level B ON mtr_item_level.item_level_id = B.item_level_parent AND B.item_level = 2").
		Joins("LEFT OUTER JOIN mtr_item_level C ON B.item_level_id = C.item_level_parent AND C.item_level = 3").
		Joins("LEFT OUTER JOIN mtr_item_level D ON C.item_level_id = D.item_level_parent AND D.item_level = 4").
		Where(masteritementities.ItemLevel{ItemLevelId: itemLevelId}).Find(&responses).Error

	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if responses == (masteritemlevelpayloads.GetItemLevelLookUp{}) {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	return responses, nil
}

// GetItemLevelLookUp implements masteritemrepository.ItemLevelRepository.
func (r *ItemLevelImpl) GetItemLevelLookUp(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination, itemClassId int) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	model := masteritementities.ItemLevel{}
	responses := []masteritemlevelpayloads.GetItemLevelLookUp{}

	query := tx.Model(model).
		Select(`mtr_item_level.item_level_code AS [item_level_1],
	mtr_item_level.item_level_name AS [item_level_1_name],
	B.item_level_code AS [item_level_2],
	B.item_level_name AS [item_level_2_name],
	C.item_level_code AS [item_level_3],
	C.item_level_name AS [item_level_3_name],
	D.item_level_code AS [item_level_4],
	D.item_level_name AS [item_level_4_name],
	mtr_item_level.item_level_id AS [item_level_id],
	mtr_item_level.is_active AS [is_active]`).Joins("LEFT OUTER JOIN mtr_item_level B ON mtr_item_level.item_level_id = B.item_level_parent AND B.item_level = 2").
		Joins("LEFT OUTER JOIN mtr_item_level C ON B.item_level_id = C.item_level_parent AND C.item_level = 3").
		Joins("LEFT OUTER JOIN mtr_item_level D ON C.item_level_id = D.item_level_parent AND D.item_level = 4").
		Where(masteritementities.ItemLevel{ItemClassId: itemClassId})

	queryFilter := utils.ApplyFilter(query, filter)

	err := queryFilter.Scopes(pagination.Paginate(&model, &pages, queryFilter)).Order("mtr_item_level.item_level_id").Scan(&responses).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.Rows = responses

	return pages, nil

}

// GetItemLevelDropDown implements masteritemrepository.ItemLevelRepository.
func (r *ItemLevelImpl) GetItemLevelDropDown(tx *gorm.DB, itemLevel string) ([]masteritemlevelpayloads.GetItemLevelDropdownResponse, *exceptions.BaseErrorResponse) {
	model := masteritementities.ItemLevel{}
	result := []masteritemlevelpayloads.GetItemLevelDropdownResponse{}

	itemlevelInt, _ := strconv.Atoi(itemLevel)

	err := tx.Model(&model).Select("mtr_item_level.*,CONCAT(item_level_code , ' - ',item_level_name)AS item_level_name").Where(masteritementities.ItemLevel{ItemLevel: strconv.Itoa(itemlevelInt - 1)}).Scan(&result).Error

	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return result, nil
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

	//GET ITEM CLASS LEVEL PARENT, IF CREATE ITEM LEVEL > 1
	itemleveltoInt := request.ItemLevel
	itemClassId := request.ItemClassId

	model := masteritementities.ItemLevel{}

	fmt.Println(itemleveltoInt, request.ItemLevelParent)

	if itemleveltoInt > 1 {
		if err := tx.Model(model).
			Select("mtr_item_level.item_class_id").
			Where(masteritementities.ItemLevel{ItemLevel: strconv.Itoa(itemleveltoInt - 1), ItemLevelId: request.ItemLevelParent}).
			First(&itemClassId).Error; err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	//

	var itemLevelEntities = masteritementities.ItemLevel{
		IsActive:        request.IsActive,
		ItemLevelId:     request.ItemLevelId,
		ItemLevel:       strconv.Itoa(request.ItemLevel),
		ItemClassId:     itemClassId,
		ItemLevelParent: request.ItemLevelParent,
		ItemLevelCode:   request.ItemLevelCode,
		ItemLevelName:   request.ItemLevelName,
	}

	err := tx.Save(&itemLevelEntities).Error

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
