package masteroperationcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads"
	"after-sales/api/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	masteroperationpayloads "after-sales/api/payloads/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"
)

type OperationModelMappingController struct {
	operationmodelmappingservice masteroperationservice.OperationModelMappingService
}

func StartOperationModelMappingRoutes(
	db *gorm.DB,
	r chi.Router,
	operationmodelmappingservice masteroperationservice.OperationModelMappingService,
) {
	handler := &OperationModelMappingController{operationmodelmappingservice: operationmodelmappingservice}

	r.Get("/operation-model-mapping/", handler.GetOperationModelMappingLookup)
	r.Get("/operation-model-mapping/{operation_model_mapping_id}", handler.GetOperationModelMappingById)
	r.Get("/operation-model-mapping-by-brand-model-operation-id/", handler.GetOperationModelMappingByBrandModelOperationCode)
	r.Post("/operation-model-mapping/", handler.SaveOperationModelMapping)
	r.Patch("/operation-model-mapping/{operation_model_mapping_id}", handler.ChangeStatusOperationModelMapping)
}

func (r *OperationModelMappingController) GetOperationModelMappingLookup(w http.ResponseWriter, req *http.Request) {
	trxHandle := req.Context().Value("db_trx").(*gorm.DB)

	queryParams := map[string]string{
		"mtr_operation_model_mapping.is_active":            req.URL.Query().Get("is_active"),
		"mtr_operation_model_mapping.operation_group_code": req.URL.Query().Get("operation_group_code"),
		"mtr_operation_code.operation_name":                req.URL.Query().Get("operation_name"),
		"mtr_operation_model_mapping.operation_code":       req.URL.Query().Get("operation_code"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	result, err := r.operationmodelmappingservice.WithTrx(trxHandle).GetOperationModelMappingLookup(criteria)
	if err != nil {
		exceptions.NotFoundException(w, err.Error())
		return
	}

	payloads.NewHandleSuccess(w, result, "success", http.StatusOK)
}

func (r *OperationModelMappingController) GetOperationModelMappingById(w http.ResponseWriter, req *http.Request) {
	trxHandle := req.Context().Value("db_trx").(*gorm.DB)
	operationModelMappingID, _ := strconv.Atoi(chi.URLParam(req, "operation_model_mapping_id"))

	result, err := r.operationmodelmappingservice.WithTrx(trxHandle).GetOperationModelMappingById(operationModelMappingID)
	if err != nil {
		exceptions.NotFoundException(w, err.Error())
		return
	}

	payloads.NewHandleSuccess(w, result, "Get Data Successfully!", http.StatusOK)
}

func (r *OperationModelMappingController) GetOperationModelMappingByBrandModelOperationCode(w http.ResponseWriter, req *http.Request) {
	trxHandle := req.Context().Value("db_trx").(*gorm.DB)

	brandID, _ := strconv.Atoi(req.URL.Query().Get("brand_id"))
	modelID, _ := strconv.Atoi(req.URL.Query().Get("model_id"))
	operationID, _ := strconv.Atoi(req.URL.Query().Get("operation_id"))

	result, err := r.operationmodelmappingservice.WithTrx(trxHandle).GetOperationModelMappingByBrandModelOperationCode(masteroperationpayloads.OperationModelModelBrandOperationCodeRequest{
		BrandId:     brandID,
		ModelId:     modelID,
		OperationId: operationID,
	})

	if err != nil {
		exceptions.NotFoundException(w, err.Error())
		return
	}

	payloads.NewHandleSuccess(w, result, "Get Data Successfully!", http.StatusOK)
}

func (r *OperationModelMappingController) SaveOperationModelMapping(w http.ResponseWriter, req *http.Request) {
	trxHandle := req.Context().Value("db_trx").(*gorm.DB)

	var request masteroperationpayloads.OperationModelMappingResponse
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		exceptions.EntityException(w, err.Error())
		return
	}

	if request.OperationModelMappingId != 0 {
		result, err := r.operationmodelmappingservice.WithTrx(trxHandle).GetOperationModelMappingById(int(request.OperationModelMappingId))
		if err != nil {
			exceptions.AppException(w, err.Error())
			return
		}

		if result.OperationModelMappingId == 0 {
			exceptions.NotFoundException(w, err.Error())
			return
		}
	}

	create, err := r.operationmodelmappingservice.WithTrx(trxHandle).SaveOperationModelMapping(request)
	if err != nil {
		exceptions.AppException(w, err.Error())
		return
	}

	message := ""
	if request.OperationModelMappingId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(w, create, message, http.StatusOK)
}

func (r *OperationModelMappingController) ChangeStatusOperationModelMapping(w http.ResponseWriter, req *http.Request) {
	trxHandle := req.Context().Value("db_trx").(*gorm.DB)
	operationModelMappingID, err := strconv.Atoi(chi.URLParam(req, "operation_model_mapping_id"))
	if err != nil {
		exceptions.EntityException(w, err.Error())
		return
	}

	result, err := r.operationmodelmappingservice.WithTrx(trxHandle).GetOperationModelMappingById(operationModelMappingID)
	if err != nil || result.OperationId == 0 {
		exceptions.NotFoundException(w, err.Error())
		return
	}

	response, err := r.operationmodelmappingservice.WithTrx(trxHandle).ChangeStatusOperationModelMapping(operationModelMappingID)
	if err != nil {
		exceptions.AppException(w, err.Error())
		return
	}

	payloads.NewHandleSuccess(w, response, "Change Status Successfully!", http.StatusOK)
}
