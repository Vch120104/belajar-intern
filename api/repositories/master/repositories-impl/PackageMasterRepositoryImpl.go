package masterrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"

	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type PackageMasterRepositoryImpl struct {
}

func StartPackageMasterRepositoryImpl() masterrepository.PackageMasterRepository {
	return &PackageMasterRepositoryImpl{}
}

func (r *PackageMasterRepositoryImpl) GetAllPackageMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var payloads []masterpayloads.PackageMasterListResponse
	var internalServiceFilter, externalServiceFilter []utils.FilterCondition
	var getProfitResponse []masterpayloads.GetProfitMaster
	var getModelResponse []masterpayloads.UnitModelResponse
	var getVariantResponse []masterpayloads.GetVariantResponse

	var modelCode string
	var profitCenter string
	var modelDescription string
	var variantCode string
	responseStruct := reflect.TypeOf(masterpayloads.PackageMasterListResponse{})

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
		if strings.Contains(externalServiceFilter[i].ColumnField, "profit_center") {
			profitCenter = externalServiceFilter[i].ColumnValue
		} else if strings.Contains(externalServiceFilter[i].ColumnField, "model_code") {
			modelCode = externalServiceFilter[i].ColumnValue
		} else if strings.Contains(externalServiceFilter[i].ColumnField, "model_desription") {
			modelDescription = externalServiceFilter[i].ColumnValue
		} else if strings.Contains(externalServiceFilter[i].ColumnField, "variant_code") {
			variantCode = externalServiceFilter[i].ColumnValue
		}
	}
	result := tx.Model(masterentities.PackageMaster{})
	whereQuery := utils.ApplyFilter(result, internalServiceFilter)
	rows, err := whereQuery.Scan(&payloads).Rows()

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	if len(payloads) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	profitCenterUrl := config.EnvConfigs.GeneralServiceUrl + "profit-center?page=0&limit=10&profit_center_code=" + profitCenter

	errProfitcenterUrl := utils.Get(profitCenterUrl, &getProfitResponse, nil)

	if errProfitcenterUrl != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errProfitcenterUrl,
		}
	}

	joinedData1, errdf := utils.DataFrameInnerJoin(payloads, getProfitResponse, "ProfitCenterId")

	if errdf != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errdf,
		}
	}

	unitModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model?page=0&limit=10&model_code=" + modelCode + "&model_description=" + modelDescription

	errUrlUnitModel := utils.Get(unitModelUrl, &getModelResponse, nil)

	if errUrlUnitModel != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlUnitModel,
		}
	}

	joinedData2, errdf := utils.DataFrameInnerJoin(joinedData1, getModelResponse, "ModelId")
	if errdf != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errdf,
		}
	}

	VariantModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-variant?page=0&limit=10&variant_code=" + variantCode

	errUrlVariantModel := utils.Get(VariantModelUrl, &getVariantResponse, nil)

	if errUrlVariantModel != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlVariantModel,
		}
	}
	joinedData3, errdf := utils.DataFrameInnerJoin(joinedData2, getVariantResponse, "VariantId")
	if errdf != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errdf,
		}
	}

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData3, &pages)
	return dataPaginate, totalPages, totalRows, nil
}

