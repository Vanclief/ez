package ez

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// NewFromGRPC wraps a GRPC Error into a standar application error
func NewFromGRPC(op string, err error) *Error {
	status := status.Convert(err)
	code := status.Code()
	return New(op, GRPCCodeToError(code), status.Message(), err)
}

// GRPCCodeToError converts a GRPC error code to a standar application error code
func GRPCCodeToError(c codes.Code) string {
	switch c {
	case codes.FailedPrecondition:
		return ECONFLICT
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
