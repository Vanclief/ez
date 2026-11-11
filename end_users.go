package ez

// ErrorMessage returns the human-readable message of the error, if available.
// Otherwise returns a generic message: timeouts and cancellations get
// code-specific text, everything else the internal-error fallback.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	msg := errorMessage(err)
	if msg != "" {
		return msg
	}
	switch ErrorCode(err) {
	case ETIMEOUT:
		return "The operation timed out."
	case ECANCELED:
		return "The operation was canceled."
	}
	return "An internal error has occurred. Please contact technical support."
}

// errorMessage returns the first non-empty message set on an ez error in the
// chain, following single-cause wrappers, or "". The text of a wrapper that
// is not an ez error is never returned: only ez messages are end-user facing.
func errorMessage(err error) string {
	e := nearest(err)
	if e != nil {
		if e.Message != "" {
			return e.Message
		}
		if e.Err != nil {
			return errorMessage(e.Err)
		}
	}
	return ""
}
