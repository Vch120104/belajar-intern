package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	"errors"
	"fmt"
	"strconv"
	"time"

	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"

	"after-sales/api/utils"
	"net/http"

	"gorm.io/gorm"
)

type WorkOrderAllocationRepositoryImpl struct {
}

func OpenWorkOrderAllocationRepositoryImpl() transactionworkshoprepository.WorkOrderAllocationRepository {
	return &WorkOrderAllocationRepositoryImpl{}
}

// uspg_atWoAllocateGrid_Select
// IF @Option = 0
// --USE FOR : * SELECT DATA FOR WO ALLOCATION GRID
func (r *WorkOrderAllocationRepositoryImpl) GetAll(tx *gorm.DB, companyCode int, foremanId int, date time.Time, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	// Delete all records from WorkOrderAllocationGrid
	if err := tx.Exec("TRUNCATE TABLE trx_work_order_allocation_grid").Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to truncate the table",
			Err:        err,
		}
	}

	// --== INSERT TECHNICIAN & SHIFT CODE FROM atAssignTech ==--
	var assignTechnicians []transactionworkshopentities.AssignTechnician
	if err := tx.Model(&transactionworkshopentities.AssignTechnician{}).
		Select("company_id, foreman_id, technician_id, shift_code, service_date").
		Where("company_id = ? AND foreman_id = ? AND CONVERT(date, service_date) = ?", companyCode, foremanId, date).
		Find(&assignTechnicians).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch technicians",
			Err:        err,
		}
	}

	// Fetch technician names from external service
	technicianNames := make(map[int]string)
	for _, assignTech := range assignTechnicians {
		TechnicianUrl := config.EnvConfigs.GeneralServiceUrl + "user-details/" + strconv.Itoa(assignTech.TechnicianId)
		var getTechnicianResponse masterwarehousepayloads.UserResponse
		if err := utils.Get(TechnicianUrl, &getTechnicianResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch technician data from external service",
				Err:        err,
			}
		}
		technicianNames[assignTech.TechnicianId] = getTechnicianResponse.EmployeeName
	}

	// Insert data into WorkOrderAllocationGrid table
	for _, assignTech := range assignTechnicians {
		workordergrid := transactionworkshopentities.WorkOrderAllocationGrid{
			ShiftCode:      assignTech.ShiftCode,
			TechnicianId:   assignTech.TechnicianId,
			TechnicianName: technicianNames[assignTech.TechnicianId],
		}

		if err := tx.Create(&workordergrid).Error; err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to insert data into WorkOrderAllocationGrid",
				Err:        err,
			}
		}
	}

	// --== UPDATE atWoAllocateGrid ==--
	timeWorkStart := 7.00
	timeWorkInterval := 0.25
	timeWorkEnd := 21.00

	timeColumns := make(map[float64]string)
	for currentTime := timeWorkStart; currentTime < timeWorkEnd; currentTime += timeWorkInterval {
		timeColumn := fmt.Sprintf("time_allocation_%04.0f", currentTime)
		if !isValidTimeColumn(timeColumn) {
			continue
		}
		timeColumns[currentTime] = timeColumn
	}

	for _, assignTech := range assignTechnicians {
		dayOfWeek := date.Weekday()
		var dayColumn string

		switch dayOfWeek {
		case time.Monday:
			dayColumn = "monday"
		case time.Tuesday:
			dayColumn = "tuesday"
		case time.Wednesday:
			dayColumn = "wednesday"
		case time.Thursday:
			dayColumn = "thursday"
		case time.Friday:
			dayColumn = "friday"
		case time.Saturday:
			dayColumn = "saturday"
		case time.Sunday:
			dayColumn = "sunday"
		}

		for currentTime, columnName := range timeColumns {
			var countAvail int64
			if err := tx.Model(&masterentities.ShiftSchedule{}).
				Select("COUNT(*)").
				Where("company_id = ? AND shift_code = ? AND "+dayColumn+" = 1 AND start_time >= ? AND end_time <= ?",
					companyCode, assignTech.ShiftCode, currentTime, currentTime).
				Count(&countAvail).Error; err != nil {
				return nil, 0, 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to fetch availability count",
					Err:        err,
				}
			}

			if countAvail == 0 {
				columnName = "time_allocation_" + fmt.Sprintf("%04.0f", currentTime)

				updateData := map[string]interface{}{
					columnName: -1.0,
				}

				if err := tx.Model(&transactionworkshopentities.WorkOrderAllocationGrid{}).
					Where("shift_code = ? AND technician_id = ?", assignTech.ShiftCode, assignTech.TechnicianId).
					Updates(updateData).Error; err != nil {
					return nil, 0, 0, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to update WorkOrderAllocationGrid",
						Err:        err,
					}
				}
			}
		}
	}

	// Query to select data
	tableStruct := transactionworkshopentities.WorkOrderAllocationGrid{}
	whereQuery := utils.ApplyFilter(tx, filterCondition)

	rows, err := whereQuery.Find(&tableStruct).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
			Err:        err,
		}
	}
	defer rows.Close()

	var convertedResponses []transactionworkshoppayloads.WorkOrderAllocationGridResponse
	for rows.Next() {
		var woResponse transactionworkshopentities.WorkOrderAllocationGrid
		if err := rows.Scan(
			&woResponse.TechnicianId,
			&woResponse.TechnicianName,
			&woResponse.ShiftCode,
			&woResponse.TimeAllocation0700,
			&woResponse.TimeAllocation0715,
			&woResponse.TimeAllocation0730,
			&woResponse.TimeAllocation0745,
			&woResponse.TimeAllocation0800,
			&woResponse.TimeAllocation0815,
			&woResponse.TimeAllocation0830,
			&woResponse.TimeAllocation0845,
			&woResponse.TimeAllocation0900,
			&woResponse.TimeAllocation0915,
			&woResponse.TimeAllocation0930,
			&woResponse.TimeAllocation0945,
			&woResponse.TimeAllocation1000,
			&woResponse.TimeAllocation1015,
			&woResponse.TimeAllocation1030,
			&woResponse.TimeAllocation1045,
			&woResponse.TimeAllocation1100,
			&woResponse.TimeAllocation1115,
			&woResponse.TimeAllocation1130,
			&woResponse.TimeAllocation1145,
			&woResponse.TimeAllocation1200,
			&woResponse.TimeAllocation1215,
			&woResponse.TimeAllocation1230,
			&woResponse.TimeAllocation1245,
			&woResponse.TimeAllocation1300,
			&woResponse.TimeAllocation1315,
			&woResponse.TimeAllocation1330,
			&woResponse.TimeAllocation1345,
			&woResponse.TimeAllocation1400,
			&woResponse.TimeAllocation1415,
			&woResponse.TimeAllocation1430,
			&woResponse.TimeAllocation1445,
			&woResponse.TimeAllocation1500,
			&woResponse.TimeAllocation1515,
			&woResponse.TimeAllocation1530,
			&woResponse.TimeAllocation1545,
			&woResponse.TimeAllocation1600,
			&woResponse.TimeAllocation1615,
			&woResponse.TimeAllocation1630,
			&woResponse.TimeAllocation1645,
			&woResponse.TimeAllocation1700,
			&woResponse.TimeAllocation1715,
			&woResponse.TimeAllocation1730,
			&woResponse.TimeAllocation1745,
			&woResponse.TimeAllocation1800,
			&woResponse.TimeAllocation1815,
			&woResponse.TimeAllocation1830,
			&woResponse.TimeAllocation1845,
			&woResponse.TimeAllocation1900,
			&woResponse.TimeAllocation1915,
			&woResponse.TimeAllocation1930,
			&woResponse.TimeAllocation1945,
			&woResponse.TimeAllocation2000,
			&woResponse.TimeAllocation2015,
			&woResponse.TimeAllocation2030,
			&woResponse.TimeAllocation2045,
			&woResponse.TimeAllocation2100,
		); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to scan rows",
				Err:        err,
			}
		}

		convertedResponses = append(convertedResponses, transactionworkshoppayloads.WorkOrderAllocationGridResponse{
			TechnicianId:   woResponse.TechnicianId,
			TechnicianName: woResponse.TechnicianName,
			ShiftCode:      woResponse.ShiftCode,
			// Include other fields as needed...
		})
	}

	var mapResponses []map[string]interface{}
	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"technician_id":   response.TechnicianId,
			"technician_name": response.TechnicianName,
			"shift_code":      response.ShiftCode,
			"0700":            response.TimeAllocation0700,
			"0715":            response.TimeAllocation0715,
			"0730":            response.TimeAllocation0730,
			"0745":            response.TimeAllocation0745,
			"0800":            response.TimeAllocation0800,
			"0815":            response.TimeAllocation0815,
			"0830":            response.TimeAllocation0830,
			"0845":            response.TimeAllocation0845,
			"0900":            response.TimeAllocation0900,
			"0915":            response.TimeAllocation0915,
			"0930":            response.TimeAllocation0930,
			"0945":            response.TimeAllocation0945,
			"1000":            response.TimeAllocation1000,
			"1015":            response.TimeAllocation1015,
			"1030":            response.TimeAllocation1030,
			"1045":            response.TimeAllocation1045,
			"1100":            response.TimeAllocation1100,
			"1115":            response.TimeAllocation1115,
			"1130":            response.TimeAllocation1130,
			"1145":            response.TimeAllocation1145,
			"1200":            response.TimeAllocation1200,
			"1215":            response.TimeAllocation1215,
			"1230":            response.TimeAllocation1230,
			"1245":            response.TimeAllocation1245,
			"1300":            response.TimeAllocation1300,
			"1315":            response.TimeAllocation1315,
			"1330":            response.TimeAllocation1330,
			"1345":            response.TimeAllocation1345,
			"1400":            response.TimeAllocation1400,
			"1415":            response.TimeAllocation1415,
			"1430":            response.TimeAllocation1430,
			"1445":            response.TimeAllocation1445,
			"1500":            response.TimeAllocation1500,
			"1515":            response.TimeAllocation1515,
			"1530":            response.TimeAllocation1530,
			"1545":            response.TimeAllocation1545,
			"1600":            response.TimeAllocation1600,
			"1615":            response.TimeAllocation1615,
			"1630":            response.TimeAllocation1630,
			"1645":            response.TimeAllocation1645,
			"1700":            response.TimeAllocation1700,
			"1715":            response.TimeAllocation1715,
			"1730":            response.TimeAllocation1730,
			"1745":            response.TimeAllocation1745,
			"1800":            response.TimeAllocation1800,
			"1815":            response.TimeAllocation1815,
			"1830":            response.TimeAllocation1830,
			"1845":            response.TimeAllocation1845,
			"1900":            response.TimeAllocation1900,
			"1915":            response.TimeAllocation1915,
			"1930":            response.TimeAllocation1930,
			"1945":            response.TimeAllocation1945,
			"2000":            response.TimeAllocation2000,
			"2015":            response.TimeAllocation2015,
			"2030":            response.TimeAllocation2030,
			"2045":            response.TimeAllocation2045,
			"2100":            response.TimeAllocation2100,
		}

		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *WorkOrderAllocationRepositoryImpl) GetAllocate(tx *gorm.DB, date time.Time, brandId int, woSysNum int) (transactionworkshoppayloads.WorkOrderAllocationResponse, *exceptions.BaseErrorResponse) {
	var response transactionworkshoppayloads.WorkOrderAllocationResponse

	err := tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Select("company_id, foreman_id, technician_id, shift_code, service_date").
		Where("service_date = ? AND brand_id = ? AND work_order_system_number = ?", date, brandId, woSysNum).
		First(&response).Error

	if err != nil {
		return transactionworkshoppayloads.WorkOrderAllocationResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
			Err:        err,
		}
	}

	return response, nil

}

