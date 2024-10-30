# ez

`ez` is a minimalistic Go package for error handling, that makes errors a first-class citizen in your application domain.

It provides a clean, easy (pun intended) and consistent way to handle errors across different consumer roles: your application logic, end users and developers.

Based on Ben Johnson [Failure is your domain](https://www.gobeyond.dev/failure-is-your-domain/) awesome post.

## Why ez?

Go's error handling can be challenging - while errors are core to the language, there's no prescribed way to handle them effectively. `ez` solves this by providing:

- **Role-Based Error Handling**: Different error information for different consumers

  - 🤖 **Application**: Clean error codes for programmatic handling
  - 👤 **End Users**: Clear, actionable error messages
  - 👨💻 **Developers**: Detailed logical stack traces for debugging

- **Domain-Centric Design**: Errors become part of your domain model, just like your `Customer` or `Order` types
- **Clean Stack Traces**: Logical operation tracking without the noise of full stack traces
- **Standard Error Codes**: Pre-defined, widely-applicable error codes inspired by HTTP/gRPC standards

## Installation

```bash
go get github.com/vanclief/ez
```

## Quick Start

```go
import "github.com/vanclief/ez"

// Create a new error
err := ez.New(
    "UserService.CreateUser",  // Operation name
    ez.EINVALID,              // Error code
    "Username cannot be empty", // User-friendly message
    nil,                      // Optional underlying error
)

// Check error codes
if ez.ErrorCode(err) == ez.EINVALID {
    // Handle validation error
}

// Get user-friendly message
message := ez.ErrorMessage(err) // "Username cannot be empty"

// Get full error trace for developers
ez.ErrorStackTrace(err) // "UserService.CreateUser: <invalid> Username cannot be empty"
```

## Core Features

### 1. Standardized Error Codes

Pre-defined error codes that cover most common scenarios:

```go
const (
    ECONFLICT   = "conflict"           // Action cannot be performed
    EINTERNAL   = "internal"           // Internal error
    EINVALID    = "invalid"            // Validation failed
    ENOTFOUND   = "not_found"          // Entity does not exist
    ENOTAUTHORIZED    = "not_authorized"     // Missing permissions
    ENOTAUTHENTICATED = "not_authenticated"  // Not authenticated
    ERESOURCEEXHAUSTED = "resource_exhausted" // Resource exhausted
    ENOTIMPLEMENTED    = "not_implemented"    // Not implemented
    EUNAVAILABLE       = "unavailable"        // System unavailable
)
```

### 2. Error Wrapping

Build logical stack traces by wrapping errors:

```go
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    const op = "UserService.CreateUser"

    // Validate user
    if user.Username == "" {
        return ez.New(op, ez.EINVALID, "Username is required", nil)
    }

    // Try to create user
    if err := s.db.CreateUser(user); err != nil {
        return ez.Wrap(op, err) // Preserves original error details
    }

    return nil
}
```

### 3. Error Information Extraction

Easy access to error details:

```go
// Get error code
code := ez.ErrorCode(err)    // e.g., "invalid"

// Get user message
msg := ez.ErrorMessage(err)  // e.g., "Username is required"

// Get full error trace (for developers)
ez.ErrorStacktrace(err)         // e.g., "UserService.CreateUser: <invalid> Username is required"

```

## Example

Here's an example showing how to handle errors with `ez`:

```go
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    const op = "UserService.CreateUser"

    // Validation error (end user focused)
    if user.Username == "" {
        return ez.New(op, ez.EINVALID, "Username is required", nil)
    }

    // Check for conflicts (application logic focused)
    exists, err := s.checkUserExists(user.Username)
    if err != nil {
        return ez.Wrap(op, err) // Wraps internal error for developers
    }
    if exists {
        return ez.New(op, ez.ECONFLICT,
            "Username is already taken. Please choose another one.", nil)
    }

    // Database error (developer focused)
    if err := s.db.CreateUser(user); err != nil {
        return ez.Wrap(op, err)
    }

    return nil
}
```

### Handling the Error

```go
user := &User{Username: ""}
err := svc.CreateUser(ctx, user)

// Application logic
switch ez.ErrorCode(err) {
case ez.EINVALID:
    // Handle validation error
case ez.ECONFLICT:
    // Handle conflict error
case ez.EINTERNAL:
    // Handle internal error
}

// End user message
if err != nil {
    fmt.Println("Error:", ez.ErrorMessage(err))
    // Output: "Error: Username is required"
}

// Developer debugging
if err != nil {
    ez.ErrorStacktrace(err)
    // Output: "UserService.CreateUser: <invalid> Username is required"
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
