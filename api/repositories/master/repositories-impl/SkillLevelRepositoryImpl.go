package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type SkillLevelRepositoryImpl struct {
}

func StartSkillLevelRepositoryImpl() masterrepository.SkillLevelRepository {
	return &SkillLevelRepositoryImpl{}
}

func (r *SkillLevelRepositoryImpl) GetAllSkillLevel(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := masterentities.SkillLevel{}
	responses := []masterpayloads.SkillLevelResponse{}

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

func (r *SkillLevelRepositoryImpl) GetSkillLevelById(tx *gorm.DB, Id int) (masterpayloads.SkillLevelResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.SkillLevel{}
	response := masterpayloads.SkillLevelResponse{}

	rows, err := tx.Model(&entities).
		Where(masterentities.SkillLevel{
			SkillLevelId: Id,
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

func (r *SkillLevelRepositoryImpl) GetSkillLevelByCode(tx *gorm.DB, Code string) (masterpayloads.SkillLevelResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.SkillLevel{}
	response := masterpayloads.SkillLevelResponse{}

	rows, err := tx.Model(&entities).
		Where(masterentities.SkillLevel{
			SkillLevelCode: Code,
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

func (r *SkillLevelRepositoryImpl) SaveSkillLevel(tx *gorm.DB, req masterpayloads.SkillLevelResponse) (bool, *exceptions.BaseErrorResponse) {
	entities := masterentities.SkillLevel{
		IsActive:              req.IsActive,
		SkillLevelId:          req.SkillLevelId,
		SkillLevelCode:        req.SkillLevelCode,
		SkillLevelDescription: req.SkillLevelDescription,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {

			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return true, nil
}

func (r *SkillLevelRepositoryImpl) ChangeStatusSkillLevel(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masterentities.SkillLevel

	result := tx.Model(&entities).
		Where(masterentities.SkillLevel{SkillLevelId: Id}).
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