func (r *WorkOrderAllocationRepositoryImpl) SaveAllocateDetail(tx *gorm.DB, date time.Time, techId int, request transactionworkshoppayloads.WorkOrderAllocationDetailRequest, foremanId int, companyId int) (transactionworkshopentities.WorkOrderAllocationDetail, *exceptions.BaseErrorResponse) {
	// Query AssignTechnician
	var assignTechnicians []transactionworkshopentities.AssignTechnician
	err := tx.Model(&transactionworkshopentities.AssignTechnician{}).
		Where("foreman_id = ? AND service_date = ?", foremanId, date).
		Find(&assignTechnicians).Error
	if err != nil {
		return transactionworkshopentities.WorkOrderAllocationDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Assign Technician not found",
			Err:        err,
		}
	}

	// Iterate over each technician and process accordingly
	for _, tech := range assignTechnicians {
		// Count service logs
		var count int64
		err = tx.Model(&transactionworkshopentities.ServiceLog{}).
			Where("company_id = ? AND technician_id = ? AND shift_schedule_id = ? AND CAST(start_datetime AS DATE) = ?",
				companyId, tech.TechnicianId, tech.ShiftCode, date).
			Where("technician_allocation_system_number IN (?)",
				tx.Model(&transactionworkshopentities.ServiceLog{}).
					Select("technician_allocation_system_number").
					Where("company_id = ? AND technician_id = ? AND shift_schedule_id = ? AND CAST(start_datetime AS DATE) = ?",
						companyId, tech.TechnicianId, tech.ShiftCode, date).
					Group("technician_allocation_system_number")).
			Count(&count).Error
		if err != nil {
			return transactionworkshopentities.WorkOrderAllocationDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to count service logs",
				Err:        err,
			}
		}

		// Check if the technician has already been allocated
		if count > 0 {
			return transactionworkshopentities.WorkOrderAllocationDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Message:    "Technician has already been allocated",
				Err:        errors.New("technician has already been allocated"),
			}
		}

		// If no allocation found, insert into WorkOrderAllocationDetail
		if count == 0 {
			err = tx.Create(&transactionworkshopentities.WorkOrderAllocationDetail{
				TechnicianId:          request.TechnicianId,
				WorkOrderSystemNumber: request.WorkOrderSystemNumber,
				ShiftCode:             request.ShiftCode,
				StartTime:             request.StartTime,
				EndTime:               request.EndTime,
			}).Error
			if err != nil {
				return transactionworkshopentities.WorkOrderAllocationDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to insert new Allocate Detail",
					Err:        err,
				}
			}
		} else {
			// Insert into WorkOrderAllocationDetail with details
			var serviceLogs []transactionworkshopentities.ServiceLog
			err = tx.Model(&transactionworkshopentities.ServiceLog{}).
				Joins("INNER JOIN comGenVariable B ON service_log.service_status_id = B.value AND B.variable LIKE ?", "SRV_STAT_%").
				Where("company_id = ? AND technician_id = ? AND shift_schedule_id = ? AND CAST(start_datetime AS DATE) = ?",
					companyId, tech.TechnicianId, tech.ShiftCode, date).
				Where("technician_allocation_system_number IN (?)",
					tx.Model(&transactionworkshopentities.ServiceLog{}).
						Select("technician_allocation_system_number").
						Where("company_id = ? AND technician_id = ? AND shift_schedule_id = ? AND CAST(start_datetime AS DATE) = ?",
							companyId, tech.TechnicianId, tech.ShiftCode, date).
						Group("technician_allocation_system_number")).
				Order("technician_allocation_system_number, technician_allocation_line").
				Find(&serviceLogs).Error
			if err != nil {
				return transactionworkshopentities.WorkOrderAllocationDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to get service logs",
					Err:        err,
				}
			}

			for _, log := range serviceLogs {
				err = tx.Create(&transactionworkshopentities.WorkOrderAllocationDetail{
					TechnicianId:          request.TechnicianId,
					WorkOrderSystemNumber: request.WorkOrderSystemNumber,
					ShiftCode:             request.ShiftCode,
					StartTime:             request.StartTime,
					EndTime:               request.EndTime,
					ServiceLogId:          log.ServiceLogSystemNumber, // Use the ServiceLogId from the log
				}).Error
				if err != nil {
					return transactionworkshopentities.WorkOrderAllocationDetail{}, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to insert new Allocate Detail with service logs",
						Err:        err,
					}
				}
			}
		}
	}

	// Return success response
	return transactionworkshopentities.WorkOrderAllocationDetail{
		TechnicianId:          request.TechnicianId,
		WorkOrderSystemNumber: request.WorkOrderSystemNumber,
		ShiftCode:             request.ShiftCode,
		StartTime:             request.StartTime,
		EndTime:               request.EndTime,
	}, nil
}

