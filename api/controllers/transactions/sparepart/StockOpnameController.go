package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type StockOpnameController interface {
	// GetAllStockOpname(writer http.ResponseWriter, request *http.Request)
	GetAllStockOpname(writer http.ResponseWriter, request *http.Request)
	GetLocationList(writer http.ResponseWriter, request *http.Request)
	GetPersonInChargeList(writer http.ResponseWriter, request *http.Request)
	GetItemList(writer http.ResponseWriter, request *http.Request)
	// GetListForOnGoing(writer http.ResponseWriter, request *http.Request)
	GetOnGoingStockOpname(writer http.ResponseWriter, request *http.Request)
	InsertNewStockOpname(writer http.ResponseWriter, request *http.Request)
	UpdateOnGoingStockOpname(writer http.ResponseWriter, request *http.Request)
}

type StockOpnameControllerImpl struct {
	Service transactionsparepartservice.StockOpnameService
}

func NewStockOpnameControllerImpl(service transactionsparepartservice.StockOpnameService) StockOpnameController {
	return &StockOpnameControllerImpl{Service: service}
}

func (c *StockOpnameControllerImpl) GetAllStockOpname(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	filterCondition := map[string]string{
		"atStockOpname0.stock_opname_doc_no": queryValues.Get("StockOpnameNo"),
		"b.description":                      queryValues.Get("WarehoseGroup"),
		"c.warehouse_name":                   queryValues.Get("WarehouseCode"),
	}

	dateParams := make(map[string]interface{})
	if queryValues.Get("DateFrom") != "" {
		stockOpnameFrom := queryValues.Get("DateFrom")
		parsedDate, err := time.Parse("2006-01-02T15:04:05Z", stockOpnameFrom)
		if err != nil {
			exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
			})
			return
		}
		dateParams["atStockOpname0.EXEC_DATE_FROM"] = parsedDate
	} else if queryValues.Get("DateTo") != "" {
		stockOpnameTo := queryValues.Get("DateTo")
		parsedDate, err := time.Parse("2006-01-02T15:04:05Z", stockOpnameTo)
		if err != nil {
			exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
			})
			return
		}
		dateParams["atStockOpname0.EXEC_DATE_TO"] = parsedDate
	}

	pages := pagination.Pagination{
		Limit: utils.NewGetQueryInt(queryValues, "limit"),
		Page:  utils.NewGetQueryInt(queryValues, "pages"),
	}

	filterConds := utils.BuildFilterCondition(filterCondition)
	companyCodeStr := chi.URLParam(request, "companyCode")
	companyCode, errA := strconv.ParseFloat(companyCodeStr, 64)
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errA,
		})
		return
	}

	res, errB := c.Service.GetAllStockOpname(filterConds, pages, companyCode, dateParams)
	if errB != nil {
		exceptions.NewNotFoundException(writer, request, errB)
		return
	}

	logrus.Debug("data retrieved: ", res.Rows)

	payloads.NewHandleSuccessPagination(
		writer,
		res.Rows,
		"Stock Opname fetched successfully",
		http.StatusOK,
		res.Limit,
		res.Page,
		int64(res.TotalRows),
		res.TotalPages,
	)
}

func (c *StockOpnameControllerImpl) GetLocationList(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	warehouseGroup := chi.URLParam(request, "warehouseGroup")
	warehouseCode := chi.URLParam(request, "warehouseCode")
	companyCodeStr := chi.URLParam(request, "companyCode")
	companyCode, errA := strconv.ParseFloat(companyCodeStr, 64)
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errA,
		})
		return
	}

	filterCondition := map[string]string{
		"location_code": queryValues.Get("locationCode"),
		"location_name": queryValues.Get("warehouseGroup"),
	}

	pages := pagination.Pagination{
		Limit: utils.NewGetQueryInt(queryValues, "limit"),
		Page:  utils.NewGetQueryInt(queryValues, "pages"),
	}

	filterConds := utils.BuildFilterCondition(filterCondition)

	res, err := c.Service.GetLocationList(filterConds, pages, companyCode, warehouseGroup, warehouseCode)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	logrus.Debug("data retrieved: ", res.Rows)

	payloads.NewHandleSuccessPagination(
		writer,
		res.Rows,
		"Location list fetched successfully",
		http.StatusOK,
		res.Limit,
		res.Page,
		int64(res.TotalRows),
		res.TotalPages,
	)
}

