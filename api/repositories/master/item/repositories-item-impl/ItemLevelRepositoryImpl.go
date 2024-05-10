package masteritemrepositoryimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritemlevelrepo "after-sales/api/repositories/master/item"
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

// GetItemLevelDropDown implements masteritemrepository.ItemLevelRepository.
func (r *ItemLevelImpl) GetItemLevelDropDown(tx *gorm.DB, itemLevel string) ([]masteritemlevelpayloads.GetItemLevelDropdownResponse, *exceptionsss_test.BaseErrorResponse) {
	model := masteritementities.ItemLevel{}
	result := []masteritemlevelpayloads.GetItemLevelDropdownResponse{}

	itemlevelInt, _ := strconv.Atoi(itemLevel)

	fmt.Print("adawd", itemlevelInt)

	err := tx.Model(&model).Where(masteritementities.ItemLevel{ItemLevel: strconv.Itoa(itemlevelInt - 1)}).Scan(&result).Error

	if err != nil {
		return result, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return result, nil
}

func (r *ItemLevelImpl) GetAll(tx *gorm.DB, request masteritemlevelpayloads.GetAllItemLevelResponse, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	var entities []masteritementities.ItemLevel
	var itemLevelResponse []masteritemlevelpayloads.GetAllItemLevelResponse

	tempRows := tx.
		Select("mtr_item_level.is_active, mtr_item_level.item_level_id, mtr_item_level.item_level, mtr_item_level.item_level_code, mtr_item_level.item_level_name, mtr_item_level.item_class_id, mtr_item_level.item_level_parent, mtr_item_class.item_class_code").
		Table("mtr_item_level").
		Joins("JOIN mtr_item_class ON mtr_item_level.item_class_id = mtr_item_class.item_class_id").
		Where("mtr_item_level.item_level LIKE ?", "%"+request.ItemLevel+"%").
		Where("mtr_item_level.item_level_code LIKE ?", "%"+request.ItemLevelCode+"%").
		Where("mtr_item_level.item_level_name LIKE ?", "%"+request.ItemLevelName+"%").
		Where("mtr_item_class.item_class_code LIKE ?", "%"+request.ItemClassCode+"%").
		Where("mtr_item_level.item_level_parent LIKE ?", "%"+request.ItemLevelParent+"%")

	if request.IsActive != "" {
		tempRows = tempRows.Where("mtr_item_level.is_active = ?", request.IsActive)
	}

	rows, err := tempRows.
		Scopes(pagination.Paginate(entities, &pages, tempRows)).
		Scan(&itemLevelResponse).
		Rows()

	if err != nil {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	pages.Rows = itemLevelResponse
	return pages, nil
}

func (r *ItemLevelImpl) GetById(tx *gorm.DB, itemLevelId int) (masteritemlevelpayloads.GetItemLevelResponseById, *exceptionsss_test.BaseErrorResponse) {

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
		return itemLevelResponse, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return itemLevelResponse, nil
}

func (r *ItemLevelImpl) Save(tx *gorm.DB, request masteritemlevelpayloads.SaveItemLevelRequest) (bool, *exceptionsss_test.BaseErrorResponse) {

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
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {

			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return true, nil
}

func (r *ItemLevelImpl) ChangeStatus(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities masteritementities.ItemLevel

	result := tx.Model(&entities).
		Where("item_level_id = ?", Id).
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
