package mastercontroller

import (
	exceptions "after-sales/api/exceptions"
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

func (r *IncentiveGroupControllerImpl) GetAllIncentiveGroupIsActive(writer http.ResponseWriter, request *http.Request) {

	result, err := r.IncentiveGroupService.GetAllIncentiveGroupIsActive()
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *IncentiveGroupControllerImpl) GetIncentiveGroupById(writer http.ResponseWriter, request *http.Request) {
	incentiveGroupId, _ := strconv.Atoi(chi.URLParam(request, "incentive_group_id"))
	incentiveGroupResponse, errors := r.IncentiveGroupService.GetIncentiveGroupById(incentiveGroupId)

	if errors != nil {
		helper.ReturnError(writer, request, errors)
		return
	}
	payloads.NewHandleSuccess(writer, incentiveGroupResponse, utils.GetDataSuccess, http.StatusOK)
}

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

func (r *IncentiveGroupControllerImpl) ChangeStatusIncentiveGroup(writer http.ResponseWriter, request *http.Request) {

	IncentiveGroupId, _ := strconv.Atoi(chi.URLParam(request, "incentive_group_id"))

	response, err := r.IncentiveGroupService.ChangeStatusIncentiveGroup(int(IncentiveGroupId))
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

func (r *IncentiveGroupControllerImpl) UpdateIncentiveGroup(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.UpdateIncentiveGroupRequest
	incentiveGroupId, _ := strconv.Atoi(chi.URLParam(request, "incentive_group_id"))
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

func (r *IncentiveGroupControllerImpl) GetAllIncentiveGroupDropDown(writer http.ResponseWriter, request *http.Request) {

	result, err := r.IncentiveGroupService.GetAllIncentiveGroupDropDown()
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}
