package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"math"
	"time"

	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type BomRepositoryImpl struct {
}

func StartBomRepositoryImpl() masteritemrepository.BomRepository {
	return &BomRepositoryImpl{}
}

func (r *BomRepositoryImpl) GetBomList(tx *gorm.DB, filterConditions []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.BomListResponse

	baseQuery := tx.Table("mtr_bom AS bom").
		Select(`
			bom_id,
			item.item_code,
			item.item_name,
			effective_date,
			qty,
			uom.uom_code,
			bom.is_active`).
		Joins("LEFT JOIN mtr_item AS item ON bom.item_id = item.item_id").
		Joins("LEFT JOIN mtr_uom AS uom ON item.unit_of_measurement_type_id = uom.uom_id")
	/*
		SELECT
			bom_id
			,item.item_code
			,item.item_name
			,effective_date
			,qty
			,uom.uom_code
			,bom.is_active
		FROM [dms_microservices_aftersales_dev].[dbo].[mtr_bom] as bom
		LEFT JOIN mtr_item AS item ON bom.item_id = item.item_id
		LEFT JOIN mtr_uom AS uom ON item.unit_of_measurement_type_id = uom.uom_id
	*/

	baseQuery = utils.ApplyFilter(baseQuery, filterConditions)

	paginatedQuery := baseQuery.Scopes(pagination.Paginate(&pages, baseQuery))

	if err := paginatedQuery.Scan(&responses).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// If empty, return empty list
	if len(responses) == 0 {
		pages.Rows = []masteritementities.MarkupMaster{}
		pages.TotalRows = 0
		pages.TotalPages = 0
		return pages, nil
	}

	pages.Rows = responses
	return pages, nil
}

func (*BomRepositoryImpl) GetBomById(tx *gorm.DB, id int) (masteritempayloads.BomResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Bom{}
	response := masteritempayloads.BomResponse{}

	// Fetch the BOM Master record
	err := tx.Model(&entities).
		Where(masteritementities.Bom{
			BomId: id,
		}).
		First(&response).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return masteritempayloads.BomResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return masteritempayloads.BomResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch BOM Master record",
			Err:        err,
		}
	}

	// If id 0, do error
	if id == 0 {
		return masteritempayloads.BomResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
			Err:        err,
		}
	}

	return response, nil
}

func (r *BomRepositoryImpl) GetBomDetailByMasterId(tx *gorm.DB, bomId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.BomDetailListResponseNew

	baseQuery := tx.Table("mtr_bom_detail AS bom").
		Select(`
			bom_id,
			bom_detail_id,
			bom.is_active,
			seq,
			class.item_class_name,
			item.item_code,
			item.item_name,
			qty,
			uom.uom_code,
			costing_percentage,
			bom.remark`).
		Joins("LEFT JOIN mtr_item AS item ON bom.item_id = item.item_id").
		Joins("LEFT JOIN mtr_item_class AS class on item.item_class_id = class.item_class_id").
		Joins("LEFT JOIN mtr_uom AS uom ON item.unit_of_measurement_type_id = uom.uom_id").
		Where("bom_id = ?", bomId)
	/*
		SELECT
			bom_id,
			bom_detail_id,
			bom.is_active,
			seq,
			itype.item_type_name,
			item.item_code,
			item.item_name,
			qty,
			uom.uom_code,
			costing_percentage,
			bom.remark
		FROM [dms_microservices_aftersales_dev].[dbo].[mtr_bom_detail] as bom
		LEFT JOIN mtr_item AS item ON bom.item_id = item.item_id
		LEFT JOIN mtr_item_type AS itype on item.item_type_id = itype.item_type_id
		LEFT JOIN mtr_uom AS uom ON item.unit_of_measurement_type_id = uom.uom_id
		WHERE bom_id = [int]
	*/

	paginatedQuery := baseQuery.Scopes(pagination.Paginate(&pages, baseQuery))

	if err := paginatedQuery.Scan(&responses).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// If id 0, return empty list
	if len(responses) == 0 {
		pages.Rows = []masteritementities.MarkupMaster{}
		pages.TotalRows = 0
		pages.TotalPages = 0
		return pages, nil
	}

	pages.Rows = responses
	return pages, nil
}

