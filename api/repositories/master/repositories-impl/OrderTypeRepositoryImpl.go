package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"
	"net/http"

	"gorm.io/gorm"
)

type OrderTypeRepositoryImpl struct {
}

func StartOrderTypeRepositoryImpl() masterrepository.OrderTypeRepository {
	return &OrderTypeRepositoryImpl{}
}

func (r *OrderTypeRepositoryImpl) GetAllOrderType(tx *gorm.DB) ([]masterpayloads.GetOrderTypeResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.OrderType{}
	response := []masterpayloads.GetOrderTypeResponse{}

	err := tx.Model(&entities).Scan(&response).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching order type data",
			Err:        err,
		}
	}

	return response, nil
}
