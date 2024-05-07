package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterpackagemasterentity "after-sales/api/entities/master/package-master"

	exceptionsss_test "after-sales/api/expectionsss"
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

func (r *PackageMasterRepositoryImpl) GetAllPackageMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
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
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	if len(payloads) == 0 {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	profitCenterUrl := "http://10.1.32.26:8000/general-service/api/general/profit-center?page=0&limit=10&profit_center_code=" + profitCenter

	errProfitcenterUrl := utils.Get(profitCenterUrl, &getProfitResponse, nil)

	if errProfitcenterUrl != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	joinedData1 := utils.DataFrameInnerJoin(payloads, getProfitResponse, "ProfitCenterId")

	unitModelUrl := "http://10.1.32.26:8000/sales-service/api/sales/unit-model?page=0&limit=10&model_code=" + modelCode + "&model_description=" + modelDescription

	errUrlUnitModel := utils.Get(unitModelUrl, &getModelResponse, nil)

	if errUrlUnitModel != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	joinedData2 := utils.DataFrameInnerJoin(joinedData1, getModelResponse, "ModelId")

	VariantModelUrl := "http://10.1.32.26:8000/sales-service/api/sales/unit-variant?page=0&limit=10&variant_code=" + variantCode

	errUrlVariantModel := utils.Get(VariantModelUrl, &getVariantResponse, nil)

	if errUrlVariantModel != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	joinedData3 := utils.DataFrameInnerJoin(joinedData2, getVariantResponse, "VariantId")
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData3, &pages)
	return dataPaginate, totalPages, totalRows, nil
}

func (r *PackageMasterRepositoryImpl) GetAllPackageMasterDetailBodyshop(tx *gorm.DB, id int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	var entityOperation masterpackagemasterentity.PackageMasterDetailOperation
	var response masterpayloads.PackageMasterDetailOperationBodyshop
	var getlinetype masterpayloads.LineTypeCode
	rows, err := tx.Model(&entityOperation).Where(masterpackagemasterentity.PackageMasterDetailOperation{
		PackageId: id,
	}).Scan(&response).Rows()
	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()
	profitCenterUrl := "http://127.0.0.1:8000/api/general/profit-center/" + strconv.Itoa(entityOperation.LineTypeId)
	errLineTypeUrl := utils.Get(profitCenterUrl, &getlinetype, nil)
	if errLineTypeUrl != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	joinedData1 := utils.DataFrameInnerJoin(rows, getlinetype, "LineTypeId")
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData1, &pages)
	return dataPaginate, totalPages, totalRows, nil
}

func (r *PackageMasterRepositoryImpl) GetAllPackageMasterDetail(tx *gorm.DB, id int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	var entities masterentities.PackageMaster
	var itementities masterpackagemasterentity.PackageMasterDetailItem
	var Operationentities masterpackagemasterentity.PackageMasterDetailOperation
	var entityItem masterpayloads.PackageMasterDetailItem
	var entityOperation masterpayloads.PackageMasterDetailOperation
	var entitybodyshop masterpayloads.PackageMasterDetailOperationBodyshop
	var payloadheader masterpayloads.PackageMasterResponse
	var getlinetype masterpayloads.LineTypeCode
	rows, err := tx.Model(&entities).Where(masterentities.PackageMaster{
		PackageId: id,
	}).Scan(&payloadheader).Rows()
	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	if payloadheader.ProfitCenterId == 14 { //harus diganti sesuai dengan i profit center bodyshop
		rows, err := tx.Model(&Operationentities).Where(masterpackagemasterentity.PackageMasterDetailOperation{
			PackageId: id,
		}).Scan(&entitybodyshop).Rows()
		if err != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		defer rows.Close()
		profitCenterUrl := "http://127.0.0.1:8000/api/general/profit-center/" + strconv.Itoa(entityOperation.LineTypeId)
		errLineTypeUrl := utils.Get(profitCenterUrl, &getlinetype, nil)
		if errLineTypeUrl != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		joinedData1 := utils.DataFrameInnerJoin(rows, getlinetype, "LineTypeId")
		dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData1, &pages)
		return dataPaginate, totalPages, totalRows, nil
	} else {
		rows, err := tx.Model(&Operationentities).Where(masterpackagemasterentity.PackageMasterDetailOperation{
			PackageId: id,
		}).Scan(&entityOperation).Rows()
		if err != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		defer rows.Close()
		profitCenterUrl := "http://127.0.0.1:8000/api/general/profit-center/" + strconv.Itoa(entityOperation.LineTypeId)
		errLineTypeUrl := utils.Get(profitCenterUrl, &getlinetype, nil)
		if errLineTypeUrl != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		joinedData1 := utils.DataFrameInnerJoin(rows, getlinetype, "LineTypeId")

		rows2, err2 := tx.Model(&itementities).Where(masterpackagemasterentity.PackageMasterDetailOperation{
			PackageId: id,
		}).Scan(&entityItem).Rows()
		if err2 != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		defer rows2.Close()
		profitCenterUrl2 := "http://127.0.0.1:8000/api/general/profit-center/" + strconv.Itoa(entityOperation.LineTypeId)
		errLineTypeUrl2 := utils.Get(profitCenterUrl2, &getlinetype, nil)
		if errLineTypeUrl2 != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		joinedData2 := utils.DataFrameInnerJoin(rows2, getlinetype, "LineTypeId")
		combinedData := append(joinedData1, joinedData2...)
		dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(combinedData, &pages)
		return dataPaginate, totalPages, totalRows, nil
	}
}