func (r *PackageMasterRepositoryImpl) GetAllPackageMasterDetail(tx *gorm.DB, id int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []masterentities.PackageMasterDetail
	var Detailpayloads []masterpayloads.PackageMasterDetail
	var operationpayloads []masterpayloads.PackageMasterDetailOperation
	var itempayloads []masterpayloads.PackageMasterDetailItem
	rows, err := tx.Model(&entities).Where(masterentities.PackageMasterDetail{
		PackageId: id,
	}).Scan(&Detailpayloads).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()
	combinedPayloads := make([]map[string]interface{}, 0)
	for _, packdetail := range Detailpayloads {
		if packdetail.LineTypeId == 5 { //operation line type id
			err := tx.Model(&entities).Where("package_id=?", id).
				Joins("join mtr_item_operation on mtr_item_operation.item_operation_id = mtr_package_detail.item_operation_id").
				Joins("JOIN mtr_operation_model_mapping ON mtr_operation_model_mapping.operation_model_mapping_id=mtr_item_operation.tem_operation_id").
				Joins("join mtr_operation_code on mtr_operation_code.operation_id=mtr_operation_model_mapping.operation_id").
				Select("mtr_package_master_detail.*,mtr_operation_code.operation_code,mtr_operation_code.operation_name").
				Scan(&operationpayloads).Error

			if err != nil {
				return nil, 0, 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Err:        err,
				}
			}
			for _, op := range operationpayloads {
				combinedPayloads = append(combinedPayloads, map[string]interface{}{
					"is_active":                     op.IsActive,
					"package_detail_operation_id":   op.PackageDetailId,
					"package_id":                    op.PackageId,
					"line_type_id":                  op.LineTypeId,
					"item_operation_id":             op.ItemOperationId,
					"operation_name":                op.OperationName,
					"operation_code":                op.OperationCode,
					"frt_quantity":                  op.FrtQuantity,
					"workorder_transaction_type_id": op.WorkorderTransactionTypeId,
					"job_type_id":                   op.JobTypeId,
				})
			}

		} else {
			err2 := tx.Model(&entities).Where(masterentities.PackageMasterDetail{
				PackageId: id,
			}).Joins("join mtr_item_operation on mtr_item_operation.item_operation_id = mtr_package_detail.item_operation_id").
				Joins("JOIN mtr_item ON mtr_item.item_id=mtr_item_operation.item_id").
				Select("mtr_package_master_detail.*,mtr_item.item_code,mtr_item.item_name").Scan(&itempayloads).Error
			if err2 != nil {
				return nil, 0, 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Err:        err2,
				}
			}
			for _, it := range itempayloads {
				combinedPayloads = append(combinedPayloads, map[string]interface{}{
					"is_active":                     it.IsActive,
					"package_detail_operation_id":   it.PackageDetailId,
					"package_id":                    it.PackageId,
					"line_type_id":                  it.LineTypeId,
					"item_operation_id":             it.ItemOperationId,
					"item_name":                     it.ItemName,
					"item_code":                     it.ItemCode,
					"frt_quantity":                  it.FrtQuantity,
					"workorder_transaction_type_id": it.WorkorderTransactionTypeId,
					"job_type_id":                   it.JobTypeId,
				})
			}
		}

	}
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(combinedPayloads, &pages)
	return dataPaginate, totalPages, totalRows, nil
}

