package masteritemcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"

	"after-sales/api/utils"

	// "after-sales/api/middlewares"

	// "strconv"

	"github.com/go-chi/chi/v5"
)

type MarkupMasterController interface {
	GetMarkupMasterList(writer http.ResponseWriter, request *http.Request)
	GetMarkupMasterByCode(writer http.ResponseWriter, request *http.Request)
	GetAllMarkupMasterIsActive(writer http.ResponseWriter, request *http.Request)
	SaveMarkupMaster(writer http.ResponseWriter, request *http.Request)
	ChangeStatusMarkupMaster(writer http.ResponseWriter, request *http.Request)
}

type MarkupMasterControllerImpl struct {
	markupMasterService masteritemservice.MarkupMasterService
}

func NewMarkupMasterController(MarkupMasterService masteritemservice.MarkupMasterService) MarkupMasterController {
	return &MarkupMasterControllerImpl{
		markupMasterService: MarkupMasterService,
	}
}

// @Summary Get All Markup Master
// @Description REST API Markup Master
// @Accept json
// @Produce json
// @Tags Master : Markup Master
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param markup_master_code query string false "markup_master_code"
// @Param markup_master_description query string false "markup_master_description"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/markup-master/ [get]
func (r *MarkupMasterControllerImpl) GetMarkupMasterList(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"markup_master_code":        queryValues.Get("markup_master_code"),
		"markup_master_description": queryValues.Get("markup_master_description"),
		"is_active":                 queryValues.Get("is_active"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result, err := r.markupMasterService.GetMarkupMasterList(filterCondition, pagination)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Markup Master Description by code
// @Description REST API Markup Master
// @Accept json
// @Produce json
// @Tags Master : Markup Master
// @Param markup_master_code path string true "markup_master_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/markup-master/by-code/{markup_master_code} [get]
func (r *MarkupMasterControllerImpl) GetMarkupMasterByCode(writer http.ResponseWriter, request *http.Request) {

	markupMasterCode := chi.URLParam(request, "markup_master_code")

	result, err := r.markupMasterService.GetMarkupMasterByCode(markupMasterCode)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *MarkupMasterControllerImpl) GetAllMarkupMasterIsActive(writer http.ResponseWriter, request *http.Request) {

	result, err := r.markupMasterService.GetAllMarkupMasterIsActive()

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Markup Master
// @Description REST API Markup Master
// @Accept json
// @Produce json
// @Tags Master : Markup Master
// @param reqBody body masteritempayloads.MarkupMasterResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/markup-master/ [post]
func (r *MarkupMasterControllerImpl) SaveMarkupMaster(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.MarkupMasterResponse
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	var message = ""

	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.markupMasterService.SaveMarkupMaster(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	if formRequest.MarkupMasterId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Markup Master
// @Description REST API Markup Master
// @Accept json
// @Produce json
// @Tags Master : Markup Master
// @param markup_master_id path int true "markup_master_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/markup-master/{markup_master_id} [patch]
func (r *MarkupMasterControllerImpl) ChangeStatusMarkupMaster(writer http.ResponseWriter, request *http.Request) {

	markupMasterId, _ := strconv.Atoi(chi.URLParam(request, "markup_master_id"))

	response, err := r.markupMasterService.ChangeStatusMasterMarkupMaster(int(markupMasterId))

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
