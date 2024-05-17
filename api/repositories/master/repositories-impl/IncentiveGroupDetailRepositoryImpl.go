package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IncentiveGroupDetailRepositoryImpl struct {
}

func StartIncentiveGroupDetailRepositoryImpl() masterrepository.IncentiveGroupDetailRepository {
	return &IncentiveGroupDetailRepositoryImpl{}
}

func (r *IncentiveGroupDetailRepositoryImpl) GetAllIncentiveGroupDetail(tx *gorm.DB, headerId int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	entities := []masterentities.IncentiveGroupDetail{}
	response := []masterpayloads.IncentiveGroupDetailResponse{}
	//define base model
	query := tx.
		Model(&entities).
		Where(masterentities.IncentiveGroupDetail{IncentiveGroupId: headerId}).
		Scan(&response)

	//apply pagination and execute
	rows, err := query.Scopes(pagination.Paginate(&entities, &pages, query)).Scan(&response).Rows()

	if err != nil {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	pages.Rows = response

	if len(response) == 0 {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	return pages, nil
}

func (r *IncentiveGroupDetailRepositoryImpl) GetIncentiveGroupDetailById(tx *gorm.DB, Id int) (masterpayloads.IncentiveGroupDetailResponse, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.IncentiveGroupDetail{}
	response := masterpayloads.IncentiveGroupDetailResponse{}

	rows, err := tx.Model(&entities).
		Where(masterentities.IncentiveGroupDetail{
			IncentiveGroupDetailId: Id,
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

func (r *IncentiveGroupDetailRepositoryImpl) SaveIncentiveGroupDetail(tx *gorm.DB, req masterpayloads.IncentiveGroupDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.IncentiveGroupDetail{
		IsActive:               req.IsActive,
		IncentiveGroupDetailId: req.IncentiveGroupDetailId,
		IncentiveGroupId:       req.IncentiveGroupId,
		IncentiveLevel:         req.IncentiveLevel,
		TargetAmount:           req.TargetAmount,
		TargetPercent:          req.TargetPercent,
	}

	err := tx.Create(&entities).Error

	if err != nil {
		logrus.Info(err)
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}

	return true, nil
}

func (r *IncentiveGroupDetailRepositoryImpl) UpdateIncentiveGroupDetail(tx *gorm.DB, req masterpayloads.UpdateIncentiveGroupDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse) {

	model := masterentities.IncentiveGroupDetail{}
	if err := tx.Model(&model).Where(masterentities.IncentiveGroupDetail{IncentiveGroupDetailId: req.IncentiveGroupDetailId}).First(&model).Error; err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if model == (masterentities.IncentiveGroupDetail{}) {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	entities := masterentities.IncentiveGroupDetail{
		IncentiveGroupDetailId: req.IncentiveGroupDetailId,
		TargetAmount:           req.TargetAmount,
		TargetPercent:          req.TargetPercent,
	}

	err := tx.Updates(&entities).Where(masterentities.IncentiveGroupDetail{IncentiveGroupDetailId: req.IncentiveGroupDetailId}).Error

	if err != nil {

		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}
