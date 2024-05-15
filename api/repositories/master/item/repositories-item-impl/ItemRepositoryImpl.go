package masteritemrepositoryimpl

import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
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

// GetUomItemDropDown implements masteritemrepository.ItemRepository.
func (r *ItemRepositoryImpl) GetUomDropDown(tx *gorm.DB, uomTypeId int) ([]masteritempayloads.UomDropdownResponse, *exceptionsss_test.BaseErrorResponse) {
	model := masteritementities.Uom{}
	responses := []masteritempayloads.UomDropdownResponse{}
	err := tx.Model(model).Where(masteritementities.Uom{UomTypeId: uomTypeId}).Scan(&responses).Error

	if err != nil {
		return responses, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return responses, nil
}

// GetUomTypeDropDown implements masteritemrepository.ItemRepository.
func (r *ItemRepositoryImpl) GetUomTypeDropDown(tx *gorm.DB) ([]masteritempayloads.UomTypeDropdownResponse, *exceptionsss_test.BaseErrorResponse) {
	model := masteritementities.UomType{}
	responses := []masteritempayloads.UomTypeDropdownResponse{}
	err := tx.Model(model).Scan(&responses).Error

	if err != nil {
		return responses, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return responses, nil
}

func (r *ItemRepositoryImpl) GetAllItem(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	var responses []masteritempayloads.ItemLookup
	tableStruct := masteritempayloads.ItemLookup{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	rows, err := joinTable.Scopes(pagination.Paginate(&tableStruct, &pages, whereQuery)).Scan(&responses).Rows()

	if err != nil {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	pages.Rows = responses

	return pages, nil
}

func (r *ItemRepositoryImpl) GetAllItemLookup(tx *gorm.DB, internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptionsss_test.BaseErrorResponse) {
	// var paginationResponse utils.APIPaginationResponse

	// var multiIds []string
	// var responses []masteritempayloads.ItemLookup
	// var getItemGroupResponse []masteritempayloads.ItemGroupResponse
	// var getSupplierMasterResponse []masteritempayloads.SupplierMasterResponse
	// tableStruct := masteritempayloads.ItemLookup{}
	// count := 0

	// var externalQueryFirst bool

	// var supplierCode string
	// var supplierName string

	// for _, value := range externalFilterCondition {

	// 	if value.ColumnField == "supplier_code" {
	// 		supplierCode = value.ColumnValue
	// 	} else if value.ColumnField == "supplier_name" {
	// 		supplierName = value.ColumnValue
	// 	}

	// 	if value.ColumnValue != "" {
	// 		externalQueryFirst = true

	// 	}
	// }

	// for i := 1; i < 10; i++ {

	// }
	// if externalQueryFirst {

	// 	supplierUrl := config.EnvConfigs.GeneralServiceUrl + "api/general/supplier-master?page=" + strconv.Itoa(pages.Page) + " &limit=" + strconv.Itoa(pages.Limit) + " &supplier_code=" + supplierCode + "&supplier_name=" + supplierName

	// } else {

	// 	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	// }

	// for _, value := range queryParams {
	// 	if value != "" {
	// 		count++
	// 	}
	// }

	// if count == 2 && queryParams["limit"] != "" && queryParams["page"] != "" {
	// 	page, _ := strconv.Atoi(queryParams["page"])
	// 	limit, _ := strconv.Atoi(queryParams["limit"])

	// 	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	// 	//execute
	// 	err := joinTable.Offset(page * limit).Limit(limit).Scan(&responses).Error

	// 	if err != nil {
	// 		fmt.Print(err)
	// 		return nil, &exceptionsss_test.BaseErrorResponse{
	// 			StatusCode: http.StatusInternalServerError,
	// 			Err:        err,
	// 		}
	// 	}

	// 	groupServiceUrl := "http://10.1.32.26:8000/general-service/api/general/filter-item-group?item_group_code=" + queryParams["item_group_code"]
	// 	errUrlItemGroup := utils.Get(groupServiceUrl, &getItemGroupResponse, nil)

	// 	if errUrlItemGroup != nil {
	// 		return nil, &exceptionsss_test.BaseErrorResponse{
	// 			StatusCode: http.StatusInternalServerError,
	// 			Err:        errUrlItemGroup,
	// 		}
	// 	}

	// 	joinedData := utils.DataFrameInnerJoin(responses, getItemGroupResponse, "ItemGroupId")

	// 	for _, item := range responses {
	// 		idStr := strconv.Itoa(item.SupplierId)
	// 		duplicate := false
	// 		for _, existingID := range multiIds {
	// 			if existingID == idStr {
	// 				duplicate = true
	// 				break
	// 			}
	// 		}
	// 		if !duplicate {
	// 			multiIds = append(multiIds, idStr)
	// 		}
	// 	}

	// 	supplierServiceUrl := "http://10.1.32.26:8000/general-service/api/general/supplier-master-multi-id/" + strings.Join(multiIds, ",")
	// 	errUrlSupplierMaster := utils.Get(supplierServiceUrl, &getSupplierMasterResponse, nil)
	// 	if errUrlSupplierMaster != nil {
	// 		return nil, &exceptionsss_test.BaseErrorResponse{
	// 			StatusCode: http.StatusInternalServerError,
	// 			Err:        errUrlSupplierMaster,
	// 		}
	// 	}

	// Convert paginated data to map format
	// var mapResponses []map[string]interface{}
	// for _, data := range responses {
	// 	// Fetch data from mtr_item_group for each response
	// 	itemGroupUrl := config.EnvConfigs.GeneralServiceUrl + "item-group/" + strconv.Itoa(data.ItemGroupId)
	// 	var itemGroupResp masteritempayloads.ItemGroupResponse
	// 	err := utils.Get(itemGroupUrl, &itemGroupResp, nil)
	// 	if err != nil {
	// 		return nil, 0, 0, err
	// 	}

	// 	// Create response map combining ItemLookup and ItemGroupResponse data
	// 	responseMap := map[string]interface{}{
	// 		"is_active":       data.IsActive,
	// 		"item_id":         data.ItemId,
	// 		"item_code":       data.ItemCode,
	// 		"item_name":       data.ItemName,
	// 		"item_group_id":   data.ItemGroupId,
	// 		"item_class_id":   data.ItemClassId,
	// 		"item_type":       data.ItemType,
	// 		"supplier_id":     data.SupplierId,
	// 		"item_group_name": itemGroupResp.ItemGroupName,
	// 		// Add more fields as needed
	// 	}
	// 	mapResponses = append(mapResponses, responseMap)
	// }

	// // Calculate total pages
	// totalRows = pages.TotalRows
	// totalPages = pages.TotalPages

	// return mapResponses, totalPages, int(totalRows), nil
	panic("unimplemented")
}

func (r *ItemRepositoryImpl) GetItemById(tx *gorm.DB, Id int) (map[string]any, *exceptionsss_test.BaseErrorResponse) {
	entities := masteritementities.Item{}
	response := masteritempayloads.ItemResponse{}

	rows, err := tx.Model(&entities).Select("mtr_item.*,u.*").
		Where(masteritementities.Item{
			ItemId: Id,
		}).InnerJoins("Join mtr_uom_item u ON mtr_item.item_id = u.item_id").
		First(&response).
		Rows()

	if err != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	supplierResponse := masteritempayloads.SupplierMasterResponse{}

	supplierUrl := config.EnvConfigs.GeneralServiceUrl + "/supplier-master/" + strconv.Itoa(response.SupplierId)

	if err := utils.Get(supplierUrl, &supplierResponse, nil); err != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	joinSupplierData := utils.DataFrameInnerJoin([]masteritempayloads.ItemResponse{response}, []masteritempayloads.SupplierMasterResponse{supplierResponse}, "SupplierId")

	// join with user details data (not yet complete)

	defer rows.Close()

	return joinSupplierData[0], nil
}

func (r *ItemRepositoryImpl) GetItemWithMultiId(tx *gorm.DB, MultiIds []string) ([]masteritempayloads.ItemResponse, *exceptionsss_test.BaseErrorResponse) {
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
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *ItemRepositoryImpl) GetItemCode(tx *gorm.DB, code string) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse) {
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
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	//FK Luar with mtr_item_group common-general service
	ItemGroupUrl := config.EnvConfigs.GeneralServiceUrl + "item-group/" + strconv.Itoa(response.ItemGroupId)
	errUrlItemGroup := utils.Get(ItemGroupUrl, &getItemGroupResponse, nil)
	fmt.Println("Fetching mtr_item_group data from:", ItemGroupUrl)
	if errUrlItemGroup != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	firstJoin := utils.DataFrameLeftJoin([]masteritempayloads.ItemResponse{response}, []masteritempayloads.ItemGroupResponse{getItemGroupResponse}, "ItemGroupId")

	//FK luar with mtr_supplier general service
	SupplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier-master/" + strconv.Itoa(response.SupplierId)
	errUrlSupplierMaster := utils.Get(SupplierUrl, &getSupplierMasterResponse, nil)
	fmt.Println("Fetching mtr_supplier data from:", SupplierUrl)
	if errUrlSupplierMaster != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	secondJoin := utils.DataFrameLeftJoin(firstJoin, []masteritempayloads.SupplierMasterResponse{getSupplierMasterResponse}, "SupplierId")
	//FK luar with storage_type general service
	StorageTypeUrl := config.EnvConfigs.GeneralServiceUrl + "storage-type/" + strconv.Itoa(response.StorageTypeId)
	errUrlStorageType := utils.Get(StorageTypeUrl, &getStorageTypeResponse, nil)
	fmt.Println("Fetching storage_type data from:", StorageTypeUrl)
	if errUrlStorageType != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	thirdJoin := utils.DataFrameLeftJoin(secondJoin, []masteritempayloads.StorageTypeResponse{getStorageTypeResponse}, "StorageTypeId")
	//FK luar with mtr_warranty_claim_type common service
	WarrantyClaimTypeUrl := config.EnvConfigs.GeneralServiceUrl + "warranty-claim-type/" + strconv.Itoa(response.AtpmWarrantyClaimTypeId)
	errUrlWarrantyClaimType := utils.Get(WarrantyClaimTypeUrl, &getAtpmWarrantyClaimTypeResponse, nil)
	fmt.Println("Fetching mtr_warranty_claim_type data from:", WarrantyClaimTypeUrl)
	if errUrlWarrantyClaimType != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	fourthJoin := utils.DataFrameLeftJoin(thirdJoin, []masteritempayloads.AtpmWarrantyClaimTypeResponse{getAtpmWarrantyClaimTypeResponse}, "AtpmWarrantyClaimTypeId")
	//FK luar with mtr_special_movement common service
	SpecialMovementUrl := config.EnvConfigs.GeneralServiceUrl + "special-movement/" + strconv.Itoa(response.SpecialMovementId)
	errUrlSpecialMovement := utils.Get(SpecialMovementUrl, &getSpecialMovementResponse, nil)
	fmt.Println("Fetching mtr_special_movement data from:", SpecialMovementUrl)
	if errUrlSpecialMovement != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	fifthJoin := utils.DataFrameLeftJoin(fourthJoin, []masteritempayloads.SpecialMovementResponse{getSpecialMovementResponse}, "SpecialMovementId")
	//FK luar with mtr_supplier general service atpm_supplier_id
	AtpmSupplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier-master/" + strconv.Itoa(response.AtpmSupplierId)
	errUrlAtpmSupplier := utils.Get(AtpmSupplierUrl, &getAtpmSupplierResponse, nil)
	fmt.Println("Fetching mtr_supplier data from:", AtpmSupplierUrl)
	if errUrlAtpmSupplier != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	sixthJoin := utils.DataFrameLeftJoin(fifthJoin, []masteritempayloads.AtpmSupplierResponse{getAtpmSupplierResponse}, "AtpmSupplierId")
	//FK luar with mtr_supplier general service atpm_supplier_code_order_id
	AtpmSupplierCodeOrderUrl := config.EnvConfigs.GeneralServiceUrl + "supplier-master/" + strconv.Itoa(response.AtpmSupplierCodeOrderId)
	errUrlAtpmSupplierCodeOrder := utils.Get(AtpmSupplierCodeOrderUrl, &getAtpmSupplierCodeOrderResponse, nil)
	fmt.Println("Fetching mtr_supplier data from:", AtpmSupplierCodeOrderUrl)
	if errUrlAtpmSupplierCodeOrder != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	seventhJoin := utils.DataFrameLeftJoin(sixthJoin, []masteritempayloads.AtpmSupplierCodeOrderResponse{getAtpmSupplierCodeOrderResponse}, "AtpmSupplierCodeOrderId")
	//FK luar with mtr_user_details general service
	PersonInChargeUrl := config.EnvConfigs.GeneralServiceUrl + "user-details-all/" + strconv.Itoa(response.PersonInChargeId)
	errUrlPersonInCharge := utils.Get(PersonInChargeUrl, &getPersonInChargeResponse, nil)
	fmt.Println("Fetching mtr_user_details data from:", PersonInChargeUrl)
	if errUrlPersonInCharge != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	eightJoin := utils.DataFrameLeftJoin(seventhJoin, getPersonInChargeResponse, "PersonInChargeId")

	// FK luar with mtr_unit_of_measurement_type
	// fk luar with mtr_atpm_order_type common service

	return eightJoin, nil
}

func (r *ItemRepositoryImpl) SaveItem(tx *gorm.DB, req masteritempayloads.ItemRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masteritementities.Item{
		ItemCode:                     req.ItemCode,
		ItemClassId:                  req.ItemClassId,
		ItemName:                     req.ItemName,
		ItemGroupId:                  req.ItemGroupId,
		ItemType:                     req.ItemType,
		ItemLevel1:                   req.ItemLevel1,
		ItemLevel2:                   req.ItemLevel2,
		ItemLevel3:                   req.ItemLevel3,
		ItemLevel4:                   req.ItemLevel4,
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
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	model := masteritementities.Item{}

	err = tx.Model(&model).Where(masteritementities.Item{ItemCode: req.ItemCode}).First(&model).Error

	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	atpmResponse := masteritempayloads.AtpmOrderTypeResponse{}

	atpmOrderTypeUrl := config.EnvConfigs.GeneralServiceUrl + "/atpm-order-type/" + strconv.Itoa(req.SourceTypeId)

	if err := utils.Get(atpmOrderTypeUrl, &atpmResponse, nil); err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	uomTypeModel := masteritementities.UomType{}

	err = tx.Model(&uomTypeModel).Where(masteritementities.UomType{UomTypeId: req.UnitOfMeasurementTypeId}).First(&uomTypeModel).Error

	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	uomItemEntities := masteritementities.UomItem{
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
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *ItemRepositoryImpl) ChangeStatusItem(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities masteritementities.Item

	result := tx.Model(&entities).
		Where("item_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
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
		return false, &exceptionsss_test.BaseErrorResponse{
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
	var responses []masteritempayloads.ItemDetailRequest

	responseStruct := reflect.TypeOf(masteritempayloads.ItemDetailRequest{})

	// Filter internal service conditions
	var internalServiceFilter []utils.FilterCondition
	for _, condition := range filterCondition {
		for j := 0; j < responseStruct.NumField(); j++ {
			if condition.ColumnField == responseStruct.Field(j).Tag.Get("parent_entity")+"."+responseStruct.Field(j).Tag.Get("json") {
				internalServiceFilter = append(internalServiceFilter, condition)
				break
			}
		}
	}

	// Apply internal service filter conditions
	tableStruct := masteritempayloads.ItemDetailRequest{}
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	whereQuery := utils.ApplyFilter(joinTable, internalServiceFilter)

	// Fetch data from database
	err := whereQuery.Find(&responses).Error
	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to fetch data from database: %w", err),
		}
	}

	// Check if responses are empty
	if len(responses) == 0 {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("no data found"),
		}
	}

	// Define a slice to hold map responses
	var mapResponses []map[string]interface{}

	// Iterate over responses and convert them to maps
	for _, response := range responses {
		responseMap := map[string]interface{}{
			"is_active":      response.IsActive,
			"item_detail_id": response.ItemDetailId,
			"item_id":        response.ItemId,
			"brand_id":       response.BrandId,
			"millage_every":  response.MillageEvery,
			"model_id":       response.ModelId,
			"return_every":   response.ReturnEvery,
			"variant_id":     response.VariantId,
			// Add other fields as needed
		}
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
