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
	"time"

	"github.com/jmoiron/sqlx"
)

// ResourceResult represents a definition of a resource that is identifiable and tracks creation/update timestamps.
type ResourceResult struct {
	ID         int64
	CreatedOn  time.Time
	ModifiedOn time.Time
}

// GetCreatedID grabs the sequential id from the result rows returned from an INSERT...RETURNING query and closes the result set.
func GetCreatedID(rows *sqlx.Rows) (int64, error) {
	defer rows.Close()
	if rows.Next() {
		var id sql.NullInt64
		err := rows.Scan(&id)
		if err != nil {
			return -1, err
		}
		if id.Valid {
			return id.Int64, nil
		}
	}
	return -1, errors.New("no rows inserted or returned")
}

// GetCreatedResourceResult grabs the sequential id created_on and modified_on timestamps from the result rows returned from an INSERT...RETURNING query and closes the result set.
func GetCreatedResourceResult(rows *sqlx.Rows) (*ResourceResult, error) {
	defer rows.Close()
	if rows.Next() {
		var id sql.NullInt64
		var createdOn time.Time
		var modifiedOn time.Time
		err := rows.Scan(&id, &createdOn, &modifiedOn)
		if err != nil {
			return nil, err
		}
		if id.Valid {
			return &ResourceResult{
				ID:         id.Int64,
				CreatedOn:  createdOn,
				ModifiedOn: modifiedOn}, nil
		}
	}
	return nil, errors.New("no rows inserted or returned")
}
