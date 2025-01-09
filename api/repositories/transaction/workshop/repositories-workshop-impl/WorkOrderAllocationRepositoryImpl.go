package transactionworkshoprepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"fmt"
	"strconv"
	"time"

	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"

	utils "after-sales/api/utils"
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
func (r *WorkOrderAllocationRepositoryImpl) GetAll(tx *gorm.DB, companyId int, foremanId int, date time.Time, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var pages pagination.Pagination
	pages.Rows = []map[string]interface{}{}

	var assignTechnicians []transactionworkshopentities.AssignTechnician
	if err := tx.Model(&transactionworkshopentities.AssignTechnician{}).
		Select("technician_id, shift_code").
		Where("company_id = ? AND foreman_id = ? AND CONVERT(date, service_date) = ?", companyId, foremanId, date).
		Find(&assignTechnicians).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch technicians",
			Err:        err,
		}
	}

	var technicianIds []int
	for _, assignTech := range assignTechnicians {
		technicianIds = append(technicianIds, assignTech.TechnicianId)
	}

	if err := tx.Where("technician_id IN (?)", technicianIds).Delete(&transactionworkshopentities.WorkOrderAllocationGrid{}).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete data from WorkOrderAllocationGrid",
			Err:        err,
		}
	}

	technicianNames := make(map[int]string)
	for _, assignTech := range assignTechnicians {
		technicianResponse, technicianErr := generalserviceapiutils.GetEmployeeById(assignTech.TechnicianId)
		if technicianErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch technician data from external service",
				Err:        technicianErr,
			}
		}
		technicianNames[assignTech.TechnicianId] = technicianResponse.EmployeeName
	}

	for _, assignTech := range assignTechnicians {
		workordergrid := transactionworkshopentities.WorkOrderAllocationGrid{
			ShiftCode:      assignTech.ShiftCode,
			TechnicianId:   assignTech.TechnicianId,
			TechnicianName: technicianNames[assignTech.TechnicianId],
		}

		if err := tx.Create(&workordergrid).Error; err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to insert data into WorkOrderAllocationGrid",
				Err:        err,
			}
		}
	}

	var shiftSchedule masterentities.ShiftSchedule
	if err := tx.Model(&masterentities.ShiftSchedule{}).
		Select("start_time, end_time, rest_start_time, rest_end_time").
		Where("company_id = ? AND shift_code = ?", companyId, assignTechnicians[0].ShiftCode).
		First(&shiftSchedule).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch shift schedule",
			Err:        err,
		}
	}

	startTimeStr := utils.Float64ToTimeString(shiftSchedule.StartTime)
	endTimeStr := utils.Float64ToTimeString(shiftSchedule.EndTime)
	restStartTimeStr := utils.Float64ToTimeString(shiftSchedule.RestStartTime)
	restEndTimeStr := utils.Float64ToTimeString(shiftSchedule.RestEndTime)

	startTimeFloat, _ := strconv.ParseFloat(startTimeStr, 64)
	endTimeFloat, _ := strconv.ParseFloat(endTimeStr, 64)
	restStartTimeFloat, _ := strconv.ParseFloat(restStartTimeStr, 64)
	restEndTimeFloat, _ := strconv.ParseFloat(restEndTimeStr, 64)

	fmt.Println("Start Time: ", startTimeFloat)
	fmt.Println("End Time: ", endTimeFloat)
	fmt.Println("Rest Start Time: ", restStartTimeFloat)
	fmt.Println("Rest End Time: ", restEndTimeFloat)

	if startTimeFloat < 700 {
		startTimeFloat = 700
	}
	if endTimeFloat > 2100 {
		endTimeFloat = 2100
	}

	// Allocate time slots
	timeWorkInterval := 15    // 15 minutes
	defaultStartTime := 700.0 // 07:00 in float format
	defaultEndTime := 2100.0  // 21:00 in float format

	for currentTime := defaultStartTime; currentTime <= defaultEndTime; currentTime += float64(timeWorkInterval) / 60.0 {

		hours := int(currentTime) / 100
		minutes := int(currentTime) % 100
		timeColumn := fmt.Sprintf("time_allocation_%02d%02d", hours, minutes)

		if !r.isValidTimeColumn(timeColumn) {
			continue
		}

		for _, assignTech := range assignTechnicians {
			var allocate int64

			countCheckAvail, err := r.CountAvailableShifts(tx, companyId, assignTech.ShiftCode, date)
			if err != nil {
				fmt.Printf("Error checking availability: %v\n", err)
				allocate = -1 // Default to not available on error
			} else if countCheckAvail == 0 {
				allocate = -1 // Not available if no effective date or day is not active
			} else {
				if currentTime <= startTimeFloat {
					allocate = -1 // Before working hours
				} else if currentTime >= restStartTimeFloat && currentTime <= restEndTimeFloat {
					allocate = -2 // Rest period
				} else if currentTime >= startTimeFloat && currentTime <= endTimeFloat {

					var exists int64
					err = tx.Model(&transactionworkshopentities.ServiceLog{}).
						Where("company_id = ? AND technician_id = ? AND shift_code = ? AND FORMAT(start_datetime, 'dd-MMM-yyyy') = ? AND FORMAT(start_datetime, 'HH:mm:ss') <= ? AND ? < FORMAT(end_datetime, 'HH:mm:ss')",
							companyId, assignTech.TechnicianId, assignTech.ShiftCode, date, currentTime, currentTime).
						Count(&exists).Error

					if err != nil {
						return pages, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to check if service log exists",
							Err:        err,
						}
					}

					if exists > 0 {
						var serviceLogs []transactionworkshopentities.ServiceLog
						err := tx.Model(&transactionworkshopentities.ServiceLog{}).
							Where("company_id = ? AND technician_id = ? AND shift_code = ? AND DATE_FORMAT(start_datetime, '%d-%b-%Y') = ? AND DATE_FORMAT(start_datetime, '%H:%i:%s') <= ? AND ? < DATE_FORMAT(end_datetime, '%H:%i:%s')",
								companyId, assignTech.TechnicianId, assignTech.ShiftCode, date, currentTime, currentTime).
							Find(&serviceLogs).Error

						if err != nil {
							return pages, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to retrieve service logs",
								Err:        err,
							}
						}

						var maxTechAllocSysNo int
						for _, log := range serviceLogs {
							if log.WorkOrderSystemNumber > maxTechAllocSysNo {
								maxTechAllocSysNo = log.WorkOrderSystemNumber
							}
						}

						var WorkOrderSystemNumber, servStatusId int

						// Pending check
						for _, log := range serviceLogs {
							if log.WorkOrderSystemNumber == maxTechAllocSysNo && log.ServiceStatusId == utils.SrvStatPending {
								WorkOrderSystemNumber = log.WorkOrderSystemNumber
								servStatusId = log.ServiceStatusId
								break
							}
						}

						// Cek status Stop or Transfer if no Pending
						if WorkOrderSystemNumber == 0 {
							for _, log := range serviceLogs {
								if log.WorkOrderSystemNumber == maxTechAllocSysNo &&
									(log.ServiceStatusId == utils.SrvStatStop || log.ServiceStatusId == utils.SrvStatTransfer) {
									WorkOrderSystemNumber = log.WorkOrderSystemNumber
									servStatusId = log.ServiceStatusId
									break
								}
							}
						}

						if WorkOrderSystemNumber == 0 {
							for _, log := range serviceLogs {
								if log.WorkOrderSystemNumber == maxTechAllocSysNo {
									WorkOrderSystemNumber = log.WorkOrderSystemNumber
									servStatusId = log.ServiceStatusId
									break
								}
							}
						}

						if WorkOrderSystemNumber != 0 {
							switch servStatusId {
							case utils.SrvStatStart:
								allocate = int64(WorkOrderSystemNumber) // 1
							case utils.SrvStatDraft:
								allocate = int64(WorkOrderSystemNumber) //3
							case utils.SrvStatStop, utils.SrvStatQcPass:
								allocate = int64(WorkOrderSystemNumber) // 4
							case utils.SrvStatTransfer:
								allocate = int64(WorkOrderSystemNumber) // 5
							default:
								allocate = 0
							}
						}

					}
				} else {
					allocate = -1 // After working hours
				}
			}

			updateData := map[string]interface{}{timeColumn: allocate}
			if err := tx.Model(&transactionworkshopentities.WorkOrderAllocationGrid{}).
				Where("technician_id = ? AND shift_code = ?", assignTech.TechnicianId, assignTech.ShiftCode).
				Updates(updateData).Error; err != nil {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update WorkOrderAllocationGrid",
					Err:        err,
				}
			}
		}
	}

	var rows []map[string]interface{}
	query := tx.Model(&transactionworkshopentities.WorkOrderAllocationGrid{})
	query = utils.ApplyFilter(query, filterCondition)

	if err := query.Find(&rows).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch data from WorkOrderAllocationGrid",
			Err:        err,
		}
	}

	if len(rows) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No data found in WorkOrderAllocationGrid",
			Err:        errors.New("no data found in WorkOrderAllocationGrid"),
		}
	}

	pages.Rows = rows
	pages.TotalRows = int64(len(rows))

	return pages, nil
}

