package transactionbodyshoprepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionbodyshoppayloads "after-sales/api/payloads/transaction/bodyshop"
	transactionbodyshoprepository "after-sales/api/repositories/transaction/bodyshop"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ServiceBodyshopRepositoryImpl struct {
}

func OpenServiceBodyshopRepositoryImpl() transactionbodyshoprepository.ServiceBodyshopRepository {
	return &ServiceBodyshopRepositoryImpl{}
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

// Convert time to float64 hours
func getFloatTimeValue(t time.Time) float64 {
	return float64(t.Hour()) + float64(t.Minute())/60 + float64(t.Second())/3600
}

// Function to calculate actual time
func calculateActualTime(startTime, endTime time.Time) float64 {
	diffMinutes := endTime.Sub(startTime).Minutes()
	hours := diffMinutes / 60
	minutes := diffMinutes - (hours * 60)
	return hours + (minutes / 60)
}

// GetTimeZone fetches the time difference from the external API and adjusts the time accordingly
func GetTimeZone(currentDate time.Time, companyCode int) (time.Time, error) {

	apiURL := fmt.Sprintf("%scompany-reference?page=0&limit=1000&company_id=%d", config.EnvConfigs.GeneralServiceUrl, companyCode)
	fmt.Println("API URL:", apiURL)
	var timeReferences []transactionbodyshoppayloads.TimeReference

	err := utils.Get(apiURL, &timeReferences, nil)
	if err != nil {
		return time.Time{}, fmt.Errorf("error making API call to %s: %w", apiURL, err)
	}

	if len(timeReferences) == 0 {
		return time.Time{}, fmt.Errorf("no time reference data found for company with code %d", companyCode)
	}

	timeVariance := timeReferences[0].TimeDiff
	fmt.Println("Time variance:", timeVariance)
	if timeVariance < -24*60 || timeVariance > 24*60 {
		return time.Time{}, fmt.Errorf("invalid time variance %d: must be within -1440 to 1440 minutes", timeVariance)
	}

	adjustedTime := currentDate.Add(time.Minute * time.Duration(timeVariance))
	fmt.Println("Adjusted time:", adjustedTime)
	return adjustedTime, nil
}

func (r *ServiceBodyshopRepositoryImpl) GetAllByTechnicianWOBodyshop(tx *gorm.DB, idTech int, idSysWo int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entity transactionbodyshoppayloads.ServiceBodyshopRequest

	joinTable := utils.CreateJoinSelectStatement(tx, entity)
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)
	whereQuery = whereQuery.Where("technician_id = ? AND work_order_system_number = ? AND service_status_id IN (?,?,?,?)", idTech, idSysWo, utils.SrvStatDraft, utils.SrvStatStart, utils.SrvStatPending, utils.SrvStatStop)

	if err := whereQuery.Find(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Work order not found",
				Err:        err,
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
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
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Service status not found",
		}
	}

	// Fetch work order data from external API
	WorkOrderUrl := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(entity.WorkOrderSystemNumber)
	var workOrderResponses transactionbodyshoppayloads.ServiceBodyshopWoResponse
	errWorkOrder := utils.Get(WorkOrderUrl, &workOrderResponses, nil)
	if errWorkOrder != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order data from the external API",
			Err:        errWorkOrder,
		}
	}

	var totalRows int64
	totalRowsQuery := tx.Model(&transactionworkshopentities.ServiceLog{}).
		Where("service_log_system_number = ?", entity.ServiceLogSystemNumber).
		Count(&totalRows).Error
	if totalRowsQuery != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count service details",
			Err:        totalRowsQuery,
		}
	}

	var serviceDetails []transactionbodyshoppayloads.ServiceBodyshopResponse
	query := tx.Model(&transactionworkshopentities.ServiceLog{}).
		Joins("INNER JOIN trx_work_order_allocation AS WTA ON trx_service_log.technician_allocation_system_number = WTA.technician_allocation_system_number").
		Joins("LEFT JOIN dms_microservices_general_dev.dbo.mtr_service_status AS stat ON trx_service_log.service_status_id = stat.service_status_id").
		Select("trx_service_log.technician_allocation_system_number, trx_service_log.start_datetime, WTA.operation_code, trx_service_log.frt, trx_service_log.service_status_id, stat.service_status_description AS service_status_description, WTA.serv_actual_time, WTA.serv_pending_time, WTA.serv_progress_time, WTA.tech_alloc_start_date, WTA.tech_alloc_start_time, WTA.tech_alloc_end_date, WTA.tech_alloc_end_time").
		Where("trx_service_log.technician_allocation_line = (SELECT TOP 1 technician_allocation_line FROM trx_service_log A WHERE A.technician_allocation_system_number = trx_service_log.technician_allocation_system_number ORDER BY technician_allocation_line DESC)").
		Offset(pages.GetOffset()).
		Limit(pages.GetLimit())

	if err := query.Find(&serviceDetails).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get service details",
			Err:        err,
		}
	}

	mapResponses := []map[string]interface{}{}
	for _, serviceDetail := range serviceDetails {
		serviceMap := map[string]interface{}{
			"technician_allocation_system_number": serviceDetail.TechnicianAllocationSystemNumber,
			"start_datetime":                      serviceDetail.StartDatetime,
			"operation_code":                      serviceDetail.OperationItemCode,
			"frt":                                 serviceDetail.Frt,
			"service_status_id":                   serviceDetail.ServiceStatusId,
			"service_status_description":          serviceDetail.ServiceStatusDescription,
			"serv_actual_time":                    serviceDetail.ServActualTime,
			"serv_pending_time":                   serviceDetail.ServPendingTime,
			"serv_progress_time":                  serviceDetail.ServProgressTime,
			"tech_alloc_start_date":               serviceDetail.TechAllocStartDate,
			"tech_alloc_start_time":               serviceDetail.TechAllocStartTime,
			"tech_alloc_end_date":                 serviceDetail.TechAllocEndDate,
			"tech_alloc_end_time":                 serviceDetail.TechAllocEndTime,
		}
		mapResponses = append(mapResponses, serviceMap)
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(pages.GetLimit())))

	pages.Rows = mapResponses
	pages.TotalRows = totalRows
	pages.TotalPages = totalPages

	return pages, nil
}

