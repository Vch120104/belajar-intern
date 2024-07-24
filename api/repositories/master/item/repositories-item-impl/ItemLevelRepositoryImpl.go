package masteritemrepositoryimpl

import (
	exceptions "after-sales/api/exceptions"
	masteritemlevelrepo "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
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
	mtr_item_level.is_active AS [is_active]`).Joins("LEFT OUTER JOIN mtr_item_level B ON mtr_item_level.item_level_code = B.item_level_parent AND B.item_level = 2").
		Joins("LEFT OUTER JOIN mtr_item_level C ON B.item_level_code = C.item_level_parent AND C.item_level = 3").
		Joins("LEFT OUTER JOIN mtr_item_level D ON C.item_level_code = D.item_level_parent AND D.item_level = 4").
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
	mtr_item_level.is_active AS [is_active]`).Joins("LEFT OUTER JOIN mtr_item_level B ON mtr_item_level.item_level_code = B.item_level_parent AND B.item_level = 2").
		Joins("LEFT OUTER JOIN mtr_item_level C ON B.item_level_code = C.item_level_parent AND C.item_level = 3").
		Joins("LEFT OUTER JOIN mtr_item_level D ON C.item_level_code = D.item_level_parent AND D.item_level = 4").
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

	err := tx.Model(&model).Where(masteritementities.ItemLevel{ItemLevel: strconv.Itoa(itemlevelInt - 1)}).Scan(&result).Error

	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return result, nil
}

func (r *ItemLevelImpl) GetAll(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemLevel{}

	var itemLevelResponse []masteritemlevelpayloads.GetAllItemLevelResponse

	query := tx.Model(entities).Select("mtr_item_level.*,mtr_item_class.*").Joins("left join mtr_item_class on mtr_item_level.item_class_id = mtr_item_class.item_class_id")

	queryFilter := utils.ApplyFilter(query, filter)

	err := queryFilter.Scopes(pagination.Paginate(&entities, &pages, query)).Scan(&itemLevelResponse).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.Rows = itemLevelResponse
	return pages, nil
}

func (r *ItemLevelImpl) GetById(tx *gorm.DB, itemLevelId int) (masteritemlevelpayloads.GetItemLevelResponseById, *exceptions.BaseErrorResponse) {

	var entities masteritementities.ItemLevel
	var itemLevelResponse masteritemlevelpayloads.GetItemLevelResponseById

	rows, err := tx.Model(&entities).
		Where(masteritemlevelpayloads.GetItemLevelResponseById{
			ItemLevelId: itemLevelId,
		}).
		Find(&itemLevelResponse).
		First(&itemLevelResponse).
		Rows()

	if err != nil {
		return itemLevelResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return itemLevelResponse, nil
}

func (r *ItemLevelImpl) Save(tx *gorm.DB, request masteritemlevelpayloads.SaveItemLevelRequest) (bool, *exceptions.BaseErrorResponse) {

	var itemLevelEntities = masteritementities.ItemLevel{
		IsActive:        request.IsActive,
		ItemLevelId:     request.ItemLevelId,
		ItemLevel:       request.ItemLevel,
		ItemClassId:     request.ItemClassId,
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

func (r *ItemLevelImpl) ChangeStatus(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.ItemLevel

	result := tx.Model(&entities).
		Where("item_level_id = ?", Id).
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
