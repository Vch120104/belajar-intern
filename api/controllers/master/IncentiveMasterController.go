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

	"github.com/julienschmidt/httprouter"
)

type IncentiveMasterController interface {
	GetAllIncentiveMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetIncentiveMasterById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SaveIncentiveMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ChangeStatusIncentiveMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type IncentiveMasterControllerImpl struct {
	IncentiveMasterService masterservice.IncentiveMasterService
}

func NewIncentiveMasterController(incentiveMasterService masterservice.IncentiveMasterService) IncentiveMasterController {
	return &IncentiveMasterControllerImpl{
		IncentiveMasterService: incentiveMasterService,
	}
}

// @Summary Get All Incentive Master
// @Description REST API Incentive Master
// @Accept json
// @Produce json
// @Tags Master : Incentive Master
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param incentive_master_id query string false "incentive_master_id"
// @Param incentive_master_description query string false "incentive_master_description"
// @Param job_position_name query string false "job_position_name"
// @Param incentive_master_percent query float64 false "incentive_master_percent"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /incentive-master [get]
func (r *IncentiveMasterControllerImpl) GetAllIncentiveMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_incentive_master.incentive_master_id":          queryValues.Get("incentive_master_id"),
		"mtr_incentive_master.incentive_master_description": queryValues.Get("incentive_master_description"),
		"job_position_name":                                 queryValues.Get("job_position_name"),
		"mtr_incentive_master.incentive_master_level":       queryValues.Get("incentive_master_level"),
		"mtr_incentive_master.is_active":                    queryValues.Get("is_active"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows := r.IncentiveMasterService.GetAllIncentiveMaster(criteria, paginate)

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully!", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// @Summary Get Incentive Master By ID
// @Description REST API  Incentive Master
// @Accept json
// @Produce json
// @Tags Master :  Incentive Master
// @Param incentive_master_id path int true "incentive_master_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /incentive-master/{incentive_master_id} [get]
func (r *IncentiveMasterControllerImpl) GetIncentiveMasterById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	incentiveMasterId, _ := strconv.Atoi(params.ByName("incentive_master_id"))

	result := r.IncentiveMasterService.GetIncentiveMasterById(incentiveMasterId)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Incentive Master
// @Description REST API Incentive Master
// @Accept json
// @Produce json
// @Tags Master : Incentive Master
// @param reqBody body masterpayloads.IncentiveMasterRequest true "Form Request"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /incentive-master [post]
func (r *IncentiveMasterControllerImpl) SaveIncentiveMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var formRequest masterpayloads.IncentiveMasterRequest
	var message = ""
	helper.ReadFromRequestBody(request, &formRequest)

	create := r.IncentiveMasterService.SaveIncentiveMaster(formRequest)

	if formRequest.IncentiveMasterId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Incentive Master
// @Description REST API Incentive Master
// @Accept json
// @Produce json
// @Tags Master : Incentive Master
// @param incentive_master_id path int true "incentive_master_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /incentive-master/{incentive_master_id} [patch]
func (r *IncentiveMasterControllerImpl) ChangeStatusIncentiveMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	incentiveMasterId, _ := strconv.Atoi(params.ByName("incentive_master_id"))

	response := r.IncentiveMasterService.ChangeStatusIncentiveMaster(int(incentiveMasterId))

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
