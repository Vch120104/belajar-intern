package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type CampaignMasterRepositoryImpl struct {
}

func StartCampaignMasterRepositoryImpl() masterrepository.CampaignMasterRepository {
	return &CampaignMasterRepositoryImpl{}
}

func (r *CampaignMasterRepositoryImpl) PostCampaignMaster(tx *gorm.DB, req masterpayloads.CampaignMasterPost) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.CampaignMaster{
		IsActive:           req.IsActive,
		CampaignCode:       req.CampaignCode,
		CampaignName:       req.CampaignName,
		BrandId:            req.BrandId,
		ModelId:            req.ModelId,
		CampaignPeriodFrom: req.CampaignPeriodFrom,
		CampaignPeriodTo:   req.CampaignPeriodTo,
		Remark:             req.Remark,
		AppointmentOnly:    req.AppointmentOnly,
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

func (r *CampaignMasterRepositoryImpl) PostCampaignDetailMaster(tx *gorm.DB, req masterpayloads.CampaignMasterDetailPayloads) (bool, *exceptionsss_test.BaseErrorResponse) {
	var result float64
	entities := masterentities.CampaignMasterDetail{
		IsActive:          req.IsActive,
		CampaignId:        req.CampaignId,
		LineTypeId:        req.LineTypeId,
		Quantity:          req.Quantity,
		OperationItemCode: req.OperationItemCode,
		OperationItemPrice: req.OperationItemPrice,
		DiscountPercent:   req.DiscountPercent,
		SharePercent:      req.SharePercent,
	}
	err := tx.Save(entities).Error

	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{}
	}

	err = tx.Table("campaign_master").Select("campaign_masters.total").Where("campaign_master.campaign_id = ?", req.CampaignId).Scan(&result).Error
	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{}
	}

	totalValue := result + float64(req.OperationItemPrice)
	success := r.UpdateTotalCampaignMaster(tx, req.CampaignId, totalValue)
	if !success {
		return false, &exceptionsss_test.BaseErrorResponse{}
	}
	return true, nil
}