func (r *WorkOrderAllocationRepositoryImpl) GetAllocateDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tableStruct := transactionworkshopentities.WorkOrderAllocationDetail{}

	whereQuery := utils.ApplyFilter(tx, filterCondition)

	rows, err := whereQuery.Find(&tableStruct).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	var convertedResponses []transactionworkshoppayloads.WorkOrderAllocationDetailResponse
	for rows.Next() {
		var woResponse transactionworkshoppayloads.WorkOrderAllocationDetailResponse
		var startTime, endTime time.Time

		if err := rows.Scan(
			&woResponse.TechnicianId,
			&woResponse.TechnicianName,
			&woResponse.WorkOrderSystemNumber,
			&woResponse.WorkOrderDocumentNumber,
			&woResponse.ShiftCode,
			&woResponse.ServiceStatus,
			&startTime,
			&endTime,
		); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		woResponse.StartTime = startTime.Format("15:04:05")
		woResponse.EndTime = endTime.Format("15:04:05")

		convertedResponses = append(convertedResponses, woResponse)
	}

	var mapResponses []map[string]interface{}
	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"technician_id":              response.TechnicianId,
			"technician_name":            response.TechnicianName,
			"work_order_system_number":   response.WorkOrderSystemNumber,
			"work_order_document_number": response.WorkOrderDocumentNumber,
			"shift_code":                 response.ShiftCode,
			"service_status":             response.ServiceStatus,
			"start_time":                 response.StartTime,
			"end_time":                   response.EndTime,
		}

		mapResponses = append(mapResponses, responseMap)

	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *WorkOrderAllocationRepositoryImpl) GetWorkOrderAllocationHeaderData(tx *gorm.DB, companyCode int, foremanId int, techallocStartDate time.Time, vehicleBrandId int) (transactionworkshoppayloads.WorkOrderAllocationHeaderResult, *exceptions.BaseErrorResponse) {
	var result transactionworkshoppayloads.WorkOrderAllocationHeaderResult

	// Get shift start time and end time
	shiftTimes, err := r.getShiftTimes(tx, companyCode, foremanId, techallocStartDate)
	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			Message: "Failed to get shift times",
			Err:     err,
		}
	}

	// Get total time
	totalTime, err := r.getTotalTime(tx, companyCode, shiftTimes.ShiftCode, techallocStartDate, shiftTimes.StartTime, shiftTimes.EndTime)
	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			Message: "Failed to get total time",
			Err:     err,
		}
	}

	// Get used time
	usedTime, err := r.getUsedTime(tx, companyCode, foremanId, techallocStartDate)
	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			Message: "Failed to get used time",
			Err:     err,
		}
	}

	// Calculate available tech time
	availTechTime := totalTime - usedTime

	// Get unallocated operations
	unallocatedOpr, err := r.getUnallocatedOpr(tx, companyCode, techallocStartDate)
	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			Message: "Failed to get unallocated operations",
			Err:     err,
		}
	}

	// Get auto-released operations
	autoReleased, err := r.getAutoReleased(tx, companyCode, techallocStartDate)
	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			Message: "Failed to get auto-released operations",
			Err:     err,
		}
	}

	// Get book allocated time
	bookAllocTime, err := r.getBookAllocTime(tx, companyCode, vehicleBrandId, techallocStartDate)
	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			Message: "Failed to get book allocated time",
			Err:     err,
		}
	}

	result = transactionworkshoppayloads.WorkOrderAllocationHeaderResult{
		TotalTechnicianTime:     totalTime,
		UsedTechnicianTime:      usedTime,
		AvailableTechnicianTime: availTechTime,
		UnallocatedOperation:    unallocatedOpr,
		AutoReleasedOperation:   autoReleased,
		BookAllocatedTime:       bookAllocTime,
	}

	return result, nil
}

