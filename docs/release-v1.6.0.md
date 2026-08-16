# v1.6.0

- Added `ETIMEOUT` and `ECANCELED` application error codes.
- Added automatic detection of context cancellation, deadline errors and top-level `Timeout()` errors in `ErrorCode` and `Wrap`.
- Added `Unwrap()` to `*Error`, enabling standard `errors.Is` and `errors.As` traversal.
- Made typed-nil `*Error` values safe in error formatting, wrapping and stack traces.
- Updated `Wrap` to derive missing codes and messages from direct ez chains, resolve fallback messages at read time and shallow-copy error data.
- Added safe default end-user messages for timeouts and cancellations.
- Expanded HTTP mappings for timeouts, cancellations, unavailable services, missing resources and conflicts.
- Expanded gRPC mappings for cancellations, deadlines, conflicts and additional status codes.
- Prevented `NewFromGRPC` from copying untrusted gRPC status descriptions into the user-facing message while retaining diagnostics in the nested error.
- Improved handling of nil, wrapped, joined and mixed standard/ez error chains.
- Added regression coverage and v1.5/v1.6 upgrade documentation.
- Added repository-wide Go style and review rules for contributors.