func (r *WorkOrderAllocationRepositoryImpl) GetAllocate(tx *gorm.DB, companyId int, date time.Time, foremanId int, brandId int, workOrderSystemNumber int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var results []transactionworkshoppayloads.WorkOrderAllocationResponse

	pages.Rows = results

	return pages, nil
}

func (r *WorkOrderAllocationRepositoryImpl) WorkOrderAllocationGR(tx *gorm.DB, companyId int, date time.Time, foremanId int, brandId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var entities []transactionworkshopentities.WorkOrder
	var responses []transactionworkshoppayloads.WorkOrderAllocationGR

	baseModelQuery := tx.Model(&entities).
		Select(`
			company_id,
			work_order_system_number,
			work_order_document_number,
			work_order_date,
			model_id,
			variant_id,
			vehicle_id,
			service_advisor_id
		`).
		Where("company_id = ? AND foreman_id = ? AND CAST(work_order_date AS DATE) = ? AND brand_id = ?", companyId, foremanId, date.Format("2006-01-02"), brandId)

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Get model, variant, vehicle, and service advisor data
	for idx, response := range responses {
		modelResponse, modelErr := salesserviceapiutils.GetUnitModelById(response.ModelId)
		if modelErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch model data from external service",
				Err:        modelErr,
			}
		}

		variantResponse, variantErr := salesserviceapiutils.GetUnitVariantById(response.VariantId)
		if variantErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch variant data from external service",
				Err:        variantErr,
			}
		}

		vehicleResponse, vehicleErr := salesserviceapiutils.GetVehicleById(response.VehicleId)
		if vehicleErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch vehicle data from external service",
				Err:        vehicleErr,
			}
		}

		serviceAdvisorResponse, serviceAdvisorErr := generalserviceapiutils.GetEmployeeById(response.ServiceAdvisorId)
		if serviceAdvisorErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch service advisor data from external service",
				Err:        serviceAdvisorErr,
			}
		}

		responses[idx].ModelDescription = modelResponse.ModelName
		responses[idx].VariantDescription = variantResponse.VariantDescription
		responses[idx].VehicleChassisNumber = vehicleResponse.Data.Master.VehicleChassisNumber
		//responses[idx].VehicleTnkb = vehicleResponse.Data.STNK.VehicleRegistrationCertificateTNKB
		responses[idx].ServiceAdvisorName = serviceAdvisorResponse.EmployeeName
	}

	pages.Rows = responses

	return pages, nil
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

	for _, tech := range assignTechnicians {
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
			var serviceLogs []transactionworkshopentities.ServiceLog
			err = tx.Model(&transactionworkshopentities.ServiceLog{}).
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
					ServiceLogId:          log.ServiceLogSystemNumber,
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

	return transactionworkshopentities.WorkOrderAllocationDetail{
		TechnicianId:          request.TechnicianId,
		WorkOrderSystemNumber: request.WorkOrderSystemNumber,
		ShiftCode:             request.ShiftCode,
		StartTime:             request.StartTime,
		EndTime:               request.EndTime,
	}, nil
}

