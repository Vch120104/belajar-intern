package masteritemserviceimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ItemLocationServiceImpl struct {
	ItemLocationRepo      masteritemrepository.ItemLocationRepository
	WarehouseMasterRepo   masterwarehouserepository.WarehouseMasterRepository
	WarehouseLocationRepo masterwarehouserepository.WarehouseLocationRepository
	ItemRepo              masteritemrepository.ItemRepository
	DB                    *gorm.DB
	RedisClient           *redis.Client // Redis client
}

func StartItemLocationService(
	ItemLocationRepo masteritemrepository.ItemLocationRepository,
	WarehouseMasterRepo masterwarehouserepository.WarehouseMasterRepository,
	WarehouseLocationRepo masterwarehouserepository.WarehouseLocationRepository,
	ItemRepo masteritemrepository.ItemRepository,
	db *gorm.DB,
	redisClient *redis.Client,
) masteritemservice.ItemLocationService {
	return &ItemLocationServiceImpl{
		ItemLocationRepo:      ItemLocationRepo,
		WarehouseMasterRepo:   WarehouseMasterRepo,
		WarehouseLocationRepo: WarehouseLocationRepo,
		ItemRepo:              ItemRepo,
		DB:                    db,
		RedisClient:           redisClient,
	}
}

func (s *ItemLocationServiceImpl) GetAllItemLocation(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.ItemLocationRepo.GetAllItemLocation(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *ItemLocationServiceImpl) SaveItemLocation(req masteritempayloads.ItemLocationRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.ItemLocationRepo.SaveItemLocation(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *ItemLocationServiceImpl) AddItemLocation(id int, req masteritempayloads.ItemLocationDetailRequest) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	err := s.ItemLocationRepo.AddItemLocation(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return err
	}
	return nil
}

func (s *ItemLocationServiceImpl) GetItemLocationById(id int) (masteritempayloads.ItemLocationRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.ItemLocationRepo.GetItemLocationById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemLocationServiceImpl) GetAllItemLocationDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.ItemLocationRepo.GetAllItemLocationDetail(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *ItemLocationServiceImpl) PopupItemLocation(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.ItemLocationRepo.PopupItemLocation(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

// DeleteItemLocation deletes an item location by ID
func (s *ItemLocationServiceImpl) DeleteItemLocation(id int) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	err := s.ItemLocationRepo.DeleteItemLocation(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return err
	}
	return nil
}

func (s *ItemLocationServiceImpl) GetAllItemLoc(filtercondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, totalpages, totalrows, err := s.ItemLocationRepo.GetAllItemLoc(tx, filtercondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, 0, 0, err
	}
	return result, totalpages, totalrows, nil
}

func (s *ItemLocationServiceImpl) GetByIdItemLoc(id int) (masteritempayloads.ItemLocationGetByIdResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ItemLocationRepo.GetByIdItemLoc(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemLocationServiceImpl) SaveItemLoc(req masteritempayloads.SaveItemlocation) (masteritementities.ItemLocation, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ItemLocationRepo.SaveItemLoc(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.ItemLocation{}, err
	}
	return result, nil
}

func (s *ItemLocationServiceImpl) DeleteItemLoc(ids []int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ItemLocationRepo.DeleteItemLoc(tx, ids)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *ItemLocationServiceImpl) GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse) {
	f := excelize.NewFile()
	sheetName := "ItemLocationMaster"
	defer func() {
		f.DeleteSheet("Sheet1")
		if err := f.Close(); err != nil {
			return
		}
	}()

	index, err := f.NewSheet(sheetName)
	if err != nil {
		return f, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	f.SetCellValue(sheetName, "A1", "Warehouse_Code")
	f.SetCellValue(sheetName, "B1", "Warehouse_Location_Code")
	f.SetCellValue(sheetName, "C1", "Item_Code")
	f.SetColWidth(sheetName, "A", "C", 25.0)

	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left"},
		Font: &excelize.Font{
			Bold: true,
		},
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "000000",
				Style: 1,
			},
		},
	})
	if err != nil {
		return f, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	for col := 'A'; col <= 'C'; col++ {
		cell := string(col) + "1"
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	sampleData := 3
	for i := 0; i < sampleData; i++ {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), fmt.Sprintf("WH00%d", i+1))
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), fmt.Sprintf("WH-G%d", i+1))
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), fmt.Sprintf("1/02/RB/BGN/00%d", i+1))
	}

	f.SetActiveSheet(index)

	return f, nil
}

func (s *ItemLocationServiceImpl) UploadPreviewFile(rows [][]string) ([]masteritempayloads.UploadItemLocationResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse
	defer helper.CommitOrRollback(tx, err)
	response := []masteritempayloads.UploadItemLocationResponse{}

	for index, value := range rows {
		data := masteritempayloads.UploadItemLocationResponse{}
		if index > 0 {
			validation := ""
			data.WarehouseCode = value[0]
			data.WarehouseLocationCode = value[1]
			data.ItemCode = value[2]

			_, warehouseErr := s.WarehouseMasterRepo.GetWarehouseMasterByCode(tx, data.WarehouseCode)
			if warehouseErr != nil {
				validation += "Warehouse Master"
			}

			_, warehouseLocErr := s.WarehouseLocationRepo.GetByCode(tx, data.WarehouseLocationCode)
			if warehouseLocErr != nil {
				if validation != "" {
					validation += ", "
				}
				validation += "Warehouse Location"
			}

			_, itemErr := s.ItemRepo.GetItemCode(tx, data.ItemCode)
			if itemErr != nil {
				if validation != "" {
					validation += ", "
				}
				validation += "Item"
			}

			if validation != "" {
				validation = "Data not found in " + validation + "."
			}

			data.Validation = validation

			response = append(response, data)
		} else {
			if value[0] != "Warehouse_Code" || value[1] != "Warehouse_Location_Code" || value[2] != "Item_Code" && len(value) == 3 {
				return response, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Err:        errors.New("make sure header is correct"),
				}
			}
		}
	}

	return response, nil
}

func (s *ItemLocationServiceImpl) UploadProcessFile(uploadPreview []masteritempayloads.UploadItemLocationResponse) ([]masteritementities.ItemLocation, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	response := []masteritementities.ItemLocation{}
	var err *exceptions.BaseErrorResponse
	defer helper.CommitOrRollback(tx, err)

	for _, data := range uploadPreview {
		warehouseData, warehouseErr := s.WarehouseMasterRepo.GetWarehouseMasterByCode(tx, data.WarehouseCode)
		if warehouseErr != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    warehouseErr.Message,
				Err:        warehouseErr.Err,
			}
		}

		warehouseLocData, warehouseLocErr := s.WarehouseLocationRepo.GetByCode(tx, data.WarehouseLocationCode)
		if warehouseLocErr != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    warehouseLocErr.Message,
				Err:        warehouseLocErr.Err,
			}
		}

		itemData, itemErr := s.ItemRepo.GetItemCode(tx, data.ItemCode)
		if itemErr != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    itemErr.Message,
				Err:        itemErr.Err,
			}
		}

		req := masteritempayloads.SaveItemlocation{
			WarehouseGroupId:    warehouseData.WarehouseGroupId,
			ItemId:              itemData.ItemId,
			WarehouseId:         warehouseData.WarehouseId,
			WarehouseLocationId: warehouseLocData.WarehouseLocationId,
		}

		result, err := s.ItemLocationRepo.SaveItemLoc(tx, req)
		if err != nil {
			return response, err
		}
		response = append(response, result)
	}

	return response, nil
}
