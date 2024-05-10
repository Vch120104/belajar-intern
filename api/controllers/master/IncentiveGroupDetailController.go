package mastercontroller

import (
	masterpayloads "after-sales/api/payloads/master"
	// masterrepository "after-sales/api/repositories/master"
	exceptionsss_test "after-sales/api/expectionsss"
	helper_test "after-sales/api/helper_testt"
	jsonchecker "after-sales/api/helper_testt/json/json-checker"
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
// @Param id path int true "incentive_group_id"  // Update this line
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/incentive-group-detail/by-header-id/{id} [get]
func (r *IncentiveGroupDetailControllerImpl) GetAllIncentiveGroupDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	IncentiveGroupId, _ := strconv.Atoi(chi.URLParam(request, "id"))

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result, err := r.IncentiveGroupDetailService.GetAllIncentiveGroupDetail(IncentiveGroupId, pagination)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Save Incentive Group Detail
// @Description REST API Incentive Group Detail
// @Accept json
// @Produce json
// @Tags Master : Incentive Group Detail
// @Param reqBody body masterpayloads.IncentiveGroupDetailResponse true "Form Request"
// @Param incentive_group_id_detail path int true "incentive_group_id_detail"  // Update this line
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/incentive-group-detail/{incentive_group_id_detail} [post]
func (r *IncentiveGroupDetailControllerImpl) SaveIncentiveGroupDetail(writer http.ResponseWriter, request *http.Request) {

	var incentiveGroupDetailRequest masterpayloads.IncentiveGroupDetailRequest
	var message string

	err := jsonchecker.ReadFromRequestBody(request, &incentiveGroupDetailRequest)
	if err != nil {
		exceptionsss_test.NewEntityException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, incentiveGroupDetailRequest)
	if err != nil {
		exceptionsss_test.NewBadRequestException(writer, request, err)
		return
	}
	create, err := r.IncentiveGroupDetailService.SaveIncentiveGroupDetail(incentiveGroupDetailRequest)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
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
// @Description REST API Incentive Group Detail
// @Accept json
// @Produce json
// @Tags Master : Incentive Group Detail
// @Param incentive_group_detail_id path string true "incentive_group_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/incentive-group-detail/by-detail-id/{incentive_group_detail_id} [get]
func (r *IncentiveGroupDetailControllerImpl) GetIncentiveGroupDetailById(writer http.ResponseWriter, request *http.Request) {
	// IncentiveGrouDetailId, _ := strconv.Atoi(params.ByName("incentive_group_detail_id"))
	IncentiveGrouDetailId, err := strconv.Atoi(chi.URLParam(request, "incentive_group_detail_id"))
	if err != nil {
		exceptionsss_test.NewAppException(writer, request, &exceptionsss_test.BaseErrorResponse{
			Err: err,
		})
		return
	}
	IncentiveGroupDetailResponse, errors := r.IncentiveGroupDetailService.GetIncentiveGroupDetailById(IncentiveGrouDetailId)

	if errors != nil {
		helper_test.ReturnError(writer, request, errors)
		return
	}
	payloads.NewHandleSuccess(writer, IncentiveGroupDetailResponse, utils.GetDataSuccess, http.StatusOK)
}