func (r *WorkOrderAllocationRepositoryImpl) GetAllocateDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []transactionworkshopentities.WorkOrderAllocationDetail

	baseModelQuery := tx.Model(&transactionworkshopentities.WorkOrderAllocationDetail{})
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, response := range responses {
		startTime := response.StartTime.Format("15:04:05")
		endTime := response.EndTime.Format("15:04:05")

		result := map[string]interface{}{
			"technician_id":              response.TechnicianId,
			"technician_name":            response.TechnicianName,
			"work_order_system_number":   response.WorkOrderSystemNumber,
			"work_order_document_number": response.WorkOrderDocumentNumber,
			"shift_code":                 response.ShiftCode,
			"service_status":             response.ServiceStatus,
			"start_time":                 startTime,
			"end_time":                   endTime,
		}

		results = append(results, result)
	}

	pages.Rows = results

	return pages, nil
}

// uspg_atWoTechAlloc_Select
// IF @Option = 2
// --USE FOR : * SELECT DATA
// --USE IN MODUL : AWS-004 PAGE 5 REQ: DANIEL 130109
func (r *WorkOrderAllocationRepositoryImpl) GetWorkOrderAllocationHeaderData(tx *gorm.DB, companyId int, foremanId int, techallocStartDate time.Time) (transactionworkshoppayloads.WorkOrderAllocationHeaderResult, *exceptions.BaseErrorResponse) {
	var result transactionworkshoppayloads.WorkOrderAllocationHeaderResult

	// Get shift start time and end time
	shiftTimes, err := r.getShiftTimes(tx, companyId, foremanId, techallocStartDate)
	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get shift times",
			Err:        err,
		}
	}

	// Get total time
	totalTime, err := r.getTotalTime(tx, companyId, shiftTimes.ShiftCode, techallocStartDate, shiftTimes.StartTime, shiftTimes.EndTime)
	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get total time",
			Err:        err,
		}
	}

	// Get used time
	usedTime, err := r.getUsedTime(tx, companyId, foremanId, techallocStartDate)
	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get used time",
			Err:        err,
		}
	}

	// Calculate available tech time
	availTechTime := totalTime - usedTime

	// Get unallocated operations
	unallocatedOpr, err := r.getUnallocatedOpr(tx, companyId, techallocStartDate)
	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get unallocated operations",
			Err:        err,
		}
	}

	// Get auto-released operations
	autoReleased, err := r.getAutoReleased(tx, companyId, techallocStartDate)
	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get auto-released operations",
			Err:        err,
		}
	}

	// Get book allocated time
	bookAllocTime, err := r.getBookAllocTime(tx, companyId, techallocStartDate)
	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get book allocated time",
			Err:        err,
		}
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	result = transactionworkshoppayloads.WorkOrderAllocationHeaderResult{
		TotalTechnicianTime:     totalTime,
		CurrentTime:             currentTime,
		UsedTechnicianTime:      usedTime,
		AvailableTechnicianTime: availTechTime,
		UnallocatedOperation:    unallocatedOpr,
		AutoReleasedOperation:   autoReleased,
		BookAllocatedTime:       bookAllocTime,
	}

	return result, nil
}

