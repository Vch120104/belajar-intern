package transactionsparepartrepositoryimpl

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	exceptionsss_test "after-sales/api/expectionsss"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"

	"net/http"

	"gorm.io/gorm"
)

type SalesOrderRepositoryImpl struct {
}

func StartSalesOrderRepositoryImpl() transactionsparepartrepository.SalesOrderRepository {
	return &SalesOrderRepositoryImpl{}
}

func (r *SalesOrderRepositoryImpl) GetSalesOrderByID(tx *gorm.DB, Id int) (transactionsparepartpayloads.SalesOrderResponse, *exceptionsss_test.BaseErrorResponse) {
	entities := transactionsparepartentities.SalesOrder{}
	response := transactionsparepartpayloads.SalesOrderResponse{}

	rows, err := tx.Model(&entities).
		Where("sales_order_system_number = ?", Id).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}
