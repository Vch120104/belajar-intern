package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemPackageDetailRepositoryImpl struct {
}

func StartItemPackageDetailRepositoryImpl() masteritemrepository.ItemPackageDetailRepository {
	return &ItemPackageDetailRepositoryImpl{}
}

func (r *ItemPackageDetailRepositoryImpl) GetItemPackageDetailByItemPackageId(tx *gorm.DB, itemPackageId int, pages pagination.Pagination) (pagination.Pagination, error) {

	entities := masteritementities.ItemPackageDetail{}

	var responses []masteritempayloads.ItemPackageDetailResponse
	tableStruct := masteritempayloads.ItemPackageDetailResponse{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	rows, err := joinTable.Scopes(pagination.Paginate(&entities, &pages, joinTable)).Scan(&responses).Rows()

	if len(responses) == 0 {
		return pages, gorm.ErrRecordNotFound
	}

	if err != nil {
		return pages, err
	}

	defer rows.Close()

	pages.Rows = responses

	return pages, nil
}
