package masterrepositoryimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type LookupRepositoryImpl struct {
}

func StartLookupRepositoryImpl() masterrepository.LookupRepository {
	return &LookupRepositoryImpl{}
}

// usp_comLookUp
// IF @strEntity = 'ItemOprCode'--OPERATION MASTER & ITEM MASTER
func (r *LookupRepositoryImpl) ItemOprCode(tx *gorm.DB, linetypeId int, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var results []map[string]interface{}

	// Default filters and variables
	const (
		ItmGrpInventory      = 1 // "IN"
		PurchaseTypeGoods    = "G"
		PurchaseTypeServices = "S"
	)

	var (
		ItmCls                   string
		year, month, companyCode = 2024, 8, 1
	)

	// Ensure valid pagination limit
	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	// Initialize base query
	baseQuery := tx.Table("")

	// Build filter string dynamically from provided filters
	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	// Build the query based on linetypeId
	switch linetypeId {
	case utils.LinetypePackage:
		combinedDetailsSubQuery := `
				(
					SELECT package_id, frt_quantity, is_active 
					FROM mtr_package_master_detail_item
					WHERE is_active = 1
					UNION ALL
					SELECT package_id, frt_quantity, is_active 
					FROM mtr_package_master_detail_operation
					WHERE is_active = 1
				) AS CombinedDetails
			`

		baseQuery = baseQuery.Table("mtr_package A").
			Select(`
				A.package_code AS package_code, 
				A.package_name AS package_name, 
				SUM(CombinedDetails.frt_quantity) AS frt, 
				B.profit_center_id AS profit_center, 
				C.model_code AS model_code, 
				C.model_description AS description, 
				A.package_price AS price
			`).
			Joins("LEFT JOIN "+combinedDetailsSubQuery+" ON A.package_id = CombinedDetails.package_id").
			Joins("LEFT JOIN dms_microservices_general_dev.dbo.mtr_profit_center B ON A.profit_center_id = B.profit_center_id").
			Joins("LEFT JOIN dms_microservices_sales_dev.dbo.mtr_unit_model C ON A.model_id = C.model_id").
			Where("A.is_active = ?", 1).
			Where(filterQuery, filterValues...).
			Group("A.package_code, A.package_name, B.profit_center_id, C.model_code, C.model_description, A.package_price")

	case utils.LinetypeOperation:
		baseQuery = baseQuery.Table("dms_microservices_aftersales_dev.dbo.mtr_operation_code AS oc").
			Select(`
        oc.operation_code AS OPERATION_CODE, 
        oc.operation_name AS OPERATION_NAME, 
        ofrt.frt_hour AS FRT_HOUR, 
        oe.operation_entries_code AS OPERATION_ENTRIES_CODE, 
        oe.operation_entries_description AS OPERATION_ENTRIES_DESCRIPTION, 
        ok.operation_key_code AS OPERATION_KEY_CODE, 
        ok.operation_key_description AS OPERATION_KEY_DESCRIPTION
    `).
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_entries AS oe ON oc.operation_entries_id = oe.operation_entries_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_key AS ok ON oc.operation_key_id = ok.operation_key_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_model_mapping AS omm ON oc.operation_id = omm.operation_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_frt AS ofrt ON omm.operation_model_mapping_id = ofrt.operation_model_mapping_id").
			Where("oc.is_active = ? ", 1).
			Where(filterQuery, filterValues...)

	case utils.LinetypeSparepart:
		ItmCls = "1"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS AvailQty, 
				A.item_level_1 AS item_level_1, 
				A.item_level_2 AS item_level_2, 
				A.item_level_3 AS item_level_3, 
				A.item_level_4 AS item_level_4
			`, year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type = ? AND A.item_class_id = ? AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, 1).
			Where(filterQuery, filterValues...)

	case utils.LinetypeOil:
		ItmCls = "2"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS AvailQty, 
				A.item_level_1 AS item_level_1, 
				A.item_level_2 AS item_level_2, 
				A.item_level_3 AS item_level_3, 
				A.item_level_4 AS item_level_4
			`, year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type = ? AND A.item_class_id = ? AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, 1).
			Where(filterQuery, filterValues...)

	case utils.LinetypeMaterial:
		ItmCls = "3"
		ItmClsSublet := "2"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				DISTINCT A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS AvailQty, 
				A.item_level_1 AS item_level_1, 
				A.item_level_2 AS item_level_2, 
				A.item_level_3 AS item_level_3, 
				A.item_level_4 AS item_level_4
			`, year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type = ? AND (A.item_class_id = ? OR A.item_class_id = ?) AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, ItmClsSublet, 1).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	case utils.LinetypeSublet:
		ItmCls = "4"
		ItmGrpOutsideJob := 4

		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				DISTINCT A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS AvailQty, 
				A.item_level_1 AS item_level_1, 
				A.item_level_2 AS item_level_2, 
				A.item_level_3 AS item_level_3, 
				A.item_level_4 AS item_level_4
			`, year, month, companyCode).
			Where("(A.item_group_id = ? OR A.item_group_id = ?) AND A.item_class_id = ? AND A.item_type = ? AND A.is_active = ?", ItmGrpOutsideJob, ItmGrpInventory, ItmCls, PurchaseTypeServices, 1).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	case utils.LinetypeAccesories:
		ItmCls = "5"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				DISTINCT A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS AvailQty, 
				A.item_level_1 AS item_level_1, 
				A.item_level_2 AS item_level_2, 
				A.item_level_3 AS item_level_3, 
				A.item_level_4 AS item_level_4
			`, year, month, companyCode).
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", ItmCls, ItmGrpInventory, 1).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	default:
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid linetype ID",
			Err:        errors.New("invalid linetype ID"),
		}
	}

	// Count total rows for pagination
	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count rows",
			Err:        err,
		}
	}

	// Apply pagination
	offset := (paginate.Page - 1) * paginate.Limit
	if err := baseQuery.Offset(offset).Limit(paginate.Limit).Find(&results).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve data",
			Err:        err,
		}
	}

	// Get total number of pages
	totalPages := int(math.Ceil(float64(totalRows) / float64(paginate.Limit)))

	fmt.Println("Final Results:", results)
	fmt.Println("Total Rows:", totalRows)
	fmt.Println("Total Pages:", totalPages)

	return results, int(totalRows), totalPages, nil
}

// usp_comLookUp
// IF @strEntity = 'ItemOprCodeWithPrice'--OPERATION MASTER & ITEM MASTER WITH PRICELIST
func (r *LookupRepositoryImpl) ItemOprCodeWithPrice(tx *gorm.DB, linetypeId int, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var results []map[string]interface{}

	// Default filters and variables
	const (
		ItmGrpInventory   = 1 // "IN"
		PurchaseTypeGoods = "G"
		ItemService       = "S"
		BillCodeNoCharge  = "N"
		BillCodeC         = "C"
		BillCodeInt       = "I"
	)

	type Period struct {
		PeriodYear  string `gorm:"column:PERIOD_YEAR"`
		PeriodMonth string `gorm:"column:PERIOD_MONTH"`
	}

	// Mengambil periode tahun dan bulan
	var (
		ItmCls      string
		year, month string
		period      Period
		companyCode = 1
	)

	result := tx.Table("dms_microservices_finance_dev.dbo.mtr_closing_period_company").
		Select("TOP 1 PERIOD_YEAR, PERIOD_MONTH").
		Where("COMPANY_CODE = ? AND MODULE_CODE = 'SP' AND PERIOD_YEAR <= ? AND PERIOD_MONTH <= ? AND PERIOD_STATUS = 'O'", companyCode, 2024, 8).
		Order("PERIOD_YEAR DESC, PERIOD_MONTH DESC").
		Scan(&period)

	if result.Error != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get period",
			Err:        result.Error,
		}
	}

	year = period.PeriodYear
	month = period.PeriodMonth

	fmt.Println("Period Year:", year)
	fmt.Println("Period Month:", month)

	// Ensure valid pagination limit
	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	// Initialize base query
	baseQuery := tx.Table("")

	// Build filter string dynamically from provided filters
	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	// Build the query based on linetypeId
	switch linetypeId {
	case utils.LinetypePackage:
		combinedDetailsSubQuery := `
				(
					SELECT package_id, frt_quantity, is_active 
					FROM mtr_package_master_detail_item
					WHERE is_active = 1
					UNION ALL
					SELECT package_id, frt_quantity, is_active 
					FROM mtr_package_master_detail_operation
					WHERE is_active = 1
				) AS CombinedDetails
			`

		baseQuery = baseQuery.Table("mtr_package A").
			Select(`
				A.package_code AS package_code, 
				A.package_name AS package_name, 
				SUM(CombinedDetails.frt_quantity) AS frt, 
				B.profit_center_id AS profit_center, 
				C.model_code AS model_code, 
				C.model_description AS description, 
				A.package_price AS price
			`).
			Joins("LEFT JOIN "+combinedDetailsSubQuery+" ON A.package_id = CombinedDetails.package_id").
			Joins("LEFT JOIN dms_microservices_general_dev.dbo.mtr_profit_center B ON A.profit_center_id = B.profit_center_id").
			Joins("LEFT JOIN dms_microservices_sales_dev.dbo.mtr_unit_model C ON A.model_id = C.model_id").
			Where("A.is_active = ?", 1).
			Where(filterQuery, filterValues...).
			Group("A.package_code, A.package_name, B.profit_center_id, C.model_code, C.model_description, A.package_price")

	case utils.LinetypeOperation:
		baseQuery = baseQuery.Table("dms_microservices_aftersales_dev.dbo.mtr_operation_code AS oc").
			Select(`
        oc.operation_code AS OPERATION_CODE, 
        oc.operation_name AS OPERATION_NAME, 
        ofrt.frt_hour AS FRT_HOUR, 
        oe.operation_entries_code AS OPERATION_ENTRIES_CODE, 
        oe.operation_entries_description AS OPERATION_ENTRIES_DESCRIPTION, 
        ok.operation_key_code AS OPERATION_KEY_CODE, 
        ok.operation_key_description AS OPERATION_KEY_DESCRIPTION
    `).
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_entries AS oe ON oc.operation_entries_id = oe.operation_entries_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_key AS ok ON oc.operation_key_id = ok.operation_key_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_model_mapping AS omm ON oc.operation_id = omm.operation_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_frt AS ofrt ON omm.operation_model_mapping_id = ofrt.operation_model_mapping_id").
			Where("oc.is_active = ? ", 1).
			Where(filterQuery, filterValues...)

	case utils.LinetypeSparepart:
		ItmCls = "1"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS AvailQty, 
				A.item_level_1 AS item_level_1, 
				A.item_level_2 AS item_level_2, 
				A.item_level_3 AS item_level_3, 
				A.item_level_4 AS item_level_4
			`, year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type = ? AND A.item_class_id = ? AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, 1).
			Where(filterQuery, filterValues...)

	case utils.LinetypeOil:
		ItmCls = "2"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS AvailQty, 
				A.item_level_1 AS item_level_1, 
				A.item_level_2 AS item_level_2, 
				A.item_level_3 AS item_level_3, 
				A.item_level_4 AS item_level_4
			`, year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type = ? AND A.item_class_id = ? AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, 1).
			Where(filterQuery, filterValues...)

	case utils.LinetypeMaterial:
		ItmCls = "3"
		ItmClsSublet := "2"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				DISTINCT A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS AvailQty, 
				A.item_level_1 AS item_level_1, 
				A.item_level_2 AS item_level_2, 
				A.item_level_3 AS item_level_3, 
				A.item_level_4 AS item_level_4
			`, year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type = ? AND (A.item_class_id = ? OR A.item_class_id = ?) AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, ItmClsSublet, 1).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	case utils.LinetypeSublet:
		ItmCls = "4"
		ItmGrpOutsideJob := 4
		PurchaseTypeServices := "S"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				DISTINCT A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS AvailQty, 
				A.item_level_1 AS item_level_1, 
				A.item_level_2 AS item_level_2, 
				A.item_level_3 AS item_level_3, 
				A.item_level_4 AS item_level_4
			`, year, month, companyCode).
			Where("(A.item_group_id = ? OR A.item_group_id = ?) AND A.item_class_id = ? AND A.item_type = ? AND A.is_active = ?", ItmGrpOutsideJob, ItmGrpInventory, ItmCls, PurchaseTypeServices, 1).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	case utils.LinetypeAccesories:
		ItmCls = "5"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				DISTINCT A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS AvailQty, 
				A.item_level_1 AS item_level_1, 
				A.item_level_2 AS item_level_2, 
				A.item_level_3 AS item_level_3, 
				A.item_level_4 AS item_level_4
			`, year, month, companyCode).
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", ItmCls, ItmGrpInventory, 1).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	default:
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid linetype ID",
			Err:        errors.New("invalid linetype ID"),
		}
	}

	// Count total rows for pagination
	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count rows",
			Err:        err,
		}
	}

	// Apply pagination
	offset := (paginate.Page - 1) * paginate.Limit
	if err := baseQuery.Offset(offset).Limit(paginate.Limit).Find(&results).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve data",
			Err:        err,
		}
	}

	// Get total number of pages
	totalPages := int(math.Ceil(float64(totalRows) / float64(paginate.Limit)))

	fmt.Println("Final Results:", results)
	fmt.Println("Total Rows:", totalRows)
	fmt.Println("Total Pages:", totalPages)

	return results, int(totalRows), totalPages, nil
}

