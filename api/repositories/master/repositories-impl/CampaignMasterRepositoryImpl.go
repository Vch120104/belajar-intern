package masterrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type CampaignMasterRepositoryImpl struct {
	lookupRepo masterrepository.LookupRepository
}

func StartCampaignMasterRepositoryImpl() masterrepository.CampaignMasterRepository {
	lookupRepo := StartLookupRepositoryImpl()
	return &CampaignMasterRepositoryImpl{
		lookupRepo: lookupRepo,
	}
}

func (r *CampaignMasterRepositoryImpl) PostCampaignMaster(tx *gorm.DB, req masterpayloads.CampaignMasterPost) (masterentities.CampaignMaster, *exceptions.BaseErrorResponse) {
	var existingCampaign masterentities.CampaignMaster

	err := tx.Model(&existingCampaign).
		Where("campaign_code = ? AND campaign_id != ?", req.CampaignCode, req.CampaignId).
		First(&existingCampaign).Error

	if err == nil {
		return masterentities.CampaignMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        errors.New("campaign code already exists"),
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return masterentities.CampaignMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	entity := masterentities.CampaignMaster{
		CampaignId:         req.CampaignId,
		CampaignCode:       req.CampaignCode,
		CampaignName:       req.CampaignName,
		BrandId:            req.BrandId,
		ModelId:            req.ModelId,
		CampaignPeriodFrom: req.CampaignPeriodFrom,
		CampaignPeriodTo:   req.CampaignPeriodTo,
		Remark:             req.Remark,
		AppointmentOnly:    req.AppointmentOnly,
		TaxId:              req.TaxId,
		CompanyId:          req.CompanyId,
		WarehouseGroupId:   req.WarehouseGroupId,
	}

	if err := tx.Save(&entity).Error; err != nil {
		return masterentities.CampaignMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error saving campaign master",
			Err:        err,
		}
	}

	return entity, nil
}

func (r *CampaignMasterRepositoryImpl) PostCampaignMasterDetailFromHistory(tx *gorm.DB, id int, idhead int) (int, *exceptions.BaseErrorResponse) {
	var entityitem []masterentities.CampaignMasterDetail

	result := tx.Model(&entityitem).Where("campaign_id = ?", id).Find(&entityitem)
	if result.Error != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        result.Error,
		}
	}
	for _, entities := range entityitem {
		newEntity := masterentities.CampaignMasterDetail{
			IsActive:         entities.IsActive,
			CampaignDetailId: 0,
			CampaignId:       idhead,
			LineTypeId:       entities.LineTypeId,
			Quantity:         entities.Quantity,
			ItemOperationId:  entities.ItemOperationId,
			ShareBillTo:      entities.ShareBillTo,
			DiscountPercent:  entities.DiscountPercent,
			SharePercent:     entities.SharePercent,
		}
		err := tx.Create(&newEntity).Error

		if err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}
	results := r.UpdateTotalCampaignMaster(tx, idhead)
	if !results {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}
	return idhead, nil
}

func (r *CampaignMasterRepositoryImpl) PostCampaignDetailMaster(tx *gorm.DB, req masterpayloads.CampaignMasterDetailPayloads, id int) (masterentities.CampaignMasterDetail, *exceptions.BaseErrorResponse) {
	var campaign masterentities.CampaignMaster
	var lastPrice float64

	// Validasi SharePercent <= DiscountPercent
	if req.SharePercent > req.DiscountPercent {
		return masterentities.CampaignMasterDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Share percent must not be higher than discount percent",
			Err:        errors.New("share percent must not be higher than discount percent"),
		}
	}

	// Validasi OperationItemId menggunakan API eksternal
	if err := r.validateOperationItemId(req.LineTypeId, req.OperationItemId); err != nil {
		return masterentities.CampaignMasterDetail{}, err
	}

	if req.LineTypeId != 9 && req.LineTypeId != 0 {
		// Non-operation line type
		if err := getLastPriceForNonOperation(tx, req.OperationItemId, &lastPrice); err != nil {
			return masterentities.CampaignMasterDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Operation item price not found",
				Err:        err,
			}
		}
	} else {
		// Operation line type
		if err := tx.Model(&campaign).Where("campaign_id = ?", req.CampaignId).First(&campaign).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return masterentities.CampaignMasterDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Campaign not found",
					Err:        err,
				}
			}
			return masterentities.CampaignMasterDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error fetching campaign",
				Err:        err,
			}
		}

		if err := getLastPriceForOperation(tx, campaign, &lastPrice); err != nil {
			return masterentities.CampaignMasterDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Operation item price not found",
				Err:        err,
			}
		}
	}

	// Create new campaign detail
	newDetail := masterentities.CampaignMasterDetail{
		CampaignId:      id,
		LineTypeId:      req.LineTypeId,
		Quantity:        req.Quantity,
		ItemOperationId: req.OperationItemId,
		ShareBillTo:     req.ShareBillTo,
		DiscountPercent: req.DiscountPercent,
		SharePercent:    req.SharePercent,
		Price:           lastPrice,
	}

	// Save new detail
	if err := tx.Save(&newDetail).Error; err != nil {
		return masterentities.CampaignMasterDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Update campaign total
	if !r.UpdateTotalCampaignMaster(tx, id) {
		return masterentities.CampaignMasterDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update total campaign master",
			Err:        errors.New("failed to update total campaign master"),
		}
	}

	return newDetail, nil
}

