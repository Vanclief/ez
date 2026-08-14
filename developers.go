package ez

// ErrorStacktrace returns a human-readable stacktrace of all nested errors,
// one per line: structured ez frames until the first foreign error, whose
// complete text is appended once. The trace is deliberately lossy past that
// point — ez frames below a foreign wrapper appear inside its flat text,
// not as structured lines.
func ErrorStacktrace(err error) string {
	if err == nil {
		return ""
	}
	e, ok := err.(*Error)
	if ok {
		if e == nil {
			return ""
		}
		nested := ErrorStacktrace(e.Err)
		if nested != "" {
			return e.String() + "\n" + nested
		}
		return e.String()
	}
	return err.Error()
}
