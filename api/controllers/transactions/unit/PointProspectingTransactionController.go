package pointprospectingtranscontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionunitpayloads "after-sales/api/payloads/transaction/unit"
	transactionunitservice "after-sales/api/services/transaction/unit"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type PointProspectingTransactionController interface {
	GetAllCompanyData(write http.ResponseWriter, request *http.Request)
	GetAllSalesRepresentative(write http.ResponseWriter, request *http.Request)
	GetSalesByCompanyCode(write http.ResponseWriter, request *http.Request)
	Process(write http.ResponseWriter, request *http.Request)
}
type PointProspectingTransactionControllerImpl struct {
	PointProspectingTransactionService transactionunitservice.PointProspectingTransactionService
}

func NewPointProspectingTransactionControllerImpl(service transactionunitservice.PointProspectingTransactionService) PointProspectingTransactionController {
	return &PointProspectingTransactionControllerImpl{
		PointProspectingTransactionService: service,
	}
}

func (c *PointProspectingTransactionControllerImpl) GetAllCompanyData(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	filteredCondition := map[string]string{
		"gmcomp0.company_code": queryValues.Get("COMPANY_CODE"),
		"gmcomp0.company_name": queryValues.Get("COMPANY_NAME"),
	}
	filterConds := utils.BuildFilterCondition(filteredCondition)

	paginate := pagination.Pagination{
		Limit: utils.NewGetQueryInt(queryValues, "limit"),
		Page:  utils.NewGetQueryInt(queryValues, "pages"),
	}

	logrus.Debug("filterConds", filterConds)
	logrus.Debug("paginate", paginate)

	res, err := c.PointProspectingTransactionService.GetAllCompanyData(filterConds, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		res.Rows,
		"All Company Data retrieved successfully",
		http.StatusOK,
		res.Limit,
		res.Page,
		int64(res.TotalRows),
		res.TotalPages,
	)
	fmt.Println(res.TotalRows)
	fmt.Println(res.TotalPages)

}

func (c *PointProspectingTransactionControllerImpl) GetAllSalesRepresentative(write http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	filteredCondition := map[string]string{
		"gmemp.employee_no":   queryValues.Get("employee_no"),
		"gmemp.employee_name": queryValues.Get("employee_name"),
	}

	pages := pagination.Pagination{
		Page:  utils.NewGetQueryInt(queryValues, "page"),
		Limit: utils.NewGetQueryInt(queryValues, "limit"),
	}

	filterConds := utils.BuildFilterCondition(filteredCondition)

	res, err := c.PointProspectingTransactionService.GetAllSalesRepresentative(filterConds, pages)
	if err != nil {
		exceptions.NewNotFoundException(write, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		write,
		res,
		"All Sales Representative Retrieved successfully",
		http.StatusOK,
		res.Limit,
		res.Page,
		int64(res.TotalRows),
		res.TotalPages,
	)

	// fmt.Println(res.TotalRows)
}

func (c *PointProspectingTransactionControllerImpl) GetSalesByCompanyCode(write http.ResponseWriter, request *http.Request) {
	companyCodeStr := chi.URLParam(request, "companyCode")
	companyCode, errA := strconv.ParseFloat(companyCodeStr, 64)

	if errA != nil {
		exceptions.NewBadRequestException(write, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("check your params"),
		})
		return
	}

	filteredCondition := map[string]string{
		"gmemp.employee_no":   request.URL.Query().Get("employee_no"),
		"gmemp.employee_name": request.URL.Query().Get("employee_name"),
	}
	filteredConds := utils.BuildFilterCondition(filteredCondition)

	pages := pagination.Pagination{
		Page:  utils.NewGetQueryInt(request.URL.Query(), "page"),
		Limit: utils.NewGetQueryInt(request.URL.Query(), "limit"),
	}

	res, err := c.PointProspectingTransactionService.GetSalesByCompanyCode(companyCode, filteredConds, pages)
	if err != nil {
		exceptions.NewNotFoundException(write, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		write,
		res,
		"Sales by Company Retrieved successfully",
		http.StatusOK,
		res.Limit,
		res.Page,
		int64(res.TotalRows),
		res.TotalPages,
	)

}

func (c *PointProspectingTransactionControllerImpl) Process(write http.ResponseWriter, request *http.Request) {
	var processRequest transactionunitpayloads.ProcessRequest

	helper.ReadFromRequestBody(request, &processRequest)

	res, err := c.PointProspectingTransactionService.Process(processRequest)
	if err != nil {
		exceptions.NewConflictException(write, request, err)
		return
	}

	payloads.NewHandleSuccess(write, res, "Process Success", http.StatusOK)
}
