package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	mastercampaignmasterentities "after-sales/api/entities/master/campaign_master"
	masteritementities "after-sales/api/entities/master/item"
	masteroperationentities "after-sales/api/entities/master/operation"
	masterpackagemasterentity "after-sales/api/entities/master/package-master"
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type CampaignMasterRepositoryImpl struct {
}

func StartCampaignMasterRepositoryImpl() masterrepository.CampaignMasterRepository {
	return &CampaignMasterRepositoryImpl{}
}

func (r *CampaignMasterRepositoryImpl) PostCampaignMaster(tx *gorm.DB, req masterpayloads.CampaignMasterPost) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entity mastercampaignmasterentities.CampaignMaster

	result, _ := tx.Model(&entity).Where("campaign_code =?", req.CampaignCode).First(&entity).Rows()
	if result != nil {
		return true, nil
	}

	entities := &mastercampaignmasterentities.CampaignMaster{
		IsActive:           req.IsActive,
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
	err := tx.Save(entities).Error
	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{}
	}
	return true, nil
}

// func (r *CampaignMasterRepositoryImpl) PostCampaignDetailMaster(tx *gorm.DB, req masterpayloads.CampaignMasterDetailPayloads) (bool, *exceptionsss_test.BaseErrorResponse) {
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
// 		return false, &exceptionsss_test.BaseErrorResponse{}
// 	}

// 	err = tx.Table("campaign_master_details").Joins("JOIN campaign_masters ON campaign_master_details.campaign_id = campaign_masters.id").Select("campaign_masters.total").Where("campaign_master.campaign_id = ?", req.CampaignId).Scan(&result).Error
// 	if err != nil{
// 		return false, &exceptionsss_test.BaseErrorResponse{}
// 	}
// 	totalValue:=result+float64(req.OperationItemPrice)
// 	result = r.UpdateTotalCampaignMaster(tx,req.CampaignId,totalValue)
// 	if result ==false{
// 		return false,&exceptionsss_test.BaseErrorResponse{}
// 	}
// 	return true, nil
// }

func (r *CampaignMasterRepositoryImpl) PostCampaignMasterDetailFromHistory(tx *gorm.DB, id int, idhead int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entityitem []mastercampaignmasterentities.CampaignMasterDetailItem
	var entityOperation []mastercampaignmasterentities.CampaignMasterOperationDetail

	result := tx.Model(&entityitem).Where("campaign_id = ?", id).Find(&entityitem)
	if result.Error != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        result.Error,
		}
	}
	for _, entities := range entityitem {
		newEntity := mastercampaignmasterentities.CampaignMasterDetailItem{
			IsActive:         entities.IsActive,
			CampaignDetailId: 0,
			CampaignId:       idhead,
			LineTypeId:       entities.LineTypeId,
			Quantity:         entities.Quantity,
			ItemId:           entities.ItemId,
			ShareBillTo:      entities.ShareBillTo,
			DiscountPercent:  entities.DiscountPercent,
			SharePercent:     entities.SharePercent,
		}
		err := tx.Create(&newEntity).Error

		if err != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}
	results := r.UpdateTotalCampaignMaster(tx, idhead)
	if !results {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	result2 := tx.Model(&entityOperation).Where("campaign_id=?", id).Find(&entityOperation)
	if result2.Error != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        result2.Error,
		}
	}
	for _, entities := range entityOperation {
		newEntity := mastercampaignmasterentities.CampaignMasterOperationDetail{
			IsActive:         entities.IsActive,
			CampaignDetailId: 0,
			CampaignId:       idhead,
			LineTypeId:       entities.LineTypeId,
			Quantity:         entities.Quantity,
			OperationId:      entities.OperationId,
			ShareBillTo:      entities.ShareBillTo,
			DiscountPercent:  entities.DiscountPercent,
			SharePercent:     entities.SharePercent,
		}
		err := tx.Create(&newEntity).Error

		if err != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}
	results2 := r.UpdateTotalCampaignMaster(tx, idhead)
	if !results2 {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result2.Error,
		}
	}
	return true, nil
}