func (r *WorkOrderAllocationRepositoryImpl) GetAssignTechnician(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tableStruct := transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianRequest{}
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	rows, err := whereQuery.Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	var convertedResponses []transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianResponse

	for rows.Next() {
		var WoRequest transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianResponse
		var serviceDate time.Time

		if err := rows.Scan(
			&WoRequest.CompanyId,
			&WoRequest.TechnicianId,
			&WoRequest.ShiftCode,
			&WoRequest.ForemanId,
			&serviceDate,
		); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		WoRequest.ServiceDate = serviceDate
		// Calculate attendance based on shift_code and technician_id
		if WoRequest.ShiftCode != "" && WoRequest.TechnicianId != 0 {
			WoRequest.Attendance = true
		} else {
			WoRequest.Attendance = false
		}
		convertedResponses = append(convertedResponses, WoRequest)
	}

	// Create maps to fetch foreman and technician names in bulk
	foremanIds := make(map[int]struct{})
	technicianIds := make(map[int]struct{})

	for _, response := range convertedResponses {
		foremanIds[response.ForemanId] = struct{}{}
		technicianIds[response.TechnicianId] = struct{}{}
	}

	foremanNames := make(map[int]string)
	for foremanId := range foremanIds {
		ForemanUrl := config.EnvConfigs.GeneralServiceUrl + "user-details/" + strconv.Itoa(foremanId)
		var getForemanResponse masterwarehousepayloads.UserResponse
		if err := utils.Get(ForemanUrl, &getForemanResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch foreman data from external service",
				Err:        err,
			}
		}
		foremanNames[foremanId] = getForemanResponse.EmployeeName
	}

	technicianNames := make(map[int]string)
	for technicianId := range technicianIds {
		TechnicianUrl := config.EnvConfigs.GeneralServiceUrl + "user-details/" + strconv.Itoa(technicianId)
		var getTechnicianResponse masterwarehousepayloads.UserResponse
		if err := utils.Get(TechnicianUrl, &getTechnicianResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch technician data from external service",
				Err:        err,
			}
		}
		technicianNames[technicianId] = getTechnicianResponse.EmployeeName
	}

	var mapResponses []map[string]interface{}
	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"service_date":    response.ServiceDate.Format("2006-01-02 15:04:05"),
			"company_id":      response.CompanyId,
			"foreman_id":      response.ForemanId,
			"foreman_name":    foremanNames[response.ForemanId],
			"technician_id":   response.TechnicianId,
			"technician_name": technicianNames[response.TechnicianId],
			"shift_code":      response.ShiftCode,
			"attendance":      response.Attendance,
		}

		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

