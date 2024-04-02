package masterserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type DeductionServiceImpl struct {
	deductionrepo masterrepository.DeductionRepository
	DB            *gorm.DB
}

func StartDeductionService(deductionRepo masterrepository.DeductionRepository, db *gorm.DB) masterservice.DeductionService {
	return &DeductionServiceImpl{
		deductionrepo: deductionRepo,
		DB:            db,
	}
}

func (s *DeductionServiceImpl) GetAllDeduction(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.deductionrepo.GetAllDeduction(tx, filterCondition, pages)

	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}

func (s *DeductionServiceImpl) GetByIdDeductionDetail(Id int) masterpayloads.DeductionDetailResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.deductionrepo.GetByIdDeductionDetail(tx, Id)

	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}

func (s *DeductionServiceImpl) PostDeductionList(req masterpayloads.DeductionListResponse) masterpayloads.DeductionListResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.deductionrepo.SaveDeductionList(tx, req)
	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}

func (s *DeductionServiceImpl) PostDeductionDetail(req masterpayloads.DeductionDetailResponse) masterpayloads.DeductionDetailResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.deductionrepo.SaveDeductionDetail(tx, req)
	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}

func (s *DeductionServiceImpl) GetByIdDeductionList(Id int, page int, limit int) payloads.ResponsePaginationHeader {

	tx := s.DB.Begin()

	defer helper.CommitOrRollback(tx)

	result, err := s.deductionrepo.GetDeductionById(tx, Id)

	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}

	if limit == 0 && page == 0 {
		limit = 10
	}

	pagination := pagination.Pagination{
		Limit: limit,
		Page:  page,
	}

	detail_result, detail_err := s.deductionrepo.GetAllDeductionDetail(tx, pagination, Id)
	
	if detail_err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}

	detail_response := payloads.ResponsePagination{
		StatusCode: 200,
		Message:    "success",
		Page:       detail_result.Page,
		Limit:      detail_result.Limit,
		TotalRows:  detail_result.TotalRows,
		TotalPages: detail_result.TotalPages,
		Data:       detail_result.Rows,
	}

	return payloads.ResponsePaginationHeader{
		Header: result,
		Data:   detail_response,
	}

}

func (s *DeductionServiceImpl) ChangeStatusDeduction(Id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.deductionrepo.GetDeductionById(tx, Id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	results, err := s.deductionrepo.ChangeStatusDeduction(tx, Id)
	if err != nil {
		return results
	}
	return true
}