func (r *CampaignMasterRepositoryImpl) PostCampaignDetailMaster(tx *gorm.DB, req masterpayloads.CampaignMasterDetailPayloads) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entityitem masteritementities.Item
	var entityoperation masteroperationentities.OperationModelMapping
	var lastprice float64
	if req.SharePercent > req.DiscountPercent {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Share percent must not be higher that discountpercent",
		}
	}

	if req.LineTypeId != 5 {
		err := tx.Model(&entityitem).
			Where("item_id=?", req.OperationItemId).
			Select("mtr_item.last_price").Scan(&lastprice).Error
		if err != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
		entities := &mastercampaignmasterentities.CampaignMasterDetailItem{
			CampaignDetailId: req.CampaignDetailId,
			CampaignId:       req.CampaignId,
			LineTypeId:       req.LineTypeId,
			Quantity:         req.Quantity,
			ItemId:           req.OperationItemId,
			ShareBillTo:      req.ShareBillTo,
			DiscountPercent:  req.DiscountPercent,
			SharePercent:     req.SharePercent,
			Price:            lastprice,
		}
		err2 := tx.Create(&entities).Error

		if err2 != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err2,
			}
		}
	} else {
		err := tx.Model(&entityoperation).
			Joins("JOIN mtr_operation_code on mtr_operation_code.operation_code_id=mtr_operation_model_mapping.operation_id").
			Joins("JOIN mtr_package_master on mtr_package_master.package_id=mtr_package_master_detail_operation.package_id").
			Joins("JOIN mtr_labour_selling_price on mtr_labour_selling_price.job_type_id=mtr_package_master_detail_operation.job_type_id AND mtr_labour_selling_price.brand_id= mtr_package_master.brand_id").
			Joins("JOIN mtr_labour_selling_price_detail on mtr_labour_selling_price_detail.labour_selling_price_id=mtr_labour_selling_price.labour_selling_price_id").
			Where("operation_model_mapping_id=?", req.OperationItemId).Select("mtr_labour_selling_price.selling_price").Scan(&lastprice).Error
		if err!=nil{
			return false,&exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		entities2 := &mastercampaignmasterentities.CampaignMasterOperationDetail{
			CampaignDetailId: req.CampaignDetailId,
			CampaignId:       req.CampaignId,
			LineTypeId:       req.LineTypeId,
			Quantity:         req.Quantity,
			OperationId:      req.OperationItemId,
			ShareBillTo:      req.ShareBillTo,
			DiscountPercent:  req.DiscountPercent,
			SharePercent:     req.SharePercent,
			Price:            lastprice,
		}
		err2 := tx.Create(&entities2).Error
		if err2 != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err2,
			}
		}
	}
	result := r.UpdateTotalCampaignMaster(tx, req.CampaignId)
	if !result {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error totalling price",
		}
	}

	return true, nil
}