func (r *PackageMasterRepositoryImpl) GetByIdPackageMaster(tx *gorm.DB, id int) (map[string]interface{}, *exceptionsss_test.BaseErrorResponse) {
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
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	profitCenterUrl := "http://10.1.32.26:8000/general-service/api/general/profit-center/" + strconv.Itoa(payloads.ProfitCenterId)

	errProfitcenterUrl := utils.Get(profitCenterUrl, &getProfitResponse, nil)

	if errProfitcenterUrl != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errProfitcenterUrl,
		}
	}

	joinedData1 := utils.DataFrameInnerJoin([]masterpayloads.PackageMasterResponse{payloads}, []masterpayloads.GetProfitMaster{getProfitResponse}, "ProfitCenterId")

	unitModelUrl := "http://10.1.32.26:8000/sales-service/api/sales/unit-model/" + strconv.Itoa(payloads.ModelId)

	errUrlUnitModel := utils.Get(unitModelUrl, &getModelResponse, nil)

	if errUrlUnitModel != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlUnitModel,
		}
	}

	joinedData2 := utils.DataFrameInnerJoin(joinedData1, []masterpayloads.UnitModelResponse{getModelResponse}, "ModelId")

	VariantModelUrl := "http://10.1.32.26:8000/sales-service/api/sales/unit-variant/" + strconv.Itoa(payloads.VariantId)

	errUrlVariantModel := utils.Get(VariantModelUrl, &getUnitVariantResponse, nil)

	if errUrlVariantModel != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlVariantModel,
		}
	}

	joinedData3 := utils.DataFrameInnerJoin(joinedData2, []masterpayloads.UnitVariantResponse{getUnitVariantResponse}, "VariantId")

	BrandUrl := "http://10.1.32.26:8000/sales-service/api/sales/unit-brand/" + strconv.Itoa(payloads.BrandId)

	errUrlBrandModel := utils.Get(BrandUrl, &getBrandResponse, nil)

	if errUrlBrandModel != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlBrandModel,
		}
	}

	joinedData4 := utils.DataFrameInnerJoin(joinedData3, []masterpayloads.BrandResponse{getBrandResponse}, "BrandId")

	CurrencyUrl := "http://10.1.32.26:8000/finance-service/v1/currency-code/" + strconv.Itoa(payloads.CurrencyId)

	errUrlCurrencyModel := utils.Get(CurrencyUrl, &getCurrencyResponse, nil)

	if errUrlCurrencyModel != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlCurrencyModel,
		}
	}

	joinedData5 := utils.DataFrameInnerJoin(joinedData4, []masterpayloads.CurrencyResponse{getCurrencyResponse}, "CurrencyId")
	result := joinedData5[0]
	return result, nil
}

