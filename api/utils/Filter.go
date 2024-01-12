package utils

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type FilterCondition struct {
	ColumnValue string
	ColumnField string
}

// ApplyFilter function is to generate where conditions to match values from columnValue/values agains columnName/columns
func ApplyFilter(db *gorm.DB, criteria []FilterCondition) *gorm.DB {
	var queryWhere []string
	var columnValue, columnName []string

	for _, c := range criteria {
		columnValue, columnName = append(columnValue, c.ColumnValue), append(columnName, c.ColumnField)
	}

	for i := 0; i < len(columnValue); i++ {
		if strings.Contains(columnName[i], "is_active") {
			n := map[string]string{"true": "1", "false": "0", "Active": "1"} //kurang inactive/disable apalah
			columnValue[i] = n[columnValue[i]]
		}
		condition := columnName[i] + " LIKE " + "'%" + columnValue[i] + "%'"
		queryWhere = append(queryWhere, condition)
	}
	queryFinal := db.Where(strings.Join(queryWhere, " AND "))

	return queryFinal
}

func DefineInternalExternalFilter(filterCondition []FilterCondition, tableStruct interface{}) ([]FilterCondition, []FilterCondition) {
	var internalFilter, externalFilter []FilterCondition
	responseStruct := reflect.TypeOf(tableStruct)
	for i := 0; i < len(filterCondition); i++ {
		flag := false
		for j := 0; j < responseStruct.NumField(); j++ {
			if filterCondition[i].ColumnField == responseStruct.Field(j).Tag.Get("parent_entity")+"."+responseStruct.Field(j).Tag.Get("json") {
				internalFilter = append(internalFilter, filterCondition[i])
				flag = true
				break
			}
		}
		if !flag {
			externalFilter = append(externalFilter, filterCondition[i])
		}
	}
	return internalFilter, externalFilter
}

func CreateColumnExternalFilter(externalFilter []FilterCondition, column ...string) []string {
	var columnExternalFilter []string
	for i := 0; i < len(column); i++ {
		for j := 0; j < len(externalFilter); j++ {
			if externalFilter[j].ColumnField == column[i] {
				columnExternalFilter = append(columnExternalFilter, externalFilter[j].ColumnValue)
				break
			}
		}
	}
	return columnExternalFilter
}
