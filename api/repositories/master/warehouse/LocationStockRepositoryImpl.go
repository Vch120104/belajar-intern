package masterwarehouserepository

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/errors"
	"gorm.io/gorm"
	"net/http"
)

type LocationStockRepositoryImpl struct {
}

func NewLocationStockRepositoryImpl() masterrepository.LocationStockRepository {
	return &LocationStockRepositoryImpl{}
}

func (repo *LocationStockRepositoryImpl) GetAllStock(db *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	//var
	var response []masterwarehousepayloads.LocationStockDBResponse
	entities := masterentities.LocationStock{}
	Jointable := db.Table("mtr_location_stock a").Select("a.company_id," +
		"a.period_year," +
		"a.period_month," +
		"a.warehouse_id," +
		"a.location_id," +
		"a.item_id," +
		"a.warehouse_group," +
		"a.quantity_begin," +
		"a.quantity_sales," +
		"a.quantity_sales_return," +
		"a.quantity_purchase," +
		"a.quantity_purchase_return," +
		"a.quantity_transfer_in," +
		"a.quantity_transfer_out," +
		"a.quantity_claim_in," +
		"a.quantity_claim_out," +
		"a.quantity_robbing_in," +
		"a.quantity_robbing_out," +
		"a.quantity_adjustment," +
		"a.quantity_allocated," +
		"a.quantity_in_transit," +
		"a.quantity_ending," +
		"b.warehouse_costing_type," +
		"b.brand_id," +
		"(ISNULL(a.quantity_begin,0) + ISNULL(a.quantity_purchase,0)-ISNULL(a.quantity_purchase_return,0)" +
		"+ ISNULL(A.quantity_transfer_in, 0) + ISNULL(A.quantity_claim_in, 0) + ISNULL(A.quantity_robbing_in, 0) +" +
		"ISNULL(A.quantity_adjustment, 0) + ISNULL(A.quantity_sales_return, 0)" +
		"+ ISNULL(A.quantity_assembly_in, 0)) - (ISNULL(A.quantity_sales, 0) + ISNULL(A.quantity_transfer_out, 0)" +
		"+ ISNULL(A.quantity_claim_in, 0) + ISNULL(A.quantity_robbing_out, 0) + ISNULL(A.quantity_assembly_out, 0))" +
		"  AS quantity_on_hand," +
		" (ISNULL(A.quantity_begin, 0) + ISNULL(A.quantity_purchase, 0) - ISNULL(A.quantity_purchase_return, 0) +" +
		"ISNULL(A.quantity_transfer_in, 0) + ISNULL(A.quantity_robbing_in, 0) + ISNULL(A.quantity_adjustment, 0)" +
		" + ISNULL(A.quantity_sales_return, 0) + ISNULL(A.quantity_assembly_in, 0)) - (ISNULL(A.quantity_sales, 0) " +
		" + ISNULL(A.quantity_transfer_out, 0) + ISNULL(A.quantity_robbing_out, 0) + ISNULL(A.quantity_assembly_out, 0) + " +
		"ISNULL(A.quantity_allocated, 0))" +
		"AS quantity_available").Joins("left outer join mtr_warehouse_master b ON a.company_id = b.company_id AND a.warehouse_id = b.warehouse_id")
	whereQuaery := utils.ApplyFilter(Jointable, filter)
	err := whereQuaery.Scopes(pagination.Paginate(&entities, &pages, whereQuaery)).Scan(&response).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch",
			Data:       nil,
			Err:        errors.New("Failed to fetch"),
		}
	}
	pages.Rows = response
	//page := pagination.Pagination{}

	return pages, nil
}
