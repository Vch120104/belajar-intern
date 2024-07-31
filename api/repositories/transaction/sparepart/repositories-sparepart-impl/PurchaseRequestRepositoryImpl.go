package transactionsparepartrepositoryimpl

import (
	"after-sales/api/config"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type PurchaseRequestRepositoryImpl struct {
}

func NewPurchaseRequestRepositoryImpl() transactionsparepartrepository.PurchaseRequestRepository {
	return &PurchaseRequestRepositoryImpl{}
}

func (p *PurchaseRequestRepositoryImpl) GetAllPurchaseRequest(db *gorm.DB, conditions []utils.FilterCondition, paginationResponses pagination.Pagination, Dateparams map[string]string) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []transactionsparepartpayloads.PurchaseRequestResponses
	entities := transactionsparepartentities.PurchaseRequestEntities{}
	Jointable := db.Table("trx_purchase_request").
		Select("purchase_request_system_number," +
			"purchase_request_document_number," +
			"purchase_request_document_status_id," +
			"item_group_id,purchase_request_document_number," +
			"reference_system_number,order_type_id," +
			"expected_arrival_date,created_by_user_id,reference_document_number") //.Scan(&responses).Error
	WhereQuery := utils.ApplyFilter(Jointable, conditions)
	var strDateFilter string
	if Dateparams["purchase_request_date_from"] == "" {
		Dateparams["purchase_request_date_from"] = "19000101"
	}
	if Dateparams["purchase_request_date_to"] == "" {
		Dateparams["purchase_request_date_to"] = "99991212"
	}
	strDateFilter = "purchase_request_document_date >='" + Dateparams["purchase_request_date_from"] + "' AND purchase_request_document_date <= '" + Dateparams["purchase_request_date_to"] + "'"
	err := WhereQuery.Scopes(pagination.Paginate(&entities, &paginationResponses, WhereQuery)).Where(strDateFilter).Scan(&responses).Error
	if err != nil {
		return paginationResponses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	if len(responses) == 0 {
		return paginationResponses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	var result []transactionsparepartpayloads.PurchaseRequestGetAllListResponses
	for _, res := range responses {
		var purchaseRequestStatusDesc transactionsparepartpayloads.PurchaseRequestStatusResponse
		StatusURL := config.EnvConfigs.GeneralServiceUrl + "document-status/" + strconv.Itoa(res.PurchaseRequestDocumentStatusId)
		if err := utils.Get(StatusURL, &purchaseRequestStatusDesc, nil); err != nil {
			return paginationResponses, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Status data from external service",
				Err:        err,
			}
		}

		var RequestBy transactionsparepartpayloads.PurchaseRequestRequestedByResponse
		RequestByURL := config.EnvConfigs.GeneralServiceUrl + "user-details/" + strconv.Itoa(res.CreatedByUserId)
		if err := utils.Get(RequestByURL, &RequestBy, nil); err != nil {
			return paginationResponses, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Requested By data from external service",
				Err:        err,
			}
		}

		var ItemGroup transactionsparepartpayloads.PurchaseRequestItemGroupResponse
		ItemGroupURL := config.EnvConfigs.GeneralServiceUrl + "item-group/" + strconv.Itoa(res.ItemGroupId)
		if err := utils.Get(ItemGroupURL, &ItemGroup, nil); err != nil {
			return paginationResponses, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item Group data from external service",
				Err:        err,
			}
		}

		var OrderType transactionsparepartpayloads.PurchaseRequestOrderTypeResponse
		OrderTypeURL := config.EnvConfigs.GeneralServiceUrl + "order-type/" + strconv.Itoa(res.OrderTypeId)
		if err := utils.Get(OrderTypeURL, &OrderType, nil); err != nil {
			return paginationResponses, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch customer data from external service",
				Err:        err,
			}
		}
		tempRes := transactionsparepartpayloads.PurchaseRequestGetAllListResponses{
			PurchaseRequestDocumentNumber: res.PurchaseRequestDocumentNumber,
			PurchaseRequestDocumentDate:   res.PurchaseRequestDocumentDate,
			ItemGroup:                     ItemGroup.ItemGroupName,
			OrderType:                     OrderType.OrderTypeName,
			ReferenceNo:                   res.ReferenceDocumentNumber,
			ExpectedArrivalDate:           res.ExpectedArrivalDate,
			Status:                        purchaseRequestStatusDesc.PurchaseRequestStatusDescription,
			RequestBy:                     RequestBy.UserEmployeeName,
		}
		result = append(result, tempRes)
	}
	paginationResponses.Rows = result
	return paginationResponses, nil
}

