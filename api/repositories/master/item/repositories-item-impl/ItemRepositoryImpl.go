package masteritemrepositoryimpl

import (
	config "after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ItemRepositoryImpl struct {
}

func StartItemRepositoryImpl() masteritemrepository.ItemRepository {
	return &ItemRepositoryImpl{}
}

func (r *ItemRepositoryImpl) GetAllItem(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, error) {
	// Define variables
	var (
		responses    []masteritempayloads.ItemResponse
		tableStruct  = masteritempayloads.ItemLookup{}
		baseQuery    = tx.Model(&masteritempayloads.ItemLookup{})
		totalRows    int64
		mapResponses []map[string]interface{}
		totalPages   int
	)

	// Apply joins and filters
	joinTable := utils.CreateJoinSelectStatement(baseQuery, tableStruct)
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	// Get total number of rows for pagination
	if err := whereQuery.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, err
	}

	// Apply pagination
	limit := pages.GetLimit()
	offset := pages.GetOffset()
	whereQuery = whereQuery.Offset(offset).Limit(limit)

	// Execute query to fetch responses
	if err := whereQuery.Find(&responses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, gorm.ErrRecordNotFound
		}
		return nil, 0, 0, err
	}

	// Calculate total pages
	totalPages = int(math.Ceil(float64(totalRows) / float64(limit)))

	// Convert paginated data to map format
	for _, data := range responses {
		responseMap := map[string]interface{}{
			"is_active":                        data.IsActive,
			"item_id":                          data.ItemId,
			"item_code":                        data.ItemCode,
			"item_name":                        data.ItemName,
			"item_group_id":                    data.ItemGroupId,
			"item_class_id":                    data.ItemClassId,
			"item_type":                        data.ItemType,
			"item_level_1":                     data.ItemLevel_1,
			"item_level_2":                     data.ItemLevel_2,
			"item_level_3":                     data.ItemLevel_3,
			"item_level_4":                     data.ItemLevel_4,
			"supplier_id":                      data.SupplierId,
			"unit_of_measurement_type_id":      data.UnitOfMeasurementTypeId,
			"unit_of_measurement_selling_id":   data.UnitOfMeasurementSellingId,
			"unit_of_measurement_purchase_id":  data.UnitOfMeasurementPurchaseId,
			"unit_of_measurement_stock_id":     data.UnitOfMeasurementStockId,
			"sales_item":                       data.SalesItem,
			"lottable":                         data.Lottable,
			"inspection":                       data.Inspection,
			"price_list_item":                  data.PriceListItem,
			"stock_keeping":                    data.StockKeeping,
			"discount_id":                      data.DiscountId,
			"markup_master_id":                 data.MarkupMasterId,
			"dimension_of_length":              data.DimensionOfLength,
			"dimension_of_width":               data.DimensionOfWidth,
			"dimension_of_height":              data.DimensionOfHeight,
			"dimension_unit_of_measurement_id": data.DimensionUnitOfMeasurementId,
			"weight":                           data.Weight,
			"unit_of_measurement_weight":       data.UnitOfMeasurementWeight,
			"storage_type_id":                  data.StorageTypeId,
			"remark":                           data.Remark,
			"atpm_warranty_claim_type_id":      data.AtpmWarrantyClaimTypeId,
			"last_price":                       data.LastPrice,
			"use_disc_decentralize":            data.UseDiscDecentralize,
			"common_pricelist":                 data.CommonPricelist,
			"is_removable":                     data.IsRemovable,
			"is_material_plus":                 data.IsMaterialPlus,
			"special_movement_id":              data.SpecialMovementId,
			"is_item_regulation":               data.IsItemRegulation,
			"is_technical_defect":              data.IsTechnicalDefect,
			"is_mandatory":                     data.IsMandatory,
			"minimum_order_qty":                data.MinimumOrderQty,
			"harmonized_no":                    data.HarmonizedNo,
			"atpm_supplier_id":                 data.AtpmSupplierId,
			"atpm_vendor_suppliability":        data.AtpmVendorSuppliability,
			"pms_item":                         data.PmsItem,
			"regulation":                       data.Regulation,
			"auto_pick_wms":                    data.AutoPickWms,
			"gmm_catalog_code":                 data.GmmCatalogCode,
			"principal_brand_parent_id":        data.PrincipalBrandParentId,
			"proportional_supply_wms":          data.ProportionalSupplyWms,
			"remark2":                          data.Remark2,
			"remark3":                          data.Remark3,
			"source_type_id":                   data.SourceTypeId,
			"atpm_supplier_code_order_id":      data.AtpmSupplierCodeOrderId,
			"person_in_charge_id":              data.PersonInChargeId,

			// Add other fields as needed
		}
		mapResponses = append(mapResponses, responseMap)
	}

	return mapResponses, totalPages, int(totalRows), nil
}

