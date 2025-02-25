package transactionjpcbcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbservice "after-sales/api/services/transaction/JPCB"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type TechnicianAttendanceController interface {
	GetAllTechnicianAttendance(writer http.ResponseWriter, request *http.Request)
	GetAddLineTechnician(writer http.ResponseWriter, request *http.Request)
	SaveTechnicianAttendance(writer http.ResponseWriter, request *http.Request)
	ChangeStatusTechnicianAttendance(writer http.ResponseWriter, request *http.Request)
}

type TechnicianAttendanceControllerImpl struct {
	TechnicianAttendanceService transactionjpcbservice.TechnicianAttendanceService
}

func NewTechnicianAttendanceController(technicianAttendanceService transactionjpcbservice.TechnicianAttendanceService) TechnicianAttendanceController {
	return &TechnicianAttendanceControllerImpl{
		TechnicianAttendanceService: technicianAttendanceService,
	}
}

// @Summary Get All Technician Attendance
// @Description Get All Technician Attendance
// @Tags Transaction : JPCB Technician Attendance
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param company_id query string false "Company ID"
// @Param service_date query string false "Service Date"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/technician-attendance [get]
func (r *TechnicianAttendanceControllerImpl) GetAllTechnicianAttendance(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"company_id":   queryValues.Get("company_id"),
		"service_date": queryValues.Get("service_date"),
	}

	paginate := pagination.Pagination{
		Limit: utils.NewGetQueryInt(queryValues, "limit"),
		Page:  utils.NewGetQueryInt(queryValues, "page"),
	}
	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.TechnicianAttendanceService.GetAllTechnicianAttendance(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully", http.StatusOK, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Add Line Technician
// @Description Get Add Line Technician
// @Tags Transaction : JPCB Technician Attendance
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param service_date query string true "Service Date"
// @Param company_id query string true "Company ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/technician-attendance/add-line [get]
func (r *TechnicianAttendanceControllerImpl) GetAddLineTechnician(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"service_date": queryValues.Get("service_date"),
		"company_id":   queryValues.Get("company_id"),
	}

	if queryParams["service_date"] == "" {
		payloads.NewHandleError(writer, "service_date is required", http.StatusBadRequest)
		return
	}

	parsedDate, parseErr := time.Parse(time.RFC3339, queryParams["service_date"])
	if parseErr != nil {
		payloads.NewHandleError(writer, "Error parsing 'service_date'. please use RFC3339 format", http.StatusBadRequest)
		return
	}
	queryParams["service_date"] = parsedDate.Format("2006-01-02")

	paginate := pagination.Pagination{
		Limit: utils.NewGetQueryInt(queryValues, "limit"),
		Page:  utils.NewGetQueryInt(queryValues, "page"),
	}
	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.TechnicianAttendanceService.GetAddLineTechnician(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully", http.StatusOK, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Save Technician Attendance
// @Description Save Technician Attendance
// @Tags Transaction : JPCB Technician Attendance
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param body body transactionjpcbpayloads.TechnicianAttendanceSaveRequest true "Technician Attendance Save Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/technician-attendance [post]
func (r *TechnicianAttendanceControllerImpl) SaveTechnicianAttendance(writer http.ResponseWriter, request *http.Request) {
	formRequest := transactionjpcbpayloads.TechnicianAttendanceSaveRequest{}
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.TechnicianAttendanceService.SaveTechnicianAttendance(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Create Data Successfully", http.StatusOK)
}

// @Summary Change Status Technician Attendance
// @Description Change Status Technician Attendance
// @Tags Transaction : JPCB Technician Attendance
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param technician_attendance_id path int true "Technician Attendance ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/technician-attendance/{technician_attendance_id} [patch]
func (r *TechnicianAttendanceControllerImpl) ChangeStatusTechnicianAttendance(writer http.ResponseWriter, request *http.Request) {
	technicianAttendanceId, _ := strconv.Atoi(chi.URLParam(request, "technician_attendance_id"))

	update, err := r.TechnicianAttendanceService.ChangeStatusTechnicianAttendance(technicianAttendanceId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, update, "Update Data Successfully", http.StatusOK)
}
