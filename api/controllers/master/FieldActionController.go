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

// @Summary Get All Operation Group
// @Description REST API Operation Group
// @Accept json
// @Produce json
// @Tags Master : Operation Group
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param operation_group_code query string false "operation_group_code"
// @Param operation_group_description query string false "operation_group_description"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-group [get]
func (r *FieldActionControllerImpl) GetAllFieldAction(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	queryParams := map[string]string{
		"operation_group_code":        params.ByName("operation_group_code"),
		"operation_group_description": params.ByName("operation_group_description"),
		"is_active":                   params.ByName("is_active"),
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
