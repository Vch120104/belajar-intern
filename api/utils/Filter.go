package utils

import (
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
)

type FilterCondition struct {
	ColumnValue string
	ColumnField string
}

// ApplyFilter generates WHERE conditions based on a set of filter criteria and applies them to a GORM query.
//
// Parameters:
//   - db: A pointer to a GORM database query to which the WHERE conditions will be applied.
//   - criteria: A slice of FilterCondition representing the filter criteria to be applied.
//
// Returns:
//   - result: A modified GORM database query with WHERE conditions based on the provided filter criteria.
//
// Example Usage:
//
//		type FilterCondition struct {
//		    ColumnField string // e.g., "name", "is_active"
//		    ColumnValue string // e.g., "John", "true"
//		}
//		criteria := []FilterCondition{
//		   {ColumnField: "name", ColumnValue: "John"},
//	   	   {ColumnField: "is_active", ColumnValue: "true"},
//	       {ColumnField: "created_at_from", ColumnValue: "2024-01-01"},
//	       {ColumnField: "id#multiple", ColumnValue: "1,2,3"},
//		}
//		query := ApplyFilter(db, criteria)
//
// Notes:
//   - The function takes a GORM database query and a slice of filter criteria as input.
//   - It iterates through the filter criteria, constructing WHERE conditions for each criterion.
//   - If the column name contains "is_active," it maps values to standard boolean representations ("true": "1", "false": "0", "Active": "1").
//   - The generated WHERE conditions are joined using "AND" and applied to the input database query.
//   - The modified GORM database query is then returned as the result.
func ApplyFilter(db *gorm.DB, criteria []FilterCondition) *gorm.DB {
	var queryWhere []string

	for _, c := range criteria {
		// Abaikan filter jika nilai kosong
		if c.ColumnValue == "" {
			continue
		}

		column := c.ColumnField
		value := c.ColumnValue
		var condition string

		if strings.ToLower(value) == "true" || strings.ToLower(value) == "false" {
			// Handle boolean values
			boolMap := map[string]string{"true": "1", "false": "0"}
			condition = column + " = " + boolMap[strings.ToLower(value)]
		} else if strings.HasSuffix(column, "_from") {
			// Handle range filter (>=)
			key := strings.TrimSuffix(column, "_from")
			condition = key + " >= '" + value + "'"
		} else if strings.HasSuffix(column, "_to") {
			// Handle range filter (<=)
			key := strings.TrimSuffix(column, "_to")
			condition = key + " <= '" + value + "'"
		} else if strings.Contains(column, "date") {
			// Handle exact date match
			parsedDate, err := time.Parse("2006-01-02", value)
			if err == nil {
				condition = column + " = '" + parsedDate.Format("2006-01-02") + "'"
			}
		} else if strings.HasSuffix(column, "#multiple") {
			// Handle IN condition
			key := strings.TrimSuffix(column, "#multiple")
			values := strings.Split(value, ",")
			quotedValues := "'" + strings.Join(values, "','") + "'"
			condition = key + " IN (" + quotedValues + ")"
		} else {
			// Default to LIKE condition
			condition = column + " LIKE '%" + value + "%'"
		}

		if condition != "" {
			queryWhere = append(queryWhere, condition)
		}
	}

	return db.Where(strings.Join(queryWhere, " AND "))
}

// DefineInternalExternalFilter categorizes filter conditions into internal and external filters based on the provided table structure.
//
// Parameters:
//   - filterCondition: A slice of FilterCondition representing the filter conditions to be categorized.
//   - tableStruct: An instance of a structure representing the table's fields and their attributes.
//
// Returns:
//   - internalFilter: A slice of FilterCondition containing filter conditions that match fields within the provided table structure.
//   - externalFilter: A slice of FilterCondition containing filter conditions that do not match fields within the provided table structure.
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
