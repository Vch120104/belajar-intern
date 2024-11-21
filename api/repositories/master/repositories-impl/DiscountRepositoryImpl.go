package masterrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"errors"
	"net/http"

	exceptions "after-sales/api/exceptions"

	"after-sales/api/utils"

	"gorm.io/gorm"
)

type DiscountRepositoryImpl struct {
}

func StartDiscountRepositoryImpl() masterrepository.DiscountRepository {
	return &DiscountRepositoryImpl{}
}

func (r *DiscountRepositoryImpl) GetAllDiscount(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Discount{}
	responses := []masterpayloads.DiscountResponse{}

	//define base model
	baseModelQuery := tx.Model(&entities).Select("mtr_discount.*, discount_code + ' - ' + discount_description AS discount_code_description")
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	//apply pagination and execute
	err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&responses).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.Rows = responses

	return pages, nil
}

func (r *DiscountRepositoryImpl) GetAllDiscountIsActive(tx *gorm.DB) ([]masterpayloads.DiscountResponse, *exceptions.BaseErrorResponse) {
	var Discounts []masteritementities.Discount
	response := []masterpayloads.DiscountResponse{}

	err := tx.
		Model(&Discounts).
		Select("is_active, discount_code_id, discount_code, discount_description, CONCAT(discount_code, ' - ', discount_description) as discount_code_description").
		Where("is_active = ?", true).
		Scan(&response).Error

	if len(response) == 0 {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if err != nil {

		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

func (r *DiscountRepositoryImpl) GetDiscountById(tx *gorm.DB, Id int) (masterpayloads.DiscountResponse, *exceptions.BaseErrorResponse) {
	var entities masteritementities.Discount
	var response masterpayloads.DiscountResponse

	rows, err := tx.Model(&entities).
		Select("mtr_discount.*, discount_code + ' - ' + discount_description AS discount_code_description").
		Where(masteritementities.Discount{
			DiscountCodeId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *DiscountRepositoryImpl) GetDiscountByCode(tx *gorm.DB, Code string) (masterpayloads.DiscountResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Discount{}
	response := masterpayloads.DiscountResponse{}

	err := tx.Model(&entities).
		Select("mtr_discount.*, discount_code + ' - ' + discount_description AS discount_code_description").
		Where(masteritementities.Discount{
			DiscountCode: Code,
		}).
		First(&response).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "discount code not found",
				Err:        err,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

func (r *DiscountRepositoryImpl) SaveDiscount(tx *gorm.DB, req masterpayloads.DiscountResponse) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Discount{
		IsActive:            req.IsActive,
		DiscountCodeId:      req.DiscountCodeId,
		DiscountCode:        req.DiscountCode,
		DiscountDescription: req.DiscountDescription,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *DiscountRepositoryImpl) ChangeStatusDiscount(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.Discount

	result := tx.Model(&entities).
		Where(masteritementities.Discount{DiscountCodeId: Id}).
		First(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	if entities.IsActive {
		entities.IsActive = false
	} else {
		entities.IsActive = true
	}

	result = tx.Save(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}
