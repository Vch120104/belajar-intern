package masterrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"errors"
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
	var entities masterentities.CampaignMaster

	result, _ := tx.Model(&entities).Where("campaign_code = ? AND campaign_id != ?", req.CampaignCode, req.CampaignId).First(&entities).Rows()
	if result != nil {
		return masterentities.CampaignMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        errors.New("code exists"),
		}
	}

	entity := &masterentities.CampaignMaster{
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
	err := tx.Save(&entity).Error
	if err != nil {
		return masterentities.CampaignMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return *entity, nil
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
	var entity masterentities.CampaignMaster
	var lastprice *float64
	if req.SharePercent > req.DiscountPercent {
		return masterentities.CampaignMasterDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Share percent must not be higher that discount percent",
		}
	}

	if req.LineTypeId != 9 && req.LineTypeId != 0 { //not operation line type id
		err := tx.Select("mtr_item_price_list.price_list_amount").Table("mtr_item_price_list").
			Joins("JOIN mtr_item on mtr_item.item_id=mtr_item_price_list.item_id").
			Joins("Join mtr_item_operation on mtr_item.item_id=mtr_item_operation.item_operation_model_mapping_id").
			Where("item_operation_id=?", req.OperationItemId).
			Scan(&lastprice).Error
		if err != nil {
			return masterentities.CampaignMasterDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
		if lastprice == nil {
			lastprice = new(float64)
			*lastprice = 0.0
		}

	} else {
		err2 := tx.Model(&entity).Where("campaign_id=?", id).Scan(&entity).Error
		if err2 != nil {
			return masterentities.CampaignMasterDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err2,
			}
		}
		err := tx.Select("mtr_labour_selling_price_detail.selling_price").
			Table("mtr_labour_selling_price_detail").
			Joins("Join mtr_labour_selling_price on mtr_labour_selling_price.labour_selling_price_id = mtr_labour_selling_price_detail.labour_selling_price_id").
			Where("mtr_labour_selling_price.brand_id =?", entity.BrandId).
			Where("mtr_labour_selling_price_detail.model_id=?", entity.ModelId).
			Where("mtr_labour_selling_price.company_id = ?", entity.CompanyId).
			Where("mtr_labour_selling_price.effective_date < ?", time.Now()).
			Scan(&lastprice).Error
		if err != nil {
			return masterentities.CampaignMasterDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		if lastprice == nil {
			lastprice = new(float64)
			*lastprice = 0.0
		}
	}
	entities := &masterentities.CampaignMasterDetail{
		CampaignId:      id,
		LineTypeId:      req.LineTypeId,
		Quantity:        req.Quantity,
		ItemOperationId: req.OperationItemId,
		ShareBillTo:     req.ShareBillTo,
		DiscountPercent: req.DiscountPercent,
		SharePercent:    req.SharePercent,
		Price:           *lastprice,
	}
	err2 := tx.Save(&entities).Error

	if err2 != nil {
		return masterentities.CampaignMasterDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err2,
		}
	}
	results := r.UpdateTotalCampaignMaster(tx, id)
	if !results {
		return masterentities.CampaignMasterDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to update"),
		}
	}
	return *entities, nil
}

