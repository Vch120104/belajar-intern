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
	var columnValue, columnName []string
	var condition string
	var key string

	// Iterate through the criteria and prepare the column values and column names
	for _, c := range criteria {
		columnValue = append(columnValue, c.ColumnValue)
		columnName = append(columnName, c.ColumnField)
	}

	// Apply conditions based on column values and field names
	for i := 0; i < len(columnValue); i++ {
		// Handle boolean-like values (true, false, Active)
		if strings.Contains(columnValue[i], "true") || strings.Contains(columnValue[i], "false") || strings.Contains(columnValue[i], "Active") {
			n := map[string]string{"true": "1", "false": "0", "Active": "1"}
			columnValue[i] = n[columnValue[i]]
		}

		// Handle date ranges (date_from and date_to)
		if strings.Contains(columnName[i], "date_from") {
			key = strings.Split(columnName[i], "_from")[0]
			condition = key + " >= '" + columnValue[i] + "'"
		} else if strings.Contains(columnName[i], "date_to") {
			key = strings.Split(columnName[i], "_to")[0]
			condition = key + " <= '" + columnValue[i] + "'"
		} else if strings.Contains(columnName[i], "date") {
			// Handle date fields using RFC3339 format ("2025-08-15T00:00:00Z")
			parsedDate, err := time.Parse(time.RFC3339, columnValue[i])
			if err != nil {
				parsedDate, err = time.Parse("2006-01-02", columnValue[i])
				if err != nil {
					continue
				}
			}
			condition = "CONVERT(DATE, " + columnName[i] + ") = '" + parsedDate.Format("2006-01-02") + "'"
		} else if strings.Contains(columnName[i], "id") {
			// Handle multiple IDs and single ID filtering
			if strings.Contains(columnName[i], "#multiple") {
				condition = columnName[i] + " IN (" + columnValue[i] + ")"
			} else {
				condition = columnName[i] + " LIKE '" + columnValue[i] + "'"
			}
		} else {
			// Default condition (LIKE match for text)
			condition = columnName[i] + " LIKE '%" + columnValue[i] + "%'"
		}

		// Add the condition to the WHERE clause
		queryWhere = append(queryWhere, condition)
	}

	// Combine all conditions using AND and apply them to the GORM query
	queryFinal := db.Where(strings.Join(queryWhere, " AND "))

	return queryFinal
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
