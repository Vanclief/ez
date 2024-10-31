package ez

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

// ErrorData returns the data of the root error, if available.
func ErrorData(err error) map[string]interface{} {
	if err == nil {
		return nil
	} else if e, ok := err.(*Error); ok && e.Data != nil {
		return e.Data
	} else if ok && e.Err != nil {
		return ErrorData(e.Err)
	}
	return nil
}