func (r *PackageMasterRepositoryImpl) GetByIdPackageMaster(tx *gorm.DB, id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	entity := masterentities.PackageMaster{}
	payloads := masterpayloads.PackageMasterResponse{}

	var getBrandResponse masterpayloads.BrandResponse
	var getModelResponse masterpayloads.UnitModelResponse
	var getUnitVariantResponse masterpayloads.UnitVariantResponse
	var getProfitResponse masterpayloads.GetProfitMaster
	var getCurrencyResponse masterpayloads.CurrencyResponse

	rows, err := tx.Model(&entity).Where(
		masterentities.PackageMaster{
			PackageId: id,
		}).First(&payloads).Rows()

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	profitCenterUrl := config.EnvConfigs.GeneralServiceUrl + "profit-center/" + strconv.Itoa(payloads.ProfitCenterId)

	errProfitcenterUrl := utils.Get(profitCenterUrl, &getProfitResponse, nil)

	if errProfitcenterUrl != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errProfitcenterUrl,
		}
	}

	joinedData1, errdf := utils.DataFrameInnerJoin([]masterpayloads.PackageMasterResponse{payloads}, []masterpayloads.GetProfitMaster{getProfitResponse}, "ProfitCenterId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errdf,
		}
	}

	unitModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(payloads.ModelId)

	errUrlUnitModel := utils.Get(unitModelUrl, &getModelResponse, nil)

	if errUrlUnitModel != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlUnitModel,
		}
	}

	joinedData2, errdf := utils.DataFrameInnerJoin(joinedData1, []masterpayloads.UnitModelResponse{getModelResponse}, "ModelId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errdf,
		}
	}

	VariantModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-variant/" + strconv.Itoa(payloads.VariantId)

	errUrlVariantModel := utils.Get(VariantModelUrl, &getUnitVariantResponse, nil)

	if errUrlVariantModel != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlVariantModel,
		}
	}

	joinedData3, errdf := utils.DataFrameInnerJoin(joinedData2, []masterpayloads.UnitVariantResponse{getUnitVariantResponse}, "VariantId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errdf,
		}
	}

	BrandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(payloads.BrandId)

	errUrlBrandModel := utils.Get(BrandUrl, &getBrandResponse, nil)

	if errUrlBrandModel != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlBrandModel,
		}
	}

	joinedData4, errdf := utils.DataFrameInnerJoin(joinedData3, []masterpayloads.BrandResponse{getBrandResponse}, "BrandId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errdf,
		}
	}

	CurrencyUrl := config.EnvConfigs.FinanceServiceUrl + "currency-code/" + strconv.Itoa(payloads.CurrencyId)

	errUrlCurrencyModel := utils.Get(CurrencyUrl, &getCurrencyResponse, nil)

	if errUrlCurrencyModel != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlCurrencyModel,
		}
	}

	joinedData5, errdf := utils.DataFrameInnerJoin(joinedData4, []masterpayloads.CurrencyResponse{getCurrencyResponse}, "CurrencyId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errdf,
		}
	}

	result := map[string]interface{}{
		"brand_code":          joinedData5[0]["BrandCode"],
		"brand_id":            joinedData5[0]["BrandId"],
		"brand_name":          joinedData5[0]["BrandName"],
		"currency_code":       joinedData5[0]["CurrencyCode"],
		"currency_id":         joinedData5[0]["CurrencyId"],
		"is_active":           joinedData5[0]["IsActive"],
		"item_group_id":       joinedData5[0]["ItemGroupId"],
		"model_code":          joinedData5[0]["ModelCode"],
		"model_description":   joinedData5[0]["ModelDescription"],
		"model_id":            joinedData5[0]["ModelId"],
		"package_code":        joinedData5[0]["PackageCode"],
		"package_name":        joinedData5[0]["PackageName"],
		"package_id":          joinedData5[0]["PackageId"],
		"package_price":       joinedData5[0]["PackagePrice"],
		"package_remark":      joinedData5[0]["PackageRemark"],
		"package_set":         joinedData5[0]["PackageSet"],
		"profit_center_id":    joinedData5[0]["ProfitCenterId"],
		"profit_center_name":  joinedData5[0]["ProfitCenterName"],
		"tax_type_id":         joinedData5[0]["TaxTypeId"],
		"variant_description": joinedData5[0]["VariantDescription"],
		"variant_id":          joinedData5[0]["VariantId"],
	}
	return result, nil
}

