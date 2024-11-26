package transactionjpcbrepositoryimpl

import (
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type BayMasterImpl struct {
}

func NewCarWashBayRepositoryImpl() transactionjpcbrepository.BayMasterRepository {
	return &BayMasterImpl{}
}

func (*BayMasterImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	responses := []transactionjpcbpayloads.CarWashBayGetAllResponse{}

	keyAttributes := []string{
		"car_wash_bay_id",
		"car_wash_bay_code",
		"is_active",
		"car_wash_bay_description",
	}

	query := tx.Model(&transactionjpcbentities.BayMaster{}).Select(keyAttributes)
	whereQuery := utils.ApplyFilter(query, filterCondition)

	rows, err := whereQuery.Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	for rows.Next() {
		var carWashBayId int
		var carWashBayCode, carWashBayDescription string
		var isActive bool

		err := rows.Scan(&carWashBayId, &carWashBayCode, &isActive, &carWashBayDescription)

		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		responseMap := transactionjpcbpayloads.CarWashBayGetAllResponse{
			CarWashBayId:          carWashBayId,
			CarWashBayCode:        carWashBayCode,
			IsActive:              isActive,
			CarWashBayDescription: carWashBayDescription,
		}
		responses = append(responses, responseMap)
	}

	var mapResponses []map[string]interface{}

	for _, response := range responses {
		responseMap := map[string]interface{}{
			"car_wash_bay_id":          response.CarWashBayId,
			"car_wash_bay_code":        response.CarWashBayCode,
			"is_active":                response.IsActive,
			"car_wash_bay_description": response.CarWashBayDescription,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)
	return paginatedData, totalPages, totalRows, nil
}

func (r *BayMasterImpl) GetCarWashBayById(tx *gorm.DB, carWashBayId int) (transactionjpcbentities.BayMaster, *exceptions.BaseErrorResponse) {
	var entity transactionjpcbentities.BayMaster

	err := tx.Model(&entity).Where(transactionjpcbentities.BayMaster{
		CarWashBayId: carWashBayId,
	}).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not Found",
				Err:        err,
			}
		}
		return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch BOM Master record",
			Err:        err,
		}
	}

	return entity, nil
}

