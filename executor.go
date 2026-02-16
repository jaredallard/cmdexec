// Copyright (C) 2026 cmdexec contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this program. If not, see
// <https://www.gnu.org/licenses/>.
//
// SPDX-License-Identifier: LGPL-3.0

package cmdexec

import (
	"context"
	"sync"
)

// Contains package globals to control which executor is used by the
// package as well as locks to ensure this package is thread-safe.
var (
	// executor is the function used to create new commands. By default,
	// this is set to [stdExecutor], but can be replaced with a mock
	// executor using [UseMockExecutor].
	executor executorFn = stdExecutor

	// Locks to control the accessing of the executor variable. We don't
	// use a [sync.RWMutex] here because we want to be able to lock the
	// read and write operations separately.
	executorRLock = new(sync.Mutex)
	executorWLock = new(sync.Mutex)
)

// executorFn is a function that returns a new Cmd based on the given
// arguments.
type executorFn func(context.Context, string, ...string) Cmd