func (r *PackageMasterRepositoryImpl) GetByIdPackageMasterDetail(tx *gorm.DB, id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	var entity masterentities.PackageMasterDetail
	var detailpayloads masterpayloads.PackageMasterDetail
	var getLineType masteritempayloads.LineTypeResponse
	err := tx.Model(&entity).Where("package_master_detail_id=?", id).Scan(entity).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	if entity.LineTypeId == 5 {
		err := tx.Model(&entity).Where("package_id =?", id).
			Joins("JOIN mtr_item ON mtr_item.item_id=mtr_package_master_detail_item.item_id").
			Select("mtr_package_master_detail_item.*,mtr_item.item_code,mtr_item.item_name").
			First(&detailpayloads).Error
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
		LineTypeUrl := config.EnvConfigs.GeneralServiceUrl + "line-type/" + strconv.Itoa(detailpayloads.LineTypeId)

		errProfitcenterUrl := utils.Get(LineTypeUrl, &getLineType, nil)

		if errProfitcenterUrl != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
		joinedData1, errdf := utils.DataFrameInnerJoin([]masterpayloads.PackageMasterDetail{detailpayloads}, []masteritempayloads.LineTypeResponse{getLineType}, "LineTypeId")

		if errdf != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errdf,
			}
		}
		response := map[string]interface{}{
			"is_active":                     joinedData1[0]["IsActive"],
			"package_detail_operation_id":   joinedData1[0]["PackageDetailId"],
			"package_id":                    joinedData1[0]["PackageId"],
			"line_type_id":                  joinedData1[0]["LineTypeId"],
			"item_operation_id":             joinedData1[0]["ItemOperationId"],
			"item_name":                     joinedData1[0]["ItemName"],
			"item_code":                     joinedData1[0]["ItemCode"],
			"frt_quantity":                  joinedData1[0]["FrtQuantity"],
			"workorder_transaction_type_id": joinedData1[0]["WorkorderTransactionTypeId"],
			"job_type_id":                   joinedData1[0]["JobTypeId"],
		}
		return response, nil
	} else {
		err := tx.Model(&entity).Where("package_id=?", id).
			Joins("join mtr_item_operation on mtr_item_operation.item_operation_id = mtr_package_detail.item_operation_id").
			Joins("JOIN mtr_operation_model_mapping ON mtr_operation_model_mapping.operation_model_mapping_id=mtr_item_operation.tem_operation_id").
			Joins("join mtr_operation_code on mtr_operation_code.operation_id=mtr_operation_model_mapping.operation_id").
			Select("mtr_package_master_detail.*,mtr_operation_code.operation_code,mtr_operation_code.operation_name").
			First(&detailpayloads).Error
		LineTypeUrl := config.EnvConfigs.GeneralServiceUrl + "line-type/" + strconv.Itoa(detailpayloads.LineTypeId)

		errProfitcenterUrl := utils.Get(LineTypeUrl, &getLineType, nil)

		if errProfitcenterUrl != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
		joinedData1, errdf := utils.DataFrameInnerJoin([]masterpayloads.PackageMasterDetail{detailpayloads}, []masteritempayloads.LineTypeResponse{getLineType}, "LineTypeId")

		if errdf != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errdf,
			}
		}

		response := map[string]interface{}{
			"is_active":                     joinedData1[0]["IsActive"],
			"package_detail_operation_id":   joinedData1[0]["PackageDetailId"],
			"package_id":                    joinedData1[0]["PackageId"],
			"line_type_id":                  joinedData1[0]["LineTypeId"],
			"item_operation_id":             joinedData1[0]["ItemOperationId"],
			"operation_name":                joinedData1[0]["OperationName"],
			"operation_code":                joinedData1[0]["OperationCode"],
			"frt_quantity":                  joinedData1[0]["FrtQuantity"],
			"workorder_transaction_type_id": joinedData1[0]["WorkorderTransactionTypeId"],
			"job_type_id":                   joinedData1[0]["JobTypeId"],
		}
		return response, nil
	}
}

