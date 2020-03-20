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

import "errors"

// OAuth2ServiceConfig represents configuration for setting up OAuth2 service flows
type OAuth2ServiceConfig struct {
	AuthorizationURL      string `json:"authorization_url"`
	TokenURL              string `json:"token_url"`
	AccessTokenExpiration int    `json:"access_token_expiration"`
	AccessTokenPrivateKey string `json:"access_token_private_key"`
}

// Validate validates whether the configuration is valid for an OAuth2 service flow (redirect/grant and implicit flows supported)
func (o *OAuth2ServiceConfig) Validate() error {
	if o.AuthorizationURL == "" {
		return errors.New("authorization URL must be set")
	}
	// TokenURL is optional to support OAuth2 implicit flows
	if o.AccessTokenExpiration <= 0 {
		return errors.New("access token expiration must be set to a positive integer (seconds)")
	}
	if o.AccessTokenPrivateKey == "" {
		return errors.New("access token private key must be set")
	}
	return nil
}
