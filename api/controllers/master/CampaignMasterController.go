package mastercontroller

import (
	exceptions "after-sales/api/exceptions"
	helper "after-sales/api/helper"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CampaignMasterController interface {
	SaveCampaignMaster(writer http.ResponseWriter, request *http.Request)
	SaveCampaignMasterDetail(writer http.ResponseWriter, request *http.Request)
	SaveCampaignMasterDetailFromHistory(writer http.ResponseWriter, request *http.Request)
	SaveCampaignMasterDetailFromPackage(writer http.ResponseWriter, request *http.Request)
	ChangeStatusCampaignMaster(writer http.ResponseWriter, request *http.Request)
	ActivateCampaignMasterDetail(writer http.ResponseWriter, request *http.Request)
	DeactivateCampaignMasterDetail(writer http.ResponseWriter, request *http.Request)
	GetByIdCampaignMaster(writer http.ResponseWriter, request *http.Request)
	GetByIdCampaignMasterDetail(writer http.ResponseWriter, request *http.Request)
	GetByCodeCampaignMaster(writer http.ResponseWriter, request *http.Request)
	GetAllCampaignMasterCodeAndName(writer http.ResponseWriter, request *http.Request)
	GetAllCampaignMaster(writer http.ResponseWriter, request *http.Request)
	GetAllCampaignMasterDetail(writer http.ResponseWriter, request *http.Request)
	UpdateCampaignMasterDetail(writer http.ResponseWriter, request *http.Request)
	GetAllPackageMasterToCopy(writer http.ResponseWriter, request *http.Request)
	SelectFromPackageMaster(writer http.ResponseWriter, request *http.Request)
}

type CampaignMasterControllerImpl struct {
	CampaignMasterService masterservice.CampaignMasterService
}

func NewCampaignMasterController(campaignmasterservice masterservice.CampaignMasterService) CampaignMasterController {
	return &CampaignMasterControllerImpl{
		CampaignMasterService: campaignmasterservice,
	}
}

// @Summary Save Campaign Master
// @Description Save Campaign Master
// @Tags Master : Campaign Master
// @Accept json
// @Produce json
// @Param body body masterpayloads.CampaignMasterPost true "Campaign Master"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/campaign-master [post]
func (r *CampaignMasterControllerImpl) SaveCampaignMaster(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.CampaignMasterPost
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	var message string
	var status int

	create, err := r.CampaignMasterService.PostCampaignMaster(formRequest)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}

	if formRequest.CampaignId == 0 {
		message = "Created Data Successfully!"
		status = http.StatusCreated
	} else {
		message = "Updated Data Successfully!"
		status = http.StatusOK
	}

	payloads.NewHandleSuccess(writer, create, message, status)
}

// @Summary Save Campaign Master Detail
// @Description Save Campaign Master Detail
// @Tags Master : Campaign Master
// @Accept json
// @Produce json
// @Param campaign_id path int true "Campaign ID"
// @Param body body masterpayloads.CampaignMasterDetailPayloads true "Campaign Master Detail"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/campaign-master/detail/{campaign_id} [post]
func (r *CampaignMasterControllerImpl) SaveCampaignMasterDetail(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.CampaignMasterDetailPayloads
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	campaignId, _ := strconv.Atoi(chi.URLParam(request, "campaign_id"))

	create, err := r.CampaignMasterService.PostCampaignDetailMaster(formRequest, campaignId)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Create Data Successfully!", http.StatusCreated)
}

