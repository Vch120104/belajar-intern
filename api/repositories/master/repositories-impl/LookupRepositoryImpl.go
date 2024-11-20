package masterrepositoryimpl

import (
	"after-sales/api/config"
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

// dbo.getOprItemDisc
// get DISCOUNT value base on line type in operation or item master
func (r *LookupRepositoryImpl) GetOprItemDisc(tx *gorm.DB, lineTypeId int, billCodeId int, oprItemCode int, agreementId int, profitCenterId int, minValue float64, companyId int, brandId int, contractServSysNo int, whsGroup int, orderTypeId int) (float64, *exceptions.BaseErrorResponse) {
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
				if lineTypeId != utils.LinetypeOperation && lineTypeId != utils.LinetypePackage {
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
						// Check Agreement3
						err = tx.Model(&masterentities.AgreementItemDetail{}).
							Where("agreement_id = ? AND line_type_id = ? AND agreement_item_operation_id = ? AND min_value <= ?", agreementId, lineTypeId, oprItemCode, minValue).
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
						// Check Agreement1
						err = tx.Model(&masterentities.AgreementDiscount{}).
							Where("agreement_id = ? AND line_type_id = ? AND min_value <= ?", agreementId, lineTypeId, minValue).
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

				if lineTypeId == utils.LinetypeOperation || lineTypeId == utils.LinetypePackage {
					if discount == 0 {
						// Check Agreement3 for Operations
						err = tx.Model(&masterentities.AgreementItemDetail{}).
							Where("agreement_id = ? AND line_type_id = ? AND agreement_item_operation_id = ? AND min_value <= ?", agreementId, lineTypeId, oprItemCode, minValue).
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
						// Check Agreement1 for Operations
						err = tx.Model(&masterentities.AgreementDiscount{}).
							Where("agreement_id = ? AND line_type_id = ? AND min_value <= ?", agreementId, lineTypeId, minValue).
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

		if lineTypeId != utils.LinetypeOperation && lineTypeId != utils.LinetypePackage {

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
		err = tx.Model(&transactionworkshopentities.ContractService{}).
			Joins("INNER JOIN trx_contract_service_Operation_detail ON trx_contract_service_Operation_detail.contract_service_system_number = trx_contract_service.contract_service_system_number").
			Where("trx_contract_service.contract_service_system_number = ? AND trx_contract_service_Operation_detail.line_type_id = ? AND trx_contract_service_Operation_detail.operation_id = ?", contractServSysNo, lineTypeId, oprItemCode).
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

	priceListCodeUrl := config.EnvConfigs.GeneralServiceUrl + "price-list-code-by-code/A"
	preiceListCodePayloads := masterpayloads.GetPriceListCodeResponse{}
	if err := utils.Get(priceListCodeUrl, &preiceListCodePayloads, nil); err != nil || preiceListCodePayloads.PriceListCodeId == 0 {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching price list code: A",
			Err:        errors.New("error fetching default price list code"),
		}
	}
	defaultPriceCodeId = preiceListCodePayloads.PriceListCodeId

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
	case utils.LinetypePackage:
		// Package price logic
		if err := tx.Model(&masterentities.PackageMaster{}).
			Where("package_code = ?", oprItemCode).
			Select("package_price").
			Scan(&price).Error; err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get package price",
				Err:        err,
			}
		}

	case utils.LinetypeOperation:
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

			currentPeriodUrl := config.EnvConfigs.FinanceServiceUrl + "closing-period-company/current-period?company_id=" + strconv.Itoa(companyId) + "&closing_module_detail_code" + moduleSP
			currentPeriodPayloads := masterpayloads.GetCurrentPeriodResponse{}
			if err := utils.Get(currentPeriodUrl, &currentPeriodPayloads, nil); err != nil {
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
					Where("period_year = ? AND period_month = ? AND item_code = ? AND company_id = ? AND whs_group = ?",
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
				Where("is_active = 1 AND brand_id = ? AND effective_date <= ? AND item_code = ? AND currency_id = ? AND company_id = ? AND price_list_code_id = ?",
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
	if linetypeId == utils.LinetypeOperation && billCode == utils.TrxTypeWoInternal.ID {
		price += price * markupPercentage / 100
	}

	return price, nil
}

// usp_comLookUp
// IF @strEntity = 'ItemOprCode'--OPERATION MASTER & ITEM MASTER
func (r *LookupRepositoryImpl) ItemOprCode(tx *gorm.DB, linetypeId int, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var results []map[string]interface{}

	// Default filters and variables
	const (
		ItmGrpInventory      = 2 // "IN"
		PurchaseTypeGoods    = 1 // "G"
		PurchaseTypeServices = 2 // "S"
	)

	var (
		ItmCls      int
		companyCode = 1
		currentTime = time.Now()
		year, month = currentTime.Year(), int(currentTime.Month())
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	// Filter Handling
	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	// Base Query
	baseQuery := tx.Table("")

	switch linetypeId {
	case utils.LinetypePackage:
		combinedDetailsSubQuery := `
			(
				SELECT package_id, frt_quantity, is_active 
				FROM mtr_package_master_detail
				WHERE is_active = 1
				UNION ALL
				SELECT package_id, frt_quantity, is_active 
				FROM mtr_package_master_detail
				WHERE is_active = 1
			) AS CombinedDetails
		`

		baseQuery = baseQuery.Table("mtr_package A").
			Select(`
				A.package_id AS package_id, 
				A.package_code AS package_code, 
				A.package_name AS package_name, 
				SUM(CombinedDetails.frt_quantity) AS frt, 
				B.profit_center_id AS profit_center, 
				C.model_code AS model_code, 
				C.model_description AS description, 
				A.package_price AS price
			`).
			Joins("LEFT JOIN "+combinedDetailsSubQuery+" ON A.package_id = CombinedDetails.package_id").
			Joins("LEFT JOIN dms_microservices_general_dev.dbo.mtr_profit_center B ON A.profit_center_id = B.profit_center_id").
			Joins("LEFT JOIN dms_microservices_sales_dev.dbo.mtr_unit_model C ON A.model_id = C.model_id").
			Where("A.is_active = ?", 1).
			Where(filterQuery, filterValues...).
			Group("A.package_id ,A.package_code, A.package_name, B.profit_center_id, C.model_code, C.model_description, A.package_price")

	case utils.LinetypeOperation:
		baseQuery = baseQuery.Table("dms_microservices_aftersales_dev.dbo.mtr_operation_code AS oc").
			Select(`
				oc.operation_id AS operation_id, 
				oc.operation_code AS operation_code, 
				oc.operation_name AS operation_name, 
				ofrt.frt_hour AS frt_hour, 
				oe.operation_entries_code AS operation_entries_code, 
				oe.operation_entries_description AS operation_entries_description, 
				ok.operation_key_code AS operation_key_code, 
				ok.operation_key_description AS operation_key_description
			`).
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_entries AS oe ON oc.operation_entries_id = oe.operation_entries_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_key AS ok ON oc.operation_key_id = ok.operation_key_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_model_mapping AS omm ON oc.operation_id = omm.operation_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_frt AS ofrt ON omm.operation_model_mapping_id = ofrt.operation_model_mapping_id").
			Where("oc.is_active = ?", 1).
			Where(filterQuery, filterValues...)

	case utils.LinetypeSparepart:
		ItmCls = 69 // "SP"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) 
				        FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? 
				        AND V.PERIOD_MONTH = ? 
				        AND V.company_id = ?), 0) AS Available_qty, 
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
			`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, 1).
			Where(filterQuery, filterValues...)

	case utils.LinetypeOil:
		ItmCls = 70 // "OL"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) 
				        FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? 
				        AND V.PERIOD_MONTH = ? 
				        AND V.company_id = ?), 0) AS Available_qty, 
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
			`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, 1).
			Where(filterQuery, filterValues...)

	case utils.LinetypeMaterial:
		ItmCls = 71        // "MT"
		ItmClsSublet := 72 // "SB"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name, 
					ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
							WHERE A.item_id = V.item_id 
							AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
				`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND (A.item_class_id = ? OR A.item_class_id = ?) AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, ItmClsSublet, 1).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	case utils.LinetypeConsumableMaterial:
		ItmCls = 75 // "CM"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name, 
					ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
							WHERE A.item_id = V.item_id 
							AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
				`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ?  AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, 1).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	case utils.LinetypeFee:
		ItmCls = 73           // "WF"
		ItmGrpOutsideJob := 6 // "OJ"

		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name, 
					ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
							WHERE A.item_id = V.item_id 
							AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
				`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("(A.item_group_id = ? OR A.item_group_id = ?) AND A.item_class_id = ? AND A.item_type_id = ? AND A.is_active = ?", ItmGrpOutsideJob, ItmGrpInventory, ItmCls, PurchaseTypeServices, 1).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	case utils.LinetypeAccesories:
		ItmCls = 74 // "AC"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name, 
					ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
							WHERE A.item_id = V.item_id 
							AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
				`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", ItmCls, ItmGrpInventory, 1).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	case utils.LinetypeSouvenir:
		ItmCls = 77 // "SV"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
						A.item_id AS item_id, 
						A.item_code AS item_code, 
						A.item_name AS item_name, 
						ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
								WHERE A.item_id = V.item_id 
								AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
						A.item_level_1_id AS item_level_1,
						mil1.item_level_1_code AS item_level_1_code, 
						A.item_level_2_id AS item_level_2,
						mil2.item_level_2_code AS item_level_2_code, 
						A.item_level_3_id AS item_level_3,
						mil3.item_level_3_code AS item_level_3_code, 
						A.item_level_4_id AS item_level_4,
						mil4.item_level_4_code AS item_level_4_code
					`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", ItmCls, ItmGrpInventory, 1).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	default:
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid linetype ID",
			Err:        errors.New("invalid linetype ID"),
		}
	}

	// Total rows
	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count rows",
			Err:        err,
		}
	}

	// Pagination
	offset := (paginate.Page - 1) * paginate.Limit
	if err := baseQuery.Offset(offset).Limit(paginate.Limit).Find(&results).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve data",
			Err:        err,
		}
	}
	fmt.Println(baseQuery.Statement.SQL.String()) // Prints SQL statement to check filter and pagination application
	fmt.Printf("Offset: %d, Limit: %d\n", offset, paginate.Limit)

	totalPages := int(math.Ceil(float64(totalRows) / float64(paginate.Limit)))

	return results, int(totalRows), totalPages, nil
}

// usp_comLookUp
// IF @strEntity = 'ItemOprCode'--OPERATION MASTER & ITEM MASTER
func (r *LookupRepositoryImpl) ItemOprCodeByCode(tx *gorm.DB, linetypeId int, oprItemCode string, paginate pagination.Pagination, filters []utils.FilterCondition) (map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var results map[string]interface{}

	// Default filters and variables
	const (
		ItmGrpInventory      = 2 // "IN"
		PurchaseTypeGoods    = 1 // "G"
		PurchaseTypeServices = 2 // "S"
	)

	var (
		ItmCls      int
		companyCode = 1
		currentTime = time.Now()
		year, month = currentTime.Year(), int(currentTime.Month())
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	if oprItemCode == "" {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "oprItemCode cannot be empty",
			Err:        nil,
		}
	}

	baseQuery := tx

	if len(filters) > 0 {
		filterStrings := []string{}
		filterValues := []interface{}{}
		for _, filter := range filters {
			if filter.ColumnField != "" && filter.ColumnValue != "" {
				filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
				filterValues = append(filterValues, filter.ColumnValue)
			}
		}
		filterQuery := strings.Join(filterStrings, " AND ")
		baseQuery = baseQuery.Where(filterQuery, filterValues...)
	}
	switch linetypeId {
	case utils.LinetypePackage:
		combinedDetailsSubQuery := `
				(
					SELECT package_id, frt_quantity, is_active 
					FROM mtr_package_master_detail
					WHERE is_active = 1
					UNION ALL
					SELECT package_id, frt_quantity, is_active 
					FROM mtr_package_master_detail
					WHERE is_active = 1
				) AS CombinedDetails
			`

		baseQuery = baseQuery.Table("mtr_package A").
			Select(`
				A.package_id AS package_id,
				A.package_code AS package_code, 
				A.package_name AS package_name, 
				SUM(CombinedDetails.frt_quantity) AS frt, 
				B.profit_center_id AS profit_center, 
				C.model_code AS model_code, 
				C.model_description AS description, 
				A.package_price AS price
			`).
			Joins("LEFT JOIN "+combinedDetailsSubQuery+" ON A.package_id = CombinedDetails.package_id").
			Joins("LEFT JOIN dms_microservices_general_dev.dbo.mtr_profit_center B ON A.profit_center_id = B.profit_center_id").
			Joins("LEFT JOIN dms_microservices_sales_dev.dbo.mtr_unit_model C ON A.model_id = C.model_id").
			Where("A.is_active = ?", 1).
			Where("A.package_code = ?", oprItemCode).
			Group("A.package_id, A.package_code, A.package_name, B.profit_center_id, C.model_code, C.model_description, A.package_price")

	case utils.LinetypeOperation:
		baseQuery = baseQuery.Table("dms_microservices_aftersales_dev.dbo.mtr_operation_code AS oc").
			Select(`
					oc.operation_id AS operation_id, 
					oc.operation_code AS operation_code, 
					oc.operation_name AS operation_name, 
					ofrt.frt_hour AS frt_hour, 
					oe.operation_entries_code AS operation_entries_code, 
					oe.operation_entries_description AS operation_entries_description, 
					ok.operation_key_code AS operation_key_code, 
					ok.operation_key_description AS operation_key_description
				`).
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_entries AS oe ON oc.operation_entries_id = oe.operation_entries_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_key AS ok ON oc.operation_key_id = ok.operation_key_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_model_mapping AS omm ON oc.operation_id = omm.operation_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_frt AS ofrt ON omm.operation_model_mapping_id = ofrt.operation_model_mapping_id").
			Where("oc.is_active = ? ", 1).
			Where("oc.operation_code = ?", oprItemCode)

	case utils.LinetypeSparepart:
		ItmCls = 69 // "SP"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
					A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
			`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, 1).
			Where("A.item_code = ?", oprItemCode)

	case utils.LinetypeOil:
		ItmCls = 70 // "OL"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
				A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
			`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, 1).
			Where("A.item_code = ?", oprItemCode)

	case utils.LinetypeMaterial:
		ItmCls = 71        // "MT"
		ItmClsSublet := 72 // "SB"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
				A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
			`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND (A.item_class_id = ? OR A.item_class_id = ?) AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, ItmClsSublet, 1).
			Where("A.item_code = ?", oprItemCode).
			Order("A.item_code")

	case utils.LinetypeConsumableMaterial:
		ItmCls = 75 // "CM"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
						A.item_id AS item_id, 
						A.item_code AS item_code, 
						A.item_name AS item_name, 
						ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
								WHERE A.item_id = V.item_id 
								AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
						A.item_level_1_id AS item_level_1,
						mil1.item_level_1_code AS item_level_1_code, 
						A.item_level_2_id AS item_level_2,
						mil2.item_level_2_code AS item_level_2_code, 
						A.item_level_3_id AS item_level_3,
						mil3.item_level_3_code AS item_level_3_code, 
						A.item_level_4_id AS item_level_4,
						mil4.item_level_4_code AS item_level_4_code
					`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ?  AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, 1).
			Where("A.item_code = ?", oprItemCode).
			Order("A.item_code")

	case utils.LinetypeFee:
		ItmCls = 73           // "WF"
		ItmGrpOutsideJob := 6 // "OJ"

		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
				A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
			`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("(A.item_group_id = ? OR A.item_group_id = ?) AND A.item_class_id = ? AND A.item_type_id = ? AND A.is_active = ?", ItmGrpOutsideJob, ItmGrpInventory, ItmCls, PurchaseTypeServices, 1).
			Where("A.item_code = ?", oprItemCode).
			Order("A.item_code")

	case utils.LinetypeAccesories:
		ItmCls = 74 // "AC"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
				A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
			`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", ItmCls, ItmGrpInventory, 1).
			Where("A.item_code = ?", oprItemCode).
			Order("A.item_code")

	case utils.LinetypeSouvenir:
		ItmCls = 77 // "SV"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, A.item_code AS item_code, 
					A.item_name AS item_name, 
					ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
							WHERE A.item_id = V.item_id 
							AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
					A.item_level_1_id AS item_level_1,
						mil1.item_level_1_code AS item_level_1_code, 
						A.item_level_2_id AS item_level_2,
						mil2.item_level_2_code AS item_level_2_code, 
						A.item_level_3_id AS item_level_3,
						mil3.item_level_3_code AS item_level_3_code, 
						A.item_level_4_id AS item_level_4,
						mil4.item_level_4_code AS item_level_4_code
				`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", ItmCls, ItmGrpInventory, 1).
			Where("A.item_code = ?", oprItemCode).
			Order("A.item_code")
	default:
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid linetype ID",
			Err:        errors.New("invalid linetype ID"),
		}
	}

	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count rows",
			Err:        err,
		}
	}

	if totalRows == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No records found for the given item code",
			Err:        nil,
		}
	}

	offset := (paginate.Page - 1) * paginate.Limit
	if err := baseQuery.Offset(offset).Limit(paginate.Limit).Find(&results).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve data",
			Err:        err,
		}
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(paginate.Limit)))
	return results, int(totalRows), totalPages, nil
}

