package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	masteroperationentities "after-sales/api/entities/master/operation"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

type LookupRepositoryImpl struct {
}

func StartLookupRepositoryImpl() masterrepository.LookupRepository {
	return &LookupRepositoryImpl{}
}

func (r *LookupRepositoryImpl) GetOprItemPrice(tx *gorm.DB, linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, jobTypeId int, variantId int, currencyId int, billCode string, whsGroup string) (float64, *exceptions.BaseErrorResponse) {
	var (
		price               float64
		effDate             = time.Now()
		markupPercentage    float64
		companyCodePrice    int
		commonPriceList     bool
		defaultPriceCode    = "A"
		useDiscDecentralize string
		itemService         string
		priceCount          int64
		priceCode           string
	)

	// Set markup percentage based on company ID
	markupPercentage = 0
	if companyId == 139 {
		markupPercentage = 11.00
	}

	if err := tx.Model(&masteritementities.Item{}).
		Where("item_id = ?", oprItemCode).
		Select("common_pricelist").
		Scan(&commonPriceList).Error; err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get common price list",
			Err:        err,
		}
	}

	// Set company code price based on common price list
	if commonPriceList {
		companyCodePrice = 0
	} else {
		companyCodePrice = companyId
	}

	switch linetypeId {
	case utils.LinetypePackage:
		// Package price logic
		if err := tx.Model(&masterentities.PackageMaster{}).
			Where("package_code = ?", oprItemCode).
			Select("package_price").
			Scan(&price).Error; err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get package price",
				Err:        err,
			}
		}

	case utils.LinetypeOperation:
		// Operation price logic
		query := tx.Model(&masteroperationentities.LabourSellingPriceDetail{}).
			Joins("JOIN mtr_labour_selling_price ON mtr_labour_selling_price.labour_selling_price_id = mtr_labour_selling_price_detail.labour_selling_price_id").
			Where("mtr_labour_selling_price.brand_id = ? AND mtr_labour_selling_price.effective_date <= ? AND mtr_labour_selling_price.job_type_id = ? AND mtr_labour_selling_price_detail.model_id = ? AND mtr_labour_selling_price.company_id = ?",
				brandId, effDate, jobTypeId, modelId, companyId)

		if variantId == 0 {
			query = query.Where("mtr_labour_selling_price_detail.variant_id = 0")
		} else {
			query = query.Where("mtr_labour_selling_price_detail.variant_id = ? OR mtr_labour_selling_price_detail.variant_id = 0", variantId)
		}

		if err := query.Order("mtr_labour_selling_price.effective_date DESC").Limit(1).Pluck("selling_price", &price).Error; err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get selling price",
				Err:        err,
			}
		}

	default:

		if err := tx.Model(&masteritementities.PriceList{}).
			Where("is_active = 1 AND brand_id = ? AND effective_date <= ? AND item_code = ? AND currency_id = ? AND company_id = ? AND price_list_code_id = ?",
				brandId, effDate, oprItemCode, currencyId, companyCodePrice, priceCode).Count(&priceCount).Error; err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to check price list existence",
				Err:        err,
			}
		}

		if priceCount == 0 {
			priceCode = defaultPriceCode
		}

		// Handling based on bill code
		if billCode == "N" || billCode == "C" || billCode == "I" || billCode == "SU06" || billCode == "SU05" || billCode == "SU08" {
			var periodYear, periodMonth string

			month := effDate.Format("01")
			year := effDate.Format("2006")

			// Get MODULE_SP and PERIOD_STATUS_OPEN
			moduleSP := "SP"
			periodStatusOpen := "O"

			var result struct {
				PeriodYear  *string `gorm:"column:period_year"`
				PeriodMonth *string `gorm:"column:period_month"`
			}

			if err := tx.Table("dms_microservices_finance_dev.dbo.mtr_closing_period_company").
				Where("company_id = ? AND module_code = ? AND period_year <= ? AND period_month <= ? AND period_status = ?",
					companyId, moduleSP, year, month, periodStatusOpen).
				Order("period_year DESC, period_month DESC").
				Limit(1).
				Select("period_year, period_month").
				Scan(&result).Error; err != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to get period details",
					Err:        err,
				}
			}

			if result.PeriodYear != nil {
				periodYear = *result.PeriodYear
			} else {
				periodYear = "0000"
			}

			if result.PeriodMonth != nil {
				periodMonth = *result.PeriodMonth
			} else {
				periodMonth = "00"
			}

			// Check item type
			itemTypeExists := false
			if err := tx.Model(&masteritementities.Item{}).
				Where("item_code = ? AND item_type = ?", oprItemCode, itemService).
				Select("item_type").
				Scan(&itemTypeExists).Error; err != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to check item type",
					Err:        err,
				}
			}

			if itemTypeExists {
				// Get price from gmPriceList for items
				if err := tx.Model(&masteritementities.PriceList{}).
					Where("is_active = 1 AND brand_id = ? AND effective_date <= ? AND item_code = ? AND currency_id = ? AND company_id = ? AND price_list_code_id = ?",
						brandId, effDate, oprItemCode, currencyId, companyCodePrice, priceCode).
					Order("effective_date DESC").
					Limit(1).
					Pluck("price_list_amount", &price).Error; err != nil {
					return 0, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to get price amount from gmPriceList",
						Err:        err,
					}
				}
			} else {
				// Get price from amGroupStock for other items
				if err := tx.Model(&masterentities.GroupStock{}).
					Where("period_year = ? AND period_month = ? AND item_code = ? AND company_code = ? AND whs_group = ?",
						periodYear, periodMonth, oprItemCode, companyId, whsGroup).
					Select("CASE ISNULL(price_current, 0) WHEN 0 THEN price_begin ELSE price_current END AS hpp").
					Pluck("hpp", &price).Error; err != nil {
					return 0, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to get group stock price",
						Err:        err,
					}
				}
			}

			if billCode != "I" && billCode != "SU08" && billCode != "SU05" {
				if err := tx.Model(&masteritementities.Item{}).
					Where("item_code = ?", oprItemCode).
					Pluck("use_disc_decentralize", &useDiscDecentralize).Error; err != nil {
					return 0, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to get useDiscDecentralize value",
						Err:        err,
					}
				}

				if useDiscDecentralize == "" {
					useDiscDecentralize = "Y"
				}

				if useDiscDecentralize == "N" {
					if err := tx.Model(&masteritementities.PriceList{}).
						Where("is_active = 1 AND brand_id = ? AND effective_date <= ? AND item_code = ? AND currency_id = ? AND company_id = ? AND price_code = ?",
							brandId, effDate, oprItemCode, currencyId, companyCodePrice, priceCode).
						Order("effective_date DESC").
						Limit(1).
						Pluck("price_list_amount", &price).Error; err != nil {
						return 0, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to get price amount",
							Err:        err,
						}
					}
				}
			}

		} else {
			if err := tx.Model(&masteritementities.PriceList{}).
				Where("is_active = 1 AND brand_id = ? AND effective_date <= ? AND item_code = ? AND currency_id = ? AND company_id = ? AND price_list_code_id = ?",
					brandId, effDate, oprItemCode, currencyId, companyCodePrice, priceCode).
				Order("effective_date DESC").
				Limit(1).
				Pluck("price_list_amount", &price).Error; err != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to get price amount",
					Err:        err,
				}
			}
		}
	}

	if linetypeId == utils.LinetypeOperation && billCode == "I" {
		price += price * markupPercentage / 100
	}

	return price, nil
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

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	baseQuery := tx.Table("")

	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

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

	case utils.LineTypeFee:
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

	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count rows",
			Err:        err,
		}
	}

	offset := (paginate.Page - 1) * paginate.Limit
	if err := baseQuery.Offset(offset).Limit(paginate.Limit).Find(&results).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve data",
			Err:        err,
		}
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(paginate.Limit)))

	return results, int(totalRows), totalPages, nil
}

