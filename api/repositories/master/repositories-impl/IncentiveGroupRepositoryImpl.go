package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"

	"log"

	"gorm.io/gorm"
)

type IncentiveGroupImpl struct {
	myDB *gorm.DB
}

func StartIncentiveGroupImpl(db *gorm.DB) masterrepository.IncentiveGroupRepository {
	return &IncentiveGroupImpl{myDB: db}
}

func (r *IncentiveGroupImpl) WithTrx(trxHandle *gorm.DB) masterrepository.IncentiveGroupRepository {
	if trxHandle == nil {
		log.Println("Transaction Database Not Found!")
		return r
	}
	r.myDB = trxHandle
	return r
}

func (r *IncentiveGroupImpl) GetAllIncentiveGroupIsActive() ([]masterpayloads.IncentiveGroupResponse, error) {
	var IncentiveGroups []masterentities.IncentiveGroup
	response := []masterpayloads.IncentiveGroupResponse{}

	err := r.myDB.Model(&IncentiveGroups).Where("is_active = 'true'").Scan(&response).Error

	if err != nil {
		return response, err
	}

	return response, nil
}
