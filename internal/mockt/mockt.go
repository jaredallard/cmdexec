// Copyright (C) 2026 cmdexec contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

// Package mockt implements a system for mocking [testing.T].
package mockt

import "testing"

// T is an interface for [testing.T]
type T interface {
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
