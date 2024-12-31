package masteritemserviceimpl

import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
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
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type BomServiceImpl struct {
	BomRepository  masteritemrepository.BomRepository
	ItemRepository masteritemrepository.ItemRepository
	DB             *gorm.DB
	RedisClient    *redis.Client // Redis client
}

func StartBomService(BomRepository masteritemrepository.BomRepository, ItemRepository masteritemrepository.ItemRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.BomService {
	return &BomServiceImpl{
		BomRepository:  BomRepository,
		ItemRepository: ItemRepository,
		DB:             db,
		RedisClient:    redisClient,
	}
}

func (s *BomServiceImpl) GetBomList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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

	results, err := s.BomRepository.GetBomList(tx, filterCondition, pages)
	if err != nil {
		return results, err
	}

	return results, nil
}

func (s *BomServiceImpl) GetBomById(id int) (masteritempayloads.BomResponse, *exceptions.BaseErrorResponse) {
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
	results, err := s.BomRepository.GetBomById(tx, id)
	if err != nil {
		return results, err
	}

	return results, nil
}

func (s *BomServiceImpl) ChangeStatusBomMaster(Id int) (masteritementities.Bom, *exceptions.BaseErrorResponse) {
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
	// Invert status
	entity, err := s.BomRepository.ChangeStatusBomMaster(tx, Id)
	if err != nil {
		return masteritementities.Bom{}, err
	}

	return entity, nil
}

func (s *BomServiceImpl) GetBomDetailByMasterId(bomId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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

	results, err := s.BomRepository.GetBomDetailByMasterId(tx, bomId, pages)
	if err != nil {
		return results, err
	}

	return results, nil
}

func (s *BomServiceImpl) GetBomDetailByMasterUn(itemId int, effective_date time.Time, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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

	results, err := s.BomRepository.GetBomDetailByMasterUn(tx, itemId, effective_date, pages)
	if err != nil {
		return results, err
	}

	return results, nil
}

func (s *BomServiceImpl) GetBomDetailById(id int) (masteritementities.BomDetail, *exceptions.BaseErrorResponse) {
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
	results, err := s.BomRepository.GetBomDetailById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *BomServiceImpl) GetBomDetailList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	results, err := s.BomRepository.GetBomDetailList(tx, filterCondition, pages)
	if err != nil {
		return results, err
	}

	return results, nil
}

func (s *BomServiceImpl) GetBomDetailMaxSeq(id int) (int, *exceptions.BaseErrorResponse) {
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

	results, err := s.BomRepository.GetBomDetailMaxSeq(tx, id)
	if err != nil {
		return 0, err
	}

	return results, nil
}

func (s *BomServiceImpl) UpdateBomMaster(id int, qty float64) (masteritementities.Bom, *exceptions.BaseErrorResponse) {
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

	results, err := s.BomRepository.UpdateBomMaster(tx, id, qty)
	if err != nil {
		return masteritementities.Bom{}, err
	}

	return results, nil
}

func (s *BomServiceImpl) SaveBomMaster(req masteritempayloads.BomMasterNewRequest) (masteritementities.Bom, *exceptions.BaseErrorResponse) {
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
	results, err := s.BomRepository.SaveBomMaster(tx, req)
	if err != nil {
		return masteritementities.Bom{}, err
	}
	return results, nil
}

func (s *BomServiceImpl) SaveBomDetail(req masteritempayloads.BomDetailRequest) (masteritementities.BomDetail, *exceptions.BaseErrorResponse) {
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
	// Make bom first before bom detail
	if req.BomId == 0 {
		newRequest := masteritempayloads.BomMasterNewRequest{
			Qty:           req.BomQty,
			EffectiveDate: req.BomEffectiveDate,
			ItemId:        req.BomItemId,
		}

		results, err := s.BomRepository.SaveBomMaster(tx, newRequest)
		if err != nil {
			return masteritementities.BomDetail{}, err
		}
		req.BomId = results.BomId
	}

	results, err := s.BomRepository.SaveBomDetail(tx, req)
	if err != nil {
		return masteritementities.BomDetail{}, err
	}
	return results, nil
}

