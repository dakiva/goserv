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

// OpenIDConnectClientConfig represents configuration for setting up an OpenID Connect client
type OpenIDConnectClientConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Issuer       string `json:"issuer"`
	RedirectURL  string `json:"redirect_url"`
}

// Validate validates whether the configuration is valid for an OpenID Connect client
func (o *OpenIDConnectClientConfig) Validate() error {
	if o.ClientID == "" {
		return errors.New("client ID must be set")
	}
	if o.ClientSecret == "" {
		return errors.New("client secret must be set")
	}
	if o.Issuer == "" {
		return errors.New("issuer must be set")
	}
	if o.RedirectURL == "" {
		return errors.New("redirect URL must be set")
	}
	return nil
}
