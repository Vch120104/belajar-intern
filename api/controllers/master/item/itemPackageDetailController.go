package masteritemcontroller

import (
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ItemPackageDetailController interface {
	GetItemPackageDetailByItemPackageId(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type ItemPackageDetailControllerImpl struct {
	ItemPackageDetailService masteritemservice.ItemPackageDetailService
}

func NewItemPackageDetailController(ItemPackageDetailService masteritemservice.ItemPackageDetailService) ItemPackageDetailController {
	return &ItemPackageDetailControllerImpl{
		ItemPackageDetailService: ItemPackageDetailService,
	}
}

func (r *ItemPackageDetailControllerImpl) GetItemPackageDetailByItemPackageId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	queryValues := request.URL.Query()
	itemPackageId, _ := strconv.Atoi(params.ByName("item_package_id"))

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result := r.ItemPackageDetailService.GetItemPackageDetailByItemPackageId(itemPackageId, paginate)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}
