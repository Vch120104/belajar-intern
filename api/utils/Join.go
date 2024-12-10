package utils

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type JoinTable struct {
	Table         string
	Alias         string
	ForeignKey    string
	ReferenceKey  string
	JoinCondition string
}

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
// func CreateJoinSelectStatement(db *gorm.DB, tableStruct interface{}) *gorm.DB {
// 	keyAttribute := []string{}
// 	responseType := reflect.TypeOf(tableStruct)
// 	joinTable := []string{}
// 	var joinTableId []string
// 	var mainTable string
// 	referenceTable := []string{}

// 	// Define main table
// 	for i := 0; i < responseType.NumField(); i++ {
// 		mainTable = responseType.Field(i).Tag.Get("main_table")
// 		if mainTable != "" {
// 			break
// 		}
// 	}

// 	// Define reference tables
// 	for i := 0; i < responseType.NumField(); i++ {
// 		if responseType.Field(i).Tag.Get("references") != "" {
// 			referenceTable = append(referenceTable, responseType.Field(i).Tag.Get("references"))
// 		}
// 	}

// 	// Define select from table and Join table id
// 	for i := 0; i < responseType.NumField(); i++ {
// 		for _, value := range referenceTable {
// 			if value == responseType.Field(i).Tag.Get("parent_entity") && strings.Contains(responseType.Field(i).Tag.Get("json"), "id") {
// 				joinTableId = append(joinTableId, responseType.Field(i).Tag.Get("json"))
// 			}
// 		}
// 		keyAttribute = append(keyAttribute, responseType.Field(i).Tag.Get("parent_entity")+"."+responseType.Field(i).Tag.Get("json"))
// 	}

// 	// Query Table with select
// 	query := db.Table(mainTable).Select(keyAttribute)

// 	// Join Tables
// 	if len(joinTableId) > 0 {
// 		for i := 0; i < len(referenceTable); i++ {
// 			id := joinTableId[i]
// 			reference := referenceTable[i]
// 			joinTable = append(joinTable, "join "+reference+" as "+reference+" on "+mainTable+"."+id+" = "+reference+"."+id)
// 		}
// 		query = query.Joins(strings.Join(joinTable, " "))
// 	} else {
// 		fmt.Print("Please troubleshoot tableStruct")
// 	}

//		return query
//	}
func CreateJoinSelectStatement(db *gorm.DB, tableStruct interface{}) *gorm.DB {
	keyAttribute := []string{}
	responseType := reflect.TypeOf(tableStruct)
	joinTable := []string{}
	joinTableMap := make(map[string]string)
	var mainTable string
	referenceTable := map[string]bool{} // Use map to store unique referenced tables

	// Define primary table
	for i := 0; i < responseType.NumField(); i++ {
		mainTable = responseType.Field(i).Tag.Get("main_table")
		if mainTable != "" {
			break
		}
	}

	if mainTable == "" {
		fmt.Println("Please specify main_table in the struct tags")
		return nil
	}

	// Define join reference tables
	for i := 0; i < responseType.NumField(); i++ {
		ref := responseType.Field(i).Tag.Get("references")
		if ref != "" {
			referenceTable[ref] = true // Store unique referenced tables in the map
		}
	}

	// Define select from table and join table id
	for i := 0; i < responseType.NumField(); i++ {
		for ref := range referenceTable {
			if ref == responseType.Field(i).Tag.Get("parent_entity") && strings.Contains(responseType.Field(i).Tag.Get("json"), "id") {
				joinTableMap[responseType.Field(i).Tag.Get("parent_entity")] = responseType.Field(i).Tag.Get("json")
			}
		}
		keyAttribute = append(keyAttribute, responseType.Field(i).Tag.Get("parent_entity")+"."+responseType.Field(i).Tag.Get("json"))
	}

	// Query Table with select
	query := db.Table(mainTable).Select(keyAttribute)

	// Join Tables
	for ref := range referenceTable {
		joinCondition := "join " + ref + " as " + ref + " on " + mainTable + "." + joinTableMap[ref] + " = " + ref + "." + joinTableMap[ref]
		joinTable = append(joinTable, joinCondition)
	}
	query = query.Joins(strings.Join(joinTable, " "))

	return query
}

func CreateJoinSelectStatementTransaction(db *gorm.DB, tableStruct interface{}) *gorm.DB {
	keyAttribute := []string{}
	responseType := reflect.TypeOf(tableStruct)
	joinTable := []string{}
	joinTableMap := make(map[string]string)
	var mainTable string
	referenceTable := map[string]bool{} // Use map to store unique referenced tables

	// Define primary table
	for i := 0; i < responseType.NumField(); i++ {
		mainTable = responseType.Field(i).Tag.Get("main_table")
		if mainTable != "" {
			break
		}
	}

	if mainTable == "" {
		fmt.Println("Please specify main_table in the struct tags")
		return nil
	}

	// Define join reference tables
	for i := 0; i < responseType.NumField(); i++ {
		ref := responseType.Field(i).Tag.Get("references")
		if ref != "" {
			referenceTable[ref] = true // Store unique referenced tables in the map
		}
	}

	// Define select from table and join table id
	for i := 0; i < responseType.NumField(); i++ {
		for ref := range referenceTable {
			if ref == responseType.Field(i).Tag.Get("parent_entity") && strings.Contains(responseType.Field(i).Tag.Get("json"), "system_number") {
				joinTableMap[responseType.Field(i).Tag.Get("parent_entity")] = responseType.Field(i).Tag.Get("json")
			}
		}
		keyAttribute = append(keyAttribute, responseType.Field(i).Tag.Get("parent_entity")+"."+responseType.Field(i).Tag.Get("json"))
	}

	// Query Table with select
	query := db.Table(mainTable).Select(keyAttribute)

	// Join Tables
	for ref := range referenceTable {
		joinCondition := "join " + ref + " as " + ref + " on " + mainTable + "." + joinTableMap[ref] + " = " + ref + "." + joinTableMap[ref]
		joinTable = append(joinTable, joinCondition)
	}
	query = query.Joins(strings.Join(joinTable, " "))

	return query
}
