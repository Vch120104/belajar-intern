package transactionsparepartcontroller

import (
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"net/http"
)

// SalesOrderController
type SalesOrderController struct {
	salesOrderService transactionsparepartservice.SalesOrderService
}

// StartSalesOrderController
func NewSalesOrderController(salesOrderService transactionsparepartservice.SalesOrderService) SalesOrderController {
	return SalesOrderController{
		salesOrderService: salesOrderService,
	}
}

// GetSalesOrderByID
func (c *SalesOrderController) GetSalesOrderByID(w http.ResponseWriter, r *http.Request) {

}