// uspg_atAssignTech_Insert
// IF @Option = 0
// --USE IN MODUL : AMS-054 /Assign Technician General Repair
func (r *WorkOrderAllocationRepositoryImpl) NewAssignTechnician(tx *gorm.DB, date time.Time, techId int, request transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianRequest) (transactionworkshopentities.AssignTechnician, *exceptions.BaseErrorResponse) {

	var (
		startTime, endTime, restStartTime, restEndTime float64
		err                                            *exceptions.BaseErrorResponse
	)

	cpcCodeDefault := "00002"       //tx.Raw("SELECT dbo.getVariableValue('CPC_CODE')").Scan(&cpcCode)
	refTypeAvailDefault := "ASSIGN" //tx.Raw("SELECT dbo.getVariableValue('REF_TYPE_AVAIL')").Scan(&refTypeAvail)

	// Get start and end times for the shift
	startTime, err = r.getShiftStartTime(tx, request.CompanyId, request.ShiftCode, date, false)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	endTime, err = r.getShiftEndTime(tx, request.CompanyId, request.ShiftCode, date, false)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	// Get rest start and end times for the shift
	restStartTime, err = r.getShiftStartTime(tx, request.CompanyId, request.ShiftCode, date, true)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	restEndTime, err = r.getShiftEndTime(tx, request.CompanyId, request.ShiftCode, date, true)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	var existingAssignTech transactionworkshopentities.AssignTechnician
	if err := tx.Where("foreman_id = ? AND service_date = ? AND technician_id = ? AND company_id = ?",
		request.ForemanId, date, request.TechnicianId, request.CompanyId).First(&existingAssignTech).Error; err == nil {
		return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Data already exists",
			Err:        errors.New("data already exists"),
		}
	}

	var conflictingAssignTech transactionworkshopentities.AssignTechnician
	if err := tx.Where("foreman_id = ? AND service_date = ? AND technician_id <> ? AND shift_code = ? AND company_id = ?",
		request.ForemanId, date, request.TechnicianId, request.ShiftCode, request.CompanyId).First(&conflictingAssignTech).Error; err == nil {
		return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Assign Technician is not valid",
			Err:        errors.New("assign Technician is not valid"),
		}
	}

	entity := transactionworkshopentities.AssignTechnician{
		ServiceDate:  date,
		CompanyId:    request.CompanyId,
		ForemanId:    request.ForemanId,
		TechnicianId: request.TechnicianId,
		ShiftCode:    request.ShiftCode,
		CpcCode:      cpcCodeDefault,
		CreateDate:   time.Now(),
	}

	if err := tx.Create(&entity).Error; err != nil {
		return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to insert new Assign Technician",
			Err:        err,
		}
	}

	// Update wtBookAlloc table
	// Commented out, ensure it's correctly implemented if needed.
	// if err := tx.Model(&transactionworkshopentities.BookingEstimationAllocation{}).
	// 	Where("bookalloc_date = ? AND shift_code = ? AND bookalloc_technician = ?",
	// 		request.ServiceDate, request.ShiftCode, request.TechnicianId).
	// 	Update("assign_technician", gorm.Expr("?", request.TechnicianId)).Error; err != nil {
	// 	return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Message:    "Failed to update wtBookAlloc",
	// 		Err:        err,
	// 	}
	// }

	// uspg_atWoTechAllocAvailable_Insert
	// IF @Option = 2
	// --USE IN MODUL : AMS-054 /Assign Technician General Repair

	// Calculate duration in hours
	durationBeforeRest := (restStartTime - startTime) / 60
	durationAfterRest := (endTime - restEndTime) / 60

	// Insert before rest time if needed
	if restStartTime > startTime {
		if err := tx.Create(&transactionworkshopentities.WorkOrderAllocationAvailable{
			CompanyId:             request.CompanyId,
			ServiceDateTime:       date,
			ForemanId:             request.ForemanId,
			TechnicianId:          request.TechnicianId,
			ShiftCode:             request.ShiftCode,
			StartTime:             startTime,
			EndTime:               restStartTime,
			TotalHour:             durationBeforeRest,
			ReferenceType:         refTypeAvailDefault,
			ReferenceSystemNumber: 0,
			ReferenceLine:         0,
			Remark:                "",
			CreateDate:            time.Now(),
		}).Error; err != nil {
			return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to insert allocation before rest time",
				Err:        err,
			}
		}
	}

	// Insert after rest time if needed
	if restEndTime < endTime {
		if err := tx.Create(&transactionworkshopentities.WorkOrderAllocationAvailable{
			CompanyId:             request.CompanyId,
			ServiceDateTime:       date,
			ForemanId:             request.ForemanId,
			TechnicianId:          request.TechnicianId,
			ShiftCode:             request.ShiftCode,
			StartTime:             restEndTime,
			EndTime:               endTime,
			TotalHour:             durationAfterRest,
			ReferenceType:         refTypeAvailDefault,
			ReferenceSystemNumber: 0,
			ReferenceLine:         0,
			Remark:                "",
			CreateDate:            time.Now(),
		}).Error; err != nil {
			return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to insert allocation after rest time",
				Err:        err,
			}
		}
	}

	return entity, nil
}

func (r *WorkOrderAllocationRepositoryImpl) GetAssignTechnicianById(tx *gorm.DB, date time.Time, techId int, id int) (transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianResponse, *exceptions.BaseErrorResponse) {

	var response transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianResponse

	err := tx.Model(&transactionworkshopentities.AssignTechnician{}).
		Select("company_id, foreman_id, technician_id, shift_code,cpc_code, service_date, CASE WHEN shift_code <> '' AND technician_id <> 0 THEN 1 ELSE 0 END AS attendance").
		Where("assign_technician_id = ? AND foreman_id = ? AND FORMAT(service_date, 'yyyy-MM-dd') = ?", id, techId, date).
		First(&response).Error

	if err != nil {
		return transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
			Err:        err,
		}
	}

	// fetch foreman name
	ForemanUrl := config.EnvConfigs.GeneralServiceUrl + "user-details/" + strconv.Itoa(response.ForemanId)
	var getForemanResponse masterwarehousepayloads.UserResponse
	if err := utils.Get(ForemanUrl, &getForemanResponse, nil); err != nil {
		return transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch foreman data from external service",
			Err:        err,
		}
	}

	// fetch technician name
	TechnicianUrl := config.EnvConfigs.GeneralServiceUrl + "user-details/" + strconv.Itoa(response.TechnicianId)
	var getTechnicianResponse masterwarehousepayloads.UserResponse
	if err := utils.Get(TechnicianUrl, &getTechnicianResponse, nil); err != nil {
		return transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch technician data from external service",
			Err:        err,
		}
	}

	response.ForemanName = getForemanResponse.EmployeeName
	response.TechnicianName = getTechnicianResponse.EmployeeName

	return response, nil
}

