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
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	New(request transactionworkshoppayloads.BookingEstimationNormalRequest) (transactionworkshopentities.BookingEstimation, *exceptions.BaseErrorResponse)
	GetById(id int) (transactionworkshoppayloads.BookingEstimationResponseDetail, *exceptions.BaseErrorResponse)
	Save(request transactionworkshoppayloads.BookingEstimationSaveRequest, id int) (transactionworkshopentities.BookingEstimation, *exceptions.BaseErrorResponse)
	Submit(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)
	Void(Id int) (bool, *exceptions.BaseErrorResponse)
	CloseOrder(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse
	DeleteBookEstimReq(booksysno int, id []int) (bool, *exceptions.BaseErrorResponse)
	SaveBookEstimReq(req transactionworkshoppayloads.BookEstimRemarkRequest, id int) (transactionworkshopentities.BookingEstimationRequest, *exceptions.BaseErrorResponse)
	UpdateBookEstimReq(req transactionworkshoppayloads.BookEstimRemarkRequest, booksysno int, id int) (transactionworkshopentities.BookingEstimationRequest, *exceptions.BaseErrorResponse)
	GetByIdBookEstimReq(booksysno int, id int) (transactionworkshoppayloads.BookEstimRemarkResponse, *exceptions.BaseErrorResponse)
	GetAllBookEstimReq(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)

	SaveBookEstimReminderServ(req transactionworkshoppayloads.ReminderServicePost, id int) (transactionworkshopentities.BookingEstimationServiceReminder, *exceptions.BaseErrorResponse)
	AddPackage(id int, packId int) (bool, *exceptions.BaseErrorResponse)
	AddContractService(id int, contractserviceid int) (bool, *exceptions.BaseErrorResponse)
	InputDiscount(id int, req transactionworkshoppayloads.BookEstimationPayloadsDiscount) (int, *exceptions.BaseErrorResponse)
	AddFieldAction(id int, idrecall int) (int, *exceptions.BaseErrorResponse)
	CopyFromHistory(batchid int) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	PostBookingEstimationCalculation(id int) (int, *exceptions.BaseErrorResponse)
	SaveBookingEstimationFromPDI(id int, req transactionworkshoppayloads.PdiServiceRequest) (bool, *exceptions.BaseErrorResponse)
	SaveBookingEstimationFromServiceRequest(id int, req transactionworkshoppayloads.PdiServiceRequest) (bool, *exceptions.BaseErrorResponse)
	SaveBookingEstimationAllocation(id int, req transactionworkshoppayloads.BookEstimationAllocation) (transactionworkshopentities.BookingEstimationAllocation, *exceptions.BaseErrorResponse)

	GetAllDetailBookingEstimation(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveDetailBookEstim(id int, req transactionworkshoppayloads.BookingEstimationDetailRequest) (transactionworkshopentities.BookingEstimationDetail, *exceptions.BaseErrorResponse)
	GetByIdBookEstimDetail(estimsysno int, id int) (transactionworkshoppayloads.BookingEstimationDetailResponse, *exceptions.BaseErrorResponse)
	UpdateDetailBookEstim(estimsysno int, id int, req transactionworkshoppayloads.BookingEstimationDetailRequestSave) (transactionworkshopentities.BookingEstimationDetail, *exceptions.BaseErrorResponse)
	DeleteDetailBookEstim(estimsysno int, id []int) (bool, *exceptions.BaseErrorResponse)
}
