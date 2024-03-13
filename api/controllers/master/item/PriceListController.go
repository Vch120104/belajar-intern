package masteritemcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type PriceListController interface {
	GetPriceListLookup(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetPriceList(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SavePriceList(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ChangeStatusPriceList(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type PriceListControllerImpl struct {
	pricelistservice masteritemservice.PriceListService
}

func NewPriceListController(PriceListService masteritemservice.PriceListService) PriceListController {
	return &PriceListControllerImpl{
		pricelistservice: PriceListService,
	}
}

// @Summary Get All Price List Lookup
// @Description REST API Price List
// @Param price_list_code query string false "price_list_code"
// @Param company_id query int false "company_id"
// @Param brand_id query int false "brand_id"
// @Param currency_id query int false "currency_id"
// @Param effective_date query string false "effective_date"
// @Param item_group_id query int false "item_group_id"
// @Param item_class_id query int false "item_class_id"
// @Accept json
// @Produce json
// @Tags Master : Price List
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/price-list/get-all-lookup [get]
func (r *PriceListControllerImpl) GetPriceListLookup(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	queryValues := request.URL.Query()
	PriceListCode := queryValues.Get("price_list_code")
	companyId, _ := strconv.Atoi(queryValues.Get("company_id"))
	brandId, _ := strconv.Atoi(queryValues.Get("brand_id"))
	currencyId, _ := strconv.Atoi(queryValues.Get("currency_id"))
	effectiveDate, _ := time.Parse("2006-01-02T15:04:05.000Z", queryValues.Get("effective_date"))
	itemGroupId, _ := strconv.Atoi(queryValues.Get("item_group_id"))
	itemClassId, _ := strconv.Atoi(queryValues.Get("item_class_id"))

	priceListRequest := masteritempayloads.PriceListGetAllRequest{
		PriceListCode: PriceListCode,
		CompanyId:     int32(companyId),
		BrandId:       int32(brandId),
		CurrencyId:    int32(currencyId),
		EffectiveDate: effectiveDate,
		ItemGroupId:   int32(itemGroupId),
		ItemClassId:   int32(itemClassId),
	}

	result := r.pricelistservice.GetPriceList(priceListRequest)

	payloads.NewHandleSuccess(writer, result, "success", 200)
}

// @Summary Get All Price List
// @Description REST API Price List
// @Param price_list_code query string false "price_list_code"
// @Param company_id query int false "company_id"
// @Param brand_id query int false "brand_id"
// @Param currency_id query int false "currency_id"
// @Param effective_date query string false "effective_date"
// @Param item_id query int false "item_id"
// @Param item_group_id query int false "item_group_id"
// @Param item_class_id query int false "item_class_id"
// @Param price_list_amount query string false "price_list_amount"
// @Param price_list_modifiable query string false "price_list_modifiable" Enums(true, false)
// @Param atpm_syncronize query string false "atpm_syncronize" Enums(true, false)
// @Param atpm_syncronize_time query string false "atpm_syncronize_time"
// @Accept json
// @Produce json
// @Tags Master : Price List
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/price-list/get-all [get]
func (r *PriceListControllerImpl) GetPriceList(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	queryValues := request.URL.Query()
	PriceListCode := queryValues.Get("price_list_code")
	companyId, _ := strconv.Atoi(queryValues.Get("company_id"))
	brandId, _ := strconv.Atoi(queryValues.Get("brand_id"))
	currencyId, _ := strconv.Atoi(queryValues.Get("currency_id"))
	effectiveDate, _ := time.Parse("2006-01-02T15:04:05.000Z", queryValues.Get("effective_date"))
	itemId, _ := strconv.Atoi(queryValues.Get("item_id"))
	itemGroupId, _ := strconv.Atoi(queryValues.Get("item_group_id"))
	itemClassId, _ := strconv.Atoi(queryValues.Get("item_class_id"))
	priceListAmount, _ := strconv.ParseFloat(queryValues.Get("price_list_amount"), 64)
	priceListModifiable := queryValues.Get("price_list_modifiable")
	atpmSyncronize := queryValues.Get("atpm_syncronize")
	atpmSyncronizeTime, _ := time.Parse("2006-01-02T15:04:05.000Z", queryValues.Get("atpm_syncronize_time"))

	priceListRequest := masteritempayloads.PriceListGetAllRequest{
		PriceListCode:       PriceListCode,
		CompanyId:           int32(companyId),
		BrandId:             int32(brandId),
		CurrencyId:          int32(currencyId),
		EffectiveDate:       effectiveDate,
		ItemId:              int32(itemId),
		ItemGroupId:         int32(itemGroupId),
		ItemClassId:         int32(itemClassId),
		PriceListAmount:     priceListAmount,
		PriceListModifiable: priceListModifiable,
		AtpmSyncronize:      atpmSyncronize,
		AtpmSyncronizeTime:  atpmSyncronizeTime,
	}

	result := r.pricelistservice.GetPriceList(priceListRequest)

	payloads.NewHandleSuccess(writer, result, "success", 200)
}

// @Summary Save Price List
// @Description REST API Price List
// @Accept json
// @Produce json
// @Tags Master : Price List
// @param reqBody body masteritempayloads.PriceListResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/price-list [post]
func (r *PriceListControllerImpl) SavePriceList(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var formRequest masteritempayloads.PriceListResponse
	var message = ""

	helper.ReadFromRequestBody(request, &formRequest)

	create := r.pricelistservice.SavePriceList(formRequest)

	if formRequest.PriceListId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Price List
// @Description REST API Price List
// @Accept json
// @Produce json
// @Tags Master : Price List
// @param price_list_id path int true "price_list_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/price-list/{price_list_id} [patch]
func (r *PriceListControllerImpl) ChangeStatusPriceList(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	PriceListId, _ := strconv.Atoi(params.ByName("price_list_id"))

	response := r.pricelistservice.ChangeStatusPriceList(int(PriceListId))

	payloads.NewHandleSuccess(writer, response, "Change Status Successfully!", http.StatusOK)
}
