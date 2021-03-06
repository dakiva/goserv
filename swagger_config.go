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
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/emicklei/go-restful"
	openapi "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
)

// SwaggerConfig represents configuration to enable swagger documentation.
type SwaggerConfig struct {
	APIPath         string `json:"api_path"`
	SwaggerPath     string `json:"swagger_path"`
	SwaggerFilePath string `json:"swagger_file_path"`
}

// Validate ensures the configuration is valid
func (s *SwaggerConfig) Validate() error {
	if s.APIPath == "" {
		return errors.New("API path must be specified and point to the apidocs.json")
	}
	if s.SwaggerPath == "" {
		return errors.New("swagger path must be specified and be a relative path to the apidocs")
	}
	if s.SwaggerFilePath == "" {
		return errors.New("swagger file path must be specified and point to the swagger distribution directory")
	}
	if dir, err := os.Lstat(s.SwaggerFilePath); err != nil || !dir.IsDir() {
		return fmt.Errorf("swagger file path must point to a valid directory: %w", err)
	}
	return nil
}

// InstallSwaggerService sets up and installs the swagger service
func (s *SwaggerConfig) InstallSwaggerService(info *spec.Info, securityDefinitions spec.SecurityDefinitions, security []map[string][]string, container *restful.Container) {
	config := openapi.Config{
		WebServices: container.RegisteredWebServices(),
		APIPath:     s.APIPath,
		PostBuildSwaggerObjectHandler: func(swo *spec.Swagger) {
			swo.Info = info
			swo.SecurityDefinitions = securityDefinitions
			swo.Security = security
		}}
	container.Add(openapi.NewOpenAPIService(config))
	container.ServeMux.Handle(s.SwaggerPath, http.StripPrefix(s.SwaggerPath, http.FileServer(http.Dir(s.SwaggerFilePath))))
}
