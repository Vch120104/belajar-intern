package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ServiceWorkshopRepositoryImpl struct {
}

func OpenServiceWorkshopRepositoryImpl() transactionworkshoprepository.ServiceWorkshopRepository {
	return &ServiceWorkshopRepositoryImpl{}
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
func getShiftStartTime(tx *gorm.DB, companyId int, shiftCode string, effectiveDate time.Time, rest bool) (float64, *exceptions.BaseErrorResponse) {
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

func getShiftEndTime(tx *gorm.DB, companyId int, shiftCode string, effectiveDate time.Time, rest bool) (float64, *exceptions.BaseErrorResponse) {
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

// Convert float64 hours to time.Duration
func hoursToDuration(hours float64) time.Duration {
	return time.Duration(hours * float64(time.Hour))
}

// Convert time string to float64 hours
func getTimeValue(timeStr string) float64 {
	t, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return 0
	}
	return float64(t.Hour()) + float64(t.Minute())/60 + float64(t.Second())/3600
}

// Convert time float64 to string
func getTime(time float64) string {
	hours := int(time)
	minutes := int((time - float64(hours)) * 60)
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

// Convert time to float64 hours
func getFloatTimeValue(t time.Time) float64 {
	return float64(t.Hour()) + float64(t.Minute())/60 + float64(t.Second())/3600
}

// Convert float64 hours to time.Time with a reference date
func getTimeFromFloatValue(hours float64, referenceDate time.Time) time.Time {
	return time.Date(
		referenceDate.Year(),
		referenceDate.Month(),
		referenceDate.Day(),
		int(hours),
		int((hours-float64(int(hours)))*60),
		0,
		0,
		time.UTC, // Adjust to the required time zone
	)
}

// GetTimeZone fetches the time difference from the external API and adjusts the time accordingly
func GetTimeZone(currentDate time.Time, companyCode int) (time.Time, error) {
	apiURL := config.EnvConfigs.SalesServiceUrl + "company-reference?page=0&limit=1000&company_id=" + strconv.Itoa(companyCode)

	var timeReferences []transactionworkshoppayloads.TimeReference
	err := utils.Get(apiURL, &timeReferences, nil)
	if err != nil {
		return time.Time{}, err
	}

	if len(timeReferences) == 0 {
		return time.Time{}, errors.New("no time reference data found for company")
	}

	timeVariance := timeReferences[0].TimeDiff

	adjustedTime := currentDate.Add(time.Hour * time.Duration(timeVariance))

	return adjustedTime, nil
}

func (r *ServiceWorkshopRepositoryImpl) GetAllByTechnicianWO(tx *gorm.DB, idTech int, idSysWo int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionworkshoppayloads.ServiceWorkshopDetailResponse, *exceptions.BaseErrorResponse) {

	var entity transactionworkshoppayloads.ServiceWorkshopRequest

	joinTable := utils.CreateJoinSelectStatement(tx, entity)
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)
	whereQuery = whereQuery.Where("technician_id = ? AND work_order_system_number = ? AND service_status_id IN (?,?,?,?)", idTech, idSysWo, utils.SrvStatDraft, utils.SrvStatStart, utils.SrvStatPending, utils.SrvStatStop)

	if err := whereQuery.Find(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.ServiceWorkshopDetailResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Work order not found",
				Err:        err,
			}
		}
		return transactionworkshoppayloads.ServiceWorkshopDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch entity",
			Err:        err,
		}
	}

	// Check if the service_status_id is valid
	validStatuses := map[int]bool{
		utils.SrvStatDraft:   true,
		utils.SrvStatStart:   true,
		utils.SrvStatPending: true,
		utils.SrvStatStop:    true,
	}
	if !validStatuses[entity.ServiceStatusId] {
		return transactionworkshoppayloads.ServiceWorkshopDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Service status not found",
		}
	}

	// Fetch data work order from external API
	WorkOrderUrl := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(entity.WorkOrderSystemNumber)
	var workOrderResponses transactionworkshoppayloads.ServiceWorkshopWoResponse
	errWorkOrder := utils.Get(WorkOrderUrl, &workOrderResponses, nil)
	if errWorkOrder != nil {
		return transactionworkshoppayloads.ServiceWorkshopDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order data from the external API",
			Err:        errWorkOrder,
		}
	}

	var serviceDetails []transactionworkshoppayloads.ServiceWorkshopResponse

	var totalRows int64
	totalRowsQuery := tx.Model(&transactionworkshopentities.ServiceLog{}).
		Where("service_log_system_number = ?", entity.ServiceLogSystemNumber).
		Count(&totalRows).Error
	if totalRowsQuery != nil {
		return transactionworkshoppayloads.ServiceWorkshopDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count service details",
			Err:        totalRowsQuery,
		}
	}

	// Fetch paginated service details
	query := tx.Model(&transactionworkshopentities.ServiceLog{}).
		Joins("INNER JOIN trx_work_order_allocation AS WTA ON trx_service_log.technician_allocation_system_number = WTA.technician_allocation_system_number").
		Select("trx_service_log.technician_allocation_system_number, trx_service_log.start_datetime, WTA.operation_code, trx_service_log.frt, trx_service_log.service_status_id, WTA.serv_actual_time, WTA.serv_pending_time, WTA.serv_progress_time, WTA.tech_alloc_start_date, WTA.tech_alloc_start_time, WTA.tech_alloc_end_date, WTA.tech_alloc_end_time").
		Where("trx_service_log.technician_allocation_line = (SELECT TOP 1 technician_allocation_line FROM trx_service_log A WHERE A.technician_allocation_system_number = trx_service_log.technician_allocation_system_number ORDER BY technician_allocation_line DESC)").
		Offset(pages.GetOffset()).
		Limit(pages.GetLimit())

	if err := query.Find(&serviceDetails).Error; err != nil {
		return transactionworkshoppayloads.ServiceWorkshopDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get service details",
			Err:        err,
		}
	}

	response := transactionworkshoppayloads.ServiceWorkshopDetailResponse{
		ServiceTypeName:         "Workshop",
		TechnicianId:            entity.TechnicianId,
		WorkOrderSystemNumber:   entity.WorkOrderSystemNumber,
		WorkOrderDocumentNumber: workOrderResponses.WorkOrderDocumentNumber,
		WorkOrderDate:           workOrderResponses.WorkOrderDate,
		ModelName:               workOrderResponses.ModelName,
		VariantName:             workOrderResponses.VariantName,
		VehicleCode:             workOrderResponses.VehicleCode,
		VehicleTnkb:             workOrderResponses.VehicleTnkb,
		ServiceDetails: transactionworkshoppayloads.ServiceWorkshopDetailsResponse{
			Page:       pages.GetPage(),
			Limit:      pages.GetLimit(),
			TotalPages: int(math.Ceil(float64(totalRows) / float64(pages.GetLimit()))),
			TotalRows:  int(totalRows),
			Data:       serviceDetails,
		},
	}

	return response, nil
}

