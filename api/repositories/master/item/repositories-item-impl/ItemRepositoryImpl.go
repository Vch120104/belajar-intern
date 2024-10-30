package masteritemrepositoryimpl

import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ItemRepositoryImpl struct {
}

func StartItemRepositoryImpl() masteritemrepository.ItemRepository {
	return &ItemRepositoryImpl{}
}

// CheckItemCodeExist implements masteritemrepository.ItemRepository.
func (r *ItemRepositoryImpl) CheckItemCodeExist(tx *gorm.DB, itemCode string, itemGroupId int, commonPriceList bool, brandId int) (bool, int, int, *exceptions.BaseErrorResponse) {
	model := masteritementities.Item{}

	if err := tx.Model(model).Select("mtr_item.item_code,mtr_item.item_id,mtr_item.item_class_id").
		Joins("ItemDetail", tx.Select("1")).
		Where(masteritementities.Item{ItemCode: itemCode, ItemGroupId: itemGroupId, CommonPricelist: commonPriceList}).
		Where("ItemDetail.brand_id = ?", brandId).
		First(&model).Error; err != nil {
		return false, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}
	return true, model.ItemId, model.ItemClassId, nil
}

// GetUomItemDropDown implements masteritemrepository.ItemRepository.
func (r *ItemRepositoryImpl) GetUomDropDown(tx *gorm.DB, uomTypeId int) ([]masteritempayloads.UomDropdownResponse, *exceptions.BaseErrorResponse) {
	model := masteritementities.Uom{}
	responses := []masteritempayloads.UomDropdownResponse{}
	err := tx.Model(model).Where(masteritementities.Uom{UomTypeId: uomTypeId}).Scan(&responses).Error

	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return responses, nil
}

// GetUomTypeDropDown implements masteritemrepository.ItemRepository.
func (r *ItemRepositoryImpl) GetUomTypeDropDown(tx *gorm.DB) ([]masteritempayloads.UomTypeDropdownResponse, *exceptions.BaseErrorResponse) {
	model := masteritementities.UomType{}
	responses := []masteritempayloads.UomTypeDropdownResponse{}
	err := tx.Model(model).Scan(&responses).Error

	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return responses, nil
}

func (r *ItemRepositoryImpl) GetAllItemListTransLookup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	entites := masteritementities.Item{}
	response := []masteritempayloads.ItemListTransLookUp{}

	baseModelQuery := tx.Model(&entites).
		Select(`
			mtr_item.item_id,
			mtr_item.item_code,
			mtr_item.item_name,
			mtr_item.item_class_id,
			ic.item_class_code,
			ic.item_class_name,
			mtr_item.item_type_id,
			it.item_type_code,
			mtr_item.item_level_1,
			mtr_item.item_level_2,
			mtr_item.item_level_3,
			mtr_item.item_level_4`).
		Joins("INNER JOIN mtr_item_class ic ON ic.item_class_id = mtr_item.item_class_id").
		Joins("INNER JOIN mtr_item_type it ON it.item_type_id = mtr_item.item_type_id")

	whereQuery := utils.ApplyFilterSearch(baseModelQuery, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&entites, &pages, whereQuery)).Scan(&response).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.Rows = response

	return pages, nil
}