// validasi untuk memastikan OperationItemId valid dengan memanggil API eksternal
func (r *CampaignMasterRepositoryImpl) validateOperationItemId(lineTypeId, operationItemId int) *exceptions.BaseErrorResponse {
	// Membuat URL endpoint API eksternal
	url := config.EnvConfigs.AfterSalesServiceUrl + "lookup/item-opr-code/" + strconv.Itoa(lineTypeId) + "/by-id/" + strconv.Itoa(operationItemId)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Error calling external service: %v", err),
			Err:        err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Invalid OperationItemId from external service",
			Err:        errors.New("invalid OperationItemId from external service"),
		}
	}

	// Menangani respons JSON dari API
	var responseData struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
		Data       struct {
			OperationId int `json:"operation_id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Error decoding response: %v", err),
			Err:        err,
		}
	}

	if responseData.Data.OperationId == 0 {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "OperationItemId is invalid according to the external service",
			Err:        errors.New("OperationItemId is invalid according to the external service"),
		}
	}

	return nil
}

// Helper untuk mendapatkan harga terakhir untuk Non-Operation
func getLastPriceForNonOperation(tx *gorm.DB, operationItemId int, lastPrice *float64) error {
	return tx.Select("mtr_item_price_list.price_list_amount").
		Table("mtr_item_price_list").
		Joins("JOIN mtr_item ON mtr_item.item_id = mtr_item_price_list.item_id").
		Joins("JOIN mtr_item_operation ON mtr_item.item_id = mtr_item_operation.item_operation_model_mapping_id").
		Where("item_operation_id = ?", operationItemId).
		Scan(lastPrice).Error
}

// Helper untuk mendapatkan harga terakhir untuk Operation
func getLastPriceForOperation(tx *gorm.DB, campaign masterentities.CampaignMaster, lastPrice *float64) error {
	return tx.Select("mtr_labour_selling_price_detail.selling_price").
		Table("mtr_labour_selling_price_detail").
		Joins("JOIN mtr_labour_selling_price ON mtr_labour_selling_price.labour_selling_price_id = mtr_labour_selling_price_detail.labour_selling_price_id").
		Where("mtr_labour_selling_price.brand_id = ?", campaign.BrandId).
		Where("mtr_labour_selling_price_detail.model_id = ?", campaign.ModelId).
		Where("mtr_labour_selling_price.company_id = ?", campaign.CompanyId).
		Where("mtr_labour_selling_price.effective_date < ?", time.Now()).
		Scan(lastPrice).Error
}

func (r *CampaignMasterRepositoryImpl) PostCampaignMasterDetailFromPackage(tx *gorm.DB, req masterpayloads.CampaignMasterDetailPostFromPackageRequest) (masterentities.CampaignMasterDetail, *exceptions.BaseErrorResponse) {
	campaignDetail := masterentities.CampaignMasterDetail{}
	response := masterentities.CampaignMasterDetail{}

	// Fetch Company Reference
	companyReference, err := generalserviceapiutils.GetCompanyReferenceById(req.CompanyId)
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching company reference",
			Err:        err.Err,
		}
	}
	currencyId := companyReference.CurrencyId

	// Fetch Job Type Campaign
	jobType, err := generalserviceapiutils.GetJobTransactionTypeByCode("CP")
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching job type campaign",
			Err:        err.Err,
		}
	}
	jobTypeCampaignId := jobType.JobTypeId

	// Fetch Work Order Transaction Type
	transactionType, err := generalserviceapiutils.GetWoTransactionTypeByCode("G")
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching work order transaction",
			Err:        err.Err,
		}
	}
	billCode := transactionType.WoTransactionTypeId

	// Check if details already exist
	var totalRows int64
	if err := tx.Model(&campaignDetail).
		Where("campaign_id = ? AND package_id = ?", req.CampaignId, req.PackageId).
		Count(&totalRows).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if totalRows == 0 {
		var packageDetails []masterentities.PackageMasterDetail
		if err := tx.Model(&masterentities.PackageMasterDetail{}).
			Where("package_id = ? AND is_active = ?", req.PackageId, true).
			Find(&packageDetails).Error; err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Fetch Warehouse Group
		var warehouseGroup string
		if err := tx.Model(&masterwarehouseentities.WarehouseGroup{}).
			Select("warehouse_group_type_code").
			Where("warehouse_group_mapping_id = ?", 1). //Warehouse Group For Campaign
			First(&warehouseGroup).Error; err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		for _, detail := range packageDetails {
			// Get Item Price Code
			itemPriceCode, err := r.lookupRepo.GetOprItemPrice(tx, detail.LineTypeId, req.CompanyId, detail.ItemOperationId, req.BrandId, req.ModelId, jobTypeCampaignId, 0, currencyId, billCode, warehouseGroup)
			if err != nil {
				return response, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err.Err,
				}
			}
			if itemPriceCode == 0 {
				return response, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Err:        errors.New("operation item price not found, please set the price first"),
				}
			}
		}

		// Fetch Single Package Detail
		packageDetail := masterentities.PackageMasterDetail{}
		if err := tx.Model(&masterentities.PackageMasterDetail{}).
			Where("package_id = ? AND is_active = ?", req.PackageId, true).
			First(&packageDetail).Error; err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Get Item Price Code
		itemPriceCode, err := r.lookupRepo.GetOprItemPrice(tx, packageDetail.LineTypeId, req.CompanyId, packageDetail.ItemOperationId, req.BrandId, req.ModelId, jobTypeCampaignId, 0, currencyId, billCode, warehouseGroup)
		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err.Err,
			}
		}

		campaignDetail = masterentities.CampaignMasterDetail{
			IsActive:        true,
			CampaignId:      req.CampaignId,
			LineTypeId:      packageDetail.LineTypeId,
			ItemOperationId: packageDetail.ItemOperationId,
			Quantity:        packageDetail.FrtQuantity,
			ShareBillTo:     "",
			DiscountPercent: 0,
			SharePercent:    0,
			PackageId:       req.PackageId,
			Price:           itemPriceCode,
		}

		if err := tx.Save(&campaignDetail).Error; err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		response = campaignDetail
	}

	return response, nil
}

func (r *CampaignMasterRepositoryImpl) ChangeStatusCampaignMaster(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masterentities.CampaignMaster
	result := tx.
		Model(&entities).
		Where("campaign_id = ?", id).
		First(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "CampaignMaster not found",
				Err:        result.Error,
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	entities.IsActive = !entities.IsActive

	if err := tx.Save(&entities).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update CampaignMaster status",
			Err:        err,
		}
	}

	return entities.IsActive, nil
}

func (r *CampaignMasterRepositoryImpl) DeactivateCampaignMasterDetail(tx *gorm.DB, ids string) (bool, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(ids, ",")

	result := tx.Model(&masterentities.CampaignMasterDetail{}).
		Where("campaign_detail_id IN (?)", idSlice).
		Update("is_active", false)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	var campaignIDs []int
	tx.Model(&masterentities.CampaignMasterDetail{}).
		Where("campaign_detail_id IN (?)", idSlice).
		Pluck("campaign_id", &campaignIDs)

	for _, campaignID := range campaignIDs {
		if !r.UpdateTotalCampaignMaster(tx, campaignID) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update total campaign master",
			}
		}
	}

	return true, nil
}

func (r *CampaignMasterRepositoryImpl) ActivateCampaignMasterDetail(tx *gorm.DB, ids string) (bool, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(ids, ",")

	result := tx.Model(&masterentities.CampaignMasterDetail{}).
		Where("campaign_detail_id IN (?)", idSlice).
		Update("is_active", true)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	var campaignIDs []int
	tx.Model(&masterentities.CampaignMasterDetail{}).
		Where("campaign_detail_id IN (?)", idSlice).
		Pluck("campaign_id", &campaignIDs)

	for _, campaignID := range campaignIDs {
		if !r.UpdateTotalCampaignMaster(tx, campaignID) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update total campaign master",
			}
		}
	}

	return true, nil
}

func (r *CampaignMasterRepositoryImpl) GetByIdCampaignMaster(tx *gorm.DB, id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	entities := masterentities.CampaignMaster{}
	payloads := masterpayloads.CampaignMasterResponse{}
	err := tx.Model(&entities).Where(masterentities.CampaignMaster{
		CampaignId: id,
	}).First(&payloads).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Fetch brand details using the salesserviceapiutils
	brandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(payloads.BrandId)
	if brandErr != nil {
		return nil, brandErr
	}

	// Fetch model details using the salesserviceapiutils
	modelResponse, modelErr := salesserviceapiutils.GetUnitModelById(payloads.ModelId)
	if modelErr != nil {
		return nil, modelErr
	}

	result := map[string]interface{}{
		"is_active":            payloads.IsActive,
		"campaign_id":          payloads.CampaignId,
		"campaign_code":        payloads.CampaignCode,
		"campaign_name":        payloads.CampaignName,
		"brand_id":             payloads.BrandId,
		"brand_code":           brandResponse.BrandCode,
		"brand_name":           brandResponse.BrandName,
		"model_id":             payloads.ModelId,
		"model_code":           modelResponse.ModelCode,
		"model_description":    modelResponse.ModelName,
		"campaign_period_from": payloads.CampaignPeriodFrom,
		"campaign_period_to":   payloads.CampaignPeriodTo,
		"remark":               payloads.Remark,
		"appointment_only":     payloads.AppointmentOnly,
		"total":                payloads.Total,
		"total_vat":            payloads.TotalVat,
		"total_after_vat":      payloads.TotalAfterVat,
		"company_id":           payloads.CompanyId,
	}

	return result, nil
}

func (r *CampaignMasterRepositoryImpl) GetByIdCampaignMasterDetail(tx *gorm.DB, id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	var entities masterentities.CampaignMasterDetail
	var payloads masterpayloads.CampaignMasterDetailGetPayloads
	var item masteritempayloads.BomItemNameResponse
	var operation masterpayloads.Operation
	err := tx.Model(&entities).
		Where(masterentities.CampaignMasterDetail{CampaignDetailId: id}).
		First(&payloads).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if entities.LineTypeId != 9 && entities.LineTypeId != 0 {
		err = tx.Select("mtr_item.item_name,mtr_item.item_code").Table("mtr_campaign_master_detail").
			Joins("join mtr_item_operation on mtr_item_operation.item_operation_id=mtr_campaign_master_detail.item_operation_id").
			Joins("join mtr_item on mtr_item.item_id=mtr_item_operation.item_operation_model_mapping_id").
			Where("mtr_campaign_master_detail.campaign_detail_id=?", id).
			Scan(&item).
			Error
	} else {
		err = tx.Select("operation_code.operation_name,operation_code.operation_code").Where("campaign_detail_id=?", id).
			Joins("join mtr_item_operation on mtr_item_operation.item_operation_id = mtr_campaign_master_detail.item_operation_id").
			Joins("JOIN mtr_operation_model_mapping ON mtr_operation_model_mapping.operation_model_mapping_id=mtr_item_operation.item_operation_model_mapping_id").
			Joins("join mtr_operation_code on mtr_operation_code.operation_id=mtr_operation_model_mapping.operation_id").
			Select("mtr_campaign_master_detail.*,mtr_operation_code.operation_code,mtr_operation_code.operation_name").
			Table("mtr_campaign_master_detail").
			Scan(&operation).Error
	}

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	beforeDisc := payloads.Price * payloads.Quantity
	afterDisc := beforeDisc
	if payloads.DiscountPercent > 0 {
		afterDisc = beforeDisc - (beforeDisc * payloads.DiscountPercent / 100)
	}

	response := map[string]interface{}{
		"is_active":          payloads.IsActive,
		"campaign_detail_id": payloads.CampaignDetailId,
		"campaign_id":        payloads.CampaignId,
		"line_type_id":       payloads.LineTypeId,
		"item_operation_id":  payloads.ItemOperationId,
		"frt_quantity":       payloads.Quantity,
		"price":              payloads.Price,
		"discount_percent":   payloads.DiscountPercent,
		"share_percent":      payloads.SharePercent,
		"share_bill_to":      payloads.ShareBillTo,
		"total":              afterDisc,
	}

	if entities.LineTypeId != 9 && entities.LineTypeId != 1 {
		response["item_name"] = item.ItemName
		response["item_code"] = item.ItemCode
	} else {
		response["operation_name"] = operation.OperationName
		response["operation_code"] = operation.OperationCode
	}

	return response, nil
}

func (r *CampaignMasterRepositoryImpl) GetByCodeCampaignMaster(tx *gorm.DB, code string) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	entities := masterentities.CampaignMaster{}
	payloads := masterpayloads.CampaignMasterResponse{}
	err := tx.Model(&entities).Where(masterentities.CampaignMaster{
		CampaignCode: code,
	}).First(&payloads).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Fetch brand details using salesserviceapiutils
	brandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(payloads.BrandId)
	if brandErr != nil {
		return nil, brandErr
	}

	// Fetch model details using salesserviceapiutils
	modelResponse, modelErr := salesserviceapiutils.GetUnitModelById(payloads.ModelId)
	if modelErr != nil {
		return nil, modelErr
	}

	result := map[string]interface{}{
		"is_active":            payloads.IsActive,
		"campaign_id":          payloads.CampaignId,
		"campaign_code":        payloads.CampaignCode,
		"campaign_name":        payloads.CampaignName,
		"brand_id":             payloads.BrandId,
		"brand_code":           brandResponse.BrandCode,
		"brand_name":           brandResponse.BrandName,
		"model_id":             payloads.ModelId,
		"model_code":           modelResponse.ModelCode,
		"model_description":    modelResponse.ModelName,
		"campaign_period_from": payloads.CampaignPeriodFrom,
		"campaign_period_to":   payloads.CampaignPeriodTo,
		"remark":               payloads.Remark,
		"appointment_only":     payloads.AppointmentOnly,
		"total":                payloads.Total,
		"total_vat":            payloads.TotalVat,
		"total_after_vat":      payloads.TotalAfterVat,
		"company_id":           payloads.CompanyId,
	}

	return result, nil
}

func (r *CampaignMasterRepositoryImpl) GetAllCampaignMasterCodeAndName(tx *gorm.DB, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	CampaignMasterResponse := []masterpayloads.GetHistory{}
	CampaignMasterMapping := []masterentities.CampaignMaster{}

	query := tx.Model(masterentities.CampaignMaster{}).Scan(&CampaignMasterResponse)
	err := query.Scopes(pagination.Paginate(&pages, query)).Scan(&CampaignMasterResponse).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(CampaignMasterResponse) == 0 {
		pages.Rows = CampaignMasterMapping
		return pages, nil
	}

	pages.Rows = CampaignMasterResponse

	return pages, nil
}

func (r *CampaignMasterRepositoryImpl) GetAllCampaignMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities masterentities.CampaignMaster
	var responses []masterpayloads.CampaignMasterResponse
	var mapResponses []map[string]interface{}

	var modelDescription, modelCode, campaignPeriodFrom, campaignPeriodTo string
	for _, filter := range filterCondition {
		switch filter.ColumnField {
		case "model_description":
			modelDescription = filter.ColumnValue
		case "model_code":
			modelCode = filter.ColumnValue
		case "campaign_period_from":
			campaignPeriodFrom = filter.ColumnValue
		case "campaign_period_to":
			campaignPeriodTo = filter.ColumnValue
		}
	}

	query := tx.Model(&entities).
		Select("mtr_campaign.*, model.model_code, model.model_description").
		Joins("JOIN dms_microservices_sales_dev.dbo.mtr_unit_model model ON mtr_campaign.model_id = model.model_id")

	if modelDescription != "" {
		query = query.Where("model.model_description LIKE ?", fmt.Sprintf("%%%s%%", modelDescription))
	}
	if modelCode != "" {
		query = query.Where("model.model_code LIKE ?", fmt.Sprintf("%%%s%%", modelCode))
	}
	if campaignPeriodFrom != "" {
		query = query.Where("mtr_campaign.campaign_period_from >= ?", campaignPeriodFrom)
	}
	if campaignPeriodTo != "" {
		query = query.Where("mtr_campaign.campaign_period_to <= ?", campaignPeriodTo)
	}

	whereQuery := utils.ApplyFilter(query, filterCondition)

	var totalRows int64
	err := whereQuery.Model(&entities).Count(&totalRows).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	offset := pages.GetOffset()
	limit := pages.GetLimit()

	err = whereQuery.Offset(offset).Limit(limit).Order("mtr_campaign.campaign_id").Find(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	for _, response := range responses {
		result := map[string]interface{}{
			"campaign_id":          response.CampaignId,
			"campaign_name":        response.CampaignName,
			"campaign_period_from": response.CampaignPeriodFrom,
			"campaign_period_to":   response.CampaignPeriodTo,
			"model_code":           response.ModelCode,
			"model_description":    response.ModelDescription,
			"remark":               response.Remark,
			"appointment_only":     response.AppointmentOnly,
			"total":                response.Total,
			"total_vat":            response.TotalVat,
			"total_after_vat":      response.TotalAfterVat,
			"is_active":            response.IsActive,
			"company_id":           response.CompanyId,
		}

		mapResponses = append(mapResponses, result)
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))
	pages.TotalRows = totalRows
	pages.TotalPages = totalPages

	pages.Rows = mapResponses
	return pages, nil
}

func (r *CampaignMasterRepositoryImpl) GetAllCampaignMasterDetail(tx *gorm.DB, pages pagination.Pagination, id int) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responsedetail []masterpayloads.CampaignMasterDetailGetPayloads
	var item masteritempayloads.BomItemNameResponse
	var operation masterpayloads.Operation
	combinedPayloads := make([]map[string]interface{}, 0)

	// Query Campaign Master Detail filtered by CampaignId
	err := tx.Model(&masterentities.CampaignMasterDetail{}).
		Where("campaign_id = ?", id).
		Scan(&responsedetail).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Fetch Line Type for Operation
	lineTypeOpr, lineTypeError := generalserviceapiutils.GetLineTypeByCode("1")
	if lineTypeError != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching line type operation data",
			Err:        lineTypeError.Err,
		}
	}

	for _, detail := range responsedetail {
		var packageCode string
		if detail.PackageId != 0 {
			err := tx.Table("mtr_package").Select("package_code").
				Where("package_id = ?", detail.PackageId).
				Scan(&packageCode).Error
			if err != nil {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Err:        err,
				}
			}
		}

		if detail.LineTypeId != lineTypeOpr.LineTypeId && detail.LineTypeId != 0 {
			err = tx.Table("mtr_item").
				Select("item_name, item_code").
				Joins("JOIN mtr_item_operation ON mtr_item.item_id = mtr_item_operation.item_operation_model_mapping_id").
				Joins("JOIN mtr_campaign_master_detail ON mtr_campaign_master_detail.item_operation_id = mtr_item_operation.item_operation_id").
				Where("mtr_campaign_master_detail.campaign_detail_id = ?", detail.CampaignDetailId).
				Scan(&item).Error
		} else {
			err = tx.Table("mtr_operation_code").
				Select("operation_name, operation_code").
				Joins("JOIN mtr_operation_model_mapping ON mtr_operation_code.operation_id = mtr_operation_model_mapping.operation_id").
				Joins("JOIN mtr_item_operation ON mtr_operation_model_mapping.operation_model_mapping_id = mtr_item_operation.item_operation_model_mapping_id").
				Joins("JOIN mtr_campaign_master_detail ON mtr_item_operation.item_operation_id = mtr_campaign_master_detail.item_operation_id").
				Where("mtr_campaign_master_detail.campaign_detail_id = ?", detail.CampaignDetailId).
				Scan(&operation).Error
		}

		if err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		// Calculate discounted total
		beforeDisc := detail.Price * detail.Quantity
		afterDisc := beforeDisc
		if detail.DiscountPercent > 0 {
			afterDisc = beforeDisc - (beforeDisc * detail.DiscountPercent / 100)
		}

		response := map[string]interface{}{
			"is_active":          detail.IsActive,
			"campaign_id":        detail.CampaignId,
			"campaign_detail_id": detail.CampaignDetailId,
			"package_code":       packageCode,
			"package_id":         detail.PackageId,
			"line_type_id":       detail.LineTypeId,
			"item_operation_id":  detail.ItemOperationId,
			"frt_quantity":       detail.Quantity,
			"price":              detail.Price,
			"discount_percent":   detail.DiscountPercent,
			"share_percent":      detail.SharePercent,
			"total":              afterDisc,
		}

		if detail.LineTypeId != lineTypeOpr.LineTypeId && detail.LineTypeId != 0 {
			response["item_name"] = item.ItemName
			response["item_code"] = item.ItemCode
		} else {
			response["operation_name"] = operation.OperationName
			response["operation_code"] = operation.OperationCode
		}

		combinedPayloads = append(combinedPayloads, response)
	}

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(combinedPayloads, &pages)
	pages.Rows = dataPaginate
	pages.TotalPages = totalPages
	pages.TotalRows = int64(totalRows)

	return pages, nil
}

func (r *CampaignMasterRepositoryImpl) UpdateCampaignMasterDetail(tx *gorm.DB, id int, req masterpayloads.CampaignMasterDetailPayloads) (int, *exceptions.BaseErrorResponse) {
	var entities masterentities.CampaignMasterDetail

	result := tx.Model(&entities).Where("campaign_detail_id = ?", id).First(&entities).Updates(req)
	if result.Error != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        result.Error,
		}
	}
	return entities.CampaignId, nil
}

func (r *CampaignMasterRepositoryImpl) UpdateTotalCampaignMaster(tx *gorm.DB, id int) bool {
	var headerentity masterentities.CampaignMaster
	var detailentity []masterentities.CampaignMasterDetail
	var totalValue float64

	// Fetch and calculate the total value from CampaignMasterOperationDetail
	result := tx.Model(&detailentity).Where(masterentities.CampaignMasterDetail{
		CampaignId: id,
	}).Find(&detailentity)
	if result.Error != nil {
		return false
	}
	for _, detail := range detailentity {
		if detail.IsActive {
			if detail.DiscountPercent > 0 {
				totalValue += detail.Quantity * detail.Price * (1 - (detail.DiscountPercent / 100))
			} else {
				totalValue += detail.Quantity * detail.Price
			}
		}
	}

	currentTime := time.Now().Truncate(24 * time.Hour)
	convertedCurrentTime := utils.FormatRFC3339(currentTime)
	date, time, err := utils.ConvertDateTimeFormat(convertedCurrentTime)
	if err != nil {
		return false
	}
	taxRateUrl := config.EnvConfigs.FinanceServiceUrl + "tax-fare/detail/tax-percent?tax_service_code=PPN&tax_type_code=PPN&effective_date=" + date + "T" + time + "Z"
	taxRatePayloads := masterpayloads.TaxFarePercentResponse{}
	if err := utils.Get(taxRateUrl, &taxRatePayloads, nil); err != nil {
		return false
	}
	vatRate := taxRatePayloads.TaxPercent
	totalVat := vatRate / 100 * totalValue
	totalAfterVat := totalValue + totalVat

	err = tx.Model(&headerentity).Where(masterentities.CampaignMaster{CampaignId: id}).First(&headerentity).Error
	if err != nil {
		return false
	}
	headerentity.Total = totalValue
	headerentity.TotalVat = totalVat
	headerentity.TotalAfterVat = totalAfterVat

	err = tx.Save(&headerentity).Error

	return err == nil
}

func (r *CampaignMasterRepositoryImpl) GetAllPackageMasterToCopy(tx *gorm.DB, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var packageentities masterentities.PackageMaster
	var payloads []masterpayloads.PackageMasterForCampaignMaster

	BaseModelQuery := tx.Model(&packageentities)

	err := BaseModelQuery.Scopes(pagination.Paginate(&pages, BaseModelQuery)).
		Where("profit_center_id = ?", 13).
		Scan(&payloads).Error

	if len(payloads) == 0 {
		pages.Rows = []masterpayloads.PackageMasterForCampaignMaster{}
		return pages, nil
	}

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.Rows = payloads

	return pages, nil
}

func (r *CampaignMasterRepositoryImpl) SelectFromPackageMaster(tx *gorm.DB, id int, idhead int) (int, *exceptions.BaseErrorResponse) {
	var packagedetail []masterentities.PackageMasterDetail
	var lastprice float64
	var operationpayloads masterpayloads.CampaignMasterDetailGetPayloads
	var entity masterentities.CampaignMaster
	var itemprice masteritementities.ItemPriceList

	err := tx.Model(&packagedetail).Where("Package_id=?", id).Scan(packagedetail).Error
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	for _, pack := range packagedetail {
		if pack.LineTypeId == 5 { //line type operation
			err := tx.Model(&packagedetail).Where("package_id=?", id).
				Joins("join mtr_item_operation on mtr_item_operation.item_operation_id = mtr_package_detail.item_operation_id").
				Joins("JOIN mtr_operation_model_mapping ON mtr_operation_model_mapping.operation_model_mapping_id=mtr_item_operation.tem_operation_id").
				Joins("join mtr_operation_code on mtr_operation_code.operation_id=mtr_operation_model_mapping.operation_id").
				Select("mtr_package_master_detail.*,mtr_operation_code.operation_code,mtr_operation_code.operation_name").
				Scan(&operationpayloads).Error
			if err != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Err:        err,
				}
			}
			err2 := tx.Select("mtr_labour_selling_price_detail.selling_price").
				Table("mtr_labour_selling_price_detail").
				Joins("Join mtr_labour_selling_price on mtr_labour_selling_price.labour_selling_price_id = mtr_labour_selling_price_detail.labour_selling_price_id").
				Where("mtr_labour_selling_price.brand_id =?", entity.BrandId).
				Where("mtr_labour_selling_price_detail.model_id=?", entity.ModelId).
				Where("mtr_labour_selling_price.company_id = ?", entity.CompanyId).
				Where("mtr_labour_selling_price.effective_date < ?", time.Now()).
				Scan(&lastprice).Error
			if err2 != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err2,
				}
			}
			entity := masterentities.CampaignMasterDetail{
				IsActive:        pack.IsActive,
				CampaignId:      idhead,
				LineTypeId:      pack.LineTypeId,
				Quantity:        pack.FrtQuantity,
				ItemOperationId: pack.ItemOperationId,
				ShareBillTo:     "",
				DiscountPercent: 0,
				SharePercent:    0,
				Price:           0,
			}
			err3 := tx.Save(&entity).Error
			if err3 != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Err:        err3,
				}
			}
		} else {
			err2 := tx.Model(&packagedetail).Where(masterentities.PackageMasterDetail{
				PackageId: id,
			}).Joins("join mtr_item_operation on mtr_item_operation.item_operation_id = mtr_package_detail.item_operation_id").
				Joins("JOIN mtr_item ON mtr_item.item_id=mtr_item_operation.item_id").
				Select("mtr_package_master_detail.*,mtr_item.item_code,mtr_item.item_name").Error
			if err2 != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Err:        err2,
				}
			}
			err := tx.Model(&itemprice).Select("mtr_item_price_list.price_list_amount").
				Joins("join mtr_item on mtr_item.item_id=mtr_item_price_list.item_id").
				Joins("join mtr_item_operation on mtr_item_operation.item_id=mtr_item.item_id").
				Where("item_operation_id=?").
				Scan(&lastprice).Error
			if err != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Err:        err,
				}
			}
			entity := masterentities.CampaignMasterDetail{
				IsActive:        pack.IsActive,
				CampaignId:      idhead,
				LineTypeId:      pack.LineTypeId,
				Quantity:        pack.FrtQuantity,
				ItemOperationId: pack.ItemOperationId,
				ShareBillTo:     "",
				DiscountPercent: 0,
				SharePercent:    0,
				Price:           0,
			}
			err3 := tx.Save(&entity).Error
			if err3 != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Err:        err,
				}
			}
		}
	}
	return idhead, nil
}
