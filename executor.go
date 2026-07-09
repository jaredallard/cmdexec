// Copyright (C) 2026 cmdexec contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package cmdexec

import (
	"context"
	"os/exec"
	"sync"
)

// Contains package globals to control which executor is used by the
// package as well as locks to ensure this package is thread-safe.
var (
	// executor is the function used to create new commands. By default,
	// this is set to [stdExecutor], but can be replaced with a mock
	// executor using [UseMockExecutor].
	executor executorFn = stdExecutor

	// lookPath is the function used to look up the path of a command.
	// By default, this is set to [exec.LookPath], but can be replaced
	// with a mock using [UseMockExecutor].
	lookPath lookPathFn = exec.LookPath

	// Locks to control the accessing of the executor variable. We don't
	// use a [sync.RWMutex] here because we want to be able to lock the
	// read and write operations separately.
	executorRLock = new(sync.Mutex)
	executorWLock = new(sync.Mutex)
)

// executorFn is a function that returns a new Cmd based on the given
// arguments.
type executorFn func(context.Context, string, ...string) Cmd

// lookPathFn is a function that searches for an executable named file
// in the directories named by the PATH environment variable.
type lookPathFn func(string) (string, error)
