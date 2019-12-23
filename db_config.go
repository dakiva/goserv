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
	"strings"

	"github.com/dakiva/dbx"
	"github.com/jmoiron/sqlx"
)

// DBConfig represents database configuration that points to a specific schema and allows for connection specific settings.
type DBConfig struct {
	Hostname           string `json:"hostname"`
	Port               int    `json:"port"`
	MaxIdleConnections int    `json:"max_idle_connections"`
	MaxOpenConnections int    `json:"max_open_connections"`
	DBName             string `json:"dbname"`
	SSLMode            string `json:"sslmode"`
	ConnectTimeout     int    `json:"connect_timeout"`
	Role               string `json:"role"`
	SchemaName         string `json:"schema_name"`
	RolePassword       string `json:"role_password"`
}

// ToDsn converts this configuration to a standard DSN that can be used to open a connection to a specific schema.
func (d *DBConfig) ToDsn() string {
	dsn := ""
	if d.Hostname != "" {
		dsn += "host=" + d.Hostname + " "
	}
	if d.Port > 0 {
		dsn += fmt.Sprintf("port=%d ", d.Port)
	}
	dsn += "user=" + d.Role + " "
	if d.RolePassword != "" {
		dsn += "password=" + d.RolePassword + " "
	}
	dsn += "dbname=" + d.DBName + " "
	dsn += "sslmode=" + d.SSLMode + " "
	if d.ConnectTimeout > 0 {
		dsn += fmt.Sprintf("connect_timeout=%d", d.ConnectTimeout)
	}
	return strings.TrimSpace(dsn)
}

// Validate ensures a configuration has populated all required fields.
func (d *DBConfig) Validate() error {
	if d.DBName == "" {
		return errors.New("Database name must be specified")
	}
	if d.SSLMode != "disable" &&
		d.SSLMode != "require" &&
		d.SSLMode != "verify-ca" &&
		d.SSLMode != "verify-full" {
		return errors.New("Database SSLMode must be specified")
	}
	if d.Role == "" {
		return errors.New("Role must be specified")
	}
	if d.SchemaName == "" {
		return errors.New("Schema name must be specified")
	}
	if d.RolePassword == "" {
		return errors.New("Schema role password must be specified")
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
