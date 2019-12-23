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
	"errors"
	"fmt"
	"log/syslog"
	"os"
	"strings"

	"github.com/op/go-logging"
)

// LoggingConfig contains configuration for op/go-logging
type LoggingConfig struct {
	LogLevel string          `json:"log_level"`
	Format   string          `json:"format"`
	Backends []backendConfig `json:"backends"`
}

// Validate ensures a configuration has populated all required fields.
func (l *LoggingConfig) Validate() error {
	if _, err := logging.LogLevel(l.LogLevel); err != nil {
		return err
	}
	if l.Format != "" {
		if _, err := logging.NewStringFormatter(l.Format); err != nil {
			return err
		}
	}
	if len(l.Backends) == 0 {
		return errors.New("no logging backends defined")
	}
	for _, backend := range l.Backends {
		if backend.BackendName != "STDOUT" &&
			backend.BackendName != "SYSLOG" &&
			backend.BackendName != "FILE" {
			return fmt.Errorf("invalid backend name %s", backend.BackendName)
		} else if backend.BackendName == "FILE" && backend.FilePath == "" {
			return fmt.Errorf("file backend requires a file path to be set")
		}
	}
	return nil
}

// InitializeLogging configures logging based on the logging configuration.
func (l *LoggingConfig) InitializeLogging() error {
	if l.Format != "" {
		format := logging.MustStringFormatter(l.Format)
		logging.SetFormatter(format)
	}
	level, err := logging.LogLevel(l.LogLevel)
	if err != nil {
		return err
	}

	backends := []logging.Backend{}
	for _, b := range l.Backends {
		backend, err := b.getBackend(level)
		if err != nil {
			return err
		}
		backends = append(backends, backend)
	}
	logging.SetBackend(backends...)
	logging.SetLevel(level, "")
	return nil
}

type backendConfig struct {
	BackendName string `json:"backend_name"`
	FilePath    string `json:"file_path"`
}

// Returns a suitable logging backend for the backend name or an error if a backend name does not describe a logging backend.
func (b *backendConfig) getBackend(level logging.Level) (logging.Backend, error) {
	switch {
	case strings.EqualFold(b.BackendName, "STDOUT"):
		return logging.NewLogBackend(os.Stdout, "", 0), nil
	case strings.EqualFold(b.BackendName, "FILE"):
		if b.FilePath != "" {
			file, err := os.OpenFile(b.FilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				return nil, err
			}
			return logging.NewLogBackend(file, "", 0), nil
		}
		return nil, errors.New("Error using file as a backend")
	case strings.EqualFold(b.BackendName, "SYSLOG"):
		b, _ := logging.NewSyslogBackendPriority("", toSyslogPriority(level))
		return b, nil
	}
	return nil, errors.New("Error creating the backend for logging")
}

func toSyslogPriority(level logging.Level) syslog.Priority {
	switch level {
	case logging.CRITICAL:
		return syslog.LOG_CRIT
	case logging.ERROR:
		return syslog.LOG_ERR
	case logging.WARNING:
		return syslog.LOG_WARNING
	case logging.NOTICE:
		return syslog.LOG_NOTICE
	case logging.INFO:
		return syslog.LOG_INFO
	case logging.DEBUG:
		return syslog.LOG_DEBUG
	default:
		return syslog.LOG_DEBUG
	}
}
