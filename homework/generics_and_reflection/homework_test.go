package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

type Character struct {
	Level     int    `properties:"level"`
	Class     string `properties:"class"`
	Race      string `properties:"race"`
	Name      string `properties:"name"`
	Alignment string `properties:"alignment,omitempty"`
}

func Serialize(person Person) string {
	return serializeProperties(person)
}

func serializeProperties(data any) string {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return ""
		}

		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return ""
	}

	valueType := value.Type()
	properties := make([]string, 0, value.NumField())

	for i := range value.NumField() {
		fieldType := valueType.Field(i)
		tag := fieldType.Tag.Get("properties")
		if tag == "" || tag == "-" {
			continue
		}

		name, omitempty := parsePropertiesTag(tag)
		if name == "" {
			name = fieldType.Name
		}

		fieldValue := value.Field(i)
		if omitempty && fieldValue.IsZero() {
			continue
		}

		properties = append(properties, name+"="+formatPropertiesValue(fieldValue))
	}

	return strings.Join(properties, "\n")
}

func parsePropertiesTag(tag string) (string, bool) {
	parts := strings.Split(tag, ",")
	name := parts[0]

	for _, option := range parts[1:] {
		if option == "omitempty" {
			return name, true
		}
	}

	return name, false
}

func formatPropertiesValue(value reflect.Value) string {
	if !value.CanInterface() {
		return ""
	}

	return fmt.Sprint(value.Interface())
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestSerializePropertiesWithCharacter(t *testing.T) {
	tests := map[string]struct {
		character Character
		result    string
	}{
		"test case with empty alignment": {
			character: Character{
				Level: 1,
				Class: "Wizard",
				Race:  "Human",
				Name:  "Merlin",
			},
			result: "level=1\nclass=Wizard\nrace=Human\nname=Merlin",
		},
		"test case with alignment": {
			character: Character{
				Level:     5,
				Class:     "Paladin",
				Race:      "Human",
				Name:      "Arthur",
				Alignment: "Lawful Good",
			},
			result: "level=5\nclass=Paladin\nrace=Human\nname=Arthur\nalignment=Lawful Good",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := serializeProperties(test.character)
			assert.Equal(t, test.result, result)
		})
	}
}
