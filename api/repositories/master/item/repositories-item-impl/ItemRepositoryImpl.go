package masteritemrepositoryimpl

import (
	config "after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"log"
	"net/url"
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

func (r *ItemRepositoryImpl) GetAllItem(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, error) {
	// Define variables
	var (
		responses    []masteritempayloads.ItemLookup
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

	// Calculate offset and limit
	offset := int64((pages.Page - 1) * pages.Limit)
	if offset < 0 {
		offset = 0
	}
	limit := pages.Limit

	// Convert offset to int
	intOffset := int(offset)

	// Apply pagination
	whereQuery = whereQuery.Offset(intOffset).Limit(limit)

	// Execute query
	if err := whereQuery.Find(&responses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, gorm.ErrRecordNotFound
		}
		return nil, 0, 0, err
	}

	// Convert responses to map format
	for _, response := range responses {
		responseMap := map[string]interface{}{
			"item_id":   response.ItemId,
			"item_name": response.ItemName,
			// Add other fields as needed
		}
		mapResponses = append(mapResponses, responseMap)
	}

	// Paginate the response data
	paginatedData, totalPages, _ := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, int(totalRows), nil
}

func (r *ItemRepositoryImpl) GetAllItemLookup(tx *gorm.DB, queryParams map[string]string) ([]map[string]interface{}, error) {
	var (
		paginationResponse        utils.APIPaginationResponse
		multiIds                  []string
		responses                 []masteritempayloads.ItemLookup
		getItemGroupResponse      []masteritempayloads.ItemGroupResponse
		getSupplierMasterResponse []masteritempayloads.SupplierMasterResponse
		tableStruct               = masteritempayloads.ItemLookup{}
	)

	// Count the non-empty query parameters
	count := 0
	for _, value := range queryParams {
		if value != "" {
			count++
		}
	}

	// Check if the required parameters are present for pagination
	if count == 2 && queryParams["limit"] != "" && queryParams["page"] != "" {
		page, _ := strconv.Atoi(queryParams["page"])
		limit, _ := strconv.Atoi(queryParams["limit"])

		// Apply joins and execute query
		joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
		rows, err := joinTable.Offset(page * limit).Limit(limit).Scan(&responses).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		// Get item group data
		groupServiceUrl := config.EnvConfigs.GeneralServiceUrl + "/api/general/filter-item-group?item_group_code=" + queryParams["item_group_code"]
		errUrlItemGroup := utils.Get(groupServiceUrl, &getItemGroupResponse, nil)
		if errUrlItemGroup != nil {
			return nil, errUrlItemGroup
		}

		// Join item group data with responses
		joinedData := utils.DataFrameInnerJoin(responses, getItemGroupResponse, "ItemGroupId")

		// Extract unique supplier IDs
		for _, item := range responses {
			idStr := strconv.Itoa(item.SupplierId)
			duplicate := false
			for _, existingID := range multiIds {
				if existingID == idStr {
					duplicate = true
					break
				}
			}
			if !duplicate {
				multiIds = append(multiIds, idStr)
			}
		}

		// Get supplier master data
		supplierServiceUrl := config.EnvConfigs.GeneralServiceUrl + "/api/general/supplier-master-multi-id/" + strings.Join(multiIds, ",")
		errUrlSupplierMaster := utils.Get(supplierServiceUrl, &getSupplierMasterResponse, nil)
		if errUrlSupplierMaster != nil {
			return nil, errUrlSupplierMaster
		}

		// Join supplier master data with joined data
		joinedDataSecond := utils.DataFrameInnerJoin(joinedData, getSupplierMasterResponse, "SupplierId")

		return joinedDataSecond, nil
	}

	// If pagination parameters are not present, fetch data without pagination
	supplierDescUrl := config.EnvConfigs.GeneralServiceUrl + "/api/general/supplier-master-for-item-master"
	u, err := url.Parse(supplierDescUrl)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	supplierDescUrl = u.String()

	// Fetch data with pagination from the URL
	paginationRes, errUrlSupplierMasterLookup := utils.GetWithPagination(supplierDescUrl, paginationResponse, nil)
	if errUrlSupplierMasterLookup != nil {
		return nil, errUrlSupplierMasterLookup
	}

	dataSlice, _ := paginationRes.Data.([]map[string]interface{})

	return dataSlice, nil
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
	entities := masteritementities.Item{}
	response := masteritempayloads.ItemResponse{}
	var getSupplierMasterResponse masteritempayloads.SupplierMasterResponse
	var getItemGroupResponse masteritempayloads.ItemGroupResponse
	var getStorageTypeResponse masteritempayloads.StorageTypeResponse
	var getSpecialMovementResponse masteritempayloads.SpecialMovementResponse
	var getAtpmSupplierResponse masteritempayloads.AtpmSupplierResponse
	var getAtpmSupplierCodeOrderResponse masteritempayloads.AtpmSupplierCodeOrderResponse
	// var getPersonInChargeResponse masteritempayloads.PersonInChargeResponse
	var getAtpmWarrantyClaimTypeResponse masteritempayloads.AtpmWarrantyClaimTypeResponse

	rows, err := tx.Model(&entities).
		Where(masteritementities.Item{
			ItemCode: code,
		}).First(&response).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//FK Luar with mtr_item_group common-general service
	errUrlItemGroup := utils.Get(config.EnvConfigs.GeneralServiceUrl+"/api/general/item-group/"+strconv.Itoa(response.ItemGroupId), &getItemGroupResponse, nil)

	if errUrlItemGroup != nil {
		return nil, err
	}

	firstJoin := utils.DataFrameLeftJoin([]masteritempayloads.ItemResponse{response}, []masteritempayloads.ItemGroupResponse{getItemGroupResponse}, "ItemGroupId")

	//FK luar with mtr_supplier general service
	errUrlSupplierMaster := utils.Get(config.EnvConfigs.GeneralServiceUrl+"/api/general/supplier-master/"+strconv.Itoa(response.SupplierId), &getSupplierMasterResponse, nil)

	if errUrlSupplierMaster != nil {
		return nil, err
	}

	secondJoin := utils.DataFrameLeftJoin(firstJoin, []masteritempayloads.SupplierMasterResponse{getSupplierMasterResponse}, "SupplierId")
	//FK luar with storage_type general service
	errUrlStorageType := utils.Get(config.EnvConfigs.GeneralServiceUrl+"/api/general/storage-type/"+strconv.Itoa(response.StorageTypeId), &getStorageTypeResponse, nil)

	if errUrlStorageType != nil {
		return nil, err
	}

	thirdJoin := utils.DataFrameLeftJoin(secondJoin, []masteritempayloads.StorageTypeResponse{getStorageTypeResponse}, "StorageTypeId")
	//FK luar with mtr_warranty_claim_type common service
	errUrlWarrantyClaimType := utils.Get(config.EnvConfigs.GeneralServiceUrl+"/api/general/warranty-claim-type/"+strconv.Itoa(response.AtpmWarrantyClaimTypeId), &getAtpmWarrantyClaimTypeResponse, nil)

	if errUrlWarrantyClaimType != nil {
		return thirdJoin, err
	}

	fourthJoin := utils.DataFrameLeftJoin(thirdJoin, []masteritempayloads.AtpmWarrantyClaimTypeResponse{getAtpmWarrantyClaimTypeResponse}, "AtpmWarrantyClaimTypeId")
	//FK luar with mtr_special_movement common service
	errUrlSpecialMovement := utils.Get(config.EnvConfigs.GeneralServiceUrl+"/api/general/special-movement/"+strconv.Itoa(response.SpecialMovementId), &getSpecialMovementResponse, nil)

	if errUrlSpecialMovement != nil {
		return fourthJoin, err
	}

	fifthJoin := utils.DataFrameLeftJoin(fourthJoin, []masteritempayloads.SpecialMovementResponse{getSpecialMovementResponse}, "SpecialMovementId")
	//FK luar with mtr_supplier general service atpm_supplier_id
	errUrlAtpmSupplier := utils.Get(config.EnvConfigs.GeneralServiceUrl+"/api/general/supplier-master/"+strconv.Itoa(response.AtpmSupplierId), &getAtpmSupplierResponse, nil)

	if errUrlAtpmSupplier != nil {
		return fifthJoin, err
	}

	sixthJoin := utils.DataFrameLeftJoin(fifthJoin, []masteritempayloads.AtpmSupplierResponse{getAtpmSupplierResponse}, "AtpmSupplierId")
	//FK luar with mtr_supplier general service atpm_supplier_code_order_id
	errUrlAtpmSupplierCodeOrder := utils.Get(config.EnvConfigs.GeneralServiceUrl+"/api/general/supplier-master/"+strconv.Itoa(response.AtpmSupplierCodeOrderId), &getAtpmSupplierCodeOrderResponse, nil)

	if errUrlAtpmSupplierCodeOrder != nil {
		return sixthJoin, err
	}

	seventhJoin := utils.DataFrameLeftJoin(sixthJoin, []masteritempayloads.AtpmSupplierCodeOrderResponse{getAtpmSupplierCodeOrderResponse}, "AtpmSupplierCodeOrderId")
	//FK luar with mtr_user_details general service
	// errUrlPersonInCharge := utils.Get(c, "http://10.1.32.26:8000/general-service/api/general/user-details-all/"+strconv.Itoa(response.PersonInChargeId), &getPersonInChargeResponse, nil)
	// if errUrlPersonInCharge != nil {
	// 	return seventhJoin, err
	// }

	// joinedDataPersonInCharge := utils.DataFrameLeftJoin(seventhJoin, getPersonInChargeResponse, "PersonInChargeId")

	// FK luar with mtr_unit_of_measurement_type
	// fk luar with mtr_atpm_order_type common service

	return seventhJoin, nil
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

func (r *ItemRepositoryImpl) GetAllItemDetail(tx *gorm.DB, filterCondition []utils.FilterCondition) ([]map[string]interface{}, error) {
	// var responses []masteritempayloads.ItemDetailResponse
	entities := []masteritementities.ItemClass{}
	var responses []masteritempayloads.ItemClassResponse
	// var getLineTypeResponse []masteritempayloads.LineTypeResponse
	// var getItemGroupResponse []masteritempayloads.ItemGroupResponse
	// var c *gin.Context
	// var internalServiceFilter, externalServiceFilter []utils.FilterCondition
	// var groupName, lineTypeCode string
	// responseStruct := reflect.TypeOf(masteritempayloads.ItemClassResponse{})

	// for i := 0; i < len(filterCondition); i++ {
	// 	flag := false
	// 	for j := 0; j < responseStruct.NumField(); j++ {
	// 		if filterCondition[i].ColumnField == responseStruct.Field(j).Tag.Get("parent_entity")+"."+responseStruct.Field(j).Tag.Get("json") {
	// 			internalServiceFilter = append(internalServiceFilter, filterCondition[i])
	// 			flag = true
	// 			break
	// 		}
	// 		if !flag {
	// 			externalServiceFilter = append(externalServiceFilter, filterCondition[i])
	// 		}
	// 	}
	// }

	// //apply external services filter
	// for i := 0; i < len(externalServiceFilter); i++ {
	// 	if strings.Contains(externalServiceFilter[i].ColumnField, "line_type_code") {
	// 		lineTypeCode = externalServiceFilter[i].ColumnValue
	// 	} else {
	// 		groupName = externalServiceFilter[i].ColumnValue
	// 	}
	// }

	//define base model
	baseModelQuery := tx.Model(&entities)
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	//whereQuery := utils.ApplyFilter(baseModelQuery, internalServiceFilter)
	//apply pagination and execute
	rows, err := whereQuery.Scan(&responses).Rows()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if len(responses) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// groupServiceUrl := ""

	// errUrlItemGroup := utils.Get(c, groupServiceUrl, &getItemGroupResponse, nil)

	// if errUrlItemGroup != nil {
	// 	return nil, errUrlItemGroup
	// }

	// joinedData := utils.DataFrameInnerJoin(responses, getItemGroupResponse, "ItemGroupId")

	return nil, nil
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
