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

	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type BomServiceImpl struct {
	BomRepository masteritemrepository.BomRepository
	DB            *gorm.DB
	RedisClient   *redis.Client // Redis client
}

func StartBomService(BomRepository masteritemrepository.BomRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.BomService {
	return &BomServiceImpl{
		BomRepository: BomRepository,
		DB:            db,
		RedisClient:   redisClient,
	}
}

func (s *BomServiceImpl) GetBomMasterList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.BomRepository.GetBomMasterList(tx, filterCondition, pages)
	defer func() { helper.CommitOrRollback(tx, err) }()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *BomServiceImpl) GetBomMasterById(id int, pages pagination.Pagination) (masteritempayloads.BomMasterResponseDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.BomRepository.GetBomMasterById(tx, id, pages)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, err
	}

	return results, nil
}

func (s *BomServiceImpl) SaveBomMaster(req masteritempayloads.BomMasterRequest) (masteritementities.Bom, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.BomRepository.SaveBomMaster(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.Bom{}, err
	}
	return results, nil
}

func (s *BomServiceImpl) UpdateBomMaster(id int, req masteritempayloads.BomMasterRequest) (masteritementities.Bom, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.BomRepository.UpdateBomMaster(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.Bom{}, err
	}

	return results, nil
}

func (s *BomServiceImpl) ChangeStatusBomMaster(Id int) (masteritementities.Bom, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	// Ubah status
	entity, err := s.BomRepository.ChangeStatusBomMaster(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.Bom{}, err
	}

	return entity, nil
}

func (s *BomServiceImpl) GetBomDetailList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	//log.Printf("Menerima kondisi filter: %+v", filterCondition) // Tambahkan log untuk menerima kondisi filter
	results, totalPages, totalRows, err := s.BomRepository.GetBomDetailList(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, 0, 0, err
	}

	return results, totalPages, totalRows, nil
}

func (s *BomServiceImpl) GetBomDetailById(id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.BomRepository.GetBomDetailById(tx, id, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *BomServiceImpl) SaveBomDetail(req masteritempayloads.BomDetailRequest) (masteritementities.BomDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.BomRepository.SaveBomDetail(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.BomDetail{}, err
	}
	return results, nil
}

func (s *BomServiceImpl) UpdateBomDetail(id int, req masteritempayloads.BomDetailRequest) (masteritementities.BomDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.BomRepository.UpdateBomDetail(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.BomDetail{}, err
	}
	return results, nil
}

func (s *BomServiceImpl) GetBomItemList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	//log.Printf("Menerima kondisi filter: %+v", filterCondition) // Tambahkan log untuk menerima kondisi filter
	results, totalPages, totalRows, err := s.BomRepository.GetBomItemList(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *BomServiceImpl) DeleteByIds(ids []int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	delete, err := s.BomRepository.DeleteByIds(tx, ids)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return false, err
	}
	return delete, nil
}

func (s *BomServiceImpl) GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	// Generate template file
	f := excelize.NewFile()
	sheetName := "Sheet1"
	defer func() {
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
	f.SetCellValue(sheetName, "A1", "BOM_CODE")
	f.SetCellValue(sheetName, "B1", "EFFECTIVE_DATE")
	f.SetCellValue(sheetName, "C1", "QTY")
	f.SetCellValue(sheetName, "D1", "MATERIAL_CODE")
	f.SetCellValue(sheetName, "E1", "SEQ_DETAIL")
	f.SetCellValue(sheetName, "F1", "QTY_DETAIL")
	f.SetCellValue(sheetName, "G1", "REMARK")
	f.SetCellValue(sheetName, "H1", "COSTING_PERCENTAGE")
	f.SetColWidth(sheetName, "A", "H", 21.5)

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
	for col := 'A'; col <= 'H'; col++ {
		cell := string(col) + "1"
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	// Fetch data for the template
	internalFilterCondition := []utils.FilterCondition{}
	paginate := pagination.Pagination{
		Limit: 10,
		Page:  0,
	}

	results, _, _, errResp := s.BomRepository.GetBomDetailList(tx, internalFilterCondition, paginate)
	if errResp != nil {
		return nil, errResp
	}

	if results == nil {
		results = []map[string]interface{}{}
	}

	data, err := ConvertBomMapToStruct(results)
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{Err: err, StatusCode: http.StatusInternalServerError}
	}

	for i, value := range data {
		// Check if BomDetails.Data has at least one item before accessing it
		if len(value.BomDetails.Data) > 0 {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), value.ItemCode)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), value.BomMasterEffectiveDate)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), value.BomMasterQty)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), value.BomDetails.Data[0].ItemCode)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), value.BomDetails.Data[0].BomDetailSeq)
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", i+2), value.BomDetails.Data[0].BomDetailQty)
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", i+2), value.BomDetails.Data[0].BomDetailRemark)
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", i+2), value.BomDetails.Data[0].BomDetailCostingPercent)
		} else {
			// Handle case where BomDetails.Data is empty
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), value.ItemCode)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), value.BomMasterEffectiveDate)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), value.BomMasterQty)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), "")
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), "")
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", i+2), "")
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", i+2), "")
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", i+2), "")
		}
	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	return f, nil
}

func ConvertBomMapToStruct(maps []map[string]interface{}) ([]masteritempayloads.BomMasterResponseDetail, error) {
	var result []masteritempayloads.BomMasterResponseDetail

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