func (r *ItemRepositoryImpl) GetAllItemSearch(tx *gorm.DB, filterCondition []utils.FilterCondition, itemIDs []string, supplierIDs []string, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tableStruct := masteritempayloads.ItemSearch{}

	var supplierCode, supplierName string
	newFilterCondition := []utils.FilterCondition{}

	for _, filter := range filterCondition {
		if strings.Contains(filter.ColumnField, "supplier_code") {
			supplierCode = filter.ColumnValue
			continue
		}
		if strings.Contains(filter.ColumnField, "supplier_name") {
			supplierName = filter.ColumnValue
			continue
		}
		newFilterCondition = append(newFilterCondition, filter)
	}

	// Membuat join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct).
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_supplier ON dms_microservices_general_dev.dbo.mtr_supplier.supplier_id = mtr_item.supplier_id").
		Joins("LEFT JOIN mtr_item_type AS mtr_item_type_alias ON mtr_item_type_alias.item_type_id = mtr_item.item_type_id")

	// Terapkan filter
	whereQuery := utils.ApplyFilter(joinTable, newFilterCondition)

	// Handle item_id filter
	if len(itemIDs) > 0 && itemIDs[0] != "" {
		whereQuery = whereQuery.Where("mtr_item.item_id IN (?)", itemIDs)
	}

	var supplierIds []int
	if supplierCode != "" || supplierName != "" {
		supplierName = strings.ReplaceAll(supplierName, " ", "%20")
		supplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier?page=0&limit=1000000&supplier_code=" + supplierCode + "&supplier_name=" + supplierName
		var supplierResponse []masteritempayloads.PurchasePriceSupplierResponse
		if err := utils.GetArray(supplierUrl, &supplierResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		for _, supplier := range supplierResponse {
			supplierIds = append(supplierIds, supplier.SupplierId)
		}

		if len(supplierIds) == 0 {
			supplierIds = []int{-1}
		}

		whereQuery = whereQuery.Where("mtr_item.supplier_id IN ?", supplierIds)
	}

	var responses []masteritempayloads.ItemSearch
	err := whereQuery.Scopes(pagination.Paginate(&tableStruct, &pages, whereQuery)).Scan(&responses).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch data from database",
			Err:        errors.New("failed to fetch data from database"),
		}
	}

	if len(responses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "no data found",
			Err:        errors.New("no data found"),
		}
	}

	var mapResponses []map[string]interface{}
	for _, response := range responses {
		// Panggil API eksternal untuk mengambil Supplier data
		SupplierURL := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(response.SupplierId)
		var getSupplierResponse masteritempayloads.SupplierMasterResponse
		if err := utils.Get(SupplierURL, &getSupplierResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		itemGroupUrl := config.EnvConfigs.GeneralServiceUrl + "item-group/" + strconv.Itoa(response.ItemGroupId)
		getItemGroupResponses := masteritempayloads.ItemGroupResponse{}
		errUrlItemPackage := utils.Get(itemGroupUrl, &getItemGroupResponses, nil)
		if errUrlItemPackage != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Build response map dengan data dari supplier
		responseMap := map[string]interface{}{
			"is_active":       response.IsActive,
			"item_id":         response.ItemId,
			"item_code":       response.ItemCode,
			"item_name":       response.ItemName,
			"item_group_id":   response.ItemGroupId,
			"item_class_id":   response.ItemClassId,
			"item_type_id":    response.ItemTypeId,
			"item_type":       response.ItemTypeCode,
			"supplier_id":     response.SupplierId,
			"item_class_code": response.ItemClassCode,
			"item_group_code": getItemGroupResponses.ItemGroupCode,
			"supplier_name":   getSupplierResponse.SupplierName,
			"supplier_code":   getSupplierResponse.SupplierCode,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	return mapResponses, pages.TotalPages, int(pages.TotalRows), nil
}

func (r *ItemRepositoryImpl) GetAllItem(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	var responses []masteritempayloads.ItemLookup

	tableStruct := masteritempayloads.ItemLookup{}
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	whereQuery := utils.ApplyFilterForDB(joinTable, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&tableStruct, &pages, whereQuery)).Scan(&responses).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to fetch data from database: %w", err),
		}
	}

	if len(responses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("no data found"),
		}
	}

	var mapResponses []map[string]interface{}

	for _, response := range responses {
		responseMap := map[string]interface{}{
			"is_active":       response.IsActive,
			"item_id":         response.ItemId,
			"item_code":       response.ItemCode,
			"item_name":       response.ItemName,
			"item_group_id":   response.ItemGroupId,
			"item_class_id":   response.ItemClassId,
			"item_type_id":    response.ItemTypeId,
			"supplier_id":     response.SupplierId,
			"item_class_name": response.ItemClassName,
			"item_level_1":    response.ItemLevel_1,
			"item_level_2":    response.ItemLevel_2,
			"item_level_3":    response.ItemLevel_3,
			"item_level_4":    response.ItemLevel_4,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	return mapResponses, pages.TotalPages, int(pages.TotalRows), nil
}

func (r *ItemRepositoryImpl) GetAllItemLookup(tx *gorm.DB, filter []utils.FilterCondition) (any, *exceptions.BaseErrorResponse) {

	panic("unimplemented")
}

func (r *ItemRepositoryImpl) GetItemById(tx *gorm.DB, Id int) (masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Item{}
	response := masteritempayloads.ItemResponse{}

	// Fetch the item entity from the database
	err := tx.Model(&entities).
		Select(`
			mtr_item.*,
			mil1.item_level_1_code,
			mil1.item_level_1_name,
			mil2.item_level_2_code,
			mil2.item_level_2_name,
			mil3.item_level_3_code,
			mil3.item_level_3_name,
			mil4.item_level_4_code,
			mil4.item_level_4_name
			`).
		Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = mtr_item.item_level_1_id").
		Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = mtr_item.item_level_2_id").
		Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = mtr_item.item_level_3_id").
		Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = mtr_item.item_level_4_id").
		Where(masteritementities.Item{ItemId: Id}).
		First(&response).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "item not found",
				Err:        err,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch item data",
			Err:        err,
		}
	}

	// Call external service to get Supplier details
	supplierResponse := masteritempayloads.SupplierMasterResponse{}
	supplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(response.SupplierId)
	if err := utils.Get(supplierUrl, &supplierResponse, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch supplier data",
			Err:        err,
		}
	}

	// Populate supplier data into response
	response.SupplierCode = &supplierResponse.SupplierCode
	response.SupplierName = &supplierResponse.SupplierName

	// Return the response with a populated supplier
	return response, nil
}

