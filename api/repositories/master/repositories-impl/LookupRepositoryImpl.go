package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	masteroperationentities "after-sales/api/entities/master/operation"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	financeserviceapiutils "after-sales/api/utils/finance-service"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type LookupRepositoryImpl struct {
}

func StartLookupRepositoryImpl() masterrepository.LookupRepository {
	return &LookupRepositoryImpl{}
}

// LocationGoodsReceipt implements masterrepository.LookupRepository.
func (*LookupRepositoryImpl) LocationItemGoodsReceive(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	// var locationItem masteritementities.ItemLocation
	var response []masterpayloads.LocationItemGoodsReceiveResponse

	query := tx.Table("mtr_location_item as it").
		Select(
			"distinct it.warehouse_location_id",
			"loc.warehouse_location_code",
			"loc.warehouse_location_name",
		).
		Joins("LEFT JOIN mtr_warehouse_location loc on loc.warehouse_location_id = it.warehouse_location_id").
		Joins("Inner join mtr_warehouse_master war on war.warehouse_id = it.warehouse_id").Order("it.warehouse_location_id")

	whereQ := utils.ApplyFilter(query, filterCondition)
	paginatedQuery := whereQ.Scopes(pagination.Paginate(&pages, whereQ))

	err := paginatedQuery.Scan(&response).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching lookup data 'ItemMasterForFreeAccs'",
			Err:        err,
		}
	}

	pages.Rows = response
	return pages, nil

}

// dbo.getOprItemDisc
// get DISCOUNT value base on line type in operation or item master
func (r *LookupRepositoryImpl) GetOprItemDisc(tx *gorm.DB, linetypeStr string, billCodeId int, oprItemCode int, agreementId int, profitCenterId int, minValue float64, companyId int, brandId int, contractServSysNo int, whsGroup int, orderTypeId int) (float64, *exceptions.BaseErrorResponse) {
	var discount float64
	var discCode string
	var itemTypeId int
	var companyCodePrice int
	var useDiscDecentralize string
	var hpp float64
	var outerMargin float64
	var ccy string
	var pricelist float64
	var total float64
	var agreementCount int64

	discount = 0

	if orderTypeId == 0 {
		orderTypeId = utils.EstWoOrderTypeId // Default value E
	}

	// Get Company Code for Price List
	var commonPricelist bool
	err := tx.Model(&masteritementities.Item{}).
		Where("item_id = ?", oprItemCode).
		Pluck("common_pricelist", &commonPricelist).
		Error
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get common price list",
			Err:        err,
		}
	}

	if commonPricelist {
		companyCodePrice = 0
	} else {
		companyCodePrice = companyId
	}

	// Handle BILLCODE_NOCHARGE
	if billCodeId == utils.TrxTypeWoNoCharge.ID {
		return 100, nil
	}

	// Handle BILLCODE_EXTERNAL OR INSURANCE
	if billCodeId == utils.TrxTypeWoExternal.ID ||
		billCodeId == utils.TrxTypeWoInsurance.ID ||
		billCodeId == utils.TrxTypeSoChannel.ID ||
		billCodeId == utils.TrxTypeSoDirect.ID ||
		billCodeId == utils.TrxTypeSoGSO.ID ||
		billCodeId == utils.TrxTypeWoWarranty.ID ||
		billCodeId == utils.TrxTypeWoFreeService.ID {

		if agreementId != 0 {
			// Retrieve Discount Code
			err = tx.Model(&masteritementities.Item{}).
				Where("item_id = ?", oprItemCode).
				Pluck("discount_code", &discCode).Error
			if err != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to get discount code",
					Err:        err,
				}
			}

			// Check Agreement Validity
			err = tx.Model(&masterentities.Agreement{}).
				Where("aggreement_id = ? AND profit_center_id = ? AND (GETDATE() BETWEEN agreement_date_from AND agreement_date_to)", agreementId, profitCenterId).
				Count(&agreementCount).Error
			if err != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to check agreement validity",
					Err:        err,
				}
			}

			if agreementCount > 0 {
				if linetypeStr != utils.LinetypeOperation && linetypeStr != utils.LinetypePackage {
					if discount == 0 {
						// Check Agreement2
						err = tx.Model(&masterentities.AgreementDiscountGroupDetail{}).
							Where("agreement_id = ? AND order_type_id = ? AND agreement_discount_markup_id = ? AND agreement_selection_id = ?", agreementId, orderTypeId, discCode, utils.EstWoDiscSelectionId).
							Pluck("agreement_discount", &discount).Error
						if err != nil {
							return 0, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to get agreement discount",
								Err:        err,
							}
						}
					}

					if discount == 0 {
						// fetch linetype id from line type code
						linetypecheck, linetypeErr := generalserviceapiutils.GetLineTypeByCode(linetypeStr)
						if linetypeErr != nil {
							return 0, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to get line type",
								Err:        linetypeErr,
							}
						}

						// Check Agreement3
						err = tx.Model(&masterentities.AgreementItemDetail{}).
							Where("agreement_id = ? AND line_type_id = ? AND agreement_item_operation_id = ? AND min_value <= ?", agreementId, linetypecheck.LineTypeId, oprItemCode, minValue).
							Order("min_value DESC").
							Limit(1).
							Pluck("discount_percent", &discount).Error
						if err != nil {
							return 0, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to get agreement item detail discount percent",
								Err:        err,
							}
						}
					}

					if discount == 0 {

						// fetch linetype id from line type code
						linetypecheck1, linetypeErr := generalserviceapiutils.GetLineTypeByCode(linetypeStr)
						if linetypeErr != nil {
							return 0, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to get line type",
								Err:        linetypeErr,
							}
						}

						// Check Agreement1
						err = tx.Model(&masterentities.AgreementDiscount{}).
							Where("agreement_id = ? AND line_type_id = ? AND min_value <= ?", agreementId, linetypecheck1.LineTypeId, minValue).
							Order("min_value DESC").
							Limit(1).
							Pluck("discount_percent", &discount).Error
						if err != nil {
							return 0, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to get agreement discount",
								Err:        err,
							}
						}
					}
				}

				if linetypeStr == utils.LinetypeOperation || linetypeStr == utils.LinetypePackage {
					if discount == 0 {
						// fetch linetype id from line type code
						linetypecheck2, linetypeErr := generalserviceapiutils.GetLineTypeByCode(linetypeStr)
						if linetypeErr != nil {
							return 0, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to get line type",
								Err:        linetypeErr,
							}
						}

						// Check Agreement3 for Operations
						err = tx.Model(&masterentities.AgreementItemDetail{}).
							Where("agreement_id = ? AND line_type_id = ? AND agreement_item_operation_id = ? AND min_value <= ?", agreementId, linetypecheck2.LineTypeId, oprItemCode, minValue).
							Order("min_value DESC").
							Limit(1).
							Pluck("discount_percent", &discount).Error
						if err != nil {
							return 0, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to get agreement item detail discount percent",
								Err:        err,
							}
						}
					}

					if discount == 0 {
						// fetch linetype id from line type code
						linetypecheck2, linetypeErr := generalserviceapiutils.GetLineTypeByCode(linetypeStr)
						if linetypeErr != nil {
							return 0, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to get line type",
								Err:        linetypeErr,
							}
						}

						// Check Agreement1 for Operations
						err = tx.Model(&masterentities.AgreementDiscount{}).
							Where("agreement_id = ? AND line_type_id = ? AND min_value <= ?", agreementId, linetypecheck2.LineTypeId, minValue).
							Order("min_value DESC").
							Limit(1).
							Pluck("discount_percent", &discount).Error
						if err != nil {
							return 0, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to get agreement discount",
								Err:        err,
							}
						}
					}
				}
			}
		}
	}

	// Handle BILLCODE_DECENTRALIZE
	if billCodeId == utils.TrxTypeWoDeCentralize.ID ||
		billCodeId == utils.TrxTypeSoDeCentralize.ID {

		if linetypeStr != utils.LinetypeOperation && linetypeStr != utils.LinetypePackage {

			// Get Use Disc Decentralize and Item Type
			tx.Model(&masteritementities.Item{}).
				Where("item_id = ?", oprItemCode).
				Select("use_disc_decentralize, item_type_id").
				Row().Scan(&useDiscDecentralize, &itemTypeId)

			if useDiscDecentralize == "N" {
				discount = 0
			} else {
				// Calculate HPP
				tx.Model(&masterentities.GroupStock{}).
					Where("period_year = YEAR(GETDATE()) AND period_month = RIGHT('0' + CAST(MONTH(GETDATE()) AS VARCHAR), 2) AND company_code = ? AND item_code = ? AND whs_group = ?", companyId, oprItemCode, whsGroup).
					Select("CASE ISNULL(price_current, 0) WHEN 0 THEN price_begin ELSE price_current END AS hpp").
					Row().Scan(&hpp)

				// Get Outer Margin
				err = tx.Table("dms_microservices_general_dev.dbo.mtr_company_reference").
					Where("company_code = ?", companyId).
					Pluck("margin_outer_kpp", &outerMargin).Error
				if err != nil {
					return 0, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to get outer margin",
						Err:        err,
					}
				}

				// Get Currency
				err = tx.Table("dms_microservices_general_dev.dbo.mtr_company_reference").
					Where("company_code = ?", companyId).
					Pluck("currency_id", &ccy).Error
				if err != nil {
					return 0, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to get currency",
						Err:        err,
					}
				}

				// Get Price List
				err = tx.Table("dms_microservices_sales_dev.dbo.mtr_item_price_list").
					Where("is_active = 1 AND brand_id = ? AND effective_date <= GETDATE() AND item_code = ? AND currency_id = ? AND company_id = ?", brandId, oprItemCode, ccy, companyCodePrice).
					Order("effective_date DESC").
					Limit(1).
					Pluck("price_list_amount", &pricelist).Error

				if err != nil {
					return 0, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to get price list amount",
						Err:        err,
					}
				}

				total = hpp + (hpp * outerMargin / 100)
				if pricelist != 0 {
					discount = 100 - ((total * 100) / pricelist)
					if discount < 0 {
						discount = 0
					}
				} else {
					discount = 0
				}
			}

			// Check ItemType Services
			if itemTypeId == 2 {
				discount = 0
			}
		}
	}

	// Handle BILLCODE_CONTRACT_SERVICE
	if billCodeId == utils.TrxTypeWoContractService.ID {

		// fetch linetype id from line type code
		linetypecheck3, linetypeErr := generalserviceapiutils.GetLineTypeByCode(linetypeStr)
		if linetypeErr != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get line type",
				Err:        linetypeErr,
			}
		}

		err = tx.Model(&transactionworkshopentities.ContractService{}).
			Joins("INNER JOIN trx_contract_service_Operation_detail ON trx_contract_service_Operation_detail.contract_service_system_number = trx_contract_service.contract_service_system_number").
			Where("trx_contract_service.contract_service_system_number = ? AND trx_contract_service_Operation_detail.line_type_id = ? AND trx_contract_service_Operation_detail.operation_id = ?", contractServSysNo, linetypecheck3.LineTypeId, oprItemCode).
			Pluck("trx_contract_service_Operation_detail.operation_discount_percent", &discount).Error

		if err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get contract service operation detail discount percent",
				Err:        err,
			}
		}
	}

	return discount, nil
}

// dbo.getOprItemPrice
// get price value base on line type in operation or item master
func (r *LookupRepositoryImpl) GetOprItemPrice(tx *gorm.DB, linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, jobTypeId int, variantId int, currencyId int, billCode int, whsGroup string) (float64, *exceptions.BaseErrorResponse) {
	var (
		price               float64
		effDate             = time.Now()
		markupPercentage    float64
		companyCodePrice    int
		commonPriceList     bool
		defaultPriceCodeId  int
		useDiscDecentralize string
		priceCount          int64
		priceCodeId         int
	)

	// Set markup percentage based on company ID
	markupPercentage = 0
	if companyId == 139 {
		markupPercentage = 11.00
	}

	if err := tx.Model(&masteritementities.Item{}).
		Where("item_id = ?", oprItemCode).
		Select("ISNULL(common_pricelist, ?)", false).
		Scan(&commonPriceList).Error; err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get common price list",
			Err:        err,
		}
	}

	// Set company code price based on common price list
	if commonPriceList {
		companyCodePrice = 0
	} else {
		companyCodePrice = companyId
	}

	switch linetypeId {
	case 1:
		// Package price logic
		if err := tx.Model(&masterentities.PackageMaster{}).
			Where("package_id = ?", oprItemCode).
			Select("package_price").
			Scan(&price).Error; err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get package price",
				Err:        err,
			}
		}

	case 2:
		// Operation price logic
		query := tx.Model(&masteroperationentities.LabourSellingPriceDetail{}).
			Joins("JOIN mtr_labour_selling_price ON mtr_labour_selling_price.labour_selling_price_id = mtr_labour_selling_price_detail.labour_selling_price_id").
			Where("mtr_labour_selling_price.brand_id = ? AND mtr_labour_selling_price.effective_date <= ? AND mtr_labour_selling_price.job_type_id = ? AND mtr_labour_selling_price_detail.model_id = ? AND mtr_labour_selling_price.company_id = ?",
				brandId, effDate, jobTypeId, modelId, companyId)

		if variantId == 0 {
			query = query.Where("mtr_labour_selling_price_detail.variant_id = 0")
		} else {
			query = query.Where("mtr_labour_selling_price_detail.variant_id = ? OR mtr_labour_selling_price_detail.variant_id = 0", variantId)
		}

		if err := query.Order("mtr_labour_selling_price.effective_date DESC").Limit(1).Pluck("selling_price", &price).Error; err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get selling price",
				Err:        err,
			}
		}

	default:
		defaultPriceCodeId = 1
		priceCodeId = 1
		if err := tx.Model(&masteritementities.ItemPriceList{}).
			Where("is_active = 1 AND brand_id = ? AND effective_date <= ? AND item_id = ? AND currency_id = ? AND company_id = ? AND price_list_code_id = ?",
				brandId, effDate, oprItemCode, currencyId, companyCodePrice, priceCodeId).Count(&priceCount).Error; err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to check price list existence",
				Err:        err,
			}
		}

		if priceCount == 0 {
			priceCodeId = defaultPriceCodeId
		}

		// Handling based on bill code
		if billCode == utils.TrxTypeWoNoCharge.ID || billCode == utils.TrxTypeWoCentralize.ID || billCode == utils.TrxTypeWoInternal.ID || billCode == utils.TrxTypeSoCentralize.ID || billCode == utils.TrxTypeSoInternal.ID || billCode == utils.TrxTypeSoExport.ID {
			var periodYear, periodMonth string

			month := effDate.Format("01")
			year := effDate.Format("2006")

			// Get MODULE_SP
			moduleSP := "SP"

			currentPeriodPayloads, err := financeserviceapiutils.GetOpenPeriodByCompany(companyId, moduleSP)
			if err != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        errors.New("failed to get period details"),
				}
			}

			// Add additional validation for period month and period year
			if !(currentPeriodPayloads.PeriodMonth <= month && currentPeriodPayloads.PeriodYear <= year) {
				currentPeriodPayloads.PeriodMonth = ""
				currentPeriodPayloads.PeriodYear = ""
			}

			if currentPeriodPayloads.PeriodYear != "" {
				periodYear = currentPeriodPayloads.PeriodYear
			} else {
				periodYear = "0000"
			}

			if currentPeriodPayloads.PeriodMonth != "" {
				periodMonth = currentPeriodPayloads.PeriodMonth
			} else {
				periodMonth = "00"
			}

			// Check item type
			itemTypeExists := false
			if err := tx.Model(&masteritementities.Item{}).
				Where("item_id = ? AND item_type_id = ?", oprItemCode, 2).
				Select("item_type_id").
				Scan(&itemTypeExists).Error; err != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to check item type",
					Err:        err,
				}
			}

			if itemTypeExists {
				// Get price from gmPriceList for items
				if err := tx.Model(&masteritementities.ItemPriceList{}).
					Where("is_active = 1 AND brand_id = ? AND effective_date <= ? AND item_code = ? AND currency_id = ? AND company_id = ? AND price_list_code_id = ?",
						brandId, effDate, oprItemCode, currencyId, companyCodePrice, priceCodeId).
					Order("effective_date DESC").
					Limit(1).
					Pluck("price_list_amount", &price).Error; err != nil {
					return 0, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to get price amount from gmPriceList",
						Err:        err,
					}
				}
			} else {
				// Get price from amGroupStock for other items
				if err := tx.Model(&masterentities.GroupStock{}).
					Where("period_year = ? AND period_month = ? AND item_id = ? AND company_id = ? AND whs_group = ?",
						periodYear, periodMonth, oprItemCode, companyId, whsGroup).
					Select("CASE ISNULL(price_current, 0) WHEN 0 THEN price_begin ELSE price_current END AS hpp").
					Pluck("hpp", &price).Error; err != nil {
					return 0, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to get group stock price",
						Err:        err,
					}
				}
			}

			if billCode != utils.TrxTypeWoInternal.ID && billCode != utils.TrxTypeSoExport.ID && billCode != utils.TrxTypeSoInternal.ID {
				if err := tx.Model(&masteritementities.Item{}).
					Where("item_id = ?", oprItemCode).
					Pluck("ISNULL(use_disc_decentralize, '')", &useDiscDecentralize).Error; err != nil {
					return 0, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to get useDiscDecentralize value",
						Err:        err,
					}
				}

				if useDiscDecentralize == "" || useDiscDecentralize == " " {
					useDiscDecentralize = "Y"
				}

				if useDiscDecentralize == "N" {
					if err := tx.Model(&masteritementities.ItemPriceList{}).
						Where("is_active = 1 AND brand_id = ? AND effective_date <= ? AND item_id = ? AND currency_id = ? AND company_id = ? AND price_list_code_id = ?",
							brandId, effDate, oprItemCode, currencyId, companyId, defaultPriceCodeId).
						Pluck("price_list_amount", &price).Error; err != nil {
						return 0, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to get default price list amount",
							Err:        err,
						}
					}
				}
			}
		} else {
			if err := tx.Model(&masteritementities.ItemPriceList{}).
				Where("is_active = 1 AND brand_id = ? AND effective_date <= ? AND item_id = ? AND currency_id = ? AND company_id = ? AND price_list_code_id = ?",
					brandId, effDate, oprItemCode, currencyId, companyCodePrice, priceCodeId).
				Order("effective_date DESC").
				Limit(1).
				Pluck("price_list_amount", &price).Error; err != nil {
				return 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to get price amount",
					Err:        err,
				}
			}
		}
	}

	// Apply markup percentage if applicable
	if linetypeId == 2 && billCode == utils.TrxTypeWoInternal.ID {
		price += price * markupPercentage / 100
	}

	return price, nil
}

