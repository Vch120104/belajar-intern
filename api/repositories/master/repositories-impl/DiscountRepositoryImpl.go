package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"

	"after-sales/api/utils"
	"log"

	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DiscountRepositoryImpl struct {
}

func StartDiscountRepositoryImpl() masterrepository.DiscountRepository {
	return &DiscountRepositoryImpl{}
}

func (r *DiscountRepositoryImpl) GetAllDiscount(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := masterentities.Discount{}
	var responses []masterpayloads.DiscountResponse
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	tx.Logger = newLogger
	//define base model
	baseModelQuery := tx.Model(&entities)
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	//apply pagination and execute
	rows, _ := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&responses).Rows()

	if len(responses) == 0 {
		return pages, gorm.ErrRecordNotFound
	}

	defer rows.Close()

	pages.Rows = responses

	return pages, nil
}

func (r *DiscountRepositoryImpl) GetAllDiscountIsActive(tx *gorm.DB) ([]masterpayloads.DiscountResponse, error) {
	var Discounts []masterentities.Discount
	response := []masterpayloads.DiscountResponse{}

	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	tx.Logger = newLogger

	rows, _ := tx.Model(&Discounts).Where("is_active = 'true'").Scan(&response).Rows()

	if len(response) == 0 {
		return response, gorm.ErrRecordNotFound
	}

	defer rows.Close()

	return response, nil
}

func (r *DiscountRepositoryImpl) GetDiscountById(tx *gorm.DB, Id int) (masterpayloads.DiscountResponse, error) {
	entities := masterentities.Discount{}
	response := masterpayloads.DiscountResponse{}

	rows, err := tx.Model(&entities).
		Where(masterentities.Discount{
			DiscountCodeId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *DiscountRepositoryImpl) GetDiscountByCode(tx *gorm.DB, Code string) (masterpayloads.DiscountResponse, error) {
	entities := masterentities.Discount{}
	response := masterpayloads.DiscountResponse{}

	rows, err := tx.Model(&entities).
		Where(masterentities.Discount{
			DiscountCodeValue: Code,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *DiscountRepositoryImpl) SaveDiscount(tx *gorm.DB, req masterpayloads.DiscountResponse) (bool, error) {
	entities := masterentities.Discount{
		IsActive:                req.IsActive,
		DiscountCodeId:          req.DiscountCodeId,
		DiscountCodeValue:       req.DiscountCodeValue,
		DiscountCodeDescription: req.DiscountCodeDescription,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *DiscountRepositoryImpl) ChangeStatusDiscount(tx *gorm.DB, Id int) (bool, error) {
	var entities masterentities.Discount

	result := tx.Model(&entities).
		Where("discount_code_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	if entities.IsActive {
		entities.IsActive = false
	} else {
		entities.IsActive = true
	}

	result = tx.Save(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}
