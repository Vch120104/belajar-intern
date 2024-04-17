package masteroperationcontroller

import (
	"after-sales/api/helper"
	helper_test "after-sales/api/helper_testt"
	"after-sales/api/payloads"
	"after-sales/api/utils"
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
}

type OperationModelMappingControllerImpl struct {
	operationmodelmappingservice masteroperationservice.OperationModelMappingService
}

func NewOperationModelMappingController(operationModelMappingservice masteroperationservice.OperationModelMappingService) OperationModelMappingController {
	return &OperationModelMappingControllerImpl{
		operationmodelmappingservice: operationModelMappingservice,
	}
}

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
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *OperationModelMappingControllerImpl) GetOperationModelMappingById(writer http.ResponseWriter, request *http.Request) {
	operationModelMappingID, _ := strconv.Atoi(chi.URLParam(request, "operation_model_mapping_id"))

	result, err := r.operationmodelmappingservice.GetOperationModelMappingById(operationModelMappingID)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK)
}

func (r *OperationModelMappingControllerImpl) GetOperationModelMappingByBrandModelOperationCode(writer http.ResponseWriter, request *http.Request) {

	brandID, _ := strconv.Atoi(request.URL.Query().Get("brand_id"))
	modelID, _ := strconv.Atoi(request.URL.Query().Get("model_id"))
	operationID, _ := strconv.Atoi(request.URL.Query().Get("operation_id"))

	result, err := r.operationmodelmappingservice.GetOperationModelMappingByBrandModelOperationCode(masteroperationpayloads.OperationModelModelBrandOperationCodeRequest{
		BrandId:     brandID,
		ModelId:     modelID,
		OperationId: operationID,
	})

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *OperationModelMappingControllerImpl) SaveOperationModelMapping(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteroperationpayloads.OperationModelMappingResponse
	helper.ReadFromRequestBody(request, &formRequest)
	var message string

	create, err := r.operationmodelmappingservice.SaveOperationModelMapping(formRequest)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	if formRequest.OperationModelMappingId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *OperationModelMappingControllerImpl) ChangeStatusOperationModelMapping(writer http.ResponseWriter, request *http.Request) {
	operationModelMappingID, _ := strconv.Atoi(chi.URLParam(request, "operation_model_mapping_id"))

	response, err := r.operationmodelmappingservice.ChangeStatusOperationModelMapping(operationModelMappingID)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Change Status Successfully!", http.StatusOK)
}

func (r *OperationModelMappingControllerImpl) SaveOperationModelMappingFrt(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteroperationpayloads.OperationModelMappingFrtRequest
	helper.ReadFromRequestBody(request, &formRequest)
	var message string

	create, err := r.operationmodelmappingservice.SaveOperationModelMappingFrt(formRequest)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	if formRequest.OperationModelMappingId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *OperationModelMappingControllerImpl) DeactivateOperationFrt(writer http.ResponseWriter, request *http.Request) {

	OperationFrtIds := chi.URLParam(request, "operation_frt_id")
	response, err := r.operationmodelmappingservice.DeactivateOperationFrt(OperationFrtIds)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

func (r *OperationModelMappingControllerImpl) ActivateOperationFrt(writer http.ResponseWriter, request *http.Request) {

	OperationFrtIds := chi.URLParam(request, "operation_frt_id")
	response, err := r.operationmodelmappingservice.ActivateOperationFrt(OperationFrtIds)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

func (r *OperationModelMappingControllerImpl) GetAllOperationDocumentRequirement(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	headerId, _ := strconv.Atoi(chi.URLParam(request, "operation_model_mapping_id"))

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result, err := r.operationmodelmappingservice.GetAllOperationDocumentRequirement(headerId, paginate)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

func (r *OperationModelMappingControllerImpl) GetAllOperationFrt(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	headerId, _ := strconv.Atoi(chi.URLParam(request, "operation_model_mapping_id"))

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result, err := r.operationmodelmappingservice.GetAllOperationFrt(headerId, paginate)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

func (r *OperationModelMappingControllerImpl) GetOperationDocumentRequirementById(writer http.ResponseWriter, request *http.Request) {
	operationDocumentRequirementId, _ := strconv.Atoi(chi.URLParam(request, "operation_document_requirement_id"))

	result, err := r.operationmodelmappingservice.GetOperationDocumentRequirementById(operationDocumentRequirementId)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK)
}

func (r *OperationModelMappingControllerImpl) GetOperationFrtById(writer http.ResponseWriter, request *http.Request) {
	OperationFrtId, _ := strconv.Atoi(chi.URLParam(request, "operation_frt_id"))

	result, err := r.operationmodelmappingservice.GetOperationFrtById(OperationFrtId)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK)
}

func (r *OperationModelMappingControllerImpl) SaveOperationModelMappingDocumentRequirement(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteroperationpayloads.OperationModelMappingDocumentRequirementRequest
	helper.ReadFromRequestBody(request, &formRequest)
	var message string

	create, err := r.operationmodelmappingservice.SaveOperationModelMappingDocumentRequirement(formRequest)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	if formRequest.OperationModelMappingId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *OperationModelMappingControllerImpl) DeactivateOperationDocumentRequirement(writer http.ResponseWriter, request *http.Request) {

	OperationFrtIds := chi.URLParam(request, "operation_model_mapping_id")
	response, err := r.operationmodelmappingservice.DeactivateOperationDocumentRequirement(OperationFrtIds)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

func (r *OperationModelMappingControllerImpl) ActivateOperationDocumentRequirement(writer http.ResponseWriter, request *http.Request) {

	OperationFrtIds := chi.URLParam(request, "operation_model_mapping_id")
	response, err := r.operationmodelmappingservice.ActivateOperationDocumentRequirement(OperationFrtIds)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
