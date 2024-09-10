package masterrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/crossservice/financeservice"
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type ItemCycleRepositoryImpl struct {
}

func NewItemCycleRepositoryImpl() masterrepository.ItemCycleRepository {
	return &ItemCycleRepositoryImpl{}
}
func calculatePreviousPeriod(periodMonth, periodYear string) (string, string) {
	month, _ := strconv.Atoi(periodMonth)
	year, _ := strconv.Atoi(periodYear)

	if month == 1 {
		return "12", strconv.Itoa(year - 1)
	}
	return fmt.Sprintf("%02d", month-1), periodYear
}

func ConvertToInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

func ComparePeriods(payloads masterpayloads.ItemCycleInsertPayloads, Responses financeservice.OpenPeriodPayloadResponse) (bool, *exceptions.BaseErrorResponse) {
	periodMonthPayload, err := ConvertToInt(payloads.PeriodMonth)
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid PeriodMonth in payload",
		}
	}

	periodMonthResponse, err := ConvertToInt(Responses.PeriodMonth)
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid PeriodMonth in response",
		}
	}

	periodYearPayload, err := ConvertToInt(payloads.PeriodYear)
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid PeriodYear in payload",
		}
	}
	periodYearResponse, err := ConvertToInt(Responses.PeriodYear)
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid PeriodYear in response",
		}
	}

	if periodMonthPayload > periodMonthResponse && periodYearPayload >= periodYearResponse {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Period isn't open (Item Cycle)",
		}
	}

	return true, nil
}

func (i *ItemCycleRepositoryImpl) InsertItemCycle(db *gorm.DB, payloads masterpayloads.ItemCycleInsertPayloads) (bool, *exceptions.BaseErrorResponse) {
	var PeriodResponse financeservice.OpenPeriodPayloadResponse
	PeriodUrl := config.EnvConfigs.FinanceServiceUrl + "closing-period-company/current-period?company_id=" + strconv.Itoa(payloads.CompanyId) + "&closing_module_detail_code=SP"
	if err := utils.Get(PeriodUrl, &PeriodResponse, nil); err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed to Get Current Period Response data from external service",
			Err:        err,
		}
	}
	valid, errResponse := ComparePeriods(payloads, PeriodResponse)
	if !valid {
		return false, errResponse
	}

	var entities masterentities.ItemCycle
	//take and check is already exist for period
	//err = db.Where("COMPANY_CODE = ? AND PERIOD_YEAR = ? AND PERIOD_MONTH = ? AND ITEM_CODE = ?", companyCode, periodYear, periodMonth, itemCode).First(&existingCycle).Error

	err := db.Model(&entities).Where(masterentities.ItemCycle{PeriodYear: payloads.PeriodYear, CompanyId: payloads.CompanyId, PeriodMonth: payloads.PeriodMonth}).
		First(&entities).Error
	var PrevEntities masterentities.ItemCycle

	if errors.Is(err, gorm.ErrRecordNotFound) {
		periodMonthPrev, periodYearPrev := calculatePreviousPeriod(payloads.PeriodMonth, payloads.PeriodYear)
		errOnPrev := db.Model(&PrevEntities).Where(masterentities.ItemCycle{PeriodYear: periodYearPrev, CompanyId: payloads.CompanyId, PeriodMonth: periodMonthPrev}).
			First(&PrevEntities).Error
		if errOnPrev != nil && !errors.Is(errOnPrev, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errOnPrev}
		}
		newCycle := masterentities.ItemCycle{
			IsActive: true,
			//ItemCycleId:       0,
			CompanyId:         payloads.CompanyId,
			PeriodYear:        payloads.PeriodYear,
			PeriodMonth:       payloads.PeriodMonth,
			ItemId:            payloads.ItemId,
			OrderCycle:        payloads.OrderCycle + PrevEntities.OrderCycle,
			QuantityOnOrder:   payloads.QuantityOnOrder + PrevEntities.QuantityOnOrder,
			QuantityBackOrder: payloads.QuantityBackOrder + PrevEntities.QuantityBackOrder,
		}
		err = db.Create(&newCycle).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err}
		}
	} else {
		entities.OrderCycle = entities.OrderCycle + payloads.OrderCycle
		entities.QuantityOnOrder = entities.QuantityOnOrder + payloads.QuantityOnOrder
		entities.QuantityBackOrder = entities.QuantityBackOrder + payloads.QuantityBackOrder
		err = db.Save(&entities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Failed Update Item Cycle : " + err.Error(),
				Err:        err}
		}
	}
	return true, nil

}
