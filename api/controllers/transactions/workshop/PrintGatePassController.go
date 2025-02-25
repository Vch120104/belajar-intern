package transactionworkshopcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PrintGatePassController interface {
	GetAll(writer http.ResponseWriter, request *http.Request)
	PrintById(writer http.ResponseWriter, request *http.Request)
}

type PrintGatePassControllerImpl struct {
	PrintGatePassService transactionworkshopservice.PrintGatePassService
}

func NewPrintGatePassController(service transactionworkshopservice.PrintGatePassService) PrintGatePassController {
	return &PrintGatePassControllerImpl{
		PrintGatePassService: service,
	}
}

// GetAll gets all print gate pass
// @Summary Get All Print Gate Pass
// @Description Retrieve all print gate pass with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Print Gate Pass
// @Security AuthorizationKeyAuth
// @Param work_order_system_number query string false "Work Order System Number"
// @Param work_order_date query string false "Work Order Date"
// @Param customer_id query string false "Customer ID"
// @Param tnkb query string false "TNKB"
// @Param print_option query string false "Approval Gatepass Status"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/print-gate-pass [get]
func (r *PrintGatePassControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"work_order_system_number": queryValues.Get("work_order_system_number"),
		"work_order_date":          queryValues.Get("work_order_date"),
		"customer_id":              queryValues.Get("customer_id"),
		"tnkb":                     queryValues.Get("tnkb"),
		"approval_gatepass_status": queryValues.Get("print_option"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, BaseErr := r.PrintGatePassService.GetAll(criteria, paginate)
	if BaseErr != nil {
		exceptions.NewNotFoundException(writer, request, BaseErr)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		result.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		result.Limit,
		result.Page,
		int64(result.TotalRows),
		result.TotalPages,
	)
}

// PrintById generate PDF for print gate pass by id
// @Summary Print Gate Pass By ID
// @Description Generate PDF for print gate pass by id
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Print Gate Pass
// @Security AuthorizationKeyAuth
// @Param id path string true "Print Gate Pass ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/print-gate-pass/{gate_pass_system_number} [get]
func (r *PrintGatePassControllerImpl) PrintById(writer http.ResponseWriter, request *http.Request) {
	gatePassID := chi.URLParam(request, "gate_pass_system_number")

	id, err := strconv.Atoi(gatePassID)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid request ID", http.StatusBadRequest)
		return
	}

	pdfBytes, baseErr := r.PrintGatePassService.PrintById(id)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	writer.Header().Set("Content-Type", "application/pdf")
	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=gate_pass_%s.pdf", gatePassID))
	writer.WriteHeader(http.StatusOK)

	writer.Write(pdfBytes)
}
