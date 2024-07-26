package masterwarehouseservice

import (
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"github.com/xuri/excelize/v2"
)

type WarehouseLocationService interface {
	Save(masterwarehouseentities.WarehouseLocation) (bool, *exceptions.BaseErrorResponse)
	GetById(int) (masterwarehousepayloads.GetWarehouseLocationResponse, *exceptions.BaseErrorResponse)
	GetAll([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ChangeStatus(int) (bool, *exceptions.BaseErrorResponse)
	GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse)
	UploadPreviewFile(rows [][]string, companyId int) ([]masterwarehousepayloads.GetWarehouseLocationPreviewResponse, *exceptions.BaseErrorResponse)
	ProcessWarehouseLocationTemplate(masterwarehousepayloads.ProcessWarehouseLocationTemplate, int) (bool, *exceptions.BaseErrorResponse)
}