// uspg_wtServiceLog_Insert
// --USE FOR : * INSERT NEW DATA OR UPDATE IF SERVICE STATUS IS START, PENDING OR STOP
// --USE IN MODUL :
func (r *ServiceBodyshopRepositoryImpl) StartService(tx *gorm.DB, idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse) {

	//--============================================================================================================================
	//	--Start :	Tombol untuk memulai pekerjaan
	//	--			(validasinya : teknisi bisa mengerjakan pekerjaan lain dengan catatan tidak ada status start pada kerjaan lainnya)
	//--============================================================================================================================

	// ============================ Deklarasi variabel yang dibutuhkan
	var woAlloc transactionworkshopentities.WorkOrderAllocation
	var smrWOServiceTimeExists int64
	var startDatetime, endDatetime time.Time
	var frt float64
	var oprItemCode string
	var dateTimeComp = time.Now()
	var currentTime = time.Now()
	var maxLine int
	var serviceStatusId int
	var SequenceNumber int
	var woSysNo int
	var shiftCode string
	var sumActualTime float64
	var pcBr string

	// ============================ Sesuaikan waktu dengan zona waktu perusahaan
	fmt.Println("Current time:", dateTimeComp)
	fmt.Println("Company ID:", companyId)
	dateTimeComp, err := GetTimeZone(dateTimeComp, companyId)
	if err != nil {
		fmt.Println("Error adjusting time zone:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to adjust time zone",
			Err:        err,
		}
	}

	// ============================ Periksa apakah alokasi teknisi valid
	var count int64
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Where("technician_allocation_system_number = ?", idAlloc).
		Count(&count).Error
	if err != nil {
		fmt.Println("Error counting technician allocation:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count technician allocation",
			Err:        err,
		}
	}
	if count == 0 {
		fmt.Println("Technician Allocation is not valid. ID:", idAlloc)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Technician Allocation is not valid. Please refresh your page",
			Err:        err,
		}
	}

	// ============================ Ambil data teknisi dan alokasi
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Select("technician_group_id, brand_id, profit_center_id,  technician_id, shift_code, foreman_id, operation_code, work_order_line").
		Where("technician_allocation_system_number = ?", idAlloc).
		First(&woAlloc).Error
	if err != nil {
		fmt.Println("Error retrieving technician allocation:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve technician allocation",
			Err:        err,
		}
	}

	// fetch work order from external service
	WorkOrderUrl := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(idSysWo)
	var workOrderResponses transactionbodyshoppayloads.ServiceBodyshopDetailResponse
	errWorkOrder := utils.Get(WorkOrderUrl, &workOrderResponses, nil)
	if errWorkOrder != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order data from the external API",
			Err:        errWorkOrder,
		}
	}

	// ============================ Set variabel yang diperlukan dari alokasi teknisi
	oprItemCode = woAlloc.OperationCode
	woLine := woAlloc.WorkOrderOperationItemLine
	technicianId := woAlloc.TechnicianId
	cpccode := woAlloc.ProfitCenterId
	shiftcode := woAlloc.ShiftCode
	WoDoc := workOrderResponses.WorkOrderDocumentNumber
	WoDate := workOrderResponses.WorkOrderDate
	// ============================ Ambil nilai maksimum dari TECHALLOC_LINE untuk sistem alokasi teknisi tertentu
	err = tx.Model(&transactionworkshopentities.ServiceLog{}).
		Select("COALESCE(MAX(technician_allocation_line), 0)").
		Where("technician_allocation_system_number = ?", idAlloc).
		Scan(&maxLine).Error

	if err != nil {
		fmt.Println("Error getting max line:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get maximum technician allocation line",
			Err:        err,
		}
	}

	// ============================ Tentukan nilai berikutnya dengan menambahkan 1
	nextLine := maxLine + 1
	fmt.Println("Next technician allocation line:", nextLine)

	// ============================ Periksa status layanan untuk alokasi teknisi dan baris yang relevan
	err = tx.Model(&transactionworkshopentities.ServiceLog{}).
		Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, nextLine-1).
		Pluck("service_status_id", &serviceStatusId).Error
	//fmt.Println("Service status ID:", serviceStatusId)
	if err != nil {
		//fmt.Println("Error getting service status:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get service status",
			Err:        err,
		}
	}

	// Cek jika status layanan termasuk dalam status draft atau transfer
	if serviceStatusId == utils.SrvStatDraft || serviceStatusId == utils.SrvStatTransfer {
		//fmt.Println("Service status is draft or transfer. ID:", serviceStatusId)
		// Periksa apakah log servis dengan status 'start' sudah ada
		var logExists int64
		err = tx.Model(&transactionworkshopentities.ServiceLog{}).
			Where("technician_allocation_system_number = ? AND company_id = ? AND technician_id = ? AND service_status_id = ? AND technician_allocation_line = (SELECT MAX(technician_allocation_line) FROM trx_service_log WHERE technician_allocation_system_number = ?)",
				idAlloc, companyId, technicianId, utils.SrvStatStart, idAlloc).
			Count(&logExists).Error
		if err != nil {
			//fmt.Println("Error checking service log existence:", err)
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to check service log existence",
				Err:        err,
			}
		}
		//fmt.Println("Service log exists:", logExists)

		if logExists == 0 {
			// profit center BR
			if cpccode == 00003 {
				//fmt.Println("Profit center is BR")
				var BRexists int64
				// cek allocation sysno exist
				err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
					Where("technician_allocation_system_number = ? AND service_status_id NOT IN (?,?,?,?,?) AND technician_id = ?", idAlloc, utils.SrvStatStop, utils.SrvStatTransfer, utils.SrvStatPending, utils.SrvStatDraft, utils.SrvStatQcPass, technicianId).
					Count(&BRexists).Error
				if err != nil {
					//fmt.Println("Error checking service log existence:", err)
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to check service log existence",
						Err:        err,
					}
				}
				//fmt.Println("BR exists:", BRexists)
				if BRexists > 0 {

					// --validation TECHNICIAN CANNOT START IF ALREADY START ANOTHER OPERATION
					var exists bool
					subquery := "SELECT MAX(technician_allocation_line) FROM trx_service_log WHERE technician_allocation_system_number = trx_service_log.technician_allocation_system_number"
					err := tx.Model(&transactionworkshopentities.ServiceLog{}).
						Select("1").
						Where("work_order_operation_id <> ? AND technician_id = ? AND service_status_id = ? AND technician_allocation_line = ("+subquery+")", oprItemCode, technicianId, utils.SrvStatStart).
						Limit(1).
						Scan(&exists).Error

					if err != nil {
						//fmt.Println("Error checking service log existence:", err)
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to check service log existence",
							Err:        err,
						}
					}

					if exists {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusConflict,
							Message:    "Technician cannot start if already started another operation",
						}
					}

				} else {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusConflict,
						Message:    "Previous Operation must be Stop, Transfer Or Pending",
						Err:        err,
					}
				}

				// profit center GR
			} else {
				//fmt.Println("Profit center is not BR")
				err = tx.Model(&transactionworkshopentities.ServiceLog{}).
					Where("technician_allocation_system_number = ? AND service_status_id = ?", idAlloc, utils.SrvStatDraft).
					Order("technician_allocation_line DESC").
					Limit(1).
					Select("start_datetime, end_datetime, frt").
					Scan(&struct {
						StartDatetime time.Time
						EndDatetime   time.Time
						Frt           float64
					}{
						StartDatetime: startDatetime,
						EndDatetime:   endDatetime,
						Frt:           frt,
					}).Error

				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to get service log dates",
						Err:        err,
					}
				}
				//fmt.Println("Start datetime:", startDatetime)
				//fmt.Println("End datetime:", endDatetime)
				//fmt.Println("FRT:", frt)

				// --validation TECHNICIAN CANNOT START IF ALREADY START ANOTHER OPERATION
				var exists bool
				// Define the subquery
				subquery := `
				SELECT MAX(technician_allocation_line) 
				FROM trx_service_log AS subquery
				WHERE subquery.technician_allocation_system_number = trx_service_log.technician_allocation_system_number`
				//AND subquery.sequence_number = trx_service_log.sequence_number`

				err := tx.Model(&transactionworkshopentities.ServiceLog{}).
					Select("1").
					Where(
						"operation_item_code <> ? AND technician_id = ? AND service_status_id = ? AND technician_allocation_line = ("+subquery+")",
						oprItemCode, technicianId, utils.SrvStatStart,
					).
					Limit(1).
					Scan(&exists).Error

				if err != nil {
					//fmt.Println("Error checking service log existence:", err)
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to check service log existence",
						Err:        err,
					}
				}

				if exists {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusConflict,
						Message:    "Technician cannot start if already started another operation",
						Err:        err,
					}
				}
				//--validation operation must start as specified on the allocation
				//--AND must be the FIRST LOG DRAFT
				if nextLine-1 == 1 {
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

				err = tx.Model(&transactionworkshopentities.ServiceLog{}).
					Where("technician_allocation_system_number = ? AND service_status_id = ?", idAlloc, utils.SrvStatDraft).
					Order("technician_allocation_line DESC").
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

				startRestTime, _ := getShiftStartTime(tx, companyId, shiftcode, dateTimeComp, true)
				endRestTime, _ := getShiftEndTime(tx, companyId, shiftcode, dateTimeComp, true)

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
					CompanyId:                        companyId,
					WorkOrderSystemNumber:            idSysWo,
					WorkOrderDocumentNumber:          WoDoc,
					TechnicianId:                     technicianId,
					ServiceStatusId:                  utils.SrvStatStart,
					StartDatetime:                    dateTimeComp,
					EndDatetime:                      endDatetime,
					ActualTime:                       0,
					PendingTime:                      0,
					EstimatedPendingTime:             0,
					Frt:                              frt,
					TechnicianAllocationSystemNumber: idAlloc,
					TechnicianAllocationLine:         nextLine,
					OperationItemCode:                oprItemCode,
					WorkOrderLine:                    woLine,
					ShiftCode:                        shiftcode,
					WorkOrderDate:                    WoDate,
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
				err = tx.Model(&transactionworkshopentities.WorkOrderServiceTime{}).
					Where("work_order_system_number = ? AND operation_item_code = ? ", idSysWo, oprItemCode).
					Count(&smrWOServiceTimeExists).Error

				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to check WorkOrderServiceTime existence",
						Err:        err,
					}
				}

				if smrWOServiceTimeExists > 0 {
					err = tx.Model(&transactionworkshopentities.WorkOrderServiceTime{}).
						Where("work_order_system_number = ? AND operation_item_code = ? ", idSysWo, oprItemCode).
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
						IsActive:                         true,
						TechnicianAllocationSystemNumber: idAlloc,
						TechnicianAllocationLine:         nextLine,
						WorkOrderLineId:                  woLine,
						CompanyId:                        companyId,
						WorkOrderSystemNumber:            idSysWo,
						OperationItemCode:                oprItemCode,
						StartDatetime:                    dateTimeComp, //--- :GH getdate diganti timecomp
						EndDatetime:                      endDatetime,
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
						Where("technician_allocation_system_number <> ? AND company_id = ? AND technician_id = ?", idAlloc, companyId, technicianId).
						Where("CONVERT(VARCHAR, start_datetime, 106) = CONVERT(VARCHAR, ?, 106)", dateTimeComp).
						Where(`EXISTS (
						SELECT TOP 1 technician_allocation_line 
						FROM trx_service_log 
						WHERE technician_allocation_system_number = trx_service_log.technician_allocation_system_number 
						AND service_status_id = 1 
						AND technician_allocation_line = 1 
						AND CONVERT(VARCHAR, start_datetime, 106) = CONVERT(VARCHAR, ?, 106) 
						AND CONVERT(VARCHAR, ?, 108) <= CONVERT(VARCHAR, ?, 108)
					)`, dateTimeComp, getTimeValue("start_datetime"), getFloatTimeValue(dateTimeComp)).
						Order("technician_allocation_line DESC").
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
						Select("technician_allocation_system_number, start_datetime, end_datetime, sequence_number").
						Where("technician_allocation_system_number <> ? AND sequence_number <> ? AND company_id = ? AND technician_id = ?", idAlloc, SequenceNumber, companyId, technicianId).
						Where("CONVERT(VARCHAR, start_datetime, 106) = CONVERT(VARCHAR, ?, 106)", dateTimeComp).
						Where(`EXISTS (
						SELECT TOP 1 technician_allocation_line 
						FROM trx_service_log 
						WHERE technician_allocation_system_number = trx_service_log.technician_allocation_system_number 
						AND service_status_id = 1 
						AND technician_allocation_line = 1 
						AND CONVERT(VARCHAR, start_datetime, 106) = CONVERT(VARCHAR, ?, 106) 
						AND CONVERT(VARCHAR, ?, 108) <= CONVERT(VARCHAR, ?, 108)
					)`, dateTimeComp, getTimeValue("start_datetime"), getFloatTimeValue(dateTimeComp)).
						Order("technician_allocation_line DESC").
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

				// if cpcCode == pcBr {
				// 	// Fetch existing record to check if TECHALLOC_START_TIME is 0
				// 	var techAlloc transactionworkshopentities.WorkOrderAllocation
				// 	err := tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
				// 		Where("technician_allocation_system_number = ?", idAlloc).
				// 		Select("tech_alloc_start_date").
				// 		First(&techAlloc).Error

				// 	if err != nil {
				// 		return false, &exceptions.BaseErrorResponse{
				// 			StatusCode: http.StatusInternalServerError,
				// 			Message:    "Failed to fetch techalloc record",
				// 			Err:        err,
				// 		}
				// 	}

				// 	// If TECHALLOC_START_TIME is 0, update it
				// 	if techAlloc.TechAllocStartTime == 0 {
				// 		startTime := getTimeValue(dateTimeComp.Format("15:04:05"))
				// 		updateErr := tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
				// 			Where("technician_allocation_system_number = ?", idAlloc).
				// 			Updates(map[string]interface{}{
				// 				"tech_alloc_start_time": startTime,
				// 			}).Error

				// 		if updateErr != nil {
				// 			return false, &exceptions.BaseErrorResponse{
				// 				StatusCode: http.StatusInternalServerError,
				// 				Message:    "Failed to update techalloc start time",
				// 				Err:        updateErr,
				// 			}
				// 		}
				// 	}
				// }

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
				var existslog bool
				err = tx.Model(&transactionworkshopentities.ServiceLog{}).
					Where("technician_allocation_system_number = ? AND technician_allocation_line = ? AND service_status_id = ?", idAlloc, nextLine-2, utils.SrvStatPending).
					Select("1").
					Scan(&existslog).Error

				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to check for previous pending log entry",
						Err:        err,
					}
				}

				if existslog {
					// Calculate Pending Time
					var startTime time.Time
					err = tx.Model(&transactionworkshopentities.ServiceLog{}).
						Select("start_datetime").
						Where("technician_allocation_system_number = ? AND technician_allocation_line = (SELECT MAX(technician_allocation_line) FROM trx_service_log WHERE technician_allocation_system_number = ? AND technician_allocation_line < ? AND service_status_id = ?)",
							idAlloc, idAlloc, nextLine, utils.SrvStatPending).
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
							idAlloc, idAlloc, nextLine, utils.SrvStatPending).
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

				// //--==============UPDATE CONTRACT SERVICE=========
				// var contractServSysNo int
				// var packageCodeCs int

				// // Fetch values from wtWorkOrder0 and wtWorkOrder2 tables
				// err = tx.Model(&transactionworkshopentities.WorkOrder{}).
				// 	Select("ISNULL(WO.contract_service_system_number, 0) AS contract_service_system_number, ISNULL(WO2.package_id, 0) AS package_id").
				// 	Joins("INNER JOIN trx_work_order_detail WO2 ON WO.work_order_system_number = WO2.work_order_system_number").
				// 	Where("WO.work_order_system_number = ? AND WO2.work_order_operation_item_line = ? AND WO2.BILL_CODE = ?", woSysNo, woLine, "S").
				// 	Scan(&struct {
				// 		ContractServSysNo int `gorm:"column:contract_service_system_number"`
				// 		PackageCodeCs     int `gorm:"column:package_id"`
				// 	}{
				// 		ContractServSysNo: contractServSysNo,
				// 		PackageCodeCs:     packageCodeCs,
				// 	}).Error

				// if err != nil {
				// 	return false, &exceptions.BaseErrorResponse{
				// 		StatusCode: http.StatusInternalServerError,
				// 		Message:    "Failed to fetch contract service details",
				// 		Err:        err,
				// 	}
				// }

				// // Check conditions and perform update if necessary
				// if contractServSysNo != 0 {
				// 	var exists bool
				// 	err = tx.Model(&transactionworkshopentities.ContractServiceItemDetail{}).
				// 		Where("contract_service_system_number = ? AND package_id = ? AND ISNULL(total_use_frt_quantity, 0) = 0", contractServSysNo, packageCodeCs).
				// 		Select("1").
				// 		Scan(&exists).Error

				// 	if err != nil {
				// 		return false, &exceptions.BaseErrorResponse{
				// 			StatusCode: http.StatusInternalServerError,
				// 			Message:    "Failed to check if contract service exists",
				// 			Err:        err,
				// 		}
				// 	}

				// 	if exists {
				// 		// Perform the update
				// 		updateErr := tx.Model(&transactionworkshopentities.ContractServiceItemDetail{}).
				// 			Where("contract_service_system_number = ? AND package_id = ? AND item_id = ?", contractServSysNo, packageCodeCs, oprItemCode).
				// 			Updates(map[string]interface{}{
				// 				"total_use_frt_quantity": gorm.Expr("frt_quantity"),
				// 			}).Error

				// 		if updateErr != nil {
				// 			return false, &exceptions.BaseErrorResponse{
				// 				StatusCode: http.StatusInternalServerError,
				// 				Message:    "Failed to update contract service",
				// 				Err:        updateErr,
				// 			}
				// 		}
				// 	}
				// }

			}

		} else {

			var woDocNo string
			var errMessage string

			// Fetch Wo_Doc_No
			err := tx.Model(&transactionworkshopentities.ServiceLog{}).
				Select("A.work_order_document_number").
				Joins("INNER JOIN trx_service_log B ON A.technician_allocation_system_number = B.technician_allocation_system_number").
				Where("A.technician_allocation_system_number = ? AND A.technician_id = ? AND A.service_status_id = ? AND A.technician_allocation_line = (SELECT MAX(B.technician_allocation_line) FROM trx_service_log B WHERE B.technician_allocation_system_number = A.technician_allocation_system_number)",
					idAlloc, technicianId, utils.SrvStatStart).
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

	} else {

		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Service not in Draft or Transfer",
			Err:        err,
		}
	}

	return true, nil
}