// uspg_atAssignTech_Update
// IF @Option = 0
// --USE IN MODUL : AMS-054 /Assign Technician General Repair
func (r *WorkOrderAllocationRepositoryImpl) SaveAssignTechnician(tx *gorm.DB, date time.Time, techId int, id int, request transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianRequest) (transactionworkshopentities.AssignTechnician, *exceptions.BaseErrorResponse) {

	// Declare variables
	var (
		startTime, endTime, restStartTime, restEndTime float64
		shiftCodeOld                                   string
	)

	refTypeAvailDefault := "ASSIGN"

	// Get start and end times for the shift
	startTime, err := r.getShiftStartTime(tx, request.CompanyId, request.ShiftCode, date, false)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	endTime, err = r.getShiftEndTime(tx, request.CompanyId, request.ShiftCode, date, false)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	// Get rest start and end times for the shift
	restStartTime, err = r.getShiftStartTime(tx, request.CompanyId, request.ShiftCode, date, true)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	restEndTime, err = r.getShiftEndTime(tx, request.CompanyId, request.ShiftCode, date, true)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	// Retrieve the old shift code
	if err := tx.Model(&transactionworkshopentities.AssignTechnician{}).
		Select("shift_code").
		Where("foreman_id = ? AND service_date = ? AND technician_id = ?", request.ForemanId, date, request.TechnicianId).
		First(&shiftCodeOld).Error; err != nil {
		return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Old shift code not found",
			Err:        err,
		}
	}

	// Check if there is a conflicting assignment
	var conflictingAssignTech transactionworkshopentities.AssignTechnician
	if err := tx.Where("foreman_id = ? AND service_date = ? AND technician_id <> ? AND shift_code = ? AND company_id = ?",
		request.ForemanId, date, request.TechnicianId, request.ShiftCode, request.CompanyId).First(&conflictingAssignTech).Error; err == nil {
		return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Assign Technician is not valid",
			Err:        errors.New("assign Technician is not valid"),
		}
	}

	// Check if the technician already has an allocation
	var existingTechAlloc transactionworkshopentities.WorkOrderAllocationAvailable
	if err := tx.Where("company_id = ? AND foreman_id = ? AND technician_id = ? AND CONVERT(VARCHAR, tech_alloc_start_date, 106) = CONVERT(VARCHAR, ?, 106)",
		request.CompanyId, request.ForemanId, request.TechnicianId, date).First(&existingTechAlloc).Error; err == nil {
		return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "This Technician already has allocation",
			Err:        errors.New("this Technician already has allocation"),
		}
	}

	// Update the existing record
	var entity transactionworkshopentities.AssignTechnician
	if err := tx.Where("assign_technician_id = ? AND foreman_id = ? AND service_date = ?", id, techId, date).First(&entity).Error; err != nil {
		return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
			Err:        err,
		}
	}

	if err := tx.Model(&entity).Updates(map[string]interface{}{
		"company_id":    request.CompanyId,
		"shift_code":    request.ShiftCode,
		"technician_id": request.TechnicianId,
	}).Error; err != nil {
		return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update Assign Technician",
			Err:        err,
		}
	}

	// Delete old work order tech allocation availability
	if err := tx.Where("company_id = ? AND CONVERT(date, service_date_time) = ? AND technician_id = ? AND shift_code = ?",
		request.CompanyId, date, request.TechnicianId, shiftCodeOld).Delete(&transactionworkshopentities.WorkOrderAllocationAvailable{}).Error; err != nil {
		return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete old work order tech allocation availability",
			Err:        err,
		}
	}

	// Calculate duration in hours
	durationBeforeRest := (restStartTime - startTime) / 60
	durationAfterRest := (endTime - restEndTime) / 60

	// Insert before rest time if needed
	if restStartTime > startTime {
		if err := tx.Create(&transactionworkshopentities.WorkOrderAllocationAvailable{
			CompanyId:             request.CompanyId,
			ServiceDateTime:       date,
			ForemanId:             request.ForemanId,
			TechnicianId:          request.TechnicianId,
			ShiftCode:             request.ShiftCode,
			StartTime:             startTime,
			EndTime:               restStartTime,
			TotalHour:             durationBeforeRest,
			ReferenceType:         refTypeAvailDefault,
			ReferenceSystemNumber: 0,
			ReferenceLine:         0,
			Remark:                "",
			CreateDate:            time.Now(),
		}).Error; err != nil {
			return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to insert allocation before rest time",
				Err:        err,
			}
		}
	}

	// Insert after rest time if needed
	if restEndTime < endTime {
		if err := tx.Create(&transactionworkshopentities.WorkOrderAllocationAvailable{
			CompanyId:             request.CompanyId,
			ServiceDateTime:       date,
			ForemanId:             request.ForemanId,
			TechnicianId:          request.TechnicianId,
			ShiftCode:             request.ShiftCode,
			StartTime:             restEndTime,
			EndTime:               endTime,
			TotalHour:             durationAfterRest,
			ReferenceType:         refTypeAvailDefault,
			ReferenceSystemNumber: 0,
			ReferenceLine:         0,
			Remark:                "",
			CreateDate:            time.Now(),
		}).Error; err != nil {
			return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to insert allocation after rest time",
				Err:        err,
			}
		}
	}

	return entity, nil
}

// //////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////
//
//	Support Function		   							  //
//
// //////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////
func (r *WorkOrderAllocationRepositoryImpl) getShiftStartTime(tx *gorm.DB, companyId int, shiftCode string, effectiveDate time.Time, rest bool) (float64, *exceptions.BaseErrorResponse) {
	var startTime float64

	dayOfWeek := effectiveDate.Weekday()

	var shiftSchedule masterentities.ShiftSchedule
	err := tx.Model(&masterentities.ShiftSchedule{}).
		Select("start_time, rest_start_time").
		Where("company_id = ?", companyId).
		Where("shift_code = ?", shiftCode).
		Where("effective_date <= ?", effectiveDate).
		Where("is_active = ?", true).
		Where(`
			CASE
				WHEN ? = 0 THEN sunday
				WHEN ? = 1 THEN monday
				WHEN ? = 2 THEN tuesday
				WHEN ? = 3 THEN wednesday
				WHEN ? = 4 THEN thursday
				WHEN ? = 5 THEN friday
				WHEN ? = 6 THEN saturday
			END = 1`, dayOfWeek, dayOfWeek, dayOfWeek, dayOfWeek, dayOfWeek, dayOfWeek, dayOfWeek).
		Order("effective_date DESC").
		First(&shiftSchedule).Error

	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if rest {
		startTime = shiftSchedule.RestStartTime
	} else {
		startTime = shiftSchedule.StartTime
	}

	return startTime, nil

}

