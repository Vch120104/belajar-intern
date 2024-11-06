package masteroperationcontroller

import (
	"after-sales/api/exceptions"
	helper "after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationservice "after-sales/api/services/master/operation"
)

type OperationModelMappingController interface {
	GetOperationModelMappingLookup(writer http.ResponseWriter, request *http.Request)
	GetOperationModelMappingById(writer http.ResponseWriter, request *http.Request)
	GetOperationModelMappingByBrandModelOperationCode(writer http.ResponseWriter, request *http.Request)
	SaveOperationModelMapping(writer http.ResponseWriter, request *http.Request)
	ChangeStatusOperationModelMapping(writer http.ResponseWriter, request *http.Request)
	SaveOperationModelMappingFrt(writer http.ResponseWriter, request *http.Request)
	ActivateOperationFrt(writer http.ResponseWriter, request *http.Request)
	DeactivateOperationFrt(writer http.ResponseWriter, request *http.Request)
	SaveOperationModelMappingDocumentRequirement(writer http.ResponseWriter, request *http.Request)
	DeactivateOperationDocumentRequirement(writer http.ResponseWriter, request *http.Request)
	ActivateOperationDocumentRequirement(writer http.ResponseWriter, request *http.Request)
	GetAllOperationFrt(writer http.ResponseWriter, request *http.Request)
	GetOperationFrtById(writer http.ResponseWriter, request *http.Request)
	GetAllOperationDocumentRequirement(writer http.ResponseWriter, request *http.Request)
	GetOperationDocumentRequirementById(writer http.ResponseWriter, request *http.Request)
	SaveOperationLevel(writer http.ResponseWriter, request *http.Request)
	GetAllOperationLevel(writer http.ResponseWriter, request *http.Request)
	GetOperationLevelById(writer http.ResponseWriter, request *http.Request)
	ActivateOperationLevel(writer http.ResponseWriter, request *http.Request)
	DeactivateOperationLevel(writer http.ResponseWriter, request *http.Request)
}

type OperationModelMappingControllerImpl struct {
	operationmodelmappingservice masteroperationservice.OperationModelMappingService
}

func NewOperationModelMappingController(operationModelMappingservice masteroperationservice.OperationModelMappingService) OperationModelMappingController {
	return &OperationModelMappingControllerImpl{
		operationmodelmappingservice: operationModelMappingservice,
	}
}