func (r *ItemRepositoryImpl) GetAllItemLookup(tx *gorm.DB, queryParams []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, error) {
	// Define variables
	var (
		tableStruct = masteritempayloads.ItemLookup{}
		responses   []masteritempayloads.ItemLookup
		totalRows   int64
		totalPages  int
	)

	// Apply joins and filters
	baseQuery := tx.Model(&masteritempayloads.ItemLookup{})
	joinTable := utils.CreateJoinSelectStatement(baseQuery, tableStruct)
	whereQuery := utils.ApplyFilter(joinTable, queryParams)

	// Execute query to fetch responses with pagination
	offset := pages.GetOffset()
	limit := pages.GetLimit()
	if err := whereQuery.Offset(offset).Limit(limit).Find(&responses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, gorm.ErrRecordNotFound
		}
		return nil, 0, 0, err
	}

	// Convert paginated data to map format
	var mapResponses []map[string]interface{}
	for _, data := range responses {
		// Fetch data from mtr_item_group for each response
		itemGroupUrl := config.EnvConfigs.GeneralServiceUrl + "item-group/" + strconv.Itoa(data.ItemGroupId)
		var itemGroupResp masteritempayloads.ItemGroupResponse
		err := utils.Get(itemGroupUrl, &itemGroupResp, nil)
		if err != nil {
			return nil, 0, 0, err
		}

		// Create response map combining ItemLookup and ItemGroupResponse data
		responseMap := map[string]interface{}{
			"is_active":       data.IsActive,
			"item_id":         data.ItemId,
			"item_code":       data.ItemCode,
			"item_name":       data.ItemName,
			"item_group_id":   data.ItemGroupId,
			"item_class_id":   data.ItemClassId,
			"item_type":       data.ItemType,
			"supplier_id":     data.SupplierId,
			"item_group_name": itemGroupResp.ItemGroupName,
			// Add more fields as needed
		}
		mapResponses = append(mapResponses, responseMap)
	}

	// Calculate total pages
	totalRows = pages.TotalRows
	totalPages = pages.TotalPages

	return mapResponses, totalPages, int(totalRows), nil
}

