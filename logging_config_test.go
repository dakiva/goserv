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

func TestValidLoggingConfig(t *testing.T) {
	// given
	config := &LoggingConfig{
		LogLevel: "DEBUG",
		Format:   "%{level} %{message}",
		Backends: []backendConfig{
			backendConfig{
				BackendName: "STDOUT",
			},
		},
	}

	// when
	err := config.Validate()

	// then
	assert.NoError(t, err)
}

func TestValidEmptyFormatLoggingConfig(t *testing.T) {
	// given
	config := &LoggingConfig{
		LogLevel: "DEBUG",
		Backends: []backendConfig{
			backendConfig{
				BackendName: "STDOUT",
			},
		},
	}

	// when
	err := config.Validate()

	// then
	assert.NoError(t, err)
}

func TestInvalidLogLevelLoggingConfig(t *testing.T) {
	// given
	config := &LoggingConfig{
		LogLevel: "foo",
		Backends: []backendConfig{
			backendConfig{
				BackendName: "STDOUT",
			},
		},
	}

	// when
	err := config.Validate()

	// then
	assert.Error(t, err)
}

func TestEmptyBackendsLoggingConfig(t *testing.T) {
	// given
	config := &LoggingConfig{
		LogLevel: "INFO",
	}

	// when
	err := config.Validate()

	// then
	assert.Error(t, err)
}

func TestInvalidBackendLoggingConfig(t *testing.T) {
	// given
	config := &LoggingConfig{
		LogLevel: "ERROR",
		Backends: []backendConfig{
			backendConfig{
				BackendName: "foo",
			},
		},
	}

	// when
	err := config.Validate()

	// then
	assert.Error(t, err)
}

func TestInvalidFileBackendLoggingConfig(t *testing.T) {
	// given
	config := &LoggingConfig{
		LogLevel: "WARNING",
		Backends: []backendConfig{
			backendConfig{
				BackendName: "FILE",
			},
		},
	}

	// when
	err := config.Validate()

	// then
	assert.Error(t, err)
}

func TestValidFileBackendLoggingConfig(t *testing.T) {
	// given
	config := &LoggingConfig{
		LogLevel: "NOTICE",
		Backends: []backendConfig{
			backendConfig{
				BackendName: "FILE",
				FilePath:    ".",
			},
		},
	}

	// when
	err := config.Validate()

	// then
	assert.NoError(t, err)
}