func (r *BomRepositoryImpl) GetBomDetailByMasterUn(tx *gorm.DB, itemId int, effectiveDate time.Time, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.BomDetailListResponseNew

	baseQuery := tx.Table("mtr_bom_detail AS bom").
		Select(`
			bom.bom_id,
			bom.bom_detail_id,
			bom.is_active,
			seq,
			class.item_class_name,
			item.item_code,
			item.item_name,
			bom.qty,
			uom.uom_code,
			costing_percentage,
			bom.remark`).
		Joins("INNER JOIN mtr_bom AS bom_master ON bom_master.bom_id = bom.bom_id").
		Joins("LEFT JOIN mtr_item AS item ON bom.item_id = item.item_id").
		Joins("LEFT JOIN mtr_item_class AS class on item.item_class_id = class.item_class_id").
		Joins("LEFT JOIN mtr_uom AS uom ON item.unit_of_measurement_type_id = uom.uom_id").
		Where("bom_master.item_id = ?", itemId).
		Where("bom_master.effective_date = ?", effectiveDate)
	/*
		SELECT
			bom.bom_id,
			bom.bom_detail_id,
			bom.is_active,
			seq,
			class.item_class_name,
			item.item_code,
			item.item_name,
			bom.qty,
			uom.uom_code,
			costing_percentage,
			bom.remark
		FROM [dms_microservices_aftersales_dev].[dbo].[mtr_bom_detail] as bom
		INNER JOIN mtr_bom AS bom_master ON bom_master.bom_id = bom.bom_id
		LEFT JOIN mtr_item AS item ON bom.item_id = item.item_id
		LEFT JOIN mtr_item_class AS class on item.item_class_id = class.item_class_id
		LEFT JOIN mtr_uom AS uom ON item.unit_of_measurement_type_id = uom.uom_id
		WHERE item_id = [int] AND effective_date = [datetime]
	*/

	paginatedQuery := baseQuery.Scopes(pagination.Paginate(&pages, baseQuery))

	if err := paginatedQuery.Scan(&responses).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// If id 0, return empty list
	if len(responses) == 0 {
		pages.Rows = []masteritementities.MarkupMaster{}
		pages.TotalRows = 0
		pages.TotalPages = 0
		return pages, nil
	}

	pages.Rows = responses
	return pages, nil
}

func (r *BomRepositoryImpl) GetBomDetailById(tx *gorm.DB, id int) (masteritementities.BomDetail, *exceptions.BaseErrorResponse) {
	entities := masteritementities.BomDetail{}

	// Fetch the BOM Master record
	err := tx.Model(&entities).
		Where(masteritementities.BomDetail{
			BomDetailId: id,
		}).
		First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return masteritementities.BomDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return masteritementities.BomDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch BOM detail record",
			Err:        err,
		}
	}

	// If empty, do error
	if id == 0 {
		return masteritementities.BomDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
			Err:        err,
		}
	}

	return entities, nil
}

