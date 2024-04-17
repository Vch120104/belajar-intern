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

type MovingCodeController interface {
	GetAllMovingCode(writer http.ResponseWriter, request *http.Request)
	SaveMovingCode(writer http.ResponseWriter, request *http.Request)
	ChangePriorityMovingCode(writer http.ResponseWriter, request *http.Request)
	ChangeStatusMovingCode(writer http.ResponseWriter, request *http.Request)
}
type MovingCodeControllerImpl struct {
	MovingCodeService masterservice.MovingCodeService
}

func NewMovingCodeController(MovingCodeService masterservice.MovingCodeService) MovingCodeController {
	return &MovingCodeControllerImpl{
		MovingCodeService: MovingCodeService,
	}
}

// @Summary Get All Moving Code
// @Description REST API Moving Code
// @Accept json
// @Produce json
// @Tags Master : Moving Code
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /aftersales-service/api/aftersales/moving-code [get]
func (r *MovingCodeControllerImpl) GetAllMovingCode(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: "asc",
		SortBy: "priority",
	}

	result := r.MovingCodeService.GetAllMovingCode(pagination)

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Save Moving Code
// @Description REST API Moving Code
// @Accept json
// @Produce json
// @Tags Master : Moving Code
// @param reqBody body masterpayloads.MovingCodeResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /aftersales-service/api/aftersales/moving-code [post]
func (r *MovingCodeControllerImpl) SaveMovingCode(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.MovingCodeRequest
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.MovingCodeService.SaveMovingCode(formRequest)

	if formRequest.MovingCodeId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Moving Code
// @Description REST API Moving Code
// @Accept json
// @Produce json
// @Tags Master : Moving Code
// @param moving_code_id path int true "moving_code_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /aftersales-service/api/aftersales/moving-code/priority-increase/{moving_code_id} [patch]
func (r *MovingCodeControllerImpl) ChangePriorityMovingCode(writer http.ResponseWriter, request *http.Request) {

	MovingCodeId, _ := strconv.Atoi(chi.URLParam(request, "moving_code_id"))

	response := r.MovingCodeService.ChangePriorityMovingCode(int(MovingCodeId))

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Change Status Moving Code
// @Description REST API Moving Code
// @Accept json
// @Produce json
// @Tags Master : Moving Code
// @param moving_code_id path int true "moving_code_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /aftersales-service/api/aftersales/moving-code/activation/{moving_code_id} [patch]
func (r *MovingCodeControllerImpl) ChangeStatusMovingCode(writer http.ResponseWriter, request *http.Request) {

	MovingCodeId, _ := strconv.Atoi(chi.URLParam(request, "moving_code_id"))

	response := r.MovingCodeService.ChangeStatusMovingCode(int(MovingCodeId))

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
