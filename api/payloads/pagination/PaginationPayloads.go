package pagination

import (
	"fmt"
	"math"
	"strings"

	"reflect"

	"github.com/go-gota/gota/dataframe"
	"gorm.io/gorm"
)

type Pagination struct {
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	SortOf     string      `json:"sort_of"`
	SortBy     string      `json:"sort_by"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
}

func (p *Pagination) GetOffset() int {
	return p.GetPage() * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	return p.Page
}

func (p *Pagination) GetSortOf() string {
	if p.SortOf == "" {
		p.SortOf = "asc"
	}
	return p.SortOf
}

func (p *Pagination) GetSortBy() string {
	return p.SortBy
}

func Paginate(value interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	var sort string = ""
	if pagination.GetSortBy() != "" {
		sort = pagination.GetSortBy() + " " + pagination.GetSortOf()
	}
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(sort)
	}
}

func PaginateDistinct(value interface{}, pagination *Pagination, db *gorm.DB, distinctColumn string) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	var sort string = ""
	if pagination.GetSortBy() != "" {
		sort = pagination.GetSortBy() + " " + pagination.GetSortOf()
	}
	db.Model(value).Distinct(distinctColumn).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(sort)
	}
}

func toCamelCase(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, "")
}

func NewDataFramePaginate(rows any, pagination *Pagination) (result []map[string]interface{}, totalPages int, totalRows int) {
	var df dataframe.DataFrame
	tpy, _ := reflect.TypeOf(rows), reflect.ValueOf(rows)

	if tpy.Kind() == reflect.Slice && tpy.Elem().Kind() != reflect.Struct {
		df = dataframe.LoadMaps(rows.([]map[string]interface{}))
	} else {
		df = dataframe.LoadStructs(rows)
	}

	totalRows = df.Nrow()
	if pagination.GetSortBy() != "" {
		if pagination.GetSortBy() == "desc" {
			sortOf := pagination.GetSortOf()

			if strings.Contains(sortOf, "_") {
				sortOf = toCamelCase(sortOf)
			}
			dfSorted := df.Arrange(dataframe.RevSort(sortOf))

			if dfSorted.Err != nil {
				fmt.Println("Error sorting DataFrame in descending order:", dfSorted.Err)
			} else {
				df = dfSorted
			}

		} else {
			df = df.Arrange(dataframe.Sort(pagination.GetSortOf()))
		}
	}

	start := pagination.GetPage() * pagination.GetLimit()
	end := start + pagination.GetLimit()

	if end > df.Nrow() {
		end = df.Nrow()
	}

	indices := make([]int, end-start)
	for i := start; i < end; i++ {
		indices[i-start] = i
	}

	totalPages = int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))

	df = df.Subset(indices)
	return df.Maps(), totalPages, totalRows
}
