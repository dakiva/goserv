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

// TokenAuthFilter represents a go-restful authentication filter that validates a JWT-Token passed via an http header parameter. There is an additional optionalwWhitelist that may be set, allowing specific urls, methods, or url patterns to opt-out of authentication.
type TokenAuthFilter struct {
	URLWhiteList *URLWhiteList
	TokenManager *TokenManager
}

// Filter fits in the go-restful filterchain that encapsulates a token based authentication check.
func (t *TokenAuthFilter) Filter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	if t.URLWhiteList != nil && t.URLWhiteList.Match(request.Request.URL.String(), request.Request.Method) {
		chain.ProcessFilter(request, response)
	} else {
		token := request.HeaderParameter("Authorization")
		if err := t.TokenManager.ValidateToken(token); err != nil {
			http.Error(response, "Not Authorized", http.StatusUnauthorized)
		} else {
			chain.ProcessFilter(request, response)
		}
	}
}