func (r *WorkOrderAllocationRepositoryImpl) GetAssignTechnician(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianResponse

	baseModelQuery := tx.Model(&transactionworkshopentities.AssignTechnician{}).
		Select(`
			assign_technician_id,
			company_id, 
			technician_id, 
			shift_code, 
			foreman_id, 
			service_date, 
			technician_no, 
			CASE WHEN shift_code <> '' AND technician_id <> 0 THEN 1 ELSE 0 END AS attendance
		`)

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve assign technician data from database",
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianResponse{}
		return pages, nil
	}

	for idx, response := range responses {
		foremanResponse, foremanErr := generalserviceapiutils.GetEmployeeById(response.ForemanId)
		if foremanErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch foreman data from external service",
				Err:        foremanErr,
			}
		}

		technicianResponse, technicianErr := generalserviceapiutils.GetEmployeeById(response.TechnicianId)
		if technicianErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch technician data from external service",
				Err:        technicianErr,
			}
		}

		companyResponse, companyErr := generalserviceapiutils.GetCompanyDataById(response.CompanyId)
		if companyErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch company data from external service",
				Err:        companyErr,
			}
		}

		responses[idx].ForemanName = foremanResponse.EmployeeName
		responses[idx].TechnicianName = technicianResponse.EmployeeName
		responses[idx].CompanyName = companyResponse.CompanyName
	}

	pages.Rows = responses
	return pages, nil
}