func (s *BomServiceImpl) UpdateBomDetail(id int, req masteritempayloads.BomDetailRequest) (masteritementities.BomDetail, *exceptions.BaseErrorResponse) {
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
	results, err := s.BomRepository.UpdateBomDetail(tx, id, req)
	if err != nil {
		return masteritementities.BomDetail{}, err
	}
	return results, nil
}

func (s *BomServiceImpl) GetBomItemList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	results, err := s.BomRepository.GetBomItemList(tx, filterCondition, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *BomServiceImpl) DeleteByIds(ids []int) (bool, *exceptions.BaseErrorResponse) {
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
	delete, err := s.BomRepository.DeleteByIds(tx, ids)
	if err != nil {
		return false, err
	}
	return delete, nil
}

func (s *BomServiceImpl) GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	f := excelize.NewFile()
	sheetName := "Sheet1"
	defer func() {
		if err := f.Close(); err != nil {
			log.Error(err)
		}
	}()

	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{Err: err, StatusCode: http.StatusInternalServerError}
	}

	f.SetCellValue(sheetName, "A1", "BOM_CODE") // bom_item_id
	f.SetCellValue(sheetName, "B1", "EFFECTIVE_DATE")
	f.SetCellValue(sheetName, "C1", "QTY")
	f.SetCellValue(sheetName, "D1", "MATERIAL_CODE") // bom_detail_item_id
	f.SetCellValue(sheetName, "E1", "SEQ_DETAIL")
	f.SetCellValue(sheetName, "F1", "QTY_DETAIL")
	f.SetCellValue(sheetName, "G1", "REMARK")
	f.SetCellValue(sheetName, "H1", "COSTING_PERCENTAGE")
	f.SetColWidth(sheetName, "A", "H", 21.0)

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
		return nil, &exceptions.BaseErrorResponse{Err: err, StatusCode: http.StatusInternalServerError}
	}

	for col := 'A'; col <= 'H'; col++ {
		cell := string(col) + "1"
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	internalFilterCondition := []utils.FilterCondition{}
	paginate := pagination.Pagination{
		Limit: 10,
		Page:  0,
	}

	results, errResp := s.BomRepository.GetBomDetailTemplate(tx, internalFilterCondition, paginate)
	if errResp != nil {
		return nil, errResp
	}

	for i, value := range results {
		if len(value.BomDetailItemCode) > 0 {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), value.ItemCode)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), value.EffectiveDate)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), value.Qty)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), value.BomDetailItemCode)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), value.BomDetailSeq)
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", i+2), value.BomDetailQty)
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", i+2), value.BomDetailRemark)
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", i+2), value.BomDetailCostingPercentage)
		} else { // Unlikely case where item detail code doesn't exist
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), value.ItemCode)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), value.EffectiveDate)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), value.Qty)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), "")
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), "")
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", i+2), "")
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", i+2), "")
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", i+2), "")
		}
	}

	f.SetActiveSheet(index)
	return f, nil
}