func (r *CampaignMasterRepositoryImpl) ChangeStatusCampaignMaster(tx *gorm.DB, id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities mastercampaignmasterentities.CampaignMaster

	result := tx.
		Model(&entities).
		Where(mastercampaignmasterentities.CampaignMaster{CampaignId: id}).
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

// func (r *CampaignMasterRepositoryImpl) ChangeStatusCampaignMasterDetail(tx *gorm.DB, id []string)(bool,*exceptionsss_test.BaseErrorResponse){
// 	var entities []masterentities.CampaignMasterDetail

// 	rows, err := tx.Model(&entities).
// 		Where("warehouse_id in ?", id).
// 		Scan(&entities).
// 		Rows()

// 	if err != nil {
// 		return false, &exceptionsss_test.BaseErrorResponse{
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

func (r *CampaignMasterRepositoryImpl) DeactivateCampaignMasterDetail(tx *gorm.DB, ids string, idhead int) (bool, *exceptionsss_test.BaseErrorResponse) {
	idSlice := strings.Split(ids, ",")

	for _, Id := range idSlice {
		var entityToUpdateOperation mastercampaignmasterentities.CampaignMasterOperationDetail
		var entityToUpdateItem mastercampaignmasterentities.CampaignMasterDetailItem
		result := tx.Model(&entityToUpdateOperation).Where("campaign_detail_id = ?", Id).Where("Campaign_id=?", idhead).First(&entityToUpdateOperation).Update("is_active", false)
		if result.Error != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        result.Error,
			}
		}

		results := r.UpdateTotalCampaignMaster(tx, entityToUpdateOperation.CampaignId)
		if !results {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}

		result2 := tx.Model(&entityToUpdateItem).Where("campaign_detail_id = ?", Id).Where("Campaign_id=?", idhead).First(&entityToUpdateOperation).Update("is_active", false)
		if result2.Error != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        result.Error,
			}
		}
		results2 := r.UpdateTotalCampaignMaster(tx, entityToUpdateOperation.CampaignId)
		if !results2 {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}
	return true, nil
}

func (r *CampaignMasterRepositoryImpl) ActivateCampaignMasterDetail(tx *gorm.DB, ids string, idhead int) (bool, *exceptionsss_test.BaseErrorResponse) {
	idSlice := strings.Split(ids, ",")

	for _, Id := range idSlice {
		var entityToUpdateOperation mastercampaignmasterentities.CampaignMasterOperationDetail
		var entityToUpdateItem mastercampaignmasterentities.CampaignMasterDetailItem
		result := tx.Model(&entityToUpdateOperation).Where("campaign_detail_id = ?", Id).Where("Campaign_id=?", idhead).First(&entityToUpdateOperation).Update("is_active", false)
		if result.Error != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        result.Error,
			}
		}

		results := r.UpdateTotalCampaignMaster(tx, entityToUpdateOperation.CampaignId)
		if !results {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}

		result2 := tx.Model(&entityToUpdateItem).Where("campaign_detail_id = ?", Id).Where("Campaign_id=?", idhead).First(&entityToUpdateOperation).Update("is_active", false)
		if result2.Error != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        result.Error,
			}
		}
		results2 := r.UpdateTotalCampaignMaster(tx, entityToUpdateOperation.CampaignId)
		if !results2 {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}
	return true, nil
}

func (r *CampaignMasterRepositoryImpl) GetByIdCampaignMaster(tx *gorm.DB, id int) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse) {
	entities := mastercampaignmasterentities.CampaignMaster{}
	payloads := masterpayloads.CampaignMasterResponse{}
	var modelresponse masterpayloads.GetModelResponse
	var brandresponse masterpayloads.GetBrandResponse
	err := tx.Model(&entities).Where(mastercampaignmasterentities.CampaignMaster{
		CampaignId: id,
	}).First(&payloads).Error
	if err != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	brandIdUrl := "http://10.1.32.26:8000/sales-service/api/sales/unit-brand/" + strconv.Itoa(payloads.BrandId)
	errUrlBrandId := utils.Get(brandIdUrl, &brandresponse, nil)
	if errUrlBrandId != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlBrandId,
		}
	}
	BrandJoinData := utils.DataFrameInnerJoin([]masterpayloads.CampaignMasterResponse{payloads}, []masterpayloads.GetBrandResponse{brandresponse}, "BrandId")

	modelIdUrl := "http://10.1.32.26:8000/sales-service/api/sales/unit-model/" + strconv.Itoa(payloads.ModelId)
	errUrlModelId := utils.Get(modelIdUrl, &modelresponse, nil)
	if errUrlModelId != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlModelId,
		}
	}
	ModelIdJoinData := utils.DataFrameInnerJoin(BrandJoinData, []masterpayloads.GetModelResponse{modelresponse}, "ModelId")

	fmt.Printf("BrandJoinData: %+v\n", BrandJoinData)
	fmt.Printf("ModelIdJoinData: %+v\n", ModelIdJoinData)

	return ModelIdJoinData, nil
}

