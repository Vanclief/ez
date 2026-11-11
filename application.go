package ez

import (
	"context"
	"errors"
	"os"
)

// Application error codes
const (
	ECONFLICT          = "conflict"           // request is valid but the current state of the data forbids it
	EINTERNAL          = "internal"           // unexpected internal failure where retry is not implied
	EINVALID           = "invalid"            // validation failed
	ENOTFOUND          = "not_found"          // entity does not exist
	ENOTAUTHORIZED     = "not_authorized"     // requester does not have permissions to perform action
	ENOTAUTHENTICATED  = "not_authenticated"  // requester is not authenticated
	ERESOURCEEXHAUSTED = "resource_exhausted" // rate limit / quota exhausted
	ENOTIMPLEMENTED    = "not_implemented"    // the operation has not been implemented
	EUNAVAILABLE       = "unavailable"        // the operation is unavailable and retry may help
	ETIMEOUT           = "timeout"            // the operation ran out of time
	ECANCELED          = "canceled"           // the caller gave up before the operation finished
)

// ErrorCode returns the code of the nearest ez error in the chain, if
// available. Single-cause wrappers such as fmt.Errorf("%w") are followed, so
// a wrapped ez error keeps its code. Timeouts and cancellations are detected
// through standard error wrappers so they don't get misreported as internal
// errors. Otherwise returns EINTERNAL.
func ErrorCode(err error) string {
	if err == nil {
		return ""
	}
	e := nearest(err)
	if e != nil {
		if e.Code != "" {
			return e.Code
		}
		if e.Err != nil {
			return ErrorCode(e.Err)
		}
	}
	if errors.Is(err, context.Canceled) {
		return ECANCELED
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, os.ErrDeadlineExceeded) {
		return ETIMEOUT
	}
	// Shallow top-level check only (e.g. url.Error from a TLS handshake
	// timeout). Deliberately no tree traversal: hunting Timeout()
	// implementations through wrappers and joins caused masking and
	// precedence bugs, and the sentinels above already cover the stdlib
	// timeout sources.
	t, ok := err.(interface{ Timeout() bool })
	if ok && t.Timeout() {
		return ETIMEOUT
	}
	return EINTERNAL
}

// WithData adds a single key-value pair to the error's data
func (e *Error) AddData(key string, value interface{}) *Error {
	if e.Data == nil {
		e.Data = make(map[string]interface{})
	}
	e.Data[key] = value
	return e
}

// WithDataMap adds multiple key-value pairs to the error's data
func (e *Error) AddDataMap(data map[string]interface{}) *Error {
	if e.Data == nil {
		e.Data = make(map[string]interface{})
	}
	for k, v := range data {
		e.Data[k] = v
	}
	return e
}

// ErrorData returns the data of the nearest ez error in the chain that has
// any, following single-cause wrappers the same way ErrorCode does.
func ErrorData(err error) map[string]interface{} {
	if err == nil {
		return nil
	}
	e := nearest(err)
	if e != nil {
		if e.Data != nil {
			return e.Data
		}
		if e.Err != nil {
			return ErrorData(e.Err)
		}
	}
	return nil
}

// nearest returns the first ez error reached from err by following
// single-cause wrappers, or nil. It deliberately stops at an error with
// several causes (errors.Join, fmt.Errorf with more than one %w): picking one
// of them by argument order is the masking and precedence problem ErrorCode
// avoids for timeouts, so such an error has no ez classification of its own.
func nearest(err error) *Error {
	for err != nil {
		e, ok := err.(*Error)
		if ok {
			// A typed nil ends the search: there is nothing to read from it.
			return e
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return nil
		}
		err = u.Unwrap()
	}
	return nil
}
