package mastercontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MovingCodeController interface {
	GetAllMovingCode(writer http.ResponseWriter, request *http.Request)
	PushMovingCodePriority(writer http.ResponseWriter, request *http.Request)
	CreateMovingCode(writer http.ResponseWriter, request *http.Request)
	UpdateMovingCode(writer http.ResponseWriter, request *http.Request)
	GetMovingCodebyId(writer http.ResponseWriter, request *http.Request)
	ChangeStatusMovingCode(writer http.ResponseWriter, request *http.Request)
	GetDropdownMovingCode(writer http.ResponseWriter, request *http.Request)
	ActivateMovingCode(writer http.ResponseWriter, request *http.Request)
	DeactiveMovingCode(writer http.ResponseWriter, request *http.Request)
}

type MovingCodeControllerImpl struct {
	MovingCodeService masterservice.MovingCodeService
}

// @Summary Get Activate Moving Code
// @Description REST API Activate Moving Code
// @Accept json
// @Produce json
// @Tags Master : Moving Code
// @Security AuthorizationKeyAuth
// @Param moving_code_id path string true "moving_code_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/moving-code/activate/{moving_code_id} [get]
func (r *MovingCodeControllerImpl) ActivateMovingCode(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "moving_code_id")

	response, err := r.MovingCodeService.ActivateMovingCode(id)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Activate Status Successfully!", http.StatusOK)
}

// @Summary Get Deactive Moving Code
// @Description REST API Deactive Moving Code
// @Accept json
// @Produce json
// @Tags Master : Moving Code
// @Security AuthorizationKeyAuth
// @Param moving_code_id path string true "moving_code_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/moving-code/deactive/{moving_code_id} [get]
func (r *MovingCodeControllerImpl) DeactiveMovingCode(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "moving_code_id")

	response, err := r.MovingCodeService.DeactiveMovingCode(id)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Deactive Status Successfully!", http.StatusOK)
}

// @Summary Get Dropdown Moving Code
// @Description REST API Dropdown Moving Code
// @Accept json
// @Produce json
// @Tags Master : Moving Code
// @Security AuthorizationKeyAuth
// @Param company_id path string true "company_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/moving-code/drop-down/{company_id} [get]
func (r *MovingCodeControllerImpl) GetDropdownMovingCode(writer http.ResponseWriter, request *http.Request) {
	companyId, errA := strconv.Atoi(chi.URLParam(request, "company_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.MovingCodeService.GetDropdownMovingCode(companyId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Change Status Moving Code
// @Description REST API Change Status Moving Code
// @Accept json
// @Produce json
// @Tags Master : Moving Code
// @Security AuthorizationKeyAuth
// @Param moving_code_id path string true "moving_code_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/moving-code/{moving_code_id} [patch]
func (r *MovingCodeControllerImpl) ChangeStatusMovingCode(writer http.ResponseWriter, request *http.Request) {
	id, errA := strconv.Atoi(chi.URLParam(request, "moving_code_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.MovingCodeService.ChangeStatusMovingCode(id)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Change Status Successfully!", http.StatusOK)
}

// @Summary Create Moving Code
// @Description REST API Create Moving Code
// @Accept json
// @Produce json
// @Tags Master : Moving Code
// @Security AuthorizationKeyAuth
// @Param company_id path string true "company_id"
// @Param Request body masterpayloads.MovingCodeListRequest true "Request Body"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/moving-code/{company_id} [post]
func (r *MovingCodeControllerImpl) CreateMovingCode(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.MovingCodeListRequest
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

	create, err := r.MovingCodeService.CreateMovingCode(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Create Data Successfully!", http.StatusOK)
}

// @Summary Get All Moving Code
// @Description REST API Get All Moving Code
// @Accept json
// @Produce json
// @Tags Master : Moving Code
// @Security AuthorizationKeyAuth
// @Param company_id path string true "company_id"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Success 200 {object} payloads.ResponsePagination
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/moving-code/{company_id} [get]
func (r *MovingCodeControllerImpl) GetAllMovingCode(writer http.ResponseWriter, request *http.Request) {

	companyId, errA := strconv.Atoi(chi.URLParam(request, "company_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: "asc",
		SortBy: "Priority",
	}

	result, err := r.MovingCodeService.GetAllMovingCode(companyId, paginate)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)

}

// @Summary Get Moving Code by Id
// @Description REST API Get Moving Code by Id
// @Accept json
// @Produce json
// @Tags Master : Moving Code
// @Security AuthorizationKeyAuth
// @Param moving_code_id path string true "moving_code_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/moving-code/{moving_code_id} [get]
func (r *MovingCodeControllerImpl) GetMovingCodebyId(writer http.ResponseWriter, request *http.Request) {
	movingCodeId, errA := strconv.Atoi(chi.URLParam(request, "moving_code_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.MovingCodeService.GetMovingCodebyId(movingCodeId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Push Moving Code Priority
// @Description REST API Push Moving Code Priority
// @Accept json
// @Produce json
// @Tags Master : Moving Code
// @Security AuthorizationKeyAuth
// @Param moving_code_id path string true "moving_code_id"
// @Param company_id path string true "company_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/moving-code/push-priority/{company_id}/{moving_code_id} [get]
func (r *MovingCodeControllerImpl) PushMovingCodePriority(writer http.ResponseWriter, request *http.Request) {
	itemPackageId, errA := strconv.Atoi(chi.URLParam(request, "moving_code_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	companyId, errA := strconv.Atoi(chi.URLParam(request, "company_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.MovingCodeService.PushMovingCodePriority(companyId, itemPackageId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Push Priority Successfull!", http.StatusOK)
}

// @Summary Update Moving Code
// @Description REST API Update Moving Code
// @Accept json
// @Produce json
// @Tags Master : Moving Code
// @Security AuthorizationKeyAuth
// @Param Request body masterpayloads.MovingCodeListUpdate true "Request Body"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/moving-code [put]
func (r *MovingCodeControllerImpl) UpdateMovingCode(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.MovingCodeListUpdate
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

	create, err := r.MovingCodeService.UpdateMovingCode(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Update Data Successfully!", http.StatusOK)
}

func NewMovingCodeController(MovingCodeService masterservice.MovingCodeService) MovingCodeController {
	return &MovingCodeControllerImpl{
		MovingCodeService: MovingCodeService,
	}
}
