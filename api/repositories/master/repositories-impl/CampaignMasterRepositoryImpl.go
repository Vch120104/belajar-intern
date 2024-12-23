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
			LineTypeCode:     entities.LineTypeCode,
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

	var (
		campaign  masterentities.CampaignMaster
		lastPrice float64
	)

	// Validate SharePercent <= DiscountPercent
	if req.SharePercent > req.DiscountPercent {
		return masterentities.CampaignMasterDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Share percent must not be higher than discount percent",
			Err:        errors.New("share percent must not be higher than discount percent"),
		}
	}

	// Validate OperationItemId using external API
	if err := r.validateOperationItemId(req.LineTypeCode, req.OperationItemId); err != nil {
		return masterentities.CampaignMasterDetail{}, err
	}

	// Fetch the last price based on LineTypeCode
	if req.LineTypeCode != "9" && req.LineTypeCode != "0" {
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

	newDetail := masterentities.CampaignMasterDetail{
		CampaignId:      id,
		LineTypeCode:    req.LineTypeCode,
		Quantity:        req.Quantity,
		ItemOperationId: req.OperationItemId,
		ShareBillTo:     req.ShareBillTo,
		DiscountPercent: req.DiscountPercent,
		SharePercent:    req.SharePercent,
		Price:           lastPrice,
	}

	if err := tx.Save(&newDetail).Error; err != nil {
		return masterentities.CampaignMasterDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save campaign detail",
			Err:        err,
		}
	}

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
func (r *CampaignMasterRepositoryImpl) validateOperationItemId(lineTypeStr string, operationItemId int) *exceptions.BaseErrorResponse {
	url := config.EnvConfigs.AfterSalesServiceUrl + "lookup/item-opr-code/" + lineTypeStr + "/by-id/" + strconv.Itoa(operationItemId)
	fmt.Println("Requesting URL:", url)

	// Perform HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Error calling external service: %v", err),
			Err:        err,
		}
	}
	defer resp.Body.Close()

	// Log response status
	fmt.Println("Response Status Code:", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Invalid combination linetype & OperationItemId from external service",
			Err:        errors.New("invalid combination linetype & OperationItemId from external service"),
		}
	}

	// Decode response based on lineTypeId
	switch lineTypeStr {
	case "0":
		var responseData struct {
			StatusCode int    `json:"status_code"`
			Message    string `json:"message"`
			Data       struct {
				Description      string  `json:"description"`
				FRT              float64 `json:"frt"`
				ModelCode        string  `json:"model_code"`
				PackageCode      string  `json:"package_code"`
				PackageID        int     `json:"package_id"`
				PackageName      string  `json:"package_name"`
				Price            int     `json:"price"`
				ProfitCenter     int     `json:"profit_center"`
				ProfitCenterName string  `json:"profit_center_name"`
			} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("Error decoding response: %v", err),
				Err:        err,
			}
		}

		// Validate PackageID
		if responseData.Data.PackageID != operationItemId {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "OperationItemId is invalid for linetype 0",
				Err:        errors.New("OperationItemId is invalid for linetype 0"),
			}
		}

	case "1":
		var responseData struct {
			StatusCode int    `json:"status_code"`
			Message    string `json:"message"`
			Data       struct {
				FrtHour                     int     `json:"frt_hour"`
				OperationCode               string  `json:"operation_code"`
				OperationEntriesCode        *string `json:"operation_entries_code"`
				OperationEntriesDescription *string `json:"operation_entries_description"`
				OperationID                 int     `json:"operation_id"`
				OperationKeyCode            *string `json:"operation_key_code"`
				OperationKeyDescription     *string `json:"operation_key_description"`
				OperationName               string  `json:"operation_name"`
			} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("Error decoding response: %v", err),
				Err:        err,
			}
		}

		// Validate OperationID
		if responseData.Data.OperationID != operationItemId {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "OperationItemId is invalid for linetype 1",
				Err:        errors.New("OperationItemId is invalid for linetype 1"),
			}
		}

	case "2", "3", "4", "5", "6", "7", "8", "9":
		var responseData struct {
			StatusCode int    `json:"status_code"`
			Message    string `json:"message"`
			Data       struct {
				AvailableQty   int     `json:"available_qty"`
				ItemCode       string  `json:"item_code"`
				ItemID         int     `json:"item_id"`
				ItemLevel1     int     `json:"item_level_1"`
				ItemLevel1Code string  `json:"item_level_1_code"`
				ItemLevel2     *int    `json:"item_level_2"`
				ItemLevel2Code *string `json:"item_level_2_code"`
				ItemLevel3     *int    `json:"item_level_3"`
				ItemLevel3Code *string `json:"item_level_3_code"`
				ItemLevel4     *int    `json:"item_level_4"`
				ItemLevel4Code *string `json:"item_level_4_code"`
				ItemName       string  `json:"item_name"`
			} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("Error decoding response: %v", err),
				Err:        err,
			}
		}

		// Validate ItemID
		if responseData.Data.ItemID != operationItemId {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "OperationItemId is invalid for linetype 2-9",
				Err:        errors.New("OperationItemId is invalid for linetype 2-9"),
			}
		}

	default:
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid linetype provided",
			Err:        errors.New("invalid linetype provided"),
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
			Message:    "error checking existing campaign detail",
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
				Message:    "error fetching package details",
				Err:        err,
			}
		}

		// Fetch Warehouse Group
		var warehouseGroupId int
		if err := tx.Model(&masterentities.WarehouseGroupMappingEntities{}).
			Select("warehouse_group_id").
			Where("warehouse_group_type_code = ?", "WHS_GRP_CAMPAIGN"). //Warehouse Group For Campaign
			First(&warehouseGroupId).Error; err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Fetch Warehouse Group
		var warehouseGroup string
		if err := tx.Model(&masterwarehouseentities.WarehouseGroup{}).
			Select("warehouse_group_code").
			Where("warehouse_group_id = ?", warehouseGroupId).
			First(&warehouseGroup).Error; err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		for _, detail := range packageDetails {
			// Get Item Price Code
			itemPriceCode, err := r.lookupRepo.GetOprItemPrice(tx, detail.LineTypeCode, req.CompanyId, detail.ItemOperationId, req.BrandId, req.ModelId, jobTypeCampaignId, 0, currencyId, billCode, warehouseGroup)
			if err != nil {
				return response, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "error fetching item price code",
					Err:        err.Err,
				}
			}
			if itemPriceCode == 0 {
				return response, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "item price code not found",
					Err:        err,
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
				Message:    "error fetching package detail",
				Err:        err,
			}
		}

		// Get Item Price Code
		itemPriceCode, err := r.lookupRepo.GetOprItemPrice(tx, packageDetail.LineTypeCode, req.CompanyId, packageDetail.ItemOperationId, req.BrandId, req.ModelId, jobTypeCampaignId, 0, currencyId, billCode, warehouseGroup)
		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "error fetching item price code",
				Err:        err.Err,
			}
		}

		campaignDetail = masterentities.CampaignMasterDetail{
			IsActive:        true,
			CampaignId:      req.CampaignId,
			LineTypeCode:    packageDetail.LineTypeCode,
			ItemOperationId: packageDetail.ItemOperationId,
			Quantity:        packageDetail.FrtQuantity,
			ShareBillTo:     "",
			DiscountPercent: 0,
			SharePercent:    0,
			PackageId:       req.PackageId,
			Price:           itemPriceCode,
		}

		fmt.Println("Campaign Detail:", campaignDetail)

		if err := tx.Save(&campaignDetail).Error; err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "error saving campaign detail",
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