func (r *ItemRepositoryImpl) GetItemById(tx *gorm.DB, Id int) (masteritempayloads.ItemResponse, error) {
	entities := masteritementities.Item{}
	response := masteritempayloads.ItemResponse{}

	rows, err := tx.Model(&entities).
		Where(masteritementities.Item{
			ItemId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *ItemRepositoryImpl) GetItemWithMultiId(tx *gorm.DB, MultiIds []string) ([]masteritempayloads.ItemResponse, error) {
	entities := masteritementities.Item{}
	var response []masteritempayloads.ItemResponse

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
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *ItemRepositoryImpl) GetItemCode(tx *gorm.DB, code string) ([]map[string]interface{}, error) {
	encodedCode := url.PathEscape(code)

	entities := masteritementities.Item{}
	response := masteritempayloads.ItemResponse{}
	var getSupplierMasterResponse masteritempayloads.SupplierMasterResponse
	var getItemGroupResponse masteritempayloads.ItemGroupResponse
	var getStorageTypeResponse masteritempayloads.StorageTypeResponse
	var getSpecialMovementResponse masteritempayloads.SpecialMovementResponse
	var getAtpmSupplierResponse masteritempayloads.AtpmSupplierResponse
	var getAtpmSupplierCodeOrderResponse masteritempayloads.AtpmSupplierCodeOrderResponse
	var getPersonInChargeResponse masteritempayloads.PersonInChargeResponse
	var getAtpmWarrantyClaimTypeResponse masteritempayloads.AtpmWarrantyClaimTypeResponse

	rows, err := tx.Model(&entities).
		Where(masteritementities.Item{
			ItemCode: encodedCode, // Menggunakan kode yang telah diencode
		}).First(&response).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//FK Luar with mtr_item_group common-general service
	ItemGroupUrl := config.EnvConfigs.GeneralServiceUrl + "item-group/" + strconv.Itoa(response.ItemGroupId)
	errUrlItemGroup := utils.Get(ItemGroupUrl, &getItemGroupResponse, nil)
	fmt.Println("Fetching mtr_item_group data from:", ItemGroupUrl)
	if errUrlItemGroup != nil {
		return nil, err
	}

	firstJoin := utils.DataFrameLeftJoin([]masteritempayloads.ItemResponse{response}, []masteritempayloads.ItemGroupResponse{getItemGroupResponse}, "ItemGroupId")

	//FK luar with mtr_supplier general service
	SupplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier-master/" + strconv.Itoa(response.SupplierId)
	errUrlSupplierMaster := utils.Get(SupplierUrl, &getSupplierMasterResponse, nil)
	fmt.Println("Fetching mtr_supplier data from:", SupplierUrl)
	if errUrlSupplierMaster != nil {
		return nil, err
	}

	secondJoin := utils.DataFrameLeftJoin(firstJoin, []masteritempayloads.SupplierMasterResponse{getSupplierMasterResponse}, "SupplierId")
	//FK luar with storage_type general service
	StorageTypeUrl := config.EnvConfigs.GeneralServiceUrl + "storage-type/" + strconv.Itoa(response.StorageTypeId)
	errUrlStorageType := utils.Get(StorageTypeUrl, &getStorageTypeResponse, nil)
	fmt.Println("Fetching storage_type data from:", StorageTypeUrl)
	if errUrlStorageType != nil {
		return nil, err
	}

	thirdJoin := utils.DataFrameLeftJoin(secondJoin, []masteritempayloads.StorageTypeResponse{getStorageTypeResponse}, "StorageTypeId")
	//FK luar with mtr_warranty_claim_type common service
	WarrantyClaimTypeUrl := config.EnvConfigs.GeneralServiceUrl + "warranty-claim-type/" + strconv.Itoa(response.AtpmWarrantyClaimTypeId)
	errUrlWarrantyClaimType := utils.Get(WarrantyClaimTypeUrl, &getAtpmWarrantyClaimTypeResponse, nil)
	fmt.Println("Fetching mtr_warranty_claim_type data from:", WarrantyClaimTypeUrl)
	if errUrlWarrantyClaimType != nil {
		return thirdJoin, err
	}

	fourthJoin := utils.DataFrameLeftJoin(thirdJoin, []masteritempayloads.AtpmWarrantyClaimTypeResponse{getAtpmWarrantyClaimTypeResponse}, "AtpmWarrantyClaimTypeId")
	//FK luar with mtr_special_movement common service
	SpecialMovementUrl := config.EnvConfigs.GeneralServiceUrl + "special-movement/" + strconv.Itoa(response.SpecialMovementId)
	errUrlSpecialMovement := utils.Get(SpecialMovementUrl, &getSpecialMovementResponse, nil)
	fmt.Println("Fetching mtr_special_movement data from:", SpecialMovementUrl)
	if errUrlSpecialMovement != nil {
		return fourthJoin, err
	}

	fifthJoin := utils.DataFrameLeftJoin(fourthJoin, []masteritempayloads.SpecialMovementResponse{getSpecialMovementResponse}, "SpecialMovementId")
	//FK luar with mtr_supplier general service atpm_supplier_id
	AtpmSupplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier-master/" + strconv.Itoa(response.AtpmSupplierId)
	errUrlAtpmSupplier := utils.Get(AtpmSupplierUrl, &getAtpmSupplierResponse, nil)
	fmt.Println("Fetching mtr_supplier data from:", AtpmSupplierUrl)
	if errUrlAtpmSupplier != nil {
		return fifthJoin, err
	}

	sixthJoin := utils.DataFrameLeftJoin(fifthJoin, []masteritempayloads.AtpmSupplierResponse{getAtpmSupplierResponse}, "AtpmSupplierId")
	//FK luar with mtr_supplier general service atpm_supplier_code_order_id
	AtpmSupplierCodeOrderUrl := config.EnvConfigs.GeneralServiceUrl + "supplier-master/" + strconv.Itoa(response.AtpmSupplierCodeOrderId)
	errUrlAtpmSupplierCodeOrder := utils.Get(AtpmSupplierCodeOrderUrl, &getAtpmSupplierCodeOrderResponse, nil)
	fmt.Println("Fetching mtr_supplier data from:", AtpmSupplierCodeOrderUrl)
	if errUrlAtpmSupplierCodeOrder != nil {
		return sixthJoin, err
	}

	seventhJoin := utils.DataFrameLeftJoin(sixthJoin, []masteritempayloads.AtpmSupplierCodeOrderResponse{getAtpmSupplierCodeOrderResponse}, "AtpmSupplierCodeOrderId")
	//FK luar with mtr_user_details general service
	PersonInChargeUrl := config.EnvConfigs.GeneralServiceUrl + "user-details-all/" + strconv.Itoa(response.PersonInChargeId)
	errUrlPersonInCharge := utils.Get(PersonInChargeUrl, &getPersonInChargeResponse, nil)
	fmt.Println("Fetching mtr_user_details data from:", PersonInChargeUrl)
	if errUrlPersonInCharge != nil {
		return seventhJoin, err
	}

	eightJoin := utils.DataFrameLeftJoin(seventhJoin, getPersonInChargeResponse, "PersonInChargeId")

	// FK luar with mtr_unit_of_measurement_type
	// fk luar with mtr_atpm_order_type common service

	return eightJoin, nil
}

func (r *ItemRepositoryImpl) SaveItem(tx *gorm.DB, req masteritempayloads.ItemResponse) (bool, error) {
	entities := masteritementities.Item{
		ItemCode:                     req.ItemCode,
		ItemClassId:                  req.ItemClassId,
		ItemName:                     req.ItemName,
		ItemGroupId:                  req.ItemGroupId,
		ItemType:                     req.ItemType,
		ItemLevel1:                   req.ItemLevel_1,
		ItemLevel2:                   req.ItemLevel_2,
		ItemLevel3:                   req.ItemLevel_3,
		ItemLevel4:                   req.ItemLevel_4,
		SupplierId:                   req.SupplierId,
		UnitOfMeasurementTypeId:      req.UnitOfMeasurementTypeId,
		UnitOfMeasurementSellingId:   req.UnitOfMeasurementSellingId,
		UnitOfMeasurementPurchaseId:  req.UnitOfMeasurementPurchaseId,
		UnitOfMeasurementStockId:     req.UnitOfMeasurementStockId,
		SalesItem:                    req.SalesItem,
		Lottable:                     req.Lottable,
		Inspection:                   req.Inspection,
		PriceListItem:                req.PriceListItem,
		StockKeeping:                 req.StockKeeping,
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
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *ItemRepositoryImpl) ChangeStatusItem(tx *gorm.DB, Id int) (bool, error) {
	var entities masteritementities.Item

	result := tx.Model(&entities).
		Where("item_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	if entities.IsActive {
		entities.IsActive = false
	} else {
		entities.IsActive = true
	}

	result = tx.Save(&entities)

	if result.Error != nil {
		return false, result.Error
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
		MillageEvery: request.MillageEvery,
		ReturnEvery:  request.ReturnEvery,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *ItemRepositoryImpl) GetAllItemDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	// Define a slice to hold Item Detail responses
	var responses []masteritempayloads.ItemDetailResponse

	// Define table struct
	tableStruct := masteritempayloads.ItemDetailRequest{}

	// Define join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	// Apply filters
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	// Execute query
	rows, err := whereQuery.Find(&responses).Rows()
	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	// Define a slice to hold map responses
	var mapResponses []map[string]interface{}

	// Iterate over rows
	for rows.Next() {
		// Define variables to hold row data
		var ItemDetailRes masteritempayloads.ItemDetailRequest

		// Scan the row into ItemDetailResponse struct
		if err := rows.Scan(
			&ItemDetailRes.ItemDetailId,
			&ItemDetailRes.ItemId,
			&ItemDetailRes.BrandId,
			&ItemDetailRes.MillageEvery,
			&ItemDetailRes.ModelId,
			&ItemDetailRes.IsActive,
			&ItemDetailRes.ReturnEvery,
			&ItemDetailRes.VariantId); err != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Convert DiscountGroupResponse to map
		responseMap := map[string]interface{}{
			"item_id":        ItemDetailRes.ItemId,
			"item_detail_id": ItemDetailRes.ItemDetailId,
			"brand_id":       ItemDetailRes.BrandId,
			"model_id":       ItemDetailRes.ModelId,
			"variant_id":     ItemDetailRes.VariantId,
			"is_active":      ItemDetailRes.IsActive,
			"millage_every":  ItemDetailRes.MillageEvery,
			"return_every":   ItemDetailRes.ReturnEvery,
		}

		// Append responseMap to the slice
		mapResponses = append(mapResponses, responseMap)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *ItemRepositoryImpl) GetItemDetailById(tx *gorm.DB, ItemId, ItemDetailId int) (masteritempayloads.ItemDetailRequest, *exceptionsss_test.BaseErrorResponse) {
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
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	response.ItemDetailId = entities.ItemDetailId
	response.ItemId = entities.ItemId
	response.BrandId = entities.BrandId
	response.ModelId = entities.ModelId
	response.VariantId = entities.VariantId
	response.MillageEvery = entities.MillageEvery
	response.ReturnEvery = entities.ReturnEvery
	response.IsActive = entities.IsActive

	return response, nil
}

func (r *ItemRepositoryImpl) AddItemDetail(tx *gorm.DB, ItemId int, req masteritempayloads.ItemDetailRequest) *exceptionsss_test.BaseErrorResponse {
	entities := masteritementities.ItemDetail{
		ItemId:       ItemId,
		BrandId:      req.BrandId,
		ModelId:      req.ModelId,
		VariantId:    req.VariantId,
		MillageEvery: req.MillageEvery,
		ReturnEvery:  req.ReturnEvery,
		IsActive:     req.IsActive,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return nil
}

func (r *ItemRepositoryImpl) DeleteItemDetail(tx *gorm.DB, ItemId int, ItemDetailId int) *exceptionsss_test.BaseErrorResponse {
	var entities masteritementities.ItemDetail

	result := tx.Model(&entities).
		Where("item_id = ? AND item_detail_id = ?", ItemId, ItemDetailId).
		Delete(&entities)

	if result.Error != nil {
		return &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return nil
}
