// Copyright 2019 Daniel Akiva

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