func (*BomRepositoryImpl) SaveBomMaster(tx *gorm.DB, request masteritempayloads.BomMasterRequest) (masteritementities.Bom, *exceptions.BaseErrorResponse) {
	var entities masteritementities.Bom

	entities.IsActive = request.IsActive
	entities.Qty = request.BomMasterQty
	entities.EffectiveDate = request.BomMasterEffectiveDate
	entities.ItemId = request.ItemId

	err := tx.Create(&entities).Error
	if err != nil {
		return masteritementities.Bom{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}

func (*BomRepositoryImpl) UpdateBomMaster(tx *gorm.DB, id int, request masteritempayloads.BomMasterRequest) (masteritementities.Bom, *exceptions.BaseErrorResponse) {
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

		return masteritementities.Bom{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	entities.IsActive = request.IsActive
	entities.Qty = request.BomMasterQty
	entities.EffectiveDate = request.BomMasterEffectiveDate
	entities.ItemId = request.ItemId

	result = tx.Save(&entities)

	if result.Error != nil {
		return masteritementities.Bom{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return entities, nil
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

func (r *BomRepositoryImpl) GetBomDetailList(tx *gorm.DB, filters []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	type BomDetailResponse struct {
		BomMasterId             int     `json:"bom_master_id"`
		BomDetailSeq            int     `json:"bom_detail_seq"`
		ItemCode                string  `json:"item_code"`
		ItemName                string  `json:"item_name"`
		ItemClassCode           string  `json:"item_class_code"`
		LineTypeName            string  `json:"line_type_name"`
		BomDetailCostingPercent float64 `json:"bom_detail_costing_percent"`
		BomDetailRemark         string  `json:"bom_detail_remark"`
		BomDetailQty            float64 `json:"bom_detail_qty"`
		BomDetailId             int     `json:"bom_detail_id"`
		UomDescription          string  `json:"uom_description"`
	}

	var responses []BomDetailResponse

	baseQuery := tx.Table("mtr_bom_detail AS det").
		Select(`
			det.bom_master_id,
			det.bom_detail_seq,
			item.item_code,
			item.item_name,
			iclas.item_class_code,
			lt.line_type_name,
			det.bom_detail_costing_percent,
			det.bom_detail_remark,
			det.bom_detail_qty,
			det.bom_detail_id,
			uom.uom_description`).
		Joins("LEFT JOIN mtr_item AS item ON det.bom_detail_material_id = item.item_id").
		Joins("LEFT JOIN mtr_uom AS uom ON item.unit_of_measurement_type_id = uom.uom_id").
		Joins("LEFT JOIN mtr_item_class AS iclas ON item.item_class_id = iclas.item_class_id").
		Joins("LEFT JOIN dms_microservices_general_dev.dbo.mtr_line_type AS lt ON iclas.line_type_id = lt.line_type_id")

	baseQuery = utils.ApplyFilter(baseQuery, filters)

	var totalRows int64
	err := baseQuery.Count(&totalRows).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(pages.Limit)))

	paginatedQuery := baseQuery.Scopes(pagination.Paginate(&pages, baseQuery))

	if err := paginatedQuery.Scan(&responses).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []map[string]interface{}{}
		pages.TotalRows = totalRows
		pages.TotalPages = totalPages
		return pages, nil
	}

	pages.Rows = responses
	pages.TotalRows = totalRows
	pages.TotalPages = totalPages

	return pages, nil
}

func (r *BomRepositoryImpl) SaveBomDetail(tx *gorm.DB, request masteritempayloads.BomDetailRequest) (masteritementities.BomDetail, *exceptions.BaseErrorResponse) {
	// Tentukan nilai BomDetailSeq
	var newBomDetailSeq int
	if err := tx.Model(&masteritementities.BomDetail{}).
		Where("bom_master_id = ?", request.BomMasterId).
		Select("COALESCE(MAX(bom_detail_seq), 0)").
		Row().
		Scan(&newBomDetailSeq); err != nil {
		return masteritementities.BomDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	newBomDetailSeq++ // Tambahkan 1 pada nilai maksimum untuk mendapatkan nilai BomDetailSeq yang baru

	// Buat entitas BomDetail
	newBomDetail := masteritementities.BomDetail{
		BomId:             request.BomMasterId,
		Seq:               newBomDetailSeq,
		Qty:               request.BomDetailQty,
		CostingPercentage: request.BomDetailCostingPercent,
		Remark:            request.BomDetailRemark,
		ItemId:            request.BomDetailTypeId, // TypeId -> ItemId
	}

	// Simpan entitas BomDetail
	err := tx.Create(&newBomDetail).Error
	if err != nil {
		return masteritementities.BomDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return newBomDetail, nil
}

func (r *BomRepositoryImpl) UpdateBomDetail(tx *gorm.DB, id int, request masteritempayloads.BomDetailRequest) (masteritementities.BomDetail, *exceptions.BaseErrorResponse) {
	var entities masteritementities.BomDetail

	result := tx.Model(&entities).
		Where("bom_detail_id = ?", id).
		First(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return masteritementities.BomDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        fmt.Errorf("bom detail with ID %d not found", id),
			}
		}
		return masteritementities.BomDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	entities.Qty = request.BomDetailQty
	entities.CostingPercentage = request.BomDetailCostingPercent
	entities.Remark = request.BomDetailRemark

	result = tx.Save(&entities)

	if result.Error != nil {
		return masteritementities.BomDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return entities, nil
}

func (r *BomRepositoryImpl) GetBomItemList(tx *gorm.DB, filterConditions []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	type BomItemResponse struct {
		IsActive       bool   `json:"is_active"`
		ItemId         int    `json:"item_id"`
		ItemCode       string `json:"item_code"`
		ItemName       string `json:"item_name"`
		ItemTypeId     string `json:"item_type_id"`
		ItemGroupId    int    `json:"item_group_id"`
		ItemClassId    int    `json:"item_class_id"`
		ItemClassCode  string `json:"item_class_code"`
		UomId          int    `json:"uom_id"`
		UomDescription string `json:"uom_description"`
	}

	var responses []BomItemResponse

	baseQuery := tx.Table("mtr_item AS item").
		Select(`
			item.is_active,
			item.item_id,
			item.item_code,
			item.item_name,
			item.item_type_id,
			item.item_group_id,
			item_class.item_class_id,
			item_class.item_class_code,
			uom.uom_id,
			uom.uom_description`).
		Joins("LEFT JOIN mtr_item_class AS item_class ON item.item_class_id = item_class.item_class_id").
		Joins("LEFT JOIN mtr_uom AS uom ON item.unit_of_measurement_selling_id = uom.uom_id")

	// Apply filters
	baseQuery = utils.ApplyFilter(baseQuery, filterConditions)

	// Apply pagination
	paginatedQuery := baseQuery.Scopes(pagination.Paginate(&pages, baseQuery))

	// Scan the results into the response struct
	if err := paginatedQuery.Scan(&responses).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Handle empty results
	if len(responses) == 0 {
		pages.Rows = []BomItemResponse{}
		pages.TotalRows = 0
		pages.TotalPages = 0
		return pages, nil
	}

	// Assign results to the paginated response
	pages.Rows = responses

	return pages, nil
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
