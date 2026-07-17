package ez

import (
	"bytes"
	"fmt"
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
	// Data about the error
	Data map[string]interface{} `json:"data,omitempty"`
}

// New creates and returns a new error. The operation is derived from the
// calling function: "pkg.Type.Method" or "pkg.Function".
func New(code, message string, err error) *Error {
	return &Error{Op: callerOp(), Code: code, Message: message, Err: err}
}

// Root creates a new root error. The operation is derived from the
// calling function.
func Root(code, message string) *Error {
	return &Error{Op: callerOp(), Code: code, Message: message}
}

// Wrap returns a new error that contains the passed error, useful for
// creating stacktraces. The operation is derived from the calling function.
func Wrap(err error) *Error {
	op := callerOp()
	if e, ok := err.(*Error); ok {
		return &Error{
			Op:      op,
			Code:    e.Code,
			Message: e.Message,
			Data:    e.Data,
			Err:     err,
		}
	}
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
