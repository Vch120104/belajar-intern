package transactionsparepartrepository

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemWarehouseTransferRequestRepository interface {
	InsertWhTransferRequestHeader(*gorm.DB, transactionsparepartpayloads.InsertItemWarehouseTransferRequest) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse)
	InsertWhTransferRequestDetail(*gorm.DB, transactionsparepartpayloads.InsertItemWarehouseTransferDetailRequest) (transactionsparepartentities.ItemWarehouseTransferRequestDetail, *exceptions.BaseErrorResponse)
	UpdateWhTransferRequest(*gorm.DB, transactionsparepartpayloads.UpdateItemWarehouseTransferRequest, int) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse)
	UpdateWhTransferRequestDetail(*gorm.DB, transactionsparepartpayloads.UpdateItemWarehouseTransferRequestDetailRequest, int) (transactionsparepartentities.ItemWarehouseTransferRequestDetail, *exceptions.BaseErrorResponse)
	SubmitWhTransferRequest(*gorm.DB, int, transactionsparepartpayloads.SubmitItemWarehouseTransferRequest) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse)
	GetAllDetailTransferRequest(*gorm.DB, int, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdTransferRequest(*gorm.DB, int) (transactionsparepartpayloads.GetByIdItemWarehouseTransferRequestResponse, *exceptions.BaseErrorResponse)
	GetByIdTransferRequestDetail(*gorm.DB, int) (transactionsparepartpayloads.GetByIdItemWarehouseTransferRequestDetailResponse, *exceptions.BaseErrorResponse)
	GetAllWhTransferRequest(*gorm.DB, pagination.Pagination, []utils.FilterCondition, map[string]string) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetTransferRequestLookUp(*gorm.DB, pagination.Pagination, []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetTransferRequestDetailLookUp(*gorm.DB, int, pagination.Pagination, []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	DeleteHeaderTransferRequest(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	DeleteDetail(*gorm.DB, int, transactionsparepartpayloads.DeleteDetailItemWarehouseTransferRequest) (bool, *exceptions.BaseErrorResponse)
}
