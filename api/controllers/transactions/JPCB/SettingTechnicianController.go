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

type SettingTechnicianController interface {
	GetAllSettingTechnician(writer http.ResponseWriter, request *http.Request)
	GetAllSettingTechinicianDetail(writer http.ResponseWriter, request *http.Request)
	GetSettingTechnicianById(writer http.ResponseWriter, request *http.Request)
	GetSettingTechnicianDetailById(writer http.ResponseWriter, request *http.Request)
	GetSettingTechnicianByCompanyDate(writer http.ResponseWriter, request *http.Request)
	SaveSettingTechnicianDetail(writer http.ResponseWriter, request *http.Request)
	UpdateSettingTechnicianDetail(writer http.ResponseWriter, request *http.Request)
}

type SettingTechnicianControllerImpl struct {
	SettingTechnicianService transactionjpcbservice.SettingTechnicianService
}

func NewSettingTechnicianController(SettingTechnicianServ transactionjpcbservice.SettingTechnicianService) SettingTechnicianController {
	return &SettingTechnicianControllerImpl{
		SettingTechnicianService: SettingTechnicianServ,
	}
}

// @Summary Get All Setting Technician
// @Description Get All Setting Technician
// @Tags Transaction JPCB: Setting Technician
// @Accept json
// @Produce json
// @Param company_id query string false "Company ID"
// @Param effective_date query string false "Effective Date"
// @Param setting_technician_system_number query string false "Setting Technician System Number"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/setting-technician [get]
func (r *SettingTechnicianControllerImpl) GetAllSettingTechnician(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"company_id":                       queryValues.Get("company_id"),
		"effective_date":                   queryValues.Get("effective_date"),
		"setting_technician_system_number": queryValues.Get("setting_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.SettingTechnicianService.GetAllSettingTechnician(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully", http.StatusOK, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get All Setting Technician Detail
// @Description Get All Setting Technician Detail
// @Tags Transaction JPCB: Setting Technician
// @Accept json
// @Produce json
// @Param setting_technician_system_number query string false "Setting Technician System Number"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/setting-technician/detail [get]
func (r *SettingTechnicianControllerImpl) GetAllSettingTechinicianDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"setting_technician_system_number": queryValues.Get("setting_technician_system_number"),
	}

	paginate := pagination.Pagination{
		Limit: utils.NewGetQueryInt(queryValues, "limit"),
		Page:  utils.NewGetQueryInt(queryValues, "page"),
	}
	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.SettingTechnicianService.GetAllSettingTechnicianDetail(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully", http.StatusOK, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Setting Technician By ID
// @Description Get Setting Technician By ID
// @Tags Transaction JPCB: Setting Technician
// @Accept json
// @Produce json
// @Param setting_technician_system_number path string true "Setting Technician System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/setting-technician/{setting_technician_system_number} [get]
func (r *SettingTechnicianControllerImpl) GetSettingTechnicianById(writer http.ResponseWriter, request *http.Request) {
	settingTechnicianId, _ := strconv.Atoi(chi.URLParam(request, "setting_technician_system_number"))

	result, err := r.SettingTechnicianService.GetSettingTechnicianById(settingTechnicianId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully", http.StatusOK)
}

// @Summary Get Setting Technician Detail By ID
// @Description Get Setting Technician Detail By ID
// @Tags Transaction JPCB: Setting Technician
// @Accept json
// @Produce json
// @Param setting_technician_detail_system_number path string true "Setting Technician Detail System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/setting-technician/detail/{setting_technician_detail_system_number} [get]
func (r *SettingTechnicianControllerImpl) GetSettingTechnicianDetailById(writer http.ResponseWriter, request *http.Request) {
	settingTechnicianDetailId, _ := strconv.Atoi(chi.URLParam(request, "setting_technician_detail_system_number"))

	result, err := r.SettingTechnicianService.GetSettingTechnicianDetailById(settingTechnicianDetailId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully", http.StatusOK)
}

// @Summary Get Setting Technician By Company Date
// @Description Get Setting Technician By Company Date
// @Tags Transaction JPCB: Setting Technician
// @Accept json
// @Produce json
// @Param company_id path string true "Company ID"
// @Param effective_date path string true "Effective Date"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/setting-technician/{company_id}/{effective_date} [get]
func (r *SettingTechnicianControllerImpl) GetSettingTechnicianByCompanyDate(writer http.ResponseWriter, request *http.Request) {
	companyId, _ := strconv.Atoi(chi.URLParam(request, "company_id"))
	effectiveDate, _ := time.Parse(time.RFC3339, chi.URLParam(request, "effective_date"))

	result, err := r.SettingTechnicianService.GetSettingTechnicianByCompanyDate(companyId, effectiveDate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully", http.StatusOK)
}

// @Summary Save Setting Technician Detail
// @Description Save Setting Technician Detail
// @Tags Transaction JPCB: Setting Technician
// @Accept json
// @Produce json
// @Param request body transactionjpcbpayloads.SettingTechnicianDetailSaveRequest true "Setting Technician Detail Save Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/setting-technician/detail [post]
func (r *SettingTechnicianControllerImpl) SaveSettingTechnicianDetail(writer http.ResponseWriter, request *http.Request) {
	formRequest := transactionjpcbpayloads.SettingTechnicianDetailSaveRequest{}
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

	create, err := r.SettingTechnicianService.SaveSettingTechnicianDetail(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Create Data Successfully", http.StatusOK)
}

// @Summary Update Setting Technician Detail
// @Description Update Setting Technician Detail
// @Tags Transaction JPCB: Setting Technician
// @Accept json
// @Produce json
// @Param setting_technician_detail_system_number path string true "Setting Technician Detail System Number"
// @Param request body transactionjpcbpayloads.SettingTechnicianDetailUpdateRequest true "Setting Technician Detail Update Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/setting-technician/detail/{setting_technician_detail_system_number} [put]
func (r *SettingTechnicianControllerImpl) UpdateSettingTechnicianDetail(writer http.ResponseWriter, request *http.Request) {
	settingTechnicianDetailId, _ := strconv.Atoi(chi.URLParam(request, "setting_technician_detail_system_number"))
	formRequest := transactionjpcbpayloads.SettingTechnicianDetailUpdateRequest{}
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

	update, err := r.SettingTechnicianService.UpdateSettingTechnicianDetail(settingTechnicianDetailId, formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, update, "Update Data Successfully", http.StatusOK)
}
