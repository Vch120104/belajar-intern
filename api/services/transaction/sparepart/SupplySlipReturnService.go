package transactionsparepartservice

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
)

type SupplySlipReturnService interface {
	SaveSupplySlipReturn(req transactionsparepartentities.SupplySlipReturn) (transactionsparepartentities.SupplySlipReturn, *exceptions.BaseErrorResponse)
	SaveSupplySlipReturnDetail(req transactionsparepartentities.SupplySlipReturnDetail) (transactionsparepartentities.SupplySlipReturnDetail, *exceptions.BaseErrorResponse)
	GetAllSupplySlipReturn(internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetSupplySlipReturnById(Id int, pagination pagination.Pagination) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetSupplySlipReturnDetailById(id int) (transactionsparepartpayloads.SupplySlipReturnDetailResponse, *exceptions.BaseErrorResponse)
	UpdateSupplySlipReturn(req transactionsparepartentities.SupplySlipReturn, id int)(transactionsparepartentities.SupplySlipReturn,*exceptions.BaseErrorResponse)
	UpdateSupplySlipReturnDetail(req transactionsparepartentities.SupplySlipReturnDetail, id int)(transactionsparepartentities.SupplySlipReturnDetail,*exceptions.BaseErrorResponse)
}