func (r *ItemRepositoryImpl) GetItemWithMultiId(tx *gorm.DB, MultiIds []string) ([]masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse) {
	var response []masteritempayloads.ItemResponse
	entities := masteritementities.Item{}
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	tx.Logger = newLogger

	rows, err := tx.Model(&entities).
		Where("item_id in ?", MultiIds).
		Scan(&response).
		Rows()

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *ItemRepositoryImpl) GetItemCode(tx *gorm.DB, code string) (masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Item{}
	response := masteritempayloads.ItemResponse{}

	rows, err := tx.Model(&entities).Select("mtr_item.*,u.*").
		Where(masteritementities.Item{
			ItemCode: code,
		}).InnerJoins("Join mtr_uom_item u ON mtr_item.item_id = u.item_id").
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	supplierResponse := masteritempayloads.SupplierMasterResponse{}

	supplierUrl := config.EnvConfigs.GeneralServiceUrl + "/supplier-master/" + strconv.Itoa(response.SupplierId)

	if err := utils.Get(supplierUrl, &supplierResponse, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	response.SupplierCode = &supplierResponse.SupplierCode
	response.SupplierName = &supplierResponse.SupplierName

	// joinSupplierData := utils.DataFrameInnerJoin([]masteritempayloads.ItemResponse{response}, []masteritempayloads.SupplierMasterResponse{supplierResponse}, "SupplierId")

	// IMPLEMENT PERSON IN CHARGE AFTER INTEGRATION TOKEN AUTHORIZE TO USER SERVICE!!

	defer rows.Close()

	return response, nil

}

func (r *ItemRepositoryImpl) SaveItem(tx *gorm.DB, req masteritempayloads.ItemRequest) (masteritempayloads.ItemSaveResponse, *exceptions.BaseErrorResponse) {
	response := masteritempayloads.ItemSaveResponse{}
	entities := masteritementities.Item{
		IsActive:                     req.IsActive,
		ItemId:                       req.ItemId,
		ItemCode:                     req.ItemCode,
		ItemClassId:                  req.ItemClassId,
		ItemName:                     req.ItemName,
		ItemGroupId:                  req.ItemGroupId,
		ItemTypeId:                   req.ItemTypeId,
		ItemLevel1Str:                req.ItemLevel1,
		ItemLevel2Str:                req.ItemLevel2,
		ItemLevel3Str:                req.ItemLevel3,
		ItemLevel4Str:                req.ItemLevel4,
		SupplierId:                   req.SupplierId,
		UnitOfMeasurementTypeId:      req.UnitOfMeasurementTypeId,
		UnitOfMeasurementSellingId:   req.UnitOfMeasurementSellingId,
		UnitOfMeasurementPurchaseId:  req.UnitOfMeasurementPurchaseId,
		UnitOfMeasurementStockId:     req.UnitOfMeasurementStockId,
		SalesItem:                    req.SalesItem,
		Lottable:                     req.Lottable,
		Inspection:                   req.Inspection,
		PriceListItem:                req.PriceListItem,
		StockKeeping:                 r.DetermineStockKeeping(req, req.StockKeeping),
		DiscountId:                   req.DiscountId,
		MarkupMasterId:               req.MarkupMasterId,
		DimensionOfLength:            req.DimensionOfLength,
		DimensionOfWidth:             req.DimensionOfWidth,
		DimensionOfHeight:            req.DimensionOfHeight,
		DimensionUnitOfMeasurementId: req.DimensionUnitOfMeasurementId,
		Weight:                       req.Weight,
		UnitOfMeasurementWeight:      req.UnitOfMeasurementWeight,
		StorageTypeId:                req.StorageTypeId,
		Remark:                       req.Remark,
		AtpmWarrantyClaimTypeId:      req.AtpmWarrantyClaimTypeId,
		LastPrice:                    req.LastPrice,
		UseDiscDecentralize:          req.UseDiscDecentralize,
		CommonPricelist:              req.CommonPricelist,
		IsRemovable:                  req.IsRemovable,
		IsMaterialPlus:               req.IsMaterialPlus,
		SpecialMovementId:            req.SpecialMovementId,
		IsItemRegulation:             req.IsItemRegulation,
		IsTechnicalDefect:            req.IsTechnicalDefect,
		IsMandatory:                  req.IsMandatory,
		MinimumOrderQty:              req.MinimumOrderQty,
		HarmonizedNo:                 req.HarmonizedNo,
		AtpmSupplierId:               req.AtpmSupplierId,
		AtpmVendorSuppliability:      req.AtpmVendorSuppliability,
		PmsItem:                      req.PmsItem,
		Regulation:                   req.Regulation,
		AutoPickWms:                  req.AutoPickWms,
		GmmCatalogCode:               req.GmmCatalogCode,
		PrincipalBrandParentId:       req.PrincipalBrandParentId,
		ProportionalSupplyWms:        req.ProportionalSupplyWms,
		Remark2:                      req.Remark2,
		Remark3:                      req.Remark3,
		SourceTypeId:                 req.SourceTypeId,
		AtpmSupplierCodeOrderId:      req.AtpmSupplierCodeOrderId,
		PersonInChargeId:             req.PersonInChargeId,
		IsSellable:                   req.IsSellable,
		IsAffiliatedTrx:              req.IsAffiliatedTrx,
	}

	err := tx.Save(&entities).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to save item",
			Err:        err,
		}
	}

	model := masteritementities.Item{}
	err = tx.Model(&model).Where(masteritementities.Item{ItemCode: req.ItemCode}).First(&model).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch item data",
			Err:        err,
		}
	}

	atpmResponse := masteritempayloads.AtpmOrderTypeResponse{}
	atpmOrderTypeUrl := config.EnvConfigs.GeneralServiceUrl + "/atpm-order-type/" + strconv.Itoa(req.SourceTypeId)
	if err := utils.Get(atpmOrderTypeUrl, &atpmResponse, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch atpm order type data",
			Err:        err,
		}
	}

	uomTypeModel := masteritementities.UomType{}
	err = tx.Model(&uomTypeModel).Where(masteritementities.UomType{UomTypeId: req.UnitOfMeasurementTypeId}).First(&uomTypeModel).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch uom type data",
			Err:        err,
		}
	}

	uomItemEntities := masteritementities.UomItem{
		IsActive:          req.IsActive,
		ItemId:            model.ItemId,
		UomSourceTypeCode: atpmResponse.AtpmOrderTypeCode,
		UomTypeCode:       uomTypeModel.UomTypeCode,
		SourceUomId:       req.UnitOfMeasurementPurchaseId,
		TargetUomId:       req.UnitOfMeasurementStockId,
		SourceConvertion:  float64(req.SourceConvertion),
		TargetConvertion:  float64(req.TargetConvertion),
	}

	err = tx.Save(&uomItemEntities).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to save uom item",
			Err:        err,
		}
	}

	result := masteritempayloads.ItemSaveResponse{
		IsActive:   entities.IsActive,
		ItemId:     entities.ItemId,
		ItemName:   entities.ItemName,
		ItemCode:   entities.ItemCode,
		ItemTypeId: entities.ItemTypeId,
		ItemLevel1: entities.ItemLevel1Str,
		ItemLevel2: entities.ItemLevel2Str,
		ItemLevel3: entities.ItemLevel3Str,
		ItemLevel4: entities.ItemLevel4Str,
	}

	return result, nil
}