func (r *CampaignMasterRepositoryImpl) PostCampaignMasterDetailFromPackage(tx *gorm.DB, req masterpayloads.CampaignMasterDetailPostFromPackageRequest) (masterentities.CampaignMasterDetail, *exceptions.BaseErrorResponse) {
	campaignDetailEntities := masterentities.CampaignMasterDetail{}
	response := masterentities.CampaignMasterDetail{}

	companyReferenceUrl := config.EnvConfigs.GeneralServiceUrl + "company-reference/" + strconv.Itoa(req.CompanyId)
	companyReferencePayloads := masterpayloads.CampaignMasterCompanyReferenceResponse{}
	if err := utils.Get(companyReferenceUrl, &companyReferencePayloads, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching company reference",
			Err:        err,
		}
	}
	currencyId := companyReferencePayloads.CurrencyId

	jobTypeUrl := config.EnvConfigs.GeneralServiceUrl + "job-type-by-code/CP"
	jobTypePayloads := masterpayloads.CampaignMasterJobTypeResponse{}
	if err := utils.Get(jobTypeUrl, &jobTypePayloads, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching job type campaign",
			Err:        err,
		}
	}
	jobTypeCampaignId := jobTypePayloads.JobTypeId

	WOTransactionTypeUrl := config.EnvConfigs.GeneralServiceUrl + "work-order-transaction-type-by-code/Campaign"
	WOTransactionTypePayloads := masterpayloads.CampaignMasterWOTransactionResponse{}
	if err := utils.Get(WOTransactionTypeUrl, &WOTransactionTypePayloads, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching work order transaction",
			Err:        err,
		}
	}
	billCode := WOTransactionTypePayloads.WorkOrderTransactionTypeId

	var totalRows int
	err := tx.Model(&campaignDetailEntities).
		Select("COUNT(*) total_rows").
		Where("campaign_id = ?", req.CampaignId).
		Where("package_id = ?", req.PackageId).
		Pluck("total_rows", &totalRows).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if totalRows == 0 {
		responsePackageDetail := []masterentities.PackageMasterDetail{}

		err = tx.Model(&masterentities.PackageMasterDetail{}).
			Where("package_id = ?", req.PackageId).
			Where("is_active = ?", true).
			Scan(&responsePackageDetail).Error

		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		var warehouseGroup string
		err = tx.Model(&masterentities.WarehouseGroupMappingEntities{}).
			Select("warehouse_group_mapping_description").
			Where("warehouse_group_type_code = 'WHS_GRP_CAMPAIGN'").
			First(&warehouseGroup).Error

		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		for _, data := range responsePackageDetail {
			itemPriceCode, err := r.lookupRepo.GetOprItemPrice(tx, data.LineTypeId, req.CompanyId, data.ItemOperationId, req.BrandId, req.ModelId, jobTypeCampaignId, 0, currencyId, billCode, warehouseGroup)
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

		packageMasterDetailResponse := masterentities.PackageMasterDetail{}
		err = tx.Model(&masterentities.PackageMasterDetail{}).
			Where(masterentities.PackageMasterDetail{
				PackageId: req.PackageId,
				IsActive:  true,
			}).First(&packageMasterDetailResponse).Error

		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		itemPriceCode, itemPriceErr := r.lookupRepo.GetOprItemPrice(tx, packageMasterDetailResponse.LineTypeId, req.CompanyId, packageMasterDetailResponse.ItemOperationId, req.BrandId, req.ModelId, jobTypeCampaignId, 0, currencyId, billCode, warehouseGroup)

		if itemPriceErr != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        itemPriceErr.Err,
			}
		}

		campaignDetailEntities = masterentities.CampaignMasterDetail{
			IsActive:        true,
			CampaignId:      req.CampaignId,
			LineTypeId:      packageMasterDetailResponse.LineTypeId,
			ItemOperationId: packageMasterDetailResponse.ItemOperationId,
			Quantity:        packageMasterDetailResponse.FrtQuantity,
			ShareBillTo:     "",
			DiscountPercent: 0,
			SharePercent:    0,
			PackageId:       req.PackageId,
			Price:           itemPriceCode,
		}

		err = tx.Save(&campaignDetailEntities).Error
		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		response = campaignDetailEntities
	}

	return response, nil
}

func (r *CampaignMasterRepositoryImpl) ChangeStatusCampaignMaster(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masterentities.CampaignMaster

	result := tx.
		Model(&entities).
		Where(masterentities.CampaignMaster{CampaignId: id}).
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

	return entities.IsActive, nil
}

// func (r *CampaignMasterRepositoryImpl) ChangeStatusCampaignMasterDetail(tx *gorm.DB, id []string)(bool,*exceptions.BaseErrorResponse){
// 	var entities []masterentities.CampaignMasterDetail

// 	rows, err := tx.Model(&entities).
// 		Where("warehouse_id in ?", id).
// 		Scan(&entities).
// 		Rows()

// 	if err != nil {
// 		return false, &exceptions.BaseErrorResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Err:        err,
// 		}
// 	}
// 	if rows.IsActive {
// 		rows.IsActive = false
// 	} else {
// 		rows.IsActive = true
// 	}

