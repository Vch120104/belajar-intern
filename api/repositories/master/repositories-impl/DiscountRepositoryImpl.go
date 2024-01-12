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
	myDB *gorm.DB
}

func StartDiscountRepositoryImpl(db *gorm.DB) masterrepository.DiscountRepository {
	return &DiscountRepositoryImpl{myDB: db}
}

func (r *DiscountRepositoryImpl) WithTrx(trxHandle *gorm.DB) masterrepository.DiscountRepository {
	if trxHandle == nil {
		log.Println("Transaction Database Not Found!")
		return r
	}
	r.myDB = trxHandle
	return r
}

func (r *DiscountRepositoryImpl) GetAllDiscount(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
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

	r.myDB.Logger = newLogger
	//define base model
	baseModelQuery := r.myDB.Model(&entities)
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	//apply pagination and execute
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&responses).Rows()

	if err != nil {
		return pages, err
	}

	defer rows.Close()

	pages.Rows = responses

	return pages, nil
}

func (r *DiscountRepositoryImpl) GetAllDiscountIsActive() ([]masterpayloads.DiscountResponse, error) {
	var Discounts []masterentities.Discount
	response := []masterpayloads.DiscountResponse{}

	err := r.myDB.Model(&Discounts).Where("is_active = 'true'").Scan(&response).Error

	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *DiscountRepositoryImpl) GetDiscountById(Id int) (masterpayloads.DiscountResponse, error) {
	entities := masterentities.Discount{}
	response := masterpayloads.DiscountResponse{}

	rows, err := r.myDB.Model(&entities).
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

func (r *DiscountRepositoryImpl) GetDiscountByCode(Code string) (masterpayloads.DiscountResponse, error) {
	entities := masterentities.Discount{}
	response := masterpayloads.DiscountResponse{}

	rows, err := r.myDB.Model(&entities).
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

func (r *DiscountRepositoryImpl) SaveDiscount(req masterpayloads.DiscountResponse) (bool, error) {
	entities := masterentities.Discount{
		IsActive:                req.IsActive,
		DiscountCodeId:          req.DiscountCodeId,
		DiscountCodeValue:       req.DiscountCodeValue,
		DiscountCodeDescription: req.DiscountCodeDescription,
	}

	err := r.myDB.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *DiscountRepositoryImpl) ChangeStatusDiscount(Id int) (bool, error) {
	var entities masterentities.Discount

	result := r.myDB.Model(&entities).
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

	result = r.myDB.Save(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}
