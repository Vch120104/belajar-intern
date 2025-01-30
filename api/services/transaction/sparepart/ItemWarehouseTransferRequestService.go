package transactionsparepartservice

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"

	"github.com/xuri/excelize/v2"
)

type ItemWarehouseTransferRequestService interface {
	InsertWhTransferRequestHeader(transactionsparepartpayloads.InsertItemWarehouseTransferRequest) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse)
	InsertWhTransferRequestDetail(transactionsparepartpayloads.InsertItemWarehouseTransferDetailRequest) (transactionsparepartentities.ItemWarehouseTransferRequestDetail, *exceptions.BaseErrorResponse)
	UpdateWhTransferRequest(transactionsparepartpayloads.UpdateItemWarehouseTransferRequest, int) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse)
	UpdateWhTransferRequestDetail(transactionsparepartpayloads.UpdateItemWarehouseTransferRequestDetailRequest, int) (transactionsparepartentities.ItemWarehouseTransferRequestDetail, *exceptions.BaseErrorResponse)
	SubmitWhTransferRequest(int) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse)
	GetAllDetailTransferRequest(int, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdTransferRequest(int) (transactionsparepartpayloads.GetByIdItemWarehouseTransferRequestResponse, *exceptions.BaseErrorResponse)
	GetByIdTransferRequestDetail(int) (transactionsparepartpayloads.GetByIdItemWarehouseTransferRequestDetailResponse, *exceptions.BaseErrorResponse)
	GetAllWhTransferRequest(pagination.Pagination, []utils.FilterCondition, map[string]string) (pagination.Pagination, *exceptions.BaseErrorResponse)
	DeleteHeaderTransferRequest(int) (bool, *exceptions.BaseErrorResponse)
	DeleteDetail(int) (bool, *exceptions.BaseErrorResponse)
	PreviewUploadData([][]string) ([]transactionsparepartpayloads.UploadPreviewItemWarehouseTransferRequestPayloads, *exceptions.BaseErrorResponse)
	ProcessUploadData(transactionsparepartpayloads.UploadProcessItemWarehouseTransferRequestPayloads) ([]transactionsparepartentities.ItemWarehouseTransferRequestDetail, *exceptions.BaseErrorResponse)
	GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse)
}
