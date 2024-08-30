package mastercontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type LookupController interface {
	ItemOprCode(writer http.ResponseWriter, request *http.Request)
	CampaignMaster(writer http.ResponseWriter, request *http.Request)
	ItemOprCodeWithPrice(writer http.ResponseWriter, request *http.Request)
	VehicleUnitMaster(writer http.ResponseWriter, request *http.Request)
}

type LookupControllerImpl struct {
	LookupService masterservice.LookupService
}

func NewLookupController(LookupService masterservice.LookupService) LookupController {
	return &LookupControllerImpl{
		LookupService: LookupService,
	}
}

func (r *LookupControllerImpl) ItemOprCode(writer http.ResponseWriter, request *http.Request) {
	linetypeStrId := chi.URLParam(request, "linetype_id")
	linetypeId, err := strconv.Atoi(linetypeStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Line Type ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()
	queryParams := map[string]string{}
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	criteria := utils.BuildFilterCondition(queryParams)
	lookup, totalPages, totalRows, baseErr := r.LookupService.ItemOprCode(linetypeId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(writer, lookup, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *LookupControllerImpl) ItemOprCodeWithPrice(writer http.ResponseWriter, request *http.Request) {
	linetypeStrId := chi.URLParam(request, "linetype_id")
	linetypeId, err := strconv.Atoi(linetypeStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Line Type ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()
	queryParams := map[string]string{}
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	lookup, totalPages, totalRows, baseErr := r.LookupService.ItemOprCodeWithPrice(linetypeId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(writer, lookup, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *LookupControllerImpl) VehicleUnitMaster(writer http.ResponseWriter, request *http.Request) {
	brandStrId := chi.URLParam(request, "brand_id")
	brandId, err := strconv.Atoi(brandStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Brand ID", http.StatusBadRequest)
		return
	}

	modelStrId := chi.URLParam(request, "model_id")
	modelId, err := strconv.Atoi(modelStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Model ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()
	queryParams := map[string]string{}
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	criteria := utils.BuildFilterCondition(queryParams)
	lookup, totalPages, totalRows, baseErr := r.LookupService.VehicleUnitMaster(brandId, modelId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(writer, lookup, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *LookupControllerImpl) CampaignMaster(writer http.ResponseWriter, request *http.Request) {
	companyStrId := chi.URLParam(request, "company_id")
	companyId, err := strconv.Atoi(companyStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Campaign ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()
	queryParams := map[string]string{}
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	criteria := utils.BuildFilterCondition(queryParams)
	lookup, totalPages, totalRows, baseErr := r.LookupService.CampaignMaster(companyId, paginate, criteria)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Lookup data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccessPagination(writer, lookup, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}
