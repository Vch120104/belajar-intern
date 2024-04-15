package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"net/http"

	"gorm.io/gorm"
)

type ItemPackageDetailRepositoryImpl struct {
}

func StartItemPackageDetailRepositoryImpl() masteritemrepository.ItemPackageDetailRepository {
	return &ItemPackageDetailRepositoryImpl{}
}

func (r *ItemPackageDetailRepositoryImpl) GetItemPackageDetailByItemPackageId(tx *gorm.DB, itemPackageId int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {

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
		).
		Joins("ItemPackageDetail", tx.Select("1")).
		Joins("ItemPackageDetail.Item", tx.Select("1")).
		Joins("ItemPackageDetail.Item.ItemClass", tx.Select("1"))

	rows, err := query.Scopes(pagination.Paginate(&model, &pages, query)).Scan(&responses).Rows()

	if err != nil {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	pages.Rows = responses

	return pages, nil
}

func (r *ItemPackageDetailRepositoryImpl) SaveItemPackageDetailByItemPackageId(tx *gorm.DB, itemPackageId int, req masteritempayloads.ItemSubstitutePostPayloads) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masteritementities.ItemSubstitute{
		SubstituteTypeCode: req.SubstituteTypeCode,
		ItemSubstituteId:   req.ItemSubstituteId,
		EffectiveDate:      req.EffectiveDate,
		ItemId:             req.ItemId,
	}
	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	return true, nil
}