// usp_comLookUp
// IF @strEntity = 'ItemOprCode'--OPERATION MASTER & ITEM MASTER
func (r *LookupRepositoryImpl) ItemOprCodeByID(tx *gorm.DB, linetypeId int, oprItemId int, paginate pagination.Pagination, filters []utils.FilterCondition) (map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var results map[string]interface{}

	// Default filters and variables
	const (
		ItmGrpInventory      = 2 // "IN"
		PurchaseTypeGoods    = 1 // "G"
		PurchaseTypeServices = 2 // "S"
	)

	var (
		ItmCls      int
		companyCode = 1
		currentTime = time.Now()
		year, month = currentTime.Year(), int(currentTime.Month())
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	baseQuery := tx.Table("")

	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	switch linetypeId {
	case utils.LinetypePackage:
		combinedDetailsSubQuery := `
				(
					SELECT package_id, frt_quantity, is_active 
					FROM mtr_package_master_detail
					WHERE is_active = 1
					UNION ALL
					SELECT package_id, frt_quantity, is_active 
					FROM mtr_package_master_detail
					WHERE is_active = 1
				) AS CombinedDetails
			`

		baseQuery = baseQuery.Table("mtr_package A").
			Select(`
				A.package_id AS package_id,
				A.package_code AS package_code, 
				A.package_name AS package_name, 
				SUM(CombinedDetails.frt_quantity) AS frt, 
				B.profit_center_id AS profit_center, 
				C.model_code AS model_code, 
				C.model_description AS description, 
				A.package_price AS price
			`).
			Joins("LEFT JOIN "+combinedDetailsSubQuery+" ON A.package_id = CombinedDetails.package_id").
			Joins("LEFT JOIN dms_microservices_general_dev.dbo.mtr_profit_center B ON A.profit_center_id = B.profit_center_id").
			Joins("LEFT JOIN dms_microservices_sales_dev.dbo.mtr_unit_model C ON A.model_id = C.model_id").
			Where("A.is_active = ?", 1).
			Where("A.package_id = ?", oprItemId).
			Where(filterQuery, filterValues...).
			Group("A.package_id,A.package_code, A.package_name, B.profit_center_id, C.model_code, C.model_description, A.package_price")

	case utils.LinetypeOperation:
		baseQuery = baseQuery.Table("dms_microservices_aftersales_dev.dbo.mtr_operation_code AS oc").
			Select(`
			oc.operation_id AS operation_id, 
			oc.operation_code AS operation_code, 
			oc.operation_name AS operation_name, 
			ofrt.frt_hour AS frt_hour, 
			oe.operation_entries_code AS operation_entries_code, 
			oe.operation_entries_description AS operation_entries_description, 
			ok.operation_key_code AS operation_key_code, 
			ok.operation_key_description AS operation_key_description
		`).
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_entries AS oe ON oc.operation_entries_id = oe.operation_entries_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_key AS ok ON oc.operation_key_id = ok.operation_key_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_model_mapping AS omm ON oc.operation_id = omm.operation_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_frt AS ofrt ON omm.operation_model_mapping_id = ofrt.operation_model_mapping_id").
			Where("oc.is_active = ? ", 1).
			Where("oc.operation_id = ?", oprItemId).
			Where(filterQuery, filterValues...)

	case utils.LinetypeSparepart:
		ItmCls = 69 // "SP"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
				A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
			`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, 1).
			Where("A.item_id = ?", oprItemId).
			Where(filterQuery, filterValues...)

	case utils.LinetypeOil:
		ItmCls = 70 // "OL"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
				A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
			`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, 1).
			Where("A.item_id = ?", oprItemId).
			Where(filterQuery, filterValues...)

	case utils.LinetypeMaterial:
		ItmCls = 71        // "MT"
		ItmClsSublet := 72 // "SB"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
				A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
			`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND (A.item_class_id = ? OR A.item_class_id = ?) AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, ItmClsSublet, 1).
			Where("A.item_id = ?", oprItemId).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	case utils.LinetypeConsumableMaterial:
		ItmCls = 75 // "CM"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name, 
					ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
							WHERE A.item_id = V.item_id 
							AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
					A.item_level_1_id AS item_level_1,
						mil1.item_level_1_code AS item_level_1_code, 
						A.item_level_2_id AS item_level_2,
						mil2.item_level_2_code AS item_level_2_code, 
						A.item_level_3_id AS item_level_3,
						mil3.item_level_3_code AS item_level_3_code, 
						A.item_level_4_id AS item_level_4,
						mil4.item_level_4_code AS item_level_4_code
				`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, 1).
			Where("A.item_id = ?", oprItemId).
			Where(filterQuery, filterValues...)

	case utils.LinetypeFee:
		ItmCls = 73           // "WF"
		ItmGrpOutsideJob := 6 // "OJ"

		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
				A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
			`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("(A.item_group_id = ? OR A.item_group_id = ?) AND A.item_class_id = ? AND A.item_type_id = ? AND A.is_active = ?", ItmGrpOutsideJob, ItmGrpInventory, ItmCls, PurchaseTypeServices, 1).
			Where("A.item_id = ?", oprItemId).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	case utils.LinetypeAccesories:
		ItmCls = 74 // "AC"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_id AS item_id, 
				A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
				A.item_level_1_id AS item_level_1,
					mil1.item_level_1_code AS item_level_1_code, 
					A.item_level_2_id AS item_level_2,
					mil2.item_level_2_code AS item_level_2_code, 
					A.item_level_3_id AS item_level_3,
					mil3.item_level_3_code AS item_level_3_code, 
					A.item_level_4_id AS item_level_4,
					mil4.item_level_4_code AS item_level_4_code
			`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", ItmCls, ItmGrpInventory, 1).
			Where("A.item_id = ?", oprItemId).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	case utils.LinetypeSouvenir:
		ItmCls = 77 // "SV"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
					A.item_id AS item_id, 
					A.item_code AS item_code, 
					A.item_name AS item_name, 
					ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
							WHERE A.item_id = V.item_id 
							AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
					A.item_level_1_id AS item_level_1,
						mil1.item_level_1_code AS item_level_1_code, 
						A.item_level_2_id AS item_level_2,
						mil2.item_level_2_code AS item_level_2_code, 
						A.item_level_3_id AS item_level_3,
						mil3.item_level_3_code AS item_level_3_code, 
						A.item_level_4_id AS item_level_4,
						mil4.item_level_4_code AS item_level_4_code
				`, year, month, companyCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", ItmCls, ItmGrpInventory, 1).
			Where("A.item_id = ?", oprItemId).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	default:
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid linetype ID",
			Err:        errors.New("invalid linetype ID"),
		}
	}

	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count rows",
			Err:        err,
		}
	}

	offset := (paginate.Page - 1) * paginate.Limit
	if err := baseQuery.Offset(offset).Limit(paginate.Limit).Find(&results).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve data",
			Err:        err,
		}
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(paginate.Limit)))

	return results, int(totalRows), totalPages, nil
}

// usp_comLookUp
// IF @strEntity = 'ItemOprCodeWithPrice'--OPERATION MASTER & ITEM MASTER WITH PRICELIST
func (r *LookupRepositoryImpl) ItemOprCodeWithPrice(tx *gorm.DB, linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, jobTypeId int, variantId int, currencyId int, billCode int, whsGroup string, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var results []map[string]interface{}

	const (
		ItmGrpInventory   = 1 // "IN"
		PurchaseTypeGoods = "G"
	)

	type Period struct {
		PeriodYear  string `gorm:"column:PERIOD_YEAR"`
		PeriodMonth string `gorm:"column:PERIOD_MONTH"`
	}

	var (
		defaultPriceCode = "A"
		ItmCls           string
		year, month      string
		period           Period
		companyCode      = 151
		closingModul     = 10
		yearNow          = time.Now().Format("2006")
		monthNow         = time.Now().Format("01")
	)

	result := tx.Table("dms_microservices_finance_dev.dbo.mtr_closing_period_company").
		Select("TOP 1 period_year, period_month").
		Where("company_id = ? AND closing_module_detail_id = ? AND period_year <= ? AND period_month <= ? AND is_period_closed = '0'", companyCode, closingModul, yearNow, monthNow).
		Order("period_year DESC, period_month DESC").
		Scan(&period)

	if result.Error != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get period , please check closing period company",
			Err:        result.Error,
		}
	}

	year = period.PeriodYear
	month = period.PeriodMonth

	fmt.Println("Period Year:", year)
	fmt.Println("Period Month:", month)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	baseQuery := tx.Table("")

	filterStrings := []string{}
	filterValues := []interface{}{}
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
		filterValues = append(filterValues, filter.ColumnValue)
	}
	filterQuery := strings.Join(filterStrings, " AND ")

	price, err := r.GetOprItemPrice(tx, linetypeId, companyId, oprItemCode, brandId, modelId, jobTypeId, variantId, currencyId, billCode, whsGroup)
	if err != nil {
		return nil, 0, 0, err
	}

	switch linetypeId {
	case utils.LinetypePackage:
		combinedDetailsSubQuery := `
				(
					SELECT package_id, frt_quantity, is_active 
					FROM mtr_package_master_detail
					WHERE is_active = 1
					UNION ALL
					SELECT package_id, frt_quantity, is_active 
					FROM mtr_package_master_detail
					WHERE is_active = 1
				) AS CombinedDetails
			`

		baseQuery = baseQuery.Table("mtr_package A").
			Select(`
				A.package_code AS package_code, 
				A.package_name AS package_name, 
				SUM(CombinedDetails.frt_quantity) AS frt, 
				B.profit_center_id AS profit_center, 
				C.model_code AS model_code, 
				C.model_description AS description, 
				A.package_price AS price
			`).
			Joins("LEFT JOIN "+combinedDetailsSubQuery+" ON A.package_id = CombinedDetails.package_id").
			Joins("LEFT JOIN dms_microservices_general_dev.dbo.mtr_profit_center B ON A.profit_center_id = B.profit_center_id").
			Joins("LEFT JOIN dms_microservices_sales_dev.dbo.mtr_unit_model C ON A.model_id = C.model_id").
			Where("A.is_active = ?", 1).
			Where(filterQuery, filterValues...).
			Group("A.package_code, A.package_name, B.profit_center_id, C.model_code, C.model_description, A.package_price")

	case utils.LinetypeOperation:
		baseQuery = baseQuery.Table("dms_microservices_aftersales_dev.dbo.mtr_operation_code AS oc").
			Select(`
        oc.operation_code AS operation_code, 
        oc.operation_name AS operation_name, 
        ofrt.frt_hour AS frt_hour, 
        oe.operation_entries_code AS operation_entries_code, 
        oe.operation_entries_description AS operation_entries_description, 
        ok.operation_key_code AS operation_key_code, 
        ok.operation_key_description AS operation_key_description,
		? as PRICE
    `, price).
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_entries AS oe ON oc.operation_entries_id = oe.operation_entries_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_key AS ok ON oc.operation_key_id = ok.operation_key_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_model_mapping AS omm ON oc.operation_id = omm.operation_id").
			Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.mtr_operation_frt AS ofrt ON omm.operation_model_mapping_id = ofrt.operation_model_mapping_id").
			Where("oc.is_active = ? ", 1).
			Where(filterQuery, filterValues...)

	case utils.LinetypeSparepart:
		ItmCls = "1"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_code AS item_code,
				A.item_name AS item_name,
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
						WHERE A.item_id = V.item_id 
						AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ? 
						AND V.whs_group = ?), 0) AS Available_qty,
				A.item_level_1_id AS item_level_1,
				A.item_level_2_id AS item_level_2,
				A.item_level_3_id AS item_level_3,
				A.item_level_4_id AS item_level_4,
				CASE 
					WHEN ? IN (?, ?, ?) THEN
						CASE A.item_type_id
							WHEN ? THEN
								(SELECT TOP 1 price_list_amount FROM mtr_item_price_list
								WHERE is_active = 1 
								AND brand_id = B.brand_id 
								AND effective_date <= GETDATE()
								AND item_id = A.item_id
								AND currency_id = ? 
								AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
								AND price_list_code_id = ?
								ORDER BY effective_date DESC)
							ELSE
								(SELECT CASE ISNULL(price_current, 0)
										WHEN 0 THEN price_begin 
										ELSE price_current END AS HPP
								FROM mtr_group_stock 
								WHERE period_year = ? 
								AND period_month = ? 
								AND item_code = A.item_code 
								AND company_id = ?  
								AND whs_group = ?)
						END
					ELSE
						(SELECT TOP 1 price_list_amount FROM mtr_item_price_list
						WHERE is_active = 1 
						AND brand_id = B.brand_id 
						AND effective_date <= GETDATE()
						AND item_id = A.item_id 
						AND currency_id = ? 
						AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
						AND price_list_code_id = ?
						ORDER BY effective_date DESC)
				END AS PRICE
			`, year, month, companyId, whsGroup, // Parameters for AvailQty subquery
				billCode, utils.TrxTypeWoCentralize.Code, utils.TrxTypeWoInternal.Code, utils.TrxTypeWoNoCharge.Code, // Parameters for CASE statement
				2, currencyId, companyId, defaultPriceCode, // Parameters for subquery in CASE
				year, month, companyId, whsGroup, // Parameters for ELSE subquery in CASE
				currencyId, companyId, defaultPriceCode). // Parameters for the final ELSE condition.
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, 1).
			Where(filterQuery, filterValues...)

	case utils.LinetypeOil:
		ItmCls = "2"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
				A.item_level_1_id AS item_level_1, 
				A.item_level_2_id AS item_level_2, 
				A.item_level_3_id AS item_level_3, 
				A.item_level_4_id AS item_level_4,
				CASE 
					WHEN ? IN (?, ?, ?) THEN
						CASE A.item_type_id
							WHEN ? THEN
								(SELECT TOP 1 price_list_amount FROM mtr_item_price_list
								WHERE is_active = 1 
								AND brand_id = B.brand_id 
								AND effective_date <= GETDATE()
								AND item_id = A.item_id
								AND currency_id = ? 
								AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
								AND price_list_code_id = ?
								ORDER BY effective_date DESC)
							ELSE
								(SELECT CASE ISNULL(price_current, 0)
										WHEN 0 THEN price_begin 
										ELSE price_current END AS HPP
								FROM mtr_group_stock 
								WHERE period_year = ? 
								AND period_month = ? 
								AND item_code = A.item_code 
								AND company_id = ?  
								AND whs_group = ?)
						END
					ELSE
						(SELECT TOP 1 price_list_amount FROM mtr_item_price_list
						WHERE is_active = 1 
						AND brand_id = B.brand_id 
						AND effective_date <= GETDATE()
						AND item_id = A.item_id 
						AND currency_id = ? 
						AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
						AND price_list_code_id = ?
						ORDER BY effective_date DESC)
				END AS PRICE
			`, year, month, companyId, whsGroup, // Parameters for AvailQty subquery
				billCode, utils.TrxTypeWoCentralize.Code, utils.TrxTypeWoInternal.Code, utils.TrxTypeWoNoCharge.Code, // Parameters for CASE statement
				2, currencyId, companyId, defaultPriceCode, // Parameters for subquery in CASE
				year, month, companyId, whsGroup, // Parameters for ELSE subquery in CASE
				currencyId, companyId, defaultPriceCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND A.item_class_id = ? AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, 1).
			Where(filterQuery, filterValues...)

	case utils.LinetypeMaterial:
		ItmCls = "3"
		ItmClsSublet := "2"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				DISTINCT A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
				A.item_level_1_id AS item_level_1, 
				A.item_level_2_id AS item_level_2, 
				A.item_level_3_id AS item_level_3, 
				A.item_level_4_id AS item_level_4,
				CASE 
					WHEN ? IN (?, ?, ?) THEN
						CASE A.item_type
							WHEN ? THEN
								(SELECT TOP 1 price_list_amount FROM mtr_item_price_list
								WHERE is_active = 1 
								AND brand_id = B.brand_id 
								AND effective_date <= GETDATE()
								AND item_id = A.item_id
								AND currency_id = ? 
								AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
								AND price_list_code_id = ?
								ORDER BY effective_date DESC)
							ELSE
								(SELECT CASE ISNULL(price_current, 0)
										WHEN 0 THEN price_begin 
										ELSE price_current END AS HPP
								FROM mtr_group_stock 
								WHERE period_year = ? 
								AND period_month = ? 
								AND item_code = A.item_code 
								AND company_id = ?  
								AND whs_group = ?)
						END
					ELSE
						(SELECT TOP 1 price_list_amount FROM mtr_item_price_list
						WHERE is_active = 1 
						AND brand_id = B.brand_id 
						AND effective_date <= GETDATE()
						AND item_id = A.item_id 
						AND currency_id = ? 
						AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
						AND price_list_code_id = ?
						ORDER BY effective_date DESC)
				END AS PRICE
			`, year, month, companyId, whsGroup, // Parameters for AvailQty subquery
				billCode, utils.TrxTypeWoCentralize.Code, utils.TrxTypeWoInternal.Code, utils.TrxTypeWoNoCharge.Code, // Parameters for CASE statement
				2, currencyId, companyId, defaultPriceCode, // Parameters for subquery in CASE
				year, month, companyId, whsGroup, // Parameters for ELSE subquery in CASE
				currencyId, companyId, defaultPriceCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_group_id = ? AND A.item_type_id = ? AND (A.item_class_id = ? OR A.item_class_id = ?) AND A.is_active = ?", ItmGrpInventory, PurchaseTypeGoods, ItmCls, ItmClsSublet, 1).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	case utils.LinetypeFee:
		ItmCls = "4"
		ItmGrpOutsideJob := 4
		PurchaseTypeServices := "S"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				DISTINCT A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
				A.item_level_1_id AS item_level_1, 
				A.item_level_2_id AS item_level_2, 
				A.item_level_3_id AS item_level_3, 
				A.item_level_4_id AS item_level_4,
				CASE 
					WHEN ? IN (?, ?, ?) THEN
						CASE A.item_type_id
							WHEN ? THEN
								(SELECT TOP 1 price_list_amount FROM mtr_item_price_list
								WHERE is_active = 1 
								AND brand_id = B.brand_id 
								AND effective_date <= GETDATE()
								AND item_id = A.item_id
								AND currency_id = ? 
								AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
								AND price_list_code_id = ?
								ORDER BY effective_date DESC)
							ELSE
								(SELECT CASE ISNULL(price_current, 0)
										WHEN 0 THEN price_begin 
										ELSE price_current END AS HPP
								FROM mtr_group_stock 
								WHERE period_year = ? 
								AND period_month = ? 
								AND item_code = A.item_code 
								AND company_id = ?  
								AND whs_group = ?)
						END
					ELSE
						(SELECT TOP 1 price_list_amount FROM mtr_item_price_list
						WHERE is_active = 1 
						AND brand_id = B.brand_id 
						AND effective_date <= GETDATE()
						AND item_id = A.item_id 
						AND currency_id = ? 
						AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
						AND price_list_code_id = ?
						ORDER BY effective_date DESC)
				END AS PRICE
			`, year, month, companyId, whsGroup, // Parameters for AvailQty subquery
				billCode, utils.TrxTypeWoCentralize.Code, utils.TrxTypeWoInternal.Code, utils.TrxTypeWoNoCharge.Code, // Parameters for CASE statement
				2, currencyId, companyId, defaultPriceCode, // Parameters for subquery in CASE
				year, month, companyId, whsGroup, // Parameters for ELSE subquery in CASE
				currencyId, companyId, defaultPriceCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("(A.item_group_id = ? OR A.item_group_id = ?) AND A.item_class_id = ? AND A.item_type_id = ? AND A.is_active = ?", ItmGrpOutsideJob, ItmGrpInventory, ItmCls, PurchaseTypeServices, 1).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	case utils.LinetypeAccesories:
		ItmCls = "5"
		baseQuery = baseQuery.Table("mtr_item A").
			Select(`
				DISTINCT A.item_code AS item_code, 
				A.item_name AS item_name, 
				ISNULL((SELECT SUM(V.quantity_allocated) FROM mtr_location_stock V 
				        WHERE A.item_id = V.item_id 
				        AND V.PERIOD_YEAR = ? AND V.PERIOD_MONTH = ? AND V.company_id = ?), 0) AS Available_qty, 
				A.item_level_1_id AS item_level_1, 
				A.item_level_2_id AS item_level_2, 
				A.item_level_3_id AS item_level_3, 
				A.item_level_4_id AS item_level_4,
				CASE 
					WHEN ? IN (?, ?, ?) THEN
						CASE A.item_type_id
							WHEN ? THEN
								(SELECT TOP 1 price_list_amount FROM mtr_item_price_list
								WHERE is_active = 1 
								AND brand_id = B.brand_id 
								AND effective_date <= GETDATE()
								AND item_id = A.item_id
								AND currency_id = ? 
								AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
								AND price_list_code_id = ?
								ORDER BY effective_date DESC)
							ELSE
								(SELECT CASE ISNULL(price_current, 0)
										WHEN 0 THEN price_begin 
										ELSE price_current END AS HPP
								FROM mtr_group_stock 
								WHERE period_year = ? 
								AND period_month = ? 
								AND item_code = A.item_code 
								AND company_id = ?  
								AND whs_group = ?)
						END
					ELSE
						(SELECT TOP 1 price_list_amount FROM mtr_item_price_list
						WHERE is_active = 1 
						AND brand_id = B.brand_id 
						AND effective_date <= GETDATE()
						AND item_id = A.item_id 
						AND currency_id = ? 
						AND company_id = (CASE A.COMMON_PRICELIST WHEN '1' THEN 0 ELSE ? END)
						AND price_list_code_id = ?
						ORDER BY effective_date DESC)
				END AS PRICE
			`, year, month, companyId, whsGroup, // Parameters for AvailQty subquery
				billCode, utils.TrxTypeWoCentralize.Code, utils.TrxTypeWoInternal.Code, utils.TrxTypeWoNoCharge.Code, // Parameters for CASE statement
				2, currencyId, companyId, defaultPriceCode, // Parameters for subquery in CASE
				year, month, companyId, whsGroup, // Parameters for ELSE subquery in CASE
				currencyId, companyId, defaultPriceCode).
			Joins("LEFT JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = A.item_level_1_id").
			Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = A.item_level_2_id").
			Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = A.item_level_3_id").
			Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = A.item_level_4_id").
			Where("A.item_class_id = ? AND A.item_group_id = ? AND A.is_active = ?", ItmCls, ItmGrpInventory, 1).
			Where(filterQuery, filterValues...).
			Order("A.item_code")

	default:
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid linetype ID",
			Err:        errors.New("invalid linetype ID"),
		}
	}

	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count rows",
			Err:        err,
		}
	}

	offset := (paginate.Page - 1) * paginate.Limit
	if err := baseQuery.Offset(offset).Limit(paginate.Limit).Find(&results).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve data",
			Err:        err,
		}
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(paginate.Limit)))

	fmt.Println("Final Results:", results)
	fmt.Println("Total Rows:", totalRows)
	fmt.Println("Total Pages:", totalPages)

	return results, int(totalRows), totalPages, nil
}

// usp_comLookUp
// IF @strEntity = 'Vehicle0'--VEHICLE UNIT MASTER
func (r *LookupRepositoryImpl) GetVehicleUnitMaster(tx *gorm.DB, brandId int, modelId int, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		vehicleMasters []map[string]interface{}
		totalRows      int64
		totalPages     int
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
			UM.model_variant_colour_name AS Vehicle, 
			CAST(V.vehicle_production_year AS VARCHAR) AS vehicle_production_year, 
			CONVERT(VARCHAR, V.vehicle_last_service_date, 106) AS vehicle_last_service_date, 
			V.vehicle_last_km AS vehicle_last_km, 
			CASE 
				WHEN V.is_active = 1 THEN 'Active' 
				WHEN V.is_active = 0 THEN 'Deactive' 
			END AS Status
		`).
		Joins(`LEFT JOIN dms_microservices_sales_dev.dbo.mtr_vehicle_registration_certificate RC ON V.vehicle_id = RC.vehicle_id`).
		Joins(`LEFT JOIN dms_microservices_sales_dev.dbo.mtr_model_variant_colour UM ON UM.brand_id = V.vehicle_brand_id AND 
                                       UM.model_id = V.vehicle_model_id AND 
                                       UM.colour_id = V.vehicle_colour_id AND 
                                       ISNULL(UM.accessories_option_id, '') = ISNULL(V.option_id, '')`).
		Where(filterQuery, filterValues...).
		Where("V.vehicle_brand_id = ?", brandId).
		Where("V.vehicle_model_id = ?", modelId)

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total vehicle units",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	err = query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&vehicleMasters).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get vehicle unit master data",
			Err:        err,
		}
	}

	return vehicleMasters, totalPages, int(totalRows), nil
}

