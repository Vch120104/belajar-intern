package transactionsparepartrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	financeserviceapiutils "after-sales/api/utils/finance-service"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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
		RequestByURL := config.EnvConfigs.GeneralServiceUrl + "user-detail/" + strconv.Itoa(res.CreatedByUserId)
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
		OrderType := masterentities.OrderType{}

		err = db.Model(&OrderType).Where(masterentities.OrderType{OrderTypeId: res.OrderTypeId}).
			First(&OrderType).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return paginationResponses, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Err:        err,
					Message:    "order type is not found",
				}
			}
			return paginationResponses, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "failed on getting order type",
			}
		}
		//var OrderType transactionsparepartpayloads.PurchaseRequestOrderTypeResponse
		//OrderTypeURL := config.EnvConfigs.GeneralServiceUrl + "order-type/" + strconv.Itoa(res.OrderTypeId)
		//if err := utils.Get(OrderTypeURL, &OrderType, nil); err != nil {
		//	return paginationResponses, &exceptions.BaseErrorResponse{
		//		StatusCode: http.StatusInternalServerError,
		//		Message:    "Failed to fetch customer data from external service",
		//		Err:        err,
		//	}
		//}
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
	var response transactionsparepartpayloads.PurchaseRequestGetByIdResponses
	var entities transactionsparepartentities.PurchaseRequestEntities

	// Fetch the purchase request record by ID
	if err := db.Model(&entities).
		Where(transactionsparepartentities.PurchaseRequestEntities{PurchaseRequestSystemNumber: i}).
		First(&response).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	//get company name
	CompanyReponse, CompanyReponseerr := generalserviceapiutils.GetCompanyDataById(response.CompanyId)
	if CompanyReponseerr != nil {
		return response, CompanyReponseerr
	}
	ItemGroup, ItemGroupErr := generalserviceapiutils.GetItemGroupById(response.ItemGroupId)
	if ItemGroupErr != nil {
		return response, ItemGroupErr
	}

	OrderType := masterentities.OrderType{}

	err := db.Model(&OrderType).Where(masterentities.OrderType{OrderTypeId: response.OrderTypeId}).First(&OrderType).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get order type",
			Err:        err,
		}
	}
	var purchaseRequestStatusDesc transactionsparepartpayloads.PurchaseRequestStatusResponse
	//StatusURL, StatusURLErr := generalserviceapiutils.GetDocumentStatusById(response.PurchaseRequestDocumentStatusId)
	//if StatusURLErr != nil {
	//	return response, StatusURLErr
	//}

	GetBrandName, GetBrandNameErr := salesserviceapiutils.GetUnitBrandById(response.PurchaseRequestDocumentStatusId)
	if GetBrandNameErr != nil {
		return response, GetBrandNameErr
	}

	GetDivision, GetDivisionErr := generalserviceapiutils.GetDivisionById(response.DivisionId)
	if GetDivisionErr != nil {
		return response, GetDivisionErr
	}

	//var GetCostCenterName transactionsparepartpayloads.CostCenterResponses
	CostCenter, CostCenterErr := generalserviceapiutils.GetCostCenterById(response.CostCenterId)
	if CostCenterErr != nil {
		return response, CostCenterErr
	}

	//var GetCcyName transactionsparepartpayloads.CurrencyCodeResponse
	GetCcyName, GetCcyNameErr := financeserviceapiutils.GetCurrencyId(response.CurrencyId)
	if GetCcyNameErr != nil {
		return response, GetCcyNameErr
	}

	ProfitCenter, ProfitCenterErr := generalserviceapiutils.GetProfitCenterById(response.ProfitCenterId)
	if ProfitCenterErr != nil {
		return response, ProfitCenterErr
	}
	var WarehouseGroupName masterwarehouseentities.WarehouseGroup

	err = db.Model(&WarehouseGroupName).Where(masterwarehouseentities.WarehouseGroup{WarehouseGroupId: response.WarehouseGroupId}).
		First(&WarehouseGroupName).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "warehouse group is not found please check input",
				Err:        err,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "warehouse group fetch error please check input",
			Err:        err,
		}
	}
	//var WarehouseGroupName transactionsparepartpayloads.WarehouseGroupResponses
	//WarehouseGroupURL := config.EnvConfigs.AfterSalesServiceUrl + "warehouse-group/" + strconv.Itoa(response.WarehouseGroupId)
	//if err := utils.Get(WarehouseGroupURL, &WarehouseGroupName, nil); err != nil {
	//	return response, &exceptions.BaseErrorResponse{
	//		StatusCode: http.StatusInternalServerError,
	//		Message:    "Failed to fetch Warehouse data from external service",
	//		Err:        err,
	//	}
	//}

	//var GetWarehouseResponsesName transactionsparepartpayloads.WarehouseResponses
	GetWarehouseResponses := masterwarehouseentities.WarehouseMaster{}
	err = db.Model(&GetWarehouseResponses).Where(masterwarehouseentities.WarehouseMaster{WarehouseId: response.WarehouseId}).
		First(&GetWarehouseResponses).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "warehouse master is not found please check inputt",
				Err:        err,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "warehouse master fetch error please check input",
			Err:        err,
		}
	}
	//GetWarehouseResponsesName := masterwarehouseentities.WarehouseMaster{}
	//err := db.Model(&GetWarehouseResponsesName).Where(masterwarehouseentities.WarehouseMaster{WarehouseId: response.WarehouseId}).
	//	First(&GetWarehouseResponsesName).Error
	//if err != nil {
	//	return response, &exceptions.BaseErrorResponse{
	//		StatusCode: http.StatusInternalServerError,
	//		Message:    "Failed to fetch Warehouse data from external service",
	//		Err:        err,
	//	}
	//}

	//WarehouseURL := config.EnvConfigs.AfterSalesServiceUrl + "warehouse-master/" + strconv.Itoa(response.WarehouseId)
	//if err := utils.Get(WarehouseURL, &GetWarehouseResponsesName, nil); err != nil {
	//	return response, &exceptions.BaseErrorResponse{
	//		StatusCode: http.StatusInternalServerError,
	//		Message:    "Failed to fetch Warehouse data from external service",
	//		Err:        err,
	//	}
	//}
	var RequestBy transactionsparepartpayloads.PurchaseRequestRequestedByResponse
	RequestByURL := config.EnvConfigs.GeneralServiceUrl + "user-detail/" + strconv.Itoa(response.CreatedByUserId)
	if err := utils.Get(RequestByURL, &RequestBy, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Requested By data from external service",
			Err:        err,
		}
	}

	//var UpdatedBy transactionsparepartpayloads.PurchaseRequestRequestedByResponse
	UpdatedBy, UpdatedByerr := generalserviceapiutils.GetUserDetailsByID(response.CreatedByUserId)
	if UpdatedByerr != nil {
		return response, UpdatedByerr
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

	result := transactionsparepartpayloads.PurchaseRequestGetByIdNormalizeResponses{
		Company:                       CompanyReponse.CompanyName,
		PurchaseRequestSystemNumber:   response.PurchaseRequestSystemNumber,
		PurchaseRequestDocumentNumber: response.PurchaseRequestDocumentNumber,
		PurchaseRequestDocumentDate:   response.PurchaseRequestDocumentDate,
		PurchaseRequestDocumentStatus: purchaseRequestStatusDesc.PurchaseRequestStatusDescription,
		ItemGroup:                     ItemGroup.ItemGroupName,
		Brand:                         GetBrandName.BrandName,
		ReferenceType:                 PurchaseRequestReferenceType.ReferenceTypePurchaseRequestName,
		//ReferenceDocumentNumber:       docNo,
		ReferenceDocumentNumber: response.ReferenceDocumentNumber,

		OrderType:                  OrderType.OrderTypeName,
		BudgetCode:                 response.BudgetCode,
		ProjectNo:                  response.ProjectNo,
		Division:                   GetDivision.DivisionName,
		PurchaseRequestRemark:      response.PurchaseRequestRemark,
		PurchaseRequestTotalAmount: response.PurchaseRequestTotalAmount,
		ExpectedArrivalDate:        response.ExpectedArrivalDate,
		ExpectedArrivalTime:        response.ExpectedArrivalTime,
		CostCenter:                 CostCenter.CostCenterName,
		ProfitCenter:               ProfitCenter.ProfitCenterName,
		WarehouseGroup:             WarehouseGroupName.WarehouseGroupName,
		Warehouse:                  GetWarehouseResponses.WarehouseName,
		SetOrder:                   response.SetOrder,
		Currency:                   GetCcyName.CurrencyName,
		ChangeNo:                   0,
		CreatedByUser:              RequestBy.UserEmployeeName,
		CreatedDate:                response.CreatedDate,
		UpdatedByUser:              UpdatedBy.EmployeeName,
		UpdatedDate:                response.UpdatedDate,
	}
	fmt.Println(result)
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
	UomItemResponse := masteritementities.UomItem{}
	err = db.Model(&UomItemResponse).Where(masteritementities.UomItem{ItemId: ItemResponse.ItemId, UomSourceTypeCode: "P"}).
		First(&UomItemResponse).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "uom not found",
				Err:        err,
			}
		}
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "uom item error fetch",
			Err:        err,
		}
	}
	//UomItem := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/" + strconv.Itoa(ItemResponse.ItemId) + "/P" //strconv.Itoa(response.ItemCode)
	//var UomItemResponse transactionsparepartpayloads.UomItemResponses
	//
	//if err := utils.Get(UomItem, &UomItemResponse, nil); err != nil {
	//	return result, &exceptions.BaseErrorResponse{
	//		StatusCode: http.StatusInternalServerError,
	//		Message:    "Failed to fetch Uom Item data from external service",
	//		Err:        err,
	//	}
	//}
	var UomRate float64
	var QtyRes float64
	if UomItemResponse.SourceConvertion == 0 {
		QtyRes = 0
	} else {
		QtyRes = *response.ItemQuantity * UomItemResponse.TargetConvertion

	}
	if UomItemResponse.SourceConvertion == 0 {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Uom Source Convertion From External Data",
			Err:        err,
		}
	}
	UomRate = QtyRes * UomItemResponse.SourceConvertion // QtyRes * *UomItemResponse.SourceConvertion
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
	err := db.Model(&entities).
		Where(transactionsparepartentities.PurchaseRequestEntities{PurchaseRequestSystemNumber: payloads.PurchaseRequestSystemNumber}).
		First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "failed to get purchase request hedaer data",
			}
		}
		return Response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	LineNumber := 1
	//default value
	LineStatus := "10"
	Price := 0.000
	//validation
	TotalPrice := Price * *payloads.ItemQuantity
	PRDetailEntities := transactionsparepartentities.PurchaseRequestDetail{
		//PurchaseRequestSystemNumberDetail: 0,
		ItemId:                      payloads.ItemId,
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
	//get item class
	itemMaster := masteritementities.Item{}
	err = db.Model(&itemMaster).Where(masteritementities.Item{ItemId: payloads.ItemId}).
		First(&itemMaster).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "failed to get item id",
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed on get item id",
		}
	}
	//entities.item
	entities.ItemQuantity = payloads.ItemQuantity
	entities.ItemId = payloads.ItemId
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
	//var DocResponse transactionsparepartpayloads.PurchaseRequestDocumentStatus
	//DocumentStatusUrl := config.EnvConfigs.GeneralServiceUrl + "document-status-by-code/20"
	DocResponseReady, DocResponsErr := generalserviceapiutils.GetDocumentStatusByCode("20")
	if DocResponsErr != nil {
		return res, DocResponsErr
	}
	//UomItem := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/" + res.ItemCode + "/P" //strconv.Itoa(response.ItemCode)
	//if err := utils.Get(DocumentStatusUrl, &DocResponse, nil); err != nil {
	//	return res, &exceptions.BaseErrorResponse{
	//		StatusCode: http.StatusBadRequest,
	//		Message:    "Failed to Fetch Document Status From General Service",
	//		Err:        err,
	//	}
	//}
	entities.PurchaseRequestDocumentStatusId = DocResponseReady.DocumentStatusId //status ready
	docNo, errDocNo := p.GenerateDocumentNumber(db, 10)
	if errDocNo != nil {
		return res, errDocNo
	}
	entities.PurchaseRequestDocumentNumber = docNo
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
	PeriodResponse, periodErr := financeserviceapiutils.GetOpenPeriodByCompany(companyid, "SP")
	if periodErr != nil {
		return page, periodErr
	}

	year := PeriodResponse.PeriodYear
	month := PeriodResponse.PeriodMonth
	JoinTable := db.Table("mtr_item A").Select("A.item_id,"+
		"A.item_code,"+
		"A.item_name,"+
		"A.item_name,"+
		"Z.item_class_name,"+
		"A.item_type_id,"+
		"A.item_level_1_id,"+
		"A.item_level_2_id,"+
		"A.item_level_3_id,"+
		"L1.item_level_1_code as item_level_1,"+
		"L2.item_level_2_code as item_level_2,"+
		"L3.item_level_3_code as item_level_3,"+
		"L4.item_level_4_code as item_level_4,"+
		"IT.item_type_code,"+
		"A.item_level_4_id,A.unit_of_measurement_type_id,"+
		"ISNULL(SUM(x.quantity_ending-x.quantity_allocated),0) as quantity").
		Joins("LEFT JOIN mtr_item_class Z on A.item_class_id = Z.item_class_id").
		Joins(`LEFT JOIN mtr_item_level_1 L1 ON A.item_level_1_id = L1.item_level_1_id AND A.item_level_1_id <> 0`).
		Joins(`LEFT JOIN mtr_item_level_2 L2 ON A.item_level_2_id = L2.item_level_2_id AND A.item_level_2_id <> 0`).
		Joins(`LEFT JOIN mtr_item_level_3 L3 ON A.item_level_3_id = L3.item_level_3_id AND A.item_level_3_id <> 0`).
		Joins(`LEFT JOIN mtr_item_level_4 L4 ON A.item_level_4_id = L4.item_level_4_id AND A.item_level_4_id <> 0`).
		Joins("LEFT JOIN mtr_item_type IT ON IT.item_type_id = a.item_type_id").
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
			"A.item_type_id," +
			"A.item_level_1_id," +
			"A.item_level_2_id," +
			"A.item_level_3_id," +
			"A.item_level_4_id," +
			"L1.item_level_1_code," +
			"L2.item_level_2_code," +
			"L3.item_level_3_code," +
			"L4.item_level_4_code," +
			"A.unit_of_measurement_type_id," +
			"IT.item_type_code")
	//.Order("A.item_id")
	WhereQuery := utils.ApplyFilter(JoinTable, conditions)
	err := WhereQuery.Scopes(pagination.Paginate(&entities, &page, WhereQuery)).Scan(&response).Error
	if len(response) == 0 {
		page.Rows = []string{}
		return page, nil
	}
	if err != nil {
		return page, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: err}
	}

	i := 1

	var result []transactionsparepartpayloads.PurchaseRequestItemGetAll

	for _, res := range response {
		//get item type
		//itemTypeEntities := masteritementities.ItemType{}
		//err = db.Model(&itemTypeEntities).Where(masteritementities.ItemType{ItemTypeId: res.ItemTypeId}).
		//	First(&itemTypeEntities).Error
		//if err != nil {
		//	if errors.Is(err, gorm.ErrRecordNotFound) {
		//		return page, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Err: err, Message: "item type is not found"}
		//	}
		//	return page, &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Err: err, Message: "item type error please check input"}
		//}
		UomItemResponse := masteritementities.UomItem{}
		err = db.Model(&UomItemResponse).Where(masteritementities.UomItem{ItemId: res.ItemId, UomSourceTypeCode: "P"}).
			First(&UomItemResponse).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return page, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Err: err, Message: "uom item is not found"}
			}
			return page, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "failed to get uom item",
			}
		}
		//if res.ItemLevel1Id != 0 {
		//	itemLevel1Entities := masteritementities.ItemLevel1{}
		//	err = db.Model(&itemLevel1Entities).Where(masteritementities.ItemLevel1{ItemLevel1Id: res.ItemLevel1Id}).
		//		First(&itemLevel1Entities).Error
		//	if err != nil {
		//		if errors.Is(err, gorm.ErrRecordNotFound) {
		//			return page, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Err: err, Message: "item level 1 is not found"}
		//		}
		//		return page, &exceptions.BaseErrorResponse{
		//			StatusCode: http.StatusInternalServerError,
		//			Err:        err,
		//			Message:    "Failed to get item level id 3",
		//		}
		//	}
		//	res.ItemLevel1 = itemLevel1Entities.ItemLevel1Code
		//}
		//if res.ItemLevel2Id != 0 {
		//	itemLevel2Entities := masteritementities.ItemLevel2{}
		//	err = db.Model(&itemLevel2Entities).Where(masteritementities.ItemLevel2{ItemLevel2Id: res.ItemLevel2Id}).
		//		First(&itemLevel2Entities).Error
		//	if err != nil {
		//		if errors.Is(err, gorm.ErrRecordNotFound) {
		//			return page, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Err: err, Message: "item level 2 is not found"}
		//		}
		//		return page, &exceptions.BaseErrorResponse{
		//			StatusCode: http.StatusInternalServerError,
		//			Err:        err,
		//			Message:    "Failed to get item level id 3",
		//		}
		//	}
		//	res.ItemLevel2 = itemLevel2Entities.ItemLevel2Code
		//}
		//if res.ItemLevel3Id != 0 {
		//	itemLevel3Entities := masteritementities.ItemLevel3{}
		//	err = db.Model(&itemLevel3Entities).Where(masteritementities.ItemLevel3{ItemLevel3Id: res.ItemLevel3Id}).
		//		First(&itemLevel3Entities).Error
		//	if err != nil {
		//		if errors.Is(err, gorm.ErrRecordNotFound) {
		//			return page, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Err: err, Message: "item level 3 is not found"}
		//		}
		//		return page, &exceptions.BaseErrorResponse{
		//			StatusCode: http.StatusInternalServerError,
		//			Err:        err,
		//			Message:    "Failed to get item level id 3",
		//		}
		//	}
		//	res.ItemLevel3 = itemLevel3Entities.ItemLevel3Code
		//}
		//if res.ItemLevel4Id != 0 {
		//	itemLevel4Entities := masteritementities.ItemLevel4{}
		//	err = db.Model(&itemLevel4Entities).Where(masteritementities.ItemLevel4{ItemLevel4Id: res.ItemLevel4Id}).
		//		First(&itemLevel4Entities).Error
		//	if err != nil {
		//		if errors.Is(err, gorm.ErrRecordNotFound) {
		//			return page, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Err: err, Message: "item level 4 is not found"}
		//		}
		//	}
		//	res.ItemLevel1 = itemLevel4Entities.ItemLevel4Code
		//}

		//var UomItemResponse transactionsparepartpayloads.UomItemResponses
		//UomItem := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/" + strconv.Itoa(res.ItemId) + "/P" //strconv.Itoa(response.ItemCode)
		////UomItem := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/" + res.ItemCode + "/P" //strconv.Itoa(response.ItemCode)
		//if err := utils.Get(UomItem, &UomItemResponse, nil); err != nil {
		//	return page, &exceptions.BaseErrorResponse{
		//		StatusCode: http.StatusInternalServerError,
		//		Message:    "Failed to fetch Uom Item data from external service",
		//		Err:        err,
		//	}
		//}
		var UomRate float64
		var QtyRes float64
		if UomItemResponse.SourceConvertion == 0 {
			QtyRes = 0
		} else {
			QtyRes = res.Quantity * UomItemResponse.TargetConvertion

		}
		if UomItemResponse.SourceConvertion == 0 {
			return page, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Uom Source Conversion From External Data",
				Err:        err,
			}
		}
		UomRate = QtyRes * UomItemResponse.SourceConvertion // QtyRes * *UomItemResponse.SourceConvertion
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
		//res.ItemTypeCode = itemTypeEntities.ItemTypeCode
		i++
		result = append(result, res)
	}
	page.Rows = result
	return page, nil
}
func (p *PurchaseRequestRepositoryImpl) GetByIdPurchaseRequestItemPr(db *gorm.DB, compid int, i int) (transactionsparepartpayloads.PurchaseRequestItemGetAll, *exceptions.BaseErrorResponse) {
	var response transactionsparepartpayloads.PurchaseRequestItemGetAll
	//entities := masteritementities.Item{}
	PeriodResponse, periodErr := financeserviceapiutils.GetOpenPeriodByCompany(compid, "SP")
	if periodErr != nil {
		return response, nil
	}
	year := PeriodResponse.PeriodYear
	month := PeriodResponse.PeriodMonth
	fmt.Println("year = " + year)

	fmt.Println(PeriodResponse)
	err := db.Table("mtr_item A").Select("A.item_id,"+
		"A.item_code,"+
		"A.item_name,"+
		"A.item_name,"+
		"Z.item_class_name,"+
		"A.item_type_id,"+
		"A.item_level_1_id,"+
		"A.item_level_2_id,"+
		"A.item_level_3_id,"+
		"L1.item_level_1_code as item_level_1,"+
		"L2.item_level_2_code as item_level_2,"+
		"L3.item_level_3_code as item_level_3,"+
		"L4.item_level_4_code as item_level_4,"+
		"IT.item_type_code,"+
		"A.item_level_4_id,A.unit_of_measurement_type_id,"+
		"ISNULL(SUM(x.quantity_ending-x.quantity_allocated),0) as quantity").
		Joins("LEFT JOIN mtr_item_class Z on A.item_class_id = Z.item_class_id").
		Joins(`LEFT JOIN mtr_item_level_1 L1 ON A.item_level_1_id = L1.item_level_1_id AND A.item_level_1_id <> 0`).
		Joins(`LEFT JOIN mtr_item_level_2 L2 ON A.item_level_2_id = L2.item_level_2_id AND A.item_level_2_id <> 0`).
		Joins(`LEFT JOIN mtr_item_level_3 L3 ON A.item_level_3_id = L3.item_level_3_id AND A.item_level_3_id <> 0`).
		Joins(`LEFT JOIN mtr_item_level_4 L4 ON A.item_level_4_id = L4.item_level_4_id AND A.item_level_4_id <> 0`).
		Joins("LEFT JOIN mtr_item_type IT ON IT.item_type_id = a.item_type_id").
		Joins("LEFT JOIN mtr_location_stock x ON A.item_id = x.item_id and x.company_id = ? and period_year =?"+
			" AND period_month = ? AND x.warehouse_id in (select whs.warehouse_id "+
			" from mtr_warehouse_master whs "+
			" where whs.company_id = x.company_id "+
			" AND whs.warehouse_costing_type_id <> 'NON' "+
			" AND whs.warehouse_id = x.warehouse_id) ", compid, year, month).
		//Joins("INNER JOIN mtr_uom uom ON uom.uom_type_id = A.unit_of_measurement_type_id").
		Group("A.item_id,A.item_code,"+
			"A.item_name,"+
			"Z.item_class_name,"+
			"A.item_type_id,"+
			"A.item_level_1_id,"+
			"A.item_level_2_id,"+
			"A.item_level_3_id,"+
			"A.item_level_4_id,"+
			"L1.item_level_1_code,"+
			"L2.item_level_2_code,"+
			"L3.item_level_3_code,"+
			"L4.item_level_4_code,"+
			"A.unit_of_measurement_type_id,"+
			"IT.item_type_code").Where("A.item_id = ?", i).First(&response).Error
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
	// var QtyRes float64
	// if UomItemResponse.SourceConvertion == nil {
	// 	QtyRes = 0
	// } else {
	// 	QtyRes = response.Quantity * *UomItemResponse.TargetConvertion

	// }
	// if UomItemResponse.SourceConvertion == nil {
	// 	return response, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Message:    "Failed to fetch Uom Source Convertion From External Data",
	// 		Err:        err,
	// 	}
	// }

	// var UomRate float64
	// UomRate = QtyRes * *UomItemResponse.SourceConvertion // QtyRes * *UomItemResponse.SourceConvertion
	// UomRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", UomRate), 64)
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

	PeriodResponse, periodErr := financeserviceapiutils.GetOpenPeriodByCompany(compid, "SP")
	if periodErr != nil {
		return response, periodErr
	}

	year := PeriodResponse.PeriodYear
	month := PeriodResponse.PeriodMonth
	//year := "2012"
	//month := "12"
	err := db.Table("mtr_item A").Select("A.item_id,"+
		"A.item_code,"+
		"A.item_name,"+
		"A.item_name,"+
		"Z.item_class_name,"+
		"A.item_type_id,"+
		"A.item_level_1_id,"+
		"A.item_level_2_id,"+
		"A.item_level_3_id,"+
		"L1.item_level_1_code as item_level_1,"+
		"L2.item_level_2_code as item_level_2,"+
		"L3.item_level_3_code as item_level_3,"+
		"L4.item_level_4_code as item_level_4,"+
		"IT.item_type_code,"+
		"A.item_level_4_id,A.unit_of_measurement_type_id,"+
		"ISNULL(SUM(x.quantity_ending-x.quantity_allocated),0) as quantity").
		Joins("LEFT JOIN mtr_item_class Z on A.item_class_id = Z.item_class_id").
		Joins(`LEFT JOIN mtr_item_level_1 L1 ON A.item_level_1_id = L1.item_level_1_id AND A.item_level_1_id <> 0`).
		Joins(`LEFT JOIN mtr_item_level_2 L2 ON A.item_level_2_id = L2.item_level_2_id AND A.item_level_2_id <> 0`).
		Joins(`LEFT JOIN mtr_item_level_3 L3 ON A.item_level_3_id = L3.item_level_3_id AND A.item_level_3_id <> 0`).
		Joins(`LEFT JOIN mtr_item_level_4 L4 ON A.item_level_4_id = L4.item_level_4_id AND A.item_level_4_id <> 0`).
		Joins("LEFT JOIN mtr_item_type IT ON IT.item_type_id = a.item_type_id").
		Joins("LEFT JOIN mtr_location_stock x ON A.item_id = x.item_id and x.company_id = ? and period_year =?"+
			" AND period_month = ? AND x.warehouse_id in (select whs.warehouse_id "+
			" from mtr_warehouse_master whs "+
			" where whs.company_id = x.company_id "+
			" AND whs.warehouse_costing_type_id <> 'NON' "+
			" AND whs.warehouse_id = x.warehouse_id) ", compid, year, month).
		//Joins("INNER JOIN mtr_uom uom ON uom.uom_type_id = A.unit_of_measurement_type_id").
		Group("A.item_id,A.item_code,"+
			"A.item_name,"+
			"Z.item_class_name,"+
			"A.item_type_id,"+
			"A.item_level_1_id,"+
			"A.item_level_2_id,"+
			"A.item_level_3_id,"+
			"A.item_level_4_id,"+
			"L1.item_level_1_code,"+
			"L2.item_level_2_code,"+
			"L3.item_level_3_code,"+
			"L4.item_level_4_code,"+
			"A.unit_of_measurement_type_id,"+
			"IT.item_type_code").Where("A.item_code = ?", s).First(&response).Error
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
			Message:    "Failed to fetch Uom Source Conversion From External Data",
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
		converted, errs := strconv.Atoi(i2)
		if errs != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed on parse id please check input",
			}
		}

		err := db.Model(&entities).Where(transactionsparepartentities.PurchaseRequestDetail{PurchaseRequestDetailSystemNumber: converted}).First(&entities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
		}
		HeaderEntities := transactionsparepartentities.PurchaseRequestEntities{}
		err = db.Model(HeaderEntities).Where(transactionsparepartentities.PurchaseRequestEntities{PurchaseRequestSystemNumber: entities.PurchaseRequestSystemNumber}).
			First(&HeaderEntities).
			Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
		}
		StatusDocDraft, errStatusCode := generalserviceapiutils.GetDocumentStatusByCode("10")
		if errStatusCode != nil {
			return false, errStatusCode
		}
		if HeaderEntities.PurchaseRequestDocumentStatusId != StatusDocDraft.DocumentStatusId {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Message: "Document is Not Draft"}

		}
		err = db.Where(transactionsparepartentities.PurchaseRequestDetail{PurchaseRequestDetailSystemNumber: converted}).Delete(&entities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
		}
	}
	return true, nil
}

