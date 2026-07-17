package ez

import (
	"runtime"
	"strings"
)

// callerOp derives the logical operation name from the function that called
// into ez: "pkg.Type.Method" for methods, "pkg.Function" for functions.
func callerOp() string {
	var pcs [1]uintptr
	// 3 skips runtime.Callers, callerOp and the exported ez function.
	if runtime.Callers(3, pcs[:]) == 0 {
		return ""
	}
	frame, _ := runtime.CallersFrames(pcs[:]).Next()
	return opFromFuncName(frame.Function)
}

// opFromFuncName trims a fully qualified function name like
// "github.com/user/project/pkg.(*Type).Method" down to "pkg.Type.Method".
func opFromFuncName(name string) string {
	if i := strings.LastIndexByte(name, '/'); i >= 0 {
		name = name[i+1:]
	}
	name = strings.ReplaceAll(name, "(*", "")
	name = strings.ReplaceAll(name, ")", "")
	// The runtime elides generic type arguments to a literal "[...]".
	name = strings.ReplaceAll(name, "[...]", "")
	return name
}
