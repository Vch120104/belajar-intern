package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type MovingCodeRepositoryImpl struct {
}

// ActivateMovingCode implements masterrepository.MovingCodeRepository.
func (r *MovingCodeRepositoryImpl) ActivateMovingCode(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	multiId := strings.Split(id, ",")
	entities := masterentities.MovingCode{}

	for _, value := range multiId {
		id, _ := strconv.Atoi(value)
		if err := tx.Model(entities).Where(masterentities.MovingCode{MovingCodeId: id}).Update("is_active", true).Error; err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return true, nil

}

// DeactiveMovingCode implements masterrepository.MovingCodeRepository.
func (r *MovingCodeRepositoryImpl) DeactiveMovingCode(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	multiId := strings.Split(id, ",")
	entities := masterentities.MovingCode{}

	for _, value := range multiId {
		id, _ := strconv.Atoi(value)
		if err := tx.Model(entities).Where(masterentities.MovingCode{MovingCodeId: id}).Update("is_active", false).Error; err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return true, nil
}

// GetDropdownMovingCode implements masterrepository.MovingCodeRepository.
func (r *MovingCodeRepositoryImpl) GetDropdownMovingCode(tx *gorm.DB, companyId int) ([]masterpayloads.MovingCodeDropDown, *exceptions.BaseErrorResponse) {

	entities := masterentities.MovingCode{}

	responses := []masterpayloads.MovingCodeDropDown{}

	if err := tx.Model(entities).Where(masterentities.MovingCode{CompanyId: companyId}).Order("priority asc").Scan(&responses).Error; err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {

		if err := tx.Model(entities).Where("mtr_moving_code.company_id LIKE '0'").Order("priority asc").Scan(&responses).Error; err != nil {
			return responses, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		if len(responses) == 0 {
			return responses, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New(""),
			}

		}

	}
	return responses, nil
}

// ChangeStatusMovingCode implements masterrepository.MovingCodeRepository.
func (r *MovingCodeRepositoryImpl) ChangeStatusMovingCode(tx *gorm.DB, Id int) (any, *exceptions.BaseErrorResponse) {
	var entities masterentities.MovingCode

	result := tx.Model(&entities).
		Where(masterentities.MovingCode{MovingCodeId: Id}).
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

// GetMovingCodebyId implements masterrepository.MovingCodeRepository.
func (r *MovingCodeRepositoryImpl) GetMovingCodebyId(tx *gorm.DB, Id int) (masterpayloads.MovingCodeResponse, *exceptions.BaseErrorResponse) {
	model := masterentities.MovingCode{}
	responses := masterpayloads.MovingCodeResponse{}

	err := tx.Model(&model).Where(masterentities.MovingCode{MovingCodeId: Id}).Select("mtr_moving_code.*").First(&responses).Error

	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return responses, nil

}

// CreateMovingCode implements masterrepository.MovingCodeRepository.
func (r *MovingCodeRepositoryImpl) CreateMovingCode(tx *gorm.DB, req masterpayloads.MovingCodeListRequest) (bool, *exceptions.BaseErrorResponse) {
	model := masterentities.MovingCode{}
	var responses []masterentities.MovingCode

	if req.CompanyId != 0 {
		// CHECK COMPANY
		companyResponse, errCompany := generalserviceapiutils.GetCompanyDataById(req.CompanyId)
		if errCompany != nil {
			return false, errCompany
		}

		if companyResponse.CompanyId == 0 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("failed to find company"),
			}
		}

		// CHECK COMPANY HAS MOVING CODE
		if err := tx.Model(&model).Where(masterentities.MovingCode{CompanyId: req.CompanyId}).Scan(&responses).Error; err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// IF DOESN'T HAVE ANY MOVING CODE, INSERT ALL COMPANY_CODE = 0 MOVING CODES
		if len(responses) == 0 {
			if err := tx.Model(&model).Where("mtr_moving_code.company_id LIKE '0'").Scan(&responses).Error; err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}

			for _, value := range responses {
				entities := masterentities.MovingCode{
					CompanyId:             req.CompanyId,
					MovingCode:            value.MovingCode,
					MovingCodeDescription: value.MovingCodeDescription,
					MinimumQuantityDemand: value.MinimumQuantityDemand,
					AgingMonthFrom:        value.AgingMonthFrom,
					AgingMonthTo:          value.AgingMonthTo,
					DemandExistMonthFrom:  value.AgingMonthFrom,
					Priority:              value.Priority,
					DemandExistMonthTo:    value.AgingMonthTo,
					LastMovingMonthFrom:   value.LastMovingMonthFrom,
					LastMovingMonthTo:     value.LastMovingMonthTo,
					Remark:                value.Remark,
				}

				if err := tx.Save(&entities).Error; err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Err:        err,
					}
				}
			}
		}

		// GENERATE NEW PRIORITY
		var priority int64
		if err := tx.Model(&model).Where(masterentities.MovingCode{CompanyId: req.CompanyId}).Count(&priority).Error; err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// SAVE NEW DATA
		entities := masterentities.MovingCode{
			CompanyId:             req.CompanyId,
			MovingCode:            req.MovingCode,
			MovingCodeDescription: req.MovingCodeDescription,
			MinimumQuantityDemand: req.MinimumQuantityDemand,
			AgingMonthFrom:        req.AgingMonthFrom,
			AgingMonthTo:          req.AgingMonthTo,
			DemandExistMonthFrom:  req.AgingMonthFrom,
			Priority:              float64(priority + 1),
			DemandExistMonthTo:    req.AgingMonthTo,
			LastMovingMonthFrom:   req.LastMovingMonthFrom,
			LastMovingMonthTo:     req.LastMovingMonthTo,
			Remark:                req.Remark,
		}

		if err := tx.Save(&entities).Error; err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusConflict,
					Err:        err,
				}
			}
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		return true, nil
	}

	// HANDLE CASE FOR COMPANY_ID == 0

	// GENERATE NEW PRIORITY FOR DEFAULT MOVING CODES
	var priority int64
	if err := tx.Model(&model).Where("mtr_moving_code.company_id LIKE '0'").Count(&priority).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// SAVE NEW DEFAULT DATA
	entities := masterentities.MovingCode{
		CompanyId:             req.CompanyId,
		MovingCode:            req.MovingCode,
		MovingCodeDescription: req.MovingCodeDescription,
		MinimumQuantityDemand: req.MinimumQuantityDemand,
		AgingMonthFrom:        req.AgingMonthFrom,
		AgingMonthTo:          req.AgingMonthTo,
		DemandExistMonthFrom:  req.AgingMonthFrom,
		Priority:              float64(priority + 1),
		DemandExistMonthTo:    req.AgingMonthTo,
		LastMovingMonthFrom:   req.LastMovingMonthFrom,
		LastMovingMonthTo:     req.LastMovingMonthTo,
		Remark:                req.Remark,
	}

	if err := tx.Save(&entities).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

