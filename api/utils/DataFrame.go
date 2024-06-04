package utils

import (
	"math"
	"reflect"

	"strings"

	"github.com/go-gota/gota/dataframe"
)

// DataFramePaginate is to performs pagination on a given dataset represented as a slice of maps or a slice of structs
//
// Parameters:
//   - data: The input dataset to paginate. It can be either a slice of maps or a slice of structs (ex: []struct{} or maps[string]interface{})
//   - page: The page number to retrieve (starting from 0).
//   - limit: The number of rows per page.
//   - sortOf: The column name to sort the data by. Leave empty for no sorting.
//   - sortBy: The sorting order, either "asc" for ascending or "desc" for descending. Ignored if sortOf is empty.
//
// Returns:
//   - result: A slice of maps representing the paginated data for the specified page and limit.
//   - totalPages: The total number of pages based on the given limit and the total number of rows in the dataset.
//   - totalRows: The total number of rows in the dataset.
func DataFramePaginate(data interface{}, page int, limit int, sortOf string, sortBy string) (result []map[string]interface{}, totalPages int, totalRows int) {
	var df dataframe.DataFrame
	tpy, _ := reflect.TypeOf(data), reflect.ValueOf(data)

	//read type of data interface{} to determine which data are loaded
	if tpy.Kind() == reflect.Slice && tpy.Elem().Kind() != reflect.Struct {
		df = dataframe.LoadMaps(data.([]map[string]interface{}))
	} else {
		df = dataframe.LoadStructs(data)
	}

	if strings.Contains(sortOf, "_") {
		SnaketoPascalCase(sortOf)
	}

	totalRows = df.Nrow()
	if sortOf != "" {
		if sortBy == "desc" {
			df = df.Arrange(dataframe.RevSort(sortOf))
		} else {
			df = df.Arrange(dataframe.Sort(sortOf))
		}
	}

	start := page * limit
	end := start + limit

	if end > df.Nrow() {
		end = df.Nrow()
	}

	indices := make([]int, end-start)
	for i := start; i < end; i++ {
		indices[i-start] = i
	}

	totalPages = int(math.Ceil(float64(totalRows) / float64(limit)))

	df = df.Subset(indices)
	return df.Maps(), totalPages, totalRows
}

// DataFrameLeftJoin performs a left join operation on two datasets using the dataframe package
//
// Parameters:
//   - data1: The left dataset to perform the join on. It can be either a slice of maps or a slice of structs.
//   - data2: The right dataset to join with. It should have the same type as data1.
//   - key: The column name used as the join key.
//
// Returns:
//   - result: A slice of maps representing the result of the left join operation. The resulting maps contain columns from both datasets.
//
// Note:
//   - data1 are able to accept []map[string]interface{} and []struct{}, BUT data2 are only able to accept []struct{}
//   - if the returned data are not found after running this function, make sure your data type are correct
func DataFrameLeftJoin(data1 interface{}, data2 interface{}, key string) []map[string]interface{} {
	tpy, _ := reflect.TypeOf(data1), reflect.ValueOf(data1)
	if tpy.Kind() == reflect.Slice && tpy.Elem().Kind() != reflect.Struct {
		df1 := dataframe.LoadMaps(data1.([]map[string]interface{}))
		df2 := dataframe.LoadStructs(data2)
		dfJoin := df1.LeftJoin(df2, key)
		return ConvertNullValueToEmptyString(dfJoin.Maps())
	} else {
		df1 := dataframe.LoadStructs(data1)
		df2 := dataframe.LoadStructs(data2)
		dfJoin := df1.LeftJoin(df2, key)
		return ConvertNullValueToEmptyString(dfJoin.Maps())
	}
}

// DataFrameLeftJoin performs a inner join operation on two datasets using the dataframe package
//
// Parameters:
//   - data1: The left dataset to perform the join on. It can be either a slice of maps or a slice of structs.
//   - data2: The right dataset to join with. It should have the same type as data1.
//   - key: The column name used as the join key.
//
// Returns:
//   - result: A slice of maps representing the result of the inner join operation. The resulting maps contain columns from both datasets.
//
// Note:
//   - data1 are able to accept []map[string]interface{} and []struct{}, BUT data2 are only able to accept []struct{}
//   - if the returned data are not found after running this function, make sure your data type are correct
func DataFrameInnerJoin(data1 interface{}, data2 interface{}, key string) []map[string]interface{} {
	tpy, _ := reflect.TypeOf(data1), reflect.ValueOf(data1)
	if tpy.Kind() == reflect.Slice && tpy.Elem().Kind() != reflect.Struct {
		df1 := dataframe.LoadMaps(data1.([]map[string]interface{}))
		df2 := dataframe.LoadStructs(data2)
		dfJoin := df1.InnerJoin(df2, key)
		return ConvertNullValueToEmptyString(dfJoin.Maps())
	} else {
		df1 := dataframe.LoadStructs(data1)
		df2 := dataframe.LoadStructs(data2)
		dfJoin := df1.InnerJoin(df2, key)
		return ConvertNullValueToEmptyString(dfJoin.Maps())
	}
}

// ConvertNullValueToEmptyString converts null values in a slice of maps to empty strings.
//
// Parameters:
//   - data: A slice of maps where each map represents a row of data with column names as keys.
//
// Returns:
//   - result: A modified slice of maps where any null values have been replaced with empty strings.
func ConvertNullValueToEmptyString(data []map[string]interface{}) []map[string]interface{} {
	for _, item := range data {
		for key, value := range item {
			if value == nil {
				item[key] = ""
			}
		}
	}
	return data
}
