package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type VehicleHIstoryImpl struct {
}

func NewVehicleHistoryImpl() transactionworkshoprepository.VehicleHistoryRepository {
	return &VehicleHIstoryImpl{}
}

func (r *VehicleHIstoryImpl) GetAllVehicleHistory(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []transactionworkshoppayloads.VehicleHistoryResponses
	entities := transactionworkshopentities.WorkOrder{}
	JoinTable := tx.Table("trx_work_order as wo").
		Select("wo.work_order_system_number,wo.work_order_document_number,wo.work_order_date,WO.billable_to_id,ST.work_order_status_description,WO.service_mileage,WO.company_id,WO.customer_id,WO.total_after_vat").
		Joins("Join mtr_work_order_status as ST ON WO.work_order_status_id = ST.work_order_status_id")

	whereQuery := utils.ApplyFilter(JoinTable, filterCondition)
	//whereQuery.Model(entities).Count(&totalRows)
	err := whereQuery.Scopes(pagination.Paginate(&entities, &pages, JoinTable)).Order("wo.work_order_date desc").Scan(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	if len(responses) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	var GetAllResponses []transactionworkshoppayloads.VehicleHistoryGetAllResponses
	for _, res := range responses {
		CustomerURL := config.EnvConfigs.GeneralServiceUrl + "customer/" + strconv.Itoa(res.CustomerId)

		var getCustomerResponse transactionworkshoppayloads.CustomerResponse

		if err := utils.Get(CustomerURL, &getCustomerResponse, nil); err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch customer data from external service",
				Err:        err,
			}
		}
		var getCompanyRespons []transactionworkshoppayloads.WorkOrderDropPoint
		CompanyURL := config.EnvConfigs.GeneralServiceUrl + "company-id/" + strconv.Itoa(res.CompanyId)
		if err := utils.Get(CompanyURL, &getCompanyRespons, nil); err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Company data from external service",
				Err:        err,
			}
		}
		var CompanyName string
		if len(getCompanyRespons) > 0 {
			CompanyName = getCompanyRespons[0].CompanyName
		} else {
			CompanyName = ""
		}
		GetVehicleHistoryRespons := transactionworkshoppayloads.VehicleHistoryGetAllResponses{
			WorkOrderSystemNumber:   res.WorkOrderSystemNumber,
			WorkOrderDocumentNumber: res.WorkOrderDocumentNumber,
			WorkOrderStatusDesc:     res.WorkOrderStatusDescription,
			ServiceMileage:          res.ServiceMileage,
			Company:                 CompanyName,
			Customer:                getCustomerResponse.CustomerName,
			TotalAfterVAT:           res.TotalAfterVAT,
			WorkOrderDate:           res.WorkOrderDate,
		}
		GetAllResponses = append(GetAllResponses, GetVehicleHistoryRespons)
	}
	//paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(responses, &pages)
	//return pagination.Pagination{
	//	Page:       pages.Page,
	//	Limit:      pages.Limit,
	//	SortOf:     pages.SortOf,
	//	SortBy:     pages.SortBy,
	//	TotalRows:  int64(totalRows),
	//	TotalPages: totalPages,
	//	Rows:       paginatedData,
	//}, nil
	pages.Rows = GetAllResponses
	return pages, nil
}
func (r *VehicleHIstoryImpl) GetAllVehicleHistoryChassis(tx *gorm.DB, VehicleHistoryRequest transactionworkshoppayloads.VehicleHistoryChassisRequest, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	//var getCompanyRespons []transactionworkshoppayloads.WorkOrderCompanyDetail
	//CompanyURL := config.EnvConfigs.GeneralServiceUrl + "company-id/" + strconv.Itoa(res.CompanyId)
	//if err := utils.Get(CompanyURL, &getCompanyRespons, nil); err != nil {
	//	return pages, &exceptions.BaseErrorResponse{
	//		StatusCode: http.StatusInternalServerError,
	//		Message:    "Failed to fetch Company data from external service",
	//		Err:        err,
	//	}
	//}

	return pages, nil
}
func (r *VehicleHIstoryImpl) GetVehicleHistoryById(db *gorm.DB, id int) (transactionworkshoppayloads.VehicleHistoryByIdResponses, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.WorkOrder{}
	//only sending work order system number to be used in next page in work_order_page4
	response := transactionworkshoppayloads.VehicleHistoryByIdResponses{}

	rows, err := db.Model(&entities).
		Where(transactionworkshopentities.WorkOrder{
			WorkOrderSystemNumber: id,
		}).
		First(&response).
		Rows()
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()
	return response, nil
}
