package masterrepositoryimpl

import (
	mastercampaignmasterentities "after-sales/api/entities/master/campaign_master"
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

	result,_ := tx.Model(&entity).Where("campaign_code =?",req.CampaignCode).First(&entity).Rows()
	if result!=nil{
		return true,nil
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

func (r *CampaignMasterRepositoryImpl) PostCampaignMasterDetailFromHistory(tx *gorm.DB,id int,idhead int)(bool,*exceptionsss_test.BaseErrorResponse){
	var entity []mastercampaignmasterentities.CampaignMasterDetail

	result:=tx.Model(&entity).Where("campaign_id = ?",id).Find(&entity)
	if result.Error!= nil{
		return false, &exceptionsss_test.BaseErrorResponse{
            StatusCode: http.StatusNotFound,
            Err:        result.Error,
        }
	}
	for _,entities := range entity{
		newEntity:= mastercampaignmasterentities.CampaignMasterDetail{
		IsActive:           entities.IsActive,
		CampaignDetailId:   0,
		CampaignId:         idhead,
		LineTypeId:         entities.LineTypeId,
		Quantity:           entities.Quantity,
		OperationItemCode:  entities.OperationItemCode,
		OperationItemPrice: entities.OperationItemPrice,
		ShareBillTo:        entities.ShareBillTo,
		DiscountPercent:    entities.DiscountPercent,
		SharePercent:       entities.SharePercent,
		
	}
	err := tx.Create(&newEntity).Error

	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
		}
	}
	results:= r.UpdateTotalCampaignMaster(tx,idhead)
	if !results{
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
		}
	return true,nil
}

func (r *CampaignMasterRepositoryImpl) PostCampaignDetailMaster(tx *gorm.DB, req masterpayloads.CampaignMasterDetailPayloads) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := &mastercampaignmasterentities.CampaignMasterDetail{
		IsActive:           req.IsActive,
		CampaignDetailId:   req.CampaignDetailId,
		CampaignId:         req.CampaignId,
		LineTypeId:         req.LineTypeId,
		Quantity:           req.Quantity,
		OperationItemCode:  req.OperationItemCode,
		OperationItemPrice: req.OperationItemPrice,
		ShareBillTo:        req.ShareBillTo,
		DiscountPercent:    req.DiscountPercent,
		SharePercent:       req.SharePercent,
		
	}
	err := tx.Create(entities).Error

	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	result := r.UpdateTotalCampaignMaster(tx, req.CampaignId)
	if !result {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
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

func (r *CampaignMasterRepositoryImpl) DeactivateCampaignMasterDetail(tx *gorm.DB, ids string) (bool, *exceptionsss_test.BaseErrorResponse) {
	idSlice := strings.Split(ids, ",")

	for _, Id := range idSlice {
		var entityToUpdate mastercampaignmasterentities.CampaignMasterDetail
		result := tx.Model(&entityToUpdate).Where("campaign_detail_id = ?", Id).First(&entityToUpdate).Update("is_active",false)
		if result.Error != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        result.Error,
			}
		}
		
		results := r.UpdateTotalCampaignMaster(tx, entityToUpdate.CampaignId)
		if !results {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}
	return true, nil
}

func (r *CampaignMasterRepositoryImpl) ActivateCampaignMasterDetail(tx *gorm.DB, ids string) (bool, *exceptionsss_test.BaseErrorResponse) {
	idSlice := strings.Split(ids, ",")

	for _, id := range idSlice {
		var entityToUpdate mastercampaignmasterentities.CampaignMasterDetail
		result := tx.Model(&entityToUpdate).Where("campaign_detail_id = ?", id).First(&entityToUpdate).Update("is_active",true)
		if result.Error != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        result.Error,
			}
		}
		results := r.UpdateTotalCampaignMaster(tx, entityToUpdate.CampaignId)
		if !results {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}
	return true, nil
}