// usp_comLookUp
// IF @strEntity = 'Vehicle0'--VEHICLE UNIT MASTER
func (r *LookupRepositoryImpl) VehicleUnitMaster(tx *gorm.DB, brandId int, modelId int, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		vehicleMasters []map[string]interface{}
		totalRows      int64
		totalPages     int
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	// Build filter string dynamically from provided filters
	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	// Build the initial GORM query with joins and select
	query := tx.Table("dms_microservices_sales_dev.dbo.mtr_vehicle V").
		Select(`
			V.vehicle_id AS vehicle_id,
			V.vehicle_chassis_number AS vehicle_chassis_number, 
			RC.vehicle_registration_certificate_tnkb AS vehicle_registration_certificate_tnkb, 
			RC.vehicle_registration_certificate_owner_name AS vehicle_registration_certificate_owner_name, 
			UM.model_variant_colour_name AS Vehicle, 
			CAST(V.vehicle_production_year AS VARCHAR) AS vehicle_production_year, 
			CONVERT(VARCHAR, V.vehicle_last_service_date, 106) AS vehicle_last_service_date, 
			V.vehicle_last_km AS vehicle_last_km, 
			CASE 
				WHEN V.is_active = 1 THEN 'Active' 
				WHEN V.is_active = 0 THEN 'Deactive' 
			END AS Status
		`).
		Joins(`LEFT JOIN dms_microservices_sales_dev.dbo.mtr_vehicle_registration_certificate RC ON V.vehicle_id = RC.vehicle_id`).
		Joins(`LEFT JOIN dms_microservices_sales_dev.dbo.mtr_model_variant_colour UM ON UM.brand_id = V.vehicle_brand_id AND 
                                       UM.model_id = V.vehicle_model_id AND 
                                       UM.colour_id = V.vehicle_colour_id AND 
                                       ISNULL(UM.accessories_option_id, '') = ISNULL(V.option_id, '')`).
		Where(filterQuery, filterValues...).
		Where("V.vehicle_brand_id = ?", brandId).
		Where("V.vehicle_model_id = ?", modelId)

	// Count total rows for pagination
	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total vehicle units",
			Err:        err,
		}
	}

	// Calculate total pages based on totalRows and paginate.Limit
	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	// Apply pagination and execute query
	err = query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&vehicleMasters).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get vehicle unit master data",
			Err:        err,
		}
	}

	// Return paginated data, total pages, and total rows
	return vehicleMasters, totalPages, int(totalRows), nil
}