// uspg_atAssignTech_Insert
// IF @Option = 0
// --USE IN MODUL : AMS-054 /Assign Technician General Repair
func (r *WorkOrderAllocationRepositoryImpl) NewAssignTechnician(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianRequest) (transactionworkshopentities.AssignTechnician, *exceptions.BaseErrorResponse) {

	var (
		startTime, endTime, restStartTime, restEndTime float64
		err                                            *exceptions.BaseErrorResponse
	)

	cpcCodeDefault := "00002"
	refTypeAvailDefault := "ASSIGN"

	startTime, err = r.getShiftStartTime(tx, request.CompanyId, request.ShiftCode, request.ServiceDate, false)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	endTime, err = r.getShiftEndTime(tx, request.CompanyId, request.ShiftCode, request.ServiceDate, false)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	restStartTime, err = r.getShiftStartTime(tx, request.CompanyId, request.ShiftCode, request.ServiceDate, true)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	restEndTime, err = r.getShiftEndTime(tx, request.CompanyId, request.ShiftCode, request.ServiceDate, true)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	var existingAssignTech transactionworkshopentities.AssignTechnician
	if err := tx.Where("foreman_id = ? AND service_date = ? AND technician_id = ? AND company_id = ?",
		request.ForemanId, request.ServiceDate, request.TechnicianId, request.CompanyId).First(&existingAssignTech).Error; err == nil {
		return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Data Technician already exists",
			Err:        errors.New("data Technician already exists"),
		}
	}

	var conflictingAssignTech transactionworkshopentities.AssignTechnician
	if err := tx.Where("foreman_id = ? AND service_date = ? AND technician_id <> ? AND shift_code = ? AND company_id = ?",
		request.ForemanId, request.ServiceDate, request.TechnicianId, request.ShiftCode, request.CompanyId).First(&conflictingAssignTech).Error; err == nil {
		return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Assign Technician is not valid",
			Err:        errors.New("assign Technician is not valid"),
		}
	}

	entity := transactionworkshopentities.AssignTechnician{
		ServiceDate:  request.ServiceDate,
		CompanyId:    request.CompanyId,
		ForemanId:    request.ForemanId,
		TechnicianId: request.TechnicianId,
		TechnicianNo: request.TechnicianNo,
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
	if err := tx.Model(&transactionworkshopentities.BookingAllocation{}).
		Where("booking_allocation_date = ? AND shift_code = ? AND booking_allocation_technician = ?",
			request.ServiceDate, request.ShiftCode, request.TechnicianId).
		Update("technician_id", gorm.Expr("?", request.TechnicianId)).Error; err != nil {
		return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update Booking Allocation",
			Err:        err,
		}
	}

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
			ServiceDateTime:       request.ServiceDate,
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
			ServiceDateTime:       request.ServiceDate,
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
		Select("assign_technician_id, company_id, foreman_id, technician_id, shift_code, cpc_code, service_date, CASE WHEN shift_code <> '' AND technician_id <> 0 THEN 1 ELSE 0 END AS attendance").
		Where("assign_technician_id = ? AND foreman_id = ? AND service_date = ?", id, techId, date).
		First(&response).Error

	if err != nil {
		return transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
			Err:        err,
		}
	}

	foremanResponse, foremanErr := generalserviceapiutils.GetEmployeeById(response.ForemanId)
	if foremanErr != nil {
		return transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch foreman data from external service",
			Err:        foremanErr,
		}
	}

	technicianResponse, technicianErr := generalserviceapiutils.GetEmployeeById(response.TechnicianId)
	if technicianErr != nil {
		return transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch technician data from external service",
			Err:        technicianErr,
		}
	}

	companyResponse, companyErr := generalserviceapiutils.GetCompanyDataById(response.CompanyId)
	if companyErr != nil {
		return transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch company data from external service",
			Err:        companyErr,
		}
	}

	response.ForemanName = foremanResponse.EmployeeName
	response.TechnicianName = technicianResponse.EmployeeName
	response.CompanyName = companyResponse.CompanyName

	return response, nil
}

