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

type WarrantyFreeServiceController interface {
	GetAllWarrantyFreeService(writer http.ResponseWriter, request *http.Request)
	GetWarrantyFreeServiceByID(writer http.ResponseWriter, request *http.Request)
	SaveWarrantyFreeService(writer http.ResponseWriter, request *http.Request)
	ChangeStatusWarrantyFreeService(writer http.ResponseWriter, request *http.Request)
}
type WarrantyFreeServiceControllerImpl struct {
	WarrantyFreeServiceService masterservice.WarrantyFreeServiceService
}

func NewWarrantyFreeServiceController(warrantyFreeServiceService masterservice.WarrantyFreeServiceService) WarrantyFreeServiceController {
	return &WarrantyFreeServiceControllerImpl{
		WarrantyFreeServiceService: warrantyFreeServiceService,
	}
}

func (r *WarrantyFreeServiceControllerImpl) GetAllWarrantyFreeService(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_warranty_free_service.is_active":      queryValues.Get("is_active"),
		"mtr_warranty_free_service.effective_date": queryValues.Get("effective_date"),
		"brand_code":                      queryValues.Get("brand_code"),
		"model_code":                      queryValues.Get("model_code"),
		"warranty_free_service_type_code": queryValues.Get("warranty_free_service_type_code"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows := r.WarrantyFreeServiceService.GetAllWarrantyFreeService(criteria, paginate)

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *WarrantyFreeServiceControllerImpl) GetWarrantyFreeServiceByID(writer http.ResponseWriter, request *http.Request) {

	warrantyFreeServiceId, _ := strconv.Atoi(chi.URLParam(request, "warranty_free_services_id"))

	result := r.WarrantyFreeServiceService.GetWarrantyFreeServiceById(warrantyFreeServiceId)

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK)
}

func (r *WarrantyFreeServiceControllerImpl) SaveWarrantyFreeService(writer http.ResponseWriter, request *http.Request) {

	var formRequest masterpayloads.WarrantyFreeServiceRequest
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.WarrantyFreeServiceService.SaveWarrantyFreeService(formRequest)

	if formRequest.WarrantyFreeServicesId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *WarrantyFreeServiceControllerImpl) ChangeStatusWarrantyFreeService(writer http.ResponseWriter, request *http.Request) {

	warrantyFreeServiceId, _ := strconv.Atoi(chi.URLParam(request, "warranty_free_services_id"))

	response := r.WarrantyFreeServiceService.ChangeStatusWarrantyFreeService(int(warrantyFreeServiceId))

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
