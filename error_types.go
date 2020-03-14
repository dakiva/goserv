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
	"fmt"
	"net/http"
)

// httpError represents an error that has an http status code.
type httpError interface {
	// StatusCode returns the HTTP status code appropriate for the error type
	StatusCode() int
}

// ResourceNotFoundError represents an error when a resource could not be found (404).
type ResourceNotFoundError struct {
	ResourceID       string
	ResourceTypeName string
	Err              error
}

// Error returns this error as a string
func (r *ResourceNotFoundError) Error() string {
	message := ""
	if r.Err != nil {
		message = r.Err.Error()
	}
	return fmt.Sprintf("%v resource %v not found: %v", r.ResourceTypeName, r.ResourceID, message)
}

// Unwrap returns the underlying error
func (r *ResourceNotFoundError) Unwrap() error { return r.Err }

// StatusCode returns the HTTP status code appropriate for the error type
func (r *ResourceNotFoundError) StatusCode() int { return http.StatusNotFound }

// DuplicateResourceError represents a duplicate resource signal (409) typically raised on resource creation
type DuplicateResourceError struct {
	ResourceTypeName string
	Err              error
}

// Error returns this error as a string
func (d *DuplicateResourceError) Error() string {
	message := ""
	if d.Err != nil {
		message = d.Err.Error()
	}
	return fmt.Sprintf("%v resource already exists: %v", d.ResourceTypeName, message)
}

// Unwrap returns the underlying error
func (d *DuplicateResourceError) Unwrap() error { return d.Err }

// StatusCode returns the HTTP status code appropriate for the error type
func (d *DuplicateResourceError) StatusCode() int { return http.StatusConflict }

// UnauthorizedError represents unauthorized access to the system
type UnauthorizedError struct {
	Login string
	Err   error
}

// Error returns this error as a string
func (u *UnauthorizedError) Error() string {
	message := ""
	if u.Err != nil {
		message = u.Err.Error()
	}
	return fmt.Sprintf("%v is unauthorized: %v", u.Login, message)
}

// Unwrap returns the underlying error
func (u *UnauthorizedError) Unwrap() error { return u.Err }

// StatusCode returns the HTTP status code appropriate for the error type
func (u *UnauthorizedError) StatusCode() int { return http.StatusUnauthorized }

// IllegalArgumentError represents a bad request argument (400)
type IllegalArgumentError struct {
	Argument string
	Err      error
}

// Error returns this error as as string
func (i *IllegalArgumentError) Error() string {
	message := ""
	if i.Err != nil {
		message = i.Err.Error()
	}
	return fmt.Sprintf("illegal argument error: %v :%v", i.Argument, message)
}

// Unwrap returns the underlying error
func (i *IllegalArgumentError) Unwrap() error { return i.Err }

// StatusCode returns the HTTP status code appropriate for the error type
func (i *IllegalArgumentError) StatusCode() int { return http.StatusBadRequest }

// AccessDeniedError represents unauthorized access to a specific system function or resource
type AccessDeniedError struct {
	Err error
}

// Error returns this error as a string
func (a *AccessDeniedError) Error() string {
	message := ""
	if a.Err != nil {
		message = a.Err.Error()
	}
	return fmt.Sprintf("access denied: %v", message)
}

// Unwrap returns the underlying error
func (a *AccessDeniedError) Unwrap() error { return a.Err }

// StatusCode returns the HTTP status code appropriate for the error type
func (a *AccessDeniedError) StatusCode() int { return http.StatusForbidden }
