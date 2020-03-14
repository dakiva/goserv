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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndValidateToken(t *testing.T) {
	// given
	key := []byte("f2ca1bb6c7e907d06dafe4687e579fce76b37e4e93b7605022da52e6ccc26fd2")
	manager := NewTokenManager(key, 1)

	// when
	token, err := manager.CreateToken()

	// then
	assert.NoError(t, err)
	err = manager.ValidateToken(token)
	assert.NoError(t, err)
}

func TestDifferentKeys(t *testing.T) {
	// given
	key := []byte("f2ca1bb6c7e907d06dafe4687e579fce76b37e4e93b7605022da52e6ccc26fd2")
	manager := NewTokenManager(key, 1)

	token, err := manager.CreateToken()
	assert.NoError(t, err)
	manager2 := NewTokenManager([]byte("e2ca1bb6c7e907d06dafe4687e579fce76b37e4e93b7605022da52e6ccc26fe2"), 1)

	// when
	err = manager2.ValidateToken(token)

	// then
	assert.Error(t, err)
}

func TestExpiredToken(t *testing.T) {
	// given
	key := []byte("f2ca1bb6c7e907d06dafe4687e579fce76b37e4e93b7605022da52e6ccc26fd2")
	// set the time back one second so that the expiration timestamp is in the past
	manager := NewTokenManager(key, -1)

	token, err := manager.CreateToken()
	assert.NoError(t, err)

	// when
	err = manager.ValidateToken(token)

	// then
	assert.Error(t, err)
}
