package mastercontroller

import (
	masterpayloads "after-sales/api/payloads/master"
	"errors"

	// masterrepository "after-sales/api/repositories/master"
	exceptions "after-sales/api/exceptions"
	helper "after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"after-sales/api/validation"

	// "after-sales/api/middlewares"

	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type IncentiveGroupDetailController interface {
	GetAllIncentiveGroupDetail(writer http.ResponseWriter, request *http.Request)
	GetIncentiveGroupDetailById(writer http.ResponseWriter, request *http.Request)
	SaveIncentiveGroupDetail(writer http.ResponseWriter, request *http.Request)
	UpdateIncentiveGroupDetail(writer http.ResponseWriter, request *http.Request)
}
type IncentiveGroupDetailControllerImpl struct {
	IncentiveGroupDetailService masterservice.IncentiveGroupDetailService
}

func NewIncentiveGroupDetailController(IncentiveGroupDetailService masterservice.IncentiveGroupDetailService) IncentiveGroupDetailController {
	return &IncentiveGroupDetailControllerImpl{
		IncentiveGroupDetailService: IncentiveGroupDetailService,
	}
}

// @Summary Get All Incentive Group Detail
// @Description Get All Incentive Group Detail
// @Tags Master : Incentive Group Detail
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Incentive Group Id"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/incentive-group-detail/{id} [get]
func (r *IncentiveGroupDetailControllerImpl) GetAllIncentiveGroupDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	IncentiveGroupId, errA := strconv.Atoi(chi.URLParam(request, "id"))
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

	result, err := r.IncentiveGroupDetailService.GetAllIncentiveGroupDetail(IncentiveGroupId, pagination)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Save Incentive Group Detail
// @Description Save Incentive Group Detail
// @Tags Master : Incentive Group Detail
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Incentive Group Id"
// @Param IncentiveGroupDetailRequest body masterpayloads.IncentiveGroupDetailRequest true "Incentive Group Detail Request"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/incentive-group-detail/{id} [post]
func (r *IncentiveGroupDetailControllerImpl) SaveIncentiveGroupDetail(writer http.ResponseWriter, request *http.Request) {

	var incentiveGroupDetailRequest masterpayloads.IncentiveGroupDetailRequest
	var message string

	err := jsonchecker.ReadFromRequestBody(request, &incentiveGroupDetailRequest)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, incentiveGroupDetailRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	create, err := r.IncentiveGroupDetailService.SaveIncentiveGroupDetail(incentiveGroupDetailRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	if incentiveGroupDetailRequest.IncentiveGroupDetailId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusCreated)
}

// @Summary Get Incentive Group Detail By Id
// @Description Get Incentive Group Detail By Id
// @Tags Master : Incentive Group Detail
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param incentive_group_detail_id path int true "Incentive Group Detail Id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/incentive-group-detail/{incentive_group_detail_id} [get]
func (r *IncentiveGroupDetailControllerImpl) GetIncentiveGroupDetailById(writer http.ResponseWriter, request *http.Request) {
	// IncentiveGrouDetailId, _ := strconv.Atoi(params.ByName("incentive_group_detail_id"))
	IncentiveGrouDetailId, err := strconv.Atoi(chi.URLParam(request, "incentive_group_detail_id"))
	if err != nil {
		exceptions.NewAppException(writer, request, &exceptions.BaseErrorResponse{
			Err: err,
		})
		return
	}
	IncentiveGroupDetailResponse, errors := r.IncentiveGroupDetailService.GetIncentiveGroupDetailById(IncentiveGrouDetailId)

	if errors != nil {
		helper.ReturnError(writer, request, errors)
		return
	}
	payloads.NewHandleSuccess(writer, IncentiveGroupDetailResponse, utils.GetDataSuccess, http.StatusOK)
}

// @Summary Update Incentive Group Detail
// @Description Update Incentive Group Detail
// @Tags Master : Incentive Group Detail
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param incentive_group_detail_id path int true "Incentive Group Detail Id"
// @Param IncentiveGroupDetailRequest body masterpayloads.UpdateIncentiveGroupDetailRequest true "Incentive Group Detail Request"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/incentive-group-detail/{incentive_group_detail_id} [put]
func (r *IncentiveGroupDetailControllerImpl) UpdateIncentiveGroupDetail(writer http.ResponseWriter, request *http.Request) {
	id, errA := strconv.Atoi(chi.URLParam(request, "incentive_group_detail_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	var incentiveGroupDetailRequest masterpayloads.UpdateIncentiveGroupDetailRequest
	var message string

	err := jsonchecker.ReadFromRequestBody(request, &incentiveGroupDetailRequest)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, incentiveGroupDetailRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	create, err := r.IncentiveGroupDetailService.UpdateIncentiveGroupDetail(id, incentiveGroupDetailRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	message = "Update Data Successfully!"

	payloads.NewHandleSuccess(writer, create, message, http.StatusCreated)
}
