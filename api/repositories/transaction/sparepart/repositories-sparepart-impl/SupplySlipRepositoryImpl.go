package transactionsparepartrepositoryimpl

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"

	"gorm.io/gorm"
)

type SupplySlipRepositoryImpl struct {
	myDB *gorm.DB
}

func StartSupplySlipRepositoryImpl(db *gorm.DB) transactionsparepartrepository.SupplySlipRepository {
	return &SupplySlipRepositoryImpl{myDB: db}
}

func (r *SupplySlipRepositoryImpl) GetSupplySlipById(Id int32) (transactionsparepartpayloads.SupplySlipResponse, error) {
	entities := transactionsparepartentities.SupplySlip{}
	response := transactionsparepartpayloads.SupplySlipResponse{}

	rows, err := r.myDB.Model(&entities).
		Where(transactionsparepartentities.SupplySlip{
			SupplySystemNumber: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *SupplySlipRepositoryImpl) GetSupplySlipDetailById(Id int32) (transactionsparepartpayloads.SupplySlipDetailResponse, error) {
	entities := transactionsparepartentities.SupplySlipDetail{}
	response := transactionsparepartpayloads.SupplySlipDetailResponse{}

	rows, err := r.myDB.Model(&entities).
		Where(transactionsparepartentities.SupplySlipDetail{
			SupplySlipDetailSystemNumber: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}
