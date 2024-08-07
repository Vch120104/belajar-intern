package masterrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	masterpackagemasterentity "after-sales/api/entities/master/package-master"

	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
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
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	if len(payloads) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	profitCenterUrl := config.EnvConfigs.GeneralServiceUrl + "profit-center?page=0&limit=10&profit_center_code=" + profitCenter

	errProfitcenterUrl := utils.Get(profitCenterUrl, &getProfitResponse, nil)

	if errProfitcenterUrl != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errProfitcenterUrl,
		}
	}

	joinedData1 := utils.DataFrameInnerJoin(payloads, getProfitResponse, "ProfitCenterId")

	unitModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model?page=0&limit=10&model_code=" + modelCode + "&model_description=" + modelDescription

	errUrlUnitModel := utils.Get(unitModelUrl, &getModelResponse, nil)

	if errUrlUnitModel != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlUnitModel,
		}
	}

	joinedData2 := utils.DataFrameInnerJoin(joinedData1, getModelResponse, "ModelId")

	VariantModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-variant?page=0&limit=10&variant_code=" + variantCode

	errUrlVariantModel := utils.Get(VariantModelUrl, &getVariantResponse, nil)

	if errUrlVariantModel != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlVariantModel,
		}
	}
	joinedData3 := utils.DataFrameInnerJoin(joinedData2, getVariantResponse, "VariantId")
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData3, &pages)
	return dataPaginate, totalPages, totalRows, nil
}

func (r *PackageMasterRepositoryImpl) GetAllPackageMasterDetailBodyshop(tx *gorm.DB, id int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entityOperation masterpackagemasterentity.PackageMasterDetailOperation
	var response masterpayloads.PackageMasterDetailOperationBodyshop
	var getlinetype masterpayloads.LineTypeCode
	rows, err := tx.Model(&entityOperation).Where(masterpackagemasterentity.PackageMasterDetailOperation{
		PackageId: id,
	}).Scan(&response).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()
	LineTypeUrl := config.EnvConfigs.GeneralServiceUrl + "line-type?line_type_id=" + strconv.Itoa(entityOperation.LineTypeId)
	errLineTypeUrl := utils.Get(LineTypeUrl, &getlinetype, nil)
	if errLineTypeUrl != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	joinedData1 := utils.DataFrameInnerJoin(rows, getlinetype, "LineTypeId")
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData1, &pages)
	return dataPaginate, totalPages, totalRows, nil
}

