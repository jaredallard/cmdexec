// Copyright (C) 2024 Jared Allard <jaredallard@users.noreply.github.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

// mockt implements a system for mocking [testing.T].
package mockt

import "testing"

type T interface {
	// Fail marks the test as failed, see [testing.T.Fail].
	Fail()

	// Failed returns if the test has failed or not, see
	// [testing.T.Failed].
	Failed() bool

	// Fatal is a wrapper around [testing.T.Fatal].
	Fatal(args ...interface{})

	// Cleanup is a wrapper around [testing.T.Cleanup].
	Cleanup(func())
}

type t struct {
	_ *testing.T

	// failed denotes if the test failed or not.
	failed bool

	// args are the failure arguments for [t.Fatal].
	args []any

	cleanup func()
}

// New creates a new [T] that does not actually run any tests or exist
// as part of a test.
//
//nolint:revive // Why: We're an internal package.
func New() *t {
	return &t{}
}

// Fail implements [T.Fail].
func (t *t) Fail() {
	t.failed = true
}

// Failed implements [T.Failed].
func (t *t) Failed() bool {
	return t.failed
}

// Fatal implements [T.Fatal].
func (t *t) Fatal(args ...any) {
	t.failed = true
	t.args = args
}

// Cleanup implements [T.Cleanup].
func (t *t) Cleanup(fn func()) { t.cleanup = fn }

// RunCleanup runs the last set cleanup command. This is only provided
// in the mock implementation as a means to control when the cleanup
// function is ran.
func (t *t) RunCleanup() {
	t.cleanup()
}