// usp_comLookUp
// IF @strEntity = 'ItemOprCode'--OPERATION MASTER & ITEM MASTER
func (r *LookupRepositoryImpl) ItemOprCode(tx *gorm.DB, linetypeId int, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var (
		companyCode = 473
		currentTime = time.Now()
		year, month = currentTime.Year(), int(currentTime.Month() - 1)
	)

	// Fetch item type from external service
	var itemTypeFetchGoods masteritementities.ItemType
	if err := tx.Where("item_type_code = ?", "G").First(&itemTypeFetchGoods).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Item type not found",
				Err:        fmt.Errorf("item type with code %s not found", "G"),
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item type code",
			Err:        err,
		}
	}

	var itemTypeFetchServices masteritementities.ItemType
	if err := tx.Where("item_type_code = ?", "S").First(&itemTypeFetchServices).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Item type not found",
				Err:        fmt.Errorf("item type with code %s not found", "S"),
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item type code",
			Err:        err,
		}
	}

	// Base Query
	baseQuery := tx.Session(&gorm.Session{NewDB: true})

	switch linetypeId {
	case 1:
		baseQuery = baseQuery.Table("mtr_package A").
			Select("A.package_id, A.package_code, A.package_name, "+
				"COALESCE(SUM(mtr_package_master_detail.frt_quantity), 0) AS frt, "+
				"B.profit_center_name, C.model_code, C.model_description, A.package_price, "+
				"A.model_id, A.brand_id, A.variant_id").
			Joins("INNER JOIN mtr_package_master_detail ON A.package_id = mtr_package_master_detail.package_id").
			Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_profit_center B ON A.profit_center_id = B.profit_center_id").
			Joins("INNER JOIN dms_microservices_sales_dev.dbo.mtr_unit_model C ON A.model_id = C.model_id").
			Where("A.is_active = ?", 1).
			Group("A.package_id, A.package_code, A.package_name, B.profit_center_name, " +
				"C.model_code, C.model_description, A.package_price, A.model_id, A.brand_id, A.variant_id").
			Order("A.package_id")

	case 2:
		baseQuery = baseQuery.Table("mtr_operation_model_mapping AS omm").
			Select("omm.operation_id AS operation_id, "+
				"oc.operation_code AS operation_code, oc.operation_name AS operation_name, "+
				"MAX(ofrt.frt_hour) AS frt_hour, "+
				"oe.operation_entries_code AS operation_entries_code, oe.operation_entries_description AS operation_entries_description, "+
				"ok.operation_key_code AS operation_key_code, ok.operation_key_description AS operation_key_description").
			Joins("INNER JOIN mtr_operation_frt AS ofrt ON omm.operation_model_mapping_id = ofrt.operation_model_mapping_id").
			Joins("LEFT OUTER JOIN mtr_operation_code AS oc ON omm.operation_id = oc.operation_id").
			Joins("LEFT OUTER JOIN mtr_operation_entries AS oe ON oc.operation_entries_id = oe.operation_entries_id").
			Joins("LEFT OUTER JOIN mtr_operation_key AS ok ON oc.operation_key_id = ok.operation_key_id").
			Where("omm.is_active = ?", true).
			Group("omm.operation_id, oc.operation_code, oc.operation_name, " +
				"oe.operation_entries_code, oe.operation_entries_description, " +
				"ok.operation_key_code, ok.operation_key_description").
			Order("omm.operation_id")

	case 3:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}

		// fetch item class from external service
		var itemClassResp masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SP").First(&itemClassResp).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SP"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		}
		// "SP"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name,
				B.brand_id AS brand_id,
				B.model_id AS model_id,
				B.variant_id AS variant_id,
				COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
				A.item_level_1_id AS item_level_1,
				mil1.item_level_1_code AS item_level_1_code, 
				A.item_level_2_id AS item_level_2,
				mil2.item_level_2_code AS item_level_2_code, 
				A.item_level_3_id AS item_level_3,
				mil3.item_level_3_code AS item_level_3_code, 
				A.item_level_4_id AS item_level_4,
				mil4.item_level_4_code AS item_level_4_code
			`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?",
				itemGrpFetch.ItemGroupId,
				itemTypeFetchGoods.ItemTypeId,
				itemClassResp.ItemClassId,
				true).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")
	case 4:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassOL masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "OL").First(&itemClassOL).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "OL"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "OL"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name,
				B.brand_id AS brand_id,
				B.model_id AS model_id,
				B.variant_id AS variant_id,
				COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
				A.item_level_1_id AS item_level_1,
				mil1.item_level_1_code AS item_level_1_code, 
				A.item_level_2_id AS item_level_2,
				mil2.item_level_2_code AS item_level_2_code, 
				A.item_level_3_id AS item_level_3,
				mil3.item_level_3_code AS item_level_3_code, 
				A.item_level_4_id AS item_level_4,
				mil4.item_level_4_code AS item_level_4_code
						`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassOL.ItemClassId, true).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 5:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassMT masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "MT").First(&itemClassMT).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "MT"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "MT"
		// fetch item class from external service
		var itemClassSB masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SB").First(&itemClassSB).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SB"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "SB"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name,
				B.brand_id AS brand_id,
				B.model_id AS model_id,
				B.variant_id AS variant_id,
				COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
				A.item_level_1_id AS item_level_1,
				mil1.item_level_1_code AS item_level_1_code, 
				A.item_level_2_id AS item_level_2,
				mil2.item_level_2_code AS item_level_2_code, 
				A.item_level_3_id AS item_level_3,
				mil3.item_level_3_code AS item_level_3_code, 
				A.item_level_4_id AS item_level_4,
				mil4.item_level_4_code AS item_level_4_code
						`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type_id = ? AND (A.item_class_id = ? OR A.item_class_id = ?) AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassMT.ItemClassId, itemClassSB.ItemClassId, true).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 6:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassWF masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "WF").First(&itemClassWF).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "WF"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "WF"
		// Fetch item group from external service
		var itemGrpOJFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "OJ").First(&itemGrpOJFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "OJ"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		} // "OJ"

		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name,
					B.brand_id AS brand_id,
					B.model_id AS model_id,
					B.variant_id AS variant_id,
					COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
						`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("(A.item_group_id = ? OR A.item_group_id = ?) AND A.item_class_id = ? AND A.item_type_id = ? AND A.is_active = ?", itemGrpOJFetch.ItemGroupId, itemGrpFetch.ItemGroupId, itemClassWF.ItemClassId, itemTypeFetchServices.ItemTypeId, true).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 7:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassAC masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "AC").First(&itemClassAC).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "AC"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "AC"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
						A.item_id AS item_id, 
						A.item_code AS item_code, 
						A.item_name AS item_name,
						B.brand_id AS brand_id,
						B.model_id AS model_id,
						B.variant_id AS variant_id,
						COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
						A.item_level_1_id AS item_level_1,
						mil1.item_level_1_code AS item_level_1_code, 
						A.item_level_2_id AS item_level_2,
						mil2.item_level_2_code AS item_level_2_code, 
						A.item_level_3_id AS item_level_3,
						mil3.item_level_3_code AS item_level_3_code, 
						A.item_level_4_id AS item_level_4,
						mil4.item_level_4_code AS item_level_4_code
						`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", itemClassAC.ItemClassId, itemGrpFetch.ItemGroupId, true).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 8:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassCM masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "CM").First(&itemClassCM).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "CM"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "CM"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name,
					B.brand_id AS brand_id,
					B.model_id AS model_id,
					B.variant_id AS variant_id,
					COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
							`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ?  AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassCM.ItemClassId, true).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 9:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassSV masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SV").First(&itemClassSV).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SV"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "SV"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name,
				B.brand_id AS brand_id,
				B.model_id AS model_id,
				B.variant_id AS variant_id,
				COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
				A.item_level_1_id AS item_level_1,
				mil1.item_level_1_code AS item_level_1_code, 
				A.item_level_2_id AS item_level_2,
				mil2.item_level_2_code AS item_level_2_code, 
				A.item_level_3_id AS item_level_3,
				mil3.item_level_3_code AS item_level_3_code, 
				A.item_level_4_id AS item_level_4,
				mil4.item_level_4_code AS item_level_4_code
			`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", itemClassSV.ItemClassId, itemGrpFetch.ItemGroupId, true).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	default:
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid line type",
			Err:        errors.New("invalid line type"),
		}
	}

	// Apply manual filters based on linetype
	for _, filter := range filters {
		if linetypeId == 1 {
			switch filter.ColumnField {
			case "package_id":
				baseQuery = baseQuery.Where("A.package_id = ?", filter.ColumnValue)
			case "package_code":
				baseQuery = baseQuery.Where("A.package_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "package_name":
				baseQuery = baseQuery.Where("A.package_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "profit_center_name":
				baseQuery = baseQuery.Where("B.profit_center_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_code":
				baseQuery = baseQuery.Where("C.model_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_description":
				baseQuery = baseQuery.Where("C.model_description LIKE ?", "%"+filter.ColumnValue+"%")
			case "package_price":
				baseQuery = baseQuery.Where("A.package_price = ?", filter.ColumnValue)
			case "model_id":
				baseQuery = baseQuery.Where("A.model_id = ?", filter.ColumnValue)
			case "brand_id":
				baseQuery = baseQuery.Where("A.brand_id = ?", filter.ColumnValue)
			case "variant_id":
				baseQuery = baseQuery.Where("A.variant_id = ?", filter.ColumnValue)
			}
		} else if linetypeId == 2 {
			switch filter.ColumnField {
			case "operation_id":
				baseQuery = baseQuery.Where("oc.operation_id = ?", filter.ColumnValue)
			case "operation_code":
				baseQuery = baseQuery.Where("oc.operation_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_name":
				baseQuery = baseQuery.Where("oc.operation_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "frt_hour":
				baseQuery = baseQuery.Where("ofrt.frt_hour = ?", filter.ColumnValue)
			case "operation_entries_code":
				baseQuery = baseQuery.Where("oe.operation_entries_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_entries_description":
				baseQuery = baseQuery.Where("oe.operation_entries_description LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_key_code":
				baseQuery = baseQuery.Where("ok.operation_key_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_key_description":
				baseQuery = baseQuery.Where("ok.operation_key_description LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_id":
				baseQuery = baseQuery.Where("omm.model_id = ?", filter.ColumnValue)
			case "brand_id":
				baseQuery = baseQuery.Where("omm.brand_id = ?", filter.ColumnValue)
			case "variant_id":
				baseQuery = baseQuery.Where("ofrt.variant_id = ?", filter.ColumnValue)
			}

		} else if linetypeId == 3 || linetypeId == 4 || linetypeId == 5 || linetypeId == 6 || linetypeId == 7 || linetypeId == 8 || linetypeId == 9 {
			switch filter.ColumnField {
			case "item_id":
				baseQuery = baseQuery.Where("A.item_id = ?", filter.ColumnValue)
			case "item_code":
				baseQuery = baseQuery.Where("A.item_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_name":
				baseQuery = baseQuery.Where("A.item_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "available_qty":
				baseQuery = baseQuery.Where("available_qty = ?", filter.ColumnValue)
			case "item_level_1_code":
				baseQuery = baseQuery.Where("mil1.item_level_1_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_2_code":
				baseQuery = baseQuery.Where("mil2.item_level_2_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_3_code":
				baseQuery = baseQuery.Where("mil3.item_level_3_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_4_code":
				baseQuery = baseQuery.Where("mil4.item_level_4_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_id":
				baseQuery = baseQuery.Where("B.model_id = ?", filter.ColumnValue)
			case "brand_id":
				baseQuery = baseQuery.Where("B.brand_id = ?", filter.ColumnValue)
			case "variant_id":
				baseQuery = baseQuery.Where("B.variant_id = ?", filter.ColumnValue)
			}
		}
	}

	// Calculate total rows
	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total rows",
			Err:        err,
		}
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalRows) / float64(paginate.Limit)))

	// Apply pagination
	paginateFunc := pagination.Paginate(&paginate, baseQuery)
	baseQuery = baseQuery.Scopes(paginateFunc)

	// Fetch results
	results := []map[string]interface{}{}
	if err := baseQuery.Find(&results).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch data",
			Err:        err,
		}
	}

	// Set pagination details
	paginate.TotalRows = totalRows
	paginate.TotalPages = totalPages
	paginate.Rows = results

	return paginate, nil
}

// usp_comLookUp
// IF @strEntity = 'ItemOprCode'--OPERATION MASTER & ITEM MASTER
func (r *LookupRepositoryImpl) ItemOprCodeByCode(tx *gorm.DB, linetypeId int, oprItemCode string, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var (
		companyCode = 473
		currentTime = time.Now()
		year, month = currentTime.Year(), int(currentTime.Month() - 1)
	)

	// Fetch item type from external service
	var itemTypeFetchGoods masteritementities.ItemType
	if err := tx.Where("item_type_code = ?", "G").First(&itemTypeFetchGoods).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Item type not found",
				Err:        fmt.Errorf("item type with code %s not found", "G"),
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item type code",
			Err:        err,
		}
	}

	var itemTypeFetchServices masteritementities.ItemType
	if err := tx.Where("item_type_code = ?", "S").First(&itemTypeFetchServices).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Item type not found",
				Err:        fmt.Errorf("item type with code %s not found", "S"),
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item type code",
			Err:        err,
		}
	}

	// Base Query
	baseQuery := tx.Session(&gorm.Session{NewDB: true})

	switch linetypeId {
	case 1:
		baseQuery = baseQuery.Table("mtr_package A").
			Select("A.package_id, A.package_code, A.package_name, "+
				"COALESCE(SUM(mtr_package_master_detail.frt_quantity), 0) AS frt, "+
				"B.profit_center_name, C.model_code, C.model_description, A.package_price, "+
				"A.model_id, A.brand_id, A.variant_id").
			Joins("INNER JOIN mtr_package_master_detail ON A.package_id = mtr_package_master_detail.package_id").
			Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_profit_center B ON A.profit_center_id = B.profit_center_id").
			Joins("INNER JOIN dms_microservices_sales_dev.dbo.mtr_unit_model C ON A.model_id = C.model_id").
			Where("A.is_active = ?", 1).
			Where("A.package_code = ?", oprItemCode).
			Group("A.package_id, A.package_code, A.package_name, B.profit_center_name, " +
				"C.model_code, C.model_description, A.package_price, A.model_id, A.brand_id, A.variant_id").
			Order("A.package_id")

	case 2:
		baseQuery = baseQuery.Table("mtr_operation_model_mapping AS omm").
			Select("omm.operation_id AS operation_id, "+
				"oc.operation_code AS operation_code, oc.operation_name AS operation_name, "+
				"MAX(ofrt.frt_hour) AS frt_hour, "+
				"oe.operation_entries_code AS operation_entries_code, oe.operation_entries_description AS operation_entries_description, "+
				"ok.operation_key_code AS operation_key_code, ok.operation_key_description AS operation_key_description").
			Joins("INNER JOIN mtr_operation_frt AS ofrt ON omm.operation_model_mapping_id = ofrt.operation_model_mapping_id").
			Joins("LEFT OUTER JOIN mtr_operation_code AS oc ON omm.operation_id = oc.operation_id").
			Joins("LEFT OUTER JOIN mtr_operation_entries AS oe ON oc.operation_entries_id = oe.operation_entries_id").
			Joins("LEFT OUTER JOIN mtr_operation_key AS ok ON oc.operation_key_id = ok.operation_key_id").
			Where("omm.is_active = ?", true).
			Where("oc.operation_code = ?", oprItemCode).
			Group("omm.operation_id, oc.operation_code, oc.operation_name, " +
				"oe.operation_entries_code, oe.operation_entries_description, " +
				"ok.operation_key_code, ok.operation_key_description").
			Order("omm.operation_id")

	case 3:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassResp masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SP").First(&itemClassResp).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SP"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "SP"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
		A.item_id AS item_id, 
		A.item_code AS item_code, 
		A.item_name AS item_name,
		B.brand_id AS brand_id,
		B.model_id AS model_id,
		B.variant_id AS variant_id,
		COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
		A.item_level_1_id AS item_level_1,
		mil1.item_level_1_code AS item_level_1_code, 
		A.item_level_2_id AS item_level_2,
		mil2.item_level_2_code AS item_level_2_code, 
		A.item_level_3_id AS item_level_3,
		mil3.item_level_3_code AS item_level_3_code, 
		A.item_level_4_id AS item_level_4,
		mil4.item_level_4_code AS item_level_4_code
	`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?",
				itemGrpFetch.ItemGroupId,
				itemTypeFetchGoods.ItemTypeId,
				itemClassResp.ItemClassId,
				true).
			Where("A.item_code = ?", oprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 4:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassOL masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "OL").First(&itemClassOL).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "OL"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "OL"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
		A.item_id AS item_id, 
		A.item_code AS item_code, 
		A.item_name AS item_name,
		B.brand_id AS brand_id,
		B.model_id AS model_id,
		B.variant_id AS variant_id,
		COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
		A.item_level_1_id AS item_level_1,
		mil1.item_level_1_code AS item_level_1_code, 
		A.item_level_2_id AS item_level_2,
		mil2.item_level_2_code AS item_level_2_code, 
		A.item_level_3_id AS item_level_3,
		mil3.item_level_3_code AS item_level_3_code, 
		A.item_level_4_id AS item_level_4,
		mil4.item_level_4_code AS item_level_4_code
				`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassOL.ItemClassId, true).
			Where("A.item_code = ?", oprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 5:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassMT masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "MT").First(&itemClassMT).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "MT"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "MT"
		// fetch item class from external service
		var itemClassSB masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SB").First(&itemClassSB).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SB"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "SB"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name,
				B.brand_id AS brand_id,
				B.model_id AS model_id,
				B.variant_id AS variant_id,
				COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
				A.item_level_1_id AS item_level_1,
				mil1.item_level_1_code AS item_level_1_code, 
				A.item_level_2_id AS item_level_2,
				mil2.item_level_2_code AS item_level_2_code, 
				A.item_level_3_id AS item_level_3,
				mil3.item_level_3_code AS item_level_3_code, 
				A.item_level_4_id AS item_level_4,
				mil4.item_level_4_code AS item_level_4_code
						`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type_id = ? AND (A.item_class_id = ? OR A.item_class_id = ?) AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassMT.ItemClassId, itemClassSB.ItemClassId, true).
			Where("A.item_code = ?", oprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 6:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassWF masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "WF").First(&itemClassWF).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "WF"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "WF"
		// Fetch item group from external service
		var itemGrpOJFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "OJ").First(&itemGrpOJFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "OJ"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		} // "OJ"

		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name,
					B.brand_id AS brand_id,
					B.model_id AS model_id,
					B.variant_id AS variant_id,
					COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
						`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("(A.item_group_id = ? OR A.item_group_id = ?) AND A.item_class_id = ? AND A.item_type_id = ? AND A.is_active = ?", itemGrpOJFetch.ItemGroupId, itemGrpFetch.ItemGroupId, itemClassWF.ItemClassId, itemTypeFetchServices.ItemTypeId, true).
			Where("A.item_code = ?", oprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 7:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassAC masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "AC").First(&itemClassAC).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "AC"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "AC"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
						A.item_id AS item_id, 
						A.item_code AS item_code, 
						A.item_name AS item_name,
						B.brand_id AS brand_id,
						B.model_id AS model_id,
						B.variant_id AS variant_id,
						COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
						A.item_level_1_id AS item_level_1,
						mil1.item_level_1_code AS item_level_1_code, 
						A.item_level_2_id AS item_level_2,
						mil2.item_level_2_code AS item_level_2_code, 
						A.item_level_3_id AS item_level_3,
						mil3.item_level_3_code AS item_level_3_code, 
						A.item_level_4_id AS item_level_4,
						mil4.item_level_4_code AS item_level_4_code
						`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", itemClassAC.ItemClassId, itemGrpFetch.ItemGroupId, true).
			Where("A.item_code = ?", oprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 8:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassCM masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "CM").First(&itemClassCM).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "CM"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "CM"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name,
					B.brand_id AS brand_id,
					B.model_id AS model_id,
					B.variant_id AS variant_id,
					COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
							`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ?  AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassCM.ItemClassId, true).
			Where("A.item_code = ?", oprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 9:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassSV masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SV").First(&itemClassSV).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SV"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "SV"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name,
				B.brand_id AS brand_id,
				B.model_id AS model_id,
				B.variant_id AS variant_id,
				COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
				A.item_level_1_id AS item_level_1,
				mil1.item_level_1_code AS item_level_1_code, 
				A.item_level_2_id AS item_level_2,
				mil2.item_level_2_code AS item_level_2_code, 
				A.item_level_3_id AS item_level_3,
				mil3.item_level_3_code AS item_level_3_code, 
				A.item_level_4_id AS item_level_4,
				mil4.item_level_4_code AS item_level_4_code
			`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", itemClassSV.ItemClassId, itemGrpFetch.ItemGroupId, true).
			Where("A.item_code = ?", oprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	default:
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid line type",
			Err:        errors.New("invalid line type"),
		}
	}

	// apply filter manual baseon linetype
	for _, filter := range filters {

		if linetypeId == 1 {
			switch filter.ColumnField {
			case "package_id":
				baseQuery = baseQuery.Where("A.package_id = ?", filter.ColumnValue)
			case "package_code":
				baseQuery = baseQuery.Where("A.package_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "package_name":
				baseQuery = baseQuery.Where("A.package_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "profit_center_name":
				baseQuery = baseQuery.Where("B.profit_center_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_code":
				baseQuery = baseQuery.Where("C.model_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_description":
				baseQuery = baseQuery.Where("C.model_description LIKE ?", "%"+filter.ColumnValue+"%")
			case "package_price":
				baseQuery = baseQuery.Where("A.package_price = ?", filter.ColumnValue)
			}
		} else if linetypeId == 2 {
			switch filter.ColumnField {
			case "operation_id":
				baseQuery = baseQuery.Where("oc.operation_id = ?", filter.ColumnValue)
			case "operation_code":
				baseQuery = baseQuery.Where("oc.operation_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_name":
				baseQuery = baseQuery.Where("oc.operation_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "frt_hour":
				baseQuery = baseQuery.Where("ofrt.frt_hour = ?", filter.ColumnValue)
			case "operation_entries_code":
				baseQuery = baseQuery.Where("oe.operation_entries_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_entries_description":
				baseQuery = baseQuery.Where("oe.operation_entries_description LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_key_code":
				baseQuery = baseQuery.Where("ok.operation_key_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_key_description":
				baseQuery = baseQuery.Where("ok.operation_key_description LIKE ?", "%"+filter.ColumnValue+"%")
			}
		} else if linetypeId == 3 || linetypeId == 4 || linetypeId == 5 || linetypeId == 6 || linetypeId == 7 || linetypeId == 8 || linetypeId == 9 {
			switch filter.ColumnField {
			case "item_id":
				baseQuery = baseQuery.Where("A.item_id = ?", filter.ColumnValue)
			case "item_code":
				baseQuery = baseQuery.Where("A.item_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_name":
				baseQuery = baseQuery.Where("A.item_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "available_qty":
				baseQuery = baseQuery.Where("available_qty = ?", filter.ColumnValue)
			case "item_level_1_code":
				baseQuery = baseQuery.Where("mil1.item_level_1_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_2_code":
				baseQuery = baseQuery.Where("mil2.item_level_2_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_3_code":
				baseQuery = baseQuery.Where("mil3.item_level_3_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_4_code":
				baseQuery = baseQuery.Where("mil4.item_level_4_code LIKE ?", "%"+filter.ColumnValue+"%")
			}
		}
	}

	paginateFunc := pagination.Paginate(&paginate, baseQuery)
	baseQuery = baseQuery.Scopes(paginateFunc)

	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total rows",
			Err:        err,
		}
	}

	results := map[string]interface{}{}
	paginateFunc = pagination.Paginate(&paginate, baseQuery)
	if err := baseQuery.Scopes(paginateFunc).Find(&results).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch data",
			Err:        err,
		}
	}

	if len(results) == 0 {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
			Err:        errors.New("data not found"),
		}
	}

	paginate.Rows = results

	return paginate, nil
}

// usp_comLookUp
// IF @strEntity = 'ItemOprCode'--OPERATION MASTER & ITEM MASTER
func (r *LookupRepositoryImpl) ItemOprCodeByID(tx *gorm.DB, linetypeId int, oprItemId int, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var (
		companyCode = 473
		year        string
		month       string
	)

	// fetch data last mtr_location_stock

	var mtrLocationStock masterentities.LocationStock
	err := tx.Table("mtr_location_stock").
		Where("company_id = ?", companyCode).
		Order("period_year DESC, period_month DESC").
		Offset(0).
		Limit(1).
		Find(&mtrLocationStock).Error
	if err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch data last mtr_location_stock",
			Err:        err,
		}
	}

	year = mtrLocationStock.PeriodYear
	month = mtrLocationStock.PeriodMonth

	// Fetch item type from external service
	var itemTypeFetchGoods masteritementities.ItemType
	if err := tx.Where("item_type_code = ?", "G").First(&itemTypeFetchGoods).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Item type not found",
				Err:        fmt.Errorf("item type with code %s not found", "G"),
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item type code",
			Err:        err,
		}
	}

	var itemTypeFetchServices masteritementities.ItemType
	if err := tx.Where("item_type_code = ?", "S").First(&itemTypeFetchServices).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Item type not found",
				Err:        fmt.Errorf("item type with code %s not found", "S"),
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item type code",
			Err:        err,
		}
	}

	// Base Query
	baseQuery := tx.Session(&gorm.Session{NewDB: true})

	switch linetypeId {
	case 1:
		baseQuery = baseQuery.Table("mtr_package A").
			Select("A.package_id, A.package_code, A.package_name, "+
				"COALESCE(SUM(mtr_package_master_detail.frt_quantity), 0) AS frt, "+
				"B.profit_center_name, C.model_code, C.model_description, A.package_price, "+
				"A.model_id, A.brand_id, A.variant_id").
			Joins("INNER JOIN mtr_package_master_detail ON A.package_id = mtr_package_master_detail.package_id").
			Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_profit_center B ON A.profit_center_id = B.profit_center_id").
			Joins("INNER JOIN dms_microservices_sales_dev.dbo.mtr_unit_model C ON A.model_id = C.model_id").
			Where("A.is_active = ?", 1).
			Where("A.package_id = ?", oprItemId).
			Group("A.package_id, A.package_code, A.package_name, B.profit_center_name, " +
				"C.model_code, C.model_description, A.package_price, A.model_id, A.brand_id, A.variant_id").
			Order("A.package_id")

	case 2:
		baseQuery = baseQuery.Table("mtr_operation_model_mapping AS omm").
			Select("omm.operation_id AS operation_id, "+
				"oc.operation_code AS operation_code, oc.operation_name AS operation_name, "+
				"MAX(ofrt.frt_hour) AS frt_hour, "+
				"oe.operation_entries_code AS operation_entries_code, oe.operation_entries_description AS operation_entries_description, "+
				"ok.operation_key_code AS operation_key_code, ok.operation_key_description AS operation_key_description").
			Joins("INNER JOIN mtr_operation_frt AS ofrt ON omm.operation_model_mapping_id = ofrt.operation_model_mapping_id").
			Joins("LEFT OUTER JOIN mtr_operation_code AS oc ON omm.operation_id = oc.operation_id").
			Joins("LEFT OUTER JOIN mtr_operation_entries AS oe ON oc.operation_entries_id = oe.operation_entries_id").
			Joins("LEFT OUTER JOIN mtr_operation_key AS ok ON oc.operation_key_id = ok.operation_key_id").
			Where("omm.is_active = ?", true).
			Where("omm.operation_id = ?", oprItemId).
			Group("omm.operation_id, oc.operation_code, oc.operation_name, " +
				"oe.operation_entries_code, oe.operation_entries_description, " +
				"ok.operation_key_code, ok.operation_key_description").
			Order("omm.operation_id")

	case 3:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassResp masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SP").First(&itemClassResp).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SP"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "SP"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name,
				B.brand_id AS brand_id,
				B.model_id AS model_id,
				B.variant_id AS variant_id,
				COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
				A.item_level_1_id AS item_level_1,
				mil1.item_level_1_code AS item_level_1_code, 
				A.item_level_2_id AS item_level_2,
				mil2.item_level_2_code AS item_level_2_code, 
				A.item_level_3_id AS item_level_3,
				mil3.item_level_3_code AS item_level_3_code, 
				A.item_level_4_id AS item_level_4,
				mil4.item_level_4_code AS item_level_4_code
			`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?",
				itemGrpFetch.ItemGroupId,
				itemTypeFetchGoods.ItemTypeId,
				itemClassResp.ItemClassId,
				true).
			Where("A.item_id = ?", oprItemId).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 4:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassOL masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "OL").First(&itemClassOL).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "OL"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "OL"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name,
				B.brand_id AS brand_id,
				B.model_id AS model_id,
				B.variant_id AS variant_id,
				COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
				A.item_level_1_id AS item_level_1,
				mil1.item_level_1_code AS item_level_1_code, 
				A.item_level_2_id AS item_level_2,
				mil2.item_level_2_code AS item_level_2_code, 
				A.item_level_3_id AS item_level_3,
				mil3.item_level_3_code AS item_level_3_code, 
				A.item_level_4_id AS item_level_4,
				mil4.item_level_4_code AS item_level_4_code
						`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassOL.ItemClassId, true).
			Where("A.item_id = ?", oprItemId).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 5:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassMT masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "MT").First(&itemClassMT).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "MT"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "MT"
		// fetch item class from external service
		var itemClassSB masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SB").First(&itemClassSB).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SB"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "SB"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name,
				B.brand_id AS brand_id,
				B.model_id AS model_id,
				B.variant_id AS variant_id,
				COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
				A.item_level_1_id AS item_level_1,
				mil1.item_level_1_code AS item_level_1_code, 
				A.item_level_2_id AS item_level_2,
				mil2.item_level_2_code AS item_level_2_code, 
				A.item_level_3_id AS item_level_3,
				mil3.item_level_3_code AS item_level_3_code, 
				A.item_level_4_id AS item_level_4,
				mil4.item_level_4_code AS item_level_4_code
						`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type_id = ? AND (A.item_class_id = ? OR A.item_class_id = ?) AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassMT.ItemClassId, itemClassSB.ItemClassId, true).
			Where("A.item_id = ?", oprItemId).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 6:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassWF masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "WF").First(&itemClassWF).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "WF"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "WF"
		// Fetch item group from external service
		var itemGrpOJFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "OJ").First(&itemGrpOJFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "OJ"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		} // "OJ"

		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name,
					B.brand_id AS brand_id,
					B.model_id AS model_id,
					B.variant_id AS variant_id,
					COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
						`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("(A.item_group_id = ? OR A.item_group_id = ?) AND A.item_class_id = ? AND A.item_type_id = ? AND A.is_active = ?", itemGrpOJFetch.ItemGroupId, itemGrpFetch.ItemGroupId, itemClassWF.ItemClassId, itemTypeFetchServices.ItemTypeId, true).
			Where("A.item_id = ?", oprItemId).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 7:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassAC masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "AC").First(&itemClassAC).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "AC"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "AC"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
						A.item_id AS item_id, 
						A.item_code AS item_code, 
						A.item_name AS item_name,
						B.brand_id AS brand_id,
						B.model_id AS model_id,
						B.variant_id AS variant_id,
						COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
						A.item_level_1_id AS item_level_1,
						mil1.item_level_1_code AS item_level_1_code, 
						A.item_level_2_id AS item_level_2,
						mil2.item_level_2_code AS item_level_2_code, 
						A.item_level_3_id AS item_level_3,
						mil3.item_level_3_code AS item_level_3_code, 
						A.item_level_4_id AS item_level_4,
						mil4.item_level_4_code AS item_level_4_code
						`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", itemClassAC.ItemClassId, itemGrpFetch.ItemGroupId, true).
			Where("A.item_id = ?", oprItemId).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 8:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassCM masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "CM").First(&itemClassCM).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "CM"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "CM"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name,
					B.brand_id AS brand_id,
					B.model_id AS model_id,
					B.variant_id AS variant_id,
					COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
							`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ?  AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassCM.ItemClassId, true).
			Where("A.item_id = ?", oprItemId).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 9:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassSV masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SV").First(&itemClassSV).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SV"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "SV"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name,
				B.brand_id AS brand_id,
				B.model_id AS model_id,
				B.variant_id AS variant_id,
				COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
				A.item_level_1_id AS item_level_1,
				mil1.item_level_1_code AS item_level_1_code, 
				A.item_level_2_id AS item_level_2,
				mil2.item_level_2_code AS item_level_2_code, 
				A.item_level_3_id AS item_level_3,
				mil3.item_level_3_code AS item_level_3_code, 
				A.item_level_4_id AS item_level_4,
				mil4.item_level_4_code AS item_level_4_code
			`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyCode).
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", itemClassSV.ItemClassId, itemGrpFetch.ItemGroupId, true).
			Where("A.item_id = ?", oprItemId).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	default:
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid line type",
			Err:        errors.New("invalid line type"),
		}
	}

	// apply filter manual baseon linetype
	for _, filter := range filters {

		if linetypeId == 1 {
			switch filter.ColumnField {
			case "package_id":
				baseQuery = baseQuery.Where("A.package_id = ?", filter.ColumnValue)
			case "package_code":
				baseQuery = baseQuery.Where("A.package_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "package_name":
				baseQuery = baseQuery.Where("A.package_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "profit_center_name":
				baseQuery = baseQuery.Where("B.profit_center_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_code":
				baseQuery = baseQuery.Where("C.model_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_description":
				baseQuery = baseQuery.Where("C.model_description LIKE ?", "%"+filter.ColumnValue+"%")
			case "package_price":
				baseQuery = baseQuery.Where("A.package_price = ?", filter.ColumnValue)
			}
		} else if linetypeId == 2 {
			switch filter.ColumnField {
			case "operation_id":
				baseQuery = baseQuery.Where("oc.operation_id = ?", filter.ColumnValue)
			case "operation_code":
				baseQuery = baseQuery.Where("oc.operation_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_name":
				baseQuery = baseQuery.Where("oc.operation_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "frt_hour":
				baseQuery = baseQuery.Where("ofrt.frt_hour = ?", filter.ColumnValue)
			case "operation_entries_code":
				baseQuery = baseQuery.Where("oe.operation_entries_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_entries_description":
				baseQuery = baseQuery.Where("oe.operation_entries_description LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_key_code":
				baseQuery = baseQuery.Where("ok.operation_key_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_key_description":
				baseQuery = baseQuery.Where("ok.operation_key_description LIKE ?", "%"+filter.ColumnValue+"%")
			}
		} else if linetypeId == 3 || linetypeId == 4 || linetypeId == 5 || linetypeId == 6 || linetypeId == 7 || linetypeId == 8 || linetypeId == 9 {
			switch filter.ColumnField {
			case "item_id":
				baseQuery = baseQuery.Where("A.item_id = ?", filter.ColumnValue)
			case "item_code":
				baseQuery = baseQuery.Where("A.item_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_name":
				baseQuery = baseQuery.Where("A.item_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "available_qty":
				baseQuery = baseQuery.Where("available_qty = ?", filter.ColumnValue)
			case "item_level_1_code":
				baseQuery = baseQuery.Where("mil1.item_level_1_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_2_code":
				baseQuery = baseQuery.Where("mil2.item_level_2_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_3_code":
				baseQuery = baseQuery.Where("mil3.item_level_3_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_4_code":
				baseQuery = baseQuery.Where("mil4.item_level_4_code LIKE ?", "%"+filter.ColumnValue+"%")
			}
		}
	}

	paginateFunc := pagination.Paginate(&paginate, baseQuery)
	baseQuery = baseQuery.Scopes(paginateFunc)

	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total rows",
			Err:        err,
		}
	}

	results := map[string]interface{}{}
	paginateFunc = pagination.Paginate(&paginate, baseQuery)
	if err := baseQuery.Scopes(paginateFunc).Find(&results).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch data",
			Err:        err,
		}
	}

	if len(results) == 0 {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
			Err:        errors.New("data not found"),
		}
	}

	paginate.Rows = results

	return paginate, nil
}

// usp_comLookUp
// IF @strEntity = 'ItemOprCodeWithPrice'--OPERATION MASTER & ITEM MASTER WITH PRICELIST
func (r *LookupRepositoryImpl) ItemOprCodeWithPrice(tx *gorm.DB, linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, trxTypeId int, jobTypeId int, variantId int, currencyId int, whsGroup string, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var (
		currentTime = time.Now()
		year, month = currentTime.Year(), int(currentTime.Month() - 1)
	)

	// Fetch item type from external service
	var itemTypeFetchGoods masteritementities.ItemType
	if err := tx.Where("item_type_code = ?", "G").First(&itemTypeFetchGoods).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Item type not found",
				Err:        fmt.Errorf("item type with code %s not found", "G"),
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item type code",
			Err:        err,
		}
	}

	var itemTypeFetchServices masteritementities.ItemType
	if err := tx.Where("item_type_code = ?", "S").First(&itemTypeFetchServices).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Item type not found",
				Err:        fmt.Errorf("item type with code %s not found", "S"),
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item type code",
			Err:        err,
		}
	}

	baseQuery := tx.Session(&gorm.Session{NewDB: true})

	switch linetypeId {
	case 1:
		baseQuery = baseQuery.Table("mtr_package").
			Select("mtr_package.package_id, mtr_package.package_code, mtr_package.package_name, "+
				"COALESCE(SUM(mtr_package_master_detail.frt_quantity), 0) AS frt, "+
				"mtr_package.profit_center_id as profit_center, dms_microservices_general_dev.dbo.mtr_profit_center.profit_center_name, dms_microservices_sales_dev.dbo.mtr_unit_model.model_code, dms_microservices_sales_dev.dbo.mtr_unit_model.model_description as description, mtr_package.package_price as price, "+
				"mtr_package.model_id, mtr_package.brand_id, mtr_package.variant_id").
			Joins("INNER JOIN mtr_package_master_detail ON mtr_package.package_id = mtr_package_master_detail.package_id").
			Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_profit_center ON mtr_package.profit_center_id = dms_microservices_general_dev.dbo.mtr_profit_center.profit_center_id").
			Joins("INNER JOIN dms_microservices_sales_dev.dbo.mtr_unit_model ON mtr_package.model_id = dms_microservices_sales_dev.dbo.mtr_unit_model.model_id").
			Where("mtr_package.is_active = ?", 1)

		baseQuery = utils.ApplyFilter(baseQuery, filters)

		baseQuery = baseQuery.Group("mtr_package.package_id, mtr_package.package_code, mtr_package.package_name, mtr_package.profit_center_id, dms_microservices_general_dev.dbo.mtr_profit_center.profit_center_name, " +
			"dms_microservices_sales_dev.dbo.mtr_unit_model.model_code, dms_microservices_sales_dev.dbo.mtr_unit_model.model_description, mtr_package.package_price, mtr_package.model_id, mtr_package.brand_id, mtr_package.variant_id").
			Order("mtr_package.package_id")

	case 2:
		baseQuery = baseQuery.Table("mtr_operation_model_mapping").
			Select("mtr_operation_model_mapping.operation_id AS operation_id, "+
				"mtr_operation_code.operation_code AS operation_code, mtr_operation_code.operation_name AS operation_name, "+
				"MAX(mtr_operation_frt.frt_hour) AS frt_hour, "+
				"mtr_operation_entries.operation_entries_code AS operation_entries_code, mtr_operation_entries.operation_entries_description AS operation_entries_description, "+
				"mtr_operation_key.operation_key_code AS operation_key_code, mtr_operation_key.operation_key_description AS operation_key_description").
			Joins("INNER JOIN mtr_operation_frt ON mtr_operation_model_mapping.operation_model_mapping_id = mtr_operation_frt.operation_model_mapping_id").
			Joins("LEFT OUTER JOIN mtr_operation_code ON mtr_operation_model_mapping.operation_id = mtr_operation_code.operation_id").
			Joins("LEFT OUTER JOIN mtr_operation_entries ON mtr_operation_code.operation_entries_id = mtr_operation_entries.operation_entries_id").
			Joins("LEFT OUTER JOIN mtr_operation_key ON mtr_operation_code.operation_key_id = mtr_operation_key.operation_key_id").
			Where("mtr_operation_model_mapping.is_active = ?", true)
		baseQuery = utils.ApplyFilter(baseQuery, filters)
		baseQuery = baseQuery.Group("mtr_operation_model_mapping.operation_id, mtr_operation_code.operation_code, mtr_operation_code.operation_name, " +
			"mtr_operation_entries.operation_entries_code, mtr_operation_entries.operation_entries_description, " +
			"mtr_operation_key.operation_key_code, mtr_operation_key.operation_key_description").
			Order("mtr_operation_model_mapping.operation_id")

	case 3:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassResp masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SP").First(&itemClassResp).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SP"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "SP"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name,
					B.brand_id AS brand_id,
					B.model_id AS model_id,
					B.variant_id AS variant_id,
					COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
				`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyId).
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?",
				itemGrpFetch.ItemGroupId,
				itemTypeFetchGoods.ItemTypeId,
				itemClassResp.ItemClassId,
				true).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 4:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassOL masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "OL").First(&itemClassOL).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "OL"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "OL"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
			A.item_id AS item_id, 
			A.item_code AS item_code, 
			A.item_name AS item_name,
			B.brand_id AS brand_id,
			B.model_id AS model_id,
			B.variant_id AS variant_id,
			COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
			A.item_level_1_id AS item_level_1,
			mil1.item_level_1_code AS item_level_1_code, 
			A.item_level_2_id AS item_level_2,
			mil2.item_level_2_code AS item_level_2_code, 
			A.item_level_3_id AS item_level_3,
			mil3.item_level_3_code AS item_level_3_code, 
			A.item_level_4_id AS item_level_4,
			mil4.item_level_4_code AS item_level_4_code
					`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyId).
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassOL.ItemClassId, true).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 5:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassMT masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "MT").First(&itemClassMT).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "MT"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "MT"
		// fetch item class from external service
		var itemClassSB masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SB").First(&itemClassSB).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SB"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "SB"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name,
					B.brand_id AS brand_id,
					B.model_id AS model_id,
					B.variant_id AS variant_id,
					COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
							`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyId).
			Where("A.item_group_id = ? AND A.item_type_id = ? AND (A.item_class_id = ? OR A.item_class_id = ?) AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassMT.ItemClassId, itemClassSB.ItemClassId, true).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 6:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassWF masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "WF").First(&itemClassWF).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "WF"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "WF"
		// Fetch item group from external service
		var itemGrpOJFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "OJ").First(&itemGrpOJFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "OJ"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		} // "OJ"

		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
						A.item_id AS item_id, 
						A.item_code AS item_code, 
						A.item_name AS item_name,
						B.brand_id AS brand_id,
						B.model_id AS model_id,
						B.variant_id AS variant_id,
						COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
						A.item_level_1_id AS item_level_1,
						mil1.item_level_1_code AS item_level_1_code, 
						A.item_level_2_id AS item_level_2,
						mil2.item_level_2_code AS item_level_2_code, 
						A.item_level_3_id AS item_level_3,
						mil3.item_level_3_code AS item_level_3_code, 
						A.item_level_4_id AS item_level_4,
						mil4.item_level_4_code AS item_level_4_code
							`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyId).
			Where("(A.item_group_id = ? OR A.item_group_id = ?) AND A.item_class_id = ? AND A.item_type_id = ? AND A.is_active = ?", itemGrpOJFetch.ItemGroupId, itemGrpFetch.ItemGroupId, itemClassWF.ItemClassId, itemTypeFetchServices.ItemTypeId, true).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 7:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassAC masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "AC").First(&itemClassAC).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "AC"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "AC"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
							A.item_id AS item_id, 
							A.item_code AS item_code, 
							A.item_name AS item_name,
							B.brand_id AS brand_id,
							B.model_id AS model_id,
							B.variant_id AS variant_id,
							COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
							A.item_level_1_id AS item_level_1,
							mil1.item_level_1_code AS item_level_1_code, 
							A.item_level_2_id AS item_level_2,
							mil2.item_level_2_code AS item_level_2_code, 
							A.item_level_3_id AS item_level_3,
							mil3.item_level_3_code AS item_level_3_code, 
							A.item_level_4_id AS item_level_4,
							mil4.item_level_4_code AS item_level_4_code
							`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyId).
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", itemClassAC.ItemClassId, itemGrpFetch.ItemGroupId, true).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 8:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassCM masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "CM").First(&itemClassCM).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "CM"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "CM"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
						A.item_id AS item_id, 
						A.item_code AS item_code, 
						A.item_name AS item_name,
						B.brand_id AS brand_id,
						B.model_id AS model_id,
						B.variant_id AS variant_id,
						COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
						A.item_level_1_id AS item_level_1,
						mil1.item_level_1_code AS item_level_1_code, 
						A.item_level_2_id AS item_level_2,
						mil2.item_level_2_code AS item_level_2_code, 
						A.item_level_3_id AS item_level_3,
						mil3.item_level_3_code AS item_level_3_code, 
						A.item_level_4_id AS item_level_4,
						mil4.item_level_4_code AS item_level_4_code
								`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyId).
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ?  AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassCM.ItemClassId, true).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")

	case 9:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassSV masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SV").First(&itemClassSV).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SV"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "SV"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name,
					B.brand_id AS brand_id,
					B.model_id AS model_id,
					B.variant_id AS variant_id,
					COALESCE(SUM(V.quantity_allocated), 0) AS available_qty, 
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
				`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Joins("LEFT JOIN mtr_location_stock V ON V.item_id = A.item_id AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?", year, month, companyId).
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", itemClassSV.ItemClassId, itemGrpFetch.ItemGroupId, true).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id").
			Order("A.item_id")
	default:
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid line type",
			Err:        errors.New("invalid line type"),
		}
	}

	for _, filter := range filters {
		if linetypeId == 1 {
			switch filter.ColumnField {
			case "package_id":
				baseQuery = baseQuery.Where("mtr_package.package_id = ?", filter.ColumnValue)
			case "package_code":
				baseQuery = baseQuery.Where("mtr_package.package_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "package_name":
				baseQuery = baseQuery.Where("mtr_package.package_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "profit_center_id":
				baseQuery = baseQuery.Where("mtr_package.profit_center_id = ?", filter.ColumnValue)
			case "profit_center_name":
				baseQuery = baseQuery.Where("dms_microservices_general_dev.dbo.mtr_profit_center.profit_center_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_code":
				baseQuery = baseQuery.Where("dms_microservices_sales_dev.dbo.mtr_unit_model.model_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_description":
				baseQuery = baseQuery.Where("dms_microservices_sales_dev.dbo.mtr_unit_model.model_description LIKE ?", "%"+filter.ColumnValue+"%")
			case "package_price":
				baseQuery = baseQuery.Where("mtr_package.package_price = ?", filter.ColumnValue)
			case "model_id":
				baseQuery = baseQuery.Where("mtr_package.model_id = ?", filter.ColumnValue)
			case "brand_id":
				baseQuery = baseQuery.Where("mtr_package.brand_id = ?", filter.ColumnValue)
			case "variant_id":
				baseQuery = baseQuery.Where("mtr_package.variant_id = ?", filter.ColumnValue)
			}
		} else if linetypeId == 2 {
			switch filter.ColumnField {
			case "operation_id":
				baseQuery = baseQuery.Where("oc.operation_id = ?", filter.ColumnValue)
			case "operation_code":
				baseQuery = baseQuery.Where("oc.operation_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_name":
				baseQuery = baseQuery.Where("oc.operation_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "frt_hour":
				baseQuery = baseQuery.Where("ofrt.frt_hour = ?", filter.ColumnValue)
			case "operation_entries_code":
				baseQuery = baseQuery.Where("oe.operation_entries_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_entries_description":
				baseQuery = baseQuery.Where("oe.operation_entries_description LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_key_code":
				baseQuery = baseQuery.Where("ok.operation_key_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_key_description":
				baseQuery = baseQuery.Where("ok.operation_key_description LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_id":
				baseQuery = baseQuery.Where("omm.model_id = ?", filter.ColumnValue)
			case "brand_id":
				baseQuery = baseQuery.Where("omm.brand_id = ?", filter.ColumnValue)
			case "variant_id":
				baseQuery = baseQuery.Where("ofrt.variant_id = ?", filter.ColumnValue)
			}

		} else if linetypeId == 3 || linetypeId == 4 || linetypeId == 5 || linetypeId == 6 || linetypeId == 7 || linetypeId == 8 || linetypeId == 9 {
			switch filter.ColumnField {
			case "item_id":
				baseQuery = baseQuery.Where("A.item_id = ?", filter.ColumnValue)
			case "item_code":
				baseQuery = baseQuery.Where("A.item_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_name":
				baseQuery = baseQuery.Where("A.item_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "available_qty":
				baseQuery = baseQuery.Where("available_qty = ?", filter.ColumnValue)
			case "item_level_1_code":
				baseQuery = baseQuery.Where("mil1.item_level_1_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_2_code":
				baseQuery = baseQuery.Where("mil2.item_level_2_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_3_code":
				baseQuery = baseQuery.Where("mil3.item_level_3_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_4_code":
				baseQuery = baseQuery.Where("mil4.item_level_4_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_id":
				baseQuery = baseQuery.Where("B.model_id = ?", filter.ColumnValue)
			case "brand_id":
				baseQuery = baseQuery.Where("B.brand_id = ?", filter.ColumnValue)
			case "variant_id":
				baseQuery = baseQuery.Where("B.variant_id = ?", filter.ColumnValue)
			}
		}
	}

	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total rows",
			Err:        err,
		}
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(paginate.Limit)))
	paginateFunc := pagination.Paginate(&paginate, baseQuery)
	baseQuery = baseQuery.Scopes(paginateFunc)

	type LineType0Response struct {
		Description      string  `json:"description"`
		FRT              float64 `json:"frt"`
		ModelCode        string  `json:"model_code"`
		PackageCode      string  `json:"package_code"`
		PackageID        int     `json:"package_id"`
		PackageName      string  `json:"package_name"`
		Price            float64 `json:"price"`
		ProfitCenter     int     `json:"profit_center"`
		ProfitCenterName string  `json:"profit_center_name"`
	}

	type LineType1Response struct {
		FrtHour                     float64 `json:"frt_hour"`
		OperationCode               string  `json:"operation_code"`
		OperationEntriesCode        *string `json:"operation_entries_code"`
		OperationEntriesDescription *string `json:"operation_entries_description"`
		OperationID                 int     `json:"operation_id"`
		OperationKeyCode            *string `json:"operation_key_code"`
		OperationKeyDescription     *string `json:"operation_key_description"`
		OperationName               string  `json:"operation_name"`
		Price                       float64 `json:"price"`
	}

	type LineType2To9Response struct {
		AvailableQty   int     `json:"available_qty"`
		ItemCode       string  `json:"item_code"`
		ItemID         int     `json:"item_id"`
		ItemLevel1     int     `json:"item_level_1"`
		ItemLevel1Code string  `json:"item_level_1_code"`
		ItemLevel2     int     `json:"item_level_2"`
		ItemLevel2Code string  `json:"item_level_2_code"`
		ItemLevel3     int     `json:"item_level_3"`
		ItemLevel3Code string  `json:"item_level_3_code"`
		ItemLevel4     int     `json:"item_level_4"`
		ItemLevel4Code string  `json:"item_level_4_code"`
		ItemName       string  `json:"item_name"`
		Price          float64 `json:"price"`
	}

	switch linetypeId {
	case 1:
		var results []LineType0Response
		if err := baseQuery.Find(&results).Error; err != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch data",
				Err:        err,
			}
		}
		paginate.Rows = results

	case 2:
		var results []LineType1Response
		if err := baseQuery.Find(&results).Error; err != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch data",
				Err:        err,
			}
		}

		for i := range results {
			price, err := r.GetOprItemPrice(tx, linetypeId, companyId, results[i].OperationID, 0, 0, 0, 0, 11, 0, "")
			if err != nil {
				return pagination.Pagination{}, err
			}
			results[i].Price = price
		}
		paginate.Rows = results

	case 3, 4, 5, 6, 7, 8, 9:
		var results []LineType2To9Response
		if err := baseQuery.Find(&results).Error; err != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch data",
				Err:        err,
			}
		}

		for i := range results {
			price, err := r.GetOprItemPrice(tx, linetypeId, companyId, results[i].ItemID, 0, 0, 0, 0, 11, 0, "")
			if err != nil {
				return pagination.Pagination{}, err
			}
			results[i].Price = price
		}
		paginate.Rows = results
	}

	paginate.TotalRows = totalRows
	paginate.TotalPages = totalPages

	return paginate, nil

}

// usp_comLookUp
// IF @strEntity = 'ItemOprCodeWithPrice'--OPERATION MASTER & ITEM MASTER WITH PRICELIST
func (r *LookupRepositoryImpl) ItemOprCodeWithPriceByID(tx *gorm.DB, linetypeId int, OprItemCode int, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var (
		companyId = 426
	)

	// Fetch item type from external service
	var itemTypeFetchGoods masteritementities.ItemType
	if err := tx.Where("item_type_code = ?", "G").First(&itemTypeFetchGoods).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Item type not found",
				Err:        fmt.Errorf("item type with code %s not found", "G"),
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item type code",
			Err:        err,
		}
	}

	var itemTypeFetchServices masteritementities.ItemType
	if err := tx.Where("item_type_code = ?", "S").First(&itemTypeFetchServices).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Item type not found",
				Err:        fmt.Errorf("item type with code %s not found", "S"),
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item type code",
			Err:        err,
		}
	}

	baseQuery := tx.Session(&gorm.Session{NewDB: true})

	switch linetypeId {
	case 1:
		baseQuery = baseQuery.Table("mtr_package A").
			Select("A.package_id, A.package_code, A.package_name, "+
				"COALESCE(SUM(mtr_package_master_detail.frt_quantity), 0) AS frt, "+
				"B.profit_center_name, C.model_code, C.model_description, A.package_price, "+
				"A.model_id, A.brand_id, A.variant_id").
			Joins("INNER JOIN mtr_package_master_detail ON A.package_id = mtr_package_master_detail.package_id").
			Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_profit_center B ON A.profit_center_id = B.profit_center_id").
			Joins("INNER JOIN dms_microservices_sales_dev.dbo.mtr_unit_model C ON A.model_id = C.model_id").
			Where("A.is_active = ?", 1).
			Where("A.package_id = ?", OprItemCode).
			Group("A.package_id, A.package_code, A.package_name, B.profit_center_name, " +
				"C.model_code, C.model_description, A.package_price, A.model_id, A.brand_id, A.variant_id")

	case 2:
		baseQuery = baseQuery.Table("mtr_operation_model_mapping AS omm").
			Select("omm.operation_id AS operation_id, "+
				"oc.operation_code AS operation_code, oc.operation_name AS operation_name, "+
				"MAX(ofrt.frt_hour) AS frt_hour, "+
				"oe.operation_entries_code AS operation_entries_code, oe.operation_entries_description AS operation_entries_description, "+
				"ok.operation_key_code AS operation_key_code, ok.operation_key_description AS operation_key_description").
			Joins("INNER JOIN mtr_operation_frt AS ofrt ON omm.operation_model_mapping_id = ofrt.operation_model_mapping_id").
			Joins("LEFT OUTER JOIN mtr_operation_code AS oc ON omm.operation_id = oc.operation_id").
			Joins("LEFT OUTER JOIN mtr_operation_entries AS oe ON oc.operation_entries_id = oe.operation_entries_id").
			Joins("LEFT OUTER JOIN mtr_operation_key AS ok ON oc.operation_key_id = ok.operation_key_id").
			Where("omm.is_active = ?", true).
			Where("omm.operation_id = ?", OprItemCode).
			Group("omm.operation_id, oc.operation_code, oc.operation_name, " +
				"oe.operation_entries_code, oe.operation_entries_description, " +
				"ok.operation_key_code, ok.operation_key_description")

	case 3:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassResp masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SP").First(&itemClassResp).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SP"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "SP"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name,
					B.brand_id AS brand_id,
					B.model_id AS model_id,
					B.variant_id AS variant_id,
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
				`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?",
				itemGrpFetch.ItemGroupId,
				itemTypeFetchGoods.ItemTypeId,
				itemClassResp.ItemClassId,
				true).
			Where("A.item_id = ?", OprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id")

	case 4:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassOL masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "OL").First(&itemClassOL).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "OL"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "OL"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
			A.item_id AS item_id, 
			A.item_code AS item_code, 
			A.item_name AS item_name,
			B.brand_id AS brand_id,
			B.model_id AS model_id,
			B.variant_id AS variant_id,
			A.item_level_1_id AS item_level_1,
			mil1.item_level_1_code AS item_level_1_code, 
			A.item_level_2_id AS item_level_2,
			mil2.item_level_2_code AS item_level_2_code, 
			A.item_level_3_id AS item_level_3,
			mil3.item_level_3_code AS item_level_3_code, 
			A.item_level_4_id AS item_level_4,
			mil4.item_level_4_code AS item_level_4_code
					`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassOL.ItemClassId, true).
			Where("A.item_id = ?", OprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id")

	case 5:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassMT masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "MT").First(&itemClassMT).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "MT"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "MT"
		// fetch item class from external service
		var itemClassSB masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SB").First(&itemClassSB).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SB"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "SB"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name,
					B.brand_id AS brand_id,
					B.model_id AS model_id,
					B.variant_id AS variant_id,
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
							`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND (A.item_class_id = ? OR A.item_class_id = ?) AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassMT.ItemClassId, itemClassSB.ItemClassId, true).
			Where("A.item_id = ?", OprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id")

	case 6:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassWF masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "WF").First(&itemClassWF).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "WF"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "WF"
		// Fetch item group from external service
		var itemGrpOJFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "OJ").First(&itemGrpOJFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "OJ"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		} // "OJ"

		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
						A.item_id AS item_id, 
						A.item_code AS item_code, 
						A.item_name AS item_name,
						B.brand_id AS brand_id,
						B.model_id AS model_id,
						B.variant_id AS variant_id,
						A.item_level_1_id AS item_level_1,
						mil1.item_level_1_code AS item_level_1_code, 
						A.item_level_2_id AS item_level_2,
						mil2.item_level_2_code AS item_level_2_code, 
						A.item_level_3_id AS item_level_3,
						mil3.item_level_3_code AS item_level_3_code, 
						A.item_level_4_id AS item_level_4,
						mil4.item_level_4_code AS item_level_4_code
							`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("(A.item_group_id = ? OR A.item_group_id = ?) AND A.item_class_id = ? AND A.item_type_id = ? AND A.is_active = ?", itemGrpOJFetch.ItemGroupId, itemGrpFetch.ItemGroupId, itemClassWF.ItemClassId, itemTypeFetchServices.ItemTypeId, true).
			Where("A.item_id = ?", OprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id")

	case 7:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassAC masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "AC").First(&itemClassAC).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "AC"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "AC"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
							A.item_id AS item_id, 
							A.item_code AS item_code, 
							A.item_name AS item_name,
							B.brand_id AS brand_id,
							B.model_id AS model_id,
							B.variant_id AS variant_id,
							A.item_level_1_id AS item_level_1,
							mil1.item_level_1_code AS item_level_1_code, 
							A.item_level_2_id AS item_level_2,
							mil2.item_level_2_code AS item_level_2_code, 
							A.item_level_3_id AS item_level_3,
							mil3.item_level_3_code AS item_level_3_code, 
							A.item_level_4_id AS item_level_4,
							mil4.item_level_4_code AS item_level_4_code
							`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", itemClassAC.ItemClassId, itemGrpFetch.ItemGroupId, true).
			Where("A.item_id = ?", OprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id")

	case 8:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassCM masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "CM").First(&itemClassCM).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "CM"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "CM"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
						A.item_id AS item_id, 
						A.item_code AS item_code, 
						A.item_name AS item_name,
						B.brand_id AS brand_id,
						B.model_id AS model_id,
						B.variant_id AS variant_id,
						A.item_level_1_id AS item_level_1,
						mil1.item_level_1_code AS item_level_1_code, 
						A.item_level_2_id AS item_level_2,
						mil2.item_level_2_code AS item_level_2_code, 
						A.item_level_3_id AS item_level_3,
						mil3.item_level_3_code AS item_level_3_code, 
						A.item_level_4_id AS item_level_4,
						mil4.item_level_4_code AS item_level_4_code
								`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ?  AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassCM.ItemClassId, true).
			Where("A.item_id = ?", OprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id")

	case 9:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassSV masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SV").First(&itemClassSV).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SV"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "SV"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name,
					B.brand_id AS brand_id,
					B.model_id AS model_id,
					B.variant_id AS variant_id,
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
				`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", itemClassSV.ItemClassId, itemGrpFetch.ItemGroupId, true).
			Where("A.item_id = ?", OprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id")
	default:
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid line type",
			Err:        errors.New("invalid line type"),
		}
	}

	for _, filter := range filters {
		if linetypeId == 1 {
			switch filter.ColumnField {
			case "package_id":
				baseQuery = baseQuery.Where("A.package_id = ?", filter.ColumnValue)
			case "package_code":
				baseQuery = baseQuery.Where("A.package_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "package_name":
				baseQuery = baseQuery.Where("A.package_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "profit_center_name":
				baseQuery = baseQuery.Where("B.profit_center_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_code":
				baseQuery = baseQuery.Where("C.model_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_description":
				baseQuery = baseQuery.Where("C.model_description LIKE ?", "%"+filter.ColumnValue+"%")
			case "package_price":
				baseQuery = baseQuery.Where("A.package_price = ?", filter.ColumnValue)
			case "model_id":
				baseQuery = baseQuery.Where("A.model_id = ?", filter.ColumnValue)
			case "brand_id":
				baseQuery = baseQuery.Where("A.brand_id = ?", filter.ColumnValue)
			case "variant_id":
				baseQuery = baseQuery.Where("A.variant_id = ?", filter.ColumnValue)
			}
		} else if linetypeId == 2 {
			switch filter.ColumnField {
			case "operation_id":
				baseQuery = baseQuery.Where("oc.operation_id = ?", filter.ColumnValue)
			case "operation_code":
				baseQuery = baseQuery.Where("oc.operation_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_name":
				baseQuery = baseQuery.Where("oc.operation_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "frt_hour":
				baseQuery = baseQuery.Where("ofrt.frt_hour = ?", filter.ColumnValue)
			case "operation_entries_code":
				baseQuery = baseQuery.Where("oe.operation_entries_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_entries_description":
				baseQuery = baseQuery.Where("oe.operation_entries_description LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_key_code":
				baseQuery = baseQuery.Where("ok.operation_key_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_key_description":
				baseQuery = baseQuery.Where("ok.operation_key_description LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_id":
				baseQuery = baseQuery.Where("omm.model_id = ?", filter.ColumnValue)
			case "brand_id":
				baseQuery = baseQuery.Where("omm.brand_id = ?", filter.ColumnValue)
			case "variant_id":
				baseQuery = baseQuery.Where("ofrt.variant_id = ?", filter.ColumnValue)
			}

		} else if linetypeId == 3 || linetypeId == 4 || linetypeId == 5 || linetypeId == 6 || linetypeId == 7 || linetypeId == 8 || linetypeId == 9 {
			switch filter.ColumnField {
			case "item_id":
				baseQuery = baseQuery.Where("A.item_id = ?", filter.ColumnValue)
			case "item_code":
				baseQuery = baseQuery.Where("A.item_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_name":
				baseQuery = baseQuery.Where("A.item_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "available_qty":
				baseQuery = baseQuery.Where("available_qty = ?", filter.ColumnValue)
			case "item_level_1_code":
				baseQuery = baseQuery.Where("mil1.item_level_1_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_2_code":
				baseQuery = baseQuery.Where("mil2.item_level_2_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_3_code":
				baseQuery = baseQuery.Where("mil3.item_level_3_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_4_code":
				baseQuery = baseQuery.Where("mil4.item_level_4_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_id":
				baseQuery = baseQuery.Where("B.model_id = ?", filter.ColumnValue)
			case "brand_id":
				baseQuery = baseQuery.Where("B.brand_id = ?", filter.ColumnValue)
			case "variant_id":
				baseQuery = baseQuery.Where("B.variant_id = ?", filter.ColumnValue)
			}
		}
	}

	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total rows",
			Err:        err,
		}
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(paginate.Limit)))

	paginateFunc := pagination.Paginate(&paginate, baseQuery)
	baseQuery = baseQuery.Scopes(paginateFunc)

	type LineType0Response struct {
		PackageID        int     `json:"package_id"`
		PackageName      string  `json:"package_name"`
		PackageCode      string  `json:"package_code"`
		Description      string  `json:"description"`
		FRT              float64 `json:"frt"`
		ModelId          int     `json:"model_id"`
		ModelCode        string  `json:"model_code"`
		Price            int     `json:"price"`
		ProfitCenter     int     `json:"profit_center"`
		ProfitCenterName string  `json:"profit_center_name"`
		BrandId          int     `json:"brand_id"`
	}

	type LineType1Response struct {
		FrtHour                     float64 `json:"frt_hour"`
		OperationCode               string  `json:"operation_code"`
		OperationEntriesCode        *string `json:"operation_entries_code"`
		OperationEntriesDescription *string `json:"operation_entries_description"`
		OperationID                 int     `json:"operation_id"`
		OperationKeyCode            *string `json:"operation_key_code"`
		OperationKeyDescription     *string `json:"operation_key_description"`
		OperationName               string  `json:"operation_name"`
		Price                       float64 `json:"price"`
		BrandId                     int     `json:"brand_id"`
		ModelId                     int     `json:"model_id"`
	}

	type LineType2To9Response struct {
		ItemCode       string  `json:"item_code"`
		ItemID         int     `json:"item_id"`
		ItemLevel1     int     `json:"item_level_1"`
		ItemLevel1Code string  `json:"item_level_1_code"`
		ItemLevel2     int     `json:"item_level_2"`
		ItemLevel2Code string  `json:"item_level_2_code"`
		ItemLevel3     int     `json:"item_level_3"`
		ItemLevel3Code string  `json:"item_level_3_code"`
		ItemLevel4     int     `json:"item_level_4"`
		ItemLevel4Code string  `json:"item_level_4_code"`
		ItemName       string  `json:"item_name"`
		Price          float64 `json:"price"`
		BrandId        int     `json:"brand_id"`
		ModelId        int     `json:"model_id"`
	}

	var results interface{}
	switch linetypeId {
	case 1:

		var result LineType0Response
		if err := baseQuery.Find(&result).Error; err != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch data",
				Err:        err,
			}
		}
		results = result
		paginate.Rows = results

	case 2:

		var result LineType1Response
		if err := baseQuery.Find(&result).Error; err != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch data",
				Err:        err,
			}
		}

		price, err := r.GetOprItemPrice(tx, linetypeId, companyId, result.OperationID, result.BrandId, result.ModelId, 2, 1, 11, 6, "")
		if err != nil {
			return pagination.Pagination{}, err
		}
		result.Price = price
		results = result
		paginate.Rows = results

	case 3, 4, 5, 6, 7, 8, 9:
		var result LineType2To9Response
		if err := baseQuery.Find(&result).Error; err != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch data",
				Err:        err,
			}
		}

		price, err := r.GetOprItemPrice(tx, linetypeId, companyId, result.ItemID, result.BrandId, result.ModelId, 2, 1, 11, 6, "")
		if err != nil {
			return pagination.Pagination{}, err
		}

		result.Price = price
		results = result
		paginate.Rows = results
	}

	paginate.TotalRows = totalRows
	paginate.TotalPages = totalPages

	return paginate, nil

}

// usp_comLookUp
// IF @strEntity = 'ItemOprCodeWithPrice'--OPERATION MASTER & ITEM MASTER WITH PRICELIST
func (r *LookupRepositoryImpl) ItemOprCodeWithPriceByCode(tx *gorm.DB, linetypeId int, OprItemCode string, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var (
		companyId = 426
	)

	// Fetch item type from external service
	var itemTypeFetchGoods masteritementities.ItemType
	if err := tx.Where("item_type_code = ?", "G").First(&itemTypeFetchGoods).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Item type not found",
				Err:        fmt.Errorf("item type with code %s not found", "G"),
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item type code",
			Err:        err,
		}
	}

	var itemTypeFetchServices masteritementities.ItemType
	if err := tx.Where("item_type_code = ?", "S").First(&itemTypeFetchServices).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Item type not found",
				Err:        fmt.Errorf("item type with code %s not found", "S"),
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item type code",
			Err:        err,
		}
	}

	baseQuery := tx.Session(&gorm.Session{NewDB: true})

	switch linetypeId {
	case 1:
		baseQuery = baseQuery.Table("mtr_package A").
			Select("A.package_id, A.package_code, A.package_name, "+
				"COALESCE(SUM(mtr_package_master_detail.frt_quantity), 0) AS frt, "+
				"B.profit_center_name, C.model_code, C.model_description, A.package_price, "+
				"A.model_id, A.brand_id, A.variant_id").
			Joins("INNER JOIN mtr_package_master_detail ON A.package_id = mtr_package_master_detail.package_id").
			Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_profit_center B ON A.profit_center_id = B.profit_center_id").
			Joins("INNER JOIN dms_microservices_sales_dev.dbo.mtr_unit_model C ON A.model_id = C.model_id").
			Where("A.is_active = ?", 1).
			Where("A.package_code = ?", OprItemCode).
			Group("A.package_id, A.package_code, A.package_name, B.profit_center_name, " +
				"C.model_code, C.model_description, A.package_price, A.model_id, A.brand_id, A.variant_id")

	case 2:
		baseQuery = baseQuery.Table("mtr_operation_model_mapping AS omm").
			Select("omm.operation_id AS operation_id, "+
				"oc.operation_code AS operation_code, oc.operation_name AS operation_name, "+
				"MAX(ofrt.frt_hour) AS frt_hour, "+
				"oe.operation_entries_code AS operation_entries_code, oe.operation_entries_description AS operation_entries_description, "+
				"ok.operation_key_code AS operation_key_code, ok.operation_key_description AS operation_key_description").
			Joins("INNER JOIN mtr_operation_frt AS ofrt ON omm.operation_model_mapping_id = ofrt.operation_model_mapping_id").
			Joins("LEFT OUTER JOIN mtr_operation_code AS oc ON omm.operation_id = oc.operation_id").
			Joins("LEFT OUTER JOIN mtr_operation_entries AS oe ON oc.operation_entries_id = oe.operation_entries_id").
			Joins("LEFT OUTER JOIN mtr_operation_key AS ok ON oc.operation_key_id = ok.operation_key_id").
			Where("omm.is_active = ?", true).
			Where("omm.operation_code = ?", OprItemCode).
			Group("omm.operation_id, oc.operation_code, oc.operation_name, " +
				"oe.operation_entries_code, oe.operation_entries_description, " +
				"ok.operation_key_code, ok.operation_key_description")

	case 3:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassResp masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SP").First(&itemClassResp).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SP"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "SP"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name,
					B.brand_id AS brand_id,
					B.model_id AS model_id,
					B.variant_id AS variant_id,
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
				`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?",
				itemGrpFetch.ItemGroupId,
				itemTypeFetchGoods.ItemTypeId,
				itemClassResp.ItemClassId,
				true).
			Where("A.item_code = ?", OprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id")

	case 4:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassOL masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "OL").First(&itemClassOL).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "OL"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "OL"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
			A.item_id AS item_id, 
			A.item_code AS item_code, 
			A.item_name AS item_name,
			B.brand_id AS brand_id,
			B.model_id AS model_id,
			B.variant_id AS variant_id,
			A.item_level_1_id AS item_level_1,
			mil1.item_level_1_code AS item_level_1_code, 
			A.item_level_2_id AS item_level_2,
			mil2.item_level_2_code AS item_level_2_code, 
			A.item_level_3_id AS item_level_3,
			mil3.item_level_3_code AS item_level_3_code, 
			A.item_level_4_id AS item_level_4,
			mil4.item_level_4_code AS item_level_4_code
					`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassOL.ItemClassId, true).
			Where("A.item_code = ?", OprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id")

	case 5:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassMT masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "MT").First(&itemClassMT).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "MT"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "MT"
		// fetch item class from external service
		var itemClassSB masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SB").First(&itemClassSB).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SB"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "SB"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name,
					B.brand_id AS brand_id,
					B.model_id AS model_id,
					B.variant_id AS variant_id,
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
							`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND (A.item_class_id = ? OR A.item_class_id = ?) AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassMT.ItemClassId, itemClassSB.ItemClassId, true).
			Where("A.item_code = ?", OprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id")

	case 6:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassWF masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "WF").First(&itemClassWF).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "WF"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "WF"
		// Fetch item group from external service
		var itemGrpOJFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "OJ").First(&itemGrpOJFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "OJ"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		} // "OJ"

		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
						A.item_id AS item_id, 
						A.item_code AS item_code, 
						A.item_name AS item_name,
						B.brand_id AS brand_id,
						B.model_id AS model_id,
						B.variant_id AS variant_id,
						A.item_level_1_id AS item_level_1,
						mil1.item_level_1_code AS item_level_1_code, 
						A.item_level_2_id AS item_level_2,
						mil2.item_level_2_code AS item_level_2_code, 
						A.item_level_3_id AS item_level_3,
						mil3.item_level_3_code AS item_level_3_code, 
						A.item_level_4_id AS item_level_4,
						mil4.item_level_4_code AS item_level_4_code
							`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("(A.item_group_id = ? OR A.item_group_id = ?) AND A.item_class_id = ? AND A.item_type_id = ? AND A.is_active = ?", itemGrpOJFetch.ItemGroupId, itemGrpFetch.ItemGroupId, itemClassWF.ItemClassId, itemTypeFetchServices.ItemTypeId, true).
			Where("A.item_code = ?", OprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id")

	case 7:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassAC masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "AC").First(&itemClassAC).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "AC"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "AC"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
							A.item_id AS item_id, 
							A.item_code AS item_code, 
							A.item_name AS item_name,
							B.brand_id AS brand_id,
							B.model_id AS model_id,
							B.variant_id AS variant_id,
							A.item_level_1_id AS item_level_1,
							mil1.item_level_1_code AS item_level_1_code, 
							A.item_level_2_id AS item_level_2,
							mil2.item_level_2_code AS item_level_2_code, 
							A.item_level_3_id AS item_level_3,
							mil3.item_level_3_code AS item_level_3_code, 
							A.item_level_4_id AS item_level_4,
							mil4.item_level_4_code AS item_level_4_code
							`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", itemClassAC.ItemClassId, itemGrpFetch.ItemGroupId, true).
			Where("A.item_code = ?", OprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id")

	case 8:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassCM masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "CM").First(&itemClassCM).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "CM"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "CM"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
						A.item_id AS item_id, 
						A.item_code AS item_code, 
						A.item_name AS item_name,
						B.brand_id AS brand_id,
						B.model_id AS model_id,
						B.variant_id AS variant_id,
						A.item_level_1_id AS item_level_1,
						mil1.item_level_1_code AS item_level_1_code, 
						A.item_level_2_id AS item_level_2,
						mil2.item_level_2_code AS item_level_2_code, 
						A.item_level_3_id AS item_level_3,
						mil3.item_level_3_code AS item_level_3_code, 
						A.item_level_4_id AS item_level_4,
						mil4.item_level_4_code AS item_level_4_code
								`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ?  AND A.is_active = ?", itemGrpFetch.ItemGroupId, itemTypeFetchGoods.ItemTypeId, itemClassCM.ItemClassId, true).
			Where("A.item_code = ?", OprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id")

	case 9:
		// Fetch item group from external service
		var itemGrpFetch masteritementities.ItemGroup
		if err := tx.Where("item_group_code = ?", "IN").First(&itemGrpFetch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with code %s not found", "IN"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group code",
				Err:        err,
			}
		}
		// fetch item class from external service
		var itemClassSV masteritementities.ItemClass
		if err := tx.Where("item_class_code = ?", "SV").First(&itemClassSV).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pagination.Pagination{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item class not found",
					Err:        fmt.Errorf("item class with code %s not found", "SV"),
				}
			}
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item class code",
				Err:        err,
			}
		} // "SV"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name,
					B.brand_id AS brand_id,
					B.model_id AS model_id,
					B.variant_id AS variant_id,
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
				`).
			Joins("INNER JOIN mtr_item_detail B ON B.item_id = A.item_id").
			Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", itemClassSV.ItemClassId, itemGrpFetch.ItemGroupId, true).
			Where("A.item_code = ?", OprItemCode).
			Group("A.item_id, A.item_code, A.item_name, A.item_level_1_id, mil1.item_level_1_code, A.item_level_2_id, mil2.item_level_2_code, A.item_level_3_id, mil3.item_level_3_code, A.item_level_4_id, mil4.item_level_4_code, B.brand_Id, B.model_id, B.variant_id")
	default:
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid line type",
			Err:        errors.New("invalid line type"),
		}
	}

	for _, filter := range filters {
		if linetypeId == 1 {
			switch filter.ColumnField {
			case "package_id":
				baseQuery = baseQuery.Where("A.package_id = ?", filter.ColumnValue)
			case "package_code":
				baseQuery = baseQuery.Where("A.package_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "package_name":
				baseQuery = baseQuery.Where("A.package_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "profit_center_name":
				baseQuery = baseQuery.Where("B.profit_center_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_code":
				baseQuery = baseQuery.Where("C.model_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_description":
				baseQuery = baseQuery.Where("C.model_description LIKE ?", "%"+filter.ColumnValue+"%")
			case "package_price":
				baseQuery = baseQuery.Where("A.package_price = ?", filter.ColumnValue)
			case "model_id":
				baseQuery = baseQuery.Where("A.model_id = ?", filter.ColumnValue)
			case "brand_id":
				baseQuery = baseQuery.Where("A.brand_id = ?", filter.ColumnValue)
			case "variant_id":
				baseQuery = baseQuery.Where("A.variant_id = ?", filter.ColumnValue)
			}
		} else if linetypeId == 2 {
			switch filter.ColumnField {
			case "operation_id":
				baseQuery = baseQuery.Where("oc.operation_id = ?", filter.ColumnValue)
			case "operation_code":
				baseQuery = baseQuery.Where("oc.operation_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_name":
				baseQuery = baseQuery.Where("oc.operation_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "frt_hour":
				baseQuery = baseQuery.Where("ofrt.frt_hour = ?", filter.ColumnValue)
			case "operation_entries_code":
				baseQuery = baseQuery.Where("oe.operation_entries_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_entries_description":
				baseQuery = baseQuery.Where("oe.operation_entries_description LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_key_code":
				baseQuery = baseQuery.Where("ok.operation_key_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "operation_key_description":
				baseQuery = baseQuery.Where("ok.operation_key_description LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_id":
				baseQuery = baseQuery.Where("omm.model_id = ?", filter.ColumnValue)
			case "brand_id":
				baseQuery = baseQuery.Where("omm.brand_id = ?", filter.ColumnValue)
			case "variant_id":
				baseQuery = baseQuery.Where("ofrt.variant_id = ?", filter.ColumnValue)
			}

		} else if linetypeId == 3 || linetypeId == 4 || linetypeId == 5 || linetypeId == 6 || linetypeId == 7 || linetypeId == 8 || linetypeId == 9 {
			switch filter.ColumnField {
			case "item_id":
				baseQuery = baseQuery.Where("A.item_id = ?", filter.ColumnValue)
			case "item_code":
				baseQuery = baseQuery.Where("A.item_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_name":
				baseQuery = baseQuery.Where("A.item_name LIKE ?", "%"+filter.ColumnValue+"%")
			case "available_qty":
				baseQuery = baseQuery.Where("available_qty = ?", filter.ColumnValue)
			case "item_level_1_code":
				baseQuery = baseQuery.Where("mil1.item_level_1_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_2_code":
				baseQuery = baseQuery.Where("mil2.item_level_2_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_3_code":
				baseQuery = baseQuery.Where("mil3.item_level_3_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "item_level_4_code":
				baseQuery = baseQuery.Where("mil4.item_level_4_code LIKE ?", "%"+filter.ColumnValue+"%")
			case "model_id":
				baseQuery = baseQuery.Where("B.model_id = ?", filter.ColumnValue)
			case "brand_id":
				baseQuery = baseQuery.Where("B.brand_id = ?", filter.ColumnValue)
			case "variant_id":
				baseQuery = baseQuery.Where("B.variant_id = ?", filter.ColumnValue)
			}
		}
	}

	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total rows",
			Err:        err,
		}
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(paginate.Limit)))

	paginateFunc := pagination.Paginate(&paginate, baseQuery)
	baseQuery = baseQuery.Scopes(paginateFunc)

	type LineType0Response struct {
		PackageID        int     `json:"package_id"`
		PackageName      string  `json:"package_name"`
		PackageCode      string  `json:"package_code"`
		Description      string  `json:"description"`
		FRT              float64 `json:"frt"`
		ModelId          int     `json:"model_id"`
		ModelCode        string  `json:"model_code"`
		Price            int     `json:"price"`
		ProfitCenter     int     `json:"profit_center"`
		ProfitCenterName string  `json:"profit_center_name"`
		BrandId          int     `json:"brand_id"`
	}

	type LineType1Response struct {
		FrtHour                     float64 `json:"frt_hour"`
		OperationCode               string  `json:"operation_code"`
		OperationEntriesCode        *string `json:"operation_entries_code"`
		OperationEntriesDescription *string `json:"operation_entries_description"`
		OperationID                 int     `json:"operation_id"`
		OperationKeyCode            *string `json:"operation_key_code"`
		OperationKeyDescription     *string `json:"operation_key_description"`
		OperationName               string  `json:"operation_name"`
		Price                       float64 `json:"price"`
		BrandId                     int     `json:"brand_id"`
		ModelId                     int     `json:"model_id"`
	}

	type LineType2To9Response struct {
		ItemCode       string  `json:"item_code"`
		ItemID         int     `json:"item_id"`
		ItemLevel1     int     `json:"item_level_1"`
		ItemLevel1Code string  `json:"item_level_1_code"`
		ItemLevel2     int     `json:"item_level_2"`
		ItemLevel2Code string  `json:"item_level_2_code"`
		ItemLevel3     int     `json:"item_level_3"`
		ItemLevel3Code string  `json:"item_level_3_code"`
		ItemLevel4     int     `json:"item_level_4"`
		ItemLevel4Code string  `json:"item_level_4_code"`
		ItemName       string  `json:"item_name"`
		Price          float64 `json:"price"`
		BrandId        int     `json:"brand_id"`
		ModelId        int     `json:"model_id"`
	}

	var results interface{}
	switch linetypeId {
	case 1:

		var result LineType0Response
		if err := baseQuery.Find(&result).Error; err != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch data",
				Err:        err,
			}
		}
		results = result
		paginate.Rows = results

	case 2:

		var result LineType1Response
		if err := baseQuery.Find(&result).Error; err != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch data",
				Err:        err,
			}
		}

		price, err := r.GetOprItemPrice(tx, linetypeId, companyId, result.OperationID, result.BrandId, result.ModelId, 2, 1, 11, 6, "")
		if err != nil {
			return pagination.Pagination{}, err
		}
		result.Price = price
		results = result
		paginate.Rows = results

	case 3, 4, 5, 6, 7, 8, 9:
		var result LineType2To9Response
		if err := baseQuery.Find(&result).Error; err != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch data",
				Err:        err,
			}
		}

		price, err := r.GetOprItemPrice(tx, linetypeId, companyId, result.ItemID, result.BrandId, result.ModelId, 2, 1, 11, 6, "")
		if err != nil {
			return pagination.Pagination{}, err
		}

		result.Price = price
		results = result
		paginate.Rows = results
	}

	paginate.TotalRows = totalRows
	paginate.TotalPages = totalPages

	return paginate, nil

}

// usp_comLookUp
// IF @strEntity = 'Vehicle0'--VEHICLE UNIT MASTER
func (r *LookupRepositoryImpl) GetVehicleUnitMaster(tx *gorm.DB, brandId int, modelId int, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var (
		vehicleMasters []map[string]interface{}
		totalRows      int64
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	query := tx.Table("dms_microservices_sales_dev.dbo.mtr_vehicle V").
		Select(`
			V.vehicle_id AS vehicle_id,
			V.vehicle_chassis_number AS vehicle_chassis_number, 
			RC.vehicle_registration_certificate_tnkb AS vehicle_registration_certificate_tnkb, 
			RC.vehicle_registration_certificate_owner_name AS vehicle_registration_certificate_owner_name, 
			CONCAT(B.brand_name,' ',VA.variant_description,' ')  AS vehicle, 
			CAST(V.vehicle_production_year AS VARCHAR) AS vehicle_production_year, 
			CONVERT(VARCHAR, V.vehicle_last_service_date, 106) AS vehicle_last_service_date, 
			V.vehicle_last_km AS vehicle_last_km, 
			CASE 
				WHEN V.is_active = 1 THEN 'Active' 
				WHEN V.is_active = 0 THEN 'Deactive' 
			END AS status,
			V.user_customer_id as customer_id,
			V.vehicle_variant_id as variant_id,
			V.vehicle_colour_id as colour_id
		`).
		Joins(`LEFT JOIN dms_microservices_sales_dev.dbo.mtr_vehicle_registration_certificate RC ON V.vehicle_id = RC.vehicle_id`).
		Joins(`LEFT JOIN dms_microservices_sales_dev.dbo.mtr_model_variant_colour UM ON UM.brand_id = V.vehicle_brand_id AND 
									UM.model_id = V.vehicle_model_id AND 
									UM.colour_id = V.vehicle_colour_id AND 
									ISNULL(UM.accessories_option_id, '') = ISNULL(V.option_id, '')`).
		Joins(`INNER JOIN dms_microservices_sales_dev.dbo.mtr_brand B ON B.brand_id = V.vehicle_brand_id`).
		Joins(`INNER JOIN dms_microservices_sales_dev.dbo.mtr_unit_model M ON M.model_id = V.vehicle_model_id`).
		Joins(`INNER JOIN dms_microservices_sales_dev.dbo.mtr_unit_variant VA ON VA.variant_id = V.vehicle_variant_id`).
		Where(filterQuery, filterValues...).
		Where("V.vehicle_brand_id = ?", brandId).
		Where("V.vehicle_model_id = ?", modelId)

	err := query.Count(&totalRows).Error
	if err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total vehicle units",
			Err:        err,
		}
	}

	err = query.
		Scopes(pagination.Paginate(&paginate, query)).
		Find(&vehicleMasters).Error

	if err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get vehicle unit master data",
			Err:        err,
		}
	}

	paginate.Rows = vehicleMasters
	paginate.TotalRows = totalRows
	paginate.TotalPages = int(math.Ceil(float64(totalRows) / float64(paginate.GetLimit())))

	return paginate, nil
}

// usp_comLookUp
// IF @strEntity = 'Vehicle0'--VEHICLE UNIT MASTER
func (r *LookupRepositoryImpl) GetVehicleUnitByID(tx *gorm.DB, vehicleID int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	var vehicleMaster map[string]interface{}

	query := tx.Table("dms_microservices_sales_dev.dbo.mtr_vehicle V").
		Select(`
			V.vehicle_id AS vehicle_id,
			V.vehicle_chassis_number AS vehicle_chassis_number, 
			RC.vehicle_registration_certificate_tnkb AS vehicle_registration_certificate_tnkb, 
			RC.vehicle_registration_certificate_owner_name AS vehicle_registration_certificate_owner_name, 
			CONCAT(B.brand_name,' ',VA.variant_description,' ')  AS vehicle,
			CAST(V.vehicle_production_year AS VARCHAR) AS vehicle_production_year, 
			CONVERT(VARCHAR, V.vehicle_last_service_date, 106) AS vehicle_last_service_date, 
			V.vehicle_last_km AS vehicle_last_km, 
			CASE 
				WHEN V.is_active = 1 THEN 'Active' 
				WHEN V.is_active = 0 THEN 'Deactive' 
			END AS status,
			V.user_customer_id as customer_id,
			V.vehicle_variant_id as variant_id,
			V.vehicle_colour_id as colour_id
		`).
		Joins(`LEFT JOIN dms_microservices_sales_dev.dbo.mtr_vehicle_registration_certificate RC ON V.vehicle_id = RC.vehicle_id`).
		Joins(`LEFT JOIN dms_microservices_sales_dev.dbo.mtr_model_variant_colour UM ON UM.brand_id = V.vehicle_brand_id AND 
							UM.model_id = V.vehicle_model_id AND 
							UM.colour_id = V.vehicle_colour_id AND 
							ISNULL(UM.accessories_option_id, '') = ISNULL(V.option_id, '')`).
		Joins(`INNER JOIN dms_microservices_sales_dev.dbo.mtr_brand B ON B.brand_id = V.vehicle_brand_id`).
		Joins(`INNER JOIN dms_microservices_sales_dev.dbo.mtr_unit_model M ON M.model_id = V.vehicle_model_id`).
		Joins(`INNER JOIN dms_microservices_sales_dev.dbo.mtr_unit_variant VA ON VA.variant_id = V.vehicle_variant_id`).
		Where("V.vehicle_id = ?", vehicleID)

	err := query.Take(&vehicleMaster).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get vehicle unit master data",
			Err:        err,
		}
	}

	return vehicleMaster, nil
}