// 	defer rows.Close()

// 	return true, nil
// }

func (r *CampaignMasterRepositoryImpl) DeactivateCampaignMasterDetail(tx *gorm.DB, ids string) (bool, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(ids, ",")
	for _, Id := range idSlice {
		var entityToUpdate masterentities.CampaignMasterDetail
		result := tx.Model(&entityToUpdate).Where("campaign_detail_id = ?", Id).Where("Campaign_id=?", Id).First(&entityToUpdate).Update("is_active", false)
		if result.Error != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        result.Error,
			}
		}

		results := r.UpdateTotalCampaignMaster(tx, entityToUpdate.CampaignId)
		if !results {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}
	return true, nil
}

func (r *CampaignMasterRepositoryImpl) ActivateCampaignMasterDetail(tx *gorm.DB, ids string) (bool, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(ids, ",")

	for _, Id := range idSlice {
		var entityToUpdate masterentities.CampaignMasterDetail
		result := tx.Model(&entityToUpdate).Where("campaign_detail_id = ?", Id).Where("Campaign_detail_id=?", Id).First(&entityToUpdate).Update("is_active", true)
		if result.Error != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        result.Error,
			}
		}

		results := r.UpdateTotalCampaignMaster(tx, entityToUpdate.CampaignId)
		if !results {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}

	}
	return true, nil
}

func (r *CampaignMasterRepositoryImpl) GetByIdCampaignMaster(tx *gorm.DB, id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	entities := masterentities.CampaignMaster{}
	payloads := masterpayloads.CampaignMasterResponse{}
	var modelresponse masterpayloads.GetModelResponse
	var brandresponse masterpayloads.GetBrandResponse
	err := tx.Model(&entities).Where(masterentities.CampaignMaster{
		CampaignId: id,
	}).First(&payloads).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	brandIdUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(payloads.BrandId)
	errUrlBrandId := utils.Get(brandIdUrl, &brandresponse, nil)
	if errUrlBrandId != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlBrandId,
		}
	}
	BrandJoinData, errdf := utils.DataFrameInnerJoin([]masterpayloads.CampaignMasterResponse{payloads}, []masterpayloads.GetBrandResponse{brandresponse}, "BrandId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	modelIdUrl := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(payloads.ModelId)
	errUrlModelId := utils.Get(modelIdUrl, &modelresponse, nil)
	if errUrlModelId != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlModelId,
		}
	}
	ModelIdJoinData, errdf := utils.DataFrameInnerJoin(BrandJoinData, []masterpayloads.GetModelResponse{modelresponse}, "ModelId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	return ModelIdJoinData[0], nil
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
	var modelresponse masterpayloads.GetModelResponse
	var brandresponse masterpayloads.GetBrandResponse
	err := tx.Model(&entities).Where(masterentities.CampaignMaster{
		CampaignCode: code,
	}).First(&payloads).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	brandIdUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(payloads.BrandId)
	errUrlBrandId := utils.Get(brandIdUrl, &brandresponse, nil)
	if errUrlBrandId != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlBrandId,
		}
	}
	BrandJoinData, errdf := utils.DataFrameInnerJoin([]masterpayloads.CampaignMasterResponse{payloads}, []masterpayloads.GetBrandResponse{brandresponse}, "BrandId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	modelIdUrl := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(payloads.ModelId)
	errUrlModelId := utils.Get(modelIdUrl, &modelresponse, nil)
	if errUrlModelId != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlModelId,
		}
	}
	ModelIdJoinData, errdf := utils.DataFrameInnerJoin(BrandJoinData, []masterpayloads.GetModelResponse{modelresponse}, "ModelId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	return ModelIdJoinData[0], nil
}