// uspg_wtServiceLog_Insert
// --USE FOR : * INSERT NEW DATA OR UPDATE IF SERVICE STATUS IS START, PENDING OR STOP
// --USE IN MODUL :
func (r *ServiceBodyshopRepositoryImpl) PendingService(tx *gorm.DB, idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse) {

	// ============================ Deklarasi variabel yang dibutuhkan
	var woAlloc transactionworkshopentities.WorkOrderAllocation
	var dateTimeComp = time.Now()
	var maxLine int
	var serviceStatusId int
	var woLine int
	var shiftcode string
	var WoDoc string
	var WoDate time.Time
	var frt float64
	var startDatetime time.Time
	var oprItemCode string
	var cpccode int

	// ============================ Sesuaikan waktu dengan zona waktu perusahaan
	dateTimeComp, err := GetTimeZone(dateTimeComp, companyId)
	if err != nil {
		//fmt.Println("Error adjusting time zone:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to adjust time zone",
			Err:        err,
		}
	}

	// ============================ Periksa apakah alokasi teknisi valid
	var count int64
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Where("technician_allocation_system_number = ?", idAlloc).
		Count(&count).Error
	if err != nil {
		//fmt.Println("Error counting technician allocation:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count technician allocation",
			Err:        err,
		}
	}
	if count == 0 {
		//fmt.Println("Technician Allocation is not valid. ID:", idAlloc)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Technician Allocation is not valid. Please refresh your page",
			Err:        err,
		}
	}

	// ============================ Ambil data teknisi dan alokasi
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Select("technician_group_id, brand_id, profit_center_id,  technician_id, shift_code, foreman_id, operation_code, work_order_line").
		Where("technician_allocation_system_number = ?", idAlloc).
		First(&woAlloc).Error
	if err != nil {
		//fmt.Println("Error retrieving technician allocation:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve technician allocation",
			Err:        err,
		}
	}

	// fetch work order from external service
	WorkOrderUrl := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(idSysWo)
	var workOrderResponses transactionbodyshoppayloads.ServiceBodyshopDetailResponse
	errWorkOrder := utils.Get(WorkOrderUrl, &workOrderResponses, nil)
	if errWorkOrder != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order data from the external API",
			Err:        errWorkOrder,
		}
	}

	// ============================ Set variabel yang diperlukan dari alokasi teknisi
	//oprItemCode = woAlloc.OperationCode
	woLine = woAlloc.WorkOrderOperationItemLine
	//cpccode = woAlloc.ProfitCenterId
	shiftcode = woAlloc.ShiftCode
	technicianId := woAlloc.TechnicianId
	WoDoc = workOrderResponses.WorkOrderDocumentNumber

	// ============================ Ambil nilai maksimum dari TECHALLOC_LINE untuk sistem alokasi teknisi tertentu
	err = tx.Model(&transactionworkshopentities.ServiceLog{}).
		Select("COALESCE(MAX(technician_allocation_line), 0)").
		Where("technician_allocation_system_number = ?", idAlloc).
		Scan(&maxLine).Error

	if err != nil {
		//fmt.Println("Error getting max line:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get maximum technician allocation line",
			Err:        err,
		}
	}

	// ============================ Tentukan nilai berikutnya dengan menambahkan 1
	nextLine := maxLine + 1
	//fmt.Println("Next technician allocation line:", nextLine)

	// ============================ Periksa status layanan untuk alokasi teknisi dan baris yang relevan
	err = tx.Model(&transactionworkshopentities.ServiceLog{}).
		Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, nextLine-1).
		Pluck("service_status_id", &serviceStatusId).Error
	//fmt.Println("Service status ID:", serviceStatusId)
	if err != nil {
		//fmt.Println("Error getting service status:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get service status",
			Err:        err,
		}
	}

	if serviceStatusId == utils.SrvStatStart {
		//fmt.Println("Service status is already START")
		EstPendingTime := 0.0
		if EstPendingTime == 0 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Estimate Pending Time must be filled",
				Err:        errors.New("estimate pending time must be filled"),
			}
		}

		err = tx.Model(&transactionworkshopentities.ServiceLog{}).
			Select("work_order_system_number, work_order_line, shift_code, frt, work_order_date, work_order_document_number, start_datetime, "+
				"CONVERT(VARCHAR, start_datetime, 108) AS start_time, "+
				"CONVERT(VARCHAR, end_datetime, 108) AS end_time").
			Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, nextLine-1).
			Scan(&struct {
				WoSysNo       int
				WoLine        int
				ShiftCode     string
				Frt           float64
				WoDate        time.Time
				WoDocNo       string
				StartDatetime time.Time
			}{
				WoSysNo:       idSysWo,
				WoLine:        woLine,
				ShiftCode:     shiftcode,
				Frt:           frt,
				WoDate:        WoDate,
				WoDocNo:       WoDoc,
				StartDatetime: startDatetime,
			}).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve service log data",
				Err:        err,
			}
		}

		if cpccode != 00003 {
			//--VALIDATE DAY
			//--IF CONVERT(VARCHAR,@Getdate,106) <> CONVERT(VARCHAR,@Start_Datetime,106) --- :GH getdate diganti timecomp
			dateCompFormatted := dateTimeComp.Format("02 Jan 2006")   // 106 format in Go
			startDateFormatted := startDatetime.Format("02 Jan 2006") // 106 format in Go

			if dateCompFormatted != startDateFormatted {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Message:    "Operation must be pending on . Please contact your foreman for re-allocation.",
					Err:        errors.New("operation must be pending on . please contact your foreman for re-allocation"),
				}
			}

			//--DEALLOCATE OLD ALLOCATION--
			//Remark_Avail := "Release OLD TIME from PENDING (Start Time: " + startTime + "; End Time: " + endTime + ")"
			//EXEC uspg_atWoTechAllocAvailable_Insert
			//--@Option = 0,
			//Insert record for the remark before moving
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
			// 	Remark:                Remark_Avail,
			// }

			// insertBeforeErr := tx.Create(&newAllocationBefore).Error
			// if insertBeforeErr != nil {
			// 	return false, &exceptions.BaseErrorResponse{
			// 		StatusCode: http.StatusInternalServerError,
			// 		Message:    "Failed to insert new allocation before",
			// 		Err:        insertBeforeErr,
			// 	}
			// }

		}

		// Calculate ACTUAL_TIME
		var actualTime float64
		var startTime time.Time
		var dateTimeComp time.Time // This should be set to the appropriate time value
		var estPendingTime float64 // This should be set to the estimated pending time in hours
		var endDatetime time.Time  // This should be set to the end time of the service
		var draftEndTime time.Time
		var startRestTime time.Time
		var endRestTime time.Time

		pendingDuration := time.Duration(estPendingTime) * time.Hour
		endTime := dateTimeComp.Add(pendingDuration)
		currentTime := time.Now()

		// --IF dbo.getTimeValue(CONVERT(VARCHAR,@Getdate,108)) + @Est_Pending_Time >= 24 --- :GH getdate diganti timecomp
		// Check if the end time exceeds the current day
		if endTime.Sub(currentTime).Hours() >= 24 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Service pending cannot exceed today. Please use transfer for service on a different date",
				Err:        errors.New("service pending cannot exceed today"),
			}
		}

		// Fetch the START_DATETIME from the database
		err := tx.Model(&transactionworkshopentities.ServiceLog{}).
			Select("start_datetime").
			Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, nextLine-1).
			Take(&startTime).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve start time",
				Err:        err,
			}
		}

		// Calculate ACTUAL_TIME
		actualTime = calculateActualTime(startTime, dateTimeComp)

		// Update the record
		err = tx.Model(&transactionworkshopentities.ServiceLog{}).
			Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, nextLine-1).
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

		// Update wtSmrWOServiceTime
		err = tx.Model(&transactionworkshopentities.WorkOrderServiceTime{}).
			Where("work_order_system_number = ? AND operation_item_code = ?", idSysWo, oprItemCode).
			Updates(map[string]interface{}{
				"end_datetime": dateTimeComp,
			}).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update wtSmrWOServiceTime",
				Err:        err,
			}
		}

		// Calculate @SumActualTime and new End Time for this DRAFT LOG
		var sumActualTime float64
		err = tx.Model(&transactionworkshopentities.ServiceLog{}).
			Select("SUM(actual_time)").
			Where("technician_allocation_system_number = ?", idAlloc).
			Pluck("SUM(actual_time)", &sumActualTime).Error

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

		// Convert endDatetime to a string formatted as "15:04:05"
		endDatetimeStr := endDatetime.Format("15:04:05")

		// Calculate draftEndTime
		endTimeValue := getTimeValue(endDatetimeStr) + (frt - sumActualTime)
		if endTimeValue >= 24 {
			draftEndTime = endDatetime.Truncate(24 * time.Hour).Add(23*time.Hour + 45*time.Minute)
		} else {
			draftEndTime = endDatetime.Add(time.Duration((frt - sumActualTime) * float64(time.Hour)))
		}

		// Determine rest times
		startRestTimeFloat, _ := getShiftStartTime(tx, companyId, shiftcode, dateTimeComp, true)
		endRestTimeFloat, _ := getShiftEndTime(tx, companyId, shiftcode, dateTimeComp, true)

		// Convert float64 times to time.Time
		startRestTime = time.Date(dateTimeComp.Year(), dateTimeComp.Month(), dateTimeComp.Day(), int(startRestTimeFloat), int((startRestTimeFloat-float64(int(startRestTimeFloat)))*60), 0, 0, dateTimeComp.Location())
		endRestTime = time.Date(dateTimeComp.Year(), dateTimeComp.Month(), dateTimeComp.Day(), int(endRestTimeFloat), int((endRestTimeFloat-float64(int(endRestTimeFloat)))*60), 0, 0, dateTimeComp.Location())

		// Recalculate draftEndTime if needed
		if draftEndTime.Before(startRestTime) {
			if draftEndTime.Before(endRestTime) {
				draftEndTime = endDatetime.Add(time.Duration((frt - sumActualTime) * float64(time.Hour)))
			}
		}

		if cpccode != 00003 {
			// DECLARE @CSR2 CURSOR
			// Fetch the service logs from the database
			var serviceLogs []transactionworkshopentities.ServiceLog
			err := tx.Model(&transactionworkshopentities.ServiceLog{}).
				Select("technician_allocation_system_number, start_datetime, end_datetime").
				Where("technician_allocation_system_number <> ? AND company_id = ? AND technician_id = ? AND shift_code = ?", idAlloc, companyId, technicianId, shiftcode).
				Where("CONVERT(VARCHAR, start_datetime, 106) = CONVERT(VARCHAR, ?, 106)", dateTimeComp).
				Where(`EXISTS (
							SELECT 1 
					FROM trx_service_log AS A 
					WHERE A.technician_allocation_system_number = trx_service_log.technician_allocation_system_number 
					AND A.service_status_id = 1 
					AND A.technician_allocation_line = 1 
					AND CONVERT(VARCHAR, A.start_datetime, 106) = CONVERT(VARCHAR, ?, 106)
					AND CONVERT(VARCHAR, A.start_datetime, 108) < CONVERT(VARCHAR, ?, 108)
				)`, dateTimeComp, draftEndTime).
				Order("technician_allocation_line DESC").
				Find(&serviceLogs).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve service logs",
					Err:        err,
				}
			}

			// Process the retrieved service logs
			filteredLogs := []transactionworkshopentities.ServiceLog{}
			for _, log := range serviceLogs {
				startTimeValue := getTimeValue(log.StartDatetime.Format("15:04:05"))
				endTimeValue := getTimeValue(log.EndDatetime.Format("15:04:05"))
				queryTimeValue := getTimeValue(draftEndTime.Format("15:04:05"))

				if startTimeValue <= queryTimeValue && endTimeValue >= queryTimeValue {
					filteredLogs = append(filteredLogs, log)
				}
			}

			// Delete the filtered logs and associated records
			for _, log := range filteredLogs {
				// Delete from WorkOrderAllocation
				err := tx.Where("technician_allocation_system_number = ?", log.TechnicianAllocationSystemNumber).Delete(&transactionworkshopentities.WorkOrderAllocation{}).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to delete from WorkOrderAllocation",
						Err:        err,
					}
				}

				// Delete from ServiceLog
				err = tx.Where("technician_allocation_system_number = ?", log.TechnicianAllocationSystemNumber).Delete(&transactionworkshopentities.ServiceLog{}).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to delete from ServiceLog",
						Err:        err,
					}
				}
			}

			// DECLARE @CSR_2 CURSOR
			// FETCH NEXT FROM @CSR_2 INTO @TechAllocSysNo, @StartDateTime, @EndDateTime
			// WHILE @@FETCH_STATUS = 0
			var serviceLogs_2 []transactionworkshopentities.ServiceLog

			err = tx.Model(&transactionworkshopentities.ServiceLog{}).
				Select("technician_allocation_system_number, start_datetime, end_datetime").
				Where("technician_allocation_system_number <> ? AND company_id = ? AND technician_id = ? AND shift_code = ?", idAlloc, companyId, technicianId, shiftcode).
				Where("CONVERT(VARCHAR, start_datetime, 106) = CONVERT(VARCHAR, ?, 106)", dateTimeComp).
				Where(`EXISTS (
        SELECT 1
        FROM trx_service_log AS A
        WHERE A.technician_allocation_system_number = trx_service_log.technician_allocation_system_number
        AND A.service_status_id = 1
        AND A.technician_allocation_line = 1
        AND CONVERT(VARCHAR, A.start_datetime, 106) = CONVERT(VARCHAR, ?, 106)
        AND CONVERT(VARCHAR, A.start_datetime, 108) < CONVERT(VARCHAR, ?, 108)
    )`, dateTimeComp, draftEndTime).
				Order("technician_allocation_line DESC").
				Find(&serviceLogs_2).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve service logs",
					Err:        err,
				}
			}

			// Process the retrieved service logs
			for _, log := range serviceLogs_2 {
				startTimeValue := getTimeValue(log.StartDatetime.Format("15:04:05"))
				endTimeValue := getTimeValue(log.EndDatetime.Format("15:04:05"))
				draftEndTimeValue := getTimeValue(draftEndTime.Format("15:04:05"))

				// Calculate the difference in time
				diffTime := (draftEndTimeValue - startTimeValue) // difference in hours

				// Update ServiceLog --Move the Service LOG
				err = tx.Model(&transactionworkshopentities.ServiceLog{}).
					Where("technician_allocation_system_number = ? AND service_status_id = ?", log.TechnicianAllocationSystemNumber, utils.SrvStatDraft).
					Updates(map[string]interface{}{
						"start_datetime": gorm.Expr("CONVERT(VARCHAR, start_datetime, 106) + ' ' + dbo.getTime(? + ?)", startTimeValue, diffTime),
						"end_datetime":   gorm.Expr("CONVERT(VARCHAR, end_datetime, 106) + ' ' + dbo.getTime(? + ?)", endTimeValue, diffTime),
					}).Error

				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to update ServiceLog",
						Err:        err,
					}
				}

				// Update WorkOrderAllocation --Move the Allocation
				err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
					Where("technician_allocation_system_number = ?", log.TechnicianAllocationSystemNumber).
					Updates(map[string]interface{}{
						"techalloc_start_time": gorm.Expr("techalloc_start_time + ?", diffTime),
						"techalloc_end_time":   gorm.Expr("techalloc_end_time + ?", diffTime),
					}).Error

				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to update WorkOrderAllocation",
						Err:        err,
					}
				}
			}

		}

		// --ALLOCATE THE NEW TIME
		if cpccode != 00003 {
			var startTimeValue float64
			err = tx.Model(&transactionworkshopentities.ServiceLog{}).
				Select("dbo.getTimeValue(CONVERT(VARCHAR, start_datetime, 108))").
				Where("techalloc_sys_no = ? AND techalloc_line = ?", idAlloc, nextLine-1).
				Pluck("dbo.getTimeValue(CONVERT(VARCHAR, start_datetime, 108))", &startTimeValue).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve start time",
					Err:        err,
				}
			}

			// endTimeValue := getTimeValue(dateTimeComp.Format("15:04:05"))
			//remarkAvail := fmt.Sprintf("Allocate NEW TIME from PENDING (Start Time: %v; End Time: %v)", startTimeValue, endTimeValue)

			// Insert into atWoTechAllocAvailable
			// err = tx.Create(&transactionworkshopentities.WorkOrderAllocation{
			// 	TechAllocStartDate:       startTimeValue,
			// 	TechAllocEndTime:         endTimeValue,
			// 	TechnicianId: technicianId,
			// 	CompanyId:     companyCode,
			// 	: startDateTime,
			// 	ShiftCode:       shiftCode,
			// 	ForemanEmpNo:    foremanEmpNo,
			// 	RefType:         refTypeAvail,
			// 	RefSysNo:        techallocSysNo,
			// 	RefLine:         tempAllocLine,
			// 	Remark:          remarkAvail,
			// 	CreationUserId:  creationUserId,
			// 	ChangeUserId:    changeUserId,
			// }).Error

			// if err != nil {
			// 	return false, &exceptions.BaseErrorResponse{
			// 		StatusCode: http.StatusInternalServerError,
			// 		Message:    "Failed to insert into atWoTechAllocAvailable",
			// 		Err:        err,
			// 	}
			// }

			// Update wtWorkOrder2 status
			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Where("work_order_system_number = ? AND work_order_operation_item_line = ?", idSysWo, woLine).
				Updates(map[string]interface{}{
					"work_order_status_id": utils.SrvStatPending,
				}).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update wtWorkOrder2 status",
					Err:        err,
				}
			}
		}

		// Update atWoTechAlloc Service Status
		err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
			Where("technician_allocation_system_number = ?", idAlloc).
			Updates(map[string]interface{}{
				"serv_status":        utils.SrvStatPending,
				"serv_actual_time":   gorm.Expr("serv_actual_time + ?", actualTime),
				"serv_progress_time": gorm.Expr("serv_progress_time + ?", actualTime),
			}).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update atWoTechAlloc",
				Err:        err,
			}
		}

		// Insert New Log as Pending
		serviceLogPending := transactionworkshopentities.ServiceLog{
			CompanyId:                        companyId,
			WorkOrderSystemNumber:            idSysWo,
			WorkOrderDocumentNumber:          WoDoc,
			WorkOrderLine:                    woLine,
			OperationItemCode:                oprItemCode,
			TechnicianId:                     technicianId,
			Frt:                              frt,
			WorkOrderDate:                    WoDate.Format("2006-01-02 15:04:05"),
			ShiftCode:                        shiftcode,
			ServiceStatusId:                  utils.SrvStatPending,
			StartDatetime:                    dateTimeComp, // Replacing @Getdate with dateTimeComp
			EndDatetime:                      endDatetime,
			ActualTime:                       0,
			PendingTime:                      0,
			EstimatedPendingTime:             estPendingTime,
			TechnicianAllocationSystemNumber: idAlloc,
			TechnicianAllocationLine:         nextLine,
		}
		err = tx.Create(&serviceLogPending).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to insert new log as pending",
				Err:        err,
			}
		}

		// Insert New Log as Draft
		var maxTechallocLine int
		err = tx.Model(&transactionworkshopentities.ServiceLog{}).
			Select("ISNULL(MAX(technician_allocation_line), 0) + 1").
			Where("technician_allocation_system_number = ?", idAlloc).
			Pluck("ISNULL(MAX(technician_allocation_line), 0) + 1", &maxTechallocLine).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to set TechallocLine",
				Err:        err,
			}
		}

		serviceLogDraft := transactionworkshopentities.ServiceLog{
			CompanyId:                        companyId,
			WorkOrderSystemNumber:            idSysWo,
			WorkOrderDocumentNumber:          WoDoc,
			WorkOrderLine:                    woLine,
			OperationItemCode:                oprItemCode,
			TechnicianId:                     technicianId,
			Frt:                              frt,
			WorkOrderDate:                    WoDate.Format("2006-01-02 15:04:05"),
			ShiftCode:                        shiftcode,
			ServiceStatusId:                  utils.SrvStatDraft,
			StartDatetime:                    endDatetime,
			EndDatetime:                      draftEndTime,
			ActualTime:                       0,
			PendingTime:                      0,
			EstimatedPendingTime:             0,
			TechnicianAllocationSystemNumber: idAlloc,
			TechnicianAllocationLine:         maxTechallocLine,
		}
		err = tx.Create(&serviceLogDraft).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to insert new log as draft",
				Err:        err,
			}
		}

		// Update wtSmrWOServiceTime
		err = tx.Model(&transactionworkshopentities.WorkOrderServiceTime{}).
			Where("work_order_system_number = ? AND operation_item_code = ?", idSysWo, oprItemCode).
			Updates(map[string]interface{}{
				"end_datetime": dateTimeComp, // Replacing @Getdate with dateTimeComp
			}).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update wtSmrWOServiceTime",
				Err:        err,
			}
		}

	} else {
		fmt.Println("Service status is already PENDING")
	}

	return true, nil
}

