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
