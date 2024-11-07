package session

import (
	"encoding/gob"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/gorilla/sessions"
)

// variable declaration
var (
	// global
	Store *sessions.CookieStore

	// local

)

// type declaration
type (
	M      map[string]interface{}
	Person struct {
		FirstName string
		LastName  string
		Email     string
		Age       int
	}

	User struct {
		FirstName string
		LastName  string
		Email     string
		Age       int
		Role      string
	}
)

// const (
// 	SESSION_KEY = "asdaskdhasdhgsajdgasdsadksakdhasidoajsdousahdopj"
// )

func Init() {
	registerTypes()

	// os.Setenv("SESSION_KEY", SESSION_KEY)
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	store.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
	}

	Store = store
}

// --------------------------------------------------- helper function ---------------------------------------------------
// toString function that accepts any struct and returns a string with all fields and their values.
func toString(v interface{}) string {
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

func registerTypes() {
	gob.Register(&Person{})
	gob.Register(&User{})
	gob.Register(&M{})
}
