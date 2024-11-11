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
	generalserviceapiutils "after-sales/api/utils/general-service"
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
		supplierParams := generalserviceapiutils.SupplierMasterParams{
			Page:         0,
			Limit:        100000,
			SupplierCode: supplierCode,
			SupplierName: supplierName,
		}
		supplierResponse, supplierError := generalserviceapiutils.GetAllSupplierMaster(supplierParams)
		if supplierError != nil {
			return nil, 0, 0, supplierError
		}

		for _, supplier := range supplierResponse.Data {
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
			mil4.item_level_4_name,
			uom.uom_item_id,
			uom.source_uom_id,
			uom.target_uom_id,
			uom.source_convertion,
			uom.target_convertion
		`).
		Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = mtr_item.item_level_1_id").
		Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = mtr_item.item_level_2_id").
		Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = mtr_item.item_level_3_id").
		Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = mtr_item.item_level_4_id").
		Joins("LEFT JOIN mtr_uom_item uom on uom.item_id = mtr_item.item_id and uom.uom_source_type_code = 'P'").
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
	if response.SupplierId != nil {
		supplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(*response.SupplierId)
		if err := utils.Get(supplierUrl, &supplierResponse, nil); err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to fetch supplier data",
				Err:        err,
			}
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

	rows, err := tx.Model(&entities).
		Select(`
			mtr_item.*,
			mil1.item_level_1_code,
			mil1.item_level_1_name,
			mil2.item_level_2_code,
			mil2.item_level_2_name,
			mil3.item_level_3_code,
			mil3.item_level_3_name,
			mil4.item_level_4_code,
			mil4.item_level_4_name,
			u.uom_item_id,
			u.source_uom_id,
			u.target_uom_id,
			u.source_convertion,
			u.target_convertion
		`).
		Where(masteritementities.Item{ItemCode: code}).
		Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = mtr_item.item_level_1_id").
		Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = mtr_item.item_level_2_id").
		Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = mtr_item.item_level_3_id").
		Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = mtr_item.item_level_4_id").
		Joins("LEFT JOIN mtr_uom_item u ON mtr_item.item_id = u.item_id AND u.uom_source_type_code = 'P'").
		First(&response).
		Rows()

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
			Err:        err,
		}
	}

	supplierResponse := masteritempayloads.SupplierMasterResponse{}
	if response.SupplierId != nil {
		supplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(*response.SupplierId)
		if err := utils.Get(supplierUrl, &supplierResponse, nil); err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
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

	if req.ItemId != 0 {
		var itemExist int64
		err := tx.Model(&masteritementities.Item{}).Where(masteritementities.Item{ItemId: req.ItemId}).Count(&itemExist).Error
		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "error fetching item data",
				Err:        err,
			}
		}
		if itemExist == 0 {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("item data not found"),
			}
		}
	}

	//CHECK ITEM TYPE EXISTENCE
	shouldReturn, returnValue, errorItemType := checkItemTypeExistence(tx, req, response)
	if shouldReturn {
		return returnValue, errorItemType
	}

	//CHECK ITEM LEVEL EXISTENCE
	shouldReturn1, returnValue2, errorItemLevel := checkIfItemLevelExists(tx, req, response)
	if shouldReturn1 {
		return returnValue2, errorItemLevel
	}

	//CHECK ITEM CLASS EXISTENCE
	shouldReturn2, returnValue1, errorItemClass := checkItemClasExistence(tx, req, response)
	if shouldReturn2 {
		return returnValue1, errorItemClass
	}

	//CHECK ITEM GROUP EXISTENCE
	shouldReturn3, returnValue3, errorItemGroup := checkItemGroupExistence(tx, req, response)
	if shouldReturn3 {
		return returnValue3, errorItemGroup
	}

	//CHECK UOM TYPE EXISTENCE
	shouldReturn4, returnValue4, errorUomType := checkUomTypeExistence(tx, req, response)
	if shouldReturn4 {
		return returnValue4, errorUomType
	}

	//CHECK UOM EXISTENCE
	shouldReturn5, returnValue5, errorUom := checkUomExistence(tx, req, response)
	if shouldReturn5 {
		return returnValue5, errorUom
	}

	//CHECK DISCOUNT EXISTENCE
	shouldReturn6, returnValue6, errorDiscount := checkDiscountExistence(tx, req, response)
	if shouldReturn6 {
		return returnValue6, errorDiscount
	}

	//CHECK MARKUP MASTER EXISTENCE
	shouldReturn7, returnValue7, errorMarkupMaster := checkMarkupMasterExistence(tx, req, response)
	if shouldReturn7 {
		return returnValue7, errorMarkupMaster
	}

	//CHECK PRINCIPAL CATALOG EXISTENCE
	shouldReturn8, returnValue8, errorPrincipalCatalog := checkPrincipalCatalogExistence(tx, req, response)
	if shouldReturn8 {
		return returnValue8, errorPrincipalCatalog
	}

	//CHECK PRINCIPAL BRAND PARENT EXISTENCE
	shouldReturn9, returnValue9, errorPrincipalBrandParent := checkPrincipalBrandParentExistence(tx, req, response)
	if shouldReturn9 {
		return returnValue9, errorPrincipalBrandParent
	}

	//CHECK SUPPLIER
	if req.SupplierId != nil {
		fmt.Println("call supplier api")
		supplierResponse, supplierErr := generalserviceapiutils.GetSupplierMasterByID(*req.SupplierId)
		if supplierErr != nil || supplierResponse.SupplierId == 0 {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Supplier not found",
				Err:        errors.New("supplier not found"),
			}
		}
	}

	//CHECK WARRANTY CLAIM TYPE EXISTENCE
	if req.AtpmWarrantyClaimTypeId != nil {
		warrantyClaimTypeResponse, warrantyClaimTypeError := generalserviceapiutils.GetWarrantyClaimTypeById(*req.AtpmWarrantyClaimTypeId)
		if warrantyClaimTypeError != nil || warrantyClaimTypeResponse.WarrantyClaimTypeId == 0 {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Warranty Claim Type not found",
				Err:        errors.New("warranty claim type not found"),
			}
		}
	}

	//CHECK SPECIAL MOVEMENT EXISTENCE
	if req.SpecialMovementId != nil {
		specialMovementResponse, specialMovementError := generalserviceapiutils.GetSpecialMovementById(*req.SpecialMovementId)
		if specialMovementError != nil || specialMovementResponse.SpecialMovementId == 0 {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Special Movement not found",
				Err:        errors.New("special movement not found"),
			}
		}
	}

	//CHECK ATPM SUPPLIER EXISTENCE
	if req.AtpmSupplierId != nil {
		atpmSupplierResponse, atpmSupplierError := generalserviceapiutils.GetSupplierMasterByID(*req.AtpmSupplierId)
		if atpmSupplierError != nil || atpmSupplierResponse.SupplierId == 0 {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "ATPM Supplier not found",
				Err:        errors.New("atpm supplier not found"),
			}
		}
	}

	//CHECK ITEM REGULATION EXISTENCE
	if req.ItemRegulationId != nil {
		itemRegulationResponse, itemRegulationError := generalserviceapiutils.GetItemRegulationById(*req.ItemRegulationId)
		if itemRegulationError != nil || itemRegulationResponse.ItemRegulationId == 0 {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Item Regulation not found",
				Err:        errors.New("item regulation not found"),
			}
		}
	}

	//CHECK ATPM ORDER TYPE EXISTENCE
	if req.SourceTypeId != nil {
		atpmOrderTypeResponse, atpmOrderTypeError := generalserviceapiutils.GetAtpmOrderTypeById(*req.SourceTypeId)
		if atpmOrderTypeError != nil || atpmOrderTypeResponse.AtpmOrderTypeId == 0 {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "ATPM Order Type not found",
				Err:        errors.New("atpm order type not found"),
			}
		}
	}

	//CHECK ATPM SUPPLIER CODE ORDER EXISTENCE
	if req.AtpmSupplierCodeOrderId != nil {
		atpmSupplierCodeOrderResponse, atpmSupplierCodeOrderError := generalserviceapiutils.GetSupplierMasterByID(*req.AtpmSupplierCodeOrderId)
		if atpmSupplierCodeOrderError != nil || atpmSupplierCodeOrderResponse.SupplierId == 0 {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "ATPM Supplier Order not found",
				Err:        errors.New("atpm Supplier Order not found"),
			}
		}
	}

	//CHECK PERSON IN CHARGE EXISTENCE
	if req.PersonInChargeId != nil {
		picResponse, picError := generalserviceapiutils.GetEmployeeByID(*req.PersonInChargeId)
		if picError != nil || picResponse.UserEmployeeId == 0 {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Person in Charge not found",
				Err:        errors.New("person in charge not found"),
			}
		}
	}
	entities := masteritementities.Item{
		IsActive:                     req.IsActive,
		ItemId:                       req.ItemId,
		ItemCode:                     req.ItemCode,
		ItemClassId:                  req.ItemClassId,
		ItemName:                     req.ItemName,
		ItemGroupId:                  req.ItemGroupId,
		ItemTypeId:                   req.ItemTypeId,
		ItemLevel1Id:                 req.ItemLevel1Id,
		ItemLevel2Id:                 req.ItemLevel2Id,
		ItemLevel3Id:                 req.ItemLevel3Id,
		ItemLevel4Id:                 req.ItemLevel4Id,
		SupplierId:                   req.SupplierId,
		UnitOfMeasurementTypeId:      req.UnitOfMeasurementTypeId,
		UnitOfMeasurementSellingId:   req.UnitOfMeasurementSellingId,
		UnitOfMeasurementPurchaseId:  req.UnitOfMeasurementPurchaseId,
		UnitOfMeasurementStockId:     req.UnitOfMeasurementStockId,
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
		ItemRegulationId:             req.ItemRegulationId,
		AutoPickWms:                  req.AutoPickWms,
		PrincipalCatalogId:           req.PrincipalCatalogId,
		PrincipalBrandParentId:       req.PrincipalBrandParentId,
		ProportionalSupplyWms:        req.ProportionalSupplyWms,
		Remark2:                      req.Remark2,
		Remark3:                      req.Remark3,
		SourceTypeId:                 req.SourceTypeId,
		AtpmSupplierCodeOrderId:      req.AtpmSupplierCodeOrderId,
		PersonInChargeId:             req.PersonInChargeId,
		IsSellable:                   req.IsSellable,
	}

	err := tx.Save(&entities).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	entitiyItemType := masteritementities.ItemType{}
	err = tx.Model(&entitiyItemType).Where(masteritementities.ItemType{ItemTypeId: entities.ItemTypeId}).First(&entitiyItemType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("item type not found"),
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch item type data",
			Err:        err,
		}
	}

	// only insert to mtr_uom_item if item type is Goods
	if entitiyItemType.ItemTypeCode == "G" {
		entityUomItemPurchaseReq := masteritementities.UomItem{
			IsActive:          true,
			ItemId:            entities.ItemId,
			UomSourceTypeCode: "P",
			UomTypeId:         *req.UnitOfMeasurementTypeId,
			SourceUomId:       *req.UnitOfMeasurementPurchaseId,
			TargetUomId:       *req.UnitOfMeasurementStockId,
			SourceConvertion:  req.SourceConvertion,
			TargetConvertion:  req.TargetConvertion,
		}
		entityUomItemSellingReq := masteritementities.UomItem{
			IsActive:          true,
			ItemId:            entities.ItemId,
			UomSourceTypeCode: "S",
			UomTypeId:         *req.UnitOfMeasurementTypeId,
			SourceUomId:       *req.UnitOfMeasurementStockId,
			TargetUomId:       *req.UnitOfMeasurementStockId,
			SourceConvertion:  1,
			TargetConvertion:  1,
		}

		entityUomItemPurchase := masteritementities.UomItem{}
		err = tx.Model(&entityUomItemPurchase).
			Where(masteritementities.UomItem{
				ItemId:            entities.ItemId,
				UomSourceTypeCode: "P",
			}).
			First(&entityUomItemPurchase).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "error fetching uom item purchase",
				Err:        err,
			}
		}

		entityUomItemSelling := masteritementities.UomItem{}
		err = tx.Model(&entityUomItemSelling).
			Where(masteritementities.UomItem{
				ItemId:            entities.ItemId,
				UomSourceTypeCode: "S",
			}).
			First(&entityUomItemSelling).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "error fetching uom item selling",
				Err:        err,
			}
		}

		// initiate update uom item if uom item already exist
		if entityUomItemPurchase.UomItemId != 0 {
			entityUomItemPurchaseReq.UomItemId = entityUomItemPurchase.UomItemId
		}
		if entityUomItemSelling.UomItemId != 0 {
			entityUomItemSellingReq.UomItemId = entityUomItemSelling.UomItemId
		}

		if err := tx.Save(&entityUomItemPurchaseReq).Error; err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed insert/update uom item purchase",
				Err:        err,
			}
		}

		if err := tx.Save(&entityUomItemSellingReq).Error; err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed insert/update uom item selling",
				Err:        err,
			}
		}
	}

	result := masteritempayloads.ItemSaveResponse{
		IsActive:     entities.IsActive,
		ItemId:       entities.ItemId,
		ItemName:     entities.ItemName,
		ItemCode:     entities.ItemCode,
		ItemTypeId:   entities.ItemTypeId,
		ItemLevel1Id: entities.ItemLevel1Id,
		ItemLevel2Id: entities.ItemLevel2Id,
		ItemLevel3Id: entities.ItemLevel3Id,
		ItemLevel4Id: entities.ItemLevel4Id,
	}

	return result, nil
}

func checkPrincipalBrandParentExistence(tx *gorm.DB, req masteritempayloads.ItemRequest, response masteritempayloads.ItemSaveResponse) (bool, masteritempayloads.ItemSaveResponse, *exceptions.BaseErrorResponse) {
	if req.PrincipalBrandParentId != nil {
		var countPrincipalBrandParent int64
		if err := tx.Model(&masteritementities.PrincipalBrandParent{}).
			Where(masteritementities.PrincipalBrandParent{PrincipalBrandParentId: *req.PrincipalBrandParentId}).
			Count(&countPrincipalBrandParent).Error; err != nil {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Database error on Principal Brand Parent",
				Err:        err,
			}
		}
		if countPrincipalBrandParent == 0 {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Principal Brand Parent not found",
			}
		}
	}

	return false, masteritempayloads.ItemSaveResponse{}, nil
}

func checkPrincipalCatalogExistence(tx *gorm.DB, req masteritempayloads.ItemRequest, response masteritempayloads.ItemSaveResponse) (bool, masteritempayloads.ItemSaveResponse, *exceptions.BaseErrorResponse) {
	if req.PrincipalCatalogId != nil {
		var countPrincipalCatalog int64
		if err := tx.Model(&masteritementities.PrincipalCatalog{}).
			Where(masteritementities.PrincipalCatalog{PrincipalCatalogId: *req.PrincipalCatalogId}).
			Count(&countPrincipalCatalog).Error; err != nil {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Database error on Principal Catalog",
				Err:        err,
			}
		}
		if countPrincipalCatalog == 0 {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Principal Catalog not found",
			}
		}
	}

	return false, masteritempayloads.ItemSaveResponse{}, nil
}

func checkMarkupMasterExistence(tx *gorm.DB, req masteritempayloads.ItemRequest, response masteritempayloads.ItemSaveResponse) (bool, masteritempayloads.ItemSaveResponse, *exceptions.BaseErrorResponse) {
	if req.MarkupMasterId != nil {
		var countMarkupMaster int64
		if err := tx.Model(&masteritementities.MarkupMaster{}).
			Where(masteritementities.MarkupMaster{MarkupMasterId: *req.MarkupMasterId}).
			Count(&countMarkupMaster).Error; err != nil {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Database error on Markup Master",
				Err:        err,
			}
		}
		if countMarkupMaster == 0 {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Markup Master not found",
			}
		}
	}
	return false, masteritempayloads.ItemSaveResponse{}, nil
}

func checkDiscountExistence(tx *gorm.DB, req masteritempayloads.ItemRequest, response masteritempayloads.ItemSaveResponse) (bool, masteritempayloads.ItemSaveResponse, *exceptions.BaseErrorResponse) {
	if req.DiscountId != nil {
		var countDiscount int64
		if err := tx.Model(&masteritementities.Discount{}).
			Where(masteritementities.Discount{DiscountCodeId: *req.DiscountId}).
			Count(&countDiscount).Error; err != nil {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Database error on Discount",
				Err:        err,
			}
		}
		if countDiscount == 0 {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Discount not found",
			}
		}
	}
	return false, masteritempayloads.ItemSaveResponse{}, nil
}

func checkUomExistence(tx *gorm.DB, req masteritempayloads.ItemRequest, response masteritempayloads.ItemSaveResponse) (bool, masteritempayloads.ItemSaveResponse, *exceptions.BaseErrorResponse) {
	var countUom int64
	if req.UnitOfMeasurementSellingId != nil {
		countUom = 0
		if err := tx.Model(&masteritementities.Uom{}).
			Where(masteritementities.Uom{UomId: *req.UnitOfMeasurementSellingId}).
			Count(&countUom).Error; err != nil {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Database error on Uom Selling",
				Err:        err,
			}
		}
		if countUom == 0 {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Uom Selling not found",
			}
		}
	}

	if req.UnitOfMeasurementPurchaseId != nil {
		countUom = 0
		if err := tx.Model(&masteritementities.Uom{}).
			Where(masteritementities.Uom{UomId: *req.UnitOfMeasurementPurchaseId}).
			Count(&countUom).Error; err != nil {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Database error on Uom Purchase",
				Err:        err,
			}
		}
		if countUom == 0 {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Uom Purchase not found",
			}
		}
	}

	if req.UnitOfMeasurementStockId != nil {
		countUom = 0
		if err := tx.Model(&masteritementities.Uom{}).
			Where(masteritementities.Uom{UomId: *req.UnitOfMeasurementStockId}).
			Count(&countUom).Error; err != nil {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Database error on Uom Stock",
				Err:        err,
			}
		}
		if countUom == 0 {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Uom Stock not found",
			}
		}
	}

	if req.DimensionUnitOfMeasurementId != nil {
		countUom = 0
		if err := tx.Model(&masteritementities.Uom{}).
			Where(masteritementities.Uom{UomId: *req.DimensionUnitOfMeasurementId}).
			Count(&countUom).Error; err != nil {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Database error on Dimension Uom",
				Err:        err,
			}
		}
		if countUom == 0 {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Dimension Uom not found",
			}
		}
	}

	return false, masteritempayloads.ItemSaveResponse{}, nil
}

func checkUomTypeExistence(tx *gorm.DB, req masteritempayloads.ItemRequest, response masteritempayloads.ItemSaveResponse) (bool, masteritempayloads.ItemSaveResponse, *exceptions.BaseErrorResponse) {
	if req.UnitOfMeasurementTypeId != nil {
		var countUomType int64
		if err := tx.Model(&masteritementities.UomType{}).
			Where(masteritementities.UomType{UomTypeId: *req.UnitOfMeasurementTypeId}).
			Count(&countUomType).Error; err != nil {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Database error on UomType",
				Err:        err,
			}
		}
		if countUomType == 0 {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Uom Type not found",
			}
		}
	}
	return false, masteritempayloads.ItemSaveResponse{}, nil
}

func checkItemGroupExistence(tx *gorm.DB, req masteritempayloads.ItemRequest, response masteritempayloads.ItemSaveResponse) (bool, masteritempayloads.ItemSaveResponse, *exceptions.BaseErrorResponse) {
	var countGroup int64
	if err := tx.Model(&masteritementities.ItemGroup{}).
		Where(masteritementities.ItemGroup{ItemGroupId: req.ItemGroupId}).
		Count(&countGroup).Error; err != nil {
		return true, response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Database error on ItemGroup",
			Err:        err,
		}
	}
	if countGroup == 0 {
		return true, response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Item group not found",
		}
	}
	return false, masteritempayloads.ItemSaveResponse{}, nil
}

func checkItemClasExistence(tx *gorm.DB, req masteritempayloads.ItemRequest, response masteritempayloads.ItemSaveResponse) (bool, masteritempayloads.ItemSaveResponse, *exceptions.BaseErrorResponse) {
	var countClass int64
	if err := tx.Model(&masteritementities.ItemClass{}).
		Where(masteritementities.ItemClass{ItemClassId: req.ItemClassId}).
		Count(&countClass).Error; err != nil {
		return true, response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Database error on ItemClass",
			Err:        err,
		}
	}
	if countClass == 0 {
		return true, response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Item class not found",
		}
	}
	return false, masteritempayloads.ItemSaveResponse{}, nil
}

func checkIfItemLevelExists(tx *gorm.DB, req masteritempayloads.ItemRequest, response masteritempayloads.ItemSaveResponse) (bool, masteritempayloads.ItemSaveResponse, *exceptions.BaseErrorResponse) {
	var countLevel1 int64
	if err := tx.Model(&masteritementities.ItemLevel1{}).
		Where(masteritementities.ItemLevel1{
			ItemLevel1Id: *req.ItemLevel1Id,
		}).
		Count(&countLevel1).Error; err != nil {
		return true, response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Database error on ItemLevel1",
			Err:        err,
		}
	}
	if countLevel1 == 0 {
		return true, response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Item level 1 not found",
		}
	}

	if req.ItemLevel2Id != nil {
		var countLevel2 int64
		if err := tx.Model(&masteritementities.ItemLevel2{}).
			Where(masteritementities.ItemLevel2{ItemLevel2Id: *req.ItemLevel2Id}).
			Count(&countLevel2).Error; err != nil {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Database error on ItemLevel2",
				Err:        err,
			}
		}
		if countLevel2 == 0 {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Item level 2 not found",
			}
		}
	}

	if req.ItemLevel3Id != nil {
		var countLevel3 int64
		if err := tx.Model(&masteritementities.ItemLevel3{}).
			Where(masteritementities.ItemLevel3{ItemLevel3Id: *req.ItemLevel3Id}).
			Count(&countLevel3).Error; err != nil {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Database error on ItemLevel3",
				Err:        err,
			}
		}
		if countLevel3 == 0 {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Item level 3 not found",
			}
		}
	}

	if req.ItemLevel4Id != nil {
		var countLevel4 int64
		if err := tx.Model(&masteritementities.ItemLevel4{}).
			Where(masteritementities.ItemLevel4{ItemLevel4Id: *req.ItemLevel4Id}).
			Count(&countLevel4).Error; err != nil {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Database error on ItemLevel4",
				Err:        err,
			}
		}
		if countLevel4 == 0 {
			return true, response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Item level 4 not found",
			}
		}
	}
	return false, masteritempayloads.ItemSaveResponse{}, nil
}

func checkItemTypeExistence(tx *gorm.DB, req masteritempayloads.ItemRequest, response masteritempayloads.ItemSaveResponse) (bool, masteritempayloads.ItemSaveResponse, *exceptions.BaseErrorResponse) {
	var count int64
	if err := tx.Model(&masteritementities.ItemType{}).
		Where("item_type_id = ?", req.ItemTypeId).
		Count(&count).Error; err != nil {

		return true, response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Database error",
			Err:        err,
		}
	}

	if count == 0 {
		return true, response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Item type not found",
		}
	}
	return false, masteritempayloads.ItemSaveResponse{}, nil
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

func (r *ItemRepositoryImpl) GetPrincipalBrandDropdown(tx *gorm.DB) ([]masteritempayloads.PrincipalBrandDropdownResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PrincipalBrandParent{}
	payloads := []masteritempayloads.PrincipalBrandDropdownResponse{}
	err := tx.Model(&entities).Scan(&payloads).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	return payloads, nil
}

func (r *ItemRepositoryImpl) GetPrincipalCatalog(tx *gorm.DB) ([]masteritempayloads.GetPrincipalCatalog, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PrincipalBrandParent{}
	payloads := []masteritempayloads.GetPrincipalCatalog{}

	err := tx.Model(&entities).
		Select("mpc.*").
		Joins("INNER JOIN mtr_principal_catalog mpc ON mpc.principal_catalog_id = mtr_principal_brand_parent.principal_catalog_id").
		Scan(&payloads).Error
	if err != nil {
		return payloads, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	return payloads, nil
}

func (r *ItemRepositoryImpl) GetPrincipalBrandParent(tx *gorm.DB, id int) ([]masteritempayloads.PrincipalBrandDropdownDescription, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PrincipalBrandParent{}
	payloads := []masteritempayloads.PrincipalBrandDropdownDescription{}
	err := tx.Model(&entities).
		Joins("INNER JOIN mtr_principal_catalog mpc ON mpc.principal_catalog_id = mtr_principal_brand_parent.principal_catalog_id").
		Where("mpc.principal_catalog_id = ?", id).
		Scan(&payloads).Error
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