// uspg_wtServiceLog_Insert
// --USE FOR : * INSERT NEW DATA OR UPDATE IF SERVICE STATUS IS START, PENDING OR STOP
// --USE IN MODUL :
func (r *ServiceBodyshopRepositoryImpl) TransferService(tx *gorm.DB, idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse) {

	// ============================ Deklarasi variabel yang dibutuhkan
	var woAlloc transactionworkshopentities.WorkOrderAllocation
	var dateTimeComp = time.Now()
	var maxLine int
	var serviceStatusId int
	var oprItemCode string

	// ============================ Sesuaikan waktu dengan zona waktu perusahaan
	dateTimeComp, err := GetTimeZone(dateTimeComp, companyId)
	if err != nil {
		//fmt.Println("Error adjusting time zone:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to adjust time zone",
			Err:        err,
		}
	}

	// ============================ Periksa apakah alokasi teknisi valid
	var counts int64
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Where("technician_allocation_system_number = ?", idAlloc).
		Count(&counts).Error
	if err != nil {
		//fmt.Println("Error counting technician allocation:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count technician allocation",
			Err:        err,
		}
	}
	if counts == 0 {
		//fmt.Println("Technician Allocation is not valid. ID:", idAlloc)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Technician Allocation is not valid. Please refresh your page",
			Err:        err,
		}
	}

	// ============================ Ambil data teknisi dan alokasi
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Select("technician_group_id, brand_id, profit_center_id, technician_id, shift_code, foreman_id, operation_code, work_order_line").
		Where("technician_allocation_system_number = ?", idAlloc).
		First(&woAlloc).Error
	if err != nil {
		// Handle error and return appropriate response
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve technician allocation",
			Err:        err,
		}
	}

	// Fetch work order from external service
	WorkOrderUrl := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(idSysWo)
	var workOrderResponses transactionbodyshoppayloads.ServiceBodyshopDetailResponse
	errWorkOrder := utils.Get(WorkOrderUrl, &workOrderResponses, nil)
	if errWorkOrder != nil {
		// Handle error and return appropriate response
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order data from the external API",
			Err:        errWorkOrder,
		}
	}

	// ============================ Set variabel yang diperlukan dari alokasi teknisi
	oprItemCode = woAlloc.OperationCode
	woLine := woAlloc.WorkOrderOperationItemLine
	cpccode := woAlloc.ProfitCenterId
	shiftcode := woAlloc.ShiftCode
	technicianId := woAlloc.TechnicianId
	WoDoc := workOrderResponses.WorkOrderDocumentNumber
	WoDate := workOrderResponses.WorkOrderDate

	// ============================ Ambil nilai maksimum dari TECHALLOC_LINE untuk sistem alokasi teknisi tertentu
	err = tx.Model(&transactionworkshopentities.ServiceLog{}).
		Select("COALESCE(MAX(technician_allocation_line), 0)").
		Where("technician_allocation_system_number = ?", idAlloc).
		Scan(&maxLine).Error

	if err != nil {
		//fmt.Println("Error getting max line:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get maximum technician allocation line",
			Err:        err,
		}
	}

	// ============================ Tentukan nilai berikutnya dengan menambahkan 1
	nextLine := maxLine + 1
	//fmt.Println("Next technician allocation line:", nextLine)

	// ============================ Periksa status layanan untuk alokasi teknisi dan baris yang relevan
	err = tx.Model(&transactionworkshopentities.ServiceLog{}).
		Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, nextLine-1).
		Pluck("service_status_id", &serviceStatusId).Error
	//fmt.Println("Service status ID:", serviceStatusId)
	if err != nil {
		//fmt.Println("Error getting service status:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get service status",
			Err:        err,
		}
	}

	var count int64
	// Query to check existence
	resultcheck := tx.Model(&transactionworkshopentities.ServiceLog{}).
		Where("work_order_system_number = ? AND technician_allocation_system_number = ? AND technician_allocation_line = ? AND service_status_id IN (?,?)",
			idSysWo, idAlloc, nextLine-1, utils.SrvStatStart, utils.SrvStatDraft).
		Count(&count).Error
	if resultcheck != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check service log existence",
			Err:        resultcheck,
		}
	}

	if count > 0 {
		var status int
		// Query to get the service status
		resultcheck := tx.Model(&transactionworkshopentities.ServiceLog{}).
			Where("work_order_system_number = ? AND technician_allocation_system_number = ? AND technician_allocation_line = ?",
				idSysWo, idAlloc, nextLine-1).
			Pluck("service_status_id", &status).Error
		if resultcheck != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to check service status",
				Err:        resultcheck,
			}
		}

		if status == utils.SrvStatDraft {
			type ServiceLogDetails struct {
				WorkOrderSystemNumber string    `gorm:"column:work_order_system_order"`
				WorkOrderLine         int       `gorm:"column:work_order_line"`
				ShiftCode             string    `gorm:"column:shift_code"`
				FRT                   string    `gorm:"column:frt"`
				WorkOrderDate         time.Time `gorm:"column:work_order_date"`
				WorkOrderDocNo        string    `gorm:"column:work_order_document_number"`
				EndDatetime           time.Time `gorm:"column:start_datetime"`
				StartTime             float64   // Assuming you will handle the conversion separately
				EndTime               float64   // Assuming you will handle the conversion separately
			}
			var details ServiceLogDetails

			// Query to get the service log details
			resultcheck := tx.Model(&transactionworkshopentities.ServiceLog{}).
				Select("work_order_system_order, work_order_line, shift_code, frt, work_order_date, work_order_document_number, start_datetime").
				Where("technician_allocation_system_number = ? AND technician_allocation_line = ? AND service_status_id = ?",
					idAlloc, nextLine-1, utils.SrvStatDraft).
				First(&details).Error

			if resultcheck != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve service log details",
					Err:        resultcheck,
				}
			}

			// Convert time fields to float64 values
			// details.StartTime = getTimeValue(details.EndDatetime)
			// details.EndTime = getTimeValue(details.EndDatetime)

			if cpccode != 00003 {
				var techallocLine int
				tempAllocLine := techallocLine - 1

				// Perform the delete operation
				result := tx.Model(&transactionworkshopentities.ServiceLog{}).
					Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, tempAllocLine).
					Delete(&transactionworkshopentities.ServiceLog{})

				if result.Error != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to delete service log",
						Err:        result.Error,
					}
				}
			}

		} else if status == utils.SrvStatPending {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "1.	Pada menu Service – untuk status jasa yang pending harusnya tidak bisa dialokasi ulang kepada teknisi lain (ganti teknisi (harus hanya bisa di start ulang pada teknisi yang sama)",
				Err:        errors.New("service is already pending"),
			}
		} else {
			endDateTime := dateTimeComp

			err = tx.Model(&transactionworkshopentities.ServiceLog{}).
				Select("end_datetime").
				Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, nextLine-1).
				Pluck("end_datetime", &endDateTime).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve end time",
					Err:        err,
				}
			}

			//--==UPDATE ACTUAL TIME IN THE LOG BEFORE (START)==--
			var serviceLog transactionworkshopentities.ServiceLog

			// Fetch the start datetime and other necessary data
			err := tx.Model(&transactionworkshopentities.ServiceLog{}).
				Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, nextLine-1).
				Select("start_datetime"). // Add other fields if necessary
				First(&serviceLog).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to fetch service log",
					Err:        err,
				}
			}

			// Calculate ACTUAL_TIME
			startDatetime := serviceLog.StartDatetime
			currentDatetime := time.Now()

			// Calculate ACTUAL_TIME
			diffMinutes := currentDatetime.Sub(startDatetime).Minutes()
			actualTime := diffMinutes / 60

			// Update the ServiceLog record
			err = tx.Model(&transactionworkshopentities.ServiceLog{}).
				Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, nextLine-1).
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

			// Update the WorkOrderServiceTime record
			err = tx.Model(&transactionworkshopentities.WorkOrderServiceTime{}).
				Where("work_order_system_number = ? AND operation_item_code = ?", idSysWo, oprItemCode).
				Updates(map[string]interface{}{
					"end_datetime": dateTimeComp,
				}).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update work order service time",
					Err:        err,
				}
			}

			//--Update WorkOrder2 Status
			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Where("work_order_system_number = ? AND work_order_operation_item_line = ?", idSysWo, woLine).
				Updates(map[string]interface{}{
					"service_status_id": utils.SrvStatTransfer,
				}).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update work order status",
					Err:        err,
				}
			}

			//--==Update atWoTechAlloc Service Status==--
			err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
				Where("technician_allocation_system_number = ?", idAlloc).
				Updates(map[string]interface{}{
					"service_status_id":  utils.SrvStatTransfer,
					"serv_actual_time":   gorm.Expr("serv_actual_time + ?", actualTime),
					"serv_progress_time": gorm.Expr("serv_progress_time + ?", actualTime),
				}).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update work order allocation",
					Err:        err,
				}
			}

			var startDatetimequery time.Time

			// Query to get the top 1 START_DATETIME based on TECHALLOC_SYS_NO and ordered by TECHALLOC_LINE
			err = tx.Model(&transactionworkshopentities.ServiceLog{}).
				Select("start_datetime").
				Where("technician_allocation_system_number = ?", idAlloc).
				Order("technician_allocation_line ASC").
				Limit(1).
				Pluck("start_datetime", &startDatetimequery).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to fetch start datetime",
					Err:        err,
				}
			}

			//--==Insert New Log as Transfer==--
			serviceLogTransfer := transactionworkshopentities.ServiceLog{
				CompanyId:                        companyId,
				WorkOrderSystemNumber:            idSysWo,
				WorkOrderDocumentNumber:          WoDoc,
				WorkOrderLine:                    woLine,
				OperationItemCode:                oprItemCode,
				TechnicianId:                     technicianId,
				Frt:                              woAlloc.Frt,
				WorkOrderDate:                    WoDate,
				ShiftCode:                        shiftcode,
				ServiceStatusId:                  utils.SrvStatTransfer,
				StartDatetime:                    startDatetimequery,
				EndDatetime:                      dateTimeComp,
				ActualTime:                       actualTime,
				PendingTime:                      0,
				EstimatedPendingTime:             0,
				TechnicianAllocationSystemNumber: idAlloc,
				TechnicianAllocationLine:         nextLine,
			}

			err = tx.Create(&serviceLogTransfer).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to insert new log as transfer",
					Err:        err,
				}
			}

			// Perform the update operation
			result := tx.Model(&transactionworkshopentities.WorkOrderServiceTime{}).
				Where("work_order_system_number = ? AND operation_item_code = ?", idSysWo, oprItemCode).
				Updates(map[string]interface{}{
					"end_datetime": dateTimeComp,
				})

			if result.Error != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update work order service time",
					Err:        result.Error,
				}
			}

		}

	} else {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Service is already started or in draft status",
			Err:        errors.New("service is already started or in draft status"),
		}
	}

	return true, nil
}

