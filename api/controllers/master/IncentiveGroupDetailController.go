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

func (r *IncentiveGroupDetailControllerImpl) UpdateIncentiveGroupDetail(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(request, "incentive_group_detail_id"))

	var incentiveGroupDetailRequest masterpayloads.UpdateIncentiveGroupDetailRequest
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
	create, err := r.IncentiveGroupDetailService.UpdateIncentiveGroupDetail(id, incentiveGroupDetailRequest)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	message = "Update Data Successfully!"

	payloads.NewHandleSuccess(writer, create, message, http.StatusCreated)
}
