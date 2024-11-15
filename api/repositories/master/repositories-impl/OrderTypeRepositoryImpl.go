package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"
	"errors"
	"net/http"
	"strings"

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

func (r *OrderTypeRepositoryImpl) GetOrderTypeById(tx *gorm.DB, id int) (masterpayloads.GetOrderTypeResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.OrderType{}
	response := masterpayloads.GetOrderTypeResponse{}

	err := tx.Model(&entities).Where(masterentities.OrderType{OrderTypeId: id}).First(&response).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "order type not found",
				Err:        err,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching order type data",
			Err:        err,
		}
	}

	return response, nil
}

func (r *OrderTypeRepositoryImpl) GetOrderTypeByName(tx *gorm.DB, name string) ([]masterpayloads.GetOrderTypeResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.OrderType{}
	response := []masterpayloads.GetOrderTypeResponse{}

	err := tx.Model(&entities).Where("order_type_name LIKE ?", "%"+name+"%").Scan(&response).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching order type data",
			Err:        err,
		}
	}

	if len(response) == 0 {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "order type not found",
			Err:        errors.New("order type not found"),
		}
	}

	return response, nil
}

func (r *OrderTypeRepositoryImpl) SaveOrderType(tx *gorm.DB, req masterpayloads.OrderTypeSaveRequest) (masterentities.OrderType, *exceptions.BaseErrorResponse) {
	entities := masterentities.OrderType{
		IsActive:      true,
		OrderTypeCode: req.OrderTypeCode,
		OrderTypeName: req.OrderTypeName,
	}

	err := tx.Save(&entities).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Message:    "order type code already exist",
				Err:        err,
			}
		}
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error saving order type",
			Err:        err,
		}
	}

	return entities, nil
}

func (r *OrderTypeRepositoryImpl) UpdateOrderType(tx *gorm.DB, id int, req masterpayloads.OrderTypeUpdateRequest) (masterentities.OrderType, *exceptions.BaseErrorResponse) {
	entities := masterentities.OrderType{}

	err := tx.Model(&entities).Where(masterentities.OrderType{OrderTypeId: id}).First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "order type not found",
				Err:        err,
			}
		}
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching order type data",
			Err:        err,
		}
	}

	entities.OrderTypeCode = req.OrderTypeCode
	entities.OrderTypeName = req.OrderTypeName

	err = tx.Save(&entities).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Message:    "order type code already exist",
				Err:        err,
			}
		}
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error updating order type",
			Err:        err,
		}
	}

	return entities, nil
}

func (r *OrderTypeRepositoryImpl) ChangeStatusOrderType(tx *gorm.DB, id int) (masterentities.OrderType, *exceptions.BaseErrorResponse) {
	entities := masterentities.OrderType{}

	err := tx.Model(&entities).Where(masterentities.OrderType{OrderTypeId: id}).First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "order type not found",
				Err:        err,
			}
		}
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching order type data",
			Err:        err,
		}
	}

	entities.IsActive = !entities.IsActive

	err = tx.Save(&entities).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error changing status order type",
			Err:        err,
		}
	}

	return entities, nil
}

func (r *OrderTypeRepositoryImpl) DeleteOrderType(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	entities := masterentities.OrderType{}

	err := tx.Model(&entities).Where(masterentities.OrderType{OrderTypeId: id}).First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "order type not found",
				Err:        err,
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching order type data",
			Err:        err,
		}
	}

	err = tx.Delete(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error deleting order type",
			Err:        err,
		}
	}

	return true, nil
}
