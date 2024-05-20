package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptionsss_test "after-sales/api/expectionsss"
	"errors"
	"net/http"

	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IncentiveGroupRepositoryImpl struct {
}

func StartIncentiveGroupRepositoryImpl() masterrepository.IncentiveGroupRepository {
	return &IncentiveGroupRepositoryImpl{}
}

func (r *IncentiveGroupRepositoryImpl) GetAllIncentiveGroup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	IncentiveGroupMapping := []masterentities.IncentiveGroup{}
	IncentiveGroupResponse := []masterpayloads.IncentiveGroupResponse{}
	// IncentiveGroupResponse1 := masterpayloads.IncentiveGroupResponse{}
	query := tx.
		Model(masterentities.IncentiveGroup{}).
		Scan(&IncentiveGroupResponse)
		// Select("email").
		// Where("id in (?)", userIDs).
		// Scan(email).
		// Rows()

	ApplyFilter := utils.ApplyFilter(query, filterCondition)

	err := ApplyFilter.
		Scopes(pagination.Paginate(&IncentiveGroupMapping, &pages, ApplyFilter)).
		// Order("approval.name").
		Scan(&IncentiveGroupResponse).
		Error

	if len(IncentiveGroupResponse) == 0 {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {

		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	// defer row.Close()
	pages.Rows = IncentiveGroupResponse

	return pages, nil
}

func (r *IncentiveGroupRepositoryImpl) GetAllIncentiveGroupIsActive(tx *gorm.DB) ([]masterpayloads.IncentiveGroupResponse, *exceptionsss_test.BaseErrorResponse) {
	// var IncentiveGroupResponse masterpayloads.IncentiveGroupResponse
	IncentiveGroupResponse := []masterpayloads.IncentiveGroupResponse{}

	row, err := tx.
		Model(masterentities.IncentiveGroup{}).
		Where(masterentities.IncentiveGroup{IsActive: true}).
		Scan(&IncentiveGroupResponse).
		Rows()

	if len(IncentiveGroupResponse) == 0 {
		return IncentiveGroupResponse, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {

		return IncentiveGroupResponse, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer row.Close()

	return IncentiveGroupResponse, nil
}

func (*IncentiveGroupRepositoryImpl) GetIncentiveGroupById(tx *gorm.DB, Id int) (masterpayloads.IncentiveGroupResponse, *exceptionsss_test.BaseErrorResponse) {
	var IncentiveGroupMapping masterentities.IncentiveGroup
	var IncentiveGroupResponse masterpayloads.IncentiveGroupResponse

	rows, err := tx.
		Model(&IncentiveGroupMapping).
		Where(masterentities.IncentiveGroup{IncentiveGroupId: Id}).
		First(&IncentiveGroupResponse).
		Rows()

	if err != nil {

		return IncentiveGroupResponse, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	return IncentiveGroupResponse, nil
}

func (r *IncentiveGroupRepositoryImpl) SaveIncentiveGroup(tx *gorm.DB, req masterpayloads.IncentiveGroupResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	IncentiveGroup := masterentities.IncentiveGroup{
		IsActive:           req.IsActive,
		IncentiveGroupId:   req.IncentiveGroupId,
		IncentiveGroupCode: req.IncentiveGroupCode,
		IncentiveGroupName: req.IncentiveGroupName,
		EffectiveDate:      req.EffectiveDate,
	}
	err := tx.
		Create(&IncentiveGroup).
		Error

	if err != nil {
		logrus.Info(err)
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}

	return true, nil
}

func (r *IncentiveGroupRepositoryImpl) ChangeStatusIncentiveGroup(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	// var entities masterentities.IncentiveGroup
	var IncentiveGroupMapping masterentities.IncentiveGroup
	// var IncentiveGroupResponse masterpayloads.IncentiveGroupResponse

	result := tx.
		Model(&IncentiveGroupMapping).
		Where(masterentities.IncentiveGroup{IncentiveGroupId: Id}).
		First(&IncentiveGroupMapping)

	// result := tx.Model(&entities).
	// 	Where("incentive_group_id = ?", Id).
	// 	First(&entities)

	if result.Error != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	if IncentiveGroupMapping.IsActive {
		IncentiveGroupMapping.IsActive = false
	} else {
		IncentiveGroupMapping.IsActive = true
	}

	result = tx.Save(&IncentiveGroupMapping)

	if result.Error != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}

func (r *IncentiveGroupRepositoryImpl) UpdateIncentiveGroup(tx *gorm.DB, req masterpayloads.UpdateIncentiveGroupRequest) (bool, *exceptionsss_test.BaseErrorResponse) {

	model := masterentities.IncentiveGroup{}
	if err := tx.Model(&model).Where(masterentities.IncentiveGroup{IncentiveGroupId: req.IncentiveGroupId}).First(&model).Error; err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if model == (masterentities.IncentiveGroup{}) {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	entities := masterentities.IncentiveGroup{
		IncentiveGroupId:   req.IncentiveGroupId,
		IncentiveGroupCode: req.IncentiveGroupCode,
		IncentiveGroupName: req.IncentiveGroupName,
		EffectiveDate:      req.EffectiveDate,
	}

	err := tx.Updates(&entities).Where(masterentities.IncentiveGroup{IncentiveGroupId: req.IncentiveGroupId}).Error

	if err != nil {

		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *IncentiveGroupRepositoryImpl) GetAllIncentiveGroupDropDown(tx *gorm.DB) ([]masterpayloads.IncentiveGroupDropDown, *exceptionsss_test.BaseErrorResponse) {
	// var IncentiveGroupResponse masterpayloads.IncentiveGroupResponse
	DropDownResponse := []masterpayloads.IncentiveGroupDropDown{}

	err := tx.Model(masterentities.IncentiveGroup{}).Select("mtr_incentive_group.*, CONCAT(incentive_group_code, ' - ', incentive_group_name) AS incentive_group_code_name").Where("is_active = 'true'").Find(&DropDownResponse).Error
	if err != nil {
		return DropDownResponse, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	return DropDownResponse, nil
}