func (r *LookupRepositoryImpl) ItemOprCodeByCode(tx *gorm.DB, linetypeId int, oprItemCode string, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
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

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	baseQuery := tx.Table("")

	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

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
			Where("A.package_code = ?", oprItemCode).
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
			Where("oc.operation_code = ?", oprItemCode).
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
			Where("A.item_code = ?", oprItemCode).
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
			Where("A.item_code = ?", oprItemCode).
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
			Where("A.item_code = ?", oprItemCode).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	case utils.LineTypeFee:
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
			Where("A.item_code = ?", oprItemCode).
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
			Where("A.item_code = ?", oprItemCode).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	default:
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid linetype ID",
			Err:        errors.New("invalid linetype ID"),
		}
	}

	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count rows",
			Err:        err,
		}
	}

	offset := (paginate.Page - 1) * paginate.Limit
	if err := baseQuery.Offset(offset).Limit(paginate.Limit).Find(&results).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve data",
			Err:        err,
		}
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(paginate.Limit)))

	return results, int(totalRows), totalPages, nil
}

func (r *LookupRepositoryImpl) ItemOprCodeByID(tx *gorm.DB, linetypeId int, oprItemId int, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
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

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	baseQuery := tx.Table("")

	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

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
			Where("A.package_id = ?", oprItemId).
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
			Where("oc.operation_id = ?", oprItemId).
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
			Where("A.item_id = ?", oprItemId).
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
			Where("A.item_id = ?", oprItemId).
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
			Where("A.item_id = ?", oprItemId).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	case utils.LineTypeFee:
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
			Where("A.item_id = ?", oprItemId).
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
			Where("A.item_id = ?", oprItemId).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	default:
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid linetype ID",
			Err:        errors.New("invalid linetype ID"),
		}
	}

	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count rows",
			Err:        err,
		}
	}

	offset := (paginate.Page - 1) * paginate.Limit
	if err := baseQuery.Offset(offset).Limit(paginate.Limit).Find(&results).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve data",
			Err:        err,
		}
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(paginate.Limit)))

	return results, int(totalRows), totalPages, nil
}

