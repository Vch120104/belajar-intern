package transactionsparepartrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	financeserviceapiutils "after-sales/api/utils/finance-service"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
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

	companyResponse, companyError := generalserviceapiutils.GetCompanyDataById(companyId)
	if companyError != nil || companyResponse.CompanyId == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("company does not exist"),
		}
	}
	companyCode := companyResponse.CompanyCode

	companySessionResponse, companySessionError := generalserviceapiutils.GetCompanyDataById(companySessionId)
	if companySessionError != nil || companySessionResponse.CompanyId == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("company session does not exist"),
		}
	}

	companySessionReferenceResponse, companySessionReferenceError := generalserviceapiutils.GetCompanyReferenceById(companySessionId)
	if companySessionReferenceError != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("company session reference does not exist"),
		}
	}
	currencyId := companySessionReferenceResponse.CurrencyId

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

	currentPeriodResponse, currentResponseError := financeserviceapiutils.GetOpenPeriodByCompany(companyId, "SP")
	if currentResponseError != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("current period does not exist"),
		}
	}
	periodYear := currentPeriodResponse.PeriodYear
	periodMonth := currentPeriodResponse.PeriodMonth

	var baseModelQuery *gorm.DB

	currentTime := time.Now().Truncate(24 * time.Hour)

	// View other company other than NMDI and KIA 1
	if companySessionId != companyId && !checkCompany(companyResponse.CompanyCode) {
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
				mwg.warehouse_group_id,
				mwm.warehouse_id,
				mwl.warehouse_location_id,
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
				mwg.warehouse_group_id,
				mwm.warehouse_id,
				mwl.warehouse_location_id,
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
	if companySessionId != companyId && checkCompany(companyResponse.CompanyCode) {
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
				0 warehouse_group_id,
				0 warehouse_id,
				0 warehouse_location_id,
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

	var itemGroupInventoryId int
	err = tx.Model(&masteritementities.ItemGroup{}).
		Where(masteritementities.ItemGroup{ItemGroupCode: "IN"}).
		Pluck("item_group_id", &itemGroupInventoryId).Error
	if err != nil || itemGroupInventoryId == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching item group inventory data",
			Err:        errors.New("error fetching item group inventory data"),
		}
	}

	companyBrandParams := generalserviceapiutils.CompanyBrandParams{Page: 0, Limit: 1000000}
	companyBrandResponse, companyBrandError := generalserviceapiutils.GetCompanyBrandByCompanyPagination(companyId, companyBrandParams)
	if companyBrandError != nil || len(companyBrandResponse) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "company brand does not exist",
			Err:        err,
		}
	}
	companyBrandIds := []int{}
	for _, companyBrand := range companyBrandResponse {
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
			"WarehouseGroupId":       data.WarehouseGroupId,
			"WarehouseGroupCode":     data.WarehouseGroupCode,
			"WarehouseId":            data.WarehouseId,
			"WarehouseCode":          data.WarehouseCode,
			"WarehouseLocationId":    data.WarehouseLocationId,
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

	finalJoinedData := []transactionsparepartpayloads.ItemInquiryGetAllResponse{}

	if len(paginatedData) > 0 {
		brandIds := []int{}
		for _, data := range paginatedData {
			brandIds = append(brandIds, data["BrandId"].(int))
		}

		brandResponse, brandError := salesserviceapiutils.GetUnitBrandByMultiId(brandIds)
		if brandError != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("error fetching unit brand data"),
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

			for _, dataa := range responseItemDetails {
				ttpBrandIds = append(ttpBrandIds, dataa.BrandId)
				ttpModelids = append(ttpModelids, dataa.ModelId)
			}

			ttpBrandResponse, ttpBrandError := salesserviceapiutils.GetUnitBrandByMultiId(ttpBrandIds)
			if ttpBrandError != nil {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "fail to fetch ttp unit brand data",
					Err:        err,
				}
			}
			if len(ttpBrandResponse) == 0 {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "ttp unit model doesn not exist",
				}
			}

			ttpModelResponse, ttpModelError := salesserviceapiutils.GetUnitModelByMultiId(ttpModelids)
			if ttpModelError != nil {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "fail to fetch ttp unit brand data",
					Err:        err,
				}
			}
			if len(ttpModelResponse) == 0 {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
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

		for i := 0; i < len(joinedData); i++ {
			temp := transactionsparepartpayloads.ItemInquiryGetAllResponse{
				ItemDetailId:           joinedData[i]["ItemDetailId"].(int),
				ItemId:                 joinedData[i]["ItemId"].(int),
				ItemCode:               joinedData[i]["ItemCode"].(string),
				ItemName:               joinedData[i]["ItemName"].(string),
				ItemClassCode:          joinedData[i]["ItemClassCode"].(string),
				BrandId:                joinedData[i]["BrandId"].(int),
				BrandCode:              joinedData[i]["BrandCode"].(string),
				ModelCode:              joinedData[i]["ModelCode"].(string),
				WarehouseGroupId:       joinedData[i]["WarehouseGroupId"].(int),
				WarehouseGroupCode:     joinedData[i]["WarehouseGroupCode"].(string),
				WarehouseId:            joinedData[i]["WarehouseId"].(int),
				WarehouseCode:          joinedData[i]["WarehouseCode"].(string),
				WarehouseLocationId:    joinedData[i]["WarehouseLocationId"].(int),
				WarehouseLocationCode:  joinedData[i]["WarehouseLocationCode"].(string),
				SalesPrice:             joinedData[i]["SalesPrice"],
				QuantityAvailable:      joinedData[i]["QuantityAvailable"],
				ItemSubstitute:         joinedData[i]["ItemSubstitute"].(string),
				MovingCode:             joinedData[i]["MovingCode"].(string),
				AvailableInOtherDealer: joinedData[i]["AvailableInOtherDealer"].(string),
				Tooltip:                joinedData[i]["Tooltip"].([]map[string]interface{}),
			}
			finalJoinedData = append(finalJoinedData, temp)
		}
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

// uspg_ItemInquiry_Select IF @Option = 1
func (i *ItemInquiryRepositoryImpl) GetByIdItemInquiry(tx *gorm.DB, filter transactionsparepartpayloads.ItemInquiryGetByIdFilter) (transactionsparepartpayloads.ItemInquiryGetByIdResponse, *exceptions.BaseErrorResponse) {
	response := transactionsparepartpayloads.ItemInquiryGetByIdResponse{}

	companyResponse, companyError := generalserviceapiutils.GetCompanyDataById(filter.CompanyId)
	if companyError != nil || companyResponse.CompanyId == 0 {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching company data",
			Err:        errors.New("error fetching company data"),
		}
	}

	// validate only NMDI is allowed to fetch NMDI data
	if filter.CompanyId != filter.CompanySessionId && companyResponse.CompanyCode == "3125098" {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Other company is not allowed to fetch NMDI data",
			Err:        errors.New("other company is not allowed to fetch NMDI data"),
		}
	}

	approvalResponse, approvalError := generalserviceapiutils.GetApprovalStatusByCode("20") // approved status
	if approvalError != nil || approvalResponse.ApprovalStatusId == 0 {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching approval status accepted data",
			Err:        errors.New("error fetching approval status accepted data"),
		}
	}

	woStatusDescriptions := []string{"Draft", "Closed"}
	woStatusResponse, woStatusError := generalserviceapiutils.GetWorkOrderStatusByMultiDesc(woStatusDescriptions)
	if woStatusError != nil || len(woStatusDescriptions) != len(woStatusResponse) {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching wo status data",
			Err:        errors.New("error fetching wo status data"),
		}
	}

	lineTypePackageResp, lineTypePackageError := generalserviceapiutils.GetLineTypeByCode("0")
	if lineTypePackageError != nil || lineTypePackageResp.LineTypeId == 0 {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching line type package data",
			Err:        errors.New("error fetching line type package data"),
		}
	}

	lineTypeOperationResp, lineTypeOperationError := generalserviceapiutils.GetLineTypeByCode("1")
	if lineTypeOperationError != nil || lineTypeOperationResp.LineTypeId == 1 {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching line type operation data",
			Err:        errors.New("error fetching line type operation data"),
		}
	}

	periodDemand := utils.SuggorDemandPeriod
	// periodDemandNeg := -1 * utils.SuggorDemandPeriod
	approvalApprovedId := approvalResponse.ApprovalStatusId
	woStatusDraftId := woStatusResponse[0].WorkOrderStatusId
	woStatusClosedId := woStatusResponse[1].WorkOrderStatusId
	lineTypePackageId := lineTypePackageResp.LineTypeId
	lineTypeOperationId := lineTypeOperationResp.LineTypeId

	// sixMonthsUntil := time.Now().AddDate(0, -6, 0)
	today := time.Now()

	var warehouseGroupId int
	err := tx.Model(&masteritementities.ItemLocation{}).
		Where(masteritementities.ItemLocation{ItemId: filter.ItemId, WarehouseId: filter.WarehouseId}).
		Order("item_location_id DESC").
		Limit(1).
		Pluck("warehouse_group_id", &warehouseGroupId).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching warehouse group data",
			Err:        err,
		}
	}

	quantityDemand := 1 // from query amDemandStock
	periodYear := strconv.Itoa(int(time.Now().Year()))
	periodMonth := strconv.Itoa(int(time.Now().Month()))
	if len(periodMonth) == 1 {
		periodMonth = "0" + periodMonth
	}

	var orderTypeRegularId int
	err = tx.Model(&masterentities.OrderType{}).
		Where(masterentities.OrderType{OrderTypeCode: "R"}). // regular
		Pluck("order_type_id", &orderTypeRegularId).Error
	if err != nil || orderTypeRegularId == 0 {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching order type regular data",
			Err:        errors.New("error fetching order type regular data"),
		}
	}

	var orderTypeEmergencyId int
	err = tx.Model(&masterentities.OrderType{}).
		Where(masterentities.OrderType{OrderTypeCode: "E"}). // emergency
		Pluck("order_type_id", &orderTypeEmergencyId).Error
	if err != nil || orderTypeEmergencyId == 0 {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching order type emergency data",
			Err:        errors.New("error fetching order type emergency data"),
		}
	}

	var itemGroupInventoryId int
	err = tx.Model(&masteritementities.ItemGroup{}).
		Where(masteritementities.ItemGroup{ItemGroupCode: "IN"}). // inventory
		Pluck("item_group_id", &itemGroupInventoryId).Error
	if err != nil || itemGroupInventoryId == 0 {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching item group inventory data",
			Err:        errors.New("error fetching item group inventory data"),
		}
	}

	var supplySlipStatusId int
	err = tx.Model(&masterentities.SupplySlipStatus{}).
		Where(masterentities.SupplySlipStatus{SupplySlipStatusCode: "99"}). // complete status
		Pluck("supply_slip_status_id", &supplySlipStatusId).Error
	if err != nil || supplySlipStatusId == 0 {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching supply slip status data",
			Err:        errors.New("error fetching supply slip status data"),
		}
	}

	viewLocStock := tx.Table("mtr_location_stock mls").
		Select(`
			mls.period_year,
			mls.period_month,
			mls.company_id,
			mls.warehouse_group_id,
			mls.warehouse_id,
			mls.location_id,
			mls.item_id,
			(
				ISNULL(quantity_sales, 0) +
				ISNULL(quantity_transfer_out, 0) +
				ISNULL(quantity_robbing_out, 0) +
				ISNULL(quantity_assembly_out, 0) +
				ISNULL(quantity_allocated, 0)
			) AS quantity_available,
			mls.quantity_begin,
			mls.quantity_sales,
			mls.quantity_sales_return,
			mls.quantity_purchase,
			mls.quantity_purchase_return,
			mls.quantity_transfer_in,
			mls.quantity_transfer_out,
			mls.quantity_in_transit,
			mls.quantity_claim_in,
			mls.quantity_claim_out,
			mls.quantity_adjustment,
			mls.quantity_allocated,
			(
				ISNULL(quantity_sales, 0) +
				ISNULL(quantity_transfer_out, 0) +
				ISNULL(quantity_claim_out, 0) +
				ISNULL(quantity_robbing_out, 0) +
				ISNULL(quantity_assembly_out, 0)
			) AS quantity_on_hand
		`).
		Joins("LEFT JOIN mtr_warehouse_master mwm1 ON mwm1.company_id = mls.company_id AND mwm1.warehouse_id = mls.warehouse_id")

	lastEffectiveDate := tx.Table("mtr_item_price_list mipl1").
		Select("MAX(mipl1.effective_date)").
		Where("mipl1.brand_id = mtr_item_detail.brand_id").
		Where("mipl1.currency_id = ?", filter.CurrencyId).
		Where("mipl1.item_id = mtr_item_detail.item_id").
		Where("mipl1.effective_date < ?", today)

	lastProcessDate := tx.Table("mtr_moving_code_item mmci1").
		Select("MAX(mmci1.process_date)").
		Where("mmci1.company_id = ?", filter.CompanyId).
		Where("mmci1.item_id = mmci.item_id")

	tableDiscount := tx.Table("mtr_discount md").
		Select(`
			md.discount_code_id,
			md.discount_code,
			CASE WHEN mdp.order_type_id = ? THEN md.discount_description ELSE '' END AS discount_emergency,
			CASE WHEN mdp.order_type_id = ? THEN md.discount_description ELSE '' END AS discount_regular
		`, orderTypeEmergencyId, orderTypeRegularId).
		Joins("LEFT JOIN mtr_discount_percent mdp ON mdp.discount_code_id = md.discount_code_id")

	querySalesOrder := tx.Table("trx_sales_order_detail tsod").
		Select("SUM(ISNULL(tsod.quantity_demand, 0) - ISNULL(tsod.quantity_pick, 0))").
		Joins("INNER JOIN trx_sales_order tso ON tso.sales_order_system_number = tsod.sales_order_system_number").
		Where("tso.company_id = ?", filter.CompanyId).
		Where("tsod.item_id = mtr_item_detail.item_id").
		Where("tso.sales_order_status_id = ?", approvalApprovedId)

	var baseModelQuery *gorm.DB
	if filter.WarehouseId > 0 && filter.WarehouseLocationId > 0 {
		queryWorkOrder := tx.Table("trx_work_order_detail twod").
			Select("SUM(ISNULL(twod.frt_quantity, 0) - ISNULL(twod.supply_quantity, 0))").
			Joins("INNER JOIN trx_work_order two ON two.work_order_system_number = twod.work_order_system_number").
			Where("two.company_id = ?", filter.CompanyId).
			Where("two.work_order_status_id NOT IN (?,?)", woStatusClosedId, woStatusDraftId).
			Where("twod.line_type_id NOT IN (?, ?)", lineTypeOperationId, lineTypePackageId).
			Where("twod.frt_quantity != twod.supply_quantity").
			Where("twod.operation_item_id = ?", filter.ItemId)

		baseModelQuery = tx.Model(&masteritementities.ItemDetail{}).
			Select(`
			mtr_item_detail.item_detail_id,
			ISNULL(mwm.company_id, ?) company_id,
			? AS period_year,
			? AS period_month,
			mtr_item_detail.item_id,
			mi.item_code,
			mi.item_name,
			mtr_item_detail.brand_id,
			mli.warehouse_group_id,
			mwg.warehouse_group_code,
			mwg.warehouse_group_name,
			mli.warehouse_id,
			mwm.warehouse_code,
			mwm.warehouse_name,
			mli.warehouse_location_id,
			mwl.warehouse_location_code,
			mwl.warehouse_location_name,
			mipl.price_list_amount,
			vls.quantity_available,
			vls.quantity_begin,
			vls.quantity_sales,
			vls.quantity_sales_return,
			vls.quantity_purchase,
			vls.quantity_purchase_return,
			vls.quantity_transfer_in,
			vls.quantity_transfer_out,
			vls.quantity_in_transit,
			vls.quantity_claim_in,
			vls.quantity_claim_out,
			vls.quantity_adjustment,
			vls.quantity_allocated,
			vls.quantity_on_hand,
			mic.quantity_on_order,
			mgs.price_current,
			CASE WHEN ISNULL(mmci.moving_code_id, 0) != 0 THEN ISNULL(mmc.moving_code, '') ELSE '' END AS moving_code,
			ISNULL((?), 0) + ISNULL((?), 0) AS quantity_back_order,
			ISNULL(ISNULL(?, 0) * (ISNULL(mfm.forecast_master_lead_time, 0) + ISNULL(forecast_master_order_cycle, 0)), 0) + ISNULL(ISNULL(?, 0) * ISNULL(mfm.forecast_master_safety_factor, 0), 0) AS quantity_max,
			ISNULL(ISNULL(?, 0) * ISNULL(mfm.forecast_master_safety_factor, 0), 0) AS quantity_min,
			mi.discount_id,
			td.discount_code,
			td.discount_emergency,
			td.discount_regular,
			mics.item_class_name,
			mi.is_technical_defect
		`, filter.CompanyId, periodYear, periodMonth, querySalesOrder, queryWorkOrder, quantityDemand, quantityDemand, quantityDemand).
			Joins("LEFT JOIN mtr_item mi ON mi.item_id = mtr_item_detail.item_id").
			Joins("LEFT JOIN mtr_location_item mli ON mli.item_id = mtr_item_detail.item_id").
			Joins(`LEFT JOIN (?) vls ON vls.period_year = ? AND vls.period_month = ?
														AND vls.company_id = ?
														AND vls.warehouse_group_id = ?
														AND vls.warehouse_id = mli.warehouse_id
														AND vls.location_id = mli.warehouse_location_id
														AND vls.item_id = mtr_item_detail.item_id
														`, viewLocStock, periodYear, periodMonth, filter.CompanyId, warehouseGroupId).
			Joins(`LEFT JOIN mtr_item_price_list mipl ON mipl.brand_id = mtr_item_detail.brand_id AND mipl.currency_id = ?
																							AND mipl.item_id = mtr_item_detail.item_id
																							AND mipl.company_id = CASE WHEN mi.common_pricelist = ? THEN 0 ELSE ? END
																							AND mipl.effective_date = (?)`, filter.CurrencyId, true, filter.CompanyId, lastEffectiveDate).
			Joins("LEFT JOIN mtr_warehouse_group mwg ON mwg.warehouse_group_id = ?", warehouseGroupId).
			Joins("LEFT JOIN mtr_warehouse_master mwm ON mwm.company_id = ? AND mwm.warehouse_id = mli.warehouse_id", filter.CompanyId).
			Joins("LEFT JOIN mtr_warehouse_location mwl ON mwl.warehouse_id = mli.warehouse_id AND mwl.warehouse_location_id = mli.warehouse_location_id").
			Joins(`LEFT JOIN mtr_item_cycle mic ON mic.company_id = ? AND mic.period_year = ?
																AND mic.period_month = ?
																AND mic.item_id = mtr_item_detail.item_id
																`, filter.CompanyId, periodYear, periodMonth).
			Joins(`LEFT JOIN mtr_group_stock mgs ON mgs.company_id = ? AND mgs.period_year = ?
																AND mgs.period_month = ?
																AND mgs.warehouse_group_id = ?
																AND mgs.item_id = ?`, filter.CompanyId, periodYear, periodMonth, warehouseGroupId, filter.ItemId).
			Joins(`LEFT JOIN mtr_moving_code_item mmci ON mmci.company_id = ? AND mmci.item_id = ? AND mmci.process_date = (?)`, filter.CompanyId, filter.ItemId, lastProcessDate).
			Joins("LEFT JOIN mtr_moving_code mmc ON mmc.moving_code_id = mmci.moving_code_id").
			Joins(`LEFT JOIN mtr_forecast_master mfm ON mfm.company_id = ? AND mfm.supplier_id = mi.supplier_id
																		AND mfm.order_type_id = ?
																		AND mfm.moving_code_id = mmci.moving_code_id`, filter.CompanyId, orderTypeRegularId).
			Joins("LEFT JOIN (?) td ON td.discount_code_id = mi.discount_id", tableDiscount).
			Joins("LEFT JOIN mtr_item_class mics ON mics.item_class_id = mi.item_class_id").
			Where("mi.item_group_id = ?", itemGroupInventoryId).
			Where("mtr_item_detail.item_id = ?", filter.ItemId).
			Where("mtr_item_detail.brand_id = CASE WHEN ? = 0 THEN mtr_item_detail.brand_id ELSE ? END", filter.BrandId, filter.BrandId).
			Where("mli.warehouse_id = ?", filter.WarehouseId).
			Where("mli.warehouse_location_id = ?", filter.WarehouseLocationId)
	} else {
		querySupplySlipReturn := tx.Table("trx_supply_slip_return_detail tssrd").
			Select("tssrd.quantity_return, tssd.work_order_system_number, tssd.work_order_line_number").
			Joins("INNER JOIN trx_supply_slip_detail tssd ON tssd.supply_detail_system_number = tssrd.supply_detail_system_number").
			Joins("INNER JOIN trx_supply_slip_return tssr ON tssr.supply_return_system_number = tssrd.supply_return_system_number AND tssr.supply_return_status_id = ?", supplySlipStatusId)

		queryWorkOrder := tx.Table("trx_work_order_detail twod").
			Select("SUM(ISNULL(twod.frt_quantity, 0) - (ISNULL(twod.supply_quantity, 0) - ISNULL(ssr.quantity_return, 0)))").
			Joins("INNER JOIN trx_work_order two ON two.work_order_system_number = twod.work_order_system_number").
			Joins("LEFT JOIN (?) ssr ON ssr.work_order_system_number = twod.work_order_system_number AND ssr.work_order_line_number = twod.work_order_operation_item_line", querySupplySlipReturn).
			Where("two.company_id = ?", filter.CompanyId).
			Where("two.work_order_status_id NOT IN (?,?)", woStatusClosedId, woStatusDraftId).
			Where("twod.line_type_id NOT IN(?,?)", lineTypeOperationId, lineTypePackageId).
			Where("twod.operation_item_id = ?", filter.ItemId)

		baseModelQuery = tx.Model(&masteritementities.ItemDetail{}).
			Select(`
				mtr_item_detail.item_detail_id,
				ISNULL(mwm.company_id, ?) company_id,
				? AS period_year,
				? AS period_month,
				mtr_item_detail.item_id,
				mi.item_code,
				mi.item_name,
				mtr_item_detail.brand_id,
				mli.warehouse_group_id,
				mwg.warehouse_group_code,
				mwg.warehouse_group_name,
				mli.warehouse_id,
				mwm.warehouse_code,
				mwm.warehouse_name,
				mli.warehouse_location_id,
				mwl.warehouse_location_code,
				mwl.warehouse_location_name,
				mipl.price_list_amount,
				vls.quantity_available,
				vls.quantity_begin,
				vls.quantity_sales,
				vls.quantity_sales_return,
				vls.quantity_purchase,
				vls.quantity_purchase_return,
				vls.quantity_transfer_in,
				vls.quantity_transfer_out,
				vls.quantity_in_transit,
				vls.quantity_claim_in,
				vls.quantity_claim_out,
				vls.quantity_adjustment,
				vls.quantity_allocated,
				vls.quantity_on_hand,
				mic.quantity_on_order,
				mgs.price_current,
				CASE WHEN ISNULL(mmci.moving_code_id, 0) != 0 THEN ISNULL(mmc.moving_code, '') ELSE '' END AS moving_code,
				ISNULL(ISNULL((?), 0) + ISNULL((?), 0), 0) AS quantity_back_order,
				ISNULL((SUM(ISNULL(?, 0))/?) * (ISNULL(mfm.forecast_master_lead_time, 0) + ISNULL(mfm.forecast_master_order_cycle, 0)), 0) + ISNULL((SUM(ISNULL(?,0))/?) * ISNULL(mfm.forecast_master_safety_factor, 0), 0) AS quantity_max,
				ISNULL((SUM(ISNULL(?,0))/?) * ISNULL(mfm.forecast_master_safety_factor, 0), 0) AS quantity_min,
				mi.discount_id,
				td.discount_code,
				td.discount_emergency,
				td.discount_regular,
				mics.item_class_name,
				mi.is_technical_defect
			`, filter.CompanyId, periodYear, periodMonth, querySalesOrder, queryWorkOrder, quantityDemand, periodDemand, quantityDemand, periodDemand, quantityDemand, periodDemand).
			Joins("LEFT JOIN mtr_item mi ON mi.item_id = mtr_item_detail.item_id").
			Joins("LEFT JOIN mtr_location_item mli ON mli.item_id = mtr_item_detail.item_id").
			Joins(`LEFT JOIN (?) vls ON vls.period_year = ? AND vls.period_month = ?
														AND vls.company_id = ?
														AND vls.warehouse_group_id = ?
														AND vls.warehouse_id = mli.warehouse_id
														AND vls.location_id = mli.warehouse_location_id
														AND vls.item_id = mtr_item_detail.item_id
														`, viewLocStock, periodYear, periodMonth, filter.CompanyId, warehouseGroupId).
			Joins(`LEFT JOIN mtr_item_price_list mipl ON mipl.brand_id = mtr_item_detail.brand_id AND mipl.currency_id = ?
																							AND mipl.item_id = mtr_item_detail.item_id
																							AND mipl.company_id = CASE WHEN mi.common_pricelist = ? THEN 0 ELSE ? END
																							AND mipl.effective_date = (?)`, filter.CurrencyId, true, filter.CompanyId, lastEffectiveDate).
			Joins("LEFT JOIN mtr_warehouse_group mwg ON mwg.warehouse_group_id = ?", warehouseGroupId).
			Joins("LEFT JOIN mtr_warehouse_master mwm ON mwm.company_id = ? AND mwm.warehouse_id = mli.warehouse_id", filter.CompanyId).
			Joins("LEFT JOIN mtr_warehouse_location mwl ON mwl.warehouse_id = mli.warehouse_id AND mwl.warehouse_location_id = mli.warehouse_location_id").
			Joins(`LEFT JOIN mtr_item_cycle mic ON mic.company_id = ? AND mic.period_year = ?
																AND mic.period_month = ?
																AND mic.item_id = mtr_item_detail.item_id
																`, filter.CompanyId, periodYear, periodMonth).
			Joins(`LEFT JOIN mtr_group_stock mgs ON mgs.company_id = ? AND mgs.period_year = ?
																AND mgs.period_month = ?
																AND mgs.warehouse_group_id = ?
																AND mgs.item_id = ?`, filter.CompanyId, periodYear, periodMonth, warehouseGroupId, filter.ItemId).
			Joins(`LEFT JOIN mtr_moving_code_item mmci ON mmci.company_id = ? AND mmci.item_id = ? AND mmci.process_date = (?)`, filter.CompanyId, filter.ItemId, lastProcessDate).
			Joins("LEFT JOIN mtr_moving_code mmc ON mmc.moving_code_id = mmci.moving_code_id").
			Joins(`LEFT JOIN mtr_forecast_master mfm ON mfm.company_id = ? AND mfm.supplier_id = mi.supplier_id
																		AND mfm.order_type_id = ?
																		AND mfm.moving_code_id = mmci.moving_code_id`, filter.CompanyId, orderTypeRegularId).
			Joins("LEFT JOIN (?) td ON td.discount_code_id = mi.discount_id", tableDiscount).
			Joins("LEFT JOIN mtr_item_class mics ON mics.item_class_id = mi.item_class_id").
			Where("mi.item_group_id = ?", itemGroupInventoryId).
			Where("mtr_item_detail.item_id = ?", filter.ItemId).
			Where("mtr_item_detail.brand_id = CASE WHEN ? = 0 THEN mtr_item_detail.brand_id ELSE ? END", filter.BrandId, filter.BrandId).
			Group(`
				mtr_item_detail.item_detail_id,
				mwm.company_id,
				vls.period_year,
				vls.period_month,
				mtr_item_detail.item_id,
				mi.item_code,
				mi.item_name,
				mtr_item_detail.brand_id,
				mli.warehouse_group_id,
				mwg.warehouse_group_code,
				mwg.warehouse_group_name,
				mli.warehouse_id,
				mwm.warehouse_code,
				mwm.warehouse_name,
				mli.warehouse_location_id,
				mwl.warehouse_location_code,
				mwl.warehouse_location_name,
				mipl.price_list_amount,
				vls.quantity_available,
				vls.quantity_begin,
				vls.quantity_sales,
				vls.quantity_sales_return,
				vls.quantity_purchase,
				vls.quantity_purchase_return,
				vls.quantity_transfer_in,
				vls.quantity_transfer_out,
				vls.quantity_in_transit,
				vls.quantity_claim_in,
				vls.quantity_claim_out,
				vls.quantity_adjustment,
				vls.quantity_allocated,
				vls.quantity_on_hand,
				mic.quantity_on_order,
				mgs.price_current,
				mmci.moving_code_id,
				mmc.moving_code,
				mfm.forecast_master_lead_time,
				mfm.forecast_master_order_cycle,
				mfm.forecast_master_safety_factor,
				mi.discount_id,
				td.discount_code,
				td.discount_emergency,
				td.discount_regular,
				mics.item_class_name,
				mi.is_technical_defect
			`)
	}

	err = baseModelQuery.First(&response).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "data not found",
				Err:        err,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error when fetching item inquiry data",
			Err:        err,
		}
	}

	return response, nil
}