func (r *WorkOrderAllocationRepositoryImpl) getShiftEndTime(tx *gorm.DB, companyId int, shiftCode string, effectiveDate time.Time, rest bool) (float64, *exceptions.BaseErrorResponse) {
	var endTime float64

	dayOfWeek := effectiveDate.Weekday()

	var shiftSchedule masterentities.ShiftSchedule
	err := tx.Model(&masterentities.ShiftSchedule{}).
		Select("end_time, rest_end_time").
		Where("company_id = ?", companyId).
		Where("shift_code = ?", shiftCode).
		Where("effective_date <= ?", effectiveDate).
		Where("is_active = ?", true).
		Where(`
			CASE
				WHEN ? = 0 THEN sunday
				WHEN ? = 1 THEN monday
				WHEN ? = 2 THEN tuesday
				WHEN ? = 3 THEN wednesday
				WHEN ? = 4 THEN thursday
				WHEN ? = 5 THEN friday
				WHEN ? = 6 THEN saturday
			END = 1`, dayOfWeek, dayOfWeek, dayOfWeek, dayOfWeek, dayOfWeek, dayOfWeek, dayOfWeek).
		Order("effective_date DESC").
		First(&shiftSchedule).Error

	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if rest {
		endTime = shiftSchedule.RestEndTime
	} else {
		endTime = shiftSchedule.EndTime
	}

	return endTime, nil
}

func (r *WorkOrderAllocationRepositoryImpl) getShiftTimes(tx *gorm.DB, companyCode int, foremanId int, techallocStartDate time.Time) (transactionworkshoppayloads.ShiftTimes, error) {
	var shiftTimes transactionworkshoppayloads.ShiftTimes

	// Define day-specific conditions with SQL Server's boolean representation
	dayConditions := map[time.Weekday]string{
		time.Sunday:    "mtr_shift_schedule.sunday = 1",
		time.Monday:    "mtr_shift_schedule.monday = 1",
		time.Tuesday:   "mtr_shift_schedule.tuesday = 1",
		time.Wednesday: "mtr_shift_schedule.wednesday = 1",
		time.Thursday:  "mtr_shift_schedule.thursday = 1",
		time.Friday:    "mtr_shift_schedule.friday = 1",
		time.Saturday:  "mtr_shift_schedule.saturday = 1",
	}

	dayOfWeek := techallocStartDate.Weekday()
	condition, exists := dayConditions[dayOfWeek]
	if !exists {
		return shiftTimes, errors.New("unsupported day of week")
	}

	// Format the date to yyyy-MM-dd
	techallocStartDateStr := techallocStartDate.Format("2006-01-02")

	// Create a struct to hold the result
	var result struct {
		StartTime time.Time
		EndTime   time.Time
		ShiftCode string
	}

	// Build and execute the query
	err := tx.Table("trx_assign_technician").
		Joins("LEFT JOIN mtr_shift_schedule ON "+
			"mtr_shift_schedule.COMPANY_ID = trx_assign_technician.COMPANY_ID AND "+
			"trx_assign_technician.SHIFT_CODE = mtr_shift_schedule.SHIFT_CODE AND "+
			condition).
		Where("trx_assign_technician.SERVICE_DATE = ? AND trx_assign_technician.FOREMAN_ID = ? AND trx_assign_technician.COMPANY_ID = ?", techallocStartDateStr, foremanId, companyCode).
		Select("mtr_shift_schedule.START_TIME, mtr_shift_schedule.END_TIME, trx_assign_technician.SHIFT_CODE").
		Where("mtr_shift_schedule.EFFECTIVE_DATE = (SELECT TOP 1 EFFECTIVE_DATE FROM mtr_shift_schedule WHERE COMPANY_ID = trx_assign_technician.COMPANY_ID AND SHIFT_CODE = trx_assign_technician.SHIFT_CODE AND "+
			condition+" AND EFFECTIVE_DATE <= ? ORDER BY EFFECTIVE_DATE DESC)", techallocStartDateStr).
		Scan(&result).Error

	if err != nil {
		return shiftTimes, err
	}

	// if result.ShiftCode == "" {
	// 	return shiftTimes, errors.New("shift code is empty")
	// }

	// Convert time.Time to float64 decimal hours
	startTimeDecimal := timeToDecimalHours(result.StartTime)
	endTimeDecimal := timeToDecimalHours(result.EndTime)

	shiftTimes = transactionworkshoppayloads.ShiftTimes{
		StartTime: startTimeDecimal,
		EndTime:   endTimeDecimal,
		ShiftCode: result.ShiftCode,
	}

	return shiftTimes, nil
}

// Converts time.Time to float64 representing decimal hours
func timeToDecimalHours(t time.Time) float64 {
	hours := float64(t.Hour())
	minutes := float64(t.Minute())
	return hours + (minutes / 60)
}

func (r *WorkOrderAllocationRepositoryImpl) getTotalTime(tx *gorm.DB, companyCode int, shiftCode string, techallocStartDate time.Time, shiftStartTime float64, shiftEndTime float64) (float64, error) {
	var totalTime float64

	// Determine the day of the week and corresponding condition
	dayOfWeek := techallocStartDate.Weekday()
	dayCondition := ""

	switch dayOfWeek {
	case time.Sunday:
		dayCondition = "sunday = 1"
	case time.Monday:
		dayCondition = "monday = 1"
	case time.Tuesday:
		dayCondition = "tuesday = 1"
	case time.Wednesday:
		dayCondition = "wednesday = 1"
	case time.Thursday:
		dayCondition = "thursday = 1"
	case time.Friday:
		dayCondition = "friday = 1"
	case time.Saturday:
		dayCondition = "saturday = 1"
	default:
		return 0, errors.New("invalid day of week")
	}

	// Fetch the shift schedule from the database
	var shift masterentities.ShiftSchedule
	err := tx.Table("mtr_shift_schedule").
		Where("company_id = ? AND shift_code = ? AND effective_date <= ?", companyCode, shiftCode, techallocStartDate.Format("2006-01-02")).
		Where(dayCondition).
		Order("effective_date DESC").
		Limit(1).
		Find(&shift).Error

	if err != nil {
		return 0, err
	}

	// Handle case where no shift is found
	if shift.ShiftCode == "" {
		return 0, errors.New("shift not found")
	}

	// Convert shift times to float64 if stored as minutes past midnight
	shiftStartTimeInMinutes := float64(shift.StartTime)         // Adjust if needed
	shiftEndTimeInMinutes := float64(shift.EndTime)             // Adjust if needed
	shiftRestStartTimeInMinutes := float64(shift.RestStartTime) // Adjust if needed
	shiftRestEndTimeInMinutes := float64(shift.RestEndTime)     // Adjust if needed

	var durationBeforeRest, durationAfterRest float64

	// Calculate total duration based on shift times and assignment times
	if shiftStartTime < shiftStartTimeInMinutes {
		durationBeforeRest = shiftEndTimeInMinutes - shiftStartTime
	} else {
		durationBeforeRest = shiftEndTimeInMinutes - shiftStartTimeInMinutes
	}

	if shiftEndTime > shiftEndTimeInMinutes {
		durationAfterRest = shiftEndTime - shiftEndTimeInMinutes
	} else {
		durationAfterRest = shiftEndTime - shiftStartTimeInMinutes
	}

	totalTime = durationBeforeRest + durationAfterRest - (shiftRestEndTimeInMinutes - shiftRestStartTimeInMinutes)

	return totalTime, nil
}