func (r *ItemRepositoryImpl) DetermineStockKeeping(req masteritempayloads.ItemRequest, manualStockKeeping bool) bool {
	itemGroupID := req.ItemGroupId
	itemClassID := req.ItemClassId
	itemTypeID := req.ItemTypeId

	switch itemGroupID {
	case 1: // Fixed Asset
		return false // Non-stock keeping for Fixed Asset
	case 6: // Outside Job
		if itemTypeID == 2 { // Check for Service (ID 2)
			return false // Non-stock keeping for Services in Outside Job group
		}
	case 7: // Prepaid
		return false // Non-stock keeping for Prepaid group
	case 4, 5: // OPEX, Opex Promosi
		return false // Non-stock keeping for OPEX-related groups
	case 2: // Inventory
		switch itemClassID {
		case 73: // Fee
			if itemTypeID == 2 { // Check for Service (ID 2)
				return false // Non-stock keeping for Services in Fee class
			}
		case 75, 76, 71, 70, 77, 69: // Consumable Material, Equipment, Material, Oil, Souvenir, Sparepart
			return true // Stock keeping for these item classes
		case 74: // Accessories
			if itemTypeID == 1 { // Check for Goods (ID 1)
				return true // Stock keeping for Goods in Accessories
			} else if itemTypeID == 2 { // Check for Service (ID 2)
				return false // Non-stock keeping for Services in Accessories
			}
		}
	}

	return manualStockKeeping
}