func (r *CampaignMasterRepositoryImpl) GetAllCampaignMasterCodeAndName(tx *gorm.DB, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	CampaignMasterMapping := []masterentities.CampaignMaster{}
	CampaignMasterResponse := []masterpayloads.GetHistory{}
	query := tx.Model(masterentities.CampaignMaster{}).Scan(&CampaignMasterResponse)
	err := query.Scopes(pagination.Paginate(&CampaignMasterMapping, &pages, query)).Scan(&CampaignMasterResponse).Error
	if len(CampaignMasterResponse) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {

		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	pages.Rows = CampaignMasterResponse
	return pages, nil
}

func (r *CampaignMasterRepositoryImpl) GetAllCampaignMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	model := []masterpayloads.GetModelResponse{}
	entities := masterentities.CampaignMaster{}
	payloads := []masterpayloads.CampaignMasterResponse{}
	var mapResponses []map[string]interface{}

	var modelDescription string
	var modelCode string
	var campaignPeriodFrom string
	var campaignPeriodTo string
	var newFilterCondition []utils.FilterCondition
	for _, filter := range filterCondition {
		if filter.ColumnField == "model_description" && filter.ColumnValue != "" {
			modelDescription = filter.ColumnValue
			continue
		}
		if filter.ColumnField == "model_code" && filter.ColumnValue != "" {
			modelCode = filter.ColumnValue
			continue
		}
		if filter.ColumnField == "campaign_period_from" && filter.ColumnValue != "" {
			campaignPeriodFrom = filter.ColumnValue
			continue
		}
		if filter.ColumnField == "campaign_period_to" && filter.ColumnValue != "" {
			campaignPeriodTo = filter.ColumnValue
			continue
		}
		newFilterCondition = append(newFilterCondition, filter)
	}

	baseModelQuery := tx.Model(&entities)

	if modelDescription != "" {
		modelIds := []int{}
		modelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model?page=0&limit=1000000&model_description=" + modelDescription
		modelPayloads := []masterpayloads.GetModelResponse{}
		if err := utils.GetArray(modelUrl, &modelPayloads, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		if len(modelPayloads) > 0 {
			for _, model := range modelPayloads {
				modelIds = append(modelIds, model.ModelId)
			}
		} else {
			modelIds = append(modelIds, -1)
		}
		baseModelQuery = baseModelQuery.Where("model_id IN ?", modelIds)
	}

	if modelCode != "" {
		modelIds := []int{}
		modelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model?page=0&limit=1000000&model_code=" + modelCode
		modelPayloads := []masterpayloads.GetModelResponse{}
		if err := utils.GetArray(modelUrl, &modelPayloads, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		if len(modelPayloads) > 0 {
			for _, model := range modelPayloads {
				modelIds = append(modelIds, model.ModelId)
			}
		} else {
			modelIds = append(modelIds, -1)
		}
		baseModelQuery = baseModelQuery.Where("model_id IN ?", modelIds)
	}

	if campaignPeriodFrom != "" {
		baseModelQuery.Where("FORMAT(campaign_period_from, 'dd MMM yyyy') LIKE ?", "%"+campaignPeriodFrom+"%")
	}

	if campaignPeriodTo != "" {
		baseModelQuery.Where("FORMAT(campaign_period_to, 'dd MMM yyyy') LIKE ?", "%"+campaignPeriodTo+"%")
	}

	whereQuery := utils.ApplyFilter(baseModelQuery, newFilterCondition)
	err := whereQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&payloads).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	errUrlModel := utils.Get(config.EnvConfigs.SalesServiceUrl+"unit-model?page=0&limit=1000000", &model, nil)
	if errUrlModel != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	joineddata1, errdf := utils.DataFrameInnerJoin(payloads, model, "ModelId")

	if errdf != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	for _, response := range joineddata1 {
		responseMap := map[string]interface{}{
			"appointment_only":     response["AppointmentOnly"],
			"brand_id":             response["BrandID"],
			"campaign_code":        response["CampaignCode"],
			"campaign_id":          response["CampaignId"],
			"campaign_name":        response["CampaignName"],
			"campaign_period_from": response["CampaignPeriodFrom"],
			"campaign_period_to":   response["CampaignPeriodTo"],
			"is_active":            response["IsActive"],
			"model_code":           response["ModelCode"],
			"model_description":    response["ModelDescription"],
			"model_id":             response["ModelId"],
			"remark":               response["Remark"],
			"total":                response["Total"],
			"total_after_vat":      response["TotalAfterVat"],
			"total_vat":            response["TotalVat"],
		}
		mapResponses = append(mapResponses, responseMap)
	}

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)
	return dataPaginate, totalPages, totalRows, nil
}

