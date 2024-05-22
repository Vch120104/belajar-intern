package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"net/http"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IncentiveGroupDetailRepositoryImpl struct {
}

func StartIncentiveGroupDetailRepositoryImpl() masterrepository.IncentiveGroupDetailRepository {
	return &IncentiveGroupDetailRepositoryImpl{}
}

func (r *IncentiveGroupDetailRepositoryImpl) GetAllIncentiveGroupDetail(tx *gorm.DB, headerId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := []masterentities.IncentiveGroupDetail{}
	response := []masterpayloads.IncentiveGroupDetailResponse{}
	//define base model
	query := tx.
		Model(&entities).
		Where(masterentities.IncentiveGroupDetail{IncentiveGroupId: headerId}).
		Scan(&response)

	//apply pagination and execute
	rows, err := query.Scopes(pagination.Paginate(&entities, &pages, query)).Scan(&response).Rows()

	if len(response) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	pages.Rows = response

	return pages, nil
}

func (r *IncentiveGroupDetailRepositoryImpl) GetIncentiveGroupDetailById(tx *gorm.DB, Id int) (masterpayloads.IncentiveGroupDetailResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.IncentiveGroupDetail{}
	response := masterpayloads.IncentiveGroupDetailResponse{}

	rows, err := tx.Model(&entities).
		Where(masterentities.IncentiveGroupDetail{
			IncentiveGroupDetailId: Id,
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

func (r *IncentiveGroupDetailRepositoryImpl) SaveIncentiveGroupDetail(tx *gorm.DB, req masterpayloads.IncentiveGroupDetailRequest) (bool, *exceptions.BaseErrorResponse) {
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
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}

	return true, nil
}
