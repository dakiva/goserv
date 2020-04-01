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
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/dakiva/dbx"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/op/go-logging"
)

const (
	pqUniqueViolationCode = "23505"
)

// NewDBContextProviderSQLXWrapper initializes and returns a wrapped DBContextProvider instance
func NewDBContextProviderSQLXWrapper(db *sqlx.DB, logDB bool, logger *logging.Logger) dbx.DBContextProvider {
	return &DBContextProviderSQLXWrapper{db: db, logDB: logDB, logger: logger}
}

// DBContextProviderSQLXWrapper wraps a sqlx DB as a content provider
type DBContextProviderSQLXWrapper struct {
	db     *sqlx.DB
	logDB  bool
	logger *logging.Logger
}

// GetTxContext returns a transaction context, or an error
func (d *DBContextProviderSQLXWrapper) GetTxContext() (dbx.DBTxContext, error) {
	tx, err := d.db.Beginx()
	if err != nil {
		return nil, err
	}
	return &loggableDBTxContext{tx: tx, logDB: d.logDB, logger: d.logger}, nil
}

// GetContext returns a database context
func (d *DBContextProviderSQLXWrapper) GetContext() (dbx.DBContext, error) {
	return &loggableDBContext{ctx: d.db, logDB: d.logDB, logger: d.logger}, nil
}

type loggableDBContext struct {
	ctx    dbx.DBContext
	logDB  bool
	logger *logging.Logger
}

func (l *loggableDBContext) NamedExec(query string, arg interface{}) (sql.Result, error) {
	start := time.Now()
	res, err := l.ctx.NamedExec(query, arg)
	if l.logDB {
		l.logger.Debugf("[named exec time=%s]:%s [arg]:%v", time.Now().Sub(start), query, arg)
	}
	return res, interpretDBError(err, l.logDB, l.logger)
}

func (l *loggableDBContext) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	start := time.Now()
	res, err := l.ctx.NamedQuery(query, arg)
	if l.logDB {
		l.logger.Debugf("[named query time=%s]:%s [arg]:%v", time.Now().Sub(start), query, arg)
	}
	return res, interpretDBError(err, l.logDB, l.logger)
}

func (l *loggableDBContext) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	start := time.Now()
	res, err := l.ctx.PrepareNamed(query)
	if l.logDB {
		l.logger.Debugf("[preparing time=%s]:%s", time.Now().Sub(start), query)
	}
	return res, interpretDBError(err, l.logDB, l.logger)
}

type loggableDBTxContext struct {
	tx     dbx.DBTxContext
	logDB  bool
	logger *logging.Logger
}

func (l *loggableDBTxContext) Commit() error {
	start := time.Now()
	err := l.tx.Commit()
	if l.logDB {
		l.logger.Debugf("[tx commit time=%s]", time.Now().Sub(start))
	}
	return err
}

func (l *loggableDBTxContext) Rollback() error {
	start := time.Now()
	err := l.tx.Rollback()
	if l.logDB {
		l.logger.Debugf("[tx rollback time=%s]", time.Now().Sub(start))
	}
	return err
}

func (l *loggableDBTxContext) NamedExec(query string, arg interface{}) (sql.Result, error) {
	start := time.Now()
	res, err := l.tx.NamedExec(query, arg)
	if l.logDB {
		l.logger.Debugf("[tx named exec time=%s]:%s [arg]:%v", time.Now().Sub(start), query, arg)
	}
	return res, interpretDBError(err, l.logDB, l.logger)
}

func (l *loggableDBTxContext) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	start := time.Now()
	res, err := l.tx.NamedQuery(query, arg)
	if l.logDB {
		l.logger.Debugf("[tx named query time=%s]:%s [arg]:%v", time.Now().Sub(start), query, arg)
	}
	return res, interpretDBError(err, l.logDB, l.logger)
}

func (l *loggableDBTxContext) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	start := time.Now()
	res, err := l.tx.PrepareNamed(query)
	if l.logDB {
		l.logger.Debugf("[tx preparing time=%s]:%s", time.Now().Sub(start), query)
	}
	return res, interpretDBError(err, l.logDB, l.logger)
}

func interpretDBError(err error, logDB bool, logger *logging.Logger) error {
	if err != nil {
		if logDB {
			logger.Debugf("[db error]:%v", err)
		}
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == pqUniqueViolationCode {
			// use the table name as the resource type name
			// create a new error that removes all the postgres details from the stack
			resource := strings.ReplaceAll(pqErr.Table, "_", " ")
			return &DuplicateResourceError{
				ResourceTypeName: resource,
				Err:              errors.New(pqErr.Detail),
			}
		}
	}
	return err
}
