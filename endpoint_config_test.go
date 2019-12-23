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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateBadPort(t *testing.T) {
	// given
	config := &EndpointConfig{
		Port: -1,
	}

	// when
	err := config.Validate()

	// then
	assert.NotNil(t, err)
}

func TestGetHostAddress(t *testing.T) {
	// given
	config := &EndpointConfig{
		Port: 8080,
	}

	// when
	address := config.GetHostAddress()

	// then
	assert.Equal(t, ":8080", address)

	config.Hostname = "host"
	address = config.GetHostAddress()
	assert.Equal(t, "host:8080", address)
}
