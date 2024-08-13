package transactionjpcbcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionjpcbservice "after-sales/api/services/transaction/JPCB"
	"after-sales/api/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CarWashController interface {
	GetAllCarWash(writer http.ResponseWriter, request *http.Request)
}

type CarWashControllerImpl struct {
	carWashService transactionjpcbservice.CarWashService
}

func NewCarWashController(carWashService transactionjpcbservice.CarWashService) CarWashController {
	return &CarWashControllerImpl{
		carWashService: carWashService,
	}
}

// GetAllCarWash implements CarWashController.
func (r *CarWashControllerImpl) GetAllCarWash(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		//TODO
		"trx_car_wash.company_id":                   queryValues.Get("company_id"),
		"trx_work_order.work_order_document_number": queryValues.Get("work_order_document_number"),
		"trx_work_order.promise_time":               queryValues.Get("promise_time"),
		//filter by tnkb, tnkb is from another service
		"trx_car_wash.car_wash_bay_id":    queryValues.Get("car_wash_bay_id"),
		"trx_car_wash.car_wash_status_id": queryValues.Get("car_wash_status_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: chi.URLParam(request, "sort_of"),
		SortBy: chi.URLParam(request, "sort_by"),
	}
	print(queryParams)

	criteria := utils.BuildFilterCondition(queryParams)
	paginatedData, totalPages, totalRows, err := r.carWashService.GetAll(criteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}
