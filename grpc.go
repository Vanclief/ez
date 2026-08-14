package ez

import (
	"errors"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// NewFromGRPC wraps a GRPC Error into a standar application error.
// The operation is derived from the calling function. Status descriptions
// remain in Err and are never copied into the end-user-facing Message.
func NewFromGRPC(err error) *Error {
	if err == nil {
		// Converting a nil error is a caller bug. Returning a typed nil
		// would poison error interfaces, so fail loudly instead. The
		// diagnostic goes in Data — Message is end-user facing.
		e := &Error{Op: callerOp(), Code: EINTERNAL}
		return e.AddData("reason", "NewFromGRPC called with nil error")
	}
	// Read GRPCStatus directly: status.FromError replaces the message of a
	// wrapped status with the complete outer err.Error() text.
	type statusProvider interface {
		GRPCStatus() *status.Status
	}
	var provider statusProvider
	if errors.As(err, &provider) {
		grpcStatus := provider.GRPCStatus()
		if grpcStatus != nil {
			grpcCode := grpcStatus.Code()
			return &Error{
				Op:   callerOp(),
				Code: GRPCCodeToError(grpcCode),
				Err:  err,
			}
		}
	}
	// Raw context errors do not implement GRPCStatus.
	code := ErrorCode(err)
	switch code {
	case ECANCELED, ETIMEOUT:
		return &Error{Op: callerOp(), Code: code, Err: err}
	}
	return &Error{Op: callerOp(), Code: EINTERNAL, Err: err}
}

// GRPCCodeToError converts a GRPC error code to a standar application error code
func GRPCCodeToError(c codes.Code) string {
	switch c {
	case codes.Aborted, codes.AlreadyExists, codes.FailedPrecondition, codes.OutOfRange:
		return ECONFLICT
	case codes.DeadlineExceeded:
		return ETIMEOUT
	case codes.Canceled:
		return ECANCELED
	case codes.Internal:
		return EINTERNAL
	case codes.InvalidArgument:
		return EINVALID
	case codes.NotFound:
		return ENOTFOUND
	case codes.PermissionDenied:
		return ENOTAUTHORIZED
	case codes.Unauthenticated:
		return ENOTAUTHENTICATED
	case codes.ResourceExhausted:
		return ERESOURCEEXHAUSTED
	case codes.Unimplemented:
		return ENOTIMPLEMENTED
	case codes.Unavailable:
		return EUNAVAILABLE
	default:
		return EINTERNAL
	}
}

// ErrorToGRPCCode converts an standar application error code to a GRPC error code
func ErrorToGRPCCode(err error) codes.Code {
	code := ErrorCode(err)
	switch code {
	case ECONFLICT:
		return codes.FailedPrecondition
	case ETIMEOUT:
		return codes.DeadlineExceeded
	case ECANCELED:
		return codes.Canceled
	case EINTERNAL:
		return codes.Internal
	case EINVALID:
		return codes.InvalidArgument
	case ENOTFOUND:
		return codes.NotFound
	case ENOTAUTHORIZED:
		return codes.PermissionDenied
	case ENOTAUTHENTICATED:
		return codes.Unauthenticated
	case ERESOURCEEXHAUSTED:
		return codes.ResourceExhausted
	case ENOTIMPLEMENTED:
		return codes.Unimplemented
	case EUNAVAILABLE:
		return codes.Unavailable
	default:
		return codes.Internal
	}
}