// usp_comLookUp
// IF @strEntity = 'ItemOprCodeWithPrice'--OPERATION MASTER & ITEM MASTER WITH PRICELIST
func (r *LookupRepositoryImpl) ItemOprCodeWithPrice(tx *gorm.DB, linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, jobTypeId int, variantId int, currencyId int, billCode string, whsGroup string, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var results []map[string]interface{}

	const (
		ItmGrpInventory   = 1 // "IN"
		PurchaseTypeGoods = "G"
		ItemService       = "S"
		BillCodeNoCharge  = "N"
		BillCodeC         = "C"
		BillCodeInt       = "I"
		defaultPriceCode  = "A"
	)

	type Period struct {
		PeriodYear  string `gorm:"column:PERIOD_YEAR"`
		PeriodMonth string `gorm:"column:PERIOD_MONTH"`
	}

	var (
		ItmCls       string
		year, month  string
		period       Period
		companyCode  = 151
		closingModul = 10
		yearNow      = time.Now().Format("2006")
		monthNow     = time.Now().Format("01")
	)

	result := tx.Table("dms_microservices_finance_dev.dbo.mtr_closing_period_company").
		Select("TOP 1 PERIOD_YEAR, PERIOD_MONTH").
		Where("company_id = ? AND closing_module_detail_id = ? AND PERIOD_YEAR <= ? AND PERIOD_MONTH <= ? AND is_period_closed = '0'", companyCode, closingModul, yearNow, monthNow).
		Order("PERIOD_YEAR DESC, PERIOD_MONTH DESC").
		Scan(&period)

	if result.Error != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get period , please check closing period company",
			Err:        result.Error,
		}
	}

	year = period.PeriodYear
	month = period.PeriodMonth

	fmt.Println("Period Year:", year)
	fmt.Println("Period Month:", month)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	baseQuery := tx.Table("")

	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	price, err := r.GetOprItemPrice(tx, linetypeId, companyId, oprItemCode, brandId, modelId, jobTypeId, variantId, currencyId, billCode, whsGroup)
	if err != nil {
		return nil, 0, 0, err
	}

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
        ok.operation_key_description AS OPERATION_KEY_DESCRIPTION,
		? as PRICE
    `, price).
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
						AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ? 
						AND V.whs_group = ?), 0) AS AvailQty,
				A.item_level_1 AS item_level_1,
				A.item_level_2 AS item_level_2,
				A.item_level_3 AS item_level_3,
				A.item_level_4 AS item_level_4,
				CASE 
					WHEN ? IN (?, ?, ?) THEN
						CASE A.item_type
							WHEN ? THEN
								(SELECT TOP 1 price_list_amount FROM mtr_price_list
								WHERE is_active = 1 
								AND brand_id = B.brand_id 
								AND effective_date <= GETDATE()
								AND item_id = A.item_id
								AND currency_id = ? 
								AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
								AND price_list_code_id = ?
								ORDER BY effective_date DESC)
							ELSE
								(SELECT CASE ISNULL(price_current, 0)
										WHEN 0 THEN price_begin 
										ELSE price_current END AS HPP
								FROM mtr_group_stock 
								WHERE period_year = ? 
								AND period_month = ? 
								AND item_code = A.item_code 
								AND company_id = ?  
								AND whs_group = ?)
						END
					ELSE
						(SELECT TOP 1 price_list_amount FROM mtr_price_list
						WHERE is_active = 1 
						AND brand_id = B.brand_id 
						AND effective_date <= GETDATE()
						AND item_id = A.item_id 
						AND currency_id = ? 
						AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
						AND price_list_code_id = ?
						ORDER BY effective_date DESC)
				END AS PRICE
			`, year, month, companyId, whsGroup, // Parameters for AvailQty subquery
				billCode, BillCodeC, BillCodeInt, BillCodeNoCharge, // Parameters for CASE statement
				ItemService, currencyId, companyId, defaultPriceCode, // Parameters for subquery in CASE
				year, month, companyId, whsGroup, // Parameters for ELSE subquery in CASE
				currencyId, companyId, defaultPriceCode). // Parameters for the final ELSE condition.
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
				A.item_level_4 AS item_level_4,
				CASE 
					WHEN ? IN (?, ?, ?) THEN
						CASE A.item_type
							WHEN ? THEN
								(SELECT TOP 1 price_list_amount FROM mtr_price_list
								WHERE is_active = 1 
								AND brand_id = B.brand_id 
								AND effective_date <= GETDATE()
								AND item_id = A.item_id
								AND currency_id = ? 
								AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
								AND price_list_code_id = ?
								ORDER BY effective_date DESC)
							ELSE
								(SELECT CASE ISNULL(price_current, 0)
										WHEN 0 THEN price_begin 
										ELSE price_current END AS HPP
								FROM mtr_group_stock 
								WHERE period_year = ? 
								AND period_month = ? 
								AND item_code = A.item_code 
								AND company_id = ?  
								AND whs_group = ?)
						END
					ELSE
						(SELECT TOP 1 price_list_amount FROM mtr_price_list
						WHERE is_active = 1 
						AND brand_id = B.brand_id 
						AND effective_date <= GETDATE()
						AND item_id = A.item_id 
						AND currency_id = ? 
						AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
						AND price_list_code_id = ?
						ORDER BY effective_date DESC)
				END AS PRICE
			`, year, month, companyId, whsGroup, // Parameters for AvailQty subquery
				billCode, BillCodeC, BillCodeInt, BillCodeNoCharge, // Parameters for CASE statement
				ItemService, currencyId, companyId, defaultPriceCode, // Parameters for subquery in CASE
				year, month, companyId, whsGroup, // Parameters for ELSE subquery in CASE
				currencyId, companyId, defaultPriceCode).
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
				A.item_level_4 AS item_level_4,
				CASE 
					WHEN ? IN (?, ?, ?) THEN
						CASE A.item_type
							WHEN ? THEN
								(SELECT TOP 1 price_list_amount FROM mtr_price_list
								WHERE is_active = 1 
								AND brand_id = B.brand_id 
								AND effective_date <= GETDATE()
								AND item_id = A.item_id
								AND currency_id = ? 
								AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
								AND price_list_code_id = ?
								ORDER BY effective_date DESC)
							ELSE
								(SELECT CASE ISNULL(price_current, 0)
										WHEN 0 THEN price_begin 
										ELSE price_current END AS HPP
								FROM mtr_group_stock 
								WHERE period_year = ? 
								AND period_month = ? 
								AND item_code = A.item_code 
								AND company_id = ?  
								AND whs_group = ?)
						END
					ELSE
						(SELECT TOP 1 price_list_amount FROM mtr_price_list
						WHERE is_active = 1 
						AND brand_id = B.brand_id 
						AND effective_date <= GETDATE()
						AND item_id = A.item_id 
						AND currency_id = ? 
						AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
						AND price_list_code_id = ?
						ORDER BY effective_date DESC)
				END AS PRICE
			`, year, month, companyId, whsGroup, // Parameters for AvailQty subquery
				billCode, BillCodeC, BillCodeInt, BillCodeNoCharge, // Parameters for CASE statement
				ItemService, currencyId, companyId, defaultPriceCode, // Parameters for subquery in CASE
				year, month, companyId, whsGroup, // Parameters for ELSE subquery in CASE
				currencyId, companyId, defaultPriceCode).
			Where("A.item_group_id = ? AND A.item_type = ? AND (A.item_class_id = ? OR A.item_class_id = ?) AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, ItmClsSublet, 1).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	case utils.LineTypeFee:
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
				A.item_level_4 AS item_level_4,
				CASE 
					WHEN ? IN (?, ?, ?) THEN
						CASE A.item_type
							WHEN ? THEN
								(SELECT TOP 1 price_list_amount FROM mtr_price_list
								WHERE is_active = 1 
								AND brand_id = B.brand_id 
								AND effective_date <= GETDATE()
								AND item_id = A.item_id
								AND currency_id = ? 
								AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
								AND price_list_code_id = ?
								ORDER BY effective_date DESC)
							ELSE
								(SELECT CASE ISNULL(price_current, 0)
										WHEN 0 THEN price_begin 
										ELSE price_current END AS HPP
								FROM mtr_group_stock 
								WHERE period_year = ? 
								AND period_month = ? 
								AND item_code = A.item_code 
								AND company_id = ?  
								AND whs_group = ?)
						END
					ELSE
						(SELECT TOP 1 price_list_amount FROM mtr_price_list
						WHERE is_active = 1 
						AND brand_id = B.brand_id 
						AND effective_date <= GETDATE()
						AND item_id = A.item_id 
						AND currency_id = ? 
						AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
						AND price_list_code_id = ?
						ORDER BY effective_date DESC)
				END AS PRICE
			`, year, month, companyId, whsGroup, // Parameters for AvailQty subquery
				billCode, BillCodeC, BillCodeInt, BillCodeNoCharge, // Parameters for CASE statement
				ItemService, currencyId, companyId, defaultPriceCode, // Parameters for subquery in CASE
				year, month, companyId, whsGroup, // Parameters for ELSE subquery in CASE
				currencyId, companyId, defaultPriceCode).
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
				A.item_level_4 AS item_level_4,
				CASE 
					WHEN ? IN (?, ?, ?) THEN
						CASE A.item_type
							WHEN ? THEN
								(SELECT TOP 1 price_list_amount FROM mtr_price_list
								WHERE is_active = 1 
								AND brand_id = B.brand_id 
								AND effective_date <= GETDATE()
								AND item_id = A.item_id
								AND currency_id = ? 
								AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
								AND price_list_code_id = ?
								ORDER BY effective_date DESC)
							ELSE
								(SELECT CASE ISNULL(price_current, 0)
										WHEN 0 THEN price_begin 
										ELSE price_current END AS HPP
								FROM mtr_group_stock 
								WHERE period_year = ? 
								AND period_month = ? 
								AND item_code = A.item_code 
								AND company_id = ?  
								AND whs_group = ?)
						END
					ELSE
						(SELECT TOP 1 price_list_amount FROM mtr_price_list
						WHERE is_active = 1 
						AND brand_id = B.brand_id 
						AND effective_date <= GETDATE()
						AND item_id = A.item_id 
						AND currency_id = ? 
						AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
						AND price_list_code_id = ?
						ORDER BY effective_date DESC)
				END AS PRICE
			`, year, month, companyId, whsGroup, // Parameters for AvailQty subquery
				billCode, BillCodeC, BillCodeInt, BillCodeNoCharge, // Parameters for CASE statement
				ItemService, currencyId, companyId, defaultPriceCode, // Parameters for subquery in CASE
				year, month, companyId, whsGroup, // Parameters for ELSE subquery in CASE
				currencyId, companyId, defaultPriceCode).
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

	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count rows",
			Err:        err,
		}
	}

	offset := (paginate.Page - 1) * paginate.Limit
	if err := baseQuery.Offset(offset).Limit(paginate.Limit).Find(&results).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve data",
			Err:        err,
		}
	}

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

	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

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

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total vehicle units",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

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

	return vehicleMasters, totalPages, int(totalRows), nil
}

