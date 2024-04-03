package utils

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// CreateJoinSelectStatement generates a GORM database query for selecting columns from a main table and joining with reference tables.
//
// Parameters:
//   - db: A pointer to a GORM database query to which the join and select operations will be applied.
//   - tableStruct: An instance of a structure representing the main table's fields and their attributes, including reference table information.
//
// Returns:
//   - result: A modified GORM database query with select and join operations based on the provided main table structure and reference table information.
//
// Example Usage:
//
//	type User struct {
//	    ID           int    `json:"id" main_table:"users"`
//	    Name         string `json:"name"`
//	    ParentEntity string `json:"parent_entity" references:"organizations"`
//	}
//	tableStruct := User{}
//	query := CreateJoinSelectStatement(db, tableStruct)
func CreateJoinSelectStatement(db *gorm.DB, tableStruct interface{}) *gorm.DB {
	keyAttribute := []string{}
	responseType := reflect.TypeOf(tableStruct)
	joinTable := []string{}
	var joinTableId []string
	var mainTable string
	referenceTable := []string{}

	//define primary key table
	for i := 0; i < responseType.NumField(); i++ {
		mainTable = responseType.Field(i).Tag.Get("main_table")
		if mainTable != "" {
			break
		}
	}

	//deifne join table reference
	for i := 0; i < responseType.NumField(); i++ {
		if responseType.Field(i).Tag.Get("references") != "" {
			referenceTable = append(referenceTable, responseType.Field(i).Tag.Get("references"))
		}
	}

	//define select from table and Join table id
	for i := 0; i < responseType.NumField(); i++ {
		for _, value := range referenceTable {
			if value == responseType.Field(i).Tag.Get("parent_entity") && strings.Contains(responseType.Field(i).Tag.Get("json"), "id") {
				joinTableId = append(joinTableId, responseType.Field(i).Tag.Get("json"))
			}
		}
		keyAttribute = append(keyAttribute, responseType.Field(i).Tag.Get("parent_entity")+"."+responseType.Field(i).Tag.Get("json"))
	}

	//query Table with select
	query := db.Table(mainTable).Select(keyAttribute)

	//join Table
	if len(joinTableId) > 0 {
		for i := 0; i < len(referenceTable); i++ {
			id := joinTableId[i]
			referenceTable := referenceTable[i]
			joinTable = append(joinTable, "join "+referenceTable+" on "+mainTable+"."+id+" = "+referenceTable+"."+id)
		}
		query = query.Joins(strings.Join(joinTable, " "))
	} else {
		fmt.Print("Please troubleshoot tableStruct")
	}

	return query
}

// tableStructs := []interface{}{User{}, Organization{}}
// query := CreateJoinManyTable(db, tableStructs)
func CreateJoinManyTable(db *gorm.DB, tables []interface{}) *gorm.DB {
	var joinConditions []string

	for _, tableStruct := range tables {
		responseType := reflect.TypeOf(tableStruct)
		mainTable := ""
		referenceTable := ""

		// Define main table
		for i := 0; i < responseType.NumField(); i++ {
			mainTable = responseType.Field(i).Tag.Get("main_table")
			if mainTable != "" {
				break
			}
		}

		// Define reference table
		for i := 0; i < responseType.NumField(); i++ {
			if responseType.Field(i).Tag.Get("references") != "" {
				referenceTable = responseType.Field(i).Tag.Get("references")
				break
			}
		}

		if mainTable != "" && referenceTable != "" {
			for i := 0; i < responseType.NumField(); i++ {
				joinConditions = append(joinConditions, fmt.Sprintf("%s.%s = %s.%s", mainTable, responseType.Field(i).Tag.Get("json"), referenceTable, responseType.Field(i).Tag.Get("json")))
			}
		}
	}

	// Query table with select
	query := db

	// Join tables
	if len(joinConditions) > 0 {
		joinStatement := strings.Join(joinConditions, " AND ")
		for _, condition := range joinConditions {
			query = query.Joins(fmt.Sprintf("JOIN %s ON %s", condition, joinStatement))
		}
	} else {
		fmt.Println("Please verify the tableStructs")
	}

	return query
}
