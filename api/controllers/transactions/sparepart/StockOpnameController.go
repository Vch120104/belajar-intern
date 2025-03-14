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
)

type StockOpnameController interface {
	GetAllStockOpname(http.ResponseWriter, *http.Request)
	GetAllStockOpnameDetail(http.ResponseWriter, *http.Request)
	GetStockOpnameByStockOpnameSystemNumber(http.ResponseWriter, *http.Request)
	GetStockOpnameAllDetailByStockOpnameSystemNumber(http.ResponseWriter, *http.Request)
	InsertStockOpname(http.ResponseWriter, *http.Request)
	SubmitStockOpname(http.ResponseWriter, *http.Request)
	InsertStockOpnameDetail(http.ResponseWriter, *http.Request)
	UpdateStockOpname(http.ResponseWriter, *http.Request)
	UpdateStockOpnameDetail(http.ResponseWriter, *http.Request)
	DeleteStockOpname(http.ResponseWriter, *http.Request)
	DeleteStockOpnameDetailByLineNumber(http.ResponseWriter, *http.Request)
	DeleteStockOpnameDetailBySystemNumber(http.ResponseWriter, *http.Request)

	// // GetAllStockOpname(writer http.ResponseWriter, request *http.Request)
	// GetAllStockOpname(writer http.ResponseWriter, request *http.Request)
	// GetLocationList(writer http.ResponseWriter, request *http.Request)
	// GetPersonInChargeList(writer http.ResponseWriter, request *http.Request)
	// GetItemList(writer http.ResponseWriter, request *http.Request)
	// // GetListForOnGoing(writer http.ResponseWriter, request *http.Request)
	// GetOnGoingStockOpname(writer http.ResponseWriter, request *http.Request)
	// InsertNewStockOpname(writer http.ResponseWriter, request *http.Request)
	// UpdateOnGoingStockOpname(writer http.ResponseWriter, request *http.Request)
}

type StockOpnameControllerImpl struct {
	Service transactionsparepartservice.StockOpnameService
}

func NewStockOpnameControllerImpl(service transactionsparepartservice.StockOpnameService) StockOpnameController {
	return &StockOpnameControllerImpl{Service: service}
}

func (c *StockOpnameControllerImpl) GetAllStockOpname(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	filteredCondition := map[string]string{
		"trx_stock_opname.stock_opname_document_number": queryValues.Get("stock_opname_document_number"),
		"B.warehouse_location_group_name":               queryValues.Get("warehouse_location_group_name"),
		"C.warehouse_name":                              queryValues.Get("warehouse_name"),
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

	filterConds := utils.BuildFilterCondition(filteredCondition)

	res, err := c.Service.GetAllStockOpname(filterConds, pages, dateParams)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(
		writer,
		res.Rows,
		"All Stock Opname fetched successfully",
		http.StatusOK,
		res.Limit,
		res.Page,
		int64(res.TotalRows),
		res.TotalPages,
	)

}

func (c *StockOpnameControllerImpl) GetAllStockOpnameDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	pages := pagination.Pagination{
		Limit: utils.NewGetQueryInt(queryValues, "limit"),
		Page:  utils.NewGetQueryInt(queryValues, "pages"),
	}

	res, err := c.Service.GetAllStockOpnameDetail(pages)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		res.Rows,
		"All Stock Opname Detail fetched successfully",
		http.StatusOK,
		res.Limit,
		res.Page,
		int64(pages.TotalRows),
		pages.TotalPages,
	)
}

func (c *StockOpnameControllerImpl) GetStockOpnameByStockOpnameSystemNumber(writer http.ResponseWriter, request *http.Request) {
	stockOpnameSystemNumber, errA := strconv.Atoi(chi.URLParam(request, "stock_opname_system_number"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errA,
		})
		return
	}

	res, err := c.Service.GetStockOpnameByStockOpnameSystemNumber(stockOpnameSystemNumber)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, res, "Stock Opname fetched successfully", http.StatusOK)
}

func (c *StockOpnameControllerImpl) GetStockOpnameAllDetailByStockOpnameSystemNumber(writer http.ResponseWriter, request *http.Request) {
	stockOpnameSystemNumberStr := chi.URLParam(request, "stock_opname_system_number")
	stockOpnameSystemNumber, errA := strconv.Atoi(stockOpnameSystemNumberStr)
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errA,
		})
		return
	}

	queryValues := request.URL.Query()
	pages := pagination.Pagination{
		Limit: utils.NewGetQueryInt(queryValues, "limit"),
		Page:  utils.NewGetQueryInt(queryValues, "pages"),
	}

	res, err := c.Service.GetStockOpnameAllDetailByStockOpnameSystemNumber(stockOpnameSystemNumber, pages)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		res.Rows,
		"Stock Opname Detail fetched successfully",
		http.StatusOK,
		res.Limit,
		res.Page,
		int64(res.TotalRows),
		res.TotalPages,
	)
}