// uspg_wtServiceLog_Insert
// --USE FOR : * INSERT NEW DATA OR UPDATE IF SERVICE STATUS IS START, PENDING OR STOP
// --USE IN MODUL :
func (r *ServiceWorkshopRepositoryImpl) StartService(tx *gorm.DB, idAlloc int, idSysWo int, idServLog int, companyId int) (bool, *exceptions.BaseErrorResponse) {

	//--============================================================================================================================
	//	--Start :	Tombol untuk memulai pekerjaan
	//	--			(validasinya : teknisi bisa mengerjakan pekerjaan lain dengan catatan tidak ada status start pada kerjaan lainnya)
	//--============================================================================================================================

	var techAlloc transactionworkshopentities.WorkOrderAllocation
	var cpcCode, pcBr string
	var oprItemCode string
	var idTech int
	currentTime := time.Now()

	dateTimeComp, err := GetTimeZone(time.Now(), companyId)
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get adjusted time",
			Err:        err,
		}
	}

	// Check if technician allocation exists
	var count int64
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Where("technician_allocation_system_number = ?", idAlloc).
		Count(&count).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check technician allocation",
			Err:        err,
		}
	}

	if count == 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Technician Allocation is not valid. Please refresh your page",
			Err:        errors.New("technician allocation not valid"),
		}
	}

	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Select("technician_group_id, technician_id, operation_code, brand_id, profit_center_id, work_order_line, foreman_id, shift_code").
		Where("technician_allocation_system_number = ?", idAlloc).
		First(&techAlloc).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve technician allocation",
			Err:        err,
		}
	}

	oprItemCode = techAlloc.OperationCode
	idTech = techAlloc.TechnicianId
	woLine := techAlloc.WorkOrderLine

	var maxTechallocLine int64

	// Query the maximum TECHALLOC_LINE for the given TECHALLOC_SYS_NO
	err = tx.Model(&transactionworkshopentities.ServiceLog{}).
		Select("ISNULL(MAX(technician_allocatation_line), 0)").
		Where("technician_allocation_system_number = ?", idAlloc).
		Scan(&maxTechallocLine).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve the next TECHALLOC_LINE",
			Err:        err,
		}
	}

	techallocLine := int(maxTechallocLine) + 1

	// Check if the previous service status is in draft or transfer
	var prevStatus int
	err = tx.Model(&transactionworkshopentities.ServiceLog{}).
		Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, techallocLine-1).
		Select("service_status_id").
		Scan(&prevStatus).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get previous service status",
			Err:        err,
		}
	}

	if prevStatus != utils.SrvStatDraft && prevStatus != utils.SrvStatTransfer {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Service already started",
			Err:        err,
		}
	}

	var logExists int64
	err = tx.Model(&transactionworkshopentities.ServiceLog{}).
		Where("technician_allocation_system_number = ? AND company_id = ? AND service_status_id = ? AND technician_allocation_line = (SELECT MAX(technician_allocation_line) FROM trx_service_log WHERE technician_allocation_system_number = ?)", idAlloc, companyId, utils.SrvStatStart, idAlloc).
		Count(&logExists).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check service log existence",
			Err:        err,
		}
	}

	if logExists == 0 {

		if pcBr == "00003" {
			// Check if there are any pending operations for the technician
			var pendingCount int64
			err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
				Where("service_status_id NOT IN (?) AND technician_id = ?", []int{utils.SrvStatStop, utils.SrvStatTransfer, utils.SrvStatPending, utils.SrvStatDraft, utils.SrvStatQcPass, 0}, idAlloc).
				Count(&pendingCount).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to check pending operations",
					Err:        err,
				}
			}

			if pendingCount > 0 {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusConflict,
					Message:    "Previous Operation must be Stop, Transfer Or Pending",
					Err:        err,
				}
			}

			var technicianWorking bool
			err = tx.Model(&transactionworkshopentities.ServiceLog{}).
				Where("work_order_operation_id <> ? AND technician_id = ? AND service_status_id = ? AND EXISTS (SELECT MAX(technician_allocation_line) FROM trx_service_log B WHERE B.technician_allocation_system_number = A.technician_allocation_system_number HAVING MAX(technician_allocation_line) = A.technician_allocation_line)", oprItemCode, idTech, utils.SrvStatStart).
				Select("COUNT(*) > 0").
				Scan(&technicianWorking).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to check technician's current operations",
					Err:        err,
				}
			}

			if technicianWorking {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusConflict,
					Message:    "Technician already started another operation",
					Err:        err,
				}
			}

			dateTimeComp, _ := GetTimeZone(time.Now(), companyId)
			startDatetime := techAlloc.TechAllocStartDate
			formattedDateTimeComp := dateTimeComp.Format("02 Jan 2006")
			formattedStartDatetime := startDatetime.Format("02 Jan 2006")

			if formattedDateTimeComp != formattedStartDatetime {
				errorMsg := fmt.Sprintf("Operation must be started on %s. Please contact your foreman for re-allocation.", formattedStartDatetime)
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusConflict,
					Message:    errorMsg,
					Err:        errors.New(errorMsg),
				}
			}
		} else {
			// Validation for profit center GR
			var (
				startDatetime time.Time
				endDatetime   time.Time
				frt           float64
			)

			err = tx.Model(&transactionworkshopentities.ServiceLog{}).
				Where("technician_allocation_system_number = ? AND service_status_id = ?", idAlloc, utils.SrvStatDraft).
				Order("technician_allocation_line DESC").
				Limit(1).
				Select("start_datetime, end_datetime, frt").
				Scan(&struct {
					StartDatetime *time.Time `gorm:"column:start_datetime"`
					EndDatetime   *time.Time `gorm:"column:end_datetime"`
					Frt           *float64   `gorm:"column:frt"`
				}{
					&startDatetime,
					&endDatetime,
					&frt,
				}).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to get service log dates",
					Err:        err,
				}
			}

			//--validation operation must start as specified on the allocation
			//--AND must be the FIRST LOG DRAFT
			if techallocLine-1 == 1 {
				if currentTime.Before(startDatetime) ||
					currentTime.After(endDatetime) ||
					currentTime.Format("2006-01-02") != startDatetime.Format("2006-01-02") {

					errorMsg := "Server Time : " + currentTime.Format("15:04:05") + "<br/>Operation must start on allocated time from " + startDatetime.Format("15:04:05") + " to " + endDatetime.Format("15:04:05") + " . Please contact your foreman for re-allocation."
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusConflict,
						Message:    errorMsg,
						Err:        err,
					}
				}
			}
		}

		var woSysNo int
		var startDatetime time.Time
		var endDatetime time.Time
		var shiftCode string
		var sumActualTime float64
		var frt float64

		err := tx.Model(&transactionworkshopentities.ServiceLog{}).
			Where("technician_allocation_system_number = ? AND service_status_id = ?", idAlloc, utils.SrvStatDraft).
			Order("service_log_system_number DESC").
			First(&struct {
				WoSysNo       int       `gorm:"column:work_order_system_number"`
				StartDatetime time.Time `gorm:"column:start_datetime"`
				EndDatetime   time.Time `gorm:"column:end_datetime"`
				ShiftCode     string    `gorm:"column:shift_code"`
				Frt           float64   `gorm:"column:frt"`
			}{
				WoSysNo:       woSysNo,
				StartDatetime: startDatetime,
				EndDatetime:   endDatetime,
				ShiftCode:     shiftCode,
				Frt:           frt,
			}).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve data log with status Draft",
				Err:        err,
			}
		}

		//--Delete the log with status Draft--
		err = tx.Where("technician_allocation_system_number = ? AND service_status_id = ?", idAlloc, utils.SrvStatDraft).
			Delete(&transactionworkshopentities.ServiceLog{}).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to delete data log with status Draft",
				Err:        err,
			}
		}

		// if pcBr <> "00003" {
		// 	// --Deallocate the time with DRAFT--
		// 	startTime := getTimeValue(startDatetime.Format("15:04:05"))
		// 	endTime := getTimeValue(endDatetime.Format("15:04:05"))
		// 	remarkAvail := "Release DRAFT from START (Start Time: " + startTime.String() + ";End Time: " + endTime.String() + ")"
		// 	// EXEC uspg_atWoTechAllocAvailable_Insert
		// 	// --@Option = 0,
		// 	// --@Start_Time = @StartTime,
		// 	// --@End_Time = @EndTime,
		// 	// --@Technician_Emp_No = @Technician_Id,
		// 	// --@Company_Code = @Company_Code,
		// 	// --@Service_Datetime = @Start_Datetime,
		// 	// --@Shift_Code = @Shift_Code,
		// 	// --@Foreman_Emp_No = @Foreman_Emp_No,
		// 	// --@Ref_Type = @Ref_Type_Avail,
		// 	// --@Ref_Sys_No = @Techalloc_Sys_No,
		// 	// --@Ref_Line = 0,
		// 	// --@Remark = @Remark_Avail,
		// 	// --@Creation_User_Id = @Creation_User_Id,
		// 	// --@Change_User_Id = @Change_User_Id
		// }

		err = tx.Model(&transactionworkshopentities.ServiceLog{}).
			Select("ISNULL(SUM(actual_time), 0)").
			Where("technician_allocation_system_number = ?", idAlloc).
			Scan(&sumActualTime).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate actual time",
				Err:        err,
			}
		}

		if sumActualTime >= frt {
			sumActualTime = 0.5 * frt
		}

		startRestTime, _ := getShiftStartTime(tx, companyId, techAlloc.ShiftCode, dateTimeComp, true)
		endRestTime, _ := getShiftEndTime(tx, companyId, techAlloc.ShiftCode, dateTimeComp, true)

		startRestTimeDuration := hoursToDuration(startRestTime)
		endRestTimeDuration := hoursToDuration(endRestTime)

		// Load Jakarta time zone (UTC+7)
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to load time location",
				Err:        err,
			}
		}

		additionalTime := hoursToDuration(frt - sumActualTime)

		startTimeInLoc := startDatetime.In(loc)
		startTimeValue := getTimeValue(startTimeInLoc.Format("15:04:05"))

		if startTimeValue+additionalTime.Hours() >= 24 {
			endDatetime = time.Date(
				startDatetime.Year(), startDatetime.Month(), startDatetime.Day(),
				23, 45, 0, int(time.Nanosecond), loc,
			)
		} else {
			endDatetime = startDatetime.In(loc).Add(additionalTime)
		}

		startTime := getTimeValue(startDatetime.In(loc).Format("15:04:05"))
		endTime := getTimeValue(endDatetime.Format("15:04:05"))

		if startTime+additionalTime.Hours() < 24 {
			startRestTimeValue := startRestTime
			endRestTimeValue := endRestTime

			if (startTime <= startRestTimeValue && endRestTimeValue <= endTime) ||
				(startTime <= startRestTimeValue && startRestTimeValue < endTime) ||
				(startTime < endRestTimeValue && endRestTimeValue <= endTime) {

				timeAfterRest := endRestTimeDuration - startRestTimeDuration
				if getTimeValue(endDatetime.Format("15:04:05"))+timeAfterRest.Hours() < 24 {
					endDatetime = endDatetime.Add(timeAfterRest)
				} else {
					endDatetime = time.Date(
						endDatetime.Year(), endDatetime.Month(), endDatetime.Day(),
						23, 45, 0, int(time.Nanosecond), loc,
					)
				}
			}
		}

		// inserting the new service log with status START
		newServiceLog := transactionworkshopentities.ServiceLog{
			CompanyId:             companyId,
			WorkOrderSystemNumber: idSysWo,
			TechnicianId:          idTech,
			ServiceStatusId:       utils.SrvStatStart,
			StartDatetime:         dateTimeComp,
			EndDatetime:           endDatetime,
			ActualTime:            0,
			PendingTime:           0,
			EstimatedPendingTime:  0,
		}

		err = tx.Create(&newServiceLog).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to insert new service log",
				Err:        err,
			}
		}

		// Insert or update in wtSmrWOServiceTime
		var smrWOServiceTimeExists bool
		err = tx.Model(&transactionworkshopentities.WorkOrderServiceTime{}).
			Where("work_order_system_number = ? AND operation_item_id = ? ", idSysWo, oprItemCode).
			Select("COUNT(*) > 0").
			Scan(&smrWOServiceTimeExists).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to check WorkOrderServiceTime existence",
				Err:        err,
			}
		}

		if smrWOServiceTimeExists {
			err = tx.Model(&transactionworkshopentities.WorkOrderServiceTime{}).
				Where("work_order_system_number = ? AND operation_item_id = ? ", idSysWo, oprItemCode).
				Updates(map[string]interface{}{
					"end_datetime": dateTimeComp, //--- :GH getdate diganti timecomp
				}).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update WorkOrderServiceTime",
					Err:        err,
				}
			}
		} else {
			// Insert new record in wtSmrWOServiceTime
			newWorkOrderServiceTime := transactionworkshopentities.WorkOrderServiceTime{
				CompanyId:             companyId,
				WorkOrderSystemNumber: idSysWo,
				OperationItemCode:     oprItemCode,
				StartDatetime:         dateTimeComp, //--- :GH getdate diganti timecomp
				EndDatetime:           endDatetime,
			}

			err = tx.Create(&newWorkOrderServiceTime).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to insert into wtSmrWOServiceTime",
					Err:        err,
				}
			}
		}

		//--==Release allocation after the current operation if intersect==--
		if pcBr != "00003" {
			// Fetch relevant service logs
			// DECLARE @CSR CURSOR
			var serviceLogs []transactionworkshopentities.ServiceLog
			err := tx.Model(&transactionworkshopentities.ServiceLog{}).
				Select("technician_allocation_system_number, start_datetime, end_datetime").
				Where("technician_allocation_system_number <> ? AND company_id = ? AND technician_id = ? AND CONVERT(VARCHAR, start_datetime, 106) = CONVERT(VARCHAR, ?, 106)", idAlloc, companyId, idTech, dateTimeComp).
				Where("EXISTS (SELECT TOP 1 technician_allocation_line FROM trx_service_log WHERE technician_allocation_system_number = trx_service_log.technician_allocation_system_number AND service_status_id = ? AND technician_allocation_line = 1 AND CONVERT(VARCHAR, start_datetime, 106) = CONVERT(VARCHAR, ?, 106) AND trx_service_log.technician_allocation_line = technician_allocation_line ORDER BY technician_allocation_line DESC)").
				Find(&serviceLogs).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve service logs",
					Err:        err,
				}
			}

			filteredLogs := []transactionworkshopentities.ServiceLog{}
			for _, log := range serviceLogs {

				startTimeValue := getTimeValue(log.StartDatetime.Format("15:04:05"))
				endTimeValue := getTimeValue(log.EndDatetime.Format("15:04:05"))
				queryTimeValue := getTimeValue(dateTimeComp.Format("15:04:05"))

				if startTimeValue <= queryTimeValue && endTimeValue >= queryTimeValue {
					filteredLogs = append(filteredLogs, log)
				}
			}

			for _, log := range filteredLogs {
				err := tx.Where("technician_allocation_system_number = ?", log.TechnicianAllocationSystemNumber).Delete(&transactionworkshopentities.WorkOrderAllocation{}).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to delete from WorkOrderAllocation",
						Err:        err,
					}
				}

				err = tx.Where("technician_allocation_system_number = ?", log.TechnicianAllocationSystemNumber).Delete(&transactionworkshopentities.ServiceLog{}).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to delete from ServiceLog",
						Err:        err,
					}
				}
			}

			// DECLARE @CSR_1 CURSOR
			var serviceLogs_1 []transactionworkshopentities.ServiceLog
			err = tx.Model(&transactionworkshopentities.ServiceLog{}).
				Select("technician_allocation_system_number, start_datetime, end_datetime").
				Where("technician_allocation_system_number <> ? AND company_id = ? AND technician_id = ? AND CONVERT(VARCHAR, start_datetime, 106) = CONVERT(VARCHAR, ?, 106)", idAlloc, companyId, idTech, dateTimeComp).
				Where("EXISTS (SELECT TOP 1 technician_allocation_line FROM trx_service_log WHERE technician_allocation_system_number = trx_service_log.technician_allocation_system_number AND service_status_id = ? AND technician_allocation_line = 1 AND CONVERT(VARCHAR, start_datetime, 106) = CONVERT(VARCHAR, ?, 106) AND trx_service_log.technician_allocation_line = technician_allocation_line ORDER BY technician_allocation_line DESC)").
				Find(&serviceLogs_1).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve service logs",
					Err:        err,
				}
			}

			for _, log := range serviceLogs_1 {
				// Konversi waktu
				startTimeValue := getTimeValue(log.StartDatetime.Format("15:04:05"))
				endTimeValue := getTimeValue(log.EndDatetime.Format("15:04:05"))
				queryTimeValue := getTimeValue(dateTimeComp.Format("15:04:05"))

				// Hitung selisih waktu
				diffDuration := dateTimeComp.Sub(log.StartDatetime)
				diffTime := diffDuration.Hours()

				// Jika alokasi saat ini berada di dalam alokasi sebelumnya
				if startTimeValue <= queryTimeValue && endTimeValue >= queryTimeValue {
					// Prepare remarks
					// startTimeFormatted := log.StartDatetime.Format("15:04:05")
					// endTimeFormatted := log.EndDatetime.Format("15:04:05")

					// remarkAvailBefore := fmt.Sprintf(
					// 	"Release OLD TIME (MOVE ALLOC) from START (Start Time: %s; End Time: %s)",
					// 	startTimeFormatted, endTimeFormatted,
					// )

					// remarkAvailAfter := fmt.Sprintf(
					// 	"Allocate NEW TIME (MOVE ALLOC) from START (Start Time: %s; End Time: %s)",
					// 	startTimeFormatted, endTimeFormatted,
					// )

					// Update service log --Move the Service LOG--
					updateErr := tx.Model(&transactionworkshopentities.ServiceLog{}).
						Where("technician_allocation_system_number = ? AND service_status_id = ?", log.TechnicianAllocationSystemNumber, utils.SrvStatDraft).
						Updates(map[string]interface{}{
							"start_datetime": log.StartDatetime.Add(time.Duration(diffTime) * time.Hour),
							"end_datetime":   log.EndDatetime.Add(time.Duration(diffTime) * time.Hour),
						}).Error

					if updateErr != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to update service log",
							Err:        updateErr,
						}
					}

					// Update allocation --Move the Allocation--
					updateAllocErr := tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
						Where("technician_allocation_system_number = ?", log.TechnicianAllocationSystemNumber).
						Updates(map[string]interface{}{
							"techalloc_start_time": gorm.Expr("techalloc_start_time + ?", diffTime),
							"techalloc_end_time":   gorm.Expr("techalloc_end_time + ?", diffTime),
						}).Error

					if updateAllocErr != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to update techalloc",
							Err:        updateAllocErr,
						}
					}
					// EXEC uspg_atWoTechAllocAvailable_Insert
					// --@Option = 0,
					// Insert record for the remark before moving
					// newAllocationBefore := transactionworkshopentities.WorkOrderAllocationAvailable{
					// 	StartTime:             startTimeValue,
					// 	EndTime:               endTimeValue,
					// 	TechnicianId:          idTech,
					// 	CompanyId:             companyId,
					// 	ServiceDateTime:       log.StartDatetime,
					// 	ShiftCode:             techAlloc.ShiftCode,
					// 	ForemanId:             techAlloc.ForemanId,
					// 	ReferenceType:         "0",
					// 	ReferenceSystemNumber: log.TechnicianAllocationSystemNumber,
					// 	ReferenceLine:         0,
					// 	Remark:                remarkAvailBefore,
					// }

					// insertBeforeErr := tx.Create(&newAllocationBefore).Error
					// if insertBeforeErr != nil {
					// 	return false, &exceptions.BaseErrorResponse{
					// 		StatusCode: http.StatusInternalServerError,
					// 		Message:    "Failed to insert new allocation before",
					// 		Err:        insertBeforeErr,
					// 	}
					// }

					// EXEC uspg_atWoTechAllocAvailable_Insert
					// --@Option = 1,
					// Insert record for the remark after moving
					// newAllocationAfter := transactionworkshopentities.WorkOrderAllocationAvailable{
					// 	StartTime:             startTimeValue + diffTime,
					// 	EndTime:               endTimeValue + diffTime,
					// 	TechnicianId:          idTech,
					// 	CompanyId:             companyId,
					// 	ServiceDateTime:       log.StartDatetime,
					// 	ShiftCode:             techAlloc.ShiftCode,
					// 	ForemanId:             techAlloc.ForemanId,
					// 	ReferenceType:         "0",
					// 	ReferenceSystemNumber: log.TechnicianAllocationSystemNumber,
					// 	ReferenceLine:         0,
					// 	Remark:                remarkAvailAfter,
					// }

					// insertAfterErr := tx.Create(&newAllocationAfter).Error
					// if insertAfterErr != nil {
					// 	return false, &exceptions.BaseErrorResponse{
					// 		StatusCode: http.StatusInternalServerError,
					// 		Message:    "Failed to insert new allocation after",
					// 		Err:        insertAfterErr,
					// 	}
					// }
				}
			}

		}

		if cpcCode == pcBr {
			// Fetch existing record to check if TECHALLOC_START_TIME is 0
			var techAlloc transactionworkshopentities.WorkOrderAllocation
			err := tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
				Where("technician_allocation_system_number = ?", idAlloc).
				Select("tech_alloc_start_date").
				First(&techAlloc).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to fetch techalloc record",
					Err:        err,
				}
			}

			// If TECHALLOC_START_TIME is 0, update it
			if techAlloc.TechAllocStartTime == 0 {
				startTime := getTimeValue(dateTimeComp.Format("15:04:05"))
				updateErr := tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
					Where("technician_allocation_system_number = ?", idAlloc).
					Updates(map[string]interface{}{
						"tech_alloc_start_time": startTime,
					}).Error

				if updateErr != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to update techalloc start time",
						Err:        updateErr,
					}
				}
			}
		}

		// Update wtWorkOrder2
		updateErr := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Where("work_order_system_number = ? AND work_order_operation_item_line = ?", idSysWo, woLine).
			Updates(map[string]interface{}{
				"work_order_status_id": utils.SrvStatStart,
			}).Error

		if updateErr != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update wtWorkOrder2 service status",
				Err:        updateErr,
			}
		}

		// Update wtWorkOrder0
		updateErr = tx.Model(&transactionworkshopentities.WorkOrder{}).
			Where("work_order_system_number = ?", idSysWo).
			Updates(map[string]interface{}{
				"work_order_status_id": utils.WoStatOngoing,
			}).Error

		if updateErr != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update wtWorkOrder0 service status",
				Err:        updateErr,
			}
		}

		// Update atWoTechAlloc
		updateErr = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
			Where("technician_allocation_system_number = ?", idAlloc).
			Updates(map[string]interface{}{
				"service_status_id": utils.SrvStatStart,
			}).Error

		if updateErr != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update atWoTechAlloc service status",
				Err:        updateErr,
			}
		}

		// Check if there's a previous pending log entry
		var exists bool
		err = tx.Model(&transactionworkshopentities.ServiceLog{}).
			Where("technician_allocation_system_number = ? AND technician_allocation_line = ? AND service_status_id = ?", idAlloc, techallocLine-2, utils.SrvStatPending).
			Select("1").
			Scan(&exists).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to check for previous pending log entry",
				Err:        err,
			}
		}

		if exists {
			// Calculate Pending Time
			var startTime time.Time
			err = tx.Model(&transactionworkshopentities.ServiceLog{}).
				Select("start_datetime").
				Where("technician_allocation_system_number = ? AND technician_allocation_line = (SELECT MAX(technician_allocation_line) FROM trx_service_log WHERE technician_allocation_system_number = ? AND technician_allocation_line < ? AND service_status_id = ?)",
					idAlloc, idAlloc, techallocLine, utils.SrvStatPending).
				Pluck("start_datetime", &startTime).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve start time for pending calculation",
					Err:        err,
				}
			}

			pendingTime := float64(time.Since(startTime).Minutes()) / 60.0

			// Update previous pending service log
			updateErr := tx.Model(&transactionworkshopentities.ServiceLog{}).
				Where("technician_allocation_system_number = ? AND technician_allocation_line = (SELECT MAX(technician_allocation_line) FROM trx_service_log WHERE technician_allocation_system_number = ? AND technician_allocation_line < ? AND service_status_id = ?)",
					idAlloc, idAlloc, techallocLine, utils.SrvStatPending).
				Updates(map[string]interface{}{
					"end_datetime": dateTimeComp,
					"pending_time": pendingTime,
				}).Error

			if updateErr != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update previous pending service log",
					Err:        updateErr,
				}
			}

			// Update pending time in techalloc
			updateErr = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
				Where("technician_allocation_system_number = ?", idAlloc).
				Updates(map[string]interface{}{
					"serv_pending_time":  gorm.Expr("serv_pending_time + ?", pendingTime),
					"serv_progress_time": gorm.Expr("serv_progress_time + ?", pendingTime),
				}).Error

			if updateErr != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update techalloc pending time",
					Err:        updateErr,
				}
			}
		}

		//--==============UPDATE CONTRACT SERVICE=========
		var contractServSysNo int
		var packageCodeCs int

		// Fetch values from wtWorkOrder0 and wtWorkOrder2 tables
		err = tx.Model(&transactionworkshopentities.WorkOrder{}).
			Select("ISNULL(WO.contract_service_system_number, 0) AS contract_service_system_number, ISNULL(WO2.package_id, 0) AS package_id").
			Joins("INNER JOIN trx_work_order_detail WO2 ON WO.work_order_system_number = WO2.work_order_system_number").
			Where("WO.work_order_system_number = ? AND WO2.work_order_operation_item_line = ? AND WO2.BILL_CODE = ?", woSysNo, woLine, "S").
			Scan(&struct {
				ContractServSysNo int `gorm:"column:contract_service_system_number"`
				PackageCodeCs     int `gorm:"column:package_id"`
			}{
				ContractServSysNo: contractServSysNo,
				PackageCodeCs:     packageCodeCs,
			}).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch contract service details",
				Err:        err,
			}
		}

		// Check conditions and perform update if necessary
		if contractServSysNo != 0 {
			var exists bool
			err = tx.Model(&transactionworkshopentities.ContractServiceItemDetail{}).
				Where("contract_service_system_number = ? AND package_id = ? AND ISNULL(total_use_frt_quantity, 0) = 0", contractServSysNo, packageCodeCs).
				Select("1").
				Scan(&exists).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to check if contract service exists",
					Err:        err,
				}
			}

			if exists {
				// Perform the update
				updateErr := tx.Model(&transactionworkshopentities.ContractServiceItemDetail{}).
					Where("contract_service_system_number = ? AND package_id = ? AND item_id = ?", contractServSysNo, packageCodeCs, oprItemCode).
					Updates(map[string]interface{}{
						"total_use_frt_quantity": gorm.Expr("frt_quantity"),
					}).Error

				if updateErr != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to update contract service",
						Err:        updateErr,
					}
				}
			}
		}
	} else {

		var woDocNo string
		var errMessage string

		// Fetch Wo_Doc_No
		err := tx.Model(&transactionworkshopentities.ServiceLog{}).
			Select("A.work_order_document_number").
			Joins("INNER JOIN trx_service_log B ON A.technician_allocation_system_number = B.technician_allocation_system_number").
			Where("A.technician_allocation_system_number = ? AND A.technician_id = ? AND A.service_status_id = ? AND A.technician_allocation_line = (SELECT MAX(B.technician_allocation_line) FROM trx_service_log B WHERE B.technician_allocation_system_number = A.technician_allocation_system_number)",
				idAlloc, idTech, utils.SrvStatStart).
			Pluck("A.work_order_document_number", &woDocNo).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve work order document number",
				Err:        err,
			}
		}

		// Check if Wo_Doc_No was found and generate error if not found
		if woDocNo != "" {
			errMessage = "Technician already has started another service (" + woDocNo + ")"
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Message:    errMessage,
			}
		}
	}

	return true, nil
}