// usp_comLookUp
// IF @strEntity = 'Vehicle0'--VEHICLE UNIT MASTER
func (r *LookupRepositoryImpl) GetVehicleUnitByChassisNumber(tx *gorm.DB, chassisNumber string) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	var vehicleMaster map[string]interface{}

	query := tx.Table("dms_microservices_sales_dev.dbo.mtr_vehicle V").
		Select(`
			V.vehicle_id AS vehicle_id,
			V.vehicle_chassis_number AS vehicle_chassis_number, 
			RC.vehicle_registration_certificate_tnkb AS vehicle_registration_certificate_tnkb, 
			RC.vehicle_registration_certificate_owner_name AS vehicle_registration_certificate_owner_name, 
			CONCAT(B.brand_name,' ',VA.variant_description,' ')  AS vehicle,  
			CAST(V.vehicle_production_year AS VARCHAR) AS vehicle_production_year, 
			CONVERT(VARCHAR, V.vehicle_last_service_date, 106) AS vehicle_last_service_date, 
			V.vehicle_last_km AS vehicle_last_km, 
			CASE 
				WHEN V.is_active = 1 THEN 'Active' 
				WHEN V.is_active = 0 THEN 'Deactive' 
			END AS status,
			,
			V.user_customer_id as customer_id,
			V.vehicle_variant_id as variant_id,
			V.vehicle_colour_id as colour_id
		`).
		Joins(`LEFT JOIN dms_microservices_sales_dev.dbo.mtr_vehicle_registration_certificate RC ON V.vehicle_id = RC.vehicle_id`).
		Joins(`LEFT JOIN dms_microservices_sales_dev.dbo.mtr_model_variant_colour UM ON UM.brand_id = V.vehicle_brand_id AND 
						UM.model_id = V.vehicle_model_id AND 
						UM.colour_id = V.vehicle_colour_id AND 
						ISNULL(UM.accessories_option_id, '') = ISNULL(V.option_id, '')`).
		Joins(`INNER JOIN dms_microservices_sales_dev.dbo.mtr_brand B ON B.brand_id = V.vehicle_brand_id`).
		Joins(`INNER JOIN dms_microservices_sales_dev.dbo.mtr_unit_model M ON M.model_id = V.vehicle_model_id`).
		Joins(`INNER JOIN dms_microservices_sales_dev.dbo.mtr_unit_variant VA ON VA.variant_id = V.vehicle_variant_id`).
		Where("V.vehicle_chassis_number = ?", chassisNumber)

	err := query.Take(&vehicleMaster).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get vehicle unit master data",
			Err:        err,
		}
	}

	return vehicleMaster, nil
}