// uspg_atAssignTech_Update
// IF @Option = 0
// --USE IN MODUL : AMS-054 /Assign Technician General Repair
func (r *WorkOrderAllocationRepositoryImpl) SaveAssignTechnician(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianRequest) (transactionworkshopentities.AssignTechnician, *exceptions.BaseErrorResponse) {

	// Declare variables
	var (
		startTime, endTime, restStartTime, restEndTime float64
		shiftCodeOld                                   string
	)

	refTypeAvailDefault := "ASSIGN"

	// Get start and end times for the shift
	startTime, err := r.getShiftStartTime(tx, request.CompanyId, request.ShiftCode, request.ServiceDate, false)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	endTime, err = r.getShiftEndTime(tx, request.CompanyId, request.ShiftCode, request.ServiceDate, false)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	// Get rest start and end times for the shift
	restStartTime, err = r.getShiftStartTime(tx, request.CompanyId, request.ShiftCode, request.ServiceDate, true)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	restEndTime, err = r.getShiftEndTime(tx, request.CompanyId, request.ShiftCode, request.ServiceDate, true)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	if err := tx.Model(&transactionworkshopentities.AssignTechnician{}).
		Select("shift_code").
		Where("foreman_id = ? AND service_date = ? AND technician_id = ?", request.ForemanId, request.ServiceDate, request.TechnicianId).
		First(&shiftCodeOld).Error; err != nil {
		return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Old shift code not found",
			Err:        err,
		}
	}

	var conflictingAssignTech transactionworkshopentities.AssignTechnician
	if err := tx.Where("foreman_id = ? AND service_date = ? AND technician_id <> ? AND shift_code = ? AND company_id = ?",
		request.ForemanId, request.ServiceDate, request.TechnicianId, request.ShiftCode, request.CompanyId).First(&conflictingAssignTech).Error; err == nil {
		return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Assign Technician is not valid",
			Err:        errors.New("assign Technician is not valid"),
		}
	}

	var existingTechAlloc transactionworkshopentities.WorkOrderAllocationAvailable
	if err := tx.Where("company_id = ? AND foreman_id = ? AND technician_id = ? AND CONVERT(VARCHAR, tech_alloc_start_date, 106) = CONVERT(VARCHAR, ?, 106)",
		request.CompanyId, request.ForemanId, request.TechnicianId, request.ServiceDate).First(&existingTechAlloc).Error; err == nil {
		return transactionworkshopentities.AssignTechnician{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "This Technician already has allocation",
			Err:        errors.New("this Technician already has allocation"),
		}
	}

	var entity transactionworkshopentities.AssignTechnician
	if err := tx.Where("assign_technician_id = ? AND foreman_id = ? AND service_date = ?", id, request.ForemanId, request.ServiceDate).First(&entity).Error; err != nil {
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

	if err := tx.Where("company_id = ? AND CONVERT(date, service_date_time) = ? AND technician_id = ? AND shift_code = ?",
		request.CompanyId, request.ServiceDate, request.TechnicianId, shiftCodeOld).Delete(&transactionworkshopentities.WorkOrderAllocationAvailable{}).Error; err != nil {
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
			ServiceDateTime:       request.ServiceDate,
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
			ServiceDateTime:       request.ServiceDate,
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
			Message:    "Failed to retrieve shift start time from master data shift schedule",
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
			Message:    "Failed to retrieve shift end time from master data shift schedule",
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

func (r *WorkOrderAllocationRepositoryImpl) getShiftTimes(tx *gorm.DB, companyId int, foremanId int, techallocStartDate time.Time) (transactionworkshoppayloads.ShiftTimes, error) {
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
			"mtr_shift_schedule.company_id = trx_assign_technician.company_id AND "+
			"trx_assign_technician.shift_code = mtr_shift_schedule.shift_code AND "+
			condition).
		Where("trx_assign_technician.SERVICE_DATE = ? AND trx_assign_technician.FOREMAN_ID = ? AND trx_assign_technician.company_id = ?", techallocStartDateStr, foremanId, companyId).
		Select("mtr_shift_schedule.START_TIME, mtr_shift_schedule.END_TIME, trx_assign_technician.shift_code").
		Where("mtr_shift_schedule.EFFECTIVE_DATE = (SELECT TOP 1 EFFECTIVE_DATE FROM mtr_shift_schedule WHERE company_id = trx_assign_technician.company_id AND shift_code = trx_assign_technician.shift_code AND "+
			condition+" AND EFFECTIVE_DATE <= ? ORDER BY EFFECTIVE_DATE DESC)", techallocStartDateStr).
		Scan(&result).Error

	if err != nil {
		return shiftTimes, err
	}

	if result.ShiftCode == "" {
		return shiftTimes, errors.New("shift code is empty")
	}

	// Convert time.Time to float64 decimal hours
	startTimeDecimal := utils.TimeToDecimalHours(result.StartTime)
	endTimeDecimal := utils.TimeToDecimalHours(result.EndTime)

	shiftTimes = transactionworkshoppayloads.ShiftTimes{
		StartTime: startTimeDecimal,
		EndTime:   endTimeDecimal,
		ShiftCode: result.ShiftCode,
	}

	return shiftTimes, nil
}

func (r *WorkOrderAllocationRepositoryImpl) getTotalTime(tx *gorm.DB, companyId int, shiftCode string, techallocStartDate time.Time, shiftStartTime float64, shiftEndTime float64) (float64, error) {
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
		Where("company_id = ? AND shift_code = ? AND effective_date <= ?", companyId, shiftCode, techallocStartDate.Format("2006-01-02")).
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
	shiftStartTimeInMinutes := float64(shift.StartTime)
	shiftEndTimeInMinutes := float64(shift.EndTime)
	shiftRestStartTimeInMinutes := float64(shift.RestStartTime)
	shiftRestEndTimeInMinutes := float64(shift.RestEndTime)

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

func (r *WorkOrderAllocationRepositoryImpl) getUsedTime(tx *gorm.DB, companyId int, foremanId int, techallocStartDate time.Time) (float64, error) {
	var usedTime float64

	startDateStr := techallocStartDate.Format("2006-01-02")

	err := tx.Table("trx_service_log AS S").
		Select("COALESCE(SUM(S.actual_time), 0) AS total_time").
		Joins("LEFT JOIN trx_work_order_allocation AS WTA ON WTA.technician_allocation_system_number = S.technician_allocation_system_number").
		Where("S.company_id = ? AND WTA.foreman_id = ? AND CAST(S.start_datetime AS DATE) = ?", companyId, foremanId, startDateStr).
		Scan(&usedTime).Error

	if err != nil {
		return 0, err
	}

	return usedTime, nil
}

func (r *WorkOrderAllocationRepositoryImpl) getUnallocatedOpr(tx *gorm.DB, companyId int, techallocStartDate time.Time) (int, error) {
	var unallocatedOpr int64

	err := tx.Table("trx_work_order_detail").
		Joins("LEFT JOIN trx_work_order ON trx_work_order.work_order_system_number = trx_work_order_detail.work_order_system_number").Where("trx_work_order.company_id = ?", companyId).
		Where("trx_work_order.work_order_status_id = ?", 0).
		Where("trx_work_order.work_order_date = ?", techallocStartDate).
		Count(&unallocatedOpr).Error

	if err != nil {
		return 0, err
	}

	return int(unallocatedOpr), nil
}

func (r *WorkOrderAllocationRepositoryImpl) getAutoReleased(tx *gorm.DB, companyId int, techallocStartDate time.Time) (int, error) {
	var autoReleased int64

	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Joins("LEFT JOIN trx_work_order ON trx_work_order_detail.work_order_system_number = trx_work_order.work_order_system_number").
		Where("trx_work_order.company_id = ?", companyId).
		Where("COALESCE(trx_work_order_detail.service_status_id, '') = ?", utils.SrvStatAutoRelease).
		Where("trx_work_order.work_order_date = ?", techallocStartDate).
		Count(&autoReleased).
		Error

	if err != nil {
		return 0, err
	}

	return int(autoReleased), nil
}

func (r *WorkOrderAllocationRepositoryImpl) getBookAllocTime(tx *gorm.DB, companyId int, techallocStartDate time.Time) (float64, error) {
	var bookAllocTime float64

	err := tx.Model(&transactionworkshopentities.BookingAllocation{}).
		Select("COALESCE(SUM(booking_allocation_total_hour),0)").
		Where("company_id = ? AND booking_allocation_date = ?", companyId, techallocStartDate).
		Scan(&bookAllocTime).Error

	if err != nil {
		return 0, err
	}

	return bookAllocTime, nil
}

func (r *WorkOrderAllocationRepositoryImpl) CountAvailableShifts(tx *gorm.DB, companyId int, shiftCode string, servDate time.Time) (int, error) {
	var countAvail int64

	dayOfWeek := servDate.Weekday()
	var effectiveDate time.Time
	err := tx.Model(&masterentities.ShiftSchedule{}).
		Where("company_id = ? AND shift_code = ? AND effective_date <= ? AND "+r.getDayColumn(dayOfWeek)+" = ?",
			companyId, shiftCode, servDate, true).
		Order("effective_date DESC").
		Select("effective_date").
		Limit(1).
		Pluck("effective_date", &effectiveDate).
		Error
	if err != nil {
		return 0, err
	}

	if effectiveDate.IsZero() {
		return 0, nil
	}

	err = tx.Model(&masterentities.ShiftSchedule{}).
		Where("company_id = ? AND shift_code = ? AND effective_date = ? AND "+r.getDayColumn(dayOfWeek)+" = ? AND is_active = ?",
			companyId, shiftCode, effectiveDate, true, true).
		Count(&countAvail).
		Error

	if err != nil {
		return 0, err
	}

	return int(countAvail), nil
}

// Helper function to return the day of the week
func (r *WorkOrderAllocationRepositoryImpl) getDayColumn(dayOfWeek time.Weekday) string {
	switch dayOfWeek {
	case time.Sunday:
		return "sunday"
	case time.Monday:
		return "monday"
	case time.Tuesday:
		return "tuesday"
	case time.Wednesday:
		return "wednesday"
	case time.Thursday:
		return "thursday"
	case time.Friday:
		return "friday"
	case time.Saturday:
		return "saturday"
	default:
		return ""
	}
}

// isValidTimeColumn checks if the given time column is valid
func (r *WorkOrderAllocationRepositoryImpl) isValidTimeColumn(columnName string) bool {
	validTimeColumns := map[string]bool{
		"time_allocation_0700": true, "time_allocation_0715": true, "time_allocation_0730": true,
		"time_allocation_0745": true, "time_allocation_0800": true, "time_allocation_0815": true,
		"time_allocation_0830": true, "time_allocation_0845": true, "time_allocation_0900": true,
		"time_allocation_0915": true, "time_allocation_0930": true, "time_allocation_0945": true,
		"time_allocation_1000": true, "time_allocation_1015": true, "time_allocation_1030": true,
		"time_allocation_1045": true, "time_allocation_1100": true, "time_allocation_1115": true,
		"time_allocation_1130": true, "time_allocation_1145": true, "time_allocation_1200": true,
		"time_allocation_1215": true, "time_allocation_1230": true, "time_allocation_1245": true,
		"time_allocation_1300": true, "time_allocation_1315": true, "time_allocation_1330": true,
		"time_allocation_1345": true, "time_allocation_1400": true, "time_allocation_1415": true,
		"time_allocation_1430": true, "time_allocation_1445": true, "time_allocation_1500": true,
		"time_allocation_1515": true, "time_allocation_1530": true, "time_allocation_1545": true,
		"time_allocation_1600": true, "time_allocation_1615": true, "time_allocation_1630": true,
		"time_allocation_1645": true, "time_allocation_1700": true, "time_allocation_1715": true,
		"time_allocation_1730": true, "time_allocation_1745": true, "time_allocation_1800": true,
		"time_allocation_1815": true, "time_allocation_1830": true, "time_allocation_1845": true,
		"time_allocation_1900": true, "time_allocation_1915": true, "time_allocation_1930": true,
		"time_allocation_1945": true, "time_allocation_2000": true, "time_allocation_2015": true,
		"time_allocation_2030": true, "time_allocation_2045": true, "time_allocation_2100": true,
	}

	return validTimeColumns[columnName]
}
