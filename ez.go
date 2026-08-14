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
	e, ok := err.(*Error)
	if ok && e != nil {
		// Shallow-copy the data map so AddData on the wrapper doesn't
		// mutate the wrapped error's keys. Nested mutable values stay
		// shared.
		var data map[string]interface{}
		if e.Data != nil {
			data = make(map[string]interface{}, len(e.Data))
			for k, v := range e.Data {
				data[k] = v
			}
		}
		// Derive an empty code from deeper in the chain so a hand-built
		// code-less *Error doesn't propagate `{"code": ""}`. Message is
		// never given fallback text here: it holds only what someone
		// explicitly set, and ErrorMessage resolves fallbacks at read
		// time — baking them in would freeze contradictions like a
		// not_found error carrying internal-error text.
		code, message := e.Code, e.Message
		if code == "" {
			code = ErrorCode(err)
		}
		if message == "" {
			message = errorMessage(err)
		}
		return &Error{
			Op:      op,
			Code:    code,
			Message: message,
			Data:    data,
			Err:     err,
		}
	}
	return &Error{Op: op, Code: ErrorCode(err), Message: errorMessage(err), Err: err}
}

// Unwrap returns the nested error, so errors.Is and errors.As can
// traverse chains that mix *Error with standard library wrappers.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// Error returns the string representation of the error message.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
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
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf(`%s <%s> "%s"`, e.Op, e.Code, e.Message)
}