func (p *PurchaseRequestRepositoryImpl) GenerateDocumentNumber(tx *gorm.DB, id int) (string, *exceptions.BaseErrorResponse) {
	var workOrder transactionsparepartentities.PurchaseRequestEntities

	// Get the work order based on the work order system number
	err := tx.Model(&transactionsparepartentities.PurchaseRequestEntities{}).Where("purchase_request_system_number = ?", id).First(&workOrder).Error
	if err != nil {

		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve work order from the database: %v", err)}
	}

	if workOrder.BrandId == 0 {

		return "", &exceptions.BaseErrorResponse{Message: "brand_id is missing in the work order. Please ensure the work order has a valid brand_id before generating document number."}
	}

	// Get the last work order based on the work order system number
	var lastWorkOrder transactionsparepartentities.PurchaseRequestEntities
	err = tx.Model(&transactionsparepartentities.PurchaseRequestEntities{}).
		Where("brand_id = ?", workOrder.BrandId).
		Order("purchase_request_document_number desc").
		First(&lastWorkOrder).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {

		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve last work order: %v", err)}
	}

	currentTime := time.Now()
	month := int(currentTime.Month())
	year := currentTime.Year() % 100 // Use last two digits of the year

	// fetch data brand from external api
	brandResponse, brandErr := generalserviceapiutils.GetBrandGenerateDoc(workOrder.BrandId)
	if brandErr != nil {
		return "", &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch brand data from external service",
			Err:        brandErr.Err,
		}
	}

	// Check if BrandCode is not empty before using it
	if brandResponse.BrandCode == "" {
		return "", &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Message: "Brand code is empty"}
	}

	// Get the initial of the brand code
	brandInitial := brandResponse.BrandCode[0]

	// Handle the case when there is no last work order or the format is invalid
	newDocumentNumber := fmt.Sprintf("SPPR/%c/%02d/%02d/00001", brandInitial, month, year)
	if lastWorkOrder.PurchaseRequestSystemNumber != 0 {
		lastWorkOrderDate := lastWorkOrder.PurchaseRequestDocumentDate
		lastWorkOrderYear := lastWorkOrderDate.Year() % 100

		// Check if the last work order is from the same year
		if lastWorkOrderYear == year {
			lastWorkOrderCode := lastWorkOrder.PurchaseRequestDocumentNumber
			codeParts := strings.Split(lastWorkOrderCode, "/")
			if len(codeParts) == 5 {
				lastWorkOrderNumber, err := strconv.Atoi(codeParts[4])
				if err == nil {
					newWorkOrderNumber := lastWorkOrderNumber + 1
					newDocumentNumber = fmt.Sprintf("SPPR/%c/%02d/%02d/%05d", brandInitial, month, year, newWorkOrderNumber)
				} else {
					log.Printf("Failed to parse last work order code: %v", err)
				}
			} else {
				log.Println("Invalid last work order code format")
			}
		}
	}

	log.Printf("New document number: %s", newDocumentNumber)
	return newDocumentNumber, nil
}

//das
