package utils

import (
	"encoding/json"
	"sort"
	"strconv"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm/schema"
)

func SetExcelHeader(f *excelize.File, sheetName string, s *schema.Schema) {
	colNum := 1
	for _, field := range s.Fields {
		headerName := field.DBName
		col, _ := excelize.ColumnNumberToName(colNum)
		f.SetCellValue(sheetName, col+"1", headerName)
		colNum++
	}
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#CCCCCC"}, Pattern: 1},
	})
	colTo, _ := excelize.ColumnNumberToName(colNum - 1)
	f.SetCellStyle(sheetName, "A1", colTo+"1", style)
}

func WriteExcel(f *excelize.File, sheetName string, data interface{}) {
	b, _ := json.Marshal(data)
	var jsonMap []map[string]interface{}
	json.Unmarshal([]byte(b), &jsonMap)
	for i := 0; i < len(jsonMap); i++ {
		row := jsonMap[i]
		var keys []string
		for k := range row {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		j := 1
		for _, k := range keys {
			col, _ := excelize.ColumnNumberToName(j)
			f.SetCellValue(sheetName, col+strconv.Itoa(i+2), row[k])
			j++
		}
	}
}
