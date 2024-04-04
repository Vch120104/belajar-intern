package masteroperationrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	masteroperationentities "after-sales/api/entities/master/operation"
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	"after-sales/api/utils"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type OperationCodeRepositoryImpl struct {
	myDB *gorm.DB
}

func StartOperationCodeRepositoryImpl(db *gorm.DB) masteroperationrepository.OperationCodeRepository {
	return &OperationCodeRepositoryImpl{myDB: db}
}

func (r *OperationCodeRepositoryImpl) GetAllOperationCode(tx *gorm.DB,filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	entities := []masteroperationentities.OperationCode{}
	//define base model
	baseModelQuery := tx.Model(&entities)
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	//apply pagination and execute
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&entities).Rows()

	if len(entities) == 0 {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	pages.Rows = entities

	return pages, nil
}

func (r *OperationCodeRepositoryImpl) GetOperationCodeById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationCodeResponse, *exceptionsss_test.BaseErrorResponse) {
	entities := masteroperationentities.OperationCode{}
	response := masteroperationpayloads.OperationCodeResponse{}

	rows, err := tx.Model(&entities).
		Where(masteroperationentities.OperationCode{
			OperationId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err: err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *OperationCodeRepositoryImpl) SaveOperationCode(tx *gorm.DB,req masteroperationpayloads.OperationEntriesRequest)(bool,*exceptionsss_test.BaseErrorResponse){
	entities:masteroperationentities.OperationCode{
		
	}
}