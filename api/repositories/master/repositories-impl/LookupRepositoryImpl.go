package masterrepositoryimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"fmt"
	"math"
	"net/http"

	"gorm.io/gorm"
)

type LookupRepositoryImpl struct {
}

func StartLookupRepositoryImpl() masterrepository.LookupRepository {
	return &LookupRepositoryImpl{}
}

func (r *LookupRepositoryImpl) ItemOprCode(tx *gorm.DB, linetypeId int, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var results []map[string]interface{}
	var err error

	// Default filters and variables
	ItmGrpInventory := 1 // "IN"
	PurchaseTypeGoods := "G"
	ItmCls := ""
	year, month, companyCode := 2024, 8, 1 // Placeholder for real values or dynamic fetching

	// pagination limit is set
	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	baseQuery := tx.Table("")

	// Switch based on the linetypeId
	switch linetypeId {
	case utils.LinetypePackage:
		baseQuery = baseQuery.Table("mtr_package A").
			Select("A.package_code AS package_code, A.package_name AS package_name, SUM(A1.frt_quantity) AS frt, B.profit_center_id AS profit_center, C.model_code AS model_code, C.model_description AS description, A.package_price AS price").
			Joins("LEFT JOIN mtr_package_master_detail_item A1 ON A.package_id = A1.package_id").
			Joins("LEFT JOIN dms_microservices_general_dev.dbo.mtr_profit_center B ON A.profit_center_id = B.profit_center_id").
			Joins("LEFT JOIN dms_microservices_sales_dev.dbo.mtr_unit_model C ON A.model_id = C.model_id").
			Where("A.is_active = ? AND A1.is_active = ?", 1, 1).
			Group("A.package_code, A.package_name, B.profit_center_id, C.model_code, C.model_description, A.package_price")

	case utils.LinetypeOperation:
		baseQuery = baseQuery.Table("mtr_operation_frt A").
			Select("A.OPERATION_CODE AS Code, B.OPERATION_NAME AS Description, A.FRT_HOUR AS FRT, C.OPR_ENTRIES_CODE AS OprEntriesCode, G.OPR_ENTRIES_DESC AS OprEntriesName, C.OPR_KEY_CODE AS OprKeyCode, F.OPR_KEY_DESC AS OprKeyName").
			Joins("INNER JOIN mtr_operation_model_mapping O ON O.OPERATION_CODE = A.OPERATION_CODE AND O.VEHICLE_BRAND = A.VEHICLE_BRAND AND O.MODEL_CODE = A.MODEL_CODE").
			Joins("LEFT JOIN mtr_operation_level C ON A.OPERATION_CODE = C.OPERATION_CODE AND A.VEHICLE_BRAND = C.VEHICLE_BRAND AND A.MODEL_CODE = C.MODEL_CODE").
			Joins("LEFT JOIN mtr_operation_code B ON A.OPERATION_CODE = B.OPERATION_CODE").
			Joins("LEFT JOIN mtr_operation_key F ON C.OPR_KEY_CODE = F.OPR_KEY_CODE").
			Joins("LEFT JOIN mtr_operation_entries G ON C.OPR_ENTRIES_CODE = G.OPR_ENTRIES_CODE").
			Where("A.is_active = ? AND O.is_active = ?", 1, 1)

	case utils.LinetypeSparepart:
		ItmCls = "1" // "SP"
		baseQuery = baseQuery.Table("mtr_item A").
			Select("A.item_code AS item_code, A.item_name AS item_name, ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V WHERE A.item_id = V.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS AvailQty, A.item_level_1 AS item_level_1, A.item_level_2 AS item_level_2, A.item_level_3 AS item_level_3, A.item_level_4 AS item_level_4", year, month, companyCode).
			Joins("INNER JOIN mtr_item_detail B ON A.item_id = B.item_id").
			Where("A.item_group_id = ? AND A.item_type = ? AND A.item_class_id = ? AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, 1)

	case utils.LinetypeOil:
		ItmCls = "1" // "OL"
		baseQuery = baseQuery.Table("mtr_item A").
			Select("A.item_code AS item_code, A.item_name AS item_name, ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V WHERE A.item_id = V.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS AvailQty, A.item_level_1 AS item_level_1, A.item_level_2 AS item_level_2, A.item_level_3 AS item_level_3, A.item_level_4 AS item_level_4", year, month, companyCode).
			Joins("INNER JOIN mtr_item_detail B ON A.item_id = B.item_id").
			Where("A.item_group_id = ? AND A.item_type = ? AND A.item_class_id = ? AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, 1)

	case utils.LinetypeMaterial:
		ItmCls = "1"        // "MT"
		ItmClsSublet := "2" //"Sublet_Class"
		baseQuery = baseQuery.Table("mtr_item A").
			Select("DISTINCT A.item_code AS item_code, A.item_name AS item_name, ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V WHERE A.item_id = V.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS AvailQty, A.item_level_1 AS item_level_1, A.item_level_2 AS item_level_2, A.item_level_3 AS item_level_3, A.item_level_4 AS item_level_4", year, month, companyCode).
			Joins("INNER JOIN mtr_item_detail B ON A.item_id = B.item_id").
			Where("A.item_group_id = ? AND A.item_type = ? AND (A.item_class_id = ? OR A.item_class_id = ?) AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, ItmClsSublet, 1).
			Order("A.item_code")

	case utils.LinetypeSublet:
		ItmCls = "1"                //"Fee_Class"
		ItmGrpOutsideJob := "1"     //"OutsideJob_Group"
		PurchaseTypeServices := "1" //"Service_Type"
		baseQuery = baseQuery.Table("mtr_item A").
			Select("DISTINCT A.item_code AS item_code, A.item_name AS item_name, ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V WHERE A.item_id = V.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS AvailQty, A.item_level_1 AS item_level_1, A.item_level_2 AS item_level_2, A.item_level_3 AS item_level_3, A.item_level_4 AS item_level_4", year, month, companyCode).
			Joins("INNER JOIN mtr_item_detail B ON A.item_id = B.item_id").
			Where("(A.item_group_id = ? OR (A.item_group_id = ? AND A.item_class_id = ?)) AND A.item_type = ? AND A.is_active = ?", ItmGrpOutsideJob, ItmGrpInventory, ItmCls, PurchaseTypeServices, 1).
			Order("A.item_code")

	default:
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Line Type ID",
			Err:        fmt.Errorf("invalid line type ID: %d", linetypeId),
		}
	}

	for _, filter := range filters {
		baseQuery = baseQuery.Where(fmt.Sprintf("%s = ?", filter.ColumnField), filter.ColumnValue)
	}

	err = baseQuery.Offset((paginate.Page - 1) * paginate.Limit).Limit(paginate.Limit).Find(&results).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	var totalRows int64
	countQuery := baseQuery.Session(&gorm.Session{}).Model(&results)
	countErr := countQuery.Count(&totalRows).Error
	if countErr != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        countErr,
		}
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(paginate.Limit)))

	return results, int(totalRows), totalPages, nil
}