// uspg_wtServiceLog_Insert
// --USE FOR : * INSERT NEW DATA OR UPDATE IF SERVICE STATUS IS START, PENDING OR STOP
// --USE IN MODUL :
func (r *ServiceBodyshopRepositoryImpl) StopService(tx *gorm.DB, idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse) {
	// ============================ Deklarasi variabel yang dibutuhkan
	var woAlloc transactionworkshopentities.WorkOrderAllocation
	var dateTimeComp = time.Now()
	var maxLine int
	var serviceStatusId int
	var oprItemCode string
	var timeValue float64

	// ============================ Sesuaikan waktu dengan zona waktu perusahaan
	dateTimeComp, err := GetTimeZone(dateTimeComp, companyId)
	if err != nil {
		//fmt.Println("Error adjusting time zone:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to adjust time zone",
			Err:        err,
		}
	}

	// ============================ Periksa apakah alokasi teknisi valid
	var count int64
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Where("technician_allocation_system_number = ?", idAlloc).
		Count(&count).Error
	if err != nil {
		//fmt.Println("Error counting technician allocation:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count technician allocation",
			Err:        err,
		}
	}
	if count == 0 {
		//fmt.Println("Technician Allocation is not valid. ID:", idAlloc)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Technician Allocation is not valid. Please refresh your page",
			Err:        err,
		}
	}

	// ============================ Ambil data teknisi dan alokasi
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Select("technician_group_id, brand_id, profit_center_id, technician_id, shift_code, foreman_id, operation_code, work_order_line").
		Where("technician_allocation_system_number = ?", idAlloc).
		First(&woAlloc).Error
	if err != nil {
		// Handle error and return appropriate response
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve technician allocation",
			Err:        err,
		}
	}

	// Fetch work order from external service
	WorkOrderUrl := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(idSysWo)
	var workOrderResponses transactionbodyshoppayloads.ServiceBodyshopDetailResponse
	errWorkOrder := utils.Get(WorkOrderUrl, &workOrderResponses, nil)
	if errWorkOrder != nil {
		// Handle error and return appropriate response
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order data from the external API",
			Err:        errWorkOrder,
		}
	}

	// ============================ Set variabel yang diperlukan dari alokasi teknisi
	oprItemCode = woAlloc.OperationCode
	// woLine := woAlloc.WorkOrderLine
	cpccode := woAlloc.ProfitCenterId
	// shiftcode := woAlloc.ShiftCode
	// technicianId := woAlloc.TechnicianId
	// WoDoc := workOrderResponses.WorkOrderDocumentNumber

	// ============================ Ambil nilai maksimum dari TECHALLOC_LINE untuk sistem alokasi teknisi tertentu
	err = tx.Model(&transactionworkshopentities.ServiceLog{}).
		Select("COALESCE(MAX(technician_allocation_line), 0)").
		Where("technician_allocation_system_number = ?", idAlloc).
		Scan(&maxLine).Error

	if err != nil {
		//fmt.Println("Error getting max line:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get maximum technician allocation line",
			Err:        err,
		}
	}

	// ============================ Tentukan nilai berikutnya dengan menambahkan 1
	nextLine := maxLine + 1
	//fmt.Println("Next technician allocation line:", nextLine)

	// ============================ Periksa status layanan untuk alokasi teknisi dan baris yang relevan
	err = tx.Model(&transactionworkshopentities.ServiceLog{}).
		Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, nextLine-1).
		Pluck("service_status_id", &serviceStatusId).Error
	//fmt.Println("Service status ID:", serviceStatusId)
	if err != nil {
		//fmt.Println("Error getting service status:", err)
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get service status",
			Err:        err,
		}
	}

	if serviceStatusId == utils.SrvStatStart {
		//fmt.Println("Service status is already START")

		type ServiceLogResult struct {
			WoSysNo       int       `gorm:"column:wo_sys_no"`
			WoLine        int       `gorm:"column:wo_line"`
			ShiftCode     string    `gorm:"column:shift_code"`
			Frt           float64   `gorm:"column:frt"`
			WoDate        time.Time `gorm:"column:wo_date"`
			WoDocNo       string    `gorm:"column:wo_doc_no"`
			StartDatetime time.Time `gorm:"column:start_datetime"`
			EndDatetime   time.Time `gorm:"column:end_datetime"`
			StartTime     float64   `gorm:"column:start_time"`
			EndTime       float64   `gorm:"column:end_time"`
		}

		// Initialize a variable to hold the query result
		var result ServiceLogResult

		// Execute the query and scan the result into the struct
		err := tx.Model(&transactionworkshopentities.ServiceLog{}).
			Select("wo_sys_no, wo_line, shift_code, frt, wo_date, wo_doc_no, start_datetime, end_datetime, "+
				"CAST(CONVERT(VARCHAR, start_datetime, 108) AS FLOAT) AS start_time, "+
				"CAST(CONVERT(VARCHAR, end_datetime, 108) AS FLOAT) AS end_time").
			Where("techalloc_sys_no = ? AND techalloc_line = ?", idAlloc, nextLine-1).
			Scan(&result).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve service log data",
				Err:        err,
			}
		}

		// Convert time string to float values for start_time and end_time
		result.StartTime = getTimeValue(result.StartDatetime.Format("15:04:05"))
		result.EndTime = getTimeValue(result.EndDatetime.Format("15:04:05"))

		// Validate CPC Code
		// if cpccode != 00003 {
		// 	// Validate Day
		// 	if datetimeComp.Format("02 Jan 2006") != startDatetime.Format("02 Jan 2006") {
		// 		errorMsg := fmt.Sprintf("Operation must be stopped on %s. Please contact your foreman for re-allocation.", startDatetime.Format("02 Jan 2006"))
		// 		return false, &exceptions.BaseErrorResponse{
		// 			StatusCode: http.StatusBadRequest,
		// 			Message:    errorMsg,
		// 		}
		// 	}

		// 	// Deallocate Old Allocation
		// 	//tempAllocLine := nextLine - 1
		// 	//remarkAvail := fmt.Sprintf("Release OLD TIME from STOP (Start Time: %.2f; End Time: %.2f)", startTime, endTime)

		// 	// Perform Deallocation (GORM query or insert)
		// 	// --EXEC uspg_atWoTechAllocAvailable_Insert
		// 	// --@Option = 0,
		// 	// err := tx.Model(&transactionworkshopentities.WorkOrderAllocationAvailable{}).
		// 	// 	Create(&transactionworkshopentities.WorkOrderAllocationAvailable{
		// 	// 		TechnicianId:          woAlloc.TechnicianId,
		// 	// 		CompanyId:             companyId,
		// 	// 		ServiceDateTime:       result.StartDatetime,
		// 	// 		ShiftCode:             woAlloc.ShiftCode,
		// 	// 		ForemanId:             woAlloc.ForemanId,
		// 	// 		ReferenceType:         refTypeAvail,
		// 	// 		ReferenceSystemNumber: techallocSysNo,
		// 	// 		ReferenceLine:         tempAllocLine,
		// 	// 		Remark:                remarkAvail,
		// 	// 	}).Error

		// 	// if err != nil {
		// 	// 	return false, &exceptions.BaseErrorResponse{
		// 	// 		StatusCode: http.StatusInternalServerError,
		// 	// 		Message:    "Failed to deallocate old allocation",
		// 	// 		Err:        err,
		// 	// 	}
		// 	// }
		// }

		// Calculate ACTUAL_TIME
		var actualTime float64
		err = tx.Model(&transactionworkshopentities.ServiceLog{}).
			Select("CAST(DATEDIFF(mi, start_datetime, ?) / 60 AS DECIMAL(7,2)) + "+
				"CAST(DATEDIFF(mi, start_datetime, ?) % 60 / 60 AS DECIMAL(7,2))",
				dateTimeComp, dateTimeComp).
			Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, nextLine-1).
			Pluck("CAST(DATEDIFF(mi, start_datetime, ?) / 60 AS DECIMAL(7,2)) + "+
				"CAST(DATEDIFF(mi, start_datetime, ?) % 60 / 60 AS DECIMAL(7,2))", &actualTime).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate actual time",
				Err:        err,
			}
		}

		// Update wtServiceLog
		err = tx.Model(&transactionworkshopentities.ServiceLog{}).
			Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, nextLine-1).
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

		// Update wtSmrWOServiceTime
		err = tx.Model(&transactionworkshopentities.WorkOrderServiceTime{}).
			Where("work_order_system_number = ? AND operation_item_code = ?", idSysWo, oprItemCode).
			Updates(map[string]interface{}{
				"end_datetime": dateTimeComp,
			}).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update SMR WO Service Time",
				Err:        err,
			}
		}

		if cpccode != 00003 {
			// Update other Allocation Sys No that intersect
			// Query to check if any intersecting allocation exists
			timeValue = getTimeValue(dateTimeComp.Format("15:04:05"))
			var exists bool
			err := tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
				Where("company_code = ?", companyId).
				Where("tech_alloc_start_date = CONVERT(VARCHAR, ? ,106)", dateTimeComp).
				Where("foreman_id = ?", woAlloc.ForemanId).
				Where("technician_id = ?", woAlloc.TechnicianId).
				Where("technician_allocation_system_number <> ?", idAlloc).
				Where("service_status_id = ?", utils.SrvStatDraft).
				Where("tech_alloc_start_time < ?", timeValue). // Assuming `getTimeValue` logic is handled separately
				Limit(1).                                      // Limit to 1 record to match EXISTS
				Find(&exists).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to check for intersecting allocation",
					Err:        err,
				}
			}

			if exists {
				// Calculate the time difference after the expansion
				endDateTime := time.Now()
				startDateTime := time.Now().Add(-2 * time.Hour) // Contoh, sesuaikan dengan logika yang sebenarnya
				diffMinutes := endDateTime.Sub(startDateTime).Minutes()
				diffTime := diffMinutes / 60

				var allocations []transactionworkshopentities.WorkOrderAllocation
				err := tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
					Where("company_id = ?", companyId).
					Where("tech_alloc_start_date = ?", dateTimeComp).
					Where("foreman_id = ?", woAlloc.ForemanId).
					Where("technician_id = ?", woAlloc.TechnicianId).
					Where("technician_allocation_system_number <> ?", idAlloc).
					Where("service_status_id = ?", utils.SrvStatDraft).
					Where("tech_alloc_start_time < ?", timeValue).
					Find(&allocations).Error

				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to retrieve intersecting allocations",
						Err:        err,
					}
				}

				// Process each allocation record
				for _, allocation := range allocations {
					startTime := allocation.TechAllocStartTime + diffTime
					endTime := allocation.TechAllocEndTime + diffTime

					shiftEnd, _ := getShiftEndTime(tx, companyId, woAlloc.ShiftCode, dateTimeComp, true)

					// Check if allocation is past shift end time and exceeds tolerance
					if shiftEnd < endTime && (endTime-shiftEnd) > 0.25 { //getVariableValue("AUTORELEASE_TOLERANCE")
						// Deallocate time by deleting the allocation record
						err = tx.Delete(&transactionworkshopentities.WorkOrderAllocation{}).
							Where("technician_allocation_system_number = ?", allocation.TechAllocSystemNumber).Error
						if err != nil {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to deallocate time",
								Err:        err,
							}
						}

						// Delete the associated service log
						err = tx.Delete(&transactionworkshopentities.ServiceLog{}).
							Where("technician_allocation_system_number = ?", allocation.TechAllocSystemNumber).Error
						if err != nil {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to delete service log",
								Err:        err,
							}
						}

						// Update work order status
						err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
							Where("work_order_system_number = ? AND work_order_operation_item_line = ?", idSysWo, result.WoLine).
							Updates(map[string]interface{}{
								"work_order_status_id": utils.SrvStatPending,
							}).Error
						if err != nil {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to update work order status",
								Err:        err,
							}
						}
					} else {
						// Move the service log
						err = tx.Model(&transactionworkshopentities.ServiceLog{}).
							Where("technician_allocation_system_number = ? AND service_status_id = ?", allocation.TechAllocSystemNumber, utils.SrvStatDraft).
							Updates(map[string]interface{}{
								"start_datetime": startTime,
								"end_datetime":   endTime,
							}).Error
						if err != nil {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to move service log",
								Err:        err,
							}
						}

						// Move the allocation
						err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
							Where("technician_allocation_system_number = ?", allocation.TechAllocSystemNumber).
							Updates(map[string]interface{}{
								"tech_alloc_start_time": startTime,
								"tech_alloc_end_time":   endTime,
							}).Error
						if err != nil {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to move allocation",
								Err:        err,
							}
						}
					}
				}

			}

		}

		if cpccode != 00003 {
			var logData struct {
				StartDatetime time.Time
				EndDatetime   time.Time
				StartTime     float64
			}

			err := tx.Model(&transactionworkshopentities.ServiceLog{}).
				Select("start_datetime, end_datetime, dbo.getTimeValue(CONVERT(VARCHAR, start_datetime, 108)) AS start_time").
				Where("technician_allocation_system_number = ? AND technician_allocation_line = ?", idAlloc, nextLine-1).
				Scan(&logData).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve service log data",
					Err:        err,
				}
			}

			var woData struct {
				WoSysNo int64
				WoLine  int64
			}

			err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
				Select("work_order_system_number, work_order_line").
				Where("technician_allocation_system_number = ?", idAlloc).
				Scan(&woData).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve Work Order system number and line",
					Err:        err,
				}
			}

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
					Message:    "Failed to update Work Order Allocation status",
					Err:        err,
				}
			}

			var count int64
			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Where("work_order_system_number = ? AND work_order_operation_item_line = ? AND service_status_id <> ?", woData.WoSysNo, woData.WoLine, utils.SrvStatStop).
				Count(&count).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to check work order allocation existence",
					Err:        err,
				}
			}

			if count == 0 {
				err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
					Where("work_order_system_number = ? AND work_order_operation_item_line = ?", woData.WoSysNo, woData.WoLine).
					Updates(map[string]interface{}{
						"service_status_id": utils.SrvStatStop,
					}).Error

				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to update Work Order status",
						Err:        err,
					}
				}
			}
		} else {
			// BR-specific logic
			err := tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
				Where("technician_allocation_system_number = ?", idAlloc).
				Updates(map[string]interface{}{
					"tech_alloc_end_date": dateTimeComp,
					"tech_alloc_end_time": timeValue,
					"service_status_id":   utils.SrvStatStop,
					"serv_actual_time":    woAlloc.ServActualTime + actualTime,
					"serv_progress_time":  woAlloc.ServProgressTime + actualTime,
				}).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update Work Order Allocation status",
					Err:        err,
				}
			}

			var count int64
			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Where("work_order_system_number = ? AND work_order_operation_item_line = ? AND service_status_id <> ?", idSysWo, woAlloc.WorkOrderOperationItemLine, utils.SrvStatStop).
				Count(&count).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to check work order status",
					Err:        err,
				}
			}

			if count == 0 {
				err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
					Where("work_order_system_number = ? AND work_order_operation_item_line = ?", idSysWo, woAlloc.WorkOrderOperationItemLine).
					Updates(map[string]interface{}{
						"service_status_id": utils.SrvStatStop,
					}).Error

				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to update work order status",
						Err:        err,
					}
				}

				err = tx.Model(&transactionworkshopentities.WorkOrder{}).
					Where("work_order_system_number = ?", idSysWo).
					Updates(map[string]interface{}{
						"work_order_status_id": utils.WoStatStop,
					}).Error

				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to update work order header status",
						Err:        err,
					}
				}
			}
		}

		//--Insert NEW LOG AS STOP
		// Ambil START_DATETIME dengan urutan TECHALLOC_LINE
		var startDateTime time.Time
		var pendingTime float64

		// Ambil START_DATETIME
		err = tx.Model(&transactionworkshopentities.ServiceLog{}).
			Select("start_datetime").
			Where("technician_allocation_system_number = ?", idAlloc).
			Order("technician_allocation_line ASC").
			Limit(1).
			Pluck("start_datetime", &startDateTime).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve START_DATETIME",
				Err:        err,
			}
		}

		// Hitung ACTUAL_TIME dan PENDING_TIME
		var sumresult struct {
			ActualTime  float64 `gorm:"column:ActualTime"`
			PendingTime float64 `gorm:"column:PendingTime"`
		}
		err = tx.Model(&transactionworkshopentities.ServiceLog{}).
			Select("SUM(ISNULL(actual_time, 0)) as ActualTime, SUM(ISNULL(pending_time, 0)) as PendingTime").
			Where("technician_allocation_system_number = ?", idAlloc).
			Scan(&sumresult).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate actual and pending time",
				Err:        err,
			}
		}
		actualTime = sumresult.ActualTime
		pendingTime = sumresult.PendingTime

		// Menyisipkan data ke dalam wtServiceLog
		serviceLog := transactionworkshopentities.ServiceLog{
			CompanyId:                        companyId,
			WorkOrderSystemNumber:            idSysWo,
			WorkOrderDocumentNumber:          result.WoDocNo,
			WorkOrderLine:                    result.WoLine,
			OperationItemCode:                oprItemCode,
			TechnicianId:                     woAlloc.TechnicianId,
			Frt:                              result.Frt,
			WorkOrderDate:                    result.WoDate.Format("2006-01-02 15:04:05"),
			ShiftCode:                        result.ShiftCode,
			ServiceStatusId:                  utils.SrvStatStop,
			StartDatetime:                    result.EndDatetime,
			EndDatetime:                      dateTimeComp,
			ActualTime:                       actualTime,
			PendingTime:                      pendingTime,
			EstimatedPendingTime:             0,
			TechnicianAllocationSystemNumber: idAlloc,
			TechnicianAllocationLine:         nextLine,
		}

		err = tx.Create(&serviceLog).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to insert new log as stop",
				Err:        err,
			}
		}

	} else {
		fmt.Println("Service status is already STOP")
	}

	return true, nil
}
