package ez

import (
	"bytes"
	"fmt"
)

// Application error codes
const (
	ECONFLICT          = "conflict"           // action cannot be performed
	EINTERNAL          = "internal"           // internal error
	EINVALID           = "invalid"            // validation failed
	ENOTFOUND          = "not_found"          // entity does not exist
	ENOTAUTHORIZED     = "not_authorized"     // requester does not have permissions to perform action
	ENOTAUTHENTICATED  = "not_authenticated"  // requester is not authenticated
	ERESOURCEEXHAUSTED = "resource_exhausted" // the resource has been exhausted
	ENOTIMPLEMENTED    = "not_implemented"    // the operation has not been implemented
	EUNAVAILABLE       = "unavailable"        // the system or operation is not available
)

// Error defines a standar application error
type Error struct {
	// Machine readable code
	Code string `json:"code"`
	// Human readable message
	Message string `json:"message"`
	// Logical operation
	Op string `json:"op"`
	// Nested error
	Err error `json:"err"`
}

// New creates and returns a new error
func New(op, code, message string, err error) *Error {
	return &Error{Op: op, Code: code, Message: message, Err: err}
}

// Wrap returns a new error that contains the passed error but with a different operation, useful for creating stacktraces
func Wrap(op string, err error) *Error {
	return &Error{Op: op, Code: ErrorCode(err), Message: ErrorMessage(err), Err: err}
}

// Error returns the string representation of the error message.
func (e *Error) Error() string {
	var buf bytes.Buffer

	// Print the current operation in our stack, if any.
	if e.Op != "" {
		fmt.Fprintf(&buf, "%s: ", e.Op)
	}

	// If wrapping an error, print its Error() message.
	// Otherwise print the error code & message.
	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Code != "" {
			fmt.Fprintf(&buf, "<%s> ", e.Code)
		}
		buf.WriteString(e.Message)
	}
	return buf.String()
}

// String returns a simplified string representation of the error message
func (e *Error) String() string {
	return fmt.Sprintf(`%s <%s> "%s"`, e.Op, e.Code, e.Message)
}

// ErrorCode returns the code of the root error, if available.
// Otherwise returns EINTERNAL.
func ErrorCode(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Code != "" {
		return e.Code
	} else if ok && e.Err != nil {
		return ErrorCode(e.Err)
	}
	return EINTERNAL
}

// ErrorMessage returns the human-readable message of the error, if available.
// Otherwise returns a generic error message.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Message != "" {
		return e.Message
	} else if ok && e.Err != nil {
		return ErrorMessage(e.Err)
	}
	return "An internal error has occurred. Please contact technical support."
}

// ErrorStacktrace prints a human-redable stacktrace of all nested errors.
func ErrorStacktrace(err error) {
	if err == nil {
		return
	} else if e, ok := err.(*Error); ok {
		fmt.Println(e.String())
		ErrorStacktrace(e.Err)
	} else if ok && e.Err != nil {
		fmt.Println(e.String())
	} else {
		fmt.Println(err.Error())
	}
}