func (r *PackageMasterRepositoryImpl) GetByIdPackageMasterDetail(tx *gorm.DB, id int, idheader int, LineTypeId int) (map[string]interface{}, *exceptionsss_test.BaseErrorResponse) {
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
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()
	if payloads.ProfitCenterId == 1 {
		result, err := tx.Model(&entityOperation).Where(masterpackagemasterentity.PackageMasterDetailOperation{
			PackageDetailOperationId: id,
		}).First(&PayloadsOperationBodyshop).Rows()
		if err != nil {
			return nil, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		defer rows.Close()
		LineTypeUrl := "http://10.1.32.26:8000/general-service/api/general/line-type/" + strconv.Itoa(PayloadsOperationBodyshop.LineTypeId)

		errProfitcenterUrl := utils.Get(LineTypeUrl, &getLineType, nil)

		if errProfitcenterUrl != nil {
			return nil, &exceptionsss_test.BaseErrorResponse{
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
			}).First(&PayloadsOperation).Rows()
			if err != nil {
				return nil, &exceptionsss_test.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}
			defer rows.Close()
			LineTypeUrl := "http://10.1.32.26:8000/general-service/api/general/line-type/" + strconv.Itoa(PayloadsOperation.LineTypeId)

			errProfitcenterUrl := utils.Get(LineTypeUrl, &getLineType, nil)

			if errProfitcenterUrl != nil {
				return nil, &exceptionsss_test.BaseErrorResponse{
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
			}).First(&PayloadsItem).Rows()
			if err != nil {
				return nil, &exceptionsss_test.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}
			defer result.Close()
			LineTypeUrl := "http://10.1.32.26:8000/general-service/api/general/line-type/" + strconv.Itoa(PayloadsItem.LineTypeId)

			errProfitcenterUrl := utils.Get(LineTypeUrl, &getLineType, nil)

			if errProfitcenterUrl != nil {
				return nil, &exceptionsss_test.BaseErrorResponse{
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

func (r *PackageMasterRepositoryImpl) PostpackageMaster(tx *gorm.DB, req masterpayloads.PackageMasterResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
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
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return true, nil
}

func (r *PackageMasterRepositoryImpl) PostPackageMasterDetailBodyshop(tx *gorm.DB, req masterpayloads.PackageMasterDetailOperationBodyshop, id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var rowsAffected int64
	if err := tx.Model(&masterpackagemasterentity.PackageMasterDetailOperation{}).Where(masterpackagemasterentity.PackageMasterDetailOperation{
		PackageId: req.PackageId,
	}).Count(&rowsAffected).Error; err != nil {
		tx.Rollback()
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	entities := masterpackagemasterentity.PackageMasterDetailOperation{
		IsActive:                 req.IsActive,
		PackageDetailOperationId: req.PackageDetailOperationId,
		PackageId:                id,
		LineTypeId:               req.LineTypeId,
		OperationId:              req.OperationId,
		Sequence:                 int(rowsAffected) + 1,
	}
	err := tx.Save(&entities).Error
	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return true, nil
}

func (r *PackageMasterRepositoryImpl) PostPackageMasterDetailWorkshop(tx *gorm.DB, req masterpayloads.PackageMasterDetailWorkshop) (bool, *exceptionsss_test.BaseErrorResponse) {
	if req.LineTypeId == 1 {
		entities := masterpackagemasterentity.PackageMasterDetailOperation{
			IsActive:                 req.IsActive,
			PackageId:                req.PackageId,
			PackageDetailOperationId: req.PackageDetailItemId,
			LineTypeId:               req.LineTypeId,
			OperationId:              req.ItemOperationId,
			FrtQuantity:              req.FrtQuantity,
			TransactionTypeId:        req.WorkorderTransactionTypeId,
			JobTypeId:                req.JobTypeId,
		}
		err := tx.Save(&entities).Error
		if err != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		return true, nil
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
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		return true, nil
	}
}

func (r *PackageMasterRepositoryImpl) ChangeStatusItemPackage(tx *gorm.DB, id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities masterentities.PackageMaster

	result := tx.Model(&entities).
		Where("package_id = ?", id).
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

func (r *PackageMasterRepositoryImpl) DeactivateMultiIdPackageMasterDetail(tx *gorm.DB, ids string, idHeader int) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.PackageMaster{}

	result, err := tx.Model(&entities).Where(masterentities.PackageMaster{
		PackageId: idHeader,
	}).Scan(&entities).Rows()
	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
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
				return false, &exceptionsss_test.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        result.Error,
				}
			}
		}
		return true, nil
	} else {
		idSlice := strings.Split(ids, ",")

		for _, id := range idSlice {

			var entityToUpdate masterpackagemasterentity.PackageMasterDetailItem
			result := tx.Model(&entityToUpdate).Where("package_detail_item_id = ?", id).Update("is_active", false)
			if result.Error != nil {
				return false, &exceptionsss_test.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        result.Error,
				}
			}
		}
		return true, nil
	}
}

func (r *PackageMasterRepositoryImpl) ActivateMultiIdPackageMasterDetail(tx *gorm.DB, ids string, idHeader int) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.PackageMaster{}

	result, err := tx.Model(&entities).Where(masterentities.PackageMaster{
		PackageId: idHeader,
	}).Scan(&entities).Rows()
	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
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
				return false, &exceptionsss_test.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        result.Error,
				}
			}
		}
		return true, nil
	} else {
		idSlice := strings.Split(ids, ",")

		for _, id := range idSlice {

			var entityToUpdate masterpackagemasterentity.PackageMasterDetailItem
			result := tx.Model(&entityToUpdate).Where("package_detail_item_id = ?", id).Update("is_active", true)
			if result.Error != nil {
				return false, &exceptionsss_test.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        result.Error,
				}
			}
		}
		return true, nil
	}
}

func (r *PackageMasterRepositoryImpl) CopyToOtherModel(tx *gorm.DB, id int, code string, modelId int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entity masterentities.PackageMaster
	var payloads masterpayloads.PackageMasterResponse

	// Retrieve the package master entity
	err := tx.Model(&entity).Where(masterentities.PackageMaster{
		PackageId: id,
	}).First(&payloads).Error
	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
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
		return false, &exceptionsss_test.BaseErrorResponse{
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
		return false, &exceptionsss_test.BaseErrorResponse{
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
			return false, &exceptionsss_test.BaseErrorResponse{
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
		return false, &exceptionsss_test.BaseErrorResponse{
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
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		}
	}

	return true, nil
}
