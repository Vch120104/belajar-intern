package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"log"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BomRepositoryImpl struct {
}

func StartBomRepositoryImpl() masteritemrepository.BomRepository {
	return &BomRepositoryImpl{}
}

func (r *BomRepositoryImpl) GetBomMasterList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, error) {
	var responses []masteritempayloads.BomMasterListResponse
	var getItemResponse []masteritempayloads.BomItemNameResponse
	var c *gin.Context
	var internalServiceFilter, externalServiceFilter []utils.FilterCondition
	var ItemId string
	responseStruct := reflect.TypeOf(masteritempayloads.BomMasterListResponse{})

	// Memisahkan filter internal dan eksternal
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

	// Membuat log untuk filter internal dan eksternal
	// log.Printf("Filter internal: %+v", internalServiceFilter)
	// log.Printf("Filter eksternal: %+v", externalServiceFilter)

	// Mengambil nilai ItemId dari filter eksternal (jika ada)
	for i := 0; i < len(externalServiceFilter); i++ {
		if externalServiceFilter[i].ColumnField == "item_id" { // Ubah nama kolom sesuai dengan yang digunakan di database
			ItemId = externalServiceFilter[i].ColumnValue
			break
		}
	}

	// Membuat log untuk nilai ItemId
	//log.Printf("ItemId yang diperoleh: %s", ItemId)

	// Define table struct
	tableStruct := masteritempayloads.BomMasterListResponse{}
	// Define join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	// Apply filter
	whereQuery := utils.ApplyFilter(joinTable, internalServiceFilter)
	// Apply pagination and execute
	rows, err := whereQuery.Scan(&responses).Rows()

	if err != nil {
		return nil, 0, 0, err
	}

	defer rows.Close()

	// Logging jumlah data dari tabel responses
	log.Printf("Jumlah data dari tabel responses: %d", len(responses))

	if len(responses) == 0 {
		return nil, 0, 0, gorm.ErrRecordNotFound
	}

	// Mendapatkan data dari URL ItemId
	ItemUrl := "http://localhost:8000/item/?item_id=" + ItemId
	log.Printf("Membuat permintaan ke URL: %s", ItemUrl)
	if err := utils.Get(c, ItemUrl, &getItemResponse, nil); err != nil {
		return nil, 0, 0, err
	}

	// Logging jumlah data dari URL item
	log.Printf("Jumlah data dari URL item: %d", len(getItemResponse))

	// Menggabungkan data dan melakukan paginasi
	joinedData := utils.DataFrameInnerJoin(responses, getItemResponse, "ItemId")
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return dataPaginate, totalPages, totalRows, nil
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
		ItemId:                 request.ItemId,
		BomDetailId:            request.BomDetailId,
		BomMasterQty:           request.BomMasterQty,
		BomMasterUom:           request.BomMasterUom,
		BomMasterEffectiveDate: request.BomMasterEffectiveDate,
	}

	if request.BomMasterId == 0 {
		// Jika BomMasterId == 0, ini adalah operasi membuat data baru
		err := tx.Create(&entities).Error
		if err != nil {
			return false, err
		}
	} else {
		// Jika BomMasterId != 0, ini adalah operasi memperbarui data yang sudah ada
		err := tx.Model(&masteritementities.Bom{}).
			Where("bom_master_id = ?", request.BomMasterId).
			Updates(entities).Error
		if err != nil {
			return false, err
		}
	}

	return true, nil
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