// usp_comLookUp
// IF @strEntity = 'Vehicle0'--VEHICLE UNIT MASTER
func (r *LookupRepositoryImpl) GetVehicleUnitByID(tx *gorm.DB, vehicleID int, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		vehicleMasters []map[string]interface{}
		totalRows      int64
		totalPages     int
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	// Apply filters
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
			UM.model_variant_colour_name AS Vehicle, 
			CAST(V.vehicle_production_year AS VARCHAR) AS vehicle_production_year, 
			CONVERT(VARCHAR, V.vehicle_last_service_date, 106) AS vehicle_last_service_date, 
			V.vehicle_last_km AS vehicle_last_km, 
			CASE 
				WHEN V.is_active = 1 THEN 'Active' 
				WHEN V.is_active = 0 THEN 'Deactive' 
			END AS Status
		`).
		Joins(`LEFT JOIN dms_microservices_sales_dev.dbo.mtr_vehicle_registration_certificate RC ON V.vehicle_id = RC.vehicle_id`).
		Joins(`LEFT JOIN dms_microservices_sales_dev.dbo.mtr_model_variant_colour UM ON UM.brand_id = V.vehicle_brand_id AND 
                                       UM.model_id = V.vehicle_model_id AND 
                                       UM.colour_id = V.vehicle_colour_id AND 
                                       ISNULL(UM.accessories_option_id, '') = ISNULL(V.option_id, '')`).
		Where(filterQuery, filterValues...).
		Where("V.vehicle_id = ?", vehicleID)

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total vehicle units",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	err = query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&vehicleMasters).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get vehicle unit data by ID",
			Err:        err,
		}
	}

	return vehicleMasters, totalPages, int(totalRows), nil
}

// usp_comLookUp
// IF @strEntity = 'Vehicle0'--VEHICLE UNIT MASTER
func (r *LookupRepositoryImpl) GetVehicleUnitByChassisNumber(tx *gorm.DB, chassisNumber string, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		vehicleMasters []map[string]interface{}
		totalRows      int64
		totalPages     int
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	// Apply filters
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
			UM.model_variant_colour_name AS Vehicle, 
			CAST(V.vehicle_production_year AS VARCHAR) AS vehicle_production_year, 
			CONVERT(VARCHAR, V.vehicle_last_service_date, 106) AS vehicle_last_service_date, 
			V.vehicle_last_km AS vehicle_last_km, 
			CASE 
				WHEN V.is_active = 1 THEN 'Active' 
				WHEN V.is_active = 0 THEN 'Deactive' 
			END AS Status
		`).
		Joins(`LEFT JOIN dms_microservices_sales_dev.dbo.mtr_vehicle_registration_certificate RC ON V.vehicle_id = RC.vehicle_id`).
		Joins(`LEFT JOIN dms_microservices_sales_dev.dbo.mtr_model_variant_colour UM ON UM.brand_id = V.vehicle_brand_id AND 
                                       UM.model_id = V.vehicle_model_id AND 
                                       UM.colour_id = V.vehicle_colour_id AND 
                                       ISNULL(UM.accessories_option_id, '') = ISNULL(V.option_id, '')`).
		Where(filterQuery, filterValues...).
		Where("V.vehicle_chassis_number = ?", chassisNumber)

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total vehicle units",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	err = query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&vehicleMasters).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get vehicle unit data by chassis number",
			Err:        err,
		}
	}

	return vehicleMasters, totalPages, int(totalRows), nil
}

