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
	"errors"
	"fmt"
)

// EndpointConfig represents the root configuration for the service
type EndpointConfig struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

// GetHostAddress returns the host address host:port. If the host is empty, returns a leading ':'.
func (e *EndpointConfig) GetHostAddress() string {
	return fmt.Sprintf("%v:%d", e.Hostname, e.Port)
}

// Validate ensures the service config is valid
func (e *EndpointConfig) Validate() error {
	if e.Port <= 0 || e.Port > 65535 {
		return errors.New("port value must a specified valid number between 0 and 65535")
	}
	return nil
}
