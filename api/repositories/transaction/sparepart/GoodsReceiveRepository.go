package transactionsparepartrepository

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
	"gorm.io/gorm"
)

type GoodsReceiveRepository interface {
	GetAllGoodsReceive(db *gorm.DB, filter []utils.FilterCondition, pagination pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetGoodsReceiveById(db *gorm.DB, GoodsReceiveId int) (transactionsparepartpayloads.GoodsReceivesGetByIdResponses, *exceptions.BaseErrorResponse)
	InsertGoodsReceive(db *gorm.DB, payloads transactionsparepartpayloads.GoodsReceiveInsertPayloads) (transactionsparepartentities.GoodsReceive, *exceptions.BaseErrorResponse)
	UpdateGoodsReceive(db *gorm.DB, payloads transactionsparepartpayloads.GoodsReceiveUpdatePayloads, GoodsReceiveId int) (transactionsparepartentities.GoodsReceive, *exceptions.BaseErrorResponse)
	InsertGoodsReceiveDetail(db *gorm.DB, payloads transactionsparepartpayloads.GoodsReceiveDetailInsertPayloads) (transactionsparepartentities.GoodsReceiveDetail, *exceptions.BaseErrorResponse)
	UpdateGoodsReceiveDetail(db *gorm.DB, payloads transactionsparepartpayloads.GoodsReceiveDetailUpdatePayloads, DetailId int) (bool, *exceptions.BaseErrorResponse)
	LocationItemGoodsReceive(db *gorm.DB, filter []utils.FilterCondition, pagination pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SubmitGoodsReceive(db *gorm.DB, GoodsReceiveId int) (bool, *exceptions.BaseErrorResponse)
	DeleteGoodsReceive(db *gorm.DB, goodsReceivesId int) (bool, *exceptions.BaseErrorResponse)
}
