package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	"after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type ItemGroupRepositoryImpl struct {
}

func (i *ItemGroupRepositoryImpl) GetItemGroupByMultiId(db *gorm.DB, multiId string) ([]masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	Ids := strings.Split(multiId, ",")
	var entities []masteritementities.ItemGroup
	err := db.Model(&masteritementities.ItemGroup{}).Where("item_group_id in ?", Ids).Scan(&entities).Error
	if err != nil {
		if len(entities) == 0 {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "item group with that id is not found please check input",
			}
		}
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get item group by multi id please check input",
		}
	}
	return entities, nil
}

func (i *ItemGroupRepositoryImpl) NewItemGroup(db *gorm.DB, payload masteritempayloads.NewItemGroupPayload) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	NewItemGroup := masteritementities.ItemGroup{
		IsActive: payload.IsActive,
		//ItemGroupId:     0,
		ItemGroupCode:   payload.ItemGroupCode,
		ItemGroupName:   payload.ItemGroupName,
		IsItemSparepart: payload.IsItemSparepart,
	}
	err := db.Create(&NewItemGroup).Scan(&NewItemGroup).Error
	if err != nil {
		return NewItemGroup, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to create item group please check input",
			Err:        err,
		}
	}
	return NewItemGroup, nil
}

func (i *ItemGroupRepositoryImpl) UpdateItemGroupById(db *gorm.DB, payload masteritempayloads.ItemGroupUpdatePayload, id int) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	itemGroupToUpdate := masteritementities.ItemGroup{}
	err := db.Model(&itemGroupToUpdate).Where(masteritementities.ItemGroup{ItemGroupId: id}).First(&itemGroupToUpdate).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return itemGroupToUpdate, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "item group with this id is not found",
			}
		}
		return itemGroupToUpdate, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "error on getting item group please contact administrator",
		}
	}

	if payload.ItemGroupName != "" {
		itemGroupToUpdate.ItemGroupName = payload.ItemGroupName
	}
	if payload.ItemGroupCode != "" {
		itemGroupToUpdate.ItemGroupCode = payload.ItemGroupCode
	}
	err = db.Save(&itemGroupToUpdate).Error
	if err != nil {
		return itemGroupToUpdate, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to update item group",
		}
	}
	return itemGroupToUpdate, nil
}

func (i *ItemGroupRepositoryImpl) UpdateStatusItemGroupById(tx *gorm.DB, id int) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemGroup{}
	err := tx.Model(&entities).Where(masteritementities.ItemGroup{ItemGroupId: id}).First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "item group is not found please check input",
				Err:        err,
			}
		}
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get item group please contact your administrator",
			Err:        err,
		}
	}
	entities.IsActive = !entities.IsActive
	err = tx.Save(&entities).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to save item group please check input",
			Err:        err,
		}
	}
	return entities, nil
}

func (i *ItemGroupRepositoryImpl) GetAllItemGroup(db *gorm.DB, code string) ([]masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	Entities := []masteritementities.ItemGroup{}
	filterCode := strings.Split(code, ",")
	myQuery := db.Model(&masteritementities.ItemGroup{}).
		Select("*")
	if len(filterCode) != 0 && filterCode[0] != "" {
		myQuery.Where("item_group_code IN ?", filterCode)
	}
	err := myQuery.Scan(&Entities).Error
	if err != nil {
		if len(Entities) == 0 {
			return Entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "failed to get item group data not exist",
				Err:        err,
			}
		}
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get item group please contact admin",
			Err:        err,
		}
	}
	return Entities, nil
}

func (i *ItemGroupRepositoryImpl) GetItemGroupById(db *gorm.DB, id int) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemGroup{}
	err := db.Model(&entities).Where(masteritementities.ItemGroup{ItemGroupId: id}).First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "item group with that id is not found please check input",
			}
		}
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "there is error when getting item group pleae check input",
		}
	}
	return entities, nil
}

func (i *ItemGroupRepositoryImpl) DeleteItemGroupById(db *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	err := db.Delete(&masteritementities.ItemGroup{}, masteritementities.ItemGroup{ItemGroupId: id}).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to delete item group by id please check input",
			Err:        err,
		}
	}
	return true, nil
}

func (i *ItemGroupRepositoryImpl) GetAllItemGroupWithPagination(db *gorm.DB, internalFilter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	Entities := masteritementities.ItemGroup{}
	Responds := []masteritempayloads.ItemGroupGetAllResponses{}
	joinTable := db.Model(&Entities).Select("*")
	WhereQuery := utils.ApplyFilter(joinTable, internalFilter)
	err := WhereQuery.Scopes(pagination.Paginate(&pages, WhereQuery)).Scan(&Responds).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get all item group with filter pagination.",
			Err:        err,
		}
	}
	if len(Responds) == 0 {
		pages.Rows = []string{}
		return pages, nil
	}
	pages.Rows = Responds
	return pages, nil
}

func NewItemGroupRepositoryImpl() masteritemrepository.ItemGroupRepository {
	return &ItemGroupRepositoryImpl{}
}
