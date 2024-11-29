package masterwarehouseserviceimpl

import (
	// masterwarehousepayloads "after-sales/api/payloads/master/warehouse"

	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	// "log"
	// "after-sales/api/utils"
)

type WarehouseLocationServiceImpl struct {
	warehouseLocationRepo  masterwarehouserepository.WarehouseLocationRepository
	warehouseMasterService masterwarehouseservice.WarehouseMasterService
	DB                     *gorm.DB
	RedisClient            *redis.Client // Redis client
}

func OpenWarehouseLocationService(warehouseLocation masterwarehouserepository.WarehouseLocationRepository, warehouseMasterService masterwarehouseservice.WarehouseMasterService, db *gorm.DB, redisClient *redis.Client) masterwarehouseservice.WarehouseLocationService {
	return &WarehouseLocationServiceImpl{
		warehouseLocationRepo:  warehouseLocation,
		warehouseMasterService: warehouseMasterService,
		DB:                     db,
		RedisClient:            redisClient,
	}
}

// ProcessWarehouseLocationTemplate implements masterwarehouseservice.WarehouseLocationService.
func (s *WarehouseLocationServiceImpl) ProcessWarehouseLocationTemplate(req masterwarehousepayloads.ProcessWarehouseLocationTemplate, companyId int) (bool, *exceptions.BaseErrorResponse) {

	for _, value := range req.Data {
		fmt.Println(value.Validation)
		if value.Validation != "Ok" {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: 400,
				Err:        errors.New("failed to process data, check validation"),
			}
		}
	}

	for _, value := range req.Data {
		//GET GROUP ID AND WAREHOUSE ID

		groupId, warehouseId, _ := s.warehouseMasterService.GetWarehouseGroupAndMasterbyCodeandCompanyId(companyId, value.WarehouseCode)
		entities := masterwarehouseentities.WarehouseLocation{
			WarehouseId:                   warehouseId,
			WarehouseGroupId:              groupId,
			WarehouseLocationCode:         value.WarehouseLocationCode,
			WarehouseLocationName:         value.WarehouseLocationName,
			WarehouseLocationDetailName:   "", // FROM PROCEDURE DOCUMENTATION WHICH INSERT WITH DEFAULT VALUE
			WarehouseLocationPickSequence: 0,  // FROM PROCEDURE DOCUMENTATION WHICH INSERT WITH DEFAULT VALUE
			WarehouseLocationCapacityInM3: 0,  // FROM PROCEDURE DOCUMENTATION WHICH INSERT WITH DEFAULT VALUE
		}
		_, err := s.Save(entities)

		if err != nil {
			return false, err
		}

	}

	return true, nil
}

// UploadPreviewFile implements masterwarehouseservice.WarehouseLocationService.
func (s *WarehouseLocationServiceImpl) UploadPreviewFile(rows [][]string, companyId int) ([]masterwarehousepayloads.GetWarehouseLocationPreviewResponse, *exceptions.BaseErrorResponse) {
	response := []masterwarehousepayloads.GetWarehouseLocationPreviewResponse{}
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse
	defer helper.CommitOrRollback(tx, err)

	var warehouseCodes []string
	var warehouseLocationCodes []string
	var warehouseLocationNames []string

	for index, value := range rows {
		data := masterwarehousepayloads.GetWarehouseLocationPreviewResponse{}

		if index > 0 {
			// Check each row's column not empty
			for i := 0; i < 3; i++ {
				if value[i] == "" {
					return response, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Err:        errors.New("make sure column is not empty for each row"),
					}
				}
			}

			data.WarehouseCode = value[0]
			data.WarehouseLocationCode = value[1]
			data.WarehouseLocationName = value[2]

			warehouseCodes = append(warehouseCodes, data.WarehouseCode)
			warehouseLocationCodes = append(warehouseLocationCodes, data.WarehouseLocationCode)
			warehouseLocationNames = append(warehouseLocationNames, data.WarehouseLocationName)

			response = append(response, data)
		} else {
			if value[0] != "WAREHOUSE_CODE" || value[1] != "LOCATION_CODE" || value[2] != "LOCATION_NAME" && len(value) == 3 {
				return response, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Err:        errors.New("make sure header is correct"),
				}
			}
		}
	}

	existData, err := s.warehouseMasterService.IsWarehouseMasterByCodeAndCompanyIdExist(companyId, warehouseCodes)
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error when fetching warehouse company existence",
			Err:        err.Err,
		}
	}
	var existDataWarehouseCodes []string
	for _, warehouse := range existData {
		existDataWarehouseCodes = append(existDataWarehouseCodes, warehouse.WarehouseCode)
	}

	if len(warehouseCodes)+len(warehouseLocationCodes)+len(warehouseLocationNames) == 0 {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("all data is empty"),
		}
	}
	if !(len(warehouseCodes) == len(warehouseLocationCodes) && len(warehouseLocationCodes) == len(warehouseLocationNames)) {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("each code slices has a different length"),
		}
	}

	locationData, err := s.warehouseLocationRepo.CheckIfLocationExist(tx, warehouseCodes, warehouseLocationCodes, warehouseLocationNames)
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error when fetching warehouse location existence",
			Err:        err.Err,
		}
	}

	for i, resp := range response {
		response[i].Validation = "Ok"

		// Check Warehouse Master existence
		if isNotInListString(existDataWarehouseCodes, resp.WarehouseCode) {
			response[i].Validation = "Warehouse Code is invalid"
			continue
		}

		// Fetch Warehouse Id
		var warehouseId int
		for _, data := range existData {
			if resp.WarehouseCode == data.WarehouseCode {
				warehouseId = data.WarehouseId
				break
			}
		}

		// Check Warehouse Location Data existence
		for _, loc := range locationData {
			if warehouseId == loc.WarehouseId &&
				resp.WarehouseLocationCode == loc.WarehouseLocationCode &&
				resp.WarehouseLocationName == loc.WarehouseLocationName {
				response[i].Validation = "Location is already exist"
				break
			}
		}
	}

	return response, nil
}