func (r *ItemRepositoryImpl) ChangeStatusItem(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.Item

	result := tx.Model(&entities).
		Where("item_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	if entities.IsActive {
		entities.IsActive = false
	} else {
		entities.IsActive = true
	}

	result = tx.Save(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}

func (r *ItemRepositoryImpl) SaveItemDetail(tx *gorm.DB, request masteritempayloads.ItemDetailResponse) (bool, error) {
	entities := masteritementities.ItemDetail{
		IsActive:     request.IsActive,
		ItemDetailId: request.ItemDetailId,
		ItemId:       request.ItemId,
		BrandId:      request.BrandId,
		ModelId:      request.ModelId,
		VariantId:    request.VariantId,
		MileageEvery: request.MileageEvery,
		ReturnEvery:  request.ReturnEvery,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *ItemRepositoryImpl) GetAllItemDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	var responses []masteritempayloads.ItemDetailRequest

	tableStruct := masteritempayloads.ItemDetailRequest{}
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	whereQuery := utils.ApplyFilterExact(joinTable, filterCondition)

	rows, err := whereQuery.Find(&responses).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch data from database",
			Err:        err,
		}
	}
	defer rows.Close()

	var convertedResponses []masteritempayloads.ItemDetailResponse

	for rows.Next() {
		var (
			itemDetailReq masteritempayloads.ItemDetailRequest
			itemDetailRes masteritempayloads.ItemDetailResponse
		)

		if err := rows.Scan(
			&itemDetailReq.ItemDetailId,
			&itemDetailReq.ItemId,
			&itemDetailReq.BrandId,
			&itemDetailReq.ModelId,
			&itemDetailReq.VariantId,
			&itemDetailReq.MileageEvery,
			&itemDetailReq.ReturnEvery,
			&itemDetailReq.IsActive); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to scan item detail data",
				Err:        err,
			}
		}

		// Fetch Brand data
		BrandURL := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(itemDetailReq.BrandId)
		var getBrandResponse masterpayloads.BrandResponse
		if err := utils.Get(BrandURL, &getBrandResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to fetch brand data",
				Err:        err,
			}
		}

		// Fetch Model data
		ModelURL := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(itemDetailReq.ModelId)
		var getModelResponse masterpayloads.UnitModelResponse
		if err := utils.Get(ModelURL, &getModelResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to fetch model data",
				Err:        err,
			}
		}

		// Fetch Variant data
		VariantURL := config.EnvConfigs.SalesServiceUrl + "unit-variant/" + strconv.Itoa(itemDetailReq.VariantId)
		var getVariantResponse masterpayloads.GetVariantResponse
		if err := utils.Get(VariantURL, &getVariantResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to fetch variant data",
				Err:        err,
			}
		}

		itemDetailRes = masteritempayloads.ItemDetailResponse{
			ItemDetailId:       itemDetailReq.ItemDetailId,
			ItemId:             itemDetailReq.ItemId,
			BrandId:            itemDetailReq.BrandId,
			BrandName:          getBrandResponse.BrandName,
			ModelId:            itemDetailReq.ModelId,
			ModelCode:          getModelResponse.ModelCode,
			ModelDescription:   getModelResponse.ModelDescription,
			VariantId:          itemDetailReq.VariantId,
			VariantCode:        getVariantResponse.VariantCode,
			VariantDescription: getVariantResponse.VariantDescription,
			ReturnEvery:        itemDetailReq.ReturnEvery,
			MileageEvery:       itemDetailReq.MileageEvery,
			IsActive:           itemDetailReq.IsActive,
		}

		convertedResponses = append(convertedResponses, itemDetailRes)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error in item detail rows iteration",
			Err:        err,
		}
	}

	var mapResponses []map[string]interface{}
	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"item_detail_id":      response.ItemDetailId,
			"item_id":             response.ItemId,
			"brand_id":            response.BrandId,
			"brand_name":          response.BrandName,
			"model_id":            response.ModelId,
			"model_code":          response.ModelCode,
			"model_description":   response.ModelDescription,
			"variant_id":          response.VariantId,
			"variant_code":        response.VariantCode,
			"variant_description": response.VariantDescription,
			"mileage_every":       response.MileageEvery,
			"return_every":        response.ReturnEvery,
			"is_active":           response.IsActive,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *ItemRepositoryImpl) GetItemDetailById(tx *gorm.DB, ItemId, ItemDetailId int) (masteritempayloads.ItemDetailRequest, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemDetail{}
	response := masteritempayloads.ItemDetailRequest{}

	err := tx.Model(&entities).
		Where(masteritementities.ItemDetail{
			ItemDetailId: ItemDetailId,
			ItemId:       ItemId,
		}).
		First(&entities).
		Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	response.ItemDetailId = entities.ItemDetailId
	response.ItemId = entities.ItemId
	response.BrandId = entities.BrandId
	response.ModelId = entities.ModelId
	response.VariantId = entities.VariantId
	response.MileageEvery = entities.MileageEvery
	response.ReturnEvery = entities.ReturnEvery
	response.IsActive = entities.IsActive

	return response, nil
}

