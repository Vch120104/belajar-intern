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
			mtr_item.item_code,
			mtr_item.item_name,
			mtr_item.item_class_id,
			mtr_item.item_type,
			mtr_item.item_level_1,
			mtr_item.item_level_2,
			mtr_item.item_level_3,
			mtr_item.item_level_4`).
		Joins("INNER JOIN mtr_item_class ic ON ic.item_class_id = mtr_item.item_class_id")

	whereQuery := utils.ApplyFilterExact(baseModelQuery, filterCondition)

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

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	// Handle item_id filter
	if len(itemIDs) > 0 && itemIDs[0] != "" {
		whereQuery = whereQuery.Where("mtr_item.item_id IN (?)", itemIDs)
	}

	// Handle supplier_id filter
	if len(supplierIDs) > 0 && supplierIDs[0] != "" {
		whereQuery = whereQuery.Where("mtr_item.supplier_id IN (?)", supplierIDs)
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
		responseMap := map[string]interface{}{
			"is_active":     response.IsActive,
			"item_id":       response.ItemId,
			"item_code":     response.ItemCode,
			"item_name":     response.ItemName,
			"item_group_id": response.ItemGroupId,
			"item_class_id": response.ItemClassId,
			"item_type":     response.ItemType,
			"supplier_id":   response.SupplierId,
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
			"is_active":     response.IsActive,
			"item_id":       response.ItemId,
			"item_code":     response.ItemCode,
			"item_name":     response.ItemName,
			"item_group_id": response.ItemGroupId,
			"item_class_id": response.ItemClassId,
			"item_type":     response.ItemType,
			"supplier_id":   response.SupplierId,
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

	rows, err := tx.Model(&entities).Select("u.*,mtr_item.*").
		Where(masteritementities.Item{
			ItemId: Id,
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

func (r *ItemRepositoryImpl) SaveItem(tx *gorm.DB, req masteritempayloads.ItemRequest) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Item{
		ItemId:                       req.ItemId,
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
		IsSellable:                   req.IsSellable,
		IsAffiliatedTrx:              req.IsAffiliatedTrx,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	model := masteritementities.Item{}

	err = tx.Model(&model).Where(masteritementities.Item{ItemCode: req.ItemCode}).First(&model).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	atpmResponse := masteritempayloads.AtpmOrderTypeResponse{}

	atpmOrderTypeUrl := config.EnvConfigs.GeneralServiceUrl + "/atpm-order-type/" + strconv.Itoa(req.SourceTypeId)

	if err := utils.Get(atpmOrderTypeUrl, &atpmResponse, nil); err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	uomTypeModel := masteritementities.UomType{}

	err = tx.Model(&uomTypeModel).Where(masteritementities.UomType{UomTypeId: req.UnitOfMeasurementTypeId}).First(&uomTypeModel).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
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
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
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
	// Define a slice to hold Item Detail responses
	var responses []masteritempayloads.ItemDetailRequest
	var brandpayload []masterpayloads.BrandResponse
	var modelpayloads []masterpayloads.UnitModelResponse
	var variantpayloads []masterpayloads.GetVariantResponse
	// Filter internal service conditions

	// Apply internal service filter conditions
	tableStruct := masteritempayloads.ItemDetailRequest{}
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	whereQuery := utils.ApplyFilterExact(joinTable, filterCondition)

	// Fetch data from database
	err := whereQuery.Find(&responses).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to fetch data from database: %w", err),
		}
	}

	// Check if responses are empty
	if len(responses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("no data found"),
		}
	}

	errurlbrand := utils.Get(config.EnvConfigs.SalesServiceUrl+"/unit-brand?page=0&limit=1000000", &brandpayload, nil)
	if errurlbrand != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("no brand found"),
		}
	}
	Joineddata1, errdf := utils.DataFrameInnerJoin(responses, brandpayload, "BrandId")

	if errdf != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	errurlmodel := utils.Get(config.EnvConfigs.SalesServiceUrl+"unit-model?page=0&limit=1000000", &modelpayloads, nil)
	if errurlmodel != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
		}
	}
	joineddata2, errdf := utils.DataFrameInnerJoin(Joineddata1, modelpayloads, "ModelId")

	if errdf != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	errurlvariant := utils.Get(config.EnvConfigs.SalesServiceUrl+"unit-variant?page=0&limit=1000000", &variantpayloads, nil)
	if errurlvariant != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
		}
	}
	joineddata3, errdf := utils.DataFrameInnerJoin(joineddata2, variantpayloads, "VariantId")

	if errdf != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	// Define a slice to hold map responses
	var mapResponses []map[string]interface{}

	// Iterate over responses and convert them to maps
	for _, response := range joineddata3 {
		responseMap := map[string]interface{}{
			"is_active":           response["IsActive"],
			"item_detail_id":      response["ItemDetailId"],
			"item_id":             response["ItemId"],
			"brand_id":            response["BrandId"],
			"brand_name":          response["BrandName"],
			"mileage_every":       response["MileageEvery"],
			"model_id":            response["ModelId"],
			"model_code":          response["ModelCode"],
			"model_description":   response["ModelDescription"],
			"return_every":        response["ReturnEvery"],
			"variant_id":          response["VariantId"],
			"variant_code":        response["VariantCode"],
			"variant_description": response["VariantDescription"],
		}
		mapResponses = append(mapResponses, responseMap)
	}

	// Paginate the response data
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

func (r *ItemRepositoryImpl) AddItemDetail(tx *gorm.DB, ItemId int, req masteritempayloads.ItemDetailRequest) *exceptions.BaseErrorResponse {
	entities := masteritementities.ItemDetail{
		ItemId:       ItemId,
		BrandId:      req.BrandId,
		ModelId:      req.ModelId,
		VariantId:    req.VariantId,
		MileageEvery: req.MileageEvery,
		ReturnEvery:  req.ReturnEvery,
		IsActive:     req.IsActive,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return nil
}

func (r *ItemRepositoryImpl) DeleteItemDetail(tx *gorm.DB, ItemId int, ItemDetailId int) *exceptions.BaseErrorResponse {
	var entities masteritementities.ItemDetail

	result := tx.Model(&entities).
		Where("item_id = ? AND item_detail_id = ?", ItemId, ItemDetailId).
		Delete(&entities)

	if result.Error != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return nil
}

func (r *ItemRepositoryImpl) UpdateItem(tx *gorm.DB, ItemId int, req masteritempayloads.ItemUpdateRequest) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.Item

	result := tx.Model(&entities).Where("item_id=?", ItemId).First(&entities).Updates(req)
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

func (r *ItemRepositoryImpl) UpdateItemDetail(tx *gorm.DB, ItemId int, req masteritempayloads.ItemDetailUpdateRequest) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.ItemDetail

	result := tx.Model(&entities).Where("Item_detail_id=?", ItemId).Updates(req)
	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        result.Error,
		}
	}
	return true, nil
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

func (r *ItemRepositoryImpl) GetCatalogCode(tx *gorm.DB, gmmCatalogCode int) (masteritempayloads.GetCatalogCode, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Item{}
	payloads := masteritempayloads.GetCatalogCode{}

	err := tx.Model(&entities).
		Select("gmm_catalog_code").
		Where("gmm_catalog_code = ?", gmmCatalogCode).
		Scan(&payloads).Error

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
		PrincipalBrandParentCode: code,
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
