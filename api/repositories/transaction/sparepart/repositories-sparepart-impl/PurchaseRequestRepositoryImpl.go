package transactionsparepartrepositoryimpl

import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/crossservice/financeservice"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	Jointable := db.Table("trx_purchase_request A").
		Select("purchase_request_system_number," +
			"purchase_request_document_number," +
			"purchase_request_document_status_id," +
			"item_group_id,purchase_request_document_number," +
			"reference_system_number,order_type_id,purchase_request_document_date," +
			"expected_arrival_date,created_by_user_id,reference_document_number," +
			"'A.purchase_request_system_number'") //.Scan(&responses).Error
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
			PurchaseRequestSystemNumber:   res.PurchaseRequestSystemNumber,
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

func (p *PurchaseRequestRepositoryImpl) GetByIdPurchaseRequest(db *gorm.DB, i int) (transactionsparepartpayloads.PurchaseRequestGetByIdResponses, *exceptions.BaseErrorResponse) {
	//TODO implement me
	//result := transactionsparepartpayloads.PurchaseRequestGetByIdNormalizeResponses{}
	entities := transactionsparepartentities.PurchaseRequestEntities{}
	response := transactionsparepartpayloads.PurchaseRequestGetByIdResponses{}
	rows, err := db.Model(&entities).
		Where(transactionsparepartentities.PurchaseRequestEntities{PurchaseRequestSystemNumber: i}).
		First(&response).
		Rows()
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()
	//get company name
	var CompanyReponse []transactionsparepartpayloads.PurchaseRequestCompanyResponse
	CompanyURL := config.EnvConfigs.GeneralServiceUrl + "company-id/" + strconv.Itoa(response.CompanyId)
	if err := utils.Get(CompanyURL, &CompanyReponse, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
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
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item Group data from external service",
			Err:        err,
		}
	}

	var OrderType transactionsparepartpayloads.PurchaseRequestOrderTypeResponse
	OrderTypeURL := config.EnvConfigs.GeneralServiceUrl + "order-type/" + strconv.Itoa(response.OrderTypeId)
	if err := utils.Get(OrderTypeURL, &OrderType, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch customer data from external service",
			Err:        err,
		}
	}
	var purchaseRequestStatusDesc transactionsparepartpayloads.PurchaseRequestStatusResponse
	StatusURL := config.EnvConfigs.GeneralServiceUrl + "document-status/" + strconv.Itoa(response.PurchaseRequestDocumentStatusId)
	if err := utils.Get(StatusURL, &purchaseRequestStatusDesc, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Status data from external service",
			Err:        err,
		}
	}

	var GetBrandName transactionsparepartpayloads.PurchaseRequestStatusResponse
	BrandURL := config.EnvConfigs.GeneralServiceUrl + "document-status/" + strconv.Itoa(response.PurchaseRequestDocumentStatusId)
	if err := utils.Get(BrandURL, &GetBrandName, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Brand data from external service",
			Err:        err,
		}
	}

	var GetDivisionName transactionsparepartpayloads.DivisionResponse
	DivisionURL := config.EnvConfigs.GeneralServiceUrl + "division/" + strconv.Itoa(response.DivisionId)
	if err := utils.Get(DivisionURL, &GetDivisionName, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Division data from external service",
			Err:        err,
		}
	}

	var GetCostCenterName transactionsparepartpayloads.CostCenterResponses
	CostCenterURL := config.EnvConfigs.GeneralServiceUrl + "cost-center/" + strconv.Itoa(response.CostCenterId)
	if err := utils.Get(CostCenterURL, &GetCostCenterName, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Cost Center data from external service",
			Err:        err,
		}
	}
	var GetCcyName transactionsparepartpayloads.CurrencyCodeResponse
	CurrencyURL := config.EnvConfigs.FinanceServiceUrl + "currency-code/" + strconv.Itoa(response.CurrencyId)
	if err := utils.Get(CurrencyURL, &GetCcyName, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Currency data from external service",
			Err:        err,
		}
	}

	var ProfitCenterName transactionsparepartpayloads.ProfitCenterResponses
	ProfitCenterURL := config.EnvConfigs.GeneralServiceUrl + "profit-center/" + strconv.Itoa(response.ProfitCenterId)
	if err := utils.Get(ProfitCenterURL, &ProfitCenterName, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Profit Center data from external service",
			Err:        err,
		}
	}

	var WarehouseGroupName transactionsparepartpayloads.WarehouseGroupResponses
	WarehouseGroupURL := config.EnvConfigs.AfterSalesServiceUrl + "warehouse-group/" + strconv.Itoa(response.WarehouseGroupId)
	if err := utils.Get(WarehouseGroupURL, &WarehouseGroupName, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Warehouse data from external service",
			Err:        err,
		}
	}

	var GetWarehouseResponsesName transactionsparepartpayloads.WarehouseResponses
	WarehouseURL := config.EnvConfigs.AfterSalesServiceUrl + "warehouse-master/" + strconv.Itoa(response.WarehouseId)
	if err := utils.Get(WarehouseURL, &GetWarehouseResponsesName, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Warehouse data from external service",
			Err:        err,
		}
	}
	var RequestBy transactionsparepartpayloads.PurchaseRequestRequestedByResponse
	RequestByURL := config.EnvConfigs.GeneralServiceUrl + "user-details/" + strconv.Itoa(response.CreatedByUserId)
	if err := utils.Get(RequestByURL, &RequestBy, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Requested By data from external service",
			Err:        err,
		}
	}

	var UpdatedBy transactionsparepartpayloads.PurchaseRequestRequestedByResponse
	UpdatedByURL := config.EnvConfigs.GeneralServiceUrl + "user-details/" + strconv.Itoa(response.UpdatedByUserId)
	if err := utils.Get(UpdatedByURL, &UpdatedBy, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Requested By data from external service",
			Err:        err,
		}
	}

	var PurchaseRequestReferenceType transactionsparepartpayloads.PurchaseRequestReferenceType
	PurchaseReuqestReferenceType := config.EnvConfigs.GeneralServiceUrl + "reference-type-purchase-request/" + strconv.Itoa(response.ReferenceTypeId)
	if err := utils.Get(PurchaseReuqestReferenceType, &PurchaseRequestReferenceType, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed to fetch Reference from external service",
			Err:        err,
		}
	}
	return response, nil
	//result = transactionsparepartpayloads.PurchaseRequestGetByIdNormalizeResponses{
	//	Company:                       CompanyReponse[0].CompanyName,
	//	PurchaseRequestSystemNumber:   response.PurchaseRequestSystemNumber,
	//	PurchaseRequestDocumentNumber: response.PurchaseRequestDocumentNumber,
	//	PurchaseRequestDocumentDate:   response.PurchaseRequestDocumentDate,
	//	PurchaseRequestDocumentStatus: purchaseRequestStatusDesc.PurchaseRequestStatusDescription,
	//	ItemGroup:                     ItemGroup.ItemGroupName,
	//	Brand:                         GetBrandName.PurchaseRequestStatusDescription,
	//	ReferenceType:                 PurchaseRequestReferenceType.ReferenceTypePurchaseRequestName,
	//	//ReferenceDocumentNumber:       docNo,
	//	ReferenceDocumentNumber: response.ReferenceDocumentNumber,
	//
	//	OrderType:                  OrderType.OrderTypeName,
	//	BudgetCode:                 response.BudgetCode,
	//	ProjectNo:                  response.ProjectNo,
	//	Division:                   GetDivisionName.DivisionName,
	//	PurchaseRequestRemark:      response.PurchaseRequestRemark,
	//	PurchaseRequestTotalAmount: response.PurchaseRequestTotalAmount,
	//	ExpectedArrivalDate:        response.ExpectedArrivalDate,
	//	ExpectedArrivalTime:        response.ExpectedArrivalTime,
	//	CostCenter:                 GetCostCenterName.CostCenterName,
	//	ProfitCenter:               ProfitCenterName.ProfitCenterName,
	//	WarehouseGroup:             WarehouseGroupName.WarehouseGroupName,
	//	Warehouse:                  GetWarehouseResponsesName.WarehouseName,
	//	SetOrder:                   response.SetOrder,
	//	Currency:                   GetCcyName.CurrencyName,
	//	ChangeNo:                   0,
	//	CreatedByUser:              RequestBy.UserEmployeeName,
	//	CreatedDate:                response.CreatedDate,
	//	UpdatedByUser:              UpdatedBy.UserEmployeeName,
	//	UpdatedDate:                response.UpdatedDate,
	//}
	//fmt.Println(result)
}
func (p *PurchaseRequestRepositoryImpl) GetAllPurchaseRequestDetail(db *gorm.DB, conditions []utils.FilterCondition, paginationResponses pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	//TODO implement me
	entities := transactionsparepartentities.PurchaseRequestDetail{}
	var response []transactionsparepartpayloads.PurchaseRequestDetailRequestPayloads
	Jointable := db.Table("trx_purchase_request_detail").
		Select("item_id,item_quantity,item_remark,item_unit_of_measure,purchase_request_system_number,purchase_request_line_number,reference_system_number,reference_line,purchase_request_detail_system_number")
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

		var ItemResponse masteritementities.Item
		err = db.Model(&ItemResponse).Where(masteritementities.Item{ItemId: res.ItemId}).Scan(&ItemResponse).Error
		if err != nil {
			return paginationResponses, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Item Not Found",
				Err:        err,
			}
		}
		var UomItemResponse masteritementities.UomItem
		err = db.Model(&UomItemResponse).Where(masteritementities.UomItem{ItemId: ItemResponse.ItemId, UomSourceTypeCode: "P"}).
			First(&UomItemResponse).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				UomItemResponse.SourceConvertion = 0
				UomItemResponse.TargetConvertion = 0
			} else {
				return paginationResponses, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Message:    "Cannot Get Uom target Id",
					Err:        err,
				}
			}

		}
		var UomRate float64
		var QtyRes float64
		if UomItemResponse.SourceConvertion == 0 {
			QtyRes = 0
		} else {
			QtyRes = *res.ItemQuantity * UomItemResponse.TargetConvertion

		}
		if UomItemResponse.SourceConvertion == 0 {
			return paginationResponses, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Failed to fetch Uom Source Convertion From External Data",
				Err:        err,
			}
		}
		UomRate = QtyRes * UomItemResponse.SourceConvertion // QtyRes * *UomItemResponse.SourceConvertion
		UomRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", UomRate), 64)
		result := transactionsparepartpayloads.PurchaseRequestDetailResponsesPayloads{
			PurchaseRequestDetailSystemNumber: res.PurchaseRequestDetailSystemNumber,
			PurchaseRequestSystemNumber:       res.PurchaseRequestSystemNumber,
			PurchaseRequestLineNumber:         res.PurchaseRequestLineNumber,
			ReferenceSystemNumber:             res.ReferenceSystemNumber,
			ReferenceLine:                     res.ReferenceLine,
			ItemCode:                          ItemResponse.ItemCode,
			ItemName:                          ItemResponse.ItemName,
			ItemQuantity:                      res.ItemQuantity,
			ItemUnitOfMeasure:                 res.ItemUnitOfMeasure,
			ItemUnitOfMeasureRate:             UomRate,
			ItemRemark:                        res.ItemRemark,
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
		Where(transactionsparepartentities.PurchaseRequestDetail{PurchaseRequestDetailSystemNumber: i}).
		First(&response).
		Rows()
	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	var ItemResponse masteritementities.Item
	err = db.Model(&ItemResponse).Where(masteritementities.Item{ItemId: response.ItemId}).Scan(&ItemResponse).Error
	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Item Not Found",
			Err:        err,
		}
	}
	UomItem := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/" + strconv.Itoa(ItemResponse.ItemId) + "/P" //strconv.Itoa(response.ItemCode)
	var UomItemResponse transactionsparepartpayloads.UomItemResponses

	if err := utils.Get(UomItem, &UomItemResponse, nil); err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Uom Item data from external service",
			Err:        err,
		}
	}
	var UomRate float64
	var QtyRes float64
	if UomItemResponse.SourceConvertion == nil {
		QtyRes = 0
	} else {
		QtyRes = *response.ItemQuantity * *UomItemResponse.TargetConvertion

	}
	if UomItemResponse.SourceConvertion == nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Uom Source Convertion From External Data",
			Err:        err,
		}
	}
	UomRate = QtyRes * *UomItemResponse.SourceConvertion // QtyRes * *UomItemResponse.SourceConvertion
	UomRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", UomRate), 64)

	result = transactionsparepartpayloads.PurchaseRequestDetailResponsesPayloads{
		PurchaseRequestDetailSystemNumber: response.PurchaseRequestDetailSystemNumber,
		PurchaseRequestSystemNumber:       response.PurchaseRequestSystemNumber,
		PurchaseRequestLineNumber:         response.PurchaseRequestLineNumber,
		ReferenceSystemNumber:             response.ReferenceSystemNumber,
		ItemId:                            response.ItemId,
		ReferenceLine:                     response.ReferenceLine,
		ItemCode:                          ItemResponse.ItemCode,
		ItemName:                          ItemResponse.ItemName,
		ItemQuantity:                      response.ItemQuantity,
		ItemUnitOfMeasure:                 response.ItemUnitOfMeasure,
		ItemUnitOfMeasureRate:             UomRate,
		ItemRemark:                        response.ItemRemark,
		CreatedByUserId:                   response.CreatedByUserId,
		CreatedDate:                       response.CreatedDate,
		UpdatedByUserId:                   response.UpdatedByUserId,
		UpdatedDate:                       response.UpdatedDate,
	}
	return result, nil
}

