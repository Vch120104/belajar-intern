package transactionworkshopentities

const TableNameWorkOrderAllocationGrid = "trx_work_order_allocation_grid"

type WorkOrderAllocationGrid struct {
	TechnicianId       int     `gorm:"column:technician_id;size:30" json:"technician_id"`
	TechnicianName     string  `gorm:"column:technician_name;size:30" json:"technician_name"`
	ShiftCode          string  `gorm:"column:shift_code;" json:"shift_code"`
	TimeAllocation0700 float64 `gorm:"column:time_allocation_0700" json:"time_allocation_0700"`
	TimeAllocation0715 float64 `gorm:"column:time_allocation_0715" json:"time_allocation_0715"`
	TimeAllocation0730 float64 `gorm:"column:time_allocation_0730" json:"time_allocation_0730"`
	TimeAllocation0745 float64 `gorm:"column:time_allocation_0745" json:"time_allocation_0745"`
	TimeAllocation0800 float64 `gorm:"column:time_allocation_0800" json:"time_allocation_0800"`
	TimeAllocation0815 float64 `gorm:"column:time_allocation_0815" json:"time_allocation_0815"`
	TimeAllocation0830 float64 `gorm:"column:time_allocation_0830" json:"time_allocation_0830"`
	TimeAllocation0845 float64 `gorm:"column:time_allocation_0845" json:"time_allocation_0845"`
	TimeAllocation0900 float64 `gorm:"column:time_allocation_0900" json:"time_allocation_0900"`
	TimeAllocation0915 float64 `gorm:"column:time_allocation_0915" json:"time_allocation_0915"`
	TimeAllocation0930 float64 `gorm:"column:time_allocation_0930" json:"time_allocation_0930"`
	TimeAllocation0945 float64 `gorm:"column:time_allocation_0945" json:"time_allocation_0945"`
	TimeAllocation1000 float64 `gorm:"column:time_allocation_1000" json:"time_allocation_1000"`
	TimeAllocation1015 float64 `gorm:"column:time_allocation_1015" json:"time_allocation_1015"`
	TimeAllocation1030 float64 `gorm:"column:time_allocation_1030" json:"time_allocation_1030"`
	TimeAllocation1045 float64 `gorm:"column:time_allocation_1045" json:"time_allocation_1045"`
	TimeAllocation1100 float64 `gorm:"column:time_allocation_1100" json:"time_allocation_1100"`
	TimeAllocation1115 float64 `gorm:"column:time_allocation_1115" json:"time_allocation_1115"`
	TimeAllocation1130 float64 `gorm:"column:time_allocation_1130" json:"time_allocation_1130"`
	TimeAllocation1145 float64 `gorm:"column:time_allocation_1145" json:"time_allocation_1145"`
	TimeAllocation1200 float64 `gorm:"column:time_allocation_1200" json:"time_allocation_1200"`
	TimeAllocation1215 float64 `gorm:"column:time_allocation_1215" json:"time_allocation_1215"`
	TimeAllocation1230 float64 `gorm:"column:time_allocation_1230" json:"time_allocation_1230"`
	TimeAllocation1245 float64 `gorm:"column:time_allocation_1245" json:"time_allocation_1245"`
	TimeAllocation1300 float64 `gorm:"column:time_allocation_1300" json:"time_allocation_1300"`
	TimeAllocation1315 float64 `gorm:"column:time_allocation_1315" json:"time_allocation_1315"`
	TimeAllocation1330 float64 `gorm:"column:time_allocation_1330" json:"time_allocation_1330"`
	TimeAllocation1345 float64 `gorm:"column:time_allocation_1345" json:"time_allocation_1345"`
	TimeAllocation1400 float64 `gorm:"column:time_allocation_1400" json:"time_allocation_1400"`
	TimeAllocation1415 float64 `gorm:"column:time_allocation_1415" json:"time_allocation_1415"`
	TimeAllocation1430 float64 `gorm:"column:time_allocation_1430" json:"time_allocation_1430"`
	TimeAllocation1445 float64 `gorm:"column:time_allocation_1445" json:"time_allocation_1445"`
	TimeAllocation1500 float64 `gorm:"column:time_allocation_1500" json:"time_allocation_1500"`
	TimeAllocation1515 float64 `gorm:"column:time_allocation_1515" json:"time_allocation_1515"`
	TimeAllocation1530 float64 `gorm:"column:time_allocation_1530" json:"time_allocation_1530"`
	TimeAllocation1545 float64 `gorm:"column:time_allocation_1545" json:"time_allocation_1545"`
	TimeAllocation1600 float64 `gorm:"column:time_allocation_1600" json:"time_allocation_1600"`
	TimeAllocation1615 float64 `gorm:"column:time_allocation_1615" json:"time_allocation_1615"`
	TimeAllocation1630 float64 `gorm:"column:time_allocation_1630" json:"time_allocation_1630"`
	TimeAllocation1645 float64 `gorm:"column:time_allocation_1645" json:"time_allocation_1645"`
	TimeAllocation1700 float64 `gorm:"column:time_allocation_1700" json:"time_allocation_1700"`
	TimeAllocation1715 float64 `gorm:"column:time_allocation_1715" json:"time_allocation_1715"`
	TimeAllocation1730 float64 `gorm:"column:time_allocation_1730" json:"time_allocation_1730"`
	TimeAllocation1745 float64 `gorm:"column:time_allocation_1745" json:"time_allocation_1745"`
	TimeAllocation1800 float64 `gorm:"column:time_allocation_1800" json:"time_allocation_1800"`
	TimeAllocation1815 float64 `gorm:"column:time_allocation_1815" json:"time_allocation_1815"`
	TimeAllocation1830 float64 `gorm:"column:time_allocation_1830" json:"time_allocation_1830"`
	TimeAllocation1845 float64 `gorm:"column:time_allocation_1845" json:"time_allocation_1845"`
	TimeAllocation1900 float64 `gorm:"column:time_allocation_1900" json:"time_allocation_1900"`
	TimeAllocation1915 float64 `gorm:"column:time_allocation_1915" json:"time_allocation_1915"`
	TimeAllocation1930 float64 `gorm:"column:time_allocation_1930" json:"time_allocation_1930"`
	TimeAllocation1945 float64 `gorm:"column:time_allocation_1945" json:"time_allocation_1945"`
	TimeAllocation2000 float64 `gorm:"column:time_allocation_2000" json:"time_allocation_2000"`
	TimeAllocation2015 float64 `gorm:"column:time_allocation_2015" json:"time_allocation_2015"`
	TimeAllocation2030 float64 `gorm:"column:time_allocation_2030" json:"time_allocation_2030"`
	TimeAllocation2045 float64 `gorm:"column:time_allocation_2045" json:"time_allocation_2045"`
	TimeAllocation2100 float64 `gorm:"column:time_allocation_2100" json:"time_allocation_2100"`
}

func (*WorkOrderAllocationGrid) TableName() string {
	return TableNameWorkOrderAllocationGrid
}
