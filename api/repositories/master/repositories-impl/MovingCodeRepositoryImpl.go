package masterrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"errors"
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
func (r *MovingCodeRepositoryImpl) GetDropdownMovingCode(tx *gorm.DB) ([]masterpayloads.MovingCodeDropDown, *exceptions.BaseErrorResponse) {

	entities := masterentities.MovingCode{}

	responses := []masterpayloads.MovingCodeDropDown{}

	if err := tx.Model(entities).Scan(&responses).Error; err != nil {
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
	companyResponses := masterpayloads.CompanyResponse{}

	err := tx.Model(&model).Where(masterentities.MovingCode{MovingCodeId: Id}).Select("mtr_moving_code.*").First(&responses).Error

	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	companyByIdUrl := config.EnvConfigs.GeneralServiceUrl + "/company-list/" + strconv.Itoa(responses.CompanyId)

	if errUrlCompany := utils.Get(companyByIdUrl, &companyResponses, nil); errUrlCompany != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	if companyResponses == (masterpayloads.CompanyResponse{}) {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	responses.CompanyName = &companyResponses.CompanyName

	return responses, nil

}

// CreateMovingCode implements masterrepository.MovingCodeRepository.
func (r *MovingCodeRepositoryImpl) CreateMovingCode(tx *gorm.DB, req masterpayloads.MovingCodeListRequest) (bool, *exceptions.BaseErrorResponse) {

	//CHECK COMPANY ID
	companyResponses := masterpayloads.CompanyResponse{}

	companyByIdUrl := config.EnvConfigs.GeneralServiceUrl + "/company-list/" + strconv.Itoa(req.CompanyId)

	if errUrlCompany := utils.Get(companyByIdUrl, &companyResponses, nil); errUrlCompany != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	if companyResponses == (masterpayloads.CompanyResponse{}) {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	//GENERATE PRIORITY

	var priority int64

	model := masterentities.MovingCode{}

	if err := tx.Model(&model).Count(&priority).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	//SAVE
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

	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

// GetAllMovingCode implements masterrepository.MovingCodeRepository.
func (r *MovingCodeRepositoryImpl) GetAllMovingCode(tx *gorm.DB, pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse) {
	model := masterentities.MovingCode{}
	var responses []masterentities.MovingCode
	var companyResponses []masterpayloads.CompanyResponse

	err := tx.Model(&model).Scan(&responses).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	companyUrl := config.EnvConfigs.GeneralServiceUrl + "/company-list-all"

	if errUrlCompany := utils.Get(companyUrl, &companyResponses, nil); errUrlCompany != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	joinedData := utils.DataFrameInnerJoin(responses, companyResponses, "CompanyId")

	if len(joinedData) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return dataPaginate, totalPages, totalRows, nil

}

// PushMovingCodePriority implements masterrepository.MovingCodeRepository.
func (r *MovingCodeRepositoryImpl) PushMovingCodePriority(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {

	currentModel := masterentities.MovingCode{}
	nextIndexModel := masterentities.MovingCode{}

	//Current index

	err := tx.Model(&currentModel).Where(masterentities.MovingCode{MovingCodeId: Id}).First(&currentModel).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if currentModel.Priority == 1 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	//Next priority index

	err = tx.Model(&currentModel).Where(masterentities.MovingCode{Priority: currentModel.Priority - 1}).First(&nextIndexModel).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

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