func (r *CampaignMasterRepositoryImpl) GetByIdCampaignMasterDetail(tx *gorm.DB, id int, linetypeid int) (map[string]interface{}, *exceptionsss_test.BaseErrorResponse) {
	entitiesitem := mastercampaignmasterentities.CampaignMasterDetailItem{}
	payloadsitem := masterpayloads.CampaignMasterDetailItemPayloads{}
	entitiesoperation := mastercampaignmasterentities.CampaignMasterOperationDetail{}
	payloadsoperation := masterpayloads.CampaignMasterDetailOperationPayloads{}
	if linetypeid == 1 {
		err := tx.Model(&entitiesoperation).
			Where("campaign_detail_id = ?", id).
			Joins("JOIN mtr_operation_model_mapping ON mtr_operation_model_mapping.operation_id=mtr_campaign_detail_operation.operation_id").
			Select("mtr_campaign_detail_operation.*,mtr_operation_code.operation_code,mtr_operation_code.operation_name").
			First(&payloadsoperation).Error
		responsepayload := map[string]interface{}{
			"is_active":        payloadsoperation.IsActive,
			"package_id":       payloadsoperation.PackageId,
			"package_code":     payloadsoperation.PackageCode,
			"line_type_id":     payloadsoperation.LineTypeId,
			"operation_code":   payloadsoperation.OperationCode,
			"operation_name":   payloadsoperation.OperationName,
			"quantity":         payloadsoperation.Quantity,
			"price":            payloadsoperation.Price,
			"subtotal":         int(payloadsoperation.Price * payloadsitem.Quantity),
			"discount_percent": payloadsoperation.DiscountPercent,
			"total":            float64(payloadsoperation.Price*payloadsitem.Quantity - (payloadsoperation.Price * payloadsitem.Quantity * payloadsoperation.DiscountPercent)),
			"share_percent":    payloadsoperation.SharePercent,
		}
		if err != nil {
			return responsepayload, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		return responsepayload, nil
	} else {
		err2 := tx.Model(&entitiesitem).
			Where("package_id=?", id).
			Joins("JOIN mtr_item on mtr_item.item_id=mtr_campaign_detail_item.item_id").
			Select("mtr_campaign_detail_item.*,mtr_item.item_code,mtr_item.item_name").
			First(&payloadsitem).Error
		responsepayload := map[string]interface{}{
			"is_active":        payloadsitem.IsActive,
			"package_code":     payloadsitem.PackageCode,
			"package_id":       payloadsitem.PackageId,
			"line_type_id":     payloadsitem.LineTypeId,
			"item_code":        payloadsitem.ItemCode,
			"item_name":        payloadsitem.ItemName,
			"quantity":         payloadsitem.Quantity,
			"price":            payloadsitem.Price,
			"subtotal":         int(payloadsitem.Quantity * payloadsitem.Price),
			"discount_percent": payloadsitem.DiscountPercent,
			"total":            int(payloadsitem.Quantity*payloadsitem.Price - (payloadsitem.Quantity * payloadsitem.Price * payloadsitem.DiscountPercent)),
			"share_percent":    payloadsitem.SharePercent,
		}
		if err2 != nil {
			return nil, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err2,
			}
		}
		return responsepayload, nil
	}

}