func (r *PackageMasterRepositoryImpl) GetAllPackageMasterDetail(tx *gorm.DB, id int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []masterentities.PackageMaster
	var itementities []masterpackagemasterentity.PackageMasterDetailItem
	var Operationentities []masterpackagemasterentity.PackageMasterDetailOperation
	var operationpayloads []masterpayloads.PackageMasterDetailOperation
	var itempayloads []masterpayloads.PackageMasterDetailItem
	var payloadheader masterpayloads.PackageMasterResponse
	rows, err := tx.Model(&entities).Where(masterentities.PackageMaster{
		PackageId: id,
	}).Scan(&payloadheader).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	if payloadheader.ProfitCenterId == 14 {
		rowsOperation, err := tx.Model(&Operationentities).Where(masterpackagemasterentity.PackageMasterDetailOperation{
			PackageId: id,
		}).Joins("JOIN mtr_operation_code ON mtr_operation_code.operation_id=mtr_package_master_detail_operation.operation_id").
			Select("mtr_package_master_detail_operation.*,mtr_operation_code.operation_code,mtr_operation_code.operation_name").
			Scan(&operationpayloads).Rows()
		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		defer rowsOperation.Close()

		// Use rowsOperation for further processing
	} else {
		combinedPayloads := make([]map[string]interface{}, 0)
		err := tx.Model(&Operationentities).Where("package_id=?", id).
			Joins("JOIN mtr_operation_code ON mtr_operation_code.operation_id=mtr_package_master_detail_operation.operation_id").
			Select("mtr_package_master_detail_operation.*,mtr_operation_code.operation_code,mtr_operation_code.operation_name").
			Scan(&operationpayloads).Error

		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		for _, op := range operationpayloads {
			combinedPayloads = append(combinedPayloads, map[string]interface{}{
				"is_active":                     op.IsActive,
				"package_detail_operation_id":   op.PackageDetailOperationId,
				"package_id":                    op.PackageId,
				"line_type_id":                  op.LineTypeId,
				"operation_id":                  op.OperationId,
				"operation_name":                op.OperationName,
				"operation_code":                op.OperationCode,
				"frt_quantity":                  op.FrtQuantity,
				"workorder_transaction_type_id": op.WorkorderTransactionTypeId,
				"job_type_id":                   op.JobTypeId,
			})
		}

		// Query for PackageMasterDetailItem
		err2 := tx.Model(&itementities).Where(masterpackagemasterentity.PackageMasterDetailItem{
			PackageId: id,
		}).Joins("JOIN mtr_item ON mtr_item.item_id=mtr_package_master_detail_item.item_id").
			Select("mtr_package_master_detail_item.*,mtr_item.item_code,mtr_item.item_name").Scan(&itempayloads).Error
		if err2 != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err2,
			}
		}
		for _, op := range itempayloads {
			combinedPayloads = append(combinedPayloads, map[string]interface{}{
				"is_active":                     op.IsActive,
				"package_detail_operation_id":   op.PackageDetailItemId,
				"package_id":                    op.PackageId,
				"line_type_id":                  op.LineTypeId,
				"item_id":                       op.ItemId,
				"item_name":                     op.ItemName,
				"item_code":                     op.ItemCode,
				"frt_quantity":                  op.FrtQuantity,
				"workorder_transaction_type_id": op.WorkorderTransactionTypeId,
				"job_type_id":                   op.JobTypeId,
			})
		}
		// Use the accumulated joinedData for further processing
		dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(combinedPayloads, &pages)

		return dataPaginate, totalPages, totalRows, nil
	}
	return nil, 0, 0, &exceptions.BaseErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Err:        nil,
	}
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
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	profitCenterUrl := config.EnvConfigs.GeneralServiceUrl + "profit-center/" + strconv.Itoa(payloads.ProfitCenterId)

	errProfitcenterUrl := utils.Get(profitCenterUrl, &getProfitResponse, nil)

	if errProfitcenterUrl != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errProfitcenterUrl,
		}
	}

	joinedData1 := utils.DataFrameInnerJoin([]masterpayloads.PackageMasterResponse{payloads}, []masterpayloads.GetProfitMaster{getProfitResponse}, "ProfitCenterId")

	unitModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(payloads.ModelId)

	errUrlUnitModel := utils.Get(unitModelUrl, &getModelResponse, nil)

	if errUrlUnitModel != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlUnitModel,
		}
	}

	joinedData2 := utils.DataFrameInnerJoin(joinedData1, []masterpayloads.UnitModelResponse{getModelResponse}, "ModelId")

	VariantModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-variant/" + strconv.Itoa(payloads.VariantId)

	errUrlVariantModel := utils.Get(VariantModelUrl, &getUnitVariantResponse, nil)

	if errUrlVariantModel != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlVariantModel,
		}
	}

	joinedData3 := utils.DataFrameInnerJoin(joinedData2, []masterpayloads.UnitVariantResponse{getUnitVariantResponse}, "VariantId")

	BrandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(payloads.BrandId)

	errUrlBrandModel := utils.Get(BrandUrl, &getBrandResponse, nil)

	if errUrlBrandModel != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlBrandModel,
		}
	}

	joinedData4 := utils.DataFrameInnerJoin(joinedData3, []masterpayloads.BrandResponse{getBrandResponse}, "BrandId")

	CurrencyUrl := config.EnvConfigs.FinanceServiceUrl + "currency-code/" + strconv.Itoa(payloads.CurrencyId)

	errUrlCurrencyModel := utils.Get(CurrencyUrl, &getCurrencyResponse, nil)

	if errUrlCurrencyModel != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlCurrencyModel,
		}
	}

	joinedData5 := utils.DataFrameInnerJoin(joinedData4, []masterpayloads.CurrencyResponse{getCurrencyResponse}, "CurrencyId")

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

