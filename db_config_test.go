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

func TestValidateDBConfig(t *testing.T) {
	// given
	config := DBConfig{
		DBName:       "db",
		SSLMode:      "disable",
		SchemaName:   "myschema",
		Role:         "role",
		RolePassword: "password",
	}

	// when
	err := config.Validate()

	// then
	assert.Nil(t, err)
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
		DBName:       "db",
		SSLMode:      "disable",
		SchemaName:   "myschema",
		Role:         "role",
		RolePassword: "secret",
	}

	// when
	dsn := config.ToDsn()

	// then
	assert.Equal(t, "user=role password=secret dbname=db sslmode=disable", dsn)
}