func (r *CampaignMasterRepositoryImpl) GetByIdCampaignMaster(tx *gorm.DB, id int) (masterpayloads.CampaignMasterResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.CampaignMaster{}
	payloads := masterpayloads.CampaignMasterResponse{}
	err := tx.Model(&entities).
		Select(`is_active, campaign_code, campaign_name, campaign_id, brand_id, model_id, 
            campaign_period_from, campaign_period_to, remark, appointment_only, total, 
            total_vat, total_after_vat, company_id`).
		Where("campaign_id = ?", id).
		First(&payloads).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return masterpayloads.CampaignMasterResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Campaign not found",
				Err:        err,
			}
		}

		return masterpayloads.CampaignMasterResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching campaign",
			Err:        err,
		}
	}

	// Fetch brand details
	brandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(payloads.BrandId)
	if brandErr != nil {
		return masterpayloads.CampaignMasterResponse{}, brandErr
	}

	// Fetch model details
	modelResponse, modelErr := salesserviceapiutils.GetUnitModelById(payloads.ModelId)
	if modelErr != nil {
		return masterpayloads.CampaignMasterResponse{}, modelErr
	}

	// Fetch company details
	companyResponse, companyErr := generalserviceapiutils.GetCompanyDataById(payloads.CompanyId)
	if companyErr != nil {
		return masterpayloads.CampaignMasterResponse{}, companyErr
	}

	result := masterpayloads.CampaignMasterResponse{
		IsActive:           payloads.IsActive,
		CampaignId:         payloads.CampaignId,
		CampaignCode:       payloads.CampaignCode,
		CampaignName:       payloads.CampaignName,
		BrandId:            payloads.BrandId,
		BrandCode:          brandResponse.BrandCode,
		BrandName:          brandResponse.BrandName,
		ModelId:            payloads.ModelId,
		ModelCode:          modelResponse.ModelCode,
		ModelDescription:   modelResponse.ModelName,
		CampaignPeriodFrom: payloads.CampaignPeriodFrom,
		CampaignPeriodTo:   payloads.CampaignPeriodTo,
		Remark:             payloads.Remark,
		AppointmentOnly:    payloads.AppointmentOnly,
		Total:              payloads.Total,
		TotalVat:           payloads.TotalVat,
		TotalAfterVat:      payloads.TotalAfterVat,
		CompanyId:          payloads.CompanyId,
		CompanyCode:        companyResponse.CompanyCode,
		CompanyName:        companyResponse.CompanyName,
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

	// fetch linetype id from linetype code
	linetypeId, linetypeErr := generalserviceapiutils.GetLineTypeByCode(payloads.LineTypeCode)
	if linetypeErr != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: linetypeErr.StatusCode,
			Message:    "Error fetching line type",
			Err:        linetypeErr.Err,
		}
	}

	var response map[string]interface{}
	switch payloads.LineTypeCode {
	case "2", "3", "4", "5", "6", "7":
		err = tx.Select("mtr_item.item_name, mtr_item.item_code").
			Table("mtr_campaign_master_detail").
			Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_line_type ON dms_microservices_general_dev.dbo.mtr_line_type.line_type_code = mtr_campaign_master_detail.line_type_code").
			Joins("INNER JOIN mtr_mapping_item_operation ON mtr_mapping_item_operation.line_type_id = dms_microservices_general_dev.dbo.mtr_line_type.line_type_id AND mtr_mapping_item_operation.item_id <> 0 AND mtr_mapping_item_operation.item_id = mtr_campaign_master_detail.item_operation_id").
			Joins("INNER JOIN mtr_item ON mtr_mapping_item_operation.item_id = mtr_item.item_id").
			Where("mtr_campaign_master_detail.campaign_detail_id = ?", id).
			Scan(&item).Error
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		beforeDisc := payloads.Price * payloads.Quantity
		afterDisc := beforeDisc
		if payloads.DiscountPercent > 0 {
			afterDisc -= beforeDisc * payloads.DiscountPercent / 100
		}

		response = map[string]interface{}{
			"is_active":           payloads.IsActive,
			"campaign_detail_id":  payloads.CampaignDetailId,
			"campaign_id":         payloads.CampaignId,
			"line_type_id":        linetypeId.LineTypeId,
			"line_type_code":      payloads.LineTypeCode,
			"line_type_name":      linetypeId.LineTypeName,
			"item_operation_id":   payloads.ItemOperationId,
			"frt_quantity":        payloads.Quantity,
			"price":               payloads.Price,
			"discount_percent":    payloads.DiscountPercent,
			"share_percent":       payloads.SharePercent,
			"share_bill_to":       payloads.ShareBillTo,
			"total":               afterDisc,
			"operation_item_name": item.ItemName,
			"operation_item_code": item.ItemCode,
		}

	case "1":
		err = tx.Select("mtr_campaign_master_detail.*, mtr_operation_code.operation_code, mtr_operation_code.operation_name").
			Table("mtr_campaign_master_detail").
			Joins("JOIN dms_microservices_general_dev.dbo.mtr_line_type AS linetype ON linetype.line_type_code = mtr_campaign_master_detail.line_type_code").
			Joins("JOIN mtr_mapping_item_operation ON mtr_mapping_item_operation.operation_id = mtr_campaign_master_detail.item_operation_id").
			Joins("JOIN mtr_operation_code ON mtr_operation_code.operation_id = mtr_mapping_item_operation.operation_id").
			Where("mtr_campaign_master_detail.campaign_detail_id = ?", id).
			Scan(&operation).Error
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		beforeDisc := payloads.Price * payloads.Quantity
		afterDisc := beforeDisc
		if payloads.DiscountPercent > 0 {
			afterDisc -= beforeDisc * payloads.DiscountPercent / 100
		}

		response = map[string]interface{}{
			"is_active":           payloads.IsActive,
			"campaign_detail_id":  payloads.CampaignDetailId,
			"campaign_id":         payloads.CampaignId,
			"line_type_id":        linetypeId.LineTypeId,
			"line_type_code":      payloads.LineTypeCode,
			"line_type_name":      linetypeId.LineTypeName,
			"item_operation_id":   payloads.ItemOperationId,
			"frt_quantity":        payloads.Quantity,
			"price":               payloads.Price,
			"discount_percent":    payloads.DiscountPercent,
			"share_percent":       payloads.SharePercent,
			"share_bill_to":       payloads.ShareBillTo,
			"total":               afterDisc,
			"operation_item_name": operation.OperationName,
			"operation_item_code": operation.OperationCode,
		}

	default:
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid LineTypeCode",
			Err:        errors.New("invalid line type code"),
		}
	}

	return response, nil
}

