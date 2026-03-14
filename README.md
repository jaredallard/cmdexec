# cmdexec

Go library for mocking `exec.Command` and `exec.CommandContext` calls.
Aims to provide an as close as possible [exec.Cmd]-like interface for
drop-in support.

## Differences

- Instead of setting `cmd.Stdin` (or `in/err`), use `SetStdin()`
  functions instead. This is due to those being fields on a struct,
  which cannot be on an interface.

## Usage

See our [Go docs](https://pkg.go.dev/go.rgst.io/jaredallard/cmdexec/v2).

### Asserting Output

Normally, you shouldn't need to assert anything as your function that
executes a command should give you testing signal (is it working or not
:wink:). However, you can assert certain fields with this library.
Currently, this is limited to Stdin checking.

If you set [MockCommand.Stdin] and call `SetStdin` in the function
executing a command, `Stdin` will be checked to ensure it is equal. This
is to allow greater testing if required.

## License

LGPL-3.0