// usp_comLookUp
// IF @strEntity = 'CampaignMaster'--CAMPAIGN MASTER
func (r *LookupRepositoryImpl) GetCampaignMaster(tx *gorm.DB, companyId int, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var (
		campaignMasters []map[string]interface{}
		totalRows       int64
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	query := tx.Table("dms_microservices_aftersales_dev.dbo.mtr_campaign C").
		Select(`
			C.campaign_id AS campaign_id,
			C.campaign_code AS campaign_code,
			C.campaign_name AS campaign_name,
			C.model_id AS model_id,
			C.campaign_period_from AS campaign_period_from,
			C.campaign_period_to AS campaign_period_to,
			C.total_after_vat AS total_after_vat,
			CASE 
				WHEN C.is_active = 1 THEN 'Active' 
				WHEN C.is_active = 0 THEN 'Deactive' 
			END AS Status
			`).
		Where(filterQuery, filterValues...).
		Where("C.company_id = ?", companyId)

	err := query.Count(&totalRows).Error
	if err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total campaign data",
			Err:        err,
		}
	}

	err = query.
		Scopes(pagination.Paginate(&paginate, query)).
		Find(&campaignMasters).Error

	if err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get campaign master data",
			Err:        err,
		}
	}

	paginate.Rows = campaignMasters
	paginate.TotalRows = totalRows
	paginate.TotalPages = int(math.Ceil(float64(totalRows) / float64(paginate.GetLimit())))

	return paginate, nil
}