func (r *CampaignMasterRepositoryImpl) GetByCodeCampaignMaster(tx *gorm.DB, code string) (masterpayloads.CampaignMasterResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.CampaignMaster{}
	payloads := masterpayloads.CampaignMasterResponse{}
	err := tx.Model(&entities).
		Select(`is_active, campaign_code, campaign_name, campaign_id, brand_id, model_id, 
            campaign_period_from, campaign_period_to, remark, appointment_only, total, 
            total_vat, total_after_vat, company_id`).
		Where("campaign_code = ?", code).
		First(&payloads).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return masterpayloads.CampaignMasterResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Campaign not found",
				Err:        err,
			}
		}

		return masterpayloads.CampaignMasterResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching campaign",
			Err:        err,
		}
	}

	// Fetch brand details
	brandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(payloads.BrandId)
	if brandErr != nil {
		return masterpayloads.CampaignMasterResponse{}, brandErr
	}

	// Fetch model details
	modelResponse, modelErr := salesserviceapiutils.GetUnitModelById(payloads.ModelId)
	if modelErr != nil {
		return masterpayloads.CampaignMasterResponse{}, modelErr
	}

	// Fetch company details
	companyResponse, companyErr := generalserviceapiutils.GetCompanyDataById(payloads.CompanyId)
	if companyErr != nil {
		return masterpayloads.CampaignMasterResponse{}, companyErr
	}

	result := masterpayloads.CampaignMasterResponse{
		IsActive:           payloads.IsActive,
		CampaignId:         payloads.CampaignId,
		CampaignCode:       payloads.CampaignCode,
		CampaignName:       payloads.CampaignName,
		BrandId:            payloads.BrandId,
		BrandCode:          brandResponse.BrandCode,
		BrandName:          brandResponse.BrandName,
		ModelId:            payloads.ModelId,
		ModelCode:          modelResponse.ModelCode,
		ModelDescription:   modelResponse.ModelName,
		CampaignPeriodFrom: payloads.CampaignPeriodFrom,
		CampaignPeriodTo:   payloads.CampaignPeriodTo,
		Remark:             payloads.Remark,
		AppointmentOnly:    payloads.AppointmentOnly,
		Total:              payloads.Total,
		TotalVat:           payloads.TotalVat,
		TotalAfterVat:      payloads.TotalAfterVat,
		CompanyId:          payloads.CompanyId,
		CompanyCode:        companyResponse.CompanyCode,
		CompanyName:        companyResponse.CompanyName,
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
	var companyId int
	var modelDescription, modelCode, campaignPeriodFrom, campaignPeriodTo string
	manualFilters := []utils.FilterCondition{}
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
		case "company_id":
			companyId, _ = strconv.Atoi(filter.ColumnValue)
		default:
			manualFilters = append(manualFilters, filter)
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

	if companyId != 0 {
		query = query.Where("(mtr_campaign.company_id = ? OR mtr_campaign.company_id = 0)", companyId)
	}

	query = utils.ApplyFilter(query, manualFilters)

	var totalRows int64
	err := query.Model(&entities).Count(&totalRows).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	offset := pages.GetOffset()
	limit := pages.GetLimit()

	err = query.Offset(offset).Limit(limit).Order("mtr_campaign.campaign_id").Find(&responses).Error
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
			"campaign_code":        response.CampaignCode,
			"model_code":           response.ModelCode,
			"model_description":    response.ModelDescription,
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
	combinedPayloads := make([]map[string]interface{}, 0)

	err := tx.Model(&masterentities.CampaignMasterDetail{}).
		Where("campaign_id = ?", id).
		Scan(&responsedetail).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || len(responsedetail) == 0 {
			pages.Rows = []interface{}{}
			pages.TotalPages = 0
			pages.TotalRows = 0
			return pages, nil
		}
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching campaign details",
			Err:        err,
		}
	}

	for _, detail := range responsedetail {

		// fetch linetype id from linetype code
		linetypeId, linetypeErr := generalserviceapiutils.GetLineTypeByCode(detail.LineTypeCode)
		if linetypeErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: linetypeErr.StatusCode,
				Message:    "Error fetching line type",
				Err:        linetypeErr.Err,
			}
		}

		response := map[string]interface{}{
			"is_active":           detail.IsActive,
			"campaign_id":         detail.CampaignId,
			"campaign_detail_id":  detail.CampaignDetailId,
			"package_code":        "",
			"package_id":          detail.PackageId,
			"line_type_id":        linetypeId.LineTypeId,
			"line_type_code":      detail.LineTypeCode,
			"line_type_name":      linetypeId.LineTypeName,
			"item_operation_id":   detail.ItemOperationId,
			"frt_quantity":        detail.Quantity,
			"price":               detail.Price,
			"discount_percent":    detail.DiscountPercent,
			"share_percent":       detail.SharePercent,
			"total":               detail.Price * detail.Quantity,
			"operation_item_name": "",
			"operation_item_code": "",
		}

		switch detail.LineTypeCode {
		case "0": // LineType 0 (Package Data)
			if detail.PackageId != 0 {
				var packageData struct {
					PackageCode string `gorm:"column:package_code"`
				}
				err = tx.Table("mtr_package").
					Select("package_code").
					Where("package_id = ?", detail.PackageId).
					Scan(&packageData).Error

				if err == nil {
					response["package_code"] = packageData.PackageCode
				}
			}

		case "1": // LineType 1 (Operation Data)
			var operationData struct {
				OperationName string `gorm:"column:operation_name"`
				OperationCode string `gorm:"column:operation_code"`
			}
			fmt.Printf("Looking for operation_id: %d\n", detail.ItemOperationId)

			err = tx.Select(" mtr_operation_code.operation_code, mtr_operation_code.operation_name").
				Table("mtr_campaign_master_detail").
				Joins("JOIN dms_microservices_general_dev.dbo.mtr_line_type AS linetype ON linetype.line_type_code = mtr_campaign_master_detail.line_type_code").
				Joins("JOIN mtr_mapping_item_operation ON mtr_mapping_item_operation.operation_id = mtr_campaign_master_detail.item_operation_id").
				Joins("JOIN mtr_operation_code ON mtr_operation_code.operation_id = mtr_mapping_item_operation.operation_id").
				Where("mtr_campaign_master_detail.campaign_detail_id = ?", detail.CampaignDetailId).
				Scan(&operationData).Error

			if err == nil {
				response["operation_item_name"] = operationData.OperationName
				response["operation_item_code"] = operationData.OperationCode
			}

		case "2", "3", "4", "5", "6", "7", "8", "9": // LineType 2-9 (Item Data)
			var itemData struct {
				ItemName string `gorm:"column:item_name"`
				ItemCode string `gorm:"column:item_code"`
			}
			fmt.Printf("Looking for item_id: %d\n", detail.ItemOperationId)
			err = tx.Select("mtr_item.item_name, mtr_item.item_code").
				Table("mtr_campaign_master_detail").
				Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_line_type ON dms_microservices_general_dev.dbo.mtr_line_type.line_type_code = mtr_campaign_master_detail.line_type_code").
				Joins("INNER JOIN mtr_mapping_item_operation ON mtr_mapping_item_operation.line_type_id = dms_microservices_general_dev.dbo.mtr_line_type.line_type_id AND mtr_mapping_item_operation.item_id <> 0 AND mtr_mapping_item_operation.item_id = mtr_campaign_master_detail.item_operation_id").
				Joins("INNER JOIN mtr_item ON mtr_mapping_item_operation.item_id = mtr_item.item_id").
				Where("mtr_campaign_master_detail.campaign_detail_id = ?", detail.CampaignDetailId).
				Scan(&itemData).Error

			if err == nil {
				response["operation_item_name"] = itemData.ItemName
				response["operation_item_code"] = itemData.ItemCode
			}

		default:
			response["operation_item_name"] = ""
			response["operation_item_code"] = ""
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
				LineTypeCode:    pack.LineTypeCode,
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
				LineTypeCode:    pack.LineTypeCode,
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