func (p *PurchaseRequestRepositoryImpl) NewPurchaseRequestHeader(db *gorm.DB, request transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest) (transactionsparepartentities.PurchaseRequestEntities, *exceptions.BaseErrorResponse) {
	//cek id for status code draft 10
	var DocResponse transactionsparepartpayloads.PurchaseRequestDocumentStatus
	DocumentStatusUrl := config.EnvConfigs.GeneralServiceUrl + "document-status-by-code/10"

	if err := utils.Get(DocumentStatusUrl, &DocResponse, nil); err != nil {
		return transactionsparepartentities.PurchaseRequestEntities{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to DcoumentStatusId",
			Err:        err,
		}
	}

	purchaserequestentities := transactionsparepartentities.PurchaseRequestEntities{
		CompanyId:                       request.CompanyId,
		PurchaseRequestDocumentNumber:   request.PurchaseRequestDocumentNumber,
		PurchaseRequestDocumentDate:     &request.PurchaseRequestDocumentDate,
		PurchaseRequestDocumentStatusId: DocResponse.DocumentStatusId, //request.PurchaseRequestDocumentStatusId,
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

		ChangeNo:        1,
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
func (p *PurchaseRequestRepositoryImpl) NewPurchaseRequestDetail(db *gorm.DB, payloads transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads) (transactionsparepartentities.PurchaseRequestDetail, *exceptions.BaseErrorResponse) {
	//get header data
	var Response = transactionsparepartentities.PurchaseRequestDetail{}
	entities := transactionsparepartentities.PurchaseRequestEntities{}
	rows, err := db.Model(&entities).
		Where(transactionsparepartentities.PurchaseRequestEntities{PurchaseRequestSystemNumber: payloads.PurchaseRequestSystemNumber}).
		First(&entities).
		Rows()
	if err != nil {
		return Response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()
	LineNumber := 1
	//default value
	LineStatus := "10"
	Price := 0.000
	//validation
	TotalPrice := Price * *payloads.ItemQuantity
	PRDetailEntities := transactionsparepartentities.PurchaseRequestDetail{
		//PurchaseRequestSystemNumberDetail: 0,
		PurchaseRequestSystemNumber: payloads.PurchaseRequestSystemNumber,
		PurchaseRequestLineNumber:   LineNumber,
		PurchaseRequestLineStatus:   LineStatus,
		ItemCode:                    payloads.ItemCode,
		ItemQuantity:                payloads.ItemQuantity,
		ItemUnitOfMeasure:           payloads.ItemUnitOfMeasure,
		ItemPrice:                   &Price,
		ItemTotal:                   &TotalPrice,
		ItemRemark:                  payloads.ItemRemark,
		PurchaseOrderSystemNumber:   0,
		PurchaseOrderLine:           0,
		ReferenceTypeId:             entities.ReferenceTypeId,
		ReferenceSystemNumber:       entities.ReferenceSystemNumber,
		ReferenceLine:               payloads.ReferenceLine,
		VehicleId:                   0,
		CreatedByUserId:             payloads.CreatedByUserId,
		CreatedDate:                 &payloads.CreatedDate,
		UpdatedByUserId:             payloads.CreatedByUserId,
		UpdatedDate:                 &payloads.UpdatedDate,
		ChangeNo:                    1,
	}
	errCreate := db.Create(&PRDetailEntities).Scan(&PRDetailEntities).Error
	if errCreate != nil {
		return PRDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    errCreate.Error(),
			Data:       PRDetailEntities,
			Err:        errCreate,
		}
	}
	return PRDetailEntities, nil
}
func (p *PurchaseRequestRepositoryImpl) SavePurchaseRequestHeader(db *gorm.DB, request transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, id int) (transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, *exceptions.BaseErrorResponse) {
	//TODO implement me
	res := transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest{}
	entities := transactionsparepartentities.PurchaseRequestEntities{}
	err := db.Model(&entities).Where(transactionsparepartentities.PurchaseRequestEntities{PurchaseRequestSystemNumber: id}).First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Data to Updated Is Not Found",
				Data:       res,
				Err:        err,
			}
		}
		return res, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Update Data Failed",
			Data:       res,
			Err:        err,
		}
	}
	//updating data
	entities.CompanyId = request.CompanyId
	entities.BudgetCode = request.BudgetCode
	entities.ProjectNo = request.ProjectNo
	entities.DivisionId = request.DivisionId
	entities.PurchaseRequestRemark = request.PurchaseRequestRemark
	entities.ExpectedArrivalTime = &request.ExpectedArrivalTime
	entities.ExpectedArrivalDate = &request.ExpectedArrivalDate
	entities.CostCenterId = request.CostCenterId
	entities.ProfitCenterId = request.ProfitCenterId
	entities.BackOrder = request.BackOrder
	entities.CurrencyId = request.CurrencyId
	entities.SetOrder = request.SetOrder
	entities.OrderTypeId = request.OrderTypeId
	entities.ChangeNo = entities.ChangeNo + 1
	entities.UpdatedByUserId = request.UpdatedByUserId
	entities.UpdatedDate = &request.UpdatedDate
	err = db.Save(&entities).Error
	//db.Commit()
	if err != nil {

		return res, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed To Update Data",
			Data:       res,
			Err:        err,
		}
	}
	return request, nil

}