func (*BayMasterImpl) GetAllActive(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var responses []map[string]interface{}

	keyAttributes := []string{
		"car_wash_bay_id",
		"car_wash_bay_description",
		"order_number",
	}

	var companyIdFilter int
	for _, condition := range filterCondition {
		if condition.ColumnField == "company_id" {
			result, err := strconv.Atoi(condition.ColumnValue)
			if err != nil {
				companyIdFilter = 0
			}
			companyIdFilter = result
		}
	}

	rows, err := tx.Model(&transactionjpcbentities.BayMaster{}).Select(keyAttributes).Where("company_id = ? AND is_active = 1", companyIdFilter).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	for rows.Next() {
		var carWashBayId int
		var carWashBayDescription string
		var orderNumber int

		err := rows.Scan(&carWashBayId, &carWashBayDescription, &orderNumber)

		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		responseMap := map[string]interface{}{
			"car_wash_bay_id":          carWashBayId,
			"car_wash_bay_description": carWashBayDescription,
			"order_number":             orderNumber,
		}
		responses = append(responses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(responses, &pages)
	return paginatedData, totalPages, totalRows, nil
}

func (*BayMasterImpl) GetAllDeactive(tx *gorm.DB, filterCondition []utils.FilterCondition) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	var responses []map[string]interface{}

	keyAttributes := []string{
		"car_wash_bay_id",
		"car_wash_bay_description",
		"order_number",
	}

	var companyIdFilter int
	for _, condition := range filterCondition {
		if condition.ColumnField == "company_id" {
			result, err := strconv.Atoi(condition.ColumnValue)
			if err != nil {
				companyIdFilter = 0
			}
			companyIdFilter = result
		}
	}

	rows, err := tx.Model(&transactionjpcbentities.BayMaster{}).Select(keyAttributes).Where("company_id = ? AND order_number = 0", companyIdFilter).Rows()
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	for rows.Next() {
		var carWashId int
		var carWashBayId int

		err := rows.Scan(&carWashId, &carWashBayId)
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		responseMap := map[string]interface{}{
			"car_wash_id":     carWashId,
			"car_wash_bay_id": carWashBayId,
		}
		responses = append(responses, responseMap)
	}

	return responses, nil
}

func (r *BayMasterImpl) ChangeStatus(tx *gorm.DB, request transactionjpcbpayloads.CarWashBayUpdateRequest) (transactionjpcbentities.BayMaster, *exceptions.BaseErrorResponse) {
	carWashEntities := []transactionjpcbentities.CarWash{}
	var bayEntity transactionjpcbentities.BayMaster

	result := tx.Select("work_order_system_number").Where("company_id = ? AND car_wash_bay_id = ? AND car_wash_status_id = 2", request.CompanyId, request.CarWashBayId).
		Find(&carWashEntities)

	if result.Error != nil {
		return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	if len(carWashEntities) == 0 {
		bayEntities := []transactionjpcbentities.BayMaster{}

		resultBay := tx.Model(&bayEntities).Where("company_id = ? AND car_wash_bay_id = ?", request.CompanyId, request.CarWashBayId).
			Find(&bayEntities)
		if resultBay.Error != nil {
			return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        resultBay.Error,
			}
		}

		if len(bayEntities) != 0 {
			updateQuery := tx.Model(&bayEntity).Where("car_wash_bay_id = ?", request.CarWashBayId).Update("is_active", request.IsActive)
			if updateQuery.Error != nil {
				return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        updateQuery.Error,
				}
			}

			resetErr := resetAllOrderNumber(tx, request.CompanyId)
			if resetErr != nil {
				errorReset := errors.New("reset order fail")
				return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        errorReset,
				}
			}

			reorderErr := reorderOrderNumber(tx, request.CompanyId)
			if reorderErr != nil {
				return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        reorderErr.Err,
				}
			}

			selectQuery := tx.Model(&bayEntity).Where("car_wash_bay_id = ?", request.CarWashBayId).First(&bayEntity)
			if selectQuery.Error != nil {
				return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        selectQuery.Error,
				}
			}
			return bayEntity, nil
		} else {
			bayNotFoundError := errors.New("bay not found")
			return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        bayNotFoundError,
			}
		}
	}

	return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
		StatusCode: http.StatusBadRequest,
		Err:        fmt.Errorf("already start"),
	}

}

func (r *BayMasterImpl) PostCarWashBay(tx *gorm.DB, request transactionjpcbpayloads.CarWashBayPostRequest) (transactionjpcbentities.BayMaster, *exceptions.BaseErrorResponse) {
	var entities transactionjpcbentities.BayMaster

	entities.IsActive = true
	entities.CarWashBayCode = request.CarWashBayCode
	entities.CarWashBayDescription = request.CarWashBayDescription
	entities.CompanyId = request.CompanyId

	err := tx.Create(&entities).Error
	if err != nil {
		return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	resetErr := resetAllOrderNumber(tx, request.CompanyId)
	if resetErr != nil {
		errorReset := errors.New("reset order fail")
		return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errorReset,
		}
	}

	reorderErr := reorderOrderNumber(tx, request.CompanyId)
	if reorderErr != nil {
		return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        reorderErr.Err,
		}
	}
	return entities, nil
}

func (r *BayMasterImpl) UpdateCarWashBay(tx *gorm.DB, request transactionjpcbpayloads.CarWashBayPutRequest) (transactionjpcbentities.BayMaster, *exceptions.BaseErrorResponse) {
	var entities transactionjpcbentities.BayMaster

	result := tx.Model(&entities).Where(transactionjpcbentities.BayMaster{
		CarWashBayId: request.CarWashBayID,
	}).First(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        fmt.Errorf("bay with id %d not found", request.CarWashBayID),
			}
		}
		return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	entities.CarWashBayId = request.CarWashBayID
	entities.CarWashBayCode = request.CarWashBayCode
	entities.CarWashBayDescription = request.CarWashBayDescription

	result = tx.Model(&entities).Where(transactionjpcbentities.BayMaster{
		CarWashBayId: request.CarWashBayID,
	}).Updates(map[string]interface{}{
		"car_wash_bay_id":          request.CarWashBayID,
		"car_wash_bay_code":        request.CarWashBayCode,
		"car_wash_bay_description": request.CarWashBayDescription,
	})
	if result.Error != nil {
		return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return entities, nil
}

