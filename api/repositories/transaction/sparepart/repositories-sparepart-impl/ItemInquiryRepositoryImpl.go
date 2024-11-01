package transactionsparepartrepositoryimpl

import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ItemInquiryRepositoryImpl struct {
}

func StartItemInquiryRepositoryImpl() transactionsparepartrepository.ItemInquiryRepository {
	return &ItemInquiryRepositoryImpl{}
}

// uspg_ItemInquiry_Select IF @Option = 2
func (i *ItemInquiryRepositoryImpl) GetAllItemInquiry(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entitiesItemDetail := masteritementities.ItemDetail{}
	response := []transactionsparepartpayloads.ItemInquiryGetAllPayloads{}
	var err error

	var newFilterCondition []utils.FilterCondition
	var companyId int
	var companySessionId int
	var itemId int
	var availableQuantityFrom *float64
	var availableQuantityTo *float64
	var salesPriceFrom *float64
	var salesPriceTo *float64
	for _, filter := range filterCondition {
		if strings.Contains(filter.ColumnField, "company_id") {
			companyId, _ = strconv.Atoi(filter.ColumnValue)
			continue
		}
		if strings.Contains(filter.ColumnField, "company_session_id") {
			companySessionId, _ = strconv.Atoi(filter.ColumnValue)
			continue
		}
		if strings.Contains(filter.ColumnField, "item_id") {
			itemId, _ = strconv.Atoi(filter.ColumnValue) // purposely added to newFilter
		}
		if strings.Contains(filter.ColumnField, "available_quantity_from") {
			availableQuantityFromTemp, _ := strconv.ParseFloat(filter.ColumnValue, 64)
			availableQuantityFrom = &availableQuantityFromTemp
			continue
		}
		if strings.Contains(filter.ColumnField, "available_quantity_to") {
			availableQuantityToTemp, _ := strconv.ParseFloat(filter.ColumnValue, 64)
			availableQuantityTo = &availableQuantityToTemp
			continue
		}
		if strings.Contains(filter.ColumnField, "sales_price_from") {
			salesPriceFromTemp, _ := strconv.ParseFloat(filter.ColumnValue, 64)
			salesPriceFrom = &salesPriceFromTemp
			continue
		}
		if strings.Contains(filter.ColumnField, "sales_price_to") {
			salesPriceToTemp, _ := strconv.ParseFloat(filter.ColumnValue, 64)
			salesPriceTo = &salesPriceToTemp
			continue
		}
		newFilterCondition = append(newFilterCondition, filter)
	}

	companyUrl := config.EnvConfigs.GeneralServiceUrl + "company/" + strconv.Itoa(companyId)
	companyPayloads := transactionsparepartpayloads.ItemInquiryCompanyResponse{}
	if err := utils.Get(companyUrl, &companyPayloads, nil); err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "company does not exist",
			Err:        err,
		}
	}
	companyCode := companyPayloads.CompanyCode

	companySessionUrl := config.EnvConfigs.GeneralServiceUrl + "company/" + strconv.Itoa(companySessionId)
	companySessionPayloads := transactionsparepartpayloads.ItemInquiryCompanyResponse{}
	if err := utils.Get(companySessionUrl, &companySessionPayloads, nil); err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "company session does not exist",
			Err:        err,
		}
	}

	companySessionReferenceUrl := config.EnvConfigs.GeneralServiceUrl + "company-reference/" + strconv.Itoa(companySessionId)
	companySessionReferencePayloads := transactionsparepartpayloads.ItemInquiryCompanyReferenceResponse{}
	if err := utils.Get(companySessionReferenceUrl, &companySessionReferencePayloads, nil); err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "company reference does not exist",
			Err:        err,
		}
	}
	currencyId := companySessionReferencePayloads.CurrencyId

	var priceListCodeId int
	entitiesItemPriceCode := masteritementities.ItemPriceCode{}
	err = tx.Model(&entitiesItemPriceCode).
		Select("item_price_code_id").
		Where(masteritementities.ItemPriceCode{ItemPriceCode: "A"}).
		First(&priceListCodeId).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	currentPeriodUrl := config.EnvConfigs.FinanceServiceUrl + "closing-period-company/current-period?company_id=" + strconv.Itoa(companyId) + "&closing_module_detail_code=SP"
	currentPeriodPayloads := transactionsparepartpayloads.ItemInquiryCurrentPeriodResponse{}
	if err := utils.Get(currentPeriodUrl, &currentPeriodPayloads, nil); err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "current period does not exist",
			Err:        err,
		}
	}
	periodYear := currentPeriodPayloads.PeriodYear
	periodMonth := currentPeriodPayloads.PeriodMonth

	var baseModelQuery *gorm.DB

	currentTime := time.Now().Truncate(24 * time.Hour)

	// View other company other than NMDI and KIA 1
	if companySessionId != companyId && !checkCompany(companyPayloads.CompanyCode) {
		lastEffectiveDate := tx.Table("mtr_item_price_list mipl1").
			Select("mipl1.effective_date").
			Where("mipl1.brand_id = mtr_item_detail.brand_id").
			Where("mipl1.currency_id = ?", currencyId).
			Where("mipl1.item_id = mtr_item_detail.item_id").
			Where("mipl1.company_id = mipl.company_id").
			Where("mipl1.effective_date < ?", currentTime).
			Where("mipl1.price_list_code_id = ?", priceListCodeId).
			Order("mipl1.effective_date DESC").
			Limit(1)

		maxProcessDate := tx.Table("mtr_moving_code_item mmci1").
			Select("MAX(mmci1.process_date)").
			Where("mmci1.company_id = ?", companyId).
			Where("mmci1.item_id = mtr_item_detail.item_id")

		baseModelQuery = tx.Model(&entitiesItemDetail).
			Select(`DISTINCT
				mtr_item_detail.item_detail_id,
				mtr_item_detail.item_id,
				mi.item_code,
				mi.item_name,
				mic.item_class_code,
				mtr_item_detail.brand_id,
				'' model_code,
				ISNULL(mwg.warehouse_group_code, '') warehouse_group_code,
				ISNULL(mwm.warehouse_code, '') warehouse_code,
				ISNULL(mwl.warehouse_location_code, '') warehouse_location_code,
				ISNULL(mipl.price_list_amount, 0) sales_price,
				CASE WHEN ISNULL(auto_pick_wms, '1') = '0' THEN 0 ELSE (ISNULL(mls.quantity_sales, 0) + ISNULL(mls.quantity_transfer_out, 0) + ISNULL(mls.quantity_robbing_out, 0) + ISNULL(mls.quantity_assembly_out, 0) + ISNULL(mls.quantity_allocated, 0)) END quantity_available,
				CASE WHEN ISNULL(mis.item_id, 0) = 0 THEN 'N' ELSE 'Y' END item_substitute,
				CASE WHEN ? = '200000' AND ISNULL(mmc.moving_code, '') != '' THEN ISNULL(mmc.moving_code_description, '') ELSE '' END moving_code,
				ISNULL(available_in_other_dealer, '') available_in_other_dealer
			`, companyCode).
			Joins("INNER JOIN mtr_item mi ON mi.item_id = mtr_item_detail.item_id").
			Joins("INNER JOIN mtr_item_class mic ON mic.item_class_id = mi.item_class_id").
			Joins("LEFT JOIN mtr_location_item mli ON mli.item_id = mi.item_id").
			Joins("INNER JOIN mtr_warehouse_group mwg ON mwg.warehouse_group_id = mli.warehouse_group_id").
			Joins("INNER JOIN mtr_warehouse_master mwm ON mwm.warehouse_id = mli.warehouse_id AND mwm.company_id = ?", companyId).
			Joins("INNER JOIN mtr_warehouse_location mwl ON mwl.warehouse_location_id = mli.warehouse_location_id").
			Joins(`LEFT JOIN mtr_item_price_list mipl ON mipl.brand_id = mtr_item_detail.brand_id
					AND mipl.currency_id = ?
					AND mipl.item_id = mtr_item_detail.item_id
					AND mipl.company_id = CASE WHEN ISNULL(mi.common_pricelist, ?) = ? THEN 0 ELSE ? END
					AND mipl.effective_date = (?)
					AND mipl.price_list_code_id = ?
					`, currencyId, true, true, companyId, lastEffectiveDate, priceListCodeId).
			Joins(`LEFT JOIN mtr_location_stock mls ON mls.period_year = ?
					AND mls.period_month = ?
					AND mls.warehouse_id = mli.warehouse_id
					AND mwm.company_id = ?
					AND mls.location_id = mli.warehouse_location_id
					AND mls.item_id = mli.item_id
					`, periodYear, periodMonth, companyId).
			Joins("LEFT JOIN mtr_item_substitute mis ON mis.item_id = mi.item_id AND mis.is_active = ?", true).
			Joins(`INNER JOIN mtr_warehouse_master mwm1 ON mwm1.company_id = ?
					AND mwm1.warehouse_id = mli.warehouse_id
					AND mwm1.warehouse_group_id = mli.warehouse_group_id
					AND ISNULL(mwm1.warehouse_sales_allow, ?) = ?
					AND mwm1.is_active = ?
					`, companyId, false, true, true).
			Joins(`LEFT JOIN mtr_moving_code_item mmci ON mmci.company_id = ?
					AND mmci.item_id = mtr_item_detail.item_id
					AND mmci.process_date = (?)
					`, companyId, maxProcessDate).
			Joins("LEFT JOIN mtr_moving_code mmc ON mmc.moving_code_id = mmci.moving_code_id")
	}

	// View own company
	if companySessionId == companyId {
		lastEffectiveDate := tx.Table("mtr_item_price_list mipl1").
			Select("mipl1.effective_date").
			Where("mipl1.brand_id = mtr_item_detail.brand_id").
			Where("mipl1.currency_id = ?", currencyId).
			Where("mipl1.item_id = mtr_item_detail.item_id").
			Where("mipl1.company_id = mipl.company_id").
			Where("mipl1.effective_date < ?", currentTime).
			Where("mipl1.price_list_code_id = ?", priceListCodeId).
			Order("mipl1.effective_date DESC").
			Limit(1)

		maxProcessDate := tx.Table("mtr_moving_code_item mmci1").
			Select("MAX(mmci1.process_date)").
			Where("mmci1.company_id = ?", companyId).
			Where("mmci1.item_id = mtr_item_detail.item_id")

		baseModelQuery = tx.Model(&entitiesItemDetail).
			Select(`DISTINCT
				mtr_item_detail.item_detail_id,
				mtr_item_detail.item_id,
				mi.item_code,
				mi.item_name,
				mic.item_class_code,
				mtr_item_detail.brand_id,
				'' model_code,
				ISNULL(mwg.warehouse_group_code, '') warehouse_group_code,
				ISNULL(mwm.warehouse_code, '') warehouse_code,
				ISNULL(mwl.warehouse_location_code, '') warehouse_location_code,
				ISNULL(mipl.price_list_amount, 0) sales_price,
				(ISNULL(mls.quantity_sales, 0) + ISNULL(mls.quantity_transfer_out, 0) + ISNULL(mls.quantity_robbing_out, 0) + ISNULL(mls.quantity_assembly_out, 0) + ISNULL(mls.quantity_allocated, 0)) quantity_available,
				CASE WHEN ISNULL(mis.item_id, 0) = 0 THEN 'N' ELSE 'Y' END item_substitute,
				CASE WHEN ISNULL(mmc.moving_code, '') != '' THEN ISNULL(mmc.moving_code_description, '') ELSE '' END moving_code,
				ISNULL(available_in_other_dealer, '') available_in_other_dealer
			`).
			Joins("INNER JOIN mtr_item mi ON mi.item_id = mtr_item_detail.item_id").
			Joins("INNER JOIN mtr_item_class mic ON mic.item_class_id = mi.item_class_id").
			Joins("LEFT JOIN mtr_location_item mli ON mli.item_id = mi.item_id", companyId).
			Joins("INNER JOIN mtr_warehouse_group mwg ON mwg.warehouse_group_id = mli.warehouse_group_id").
			Joins("INNER JOIN mtr_warehouse_master mwm ON mwm.warehouse_id = mli.warehouse_id AND mwm.company_id = ?", companyId).
			Joins("INNER JOIN mtr_warehouse_location mwl ON mwl.warehouse_location_id = mli.warehouse_location_id").
			Joins(`LEFT JOIN mtr_item_price_list mipl ON mipl.brand_id = mtr_item_detail.brand_id
					AND mipl.currency_id = ? 
					AND mipl.item_id = mtr_item_detail.item_id
					AND mipl.company_id = CASE WHEN ISNULL(mi.common_pricelist, ?) = ? THEN 0 ELSE ? END
					AND mipl.effective_date = (?)
					AND mipl.price_list_code_id = ?
					`, currencyId, true, true, companyId, lastEffectiveDate, priceListCodeId).
			Joins(`LEFT JOIN mtr_location_stock mls ON mls.period_year = ?
					AND mls.period_month = ?
					AND mls.warehouse_id = mli.warehouse_id
					AND mwm.company_id = ?
					AND mls.location_id = mli.warehouse_location_id
					AND mls.item_id = mli.item_id
					`, periodYear, periodMonth, companyId).
			Joins("LEFT JOIN mtr_item_substitute mis ON mis.item_id = mi.item_id AND mis.is_active = ?", true).
			Joins(`INNER JOIN mtr_warehouse_master mwm1 ON mwm1.company_id = ?
					AND mwm1.warehouse_id = mli.warehouse_id
					AND mwm1.warehouse_group_id = mli.warehouse_group_id
					AND ISNULL(mwm1.warehouse_sales_allow, ?) = ?
					AND mwm1.is_active = ?
					`, companyId, false, true, true).
			Joins(`LEFT JOIN mtr_moving_code_item mmci ON mmci.company_id = ?
					AND mmci.item_id = mtr_item_detail.item_id
					AND mmci.process_date = (?)
					`, companyId, maxProcessDate).
			Joins("LEFT JOIN mtr_moving_code mmc ON mmc.moving_code_id = mmci.moving_code_id")
	}

	// from other company view NMDI / KIA
	if companySessionId != companyId && checkCompany(companyPayloads.CompanyCode) {
		lastEffectiveDate := tx.Table("mtr_item_price_list mipl1").
			Select("mipl1.effective_date").
			Where("mipl1.brand_id = mtr_item_detail.brand_id").
			Where("mipl1.currency_id = ?", currencyId).
			Where("mipl1.item_id = mtr_item_detail.item_id").
			Where("mipl1.company_id = mipl.company_id").
			Where("mipl1.effective_date < ?", currentTime).
			Where("mipl1.price_list_code_id = ?", priceListCodeId).
			Order("mipl1.effective_date DESC").
			Limit(1)

		maxProcessDate := tx.Table("mtr_moving_code_item mmci1").
			Select("MAX(mmci1.process_date)").
			Where("mmci1.company_id = ?", companyId).
			Where("mmci1.item_id = mtr_item_detail.item_id")

		baseModelQuery = tx.Model(&entitiesItemDetail).
			Select(`DISTINCT
				mtr_item_detail.item_detail_id,
				mtr_item_detail.item_id,
				mi.item_code,
				mi.item_name,
				mic.item_class_code,
				mtr_item_detail.brand_id,
				'' model_code,
				'' warehouse_group_code,
				'' warehouse_code,
				'' warehouse_location_code,
				ISNULL(mipl.price_list_amount, 0) sales_price,
				0 quantity_available, --(ISNULL(mls.quantity_sales, 0) + ISNULL(mls.quantity_transfer_out, 0) + ISNULL(mls.quantity_robbing_out, 0) + ISNULL(mli.quantity_assembly_out, 0) + ISNULL(mli.quantity_allocated, 0)),
				CASE WHEN ISNULL(mis.item_id, 0) = 0 THEN 'N' ELSE 'Y' END item_substitute,
				CASE WHEN ? = '1516098' AND ISNULL(mmc.moving_code, '') != '' THEN ISNULL(mmc.moving_code_description, '') ELSE '' END moving_code,
				ISNULL(available_in_other_dealer, '') available_in_other_dealer
			`, companyCode).
			Joins("INNER JOIN mtr_item mi ON mi.item_id = mtr_item_detail.item_id").
			Joins("INNER JOIN mtr_item_class mic ON mic.item_class_id = mi.item_class_id").
			Joins("LEFT JOIN mtr_location_item mli ON mli.item_id = mi.item_id", companyId).
			Joins("INNER JOIN mtr_warehouse_group mwg ON mwg.warehouse_group_id = mli.warehouse_group_id").
			Joins("INNER JOIN mtr_warehouse_master mwm ON mwm.warehouse_id = mli.warehouse_id AND mwm.company_id = ?", companyId).
			Joins("INNER JOIN mtr_warehouse_location mwl ON mwl.warehouse_location_id = mli.warehouse_location_id").
			Joins(`LEFT JOIN mtr_item_price_list mipl ON mipl.brand_id = mtr_item_detail.brand_id
					AND mipl.currency_id = ? 
					AND mipl.item_id = mtr_item_detail.item_id
					AND mipl.company_id = CASE WHEN ISNULL(mi.common_pricelist, ?) = ? THEN 0 ELSE ? END
					AND mipl.effective_date = (?)
					AND mipl.price_list_code_id = ?
					`, currencyId, true, true, companyId, lastEffectiveDate, priceListCodeId).
			// Need LEFT JOIN to largo.ItemInquiry, table does not exist yet...
			Joins("LEFT JOIN mtr_item_substitute mis ON mis.item_id = mi.item_id AND mis.is_active = ?", true).
			Joins(`LEFT JOIN mtr_moving_code_item mmci ON mmci.company_id = ?
					AND mmci.item_id = mtr_item_detail.item_id
					AND mmci.process_date = (?)
					`, companyId, maxProcessDate).
			Joins("LEFT JOIN mtr_moving_code mmc ON mmc.moving_code_id = mmci.moving_code_id")
	}

	if itemId > 0 {
		processDate := tx.Table("mtr_moving_code_item mmci3").
			Select("mmci3.process_date").
			Where("mmci3.company_id = mmci2.company_id").
			Where("mmci3.item_id = mmci2.item_id").
			Order("mmci3.process_date DESC").
			Limit(1)

		movingCodeItem := tx.Table("mtr_moving_code_item mmci2").
			Select("*").
			Joins(`INNER JOIN mtr_location_stock mls1 ON mls1.company_id = mmci2.company_id
					AND mls1.item_id = mmci2.item_id
					AND mls1.period_year = ?
					AND mls1.period_month = ?
					AND mls1.company_id = ?
					`, periodYear, periodMonth, companyId).
			Joins("INNER JOIN mtr_moving_code mmc1 ON mmc1.moving_code_id = mmci2.moving_code_id").
			Where("mmc1.moving_code = '5'").
			Where("mmci2.process_date <= (?)", processDate).
			Where("(ISNULL(mls1.quantity_sales, 0) + ISNULL(mls1.quantity_transfer_out, 0) + ISNULL(mls1.quantity_robbing_out, 0) + ISNULL(mls1.quantity_assembly_out, 0) + ISNULL(mls1.quantity_allocated, 0)) > 0").
			Where("mmci2.item_id = mi1.item_id")

		availableInOtherDealer := tx.Table("mtr_item mi1").Select("mi1.item_id, CASE WHEN EXISTS(?) THEN 'Y' ELSE 'N' END available_in_other_dealer", movingCodeItem)

		baseModelQuery = baseModelQuery.Joins("LEFT JOIN (?) O ON mi.item_id = O.item_id", availableInOtherDealer)
	} else {
		baseModelQuery = baseModelQuery.Joins("LEFT JOIN (SELECT '' AS available_in_other_dealer) O ON O.available_in_other_dealer = ''")
	}

	itemGroupUrl := config.EnvConfigs.GeneralServiceUrl + "filter-item-group?item_group_code=IN"
	itemGroupPayloads := []transactionsparepartpayloads.ItemInquiryItemGroupResponse{}
	if err := utils.GetArray(itemGroupUrl, &itemGroupPayloads, nil); err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "item group does not exist",
			Err:        err,
		}
	}
	itemGroupInventoryId := itemGroupPayloads[0].ItemGroupId

	companyBrandUrl := config.EnvConfigs.GeneralServiceUrl + "company-brand/" + strconv.Itoa(companyId) + "?page=0&limit=1000000"
	companyBrandPayloads := []transactionsparepartpayloads.ItemInquiryCompanyBrandResponse{}
	if err := utils.GetArray(companyBrandUrl, &companyBrandPayloads, nil); err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "company brand does not exist",
			Err:        err,
		}
	}
	companyBrandIds := []int{}
	for _, companyBrand := range companyBrandPayloads {
		companyBrandIds = append(companyBrandIds, companyBrand.BrandId)
	}

	baseModelQuery = baseModelQuery.
		Where("mi.item_group_id = ?", itemGroupInventoryId).
		Where("mi.is_active = ?", true).
		Where("mtr_item_detail.brand_id IN ?", companyBrandIds)

	err = utils.ApplyFilter(baseModelQuery, newFilterCondition).Scan(&response).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	mapResponse := []map[string]interface{}{}
	for _, data := range response {
		if availableQuantityFrom != nil {
			if !(data.QuantityAvailable >= *availableQuantityFrom) {
				continue
			}
		}
		if availableQuantityTo != nil {
			if !(data.QuantityAvailable <= *availableQuantityTo) {
				continue
			}
		}
		if salesPriceFrom != nil {
			if !(data.SalesPrice >= *salesPriceFrom) {
				continue
			}
		}
		if salesPriceTo != nil {
			if !(data.SalesPrice <= *salesPriceTo) {
				continue
			}
		}
		temp := map[string]interface{}{
			"ItemDetailId":           data.ItemDetailId,
			"ItemId":                 data.ItemId,
			"ItemCode":               data.ItemCode,
			"ItemName":               data.ItemName,
			"ItemClassCode":          data.ItemClassCode,
			"BrandId":                data.BrandId,
			"ModelCode":              data.ModelCode,
			"WarehouseGroupCode":     data.WarehouseGroupCode,
			"WarehouseCode":          data.WarehouseCode,
			"WarehouseLocationCode":  data.WarehouseLocationCode,
			"SalesPrice":             data.SalesPrice,
			"QuantityAvailable":      data.QuantityAvailable,
			"ItemSubstitute":         data.ItemSubstitute,
			"MovingCode":             data.MovingCode,
			"AvailableInOtherDealer": data.AvailableInOtherDealer,
		}
		mapResponse = append(mapResponse, temp)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponse, &pages)

	pages.TotalPages = totalPages
	pages.TotalRows = int64(totalRows)

	finalJoinedData := []map[string]interface{}{}

	if len(paginatedData) > 0 {
		brandIds := []int{}
		brandIdsStr := ""

		for _, data := range paginatedData {
			if isNotInList(brandIds, data["BrandId"].(int)) {
				str := strconv.Itoa(data["BrandId"].(int))
				brandIdsStr += str + ","
				brandIds = append(brandIds, data["BrandId"].(int))
			}
		}

		brandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand-multi-id/" + brandIdsStr
		brandResponse := []transactionsparepartpayloads.ItemInquiryBrandResponse{}
		if err := utils.GetArray(brandUrl, &brandResponse, nil); err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "fail to fetch unit brand data",
				Err:        err,
			}
		}

		if len(brandResponse) == 0 {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNoContent,
				Err:        errors.New("unit brand does not exist"),
			}
		}

		joinedData := utils.DataFrameLeftJoin(paginatedData, brandResponse, "BrandId")

		// start usp_comToolTip @strEntity = 'ItemInquiryBrandModel'
		itemIds := []int{}

		for _, data := range joinedData {
			if isNotInList(itemIds, data["ItemId"].(int)) {
				itemIds = append(itemIds, data["ItemId"].(int))
			}
		}

		tooltips := []transactionsparepartpayloads.ItemInquiryGetAllToolTip{}
		for _, itemId := range itemIds {
			entitiesItemDetails := masteritementities.ItemDetail{}
			responseItemDetails := []transactionsparepartpayloads.ItemInquiryToolTip{}
			err = tx.Model(&entitiesItemDetails).
				Where(masteritementities.ItemDetail{ItemId: itemId}).
				Scan(&responseItemDetails).Error

			if err != nil {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNoContent,
					Err:        err,
				}
			}

			ttpBrandIds := []int{}
			ttpModelids := []int{}
			ttpBrandIdsStr := ""
			ttpModelIdsStr := ""

			for _, dataa := range responseItemDetails {
				if isNotInList(ttpBrandIds, dataa.BrandId) {
					str := strconv.Itoa(dataa.BrandId)
					ttpBrandIdsStr += str + ","
					ttpBrandIds = append(ttpBrandIds, dataa.BrandId)
				}
				if isNotInList(ttpModelids, dataa.ModelId) {
					str := strconv.Itoa(dataa.ModelId)
					ttpModelIdsStr += str + ","
					ttpModelids = append(ttpModelids, dataa.ModelId)
				}
			}

			ttpBrandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand-multi-id/" + ttpBrandIdsStr
			ttpBrandResponse := []transactionsparepartpayloads.ItemInquiryBrandResponse{}
			if err := utils.GetArray(ttpBrandUrl, &ttpBrandResponse, nil); err != nil {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "fail to fetch ttp unit brand data",
					Err:        err,
				}
			}
			if len(ttpBrandResponse) == 0 {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNoContent,
					Err:        errors.New("ttp unit brand does not exist"),
				}
			}

			ttpModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model-multi-id/" + ttpModelIdsStr
			ttpModelResponse := []transactionsparepartpayloads.ItemInquiryModelResponse{}
			if err := utils.GetArray(ttpModelUrl, &ttpModelResponse, nil); err != nil {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "fail to fetch ttp unit model data",
					Err:        err,
				}
			}
			if len(ttpModelResponse) == 0 {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNoContent,
					Err:        errors.New("ttp unit model does not exist"),
				}
			}

			ttpJoinedData := utils.DataFrameLeftJoin(responseItemDetails, ttpBrandResponse, "BrandId")
			ttpJoinedData2 := utils.DataFrameLeftJoin(ttpJoinedData, ttpModelResponse, "ModelId")

			tooltips = append(tooltips, transactionsparepartpayloads.ItemInquiryGetAllToolTip{
				ItemId:  itemId,
				Tooltip: ttpJoinedData2,
			})
		}
		// end usp_comToolTip @strEntity = 'ItemInquiryBrandModel'

		// manual left join data frame for adding tooltip resonse
		for i := 0; i < len(joinedData); i++ {
			joinedData[i]["Tooltip"] = []map[string]interface{}{}
			for j := 0; j < len(tooltips); j++ {
				if joinedData[i]["ItemId"].(int) == tooltips[j].ItemId {
					joinedData[i]["Tooltip"] = tooltips[j].Tooltip
					break
				}
			}
		}
		finalJoinedData = joinedData
	}
	pages.Rows = finalJoinedData

	return pages, nil
}

// Checks company is NMDI or KIA
func checkCompany(companyCode string) bool {
	for _, code := range []string{"3125098", "1516098"} {
		if code == companyCode {
			return true
		}
	}
	return false
}

func isNotInList(list []int, value int) bool {
	for _, v := range list {
		if v == value {
			return false
		}
	}
	return true
}
