package masteritemserviceimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type PurchasePriceServiceImpl struct {
	PurchasePriceRepo masteritemrepository.PurchasePriceRepository
	DB                *gorm.DB
	RedisClient       *redis.Client // Redis client
}

func StartPurchasePriceService(PurchasePriceRepo masteritemrepository.PurchasePriceRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.PurchasePriceService {
	return &PurchasePriceServiceImpl{
		PurchasePriceRepo: PurchasePriceRepo,
		DB:                db,
		RedisClient:       redisClient,
	}
}

func (s *PurchasePriceServiceImpl) GetAllPurchasePrice(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.PurchasePriceRepo.GetAllPurchasePrice(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *PurchasePriceServiceImpl) UpdatePurchasePrice(id int, req masteritempayloads.PurchasePriceRequest) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.PurchasePriceRepo.UpdatePurchasePrice(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.PurchasePrice{}, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) SavePurchasePrice(req masteritempayloads.PurchasePriceRequest) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.PurchasePriceRepo.SavePurchasePrice(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.PurchasePrice{}, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) GetPurchasePriceById(id int, pagination pagination.Pagination) (masteritempayloads.PurchasePriceResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.PurchasePriceRepo.GetPurchasePriceById(tx, id, pagination)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) AddPurchasePrice(req masteritempayloads.PurchasePriceDetailRequest) (masteritementities.PurchasePriceDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.PurchasePriceRepo.AddPurchasePrice(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.PurchasePriceDetail{}, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) GetAllPurchasePriceDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.PurchasePriceRepo.GetAllPurchasePriceDetail(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *PurchasePriceServiceImpl) GetPurchasePriceDetailById(id int) (masteritempayloads.PurchasePriceDetailResponses, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.PurchasePriceRepo.GetPurchasePriceDetailById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) UpdatePurchasePriceDetail(Id int, req masteritempayloads.PurchasePriceDetailRequest) (masteritementities.PurchasePriceDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.PurchasePriceRepo.UpdatePurchasePriceDetail(tx, Id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.PurchasePriceDetail{}, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) DeletePurchasePrice(id int, iddet []int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	deletemultiid, err := s.PurchasePriceRepo.DeletePurchasePrice(tx, id, iddet)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return deletemultiid, nil
}

func (s *PurchasePriceServiceImpl) ChangeStatusPurchasePrice(Id int) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()

	// Ubah status
	entity, err := s.PurchasePriceRepo.ChangeStatusPurchasePrice(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.PurchasePrice{}, err
	}
	return entity, nil
}

func (s *PurchasePriceServiceImpl) GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	// Generate template file
	f := excelize.NewFile()
	sheetName := "purchase_price"
	defer func() {
		if err := f.Close(); err != nil {
			log.Error(err) // Ensure the error is logged if closing fails
		}
	}()

	// Create a new sheet.
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{Err: err, StatusCode: http.StatusInternalServerError}
	}

	// Set value of a cell.
	f.SetCellValue(sheetName, "A1", "Item Code")
	f.SetCellValue(sheetName, "B1", "Item Name")
	f.SetCellValue(sheetName, "C1", "Purchase Price")
	f.SetColWidth(sheetName, "A", "C", 21.5)

	// Create a style with bold font and border
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
		return nil, &exceptions.BaseErrorResponse{Err: err, StatusCode: http.StatusInternalServerError}
	}

	// Apply the style to the header cells
	for col := 'A'; col <= 'C'; col++ {
		cell := string(col) + "1"
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	// Fetch data for the template
	internalFilterCondition := []utils.FilterCondition{} // Adjust as needed
	paginate := pagination.Pagination{
		Limit: 1,
		Page:  1,
	}

	// Ensure GetAllPurchasePrice returns *exceptions.BaseErrorResponse
	results, _, _, errResp := s.PurchasePriceRepo.GetAllPurchasePrice(tx, internalFilterCondition, paginate)
	if errResp != nil {
		return nil, errResp
	}

	// Check if results are nil or empty before proceeding
	if results == nil {
		results = []map[string]interface{}{}
	}

	data, err := ConvertPurchasePriceMapToStruct(results)
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{Err: err, StatusCode: http.StatusInternalServerError}
	}

	for i, value := range data {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), value.ItemCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), value.ItemName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), value.PurchasePrice)
	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	return f, nil
}

func ConvertPurchasePriceMapToStruct(maps []map[string]interface{}) ([]masteritempayloads.PurchasePriceByIdResponse, error) {
	var result []masteritempayloads.PurchasePriceByIdResponse

	// Handle nil or empty maps
	if maps == nil {
		return nil, errors.New("maps is nil")
	}

	jsonData, err := json.Marshal(maps)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PurchasePriceServiceImpl) PreviewUploadData(rows [][]string) ([]masteritempayloads.PurchasePriceDetailResponses, *exceptions.BaseErrorResponse) {
	var results []masteritempayloads.PurchasePriceDetailResponses

	for i, row := range rows {
		if i == 0 {
			// Skip header row
			continue
		}
		if len(row) < 3 {
			// Validate row length
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid row length",
			}
		}
		purchasePrice, err := strconv.Atoi(row[2])
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid purchase price format",
			}
		}
		results = append(results, masteritempayloads.PurchasePriceDetailResponses{
			ItemCode:      row[0],
			ItemName:      row[1],
			PurchasePrice: purchasePrice,
		})
	}

	return results, nil
}

func (s *PurchasePriceServiceImpl) ProcessDataUpload(req masteritempayloads.UploadRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	for _, value := range req.Data {
		// Convert `PurchasePriceDetail` to `PurchasePriceRequest`
		requestData := convertToPurchasePriceRequest(value)

		// Fetch or create PurchasePrice
		_, err := s.PurchasePriceRepo.GetPurchasePriceById(tx, requestData.PurchasePriceId, pagination.Pagination{})
		if err != nil && err.StatusCode != http.StatusNotFound {
			tx.Rollback()
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error checking purchase price existence",
				Err:        err.Err,
			}
		}

		if err != nil && err.StatusCode == http.StatusNotFound {
			// Create new PurchasePrice if it does not exist
			purchasePriceRequest := masteritempayloads.PurchasePriceRequest{
				PurchasePriceId: requestData.PurchasePriceId,
				IsActive:        requestData.IsActive,
			}
			_, err := s.PurchasePriceRepo.SavePurchasePrice(tx, purchasePriceRequest)
			if err != nil {
				tx.Rollback()
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Error creating new purchase price",
					Err:        err.Err,
				}
			}
		}

		// Add purchase price detail
		_, err = s.PurchasePriceRepo.AddPurchasePrice(tx, requestData)
		if err != nil {
			tx.Rollback()
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error adding purchase price detail",
				Err:        err.Err,
			}
		}
	}

	tx.Commit()
	return true, nil
}

func convertToPurchasePriceRequest(detail masteritementities.PurchasePriceDetail) masteritempayloads.PurchasePriceDetailRequest {
	return masteritempayloads.PurchasePriceDetailRequest{
		PurchasePriceDetailId: detail.PurchasePriceDetailId,
		PurchasePriceId:       detail.PurchasePriceId,
		ItemId:                detail.ItemId,
		PurchasePrice:         detail.PurchasePrice,
		IsActive:              detail.IsActive,
	}
}
