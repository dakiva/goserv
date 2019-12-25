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
	"time"
)

// StatsResource represents health stats that can be exposed as an API, useful for monitoring
type StatsResource struct {
	ServiceName      string    `json:"service_name" description:"The name of the service."`
	ServiceUptime    string    `json:"service_uptime" description:"The service uptime in unix format."`
	ServiceStartTime time.Time `json:"service_start_time" description:"Timestamp in UTC of when the service was first started."`
	ServiceVersion   string    `json:"service_version" description:"The version of the service binary."`
	CommitHash       string    `json:"commit_hash" description:"The git commit hash representing the sources built into the service binary."`
	APIVersion       string    `json:"api_version" description:"The version of the API (protocol)."`
}

// UpdateUptime updates the service uptime based on the current time returned in a Unix standard uptime format.
func (s *StatsResource) UpdateUptime() {
	s.ServiceUptime = formatUptime(s.ServiceStartTime)
}

// Returns the time elapsed since the start time. The returned value is in a format similar
// to the value returned by the Unix uptime command.
func formatUptime(startTime time.Time) string {
	uptime := time.Since(startTime)
	hours := int64(uptime.Hours())
	remainingMinutes := int64(uptime.Minutes()) % 60
	var days int64 = hours / 24
	remainingHours := hours % 24
	var s string
	if days > 0 {
		if days == 1 {
			s = "1 day, "
		} else if days > 1 {
			s = fmt.Sprintf("%d days, ", days)
		}
		s += fmt.Sprintf("%02d:%02d", remainingHours, remainingMinutes)
		return s
	}
	if hours > 1 {
		return fmt.Sprintf("%02d:%02d", hours, remainingMinutes)
	}
	return uptime.String()
}
