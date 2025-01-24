package transactionworkshopservice

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"
)

type AtpmClaimRegistrationService interface {
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetById(id int, pages pagination.Pagination) (transactionworkshoppayloads.AtpmClaimRegistrationResponse, *exceptions.BaseErrorResponse)
	New(request transactionworkshoppayloads.AtpmClaimRegistrationRequest) (transactionworkshopentities.AtpmClaimVehicle, *exceptions.BaseErrorResponse)
	Save(id int, request transactionworkshoppayloads.AtpmClaimRegistrationRequestSave) (transactionworkshopentities.AtpmClaimVehicle, *exceptions.BaseErrorResponse)
	Submit(id int) (bool, *exceptions.BaseErrorResponse)
	Void(id int) (bool, *exceptions.BaseErrorResponse)
	GetAllServiceHistory(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllClaimHistory(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllDetail(id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetDetailById(claimsysno int, id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	AddDetail(id int, request transactionworkshoppayloads.AtpmClaimDetailRequest) (transactionworkshopentities.AtpmClaimVehicleDetail, *exceptions.BaseErrorResponse)
	DeleteDetail(detailId int, claimsysno int) (bool, *exceptions.BaseErrorResponse)
}
