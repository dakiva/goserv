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
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEmptyDBConfig(t *testing.T) {
	// given
	config := DBConfig{}

	// when
	isEmpty := config.Empty()

	// then
	assert.True(t, isEmpty)
}

func TestNonEmptyDBConfig(t *testing.T) {
	// given
	config := DBConfig{
		Role: "role",
	}

	// when
	isEmpty := config.Empty()

	// then
	assert.False(t, isEmpty)
}

func TestValidateDBConfig(t *testing.T) {
	// given
	config := DBConfig{
		Port:         5432,
		DBName:       "db",
		SSLMode:      "disable",
		SchemaName:   "myschema",
		Role:         "role",
		RolePassword: "password",
	}

	// when
	err := config.Validate()

	// then
	assert.NoError(t, err)
}

func TestInvalidPortDBConfig(t *testing.T) {
	// given
	config := DBConfig{
		Port:         0,
		DBName:       "db",
		SSLMode:      "disable",
		SchemaName:   "myschema",
		Role:         "role",
		RolePassword: "password",
	}

	// when
	err := config.Validate()

	// then
	assert.Error(t, err)
}

func TestInvalidMaxIdleConnectionsDBConfig(t *testing.T) {
	// given
	config := DBConfig{
		Port:               5432,
		DBName:             "db",
		SSLMode:            "disable",
		SchemaName:         "myschema",
		Role:               "role",
		RolePassword:       "password",
		MaxIdleConnections: 10,
		MaxOpenConnections: 9,
	}

	// when
	err := config.Validate()

	// then
	assert.Error(t, err)
}

func TestInvalidSSLModeDBConfig(t *testing.T) {
	// given
	config := DBConfig{
		Port:         5432,
		DBName:       "db",
		SSLMode:      "foo",
		SchemaName:   "myschema",
		Role:         "role",
		RolePassword: "password",
	}

	// when
	err := config.Validate()

	// then
	assert.Error(t, err)
}

func TestDBConfigToDsn(t *testing.T) {
	// given
	config := &DBConfig{
		Hostname:       "localhost",
		Port:           5432,
		DBName:         "db",
		SSLMode:        "disable",
		SchemaName:     "myschema",
		Role:           "role",
		RolePassword:   "password",
		ConnectTimeout: 12,
	}

	// when
	dsn := config.ToDsn()

	// then
	assert.Equal(t, "host=localhost port=5432 user=role password=password dbname=db sslmode=disable connect_timeout=12", dsn)
}

func TestRequiredDBConfigToDsn(t *testing.T) {
	// given
	config := &DBConfig{
		Port:         5432,
		DBName:       "db",
		SSLMode:      "disable",
		SchemaName:   "myschema",
		Role:         "role",
		RolePassword: "secret",
	}

	// when
	dsn := config.ToDsn()

	// then
	assert.Equal(t, "port=5432 user=role password=secret dbname=db sslmode=disable", dsn)
}

func TestParseDBConfigEmptyDsn(t *testing.T) {
	// given
	dsnEnv := fmt.Sprintf("testdsn%d", time.Now().Unix())

	// when
	config, err := ParseDBConfig(dsnEnv)

	// then
	assert.NoError(t, err)
	assert.Nil(t, config)
}

func TestParseDBConfigBadEnv(t *testing.T) {
	// given
	dsnEnv := fmt.Sprintf("testdsn%d", time.Now().Unix())
	err := os.Setenv(dsnEnv, "host=localhost port=foo")
	assert.NoError(t, err)
	defer os.Unsetenv(dsnEnv)

	// when
	config, err := ParseDBConfig(dsnEnv)

	// then
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestParseDBConfig(t *testing.T) {
	// given
	dsnEnv := fmt.Sprintf("testdsn%d", time.Now().Unix())
	expectedConfig := &DBConfig{
		Hostname:           "localhost",
		Port:               5432,
		DBName:             "mydb",
		SSLMode:            "verify-full",
		ConnectTimeout:     40,
		SchemaName:         "myservice",
		Role:               "myservice",
		RolePassword:       "secret",
		MaxOpenConnections: DefaultMaxOpenConnections,
		MaxIdleConnections: DefaultMaxIdleConnections,
	}

	err := os.Setenv(dsnEnv, expectedConfig.ToDsn())
	assert.NoError(t, err)
	defer os.Unsetenv(dsnEnv)

	// when
	config, err := ParseDBConfig(dsnEnv)

	// then
	assert.NoError(t, err)
	assert.NoError(t, config.Validate())
	assert.Equal(t, expectedConfig, config)
	assert.Equal(t, expectedConfig.ToDsn(), config.ToDsn())
}
