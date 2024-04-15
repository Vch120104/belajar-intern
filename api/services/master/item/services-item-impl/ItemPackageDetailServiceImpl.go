package masteritemserviceimpl

// import (
// 	"after-sales/api/exceptions"
// 	"after-sales/api/helper"

// 	"after-sales/api/payloads/pagination"
// 	masteritemrepository "after-sales/api/repositories/master/item"
// 	masteritemservice "after-sales/api/services/master/item"

// 	"gorm.io/gorm"
// )

// type ItemPackageDetailServiceImpl struct {
// 	ItemPackageDetailRepo masteritemrepository.ItemPackageDetailRepository
// 	DB                    *gorm.DB
// }

// func StartItemPackageDetailService(ItemPackageDetailRepo masteritemrepository.ItemPackageDetailRepository, db *gorm.DB) masteritemservice.ItemPackageDetailService {
// 	return &ItemPackageDetailServiceImpl{
// 		ItemPackageDetailRepo: ItemPackageDetailRepo,
// 		DB:                    db,
// 	}
// }

// func (s *ItemPackageDetailServiceImpl) GetItemPackageDetailByItemPackageId(itemPackageId int, pages pagination.Pagination) pagination.Pagination {
// 	tx := s.DB.Begin()
// 	defer helper.CommitOrRollback(tx)
// 	results, err := s.ItemPackageDetailRepo.GetItemPackageDetailByItemPackageId(tx, itemPackageId, pages)
// 	if err != nil {
// 		panic(exceptions.NewNotFoundError(err.Error()))
// 	}
// 	return results
// }
