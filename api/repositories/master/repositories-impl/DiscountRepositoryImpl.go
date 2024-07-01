package masterrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
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
	baseModelQuery := tx.Model(&entities).Scan(&responses)
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	//apply pagination and execute
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&responses).Rows()

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	if len(responses) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	pages.Rows = responses

	return pages, nil
}

func (r *DiscountRepositoryImpl) GetAllDiscountIsActive(tx *gorm.DB) ([]masterpayloads.DiscountResponse, *exceptions.BaseErrorResponse) {
	var Discounts []masteritementities.Discount
	response := []masterpayloads.DiscountResponse{}

	rows, err := tx.
		Model(&Discounts).
		Where(masteritementities.Discount{IsActive: true}).
		Scan(&response).
		Rows()

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

	defer rows.Close()

	return response, nil
}

func (r *DiscountRepositoryImpl) GetDiscountById(tx *gorm.DB, Id int) (masterpayloads.DiscountResponse, *exceptions.BaseErrorResponse) {
	var entities masteritementities.Discount
	var response masterpayloads.DiscountResponse

	rows, err := tx.Model(&entities).
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

	rows, err := tx.Model(&entities).
		Where(masteritementities.Discount{
			DiscountCodeValue: Code,
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

func (r *DiscountRepositoryImpl) SaveDiscount(tx *gorm.DB, req masterpayloads.DiscountResponse) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Discount{
		IsActive:                req.IsActive,
		DiscountCodeId:          req.DiscountCodeId,
		DiscountCodeValue:       req.DiscountCodeValue,
		DiscountCodeDescription: req.DiscountCodeDescription,
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
