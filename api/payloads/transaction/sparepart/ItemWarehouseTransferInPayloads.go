package transactionsparepartpayloads

import "time"

type InsertItemWarehouseHeaderTransferInRequest struct {
	CompanyId               int       `json:"company_id"`
	TransferInDate          time.Time `json:"transfer_in_date"`
	TransferOutSystemNumber int       `json:"transfer_out_system_number"`
	EventId                 int       `json:"event_id"`
	WarehouseId             *int      `json:"warehouse_id"`
}

type SubmitItemWarehouseTransferOutGroupStock struct {
	PriceCurrent   float64 `json:"price_current"`
	QuantityEnding float64 `json:"quantity_ending"`
}