// @Summary Save Campaign Master Detail From History
// @Description Save Campaign Master Detail From History
// @Tags Master : Campaign Master
// @Accept json
// @Produce json
// @Param campaign_id_1 path int true "Campaign ID 1"
// @Param campaign_id_2 path int true "Campaign ID 2"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/campaign-master/detail/save-from-history/{campaign_id_1}/{campaign_id_2} [post]
func (r *CampaignMasterControllerImpl) SaveCampaignMasterDetailFromHistory(writer http.ResponseWriter, request *http.Request) {
	CampaignId1, errA := strconv.Atoi(chi.URLParam(request, "campaign_id_1"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	CampaignId2, errA := strconv.Atoi(chi.URLParam(request, "campaign_id_2"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	var message = ""
	response, err := r.CampaignMasterService.PostCampaignMasterDetailFromHistory(CampaignId1, CampaignId2)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}
	message = "Create Data Successfully!"

	payloads.NewHandleSuccess(writer, response, message, http.StatusOK)
}

// @Summary Save Campaign Master Detail From Package
// @Description Save Campaign Master Detail From Package
// @Tags Master : Campaign Master
// @Accept json
// @Produce json
// @Param body body masterpayloads.CampaignMasterDetailPostFromPackageRequest true "Campaign Master Detail From Package"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/campaign-master/detail/save-from-package [post]
func (r *CampaignMasterControllerImpl) SaveCampaignMasterDetailFromPackage(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.CampaignMasterDetailPostFromPackageRequest
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	response, err := r.CampaignMasterService.PostCampaignMasterDetailFromPackage(formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Create Data Successfully!", http.StatusCreated)
}

// @Summary Change Status Campaign Master
// @Description Change Status Campaign Master
// @Tags Master : Campaign Master
// @Accept json
// @Produce json
// @Param campaign_id path int true "Campaign ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/campaign-master/{campaign_id} [patch]
func (r *CampaignMasterControllerImpl) ChangeStatusCampaignMaster(writer http.ResponseWriter, request *http.Request) {
	CampaignId, errA := strconv.Atoi(chi.URLParam(request, "campaign_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	response, err := r.CampaignMasterService.ChangeStatusCampaignMaster(CampaignId)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
	}
	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", 200)
}

// @Summary Activate Campaign Master Detail
// @Description Activate Campaign Master Detail
// @Tags Master : Campaign Master
// @Accept json
// @Produce json
// @Param campaign_detail_id path int true "Campaign Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/campaign-master/detail/activate/{campaign_detail_id} [patch]
func (r *CampaignMasterControllerImpl) ActivateCampaignMasterDetail(writer http.ResponseWriter, request *http.Request) {
	queryId := chi.URLParam(request, "campaign_detail_id")
	id, err := r.CampaignMasterService.ActivateCampaignMasterDetail(queryId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, id, "Update Data Successfully!", http.StatusOK)
}

// @Summary Deactivate Campaign Master Detail
// @Description Deactivate Campaign Master Detail
// @Tags Master : Campaign Master
// @Accept json
// @Produce json
// @Param campaign_detail_id path int true "Campaign Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/campaign-master/detail/deactivate/{campaign_detail_id} [patch]
func (r *CampaignMasterControllerImpl) DeactivateCampaignMasterDetail(writer http.ResponseWriter, request *http.Request) {
	queryId := chi.URLParam(request, "campaign_detail_id")
	id, err := r.CampaignMasterService.DeactivateCampaignMasterDetail(queryId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, id, "Update Data Successfully!", http.StatusOK)
}

// @Summary Get Campaign Master By ID
// @Description Get Campaign Master By ID
// @Tags Master : Campaign Master
// @Accept json
// @Produce json
// @Param campaign_id path int true "Campaign ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/campaign-master/{campaign_id} [get]
func (r *CampaignMasterControllerImpl) GetByIdCampaignMaster(writer http.ResponseWriter, request *http.Request) {
	CampaignIdstr := chi.URLParam(request, "campaign_id")

	CampaignId, errA := strconv.Atoi(CampaignIdstr)
	if CampaignId <= 0 {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Err: errors.New("id cannot be 0")})
		return
	}
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.CampaignMasterService.GetByIdCampaignMaster(CampaignId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Campaign Master Detail By ID
// @Description Get Campaign Master Detail By ID
// @Tags Master : Campaign Master
// @Accept json
// @Produce json
// @Param campaign_detail_id path int true "Campaign Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/campaign-master/detail/by-id/{campaign_detail_id} [get]
func (r *CampaignMasterControllerImpl) GetByIdCampaignMasterDetail(writer http.ResponseWriter, request *http.Request) {
	CampaignDetailIdstr := chi.URLParam(request, "campaign_detail_id")
	CampaignDetailId, errA := strconv.Atoi(CampaignDetailIdstr)

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.CampaignMasterService.GetByIdCampaignMasterDetail(CampaignDetailId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Campaign Master By Code
// @Description Get Campaign Master By Code
// @Tags Master : Campaign Master
// @Accept json
// @Produce json
// @Param campaign_code path string true "Campaign Code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/campaign-master/by-code/{campaign_code} [get]
func (r *CampaignMasterControllerImpl) GetByCodeCampaignMaster(writer http.ResponseWriter, request *http.Request) {

	encodedcampaignCode := chi.URLParam(request, "*")

	if len(encodedcampaignCode) > 0 && encodedcampaignCode[0] == '/' {
		encodedcampaignCode = encodedcampaignCode[1:]
	}

	campaignCode, err := url.PathUnescape(encodedcampaignCode)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("failed to decode campaign code")})
		return
	}

	result, baseErr := r.CampaignMasterService.GetByCodeCampaignMaster(campaignCode)
	if baseErr != nil {
		helper.ReturnError(writer, request, baseErr)
		return
	}
	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK)
}

// @Summary Get All Campaign Master
// @Description Get All Campaign Master
// @Tags Master : Campaign Master
// @Accept json
// @Produce json
// @Param is_active query string false "Is Active"
// @Param brand_id query string false "Brand ID"
// @Param campaign_id query string false "Campaign ID"
// @Param campaign_code query string false "Campaign Code"
// @Param campaign_name query string false "Campaign Name"
// @Param model_id query string false "Model ID"
// @Param model_code query string false "Model Code"
// @Param model_description query string false "Model Description"
// @Param campaign_period_from query string false "Campaign Period From"
// @Param campaign_period_to query string false "Campaign Period To"
// @Param company_id query string false "Company ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.ResponsePagination
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/campaign-master [get]
func (r *CampaignMasterControllerImpl) GetAllCampaignMaster(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"is_active":            queryValues.Get("is_active"),
		"brand_id":             queryValues.Get("brand_id"),
		"campaign_id":          queryValues.Get("campaign_id"),
		"campaign_code":        queryValues.Get("campaign_code"),
		"campaign_name":        queryValues.Get("campaign_name"),
		"model_id":             queryValues.Get("model_id"),
		"model_code":           queryValues.Get("model_code"),
		"model_description":    queryValues.Get("model_description"),
		"campaign_period_from": queryValues.Get("campaign_period_from"),
		"campaign_period_to":   queryValues.Get("campaign_period_to"),
		"company_id":           queryValues.Get("company_id"),
	}
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result, err := r.CampaignMasterService.GetAllCampaignMaster(filterCondition, pagination)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		result.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		result.Limit,
		result.Page,
		int64(result.TotalRows),
		result.TotalPages,
	)
}

// @Summary Get All Campaign Master Detail
// @Description Get All Campaign Master Detail
// @Tags Master : Campaign Master
// @Accept json
// @Produce json
// @Param campaign_id path int true "Campaign ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.ResponsePagination
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/campaign-master/detail/{campaign_id} [get]
func (r *CampaignMasterControllerImpl) GetAllCampaignMasterDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	CampaignIdStr := chi.URLParam(request, "campaign_id")

	CampaignId, errA := strconv.Atoi(CampaignIdStr)
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result, err := r.CampaignMasterService.GetAllCampaignMasterDetail(pagination, CampaignId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(
		writer,
		result.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		result.Limit,
		result.Page,
		int64(result.TotalRows),
		result.TotalPages,
	)

}

// @Summary Get All Campaign Master Code And Name
// @Description Get All Campaign Master Code And Name
// @Tags Master : Campaign Master
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.ResponsePagination
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/campaign-master/history [get]
func (r *CampaignMasterControllerImpl) GetAllCampaignMasterCodeAndName(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	result, err := r.CampaignMasterService.GetAllCampaignMasterCodeAndName(pagination)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)

}

