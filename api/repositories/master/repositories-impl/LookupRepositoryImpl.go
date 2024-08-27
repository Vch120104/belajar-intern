package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	masteroperationentities "after-sales/api/entities/master/operation"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"fmt"
	"math"
	"net/http"

	"gorm.io/gorm"
)

type LookupRepositoryImpl struct {
}

func StartLookupRepositoryImpl() masterrepository.LookupRepository {
	return &LookupRepositoryImpl{}
}

func (r *LookupRepositoryImpl) ItemOprCode(tx *gorm.DB, linetypeId int, paginate pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var results []map[string]interface{}
	var err error

	// Default filters and variables
	ItmGrpInventory := "IN"
	PurchaseTypeGoods := "G"
	ItmCls := ""
	year, month, companyCode := 2024, 8, "COMP_CODE" // Placeholder for real values or dynamic fetching

	switch linetypeId {
	case utils.LinetypePackage:
		// GORM query for Line Type Package
		err = tx.Table("mtr_package A").
			Select("A.PACKAGE_CODE AS PackageCode, A.PACKAGE_NAME AS PackageName, SUM(A1.FRT_QTY) AS FRT, B.CPC_DESCRIPTION AS ProfitCenter, A.MODEL_CODE AS ModelCode, C.MODEL_DESCRIPTION AS Description, A.PACKAGE_PRICE AS Price").
			Joins("LEFT JOIN mtr_package_detail A1 ON A.PACKAGE_CODE = A1.PACKAGE_CODE").
			Joins("LEFT JOIN mtr_profit_center B ON A.CPC_CODE = B.CPC_CODE").
			Joins("LEFT JOIN mtr_unit_model C ON A.MODEL_CODE = C.MODEL_CODE").
			Where("A.RECORD_STATUS = ? AND A1.RECORD_STATUS = ?", "A", "A").
			Group("A.PACKAGE_CODE, A.PACKAGE_NAME, B.CPC_DESCRIPTION, A.MODEL_CODE, C.MODEL_DESCRIPTION, A.PACKAGE_PRICE").
			Offset((paginate.Page - 1) * paginate.Limit).Limit(paginate.Limit).Find(&results).Error

	case utils.LinetypeOperation:
		// GORM query for Line Type Operation
		err = tx.Table("amOperation2 A").
			Select("A.OPERATION_CODE AS Code, B.OPERATION_NAME AS Description, A.FRT_HOUR AS FRT, C.OPR_ENTRIES_CODE AS OprEntriesCode, G.OPR_ENTRIES_DESC AS OprEntriesName, C.OPR_KEY_CODE AS OprKeyCode, F.OPR_KEY_DESC AS OprKeyName").
			Joins("INNER JOIN amOperation0 O ON O.OPERATION_CODE = A.OPERATION_CODE AND O.VEHICLE_BRAND = A.VEHICLE_BRAND AND O.MODEL_CODE = A.MODEL_CODE").
			Joins("LEFT JOIN amOperation3 C ON A.OPERATION_CODE = C.OPERATION_CODE AND A.VEHICLE_BRAND = C.VEHICLE_BRAND AND A.MODEL_CODE = C.MODEL_CODE").
			Joins("LEFT JOIN amOperationCode B ON A.OPERATION_CODE = B.OPERATION_CODE").
			Joins("LEFT JOIN amOprKey F ON C.OPR_KEY_CODE = F.OPR_KEY_CODE").
			Joins("LEFT JOIN amOprEntries G ON C.OPR_ENTRIES_CODE = G.OPR_ENTRIES_CODE").
			Where("A.RECORD_STATUS = ? AND O.RECORD_STATUS = ?", "A", "A").
			Offset((paginate.Page - 1) * paginate.Limit).Limit(paginate.Limit).Find(&results).Error

	case utils.LinetypeSparepart:
		// GORM query for Line Type Sparepart
		ItmCls = "SP"
		err = tx.Table("gmItem0 A").
			Select("A.ITEM_CODE AS Code, A.ITEM_NAME AS Description, ISNULL((SELECT SUM(V.QTY_AVAILABLE) FROM viewLocationStock V WHERE A.ITEM_CODE = V.ITEM_CODE AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.COMPANY_CODE = ?), 0) AS AvailQty, A.ITEM_LVL_1 AS ItemLvl1, A.ITEM_LVL_2 AS ItemLvl2, A.ITEM_LVL_3 AS ItemLvl3, A.ITEM_LVL_4 AS ItemLvl4", year, month, companyCode).
			Joins("INNER JOIN gmItem1 B ON A.ITEM_CODE = B.ITEM_CODE").
			Where("A.ITEM_GROUP = ? AND A.ITEM_TYPE = ? AND A.ITEM_CLASS = ? AND A.RECORD_STATUS = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, "A").
			Offset((paginate.Page - 1) * paginate.Limit).Limit(paginate.Limit).Find(&results).Error

	case utils.LinetypeOil:
		// GORM query for Line Type Oil
		ItmCls = "OL"
		err = tx.Table("gmItem0 A").
			Select("A.ITEM_CODE AS Code, A.ITEM_NAME AS Description, ISNULL((SELECT SUM(V.QTY_AVAILABLE) FROM viewLocationStock V WHERE A.ITEM_CODE = V.ITEM_CODE AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.COMPANY_CODE = ?), 0) AS AvailQty, A.ITEM_LVL_1 AS ItemLvl1, A.ITEM_LVL_2 AS ItemLvl2, A.ITEM_LVL_3 AS ItemLvl3, A.ITEM_LVL_4 AS ItemLvl4", year, month, companyCode).
			Joins("INNER JOIN gmItem1 B ON A.ITEM_CODE = B.ITEM_CODE").
			Where("A.ITEM_GROUP = ? AND A.ITEM_TYPE = ? AND A.ITEM_CLASS = ? AND A.RECORD_STATUS = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, "A").
			Offset((paginate.Page - 1) * paginate.Limit).Limit(paginate.Limit).Find(&results).Error

	case utils.LinetypeMaterial:
		// GORM query for Line Type Material
		ItmCls = "MT"
		ItmClsSublet := "Sublet_Class"
		err = tx.Table("gmItem0 A").
			Select("DISTINCT A.ITEM_CODE AS Code, A.ITEM_NAME AS Description, ISNULL((SELECT SUM(V.QTY_AVAILABLE) FROM viewLocationStock V WHERE A.ITEM_CODE = V.ITEM_CODE AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.COMPANY_CODE = ?), 0) AS AvailQty, A.ITEM_LVL_1 AS ItemLvl1, A.ITEM_LVL_2 AS ItemLvl2, A.ITEM_LVL_3 AS ItemLvl3, A.ITEM_LVL_4 AS ItemLvl4", year, month, companyCode).
			Joins("INNER JOIN gmItem1 B ON A.ITEM_CODE = B.ITEM_CODE").
			Where("A.ITEM_GROUP = ? AND A.ITEM_TYPE = ? AND (A.ITEM_CLASS = ? OR A.ITEM_CLASS = ?) AND A.RECORD_STATUS = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, ItmClsSublet, "A").
			Offset((paginate.Page - 1) * paginate.Limit).Limit(paginate.Limit).Find(&results).Error

	case utils.LinetypeSublet:
		// GORM query for Line Type Sublet
		ItmCls = "Fee_Class"
		ItmGrpOutsideJob := "OutsideJob_Group"
		PurchaseTypeServices := "Service_Type"
		err = tx.Table("gmItem0 A").
			Select("DISTINCT A.ITEM_CODE AS Code, A.ITEM_NAME AS Description, ISNULL((SELECT SUM(V.QTY_AVAILABLE) FROM viewLocationStock V WHERE A.ITEM_CODE = V.ITEM_CODE AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.COMPANY_CODE = ?), 0) AS AvailQty, A.ITEM_LVL_1 AS ItemLvl1, A.ITEM_LVL_2 AS ItemLvl2, A.ITEM_LVL_3 AS ItemLvl3, A.ITEM_LVL_4 AS ItemLvl4", year, month, companyCode).
			Joins("INNER JOIN gmItem1 B ON A.ITEM_CODE = B.ITEM_CODE").
			Where("(A.ITEM_GROUP = ? OR (A.ITEM_GROUP = ? AND A.ITEM_CLASS = ?)) AND A.ITEM_TYPE = ? AND A.RECORD_STATUS = ?", ItmGrpOutsideJob, ItmGrpInventory, ItmCls, PurchaseTypeServices, "A").
			Offset((paginate.Page - 1) * paginate.Limit).Limit(paginate.Limit).Find(&results).Error

	default:
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Line Type ID",
			Err:        fmt.Errorf("invalid line type ID: %d", linetypeId),
		}
	}

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Calculate Total Rows and Pages
	var totalRows int64
	var countErr error
	switch linetypeId {
	case utils.LinetypePackage:
		countErr = tx.Model(&masterentities.PackageMaster{}).Where("RECORD_STATUS = ?", "A").Count(&totalRows).Error
	case utils.LinetypeOperation:
		countErr = tx.Model(&masteroperationentities.OperationFrt{}).Where("RECORD_STATUS = ?", "A").Count(&totalRows).Error
	case utils.LinetypeSparepart:
		countErr = tx.Model(&masteritementities.Item{}).Where("ITEM_GROUP = ? AND ITEM_TYPE = ? AND ITEM_CLASS = ? AND RECORD_STATUS = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, "A").Count(&totalRows).Error
	case utils.LinetypeOil:
		countErr = tx.Model(&masteritementities.Item{}).Where("ITEM_GROUP = ? AND ITEM_TYPE = ? AND ITEM_CLASS = ? AND RECORD_STATUS = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, "A").Count(&totalRows).Error
	case utils.LinetypeMaterial:
		countErr = tx.Model(&masteritementities.Item{}).Where("ITEM_GROUP = ? AND ITEM_TYPE = ? AND (ITEM_CLASS = ? OR ITEM_CLASS = ?) AND RECORD_STATUS = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, "Sublet_Class", "A").Count(&totalRows).Error
	case utils.LinetypeSublet:
		ItmGrpOutsideJob := "OutsideJob_Group"
		PurchaseTypeServices := "Service_Type"
		countErr = tx.Model(&masteritementities.Item{}).Where("(ITEM_GROUP = ? OR (ITEM_GROUP = ? AND ITEM_CLASS = ?)) AND ITEM_TYPE = ? AND RECORD_STATUS = ?", ItmGrpOutsideJob, ItmGrpInventory, ItmCls, PurchaseTypeServices, "A").Count(&totalRows).Error
	}

	if countErr != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(paginate.Limit)))

	return results, int(totalRows), totalPages, nil
}