func (r *LookupRepositoryImpl) GetVehicleUnitByID(tx *gorm.DB, vehicleID int, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		vehicleMasters []map[string]interface{}
		totalRows      int64
		totalPages     int
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	// Apply filters
	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

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
		Where("V.vehicle_id = ?", vehicleID)

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total vehicle units",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	err = query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&vehicleMasters).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get vehicle unit data by ID",
			Err:        err,
		}
	}

	return vehicleMasters, totalPages, int(totalRows), nil
}

func (r *LookupRepositoryImpl) GetVehicleUnitByChassisNumber(tx *gorm.DB, chassisNumber string, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		vehicleMasters []map[string]interface{}
		totalRows      int64
		totalPages     int
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	// Apply filters
	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

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
		Where("V.vehicle_chassis_number = ?", chassisNumber)

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total vehicle units",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	err = query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&vehicleMasters).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get vehicle unit data by chassis number",
			Err:        err,
		}
	}

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

	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

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

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total vehicle units",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

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

	return campaignMasters, totalPages, int(totalRows), nil
}

// usp_comLookUp
// IF @strEntity = 'WorkOrderService'--WO SERVICE
func (r *LookupRepositoryImpl) WorkOrderService(tx *gorm.DB, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		results []struct {
			WorkOrderNo    string
			WorkOrderDate  time.Time
			NoPolisi       string
			ChassisNo      string
			Brand          int
			Model          int
			Variant        int
			WorkOrderSysNo int
		}
		totalRows  int64
		totalPages int
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	filterStrings := []string{}
	filterValues := []interface{}{}
	if len(filters) > 0 {
		for _, filter := range filters {
			filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
			filterValues = append(filterValues, filter.ColumnValue)
		}
	}

	filterQuery := strings.Join(filterStrings, " AND ")
	if len(filterStrings) > 0 {
		tx = tx.Where(filterQuery, filterValues...)
	}

	query := tx.Table("trx_work_order_allocation AS A").
		Select("A.work_order_document_number AS WorkOrderNo, B.work_order_date AS WorkOrderDate, "+
			"B.vehicle_chassis_number AS ChassisNo, B.brand_id AS Brand, B.model_id AS Model, "+
			"B.variant_id AS Variant, A.work_order_system_number AS WorkOrderSysNo").
		Joins("LEFT JOIN trx_work_order AS B ON B.work_order_system_number = A.work_order_system_number").
		Where("A.service_status_id NOT IN (?, ?, ?, ?)", utils.SrvStatStop, utils.SrvStatAutoRelease, utils.SrvStatTransfer, utils.SrvStatQcPass)

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total vehicle units",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	err = query.
		Order("A.work_order_document_number").
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&results).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get vehicle unit master data",
			Err:        err,
		}
	}

	mappedResults := make([]map[string]interface{}, len(results))
	for i, result := range results {
		mappedResults[i] = map[string]interface{}{
			"work_order_document_number": result.WorkOrderNo,
			"work_order_date":            result.WorkOrderDate,
			"vehicle_tnkb":               result.NoPolisi,
			"vehicle_chassis_number":     result.ChassisNo,
			"brand_id":                   result.Brand,
			"model_id":                   result.Model,
			"variant_id":                 result.Variant,
			"work_order_system_number":   result.WorkOrderSysNo,
		}
	}

	return mappedResults, totalPages, int(totalRows), nil
}

