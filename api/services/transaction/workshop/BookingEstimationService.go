package transactionworkshopservice

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type BookingEstimationService interface {
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	New(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (bool, *exceptions.BaseErrorResponse)
	GetById(id int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	Save(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (bool, *exceptions.BaseErrorResponse)
	Submit(tx *gorm.DB, Id int) (bool,*exceptions.BaseErrorResponse)
	Void(tx *gorm.DB, Id int) (bool,*exceptions.BaseErrorResponse)
	CloseOrder(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse
	SaveBookEstimReq(req transactionworkshoppayloads.BookEstimRemarkRequest, id int) (int, *exceptions.BaseErrorResponse)
	UpdateBookEstimReq(req transactionworkshoppayloads.BookEstimRemarkRequest, id int) (int, *exceptions.BaseErrorResponse)
	GetByIdBookEstimReq( id int) (transactionworkshoppayloads.BookEstimRemarkRequest, *exceptions.BaseErrorResponse)
	GetAllBookEstimReq(pages *pagination.Pagination, id int) ([]transactionworkshoppayloads.BookEstimRemarkRequest, *exceptions.BaseErrorResponse)
	SaveBookEstimReminderServ(req transactionworkshoppayloads.ReminderServicePost, id int) (int, *exceptions.BaseErrorResponse)
	SaveDetailBookEstim(req transactionworkshoppayloads.BookEstimDetailReq) (int, *exceptions.BaseErrorResponse)
	AddPackage(id int, packId int) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	AddContractService(id int, contractserviceid int) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	InputDiscount(id int, req transactionworkshoppayloads.BookEstimationPayloadsDiscount) (int, *exceptions.BaseErrorResponse)
	AddFieldAction(id int, idrecall int) (int, *exceptions.BaseErrorResponse)
	CopyFromHistory(id int)([]map[string]interface{},*exceptions.BaseErrorResponse)
	GetByIdBookEstimDetail (id int ,LineTypeID int)(map[string]interface{},*exceptions.BaseErrorResponse)
	PostBookingEstimationCalculation(id int)(int,*exceptions.BaseErrorResponse)
	SaveBookingEstimationFromPDI( id int) (transactionworkshopentities.BookingEstimation, *exceptions.BaseErrorResponse)
	SaveBookingEstimationFromServiceRequest(id int)(transactionworkshopentities.BookingEstimation,*exceptions.BaseErrorResponse)
}
