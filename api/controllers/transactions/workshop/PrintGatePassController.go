package transactionworkshopcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"net/http"
)

type PrintGatePassController struct {
	PrintGatePassService transactionworkshopservice.PrintGatePassService
}

func NewPrintGatePassController(service transactionworkshopservice.PrintGatePassService) *PrintGatePassController {
	return &PrintGatePassController{
		PrintGatePassService: service,
	}
}

func (ctrl *PrintGatePassController) GetAll(writer http.ResponseWriter, request *http.Request) {
	// Parsing query params from URL
	queryValues := request.URL.Query()

	// Defining query parameters as key-value pairs
	queryParams := map[string]string{
		"work_order_system_number": queryValues.Get("work_order_system_number"),
		"work_order_date":          queryValues.Get("work_order_date"),
		"customer_id":              queryValues.Get("customer_id"),
		"tnkb":                     queryValues.Get("tnkb"),
		"approval_gatepass_status": queryValues.Get("print_option"),
	}

	// Parsing pagination params using utils
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	// Build filter conditions using the defined queryParams
	criteria := utils.BuildFilterCondition(queryParams)

	// Fetch data using the service layer
	result, err := ctrl.PrintGatePassService.GetAll(criteria, paginate)
	if err != nil {
		// Handle error if the data fetch fails
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	// Respond with paginated data
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