// usp_comLookUp
// IF @strEntity =  'CustomerByTypeAndAddress'--CUSTOMER MASTER
func (r *LookupRepositoryImpl) CustomerByTypeAndAddress(tx *gorm.DB, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		customerMasters []map[string]interface{}
		totalRows       int64
		totalPages      int
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}

	filterQuery := strings.Join(filterStrings, " AND ")

	query := tx.Table("dms_microservices_general_dev.dbo.mtr_customer C").
		Select(`
			C.customer_id AS customer_id,
			C.customer_code AS customer_code,
			C.customer_name AS customer_name,
			CA.client_type_description AS client_type_description,
			A.address_street_1 AS address_1,
			A.address_street_2 AS address_2,
			A.address_street_3 AS address_3,
			C.id_phone_no AS id_phone_no
		`).
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_client_type CA ON C.client_type_id = CA.client_type_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_address AS A ON C.id_address_id = A.address_id").
		Where(filterQuery, filterValues...).
		Where("C.is_active = 1")

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total vehicle units",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	err = query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&customerMasters).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get vehicle unit master data",
			Err:        err,
		}
	}

	return customerMasters, totalPages, int(totalRows), nil
}

func (r *LookupRepositoryImpl) CustomerByTypeAndAddressByID(tx *gorm.DB, customerId int, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		customerMasters []map[string]interface{}
		totalRows       int64
		totalPages      int
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	query := tx.Table("dms_microservices_general_dev.dbo.mtr_customer C").
		Select(`
			C.customer_id AS customer_id,
			C.customer_code AS customer_code,
			C.customer_name AS customer_name,
			CA.client_type_description AS client_type_description,
			A.address_street_1 AS address_1,
			A.address_street_2 AS address_2,
			A.address_street_3 AS address_3,
			C.id_phone_no AS id_phone_no
		`).
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_client_type CA ON C.client_type_id = CA.client_type_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_address AS A ON C.id_address_id = A.address_id").
		Where(filterQuery, filterValues...).
		Where("C.customer_id = ?", customerId)

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total customers",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	err = query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&customerMasters).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get customer data",
			Err:        err,
		}
	}

	return customerMasters, totalPages, int(totalRows), nil
}

func (r *LookupRepositoryImpl) CustomerByTypeAndAddressByCode(tx *gorm.DB, customerCode string, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		customerMasters []map[string]interface{}
		totalRows       int64
		totalPages      int
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	query := tx.Table("dms_microservices_general_dev.dbo.mtr_customer C").
		Select(`
			C.customer_id,
			C.customer_code,
			C.customer_name,
			CA.client_type_description,
			A.address_street_1 AS address_1,
			A.address_street_2 AS address_2,
			A.address_street_3 AS address_3,
			C.id_phone_no
		`).
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_client_type CA ON C.client_type_id = CA.client_type_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_address AS A ON C.id_address_id = A.address_id").
		Where(filterQuery, filterValues...).
		Where("C.customer_code = ?", customerCode)

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total customers",
			Err:        err,
		}
	}

	if totalRows > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	if err := query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&customerMasters).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get customer data",
			Err:        err,
		}
	}

	return customerMasters, totalPages, int(totalRows), nil
}
