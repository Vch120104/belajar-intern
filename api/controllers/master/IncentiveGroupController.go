package mastercontroller

import (
	exceptionsss_test "after-sales/api/expectionsss"
	// "after-sales/api/helper"
	helper_test "after-sales/api/helper_testt"
	jsonchecker "after-sales/api/helper_testt/json/json-checker"
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
	result, err := r.IncentiveGroupService.GetAllIncentiveGroup(filterCondition, pagination)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

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

	result, err := r.IncentiveGroupService.GetAllIncentiveGroupIsActive()
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *IncentiveGroupControllerImpl) GetIncentiveGroupById(writer http.ResponseWriter, request *http.Request) {
	incentiveGroupId, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		exceptionsss_test.NewAppException(writer, request, &exceptionsss_test.BaseErrorResponse{
			Err: err,
		})
		return
	}
	incentiveGroupResponse, errors := r.IncentiveGroupService.GetIncentiveGroupById(incentiveGroupId)

	if errors != nil {
		helper_test.ReturnError(writer, request, errors)
		return
	}
	payloads.NewHandleSuccess(writer, incentiveGroupResponse, utils.GetDataSuccess, http.StatusOK)
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
	// var formRequest masterpayloads.IncentiveGroupResponse
	// helper_test.ReadFromRequestBody(request, &formRequest)
	// var message = ""

	// create := r.IncentiveGroupService.SaveIncentiveGroup(formRequest)

	// if formRequest.IncentiveGroupId == 0 {
	// 	message = "Create Data Successfully!"
	// } else {
	// 	message = "Update Data Successfully!"
	// }

	// payloads.NewHandleSuccess(writer, create, message, http.StatusOK)

	IncentiveGroupRequest := masterpayloads.IncentiveGroupResponse{}

	err := jsonchecker.ReadFromRequestBody(request, &IncentiveGroupRequest)
	if err != nil {
		exceptionsss_test.NewEntityException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, IncentiveGroupRequest)
	if err != nil {
		exceptionsss_test.NewBadRequestException(writer, request, err)
		return
	}
	create, err := r.IncentiveGroupService.SaveIncentiveGroup(IncentiveGroupRequest)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Create Approval Success", http.StatusCreated)
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

	response, err := r.IncentiveGroupService.ChangeStatusIncentiveGroup(int(IncentiveGroupId))
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
