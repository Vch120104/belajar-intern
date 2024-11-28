package masteritemserviceimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ItemGroupServiceImpl struct {
	repository masteritemrepository.ItemGroupRepository
	DB         *gorm.DB
	rdb        *redis.Client
}

func (i *ItemGroupServiceImpl) GetAllItemGroupWithPagination(internalFilter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
	results, err := i.repository.GetAllItemGroupWithPagination(tx, internalFilter, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (i *ItemGroupServiceImpl) GetAllItemGroup(code string) ([]masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
	results, err := i.repository.GetAllItemGroup(tx, code)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (i *ItemGroupServiceImpl) GetItemGroupById(id int) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
	results, err := i.repository.GetItemGroupById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (i *ItemGroupServiceImpl) DeleteItemGroupById(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
	results, err := i.repository.DeleteItemGroupById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (i *ItemGroupServiceImpl) UpdateItemGroupById(payload masteritempayloads.ItemGroupUpdatePayload, id int) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
	results, err := i.repository.UpdateItemGroupById(tx, payload, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (i *ItemGroupServiceImpl) UpdateStatusItemGroupById(id int) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
	results, err := i.repository.UpdateStatusItemGroupById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (i *ItemGroupServiceImpl) GetItemGroupByMultiId(multiId string) ([]masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
	results, err := i.repository.GetItemGroupByMultiId(tx, multiId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (i *ItemGroupServiceImpl) NewItemGroup(payload masteritempayloads.NewItemGroupPayload) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
	results, err := i.repository.NewItemGroup(tx, payload)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func NewItemGroupServiceImpl(repo masteritemrepository.ItemGroupRepository, DB *gorm.DB, rdb *redis.Client) masteritemservice.ItemGroupService {
	return &ItemGroupServiceImpl{
		repository: repo,
		DB:         DB,
		rdb:        rdb,
	}
}