func (r *CampaignMasterRepositoryImpl) GetAllCampaignMasterDetail(tx *gorm.DB, pages pagination.Pagination, id int) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []masterentities.CampaignMasterDetail
	var responsedetail []masterpayloads.CampaignMasterDetailGetPayloads
	var item masteritempayloads.BomItemNameResponse
	var operation masterpayloads.Operation
	var packagecode string
	combinedPayloads := make([]map[string]interface{}, 0)

	err := tx.Model(&entities).
		Where(masterentities.CampaignMasterDetail{
			CampaignId: id,
		}).Scan(&responsedetail).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	for _, op := range responsedetail {
		if op.PackageId != 0 {
			err := tx.Select("mtr_package.package_code").Table("mtr_package").Where("mtr_package.package_id=?", op.PackageId).Scan(packagecode).Error
			if err != nil {
				return nil, 0, 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
				}
			}
		}
		if op.LineTypeId != 9 && op.LineTypeId != 0 {
			err = tx.Select("mtr_item.item_name,mtr_item.item_code").Table("mtr_campaign_master_detail").
				Joins("join mtr_item_operation on mtr_item_operation.item_operation_id=mtr_campaign_master_detail.item_operation_id").
				Joins("join mtr_item on mtr_item.item_id=mtr_item_operation.item_operation_model_mapping_id").
				Where("mtr_campaign_master_detail.campaign_detail_id=?", op.CampaignDetailId).
				Scan(&item).
				Error
		} else {
			err = tx.Select("operation_code.operation_name,operation_code.operation_code").Where("campaign_detail_id=?", op.CampaignDetailId).
				Joins("join mtr_item_operation on mtr_item_operation.item_operation_id = mtr_campaign_master_detail.item_operation_id").
				Joins("JOIN mtr_operation_model_mapping ON mtr_operation_model_mapping.operation_model_mapping_id=mtr_item_operation.item_operation_model_mapping_id").
				Joins("join mtr_operation_code on mtr_operation_code.operation_id=mtr_operation_model_mapping.operation_id").
				Select("mtr_campaign_master_detail.*,mtr_operation_code.operation_code,mtr_operation_code.operation_name").
				Table("mtr_campaign_master_detail").
				Scan(&operation).Error
		}
		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		beforeDisc := op.Price * op.Quantity
		afterDisc := beforeDisc
		if op.DiscountPercent > 0 {
			afterDisc = beforeDisc - (beforeDisc * op.DiscountPercent / 100)
		}

		response := map[string]interface{}{
			"is_active":          op.IsActive,
			"campaign_id":        op.CampaignId,
			"campaign_detail_id": op.CampaignDetailId,
			"package_code":       packagecode,
			"package_id":         op.PackageId,
			"line_type_id":       op.LineTypeId,
			"item_operation_id":  op.ItemOperationId,
			"frt_quantity":       op.Quantity,
			"price":              op.Price,
			"discount_percent":   op.DiscountPercent,
			"share_percent":      op.SharePercent,
			"total":              afterDisc,
		}

		if op.LineTypeId != 9 && op.LineTypeId != 1 {
			response["item_name"] = item.ItemName
			response["item_code"] = item.ItemCode
		} else {
			response["operation_name"] = operation.OperationName
			response["operation_code"] = operation.OperationCode
		}
		combinedPayloads = append(combinedPayloads, response)
	}
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(combinedPayloads, &pages)
	return dataPaginate, totalPages, totalRows, nil
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
	rows, err := BaseModelQuery.Scopes(pagination.Paginate(&packageentities, &pages, BaseModelQuery)).Where("profit_center_id=?", 13).Scan(payloads).Rows()
	if len(payloads) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	pages.Rows = payloads

	return pages, nil
}

func (r *CampaignMasterRepositoryImpl) SelectFromPackageMaster(tx *gorm.DB, id int, idhead int) (int, *exceptions.BaseErrorResponse) {
	var packagedetail []masterentities.PackageMasterDetail
	var lastprice float64
	var operationpayloads masterpayloads.CampaignMasterDetailGetPayloads
	var entity masterentities.CampaignMaster
	var itemprice masteritementities.PriceList

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
