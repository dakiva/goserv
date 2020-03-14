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

func TestMatchByAcceptAllMethod(t *testing.T) {
	// given
	list := NewURLWhiteList()
	list.AcceptAll("OPTIONS")

	// when
	match := list.Match("/foo", "OPTIONS")

	// then
	assert.True(t, match)
}

func TestMatchURL(t *testing.T) {
	// given
	list := NewURLWhiteList()
	list.AddURL("/foo/Bar", "POST", "GET")
	list.AddURL("/bar", "PUT")

	// when
	match := list.Match("/foo/bar", "GET")

	// then
	assert.True(t, match)
}

func TestMatchURLWithQuery(t *testing.T) {
	// given
	list := NewURLWhiteList()
	list.AddURL("/foo/Bar", "POST", "GET")
	list.AddURL("/bar", "PUT")

	// when
	match := list.Match("/foo/bar?param=value&param2=value2", "GET")

	// then
	assert.True(t, match)
}

func TestURLWithPathVariable(t *testing.T) {
	// given
	list := NewURLWhiteList()
	list.AddURL("/foo/:id/bar/:id2", "POST")

	// when
	match := list.Match("/foo/123/bar/abc123/", "post")

	// then
	assert.True(t, match)
}

func TestURLWithCustomPathVariable(t *testing.T) {
	// given
	list := NewURLWhiteList()
	list.Prefix = "{"
	list.Suffix = "}"
	list.AddURL("/foo/{id}/bar/{id2}", "POST")

	// when
	match := list.Match("/foo/123/bar/abc123/", "post")

	// then
	assert.True(t, match)
}

func TestURLWithPathVariableAndQuery(t *testing.T) {
	// given
	list := NewURLWhiteList()
	list.AddURL("/foo/:id/bar/:id2", "POST")

	// when
	match := list.Match("/foo/123/bar/abc123/?param=value&param2=value2", "post")

	// then
	assert.True(t, match)
}

func TestURLMatchingAllRequestMethods(t *testing.T) {
	// given
	list := NewURLWhiteList()
	list.AddURL("/foo")

	// when
	match := list.Match("/foo", "POST")

	// then
	assert.True(t, match)
}

func TestPathVariableURLMatchingAllRequestMethods(t *testing.T) {
	// given
	list := NewURLWhiteList()
	list.AddURL("/foo/:id")

	// when
	match := list.Match("/foo/123", "POST")

	// then
	assert.True(t, match)
}

func TestNonMatchingURL(t *testing.T) {
	// given
	list := NewURLWhiteList()
	list.AddURL("/foo", "GET")

	// when
	match := list.Match("/bar", "GET")

	// then
	assert.False(t, match)
}

func TestNonMatchingPathVariable(t *testing.T) {
	// given
	list := NewURLWhiteList()
	list.AddURL("/foo/:id", "POST")

	// when
	match := list.Match("/foo", "POST")

	// then
	assert.False(t, match)
}

func TestNonMatchingRequestMethod(t *testing.T) {
	// given
	list := NewURLWhiteList()
	list.AddURL("/foo/:id/bar/:id2", "GET", "POST")

	// when
	match := list.Match("/foo/123/bar/abc123", "PUT")

	// then
	assert.False(t, match)
}

func TestNonMatchingSegmentLength(t *testing.T) {
	// given
	list := NewURLWhiteList()
	list.AddURL("/foo/:id", "POST")

	// when
	match := list.Match("/foo/123/bar", "POST")

	// then
	assert.False(t, match)
}
