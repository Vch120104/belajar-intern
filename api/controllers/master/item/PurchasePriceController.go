package masteritemcontroller

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	helper_test "after-sales/api/helper_testt"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PurchasePriceController interface {
	GetAllPurchasePrice(writer http.ResponseWriter, request *http.Request)
	SavePurchasePrice(writer http.ResponseWriter, request *http.Request)
	GetPurchasePriceById(writer http.ResponseWriter, request *http.Request)
	//GetAllPurchasePriceDetail(writer http.ResponseWriter, request *http.Request)
	//PopupPurchasePrice(writer http.ResponseWriter, request *http.Request)
	//AddPurchasePrice(writer http.ResponseWriter, request *http.Request)
	//DeletePurchasePrice(writer http.ResponseWriter, request *http.Request)
}

type PurchasePriceControllerImpl struct {
	PurchasePriceService masteritemservice.PurchasePriceService
}

func NewPurchasePriceController(PurchasePriceService masteritemservice.PurchasePriceService) PurchasePriceController {
	return &PurchasePriceControllerImpl{
		PurchasePriceService: PurchasePriceService,
	}
}

// @Summary Get All Item Location
// @Description REST API Item Location
// @Accept json
// @Produce json
// @Tags Master : Item Location
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param item_name query int false "item_name"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /item-location [get]
func (r *PurchasePriceControllerImpl) GetAllPurchasePrice(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_purchase_price.purchase_price_id": queryValues.Get("purchase_price_id"),
		"mtr_purchase_price.supplier_id":       queryValues.Get("supplier_id"),
		"mtr_purchase_price.currency_id":       queryValues.Get("currency_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.PurchasePriceService.GetAllPurchasePrice(criteria, paginate)
	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// @Summary Save Item Location
// @Description REST API Item Location
// @Accept json
// @Produce json
// @Tags Master :Item Location
// @param reqBody body masteritempayloads.PurchasePriceResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router / [post]
func (r *PurchasePriceControllerImpl) SavePurchasePrice(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.PurchasePriceRequest
	var message = ""
	helper.ReadFromRequestBody(request, &formRequest)

	create, err := r.PurchasePriceService.SavePurchasePrice(formRequest)
	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	}
	if formRequest.PurchasePriceId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Get Item Location By ID
// @Description REST API  Item Location
// @Accept json
// @Produce json
// @Tags Master :  Item Location
// @Param purchase_price_id path int true "purchase_price_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /{purchase_price_id} [get]
func (r *PurchasePriceControllerImpl) GetPurchasePriceById(writer http.ResponseWriter, request *http.Request) {

	PurchasePriceIds, _ := strconv.Atoi(chi.URLParam(request, "purchase_price_id"))

	result, err := r.PurchasePriceService.GetPurchasePriceById(PurchasePriceIds)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}
