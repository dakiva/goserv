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
	"os"
)

// LoadServiceConfig loads configuration from a file containing configuration in JSON.
func LoadServiceConfig(fileName string, output *ServiceConfig) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	return decoder.Decode(output)
}

// ServiceConfig represents a configuration suitable for fully configuration a service.
type ServiceConfig struct {
	Endpoint        *EndpointConfig        `json:"endpoint"`
	DB              *DBConfig              `json:"db"`
	MigrationDB     *DBConfig              `json:"migration_db"`
	Swagger         *SwaggerConfig         `json:"swagger"`
	Logging         *LoggingConfig         `json:"logging"`
	Custom          map[string]interface{} `json:"custom"`
	customValidator func(map[string]interface{}) error
}

// NewServiceConfig intializes a new instance
func NewServiceConfig() *ServiceConfig {
	return &ServiceConfig{Custom: make(map[string]interface{})}
}

// Validate validates a configuration, returning an error signaling invalid configuration
func (s *ServiceConfig) Validate() error {
	if s.Endpoint != nil {
		if err := s.Endpoint.Validate(); err != nil {
			return err
		}
	}
	if s.DB != nil {
		if err := s.DB.Validate(); err != nil {
			return err
		}
	}
	if s.MigrationDB != nil {
		if err := s.MigrationDB.Validate(); err != nil {
			return err
		}
	}
	if s.Swagger != nil {
		if err := s.Swagger.Validate(); err != nil {
			return err
		}
	}
	if s.customValidator != nil {
		return s.customValidator(s.Custom)
	}
	return nil
}

// SetCustomValidator sets a function that is executed during validation that validates any custom configuration
func (s *ServiceConfig) SetCustomValidator(f func(map[string]interface{}) error) {
	s.customValidator = f
}