func (p *PurchaseRequestRepositoryImpl) GetByIdPurchaseRequest(db *gorm.DB, i int) (transactionsparepartpayloads.PurchaseRequestGetByIdNormalizeResponses, *exceptions.BaseErrorResponse) {
	//TODO implement me
	result := transactionsparepartpayloads.PurchaseRequestGetByIdNormalizeResponses{}
	entities := transactionsparepartentities.PurchaseRequestEntities{}
	response := transactionsparepartpayloads.PurchaseRequestGetByIdResponses{}
	rows, err := db.Model(&entities).
		Where(transactionsparepartentities.PurchaseRequestEntities{PurchaseRequestSystemNumber: i}).
		First(&response).
		Rows()
	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()
	//get company name
	var CompanyReponse []transactionsparepartpayloads.PurchaseRequestCompanyResponse
	CompanyURL := config.EnvConfigs.GeneralServiceUrl + "company-id/" + strconv.Itoa(response.CompanyId)
	if err := utils.Get(CompanyURL, &CompanyReponse, nil); err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Company data from external service",
			Err:        err,
		}
	}
	if len(CompanyReponse) == 0 {
		CompanyReponse = append(CompanyReponse, transactionsparepartpayloads.PurchaseRequestCompanyResponse{
			CompanyId:   0,
			CompanyCode: "",
			CompanyName: "",
		})
	}
	var ItemGroup transactionsparepartpayloads.PurchaseRequestItemGroupResponse
	ItemGroupURL := config.EnvConfigs.GeneralServiceUrl + "item-group/" + strconv.Itoa(response.ItemGroupId)
	if err := utils.Get(ItemGroupURL, &ItemGroup, nil); err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item Group data from external service",
			Err:        err,
		}
	}

	var OrderType transactionsparepartpayloads.PurchaseRequestOrderTypeResponse
	OrderTypeURL := config.EnvConfigs.GeneralServiceUrl + "order-type/" + strconv.Itoa(response.OrderTypeId)
	if err := utils.Get(OrderTypeURL, &OrderType, nil); err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch customer data from external service",
			Err:        err,
		}
	}
	var purchaseRequestStatusDesc transactionsparepartpayloads.PurchaseRequestStatusResponse
	StatusURL := config.EnvConfigs.GeneralServiceUrl + "document-status/" + strconv.Itoa(response.PurchaseRequestDocumentStatusId)
	if err := utils.Get(StatusURL, &purchaseRequestStatusDesc, nil); err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Status data from external service",
			Err:        err,
		}
	}

	var GetBrandName transactionsparepartpayloads.PurchaseRequestStatusResponse
	BrandURL := config.EnvConfigs.GeneralServiceUrl + "document-status/" + strconv.Itoa(response.PurchaseRequestDocumentStatusId)
	if err := utils.Get(BrandURL, &GetBrandName, nil); err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Brand data from external service",
			Err:        err,
		}
	}

	var GetDivisionName transactionsparepartpayloads.DivisionResponse
	DivisionURL := config.EnvConfigs.GeneralServiceUrl + "division/" + strconv.Itoa(response.DivisionId)
	if err := utils.Get(DivisionURL, &GetDivisionName, nil); err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Division data from external service",
			Err:        err,
		}
	}

	var GetCostCenterName transactionsparepartpayloads.CostCenterResponses
	CostCenterURL := config.EnvConfigs.GeneralServiceUrl + "cost-center/" + strconv.Itoa(response.CostCenterId)
	if err := utils.Get(CostCenterURL, &GetCostCenterName, nil); err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Cost Center data from external service",
			Err:        err,
		}
	}
	var GetCcyName transactionsparepartpayloads.CurrencyCodeResponse
	CurrencyURL := config.EnvConfigs.FinanceServiceUrl + "currency-code/" + strconv.Itoa(response.CurrencyId)
	if err := utils.Get(CurrencyURL, &GetCcyName, nil); err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Currency data from external service",
			Err:        err,
		}
	}

	var ProfitCenterName transactionsparepartpayloads.ProfitCenterResponses
	ProfitCenterURL := config.EnvConfigs.GeneralServiceUrl + "profit-center/" + strconv.Itoa(response.ProfitCenterId)
	if err := utils.Get(ProfitCenterURL, &ProfitCenterName, nil); err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Profit Center data from external service",
			Err:        err,
		}
	}

	var WarehouseGroupName transactionsparepartpayloads.WarehouseGroupResponses
	WarehouseGroupURL := config.EnvConfigs.AfterSalesServiceUrl + "warehouse-group/" + strconv.Itoa(response.WarehouseGroupId)
	if err := utils.Get(WarehouseGroupURL, &WarehouseGroupName, nil); err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Warehouse data from external service",
			Err:        err,
		}
	}

	var GetWarehouseResponsesName transactionsparepartpayloads.WarehouseResponses
	WarehouseURL := config.EnvConfigs.AfterSalesServiceUrl + "warehouse-master/" + strconv.Itoa(response.WarehouseId)
	if err := utils.Get(WarehouseURL, &GetWarehouseResponsesName, nil); err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Warehouse data from external service",
			Err:        err,
		}
	}
	var RequestBy transactionsparepartpayloads.PurchaseRequestRequestedByResponse
	RequestByURL := config.EnvConfigs.GeneralServiceUrl + "user-details/" + strconv.Itoa(response.CreatedByUserId)
	if err := utils.Get(RequestByURL, &RequestBy, nil); err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Requested By data from external service",
			Err:        err,
		}
	}

	var UpdatedBy transactionsparepartpayloads.PurchaseRequestRequestedByResponse
	UpdatedByURL := config.EnvConfigs.GeneralServiceUrl + "user-details/" + strconv.Itoa(response.UpdatedByUserId)
	if err := utils.Get(UpdatedByURL, &UpdatedBy, nil); err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Requested By data from external service",
			Err:        err,
		}
	}

	refEntities := transactionsparepartentities.PurchaseRequestReferenceType{}
	RefTypeRespons := transactionsparepartpayloads.PurchaseRequestReferenceResponses{}
	row, errs := db.Model(&refEntities).
		Where(transactionsparepartentities.PurchaseRequestReferenceType{ReferencesTypeId: response.ReferenceTypeId}).
		First(&RefTypeRespons).
		Rows()
	row.Close()
	if errs != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	var docNo string

	if response.ReferenceTypeId == 7 {
		var WorkOrder transactionsparepartpayloads.WorkOrderDocNoResponses
		WorkOrderURL := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(response.ReferenceSystemNumber)
		if err := utils.Get(WorkOrderURL, &WorkOrder, nil); err != nil {
			return result, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Requested By data from external service",
				Err:        err,
			}
		}
		docNo = WorkOrder.WorkOrderDocumentNumber
		fmt.Println(docNo)
	}
	//reference type pendiong
	result = transactionsparepartpayloads.PurchaseRequestGetByIdNormalizeResponses{
		Company:                       CompanyReponse[0].CompanyName,
		PurchaseRequestSystemNumber:   response.PurchaseRequestSystemNumber,
		PurchaseRequestDocumentNumber: response.PurchaseRequestDocumentNumber,
		PurchaseRequestDocumentDate:   response.PurchaseRequestDocumentDate,
		PurchaseRequestDocumentStatus: purchaseRequestStatusDesc.PurchaseRequestStatusDescription,
		ItemGroup:                     ItemGroup.ItemGroupName,
		Brand:                         GetBrandName.PurchaseRequestStatusDescription,
		ReferenceType:                 RefTypeRespons.ReferenceTypeName,
		//ReferenceDocumentNumber:       docNo,
		ReferenceDocumentNumber: response.ReferenceDocumentNumber,

		OrderType:                  OrderType.OrderTypeName,
		BudgetCode:                 response.BudgetCode,
		ProjectNo:                  response.ProjectNo,
		Division:                   GetDivisionName.DivisionName,
		PurchaseRequestRemark:      response.PurchaseRequestRemark,
		PurchaseRequestTotalAmount: response.PurchaseRequestTotalAmount,
		ExpectedArrivalDate:        response.ExpectedArrivalDate,
		ExpectedArrivalTime:        response.ExpectedArrivalTime,
		CostCenter:                 GetCostCenterName.CostCenterName,
		ProfitCenter:               ProfitCenterName.ProfitCenterName,
		WarehouseGroup:             WarehouseGroupName.WarehouseGroupName,
		Warehouse:                  GetWarehouseResponsesName.WarehouseName,
		SetOrder:                   response.SetOrder,
		Currency:                   GetCcyName.CurrencyName,
		ChangeNo:                   0,
		CreatedByUser:              RequestBy.UserEmployeeName,
		CreatedDate:                response.CreatedDate,
		UpdatedByUser:              UpdatedBy.UserEmployeeName,
		UpdatedDate:                response.UpdatedDate,
	}
	//panic(tempResult)
	return result, nil
}
func (p *PurchaseRequestRepositoryImpl) GetAllPurchaseRequestDetail(db *gorm.DB, conditions []utils.FilterCondition, paginationResponses pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	//TODO implement me
	entities := transactionsparepartentities.PurchaseRequestDetail{}
	var response []transactionsparepartpayloads.PurchaseRequestDetailRequestPayloads
	Jointable := db.Table("trx_purchase_request_detail").
		Select("item_code,item_quantity,item_remark,item_unit_of_measure")
	WhereQuery := utils.ApplyFilter(Jointable, conditions)
	err := WhereQuery.Scopes(pagination.Paginate(&entities, &paginationResponses, WhereQuery)).Scan(&response).Error
	if err != nil {
		return paginationResponses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	var NormalResponses []transactionsparepartpayloads.PurchaseRequestDetailResponsesPayloads
	for _, res := range response {

		var ItemResponse transactionsparepartpayloads.PurchaseRequestItemResponse
		ItemURL := config.EnvConfigs.AfterSalesServiceUrl + "item/by-code/" + res.ItemCode //strconv.Itoa(response.ItemCode)
		if err := utils.Get(ItemURL, &ItemResponse, nil); err != nil {
			return paginationResponses, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item Group data from external service",
				Err:        err,
			}
		}

		var UomItemResponse transactionsparepartpayloads.UomItemResponses
		UomItem := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/" + res.ItemCode + "/P" //strconv.Itoa(response.ItemCode)
		if err := utils.Get(UomItem, &UomItemResponse, nil); err != nil {
			return paginationResponses, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Uom Item data from external service",
				Err:        err,
			}
		}
		var UomRate float64
		QtyRes := *res.ItemQuantity * *UomItemResponse.TargetConvertion
		if UomItemResponse.SourceConvertion == nil {
			*UomItemResponse.SourceConvertion = 0
		}
		UomRate = QtyRes * *UomItemResponse.SourceConvertion
		UomRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", UomRate), 64)
		result := transactionsparepartpayloads.PurchaseRequestDetailResponsesPayloads{
			ItemCode:              res.ItemCode,
			ItemName:              ItemResponse.ItemName,
			ItemQuantity:          res.ItemQuantity,
			ItemUnitOfMeasure:     res.ItemUnitOfMeasure,
			ItemUnitOfMeasureRate: UomRate,
			ItemRemark:            res.ItemRemark,
		}
		NormalResponses = append(NormalResponses, result)
	}

	paginationResponses.Rows = NormalResponses
	return paginationResponses, nil
}
func (p *PurchaseRequestRepositoryImpl) GetByIdPurchaseRequestDetail(db *gorm.DB, i int) (transactionsparepartpayloads.PurchaseRequestDetailResponsesPayloads, *exceptions.BaseErrorResponse) {
	//TODO implement me

	result := transactionsparepartpayloads.PurchaseRequestDetailResponsesPayloads{}
	entities := transactionsparepartentities.PurchaseRequestDetail{}
	response := transactionsparepartpayloads.PurchaseRequestDetailRequestPayloads{}
	rows, err := db.Model(&entities).
		Where(transactionsparepartentities.PurchaseRequestEntities{PurchaseRequestSystemNumber: i}).
		First(&response).
		Rows()
	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	var ItemResponse transactionsparepartpayloads.PurchaseRequestItemResponse
	ItemURL := config.EnvConfigs.AfterSalesServiceUrl + "item/by-code/" + response.ItemCode //strconv.Itoa(response.ItemCode)
	if err := utils.Get(ItemURL, &ItemResponse, nil); err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item Group data from external service",
			Err:        err,
		}
	}

	result = transactionsparepartpayloads.PurchaseRequestDetailResponsesPayloads{
		ItemCode:              response.ItemCode,
		ItemName:              ItemResponse.ItemName,
		ItemQuantity:          response.ItemQuantity,
		ItemUnitOfMeasure:     response.ItemUnitOfMeasure,
		ItemUnitOfMeasureRate: 0,
		ItemRemark:            response.ItemRemark,
	}
	return result, nil
}