func (r *PackageMasterRepositoryImpl) GetByIdPackageMasterDetail(tx *gorm.DB, id int, idheader int, LineTypeId int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	entityItem := masterpackagemasterentity.PackageMasterDetailItem{}
	entityOperation := masterpackagemasterentity.PackageMasterDetailOperation{}
	entity := masterentities.PackageMaster{}
	payloads := masterpayloads.PackageMasterResponse{}
	PayloadsOperationBodyshop := masterpayloads.PackageMasterDetailOperationBodyshop{}
	PayloadsOperation := masterpayloads.PackageMasterDetailOperation{}
	PayloadsItem := masterpayloads.PackageMasterDetailItem{}

	var getLineType masterpayloads.LineTypeCode

	rows, err := tx.Model(&entity).Where(masterentities.PackageMaster{
		PackageId: idheader,
	}).First(&payloads).Rows()

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()
	if payloads.ProfitCenterId == 1 {
		result, err := tx.Model(&entityOperation).Where(masterpackagemasterentity.PackageMasterDetailOperation{
			PackageDetailOperationId: id,
		}).Joins("JOIN mtr_operation_code ON mtr_operation_code.operation_id=mtr_package_master_detail_operation.operation_id").
			Select("mtr_package_master_detail_operation.*,mtr_operation_code.operation_code,mtr_operation_code.operation_name").
			First(&PayloadsOperationBodyshop).Rows()
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		defer rows.Close()
		LineTypeUrl := "http://172.16.5.101/general-service/v1/line-type/" + strconv.Itoa(PayloadsOperationBodyshop.LineTypeId)

		errProfitcenterUrl := utils.Get(LineTypeUrl, &getLineType, nil)

		if errProfitcenterUrl != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		joinedData1 := utils.DataFrameInnerJoin(result, getLineType, "LineTypeCode")
		result1 := joinedData1[0]

		return result1, nil
	} else {
		if LineTypeId == 1 {
			result, err := tx.Model(&entityOperation).Where(masterpackagemasterentity.PackageMasterDetailOperation{
				PackageId:                idheader,
				PackageDetailOperationId: id,
			}).Joins("JOIN mtr_operation_code ON mtr_operation_code.operation_id=mtr_package_master_detail_operation.operation_id").
				Select("mtr_package_detail_operation.*,mtr_operation_code.operation_code,mtr_operation_code.operation_name").First(&PayloadsOperation).Rows()
			if err != nil {
				return nil, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}
			defer rows.Close()
			LineTypeUrl := "http://172.16.5.101/general-service/v1/line-type/" + strconv.Itoa(PayloadsOperation.LineTypeId)

			errProfitcenterUrl := utils.Get(LineTypeUrl, &getLineType, nil)

			if errProfitcenterUrl != nil {
				return nil, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}
			joinedData1 := utils.DataFrameInnerJoin(result, getLineType, "LineTypeId")
			result1 := joinedData1[0]

			return result1, nil
		} else {
			result, err := tx.Model(&entityItem).Where(masterpackagemasterentity.PackageMasterDetailItem{
				PackageId:           idheader,
				PackageDetailItemId: id,
			}).Joins("JOIN mtr_item ON mtr_item.item_id=mtr_package_master_detail_item.item_id").
				Select("mtr_package_master_detail_item.*,mtr_item.item_code,mtr_item.item_name").
				First(&PayloadsItem).Rows()
			if err != nil {
				return nil, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}
			defer result.Close()
			LineTypeUrl := "http://172.16.5.101/general-service/v1/line-type/" + strconv.Itoa(PayloadsItem.LineTypeId)

			errProfitcenterUrl := utils.Get(LineTypeUrl, &getLineType, nil)

			if errProfitcenterUrl != nil {
				return nil, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}
			joinedData1 := utils.DataFrameInnerJoin([]masterpayloads.PackageMasterDetailItem{PayloadsItem}, []masterpayloads.LineTypeCode{getLineType}, "LineTypeId")

			result1 := joinedData1[0]
			return result1, nil
		}
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
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return entities, nil
}

func (r *PackageMasterRepositoryImpl) PostPackageMasterDetailWorkshop(tx *gorm.DB, req masterpayloads.PackageMasterDetailWorkshop) (int, *exceptions.BaseErrorResponse) {
	if req.LineTypeId == 1 {
		var rowsAffected int64
		err := tx.Model(&masterpackagemasterentity.PackageMasterDetailOperation{}).Where("package_id = ?", req.PackageId).Count(&rowsAffected).Error
		if err != nil {
			tx.Rollback()
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		entities := masterpackagemasterentity.PackageMasterDetailOperation{
			IsActive:                 req.IsActive,
			PackageDetailOperationId: req.PackageDetailItemId,
			PackageId:                req.PackageId,
			LineTypeId:               req.LineTypeId,
			OperationId:              req.PackageDetailItemId,
			Sequence:                 int(rowsAffected) + 1,
		}
		err2 := tx.Save(&entities).Error
		if err2 != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err2,
			}
		}
		return entities.PackageId, nil
	} else {
		entities := masterpackagemasterentity.PackageMasterDetailItem{
			IsActive:                   req.IsActive,
			PackageId:                  req.PackageId,
			PackageDetailItemId:        req.PackageDetailItemId,
			LineTypeId:                 req.LineTypeId,
			ItemId:                     req.ItemOperationId,
			FrtQuantity:                req.FrtQuantity,
			WorkorderTransactionTypeId: req.WorkorderTransactionTypeId,
			JobTypeId:                  req.JobTypeId,
		}
		err := tx.Save(&entities).Error
		if err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		return entities.PackageId, nil
	}
}

