package transactionsparepartservice

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
)

type GoodsReceiveService interface {
	GetAllGoodsReceive(filter []utils.FilterCondition, paginations pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetGoodsReceiveById(GoodsReceiveId int) (transactionsparepartpayloads.GoodsReceivesGetByIdResponses, *exceptions.BaseErrorResponse)
	InsertGoodsReceive(payloads transactionsparepartpayloads.GoodsReceiveInsertPayloads) (transactionsparepartentities.GoodsReceive, *exceptions.BaseErrorResponse)
	UpdateGoodsReceive(payloads transactionsparepartpayloads.GoodsReceiveUpdatePayloads, GoodsReceiveId int) (transactionsparepartentities.GoodsReceive, *exceptions.BaseErrorResponse)
	InsertGoodsReceiveDetail(payloads transactionsparepartpayloads.GoodsReceiveDetailInsertPayloads) (transactionsparepartentities.GoodsReceiveDetail, *exceptions.BaseErrorResponse)
}
