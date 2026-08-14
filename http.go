package ez

import (
	"net/http"
)

// StatusClientClosedRequest is nginx's non-standard status for requests
// canceled by the client. net/http has no constant for it.
const StatusClientClosedRequest = 499

// HTTPStatusToError converts a HTTP status code to a standar application error
// code. Only statuses that indicate an unavailable service map to EUNAVAILABLE.
// Other 5xx statuses remain EINTERNAL because retryability is not implied.
func HTTPStatusToError(status int) string {
	switch status {
	case http.StatusConflict, http.StatusPreconditionFailed:
		return ECONFLICT
	case http.StatusNotFound, http.StatusGone:
		return ENOTFOUND
	case http.StatusForbidden:
		return ENOTAUTHORIZED
	case http.StatusUnauthorized:
		return ENOTAUTHENTICATED
	case http.StatusTooManyRequests:
		return ERESOURCEEXHAUSTED
	case http.StatusNotImplemented:
		return ENOTIMPLEMENTED
	case http.StatusBadGateway, http.StatusServiceUnavailable:
		return EUNAVAILABLE
	case http.StatusRequestTimeout, http.StatusGatewayTimeout:
		return ETIMEOUT
	case StatusClientClosedRequest:
		return ECANCELED
	}
	switch {
	case status >= 400 && status <= 499:
		return EINVALID
	case status >= 500 && status <= 599:
		return EINTERNAL
	}
	return EINTERNAL
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
		return http.StatusForbidden
	case ENOTAUTHENTICATED:
		return http.StatusUnauthorized
	case ERESOURCEEXHAUSTED:
		return http.StatusTooManyRequests
	case ENOTIMPLEMENTED:
		return http.StatusNotImplemented
	case EUNAVAILABLE:
		return http.StatusServiceUnavailable
	case ETIMEOUT:
		return http.StatusGatewayTimeout
	case ECANCELED:
		return StatusClientClosedRequest
	default:
		return http.StatusInternalServerError
	}
}