// uspg_wtServiceLog_Insert
// --USE FOR : * INSERT NEW DATA OR UPDATE IF SERVICE STATUS IS START, PENDING OR STOP
// --USE IN MODUL :
func (r *ServiceWorkshopRepositoryImpl) PendingService(tx *gorm.DB, idAlloc int, idSysWo int, idServLog int, companyId int) (bool, *exceptions.BaseErrorResponse) {

	var (
		techAlloc        transactionworkshopentities.WorkOrderAllocation
		woSysNo          int
		estPendingTime   float64
		frt              float64
		shiftCode        string
		startDatetime    time.Time
		endDatetime      time.Time
		maxTechallocLine int
		serviceStatus    int
		pcBr             string
		oprItemCode      string
		woLine           int
		idTech           int
	)

	// Get the maximum technician allocation line
	err := tx.Model(&transactionworkshopentities.ServiceLog{}).
		Select("COALESCE(MAX(technician_allocation_line), 0)").
		Where("technician_allocation_system_number = ?", idAlloc).
		Scan(&maxTechallocLine).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve the next TECHALLOC_LINE",
			Err:        err,
		}
	}

	// Retrieve the required fields
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Select("technician_group, technician_id, operation_code, brand_id, profit_center_id, work_order_line, foreman_id, shift_code").
		Where("technician_allocation_system_number = ?", idAlloc).
		First(&techAlloc).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve technician allocation",
			Err:        err,
		}
	}

	techallocLine := maxTechallocLine + 1
	oprItemCode = techAlloc.OperationCode
	idTech = techAlloc.TechnicianId
	woLine = techAlloc.WorkOrderLine

	// Check if the SERVICE_STATUS condition is met
	err = tx.Model(&transactionworkshopentities.ServiceLog{}).
		Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, techallocLine-1).
		Pluck("service_status_id", &serviceStatus).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check service status",
			Err:        err,
		}
	}

	if serviceStatus == utils.SrvStatStart {
		if estPendingTime == 0 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Estimate Pending Time must be filled",
			}
		}

		// Get current time with company timezone
		dateTimeComp, _ := GetTimeZone(time.Now(), companyId)

		// Check if dates match
		if dateTimeComp.Format("2006-01-02") != startDatetime.Format("2006-01-02") {
			errorMsg := fmt.Sprintf("Operation must be pending on %s. Please contact your foreman for re-allocation.", startDatetime.Format("2006-01-02"))
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    errorMsg,
			}
		}

		// Validate if pending + estimated time exceed the current day
		timeValueComp := getTimeValue(dateTimeComp.Format("15:04:05"))
		if timeValueComp+estPendingTime >= 24 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Service pending cannot pass today. Please use transfer for service on different date",
			}
		}

		// Calculate pending end datetime with estimated time
		endDatetimeValue := timeValueComp + estPendingTime
		endDatetimeStr := dateTimeComp.Format("2006-01-02") + " " + getTime(endDatetimeValue)

		// Update actual time in the log
		var actualTime float64
		err = tx.Model(&transactionworkshopentities.ServiceLog{}).
			Select("(DATEDIFF(minute, start_datetime, ?) / 60.0) + (DATEDIFF(minute, start_datetime, ?) % 60 / 60.0)", dateTimeComp, dateTimeComp).
			Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, techallocLine-1).
			Scan(&actualTime).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate actual time",
				Err:        err,
			}
		}

		err = tx.Model(&transactionworkshopentities.ServiceLog{}).
			Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, techallocLine-1).
			Updates(map[string]interface{}{
				"actual_time":  actualTime,
				"end_datetime": dateTimeComp,
			}).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update service log",
				Err:        err,
			}
		}

		// Update the related service time log
		err = tx.Model(&transactionworkshopentities.WorkOrderServiceTime{}).
			Where("wo_sys_no = ? AND opr_item_code = ?", woSysNo, shiftCode).
			Updates(map[string]interface{}{
				"end_datetime": dateTimeComp,
			}).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update service time log",
				Err:        err,
			}
		}

		// Calculate new draft end time
		var sumActualTime float64
		err = tx.Model(&transactionworkshopentities.ServiceLog{}).
			Select("SUM(actual_time)").
			Where("technician_allocation_system_number = ?", idAlloc).
			Scan(&sumActualTime).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate sum of actual time",
				Err:        err,
			}
		}

		if sumActualTime >= frt {
			sumActualTime = 0.5 * frt
		}

		var draftEndTime time.Time
		if getTimeValue(endDatetimeStr)+frt-sumActualTime >= 24 {
			draftEndTime, _ = time.Parse("2006-01-02 15:04:05", endDatetime.Format("2006-01-02")+" 23:45:00")
		} else {
			// Add (frt - sumActualTime) hours to endDatetime
			draftEndTime = endDatetime.Add(time.Duration((frt - sumActualTime) * float64(time.Hour)))
		}

		// Update draft end time
		startRestTime, _ := getShiftStartTime(tx, companyId, shiftCode, dateTimeComp, true) // Implement this function
		endRestTime, _ := getShiftEndTime(tx, companyId, shiftCode, dateTimeComp, true)     // Implement this function

		startDatetimeValue := getTimeValue(startDatetime.Format("15:04:05"))
		startRestTimeValue := startRestTime
		endRestTimeValue := endRestTime
		draftEndTimeValue := getFloatTimeValue(draftEndTime)

		// Perform the checks with float64 values
		if startDatetimeValue+frt-sumActualTime < 24 {
			if (startDatetimeValue <= startRestTimeValue && endRestTimeValue <= draftEndTimeValue) ||
				(startDatetimeValue <= startRestTimeValue && startRestTimeValue < draftEndTimeValue) ||
				(startDatetimeValue < endRestTimeValue && endRestTimeValue <= draftEndTimeValue) {

				// Convert endRestTime float64 to time.Time
				endRestTimeAsTime := getTimeFromFloatValue(endRestTimeValue, dateTimeComp)
				if draftEndTimeValue < endRestTimeValue {
					draftEndTime = endRestTimeAsTime.Add(time.Duration((frt - sumActualTime) * float64(time.Hour)))
					draftEndTimeValue = getFloatTimeValue(draftEndTime)
				}
			}
		}
		//--==Release allocation that intersect with the current operation==--
		if pcBr == "00003" {
			// Cursor equivalent logic
			var intersectingRecords []struct {
				TechallocSysNo int
				StartDatetime  time.Time
				EndDatetime    time.Time
			}

			err = tx.Model(&transactionworkshopentities.ServiceLog{}).
				Select("SVC.technician_allocation_system_number, SVC.start_datetime, SVC.end_datetime").
				Joins("JOIN wtServiceLog A ON A.technician_allocation_system_number = SVC.technician_allocation_system_number").
				Where("SVC.technician_allocation_system_number <> ? AND SVC.company_id = ? AND SVC.technician_id = ? AND SVC.shift_code = ? AND CONVERT(VARCHAR,SVC.start_datetime,106) = CONVERT(VARCHAR,?,106) AND A.service_status_id = ? AND CONVERT(VARCHAR,A.start_datetime,106) = CONVERT(VARCHAR,?,106) AND A.TECHALLOC_LINE = 1 AND dbo.getTimeValue(CONVERT(VARCHAR,A.start_datetime,108)) < dbo.getTimeValue(CONVERT(VARCHAR,?,108)) AND SVC.TECHALLOC_LINE = A.TECHALLOC_LINE",
					idAlloc, companyId, idServLog, shiftCode, dateTimeComp, utils.SrvStatStart, dateTimeComp, draftEndTimeValue).
				Scan(&intersectingRecords).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to check for intersecting operations",
					Err:        err,
				}
			}

			for _, record := range intersectingRecords {
				// startTime := getTimeValue(record.StartDatetime.Format("15:04:05"))
				// endTime := getTimeValue(record.EndDatetime.Format("15:04:05"))

				// remarkAvail := fmt.Sprintf("Release INTERSECT OPERATION from PENDING (Start Time: %v; End Time: %v)", startTime, endTime)
				// // Execute your custom logic for handling intersecting operations here
				// // --EXEC uspg_atWoTechAllocAvailable_Insert
				// // --@Option = 0,
				// // Insert record for the remark before moving
				// newAllocationBefore := transactionworkshopentities.WorkOrderAllocationAvailable{
				// 	StartTime:     startTime,
				// 	EndTime:       endTime,
				// 	CompanyId:     companyId,
				// 	ReferenceLine: 0,
				// 	Remark:        remarkAvail,
				// }

				// insertBeforeErr := tx.Create(&newAllocationBefore).Error
				// if insertBeforeErr != nil {
				// 	return false, &exceptions.BaseErrorResponse{
				// 		StatusCode: http.StatusInternalServerError,
				// 		Message:    "Failed to insert availability record",
				// 		Err:        insertBeforeErr,
				// 	}
				// }

				err = tx.Delete(&transactionworkshopentities.WorkOrderServiceTime{}, "wo_sys_no = ? AND opr_item_code = ?", woSysNo, shiftCode).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to delete work order service time",
						Err:        err,
					}
				}

				err = tx.Delete(&transactionworkshopentities.ServiceLog{}, "techalloc_sys_no = ?", record.TechallocSysNo).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to delete service log",
						Err:        err,
					}
				}
			}

		}

		if pcBr != "00003" {
			var startTime string

			// Fetch the Start Time
			err := tx.Model(&transactionworkshopentities.ServiceLog{}).
				Select("CONVERT(VARCHAR, START_DATETIME, 108) AS StartTime").
				Where("TECHALLOC_SYS_NO = ? AND TECHALLOC_LINE = ?", idAlloc, techallocLine-1).
				Pluck("StartTime", &startTime).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to fetch start time",
					Err:        err,
				}
			}

			// Set the End Time (Update this with actual calculation or function call)
			// endTime = "dbo.getTimeValue(CONVERT(VARCHAR, ?, 108))" // Replace with actual function call or calculation
			// remarkAvail = fmt.Sprintf("Allocate NEW TIME from PENDING (Start Time: %s; End Time: %s)", startTime, endTime)

			// uspg_atWoTechAllocAvailable_Insert
			//	@Option = 1,
			//	@Start_Time = ?,
			//	@End_Time = ?,
			//	@Technician_Emp_No = ?,
			//	@Company_Code = ?,
			//	@Service_Datetime = ?,
			//	@Shift_Code = ?,
			//	@Foreman_Emp_No = ?,
			//	@Ref_Type = ?,
			//	@Ref_Sys_No = ?,
			//	@Ref_Line = ?,
			//	@Remark = ?,
			//	@Creation_User_Id = ?,
			//	@Change_User_Id = ?`,
			//	startTime, endTime, technicianId, companyCode, startDatetime, shiftCode, foremanEmpNo, refTypeAvail, techallocSysNo, tempAllocLine, remarkAvail, creationUserId, changeUserId).Error
			// if err != nil {
			//	return fmt.Errorf("failed to insert into atWoTechAllocAvailable: %w", err)
			// }

			// Update WorkOrder2 Status
			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Where("work_order_system_number = ? AND work_order_operation_item_line = ?", idSysWo, woLine).
				Updates(map[string]interface{}{
					"WO_OPR_STATUS": utils.SrvStatPending,
				}).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update workorder detail",
					Err:        err,
				}
			}

			// Update atWoTechAlloc Service Status
			err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
				Where("technician_allocation_system_number = ?", idAlloc).
				Updates(map[string]interface{}{
					"service_status_id":  utils.SrvStatPending,
					"serv_actual_time":   gorm.Expr("serv_actual_time + ?", actualTime),
					"serv_progress_time": gorm.Expr("serv_progress_time + ?", actualTime),
				}).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update workorder allocation",
					Err:        err,
				}
			}

			// Insert New Log as Pending
			err = tx.Create(&transactionworkshopentities.ServiceLog{
				CompanyId:                companyId,
				WorkOrderSystemNumber:    idSysWo,
				Frt:                      frt,
				ShiftCode:                shiftCode,
				ServiceStatusId:          utils.SrvStatPending,
				StartDatetime:            dateTimeComp, // Replace with actual value if needed
				EndDatetime:              endDatetime,
				ActualTime:               0,
				PendingTime:              0,
				EstimatedPendingTime:     estPendingTime,
				TechnicianAllocationLine: techallocLine,
			}).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to new loag as pending",
					Err:        err,
				}
			}

			// Insert New Log as Draft
			row := tx.Model(&transactionworkshopentities.ServiceLog{}).
				Where("technician_allocation_system_number = ?", idAlloc).
				Select("ISNULL(MAX(technician_allocation_line), 0) + 1").
				Row()

			err = row.Scan(&techallocLine)
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to get next technician allocation line",
					Err:        err,
				}
			}

			// Insert new log as draft
			err = tx.Create(&transactionworkshopentities.ServiceLog{
				CompanyId:                companyId,
				WorkOrderSystemNumber:    woSysNo,
				Frt:                      frt,
				ShiftCode:                shiftCode,
				ServiceStatusId:          utils.SrvStatDraft,
				StartDatetime:            endDatetime,
				EndDatetime:              draftEndTime,
				ActualTime:               0,
				PendingTime:              0,
				EstimatedPendingTime:     0,
				TechnicianAllocationLine: techallocLine,
			}).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to insert new log as draft",
					Err:        err,
				}
			}

			// Update WorkOrderServiceTime
			err = tx.Model(&transactionworkshopentities.WorkOrderServiceTime{}).
				Where("work_order_system_number = ? AND operation_item_code = ?", woSysNo, oprItemCode).
				Updates(map[string]interface{}{
					"end_datetime": dateTimeComp,
				}).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update work order summary service time",
					Err:        err,
				}
			}

		}

		if pcBr != "00003" {

			var serviceLogs []transactionworkshopentities.ServiceLog
			err := tx.Model(&transactionworkshopentities.ServiceLog{}).
				Select("technician_allocation_system_number, start_datetime, end_datetime").
				Where("technician_allocation_system_number <> ? AND company_id = ? AND technician_id = ? AND shift_code = ? AND CONVERT(VARCHAR, start_datetime, 106) = CONVERT(VARCHAR, ?, 106)",
					idAlloc, companyId, idTech, shiftCode, dateTimeComp).
				Where("EXISTS (SELECT MAX(A.TECHALLOC_LINE) FROM git  A WHERE A.technician_allocation_system_number = wtServiceLog.technician_allocation_system_number AND CONVERT(VARCHAR, A.start_datetime, 106) = CONVERT(VARCHAR, ?, 106) AND A.SERVICE_STATUS = ? AND dbo.getTimeValue(CONVERT(VARCHAR, A.START_DATETIME, 108)) < dbo.getTimeValue(CONVERT(VARCHAR, ?, 108)) HAVING wtServiceLog.TECHALLOC_LINE = MAX(A.TECHALLOC_LINE))",
					dateTimeComp, utils.SrvStatDraft, draftEndTime).
				Find(&serviceLogs).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to fetch service logs",
					Err:        err,
				}
			}

			for _, log := range serviceLogs {
				// startTime := log.StartDatetime // Implement getTimeValue function as needed
				// endTime := log.EndDatetime     // Implement getTimeValue function as needed
				// remarkAvail := fmt.Sprintf("Release INTERSECT OPERATION (Because DRAFT) from PENDING (Start Time: %s; End Time: %s)", startTime, endTime)

				// Optionally, insert into atWoTechAllocAvailable
				// err = tx.Exec(`EXEC uspg_atWoTechAllocAvailable_Insert
				//    @Option = 0,
				//    @Start_Time = ?,
				//    @End_Time = ?,
				//    @Technician_Emp_No = ?,
				//    @Company_Code = ?,
				//    @Service_Datetime = ?,
				//    @Shift_Code = ?,
				//    @Foreman_Emp_No = ?,
				//    @Ref_Type = ?,
				//    @Ref_Sys_No = ?,
				//    @Ref_Line = ?,
				//    @Remark = ?,
				//    @Creation_User_Id = ?,
				//    @Change_User_Id = ?`,
				//    startTime, endTime, technicianId, companyCode, log.StartDatetime, shiftCode, foremanEmpNo, refTypeAvail, log.TechallocSysNo, 0, remarkAvail, creationUserId, changeUserId).Error
				// if err != nil {
				//    return fmt.Errorf("failed to insert into atWoTechAllocAvailable: %w", err)
				// }

				// Delete records
				if err := tx.Where("technician_allocation_system_number = ?", log.TechnicianAllocationSystemNumber).Delete(&transactionworkshopentities.WorkOrderAllocation{}).Error; err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to delete from wtWorkOrderTechAlloc",
						Err:        err,
					}
				}

				if err := tx.Where("technician_allocation_system_number = ?", idAlloc).Delete(&transactionworkshopentities.ServiceLog{}).Error; err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to delete from wtServiceLog",
						Err:        err,
					}
				}
			}

		}

	}

	return true, nil
}