// usp_comLookUp
// IF @strEntity = 'CampaignMaster'--CAMPAIGN MASTER
func (r *LookupRepositoryImpl) CampaignMaster(tx *gorm.DB, companyId int, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	var (
		campaignMasters []map[string]interface{}
		totalRows       int64
		totalPages      int
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	// Build filter string dynamically from provided filters
	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	// Build the initial GORM query with joins and select
	query := tx.Table("dms_microservices_aftersales_dev.dbo.mtr_campaign C").
		Select(`
			C.campaign_id AS campaign_id,
			C.campaign_code AS campaign_code,
			C.campaign_name AS campaign_name,
			C.model_id AS model_id,
			C.campaign_period_from AS campaign_period_from,
			C.campaign_period_to AS campaign_period_to,
			C.total_after_vat AS total_after_vat,
			CASE 
				WHEN C.is_active = 1 THEN 'Active' 
				WHEN C.is_active = 0 THEN 'Deactive' 
			END AS Status
			`).
		//Joins(`LEFT JOIN dms_microservices_sales_dev.dbo.mtr_model_variant_colour VC ON C.model_id = VC.model_id`).
		Where(filterQuery, filterValues...).
		Where("C.company_id = ?", companyId)

		// Count total rows for pagination
	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total vehicle units",
			Err:        err,
		}
	}

	// Calculate total pages based on totalRows and paginate.Limit
	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	// Apply pagination and execute query
	err = query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&campaignMasters).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get vehicle unit master data",
			Err:        err,
		}
	}

	// Return paginated data, total pages, and total rows
	return campaignMasters, totalPages, int(totalRows), nil
}
