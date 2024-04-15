package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"net/http"

	"gorm.io/gorm"
)

type SkillLevelRepositoryImpl struct {
}

func StartSkillLevelRepositoryImpl() masterrepository.SkillLevelRepository {
	return &SkillLevelRepositoryImpl{}
}

func (r *SkillLevelRepositoryImpl) GetAllSkillLevel(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.SkillLevel{}
	responses := []masterpayloads.SkillLevelResponse{}

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

func (r *SkillLevelRepositoryImpl) GetSkillLevelById(tx *gorm.DB, Id int) (masterpayloads.SkillLevelResponse, *exceptionsss_test.BaseErrorResponse) {
	var entities masterentities.SkillLevel
	var response masterpayloads.SkillLevelResponse

	rows, err := tx.Model(&entities).
		Where(masterentities.SkillLevel{
			SkillLevelId: Id,
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

func (r *SkillLevelRepositoryImpl) GetSkillLevelByCode(tx *gorm.DB, Code string) (masterpayloads.SkillLevelResponse, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.SkillLevel{}
	response := masterpayloads.SkillLevelResponse{}

	rows, err := tx.Model(&entities).
		Where(masterentities.SkillLevel{
			SkillLevelCode: Code,
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

func (r *SkillLevelRepositoryImpl) SaveSkillLevel(tx *gorm.DB, req masterpayloads.SkillLevelResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.SkillLevel{
		IsActive:              req.IsActive,
		SkillLevelId:          req.SkillLevelCodeId,
		SkillLevelCode:        req.SkillLevelCodeValue,
		SkillLevelDescription: req.SkillLevelCodeDescription,
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

func (r *SkillLevelRepositoryImpl) ChangeStatusSkillLevel(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities masterentities.SkillLevel

	result := tx.Model(&entities).
		Where(masterentities.SkillLevel{SkillLevelId: Id}).
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