func (r *PackageMasterRepositoryImpl) PostpackageMaster(tx *gorm.DB, req masterpayloads.PackageMasterResponse) (masterentities.PackageMaster, *exceptions.BaseErrorResponse) {
	entities := masterentities.PackageMaster{
		IsActive:       req.IsActive,
		PackageId:      req.PackageId,
		ItemGroupId:    req.ItemGroupId,
		PackageName:    req.PackageName,
		PackageCode:    req.PackageCode,
		BrandId:        req.BrandId,
		ModelId:        req.ModelId,
		VariantId:      req.VariantId,
		ProfitCenterId: req.ProfitCenterId,
		PackageSet:     req.PackageSet,
		CurrencyId:     req.CurrencyId,
		PackagePrice:   req.PackagePrice,
		TaxTypeId:      req.TaxTypeId,
		PackageRemark:  req.PackageRemark,
	}
	err := tx.Save(&entities).Error
	if err != nil {
		return masterentities.PackageMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return entities, nil
}

func (r *PackageMasterRepositoryImpl) PostPackageMasterDetail(tx *gorm.DB, req masterpayloads.PackageMasterDetail) (masterentities.PackageMasterDetail, *exceptions.BaseErrorResponse) {
	entities := masterentities.PackageMasterDetail{
		IsActive:                   req.IsActive,
		PackageId:                  req.PackageId,
		LineTypeId:                 req.LineTypeId,
		ItemOperationId:            req.ItemOperationId,
		FrtQuantity:                req.FrtQuantity,
		WorkorderTransactionTypeId: req.WorkorderTransactionTypeId,
		JobTypeId:                  req.JobTypeId,
	}
	err := tx.Save(&entities).Error
	if err != nil {
		return masterentities.PackageMasterDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return entities, nil
}

func (r *PackageMasterRepositoryImpl) ChangeStatusItemPackage(tx *gorm.DB, id int) (masterentities.PackageMaster, *exceptions.BaseErrorResponse) {
	var entities masterentities.PackageMaster

	result := tx.Model(&entities).
		Where("package_id = ?", id).
		First(&entities)

	if result.Error != nil {
		return masterentities.PackageMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
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
		return masterentities.PackageMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        result.Error,
		}
	}

	return entities, nil
}

func (r *PackageMasterRepositoryImpl) DeactivateMultiIdPackageMasterDetail(tx *gorm.DB, ids string) (bool, *exceptions.BaseErrorResponse) {
	entities := masterentities.PackageMasterDetail{}
	id := strings.Split(ids, ",")
	for _, ides := range id {
		err := tx.Select(&entities).Where("package_detail_id=?", ides).Update("is_active", false).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
			}
		}
	}
	return true, nil
}

func (r *PackageMasterRepositoryImpl) ActivateMultiIdPackageMasterDetail(tx *gorm.DB, ids string) (bool, *exceptions.BaseErrorResponse) {
	entities := masterentities.PackageMasterDetail{}
	id := strings.Split(ids, ",")
	for _, ides := range id {
		err := tx.Select(&entities).Where("package_detail_id=?", ides).Update("is_active", true).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
			}
		}
	}
	return true, nil
}

func (r *PackageMasterRepositoryImpl) CopyToOtherModel(tx *gorm.DB, id int, code string, modelId int) (int, *exceptions.BaseErrorResponse) {
	var entity masterentities.PackageMaster
	var payloads masterpayloads.PackageMasterResponse

	// Retrieve the package master entity
	err := tx.Model(&entity).Where(masterentities.PackageMaster{
		PackageId: id,
	}).First(&payloads).Error
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	// Save the new package master entity
	newEntity := masterentities.PackageMaster{
		IsActive:       payloads.IsActive,
		PackageCode:    code,
		ItemGroupId:    payloads.ItemGroupId,
		PackageName:    payloads.PackageName,
		BrandId:        payloads.BrandId,
		ModelId:        modelId,
		VariantId:      payloads.VariantId,
		ProfitCenterId: payloads.ProfitCenterId,
		PackageSet:     payloads.PackageSet,
		CurrencyId:     payloads.CurrencyId,
		PackagePrice:   payloads.PackagePrice,
		TaxTypeId:      payloads.TaxTypeId,
		PackageRemark:  payloads.PackageRemark,
	}

	err = tx.Save(&newEntity).Error
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}

	// Copy detailEntitiesoperation
	var detailentities []masterentities.PackageMasterDetail
	err2 := tx.Model(&detailentities).Where("package_id=?", id).Scan(detailentities).Error
	if err2 != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err2,
		}
	}
	for _, detail := range detailentities {
		entities := masterentities.PackageMasterDetail{
			IsActive:                   detail.IsActive,
			PackageId:                  newEntity.PackageId,
			LineTypeId:                 detail.LineTypeId,
			ItemOperationId:            detail.ItemOperationId,
			FrtQuantity:                detail.FrtQuantity,
			WorkorderTransactionTypeId: detail.WorkorderTransactionTypeId,
			JobTypeId:                  detail.JobTypeId,
		}
		err := tx.Save(&entities).Error
		if err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
			}
		}
	}
	return id, nil
}
