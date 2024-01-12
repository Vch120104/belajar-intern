package masteroperationrepositoryimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	"after-sales/api/utils"
	"log"

	"gorm.io/gorm"
)

type OperationCodeRepositoryImpl struct {
	myDB *gorm.DB
}

func StartOperationCodeRepositoryImpl(db *gorm.DB) masteroperationrepository.OperationCodeRepository {
	return &OperationCodeRepositoryImpl{myDB: db}
}

func (r *OperationCodeRepositoryImpl) GetAllOperationCode(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := []masteroperationentities.OperationCode{}
	//define base model
	baseModelQuery := r.myDB.Model(&entities)
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	//apply pagination and execute
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&entities).Rows()

	if len(entities) == 0 {
		return pages, gorm.ErrRecordNotFound
	}

	if err != nil {
		return pages, err
	}

	defer rows.Close()

	pages.Rows = entities

	return pages, nil
}

func (r *OperationCodeRepositoryImpl) WithTrx(trxHandle *gorm.DB) masteroperationrepository.OperationCodeRepository {
	if trxHandle == nil {
		log.Println("Transaction Database Not Found!")
		return r
	}
	r.myDB = trxHandle
	return r
}

func (r *OperationCodeRepositoryImpl) GetOperationCodeById(Id int32) (masteroperationpayloads.OperationCodeResponse, error) {
	entities := masteroperationentities.OperationCode{}
	response := masteroperationpayloads.OperationCodeResponse{}

	rows, err := r.myDB.Model(&entities).
		Where(masteroperationentities.OperationCode{
			OperationId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}
