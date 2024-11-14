package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// --------------------------------------------------- helper function ---------------------------------------------------
// toString function that accepts any struct and returns a string with all fields and their values.
func ToString(v interface{}) string {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return "toString function requires a struct type"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s: { ", val.Type().Name()))

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i).Interface()
		sb.WriteString(fmt.Sprintf("%s: %v, ", field.Name, value))
	}

	// Remove the last comma and space, and add closing bracket
	result := sb.String()
	return result[:len(result)-2] + " }"
}

// GetColumnNames builds a column list based on struct tags
func GetColumnNames(model interface{}) string {
	// Check if model is a pointer and get its underlying element if so
	typ := reflect.TypeOf(model)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// Ensure we're working with a struct
	if typ.Kind() != reflect.Struct {
		panic("model must be a struct or a pointer to a struct")
	}
	var columns []string

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if dbTag := field.Tag.Get("json"); dbTag != "" {
			columns = append(columns, dbTag)
		}
	}
	return strings.Join(columns, ", ")
}

// BuildSelectQuery dynamically builds a SELECT query using the struct fields
func BuildSelectQuery(tableName string, model interface{}, whereClause string) string {
	columns := GetColumnNames(model)
	query := fmt.Sprintf("SELECT %s FROM %s", columns, tableName)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}
	return query
}

// BuildInsertQuery dynamically builds an INSERT statement
func BuildInsertQuery(tableName string, model interface{}) (string, []interface{}) {
	typ := reflect.TypeOf(model)
	val := reflect.ValueOf(model)

	var columns []string
	var placeholders []string
	var values []interface{}

	for i := 0; i < typ.NumField(); i++ {
		if i == 0 {
			continue
		}
		field := typ.Field(i)
		if dbTag := field.Tag.Get("json"); dbTag != "" {
			columns = append(columns, dbTag)
			placeholders = append(placeholders, "?")
			values = append(values, val.Field(i).Interface())
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	return query, values
}

// BuildUpdateQuery dynamically builds an UPDATE SQL statement
func BuildUpdateQuery(tableName string, model interface{}, whereClause string) (string, []interface{}) {
	typ := reflect.TypeOf(model)
	val := reflect.ValueOf(model)

	// Ensure model is a pointer to a struct
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() != reflect.Struct {
		panic("model must be a struct or a pointer to a struct")
	}

	// Collect the columns and their new values based on the json tags
	var setClauses []string
	var args []interface{}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		dbTag := field.Tag.Get("json")   // Using "json" tag as you requested
		if dbTag != "" && dbTag != "-" { // Only consider fields with a "json" tag and not omitted
			// Skip the first field (index 0) as it is typically productId
			if i == 0 {
				continue
			}

			setClauses = append(setClauses, fmt.Sprintf("%s = ?", dbTag))
			args = append(args, val.Field(i).Interface())
		}
	}

	// Construct the SET part of the query
	setClause := strings.Join(setClauses, ", ")

	// Construct the full UPDATE query
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, setClause, whereClause)

	return query, args
}

// BuildDeleteQuery dynamically builds a DELETE SQL statement
func BuildDeleteQuery(tableName string, whereClause string) string {
	// Construct the DELETE query
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, whereClause)

	// Return the query
	return query
}