func reorderOrderNumber(tx *gorm.DB, companyId int) *exceptions.BaseErrorResponse {
	keyAttributes := []string{
		"company_id",
		"car_wash_bay_id",
	}

	query := tx.Model(&transactionjpcbentities.BayMaster{}).Select(keyAttributes).Where("company_id = ? AND is_active = 1", companyId)

	rows, err := query.Rows()
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	for rows.Next() {
		var companyId, carWashBayId int

		err := rows.Scan(&companyId, &carWashBayId)
		if err != nil {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		highestOrderNumber, errHighestNumber := getHighestOrderNumber(tx, companyId)
		if errHighestNumber != nil {
			err := errors.New("error highest number")
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		setOrderQuery := tx.Model(&transactionjpcbentities.BayMaster{}).Where("car_wash_bay_id = ?", carWashBayId).Update("order_number", highestOrderNumber+1)
		if setOrderQuery.Error != nil {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        setOrderQuery.Error,
			}
		}
	}

	return nil
}

func getHighestOrderNumber(tx *gorm.DB, companyId int) (int, *exceptions.BaseErrorResponse) {
	highestOrder := 0

	joinQuery := tx.Model(&transactionjpcbentities.BayMaster{}).Select("MAX(order_number)").Where("company_id = ?", companyId).Find(&highestOrder)

	if joinQuery.Error != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        joinQuery.Error,
		}
	}

	return highestOrder, nil
}

func resetAllOrderNumber(tx *gorm.DB, companyId int) *exceptions.BaseErrorResponse {
	var bayEntities transactionjpcbentities.BayMaster
	keyAttributes := []string{
		"car_wash_bay_id",
		"car_wash_bay_code",
		"is_active",
		"car_wash_bay_description",
	}

	joinQuery := tx.Model(&bayEntities).Select(keyAttributes).Where("company_id = ?", companyId)

	rows, err := joinQuery.Rows()
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	for rows.Next() {
		var carWashBayId int
		var carWashBayCode, carWashBayDescription string
		var isActive bool

		err := rows.Scan(&carWashBayId, &carWashBayCode, &isActive, &carWashBayDescription)

		if err != nil {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		resetOrderNumberQuery := tx.Model(&bayEntities).Where("car_wash_bay_id = ?", carWashBayId).Update("order_number", 0)

		if resetOrderNumberQuery.Error != nil {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        resetOrderNumberQuery.Error,
			}
		}
	}

	return nil
}

func (*BayMasterImpl) CarWashBayDropDown(tx *gorm.DB, filterCondition []utils.FilterCondition) ([]transactionjpcbpayloads.CarWashBayDropDownResponse, *exceptions.BaseErrorResponse) {
	var responses []transactionjpcbpayloads.CarWashBayDropDownResponse

	keyAttributes := []string{
		"car_wash_bay_id",
		"car_wash_bay_code",
		"car_wash_bay_description",
		"is_active",
	}

	query := tx.Model(&transactionjpcbentities.BayMaster{}).Select(keyAttributes)
	whereQuery := utils.ApplyFilter(query, filterCondition)

	rows, err := whereQuery.Rows()
	if err != nil {
		return []transactionjpcbpayloads.CarWashBayDropDownResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	for rows.Next() {
		var carWashBayId int
		var carWashBayCode, carWashBayDescription string
		var isActive bool

		err := rows.Scan(&carWashBayId, &carWashBayCode, &carWashBayDescription, &isActive)

		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		responseMap := transactionjpcbpayloads.CarWashBayDropDownResponse{
			CarWashBayId:          carWashBayId,
			CarWashBayCode:        carWashBayCode,
			CarWashBayDescription: carWashBayDescription,
			IsActive:              isActive,
		}
		responses = append(responses, responseMap)
	}

	return responses, nil
}
