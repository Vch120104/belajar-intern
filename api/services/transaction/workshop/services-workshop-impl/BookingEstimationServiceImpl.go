package transactionworkshopserviceimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"errors"
	"net/http"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BookingEstimationServiceImpl struct {
	structBookingEstimationRepo transactionworkshoprepository.BookingEstimationRepository
	DB                          *gorm.DB
	RedisClient                 *redis.Client // Redis client
}

func OpenBookingEstimationServiceImpl(bookingEstimationRepo transactionworkshoprepository.BookingEstimationRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.BookingEstimationService {
	return &BookingEstimationServiceImpl{
		structBookingEstimationRepo: bookingEstimationRepo,
		DB:                          db,
		RedisClient:                 redisClient,
	}
}

func (s *BookingEstimationServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.structBookingEstimationRepo.GetAll(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *BookingEstimationServiceImpl) New(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (bool, *exceptions.BaseErrorResponse) {
	
	_, err := s.structBookingEstimationRepo.Save(tx, request)
	defer helper.CommitOrRollback(tx,err)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *BookingEstimationServiceImpl) GetById(id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.structBookingEstimationRepo.GetById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *BookingEstimationServiceImpl) Save(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (bool, *exceptions.BaseErrorResponse) {
	// Menggunakan "=" untuk menginisialisasi tx dengan transaksi yang dimulai
	_, err := s.structBookingEstimationRepo.Save(tx, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *BookingEstimationServiceImpl) Submit(tx *gorm.DB, id int) (bool,*exceptions.BaseErrorResponse) {
	
	result,err := s.structBookingEstimationRepo.Submit(tx, id)
	defer helper.CommitOrRollback(tx,err)
	if err != nil {
		return false,err
	}
	return result,nil
}

func (s *BookingEstimationServiceImpl) Void(tx *gorm.DB, id int) (bool,*exceptions.BaseErrorResponse) {
	result,err := s.structBookingEstimationRepo.Void(tx, id)
	defer helper.CommitOrRollback(tx,err)
	if err != nil {
		return false,err
	}
	return result,nil
}

func (s *BookingEstimationServiceImpl) CloseOrder(tx *gorm.DB, id int) *exceptions.BaseErrorResponse {
	err := s.structBookingEstimationRepo.CloseOrder(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return err
	}
	return nil
}

func (s *BookingEstimationServiceImpl) SaveBookEstimReq(req transactionworkshoppayloads.BookEstimRemarkRequest, id int) (int, *exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	
	result,err := s.structBookingEstimationRepo.SaveBookEstimReq(tx,req,id)
	defer helper.CommitOrRollback(tx,err)
	if err != nil{
		return 0,err
	}
	return result,nil
}

func(s *BookingEstimationServiceImpl) UpdateBookEstimReq(req transactionworkshoppayloads.BookEstimRemarkRequest, id int) (int, *exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	
	result,err := s.structBookingEstimationRepo.UpdateBookEstimReq(tx,req,id)
	defer helper.CommitOrRollback(tx,err)
	if err != nil{
		return 0,err
	}
	return result,nil
}

func (s *BookingEstimationServiceImpl) GetByIdBookEstimReq( id int) (transactionworkshoppayloads.BookEstimRemarkRequest, *exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	
	result,err := s.structBookingEstimationRepo.GetByIdBookEstimReq(tx,id)
	defer helper.CommitOrRollback(tx,err)
	if err != nil{
		return transactionworkshoppayloads.BookEstimRemarkRequest{},err
	}
	return result,nil
}

func (s *BookingEstimationServiceImpl) GetAllBookEstimReq(pages *pagination.Pagination, id int) ([]transactionworkshoppayloads.BookEstimRemarkRequest, *exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	
	result,err := s.structBookingEstimationRepo.GetAllBookEstimReq(tx,pages,id)
	defer helper.CommitOrRollback(tx,err)
	if err != nil{
		return []transactionworkshoppayloads.BookEstimRemarkRequest{},err
	}
	return result,nil
}

func (s *BookingEstimationServiceImpl) SaveBookEstimReminderServ(req transactionworkshoppayloads.ReminderServicePost, id int) (int, *exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	result,err := s.structBookingEstimationRepo.SaveBookEstimReminderServ(tx,req,id)
	defer helper.CommitOrRollback(tx,err)
	if err != nil{
		return 0,err
	}
	return result,nil
}

func (s *BookingEstimationServiceImpl) SaveDetailBookEstim(req transactionworkshoppayloads.BookEstimDetailReq) (int, *exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	
	result,err := s.structBookingEstimationRepo.SaveDetailBookEstim(tx,req)
	defer helper.CommitOrRollback(tx,err)
	if err != nil{
		return 0,err
	}
	_,err2:= s.structBookingEstimationRepo.PutBookingEstimationCalculation(tx,req.EstimationSystemNumber,req.LineTypeID)
	if err2 !=nil{
		return 0,err2
	}
	return result,nil
}

func (s *BookingEstimationServiceImpl) AddPackage(id int, packId int) ([]map[string]interface{}, *exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	result,err := s.structBookingEstimationRepo.AddPackage(tx,id,packId)
	defer helper.CommitOrRollback(tx,err)
	if err !=nil{
		return nil,err
	}
	for _,detail := range result{
		estimationSystemNumberFloat, ok1 := detail["estimation_system_number"].(float64)
		lineTypeIDFloat, ok2 := detail["line_type_id"].(float64)

		if !ok1 || !ok2 {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("invalid type assertion for map values"),
			}
    }

    // Convert float64 to int
		estimationSystemNumber := int(estimationSystemNumberFloat)
		lineTypeID := int(lineTypeIDFloat)
		_,err:= s.structBookingEstimationRepo.PutBookingEstimationCalculation(tx,estimationSystemNumber,lineTypeID)
		if err != nil{
			return nil,err
		}
	}
	return result,nil
}

func (s *BookingEstimationServiceImpl) AddContractService(id int, contractserviceid int) ([]map[string]interface{}, *exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	result,err := s.structBookingEstimationRepo.AddContractService(tx,id,contractserviceid)
	defer helper.CommitOrRollback(tx,err)
	if err !=nil{
		return nil,err
	}
	for _,detail := range result{
		estimationSystemNumberFloat, ok1 := detail["estimation_system_number"].(float64)
		lineTypeIDFloat, ok2 := detail["line_type_id"].(float64)

		if !ok1 || !ok2 {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("invalid type assertion for map values"),
			}
    }

    // Convert float64 to int
		estimationSystemNumber := int(estimationSystemNumberFloat)
		lineTypeID := int(lineTypeIDFloat)
		_,err:= s.structBookingEstimationRepo.PutBookingEstimationCalculation(tx,estimationSystemNumber,lineTypeID)
		if err != nil{
			return nil,err
		}
	}
	return result,nil
}

func (s *BookingEstimationServiceImpl) CopyFromHistory(id int)([]map[string]interface{},*exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	result,err := s.structBookingEstimationRepo.CopyFromHistory(tx,id)
	if err != nil{
		return nil,err
	}
	for _,detail := range result{
		estimationSystemNumberFloat, ok1 := detail["estimation_system_number"].(float64)
		lineTypeIDFloat, ok2 := detail["line_type_id"].(float64)

		if !ok1 || !ok2 {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("invalid type assertion for map values"),
			}
    }

    // Convert float64 to int
		estimationSystemNumber := int(estimationSystemNumberFloat)
		lineTypeID := int(lineTypeIDFloat)
		_,err:= s.structBookingEstimationRepo.PutBookingEstimationCalculation(tx,estimationSystemNumber,lineTypeID)
		if err != nil{
			return nil,err
		}
	}
	return result,nil
}

func (s *BookingEstimationServiceImpl) InputDiscount(id int, req transactionworkshoppayloads.BookEstimationPayloadsDiscount) (int, *exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	result,err := s.structBookingEstimationRepo.InputDiscount(tx,id,req)
	defer helper.CommitOrRollback(tx,err)
	if err != nil{
		return 0,err
	}
	return result,nil
}

func (s *BookingEstimationServiceImpl) AddFieldAction(id int, idrecall int) (int, *exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	result,err:= s.structBookingEstimationRepo.AddFieldAction(tx,id,idrecall)
	defer helper.CommitOrRollback(tx,err)
	if err !=nil{
		return 0,err
	}
	return result,nil
}

func (s *BookingEstimationServiceImpl) GetByIdBookEstimDetail (id int ,LineTypeID int)(map[string]interface{},*exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	result,err := s.structBookingEstimationRepo.GetByIdBookEstimDetail(tx,id,LineTypeID)
	defer helper.CommitOrRollback(tx,err)
	if err != nil{
		return nil,err
	}
	return result,nil
}

func (s *BookingEstimationServiceImpl) PostBookingEstimationCalculation(id int)(int,*exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	result,err := s.structBookingEstimationRepo.PostBookingEstimationCalculation(tx,id)
	defer helper.CommitOrRollback(tx,err)
	if err != nil{
		return 0,err
	}
	return result,nil
}



func (s *BookingEstimationServiceImpl) SaveBookingEstimationFromPDI( id int) (transactionworkshopentities.BookingEstimation, *exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	result,err := s.structBookingEstimationRepo.SaveBookingEstimationFromPDI(tx,id)
	defer helper.CommitOrRollback(tx,err)
	if err != nil{
		return transactionworkshopentities.BookingEstimation{},err
	}
	return result,nil
}

func (s *BookingEstimationServiceImpl) SaveBookingEstimationFromServiceRequest(id int)(transactionworkshopentities.BookingEstimation,*exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	result,err := s.structBookingEstimationRepo.SaveBookingEstimationFromServiceRequest(tx,id)
	defer helper.CommitOrRollback(tx,err)
	if err != nil{
		return transactionworkshopentities.BookingEstimation{},err
	}
	return result,nil
}