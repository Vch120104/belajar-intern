package masteritemrepositoryimpl

import (
	masteritemlevelrepo "after-sales/api/repositories/master/item"
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

func (r *ItemLevelImpl) GetAll(tx *gorm.DB, request masteritemlevelpayloads.GetAllItemLevelResponse, pages pagination.Pagination) (pagination.Pagination, error) {
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
		return pages, err
	}

	defer rows.Close()

	pages.Rows = itemLevelResponse
	return pages, nil
}

func (r *ItemLevelImpl) GetById(tx *gorm.DB, itemLevelId int) (masteritemlevelpayloads.GetItemLevelResponseById, error) {

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
		return itemLevelResponse, err
	}

	defer rows.Close()

	return itemLevelResponse, nil
}

func (r *ItemLevelImpl) Save(tx *gorm.DB, request masteritemlevelpayloads.SaveItemLevelRequest) (bool, error) {

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
		return false, err
	}

	return true, nil
}

func (r *ItemLevelImpl) ChangeStatus(tx *gorm.DB, Id int) (bool, error) {
	var entities masteritementities.ItemLevel

	result := tx.Model(&entities).
		Where("item_level_id = ?", Id).
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