func (r *WorkOrderAllocationRepositoryImpl) getUsedTime(tx *gorm.DB, companyCode int, foremanId int, techallocStartDate time.Time) (float64, error) {
	var usedTime float64

	startDateStr := techallocStartDate.Format("2006-01-02")

	err := tx.Model(&transactionworkshopentities.ServiceLog{}).
		Select("SUM(actual_time) AS total_time").
		Joins("LEFT JOIN ? AS WTA ON WTA.tech_alloc_system_number = S.technician_allocation_system_number", &transactionworkshopentities.WorkOrderAllocation{}).
		Where("S.company_id = ? AND WTA.foreman_id = ? AND DATE(S.start_datetime) = ?", companyCode, foremanId, startDateStr).
		Scan(&usedTime).Error

	if err != nil {
		return 0, err
	}

	return usedTime, nil
}

func (r *WorkOrderAllocationRepositoryImpl) getUnallocatedOpr(tx *gorm.DB, companyCode int, techallocStartDate time.Time) (int, error) {
	var unallocatedOpr int64

	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Joins("LEFT JOIN ? ON ? = ?",
			&transactionworkshopentities.WorkOrder{},
			"trx_work_order_detail.work_order_system_number",
			"trx_work_order.work_order_system_number").
		Where("trx_work_order.company_id = ?", companyCode).
		Where("trx_work_order.work_order_status_id = ?", 0). // Assuming empty status is represented as 0, adjust as necessary
		Where("trx_work_order.work_order_date = ?", techallocStartDate).
		Count(&unallocatedOpr).Error

	if err != nil {
		return 0, err
	}

	return int(unallocatedOpr), nil
}

func (r *WorkOrderAllocationRepositoryImpl) getAutoReleased(tx *gorm.DB, companyCode int, techallocStartDate time.Time) (int, error) {
	var autoReleased int64

	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Joins("LEFT JOIN ? ON ? = ?",
			&transactionworkshopentities.WorkOrder{},
			"trx_work_order_detail.work_order_system_number",
			"trx_work_order.work_order_system_number").
		Where("trx_work_order.company_id = ?", companyCode).
		Where("COALESCE(trx_work_order.work_order_status, '') = ?", "SRV_STAT_AUTORELEASE").
		Where("trx_work_order.work_order_date = ?", techallocStartDate).
		Count(&autoReleased). // Use the integer pointer for Count
		Error

	if err != nil {
		return 0, err
	}

	return int(autoReleased), nil
}

func (r *WorkOrderAllocationRepositoryImpl) getBookAllocTime(tx *gorm.DB, companyCode int, vehicleBrandId int, techallocStartDate time.Time) (float64, error) {
	var bookAllocTime float64

	err := tx.Model(&transactionworkshopentities.BookingEstimationAllocation{}).
		Select("SUM(bookalloc_total_hour)").
		Where("company_code = ? AND vehicle_brand = ? AND bookalloc_date = ?", companyCode, vehicleBrandId, techallocStartDate).
		Scan(&bookAllocTime).Error

	if err != nil {
		return 0, err
	}

	return bookAllocTime, nil
}

func isValidTimeColumn(columnName string) bool {
	validColumns := []string{
		"time_allocation_0700", "time_allocation_0715", "time_allocation_0730", "time_allocation_0745",
		"time_allocation_0800", "time_allocation_0815", "time_allocation_0830", "time_allocation_0845",
		"time_allocation_0900", "time_allocation_0915", "time_allocation_0930", "time_allocation_0945",
		"time_allocation_1000", "time_allocation_1015", "time_allocation_1030", "time_allocation_1045",
		"time_allocation_1100", "time_allocation_1115", "time_allocation_1130", "time_allocation_1145",
		"time_allocation_1200", "time_allocation_1215", "time_allocation_1230", "time_allocation_1245",
		"time_allocation_1300", "time_allocation_1315", "time_allocation_1330", "time_allocation_1345",
		"time_allocation_1400", "time_allocation_1415", "time_allocation_1430", "time_allocation_1445",
		"time_allocation_1500", "time_allocation_1515", "time_allocation_1530", "time_allocation_1545",
		"time_allocation_1600", "time_allocation_1615", "time_allocation_1630", "time_allocation_1645",
		"time_allocation_1700", "time_allocation_1715", "time_allocation_1730", "time_allocation_1745",
		"time_allocation_1800", "time_allocation_1815", "time_allocation_1830", "time_allocation_1845",
		"time_allocation_1900", "time_allocation_1915", "time_allocation_1930", "time_allocation_1945",
		"time_allocation_2000", "time_allocation_2015", "time_allocation_2030", "time_allocation_2045",
		"time_allocation_2100", "time_allocation_2115", "time_allocation_2130", "time_allocation_2145",
		"time_allocation_2200", "time_allocation_2215", "time_allocation_2230", "time_allocation_2245",
		"time_allocation_2300", "time_allocation_2315", "time_allocation_2330", "time_allocation_2345",
	}
	for _, col := range validColumns {
		if columnName == col {
			return true
		}
	}
	return false
}
