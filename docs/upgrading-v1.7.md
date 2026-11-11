# Upgrading to v1.7

`ErrorCode`, `ErrorMessage` and `ErrorData` now find an ez error behind a
single-cause wrapper that is not an ez error, such as
`fmt.Errorf("...: %w", err)`. Before, they read direct ez chains only, so an
ez error behind a foreign wrapper was not seen, and the error was usually
classified `EINTERNAL` with the generic message.

```go
err := fmt.Errorf("load profile: %w", ez.New(ez.ENOTFOUND, "User not found", nil))

ez.ErrorCode(err)         // before: "internal"   after: "not_found"
ez.ErrorMessage(err)      // before: generic text  after: "User not found"
ez.ErrorToHTTPStatus(err) // before: 500           after: 404
```

There are no API changes, but classification behavior changes:

- An error that reached a handler as `EINTERNAL` only because a foreign
  wrapper sat in the chain now reports the code, message and data of the
  nearest ez error. This covers your own `fmt.Errorf` calls and libraries
  that wrap the errors your callbacks return. Audit code that depends on such
  an error being `EINTERNAL`.
- `ErrorToHTTPStatus` and `ErrorToGRPCCode` follow `ErrorCode`, so the
  transport status changes with it.
- `NewFromGRPC` falls back to `ErrorCode` when the error carries no gRPC
  status, so its result changes for an ez error behind a foreign wrapper too.
  It keeps `ETIMEOUT` and `ECANCELED` and maps every other code to
  `EINTERNAL`, as it already did for a direct ez error.
- `Wrap` on such a wrapper now takes the code and message found behind it
  instead of `EINTERNAL` and an empty message.
- An explicit ez code wins over a cancellation or timeout deeper in the
  chain whether or not a foreign wrapper surrounds it. Before, the wrapper
  made the cancellation win.
- An explicit ez code also wins over a foreign wrapper above it that is
  itself a timeout or cancellation, such as a custom type whose `Timeout()`
  returns true. Before, such an error was classified `ETIMEOUT` or
  `ECANCELED`.
- The end-user message still comes only from an ez error. The text of a
  foreign wrapper is never returned by `ErrorMessage`.
- An error with several causes (`errors.Join`, or `fmt.Errorf` with more than
  one `%w`) is not searched, as before: it stays `EINTERNAL` unless it is a
  timeout or cancellation.
- `ErrorStacktrace` is unchanged: it prints structured ez lines until the
  first foreign error and then that error's complete text.

`ez.Wrap` is still the way to propagate an error, because it records the
operation as a structured line in the stack trace.
