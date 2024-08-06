package transactionjpcbrepositoryimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	"after-sales/api/utils"
	"net/http"

	"gorm.io/gorm"
)

type BayMasterImpl struct {
}

func OpenBayMasterRepositoryImpl() transactionjpcbrepository.BayMasterRepository {
	return &BayMasterImpl{}
}

func (*BayMasterImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	responses := []transactionjpcbpayloads.BayMasterResponse{}
	// find carwash by company id
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

	//appending data for response
	for rows.Next() {
		var companyId int
		var carWashId int
		var carWashBayId int
		var carWashBayCode string
		var isActive bool
		var carWashBayDescription string
		var carWashStatusId int
		var workOrderSystemNumber int

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

// GetAllActive implements transactionjpcbrepository.BayMasterRepository.
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

	// companyId = filterCondition[]
	joinQuery = joinQuery.Select(keyAttributes)
	filterCondition = append(filterCondition, utils.FilterCondition{ColumnValue: "2", ColumnField: "car_wash_status_id"})

	whereQuery := utils.ApplyFilter(joinQuery, filterCondition)

	rows, err := whereQuery.Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	//appending data for response
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

// GetAllDeactive implements transactionjpcbrepository.BayMasterRepository.
func (*BayMasterImpl) GetAllDeactive(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	panic("unimplemented")
}
