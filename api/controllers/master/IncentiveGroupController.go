package mastercontroller

import (
	exceptions "after-sales/api/exceptions"
	"errors"

	// "after-sales/api/helper"
	helper "after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	// "after-sales/api/utils/validation"

	"github.com/go-chi/chi/v5"
	// "github.com/julienschmidt/httprouter"
)

type IncentiveGroupController interface {
	GetAllIncentiveGroup(writer http.ResponseWriter, request *http.Request)
	GetAllIncentiveGroupIsActive(writer http.ResponseWriter, request *http.Request)
	GetIncentiveGroupById(writer http.ResponseWriter, request *http.Request)
	SaveIncentiveGroup(writer http.ResponseWriter, request *http.Request)
	ChangeStatusIncentiveGroup(writer http.ResponseWriter, request *http.Request)
	UpdateIncentiveGroup(writer http.ResponseWriter, request *http.Request)
	GetAllIncentiveGroupDropDown(writer http.ResponseWriter, request *http.Request)
}

type IncentiveGroupControllerImpl struct {
	IncentiveGroupService masterservice.IncentiveGroupService
}

func NewIncentiveGroupController(IncentiveGroupService masterservice.IncentiveGroupService) IncentiveGroupController {
	return &IncentiveGroupControllerImpl{
		IncentiveGroupService: IncentiveGroupService,
	}
}

// @Summary Get All Incentive Group
// @Description REST API Incentive Group
// @Accept json
// @Produce json
// @Tags Master : Incentive Group
// @Security AuthorizationKeyAuth
// @Param incentive_group_id query string false "incentive_group_id"
// @Param incentive_group_code query string false "incentive_group_code"
// @Param incentive_group_name query string false "incentive_group_name"
// @Param effective_date query string false "effective_date"
// @Param is_active query string false "is_active"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param sort_of query string false "sort_of"
// @Param sort_by query string false "sort_by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/incentive-group [get]
func (r *IncentiveGroupControllerImpl) GetAllIncentiveGroup(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"incentive_group_id":   queryValues.Get("incentive_group_id"),
		"incentive_group_code": queryValues.Get("incentive_group_code"),
		"incentive_group_name": queryValues.Get("incentive_group_name"),
		"effective_date":       queryValues.Get("effective_date"),
		"is_active":            queryValues.Get("is_active"),
	}
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)
	result, err := r.IncentiveGroupService.GetAllIncentiveGroup(filterCondition, pagination)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get All Incentive Group Is Active
// @Description REST API Incentive Group
// @Accept json
// @Produce json
// @Tags Master : Incentive Group
// @Security AuthorizationKeyAuth
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/incentive-group/is-active [get]
func (r *IncentiveGroupControllerImpl) GetAllIncentiveGroupIsActive(writer http.ResponseWriter, request *http.Request) {

	result, err := r.IncentiveGroupService.GetAllIncentiveGroupIsActive()
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Incentive Group By Id
// @Description REST API Incentive Group
// @Accept json
// @Produce json
// @Tags Master : Incentive Group
// @Security AuthorizationKeyAuth
// @Param incentive_group_id path string true "incentive_group_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/incentive-group/by-id/{incentive_group_id} [get]
func (r *IncentiveGroupControllerImpl) GetIncentiveGroupById(writer http.ResponseWriter, request *http.Request) {
	incentiveGroupId, errA := strconv.Atoi(chi.URLParam(request, "incentive_group_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	incentiveGroupResponse, errors := r.IncentiveGroupService.GetIncentiveGroupById(incentiveGroupId)

	if errors != nil {
		helper.ReturnError(writer, request, errors)
		return
	}
	payloads.NewHandleSuccess(writer, incentiveGroupResponse, utils.GetDataSuccess, http.StatusOK)
}

// @Summary Save Incentive Group
// @Description REST API Incentive Group
// @Accept json
// @Produce json
// @Tags Master : Incentive Group
// @Security AuthorizationKeyAuth
// @Param incentive_group_id body object true "incentive_group_id"
// @Param incentive_group_code body string true "incentive_group_code"
// @Param incentive_group_name body string true "incentive_group_name"
// @Param effective_date body string true "effective_date"
// @Param is_active body string true "is_active"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/incentive-group [post]
func (r *IncentiveGroupControllerImpl) SaveIncentiveGroup(writer http.ResponseWriter, request *http.Request) {
	IncentiveGroupRequest := masterpayloads.IncentiveGroupResponse{}
	var message string

	err := jsonchecker.ReadFromRequestBody(request, &IncentiveGroupRequest)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, IncentiveGroupRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	create, err := r.IncentiveGroupService.SaveIncentiveGroup(IncentiveGroupRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	if IncentiveGroupRequest.IncentiveGroupId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusCreated)
}

// @Summary Change Status Incentive Group
// @Description REST API Incentive Group
// @Accept json
// @Produce json
// @Tags Master : Incentive Group
// @Security AuthorizationKeyAuth
// @Param incentive_group_id path string true "incentive_group_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/incentive-group/{incentive_group_id} [patch]
func (r *IncentiveGroupControllerImpl) ChangeStatusIncentiveGroup(writer http.ResponseWriter, request *http.Request) {

	IncentiveGroupId, errA := strconv.Atoi(chi.URLParam(request, "incentive_group_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.IncentiveGroupService.ChangeStatusIncentiveGroup(int(IncentiveGroupId))
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Update Incentive Group
// @Description REST API Incentive Group
// @Accept json
// @Produce json
// @Tags Master : Incentive Group
// @Security AuthorizationKeyAuth
// @Param incentive_group_id path string true "incentive_group_id"
// @Param incentive_group_code body string true "incentive_group_code"
// @Param incentive_group_name body string true "incentive_group_name"
// @Param effective_date body string true "effective_date"
// @Param is_active body string true "is_active"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/incentive-group/{incentive_group_id} [put]
func (r *IncentiveGroupControllerImpl) UpdateIncentiveGroup(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.UpdateIncentiveGroupRequest
	incentiveGroupId, errA := strconv.Atoi(chi.URLParam(request, "incentive_group_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.IncentiveGroupService.UpdateIncentiveGroup(formRequest, incentiveGroupId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Update Data Successfully!", http.StatusOK)
}

// @Summary Get All Incentive Group Drop Down
// @Description REST API Incentive Group
// @Accept json
// @Produce json
// @Tags Master : Incentive Group
// @Security AuthorizationKeyAuth
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/incentive-group/dropdown [get]
func (r *IncentiveGroupControllerImpl) GetAllIncentiveGroupDropDown(writer http.ResponseWriter, request *http.Request) {

	result, err := r.IncentiveGroupService.GetAllIncentiveGroupDropDown()
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}