func (r *PackageMasterRepositoryImpl) ChangeStatusItemPackage(tx *gorm.DB, id int) (masterentities.PackageMaster, *exceptions.BaseErrorResponse) {
	var entities masterentities.PackageMaster

	result := tx.Model(&entities).
		Where("package_id = ?", id).
		First(&entities)

	if result.Error != nil {
		return masterentities.PackageMaster{}, &exceptions.BaseErrorResponse{
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
		return masterentities.PackageMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return entities, nil
}

func (r *PackageMasterRepositoryImpl) DeactivateMultiIdPackageMasterDetail(tx *gorm.DB, ids string, idHeader int) (int, *exceptions.BaseErrorResponse) {
	entities := masterentities.PackageMaster{}

	result, err := tx.Model(&entities).Where(masterentities.PackageMaster{
		PackageId: idHeader,
	}).Scan(&entities).Rows()
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer result.Close()
	if entities.ProfitCenterId == 14 {
		idSlice := strings.Split(ids, ",")

		for _, id := range idSlice {

			var entityToUpdate masterpackagemasterentity.PackageMasterDetailOperation
			result := tx.Model(&entityToUpdate).Where("package_detail_operation_id = ?", id).Update("is_active", false)
			if result.Error != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        result.Error,
				}
			}
		}
		return idHeader, nil
	} else {
		idSlice := strings.Split(ids, ",")

		for _, id := range idSlice {

			var entityToUpdate masterpackagemasterentity.PackageMasterDetailItem
			result := tx.Model(&entityToUpdate).Where("package_detail_item_id = ?", id).Update("is_active", false)
			if result.Error != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        result.Error,
				}
			}
		}
		return idHeader, nil
	}
}

func (r *PackageMasterRepositoryImpl) ActivateMultiIdPackageMasterDetail(tx *gorm.DB, ids string, idHeader int) (int, *exceptions.BaseErrorResponse) {
	entities := masterentities.PackageMaster{}

	result, err := tx.Model(&entities).Where(masterentities.PackageMaster{
		PackageId: idHeader,
	}).Scan(&entities).Rows()
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer result.Close()
	if entities.ProfitCenterId == 14 {
		idSlice := strings.Split(ids, ",")

		for _, id := range idSlice {

			var entityToUpdate masterpackagemasterentity.PackageMasterDetailOperation
			result := tx.Model(&entityToUpdate).Where("package_detail_id = ?", id).Update("is_active", true)
			if result.Error != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        result.Error,
				}
			}
		}
		return idHeader, nil
	} else {
		idSlice := strings.Split(ids, ",")

		for _, id := range idSlice {

			var entityToUpdate masterpackagemasterentity.PackageMasterDetailItem
			result := tx.Model(&entityToUpdate).Where("package_detail_item_id = ?", id).Update("is_active", true)
			if result.Error != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        result.Error,
				}
			}
		}
		return idHeader, nil
	}
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
	var detailEntitiesoperation []masterpackagemasterentity.PackageMasterDetailOperation
	err = tx.Model(&detailEntitiesoperation).Where(masterpackagemasterentity.PackageMasterDetailOperation{
		PackageId: id,
	}).Find(&detailEntitiesoperation).Error
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	for _, details := range detailEntitiesoperation {
		operationdetailentities := masterpackagemasterentity.PackageMasterDetailOperation{
			IsActive:          details.IsActive,
			PackageId:         newEntity.PackageId, // Use newEntity.PackageId
			LineTypeId:        details.LineTypeId,
			OperationId:       details.OperationId,
			FrtQuantity:       details.FrtQuantity,
			Sequence:          details.Sequence,
			TransactionTypeId: details.TransactionTypeId,
			JobTypeId:         details.JobTypeId,
		}

		err := tx.Save(&operationdetailentities).Error
		if err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		}
	}

	// Copy detailEntitiesitem
	var detailEntitiesitem []masterpackagemasterentity.PackageMasterDetailItem
	err = tx.Model(&detailEntitiesitem).Where(masterpackagemasterentity.PackageMasterDetailItem{
		PackageId: id,
	}).Find(&detailEntitiesitem).Error
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	for _, details := range detailEntitiesitem {
		operationdetailentities := masterpackagemasterentity.PackageMasterDetailItem{
			IsActive:                   details.IsActive,
			PackageId:                  newEntity.PackageId, // Use newEntity.PackageId
			LineTypeId:                 details.LineTypeId,
			ItemId:                     details.ItemId,
			FrtQuantity:                details.FrtQuantity,
			WorkorderTransactionTypeId: details.WorkorderTransactionTypeId,
			JobTypeId:                  details.JobTypeId,
		}

		err := tx.Save(&operationdetailentities).Error
		if err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		}
	}

	return id, nil
}
