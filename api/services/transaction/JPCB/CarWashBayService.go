package transactionjpcbservice

import (
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	"after-sales/api/utils"
)

type BayMasterService interface {
	GetAllCarWashBay(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllActiveCarWashBay(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllDeactiveCarWashBay(filterCondition []utils.FilterCondition) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	GetAllCarWashBayDropDown(filterCondition []utils.FilterCondition) ([]transactionjpcbpayloads.CarWashBayDropDownResponse, *exceptions.BaseErrorResponse)
	ChangeStatusCarWashBay(request transactionjpcbpayloads.CarWashBayUpdateRequest) (transactionjpcbentities.BayMaster, *exceptions.BaseErrorResponse)
	PostCarWashBay(request transactionjpcbpayloads.CarWashBayPostRequest) (transactionjpcbentities.BayMaster, *exceptions.BaseErrorResponse)
	PutCarWashBay(request transactionjpcbpayloads.CarWashBayPutRequest) (transactionjpcbentities.BayMaster, *exceptions.BaseErrorResponse)
	GetCarWashBayById(carWashBayId int) (transactionjpcbentities.BayMaster, *exceptions.BaseErrorResponse)
}
