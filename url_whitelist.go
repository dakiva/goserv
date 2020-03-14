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

import "strings"

// URLWhiteList represents a list of URLs and acceptable request methods per URL that can be matched against. A match is defined as one of the following
// a) the request method itself matches the global acceptAllMethods: (ie match  "All OPTIONS requests")
// b) the url is an exact match and the request method matches one of the accepted methds mapped to the url. Url matching is done by matching each segment. Segments are separated by the slash character. '/'
// Url path variable tokens are supported, defaulted to using the : prefix to denote the path variable. Path variables match any value for the specific url segment in the Url. This can be customized by supplying custom Prefix and, optionally a custom Suffix
// All matching is case insensitve. Query parameters are stripped off prior to matching.
type URLWhiteList struct {
	urlMappings      map[string][]string
	acceptAllMethods []string
	Prefix           string
	Suffix           string
}

// NewURLWhiteList returns a new whitelist
func NewURLWhiteList() *URLWhiteList {
	list := &URLWhiteList{urlMappings: make(map[string][]string), acceptAllMethods: make([]string, 0)}
	list.Prefix = ":"
	return list
}

// AddURL adds an URL and its mapped request methods. If no methods are specified, any request method is allowed.
func (u *URLWhiteList) AddURL(url string, requestMethods ...string) {
	normalizedMethods := make([]string, len(requestMethods))
	for i := 0; i < len(requestMethods); i++ {
		normalizedMethods[i] = strings.ToLower(requestMethods[i])
	}
	u.urlMappings[normalizeURL(url)] = normalizedMethods
}

// AcceptAll accepts all URLs for request methods. (ie if passing "OPTIONS", any options request is whitelisted)
func (u *URLWhiteList) AcceptAll(requestMethods ...string) {
	normalizedMethods := make([]string, len(requestMethods))
	for i := 0; i < len(requestMethods); i++ {
		normalizedMethods[i] = strings.ToLower(requestMethods[i])
	}
	u.acceptAllMethods = normalizedMethods
}

// Match returns true if the url and requestMethod specified is a match, false otherwise.
func (u *URLWhiteList) Match(url string, requestMethod string) bool {
	normalizedMethod := strings.ToLower(requestMethod)
	if methodMatch(u.acceptAllMethods, normalizedMethod) {
		return true
	}
	url = normalizeURL(url)
	// first try brute force match
	if methods, ok := u.urlMappings[url]; ok {
		return len(methods) == 0 || methodMatch(methods, normalizedMethod)
	}
	urlValues := strings.Split(url, "/")
	numValues := len(urlValues)
	for k, methods := range u.urlMappings {
		tokens := strings.Split(k, "/")
		if numValues != len(tokens) {
			continue
		}
		// loop until a segment is not matched, or all segments match
		allMatch := true
		for i := 0; i < numValues && allMatch; i++ {
			// if a token sgement is a path variable, denoted by a : prefix, any url segment matches
			allMatch = (strings.HasPrefix(tokens[i], u.Prefix) && strings.HasSuffix(tokens[i], u.Suffix)) || urlValues[i] == tokens[i]
		}
		if allMatch {
			// if no methods specified, accept any method
			return len(methods) == 0 || methodMatch(methods, normalizedMethod)
		}
	}
	return false
}

func methodMatch(requestMethods []string, requestMethod string) bool {
	for _, m := range requestMethods {
		if requestMethod == m {
			return true
		}
	}
	return false
}

func normalizeURL(url string) string {
	idx := strings.Index(url, "?")
	if idx >= 0 {
		url = url[:idx]
	}
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}
	return strings.ToLower(url)
}