// usp_comLookUp
// IF @strEntity = 'WorkOrderService'--WO SERVICE
func (r *LookupRepositoryImpl) WorkOrderService(tx *gorm.DB, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var (
		results   []map[string]interface{}
		totalRows int64
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	// Prepare filters
	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	query := tx.Table("trx_work_order_allocation AS A").
		Select(`
			A.work_order_document_number AS work_order_document_number,
			B.work_order_date AS work_order_date,
			B.vehicle_tnkb AS vehicle_tnkb,
			B.vehicle_chassis_number AS vehicle_chassis_number,
			B.brand_id AS brand_id,
			B.model_id AS model_id,
			B.variant_id AS variant_id,
			A.work_order_system_number AS work_order_system_number
		`).
		Joins("LEFT JOIN trx_work_order AS B ON B.work_order_system_number = A.work_order_system_number").
		Where("A.service_status_id NOT IN (?, ?, ?, ?)", utils.SrvStatStop, utils.SrvStatAutoRelease, utils.SrvStatTransfer, utils.SrvStatQcPass)

	// Apply additional filters
	if len(filterStrings) > 0 {
		query = query.Where(filterQuery, filterValues...)
	}

	// Count total rows
	err := query.Count(&totalRows).Error
	if err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total work orders",
			Err:        err,
		}
	}

	// Paginate the results
	err = query.
		Scopes(pagination.Paginate(&paginate, query)).
		Order("A.work_order_document_number").
		Find(&results).Error

	if err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get work order data",
			Err:        err,
		}
	}

	// Finalize pagination details
	paginate.Rows = results
	paginate.TotalRows = totalRows
	paginate.TotalPages = int(math.Ceil(float64(totalRows) / float64(paginate.GetLimit())))

	return paginate, nil
}

// usp_comLookUp
// IF @strEntity = 'WoAtpmRegistration'--AWS-018 - ATPM Registration
func (r *LookupRepositoryImpl) WorkOrderAtpmRegistration(tx *gorm.DB, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	baseQuery := tx.Table("trx_work_order A").
		Select(`A.work_order_document_number, A.work_order_date`).
		Joins("INNER JOIN trx_work_order_detail B ON A.work_order_system_number = B.work_order_system_number").
		Where("B.work_order_status_id IN (?, ?, ?)", utils.WoStatDraft, utils.WoStatCancel, utils.WoStatClosed).
		Where("COALESCE(A.atpm_claim_number, 0) != 0").
		Where("A.claim_system_number != 0")

	for _, filter := range filters {
		baseQuery = baseQuery.Where(fmt.Sprintf("%s = ?", filter.ColumnField), filter.ColumnValue)
	}

	paginateFunc := pagination.Paginate(&paginate, baseQuery)
	baseQuery = baseQuery.Scopes(paginateFunc)

	type WoAtpmRegistrationResponse struct {
		WorkOrderDocumentNumber string    `json:"work_order_document_number"`
		WorkOrderDate           time.Time `json:"work_order_date"`
		WorkOrderSystemNumber   int       `json:"work_order_system_number"`
	}

	var results []WoAtpmRegistrationResponse
	if err := baseQuery.Find(&results).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch data",
			Err:        err,
		}
	}

	// Convert the results to the map format required by the pagination package
	var response []map[string]interface{}
	for _, result := range results {
		data := map[string]interface{}{
			"work_order_document_number": result.WorkOrderDocumentNumber,
			"work_order_date":            result.WorkOrderDate,
		}
		response = append(response, data)
	}

	// Calculate the total number of rows
	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total rows",
			Err:        err,
		}
	}

	// Calculate the total number of pages
	totalPages := int(math.Ceil(float64(totalRows) / float64(paginate.Limit)))

	// Set the pagination details
	paginate.Rows = response
	paginate.TotalRows = totalRows
	paginate.TotalPages = totalPages

	// Return the paginated result
	return paginate, nil
}

