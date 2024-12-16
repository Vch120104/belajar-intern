package mastercontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"time"
)

type LocationStockController interface {
	GetAllLocationStock(writer http.ResponseWriter, request *http.Request)
	UpdateLocationStock(writer http.ResponseWriter, request *http.Request)
	GetAvailableQuantity(writer http.ResponseWriter, request *http.Request)
}
type LocationStockControlerImpl struct {
	LocationStockService masterservice.LocationStockService
}

func NewLocationStockController(LocationStockService masterservice.LocationStockService) LocationStockController {
	return &LocationStockControlerImpl{LocationStockService: LocationStockService}
}

// GetAllLocationStock
//
//	@Summary		Get All Location Stock
//	@Description	REST API Location Stock
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Purchase Request
//	@Param			page					query		string	true	"page"
//	@Param			limit					query		string	true	"limit"
//	@Param			company_id				query		string	false	"company_id"
//	@Param			period_year				query		string	false	"period_year"
//	@Param			period_month			query		string	false	"period_month"
//	@Param			warehouse_id			query		string	false	"warehouse_id"
//	@Param			warehouse_group			query		string	false	"warehouse_group"
//	@Param			quantity_ending			query		string	false	"quantity_ending"
//	@Param			location_id				query		string	false	"location_id"
//	@Param			item_id					query		string	false	"item_id"
//	@Param			sort_by					query		string	false	"sort_by"
//	@Param			sort_of					query		string	false	"sort_of"
//	@Success		200									{object}	[]masterwarehousepayloads.LocationStockDBResponse
//	@Failure		500,400,401,404,403,422				{object}	exceptions.BaseErrorResponse
//	@Router			/v1/location-stock/ [get]
func (l *LocationStockControlerImpl) GetAllLocationStock(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"company_id":      queryValues.Get("company_id"),
		"period_year":     queryValues.Get("period_year"),
		"period_month":    queryValues.Get("period_month"),
		"warehouse_id":    queryValues.Get("warehouse_id"),
		"warehouse_group": queryValues.Get("warehouse_group"),
		"quantity_ending": queryValues.Get("quantity_ending"),
		"location_id":     queryValues.Get("location_id"),
		"item_id":         queryValues.Get("item_id"),
	}
	pageninateparam := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	fiter := utils.BuildFilterCondition(queryParams)
	result, err := l.LocationStockService.GetAllLocationStock(fiter, pageninateparam)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfull", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// location-stock [put]
func (l *LocationStockControlerImpl) UpdateLocationStock(writer http.ResponseWriter, request *http.Request) {
	var locationStockPayloads masterwarehousepayloads.LocationStockUpdatePayloads

	helper.ReadFromRequestBody(request, &locationStockPayloads)
	if validationErr := validation.ValidationForm(writer, request, &locationStockPayloads); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	//if err != nil {
	//	helper.ReturnError(writer, request, err)
	//	return
	//}
	res, err := l.LocationStockService.UpdateLocationStock(locationStockPayloads)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "success to update location stock", http.StatusOK)
}
func (l *LocationStockControlerImpl) GetAvailableQuantity(writer http.ResponseWriter, request *http.Request) {
	filter := masterwarehousepayloads.GetAvailableQuantityPayload{}
	queryValues := request.URL.Query()
	periodDate := queryValues.Get("period_date")
	periodDateParse, errParseDate := time.Parse("2006-01-02T15:04:05.000Z", periodDate)
	if errParseDate != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to parse date",
			//Data:       errParseDate,
			Err: errParseDate,
		})
		return
	}
	filter = masterwarehousepayloads.GetAvailableQuantityPayload{
		CompanyId:        utils.NewGetQueryInt(queryValues, "company_id"),
		PeriodDate:       periodDateParse,
		WarehouseId:      utils.NewGetQueryInt(queryValues, "warehouse_id"),
		LocationId:       utils.NewGetQueryInt(queryValues, "location_id"),
		ItemId:           utils.NewGetQueryInt(queryValues, "item_id"),
		WarehouseGroupId: utils.NewGetQueryInt(queryValues, "warehouse_group_id"),
		UomTypeId:        utils.NewGetQueryInt(queryValues, "uom_id"),
	}
	res, err := l.LocationStockService.GetAvailableQuantity(filter)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "success to get available quantity", http.StatusOK)
}
