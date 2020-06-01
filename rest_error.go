package ez

import (
	"net/http"
)

// HTTPStatusToError converts a HTTP error code to a standar application error code
func HTTPStatusToError(status int) string {
	switch status {
	case http.StatusConflict:
		return ECONFLICT
	case http.StatusInternalServerError:
		return EINTERNAL
	case http.StatusBadRequest:
		return EINVALID
	case http.StatusNotFound:
		return ENOTFOUND
	case http.StatusUnauthorized:
		return ENOTAUTHORIZED
	case http.StatusForbidden:
		return ENOTAUTHENTICATED
	case http.StatusTooManyRequests:
		return ERESOURCEEXHAUSTED
	case http.StatusNotImplemented:
		return ENOTIMPLEMENTED
	case http.StatusServiceUnavailable:
		return EUNAVAILABLE
	default:
		return EINTERNAL
	}
}

// ErrorToHTTPStatus converts an standar application error code to a HTTP status
func ErrorToHTTPStatus(err error) int {
	code := ErrorCode(err)
	switch code {
	case ECONFLICT:
		return http.StatusConflict
	case EINTERNAL:
		return http.StatusInternalServerError
	case EINVALID:
		return http.StatusBadRequest
	case ENOTFOUND:
		return http.StatusNotFound
	case ENOTAUTHORIZED:
		return http.StatusUnauthorized
	case ENOTAUTHENTICATED:
		return http.StatusForbidden
	case ERESOURCEEXHAUSTED:
		return http.StatusTooManyRequests
	case ENOTIMPLEMENTED:
		return http.StatusNotImplemented
	case EUNAVAILABLE:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