func ConvertBomMapToStruct(maps []map[string]interface{}) ([]masteritempayloads.BomDetailTemplate, error) {
	var result []masteritempayloads.BomDetailTemplate

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

func (s *BomServiceImpl) FetchItemId(itemCode string) (int, *exceptions.BaseErrorResponse) {
	resp, err := http.Get(config.EnvConfigs.AfterSalesServiceUrl + "item?item_code=" + itemCode + "&limit=1&page=0")
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
		Data       []struct {
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

	if len(result.Data) == 0 {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Item not found for item code: " + itemCode,
		}
	}

	return result.Data[0].ItemId, nil
}

func (s *BomServiceImpl) PreviewUploadData(rows [][]string) ([]masteritempayloads.BomDetailTemplate, *exceptions.BaseErrorResponse) {
	var results []masteritempayloads.BomDetailTemplate
	// var numericRegex = regexp.MustCompile(`^\d*\.?\d+$`)

	// var regexCheckInput = regexp.MustCompile(`^(0[,.]?\d{1,2}|\d*(,\d{3})*(\.\d{1,2})?)$`)  // Handles 12,345,678.99
	// var regexCheckInput2 = regexp.MustCompile(`^(0[,.]?\d{1,2}|\d*(\.\d{3})*(,\d{1,2})?)$`) // Handles 12.345.678,99

	for i, row := range rows {
		if i == 0 {
			// Skip header row
			continue
		}
		if len(row) != 8 {
			// Validate row length
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid row length",
			}
		}

		// Debugging row data
		//fmt.Printf("Debugging Row: %v\n", row)

		// Effective date
		row[1] = strings.TrimSpace(row[1])
		row[1] = strings.ReplaceAll(row[1], "/", "-")
		effectiveDate, err := time.Parse("1-2-06 15:04", row[1])
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid effective date format",
			}
		}

		// Bom qty
		row[2] = strings.ReplaceAll(row[2], ",", ".") // Replace comma with dot
		bomQty, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid bom quantity format",
			}
		}

		// Detail seq
		detailSeq, err := strconv.Atoi(row[4])
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid sequence format",
			}
		}

		// Detail qty
		row[5] = strings.ReplaceAll(row[5], ",", ".") // Replace comma with dot
		detailQty, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid detail quantity format",
			}
		}

		// Cost percentage
		row[7] = strings.ReplaceAll(row[7], ",", ".") // Replace comma with dot
		costPercentage, err := strconv.ParseFloat(row[7], 64)
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid costing percentage format",
			}
		}

		results = append(results, masteritempayloads.BomDetailTemplate{
			ItemCode:                   strings.TrimSpace(row[0]),
			EffectiveDate:              effectiveDate,
			Qty:                        bomQty,
			BomDetailItemCode:          strings.TrimSpace(row[3]),
			BomDetailSeq:               detailSeq,
			BomDetailQty:               detailQty,
			BomDetailRemark:            strings.TrimSpace(row[6]),
			BomDetailCostingPercentage: costPercentage,
		})
	}
	return results, nil
}

func (s *BomServiceImpl) ProcessDataUpload(request masteritempayloads.BomDetailUpload) ([]masteritementities.BomDetail, *exceptions.BaseErrorResponse) {
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

	/// Process bom header
	// Put header to map - any duplicate will not be taken
	bomDetails := request.BomDetails
	type headerKey struct {
		ItemCode      string
		EffectiveDate time.Time
	}
	type headerVal struct {
		Qty    float64
		ItemId int
		BomId  int
	}
	header := map[headerKey]headerVal{}
	for _, bomDetail := range bomDetails {
		key := headerKey{bomDetail.ItemCode, bomDetail.EffectiveDate}
		_, ok := header[key]
		if !ok {
			header[key] = headerVal{bomDetail.Qty, 0, 0}
		}
	}

	for k, v := range header {
		// Get header itemId
		results, err := s.ItemRepository.GetItemCode(tx, k.ItemCode)
		if err != nil {
			return []masteritementities.BomDetail{}, err
		}

		// Insert and get header bomId
		req := masteritempayloads.BomMasterNewRequest{
			Qty:           v.Qty,
			EffectiveDate: k.EffectiveDate,
			ItemId:        results.ItemId,
		}
		bomId, err := s.BomRepository.FirstOrCreateBom(tx, req)
		if err != nil {
			return []masteritementities.BomDetail{}, err
		}

		header[k] = headerVal{v.Qty, results.ItemId, bomId}
	}

	/// Process bom detail
	results := []masteritementities.BomDetail{}
	for _, bomDetail := range bomDetails {
		// Get detail itemId
		itemQuery, err := s.ItemRepository.GetItemCode(tx, bomDetail.BomDetailItemCode)
		if err != nil {
			return []masteritementities.BomDetail{}, err
		}

		// Insert detail
		key := headerKey{bomDetail.ItemCode, bomDetail.EffectiveDate}
		req := masteritempayloads.BomDetailRequest{
			BomId:            header[key].BomId,
			Seq:              bomDetail.BomDetailSeq,
			ItemId:           itemQuery.ItemId,
			Qty:              bomDetail.Qty,
			Remark:           bomDetail.BomDetailRemark,
			CostingPercent:   bomDetail.BomDetailCostingPercentage,
			BomQty:           header[key].Qty,
			BomEffectiveDate: bomDetail.EffectiveDate,
			BomItemId:        header[key].ItemId,
		}
		result, err := s.BomRepository.SaveBomDetail(tx, req)
		if err != nil {
			return []masteritementities.BomDetail{}, err
		}
		results = append(results, result)
	}

	return results, nil
}