func (r *CampaignMasterRepositoryImpl) ChangeStatusCampaignMaster(tx *gorm.DB, id int)(bool,*exceptionsss_test.BaseErrorResponse){
	var entities masterentities.CampaignMaster

	result := tx.
		Model(&entities).
		Where(masterentities.CampaignMaster{CampaignId: id}).
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

	
func (r *CampaignMasterRepositoryImpl) DeactivateCampaignMasterDetail(tx *gorm.DB, ids string) (bool, *exceptionsss_test.BaseErrorResponse) {
    var result float64
	idSlice := strings.Split(ids, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masterentities.CampaignMasterDetail
		err := tx.Model(&entityToUpdate).Where("campaign_detail_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		err = tx.Table("campaign_master").Select("campaign_masters.total").Where("campaign_master_details.id = ?", entityToUpdate.CampaignId).Scan(&result).Error
		if err != nil {
			return false, &exceptionsss_test.BaseErrorResponse{}
		}

		totalValue := result - float64(entityToUpdate.OperationItemPrice)
		success := r.UpdateTotalCampaignMaster(tx, entityToUpdate.CampaignId, totalValue)
		if !success {
			return false, &exceptionsss_test.BaseErrorResponse{}
		}




		entityToUpdate.IsActive = false
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}

func (r *CampaignMasterRepositoryImpl) ActivateCampaignMasterDetail(tx *gorm.DB, ids string) (bool, *exceptionsss_test.BaseErrorResponse) {
    var result float64
	idSlice := strings.Split(ids, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masterentities.CampaignMasterDetail
		err := tx.Model(&entityToUpdate).Where("campaign_detail_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		err = tx.Table("campaign_master").Select("campaign_masters.total").Where("campaign_master_details.id = ?", entityToUpdate.CampaignId).Scan(&result).Error
		if err != nil {
			return false, &exceptionsss_test.BaseErrorResponse{}
		}

		totalValue := result + float64(entityToUpdate.OperationItemPrice)
		success := r.UpdateTotalCampaignMaster(tx, entityToUpdate.CampaignId, totalValue)
		if !success {
			return false, &exceptionsss_test.BaseErrorResponse{}
		}

		entityToUpdate.IsActive = true
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}
func (r *CampaignMasterRepositoryImpl) GetByIdCampaignMaster(tx *gorm.DB, id int) (masterpayloads.CampaignMasterResponse, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.CampaignMaster{}
	payloads := masterpayloads.CampaignMasterResponse{}
	rows, err := tx.Model(&entities).Where(masterentities.CampaignMaster{
		CampaignId: id,
	}).First(&payloads).Rows()
	if err != nil {
		return payloads, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return payloads, nil
}

func (r *CampaignMasterRepositoryImpl) GetByIdCampaignMasterDetail(tx *gorm.DB, id int) (masterpayloads.CampaignMasterDetailPayloads, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.CampaignMasterDetail{}
	payloads := masterpayloads.CampaignMasterDetailPayloads{}
	rows, err := tx.Model(&entities).Where(masterentities.CampaignMasterDetail{
		CampaignDetailId: id,
	}).First(&payloads).Rows()
	if err != nil {
		return payloads, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return payloads, nil
}

func(r *CampaignMasterRepositoryImpl) GetAllCampaignMasterCodeAndName(tx *gorm.DB)([]masterpayloads.GetHistory,*exceptionsss_test.BaseErrorResponse){
	var responses []masterpayloads.GetHistory

	err:=tx.Find(&responses).Error
	if err != nil{
		return responses,&exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return responses,nil
}

func (r *CampaignMasterRepositoryImpl) GetAllCampaignMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.CampaignMaster{}
	var GetModelResponse []masterpayloads.GetModelResponse
	var responses []masterpayloads.CampaignMasterResponse
	var internalServiceFilter, externalServiceFilter []utils.FilterCondition
	var modelCode string
	var modelName string
	responseStruct := reflect.TypeOf(masterpayloads.CampaignMasterListReponse{})

	for i := 0; i < len(filterCondition); i++ {
		flag := false
		for j := 0; j < responseStruct.NumField(); j++ {
			if filterCondition[i].ColumnField == responseStruct.Field(j).Tag.Get("parent_entity")+"."+responseStruct.Field(j).Tag.Get("json") {
				internalServiceFilter = append(internalServiceFilter, filterCondition[i])
				flag = true
				break
			}
		}
		if !flag {
			externalServiceFilter = append(externalServiceFilter, filterCondition[i])
		}
	}

	for i := 0; i < len(externalServiceFilter); i++ {
		if externalServiceFilter[i].ColumnField == "model_code" {
			modelCode = externalServiceFilter[i].ColumnValue
		}
		if externalServiceFilter[i].ColumnField == "model_name" {
			modelName = externalServiceFilter[i].ColumnValue
		}
	}

	BaseModelQuery := tx.Model(&entities).Scan(&responses)
	WhereQuery := utils.ApplyFilter(BaseModelQuery, internalServiceFilter)

	rows, err := BaseModelQuery.Scopes(pagination.Paginate(&entities, &pages, WhereQuery)).Scan(&responses).Rows()

	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	defer rows.Close()
	modelUrl := "http://127.0.0.1:8000/api/sales/unit-model-codename/?"
	if modelCode != "" {
		modelUrl += "model_code=" + modelCode
		if modelName != "" {
			modelUrl += "&"
		}
	}
	if modelName != "" {
		modelUrl += "model_name=" + modelName
	}

	errUrlmodel := utils.Get(modelUrl, &GetModelResponse, nil)

	if errUrlmodel != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlmodel,
		}
	}

	joinedData := utils.DataFrameInnerJoin(responses, GetModelResponse, "OrderTypeId")

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *CampaignMasterRepositoryImpl) GetAllCampaignMasterDetail(tx *gorm.DB, pages pagination.Pagination, id int) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	entities := []masterentities.CampaignMasterDetail{}
	response := []masterpayloads.CampaignMasterDetailPayloads{}

	BaseModelQuery := tx.Model(&entities).Where(masterentities.CampaignMasterDetail{
		CampaignId: id,
	})
	rows, err := BaseModelQuery.Scopes(pagination.Paginate(&entities, &pages, BaseModelQuery)).Scan(&response).Rows()
	if len(response) == 0 {
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

	pages.Rows = response

	return pages, nil
}

func (r *CampaignMasterRepositoryImpl) UpdateTotalCampaignMaster(tx *gorm.DB, id int, value float64) (bool) {
	err := tx.Model(&masterentities.CampaignMaster{}).Where(masterentities.CampaignMaster{
		CampaignId: id,
	}).UpdateColumns(map[string]interface{}{
		"total":           value,
		"total_vat":       value / 10,
		"total_after_vat": value + value/10,
	}).Error
	if err != nil {
		return false
	}
	return true
}
