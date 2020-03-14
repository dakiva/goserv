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
	"fmt"
	"strings"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/op/go-logging"
)

const (
	// DefaultRequestLogFormat format
	DefaultRequestLogFormat = "[%s %s] request data\n\trem_ip='%s'\n\treceived_cookies='%s'\n\treferer='%s'\n\tuser_agent='%s'\n\tcontent_length='%d'\n\tcontent_type='%s'\n\tform='%s'\nrequest headers:\n%s\n"
	// DefaultResponseLogFormat format
	DefaultResponseLogFormat = "[%s %s] response\n\tstatus=%d\n\ttime=%s\nheaders:\n%s\n"

	timeRequestAttribute = "request_time_attr"
)

// RestfulLoggingFilter is middleware that logs restful requests and response payloads
type RestfulLoggingFilter struct {
	logger            *logging.Logger
	requestLogFormat  string
	responseLogFormat string
}

// NewRestfulLoggingFilter initializes a new logging filter instance
func NewRestfulLoggingFilter(logger *logging.Logger) *RestfulLoggingFilter {
	return &RestfulLoggingFilter{
		logger:            logger,
		requestLogFormat:  DefaultRequestLogFormat,
		responseLogFormat: DefaultResponseLogFormat,
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

// Filter is a filter function that logs before and after processing the request
func (r *RestfulLoggingFilter) Filter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	start := time.Now()
	request.SetAttribute(timeRequestAttribute, start)
	if isRoot(request) {
		// don't log the root call, as it is typically used as a server ping and can clutter the log
		chain.ProcessFilter(request, response)
		return
	}
	if r.logger.IsEnabledFor(logging.DEBUG) {
		headerStr := ""
		for k, v := range request.Request.Header {
			if strings.EqualFold(k, "access-token") {
				headerStr = headerStr + fmt.Sprintf("\t%v:\t[hidden]\n", k)
			} else {
				headerStr = headerStr + fmt.Sprintf("\t%v:\t%v\n", k, v)
			}
		}
		r.logger.Debugf(r.requestLogFormat, request.Request.Method, request.Request.URL, request.HeaderParameter("X-Forwarded-For"), request.Request.Cookies(), request.Request.Referer(), request.Request.UserAgent(), request.Request.ContentLength, request.HeaderParameter("Content-Type"), request.Request.Form, headerStr)
	}

	chain.ProcessFilter(request, response)
	if response.StatusCode() >= 400 {
		if r.logger.IsEnabledFor(logging.ERROR) {
			r.logger.Errorf(r.responseLogFormat, request.Request.Method, request.Request.URL, response.StatusCode(), time.Now().Sub(start))
		}
	} else if r.logger.IsEnabledFor(logging.DEBUG) {
		headerStr := ""
		for k, v := range response.Header() {
			headerStr = headerStr + fmt.Sprintf("\t%v:\t%v\n", k, v)
		}
		r.logger.Debugf(r.responseLogFormat, request.Request.Method, request.Request.URL, response.StatusCode(), time.Now().Sub(start), headerStr)
	}
}

func isRoot(request *restful.Request) bool {
	return request.Request.Method == "GET" && request.Request.URL != nil && request.Request.URL.String() == "/"
}
