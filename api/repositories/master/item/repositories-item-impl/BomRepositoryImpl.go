package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
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
	responses := []masteritempayloads.BomListResponse{}

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
		Joins("LEFT JOIN mtr_uom AS uom ON item.unit_of_measurement_stock_id = uom.uom_id")
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
		LEFT JOIN mtr_uom AS uom ON item.unit_of_measurement_stock_id = uom.uom_id
	*/

	baseQuery = utils.ApplyFilter(baseQuery, filterConditions)

	paginatedQuery := baseQuery.Scopes(pagination.Paginate(&pages, baseQuery))

	if err := paginatedQuery.Scan(&responses).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching BOM list",
			Err:        err,
		}
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

func (r *BomRepositoryImpl) ChangeStatusBomMaster(tx *gorm.DB, id int) (masteritementities.Bom, *exceptions.BaseErrorResponse) {
	var entities masteritementities.Bom

	err := tx.Model(&entities).
		Where("bom_id = ?", id).
		First(&entities).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return masteritementities.Bom{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        fmt.Errorf("BOM with ID %d not found", id),
			}
		}
		// Jika ada galat lain, kembalikan galat internal server
		return masteritementities.Bom{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error updating BOM",
			Err:        err,
		}
	}

	// Ubah status entitas
	entities.IsActive = !entities.IsActive

	err = tx.Save(&entities).Error

	if err != nil {
		return masteritementities.Bom{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error updating BOM",
			Err:        err,
		}
	}

	return entities, nil
}

func (r *BomRepositoryImpl) GetBomDetailByMasterId(tx *gorm.DB, bomId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	responses := []masteritempayloads.BomDetailListResponse{}

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
		Joins("LEFT JOIN mtr_uom AS uom ON item.unit_of_measurement_stock_id = uom.uom_id").
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
		LEFT JOIN mtr_uom AS uom ON item.unit_of_measurement_stock_id = uom.uom_id
		WHERE bom_id = [int]
	*/

	paginatedQuery := baseQuery.Scopes(pagination.Paginate(&pages, baseQuery))

	if err := paginatedQuery.Scan(&responses).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching BOM detail record",
			Err:        err,
		}
	}

	pages.Rows = responses
	return pages, nil
}

func (r *BomRepositoryImpl) GetBomDetailByMasterUn(tx *gorm.DB, itemId int, effectiveDate time.Time, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	responses := []masteritempayloads.BomDetailListResponse{}

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
		Joins("LEFT JOIN mtr_uom AS uom ON item.unit_of_measurement_stock_id = uom.uom_id").
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
		LEFT JOIN mtr_uom AS uom ON item.unit_of_measurement_stock_id = uom.uom_id
		WHERE item_id = [int] AND effective_date = [datetime]
	*/

	paginatedQuery := baseQuery.Scopes(pagination.Paginate(&pages, baseQuery))

	if err := paginatedQuery.Scan(&responses).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching BOM detail record",
			Err:        err,
		}
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

func (*BomRepositoryImpl) UpdateBomMaster(tx *gorm.DB, id int, qty float64) (masteritementities.Bom, *exceptions.BaseErrorResponse) {
	// Check if exists
	var entities masteritementities.Bom
	err := tx.Model(&entities).Select("qty").
		Where("bom_id = ?", id).
		First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return masteritementities.Bom{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        fmt.Errorf("BOM with ID %d not found", id),
			}
		}
		// Jika ada galat lain, kembalikan galat internal server
		return masteritementities.Bom{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error updating BOM",
			Err:        err,
		}
	}

	// Update
	err = tx.Model(&entities).Select("qty").
		Where("bom_id = ?", id).
		Updates(masteritementities.Bom{
			Qty: qty,
		}).Error
	if err != nil {
		return masteritementities.Bom{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error updating BOM",
			Err:        err,
		}
	}

	return entities, nil
}

func (*BomRepositoryImpl) SaveBomMaster(tx *gorm.DB, request masteritempayloads.BomMasterNewRequest) (masteritementities.Bom, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Bom{
		IsActive:      true,
		Qty:           request.Qty,
		EffectiveDate: request.EffectiveDate,
		ItemId:        request.ItemId,
	}
	err := tx.Create(&entities).Error
	if err != nil {
		return masteritementities.Bom{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Data already exists",
			Err:        err,
		}
	}

	return entities, nil
}

func (*BomRepositoryImpl) FirstOrCreateBom(tx *gorm.DB, request masteritempayloads.BomMasterNewRequest) (int, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Bom{
		IsActive:      true,
		EffectiveDate: request.EffectiveDate,
		ItemId:        request.ItemId,
	}

	// Check if item id the same
	var check int64
	errA := tx.Model(&masteritementities.Bom{}).
		Where(entities).
		Count(&check).Error
	if errA != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching BOM record",
			Err:        errA,
		}
	}
	if check == 0 { // Create because null
		err := tx.Create(&entities).Error
		if err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error creating BOM",
				Err:        err,
			}
		}

		return entities.BomId, nil
	}

	var bomId int
	err := tx.Model(&masteritementities.Bom{}).
		Where(entities).
		Pluck("bom_id", &bomId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch BOM record",
			Err:        err,
		}
	}
	return bomId, nil
}

