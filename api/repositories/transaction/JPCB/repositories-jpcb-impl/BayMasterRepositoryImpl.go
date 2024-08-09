package transactionjpcbrepositoryimpl

import (
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type BayMasterImpl struct {
}

func OpenBayMasterRepositoryImpl() transactionjpcbrepository.BayMasterRepository {
	return &BayMasterImpl{}
}

func (*BayMasterImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	responses := []transactionjpcbpayloads.BayMasterResponse{}

	mainTable := "trx_car_wash"
	mainAlias := "carwash"
	mainAliasBay := "bay"

	joinTables := []utils.JoinTable{
		{Table: "mtr_car_wash_bay", Alias: "bay", ForeignKey: mainAlias + ".car_wash_bay_id", ReferenceKey: "bay.car_wash_bay_id"},
	}

	joinQuery := utils.CreateJoin(tx, mainTable, mainAlias, joinTables...)

	keyAttributes := []string{
		mainAlias + ".company_id",
		mainAlias + ".car_wash_id",
		mainAlias + ".car_wash_status_id",
		mainAlias + ".work_order_system_number",
		mainAliasBay + ".car_wash_bay_id",
		mainAliasBay + ".car_wash_bay_code",
		mainAliasBay + ".is_active",
		mainAliasBay + ".car_wash_bay_description",
	}

	joinQuery = joinQuery.Select(keyAttributes)
	whereQuery := utils.ApplyFilter(joinQuery, filterCondition)

	rows, err := whereQuery.Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	for rows.Next() {
		var companyId, carWashId, carWashBayId, carWashStatusId, workOrderSystemNumber int
		var carWashBayCode, carWashBayDescription string
		var isActive bool

		err := rows.Scan(&companyId, &carWashId, &carWashStatusId, &workOrderSystemNumber, &carWashBayId, &carWashBayCode, &isActive, &carWashBayDescription)

		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		responseMap := transactionjpcbpayloads.BayMasterResponse{
			CompanyId:             companyId,
			CarWashId:             carWashId,
			CarWashBayId:          carWashBayId,
			CarWashBayCode:        carWashBayCode,
			IsActive:              isActive,
			CarWashBayDescription: carWashBayDescription,
			CarWashStatusId:       carWashStatusId,
			WorkOrderSystemNumber: workOrderSystemNumber,
		}
		responses = append(responses, responseMap)
	}

	var mapResponses []map[string]interface{}

	for _, response := range responses {
		responseMap := map[string]interface{}{
			"company_id":               response.CompanyId,
			"car_wash_id":              response.CarWashId,
			"car_wash_bay_id":          response.CarWashBayId,
			"car_wash_bay_code":        response.CarWashBayCode,
			"is_active":                response.IsActive,
			"car_wash_bay_description": response.CarWashBayDescription,
			"car_wash_status_id":       response.CarWashStatusId,
			"work_order_system_number": response.WorkOrderSystemNumber,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)
	return paginatedData, totalPages, totalRows, nil
}

func (*BayMasterImpl) GetAllActive(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var responses []map[string]interface{}

	mainTable := "trx_car_wash"
	mainAlias := "carwash"

	joinTables := []utils.JoinTable{
		{Table: "mtr_car_wash_bay", Alias: "bay", ForeignKey: "bay.car_wash_bay_id", ReferenceKey: "bay.car_wash_bay_id"},
	}

	joinQuery := utils.CreateJoin(tx, mainTable, mainAlias, joinTables...)

	keyAttributes := []string{
		"carwash.car_wash_id",
		"bay.car_wash_bay_id",
		"bay.car_wash_bay_description",
		"carwash.car_wash_status_id",
		"carwash.work_order_system_number",
	}

	joinQuery = joinQuery.Select(keyAttributes)

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

	rows, err := joinQuery.Where("company_id = ? AND car_wash_status_id = 1", companyIdFilter).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	for rows.Next() {
		var carWashId int
		var carWashBayId int
		var carWashBayDescription string
		var carWashStatusId string
		var WorkOrderSystemNumber int

		err := rows.Scan(&carWashId, &carWashBayId, &carWashBayDescription, &carWashStatusId, &WorkOrderSystemNumber)

		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		responseMap := map[string]interface{}{
			"car_wash_id":              carWashId,
			"car_wash_bay_id":          carWashBayId,
			"car_wash_bay_description": carWashBayDescription,
			"car_wash_bay_status_id":   carWashStatusId,
			"work_order_system_number": WorkOrderSystemNumber,
		}
		responses = append(responses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(responses, &pages)
	return paginatedData, totalPages, totalRows, nil
}

func (*BayMasterImpl) GetAllDeactive(tx *gorm.DB, filterCondition []utils.FilterCondition) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	var responses []map[string]interface{}

	mainTable := "trx_car_wash"
	mainAlias := "carwash"

	joinTables := []utils.JoinTable{
		{Table: "mtr_car_wash_bay", Alias: "bay", ForeignKey: "bay.car_wash_bay_id", ReferenceKey: "bay.car_wash_bay_id"},
	}

	joinQuery := utils.CreateJoin(tx, mainTable, mainAlias, joinTables...)

	keyAttributes := []string{
		"carwash.car_wash_id",
		"bay.car_wash_bay_id",
	}

	joinQuery = joinQuery.Select(keyAttributes)

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

	rows, err := joinQuery.Where("company_id = ? AND work_order_system_number = 0", companyIdFilter).Rows()
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	//appending data for response
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

func (r *BayMasterImpl) Update(tx *gorm.DB, request transactionjpcbpayloads.BayMasterUpdateRequest) (transactionjpcbentities.BayMaster, *exceptions.BaseErrorResponse) {
	carWashEntities := []transactionjpcbentities.CarWash{}
	var bayEntities transactionjpcbentities.BayMaster

	result := tx.Select("work_order_system_number").Where("company_id = ? AND car_wash_bay_id = ? AND car_wash_status_id = 3", request.CompanyId, request.CarWashBayId).
		Find(&carWashEntities)

	if result.Error != nil {
		return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	if len(carWashEntities) == 0 {
		carWashEntities2 := []transactionjpcbentities.CarWash{}
		resultBay := tx.Model(&carWashEntities2).Where("company_id = ? AND car_wash_bay_id = ?", request.CompanyId, request.CarWashBayId).
			Find(&carWashEntities2)

		if resultBay.Error != nil {
			return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        resultBay.Error,
			}
		}

		if len(carWashEntities2) != 0 {
			updateQuery := tx.Model(&bayEntities).Where("car_wash_bay_id = ?", request.CarWashBayId).Update("is_active", request.RecordStatus)
			if updateQuery.Error != nil {
				return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        updateQuery.Error,
				}
			}

			// reorder order_no (column isnt yet present in current entity)

			//reset all bay in a company
			resetErr := reset(tx, request.CompanyId)

			if resetErr != nil {
				errorReset := errors.New("reset order fail")
				return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        errorReset,
				}
			}

			reorderErr := reorder(tx, request.CompanyId)

			if reorderErr != nil {
				errorReorder := errors.New("reorder fail")
				return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        errorReorder,
				}
			}

			return bayEntities, nil
		} else {
			error4 := errors.New("bay not found")
			return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        error4,
			}
		}
	}

	newError := errors.New("already start")
	return transactionjpcbentities.BayMaster{}, &exceptions.BaseErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Err:        newError,
	}

}