// usp_comLookUp
// IF @strEntity =  'CustomerByTypeAndAddress'--CUSTOMER MASTER
func (r *LookupRepositoryImpl) CustomerByTypeAndAddress(tx *gorm.DB, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var (
		customerMasters []map[string]interface{}
		totalRows       int64
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	// Prepare filters
	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	// Query construction
	query := tx.Table("dms_microservices_general_dev.dbo.mtr_customer C").
		Select(`
			C.customer_id AS customer_id,
			C.customer_code AS customer_code,
			C.customer_name AS customer_name,
			CA.client_type_description AS client_type_description,
			A.address_street_1 AS address_1,
			A.address_street_2 AS address_2,
			A.address_street_3 AS address_3,
			C.id_phone_no AS id_phone_no
		`).
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_client_type CA ON C.client_type_id = CA.client_type_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_address AS A ON C.id_address_id = A.address_id").
		Where("C.is_active = 1")

	// Apply filters if provided
	if len(filterStrings) > 0 {
		query = query.Where(filterQuery, filterValues...)
	}

	// Count total rows
	err := query.Count(&totalRows).Error
	if err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total customers",
			Err:        err,
		}
	}

	// Fetch data with pagination
	err = query.
		Scopes(pagination.Paginate(&paginate, query)).
		Order("C.customer_id").
		Find(&customerMasters).Error

	if err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get customer data",
			Err:        err,
		}
	}

	// Set pagination details
	paginate.Rows = customerMasters
	paginate.TotalRows = totalRows
	paginate.TotalPages = int(math.Ceil(float64(totalRows) / float64(paginate.GetLimit())))

	return paginate, nil
}

// usp_comLookUp
// IF @strEntity =  'CustomerByTypeAndAddress'--CUSTOMER MASTER
func (r *LookupRepositoryImpl) CustomerByTypeAndAddressByID(tx *gorm.DB, customerId int, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var (
		customerMasters []map[string]interface{}
		totalRows       int64
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	query := tx.Table("dms_microservices_general_dev.dbo.mtr_customer C").
		Select(`
			C.customer_id AS customer_id,
			C.customer_code AS customer_code,
			C.customer_name AS customer_name,
			CA.client_type_description AS client_type_description,
			A.address_street_1 AS address_1,
			A.address_street_2 AS address_2,
			A.address_street_3 AS address_3,
			C.id_phone_no AS id_phone_no
		`).
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_client_type CA ON C.client_type_id = CA.client_type_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_address AS A ON C.id_address_id = A.address_id").
		Where("C.customer_id = ?", customerId)

	if len(filterStrings) > 0 {
		query = query.Where(filterQuery, filterValues...)
	}

	err := query.Count(&totalRows).Error
	if err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total customers",
			Err:        err,
		}
	}

	err = query.
		Scopes(pagination.Paginate(&paginate, query)).
		Order("C.customer_id").
		Find(&customerMasters).Error

	if err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get customer data",
			Err:        err,
		}
	}

	paginate.Rows = customerMasters
	paginate.TotalRows = totalRows
	paginate.TotalPages = int(math.Ceil(float64(totalRows) / float64(paginate.GetLimit())))

	return paginate, nil
}

// usp_comLookUp
// IF @strEntity =  'CustomerByTypeAndAddress'--CUSTOMER MASTER
func (r *LookupRepositoryImpl) CustomerByTypeAndAddressByCode(tx *gorm.DB, customerCode string, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var (
		customerMasters []map[string]interface{}
		totalRows       int64
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	query := tx.Table("dms_microservices_general_dev.dbo.mtr_customer C").
		Select(`
			C.customer_id AS customer_id,
			C.customer_code AS customer_code,
			C.customer_name AS customer_name,
			CA.client_type_description AS client_type_description,
			A.address_street_1 AS address_1,
			A.address_street_2 AS address_2,
			A.address_street_3 AS address_3,
			C.id_phone_no AS id_phone_no
		`).
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_client_type CA ON C.client_type_id = CA.client_type_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_address AS A ON C.id_address_id = A.address_id").
		Where("C.customer_code = ?", customerCode)

	if len(filterStrings) > 0 {
		query = query.Where(filterQuery, filterValues...)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total customers",
			Err:        err,
		}
	}

	if err := query.
		Scopes(pagination.Paginate(&paginate, query)).
		Order("C.customer_code").
		Find(&customerMasters).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get customer data",
			Err:        err,
		}
	}

	paginate.Rows = customerMasters
	paginate.TotalRows = totalRows
	paginate.TotalPages = int(math.Ceil(float64(totalRows) / float64(paginate.GetLimit())))

	return paginate, nil
}

// dbo.FCT_getBillCode
// FctGetBillCode retrieves the bill code based on companyCode, supcusCode, and supcusType
func (r *LookupRepositoryImpl) FctGetBillCode(tx *gorm.DB, companyCode int, supcusCode, supcusType string) (string, *exceptions.BaseErrorResponse) {
	var (
		billCode       string
		npwpSupcus     string
		npwpCompany    string
		typeMap        string
		supplierExists bool
		customerExists bool
	)

	// Hardcoded bill code values

	billcodeRelatedParties := "R"
	supcusDealer := "00"
	supcusImsi := "51"
	supcusAtpm := "61"
	supcusSalim := "71"
	supcusMaintained := "81"

	// If supcusType is 'F'
	if strings.ToUpper(supcusType) == "F" {
		if supcusCode == strconv.Itoa(companyCode) {
			billCode = utils.TrxTypeWoInternal.Code
		} else {
			// Check if supcusCode exists in gmSupplier0
			if err := tx.Table("dms_microservices_general_dev.dbo.mtr_supplier").
				Select("1").
				Where("supplier_code = ?", supcusCode).
				Find(&supplierExists).Error; err != nil {
				return "", &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to check supplier existence",
					Err:        err,
				}
			}

			if !supplierExists {
				// Check if supcusCode exists in gmCust0
				if err := tx.Table("dms_microservices_general_dev.dbo.mtr_customer").
					Select("1").
					Where("customer_code = ?", supcusCode).
					Find(&customerExists).Error; err != nil {
					return "", &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to check customer existence",
						Err:        err,
					}
				}

				if !customerExists {
					billCode = ""
				} else {
					var typeMap string
					// check company_type in mtr_company_type_map_from_customer
					if err := tx.Table("dms_microservices_general_dev.dbo.mtr_company_type_map_from_customer").
						Select("client_type_id").
						Where("company_from_id = ? AND company_to_id = ?", companyCode, supcusCode).
						Scan(&typeMap).Error; err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							// Fallback query: check customer_type in mtr_customer if no record is found
							if err := tx.Table("dms_microservices_general_dev.dbo.mtr_customer").
								Select("client_type_id").
								Where("customer_code = ?", supcusCode).
								Scan(&typeMap).Error; err != nil {
								return "", &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to get customer type",
									Err:        err,
								}
							}
						} else {
							return "", &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to get company type",
								Err:        err,
							}
						}
					}

					switch typeMap {
					case supcusDealer:
						// Get VAT_REG_NO for company and supcus
						if err := r.getVatRegNo(tx, companyCode, supcusCode, &npwpCompany, &npwpSupcus); err != nil {
							return "", err
						}

						if npwpCompany == npwpSupcus {
							billCode = utils.TrxTypeWoCentralize.Code
						} else {
							billCode = utils.TrxTypeWoDeCentralize.Code
						}
					case supcusImsi:
						billCode = utils.TrxTypeWoDeCentralize.Code
					case supcusAtpm, supcusSalim, supcusMaintained:
						billCode = billcodeRelatedParties
					default:
						billCode = utils.TrxTypeWoExternal.Code
					}
				}
			} else {
				// Handle supplier logic
				// Query the mtr_company_type_map_from_customer table for company type
				err := tx.Table("dms_microservices_general_dev.dbo.mtr_company_type_map_from_customer").
					Select("client_type_id").
					Where("company_from_id = ? AND company_to_id = ?", companyCode, supcusCode).
					First(&typeMap).Error

				if err == gorm.ErrRecordNotFound {
					err = tx.Table("dms_microservices_general_dev.dbo.mtr_supplier").
						Select("client_type_id").
						Where("supplier_code = ?", supcusCode).
						First(&typeMap).Error
				}

				if err != nil {
					return "", &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to get supplier type",
						Err:        err,
					}
				}

				switch typeMap {
				case supcusDealer:
					if err := r.getVatRegNo(tx, companyCode, supcusCode, &npwpCompany, &npwpSupcus); err != nil {
						return "", err
					}

					if npwpCompany == npwpSupcus {
						billCode = utils.TrxTypeWoCentralize.Code
					} else {
						billCode = utils.TrxTypeWoDeCentralize.Code
					}
				case supcusImsi:
					billCode = utils.TrxTypeWoDeCentralize.Code
				case supcusAtpm, supcusSalim, supcusMaintained:
					billCode = billcodeRelatedParties
				default:
					billCode = utils.TrxTypeWoExternal.Code
				}
			}
		}
	}

	return billCode, nil
}

// Helper function to get VAT_REG_NO for company and supcus
func (r *LookupRepositoryImpl) getVatRegNo(tx *gorm.DB, companyCode int, supcusCode string, npwpCompany, npwpSupcus *string) *exceptions.BaseErrorResponse {
	if err := tx.Table("dms_microservices_general_dev.dbo.mtr_customer").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_tax_data ON mtr_customer.tax_customer_id = mtr_tax_data.tax_id").
		Select("npwp_no").
		Where("mtr_customer.company_id = ?", companyCode).Scan(npwpCompany).Error; err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get company VAT_REG_NO",
			Err:        err,
		}
	}

	if err := tx.Table("dms_microservices_general_dev.dbo.mtr_customer").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_tax_data ON mtr_customer.tax_customer_id = mtr_tax_data.tax_id").
		Select("npwp_no").
		Where("company_id = ?", supcusCode).
		Scan(npwpSupcus).Error; err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get supplier/customer VAT_REG_NO",
			Err:        err,
		}
	}

	return nil
}

// dbo.getLineTypebyItemCode
// GetLineTypeByItemCode retrieves the line type based on the item code
func (r *LookupRepositoryImpl) GetLineTypeByItemCode(tx *gorm.DB, itemCode string) (string, *exceptions.BaseErrorResponse) {
	var (
		linetypeStr    string
		itemGrp        int
		itemTypeId     int
		itemCls        int
		lineTypeSublet = utils.LinetypeSublet
	)

	// Retrieve item details
	var itemDetails struct {
		ItemGroupId int
		ItemTypeId  int
		ItemClassId int
	}

	if err := tx.Model(&masteritementities.Item{}).
		Select("item_group_id, item_type_id, item_class_id").
		Where("item_code = ?", itemCode).
		Scan(&itemDetails).Error; err != nil {
		return "", &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get item details",
			Err:        err,
		}
	}

	itemGrp = itemDetails.ItemGroupId
	itemTypeId = itemDetails.ItemTypeId
	itemCls = itemDetails.ItemClassId

	// Determine line type based on item details
	if itemGrp == 2 {
		if itemTypeId == 1 {
			switch itemCls {
			case 69: // "SP"
				linetypeStr = utils.LinetypeSparepart
			case 70: // "OL"
				linetypeStr = utils.LinetypeOil
			case 71, 72: // "MT"
				linetypeStr = utils.LinetypeMaterial
			case 75: // "CM"
				linetypeStr = utils.LinetypeConsumableMaterial
			case 74: // "SR"
				linetypeStr = utils.LinetypeAccesories
			case 77: // "SV"
				linetypeStr = utils.LinetypeSublet
			default:
				linetypeStr = utils.LinetypeAccesories
			}
		} else if itemCls == 73 { // "WF"
			linetypeStr = lineTypeSublet
		} else if itemCls == 74 && itemTypeId == 2 {
			linetypeStr = utils.LinetypeOperation
		}
	} else if itemGrp == 6 || (itemGrp == 2 && itemTypeId == 2 && itemCls == 73) {
		linetypeStr = lineTypeSublet
	}

	// Check item existence
	var itemExists int64
	if err := tx.Model(&masteritementities.Item{}).
		Where("item_code = ?", itemCode).
		Count(&itemExists).Error; err != nil {
		return "0", &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get item details",
			Err:        err,
		}
	}

	// Handle package existence
	if itemExists == 0 {
		var packageExists int64
		if err := tx.Model(&masterentities.PackageMaster{}).
			Where("package_code = ?", itemCode).
			Count(&packageExists).Error; err != nil {
			return "0", &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get item details",
				Err:        err,
			}
		}
		if packageExists > 0 {
			linetypeStr = utils.LinetypePackage
		} else {
			linetypeStr = utils.LinetypeOperation
		}
	}

	return linetypeStr, nil
}

func (r *LookupRepositoryImpl) GetWhsGroup(tx *gorm.DB, companyCode int) (int, *exceptions.BaseErrorResponse) {
	var (
		whsGroup int
		err      error
	)

	if err = tx.Table("dms_microservices_aftersales_dev.dbo.mtr_warehouse_group_mapping").
		Select("warehouse_group_type_id").
		Where("company_id = ?", companyCode).
		Scan(&whsGroup).Error; err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get warehouse group",
			Err:        err,
		}
	}

	return whsGroup, nil
}

// dbo.getCampaignDiscForWO
// GetCampaignDiscForWO retrieves the campaign discount for work order
func (r *LookupRepositoryImpl) GetCampaignDiscForWO(tx *gorm.DB, campaignId int, linetypeId int, oprItemId int, frtQty float64, markupAmount float64, markupPercentage float64, millage float64) (masterpayloads.CampaignDiscount, *exceptions.BaseErrorResponse) {

	// Nested query for determining VAT_TAX_CODE based on certain conditions
	var vatTaxCode float64

	if err := tx.Table("dms_microservices_sales_dev.dbo.mtr_campaign").
		Select("vat_tax_code").
		Where("campaign_id = ?", campaignId).
		Scan(&vatTaxCode).Error; err != nil {
		return masterpayloads.CampaignDiscount{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get VAT_TAX_CODE",
			Err:        err,
		}
	}

	var campaignDiscount masterpayloads.CampaignDiscount

	if err := tx.Table("mtr_campaign").
		Select("(mtr_campaign_detail.price + ? + (mtr_campaign_detail.price * (? / 100))) AS operation_item_price, "+
			"mtr_campaign_detail.discount_percent as operation_item_discount_percent, "+
			"ROUND(((mtr_campaign_detail.price + ? + (mtr_campaign_detail.price * (? / 100))) * (mtr_campaign_detail.discount_percent / 100)), 0) AS operation_item_discount_amount, "+
			"'5' AS transaction_type_id",
			markupAmount, markupPercentage, markupAmount, markupPercentage).
		Joins("LEFT JOIN mtr_campaign_detail ON mtr_campaign.campaign_id = mtr_campaign_detail.campaign_id").
		Where("mtr_campaign.campaign_id = ? AND mtr_campaign_detail.line_type_id = ? AND mtr_campaign_detail.item_operation_id = ? AND ? >= mtr_campaign_detail.quantity AND ? >= mtr_campaign_detail.millage",
			campaignId, linetypeId, oprItemId, frtQty, millage).
		Scan(&campaignDiscount).Error; err != nil {
		return masterpayloads.CampaignDiscount{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get campaign discount",
			Err:        err,
		}
	}

	return campaignDiscount, nil
}

// usp_comLookUp IF @strEntity = 'ListItemLocation'
func (r *LookupRepositoryImpl) ListItemLocation(tx *gorm.DB, companyId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	response := []masterpayloads.WarehouseMasterForItemLookupResponse{}

	entities := masterwarehouseentities.WarehouseMaster{}

	baseModelQuery := tx.Model(&entities).
		Distinct(`
			mli.warehouse_id,
			mwg.warehouse_group_code,
			mwg.warehouse_group_name,
			mtr_warehouse_master.warehouse_code,
			mtr_warehouse_master.warehouse_name
		`).
		Joins("INNER JOIN mtr_location_item mli ON mtr_warehouse_master.warehouse_id = mli.warehouse_id").
		Joins("INNER JOIN mtr_warehouse_group mwg ON mwg.warehouse_group_id = mtr_warehouse_master.warehouse_group_id").
		Where("mtr_warehouse_master.company_id = ?", companyId)

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&response).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(response) == 0 {
		pages.Rows = []masterpayloads.WarehouseMasterForItemLookupResponse{}
		return pages, nil
	}

	pages.Rows = response

	return pages, nil
}

// IF @strEntity = 'WarehouseGroupByCompany'
func (r *LookupRepositoryImpl) WarehouseGroupByCompany(tx *gorm.DB, companyId int) ([]masterpayloads.WarehouseGroupByCompanyResponse, *exceptions.BaseErrorResponse) {
	entities := masterwarehouseentities.WarehouseMaster{}
	response := []masterpayloads.WarehouseGroupByCompanyResponse{}

	err := tx.Model(&entities).
		Select("DISTINCT mwg.warehouse_group_id, mwg.warehouse_group_code + ' - ' + mwg.warehouse_group_name AS warehouse_group_code_name").
		Joins("INNER JOIN mtr_warehouse_group mwg ON mtr_warehouse_master.warehouse_group_id = mwg.warehouse_group_id").
		Where("mtr_warehouse_master.is_active = ?", true).
		Where("mtr_warehouse_master.company_id = ?", companyId).
		Scan(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

// usp_comLookUp IF @strEntity = 'ItemListTrans'
func (r *LookupRepositoryImpl) ItemListTrans(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Item{}
	responses := []masterpayloads.ItemListTransResponse{}

	baseModelQuery := tx.Model(&entities).
		Select(`
			mtr_item.item_id,
			mtr_item.item_code,
			mtr_item.item_name,
			mil1.item_level_1_code,
			mil2.item_level_2_code,
			mil3.item_level_3_code,
			mil4.item_level_4_code,
			mic.item_class_code,
			mit.item_type_code
		`).
		Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = mtr_item.item_level_1_id").
		Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = mtr_item.item_level_2_id").
		Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = mtr_item.item_level_3_id").
		Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = mtr_item.item_level_4_id").
		Joins("INNER JOIN mtr_item_class mic ON mic.item_class_id = mtr_item.item_class_id").
		Joins("INNER JOIN mtr_item_type mit ON mit.item_type_id = mtr_item.item_type_id").
		Where("mtr_item.is_active = ?", true)

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.Rows = responses

	return pages, nil
}

// usp_comLookUp IF @strEntity = 'ItemListTransPL'
func (r *LookupRepositoryImpl) ItemListTransPL(tx *gorm.DB, companyId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Item{}
	responses := []masterpayloads.ItemListTransPLResponse{}

	baseModelQuery := tx.Model(&entities).
		Select(`
			MIN(mtr_item.item_id) AS item_id,
			mtr_item.item_code,
			mtr_item.item_name,
			mil1.item_level_1_code,
			mil2.item_level_2_code,
			mil3.item_level_3_code,
			mil4.item_level_4_code,
			mic.item_class_code,
			mit.item_type_code
		`).
		Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = mtr_item.item_level_1_id").
		Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = mtr_item.item_level_2_id").
		Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = mtr_item.item_level_3_id").
		Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = mtr_item.item_level_4_id").
		Joins("INNER JOIN mtr_item_class mic ON mic.item_class_id = mtr_item.item_class_id").
		Joins("INNER JOIN mtr_item_type mit ON mit.item_type_id = mtr_item.item_type_id").
		Joins("INNER JOIN mtr_item_detail mid ON mid.item_id = mtr_item.item_id").
		Where("mtr_item.is_active = ?", true).
		Where("mtr_item.price_list_item = ?", true).
		Group(`
			mtr_item.item_code,
			mtr_item.item_name,
			mil1.item_level_1_code,
			mil2.item_level_2_code,
			mil3.item_level_3_code,
			mil4.item_level_4_code,
			mic.item_class_code,
			mit.item_type_code
		`)

	if companyId == 0 {
		baseModelQuery = baseModelQuery.Where("mtr_item.common_pricelist = ?", true)
	}
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&responses).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.Rows = responses

	return pages, nil
}

func (r *LookupRepositoryImpl) SelectLocationStockItem(tx *gorm.DB, option int, companyId int, periodDate time.Time, whsCode int, locCode string, itemId int, whsGroup int, uomType string) (float64, *exceptions.BaseErrorResponse) {
	var qtyResult float64
	var qtyTemp float64
	var periodYear, periodMonth string

	moduleCode := "MODULE_SP"
	periodStatusClose := 3

	// Validate Item Code
	var itemCount int64
	if err := tx.Table("mtr_item").
		Where("item_id = ?", itemId).
		Count(&itemCount).Error; err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count item",
			Err:        err,
		}
	}

	if itemCount == 0 {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Item not found",
			Err:        errors.New("item not found"),
		}
	}

	// Handle period date if not provided
	if periodDate.IsZero() {
		periodDate = time.Now()
	}

	// Get current period year and month
	periodYear = periodDate.Format("2006")
	periodMonth = periodDate.Format("01")

	// Check if the period is closed
	var periodStatusId int
	if err := tx.Table("dms_microservices_finance_dev.dbo.mtr_period_audit").
		Select("period_status_id").
		Where("company_id = ? AND period_year = ? AND period_month = ?", companyId, periodYear, periodMonth).
		Scan(&periodStatusId).Error; err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get period status",
			Err:        err,
		}
	}

	if periodStatusId == periodStatusClose {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Period is closed",
			Err:        errors.New("period is closed"),
		}
	}

	// Build the query based on the option: ON HAND QTY & AVAILABLE QTY
	query := tx.Model(&masterentities.LocationStock{}).
		Where("company_id = ? AND period_year = ? AND period_month = ? AND warehouse_id = ? AND location_id = ? AND item_id = ? AND warehouse_group = ?",
			companyId, periodYear, periodMonth, whsCode, locCode, itemId, whsGroup)

	switch option {
	case 1:
		if err := query.Select("SUM(quantity_ending)").Scan(&qtyTemp).Error; err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get on hand quantity for item",
				Err:        err,
			}
		}
	case 2:
		if err := query.
			Where("module_code = ?", moduleCode).
			Select("SUM(quantity_available)").Scan(&qtyTemp).Error; err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get available quantity for item",
				Err:        err,
			}
		}
	default:
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid option provided",
			Err:        errors.New("invalid option"),
		}
	}

	qtyResult = qtyTemp

	return qtyResult, nil
}

