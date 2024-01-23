package masteritemrepositoryimpl

import (
	masteritemlevelrepo "after-sales/api/repositories/master/item"
	// masteritemlevelservice "after-sales/api/services/master/item_level"
	masteritementities "after-sales/api/entities/master/item"
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"

	"log"

	"gorm.io/gorm"
)

type ItemLevelImpl struct {
	DB *gorm.DB
}

func StartItemLevelRepositoryImpl(db *gorm.DB) masteritemlevelrepo.ItemLevelRepository {
	return &ItemLevelImpl{DB: db}
}

func (r *ItemLevelImpl) WithTrx(Trxhandle *gorm.DB) masteritemlevelrepo.ItemLevelRepository {
	if Trxhandle == nil {
		log.Println("Transaction Database Not Found")
		return r
	}
	r.DB = Trxhandle
	return r
}

func (r *ItemLevelImpl) GetAll(request masteritemlevelpayloads.GetAllItemLevelResponse, pages pagination.Pagination) (pagination.Pagination, error) {
	var entities []masteritementities.ItemLevel
	var itemLevelResponse []masteritemlevelpayloads.GetAllItemLevelResponse

	tempRows := r.DB.
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

func (r *ItemLevelImpl) GetById(itemLevelId int) (masteritemlevelpayloads.GetItemLevelResponseById, error) {

	var entities masteritementities.ItemLevel
	var itemLevelResponse masteritemlevelpayloads.GetItemLevelResponseById

	rows, err := r.DB.Model(&entities).
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

func (r *ItemLevelImpl) Save(request masteritemlevelpayloads.SaveItemLevelRequest) (bool, error) {

	var itemLevelEntities = masteritementities.ItemLevel{
		IsActive:        request.IsActive,
		ItemLevelId:     request.ItemLevelId,
		ItemLevel:       request.ItemLevel,
		ItemClassId:     request.ItemClassId,
		ItemLevelParent: request.ItemLevelParent,
		ItemLevelCode:   request.ItemLevelCode,
		ItemLevelName:   request.ItemLevelName,
	}

	err := r.DB.Save(&itemLevelEntities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *ItemLevelImpl) ChangeStatus(Id int) (bool, error) {
	var entities masteritementities.ItemLevel

	result := r.DB.Model(&entities).
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

	result = r.DB.Save(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}
