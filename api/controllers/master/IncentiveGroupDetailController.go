package mastercontroller

import (
	masterpayloads "after-sales/api/payloads/master"
	// masterrepository "after-sales/api/repositories/master"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	// "after-sales/api/middlewares"

	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type IncentiveGroupDetailController interface {
	GetAllIncentiveGroupDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetIncentiveGroupDetailById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SaveIncentiveGroupDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
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
// @Description REST API Incentive Group Detail
// @Accept json
// @Produce json
// @Tags Master : Incentive Group Detail
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param incentive_group_detail_code query string false "incentive_group_detail_code"
// @Param incentive_group_detail_description query string false "incentive_group_detail_description"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/incentive-group-detail [get]
func (r *IncentiveGroupDetailControllerImpl) GetAllIncentiveGroupDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	queryParams := map[string]string{
		"incentive_group_detail_code":        params.ByName("incentive_group_detail_code"),
		"incentive_group_detail_description": params.ByName("incentive_group_detail_description"),
		"is_active":                          params.ByName("is_active"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(params, "limit"),
		Page:   utils.NewGetQueryInt(params, "page"),
		SortOf: params.ByName("sort_of"),
		SortBy: params.ByName("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result := r.IncentiveGroupDetailService.GetAllIncentiveGroupDetail(filterCondition, pagination)

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Save Incentive Group Detail
// @Description REST API Incentive Group Detail
// @Accept json
// @Produce json
// @Tags Master : Incentive Group Detail
// @param reqBody body masterpayloads.IncentiveGroupDetailResponse true "Form Request"
// @param incentive_group_id_detail path int true "incentive_group_id_detail"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/incentive-group-detail [post]
func (r *IncentiveGroupDetailControllerImpl) SaveIncentiveGroupDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var formRequest masterpayloads.IncentiveGroupDetailResponse
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	IncentiveGroupDetailId, _ := strconv.Atoi(params.ByName("incentive_group_id_detail"))

	create := r.IncentiveGroupDetailService.SaveIncentiveGroupDetail(IncentiveGroupDetailId, formRequest)

	if formRequest.IncentiveGroupDetailId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Get Incentive Group Detail By Id
// @Description REST API Incentive Group Detail
// @Accept json
// @Produce json
// @Tags Master : Incentive Group Detail
// @Param incentive_group_detail_id path string true "incentive_group_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/incentive-group-detail-by-id/{incentive_group_detail_id} [get]
func (r *IncentiveGroupDetailControllerImpl) GetIncentiveGroupDetailById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	IncentiveGrouDetailId,_ := strconv.Atoi(params.ByName("incentive_group_detail_id"))

	result := r.IncentiveGroupDetailService.GetIncentiveGroupDetailById(IncentiveGrouDetailId)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}