func (p *PurchaseRequestRepositoryImpl) SavePurchaseRequestDetail(db *gorm.DB, payloads transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, id int) (transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, *exceptions.BaseErrorResponse) {
	//TODO implement me
	response := transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads{}
	entities := transactionsparepartentities.PurchaseRequestDetail{}

	err := db.Model(&entities).Where(transactionsparepartentities.PurchaseRequestDetail{PurchaseRequestDetailSystemNumber: id}).First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data To Updated Is Not Found Try Insert Instead",
				Data:       nil,
				Err:        err,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Data Updated Failed",
			Data:       nil,
			Err:        err,
		}
	}
	entities.ItemQuantity = payloads.ItemQuantity
	entities.ItemRemark = payloads.ItemRemark
	entities.ChangeNo = entities.ChangeNo + 1
	entities.UpdatedDate = &payloads.UpdatedDate
	entities.UpdatedByUserId = payloads.UpdatedByUserId
	err = db.Save(&entities).Error
	fmt.Println()
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Update Data Failed",
			Data:       nil,
			Err:        err,
		}
	}
	return payloads, nil
}

func (p *PurchaseRequestRepositoryImpl) VoidPurchaseRequest(db *gorm.DB, i int) (bool, *exceptions.BaseErrorResponse) {
	entities := transactionsparepartentities.PurchaseRequestEntities{}
	err := db.Model(&entities).Where(transactionsparepartentities.PurchaseRequestEntities{PurchaseRequestSystemNumber: i}).First(&entities).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Header to delete Is Not Found",
				Data:       i,
				Err:        err,
			}
		} else {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Delete Data Failed",
				Data:       i,
				Err:        err,
			}
		}
	}
	var DocResponse transactionsparepartpayloads.PurchaseRequestDocumentStatus
	//draft = 10
	DocumentStatusUrl := config.EnvConfigs.GeneralServiceUrl + "document-status-by-code/10"
	//UomItem := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/" + res.ItemCode + "/P" //strconv.Itoa(response.ItemCode)
	if err := utils.Get(DocumentStatusUrl, &DocResponse, nil); err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed to Fetch Document Status From General Service",
			Err:        err,
		}
	}

	DraftPurchaseRequestDocumentStatusId := DocResponse.DocumentStatusId
	DocumentStatusUrl = config.EnvConfigs.GeneralServiceUrl + "document-status-by-code/99"
	//UomItem := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/" + res.ItemCode + "/P" //strconv.Itoa(response.ItemCode)
	if err := utils.Get(DocumentStatusUrl, &DocResponse, nil); err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed to Fetch Document Status From General Service",
			Err:        err,
		}
	}
	ClosedDocumentId := DocResponse.DocumentStatusId
	if entities.PurchaseRequestDocumentStatusId != DraftPurchaseRequestDocumentStatusId {
		entities.PurchaseRequestDocumentStatusId = ClosedDocumentId
		err = db.Save(&entities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Delete Data Failed",
				Data:       i,
				Err:        err,
			}
		}
		return true, nil
	}
	detailentities := transactionsparepartentities.PurchaseRequestDetail{}
	errDetail := db.Where("purchase_request_system_number = ?", i).Delete(&detailentities).Error
	//errDetail := db.Delete(&detailentities).Where("purchase_request_system_number = ?", i).Error
	if errDetail != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Delete Data Failed",
			Data:       i,
			Err:        errDetail,
		}
	}
	err = db.Delete(&entities).Where(transactionsparepartentities.PurchaseRequestEntities{PurchaseRequestSystemNumber: i}).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Delete Data Failed",
			Data:       i,
			Err:        err,
		}
	}
	return true, nil
}
func (p *PurchaseRequestRepositoryImpl) SubmitPurchaseRequest(db *gorm.DB, request transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, id int) (transactionsparepartpayloads.PurchaseRequestGetByIdResponses, *exceptions.BaseErrorResponse) {
	//fix normalize response
	var count int64
	var res transactionsparepartpayloads.PurchaseRequestGetByIdResponses
	entities := transactionsparepartentities.PurchaseRequestEntities{}
	err := db.Model(&entities).
		Where(transactionsparepartentities.PurchaseRequestEntities{PurchaseRequestSystemNumber: id}).
		First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Data to Updated Is Not Found",
				Data:       res,
				Err:        err,
			}
		}
		return res, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Update Data Failed",
			Data:       res,
			Err:        err,
		}
	}
	//this is logic for getting doc no
	//CEK DOC STATUS ID FOR READY CODE = 20 status ready code = 20
	var DocResponse transactionsparepartpayloads.PurchaseRequestDocumentStatus
	DocumentStatusUrl := config.EnvConfigs.GeneralServiceUrl + "document-status-by-code/20"

	//UomItem := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/" + res.ItemCode + "/P" //strconv.Itoa(response.ItemCode)
	if err := utils.Get(DocumentStatusUrl, &DocResponse, nil); err != nil {
		return res, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed to Fetch Document Status From General Service",
			Err:        err,
		}
	}
	entities.PurchaseRequestDocumentStatusId = DocResponse.DocumentStatusId //status ready
	entities.PurchaseRequestDocumentNumber = "SPPR/N/10/19/11111"
	entities.ChangeNo = entities.ChangeNo + 1
	entities.UpdatedDate = &request.UpdatedDate
	entities.UpdatedByUserId = request.UpdatedByUserId

	err = db.Save(&entities).Error
	if err != nil {
		return res, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed Update Data",
			Data:       res,
			Err:        err,
		}
	}
	err = db.Table("trx_purchase_request_detail pr").
		Joins("inner join trx_work_order_detail w on pr.reference_system_number = w.work_order_status_id and pr.reference_line = w.work_order_operation_item_line").
		Where("pr.purchase_request_system_number = ? and pr.item_quantity <> w.frt_quantity", id).Count(&count).Error
	if err != nil {
		return res, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
			Err:        err,
		}
	}

	result, errs := p.GetByIdPurchaseRequest(db, id)
	return result, errs
}
func (p *PurchaseRequestRepositoryImpl) InsertPurchaseRequestDetail(db *gorm.DB, payloads transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, i int) (transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, *exceptions.BaseErrorResponse) {
	res, err := p.SavePurchaseRequestDetail(db, payloads, i)
	return res, err
}
func (p *PurchaseRequestRepositoryImpl) GetAllItemTypePrRequest(db *gorm.DB, conditions []utils.FilterCondition, page pagination.Pagination, companyid int) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var response []transactionsparepartpayloads.PurchaseRequestItemGetAll
	entities := masteritementities.Item{}
	var PeriodResponse financeservice.OpenPeriodPayloadResponse
	PeriodUrl := config.EnvConfigs.FinanceServiceUrl + "closing-period-company/current-period?company_id" + strconv.Itoa(companyid) + "&closing_module_detail_code=SP" //strconv.Itoa(response.ItemCode)

	//UomItem := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/" + res.ItemCode + "/P" //strconv.Itoa(response.ItemCode)
	if err := utils.Get(PeriodUrl, &PeriodResponse, nil); err != nil {
		return page, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to Period Response data from external service",
			Err:        err,
		}
	}

	year := PeriodResponse.PeriodYear
	month := PeriodResponse.PeriodMonth
	JoinTable := db.Table("mtr_item A").Select("A.item_id,"+
		"A.item_code,"+
		"A.item_name,"+
		"Z.item_class_name,"+
		"A.item_type,"+
		"A.item_level_1,"+
		"A.item_level_2,"+
		"A.item_level_3,"+
		"A.item_level_4,A.unit_of_measurement_type_id,"+
		"ISNULL(SUM(x.quantity_ending-x.quantity_allocated),0) as quantity").
		Joins("LEFT JOIN mtr_item_class Z on A.item_class_id = Z.item_class_id").
		Joins("LEFT JOIN mtr_location_stock x ON A.item_id = x.item_id and x.company_id = ? and period_year =?"+
			" AND period_month = ? AND x.warehouse_id in (select whs.warehouse_id "+
			" from mtr_warehouse_master whs "+
			" where whs.company_id = x.company_id "+
			" AND whs.warehouse_costing_type_id <> 'NON' "+
			" AND whs.warehouse_id = x.warehouse_id) ", companyid, year, month).
		//Joins("INNER JOIN mtr_uom uom ON uom.uom_type_id = A.unit_of_measurement_type_id").
		Group("A.item_id,A.item_code," +
			"A.item_name," +
			"Z.item_class_name," +
			"A.item_type," +
			"A.item_level_1," +
			"A.item_level_2," +
			"A.item_level_3," +
			"A.item_level_4," +
			"A.unit_of_measurement_type_id") //.Order("A.item_id")
	WhereQuery := utils.ApplyFilter(JoinTable, conditions)
	err := WhereQuery.Scopes(pagination.Paginate(&entities, &page, WhereQuery)).Scan(&response).Error
	if err != nil {
		return page, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: err}
	}
	if len(response) == 0 {
		return page, &exceptions.BaseErrorResponse{
			Message:    "Data not found",
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	i := 1
	var result []transactionsparepartpayloads.PurchaseRequestItemGetAll

	for _, res := range response {
		var UomItemResponse transactionsparepartpayloads.UomItemResponses
		UomItem := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/" + strconv.Itoa(res.ItemId) + "/P" //strconv.Itoa(response.ItemCode)
		//UomItem := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/" + res.ItemCode + "/P" //strconv.Itoa(response.ItemCode)
		if err := utils.Get(UomItem, &UomItemResponse, nil); err != nil {
			return page, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Uom Item data from external service",
				Err:        err,
			}
		}
		var UomRate float64
		var QtyRes float64
		if UomItemResponse.SourceConvertion == nil {
			QtyRes = 0
		} else {
			QtyRes = res.Quantity * *UomItemResponse.TargetConvertion

		}
		if UomItemResponse.SourceConvertion == nil {
			return page, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Uom Source Convertion From External Data",
				Err:        err,
			}
		}
		UomRate = QtyRes * *UomItemResponse.SourceConvertion // QtyRes * *UomItemResponse.SourceConvertion
		UomRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", UomRate), 64)

		uomentities := masteritementities.UomItem{}
		err = db.Model(&uomentities).Where(masteritementities.UomItem{ItemId: res.ItemId}).Scan(&uomentities).Error
		res.UnitOfMeasurementCode = uomentities.UomTypeCode
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				uomentities.UomTypeCode = ""
			} else {
				return page, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed Get Uom Code",
					Err:        err,
				}
			}

		}
		res.UnitOfMeasurementRate = UomRate
		res.Sequence = i
		i++
		result = append(result, res)
	}
	page.Rows = result
	return page, nil
}
func (p *PurchaseRequestRepositoryImpl) GetByIdPurchaseRequestItemPr(db *gorm.DB, compid int, i int) (transactionsparepartpayloads.PurchaseRequestItemGetAll, *exceptions.BaseErrorResponse) {
	var response transactionsparepartpayloads.PurchaseRequestItemGetAll
	//entities := masteritementities.Item{}
	var PeriodResponse financeservice.OpenPeriodPayloadResponse
	PeriodUrl := config.EnvConfigs.FinanceServiceUrl + "closing-period-company/current-period?company_id=" + strconv.Itoa(compid) + "&closing_module_detail_code=SP" //strconv.Itoa(response.ItemCode)
	if err := utils.Get(PeriodUrl, &PeriodResponse, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to Period Response data from external service",
			Err:        err,
		}
	}

	year := PeriodResponse.PeriodYear
	month := PeriodResponse.PeriodMonth
	fmt.Println("year = " + year)

	fmt.Println(PeriodResponse)
	err := db.Table("mtr_item A").Select("A.item_id,"+
		"A.item_code,"+
		"A.item_name,"+
		"Z.item_class_name,"+
		"A.item_type,"+
		"A.item_level_1,"+
		"A.item_level_2,"+
		"A.item_level_3,"+
		"A.item_level_4,A.unit_of_measurement_type_id,"+
		"ISNULL(SUM(x.quantity_ending-x.quantity_allocated),0) as quantity").
		Joins("LEFT JOIN mtr_item_class Z on A.item_class_id = Z.item_class_id").
		Joins("LEFT JOIN mtr_location_stock x ON A.item_id = x.item_id and x.company_id = ? and period_year =?"+
			" AND period_month = ? AND x.warehouse_id in (select whs.warehouse_id "+
			" from mtr_warehouse_master whs "+
			" where whs.company_id = x.company_id "+
			" AND whs.warehouse_costing_type_id <> 'NON' "+
			" AND whs.warehouse_id = x.warehouse_id) ", compid, year, month).
		Joins("INNER JOIN mtr_uom uom ON uom.uom_type_id = A.unit_of_measurement_type_id").
		Group("A.item_id,A.item_code,"+
			"A.item_name,A.unit_of_measurement_type_id,"+
			"Z.item_class_name,"+
			"A.item_type,"+
			"A.item_level_1,"+
			"A.item_level_2,"+
			"A.item_level_3,"+
			"A.item_level_4").Where("A.item_id = ?", i).First(&response).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Record Not Found",
			Data:       nil,
			Err:        err,
		}
	}
	var UomItemResponse transactionsparepartpayloads.UomItemResponses
	UomItem := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/" + strconv.Itoa(response.ItemId) + "/P" //strconv.Itoa(response.ItemCode)

	if err := utils.Get(UomItem, &UomItemResponse, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Uom Item data from external service",
			Err:        err,
		}
	}
	var UomRate float64
	var QtyRes float64
	if UomItemResponse.SourceConvertion == nil {
		QtyRes = 0
	} else {
		QtyRes = response.Quantity * *UomItemResponse.TargetConvertion

	}
	if UomItemResponse.SourceConvertion == nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Uom Source Convertion From External Data",
			Err:        err,
		}
	}
	UomRate = QtyRes * *UomItemResponse.SourceConvertion // QtyRes * *UomItemResponse.SourceConvertion
	UomRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", UomRate), 64)
	response.UnitOfMeasurementRate = UomRate
	uomentities := masteritementities.UomItem{}
	err = db.Model(&uomentities).Where(masteritementities.UomItem{ItemId: response.ItemId}).Scan(&uomentities).Error
	response.UnitOfMeasurementCode = uomentities.UomTypeCode
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uomentities.UomTypeCode = ""
		} else {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed Get Uom Code",
				Err:        err,
			}
		}

	}
	response.UnitOfMeasurementCode = uomentities.UomTypeCode
	return response, nil
}

