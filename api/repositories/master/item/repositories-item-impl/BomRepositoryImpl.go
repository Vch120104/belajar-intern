package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"

	"gorm.io/gorm"
)

type BomRepositoryImpl struct {
}

func StartBomRepositoryImpl() masteritemrepository.BomRepository {
	return &BomRepositoryImpl{}
}

func (r *BomRepositoryImpl) GetBomMasterList(tx *gorm.DB, filters []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, error) {
	var responses []masteritempayloads.BomMasterListResponse

	// Define table struct
	tableStruct := masteritempayloads.BomMasterListResponse{}
	// Define join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	// Apply filters
	whereQuery := utils.ApplyFilter(joinTable, filters)

	// Execute query
	rows, err := whereQuery.Find(&responses).Rows()
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	// Convert responses to maps
	responseMaps := make([]map[string]interface{}, 0)
	for _, response := range responses {
		responseMap := map[string]interface{}{
			"is_active":                 response.IsActive,
			"bom_master_id":             response.BomMasterId,
			"bom_master_qty":            response.BomMasterQty,
			"bom_master_uom":            response.BomMasterUom,
			"bom_master_effective_date": response.BomMasterEffectiveDate,
			"item_code":                 response.ItemCode,
			"item_name":                 response.ItemName,
		}
		responseMaps = append(responseMaps, responseMap)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(responseMaps, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (*BomRepositoryImpl) GetBomMasterById(tx *gorm.DB, id int) (masteritempayloads.BomMasterRequest, error) {
	entities := masteritementities.Bom{}
	response := masteritempayloads.BomMasterRequest{}

	err := tx.Model(&entities).
		Where(masteritementities.Bom{
			BomMasterId: id,
		}).
		First(&response).
		Error

	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *BomRepositoryImpl) SaveBomMaster(tx *gorm.DB, request masteritempayloads.BomMasterRequest) (bool, error) {

	entities := masteritementities.Bom{
		BomMasterId:            request.BomMasterId,
		BomMasterSeq:           request.BomMasterSeq,
		BomMasterQty:           request.BomMasterQty,
		BomMasterUom:           request.BomMasterUom,
		BomMasterEffectiveDate: request.BomMasterEffectiveDate,
		BomMasterChangeNumber:  request.BomMasterChangeNumber,
		ItemId:                 request.ItemId,
	}

	if request.BomMasterId == 0 {
		err := tx.Create(&entities).Error
		if err != nil {
			return false, err // Mengembalikan pesan kesalahan jika terjadi error saat membuat data baru
		}
	} else {
		err := tx.Model(&masteritementities.Bom{}).
			Where("bom_master_id = ?", request.BomMasterId).
			Updates(entities).Error
		if err != nil {
			return false, err // Mengembalikan pesan kesalahan jika terjadi error saat memperbarui data yang sudah ada
		}
	}

	return true, nil // Mengembalikan true jika operasi berhasil tanpa error
}

func (r *BomRepositoryImpl) ChangeStatusBomMaster(tx *gorm.DB, id int) (bool, error) {
	var entities masteritementities.Bom

	result := tx.Model(&entities).
		Where("bom_master_id = ?", id).
		First(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	if entities.IsActive {
		entities.IsActive = false
	} else {
		entities.IsActive = true
	}

	result = tx.Save(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

func (r *BomRepositoryImpl) GetBomDetailList(tx *gorm.DB, filters []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, error) {
	var responses []masteritempayloads.BomDetailListResponse

	// Define table struct
	tableStruct := masteritempayloads.BomDetailListResponse{}
	// Define join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	// Apply filters
	whereQuery := utils.ApplyFilter(joinTable, filters)

	// Execute query
	rows, err := whereQuery.Find(&responses).Rows()
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	// Convert responses to maps
	responseMaps := make([]map[string]interface{}, 0)
	for _, response := range responses {
		responseMap := map[string]interface{}{
			"bom_detail_id":              response.BomDetailId,
			"bom_detail_qty":             response.BomDetailQty,
			"bom_detail_uom":             response.BomDetailUom,
			"bom_detail_seq":             response.BomDetailSeq,
			"bom_detail_remark":          response.BomDetailRemark,
			"bom_detail_costing_percent": response.BomDetailCostingPercent,
		}
		responseMaps = append(responseMaps, responseMap)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(responseMaps, &pages)

	return paginatedData, totalPages, totalRows, nil
}

// Implementasi yang diperbaiki di BomRepositoryImpl
func (*BomRepositoryImpl) GetBomDetailById(tx *gorm.DB, id int) ([]masteritempayloads.BomDetailRequest, error) {
	var entities []masteritementities.BomDetail

	// Mengambil data berdasarkan id
	err := tx.Where(&masteritementities.BomDetail{BomMasterId: id}).Find(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Jika tidak ada data yang ditemukan, kembalikan error not found
			return nil, errors.New("data not found")
		}
		// Jika terjadi error lainnya, kembalikan error tersebut
		return nil, err
	}

	// Jika tidak ada data yang ditemukan, kembalikan error not found
	if len(entities) == 0 {
		return nil, errors.New("data not found")
	}

	// Mengonversi data menjadi slice dari response
	var response []masteritempayloads.BomDetailRequest
	for _, entity := range entities {
		response = append(response, masteritempayloads.BomDetailRequest{
			BomDetailId:             entity.BomDetailId,
			BomDetailSeq:            entity.BomDetailSeq,
			BomDetailQty:            entity.BomDetailQty,
			BomDetailUom:            entity.BomDetailUom,
			BomDetailRemark:         entity.BomDetailRemark,
			BomDetailCostingPercent: entity.BomDetailCostingPct,
			// Isilah properti lainnya sesuai kebutuhan
		})
	}

	// Mengembalikan response
	return response, nil
}
