package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type BomRepositoryImpl struct {
}

func StartBomRepositoryImpl() masteritemrepository.BomRepository {
	return &BomRepositoryImpl{}
}

func (r *BomRepositoryImpl) GetBomMasterList(tx *gorm.DB, filters []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var responses []map[string]interface{}

	// Define main table
	mainTable := "mtr_bom"
	mainAlias := "bom"

	// Define join tables
	joinTables := []utils.JoinTable{
		{Table: "mtr_item", Alias: "item", ForeignKey: "bom.item_id", ReferenceKey: "item.item_id"},
		{Table: "mtr_uom", Alias: "uom", ForeignKey: "item.unit_of_measurement_type_id", ReferenceKey: "uom.uom_id"},
	}

	// Create join query
	joinQuery := utils.CreateJoin(tx, mainTable, mainAlias, joinTables...)

	// Define key attributes to be selected
	keyAttributes := []string{
		"bom.is_active",
		"bom.bom_master_id",
		"bom.bom_master_qty",
		"bom.bom_master_effective_date",
		"item.item_code",
		"item.item_name",
		"item.item_id",
		"uom.uom_id",
		"uom.uom_description",
	}

	// Apply key attributes selection
	joinQuery = joinQuery.Select(keyAttributes)

	// Apply filters
	whereQuery := utils.ApplyFilter(joinQuery, filters)

	// Execute query
	rows, err := whereQuery.Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	// Fetch data and append to response slice
	for rows.Next() {
		var isActive bool
		var bomMasterId, bomMasterQty int
		var bomMasterEffectiveDate time.Time
		var itemId, uomId int
		var itemCode, itemName, uomDescription string

		err := rows.Scan(&isActive,
			&bomMasterId,
			&bomMasterQty,
			&bomMasterEffectiveDate,
			&itemCode,
			&itemName,
			&itemId,
			&uomId,
			&uomDescription)

		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		responseMap := map[string]interface{}{
			"is_active":                 isActive,
			"bom_master_id":             bomMasterId,
			"bom_master_qty":            bomMasterQty,
			"bom_master_effective_date": bomMasterEffectiveDate,
			"item_code":                 itemCode,
			"item_name":                 itemName,
			"uom_description":           uomDescription,
		}
		responses = append(responses, responseMap)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(responses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (*BomRepositoryImpl) GetBomMasterById(tx *gorm.DB, id int) (masteritempayloads.BomMasterRequest, *exceptions.BaseErrorResponse) {
	var response masteritempayloads.BomMasterRequest

	err := tx.Table("mtr_bom").
		Select("mtr_bom.bom_master_id, mtr_bom.is_active, mtr_bom.bom_master_qty,  mtr_bom.bom_master_effective_date, mtr_bom.bom_master_change_number, mtr_item.item_code, mtr_item.item_name, mtr_item.item_id,mtr_uom.uom_description").
		Joins("JOIN mtr_item ON mtr_bom.item_id = mtr_item.item_id").
		Joins("JOIN mtr_uom ON mtr_uom.uom_id = mtr_item.unit_of_measurement_type_id").
		Where("mtr_bom.bom_master_id = ?", id).
		First(&response).
		Error

	if err != nil {
		// notFoundErr := exceptions.NewNotFoundError("Bom master not found")
		return masteritempayloads.BomMasterRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	return response, nil
}

func (r *BomRepositoryImpl) SaveBomMaster(tx *gorm.DB, request masteritempayloads.BomMasterRequest) (bool, *exceptions.BaseErrorResponse) {

	entities := masteritementities.Bom{
		BomMasterId:            request.BomMasterId,
		BomMasterQty:           request.BomMasterQty,
		BomMasterEffectiveDate: request.BomMasterEffectiveDate,
		BomMasterChangeNumber:  request.BomMasterChangeNumber,
		ItemId:                 request.ItemId,
	}

	if request.BomMasterId == 0 {
		err := tx.Create(&entities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			} // Mengembalikan pesan kesalahan jika terjadi error saat membuat data baru
		}
	} else {
		err := tx.Model(&masteritementities.Bom{}).
			Where("bom_master_id = ?", request.BomMasterId).
			Updates(entities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			} // Mengembalikan pesan kesalahan jika terjadi error saat memperbarui data yang sudah ada
		}
	}

	return true, nil // Mengembalikan true jika operasi berhasil tanpa error
}

func (r *BomRepositoryImpl) ChangeStatusBomMaster(tx *gorm.DB, id int) (masteritementities.Bom, *exceptions.BaseErrorResponse) {
	var entities masteritementities.Bom

	result := tx.Model(&entities).
		Where("bom_master_id = ?", id).
		First(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return masteritementities.Bom{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        fmt.Errorf("bom with ID %d not found", id),
			}
		}
		// Jika ada galat lain, kembalikan galat internal server
		return masteritementities.Bom{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	// Ubah status entitas
	entities.IsActive = !entities.IsActive

	result = tx.Save(&entities)

	if result.Error != nil {
		return masteritementities.Bom{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        result.Error,
		}
	}

	return entities, nil
}

func (r *BomRepositoryImpl) GetBomDetailList(tx *gorm.DB, filters []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.BomDetailListResponse

	// Define join table
	joinTable := tx.Table("mtr_bom as bom").
		Select("bom.bom_master_id, bom.is_active, bom.bom_master_effective_date, bom.bom_master_qty, det.bom_detail_seq, item.item_code, item.item_name, iclas.item_class_code, lt.line_type_name, det.bom_detail_costing_percent, det.bom_detail_remark, det.bom_detail_qty , det.bom_detail_id,uom.uom_description").
		Joins("left join mtr_bom_detail as det ON bom.bom_master_id = det.bom_master_id").
		Joins("INNER join mtr_item as item ON bom.item_id = item.item_id").
		Joins("INNER join mtr_uom as uom ON item.unit_of_measurement_type_id  = uom.uom_id").
		Joins("INNER join mtr_item_class as iclas ON item.item_class_id = iclas.item_class_id").
		Joins("INNER join dms_microservices_general_dev.dbo.mtr_line_type as lt ON iclas.line_type_id = lt.line_type_id")

	// Apply filters
	whereQuery := utils.ApplyFilter(joinTable, filters)

	// Execute query
	rows, err := whereQuery.Find(&responses).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	// Convert responses to maps
	responseMaps := make([]map[string]interface{}, 0)
	for _, response := range responses {
		responseMap := map[string]interface{}{
			"bom_master_id":              response.BomMasterId,
			"is_active":                  response.IsActive,
			"bom_master_effective_date":  response.BomMasterEffectiveDate,
			"bom_master_qty":             response.BomMasterQty,
			"bom_detail_seq":             response.BomDetailSeq,
			"item_code":                  response.ItemCode,
			"item_name":                  response.ItemName,
			"item_class_code":            response.ItemClassCode,
			"line_type_name":             response.LineTypeName,
			"bom_detail_costing_percent": response.BomDetailCostingPercent,
			"bom_detail_remark":          response.BomDetailRemark,
			"bom_detail_qty":             response.BomDetailQty,
			"bom_detail_id":              response.BomDetailId,
			"uom_description":            response.UomDescription,
		}
		responseMaps = append(responseMaps, responseMap)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(responseMaps, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *BomRepositoryImpl) GetBomDetailById(tx *gorm.DB, id int) ([]masteritempayloads.BomDetailListResponse, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.BomDetailListResponse

	// Execute query
	err := tx.Table("mtr_bom as bom").
		Select("bom.bom_master_id, bom.is_active, bom.bom_master_effective_date, bom.bom_master_qty, det.bom_detail_seq, item.item_code, item.item_name, iclas.item_class_code, lt.line_type_name, det.bom_detail_costing_percent, det.bom_detail_remark , det.bom_detail_qty , det.bom_detail_id,uom.uom_description").
		Joins("left join mtr_bom_detail as det ON bom.bom_master_id = det.bom_master_id").
		Joins("INNER join mtr_item as item ON bom.item_id = item.item_id").
		Joins("INNER join mtr_uom as uom ON item.unit_of_measurement_type_id  = uom.uom_id").
		Joins("INNER join mtr_item_class as iclas ON item.item_class_id = iclas.item_class_id").
		Joins("INNER join dms_microservices_general_dev.dbo.mtr_line_type as lt ON iclas.line_type_id = lt.line_type_id").
		Where("bom.bom_master_id = ?", id).
		Find(&responses).Error
	if err != nil {
		// notFoundErr := exceptions.NewNotFoundError("Bom master not found")
		return []masteritempayloads.BomDetailListResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	// Mengembalikan response
	return responses, nil
}

func (r *BomRepositoryImpl) GetBomDetailByIds(tx *gorm.DB, id int) ([]masteritempayloads.BomDetailListResponse, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.BomDetailListResponse

	// Execute query
	err := tx.Table("mtr_bom as bom").
		Select("bom.bom_master_id, bom.is_active, bom.bom_master_effective_date, bom.bom_master_qty, det.bom_detail_seq, item.item_code, item.item_name, iclas.item_class_code, lt.line_type_name, det.bom_detail_costing_percent, det.bom_detail_remark , det.bom_detail_qty , det.bom_detail_id,uom.uom_description").
		Joins("left join mtr_bom_detail as det ON bom.bom_master_id = det.bom_master_id").
		Joins("INNER join mtr_item as item ON bom.item_id = item.item_id").
		Joins("INNER join mtr_uom as uom ON item.unit_of_measurement_type_id  = uom.uom_id").
		Joins("INNER join mtr_item_class as iclas ON item.item_class_id = iclas.item_class_id").
		Joins("INNER join dms_microservices_general_dev.dbo.mtr_line_type as lt ON iclas.line_type_id = lt.line_type_id").
		Where("det.bom_detail_id = ?", id).
		Find(&responses).Error
	if err != nil {
		// notFoundErr := exceptions.NewNotFoundError("Bom detail not found")
		return []masteritempayloads.BomDetailListResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	// Mengembalikan response
	return responses, nil
}

func (r *BomRepositoryImpl) SaveBomDetail(tx *gorm.DB, request masteritempayloads.BomDetailRequest) (bool, *exceptions.BaseErrorResponse) {
	// Tentukan nilai BomDetailSeq
	var newBomDetailSeq int
	if err := tx.Model(&masteritementities.BomDetail{}).
		Where("bom_master_id = ?", request.BomMasterId).
		Select("COALESCE(MAX(bom_detail_seq), 0)").
		Row().
		Scan(&newBomDetailSeq); err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		} // Mengembalikan pesan kesalahan jika terjadi error saat mengambil nilai maksimum bom_detail_seq
	}
	newBomDetailSeq++ // Tambahkan 1 pada nilai maksimum untuk mendapatkan nilai BomDetailSeq yang baru

	entities := masteritementities.BomDetail{
		BomDetailId:             request.BomDetailId,
		BomMasterId:             request.BomMasterId,
		BomDetailSeq:            newBomDetailSeq, // Gunakan nilai BomDetailSeq yang baru
		BomDetailQty:            request.BomDetailQty,
		BomDetailUom:            request.BomDetailUom,
		BomDetailRemark:         request.BomDetailRemark,
		BomDetailCostingPercent: request.BomDetailCostingPercent,
		BomDetailType:           request.BomDetailType,
		BomDetailMaterialCode:   request.BomDetailMaterialCode,
		BomDetailMaterialName:   request.BomDetailMaterialName,
	}

	if request.BomDetailId == 0 {
		err := tx.Create(&entities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			} // Mengembalikan pesan kesalahan jika terjadi error saat membuat data baru
		}
	} else {
		err := tx.Model(&masteritementities.BomDetail{}).
			Where("bom_detail_id = ?", request.BomDetailId).
			Updates(entities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			} // Mengembalikan pesan kesalahan jika terjadi error saat memperbarui data yang sudah ada
		}
	}

	return true, nil // Mengembalikan true jika operasi berhasil tanpa error
}

func (r *BomRepositoryImpl) GetBomItemList(tx *gorm.DB, filters []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.BomItemLookup

	// Define table struct
	tableStruct := masteritempayloads.BomItemLookup{}
	// Define join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	// Apply filters
	whereQuery := utils.ApplyFilter(joinTable, filters)

	// Execute query
	rows, err := whereQuery.Find(&responses).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	// Convert responses to maps
	responseMaps := make([]map[string]interface{}, 0)
	for _, response := range responses {
		responseMap := map[string]interface{}{
			"item_code":       response.ItemCode,
			"item_name":       response.ItemName,
			"item_type":       response.ItemType,
			"item_group_code": response.ItemGroupId,
			"item_class_code": response.ItemClassCode,
			"is_active":       response.IsActive,
		}
		responseMaps = append(responseMaps, responseMap)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(responseMaps, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *BomRepositoryImpl) DeleteByIds(tx *gorm.DB, ids []int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.BomDetail

	if err := tx.Delete(&entities, ids).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}
