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

	"github.com/go-chi/chi/v5"
)

type SettingTechnicianController interface {
	GetAllSettingTechnician(writer http.ResponseWriter, request *http.Request)
	GetAllSettingTechinicianDetail(writer http.ResponseWriter, request *http.Request)
	GetSettingTechnicianById(writer http.ResponseWriter, request *http.Request)
	GetSettingTechnicianDetailById(writer http.ResponseWriter, request *http.Request)
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

func (r *SettingTechnicianControllerImpl) GetAllSettingTechnician(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"company_id": queryValues.Get("company_id"),
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

func (r *SettingTechnicianControllerImpl) GetSettingTechnicianById(writer http.ResponseWriter, request *http.Request) {
	settingTechnicianId, _ := strconv.Atoi(chi.URLParam(request, "setting_technician_system_number"))

	result, err := r.SettingTechnicianService.GetSettingTechnicianById(settingTechnicianId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully", http.StatusOK)
}

func (r *SettingTechnicianControllerImpl) GetSettingTechnicianDetailById(writer http.ResponseWriter, request *http.Request) {
	settingTechnicianDetailId, _ := strconv.Atoi(chi.URLParam(request, "setting_technician_detail_system_number"))

	result, err := r.SettingTechnicianService.GetSettingTechnicianDetailById(settingTechnicianDetailId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully", http.StatusOK)
}

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
	payloads.NewHandleSuccess(writer, update, "Create Data Successfully", http.StatusOK)
}