// dbo.getOprItemFrt
// GetOprItemFrt retrieves the FRT hour based on the input parameters
func (r *LookupRepositoryImpl) GetOprItemFrt(tx *gorm.DB, oprItemId int, brandId int, modelId int, variantId int, vehicleChassisNo string) (float64, *exceptions.BaseErrorResponse) {
	var frt float64

	if brandId != 0 && oprItemId != 0 && modelId != 0 {
		// Check if the variant code exists in amOperation2
		var variantExists bool
		if err := tx.Table("mtr_operation_model_mapping").
			Joins("INNER JOIN mtr_operation_frt ON mtr_operation_model_mapping.operation_model_mapping_id = mtr_operation_frt.operation_model_mapping_id").
			Select("1").
			Where("mtr_operation_model_mapping.brand_id = ? AND mtr_operation_model_mapping.model_id = ? AND mtr_operation_model_mapping.operation_id = ? ", brandId, modelId, oprItemId).
			Where("mtr_operation_frt.variant_id = ?", variantId).
			Scan(&variantExists).Error; err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error checking variant code",
				Err:        err,
			}
		}

		// Fetch FRT Hour value
		if err := tx.Table("mtr_operation_model_mapping").
			Joins("INNER JOIN mtr_operation_frt ON mtr_operation_model_mapping.operation_model_mapping_id = mtr_operation_frt.operation_model_mapping_id").
			Select("mtr_operation_frt.frt_hour").
			Where("mtr_operation_model_mapping.brand_id = ? AND mtr_operation_model_mapping.model_id = ? AND mtr_operation_frt.variant_id = ? AND mtr_operation_model_mapping.operation_id = ?",
				brandId, modelId, variantId, oprItemId).
			Row().Scan(&frt); err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error fetching FRT hour",
				Err:        err,
			}
		}

	} else if vehicleChassisNo != "" && oprItemId != 0 {
		// Get vehicle details from umVehicle0 using chassis number
		var vehicle struct {
			variantId int
			brandId   int
			modelId   int
		}
		if err := tx.Table("dms_microservices_sales_dev.dbo.mtr_vehicle").
			Select("vehicle_variant_id, vehicle_brand_id, vehicle_model_id").
			Where("vehicle_chassis_number = ?", vehicleChassisNo).
			Row().Scan(&vehicle.variantId, &vehicle.brandId, &vehicle.modelId); err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error fetching vehicle details",
				Err:        err,
			}
		}

		// Check if the variant code exists in amOperation2
		var variantExists bool
		if err := tx.Table("mtr_operation_model_mapping").
			Joins("INNER JOIN mtr_operation_frt ON mtr_operation_model_mapping.operation_model_mapping_id = mtr_operation_frt.operation_model_mapping_id").
			Select("1").
			Where("mtr_operation_model_mapping.brand_id = ? AND mtr_operation_model_mapping.model_id = ? AND mtr_operation_model_mapping.operation_id = ?", vehicle.brandId, vehicle.modelId, oprItemId).
			Where("mtr_operation_frt.variant_id = ?", vehicle.variantId).
			Scan(&variantExists).Error; err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error checking variant code",
				Err:        err,
			}
		}

		// Fetch FRT Hour value
		if err := tx.Table("mtr_operation_model_mapping").
			Joins("INNER JOIN mtr_operation_frt ON mtr_operation_model_mapping.operation_model_mapping_id = mtr_operation_frt.operation_model_mapping_id").
			Select("mtr_operation_frt.frt_hour").
			Where("mtr_operation_model_mapping.brand_id = ? AND mtr_operation_model_mapping.model_id = ? AND mtr_operation_frt.variant_id = ? AND mtr_operation_model_mapping.operation_id = ?",
				vehicle.brandId, vehicle.modelId, vehicle.variantId, oprItemId).
			Row().Scan(&frt); err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error fetching FRT hour",
				Err:        err,
			}
		}
	}

	return frt, nil
}

// usp_comLookUp
// IF @strEntity = 'ServiceReqRefTypeWO'--SERVICE REQUEST REF TYPE WO
func (r *LookupRepositoryImpl) ReferenceTypeWorkOrder(tx *gorm.DB, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var (
		results []struct {
			WorkOrderDocumentNumber    string    `json:"work_order_document_number"`
			WorkOrderDate              time.Time `json:"work_order_date"`
			WorkOrderStatusId          int       `json:"work_order_status_id"`
			WorkOrderStatusDescription string    `json:"work_order_status_description"`
			WorkOrderSystemNumber      int       `json:"work_order_system_number"`
		}
		totalRows int64
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	query := tx.Table("trx_work_order AS A").
		Select(`
			A.work_order_document_number AS work_order_document_number,
			A.work_order_date AS work_order_date,
			B.work_order_status_id AS work_order_status_id,
			E.work_order_status_description AS work_order_status_description,
			A.work_order_system_number AS work_order_system_number
		`).
		Joins("INNER JOIN trx_work_order_detail AS B ON B.work_order_system_number = A.work_order_system_number").
		Joins(`LEFT OUTER JOIN trx_service_request AS C 
				ON A.service_request_system_number = C.service_request_system_number 
				AND C.reference_type_id = 1 
				AND C.service_request_status_id NOT IN (4, 5) 
				AND NOT (C.service_request_status_id = 8 AND COALESCE(C.booking_system_number, 0) != 0 AND COALESCE(C.work_order_system_number, 0) != 0)`).
		Joins("LEFT OUTER JOIN trx_service_request_detail AS D ON C.service_request_system_number = D.service_request_system_number AND D.operation_item_id = B.operation_item_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_work_order_status AS E ON A.work_order_status_id = E.work_order_status_id").
		Where("B.work_order_status_id NOT IN (?, ?, ?)", utils.WoStatDraft, utils.WoStatClosed, utils.WoStatCancel).
		Where("COALESCE(A.work_order_system_number, 0) != 0").
		Where("COALESCE(D.service_request_line_number, 0) != 0")

	if len(filterStrings) > 0 {
		query = query.Where(filterQuery, filterValues...)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total data",
			Err:        err,
		}
	}

	if err := query.
		Scopes(pagination.Paginate(&paginate, query)).
		Order("A.work_order_document_number").
		Find(&results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get data",
			Err:        err,
		}
	}

	mappedResults := make([]map[string]interface{}, len(results))
	for i, result := range results {
		mappedResults[i] = map[string]interface{}{
			"work_order_document_number":    result.WorkOrderDocumentNumber,
			"work_order_date":               result.WorkOrderDate.Format("2006-01-02"),
			"work_order_status_id":          result.WorkOrderStatusId,
			"work_order_status_description": result.WorkOrderStatusDescription,
			"work_order_system_number":      result.WorkOrderSystemNumber,
		}
	}

	paginate.Rows = mappedResults
	paginate.TotalRows = totalRows
	paginate.TotalPages = int(math.Ceil(float64(totalRows) / float64(paginate.GetLimit())))

	return paginate, nil
}

// usp_comLookUp
// IF @strEntity = 'ServiceReqRefTypeWO'--SERVICE REQUEST REF TYPE WO
func (r *LookupRepositoryImpl) ReferenceTypeWorkOrderByID(tx *gorm.DB, referenceId int, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var (
		results []struct {
			WorkOrderDocumentNumber    string    `json:"work_order_document_number"`
			WorkOrderDate              time.Time `json:"work_order_date"`
			WorkOrderStatusId          int       `json:"work_order_status_id"`
			WorkOrderStatusDescription string    `json:"work_order_status_description"`
			WorkOrderSystemNumber      int       `json:"work_order_system_number"`
		}
		totalRows int64
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	query := tx.Table("trx_work_order AS A").
		Select(`
			A.work_order_document_number AS work_order_document_number,
			A.work_order_date AS work_order_date,
			B.work_order_status_id AS work_order_status_id,
			E.work_order_status_description AS work_order_status_description,
			A.work_order_system_number AS work_order_system_number
		`).
		Joins("INNER JOIN trx_work_order_detail AS B ON B.work_order_system_number = A.work_order_system_number").
		Joins(`LEFT OUTER JOIN trx_service_request AS C 
				ON A.service_request_system_number = C.service_request_system_number 
				AND C.reference_type_id = 1 
				AND C.service_request_status_id NOT IN (4, 5) 
				AND NOT (C.service_request_status_id = 8 AND COALESCE(C.booking_system_number, 0) != 0 AND COALESCE(C.work_order_system_number, 0) != 0)`).
		Joins("LEFT OUTER JOIN trx_service_request_detail AS D ON C.service_request_system_number = D.service_request_system_number AND D.operation_item_id = B.operation_item_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_work_order_status AS E ON A.work_order_status_id = E.work_order_status_id").
		Where("B.work_order_status_id NOT IN (?, ?, ?)", utils.WoStatDraft, utils.WoStatClosed, utils.WoStatCancel).
		Where("COALESCE(A.work_order_system_number, 0) != 0").
		Where("COALESCE(D.service_request_line_number, 0) != 0").
		Where("C.service_request_system_number = ?", referenceId)

	if len(filterStrings) > 0 {
		query = query.Where(filterQuery, filterValues...)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total data",
			Err:        err,
		}
	}

	if err := query.
		Scopes(pagination.Paginate(&paginate, query)).
		Order("A.work_order_document_number").
		Find(&results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get data",
			Err:        err,
		}
	}

	mappedResults := make([]map[string]interface{}, len(results))
	for i, result := range results {
		mappedResults[i] = map[string]interface{}{
			"work_order_document_number":    result.WorkOrderDocumentNumber,
			"work_order_date":               result.WorkOrderDate.Format("2006-01-02"),
			"work_order_status_id":          result.WorkOrderStatusId,
			"work_order_status_description": result.WorkOrderStatusDescription,
			"work_order_system_number":      result.WorkOrderSystemNumber,
		}
	}

	paginate.Rows = mappedResults
	paginate.TotalRows = totalRows
	paginate.TotalPages = int(math.Ceil(float64(totalRows) / float64(paginate.GetLimit())))

	return paginate, nil
}

// usp_comLookUp
// IF @strEntity = 'ServiceReqRefTypeSO'--SERVICE REQUEST REF TYPE SO
func (r *LookupRepositoryImpl) ReferenceTypeSalesOrder(tx *gorm.DB, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var (
		results []struct {
			SalesOrderDocumentNumber    string    `json:"sales_order_document_number"`
			SalesOrderDate              time.Time `json:"sales_order_date"`
			SalesOrderStatusId          int       `json:"sales_order_status_id"`
			SalesOrderStatusDescription string    `json:"sales_order_status_description"`
			SalesOrderSystemNumber      int       `json:"sales_order_system_number"`
		}
		totalRows int64
	)

	// Set default pagination limit
	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	// Process filters
	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	// Build query
	query := tx.Table("trx_sales_order AS A").
		Select(`
			A.sales_order_document_number AS sales_order_document_number,
			A.sales_order_date AS sales_order_date,
			B.sales_order_status_id AS sales_order_status_id,
			E.sales_order_status_description AS sales_order_status_description,
			A.sales_order_system_number AS sales_order_system_number
		`).
		Joins("INNER JOIN trx_sales_order_detail AS B ON B.sales_order_system_number = A.sales_order_system_number").
		Joins(`LEFT OUTER JOIN trx_service_request AS C 
				ON A.service_request_system_number = C.service_request_system_number 
				AND C.reference_type_id = 1 
				AND C.service_request_status_id NOT IN (4, 5) 
				AND NOT (C.service_request_status_id = 8 AND COALESCE(C.booking_system_number, 0) != 0 AND COALESCE(C.sales_order_system_number, 0) != 0)`).
		Joins("LEFT OUTER JOIN trx_service_request_detail AS D ON C.service_request_system_number = D.service_request_system_number AND D.operation_item_id = B.operation_item_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_sales_order_status AS E ON A.sales_order_status_id = E.sales_order_status_id").
		Where("B.sales_order_status_id NOT IN (?, ?, ?)", utils.WoStatDraft, utils.WoStatClosed, utils.WoStatCancel).
		Where("COALESCE(A.sales_order_system_number, 0) != 0").
		Where("COALESCE(D.service_request_line_number, 0) != 0")

	// Apply filters if any
	if len(filterStrings) > 0 {
		query = query.Where(filterQuery, filterValues...)
	}

	// Count total rows
	if err := query.Count(&totalRows).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total data",
			Err:        err,
		}
	}

	// Apply pagination and fetch results
	if err := query.
		Scopes(pagination.Paginate(&paginate, query)).
		Order("A.sales_order_document_number").
		Find(&results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get data",
			Err:        err,
		}
	}

	// Map results
	mappedResults := make([]map[string]interface{}, len(results))
	for i, result := range results {
		mappedResults[i] = map[string]interface{}{
			"sales_order_document_number":    result.SalesOrderDocumentNumber,
			"sales_order_date":               result.SalesOrderDate.Format("2006-01-02"),
			"sales_order_status_id":          result.SalesOrderStatusId,
			"sales_order_status_description": result.SalesOrderStatusDescription,
			"sales_order_system_number":      result.SalesOrderSystemNumber,
		}
	}

	// Set pagination response
	paginate.Rows = mappedResults
	paginate.TotalRows = totalRows
	paginate.TotalPages = int(math.Ceil(float64(totalRows) / float64(paginate.GetLimit())))

	return paginate, nil
}

// usp_comLookUp
// IF @strEntity = 'ServiceReqRefTypeSO'--SERVICE REQUEST REF TYPE SO
func (r *LookupRepositoryImpl) ReferenceTypeSalesOrderByID(tx *gorm.DB, referenceId int, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var (
		result struct {
			SalesOrderDocumentNumber    string    `json:"sales_order_document_number"`
			SalesOrderDate              time.Time `json:"sales_order_date"`
			SalesOrderStatusId          int       `json:"sales_order_status_id"`
			SalesOrderStatusDescription string    `json:"sales_order_status_description"`
			SalesOrderSystemNumber      int       `json:"sales_order_system_number"`
		}
		totalRows int64
	)

	// Set default pagination limit
	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	// Build filters dynamically
	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}

	// Build query
	query := tx.Table("trx_sales_order AS A").
		Select(`
			A.sales_order_document_number AS sales_order_document_number,
			A.sales_order_date AS sales_order_date,
			B.sales_order_status_id AS sales_order_status_id,
			E.sales_order_status_description AS sales_order_status_description,
			A.sales_order_system_number AS sales_order_system_number
		`).
		Joins("INNER JOIN trx_sales_order_detail AS B ON B.sales_order_system_number = A.sales_order_system_number").
		Joins(`LEFT OUTER JOIN trx_service_request AS C
			ON A.service_request_system_number = C.service_request_system_number
			AND C.reference_type_id = 1
			AND C.service_request_status_id NOT IN (4, 5)
			AND NOT (C.service_request_status_id = 8
				AND COALESCE(C.booking_system_number, 0) != 0
				AND COALESCE(C.sales_order_system_number, 0) != 0)
		`).
		Joins("LEFT OUTER JOIN trx_service_request_detail AS D ON C.service_request_system_number = D.service_request_system_number AND D.operation_item_id = B.operation_item_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_sales_order_status AS E ON A.sales_order_status_id = E.sales_order_status_id").
		Where("B.sales_order_status_id NOT IN (?, ?, ?)", utils.WoStatDraft, utils.WoStatClosed, utils.WoStatCancel).
		Where("COALESCE(A.sales_order_system_number, 0) != 0").
		Where("COALESCE(D.service_request_line_number, 0) != 0").
		Where("A.sales_order_system_number = ?", referenceId)

	// Apply filters if present
	if len(filterStrings) > 0 {
		query = query.Where(strings.Join(filterStrings, " AND "), filterValues...)
	}

	// Count total rows
	if err := query.Count(&totalRows).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total data",
			Err:        err,
		}
	}

	// Calculate total pages
	paginate.TotalRows = totalRows
	paginate.TotalPages = int(math.Ceil(float64(totalRows) / float64(paginate.Limit)))

	// Fetch data with pagination
	if err := query.Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get data",
			Err:        err,
		}
	}

	// Map result to pagination response
	paginate.Rows = []map[string]interface{}{
		{
			"sales_order_document_number":    result.SalesOrderDocumentNumber,
			"sales_order_date":               result.SalesOrderDate.Format("2006-01-02"),
			"sales_order_status_id":          result.SalesOrderStatusId,
			"sales_order_status_description": result.SalesOrderStatusDescription,
			"sales_order_system_number":      result.SalesOrderSystemNumber,
		},
	}

	return paginate, nil
}

func (r *LookupRepositoryImpl) GetLineTypeByReferenceType(tx *gorm.DB, referenceTypeId int, paginate pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var lineTypes []map[string]interface{}
	var excludedIds []int

	switch referenceTypeId {
	case 1:
		excludedIds = []int{1, 8}
	case 2:
		excludedIds = []int{2, 8}
	default:
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid reference type ID",
			Err:        fmt.Errorf("unsupported reference type ID: %d", referenceTypeId),
		}
	}

	tx = tx.Table("dms_microservices_general_dev.dbo.mtr_line_type").
		Select("line_type_id, line_type_code, line_type_name").
		Where("line_type_id NOT IN ?", excludedIds)

	// Menghitung total rows
	var totalRows int64
	if err := tx.Count(&totalRows).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total data",
			Err:        err,
		}
	}

	// Menerapkan paginasi
	tx = tx.Scopes(pagination.Paginate(&paginate, tx))

	if err := tx.Find(&lineTypes).Error; err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get line type",
			Err:        err,
		}
	}

	paginate.Rows = lineTypes
	paginate.TotalRows = totalRows
	paginate.TotalPages = int(math.Ceil(float64(totalRows) / float64(paginate.GetLimit())))

	return paginate, nil
}

// usp_comLookUp
// IF @strEntity = 'LocationAvailable'
// Used for insert item location detail
func (r *LookupRepositoryImpl) LocationAvailable(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := masterwarehouseentities.WarehouseLocation{}
	response := []masterpayloads.LocationAvailableResponse{}
	var err error

	var newFilter []utils.FilterCondition
	var companyId int
	var warehouseId int
	for _, filter := range filterCondition {
		if strings.Contains(filter.ColumnField, "company_id") {
			companyId, _ = strconv.Atoi(filter.ColumnValue)
			continue
		}
		if strings.Contains(filter.ColumnField, "warehouse_id") {
			warehouseId, _ = strconv.Atoi(filter.ColumnValue)
			continue
		}
		newFilter = append(newFilter, filter)
	}

	periodResponse, periodError := financeserviceapiutils.GetOpenPeriodByCompany(companyId, "SP")
	if periodError != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching company current period",
			Err:        periodError.Err,
		}
	}

	var existingWarehouseLocIds []int
	err = tx.Model(&masteritementities.ItemLocation{}).
		Joins("INNER JOIN mtr_item mi ON mi.item_id = mtr_location_item.item_id").
		Joins("INNER JOIN mtr_item_group mig ON mig.item_group_id = mi.item_group_id").
		Joins("INNER JOIN mtr_warehouse_master mwm ON mwm.warehouse_id = mtr_location_item.warehouse_id").
		Where("mwm.company_id = ?", companyId).
		Where("mtr_location_item.warehouse_id = ?", warehouseId).
		Where("mig.item_group_code = 'IN'").
		Pluck("warehouse_location_id", &existingWarehouseLocIds).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching existing item location data",
			Err:        err,
		}
	}
	if len(existingWarehouseLocIds) == 0 {
		existingWarehouseLocIds = []int{-1}
	}
	existingWarehouseLocIds = utils.RemoveDuplicateIds(existingWarehouseLocIds)

	viewLocStock := tx.Table("mtr_location_stock mls").
		Select(`
			mls.item_inquiry_id,
			mls.period_year,
			mls.period_month,
			mls.company_id,
			mls.warehouse_group_id,
			mls.warehouse_id,
			mls.location_id,
			(
				ISNULL(quantity_sales, 0) +
				ISNULL(quantity_transfer_out, 0) +
				ISNULL(quantity_claim_out, 0) +
				ISNULL(quantity_robbing_out, 0) +
				ISNULL(quantity_assembly_out, 0)
			) AS quantity_on_hand
		`).
		Joins("LEFT JOIN mtr_warehouse_master mwm ON mwm.company_id = mls.company_id AND mwm.warehouse_id = mls.warehouse_id")

	baseModelQuery := tx.Model(&entities).
		Select(`
			mtr_warehouse_location.is_active,
			mtr_warehouse_location.warehouse_location_id,
			mtr_warehouse_location.warehouse_location_code,
			mtr_warehouse_location.warehouse_location_name,
			mwm.company_id,
			mtr_warehouse_location.warehouse_group_id,
			mtr_warehouse_location.warehouse_id,
			SUM(ISNULL(view_stock.quantity_on_hand, 0)) AS quantity_on_hand
			`).
		Joins("INNER JOIN mtr_warehouse_master mwm ON mwm.warehouse_id = mtr_warehouse_location.warehouse_id").
		Joins(`LEFT JOIN (?) view_stock ON view_stock.period_year = ?
										AND view_stock.period_month = ?
										AND view_stock.company_id = mwm.company_id
										AND view_stock.warehouse_group_id = mtr_warehouse_location.warehouse_group_id
										AND view_stock.warehouse_id = mtr_warehouse_location.warehouse_id
										AND view_stock.location_id = mtr_warehouse_location.warehouse_location_id
										`, viewLocStock, periodResponse.PeriodYear, periodResponse.PeriodMonth).
		Where("warehouse_location_id NOT IN ?", existingWarehouseLocIds).
		Where("mwm.company_id = ?", companyId).
		Where("mwm.warehouse_id = ?", warehouseId).
		Group(`
			mtr_warehouse_location.is_active,
			mtr_warehouse_location.warehouse_location_id,
			mtr_warehouse_location.warehouse_location_code,
			mtr_warehouse_location.warehouse_location_name,
			mwm.company_id,
			mtr_warehouse_location.warehouse_group_id,
			mtr_warehouse_location.warehouse_id
		`).
		Order("warehouse_location_id ASC")

	whereQuery := utils.ApplyFilter(baseModelQuery, newFilter)
	err = whereQuery.Scan(&response).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fethching location available data",
			Err:        err,
		}
	}

	paginateData, totalPages, totalRows := pagination.NewDataFramePaginate(response, &pages)

	pages.Rows = utils.ModifyKeysInResponse(paginateData)
	pages.TotalPages = totalPages
	pages.TotalRows = int64(totalRows)

	return pages, nil
}

// uspg_gmItem1_Select
// IF @Option = 5
func (r *LookupRepositoryImpl) ItemDetailForItemInquiry(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemDetail{}
	payloads := []masterpayloads.ItemDetailForItemInquiryPayload{}

	baseModelQuery := tx.Model(&entities).Select("model_id, variant_id")
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&payloads).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching item detail for item inquiry",
			Err:        err,
		}
	}

	finalJoinedData := []map[string]interface{}{}
	if len(payloads) > 0 {
		var modelIds []int
		var variantIds []int

		for _, data := range payloads {
			modelIds = append(modelIds, data.ModelId)
			variantIds = append(variantIds, data.VariantId)
		}

		modelResponse, modelError := salesserviceapiutils.GetUnitModelByMultiId(modelIds)
		if modelError != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error fetching model data",
				Err:        modelError.Err,
			}
		}

		variantResponse, variantError := salesserviceapiutils.GetUnitVariantByMultiId(variantIds)
		if variantError != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error fetching model data",
				Err:        variantError.Err,
			}
		}

		joinedData := utils.DataFrameLeftJoin(payloads, modelResponse, "ModelId")

		finalJoinedData = utils.DataFrameLeftJoin(joinedData, variantResponse, "VariantId")
	}

	pages.Rows = utils.ModifyKeysInResponse(finalJoinedData)

	return pages, nil
}

