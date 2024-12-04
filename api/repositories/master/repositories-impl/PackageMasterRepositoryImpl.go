package masterrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	"math"

	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type PackageMasterRepositoryImpl struct {
}

func StartPackageMasterRepositoryImpl() masterrepository.PackageMasterRepository {
	return &PackageMasterRepositoryImpl{}
}

func (r *PackageMasterRepositoryImpl) GetAllPackageMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var payloads []masterentities.PackageMaster
	var profitCenterName, variantCode, modelDescription, modelCode string

	internalFilter, externalFilter := utils.DefineInternalExternalFilter(filterCondition, masterentities.PackageMaster{})

	result := tx.Model(&masterentities.PackageMaster{})

	whereQuery := utils.ApplyFilter(result, internalFilter)

	if len(externalFilter) > 0 {
		for _, filter := range externalFilter {

			if filter.ColumnField == "profit_center_name" {
				profitCenterName = filter.ColumnValue
			}
		}

		if profitCenterName != "" {
			profitCenterParams := generalserviceapiutils.ProfitCenterParams{
				ProfitCenterName: profitCenterName,
			}

			profitCenterResponse, profitCenterErr := generalserviceapiutils.GetAllProfitCenter(profitCenterParams)
			if profitCenterErr != nil {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: profitCenterErr.StatusCode,
					Message:    "Error fetching profit center data",
					Err:        profitCenterErr.Err,
				}
			}

			var profitCenterIds []int
			for _, pc := range profitCenterResponse {
				profitCenterIds = append(profitCenterIds, pc.ProfitCenterId)
			}

			if len(profitCenterIds) > 0 {
				whereQuery = whereQuery.Where("profit_center_id IN ?", profitCenterIds)
			} else {
				pages.Rows = []map[string]interface{}{}
				return pages, nil
			}
		}
	}

	if variantCode != "" {
		variantParams := salesserviceapiutils.UnitVariantParams{
			VariantCode: variantCode,
		}

		variantResponse, variantErr := salesserviceapiutils.GetAllUnitVariant(variantParams)
		if variantErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: variantErr.StatusCode,
				Message:    "Error fetching variant data",
				Err:        variantErr.Err,
			}
		}

		var variantIds []int
		for _, v := range variantResponse {
			variantIds = append(variantIds, v.VariantId)
		}

		if len(variantIds) > 0 {
			whereQuery = whereQuery.Where("variant_id IN ?", variantIds)
		} else {
			pages.Rows = []map[string]interface{}{}
			return pages, nil
		}
	}

	if modelDescription != "" || modelCode != "" {
		modelParams := salesserviceapiutils.UnitModelParams{
			ModelDescription: modelDescription,
			ModelCode:        modelCode,
		}

		modelResponse, modelErr := salesserviceapiutils.GetAllUnitModel(modelParams)
		if modelErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: modelErr.StatusCode,
				Message:    "Error fetching model data",
				Err:        modelErr.Err,
			}
		}

		var modelIds []int
		for _, m := range modelResponse {
			modelIds = append(modelIds, m.ModelId)
		}

		if len(modelIds) > 0 {
			whereQuery = whereQuery.Where("model_id IN ?", modelIds)
		} else {
			pages.Rows = []map[string]interface{}{}
			return pages, nil
		}
	}

	var totalRows int64
	if err := whereQuery.Count(&totalRows).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pages.GetLimit())))
	pages.TotalPages = totalPages

	err := whereQuery.Order("package_id ASC").Offset(pages.GetOffset()).Limit(pages.GetLimit()).Find(&payloads).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(payloads) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, response := range payloads {
		// Fetch profit center data
		profitCenterData, profitCenterErr := generalserviceapiutils.GetProfitCenterById(response.ProfitCenterId)
		if profitCenterErr != nil {
			return pages, profitCenterErr
		}

		// Fetch model data
		unitModelData, unitModelErr := salesserviceapiutils.GetUnitModelById(response.ModelId)
		if unitModelErr != nil {
			return pages, unitModelErr
		}

		// Fetch variant data
		unitVariantData, unitVariantErr := salesserviceapiutils.GetUnitVariantById(response.VariantId)
		if unitVariantErr != nil {
			return pages, unitVariantErr
		}

		result := map[string]interface{}{
			"package_id":          response.PackageId,
			"package_code":        response.PackageCode,
			"package_name":        response.PackageName,
			"profit_center_id":    response.ProfitCenterId,
			"profit_center_name":  profitCenterData.ProfitCenterName,
			"model_description":   unitModelData.ModelName,
			"model_code":          unitModelData.ModelCode,
			"variant_code":        unitVariantData.VariantCode,
			"variant_description": unitVariantData.VariantDescription,
			"package_price":       response.PackagePrice,
			"is_active":           response.IsActive,
		}

		results = append(results, result)
	}

	pages.Rows = results
	return pages, nil
}

