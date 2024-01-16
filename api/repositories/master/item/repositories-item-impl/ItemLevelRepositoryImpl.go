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

func OpenItemLevelImpl(db *gorm.DB) masteritemlevelrepo.ItemLevelRepository {
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

func (r *ItemLevelImpl) Save(request masteritemlevelpayloads.SaveItemLevelRequest) (bool, error) {

	var itemLevelEntities = masteritementities.ItemLevel{
		IsActive:        request.IsActive,
		ItemLevel:       request.ItemLevel,
		ItemClassCode:   request.ItemClassCode,
		ItemLevelParent: request.ItemLevelParent,
		ItemLevelCode:   request.ItemLevelCode,
		ItemLevelName:   request.ItemLevelName,
	}

	rows, err := r.DB.Model(&itemLevelEntities).
		Create(&itemLevelEntities).
		Rows()

	if err != nil {
		return false, err
	}

	defer rows.Close()

	return true, nil
}

func (r *ItemLevelImpl) Update(request masteritemlevelpayloads.SaveItemLevelRequest) (bool, error) {

	var itemLevelEntities = masteritementities.ItemLevel{
		IsActive:        request.IsActive,
		ItemLevel:       request.ItemLevel,
		ItemClassCode:   request.ItemClassCode,
		ItemLevelParent: request.ItemLevelParent,
		ItemLevelCode:   request.ItemLevelCode,
		ItemLevelName:   request.ItemLevelName,
	}

	rows, err := r.DB.Model(&itemLevelEntities).
		Where(masteritemlevelpayloads.SaveItemLevelRequest{
			ItemLevelId: request.ItemLevelId,
		}).
		Updates(&itemLevelEntities).
		Rows()

	if err != nil {
		return false, err
	}

	defer rows.Close()

	return true, nil
}

func (r *ItemLevelImpl) GetById(itemLevelId int) (masteritemlevelpayloads.GetItemLevelResponse, error) {

	var entities masteritementities.ItemLevel
	var itemLevelResponse masteritemlevelpayloads.GetItemLevelResponse

	rows, err := r.DB.Model(&entities).
		Where(masteritemlevelpayloads.GetItemLevelResponse{
			ItemLevelId: itemLevelId,
		}).
		Find(&itemLevelResponse).
		Scan(&itemLevelResponse).
		Rows()

	if err != nil {
		return itemLevelResponse, err
	}

	defer rows.Close()

	return itemLevelResponse, nil
}

func (r *ItemLevelImpl) GetAll(request masteritemlevelpayloads.GetAllItemLevelResponse, pages pagination.Pagination) (pagination.Pagination, error) {
	var entities []masteritementities.ItemLevel
	var itemLevelResponse []masteritemlevelpayloads.GetAllItemLevelResponse

	tempRows := r.DB.
		Model(&masteritementities.ItemLevel{}).
		Where("item_level like ?", "%"+request.ItemLevel+"%").
		Where("item_level_code like ?", "%"+request.ItemLevelCode+"%").
		Where("item_level_name like ?", "%"+request.ItemLevelName+"%").
		Where("item_class like ?", "%"+request.ItemClassCode+"%").
		Where("item_level_parent like ?", "%"+request.ItemLevelParent+"%")

	if request.IsActive != "" {
		tempRows = tempRows.Where("is_active = ?", request.IsActive)
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

func (r *ItemLevelImpl) ChangeStatus(itemLevelId int) (masteritemlevelpayloads.GetItemLevelResponse, error) {
	var entities masteritementities.ItemLevel
	var itemLevelPayloads masteritemlevelpayloads.GetItemLevelResponse

	rows, err := r.DB.Model(&entities).
		Where(masteritemlevelpayloads.GetItemLevelResponse{
			ItemLevelId: itemLevelId,
		}).
		Update("is_active", gorm.Expr("1 ^ is_active")).
		Rows()

	if err != nil {
		log.Panic((err.Error()))
	}

	rows, err = r.DB.Model(&entities).
		Where(masteritemlevelpayloads.GetItemLevelResponse{
			ItemLevelId: itemLevelId,
		}).
		Find(&itemLevelPayloads).
		Scan(&itemLevelPayloads).
		Rows()

	if err != nil {
		return itemLevelPayloads, err
	}

	defer rows.Close()

	return itemLevelPayloads, nil
}