// USPG_SMSUBSTITUTE1_SELECT
// IF @Option=3
func (r *LookupRepositoryImpl) ItemSubstituteDetailForItemInquiry(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var err error
	entities := masteritementities.ItemSubstitute{}
	entitiesDetail := masteritementities.ItemSubstituteDetail{}

	var itemId int
	var companyId int
	for _, filter := range filterCondition {
		if strings.Contains(filter.ColumnField, "item_id") {
			itemId, _ = strconv.Atoi(filter.ColumnValue)
			continue
		}
		if strings.Contains(filter.ColumnField, "company_id") {
			companyId, _ = strconv.Atoi(filter.ColumnValue)
			continue
		}
	}

	currentYear := strconv.Itoa(int(time.Now().Year()))
	currentMonth := strconv.Itoa(int(time.Now().Month()))
	if len(currentMonth) == 1 {
		currentMonth = "0" + currentMonth
	}

	effectiveDate := time.Now().Truncate(24 * time.Hour)
	err = tx.Model(&entities).
		Select("MAX(effective_date)").
		Where(masteritementities.ItemSubstitute{ItemId: itemId, IsActive: true}).
		Pluck("effective_date", &effectiveDate).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error when fetching item substitute effective date data",
			Err:        err,
		}
	}

	viewLocStock := tx.Table("mtr_location_stock mls").
		Select(`
			mls.item_inquiry_id,
			mls.period_year,
			mls.period_month,
			mls.company_id,
			mls.item_id,
			(
				ISNULL(quantity_sales, 0) +
				ISNULL(quantity_transfer_out, 0) +
				ISNULL(quantity_robbing_out, 0) +
				ISNULL(quantity_assembly_out, 0) +
				ISNULL(quantity_allocated, 0)
			) AS quantity_available
		`).
		Joins("LEFT JOIN mtr_warehouse_master mwm ON mwm.company_id = mls.company_id AND mwm.warehouse_id = mls.warehouse_id")

	query1 := tx.Model(&entitiesDetail).
		Select(`
			mtr_item_substitute_detail.item_substitute_detail_id,
			mtr_item_substitute_detail.is_active,
			mtr_item_substitute_detail.item_id,
			mi.item_name,
			SUM(view_location.quantity_available) AS quantity,
			mtr_item_substitute_detail.sequence
		`).
		Joins("INNER JOIN mtr_item_substitute mis ON mis.item_substitute_id = mtr_item_substitute_detail.item_substitute_id").
		Joins("LEFT JOIN mtr_item mi ON mi.item_id = mtr_item_substitute_detail.item_id").
		Joins("INNER JOIN (?) view_location ON view_location.item_id = mtr_item_substitute_detail.item_id", viewLocStock).
		Where("mtr_item_substitute_detail.item_id = ?", itemId).
		Where("view_location.period_year = ?", currentYear).
		Where("view_location.period_month = ?", currentMonth).
		Where("view_location.company_id = ?", companyId).
		Where("mis.effective_date = ?", effectiveDate).
		Group(`
			mtr_item_substitute_detail.item_substitute_detail_id,
			mtr_item_substitute_detail.is_active,
			mtr_item_substitute_detail.item_id,
			mi.item_name,
			mtr_item_substitute_detail.sequence
		`)

	response1 := []masterpayloads.ItemSubstituteDetailForItemInquiryResponse{}
	err = query1.Scan(&response1).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error when fetching item substitute detail id first query",
			Err:        err,
		}
	}

	query1Ids := []int{-1} // set default value to prevent error
	for _, resp := range response1 {
		query1Ids = append(query1Ids, resp.ItemSubstituteDetailId)
	}

	query2 := tx.Model(&entitiesDetail).
		Select(`
			mtr_item_substitute_detail.item_substitute_detail_id,
			mtr_item_substitute_detail.is_active,
			mtr_item_substitute_detail.item_id,
			mi.item_name,
			0 AS quantity,
			mtr_item_substitute_detail.sequence
		`).
		Joins("INNER JOIN mtr_item_substitute mis ON mis.item_substitute_id = mtr_item_substitute_detail.item_substitute_id").
		Joins("LEFT JOIN mtr_item mi ON mi.item_id = mtr_item_substitute_detail.item_id").
		Where("mtr_item_substitute_detail.item_id = ?", itemId).
		Where("mis.effective_date = ?", effectiveDate).
		Where("mtr_item_substitute_detail.item_substitute_id NOT IN ?", query1Ids)

	response2 := []masterpayloads.ItemSubstituteDetailForItemInquiryResponse{}
	err = query2.Scan(&response2).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error when fetching item substitute detail id second query",
			Err:        err,
		}
	}

	response1 = append(response1, response2...)

	joinedResponse, totalPages, totalRows := pagination.NewDataFramePaginate(response1, &pages)

	pages.Rows = utils.ModifyKeysInResponse(joinedResponse)
	pages.TotalPages = totalPages
	pages.TotalRows = int64(totalRows)

	return pages, nil
}

// new req used for part number lookup in item import
func (r *LookupRepositoryImpl) GetPartNumberItemImport(tx *gorm.DB, internalCondition []utils.FilterCondition, externalCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Item{}
	response := []masterpayloads.GetPartNumberItemImportResponse{}

	var supplierCode string
	var supplierName string
	for _, filter := range externalCondition {
		if strings.Contains(filter.ColumnField, "supplier_code") {
			supplierCode = filter.ColumnValue
		}
		if strings.Contains(filter.ColumnField, "supplier_name") {
			supplierName = filter.ColumnValue
		}
	}

	baseModelQuery := tx.Model(&entities).
		Joins("INNER JOIN mtr_item_group mig ON mig.item_group_id = mtr_item.item_group_id").
		Joins("INNER JOIN mtr_item_type mit ON mit.item_type_id = mtr_item.item_type_id").
		Where("mig.item_group_code = 'IN' AND mit.item_type_code = 'G'")

	if supplierCode != "" || supplierName != "" {
		supplierParams := generalserviceapiutils.SupplierMasterParams{
			Page:         0,
			Limit:        1000,
			SupplierCode: supplierCode,
			SupplierName: supplierName,
		}

		supplierResponse, supplierError := generalserviceapiutils.GetAllSupplierMaster(supplierParams)
		if supplierError != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: supplierError.StatusCode,
				Message:    "Error fetching supplier data",
				Err:        supplierError.Err,
			}
		}

		if len(supplierResponse) == 0 {
			internalCondition = append(internalCondition, utils.FilterCondition{
				ColumnField: "mtr_item.supplier_id",
				ColumnValue: "-1",
			})
		} else {
			var supplierIdFilter []int
			for _, supplier := range supplierResponse {
				supplierIdFilter = append(supplierIdFilter, supplier.SupplierId)
			}
			baseModelQuery = baseModelQuery.Where("mtr_item.supplier_id IN ?", supplierIdFilter)
		}
	}

	whereQuery := utils.ApplyFilter(baseModelQuery, internalCondition)
	err := whereQuery.Scopes(pagination.Paginate(&pages, baseModelQuery)).Scan(&response).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching item master data",
			Err:        err,
		}
	}

	supplierIds := []int{}
	for _, resp := range response {
		if resp.SupplierId > 0 {
			supplierIds = append(supplierIds, resp.SupplierId)
		}
	}

	supplierResponse := []masterpayloads.ItemImportSupplierResponse{}
	if len(supplierIds) > 0 {
		supplierError := generalserviceapiutils.GetSupplierMasterByMultiId(supplierIds, &supplierResponse)
		if supplierError != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "error fetching supplier data",
				Err:        supplierError.Err,
			}
		}
	}

	for i := 0; i < len(response); i++ {
		for j := 0; j < len(supplierResponse); j++ {
			if response[i].SupplierId == supplierResponse[j].SupplierId {
				response[i].SupplierCode = supplierResponse[j].SupplierCode
				response[i].SupplierName = supplierResponse[j].SupplierName
				break
			}
		}
	}

	pages.Rows = response

	return pages, nil
}

// usp_comLookUp
// IF @strEntity = 'LocationItem'
func (r *LookupRepositoryImpl) LocationItem(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemLocation{}
	response := []masterpayloads.LocationItemResponse{}

	baseModelQuery := tx.Model(&entities).
		Select(`
			MIN(mwl.warehouse_location_id) warehouse_location_id,
			mwl.warehouse_location_code,
			mwl.warehouse_location_name
		`).
		Joins("LEFT JOIN mtr_warehouse_location mwl ON mwl.warehouse_id = mtr_location_item.warehouse_id AND mwl.warehouse_location_id = mtr_location_item.warehouse_location_id").
		Joins("LEFT JOIN mtr_warehouse_master mwm ON mwm.warehouse_id = mtr_location_item.warehouse_id").
		Where("mwl.is_active = ?", true)
	whereCondition := utils.ApplyFilter(baseModelQuery, filterCondition).Group("mwl.warehouse_location_id, mwl.warehouse_location_code, mwl.warehouse_location_name")
	err := whereCondition.Scopes(pagination.Paginate(&pages, baseModelQuery)).Scan(&response).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching lookup data 'LocationItem'",
			Err:        err,
		}
	}

	pages.Rows = response

	return pages, nil
}

// usp_comLookUp
// IF @strEntity = 'ItemLocUOM'
func (r *LookupRepositoryImpl) ItemLocUOM(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var newFilterCondition []utils.FilterCondition
	var companyId int
	for _, filter := range filterCondition {
		if strings.Contains(filter.ColumnField, "company_id") {
			companyId, _ = strconv.Atoi(filter.ColumnValue)
			continue
		}
		newFilterCondition = append(newFilterCondition, filter)
	}

	periodResponse, periodError := financeserviceapiutils.GetOpenPeriodByCompany(companyId, "SP")
	if periodError != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching current period company",
			Err:        periodError.Err,
		}
	}

	var periodYear string = periodResponse.PeriodYear
	var periodMonth string = periodResponse.PeriodMonth

	entities := masteritementities.ItemLocation{}
	response := []masterpayloads.ItemLocUOMResponse{}

	viewLocStock := tx.Table("mtr_location_stock mls").
		Select(`
			mls.item_inquiry_id,
			mls.period_year,
			mls.period_month,
			mls.company_id,
			mls.item_id,
			mls.warehouse_group_id,
			mls.warehouse_id,
			mls.location_id,
			(
				ISNULL(quantity_sales, 0) +
				ISNULL(quantity_transfer_out, 0) +
				ISNULL(quantity_robbing_out, 0) +
				ISNULL(quantity_assembly_out, 0) +
				ISNULL(quantity_allocated, 0)
			) AS quantity_available
		`).
		Joins("LEFT JOIN mtr_warehouse_master mwm ON mwm.company_id = mls.company_id AND mwm.warehouse_id = mls.warehouse_id")

	baseModelQuery := tx.Model(&entities).
		Select(`
			MIN(mi.item_id) item_id,
			mi.item_code,
			mi.item_name,
			mu.uom_code,
			SUM(ISNULL(vls.quantity_available, 0)) quantity_available,
			mi.is_active
		`).
		Joins("INNER JOIN mtr_item mi ON mi.item_id = mtr_location_item.item_id").
		Joins("LEFT JOIN mtr_uom mu ON mu.uom_id = mi.unit_of_measurement_stock_id").
		Joins(`LEFT JOIN (?) vls ON vls.item_id = mi.item_id AND vls.warehouse_group_id = mtr_location_item.warehouse_group_id
															AND vls.warehouse_id = mtr_location_item.warehouse_id
															AND vls.location_id = mtr_location_item.warehouse_location_id
															AND vls.period_year = ?
															AND vls.period_month = ?`, viewLocStock, periodYear, periodMonth)
	whereQuery := utils.ApplyFilter(baseModelQuery, newFilterCondition).Group("mi.item_code, mi.item_name, mu.uom_code, mi.is_active")
	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&response).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching lookup data 'ItemLocUOM'",
			Err:        err,
		}
	}

	pages.Rows = response

	return pages, nil
}

func (r *LookupRepositoryImpl) ItemLocUOMById(tx *gorm.DB, companyId int, itemId int) (masterpayloads.ItemLocUOMResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemLocation{}
	response := masterpayloads.ItemLocUOMResponse{}

	periodResponse, periodError := financeserviceapiutils.GetOpenPeriodByCompany(companyId, "SP")
	if periodError != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching current period company",
			Err:        periodError.Err,
		}
	}

	var periodYear string = periodResponse.PeriodYear
	var periodMonth string = periodResponse.PeriodMonth

	viewLocStock := tx.Table("mtr_location_stock mls").
		Select(`
			mls.item_inquiry_id,
			mls.period_year,
			mls.period_month,
			mls.company_id,
			mls.item_id,
			mls.warehouse_group_id,
			mls.warehouse_id,
			mls.location_id,
			(
				ISNULL(quantity_sales, 0) +
				ISNULL(quantity_transfer_out, 0) +
				ISNULL(quantity_robbing_out, 0) +
				ISNULL(quantity_assembly_out, 0) +
				ISNULL(quantity_allocated, 0)
			) AS quantity_available
		`).
		Joins("LEFT JOIN mtr_warehouse_master mwm ON mwm.company_id = mls.company_id AND mwm.warehouse_id = mls.warehouse_id")

	baseModelQuery := tx.Model(&entities).
		Select(`TOP 1
			MIN(mi.item_id) item_id,
			mi.item_code,
			mi.item_name,
			mu.uom_code,
			SUM(ISNULL(vls.quantity_available, 0)) quantity_available,
			mi.is_active
		`).
		Joins("INNER JOIN mtr_item mi ON mi.item_id = mtr_location_item.item_id").
		Joins("LEFT JOIN mtr_uom mu ON mu.uom_id = mi.unit_of_measurement_stock_id").
		Joins(`LEFT JOIN (?) vls ON vls.item_id = mi.item_id AND vls.warehouse_group_id = mtr_location_item.warehouse_group_id
															AND vls.warehouse_id = mtr_location_item.warehouse_id
															AND vls.location_id = mtr_location_item.warehouse_location_id
															AND vls.period_year = ?
															AND vls.period_month = ?`, viewLocStock, periodYear, periodMonth).
		Group("mi.item_code, mi.item_name, mu.uom_code, mi.is_active").
		Order("item_id")

	err := baseModelQuery.Where("mi.item_id = ?", itemId).Scan(&response).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "data not found",
				Err:        err,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fethcing lookup data 'ItemLocUOMById'",
			Err:        err,
		}
	}

	return response, nil
}

func (r *LookupRepositoryImpl) ItemLocUOMByCode(tx *gorm.DB, companyId int, itemCode string) (masterpayloads.ItemLocUOMResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemLocation{}
	response := masterpayloads.ItemLocUOMResponse{}

	periodResponse, periodError := financeserviceapiutils.GetOpenPeriodByCompany(companyId, "SP")
	if periodError != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching current period company",
			Err:        periodError.Err,
		}
	}

	var periodYear string = periodResponse.PeriodYear
	var periodMonth string = periodResponse.PeriodMonth

	viewLocStock := tx.Table("mtr_location_stock mls").
		Select(`
			mls.item_inquiry_id,
			mls.period_year,
			mls.period_month,
			mls.company_id,
			mls.item_id,
			mls.warehouse_group_id,
			mls.warehouse_id,
			mls.location_id,
			(
				ISNULL(quantity_sales, 0) +
				ISNULL(quantity_transfer_out, 0) +
				ISNULL(quantity_robbing_out, 0) +
				ISNULL(quantity_assembly_out, 0) +
				ISNULL(quantity_allocated, 0)
			) AS quantity_available
		`).
		Joins("LEFT JOIN mtr_warehouse_master mwm ON mwm.company_id = mls.company_id AND mwm.warehouse_id = mls.warehouse_id")

	baseModelQuery := tx.Model(&entities).
		Select(`TOP 1
			MIN(mi.item_id) item_id,
			mi.item_code,
			mi.item_name,
			mu.uom_code,
			SUM(ISNULL(vls.quantity_available, 0)) quantity_available,
			mi.is_active
		`).
		Joins("INNER JOIN mtr_item mi ON mi.item_id = mtr_location_item.item_id").
		Joins("LEFT JOIN mtr_uom mu ON mu.uom_id = mi.unit_of_measurement_stock_id").
		Joins(`LEFT JOIN (?) vls ON vls.item_id = mi.item_id AND vls.warehouse_group_id = mtr_location_item.warehouse_group_id
															AND vls.warehouse_id = mtr_location_item.warehouse_id
															AND vls.location_id = mtr_location_item.warehouse_location_id
															AND vls.period_year = ?
															AND vls.period_month = ?`, viewLocStock, periodYear, periodMonth).
		Group("mi.item_code, mi.item_name, mu.uom_code, mi.is_active").
		Order("item_id")

	err := baseModelQuery.Where("mi.item_code = ?", itemCode).Scan(&response).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "data not found",
				Err:        err,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fethcing lookup data 'ItemLocUOMById'",
			Err:        err,
		}
	}

	return response, nil
}

// usp_comLookUp
// IF @strEntity = 'ItemMasterForFreeAccs'
func (r *LookupRepositoryImpl) ItemMasterForFreeAccs(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	response := []masterpayloads.ItemMasterForFreeAccsResponse{}

	baseModelQuery := tx.Table("mtr_item").
		Select(`mtr_item.item_id, mtr_item.item_code, mtr_item.item_name, 
                 mtr_item.item_class_id, mtr_item_class.item_class_code, mtr_item_class.item_class_name, 
                 mtr_uom.uom_id, mtr_uom.uom_code, mtr_uom.uom_description, 
                 mtr_item_price_list.price_list_amount as price, mtr_item.is_active`).
		Joins("INNER JOIN mtr_item_class ON mtr_item.item_class_id = mtr_item_class.item_class_id").
		Joins("INNER JOIN mtr_item_detail ON mtr_item.item_id = mtr_item_detail.item_id").
		Joins("INNER JOIN mtr_uom ON mtr_item.unit_of_measurement_stock_id = mtr_uom.uom_id").
		Joins("INNER JOIN mtr_item_price_list ON mtr_item.item_id = mtr_item_price_list.item_id " +
			"AND mtr_item_detail.brand_id = mtr_item_price_list.brand_id " +
			"AND mtr_item_price_list.price_list_code_id = 1")

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	paginatedQuery := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery))

	err := paginatedQuery.Scan(&response).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching lookup data 'ItemMasterForFreeAccs'",
			Err:        err,
		}
	}

	pages.Rows = response
	return pages, nil
}

// usp_comLookUp
// IF @strEntity = 'ItemMasterForFreeAccs'
func (r *LookupRepositoryImpl) ItemMasterForFreeAccsById(tx *gorm.DB, companyId int, itemId int) (masterpayloads.ItemMasterForFreeAccsResponse, *exceptions.BaseErrorResponse) {
	response := masterpayloads.ItemMasterForFreeAccsResponse{}

	baseModelQuery := tx.Table("mtr_item").
		Select(`mtr_item.item_id, mtr_item.item_code, mtr_item.item_name, 
				 mtr_item.item_class_id, mtr_item_class.item_class_code, mtr_item_class.item_class_name, 
				 mtr_uom.uom_id, mtr_uom.uom_code, mtr_uom.uom_description, 
				 mtr_item_price_list.price_list_amount as price, mtr_item.is_active`).
		Joins("INNER JOIN mtr_item_class ON mtr_item.item_class_id = mtr_item_class.item_class_id").
		Joins("INNER JOIN mtr_item_detail ON mtr_item.item_id = mtr_item_detail.item_id").
		Joins("INNER JOIN mtr_uom ON mtr_item.unit_of_measurement_stock_id = mtr_uom.uom_id").
		Joins("INNER JOIN mtr_item_price_list ON mtr_item.item_id = mtr_item_price_list.item_id "+
			"AND mtr_item_detail.brand_id = mtr_item_price_list.brand_id "+
			"AND mtr_item_price_list.price_list_code_id = 1").
		Where("mtr_item.item_id = ?", itemId)

	err := baseModelQuery.First(&response).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "data not found",
				Err:        err,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching lookup data 'ItemMasterForFreeAccsById'",
			Err:        err,
		}
	}

	return response, nil
}

// usp_comLookUp
// IF @strEntity = 'ItemMasterForFreeAccs'
func (r *LookupRepositoryImpl) ItemMasterForFreeAccsByCode(tx *gorm.DB, companyId int, itemCode string) (masterpayloads.ItemMasterForFreeAccsResponse, *exceptions.BaseErrorResponse) {
	response := masterpayloads.ItemMasterForFreeAccsResponse{}

	baseModelQuery := tx.Table("mtr_item").
		Select(`mtr_item.item_id, mtr_item.item_code, mtr_item.item_name, 
				 mtr_item.item_class_id, mtr_item_class.item_class_code, mtr_item_class.item_class_name, 
				 mtr_uom.uom_id, mtr_uom.uom_code, mtr_uom.uom_description, 
				 mtr_item_price_list.price_list_amount as price, mtr_item.is_active`).
		Joins("INNER JOIN mtr_item_class ON mtr_item.item_class_id = mtr_item_class.item_class_id").
		Joins("INNER JOIN mtr_item_detail ON mtr_item.item_id = mtr_item_detail.item_id").
		Joins("INNER JOIN mtr_uom ON mtr_item.unit_of_measurement_stock_id = mtr_uom.uom_id").
		Joins("INNER JOIN mtr_item_price_list ON mtr_item.item_id = mtr_item_price_list.item_id "+
			"AND mtr_item_detail.brand_id = mtr_item_price_list.brand_id "+
			"AND mtr_item_price_list.price_list_code_id = 1").
		Where("mtr_item.item_code = ?", itemCode)

	err := baseModelQuery.First(&response).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "data not found",
				Err:        err,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching lookup data 'ItemMasterForFreeAccsByCode'",
			Err:        err,
		}
	}

	return response, nil
}

// usp_comLookUp
// IF @strEntity = 'ItemMasterForFreeAccs'
func (r *LookupRepositoryImpl) ItemMasterForFreeAccsByBrand(tx *gorm.DB, companyId int, itemId int, brandId int) (masterpayloads.ItemMasterForFreeAccsBrandResponse, *exceptions.BaseErrorResponse) {

	response := masterpayloads.ItemMasterForFreeAccsBrandResponse{}

	baseModelQuery := tx.Table("mtr_item").
		Select(`mtr_item.item_id, mtr_item.item_code, mtr_item.item_name, 
				 mtr_item.item_class_id, mtr_item_class.item_class_code, mtr_item_class.item_class_name, 
				 mtr_uom.uom_id, mtr_uom.uom_code, mtr_uom.uom_description, 
				 mtr_item_price_list.price_list_amount as price, mtr_item.is_active`).
		Joins("INNER JOIN mtr_item_class ON mtr_item.item_class_id = mtr_item_class.item_class_id").
		Joins("INNER JOIN mtr_item_detail ON mtr_item.item_id = mtr_item_detail.item_id").
		Joins("INNER JOIN mtr_uom ON mtr_item.unit_of_measurement_stock_id = mtr_uom.uom_id").
		Joins("INNER JOIN mtr_item_price_list ON mtr_item.item_id = mtr_item_price_list.item_id "+
			"AND mtr_item_detail.brand_id = mtr_item_price_list.brand_id "+
			"AND mtr_item_price_list.price_list_code_id = 1").
		Where("mtr_item_price_list.brand_id = ?", brandId).
		Where("mtr_item.item_id = ?", itemId).
		Where("mtr_item_price_list.company_id = ?", companyId)

	err := baseModelQuery.First(&response).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "data not found",
				Err:        err,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching lookup data 'ItemMasterForFreeAccsByBrand'",
			Err:        err,
		}
	}

	return response, nil
}
