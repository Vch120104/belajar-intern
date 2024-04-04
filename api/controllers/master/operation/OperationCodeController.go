package masteroperationcontroller

import (
	exceptionsss_test "after-sales/api/expectionsss"
	helper_test "after-sales/api/helper_testt"
	jsonchecker "after-sales/api/helper_testt/json/json-checker"
	"after-sales/api/payloads"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type OperationCodeController interface {
	GetAllOperationCode(writer http.ResponseWriter, request *http.Request)
	GetByIdOperationCode(writer http.ResponseWriter, request *http.Request)
	SaveOperationCode(writer http.ResponseWriter, request *http.Request)
	ChangeStatusOperationCode(writer http.ResponseWriter, request *http.Request)
}

type OperationCodeControllerImpl struct {
	operationCodeService masteroperationservice.OperationCodeService
}

func NewOperationCodeController(operationCodeservice masteroperationservice.OperationCodeService) OperationCodeController {
	return &OperationCodeControllerImpl{
		operationCodeService: operationCodeservice,
	}
}

// @Summary Get All OPeration Code
// @Description REST API Operation Code
// @Accept json
// @Produce json
// @Tags Master : Operation Code
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param operation_code query string false "operation_code"
// @Param operation_name query string false "operation_name"
// @Param is_active query string false "is_active" Enums(true,false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-code [get]
func (r *OperationCodeControllerImpl) GetAllOperationCode(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"is_active":      queryValues.Get("is_active"),
		"operation_code": queryValues.Get("operation_code"),
		"operation_name": queryValues.Get("operation_name"),
	}
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result, err := r.operationCodeService.GetAllOperationCode(filterCondition, pagination)

	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Operation Code By Id
// @Description REST API Operation Code
// @Accept json
// @Produce json
// @Tags Master : Operation Code
// @Param operation_id path int true "operation_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-code/by-id/{operation_id} [get]
func (r *OperationCodeControllerImpl) GetByIdOperationCode(writer http.ResponseWriter, request *http.Request) {
	OperationIdStr, _ := strconv.Atoi(chi.URLParam(request, "operation_id"))

	result, err := r.operationCodeService.GetOperationCodeById(int(OperationIdStr))

	if err != nil {
		exceptionsss_test.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Update Data Successfully!", http.StatusOK)
}

// @Summary Save Operation Code
// @Description REST API Operation Code
// @Accept json
// @Produce json
// @Tags Master : Operation Code
// @param reqBody body masteroperationpayloads.OperationCodeSave true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-code/ [post]
func (r *OperationCodeControllerImpl) SaveOperationCode(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteroperationpayloads.OperationCodeSave
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptionsss_test.NewBadRequestException(writer, request, err)
		return
	}
	var message = ""

	create, err := r.operationCodeService.SaveOperationCode(formRequest)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	if formRequest.OperationId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Patch Operation Code
// @Description REST API Patch Operation Code
// @Accept json
// @Produce json
// @Tags Master : Operation Code
// @param operation_id path int true "operation_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-code/{operation_id} [patch]
func (r *OperationCodeControllerImpl) ChangeStatusOperationCode(writer http.ResponseWriter, request *http.Request) {

	OperationId, _ := strconv.Atoi(chi.URLParam(request,"operation_id"))

	response, err := r.operationCodeService.ChangeStatusOperationCode(OperationId)

	if err != nil {
		exceptionsss_test.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
