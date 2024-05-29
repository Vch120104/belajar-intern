package transactionsparepartrepositoryimpl

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	exceptions "after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"net/http"

	"gorm.io/gorm"
)

type SupplySlipRepositoryImpl struct {
}

func StartSupplySlipRepositoryImpl() transactionsparepartrepository.SupplySlipRepository {
	return &SupplySlipRepositoryImpl{}
}

func (r *SupplySlipRepositoryImpl) GetSupplySlipById(tx *gorm.DB, Id int) (transactionsparepartpayloads.SupplySlipResponse, *exceptions.BaseErrorResponse) {
	entities := transactionsparepartentities.SupplySlip{}
	response := transactionsparepartpayloads.SupplySlipResponse{}

	rows, err := tx.Model(&entities).
		Where("supply_system_number = ?", Id).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *SupplySlipRepositoryImpl) GetSupplySlipDetailById(tx *gorm.DB, Id int) (transactionsparepartpayloads.SupplySlipDetailResponse, *exceptions.BaseErrorResponse) {
	entities := transactionsparepartentities.SupplySlipDetail{}
	response := transactionsparepartpayloads.SupplySlipDetailResponse{}

	rows, err := tx.Model(&entities).
		Where("supply_slip_detail_system_number = ?", Id).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}
