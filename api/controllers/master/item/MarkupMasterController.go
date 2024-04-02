package masteritemcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"net/http"
	"strconv"

	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"

	"after-sales/api/utils"

	"github.com/go-chi/chi/v5"
	// "after-sales/api/middlewares"
	// "strconv"
)

type MarkupMasterController interface {
	GetMarkupMasterList(writer http.ResponseWriter, request *http.Request)
	GetMarkupMasterByCode(writer http.ResponseWriter, request *http.Request)
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
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/markup-master [get]
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

	result := r.markupMasterService.GetMarkupMasterList(filterCondition, pagination)

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Markup Master Description by code
// @Description REST API Markup Master
// @Accept json
// @Produce json
// @Tags Master : Markup Master
// @Param markup_master_code path string true "markup_master_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/markup-master-by-code/{markup_master_code} [get]
func (r *MarkupMasterControllerImpl) GetMarkupMasterByCode(writer http.ResponseWriter, request *http.Request) {

	markupMasterCode := chi.URLParam(request, "markup_master_code")

	result := r.markupMasterService.GetMarkupMasterByCode(markupMasterCode)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Markup Master
// @Description REST API Markup Master
// @Accept json
// @Produce json
// @Tags Master : Markup Master
// @param reqBody body masteritempayloads.MarkupMasterResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/markup-master [post]
func (r *MarkupMasterControllerImpl) SaveMarkupMaster(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.MarkupMasterResponse
	var message = ""

	helper.ReadFromRequestBody(request, &formRequest)

	create := r.markupMasterService.SaveMarkupMaster(formRequest)

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
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/markup-master/{markup_master_id} [patch]
func (r *MarkupMasterControllerImpl) ChangeStatusMarkupMaster(writer http.ResponseWriter, request *http.Request) {

	markupMasterId, _ := strconv.Atoi(chi.URLParam(request, "markup_master_id"))

	response := r.markupMasterService.ChangeStatusMasterMarkupMaster(int(markupMasterId))

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
