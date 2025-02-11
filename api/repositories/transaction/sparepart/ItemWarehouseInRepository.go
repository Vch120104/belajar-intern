package transactionsparepartrepository

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"

	"gorm.io/gorm"
)

type ItemWarehouseTransferInRepository interface {
	InsertDetail(*gorm.DB, transactionsparepartpayloads.InsertItemWarehouseHeaderTransferInRequest) (transactionsparepartentities.ItemWarehouseTransferIn, *exceptions.BaseErrorResponse)
	Submit(*gorm.DB, int) (transactionsparepartentities.ItemWarehouseTransferIn, *exceptions.BaseErrorResponse)
}

// TableNameItemWarehouseTransferInDetail
