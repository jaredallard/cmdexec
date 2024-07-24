# cmdexec

<a href="https://github.com/jaredallard/cmdexec/releases">
	<img alt="Latest Version" src="https://img.shields.io/github/v/release/jaredallard/cmdexec?style=for-the-badge">
</a>
<a href="https://github.com/jaredallard/cmdexec/blob/main/LICENSE">
	<img alt="License" src="https://img.shields.io/github/license/jaredallard/cmdexec?style=for-the-badge">
</a>
<a href="https://github.com/jaredallard/cmdexec/actions/workflows/tests.yaml">
	<img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/jaredallard/cmdexec/tests.yaml?style=for-the-badge">
</a>
<a href="https://app.codecov.io/gh/jaredallard/cmdexec">
	<img alt="Codecov" src="https://img.shields.io/codecov/c/github/jaredallard/cmdexec?style=for-the-badge">
</a>

<br />

Go library for mocking `exec.Command` and `exec.CommandContext` calls.
Aims to provide an as close as possible [exec.Cmd]-like interface for
drop-in support.

## Differences

- Instead of setting `cmd.Stdin` (or `in/err`), use `SetStdin()`
functions instead. This is due to those being fields on a struct,
which cannot be on an interface.

## Usage

See our [Go docs](https://pkg.go.dev/github.com/jaredallard/cmdexec).

### Asserting Output

Normally, you shouldn't need to assert anything as your function that
executes a command should give you testing signal (is it working or not
:wink:). However, you can assert certain fields with this library.
Currently, this is limited to Stdin checking.

If you set [MockCommand.Stdin] and call `SetStdin` in the function
executing a command, `Stdin` will be checked to ensure it is equal. This
is to allow greater testing if required.

## License

GPL-3.0
