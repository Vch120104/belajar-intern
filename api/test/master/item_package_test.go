package test

import (
	"after-sales/api/config"
	masteritempayloads "after-sales/api/payloads/master/item"
	"fmt"

	masteritementities "after-sales/api/entities/master/item"
	"testing"
)

func TestGetItemPackageDetailByItemPackageId(t *testing.T) {

	config.InitEnvConfigs(true, "")
	db := config.InitDB()

	model := masteritementities.ItemPackage{}
	var responses []masteritempayloads.ItemPackageDetailResponse

	err := db.Model(&model).
		Select(
			"item_package_detail_id",
			"ItemPackageDetail.is_active is_active",
			"ItemPackageDetail.item_package_id item_package_id",
			"ItemPackageDetail__Item.item_id item_id",
			"ItemPackageDetail__Item.item_code item_code",
			"ItemPackageDetail__Item.item_name item_name",
			"ItemPackageDetail__Item.item_class_id item_class_id",
			"ItemPackageDetail__Item__ItemClass.item_class_code item_class_code",
			"ItemPackageDetail.quantity quantity",
		).
		Joins("ItemPackageDetail", db.Select("1")).
		Joins("ItemPackageDetail.Item", db.Select("1")).
		Joins("ItemPackageDetail.Item.ItemClass", db.Select("1")).
		Scan(&responses).
		Error

	if err != nil {
		fmt.Println(responses, err)
	}

	fmt.Println(responses)

}

func TestGetItemPackageDetailByItemPackageIdPreload(t *testing.T) {

	config.InitEnvConfigs(true, "")
	db := config.InitDB()

	model := masteritementities.ItemPackage{}
	var detailResponse []masteritempayloads.ItemPackageDetailResponse

	err := db.
		Preload("ItemPackageDetail.Item").
		Preload("ItemPackageDetail.Item.ItemClass").
		Find(&model).Scan(&detailResponse)

	header := masteritempayloads.GetAllItemPackageResponse{
		ItemPackageId:   model.ItemPackageId,
		IsActive:        model.IsActive,
		ItemGroupId:     model.ItemGroupId,
		ItemPackageCode: model.ItemPackageCode,
		ItemPackageName: model.ItemPackageName,
		ItemPackageSet:  model.ItemPackageSet,
		Description:     model.Description,
	}

	for i := 0; i < len(model.ItemPackageDetail); i++ {
		data := model.ItemPackageDetail[i]
		detail := masteritempayloads.ItemPackageDetailResponse{
			ItemPackageDetailId: data.ItemPackageDetailId,
			IsActive:            data.IsActive,
			ItemPackageId:       data.ItemPackageId,
			ItemId:              data.ItemId,
			ItemCode:            data.Item.ItemCode,
			ItemName:            data.Item.ItemName,
			ItemClassId:         data.Item.ItemClassId,
			ItemClassCode:       data.Item.ItemClass.ItemClassCode,
			Quantity:            data.Quantity,
		}
		detailResponse = append(detailResponse, detail)
	}

	response := masteritempayloads.ItemPackageDetailPayload{
		GetAllItemPackageResponse: header,
		ItemPackageDetailResponse: detailResponse,
	}

	if err != nil {
		fmt.Println(response, err)
	}

}