func (r *BomRepositoryImpl) GetBomDetailTemplate(tx *gorm.DB, filters []utils.FilterCondition, pages pagination.Pagination) ([]masteritempayloads.BomDetailTemplate, *exceptions.BaseErrorResponse) {
	responses := []masteritempayloads.BomDetailTemplate{}

	baseQuery := tx.Table("mtr_bom_detail AS detail").
		Select(`
			item2.item_code,
			bom.effective_date,
			bom.qty,
			item1.item_code bom_detail_item_code,
			detail.seq bom_detail_seq,
			detail.qty bom_detail_qty,
			detail.remark bom_detail_remark,
			detail.costing_percentage bom_detail_costing_percentage`).
		Joins("INNER JOIN mtr_bom bom on bom.bom_id = detail.bom_id").
		Joins("INNER JOIN mtr_item item1 on item1.item_id = detail.item_id").
		Joins("INNER JOIN mtr_item item2 on item2.item_id = bom.item_id")

	baseQuery = utils.ApplyFilter(baseQuery, filters)

	var totalRows int64
	err := baseQuery.Count(&totalRows).Error
	if err != nil {
		return []masteritempayloads.BomDetailTemplate{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching BOM detail record",
			Err:        err,
		}
	}

	paginatedQuery := baseQuery.Scopes(pagination.Paginate(&pages, baseQuery))

	if err := paginatedQuery.Scan(&responses).Error; err != nil {
		return []masteritempayloads.BomDetailTemplate{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching BOM detail record",
			Err:        err,
		}
	}

	return responses, nil
}

func (*BomRepositoryImpl) GetBomDetailMaxSeq(tx *gorm.DB, id int) (int, *exceptions.BaseErrorResponse) {
	var maxSeq int

	err := tx.Model(&masteritementities.BomDetail{}).
		Where("bom_id = ?", id).
		Select("COALESCE(MAX(seq), 0)").
		Row().
		Scan(&maxSeq)
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching BOM detail record",
			Err:        err,
		}
	}

	return maxSeq, nil
}

func (r *BomRepositoryImpl) SaveBomDetail(tx *gorm.DB, request masteritempayloads.BomDetailRequest) (masteritementities.BomDetail, *exceptions.BaseErrorResponse) {
	// Check if percentage goes above 100%
	var curBomPercentage float64
	errC := tx.Model(&masteritementities.BomDetail{}).
		Select("COALESCE(SUM(costing_percentage), 0) aa").
		Where("bom_id = ?", request.BomId).
		Where("item_id != ?", request.ItemId).
		Pluck("aa", &curBomPercentage).Error
	if errC != nil {
		return masteritementities.BomDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching BOM detail record",
			Err:        errC,
		}
	}
	if curBomPercentage+request.CostingPercent > 100.0 {
		return masteritementities.BomDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid input BOM detail record",
			Err:        fmt.Errorf("total cost percentage more than 100: %f", curBomPercentage+request.CostingPercent),
		}
	}

	// Find next BomDetailSeq
	var newBomDetailSeq int
	if request.Seq == 0 {
		newBomDetailSeq, errB := r.GetBomDetailMaxSeq(tx, request.BomId)
		if errB != nil {
			return masteritementities.BomDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    errB.Message,
				Err:        errB,
			}
		}
		newBomDetailSeq++ // Tambahkan 1 pada nilai maksimum untuk mendapatkan nilai BomDetailSeq yang baru
	} else {
		newBomDetailSeq = request.Seq
	}
	// Buat entitas BomDetail
	newBomDetail := masteritementities.BomDetail{
		IsActive:          true,
		BomId:             request.BomId,
		Seq:               newBomDetailSeq,
		ItemId:            request.ItemId,
		Qty:               request.Qty,
		Remark:            request.Remark,
		CostingPercentage: request.CostingPercent,
	}

	// Check if incoming request is unique
	var check int64
	errA := tx.Model(&masteritementities.BomDetail{}).
		Where("bom_id = ?", request.BomId).
		Where("item_id = ?", request.ItemId).
		Count(&check).Error
	if errA != nil {
		return masteritementities.BomDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching BOM detail record",
			Err:        errA,
		}
	}
	if check == 0 {
		/// Insert
		err := tx.Create(&newBomDetail).Error
		if err != nil {
			return masteritementities.BomDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error creating BOM detail",
				Err:        err,
			}
		}

		return newBomDetail, nil
	}

	/// Update
	err := tx.Model(&newBomDetail).Select("qty", "remark", "costing_percentage").
		Where("bom_id = ?", request.BomId).
		Where("item_id = ?", request.ItemId).
		Updates(masteritementities.BomDetail{
			Qty:               newBomDetail.Qty,
			Remark:            newBomDetail.Remark,
			CostingPercentage: newBomDetail.CostingPercentage,
		}).Error
	if err != nil {
		return masteritementities.BomDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error updating BOM detail",
			Err:        err,
		}
	}

	return newBomDetail, nil
}

func (r *BomRepositoryImpl) DeleteByIds(tx *gorm.DB, ids []int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.BomDetail

	// Update isActive to false (if needed at some point in time)
	/*
		err := tx.Model(&entities).Select("is_active").
			Where("bom_detail_id IN ?", ids).
			Updates(masteritementities.BomDetail{
				IsActive: false,
			}).Error
	*/

	err := tx.Delete(&entities, ids).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error deleting BOM detail record",
			Err:        err,
		}
	}

	return true, nil
}