func (p *PurchaseRequestRepositoryImpl) GetByCodePurchaseRequestItemPr(db *gorm.DB, compid int, s string) (transactionsparepartpayloads.PurchaseRequestItemGetAll, *exceptions.BaseErrorResponse) {
	var response transactionsparepartpayloads.PurchaseRequestItemGetAll

	var PeriodResponse financeservice.OpenPeriodPayloadResponse
	PeriodUrl := config.EnvConfigs.FinanceServiceUrl + "closing-period-company/current-period?company_id" + strconv.Itoa(compid) + "&closing_module_detail_code=SP" //strconv.Itoa(response.ItemCode)

	if err := utils.Get(PeriodUrl, &PeriodResponse, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to Period Response data from external service",
			Err:        err,
		}
	}

	year := PeriodResponse.PeriodYear
	month := PeriodResponse.PeriodMonth
	//year := "2012"
	//month := "12"
	err := db.Table("mtr_item A").Select("A.item_id,"+
		"A.item_code,"+
		"A.item_name,"+
		"Z.item_class_name,"+
		"A.item_type,"+
		"A.item_level_1,"+
		"A.item_level_2,"+
		"A.item_level_3,"+
		"A.item_level_4, A.unit_of_measurement_type_id,"+
		"ISNULL(SUM(x.quantity_ending-x.quantity_allocated),0) as quantity").
		Joins("LEFT JOIN mtr_item_class Z on A.item_class_id = Z.item_class_id").
		Joins("LEFT JOIN mtr_location_stock x ON A.item_id = x.item_id and x.company_id = ? and period_year =?"+
			" AND period_month = ? AND x.warehouse_id in (select whs.warehouse_id "+
			" from mtr_warehouse_master whs "+
			" where whs.company_id = x.company_id "+
			" AND whs.warehouse_costing_type_id <> 'NON' "+
			" AND whs.warehouse_id = x.warehouse_id) ", compid, year, month).
		Joins("INNER JOIN mtr_uom uom ON uom.uom_type_id = A.unit_of_measurement_type_id").
		Group("A.item_id,A.item_code,"+
			"A.item_name,A.unit_of_measurement_type_id,"+
			"Z.item_class_name,"+
			"A.item_type,"+
			"A.item_level_1,"+
			"A.item_level_2,"+
			"A.item_level_3,"+
			"A.item_level_4").Where("A.item_code = ?", s).First(&response).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Record Not Found",
			Data:       nil,
			Err:        err,
		}
	}
	var UomItemResponse transactionsparepartpayloads.UomItemResponses
	UomItem := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/" + strconv.Itoa(response.ItemId) + "/P" //strconv.Itoa(response.ItemCode)
	if err := utils.Get(UomItem, &UomItemResponse, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Uom Item data from external service",
			Err:        err,
		}
	}
	var UomRate float64
	var QtyRes float64
	if UomItemResponse.SourceConvertion == nil {
		QtyRes = 0
	} else {
		QtyRes = response.Quantity * *UomItemResponse.TargetConvertion

	}
	if UomItemResponse.SourceConvertion == nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Uom Source Convertion From External Data",
			Err:        err,
		}
	}
	UomRate = QtyRes * *UomItemResponse.SourceConvertion // QtyRes * *UomItemResponse.SourceConvertion
	UomRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", UomRate), 64)
	uomentities := masteritementities.UomItem{}
	err = db.Model(&uomentities).Where(masteritementities.UomItem{ItemId: response.ItemId}).Scan(&uomentities).Error
	response.UnitOfMeasurementCode = uomentities.UomTypeCode
	response.UnitOfMeasurementRate = UomRate
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uomentities.UomTypeCode = ""
		} else {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed Get Uom Code",
				Err:        err,
			}
		}

	}
	response.UnitOfMeasurementCode = uomentities.UomTypeCode
	return response, nil
}

