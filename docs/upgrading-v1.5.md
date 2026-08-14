# Upgrading to v1.5

Since v1.5.0 the constructors no longer take an operation argument — it is
derived automatically from the calling function. Drop the first argument at
every call site:

```go
ez.New(op, ez.EINVALID, "Username is required", nil)  // before
ez.New(ez.EINVALID, "Username is required", nil)      // after
```

The same applies to `ez.Root` and `ez.Wrap`. If you need a custom operation
name (for example in a gRPC interceptor, where the method name is better than
any function name), set the exported field directly:

```go
e := ez.NewFromGRPC(err)
e.Op = method
```

`ez.ErrorStacktrace` now returns the stacktrace as a string instead of
printing it to stdout, so it can be routed to your logger:

```go
slog.Error("create user failed", "stacktrace", ez.ErrorStacktrace(err))
```