func (r *ItemRepositoryImpl) AddItemDetail(tx *gorm.DB, ItemId int, req masteritempayloads.ItemDetailRequest) (masteritementities.ItemDetail, *exceptions.BaseErrorResponse) {
	var entities masteritementities.ItemDetail

	// Cek apakah detail item sudah ada untuk kombinasi ItemId, BrandId, ModelId, VariantId
	err := tx.Where("item_id = ? AND brand_id = ? AND model_id = ? AND variant_id = ?", ItemId, req.BrandId, req.ModelId, req.VariantId).First(&entities).Error
	if err == nil {
		return masteritementities.ItemDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Item detail already exists",
			Err:        errors.New("item detail already exists"),
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return masteritementities.ItemDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error checking existing item detail",
			Err:        err,
		}
	}

	// Jika IsActive adalah false, set MileageEvery dan ReturnEvery ke 0
	if !req.IsActive {
		req.MileageEvery = 0
		req.ReturnEvery = 0
	}

	entities = masteritementities.ItemDetail{
		ItemId:       ItemId,
		BrandId:      req.BrandId,
		ModelId:      req.ModelId,
		VariantId:    req.VariantId,
		MileageEvery: req.MileageEvery,
		ReturnEvery:  req.ReturnEvery,
		IsActive:     req.IsActive,
	}

	err = tx.Save(&entities).Error
	if err != nil {
		return masteritementities.ItemDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save item detail",
			Err:        err,
		}
	}

	return entities, nil
}

func (r *ItemRepositoryImpl) DeleteItemDetails(tx *gorm.DB, ItemId int, itemDetailIDs []int) (masteritempayloads.DeleteItemResponse, *exceptions.BaseErrorResponse) {
	var entities masteritementities.ItemDetail

	result := tx.Model(&entities).
		Where("item_id = ? AND item_detail_id IN (?)", ItemId, itemDetailIDs).
		Delete(&entities)

	if result.Error != nil {
		return masteritempayloads.DeleteItemResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete item details",
			Err:        result.Error,
		}
	}

	return masteritempayloads.DeleteItemResponse{}, nil
}