func (c *StockOpnameControllerImpl) GetPersonInChargeList(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	filterCondition := map[string]string{
		"gmemp.employee_no":   queryValues.Get("EmployeeNo"),
		"gmemp.employee_name": queryValues.Get("EmployeeName"),
		"b.description":       queryValues.Get("Position"),
	}

	pages := pagination.Pagination{
		Limit: utils.NewGetQueryInt(queryValues, "limit"),
		Page:  utils.NewGetQueryInt(queryValues, "pages"),
	}

	filterConds := utils.BuildFilterCondition(filterCondition)
	companyCodeStr := chi.URLParam(request, "companyCode")
	companyCode, errA := strconv.ParseFloat(companyCodeStr, 64)
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errA,
		})
		return
	}

	res, errB := c.Service.GetPersonInChargeList(filterConds, pages, companyCode)
	if errB != nil {
		exceptions.NewNotFoundException(writer, request, errB)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		res.Rows,
		"Person in charge list fetched successfully",
		http.StatusOK,
		res.Limit,
		res.Page,
		int64(res.TotalRows),
		res.TotalPages,
	)
}

func (c *StockOpnameControllerImpl) GetItemList(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	whsCode := chi.URLParam(request, "whsCode")
	itemGroup := chi.URLParam(request, "itemGroup")

	pages := pagination.Pagination{
		Limit: utils.NewGetQueryInt(queryValues, "limit"),
		Page:  utils.NewGetQueryInt(queryValues, "pages"),
	}

	res, err := c.Service.GetItemList(pages, whsCode, itemGroup)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		res.Rows,
		"Item list fetched successfully",
		http.StatusOK,
		res.Limit,
		res.Page,
		int64(res.TotalRows),
		res.TotalPages,
	)
}

// func (c *StockOpnameControllerImpl) GetListForOnGoing(writer http.ResponseWriter, request *http.Request) {

// 	sysNo := chi.URLParam(request, "sysNo")

// 	err := c.Service.GetListForOnGoing(sysNo)
// 	if err != nil {
// 		exceptions.NewNotFoundException(writer, request, err)
// 		return
// 	}
// 	payloads.NewHandleSuccess(writer, nil, "Stock Opname is ongoing", http.StatusOK)
// }

func (c *StockOpnameControllerImpl) GetOnGoingStockOpname(writer http.ResponseWriter, request *http.Request) {

	companyCodeStr := chi.URLParam(request, "companyCode")
	companyCode, errA := strconv.ParseFloat(companyCodeStr, 64)

	sysNoStr := chi.URLParam(request, "sysNo")
	sysNo, errB := strconv.ParseFloat(sysNoStr, 64)

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errA,
		})
		return
	}

	if errB != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errB,
		})
		return
	}
	data, err := c.Service.GetOnGoingStockOpname(companyCode, sysNo)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, data, "Stock Opname is ongoing", http.StatusOK)
}

func (c *StockOpnameControllerImpl) InsertNewStockOpname(writer http.ResponseWriter, request *http.Request) {
	var newRequest transactionsparepartpayloads.InsertNewStockOpnameRequest

	helper.ReadFromRequestBody(request, &newRequest)
	if validationErr := validation.ValidationForm(writer, request, &newRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	isTrue, err := c.Service.InsertNewStockOpname(newRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, isTrue, "Stock Opname inserted successfully", http.StatusOK)
}

func (c *StockOpnameControllerImpl) UpdateOnGoingStockOpname(writer http.ResponseWriter, request *http.Request) {
	var newRequest transactionsparepartpayloads.InsertNewStockOpnameRequest

	helper.ReadFromRequestBody(request, &newRequest)
	if validationErr := validation.ValidationForm(writer, request, &newRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	sysNoStr := chi.URLParam(request, "sysNo")
	sysNo, errA := strconv.ParseFloat(sysNoStr, 64)
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errA,
		})
		return
	}

	isTrue, err := c.Service.UpdateOnGoingStockOpname(sysNo, newRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, isTrue, "Stock Opname updated successfully", http.StatusOK)
}