// @Summary Get Operation Model Mapping Lookup
// @Description Retrieve operation model mapping lookup with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @Param is_active query string false "Is Active"
// @Param operation_group_code query string false "Operation Group Code"
// @Param operation_name query string false "Operation Name"
// @Param operation_code query string false "Operation Code"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_by query string false "Field to sort by"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-model-mapping/ [get]
func (r *OperationModelMappingControllerImpl) GetOperationModelMappingLookup(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"mtr_operation_model_mapping.is_active":            request.URL.Query().Get("is_active"),
		"mtr_operation_model_mapping.operation_group_code": request.URL.Query().Get("operation_group_code"),
		"mtr_operation_code.operation_name":                request.URL.Query().Get("operation_name"),
		"mtr_operation_model_mapping.operation_code":       request.URL.Query().Get("operation_code"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	paginatedData, totalPages, totalRows, err := r.operationmodelmappingservice.GetOperationModelMappingLookup(criteria, paginate)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// @Summary Get Operation Model Mapping By ID
// @Description Retrieve an operation model mapping by its ID
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @Param operation_model_mapping_id path int true "Operation Model Mapping ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-model-mapping/{operation_model_mapping_id} [get]
func (r *OperationModelMappingControllerImpl) GetOperationModelMappingById(writer http.ResponseWriter, request *http.Request) {
	operationModelMappingID, errA := strconv.Atoi(chi.URLParam(request, "operation_model_mapping_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.operationmodelmappingservice.GetOperationModelMappingById(operationModelMappingID)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Operation Model Mapping By Brand Model Operation Code
// @Description Retrieve an operation model mapping by brand, model, and operation codes
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @Param brand_id query int true "Brand ID"
// @Param model_id query int true "Model ID"
// @Param operation_id query int true "Operation ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-model-mapping/lookup [get]
func (r *OperationModelMappingControllerImpl) GetOperationModelMappingByBrandModelOperationCode(writer http.ResponseWriter, request *http.Request) {

	brandID, errA := strconv.Atoi(request.URL.Query().Get("brand_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	modelID, errA := strconv.Atoi(request.URL.Query().Get("model_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	operationID, errA := strconv.Atoi(request.URL.Query().Get("operation_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.operationmodelmappingservice.GetOperationModelMappingByBrandModelOperationCode(masteroperationpayloads.OperationModelModelBrandOperationCodeRequest{
		BrandId:     brandID,
		ModelId:     modelID,
		OperationId: operationID,
	})

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Operation Model Mapping
// @Description Create or update an operation model mapping
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @Param reqBody body masteroperationpayloads.OperationModelMappingResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-model-mapping/ [post]
func (r *OperationModelMappingControllerImpl) SaveOperationModelMapping(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteroperationpayloads.OperationModelMappingResponse
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	var message string

	create, err := r.operationmodelmappingservice.SaveOperationModelMapping(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	if formRequest.OperationModelMappingId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Operation Model Mapping
// @Description Change the status of an operation model mapping by its ID
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @Param operation_model_mapping_id path int true "Operation Model Mapping ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-model-mapping/{operation_model_mapping_id} [patch]
func (r *OperationModelMappingControllerImpl) ChangeStatusOperationModelMapping(writer http.ResponseWriter, request *http.Request) {
	operationModelMappingID, errA := strconv.Atoi(chi.URLParam(request, "operation_model_mapping_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.operationmodelmappingservice.ChangeStatusOperationModelMapping(operationModelMappingID)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Change Status Successfully!", http.StatusOK)
}

// @Summary Save Operation Model Mapping FRT
// @Description Create or update an operation model mapping FRT (Fixed Repair Time)
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @Param reqBody body masteroperationpayloads.OperationModelMappingFrtRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-model-mapping/frt [post]
func (r *OperationModelMappingControllerImpl) SaveOperationModelMappingFrt(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteroperationpayloads.OperationModelMappingFrtRequest
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	var message string

	create, err := r.operationmodelmappingservice.SaveOperationModelMappingFrt(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	if formRequest.OperationFrtId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Deactivate Operation FRT
// @Description Deactivate one or more operation FRTs by their IDs
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @Param operation_frt_id path string true "Operation FRT ID(s) to deactivate, comma-separated"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-model-mapping/frt/deactivate/{operation_frt_id} [patch]
func (r *OperationModelMappingControllerImpl) DeactivateOperationFrt(writer http.ResponseWriter, request *http.Request) {

	OperationFrtIds := chi.URLParam(request, "operation_frt_id")
	response, err := r.operationmodelmappingservice.DeactivateOperationFrt(OperationFrtIds)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Activate Operation FRT
// @Description Activate one or more deactivated operation FRTs by their IDs
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @Param operation_frt_id path string true "Operation FRT ID(s) to activate, comma-separated"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-model-mapping/frt/activate/{operation_frt_id} [patch]
func (r *OperationModelMappingControllerImpl) ActivateOperationFrt(writer http.ResponseWriter, request *http.Request) {

	OperationFrtIds := chi.URLParam(request, "operation_frt_id")
	response, err := r.operationmodelmappingservice.ActivateOperationFrt(OperationFrtIds)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Get All Operation Document Requirement
// @Description Retrieve all operation document requirements associated with a specific operation model mapping
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @Param operation_model_mapping_id path int true "Operation Model Mapping ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_by query string false "Field to sort by"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-model-mapping/{operation_model_mapping_id}/document-requirements [get]
func (r *OperationModelMappingControllerImpl) GetAllOperationDocumentRequirement(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	headerId, errA := strconv.Atoi(chi.URLParam(request, "operation_model_mapping_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result, err := r.operationmodelmappingservice.GetAllOperationDocumentRequirement(headerId, paginate)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get All Operation FRT
// @Description Retrieve all operation FRTs associated with a specific operation model mapping
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @Param operation_model_mapping_id path int true "Operation Model Mapping ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_by query string false "Field to sort by"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-model-mapping/{operation_model_mapping_id}/frt [get]
func (r *OperationModelMappingControllerImpl) GetAllOperationFrt(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	headerId, errA := strconv.Atoi(chi.URLParam(request, "operation_model_mapping_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	results, totalPages, totalRows, err := r.operationmodelmappingservice.GetAllOperationFrt(headerId, paginate)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(results), "Get Data Successfully!", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// @Summary Get Operation Document Requirement By ID
// @Description Retrieve an operation document requirement by its ID
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @Param operation_document_requirement_id path int true "Operation Document Requirement ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-model-mapping/document-requirements/{operation_document_requirement_id} [get]
func (r *OperationModelMappingControllerImpl) GetOperationDocumentRequirementById(writer http.ResponseWriter, request *http.Request) {
	operationDocumentRequirementId, errA := strconv.Atoi(chi.URLParam(request, "operation_document_requirement_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.operationmodelmappingservice.GetOperationDocumentRequirementById(operationDocumentRequirementId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Operation FRT By ID
// @Description Retrieve an operation FRT by its ID
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @Param operation_frt_id path int true "Operation FRT ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-model-mapping/frt/{operation_frt_id} [get]
func (r *OperationModelMappingControllerImpl) GetOperationFrtById(writer http.ResponseWriter, request *http.Request) {
	OperationFrtId, errA := strconv.Atoi(chi.URLParam(request, "operation_frt_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.operationmodelmappingservice.GetOperationFrtById(OperationFrtId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Operation Model Mapping Document Requirement
// @Description Create or update an operation model mapping document requirement
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @Param reqBody body masteroperationpayloads.OperationModelMappingDocumentRequirementRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-model-mapping/document-requirements [post]
func (r *OperationModelMappingControllerImpl) SaveOperationModelMappingDocumentRequirement(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteroperationpayloads.OperationModelMappingDocumentRequirementRequest
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	var message string

	create, err := r.operationmodelmappingservice.SaveOperationModelMappingDocumentRequirement(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	if formRequest.OperationDocumentRequirementId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}
func (r *OperationModelMappingControllerImpl) DeactivateOperationDocumentRequirement(writer http.ResponseWriter, request *http.Request) {

	OperationDocReqIds := chi.URLParam(request, "operation_document_requirement_id")
	response, err := r.operationmodelmappingservice.DeactivateOperationDocumentRequirement(OperationDocReqIds)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
func (r *OperationModelMappingControllerImpl) ActivateOperationDocumentRequirement(writer http.ResponseWriter, request *http.Request) {

	OperationDocReqIds := chi.URLParam(request, "operation_document_requirement_id")
	response, err := r.operationmodelmappingservice.ActivateOperationDocumentRequirement(OperationDocReqIds)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

func (r *OperationModelMappingControllerImpl) SaveOperationLevel(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteroperationpayloads.OperationLevelRequest
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	var message string

	create, err := r.operationmodelmappingservice.SaveOperationLevel(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	if formRequest.OperationLevelId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *OperationModelMappingControllerImpl) GetAllOperationLevel(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	headerId, errA := strconv.Atoi(chi.URLParam(request, "operation_model_mapping_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	results, err := r.operationmodelmappingservice.GetAllOperationLevel(headerId, paginate)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, results.Rows, "Get Data Successfully!", 200, paginate.Limit, paginate.Page, results.TotalRows, results.TotalPages)
}

func (r *OperationModelMappingControllerImpl) GetOperationLevelById(writer http.ResponseWriter, request *http.Request) {
	operationlevelid, errA := strconv.Atoi(chi.URLParam(request, "operation_level_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.operationmodelmappingservice.GetOperationLevelById(operationlevelid)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK)
}

func (r *OperationModelMappingControllerImpl) DeactivateOperationLevel(writer http.ResponseWriter, request *http.Request) {

	OperationLevelIds := chi.URLParam(request, "operation_level_id")
	response, err := r.operationmodelmappingservice.DeactivateOperationLevel(OperationLevelIds)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
func (r *OperationModelMappingControllerImpl) ActivateOperationLevel(writer http.ResponseWriter, request *http.Request) {

	OperationLevelIds := chi.URLParam(request, "operation_level_id")
	response, err := r.operationmodelmappingservice.ActivateOperationLevel(OperationLevelIds)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
