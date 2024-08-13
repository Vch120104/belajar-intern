package transactionsparepartservice

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
)

type SupplySlipService interface {
	GetSupplySliptById(Id int, pagination pagination.Pagination) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetSupplySlipDetailById(Id int) (transactionsparepartpayloads.SupplySlipDetailResponse, *exceptions.BaseErrorResponse)
	SaveSupplySlip(req transactionsparepartentities.SupplySlip) (transactionsparepartentities.SupplySlip, *exceptions.BaseErrorResponse)
	SaveSupplySlipDetail(req transactionsparepartentities.SupplySlipDetail) (transactionsparepartentities.SupplySlipDetail, *exceptions.BaseErrorResponse)
	GetAllSupplySlip(internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	UpdateSupplySlip(req transactionsparepartentities.SupplySlip, id int) (transactionsparepartentities.SupplySlip, *exceptions.BaseErrorResponse)
	UpdateSupplySlipDetail(req transactionsparepartentities.SupplySlipDetail, id int)(transactionsparepartentities.SupplySlipDetail,*exceptions.BaseErrorResponse)
	SubmitSupplySlip(id int) (bool, string, *exceptions.BaseErrorResponse)
}
