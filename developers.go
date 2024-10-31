package ez

import "fmt"

// ErrorStacktrace prints a human-redable stacktrace of all nested errors.
func ErrorStacktrace(err error) {
	if err == nil {
		return
	} else if e, ok := err.(*Error); ok {
		fmt.Println(e.String())
		ErrorStacktrace(e.Err)
	} else if ok && e.Err != nil {
		fmt.Println(e.String())
	} else {
		fmt.Println(err.Error())
	}
}
