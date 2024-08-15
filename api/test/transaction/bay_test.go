package test

import (
	"after-sales/api/config"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	"after-sales/api/utils"
	"fmt"
	"testing"
)

func TestGetByIdSupplySlipDetail(t *testing.T) {
	tx := config.InitDB()

	mainTable := "trx_car_wash"
	mainAlias := "carwash"

	joinTables := []utils.JoinTable{
		{Table: "mtr_car_wash_bay", Alias: "bay", ForeignKey: mainAlias + ".car_wash_bay_id", ReferenceKey: "bay.car_wash_bay_id"},
		{Table: "mtr_car_wash_status", Alias: "status", ForeignKey: mainAlias + ".car_wash_status_id", ReferenceKey: "status.car_wash_status_id"},
		{Table: "trx_work_order", Alias: "wo", ForeignKey: mainAlias + ".work_order_system_number", ReferenceKey: "wo.work_order_system_number"},
	}

	joinQuery := utils.CreateJoin(tx, mainTable, mainAlias, joinTables...)

	keyAttributes := []string{
		"wo.work_order_document_number",
		"bay.car_wash_bay_description",
		"carwash.car_wash_status_id",
		"status.car_wash_status_description",
	}

	var result transactionjpcbpayloads.CarWashErrorDetail
	_ = joinQuery.Select(keyAttributes).Where("wo.work_order_system_number = ?", 1).
		Scan(&result)
	// if joinQuery.Error != nil {
	// 	panic(joinQuery.Error)
	// }
	fmt.Print(result)
}