func (r *PackageMasterRepositoryImpl) GetAllPackageMasterDetail(tx *gorm.DB, id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responsedetail []masterpayloads.PackageMasterDetail
	var item masteritempayloads.BomItemNameResponse
	var operation masterpayloads.Operation
	combinedPayloads := make([]map[string]interface{}, 0)

	// Query Package Master Detail filtered by PackageId
	err := tx.Model(&masterentities.PackageMasterDetail{}).
		Where("package_id = ?", id).
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

		// Fetch Item or Operation data based on LineTypeId
		if detail.LineTypeId != lineTypeOpr.LineTypeId && detail.LineTypeId != 0 {
			err = tx.Table("mtr_item").
				Select("item_name, item_code").
				Joins("JOIN mtr_item_operation ON mtr_item.item_id = mtr_item_operation.item_operation_model_mapping_id").
				Joins("JOIN mtr_package_master_detail ON mtr_package_master_detail.item_operation_id = mtr_item_operation.item_operation_id").
				Where("mtr_package_master_detail.package_detail_id = ?", detail.PackageDetailId).
				Scan(&item).Error
		} else {
			err = tx.Table("mtr_operation_code").
				Select("operation_name, operation_code").
				Joins("JOIN mtr_operation_model_mapping ON mtr_operation_code.operation_id = mtr_operation_model_mapping.operation_id").
				Joins("JOIN mtr_item_operation ON mtr_operation_model_mapping.operation_model_mapping_id = mtr_item_operation.item_operation_model_mapping_id").
				Joins("JOIN mtr_package_master_detail ON mtr_item_operation.item_operation_id = mtr_package_master_detail.item_operation_id").
				Where("mtr_package_master_detail.package_detail_id = ?", detail.PackageDetailId).
				Scan(&operation).Error
		}

		if err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		// Build the response map
		response := map[string]interface{}{
			"is_active":                     detail.IsActive,
			"package_detail_id":             detail.PackageDetailId,
			"package_id":                    detail.PackageId,
			"line_type_id":                  detail.LineTypeId,
			"item_operation_id":             detail.ItemOperationId,
			"frt_quantity":                  detail.FrtQuantity,
			"workorder_transaction_type_id": detail.WorkorderTransactionTypeId,
			"job_type_id":                   detail.JobTypeId,
		}

		// Add item or operation data based on line type
		if detail.LineTypeId != lineTypeOpr.LineTypeId && detail.LineTypeId != 0 {
			response["item_name"] = item.ItemName
			response["item_code"] = item.ItemCode
		} else {
			response["operation_name"] = operation.OperationName
			response["operation_code"] = operation.OperationCode
		}

		combinedPayloads = append(combinedPayloads, response)
	}

	// Paginate the results
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(combinedPayloads, &pages)
	pages.Rows = dataPaginate
	pages.TotalPages = totalPages
	pages.TotalRows = int64(totalRows)

	return pages, nil
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
	var item masteritempayloads.BomItemNameResponse
	var operation masterpayloads.Operation
	err := tx.Model(&entity).Where("package_detail_id=?", id).Scan(&detailpayloads).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	if detailpayloads.LineTypeId != 9 && detailpayloads.LineTypeId != 0 {
		err = tx.Select("mtr_item.item_name,mtr_item.item_code").Table("mtr_package_master_detail").
			Joins("join mtr_item_operation on mtr_item_operation.item_operation_id=mtr_package_master_detail.item_operation_id").
			Joins("join mtr_item on mtr_item.item_id=mtr_item_operation.item_operation_model_mapping_id").
			Where("mtr_package_master_detail.package_detail_id=?", id).
			Scan(&item).
			Error
	} else {
		err = tx.Select("operation_code.operation_name,operation_code.operation_code").Where("package_detail_id=?", id).
			Joins("join mtr_item_operation on mtr_item_operation.item_operation_id = mtr_package_master_detail.item_operation_id").
			Joins("JOIN mtr_operation_model_mapping ON mtr_operation_model_mapping.operation_model_mapping_id=mtr_item_operation.item_operation_model_mapping_id").
			Joins("join mtr_operation_code on mtr_operation_code.operation_id=mtr_operation_model_mapping.operation_id").
			Select("mtr_operation_code.operation_code,mtr_operation_code.operation_name").
			Table("mtr_package_master_detail").
			Scan(&operation).Error
	}

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
			Err:        errProfitcenterUrl,
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
		"package_detail_id":             joinedData1[0]["PackageDetailId"],
		"package_id":                    joinedData1[0]["PackageId"],
		"line_type_id":                  joinedData1[0]["LineTypeId"],
		"item_operation_id":             joinedData1[0]["ItemOperationId"],
		"frt_quantity":                  joinedData1[0]["FrtQuantity"],
		"workorder_transaction_type_id": joinedData1[0]["WorkorderTransactionTypeId"],
		"job_type_id":                   joinedData1[0]["JobTypeId"],
	}

	if detailpayloads.LineTypeId != 9 && detailpayloads.LineTypeId != 1 {
		response["item_name"] = item.ItemName
		response["item_code"] = item.ItemCode
	} else {
		response["operation_name"] = operation.OperationName
		response["operation_code"] = operation.OperationCode
	}

	return response, nil
}

func (r *PackageMasterRepositoryImpl) GetByCodePackageMaster(tx *gorm.DB, code string) (masterentities.PackageMaster, *exceptions.BaseErrorResponse) {
	entities := masterentities.PackageMaster{}

	err := tx.Model(&entities).Where(masterentities.PackageMaster{PackageCode: code}).First(&entities).Error

	if err != nil {
		return masterentities.PackageMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *PackageMasterRepositoryImpl) PostpackageMaster(tx *gorm.DB, req masterpayloads.PackageMasterResponse) (masterentities.PackageMaster, *exceptions.BaseErrorResponse) {
	entities := masterentities.PackageMaster{
		IsActive:       req.IsActive,
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

func (r *PackageMasterRepositoryImpl) PostPackageMasterDetail(tx *gorm.DB, req masterpayloads.PackageMasterDetail, id int) (masterentities.PackageMasterDetail, *exceptions.BaseErrorResponse) {
	entities := masterentities.PackageMasterDetail{
		IsActive:                   req.IsActive,
		PackageId:                  id,
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
