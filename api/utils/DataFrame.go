package utils

import (
	"math"
	"reflect"

	"github.com/go-gota/gota/dataframe"
)


func DataFramePaginate(data interface{}, page int, limit int, sortOf string, sortBy string) (result []map[string]interface{}, totalPages int, totalRows int) {
	var df dataframe.DataFrame
	tpy, _ := reflect.TypeOf(data), reflect.ValueOf(data)

	if tpy.Kind() == reflect.Slice && tpy.Elem().Kind() != reflect.Struct {
		df = dataframe.LoadMaps(data.([]map[string]interface{}))
	} else {
		df = dataframe.LoadStructs(data)
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
