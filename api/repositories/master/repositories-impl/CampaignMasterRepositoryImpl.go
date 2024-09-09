package masterrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
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
}

func StartCampaignMasterRepositoryImpl() masterrepository.CampaignMasterRepository {
	return &CampaignMasterRepositoryImpl{}
}

func (r *CampaignMasterRepositoryImpl) PostCampaignMaster(tx *gorm.DB, req masterpayloads.CampaignMasterPost) (masterentities.CampaignMaster, *exceptions.BaseErrorResponse) {
	var entities masterentities.CampaignMaster

	result, _ := tx.Model(&entities).Where("campaign_code =?", req.CampaignCode).First(&entities).Rows()
	if result != nil {
		return masterentities.CampaignMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        errors.New("code exists"),
		}
	}

	entity := &masterentities.CampaignMaster{
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

// func (r *CampaignMasterRepositoryImpl) PostCampaignDetailMaster(tx *gorm.DB, req masterpayloads.CampaignMasterDetailPayloads) (bool, *exceptions.BaseErrorResponse) {
// 	var result float64
// 	entities := masterentities.CampaignMasterDetail{
// 		IsActive:          req.IsActive,
// 		CampaignId:        req.CampaignId,
// 		LineTypeId:        req.LineTypeId,
// 		Quantity:          req.Quantity,
// 		OperationItemCode: req.OperationItemCode,
// 		OperationItemPrice: req.OperationItemPrice,
// 		DiscountPercent:   req.DiscountPercent,
// 		SharePercent:      req.SharePercent,
// 	}
// 	err := tx.Save(entities).Error

// 	if err != nil {
// 		return false, &exceptions.BaseErrorResponse{}
// 	}

// 	err = tx.Table("campaign_master_details").Joins("JOIN campaign_masters ON campaign_master_details.campaign_id = campaign_masters.id").Select("campaign_masters.total").Where("campaign_master.campaign_id = ?", req.CampaignId).Scan(&result).Error
// 	if err != nil{
// 		return false, &exceptions.BaseErrorResponse{}
// 	}
// 	totalValue:=result+float64(req.OperationItemPrice)
// 	result = r.UpdateTotalCampaignMaster(tx,req.CampaignId,totalValue)
// 	if result ==false{
// 		return false,&exceptions.BaseErrorResponse{}
// 	}
// 	return true, nil
// }

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

func (r *CampaignMasterRepositoryImpl) PostCampaignDetailMaster(tx *gorm.DB, req masterpayloads.CampaignMasterDetailPayloads) (int, *exceptions.BaseErrorResponse) {
	var itemprice masteritementities.PriceList
	var entity masterentities.CampaignMaster
	var lastprice *float64
	if req.SharePercent > req.DiscountPercent {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Share percent must not be higher that discount percent",
		}
	}

	if req.LineTypeId != 5 {//not operation line type id
		err := tx.Model(&itemprice).Select("mtr_price_list.price_list_amount").
			Where("item_id=?", req.OperationItemId).
			Scan(&lastprice).Error
		if err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
		if lastprice == nil {
			lastprice = new(float64)
			*lastprice = 0.0
		}
		
	} else {
		err1 := tx.Model(&entity).Where("campaign_id = ?", req.CampaignId).Scan(&entity).Error
		if err1 != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err1,
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
			return 0, &exceptions.BaseErrorResponse{
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
			CampaignId:      req.CampaignId,
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
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err2,
			}
		}
		return entities.CampaignDetailId, nil
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
	entities := masterentities.CampaignMasterDetail{}
	err := tx.Model(&entities).Where("campaign_detail_id = ?", id).Scan(&entities).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	var payload map[string]interface{}
	if entities.LineTypeId == 5 {
		payload, err = getOperationPayload(tx, id)
	} else {
		payload, err = getItemPayload(tx, id)
	}

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return payload, nil
}

func getOperationPayload(tx *gorm.DB, id int) (map[string]interface{}, error) {
	payloadsoperation := masterpayloads.CampaignMasterDetailOperationPayloads{}
	err := tx.Model(&payloadsoperation).
		Where("campaign_detail_id = ?", id).
		Joins("JOIN mtr_item_operation ON mtr_item_operation.item_operation_id = mtr_campaign_detail.item_operation_id").
		Joins("JOIN mtr_operation_model_mapping ON mtr_operation_model_mapping.operation_id = mtr_item_operation.item_operation_id").
		Select("mtr_campaign_master_detail.*, mtr_operation_code.operation_code, mtr_operation_code.operation_name").
		First(&payloadsoperation).Error

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"is_active":        payloadsoperation.IsActive,
		"package_id":       payloadsoperation.PackageId,
		"package_code":     payloadsoperation.PackageCode,
		"line_type_id":     payloadsoperation.LineTypeId,
		"operation_code":   payloadsoperation.OperationCode,
		"operation_name":   payloadsoperation.OperationName,
		"quantity":         payloadsoperation.Quantity,
		"price":            payloadsoperation.Price,
		"subtotal":         int(payloadsoperation.Price * payloadsoperation.Quantity),
		"discount_percent": payloadsoperation.DiscountPercent,
		"total":            float64(payloadsoperation.Price*payloadsoperation.Quantity - (payloadsoperation.Price * payloadsoperation.Quantity * payloadsoperation.DiscountPercent)),
		"share_percent":    payloadsoperation.SharePercent,
	}, nil
}

