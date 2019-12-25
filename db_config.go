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
	"os"
	"strconv"
	"strings"

	"github.com/dakiva/dbx"
	"github.com/jmoiron/sqlx"
)

const (
	// DefaultMaxOpenConnections represents the default number of open database connections
	DefaultMaxOpenConnections = 10
	// DefaultMaxIdleConnections represents the default number of idle database connections
	DefaultMaxIdleConnections = 10
)

// ParseDBConfig represents a configuration that is parsed from an environment variable. Returns
// an error onky if an environment variable exists and parsing data from that environment fails. If an environment variable is not set nil is returned.
func ParseDBConfig(dsnEnv string) (*DBConfig, error) {
	if dsn, exists := os.LookupEnv(dsnEnv); exists {
		m := dbx.ParseDsn(dsn)
		port := 0
		if p, err := strconv.Atoi(m["port"]); err == nil && p > 0 {
			port = p
		} else {
			return nil, fmt.Errorf("invalid database port %s: %w", m["port"], err)
		}
		connectTimeout := 0
		if timeout, err := strconv.Atoi(m["connect_timeout"]); err == nil {
			connectTimeout = timeout
		} else {
			return nil, fmt.Errorf("invalid database connect_timeout %s: %w", m["connect_timeout"], err)
		}
		return &DBConfig{
			Hostname:           m["host"],
			Port:               port,
			DBName:             m["dbname"],
			SSLMode:            m["sslmode"],
			ConnectTimeout:     connectTimeout,
			SchemaName:         m["user"],
			Role:               m["user"],
			RolePassword:       m["password"],
			MaxIdleConnections: DefaultMaxOpenConnections,
			MaxOpenConnections: DefaultMaxIdleConnections,
		}, nil
	}
	return nil, nil
}

// DBConfig represents database configuration that points to a specific schema and allows for connection specific settings.
type DBConfig struct {
	Hostname           string `json:"hostname"`
	Port               int    `json:"port"`
	MaxIdleConnections int    `json:"max_idle_connections"`
	MaxOpenConnections int    `json:"max_open_connections"`
	DBName             string `json:"dbname"`
	SSLMode            string `json:"sslmode"`
	ConnectTimeout     int    `json:"connect_timeout"`
	SchemaName         string `json:"schema_name"`
	Role               string `json:"role"`
	RolePassword       string `json:"role_password"`
}

// ToDsn converts this configuration to a standard DSN that can be used to open a connection to a specific schema.
func (d *DBConfig) ToDsn() string {
	dsn := ""
	if d.Hostname != "" {
		dsn += "host=" + d.Hostname + " "
	}
	dsn += fmt.Sprintf("port=%d ", d.Port)
	dsn += "user=" + d.Role + " "
	dsn += "password=" + d.RolePassword + " "
	dsn += "dbname=" + d.DBName + " "
	dsn += "sslmode=" + d.SSLMode + " "
	if d.ConnectTimeout > 0 {
		dsn += fmt.Sprintf("connect_timeout=%d", d.ConnectTimeout)
	}
	return strings.TrimSpace(dsn)
}

// Validate ensures a configuration has populated all required fields.
func (d *DBConfig) Validate() error {
	if d.Port <= 0 {
		return errors.New("database port value must be a positive number")
	}
	if d.DBName == "" {
		return errors.New("database name must be specified")
	}
	if d.SSLMode != "disable" &&
		d.SSLMode != "require" &&
		d.SSLMode != "verify-ca" &&
		d.SSLMode != "verify-full" {
		return errors.New("database sslmode must be specified with [disable, require, verify-ca, verify-full]")
	}
	if d.Role == "" {
		return errors.New("role must be specified")
	}
	if d.SchemaName == "" {
		return errors.New("schema name must be specified")
	}
	if d.RolePassword == "" {
		return errors.New("role password must be specified")
	}
	if d.MaxIdleConnections > d.MaxOpenConnections {
		return errors.New("max idle connections cannot exceed the max number of open connections")
	}
	return nil
}

// OpenDB creates a pool of open connections to the database
func (d *DBConfig) OpenDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect(dbx.PostgresType, d.ToDsn())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(d.MaxOpenConnections)
	db.SetMaxIdleConns(d.MaxIdleConnections)
	return db, nil
}

// Empty returns true if this DBConfig represents an empty configuration
func (d *DBConfig) Empty() bool {
	return d.Hostname == "" &&
		d.Port == 0 &&
		d.MaxIdleConnections == 0 &&
		d.MaxOpenConnections == 0 &&
		d.DBName == "" &&
		d.SSLMode == "" &&
		d.ConnectTimeout == 0 &&
		d.SchemaName == "" &&
		d.Role == "" &&
		d.RolePassword == ""
}
