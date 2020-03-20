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

func TestLoadServiceConfig(t *testing.T) {
	// given
	expectedConfig := &ServiceConfig{}
	expectedConfig.Endpoint = &EndpointConfig{Hostname: "localhost", Port: 8080}
	expectedConfig.DB = &DBConfig{
		Hostname:           "dbhost",
		Port:               5432,
		SchemaName:         "myschema",
		Role:               "role",
		RolePassword:       "secret",
		DBName:             "mydb",
		SSLMode:            "require",
		MaxOpenConnections: 10,
		MaxIdleConnections: 5,
	}
	expectedConfig.MigrationDB = &DBConfig{
		Hostname:           "dbhost",
		Port:               5432,
		SchemaName:         "migrationschema",
		Role:               "migrationrole",
		RolePassword:       "secret",
		DBName:             "mydb",
		SSLMode:            "require",
		MaxOpenConnections: 10,
		MaxIdleConnections: 5,
	}
	expectedConfig.Swagger = &SwaggerConfig{
		APIPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: ".",
	}
	expectedConfig.Logging = &LoggingConfig{
		LogLevel: "DEBUG",
		Format:   "%{time} %{shortfile} %{level} %{message}",
		Backends: []BackendConfig{
			BackendConfig{
				BackendName: "FILE",
				FilePath:    "/home/centos/temp.log",
			},
		},
	}
	expectedConfig.OAuth2Service = &OAuth2ServiceConfig{
		AuthorizationURL:      "http://localhost/authorize",
		TokenURL:              "http://localhost/token",
		AccessTokenExpiration: 3600,
		AccessTokenPrivateKey: "secret",
	}
	expectedConfig.OpenIDConnectClient = &OpenIDConnectClientConfig{
		ClientID:     "client123",
		ClientSecret: "secret",
		Issuer:       "http://account.example.com",
		RedirectURL:  "http://localhost/authorize/callback",
	}
	// when
	config := &ServiceConfig{}
	err := LoadServiceConfig("service_config_test.json", config)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, config)

	err = config.Validate()
	assert.NoError(t, err)
}