func isNotInListString(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return false
		}
	}
	return true
}

// GenerateTemplateFile implements masterwarehouseservice.WarehouseLocationService.
func (s *WarehouseLocationServiceImpl) GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse) {
	//Initiate db for get data example
	tx := s.DB.Begin()

	f := excelize.NewFile()
	sheetName := "WarehouseLocation"
	defer func() {
		f.DeleteSheet("Sheet1")
		if err := f.Close(); err != nil {
			return
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return f, &exceptions.BaseErrorResponse{Err: err, StatusCode: 500}
	}
	// Set HEADER
	f.SetCellValue(sheetName, "A1", "WAREHOUSE_CODE")
	f.SetCellValue(sheetName, "B1", "LOCATION_CODE")
	f.SetCellValue(sheetName, "C1", "LOCATION_NAME")
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
		return f, &exceptions.BaseErrorResponse{Err: err, StatusCode: 500}
	}

	// Apply the style to the header cells
	for col := 'A'; col <= 'C'; col++ {
		cell := string(col) + "1"
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	//Get Data Example

	filter := map[string]string{}

	paginate := pagination.Pagination{
		Limit: 3,
		Page:  0,
	}

	filtercondition := utils.BuildFilterCondition(filter)

	result, errGetData := s.warehouseLocationRepo.GetAll(tx, filtercondition, paginate)
	defer helper.CommitOrRollback(tx, errGetData)

	data := result.Rows

	if sliceRows, ok := data.([]masterwarehousepayloads.GetAllWarehouseLocationResponse); ok {
		fmt.Print(data)
		for i := 0; i < len(sliceRows); i++ {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), sliceRows[i].WarehouseCode)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), sliceRows[i].WarehouseLocationCode)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), sliceRows[i].WarehouseLocationName)

		}
	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	return f, nil

}

func (s *WarehouseLocationServiceImpl) Save(request masterwarehouseentities.WarehouseLocation) (bool, *exceptions.BaseErrorResponse) {
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()

	if request.WarehouseLocationId != 0 {
		_, err := s.warehouseLocationRepo.GetById(tx, request.WarehouseLocationId)

		if err != nil {
			return false, err
		}
	}

	save, err := s.warehouseLocationRepo.Save(tx, request)

	if err != nil {
		return false, err
	}
	return save, err
}

func (s *WarehouseLocationServiceImpl) GetById(warehouseLocationId int) (masterwarehousepayloads.GetAllWarehouseLocationResponse, *exceptions.BaseErrorResponse) {
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	get, err := s.warehouseLocationRepo.GetById(tx, warehouseLocationId)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *WarehouseLocationServiceImpl) GetByCode(warehouseLocationCode string) (masterwarehousepayloads.GetAllWarehouseLocationResponse, *exceptions.BaseErrorResponse) {
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	get, err := s.warehouseLocationRepo.GetByCode(tx, warehouseLocationCode)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *WarehouseLocationServiceImpl) GetAll(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	get, err := s.warehouseLocationRepo.GetAll(tx, filter, pages)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *WarehouseLocationServiceImpl) ChangeStatus(warehouseLocationId int) (bool, *exceptions.BaseErrorResponse) {
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()

	_, err = s.warehouseLocationRepo.GetById(tx, warehouseLocationId)

	if err != nil {
		return false, err
	}

	change_status, err := s.warehouseLocationRepo.ChangeStatus(tx, warehouseLocationId)

	if err != nil {
		return change_status, err
	}
	return change_status, nil
}
