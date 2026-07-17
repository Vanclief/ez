package ez

// ErrorStacktrace returns a human-readable stacktrace of all nested errors,
// one per line.
func ErrorStacktrace(err error) string {
	if err == nil {
		return ""
	}
	if e, ok := err.(*Error); ok {
		if nested := ErrorStacktrace(e.Err); nested != "" {
			return e.String() + "\n" + nested
		}
		return e.String()
	}
	return err.Error()
}
