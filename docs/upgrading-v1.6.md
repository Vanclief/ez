# Upgrading to v1.6

v1.6.0 adds `ETIMEOUT` and `ECANCELED` codes, the `StatusClientClosedRequest`
(499) constant and an `Unwrap` method on `*Error` (so `errors.Is`/`errors.As`
traverse ez chains). The API changes are additive only, but classification
behavior is not:

- `ErrorCode`/`Wrap` now classify wrapped timeouts as `ETIMEOUT` and
  cancellations as `ECANCELED` instead of `EINTERNAL`. `NewFromGRPC` also
  recognizes raw or wrapped context cancellation and deadline errors.
- `HTTPStatusToError` now maps 502/503 to `EUNAVAILABLE`, unmapped 4xx to
  `EINVALID`, 410 to `ENOTFOUND`, 412 to `ECONFLICT`, 504/408 to `ETIMEOUT`
  and 499 to `ECANCELED`. HTTP 500 and other unmapped 5xx statuses remain
  `EINTERNAL` because they do not imply that retrying will help.
- `GRPCCodeToError` now maps inbound `Aborted`, `AlreadyExists` and `OutOfRange`
  to `ECONFLICT`, alongside `FailedPrecondition`. Outbound `ECONFLICT` remains
  `FailedPrecondition`.
- gRPC `Internal`, `Unknown` and `DataLoss` remain `EINTERNAL`; only
  `Unavailable` maps to `EUNAVAILABLE`. A valid explicit gRPC status takes
  precedence over context sentinels elsewhere in the error chain.
- `NewFromGRPC` no longer copies any upstream status description into the
  end-user-facing `Message`. Diagnostics remain in `Err`; callers must
  explicitly provide a trusted, sanitized message when needed.
- `Wrap` derives an empty `Code` from deeper in a direct ez chain instead of
  copying it verbatim, and no longer bakes fallback text into `Message`:
  the field holds only what was explicitly set, and `ErrorMessage` resolves
  fallbacks at read time.
- `*Error` now implements `Unwrap`, so existing `errors.Is`/`errors.As`
  calls traverse through ez errors into their nested causes where they
  previously stopped at the `*Error` boundary. Audit guards like
  `errors.Is(err, sql.ErrNoRows)` — they may start returning true.
- `ErrorMessage` returns new end-user strings for timeout and cancellation
  codes ("The operation timed out." / "The operation was
  canceled.") instead of the generic internal-error text.
- `Wrap` now shallow-copies an existing data map instead of sharing it, so
  `AddData` on the wrapper no longer mutates the original error. Nested mutable
  values remain shared; changes made directly to the original are not isolated.