func (p *PurchaseRequestRepositoryImpl) GetAllItemTypePr(db *gorm.DB, payloads transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, i int) (transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, *exceptions.BaseErrorResponse) {
	//resultitem := transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads{}
	//
	//err := db.Table("gmItem0 A").
	//	Select("DISTINCT A.ITEM_CODE AS Code, A.ITEM_NAME AS Description, A.ITEM_CLASS AS ItemClass, "+
	//		"A.ITEM_TYPE AS ItemType, A.ITEM_LVL_1 AS ItemLvl1, A.ITEM_LVL_2 AS ItemLvl2, "+
	//		"A.ITEM_LVL_3 AS ItemLvl3, A.ITEM_LVL_4 AS ItemLvl4, "+
	//		"ISNULL(C.QTY_AVAILABLE, 0) AS AvailQty, D.MOVING_CODE AS MovingCode, "+
	//		"ISNULL(AmIc.QTY_ON_ORDER, 0) AS QtyOnOrder").
	//	Joins("LEFT OUTER JOIN (SELECT DISTINCT ITEM_CODE, VEHICLE_BRAND FROM gmItem1) B ON A.ITEM_CODE = B.ITEM_CODE").
	//	Joins("LEFT JOIN (SELECT F.ITEM_CODE, SUM(F.QTY_AVAILABLE) AS QTY_AVAILABLE "+
	//		"FROM viewLocationStock F "+
	//		"INNER JOIN gmLoc1 G ON F.COMPANY_CODE = G.COMPANY_CODE AND F.WHS_CODE = G.WAREHOUSE_CODE "+
	//		"WHERE F.PERIOD_YEAR = ? AND F.PERIOD_MONTH = ? AND F.COMPANY_CODE = ? AND F.WHS_GROUP = ? "+
	//		"AND G.COSTING_TYPE <> dbo.getVariableValue('HPP_WH_TYPE_NON') "+
	//		"GROUP BY F.ITEM_CODE) AS C ON C.ITEM_CODE = A.ITEM_CODE", year, month, companyCode, varValue2).
	//	Joins("LEFT JOIN (SELECT Z.ITEM_CODE, Z.QTY_ON_ORDER "+
	//		"FROM amItemCycle Z "+
	//		"WHERE Z.PERIOD_YEAR = ? AND Z.PERIOD_MONTH = ? AND Z.COMPANY_CODE = ?) AS AmIc ON AmIc.ITEM_CODE = A.ITEM_CODE", year, month, companyCode).
	//	Joins("LEFT JOIN amMovingCodeItem D ON D.COMPANY_CODE = ? AND D.ITEM_CODE = A.ITEM_CODE AND "+
	//		"D.PROCESS_DATE = (SELECT TOP 1 E.PROCESS_DATE FROM amMovingCodeItem E "+
	//		"WHERE E.COMPANY_CODE = ? AND E.ITEM_CODE = D.ITEM_CODE ORDER BY E.PROCESS_DATE DESC)", companyCode, companyCode).
	//	Find(&items).Error
	//pertanyuaan 1 dbo getvariable value itu belum ada
	//gmLoc1 gada costingril
	return transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads{}, nil
}
func (p *PurchaseRequestRepositoryImpl) VoidPurchaseRequestDetailMultiId(db *gorm.DB, s string) (bool, *exceptions.BaseErrorResponse) {
	ids := strings.Split(s, ",")
	for _, i2 := range ids {
		entities := transactionsparepartentities.PurchaseRequestDetail{}
		converted, _ := strconv.Atoi(i2)

		err := db.Model(&entities).Where(transactionsparepartentities.PurchaseRequestDetail{PurchaseRequestDetailSystemNumber: converted}).First(&entities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
		}
		HeaderEntities := transactionsparepartentities.PurchaseRequestEntities{}
		err = db.Model(HeaderEntities).Where(transactionsparepartentities.PurchaseRequestEntities{PurchaseRequestSystemNumber: entities.PurchaseRequestSystemNumber}).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
		}
		if HeaderEntities.PurchaseRequestDocumentStatusId != 10 {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Message: "Document is Not Draf"}

		}
		err = db.Where(transactionsparepartentities.PurchaseRequestDetail{PurchaseRequestDetailSystemNumber: converted}).Delete(&entities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
		}
	}
	return true, nil
}
