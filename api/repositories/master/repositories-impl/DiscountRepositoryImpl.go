package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"net/http"

	exceptionsss_test "after-sales/api/expectionsss"

	"after-sales/api/utils"

	"gorm.io/gorm"
)

type DiscountRepositoryImpl struct {
}

func StartDiscountRepositoryImpl() masterrepository.DiscountRepository {
	return &DiscountRepositoryImpl{}
}

func (r *DiscountRepositoryImpl) GetAllDiscount(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.Discount{}
	responses := []masterpayloads.DiscountResponse{}

	//define base model
	baseModelQuery := tx.Model(&entities).Scan(&responses)
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	//apply pagination and execute
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&responses).Rows()

	if err != nil {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	if len(responses) == 0 {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	pages.Rows = responses

	return pages, nil
}

func (r *DiscountRepositoryImpl) GetAllDiscountIsActive(tx *gorm.DB) ([]masterpayloads.DiscountResponse, *exceptionsss_test.BaseErrorResponse) {
	var Discounts []masterentities.Discount
	response := []masterpayloads.DiscountResponse{}

	rows, err := tx.
		Model(&Discounts).
		Where(masterentities.Discount{IsActive: true}).
		Scan(&response).
		Rows()

	if len(response) == 0 {
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if err != nil {

		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *DiscountRepositoryImpl) GetDiscountById(tx *gorm.DB, Id int) (masterpayloads.DiscountResponse, *exceptionsss_test.BaseErrorResponse) {
	var entities masterentities.Discount
	var response masterpayloads.DiscountResponse

	rows, err := tx.Model(&entities).
		Where(masterentities.Discount{
			DiscountCodeId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *DiscountRepositoryImpl) GetDiscountByCode(tx *gorm.DB, Code string) (masterpayloads.DiscountResponse, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.Discount{}
	response := masterpayloads.DiscountResponse{}

	rows, err := tx.Model(&entities).
		Where(masterentities.Discount{
			DiscountCodeValue: Code,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *DiscountRepositoryImpl) SaveDiscount(tx *gorm.DB, req masterpayloads.DiscountResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.Discount{
		IsActive:                req.IsActive,
		DiscountCodeId:          req.DiscountCodeId,
		DiscountCodeValue:       req.DiscountCodeValue,
		DiscountCodeDescription: req.DiscountCodeDescription,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *DiscountRepositoryImpl) ChangeStatusDiscount(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities masterentities.Discount

	result := tx.Model(&entities).
		Where(masterentities.Discount{DiscountCodeId: Id}).
		First(&entities)

	if result.Error != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
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
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}