func getItemPayload(tx *gorm.DB, id int) (map[string]interface{}, error) {
	payloadsitem := masterpayloads.CampaignMasterDetailItemPayloads{}
	err := tx.Model(&payloadsitem).
		Where("campaign_detail_id = ?", id).
		Joins("JOIN mtr_item_operation ON mtr_item_operation.item_operation_id = mtr_campaign_detail.item_operation_id").
		Joins("JOIN mtr_item ON mtr_item.item_id = mtr_item_operation.item_id").
		Select("mtr_campaign_master_detail_items.*, mtr_item.item_code, mtr_item.item_name").
		First(&payloadsitem).Error

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"is_active":        payloadsitem.IsActive,
		"package_id":       payloadsitem.PackageId,
		"package_code":     payloadsitem.PackageCode,
		"line_type_id":     payloadsitem.LineTypeId,
		"item_code":        payloadsitem.ItemCode,
		"item_name":        payloadsitem.ItemName,
		"quantity":         payloadsitem.Quantity,
		"price":            payloadsitem.Price,
		"subtotal":         int(payloadsitem.Quantity * payloadsitem.Price),
		"discount_percent": payloadsitem.DiscountPercent,
		"total":            int(payloadsitem.Quantity*payloadsitem.Price - (payloadsitem.Quantity * payloadsitem.Price * payloadsitem.DiscountPercent)),
		"share_percent":    payloadsitem.SharePercent,
	}, nil
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
	baseModelQuery := tx.Model(&entities).Scan(&payloads)
	var mapResponses []map[string]interface{}

	Wherequery := utils.ApplyFilter(baseModelQuery, filterCondition)

	_, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, Wherequery)).Scan(&payloads).Rows()

	if len(payloads) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

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
			"camapign_period_to":   response["CampaignPeriodTo"],
			"is_active":            response["IsActive"],
			"model_code":           response["ModelCode"],
			"model_description":    response["ModelDescription"],
			"model_id":             response["ModelId"],
			"remark":               response["Remark"],
			"total":                response["Total"],
			"total_after_vat":      response["TotalAfterVAT"],
			"total_vat":            response["TotalVAT"],
		}
		mapResponses = append(mapResponses, responseMap)
	}

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)
	return dataPaginate, totalPages, totalRows, nil
}

func (r *CampaignMasterRepositoryImpl) GetAllCampaignMasterDetail(tx *gorm.DB, pages pagination.Pagination, id int) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	entities := []masterentities.CampaignMasterDetail{}
	responseoperation := []masterpayloads.CampaignMasterDetailOperationPayloads{}
	responseitem := []masterpayloads.CampaignMasterDetailItemPayloads{}
	combinedPayloads := make([]map[string]interface{}, 0)

	err := tx.Model(&entities).Where(masterentities.CampaignMasterDetail{
		CampaignId: id,
	}).Joins("join mtr_item_operation on mtr_item_operation.item_operation_id=mtr_campaign_detail.item_operation_id").
		Joins("JOIN mtr_operation_model_mapping ON mtr_operation_model_mapping.operation_model_mapping_id=mtr_item_operation.operation_model_mapping_id").
		Joins("JOIN mtr_operation_code ON mtr_operation_code.operation_id = mtr_operation_model_mapping.operation_id").
		Select("mtr_campaign_details.*,mtr_operation_code.operation_code,mtr_operation_code.operation_name").
		Scan(&responseoperation).
		Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	for _, op := range responseoperation {
		combinedPayloads = append(combinedPayloads, map[string]interface{}{
			"is_active":        op.IsActive,
			"package_code":     op.PackageCode,
			"package_id":       op.PackageId,
			"line_type_id":     op.LineTypeId,
			"operation_code":   op.OperationCode,
			"operation_name":   op.OperationName,
			"quantity":         op.Quantity,
			"price":            op.Price,
			"discount_percent": op.DiscountPercent,
			"share_percent":    op.SharePercent,
			"total":            float64(op.Price * op.Quantity * (1 - (op.DiscountPercent / 100))),
		})
	}
	err2 := tx.Model(&entities).Where(masterentities.CampaignMasterDetail{
		CampaignId: id,
	}).Joins("Join mtr_item_operation on mtr_item_operation.item_operation_id=mtr_campaign_detail.item_operation_id").
		Joins("JOIN mtr_item ON mtr_item.item_id=mtr_campaign_master_detail_items.item_id").
		Select("mtr_campaign_master_detail_items.*,mtr_item.item_code,mtr_item.item_name").
		Scan(&responseitem).
		Error
	if err2 != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err2,
		}
	}
	for _, it := range responseitem {
		combinedPayloads = append(combinedPayloads, map[string]interface{}{
			"is_active":        it.IsActive,
			"package_code":     it.PackageCode,
			"package_id":       it.PackageId,
			"line_type_id":     it.LineTypeId,
			"item_code":        it.ItemName,
			"item_name":        it.ItemName,
			"quantity":         it.Quantity,
			"price":            it.Price,
			"discount_percent": it.DiscountPercent,
			"share_percent":    it.SharePercent,
			"total":            float64(it.Quantity * it.Price * (1 - (it.DiscountPercent / 100))),
		})
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
			totalValue += detail.Quantity * detail.Price * (1 - (detail.DiscountPercent / 100))
		}
	}
	return true
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
	var operationpayloads masterpayloads.CampaignMasterDetailOperationPayloads
	var itempayloads masterpayloads.CampaignMasterDetailItemPayloads
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
				Select("mtr_package_master_detail.*,mtr_item.item_code,mtr_item.item_name").Scan(&itempayloads).Error
			if err2 != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Err:        err2,
				}
			}
			err := tx.Model(&itemprice).Select("mtr_price_list.price_list_amount").
				Where("item_id=?", itempayloads.ItemOperationId).
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
