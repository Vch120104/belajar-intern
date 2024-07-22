package transactionworkshoprepository

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type BookingEstimationRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	Post(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (transactionworkshopentities.BookingEstimation, *exceptions.BaseErrorResponse)
	GetById(tx *gorm.DB, Id int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	Save(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (bool, *exceptions.BaseErrorResponse)
	Submit(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse
	Void(tx *gorm.DB, Id int) (bool,*exceptions.BaseErrorResponse)
	CloseOrder(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse
	SaveBookEstimReq(tx *gorm.DB, req transactionworkshoppayloads.BookEstimRemarkRequest, id int) (int, *exceptions.BaseErrorResponse)
	UpdateBookEstimReq(tx *gorm.DB, req transactionworkshoppayloads.BookEstimRemarkRequest, id int) (int, *exceptions.BaseErrorResponse)
	GetByIdBookEstimReq(tx *gorm.DB, id int) (transactionworkshoppayloads.BookEstimRemarkRequest, *exceptions.BaseErrorResponse)
	GetAllBookEstimReq(tx *gorm.DB, pages *pagination.Pagination, id int) ([]transactionworkshoppayloads.BookEstimRemarkRequest, *exceptions.BaseErrorResponse)
	SaveBookEstimReminderServ(tx *gorm.DB, req transactionworkshoppayloads.ReminderServicePost, id int) (int, *exceptions.BaseErrorResponse)
	SaveDetailBookEstim(tx *gorm.DB, req transactionworkshoppayloads.BookEstimDetailReq) (int, *exceptions.BaseErrorResponse)
	AddPackage(tx *gorm.DB, id int, packId int) (int, *exceptions.BaseErrorResponse)
	AddContractService(tx *gorm.DB, id int, contractserviceid int) (int, *exceptions.BaseErrorResponse)
	InputDiscount(tx *gorm.DB, id int, req transactionworkshoppayloads.BookEstimationPayloadsDiscount) (int, *exceptions.BaseErrorResponse)
	AddFieldAction(tx *gorm.DB, id int, idrecall int) (int, *exceptions.BaseErrorResponse)
	GetByIdBookEstimDetail (tx *gorm.DB ,id int ,LineTypeID int)(map[string]interface{},*exceptions.BaseErrorResponse)
	PostBookingEstimationCalculation(tx*gorm.DB,id int)(int,*exceptions.BaseErrorResponse)
	PutBookingEstimationCalculation (tx *gorm.DB, id int, linetypeid int, req transactionworkshoppayloads.BookingEstimationCalculationPayloads)(int,*exceptions.BaseErrorResponse)
}
