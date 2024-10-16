package transactionsparepartservice

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
	"context"
)

type BinningListService interface {
	GetBinningListById(BinningStockId int) (transactionsparepartpayloads.BinningListGetByIdResponse, *exceptions.BaseErrorResponse)
	GetAllBinningListWithPagination(filter []utils.FilterCondition, pagination pagination.Pagination, ctx context.Context) (pagination.Pagination, *exceptions.BaseErrorResponse)
	InsertBinningListHeader(payloads transactionsparepartpayloads.BinningListInsertPayloads) (transactionsparepartentities.BinningStock, *exceptions.BaseErrorResponse)
	UpdateBinningListHeader(payloads transactionsparepartpayloads.BinningListSavePayload) (transactionsparepartentities.BinningStock, *exceptions.BaseErrorResponse)
	GetBinningListDetailById(BinningDetailSystemNumber int) (transactionsparepartpayloads.BinningListGetByIdResponses, *exceptions.BaseErrorResponse)
	GetAllBinningListDetailWithPagination(filter []utils.FilterCondition, pagination pagination.Pagination, binningListId int) (pagination.Pagination, *exceptions.BaseErrorResponse)
}
