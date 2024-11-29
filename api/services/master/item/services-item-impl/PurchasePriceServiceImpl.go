package masteritemserviceimpl

import (
	"after-sales/api/config"
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
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
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
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	results, totalPages, totalRows, err := s.PurchasePriceRepo.GetAllPurchasePrice(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *PurchasePriceServiceImpl) UpdatePurchasePrice(id int, req masteritempayloads.PurchasePriceRequest) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	results, err := s.PurchasePriceRepo.UpdatePurchasePrice(tx, id, req)
	if err != nil {
		return masteritementities.PurchasePrice{}, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) SavePurchasePrice(req masteritempayloads.PurchasePriceRequest) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	results, err := s.PurchasePriceRepo.SavePurchasePrice(tx, req)
	if err != nil {
		return masteritementities.PurchasePrice{}, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) GetPurchasePriceById(id int, pagination pagination.Pagination) (masteritempayloads.PurchasePriceResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	results, err := s.PurchasePriceRepo.GetPurchasePriceById(tx, id, pagination)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) AddPurchasePrice(req masteritempayloads.PurchasePriceDetailRequest) (masteritementities.PurchasePriceDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	results, err := s.PurchasePriceRepo.AddPurchasePrice(tx, req)
	if err != nil {
		return masteritementities.PurchasePriceDetail{}, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) GetAllPurchasePriceDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	results, totalPages, totalRows, err := s.PurchasePriceRepo.GetAllPurchasePriceDetail(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *PurchasePriceServiceImpl) GetPurchasePriceDetailById(id int) (masteritempayloads.PurchasePriceDetailResponses, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	results, err := s.PurchasePriceRepo.GetPurchasePriceDetailById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) UpdatePurchasePriceDetail(Id int, req masteritempayloads.PurchasePriceDetailRequest) (masteritementities.PurchasePriceDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	results, err := s.PurchasePriceRepo.UpdatePurchasePriceDetail(tx, Id, req)
	if err != nil {
		return masteritementities.PurchasePriceDetail{}, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) DeletePurchasePrice(id int, iddet []int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	deletemultiid, err := s.PurchasePriceRepo.DeletePurchasePrice(tx, id, iddet)
	if err != nil {
		return false, err
	}
	return deletemultiid, nil
}

func (s *PurchasePriceServiceImpl) ActivatePurchasePriceDetail(id int, iddet []int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	activatemultiid, err := s.PurchasePriceRepo.ActivatePurchasePriceDetail(tx, id, iddet)
	if err != nil {
		return false, err
	}
	return activatemultiid, nil
}

func (s *PurchasePriceServiceImpl) DeactivatePurchasePriceDetail(id int, iddet []int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	deactivatemultiid, err := s.PurchasePriceRepo.DeactivatePurchasePriceDetail(tx, id, iddet)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return deactivatemultiid, nil
}

func (s *PurchasePriceServiceImpl) ChangeStatusPurchasePrice(Id int) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

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
		f.DeleteSheet("Sheet1")
		if err := f.Close(); err != nil {
			log.Error(err)
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
	internalFilterCondition := []utils.FilterCondition{}
	paginate := pagination.Pagination{
		Limit: 10,
		Page:  0,
	}

	results, _, _, errResp := s.PurchasePriceRepo.GetAllPurchasePriceDetail(tx, internalFilterCondition, paginate)
	if errResp != nil {
		return nil, errResp
	}

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

func (s *PurchasePriceServiceImpl) FetchItemId(itemCode string) (int, *exceptions.BaseErrorResponse) {
	resp, err := http.Get(config.EnvConfigs.AfterSalesServiceUrl + "item/by-code?item_code=" + itemCode)
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching item ID",
			Err:        err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: resp.StatusCode,
			Message:    "Error fetching item ID, item code: " + itemCode + " not found in master item service",
		}
	}

	var result struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
		Data       struct {
			ItemId int `json:"item_id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error decoding item ID response",
			Err:        err,
		}
	}

	if result.Data.ItemId == 0 {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Item not found for item code: " + itemCode,
		}
	}

	return result.Data.ItemId, nil
}

func (s *PurchasePriceServiceImpl) PreviewUploadData(rows [][]string, id int) ([]masteritempayloads.PurchasePriceDetailResponses, *exceptions.BaseErrorResponse) {
	var results []masteritempayloads.PurchasePriceDetailResponses
	var numericRegex = regexp.MustCompile(`^\d*\.?\d+$`)

	var regexCheckInput = regexp.MustCompile(`^\d*(,\d{3})*(\.\d{2})?$`)  // handle 12,345,678.99
	var regexCheckInput2 = regexp.MustCompile(`^\d*(\.\d{3})*(,\d{2})?$`) // handle 12.345.678,99

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

		// Debugging row data
		//fmt.Printf("Debugging Row: %v\n", row)

		// Preprocessing purchase price
		purchasePriceStr := strings.TrimSpace(row[2])      // Trim whitespace
		if regexCheckInput.MatchString(purchasePriceStr) { // check if format 12,345,678.99
			purchasePriceStr = strings.ReplaceAll(purchasePriceStr, ",", "") // Delete comma
		} else if regexCheckInput2.MatchString(purchasePriceStr) { //check if format 12.345.678,99
			purchasePriceStr = strings.ReplaceAll(purchasePriceStr, ".", "")  // Delete dot
			purchasePriceStr = strings.ReplaceAll(purchasePriceStr, ",", ".") // Replace comma with dot
		} else {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "format must be either '12,345,678.99' or '12.345.678,99'",
			}
		}

		if !numericRegex.MatchString(purchasePriceStr) {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Purchase price must be numeric (integer or decimal)",
			}
		}

		// Debugging PurchasePrice before validation
		//fmt.Printf("Debugging PurchasePrice: '%s'\n", purchasePriceStr)

		purchasePrice, err := strconv.ParseFloat(purchasePriceStr, 64)
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid purchase price format",
			}
		}

		purchasePriceFloat := float64(purchasePrice)
		results = append(results, masteritempayloads.PurchasePriceDetailResponses{
			ItemCode:        row[0], // Include ItemCode here
			ItemName:        row[1],
			PurchasePrice:   purchasePriceFloat,
			PurchasePriceId: id,
		})
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) ProcessDataUpload(req masteritempayloads.UploadRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	var dupCheck []string

	for i, value := range req.Data {
		itemCode := value.ItemCode
		if len(dupCheck) > 0 && isDuplicate(dupCheck, itemCode) {
			tx.Rollback()
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "item code: " + itemCode + " already exist in excel file at line: " + strconv.Itoa(i+1),
			}
		}
		dupCheck = append(dupCheck, itemCode)

		itemId, errResp := s.FetchItemId(itemCode)
		if errResp != nil {
			tx.Rollback()
			return false, errResp
		}

		requestData := convertToPurchasePriceRequest(value, itemId)

		isExist, errExist := s.PurchasePriceRepo.CheckPurchasePriceDetailExistence(tx, requestData.PurchasePriceId, requestData.ItemId)
		if errExist != nil {
			tx.Rollback()
			return false, errExist
		}
		if isExist {
			tx.Rollback()
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "item code: " + itemCode + " already exist in the current header.",
				Err:        errors.New("item code: " + itemCode + " already exist in the current header"),
			}
		}

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

func isDuplicate(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func convertToPurchasePriceRequest(detail masteritempayloads.PurchasePriceDetailResponses, itemId int) masteritempayloads.PurchasePriceDetailRequest {
	return masteritempayloads.PurchasePriceDetailRequest{
		PurchasePriceDetailId: detail.PurchasePriceDetailId,
		PurchasePriceId:       detail.PurchasePriceId,
		ItemId:                itemId,
		PurchasePrice:         detail.PurchasePrice,
		IsActive:              detail.IsActive,
	}
}

func (s *PurchasePriceServiceImpl) DownloadData(id int) (string, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer func() {
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			log.Error(err)
		}
	}()

	// Fetch Purchase Price data
	purchasePriceData, errResp := s.GetPurchasePriceById(id, pagination.Pagination{Limit: 1000, Page: 0})
	if errResp != nil {
		return "", errResp
	}

	// Log fetched data for debugging
	//log.Infof("Fetched PurchasePriceData: %+v", purchasePriceData)
	//log.Infof("Fetched PurchasePriceDetails: %+v", purchasePriceData.PurchasePriceDetails.Data)

	// Generate a new Excel file
	f := excelize.NewFile()
	sheetName := "purchase_price_detail"
	defer func() {
		if err := f.Close(); err != nil {
			log.Error(err) // Ensure the error is logged if closing fails
		}
	}()

	// Create a new sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", &exceptions.BaseErrorResponse{Err: err, StatusCode: http.StatusInternalServerError}
	}

	// Set cell values
	f.SetCellValue(sheetName, "A1", "Purchase Price Master")
	f.SetCellValue(sheetName, "A2", "Supplier Code")
	f.SetCellValue(sheetName, "B2", "Supplier Name")
	f.SetCellValue(sheetName, "C2", "Item Code")
	f.SetCellValue(sheetName, "D2", "Currency Code")
	f.SetCellValue(sheetName, "E2", "Effective Date")
	f.SetCellValue(sheetName, "F2", "Item Name")
	f.SetCellValue(sheetName, "G2", "Purchase Price")
	f.SetCellValue(sheetName, "H2", "Is Active")
	f.SetColWidth(sheetName, "A", "H", 25)

	// Create a style with bold font and border
	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left"},
		Font: &excelize.Font{
			Bold: true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return "", &exceptions.BaseErrorResponse{Err: err, StatusCode: http.StatusInternalServerError}
	}

	// Apply the style to the header cells
	for col := 'A'; col <= 'H'; col++ {
		cell := string(col) + "1"
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	// Populate the data
	rowNum := 3
	for _, detail := range purchasePriceData.PurchasePriceDetails.Data {

		log.Infof("Populating row %d with detail: %+v", rowNum, detail)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), purchasePriceData.SupplierCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), purchasePriceData.SupplierName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), detail.ItemCode)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), purchasePriceData.CurrencyCode)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), purchasePriceData.PurchasePriceEffectiveDate.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), detail.ItemName)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), detail.PurchasePrice)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), detail.IsActive)
		rowNum++
	}

	f.SetActiveSheet(index)

	tempFilePath := filepath.Join(os.TempDir(), fmt.Sprintf("PurchasePrice_%d.xlsx", id))
	if err := f.SaveAs(tempFilePath); err != nil {
		return "", &exceptions.BaseErrorResponse{Err: err, StatusCode: http.StatusInternalServerError}
	}

	//log.Infof("Excel file saved to: %s", tempFilePath) // Log the file path

	return tempFilePath, nil
}

func ConvertPurchasePriceDetailMapToStruct(maps []map[string]interface{}) ([]masteritempayloads.PurchasePriceDetailResponses, error) {
	var result []masteritempayloads.PurchasePriceDetailResponses

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

func (s *PurchasePriceServiceImpl) GetPurchasePriceDetailByParam(curId int, supId int, effectiveDate string, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	results, err := s.PurchasePriceRepo.GetPurchasePriceDetailByParam(tx, curId, supId, effectiveDate, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}