func (r *CampaignMasterRepositoryImpl) GetByIdCampaignMaster(tx *gorm.DB, id int) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse) {
	entities:=  mastercampaignmasterentities.CampaignMaster{}
	payloads:=masterpayloads.CampaignMasterResponse{}
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
	brandIdUrl := "http://10.1.32.26:8000/sales-service/api/sales/unit-brand/"+strconv.Itoa(payloads.BrandId)
	errUrlBrandId := utils.Get(brandIdUrl,&brandresponse,nil)
	if errUrlBrandId != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlBrandId,
		}
	}
	BrandJoinData := utils.DataFrameInnerJoin([]masterpayloads.CampaignMasterResponse{payloads},[]masterpayloads.GetBrandResponse{brandresponse},"BrandId")

	modelIdUrl := "http://10.1.32.26:8000/sales-service/api/sales/unit-model/"+strconv.Itoa(payloads.ModelId)
	errUrlModelId:= utils.Get(modelIdUrl,&modelresponse,nil)
	if errUrlModelId != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlModelId,
		}
	}
	ModelIdJoinData := utils.DataFrameInnerJoin(BrandJoinData,[]masterpayloads.GetModelResponse{modelresponse},"ModelId")

	fmt.Printf("BrandJoinData: %+v\n", BrandJoinData)
	fmt.Printf("ModelIdJoinData: %+v\n", ModelIdJoinData)

	return ModelIdJoinData, nil
}

func (r *CampaignMasterRepositoryImpl) GetByIdCampaignMasterDetail(tx *gorm.DB, id int) (masterpayloads.CampaignMasterDetailPayloads, *exceptionsss_test.BaseErrorResponse) {
	entities := mastercampaignmasterentities.CampaignMasterDetail{}
	payloads := masterpayloads.CampaignMasterDetailPayloads{}
	rows, err := tx.Model(&entities).Where("campaign_detail_id = ?",id).First(&payloads).Rows()
	if err != nil {
		return payloads, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return payloads, nil
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

func (r *CampaignMasterRepositoryImpl) GetAllCampaignMasterDetail(tx *gorm.DB, pages pagination.Pagination, id int) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	entities := []mastercampaignmasterentities.CampaignMasterDetail{}
	response := []masterpayloads.CampaignMasterDetailPayloads{}

	BaseModelQuery := tx.Model(&entities).Where(mastercampaignmasterentities.CampaignMasterDetail{
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

func(r *CampaignMasterRepositoryImpl) UpdateCampaignMasterDetail(tx *gorm.DB,id int,req masterpayloads.CampaignMasterDetailPayloads)(bool,*exceptionsss_test.BaseErrorResponse){
	var entities mastercampaignmasterentities.CampaignMasterDetail

	result := tx.Model(&entities).Where("campaign_detail_id = ?",id).First(&entities)
	if result.Error != nil{
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        result.Error,
		}
	}

	update:= tx.Model(&entities).Where("campaign_detail_id = ?",id).Updates(req)
	if update.Error != nil {
        return false, &exceptionsss_test.BaseErrorResponse{
            StatusCode: http.StatusInternalServerError,
            Err:        update.Error,
        }
    }

	results:= r.UpdateTotalCampaignMaster(tx,req.CampaignId)
	if !results{
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        result.Error,
		}
	}

    return true, nil

}

func (r *CampaignMasterRepositoryImpl) UpdateTotalCampaignMaster(tx *gorm.DB, id int) bool {
	var entities []mastercampaignmasterentities.CampaignMasterDetail
	var entity mastercampaignmasterentities.CampaignMaster
	var value float64
	result := tx.Model(&entities).Where(mastercampaignmasterentities.CampaignMasterDetail{
		CampaignId: id,
	}).Find(&entities)
	if result.Error != nil {
		return false
	}
	for _, detail := range entities {
		if !detail.IsActive {
			continue
		} else {
			value += detail.Quantity * detail.OperationItemPrice * ((1 - (detail.DiscountPercent/100)))
		}
	}
	result = tx.Model(&entity).Where(mastercampaignmasterentities.CampaignMaster{
		CampaignId: id,
	}).Update("total",value)
	return result.Error == nil
}