func (r *ItemRepositoryImpl) UpdateItem(tx *gorm.DB, ItemId int, req masteritempayloads.ItemUpdateRequest) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.Item

	result := tx.Model(&entities).Where("item_id = ?", ItemId).First(&entities).Updates(req)
	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        result.Error,
		}
	}

	if req.SourceConvertion != 0 || req.TargetConvertion != 0 {
		uomItemModel := masteritementities.UomItem{}

		uomItemEntities := masteritementities.UomItem{
			SourceConvertion: float64(req.SourceConvertion),
			TargetConvertion: float64(req.TargetConvertion),
		}

		err := tx.Model(&uomItemModel).Where(masteritementities.UomItem{ItemId: entities.ItemId}).Updates(&uomItemEntities).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

	}
	return true, nil
}

func (r *ItemRepositoryImpl) UpdateItemDetail(tx *gorm.DB, Id int, itemDetailId int, req masteritempayloads.ItemDetailUpdateRequest) (masteritementities.ItemDetail, *exceptions.BaseErrorResponse) {
	var entities masteritementities.ItemDetail

	// Fetch the existing record to update
	err := tx.Where("item_detail_id = ? AND item_id = ?", itemDetailId, Id).First(&entities).Error
	if err != nil {
		return masteritementities.ItemDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Item detail not found",
			Err:        err,
		}
	}

	result := tx.Model(&entities).Where("item_detail_id = ? AND item_id = ?", itemDetailId, Id).Updates(req)
	if result.Error != nil {
		return masteritementities.ItemDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Failed to update item detail",
			Err:        result.Error,
		}
	}

	return entities, nil
}

func (r *ItemRepositoryImpl) GetPrincipleBrandDropdown(tx *gorm.DB) ([]masteritempayloads.PrincipleBrandDropdownResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PrincipleBrandParent{}
	payloads := []masteritempayloads.PrincipleBrandDropdownResponse{}
	err := tx.Model(&entities).Scan(&payloads).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	return payloads, nil
}

func (r *ItemRepositoryImpl) GetCatalogCode(tx *gorm.DB) ([]masteritempayloads.GetCatalogCode, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PrincipleBrandParent{}
	payloads := []masteritempayloads.GetCatalogCode{}

	err := tx.Model(&entities).Scan(&payloads).Error
	if err != nil {
		return payloads, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	return payloads, nil
}

func (r *ItemRepositoryImpl) GetPrincipleBrandParent(tx *gorm.DB, code string) ([]masteritempayloads.PrincipleBrandDropdownDescription, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PrincipleBrandParent{}
	payloads := []masteritempayloads.PrincipleBrandDropdownDescription{}
	err := tx.Model(&entities).Where(masteritementities.PrincipleBrandParent{
		CatalogueCode: code,
	}).Scan(&payloads).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	return payloads, nil
}

func (r *ItemRepositoryImpl) AddItemDetailByBrand(tx *gorm.DB, id string, itemId int) ([]masteritempayloads.ItemDetailResponse, *exceptions.BaseErrorResponse) {
	var itemDetails []masteritempayloads.ItemDetailResponse
	brandid := strings.Split(id, ",")

	for _, id := range brandid {
		var getdatabybrand []masteritempayloads.BrandModelVariantResponse
		err := utils.Get(config.EnvConfigs.SalesServiceUrl+"unit-variant-by-brand/"+id, &getdatabybrand, nil)
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        errors.New("brand has no variant and model"),
			}
		}

		for _, detail := range getdatabybrand {
			entities := masteritementities.ItemDetail{
				IsActive:  true,
				ItemId:    itemId,
				BrandId:   detail.BrandId,
				ModelId:   detail.ModelId,
				VariantId: detail.VariantId,
			}

			err = tx.Save(&entities).Error
			if err != nil {
				return nil, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusConflict,
					Err:        err,
				}
			}

			itemDetails = append(itemDetails, masteritempayloads.ItemDetailResponse{
				ItemDetailId: entities.ItemDetailId,
				IsActive:     entities.IsActive,
				ItemId:       itemId,
				BrandId:      detail.BrandId,
				ModelId:      detail.ModelId,
				VariantId:    detail.VariantId,
			})
		}
	}

	return itemDetails, nil
}
