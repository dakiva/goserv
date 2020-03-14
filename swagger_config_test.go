// Copyright 2020 Daniel Akiva

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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSwaggerConfig(t *testing.T) {
	// given
	config := &SwaggerConfig{
		APIPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: ".",
	}

	// when
	err := config.Validate()

	// then
	assert.NoError(t, err)
}

func TestInvalidAPIPathSwaggerConfig(t *testing.T) {
	// given
	config := &SwaggerConfig{
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "/home/test/dist",
	}

	// when
	err := config.Validate()

	// then
	assert.Error(t, err)
}

func TestInvalidSwaggerPathSwaggerConfig(t *testing.T) {
	// given
	config := &SwaggerConfig{
		APIPath:         "/apidocs.json",
		SwaggerFilePath: "/home/test/dist",
	}

	// when
	err := config.Validate()

	// then
	assert.Error(t, err)
}

func TestInvalidSwaggerFilePathSwaggerConfig(t *testing.T) {
	// given
	config := &SwaggerConfig{
		APIPath:     "/apidocs.json",
		SwaggerPath: "/apidocs/",
	}

	// when
	err := config.Validate()

	// then
	assert.Error(t, err)
}
