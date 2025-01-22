package transactionworkshoprepository

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type AtpmClaimRegistrationRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetById(tx *gorm.DB, id int, pages pagination.Pagination) (transactionworkshoppayloads.AtpmClaimRegistrationResponse, *exceptions.BaseErrorResponse)
	New(tx *gorm.DB, request transactionworkshoppayloads.AtpmClaimRegistrationRequest) (transactionworkshopentities.AtpmClaimVehicle, *exceptions.BaseErrorResponse)
	Save(tx *gorm.DB, id int, request transactionworkshoppayloads.AtpmClaimRegistrationRequestSave) (transactionworkshopentities.AtpmClaimVehicle, *exceptions.BaseErrorResponse)
	Submit(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse)
	Void(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse)
	GetAllServiceHistory(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllClaimHistory(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllDetail(tx *gorm.DB, claimsysno int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetDetailById(tx *gorm.DB, claimsysno int, detailId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	AddDetail(tx *gorm.DB, claimsysno int, request transactionworkshoppayloads.AtpmClaimDetailRequest) (transactionworkshopentities.AtpmClaimVehicleDetail, *exceptions.BaseErrorResponse)
}