func (r *CampaignMasterRepositoryImpl) GetAllCampaignMasterCodeAndName(tx *gorm.DB, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	CampaignMasterMapping := []mastercampaignmasterentities.CampaignMaster{}
	CampaignMasterResponse := []masterpayloads.GetHistory{}
	query := tx.Model(mastercampaignmasterentities.CampaignMaster{}).Scan(&CampaignMasterResponse)
	err := query.Scopes(pagination.Paginate(&CampaignMasterMapping, &pages, query)).Scan(&CampaignMasterResponse).Error
	if len(CampaignMasterResponse) == 0 {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {

		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	pages.Rows = CampaignMasterResponse
	return pages, nil
}

func (r *CampaignMasterRepositoryImpl) GetAllCampaignMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	entities := mastercampaignmasterentities.CampaignMaster{}
	payloads := []masterpayloads.CampaignMasterResponse{}
	baseModelQuery := tx.Model(&entities).Scan(&payloads)

	Wherequery := utils.ApplyFilter(baseModelQuery, filterCondition)

	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, Wherequery)).Scan(&payloads).Rows()

	if len(payloads) == 0 {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()
	pages.Rows = payloads
	return pages, nil
}

func (r *CampaignMasterRepositoryImpl) GetAllCampaignMasterDetail(tx *gorm.DB, pages pagination.Pagination, id int) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	entitiesoperation := []mastercampaignmasterentities.CampaignMasterOperationDetail{}
	entitiesitem := []mastercampaignmasterentities.CampaignMasterDetailItem{}
	responseoperation := []masterpayloads.CampaignMasterDetailOperationPayloads{}
	responseitem := []masterpayloads.CampaignMasterDetailItemPayloads{}
	combinedPayloads := make([]map[string]interface{}, 0)

	err := tx.Model(&entitiesoperation).Where(mastercampaignmasterentities.CampaignMasterOperationDetail{
		CampaignId: id,
	}).Joins("JOIN mtr_operation_model_mapping ON mtr_operaton_model_mapping.operation_id=mtr_campaign_master_detail_operation.operation_id").
		Select("mtr_campaign_master_detail_item.*,mtr_operation_model_mapping.operation_code,mtr_operation_model_mapping.operation_name").
		Scan(&responseoperation).
		Error

	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
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
			"total":            float64(op.Price * op.Quantity),
		})
	}
	err2 := tx.Model(&entitiesitem).Where(mastercampaignmasterentities.CampaignMasterDetailItem{
		CampaignId: id,
	}).Joins("JOIN mtr_item ON mtr_item.item_id=mtr_campaign_master_detail_item.item_id").
		Select("mtr_campaign_master_detail_item.*,mtr_item.item_code,mtr_item.item_name").
		Scan(&responseitem).
		Error
	if err2 != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
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
			"total":            float64(it.Quantity * it.Price),
		})
	}
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(combinedPayloads, &pages)
	return dataPaginate, totalPages, totalRows, nil
}

func (r *CampaignMasterRepositoryImpl) UpdateCampaignMasterDetail(tx *gorm.DB, id int, req masterpayloads.CampaignMasterDetailPayloads) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities mastercampaignmasterentities.CampaignMasterOperationDetail

	result := tx.Model(&entities).Where("campaign_detail_id = ?", id).First(&entities)
	if result.Error != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        result.Error,
		}
	}

	update := tx.Model(&entities).Where("campaign_detail_id = ?", id).Updates(req)
	if update.Error != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        update.Error,
		}
	}

	results := r.UpdateTotalCampaignMaster(tx, req.CampaignId)
	if !results {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        result.Error,
		}
	}

	return true, nil

}

