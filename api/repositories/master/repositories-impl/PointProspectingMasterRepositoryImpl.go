package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"net/http"

	"gorm.io/gorm"
)

type PointProspectingRepositoryImpl struct {
}

func NewPointProspectingRepositoryImpl() masterrepository.PointProspectingRepository {
	return &PointProspectingRepositoryImpl{}
}

func (r *PointProspectingRepositoryImpl) GetAllPointProspecting(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var datas []masterpayloads.PointProspectingResponse

	query := tx.Model(&masterentities.PointProspecting{}).Select("RECORD_STATUS", "POINT_VARIABLE", "POINT_VALUE", "EFFECTIVE_DATE")

	whereQ := utils.ApplyFilter(query, filterCondition)
	paginatedQuery := whereQ.Scopes(pagination.Paginate(&pages, whereQ))

	err := paginatedQuery.Scan(&datas).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error when scanning data",
			Err:        err,
		}
	}
	// logrus.Debug("datas", datas)

	pages.Rows = datas
	return pages, nil
}

func (r *PointProspectingRepositoryImpl) GetOnePointProspecting(tx *gorm.DB, pointVariable string, pointValue int) (masterpayloads.PointProspectingResponse, *exceptions.BaseErrorResponse) {
	var data masterpayloads.PointProspectingResponse

	err := tx.Model(&masterentities.PointProspecting{}).Select("RECORD_STATUS", "POINT_VARIABLE", "POINT_VALUE", "EFFECTIVE_DATE").
		Where("POINT_VARIABLE = ? AND POINT_VALUE = ?", pointVariable, pointValue).Scan(&data).Error

	if err != nil {
		return masterpayloads.PointProspectingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error when scanning data",
			Err:        err,
		}
	}
	return data, nil
}

func (r *PointProspectingRepositoryImpl) CreatePointProspecting(tx *gorm.DB, req masterpayloads.PointProspectingRequest) (bool, *exceptions.BaseErrorResponse) {
	entities := masterentities.PointProspecting{
		RecordStatus:  req.RecordStatus,
		PointVariable: req.PointVariable,
		PointValue:    req.PointValue,
		EffectiveDate: req.EffectiveDate,
		UserIdCreated: req.UserIdCreated,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error when inserting data",
			Err:        err,
		}
	}
	return true, nil
}

func (r *PointProspectingRepositoryImpl) UpdatePointProspectingStatus(tx *gorm.DB, pointVariable string, pointValue int, req masterpayloads.PointProspectingUpdateStatus) (bool, *exceptions.BaseErrorResponse) {
	err := tx.Model(&masterentities.PointProspecting{}).Where("POINT_VARIABLE = ? AND POINT_VALUE = ?", pointVariable, pointValue).Updates(map[string]interface{}{"RECORD_STATUS": req.RecordStatus}).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error when updating data",
			Err:        err,
		}
	}
	return true, nil
}

func (r *PointProspectingRepositoryImpl) UpdatePointProspectingData(tx *gorm.DB, pointVariable string, pointValue int, req masterpayloads.PointProspectingUpdateRequest) (bool, *exceptions.BaseErrorResponse) {
	err := tx.Model(&masterentities.PointProspecting{}).Where("POINT_VARIABLE = ? AND POINT_VALUE = ?", pointVariable, pointValue).
		Updates(map[string]interface{}{
			"POINT_VARIABLE": req.PointVariable,
			"POINT_VALUE":    req.PointValue,
			"EFFECTIVE_DATE": req.EffectiveDate,
		}).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error when updating data",
			Err:        err,
		}
	}
	return true, nil
}


