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
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/google/uuid"
	"github.com/op/go-logging"
)

const (
	// DefaultRequestLogFormat format
	DefaultRequestLogFormat = "[Request %s %s] trace=%s content_length=%d form=%s headers=[%s]\n"
	// DefaultResponseLogFormat format
	DefaultResponseLogFormat = "[Response %s %s] trace=%s status=%d time=%s headers=[%s]\n"

	// RequestTimestampAttribute is an attribute representing the timestamp at which the request was received by the service
	RequestTimestampAttribute = "request_timestamp_attr"
	// TraceIDAttribute is an attribute representing a unique identifier for the request useful for correlating log statements
	TraceIDAttribute = "trace_id_attr"
)

// RestfulLoggingFilter is middleware that logs restful requests and response payloads
type RestfulLoggingFilter struct {
	logger            *logging.Logger
	requestLogFormat  string
	responseLogFormat string
	logRoot           bool
}

// NewRestfulLoggingFilter initializes a new logging filter instance
func NewRestfulLoggingFilter(logger *logging.Logger) *RestfulLoggingFilter {
	return &RestfulLoggingFilter{
		logger:            logger,
		requestLogFormat:  DefaultRequestLogFormat,
		responseLogFormat: DefaultResponseLogFormat,
		logRoot:           false,
	}
}

// SetRequestLogFormat sets the verbose request log format
func (r *RestfulLoggingFilter) SetRequestLogFormat(format string) {
	r.requestLogFormat = format
}

// SetResponseLogFormat sets the verbose response log format
func (r *RestfulLoggingFilter) SetResponseLogFormat(format string) {
	r.responseLogFormat = format
}

// SetLogRoot toggles logging the root call (ping)
func (r *RestfulLoggingFilter) SetLogRoot(logRoot bool) {
	r.logRoot = logRoot
}

// Filter is a filter function that logs before and after processing the request
func (r *RestfulLoggingFilter) Filter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	start := time.Now()
	request.SetAttribute(RequestTimestampAttribute, start)

	traceID, err := uuid.NewRandom()
	if err != nil {
		request.SetAttribute(TraceIDAttribute, traceID)
	}
	if !r.logRoot && isRoot(request) {
		// don't log the root call, as it is typically used as a server ping and can clutter the log
		chain.ProcessFilter(request, response)
		return
	}
	if r.logger.IsEnabledFor(logging.DEBUG) {
		headerStr := flattenHeader(request.Request.Header)
		r.logger.Debugf(r.requestLogFormat, request.Request.Method, request.Request.URL, traceID, request.Request.ContentLength, request.Request.Form, headerStr)
	}

	chain.ProcessFilter(request, response)
	if response.StatusCode() >= 400 && r.logger.IsEnabledFor(logging.ERROR) {
		headerStr := flattenHeader(response.Header())
		r.logger.Errorf(r.responseLogFormat, request.Request.Method, request.Request.URL, traceID, response.StatusCode(), time.Now().Sub(start), headerStr)
	} else if r.logger.IsEnabledFor(logging.DEBUG) {
		headerStr := flattenHeader(response.Header())
		r.logger.Debugf(r.responseLogFormat, request.Request.Method, request.Request.URL, traceID, response.StatusCode(), time.Now().Sub(start), headerStr)
	}
}

func isRoot(request *restful.Request) bool {
	return request.Request.Method == "GET" && request.Request.URL != nil && request.Request.URL.String() == "/"
}

func flattenHeader(h http.Header) string {
	headerStr := ""
	for k, v := range h {
		val := ""
		if strings.EqualFold(k, authorizationHeader) {
			val = "[hidden]"
		} else {
			val = flatten(v)
		}
		if headerStr != "" {
			headerStr += ","
		}
		headerStr = headerStr + fmt.Sprintf("%v:%v", k, val)
	}
	return headerStr
}

func flatten(value []string) string {
	str := ""
	for _, v := range value {
		if str == "" {
			str = "["
		} else {
			str += ","
		}
		str += v
	}
	str += "]"
	return str
}