// usp_comLookUp
// IF @strEntity = 'CampaignMaster'--CAMPAIGN MASTER
func (r *LookupRepositoryImpl) GetCampaignMaster(tx *gorm.DB, companyId int, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	var (
		campaignMasters []map[string]interface{}
		totalRows       int64
		totalPages      int
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
		//Joins(`LEFT JOIN dms_microservices_sales_dev.dbo.mtr_model_variant_colour VC ON C.model_id = VC.model_id`).
		Where(filterQuery, filterValues...).
		Where("C.company_id = ?", companyId)

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total vehicle units",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	err = query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&campaignMasters).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get vehicle unit master data",
			Err:        err,
		}
	}

	return campaignMasters, totalPages, int(totalRows), nil
}

// usp_comLookUp
// IF @strEntity = 'WorkOrderService'--WO SERVICE
func (r *LookupRepositoryImpl) WorkOrderService(tx *gorm.DB, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		results []struct {
			WorkOrderNo    string
			WorkOrderDate  time.Time
			NoPolisi       string
			ChassisNo      string
			Brand          int
			Model          int
			Variant        int
			WorkOrderSysNo int
		}
		totalRows  int64
		totalPages int
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	filterStrings := []string{}
	filterValues := []interface{}{}
	if len(filters) > 0 {
		for _, filter := range filters {
			filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
			filterValues = append(filterValues, filter.ColumnValue)
		}
	}

	filterQuery := strings.Join(filterStrings, " AND ")
	if len(filterStrings) > 0 {
		tx = tx.Where(filterQuery, filterValues...)
	}

	query := tx.Table("trx_work_order_allocation AS A").
		Select("A.work_order_document_number AS WorkOrderNo, B.work_order_date AS WorkOrderDate, "+
			"B.vehicle_chassis_number AS ChassisNo, B.brand_id AS Brand, B.model_id AS Model, "+
			"B.variant_id AS Variant, A.work_order_system_number AS WorkOrderSysNo").
		Joins("LEFT JOIN trx_work_order AS B ON B.work_order_system_number = A.work_order_system_number").
		Where("A.service_status_id NOT IN (?, ?, ?, ?)", utils.SrvStatStop, utils.SrvStatAutoRelease, utils.SrvStatTransfer, utils.SrvStatQcPass)

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total vehicle units",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	err = query.
		Order("A.work_order_document_number").
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&results).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get vehicle unit master data",
			Err:        err,
		}
	}

	mappedResults := make([]map[string]interface{}, len(results))
	for i, result := range results {
		mappedResults[i] = map[string]interface{}{
			"work_order_document_number": result.WorkOrderNo,
			"work_order_date":            result.WorkOrderDate,
			"vehicle_tnkb":               result.NoPolisi,
			"vehicle_chassis_number":     result.ChassisNo,
			"brand_id":                   result.Brand,
			"model_id":                   result.Model,
			"variant_id":                 result.Variant,
			"work_order_system_number":   result.WorkOrderSysNo,
		}
	}

	return mappedResults, totalPages, int(totalRows), nil
}

