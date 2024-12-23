package transactionsparepartrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	financeserviceapiutils "after-sales/api/utils/finance-service"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ItemQueryAllCompanyRepositoryImpl struct {
}

func NewItemQueryAllCompanyRepositoryImpl() transactionsparepartrepository.ItemQueryAllCompanyRepository {
	return &ItemQueryAllCompanyRepositoryImpl{}
}

// uspg_ItemInquiryAllComp_Select
// IF @Option = 0
func (r *ItemQueryAllCompanyRepositoryImpl) GetAllItemQueryAllCompany(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []transactionsparepartpayloads.GetAllItemqueryAllCompanyResponse
	var companyId int
	var itemCodeList []string
	var movingCodeList []string
	var movingCode6 bool

	for _, data := range filterCondition {
		if data.ColumnField == "company_id" {
			tempCompanyId, errConvert := strconv.Atoi(data.ColumnValue)
			if errConvert != nil {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        errConvert,
				}
			}
			companyId = tempCompanyId
		} else if strings.HasPrefix(data.ColumnField, "item_code_") && data.ColumnValue != "" {
			itemCodeList = append(itemCodeList, data.ColumnValue)
		} else if data.ColumnField == "moving_code_6" && data.ColumnValue == "true" {
			movingCode6 = true
		} else if strings.HasPrefix(data.ColumnField, "moving_code_") && data.ColumnValue == "true" {
			movingCodeList = append(movingCodeList, strings.TrimPrefix(data.ColumnField, "moving_code_"))
		}
	}

	if movingCode6 && len(movingCodeList) != 0 && movingCodeList[len(movingCodeList)-1] != "5" {
		movingCodeList = append(movingCodeList, "5")
	}

	periodResponse, periodError := financeserviceapiutils.GetOpenPeriodByCompany(companyId, "SP")
	if periodError != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching company current period",
			Err:        periodError.Err,
		}
	}

	query := tx.
		Table("(?) as query_item", tx.
			Model(&masterentities.LocationStock{}).
			Select(
				"mtr_location_stock.company_id",
				"mtr_location_stock.item_id",
				"Item.item_name",
				"(mtr_location_stock.quantity_ending - mtr_location_stock.quantity_allocated) as quantity_on_hand",
				"ISNULL((SELECT TOP 1 moving_code_id FROM mtr_moving_code_item WHERE mtr_moving_code_item.company_id = mtr_location_stock.company_id AND mtr_moving_code_item.item_id = mtr_location_stock.item_id ORDER BY mtr_moving_code_item.process_date DESC),0) AS moving_code_id",
				"mtr_location_stock.period_year",
				"mtr_location_stock.period_month",
			).
			Joins("INNER JOIN mtr_item as Item ON Item.item_id = mtr_location_stock.item_id").
			Where("mtr_location_stock.period_year = ?", periodResponse.PeriodYear).
			Where("mtr_location_stock.period_month = ?", periodResponse.PeriodMonth).
			Where("Item.item_code IN (?)", itemCodeList).
			Where("(mtr_location_stock.quantity_ending - mtr_location_stock.quantity_allocated) > 0").
			Where("mtr_location_stock.warehouse_id NOT IN (?)", tx.
				Model(&masterwarehouseentities.WarehouseMaster{}).
				Distinct("mtr_warehouse_master.warehouse_id").
				Joins("INNER JOIN mtr_warehouse_costing_type as WarehouseCostingType ON WarehouseCostingType.warehouse_costing_type_id = mtr_warehouse_master.warehouse_costing_type_id").
				Where("WarehouseCostingType.warehouse_costing_type_code = ?", "NON"),
			),
		).
		Select(
			"query_item.company_id",
			"query_item.item_id",
			"query_item.item_name",
			"query_item.quantity_on_hand",
			"query_item.moving_code_id",
			"MovingCode.moving_code",
			"query_item.period_year",
			"query_item.period_month",
		).
		Joins("INNER JOIN mtr_moving_code as MovingCode ON MovingCode.moving_code_id = query_item.moving_code_id")

	if len(movingCodeList) != 0 {
		query = query.Where("MovingCode.moving_code IN (?)", movingCodeList)
	} else if len(movingCodeList) == 0 && movingCode6 {
		query = query.
			Where("MovingCode.moving_code IN = ?", "5").
			Where("ISNULL((?),0) > 24", tx.
				Model(&transactionsparepartentities.GoodsReceive{}).
				Select("TOP 1 DATEDIFF(MONTH,trx_goods_receive.goods_receive_document_date, ?)", time.Now()).
				Joins("INNER JOIN trx_goods_receive_detail as GoodsReceiveDetail ON GoodsReceiveDetail.goods_receive_system_number = trx_goods_receive.goods_receive_system_number").
				Joins("INNER JOIN mtr_item_goods_receive_status as GoodsReceiveStatus ON GoodsReceiveStatus.item_goods_receive_status_id = trx_goods_receive.goods_receive_status_id").
				Where("trx_goods_receive.company_id = query_item.company_id").
				Where("GoodsReceiveDetail.item_id = query_item.item_id").
				Where("GoodsReceiveStatus.item_goods_receive_status_code <> ?", "80").
				Order("trx_goods_receive.goods_receive_document_date DESC"),
			)
	}

	query = query.
		Group("query_item.company_id").
		Group("query_item.item_id").
		Group("query_item.item_name").
		Group("query_item.quantity_on_hand").
		Group("query_item.moving_code_id").
		Group("MovingCode.moving_code").
		Group("query_item.period_year").
		Group("query_item.period_month")

	err := query.Scopes(pagination.Paginate(&pages, query)).Scan(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.Rows = responses

	return pages, nil

}
