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
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenManager generates and validates authentication tokens
type TokenManager struct {
	signingKey      []byte
	tokenExpiration int
}

// NewTokenManager creates a new token manager with the given private key and token expiration (in seconds)
func NewTokenManager(signingKey []byte, tokenExpiration int) *TokenManager {
	generator := &TokenManager{
		signingKey:      signingKey,
		tokenExpiration: tokenExpiration,
	}
	return generator
}

// ValidateToken validates the token, returning an error if validation fails.
func (t *TokenManager) ValidateToken(token string) error {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.signingKey, nil
	}
	parsedToken, err := jwt.Parse(token, keyFunc)
	if err != nil || !parsedToken.Valid {
		return errors.New("invalid token")
	}
	return nil
}

// CreateToken creates, signs and returns a new JSON Web Token using the signing key and expiration provided.
func (t *TokenManager) CreateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Second * time.Duration(t.tokenExpiration)).Unix(),
	})
	tokenString, err := token.SignedString(t.signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
