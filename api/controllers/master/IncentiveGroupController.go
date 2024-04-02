package mastercontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type IncentiveGroupController interface {
	GetAllIncentiveGroup(writer http.ResponseWriter, request *http.Request)
	GetAllIncentiveGroupIsActive(writer http.ResponseWriter, request *http.Request)
	GetIncentiveGroupById(writer http.ResponseWriter, request *http.Request)
	SaveIncentiveGroup(writer http.ResponseWriter, request *http.Request)
	ChangeStatusIncentiveGroup(writer http.ResponseWriter, request *http.Request)
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
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param incentive_group_code query string false "incentive_group_code"
// @Param incentive_group_name query string false "incentive_group_name"
// @Param effective_date query string false "effective_date"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/incentive-group [get]
func (r *IncentiveGroupControllerImpl) GetAllIncentiveGroup(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"incentive_group_code": queryValues.Get("incentive_group_code"),
		"incentive_group_name": queryValues.Get("incentive_group_name"),
		"effective_date":       queryValues.Get("effective_date"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result := r.IncentiveGroupService.GetAllIncentiveGroup(filterCondition, pagination)

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get All Incentive Group Drop Down
// @Description REST API Incentive Group
// @Accept json
// @Produce json
// @Tags Master : Incentive Group
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/incentive-group/drop-down [get]
func (r *IncentiveGroupControllerImpl) GetAllIncentiveGroupIsActive(writer http.ResponseWriter, request *http.Request) {

	result := r.IncentiveGroupService.GetAllIncentiveGroupIsActive()

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Incentive Group By Id
// @Description REST API Incentive Group
// @Accept json
// @Produce json
// @Tags Master : Incentive Group
// @Param incentive_group_id path string true "incentive_group_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/incentive-group/by-id/{incentive_group_id} [get]
func (r *IncentiveGroupControllerImpl) GetIncentiveGroupById(writer http.ResponseWriter, request *http.Request) {
	IncentiveGroupId, _ := strconv.Atoi(chi.URLParam(request, "incentive_group_id"))

	result := r.IncentiveGroupService.GetIncentiveGroupById(IncentiveGroupId)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Incentive Group
// @Description REST API Incentive Group
// @Accept json
// @Produce json
// @Tags Master : Incentive Group
// @param reqBody body masterpayloads.IncentiveGroupResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/incentive-group [post]
func (r *IncentiveGroupControllerImpl) SaveIncentiveGroup(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.IncentiveGroupResponse
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.IncentiveGroupService.SaveIncentiveGroup(formRequest)

	if formRequest.IncentiveGroupId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Incentive Group
// @Description REST API Incentive Group
// @Accept json
// @Produce json
// @Tags Master : Incentive Group
// @param incentive_group_id path int true "incentive_group_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/incentive-group/{incentive_group_id} [patch]
func (r *IncentiveGroupControllerImpl) ChangeStatusIncentiveGroup(writer http.ResponseWriter, request *http.Request) {

	IncentiveGroupId, _ := strconv.Atoi(chi.URLParam(request, "incentive_group_id"))

	response := r.IncentiveGroupService.ChangeStatusIncentiveGroup(int(IncentiveGroupId))

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
