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

import "github.com/op/go-logging"

// RestfulLogAdapter adapts go-restful logging to go-logging
type RestfulLogAdapter struct {
	logger      *logging.Logger
	logEndpoint bool
}

// NewRestfulLogAdapter constructs and returns a new instance
func NewRestfulLogAdapter(logger *logging.Logger, logEndpoint bool) *RestfulLogAdapter {
	return &RestfulLogAdapter{
		logger:      logger,
		logEndpoint: logEndpoint,
	}
}

// Print adapts the print call to a debug log
func (r *RestfulLogAdapter) Print(v ...interface{}) {
	if r.logEndpoint {
		r.logger.Debug(v)
	}
}

// Printf adapts the printf call to a debugf log
func (r *RestfulLogAdapter) Printf(format string, v ...interface{}) {
	if r.logEndpoint {
		r.logger.Debugf(format, v...)
	}
}
