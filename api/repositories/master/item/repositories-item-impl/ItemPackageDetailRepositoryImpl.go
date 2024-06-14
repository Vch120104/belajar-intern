package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type ItemPackageDetailRepositoryImpl struct {
}

// ChangeStatusItemPackageDetail implements masteritemrepository.ItemPackageDetailRepository.
func (r *ItemPackageDetailRepositoryImpl) ChangeStatusItemPackageDetail(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.ItemPackageDetail

	result := tx.Model(&entities).
		Where("item_package_detail_id = ?", id).
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

func StartItemPackageDetailRepositoryImpl() masteritemrepository.ItemPackageDetailRepository {
	return &ItemPackageDetailRepositoryImpl{}
}

func (r *ItemPackageDetailRepositoryImpl) GetItemPackageDetailByItemPackageId(tx *gorm.DB, itemPackageId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	// entities := masteritementities.ItemPackageDetail{}
	model := masteritementities.ItemPackage{}
	var responses []masteritempayloads.ItemPackageDetailResponse
	// tableStruct := masteritempayloads.ItemPackageDetailResponse{}

	query := tx.Model(&model).
		Select(
			"item_package_detail_id",
			"ItemPackageDetail.is_active is_active",
			"ItemPackageDetail.item_package_id item_package_id",
			"ItemPackageDetail__Item.item_id item_id",
			"ItemPackageDetail__Item.item_code item_code",
			"ItemPackageDetail__Item.item_name item_name",
			"ItemPackageDetail__Item.item_class_id item_class_id",
			"ItemPackageDetail__Item__ItemClass.item_class_code item_class_code",
			"ItemPackageDetail.quantity quantity",
		).Where(masteritementities.ItemPackage{ItemPackageId: itemPackageId}).
		InnerJoins("ItemPackageDetail", tx.Select("1")).
		InnerJoins("ItemPackageDetail.Item", tx.Select("1")).
		InnerJoins("ItemPackageDetail.Item.ItemClass", tx.Select("1"))

	rows, err := query.Scopes(pagination.Paginate(&model, &pages, query)).Scan(&responses).Rows()

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	defer rows.Close()

	pages.Rows = responses

	return pages, nil
}

func (r *ItemPackageDetailRepositoryImpl) GetItemPackageDetailById(tx *gorm.DB, itemPackageDetailId int) (masteritempayloads.ItemPackageDetailResponse, *exceptions.BaseErrorResponse) {
	model := masteritementities.ItemPackage{}
	response := masteritempayloads.ItemPackageDetailResponse{}

	query := tx.Model(&model).
		Select(
			"item_package_detail_id",
			"ItemPackageDetail.is_active is_active",
			"ItemPackageDetail.item_package_id item_package_id",
			"ItemPackageDetail__Item.item_id item_id",
			"ItemPackageDetail__Item.item_code item_code",
			"ItemPackageDetail__Item.item_name item_name",
			"ItemPackageDetail__Item.item_class_id item_class_id",
			"ItemPackageDetail__Item__ItemClass.item_class_code item_class_code",
			"ItemPackageDetail.quantity quantity",
		).
		InnerJoins("ItemPackageDetail", tx.Where(masteritementities.ItemPackageDetail{
			ItemPackageDetailId: itemPackageDetailId,
		})).
		InnerJoins("ItemPackageDetail.Item", tx.Select("1")).
		InnerJoins("ItemPackageDetail.Item.ItemClass", tx.Select("1"))

	rows, err := query.First(&response).Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	return response, nil
}

func (r *ItemPackageDetailRepositoryImpl) CreateItemPackageDetailByItemPackageId(tx *gorm.DB, req masteritempayloads.SaveItemPackageDetail) (bool, *exceptions.BaseErrorResponse) {

	entities := masteritementities.ItemPackageDetail{
		IsActive:      req.IsActive,
		ItemPackageId: req.ItemPackageId,
		ItemId:        req.ItemId,
		Quantity:      req.Quantity,
	}

	err := tx.Create(&entities).Error

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

func (r *ItemPackageDetailRepositoryImpl) UpdateItemPackageDetail(tx *gorm.DB, req masteritempayloads.SaveItemPackageDetail) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemPackageDetail{
		ItemPackageDetailId: req.ItemPackageDetailId,
		Quantity:            req.Quantity,
	}

	err := tx.Updates(&entities).Error

	if err != nil {

		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}