// GetAllMovingCode implements masterrepository.MovingCodeRepository.
func (r *MovingCodeRepositoryImpl) GetAllMovingCode(tx *gorm.DB, companyId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	model := masterentities.MovingCode{}
	var responses []masterentities.MovingCode

	whereQuery := tx.Model(&model).Where("mtr_moving_code.company_id = ?", companyId)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []masterentities.MovingCode{}
		return pages, nil
	}

	pages.Rows = responses

	return pages, nil
}

// PushMovingCodePriority implements masterrepository.MovingCodeRepository.
func (r *MovingCodeRepositoryImpl) PushMovingCodePriority(tx *gorm.DB, companyId int, Id int) (bool, *exceptions.BaseErrorResponse) {

	currentModel := masterentities.MovingCode{}
	nextIndexModel := masterentities.MovingCode{}

	//Current index

	err := tx.Model(&currentModel).Where(masterentities.MovingCode{CompanyId: companyId}).Where(masterentities.MovingCode{MovingCodeId: Id}).First(&currentModel).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if currentModel.Priority == 1 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("priority already 1"),
		}
	}
	fmt.Println(currentModel)

	//Next priority index

	err = tx.Model(&currentModel).Where(masterentities.MovingCode{CompanyId: companyId}).Where(masterentities.MovingCode{Priority: currentModel.Priority - 1}).First(&nextIndexModel).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	fmt.Println(nextIndexModel)
	//PUSH PRIORITY

	currentModel.Priority -= 1

	pushPriority := tx.Save(&currentModel)

	if pushPriority.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        pushPriority.Error,
		}
	}

	//DECREASE NEXT PRIORITY

	nextIndexModel.Priority += 1

	decreasePriority := tx.Save(&nextIndexModel)

	if decreasePriority.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        decreasePriority.Error,
		}
	}

	return true, nil
}

// UpdateMovingCode implements masterrepository.MovingCodeRepository.
func (r *MovingCodeRepositoryImpl) UpdateMovingCode(tx *gorm.DB, req masterpayloads.MovingCodeListUpdate) (bool, *exceptions.BaseErrorResponse) {

	model := masterentities.MovingCode{}
	if err := tx.Model(&model).Where(masterentities.MovingCode{MovingCodeId: req.MovingCodeId}).First(&model).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if model == (masterentities.MovingCode{}) {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	entities := masterentities.MovingCode{
		MovingCodeId:          req.MovingCodeId,
		MovingCodeDescription: req.MovingCodeDescription,
		MinimumQuantityDemand: req.MinimumQuantityDemand,
		AgingMonthFrom:        req.AgingMonthFrom,
		AgingMonthTo:          req.AgingMonthTo,
		DemandExistMonthFrom:  req.AgingMonthFrom,
		DemandExistMonthTo:    req.AgingMonthTo,
		LastMovingMonthFrom:   req.LastMovingMonthFrom,
		LastMovingMonthTo:     req.LastMovingMonthTo,
		Remark:                req.Remark,
	}

	err := tx.Updates(&entities).Where(masterentities.MovingCode{MovingCodeId: req.MovingCodeId}).Error

	if err != nil {

		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func StartMovingCodeRepositoryImpl() masterrepository.MovingCodeRepository {
	return &MovingCodeRepositoryImpl{}
}
