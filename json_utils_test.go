package goserv

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertNilToFloat(t *testing.T) {
	// when
	f, err := ConvertToFloat(nil)

	// then
	assert.Error(t, err)
	assert.Empty(t, f)
}

func TestConvertBadTypeToFloat(t *testing.T) {
	val := "12345.55"

	// when
	f, err := ConvertToFloat(val)

	// then
	assert.Error(t, err)
	assert.Empty(t, f)
}

func TestConvertFloatValToFloat(t *testing.T) {
	expected := 12404045.55

	// when
	f, err := ConvertToFloat(expected)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expected, f)
}

func TestConvertIntValToFloat(t *testing.T) {
	expected := 12404045

	// when
	f, err := ConvertToFloat(expected)

	// then
	assert.NoError(t, err)
	assert.Equal(t, float64(expected), f)
}

func TestConvertNumberValToFloat(t *testing.T) {
	expected := 12404045.55
	num := json.Number("12404045.55")

	// when
	f, err := ConvertToFloat(num)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expected, f)
}

func TestConvertNilToInteger(t *testing.T) {
	// when
	f, err := ConvertToInteger(nil)

	// then
	assert.Error(t, err)
	assert.Empty(t, f)
}

func TestConvertBadTypeToInteger(t *testing.T) {
	val := "12345.55"

	// when
	f, err := ConvertToInteger(val)

	// then
	assert.Error(t, err)
	assert.Empty(t, f)
}

func TestConvertFloatValToInteger(t *testing.T) {
	val := float64(12404045)
	expected := 12404045

	// when
	f, err := ConvertToInteger(val)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expected, f)
}

func TestConvertIntValToInteger(t *testing.T) {
	expected := 12404045

	// when
	f, err := ConvertToInteger(expected)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expected, f)
}

func TestConvertFloatValToIntegerWithFlooring(t *testing.T) {
	val := 12404045.55
	expected := 12404045

	// when
	f, err := ConvertToInteger(val)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expected, f)
}

func TestConvertNumberValToInteger(t *testing.T) {
	expected := 12404045
	num := json.Number("12404045.55")

	// when
	f, err := ConvertToInteger(num)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expected, f)
}
