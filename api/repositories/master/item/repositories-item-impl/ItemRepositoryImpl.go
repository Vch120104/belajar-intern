package masteritemrepositoryimpl

import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"log"
	"net/http"
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

	// 	joinedDataSecond := utils.DataFrameInnerJoin(joinedData, getSupplierMasterResponse, "SupplierId")

	// 	return joinedDataSecond, nil
	// }

	// supplierDescUrl := "http://10.1.32.26:8000/general-service/api/general/supplier-master-for-item-master"

	// u, err := url.Parse(supplierDescUrl)
	// if err != nil {
	// 	return nil, &exceptionsss_test.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Err:        err,
	// 	}
	// }
	// q := u.Query()
	// for key, value := range queryParams {
	// 	q.Set(key, value)
	// }
	// u.RawQuery = q.Encode()
	// supplierDescUrl = u.String()

	// paginationRes, errUrlSupplierMasterLookup := utils.GetWithPagination(supplierDescUrl, paginationResponse, nil)
	// if errUrlSupplierMasterLookup != nil {
	// 	return nil, &exceptionsss_test.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Err:        err,
	// 	}
	// }

	// dataSlice, _ := paginationRes.Data.([]map[string]interface{})

	// return dataSlice, nil
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
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *ItemRepositoryImpl) GetItemCode(tx *gorm.DB, code string) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse) {
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
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	//FK Luar with mtr_item_group common-general service
	errUrlItemGroup := utils.Get("http://10.1.32.26:8000/general-service/api/general/item-group/"+strconv.Itoa(response.ItemGroupId), &getItemGroupResponse, nil)

	if errUrlItemGroup != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	firstJoin := utils.DataFrameLeftJoin([]masteritempayloads.ItemResponse{response}, []masteritempayloads.ItemGroupResponse{getItemGroupResponse}, "ItemGroupId")

	//FK luar with mtr_supplier general service
	errUrlSupplierMaster := utils.Get("http://10.1.32.26:8000/general-service/api/general/supplier-master/"+strconv.Itoa(response.SupplierId), &getSupplierMasterResponse, nil)

	if errUrlSupplierMaster != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	secondJoin := utils.DataFrameLeftJoin(firstJoin, []masteritempayloads.SupplierMasterResponse{getSupplierMasterResponse}, "SupplierId")
	//FK luar with storage_type general service
	errUrlStorageType := utils.Get("http://10.1.32.26:8000/general-service/api/general/storage-type/"+strconv.Itoa(response.StorageTypeId), &getStorageTypeResponse, nil)

	if errUrlStorageType != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	thirdJoin := utils.DataFrameLeftJoin(secondJoin, []masteritempayloads.StorageTypeResponse{getStorageTypeResponse}, "StorageTypeId")
	//FK luar with mtr_warranty_claim_type common service
	errUrlWarrantyClaimType := utils.Get("http://10.1.32.26:8000/general-service/api/general/warranty-claim-type/"+strconv.Itoa(response.AtpmWarrantyClaimTypeId), &getAtpmWarrantyClaimTypeResponse, nil)

	if errUrlWarrantyClaimType != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	fourthJoin := utils.DataFrameLeftJoin(thirdJoin, []masteritempayloads.AtpmWarrantyClaimTypeResponse{getAtpmWarrantyClaimTypeResponse}, "AtpmWarrantyClaimTypeId")
	//FK luar with mtr_special_movement common service
	errUrlSpecialMovement := utils.Get("http://10.1.32.26:8000/general-service/api/general/special-movement/"+strconv.Itoa(response.SpecialMovementId), &getSpecialMovementResponse, nil)

	if errUrlSpecialMovement != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	fifthJoin := utils.DataFrameLeftJoin(fourthJoin, []masteritempayloads.SpecialMovementResponse{getSpecialMovementResponse}, "SpecialMovementId")
	//FK luar with mtr_supplier general service atpm_supplier_id
	errUrlAtpmSupplier := utils.Get("http://10.1.32.26:8000/general-service/api/general/supplier-master/"+strconv.Itoa(response.AtpmSupplierId), &getAtpmSupplierResponse, nil)

	if errUrlAtpmSupplier != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	sixthJoin := utils.DataFrameLeftJoin(fifthJoin, []masteritempayloads.AtpmSupplierResponse{getAtpmSupplierResponse}, "AtpmSupplierId")
	//FK luar with mtr_supplier general service atpm_supplier_code_order_id
	errUrlAtpmSupplierCodeOrder := utils.Get("http://10.1.32.26:8000/general-service/api/general/supplier-master/"+strconv.Itoa(response.AtpmSupplierCodeOrderId), &getAtpmSupplierCodeOrderResponse, nil)

	if errUrlAtpmSupplierCodeOrder != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
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

func (r *ItemRepositoryImpl) GetAllItemDetail(tx *gorm.DB, filterCondition []utils.FilterCondition) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse) {
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
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	if len(responses) == 0 {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	// groupServiceUrl := ""

	// errUrlItemGroup := utils.Get(c, groupServiceUrl, &getItemGroupResponse, nil)

	// if errUrlItemGroup != nil {
	// 	return nil, errUrlItemGroup
	// }

	// joinedData := utils.DataFrameInnerJoin(responses, getItemGroupResponse, "ItemGroupId")

	return nil, nil
}

func (r *ItemRepositoryImpl) SaveItemDetail(tx *gorm.DB, request masteritempayloads.ItemDetailResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
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
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}
