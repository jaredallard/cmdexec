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

package cmdexec_test

import (
	"fmt"
	"runtime"

	"github.com/jaredallard/cmdexec"
)

func ExampleMockCommand_SetDir() {
	// Specific example doesn't work on Windows, but the functionality
	// does!
	if runtime.GOOS == "windows" {
		return
	}

	cmd := cmdexec.Command("pwd")
	cmd.SetDir("/tmp")

	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))

	// Output:
	// /tmp
}

func ExampleMockCommand_UseOSStreams() {
	// This example demonstrates how to use the UseOSStreams function to
	// set the stdin, stdout, and stderr of a command to the OS streams.
	cmd := cmdexec.Command("echo", "hello")
	cmd.UseOSStreams(false)
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	// Output:
	// hello
}