func reorder(tx *gorm.DB, companyId int) *exceptions.BaseErrorResponse {
	mainTable := "trx_car_wash"
	mainAlias := "carwash"
	mainAliasBay := "bay"

	joinTables := []utils.JoinTable{
		{Table: "mtr_car_wash_bay", Alias: "bay", ForeignKey: mainAlias + ".car_wash_bay_id", ReferenceKey: "bay.car_wash_bay_id"},
	}

	joinQuery := utils.CreateJoin(tx, mainTable, mainAlias, joinTables...)

	keyAttributes := []string{
		mainAlias + ".company_id",
		mainAlias + ".car_wash_id",
		mainAlias + ".car_wash_status_id",
		mainAlias + ".work_order_system_number",
		mainAliasBay + ".car_wash_bay_id",
		mainAliasBay + ".car_wash_bay_code",
		mainAliasBay + ".is_active",
		mainAliasBay + ".car_wash_bay_description",
	}

	//find the active only bay
	joinQuery = joinQuery.Select(keyAttributes).Where("company_id = ? AND bay.is_active = 1", companyId)

	rows, err := joinQuery.Rows()
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	for rows.Next() {
		var companyId int
		var carWashId int
		var carWashStatusId int
		var workOrderSystemNumber int
		var carWashBayId int
		var carWashBayCode string
		var isActive bool
		var carWashBayDescription string

		err := rows.Scan(
			&companyId,
			&carWashId,
			&carWashStatusId,
			&workOrderSystemNumber,
			&carWashBayId,
			&carWashBayCode,
			&isActive,
			&carWashBayDescription,
		)

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
	//TODO get max order_number bay in one company

	mainTable := "trx_car_wash"
	mainAlias := "carwash"

	joinTables := []utils.JoinTable{
		{Table: "mtr_car_wash_bay", Alias: "bay", ForeignKey: mainAlias + ".car_wash_bay_id", ReferenceKey: "bay.car_wash_bay_id"},
	}

	joinQuery := utils.CreateJoin(tx, mainTable, mainAlias, joinTables...)

	joinQuery = joinQuery.Select("MAX(bay.order_number)").Where("company_id = ?", companyId).Find(&highestOrder)

	if joinQuery.Error != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        joinQuery.Error,
		}
	}

	return highestOrder, nil
}

func reset(tx *gorm.DB, companyId int) *exceptions.BaseErrorResponse {
	var bayEntities transactionjpcbentities.BayMaster

	mainTable := "trx_car_wash"
	mainAlias := "carwash"
	mainAliasBay := "bay"

	joinTables := []utils.JoinTable{
		{Table: "mtr_car_wash_bay", Alias: "bay", ForeignKey: mainAlias + ".car_wash_bay_id", ReferenceKey: "bay.car_wash_bay_id"},
	}

	joinQuery := utils.CreateJoin(tx, mainTable, mainAlias, joinTables...)

	keyAttributes := []string{
		mainAlias + ".company_id",
		mainAlias + ".car_wash_id",
		mainAlias + ".car_wash_status_id",
		mainAlias + ".work_order_system_number",
		mainAliasBay + ".car_wash_bay_id",
		mainAliasBay + ".car_wash_bay_code",
		mainAliasBay + ".is_active",
		mainAliasBay + ".car_wash_bay_description",
	}

	joinQuery = joinQuery.Select(keyAttributes).Where("company_id = ?", companyId)

	rows, err := joinQuery.Rows()
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	for rows.Next() {
		var companyId, carWashId, carWashBayId, carWashStatusId, workOrderSystemNumber int
		var carWashBayCode, carWashBayDescription string
		var isActive bool

		err := rows.Scan(&companyId, &carWashId, &carWashStatusId, &workOrderSystemNumber, &carWashBayId, &carWashBayCode, &isActive, &carWashBayDescription)

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
