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
	"net/url"
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
)

func TestWhiteListedURL(t *testing.T) {
	// given
	urlList := NewURLWhiteList()
	urlList.AddURL("/", "GET")
	authFilter := &TokenAuthFilter{
		URLWhiteList: urlList,
	}

	url, _ := url.Parse("/")
	req := restful.NewRequest(&http.Request{
		URL:    url,
		Method: "GET",
	})

	// when
	targetFuncCalled := false
	chain := &restful.FilterChain{
		Target: func(*restful.Request, *restful.Response) {
			targetFuncCalled = true
		},
	}
	authFilter.Filter(req, nil, chain)

	// then
	assert.True(t, targetFuncCalled)
}

func TestInvalidAccess(t *testing.T) {
	// given
	urlList := NewURLWhiteList()
	urlList.AddURL("/", "GET")
	authFilter := &TokenAuthFilter{
		URLWhiteList: urlList,
	}

	url, _ := url.Parse("/home")
	req := restful.NewRequest(&http.Request{
		URL:    url,
		Method: "GET",
	})
	resp := &restful.Response{
		ResponseWriter: NoopResponseWriter{},
	}

	// when
	targetFuncCalled := false
	chain := &restful.FilterChain{
		Target: func(*restful.Request, *restful.Response) {
			targetFuncCalled = true
		},
	}
	authFilter.Filter(req, resp, chain)

	// then
	assert.False(t, targetFuncCalled)
}

func TestValidAccess(t *testing.T) {
	// given
	urlList := NewURLWhiteList()
	urlList.AddURL("/", "GET")
	manager := NewTokenManager([]byte("8831dcf1c522debbdc187f909f52b743f0028777c29517ab12938a624fc4ed12"), 60)
	token, err := manager.CreateToken()
	assert.NoError(t, err)
	authFilter := &TokenAuthFilter{
		URLWhiteList: urlList,
		TokenManager: manager,
	}

	url, _ := url.Parse("/home")
	req := restful.NewRequest(&http.Request{
		URL:    url,
		Method: "GET",
	})
	req.Request.Header = make(map[string][]string, 0)
	req.Request.Header.Add("Authorization", token)

	// when
	targetFuncCalled := false
	chain := &restful.FilterChain{
		Target: func(*restful.Request, *restful.Response) {
			targetFuncCalled = true
		},
	}
	authFilter.Filter(req, nil, chain)

	// then
	assert.True(t, targetFuncCalled)
}

type NoopResponseWriter struct {
	headers map[string][]string
}

func (n NoopResponseWriter) Header() http.Header {
	if n.headers == nil {
		n.headers = make(map[string][]string, 0)
	}
	return n.headers
}

func (n NoopResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (n NoopResponseWriter) WriteHeader(int) {
}