func (p *PurchaseRequestRepositoryImpl) PurchaseRequestSaveHeader(db *gorm.DB, request transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest) (transactionsparepartentities.PurchaseRequestEntities, *exceptions.BaseErrorResponse) {
	purchaserequestentities := transactionsparepartentities.PurchaseRequestEntities{
		CompanyId: request.CompanyId,
		//PurchaseRequestSystemNumber:     request.PurchaseRequestSystemNumber,
		PurchaseRequestDocumentNumber:   request.PurchaseRequestDocumentNumber,
		PurchaseRequestDocumentDate:     &request.PurchaseRequestDocumentDate,
		PurchaseRequestDocumentStatusId: request.PurchaseRequestDocumentStatusId,
		ItemGroupId:                     request.ItemGroupId,
		BrandId:                         request.BrandId,
		ReferenceTypeId:                 request.ReferenceTypeId,
		ReferenceSystemNumber:           request.ReferenceSystemNumber,
		ReferenceDocumentNumber:         request.ReferenceDocumentNumber,
		OrderTypeId:                     request.OrderTypeId,
		BudgetCode:                      request.BudgetCode,
		ProjectNo:                       request.ProjectNo,
		DivisionId:                      request.DivisionId,
		PurchaseRequestRemark:           request.PurchaseRequestRemark,
		PurchaseRequestTotalAmount:      request.PurchaseRequestTotalAmount,
		ExpectedArrivalDate:             &request.ExpectedArrivalTime,
		ExpectedArrivalTime:             &request.ExpectedArrivalTime,
		CostCenterId:                    request.CostCenterId,
		ProfitCenterId:                  request.ProfitCenterId,
		WarehouseGroupId:                request.WarehouseGroupId,
		WarehouseId:                     request.WarehouseId,
		BackOrder:                       request.BackOrder,
		SetOrder:                        request.SetOrder,
		CurrencyId:                      request.CurrencyId,
		ItemClassId:                     request.ItemClassId,
		//PurchaseRequestDetail:           request.purchaserequestde,
		//ChangeNo:        request.ChangeNo,
		CreatedByUserId: request.CreatedByUserId,
		CreatedDate:     &request.CreatedDate,
		UpdatedByUserId: request.UpdatedByUserId,
		UpdatedDate:     &request.UpdatedDate,
	}
	err := db.Create(&purchaserequestentities).Scan(&purchaserequestentities).Error
	if err != nil {
		return purchaserequestentities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Data:       purchaserequestentities,
			Err:        err,
		}
	}
	return purchaserequestentities, nil
}