// @Summary Update Campaign Master Detail
// @Description Update Campaign Master Detail
// @Tags Master : Campaign Master
// @Accept json
// @Produce json
// @Param campaign_detail_id path int true "Campaign Detail ID"
// @Param body body masterpayloads.CampaignMasterDetailPayloads true "Campaign Master Detail"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/campaign-master/detail/update/{campaign_detail_id} [put]
func (r *CampaignMasterControllerImpl) UpdateCampaignMasterDetail(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.CampaignMasterDetailPayloads
	CampaignDetailIdstr := chi.URLParam(request, "campaign_detail_id")
	CampaignDetailId, errA := strconv.Atoi(CampaignDetailIdstr)
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	var message = ""
	result, err := r.CampaignMasterService.UpdateCampaignMasterDetail(CampaignDetailId, formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	message = "Update Data Successfully!"
	payloads.NewHandleSuccess(writer, result, message, http.StatusOK)
}

// @Summary Get All Package Master To Copy
// @Description Get All Package Master To Copy
// @Tags Master : Campaign Master
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.ResponsePagination
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/campaign-master/package [get]
func (r *CampaignMasterControllerImpl) GetAllPackageMasterToCopy(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	result, err := r.CampaignMasterService.GetAllPackageMasterToCopy(pagination)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Select From Package Master
// @Description Select From Package Master
// @Tags Master : Campaign Master
// @Accept json
// @Produce json
// @Param package_id path int true "Package ID"
// @Param campaign_detail_id path int true "Campaign Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/campaign-master/package/{package_id}/{campaign_detail_id} [get]
func (r *CampaignMasterControllerImpl) SelectFromPackageMaster(writer http.ResponseWriter, request *http.Request) {
	var message = ""
	PackageMaster, errA := strconv.Atoi(chi.URLParam(request, "package_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	CampaignMasterId, errA := strconv.Atoi(chi.URLParam(request, "campaign_detail_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.CampaignMasterService.SelectFromPackageMaster(PackageMaster, CampaignMasterId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, message, http.StatusOK)
}
