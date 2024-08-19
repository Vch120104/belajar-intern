package test

import (
	"after-sales/api/config"
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	"testing"
	"time"
)

func Test_CarWash(t *testing.T) {
	config.InitEnvConfigs(true, "")
	tx := config.InitDB()

	workOrderSystemNumber := 111

	var workOrderEntity transactionworkshopentities.WorkOrder

	var workOrderResponse transactionjpcbpayloads.CarWashWorkOrder
	err := tx.Model(&workOrderEntity).Select("car_wash, company_id, work_order_status_id").Where("work_order_system_number = ?", workOrderSystemNumber).Scan(&workOrderResponse).Error
	if err != nil {
		panic("Error")
	}

	const (
		QC_PASS = 6
	)

	if true { //TODO check if company use jpcb
		if workOrderResponse.WorkOrderStatusId == QC_PASS {
			if workOrderResponse.CarWash {
				// cehck if work order doesnt exist in trx_car_Wash
				var workOrder int
				result := tx.Model(&transactionjpcbentities.CarWash{}).Select("work_order_system_number").
					Where("work_order_system_number = ?", workOrderSystemNumber).Scan(&workOrder)
				if result.Error != nil {
					panic("Error")
				}
				if result.RowsAffected == 0 {
					newCarWash := transactionjpcbentities.CarWash{
						CompanyId:             workOrderResponse.CompanyId,
						WorkOrderSystemNumber: workOrderSystemNumber,
						StatusId:              1, //Draft
						PriorityId:            2, //Normal
						CarWashDate:           time.Now(),
					}

					err := tx.Create(&newCarWash).Error
					if err != nil {
						panic("error")
					}
				}
			}
		} else {
			if workOrderResponse.CarWash {
				lineTypeOperationId := 1
				resultLineTypeOperation := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).Select("work_order_system_number").
					Where("work_order_system_number = ? AND line_type_id = ?", workOrderSystemNumber, lineTypeOperationId)
				if resultLineTypeOperation.Error != nil {
					panic("error")
				}

				if resultLineTypeOperation.RowsAffected == 0 {
					result := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).Select("work_order_system_number").
						Where("work_order_system_number = ? AND frt_quantity <> supply_quantity", workOrderSystemNumber)
					if result.Error != nil {
						panic("error")
					}
					if result.RowsAffected == 0 {
						result := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).Select("work_order_system_number").
							Where("work_order_system_number = ?", workOrderSystemNumber)
						if result.Error != nil {
							panic("error")
						}
						if result.RowsAffected == 0 {
							newCarWash := transactionjpcbentities.CarWash{
								CompanyId:             workOrderResponse.CompanyId,
								WorkOrderSystemNumber: workOrderSystemNumber,
								StatusId:              1, //Draft
								PriorityId:            2, //Normal
								CarWashDate:           time.Now(),
							}
							err := tx.Create(&newCarWash).Error
							if err != nil {
								panic("error")
							}

						}
					}
				} else {
					var deleteCarWash transactionjpcbentities.CarWash
					result := tx.Model(&deleteCarWash).Select("work_order_system_number").
						Where("work_order_system_number = ?", workOrderSystemNumber).First(&deleteCarWash)
					if result.Error != nil {
						panic("error")
					}

					err = tx.Delete(&deleteCarWash).Error
					if err != nil {
						panic("error")
					}
				}
			}
		}
	}

	panic("error")
}