func (r *CampaignMasterRepositoryImpl) UpdateTotalCampaignMaster(tx *gorm.DB, id int) bool {
	var operations []mastercampaignmasterentities.CampaignMasterOperationDetail
	var items []mastercampaignmasterentities.CampaignMasterDetailItem
	var entity mastercampaignmasterentities.CampaignMaster
	var totalValue float64

	// Fetch and calculate the total value from CampaignMasterOperationDetail
	result := tx.Model(&operations).Where(mastercampaignmasterentities.CampaignMasterOperationDetail{
		CampaignId: id,
	}).Find(&operations)
	if result.Error != nil {
		return false
	}
	for _, detail := range operations {
		if detail.IsActive {
			totalValue += detail.Quantity * detail.Price * (1 - (detail.DiscountPercent / 100))
		}
	}

	// Fetch and calculate the total value from CampaignMasterDetailItem
	result = tx.Model(&items).Where(mastercampaignmasterentities.CampaignMasterDetailItem{
		CampaignId: id,
	}).Find(&items)
	if result.Error != nil {
		return false
	}
	for _, item := range items {
		if item.IsActive { // Assuming there's an IsActive field similar to operations
			totalValue += item.Quantity * item.Price * (1 - (item.DiscountPercent / 100))
		}
	}
	result = tx.Model(&entity).Where(mastercampaignmasterentities.CampaignMaster{
		CampaignId: id,
	}).Update("total", totalValue)
	return result.Error == nil
}

func (r *CampaignMasterRepositoryImpl) GetAllPackageMasterToCopy(tx *gorm.DB, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	var packageentities masterentities.PackageMaster
	var payloads []masterpayloads.PackageMasterForCampaignMaster

	BaseModelQuery := tx.Model(&packageentities)
	rows, err := BaseModelQuery.Scopes(pagination.Paginate(&packageentities, &pages, BaseModelQuery)).Where("profit_center_id=?", 13).Scan(payloads).Rows()
	if len(payloads) == 0 {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	pages.Rows = payloads

	return pages, nil
}

func (r *CampaignMasterRepositoryImpl) SelectFromPackageMaster(tx *gorm.DB, id int, idhead int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entityitem masterpackagemasterentity.PackageMasterDetailItem
	var payloadsitem []masterpayloads.PackageMasterDetailItem
	var entityOperation masterpackagemasterentity.PackageMasterDetailOperation
	var payloadsoperation []masterpayloads.PackageMasterDetailOperation

	err := tx.Model(&entityitem).
		Joins("JOIN mtr_item on mtr_item.item_id=mtr_package_master_detail_item.item_id").
		Where("package_id=?", id).
		Select("mtr_package_master_detail_item.*,mtr_item.price_list_item").
		Scan(&payloadsitem).Error
	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	for _, it := range payloadsitem {
		entities := mastercampaignmasterentities.CampaignMasterDetailItem{
			CampaignId: idhead,
			LineTypeId: it.LineTypeId,
			Quantity:   it.FrtQuantity,
			ItemId:     it.ItemId,
		}
		err := tx.Create(entities).Error

		if err != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
	}
	err2 := tx.Model(entityOperation).
		Joins("JOIN mtr_operation_model_mapping on mtr_operation_model_mapping.operation_id=mtr_package_master_detail_operation.operation_id").
		Joins("JOIN mtr_operation_code on mtr_operation_code.operation_code_id=mtr_operation_model_mapping.operation_id").
		Joins("JOIN mtr_package_master on mtr_package_master.package_id=mtr_package_master_detail_operation.package_id").
		Joins("JOIN mtr_labour_selling_price on mtr_labour_selling_price.job_type_id=mtr_package_master_detail_operation.job_type_id AND mtr_labour_selling_price.brand_id= mtr_package_master.brand_id").
		Joins("JOIN mtr_labour_selling_price_detail on mtr_labour_selling_price_detail.labour_selling_price_id=mtr_labour_selling_price.labour_selling_price_id").
		Where("package_id=?", id).Select("mtr_package_master_detail.*,mtr_labour_selling_price_detail.selling_price,mtr_operation_code.operation_code,mtr_operation_code.operation_name").Scan(&payloadsoperation).Error
	if err2 != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err2,
		}
	}
	for _, op := range payloadsoperation {
		entities := mastercampaignmasterentities.CampaignMasterOperationDetail{
			CampaignId:  idhead,
			LineTypeId:  op.LineTypeId,
			Quantity:    op.FrtQuantity,
			OperationId: op.OperationId,
		}
		err := tx.Create(entities).Error
		if err != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
	}
	return true, nil
}
