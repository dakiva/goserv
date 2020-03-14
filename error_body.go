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
	"net/http"

	"github.com/emicklei/go-restful"
)

// ErrorBody struct is used when constructing a single error response body.
type ErrorBody struct {
	ErrorMessage string `json:"error" description:"The error message."`
	StatusCode   int    `json:"status_code" description:"The status code."`
}

// WriteError ensures the error is appropriately handled by ensuring the correct http status code is assigned and the error message is logged.
func WriteError(response *restful.Response, err error) {
	errStatusCode := http.StatusInternalServerError
	if httpErr, ok := err.(httpError); ok {
		errStatusCode = httpErr.StatusCode()
	}
	errBody := &ErrorBody{ErrorMessage: err.Error(), StatusCode: errStatusCode}
	response.WriteHeaderAndJson(errStatusCode, errBody, restful.MIME_JSON)
}