// usp_comLookUp
// IF @strEntity =  'CustomerByTypeAndAddress'--CUSTOMER MASTER
func (r *LookupRepositoryImpl) CustomerByTypeAndAddress(tx *gorm.DB, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		customerMasters []map[string]interface{}
		totalRows       int64
		totalPages      int
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
		Where(filterQuery, filterValues...).
		Where("C.is_active = 1")

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total vehicle units",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	err = query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&customerMasters).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get vehicle unit master data",
			Err:        err,
		}
	}

	return customerMasters, totalPages, int(totalRows), nil
}

// usp_comLookUp
// IF @strEntity =  'CustomerByTypeAndAddress'--CUSTOMER MASTER
func (r *LookupRepositoryImpl) CustomerByTypeAndAddressByID(tx *gorm.DB, customerId int, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		customerMasters []map[string]interface{}
		totalRows       int64
		totalPages      int
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
		Where(filterQuery, filterValues...).
		Where("C.customer_id = ?", customerId)

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total customers",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	err = query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&customerMasters).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get customer data",
			Err:        err,
		}
	}

	return customerMasters, totalPages, int(totalRows), nil
}

// usp_comLookUp
// IF @strEntity =  'CustomerByTypeAndAddress'--CUSTOMER MASTER
func (r *LookupRepositoryImpl) CustomerByTypeAndAddressByCode(tx *gorm.DB, customerCode string, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		customerMasters []map[string]interface{}
		totalRows       int64
		totalPages      int
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
			C.customer_id,
			C.customer_code,
			C.customer_name,
			CA.client_type_description,
			A.address_street_1 AS address_1,
			A.address_street_2 AS address_2,
			A.address_street_3 AS address_3,
			C.id_phone_no
		`).
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_client_type CA ON C.client_type_id = CA.client_type_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_address AS A ON C.id_address_id = A.address_id").
		Where(filterQuery, filterValues...).
		Where("C.customer_code = ?", customerCode)

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total customers",
			Err:        err,
		}
	}

	if totalRows > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	if err := query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&customerMasters).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get customer data",
			Err:        err,
		}
	}

	return customerMasters, totalPages, int(totalRows), nil
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
					// Handle customer logic
					var typeMap string

					// First query: check company_type in mtr_company_type_map_from_customer
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

					// Assign billCode based on customer type
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

				// If no result found in company_type_map_from_customer, query the mtr_supplier table
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

				// Assign billCode based on supplier type
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

	// Additional logic for other supcusType cases ('S', 'P', 'C', 'W') can be added here.

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
func (r *LookupRepositoryImpl) GetLineTypeByItemCode(tx *gorm.DB, itemCode string) (int, *exceptions.BaseErrorResponse) {
	var (
		lineType         int
		itemGrp          string
		itemTypeId       int
		itemCls          string
		itmClsSublet     = "SB" // Assuming these are constants in your utils
		lineTypeSublet   = utils.LinetypeSublet
		itemClsFee       = "WF"
		itemGrpInventory = "IN"
		itemClsAccs      = "AC"
		itemClsSv        = "SV"
	)

	// Retrieve item details
	var itemDetails struct {
		ItemGroupId string
		ItemTypeId  int
		ItemClassId string
	}

	if err := tx.Model(&masteritementities.Item{}).
		Select("item_group_id, item_type, item_class_id").
		Where("item_code = ? ", itemCode). // Add record status check
		Scan(&itemDetails).Error; err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get item details",
			Err:        err,
		}
	}

	itemGrp = itemDetails.ItemGroupId
	itemTypeId = itemDetails.ItemTypeId
	itemCls = itemDetails.ItemClassId

	// Determine line type based on the item details
	if itemGrp == itemGrpInventory {
		if itemTypeId == 1 {
			switch itemCls {
			case "SP":
				lineType = utils.LinetypeSparepart
			case "OL":
				lineType = utils.LinetypeOil
			case "MT", itmClsSublet:
				lineType = utils.LinetypeMaterial
			case "CM":
				lineType = utils.LinetypeConsumableMaterial
			case itemClsAccs:
				lineType = utils.LinetypeAccesories
			case itemClsSv:
				lineType = utils.LinetypeSublet
			default:
				lineType = utils.LinetypeAccesories
			}
		} else if itemCls == itemClsFee {
			lineType = lineTypeSublet
		} else if itemCls == itemClsAccs && itemTypeId == 2 {
			lineType = utils.LinetypeOperation
		}
	} else if itemGrp == "OJ" || (itemGrp == itemGrpInventory && itemTypeId == 2 && itemCls == itemClsFee) {
		lineType = lineTypeSublet
	}

	// Check if the item exists
	var itemExists int64
	if err := tx.Model(&masteritementities.Item{}).
		Where("item_code = ? ", itemCode).
		Count(&itemExists).Error; err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count item",
			Err:        err,
		}
	}

	if itemExists == 0 {
		var packageExists int64
		if err := tx.Model(&masterentities.PackageMaster{}).
			Where("package_code = ? ", itemCode).
			Count(&packageExists).Error; err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to count package",
				Err:        err,
			}
		}
		if packageExists > 0 {
			lineType = utils.LinetypePackage
		} else {
			lineType = utils.LinetypeOperation
		}
	}

	return lineType, nil
}

func (r *LookupRepositoryImpl) GetWhsGroup(tx *gorm.DB, companyCode int) (int, *exceptions.BaseErrorResponse) {
	var (
		whsGroup int
		err      error
	)

	// Execute the GORM query
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

	// Main query to fetch campaign discount details based on the input parameters
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
	err := whereQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&response).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
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
	err := whereQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&responses).Error
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
	err := whereQuery.Scopes(pagination.PaginateDistinct(&pages, whereQuery)).Scan(&responses).Error

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
func (r *LookupRepositoryImpl) ReferenceTypeWorkOrder(tx *gorm.DB, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		results []struct {
			WorkOrderDocumentNumber    string
			WorkOrderDate              time.Time
			WorkOrderStatusId          int
			WorkOrderStatusDescription string
			WorkOrderSystemNumber      int
		}
		totalRows  int64
		totalPages int
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	filterStrings := []string{}
	filterValues := []interface{}{}
	if len(filters) > 0 {
		for _, filter := range filters {
			filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
			filterValues = append(filterValues, filter.ColumnValue)
		}
	}

	filterQuery := strings.Join(filterStrings, " AND ")
	if len(filterStrings) > 0 {
		tx = tx.Where(filterQuery, filterValues...)
	}

	query := tx.Table("trx_work_order AS A").
		Select("A.work_order_document_number AS work_order_document_number, A.work_order_date AS work_order_date, "+
			"B.work_order_status_id AS work_order_status_id, E.work_order_status_description AS work_order_status_description, A.work_order_system_number AS work_order_system_number").
		Joins("INNER JOIN trx_work_order_detail AS B ON B.work_order_system_number = A.work_order_system_number").
		Joins("LEFT OUTER JOIN trx_service_request AS C ON A.service_request_system_number = C.service_request_system_number AND C.reference_type_id = 1 AND C.service_request_status_id NOT IN (4, 5) AND NOT (C.service_request_status_id = 8 AND COALESCE(C.booking_system_number, 0) != 0 AND COALESCE(C.work_order_system_number, 0) != 0)").
		Joins("LEFT OUTER JOIN trx_service_request_detail AS D ON C.service_request_system_number = D.service_request_system_number AND D.operation_item_id = B.operation_item_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_work_order_status AS E ON A.work_order_status_id = E.work_order_status_id").
		Where("B.work_order_status_id NOT IN (?, ?, ?)", utils.WoStatDraft, utils.WoStatClosed, utils.WoStatCancel).
		Where("COALESCE(A.work_order_system_number, 0) != 0").
		Where("COALESCE(D.service_request_line_number, 0) != 0")

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total Data",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	err = query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&results).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return nil, 0, 0, &exceptions.BaseErrorResponse{
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

	return mappedResults, totalPages, int(totalRows), nil
}

// usp_comLookUp
// IF @strEntity = 'ServiceReqRefTypeWO'--SERVICE REQUEST REF TYPE WO
func (r *LookupRepositoryImpl) ReferenceTypeWorkOrderByID(tx *gorm.DB, referenceId int, paginate pagination.Pagination, filters []utils.FilterCondition) (map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		result struct {
			WorkOrderDocumentNumber string
			WorkOrderDate           time.Time
			WoStatusId              int
			WorkOrderStatus         string
			WorkOrderSystemNumber   int
		}
		totalRows  int64
		totalPages int
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	filterStrings := []string{}
	filterValues := []interface{}{}
	if len(filters) > 0 {
		for _, filter := range filters {
			filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
			filterValues = append(filterValues, filter.ColumnValue)
		}
	}

	filterQuery := strings.Join(filterStrings, " AND ")
	if len(filterStrings) > 0 {
		tx = tx.Where(filterQuery, filterValues...)
	}

	query := tx.Table("trx_work_order AS A").
		Select("A.work_order_document_number AS WorkOrderNo, A.work_order_date AS WorkOrderDate, "+
			"B.work_order_status_id AS WoStatusId, E.work_order_status_description AS WorkOrderStatus, A.work_order_system_number AS WorkOrderSysNo").
		Joins("INNER JOIN trx_work_order_detail AS B ON B.work_order_system_number = A.work_order_system_number").
		Joins("LEFT OUTER JOIN trx_service_request AS C ON A.service_request_system_number = C.service_request_system_number AND C.reference_type_id = 1 AND C.service_request_status_id NOT IN (4, 5) AND NOT (C.service_request_status_id = 8 AND COALESCE(C.booking_system_number, 0) != 0 AND COALESCE(C.work_order_system_number, 0) != 0)").
		Joins("LEFT OUTER JOIN trx_service_request_detail AS D ON C.service_request_system_number = D.service_request_system_number AND D.operation_item_id = B.operation_item_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_work_order_status AS E ON A.work_order_status_id = E.work_order_status_id").
		Where("B.work_order_status_id NOT IN (?, ?, ?)", utils.WoStatDraft, utils.WoStatClosed, utils.WoStatCancel).
		Where("COALESCE(A.work_order_system_number, 0) != 0").
		Where("COALESCE(D.service_request_line_number, 0) != 0").
		Where("A.work_order_system_number = ?", referenceId)

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total Data",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	err = query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get data",
			Err:        err,
		}
	}

	mappedResult := map[string]interface{}{
		"work_order_document_number": result.WorkOrderDocumentNumber,
		"work_order_date":            result.WorkOrderDate,
		"work_order_status_id":       result.WoStatusId,
		"work_order_status":          result.WorkOrderStatus,
		"work_order_system_number":   result.WorkOrderSystemNumber,
	}

	return mappedResult, totalPages, int(totalRows), nil
}

// usp_comLookUp
// IF @strEntity = 'ServiceReqRefTypeSO'--SERVICE REQUEST REF TYPE SO
func (r *LookupRepositoryImpl) ReferenceTypeSalesOrder(tx *gorm.DB, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		results []struct {
			SalesOrderDocumentNumber    string
			SalesOrderDate              time.Time
			SalesOrderStatusId          int
			SalesOrderStatusDescription string
			SalesOrderSystemNumber      int
		}
		totalRows  int64
		totalPages int
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	filterStrings := []string{}
	filterValues := []interface{}{}
	if len(filters) > 0 {
		for _, filter := range filters {
			filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
			filterValues = append(filterValues, filter.ColumnValue)
		}
	}

	filterQuery := strings.Join(filterStrings, " AND ")
	if len(filterStrings) > 0 {
		tx = tx.Where(filterQuery, filterValues...)
	}

	query := tx.Table("trx_sales_order AS A").
		Select("A.sales_order_document_number AS sales_order_document_number, A.sales_order_date AS sales_order_date, "+
			"B.sales_order_status_id AS work_order_status_id, E.sales_order_status_description AS sales_order_status_description, A.sales_order_system_number AS sales_order_system_number").
		Joins("INNER JOIN trx_sales_order_detail AS B ON B.sales_order_system_number = A.sales_order_system_number").
		Joins("LEFT OUTER JOIN trx_service_request AS C ON A.service_request_system_number = C.service_request_system_number AND C.reference_type_id = 1 AND C.service_request_status_id NOT IN (4, 5) AND NOT (C.service_request_status_id = 8 AND COALESCE(C.booking_system_number, 0) != 0 AND COALESCE(C.sales_order_system_number, 0) != 0)").
		Joins("LEFT OUTER JOIN trx_service_request_detail AS D ON C.service_request_system_number = D.service_request_system_number AND D.operation_item_id = B.operation_item_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_sales_order_status AS E ON A.sales_order_status_id = E.work_order_status_id").
		Where("B.sales_order_status_id NOT IN (?, ?, ?)", utils.WoStatDraft, utils.WoStatClosed, utils.WoStatCancel).
		Where("COALESCE(A.sales_order_system_number, 0) != 0").
		Where("COALESCE(D.service_request_line_number, 0) != 0")

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total Data",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	err = query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		Find(&results).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get data",
			Err:        err,
		}
	}

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

	return mappedResults, totalPages, int(totalRows), nil
}

// usp_comLookUp
// IF @strEntity = 'ServiceReqRefTypeSO'--SERVICE REQUEST REF TYPE SO
func (r *LookupRepositoryImpl) ReferenceTypeSalesOrderByID(tx *gorm.DB, referenceId int, paginate pagination.Pagination, filters []utils.FilterCondition) (map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var (
		result struct {
			SalesOrderDocumentNumber    string
			SalesOrderDate              time.Time
			SalesOrderStatusId          int
			SalesOrderStatusDescription string
			SalesOrderSystemNumber      int
		}
		totalRows  int64
		totalPages int
	)

	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	filterStrings := []string{}
	filterValues := []interface{}{}
	if len(filters) > 0 {
		for _, filter := range filters {
			filterStrings = append(filterStrings, fmt.Sprintf("%s = ?", filter.ColumnField))
			filterValues = append(filterValues, filter.ColumnValue)
		}
	}

	filterQuery := strings.Join(filterStrings, " AND ")
	if len(filterStrings) > 0 {
		tx = tx.Where(filterQuery, filterValues...)
	}

	query := tx.Table("trx_sales_order AS A").
		Select("A.sales_order_document_number AS sales_order_document_number, A.sales_order_date AS sales_order_date, "+
			"B.sales_order_status_id AS work_order_status_id, E.sales_order_status_description AS sales_order_status_description, A.sales_order_system_number AS sales_order_system_number").
		Joins("INNER JOIN trx_sales_order_detail AS B ON B.sales_order_system_number = A.sales_order_system_number").
		Joins("LEFT OUTER JOIN trx_service_request AS C ON A.service_request_system_number = C.service_request_system_number AND C.reference_type_id = 1 AND C.service_request_status_id NOT IN (4, 5) AND NOT (C.service_request_status_id = 8 AND COALESCE(C.booking_system_number, 0) != 0 AND COALESCE(C.sales_order_system_number, 0) != 0)").
		Joins("LEFT OUTER JOIN trx_service_request_detail AS D ON C.service_request_system_number = D.service_request_system_number AND D.operation_item_id = B.operation_item_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_sales_order_status AS E ON A.sales_order_status_id = E.work_order_status_id").
		Where("B.sales_order_status_id NOT IN (?, ?, ?)", utils.WoStatDraft, utils.WoStatClosed, utils.WoStatCancel).
		Where("COALESCE(A.sales_order_system_number, 0) != 0").
		Where("COALESCE(D.service_request_line_number, 0) != 0").
		Where("A.sales_order_system_number = ?", referenceId)

	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count total Data",
			Err:        err,
		}
	}

	if paginate.Limit > 0 {
		totalPages = int(totalRows) / paginate.Limit
		if int(totalRows)%paginate.Limit != 0 {
			totalPages++
		}
	}

	err = query.
		Offset((paginate.Page - 1) * paginate.Limit).
		Limit(paginate.Limit).
		First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get data",
			Err:        err,
		}
	}

	mappedResult := map[string]interface{}{
		"sales_order_document_number":    result.SalesOrderDocumentNumber,
		"sales_order_date":               result.SalesOrderDate,
		"sales_order_status_id":          result.SalesOrderStatusId,
		"sales_order_status_description": result.SalesOrderStatusDescription,
		"sales_order_system_number":      result.SalesOrderSystemNumber,
	}

	return mappedResult, totalPages, int(totalRows), nil
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
