package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPointerValue_NilValues(t *testing.T) {
	type x struct{}
	var nilInt *int
	var nilString *string
	var nilStruct *x

	zeroInt := GetPointerValue(nilInt)
	assert.Equal(t, zeroInt, int(0))

	zeroString := GetPointerValue(nilString)
	assert.Equal(t, zeroString, "")

	zeroStruct := GetPointerValue(nilStruct)
	y := x{}
	assert.Equal(t, zeroStruct, y)
}

func TestGetPointerValue_Pointed(t *testing.T) {
	type x struct{ abc int }
	var valInt int = 123
	var valString string = "123"
	var valStruct x = x{abc: 123}

	vv := GetPointerValue(&valInt)
	assert.Equal(t, vv, int(123))

	vvstring := GetPointerValue(&valString)
	assert.Equal(t, vvstring, "123")

	vvstruct := GetPointerValue(&valStruct)
	y := x{abc: 123}
	assert.Equal(t, vvstruct, y)
}

func TestContains(t *testing.T) {
	testCases := []struct {
		name     string
		source   []any
		target   any
		expected bool
	}{
		{"int found", []any{1, 2, 3}, 2, true},
		{"int not found", []any{1, 2, 3}, 4, false},
		{"string found", []any{"apple", "banana", "cherry"}, "banana", true},
		{"string not found", []any{"apple", "banana", "cherry"}, "grape", false},
		{"empty slice", []any{}, 1, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Contains(tc.source, tc.target)
			assert.Equal(t, tc.expected, result, "error expected %v, but actual: %v", tc.expected, result)
		})
	}
}

func TestToPointerValue(t *testing.T) {
	type Custom struct {
		Field any
	}

	testCases := []struct {
		name     string
		input    any
		expected any
	}{
		{"int value", 123, 123},
		{"string value", "string test", "string test"},
		{"bool value", true, true},
		{"custom struct value", Custom{Field: 10}, Custom{Field: 10}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ToPointerValue(tc.input)
			assert.Equal(t, tc.expected, *result)
		})
	}
}
