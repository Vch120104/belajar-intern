package transactionsparepartrepository

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BinningListRepository interface {
	GetBinningListById(db *gorm.DB, BinningStockId int) (transactionsparepartpayloads.BinningListGetByIdResponse, *exceptions.BaseErrorResponse)
	GetAllBinningListWithPagination(db *gorm.DB, rdb *redis.Client, filter []utils.FilterCondition, pagination pagination.Pagination, ctx context.Context) (pagination.Pagination, *exceptions.BaseErrorResponse)
	InsertBinningListHeader(db *gorm.DB, payloads transactionsparepartpayloads.BinningListInsertPayloads) (transactionsparepartentities.BinningStock, *exceptions.BaseErrorResponse)
	UpdateBinningListHeader(db *gorm.DB, payloads transactionsparepartpayloads.BinningListSavePayload) (transactionsparepartentities.BinningStock, *exceptions.BaseErrorResponse)
	GetBinningListDetailById(db *gorm.DB, BinningDetailId int) (transactionsparepartpayloads.BinningListGetByIdResponses, *exceptions.BaseErrorResponse)
	GetAllBinningListDetailWithPagination(db *gorm.DB, filter []utils.FilterCondition, pagination pagination.Pagination, binningListId int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	InsertBinningListDetail(db *gorm.DB, payloads transactionsparepartpayloads.BinningListDetailInsertPayloads) (transactionsparepartentities.BinningStockDetail, *exceptions.BaseErrorResponse)
}
