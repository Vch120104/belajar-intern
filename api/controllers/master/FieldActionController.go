package mastercontroller

import (
	// "after-sales/api/helper"
	"after-sales/api/payloads"
	// masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	// "after-sales/api/middlewares"

	"net/http"
	// "strconv"

	"github.com/julienschmidt/httprouter"
)

type FieldActionController interface {
	GetAllFieldAction(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	// GetAllFieldActionIsActive(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	// GetFieldActionByCode(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	// SaveFieldAction(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	// ChangeStatusFieldAction(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
type FieldActionControllerImpl struct {
	FieldActionService masterservice.FieldActionService
}

func NewFieldActionController(FieldActionService masterservice.FieldActionService) FieldActionController {
	return &FieldActionControllerImpl{
		FieldActionService: FieldActionService,
	}
}

// @Summary Get All Field Action
// @Description REST API Field Action
// @Accept json
// @Produce json
// @Tags Master : Field Action
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param field_action_system_number query string false "field_action_system_number"
// @Param field_action_document_number query string false "field_action_document_number"
// @Param approval_value query string false "approval_value"
// @Param brand_id query string false "brand_id"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-group [get]
func (r *FieldActionControllerImpl) GetAllFieldAction(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	queryParams := map[string]string{
		"is_active":                   params.ByName("is_active"),
		"field_action_system_number":        params.ByName("field_action_system_number"),
		"field_action_document_number": params.ByName("field_action_document_number"),
		"brand_id": params.ByName("brand_id"),
		"approval_value": params.ByName("approval_value"),

	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(params, "limit"),
		Page:   utils.NewGetQueryInt(params, "page"),
		SortOf: params.ByName("sort_of"),
		SortBy: params.ByName("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result := r.FieldActionService.GetAllFieldAction(filterCondition, pagination)

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}