func (c *StockOpnameControllerImpl) InsertStockOpname(writer http.ResponseWriter, request *http.Request) {
	var transferRequest transactionsparepartpayloads.StockOpnameInsertRequest

	helper.ReadFromRequestBody(request, &transferRequest)
	if validationErr := validation.ValidationForm(writer, request, &transferRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	isTrue, err := c.Service.InsertStockOpname(transferRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, isTrue, "Stock Opname inserted successfully", http.StatusOK)
}

func (c *StockOpnameControllerImpl) SubmitStockOpname(writer http.ResponseWriter, request *http.Request) {
	var submitRequest transactionsparepartpayloads.StockOpnameSubmitRequest

	helper.ReadFromRequestBody(request, &submitRequest)
	if validationErr := validation.ValidationForm(writer, request, &submitRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	systemNumberStr := chi.URLParam(request, "stock_opname_system_number")
	systemNumber, errA := strconv.Atoi(systemNumberStr)
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errA,
		})
		return
	}

	isTrue, err := c.Service.SubmitStockOpname(systemNumber, submitRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, isTrue, "Stock Opname submitted successfully", http.StatusOK)
}

func (c *StockOpnameControllerImpl) InsertStockOpnameDetail(writer http.ResponseWriter, request *http.Request) {
	var insertRequest transactionsparepartpayloads.StockOpnameInsertDetailRequest

	helper.ReadFromRequestBody(request, &insertRequest)
	if validationErr := validation.ValidationForm(writer, request, &insertRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	isTrue, err := c.Service.InsertStockOpnameDetail(insertRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, isTrue, "Stock Opname Detail inserted successfully", http.StatusOK)
}

func (c *StockOpnameControllerImpl) UpdateStockOpname(writer http.ResponseWriter, request *http.Request) {
	var updateRequest transactionsparepartpayloads.StockOpnameInsertRequest
	helper.ReadFromRequestBody(request, &updateRequest)

	systemNumberStr := chi.URLParam(request, "stock_opname_system_number")
	systemNumber, errA := strconv.Atoi(systemNumberStr)

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errA,
		})
		return
	}

	istrue, err := c.Service.UpdateStockOpname(updateRequest, systemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, istrue, "Stock Opname Updated", http.StatusOK)
}

func (c *StockOpnameControllerImpl) UpdateStockOpnameDetail(writer http.ResponseWriter, request *http.Request) {
	var updateRequest transactionsparepartpayloads.StockOpnameUpdateDetailRequest
	helper.ReadFromRequestBody(request, &updateRequest)

	systemNumberStr := chi.URLParam(request, "stock_opname_system_number")
	systemNumber, errA := strconv.Atoi(systemNumberStr)

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errA,
		})
		return
	}

	istrue, err := c.Service.UpdateStockOpnameDetail(updateRequest, systemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, istrue, "Stock Opname Detail Updated", http.StatusOK)
}

func (c *StockOpnameControllerImpl) DeleteStockOpname(writer http.ResponseWriter, request *http.Request) {
	systemNumberStr := chi.URLParam(request, "stock_opname_system_number")
	systemNumber, errA := strconv.Atoi(systemNumberStr)

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errA,
		})
		return
	}

	istrue, err := c.Service.DeleteStockOpname(systemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, istrue, "Stock Opname Deleted", http.StatusOK)
}

func (c *StockOpnameControllerImpl) DeleteStockOpnameDetailByLineNumber(writer http.ResponseWriter, request *http.Request) {
	systemNumberStr := chi.URLParam(request, "stock_opname_system_number")
	systemNumber, errA := strconv.Atoi(systemNumberStr)
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errA,
		})
		return
	}
	systemNumberLineStr := chi.URLParam(request, "line_number")
	systemNumberLine, errB := strconv.Atoi(systemNumberLineStr)
	if errB != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errB,
		})
		return
	}

	istrue, err := c.Service.DeleteStockOpnameDetailByLineNumber(systemNumber, systemNumberLine)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, istrue, "Stock Opname Line Deleted", http.StatusOK)
}

func (c *StockOpnameControllerImpl) DeleteStockOpnameDetailBySystemNumber(writer http.ResponseWriter, request *http.Request) {
	systemNumberStr := chi.URLParam(request, "stock_opname_system_number")
	systemNumber, errA := strconv.Atoi(systemNumberStr)
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errA,
		})
		return
	}

	istrue, err := c.Service.DeleteStockOpnameDetailBySystemNumber(systemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, istrue, "Stock Opname System Number Deleted", http.StatusOK)
}
