<!-- BEGIN golang: generated from Vanclief/skills golang/SKILL.md, do not edit here -->
## Go conventions

Apply these Go conventions when writing code and report violations in reviews, even when review guidance says to skip routine lint or formatting nits.

Write one statement per line. Never use an explicit `;` in Go syntax except the two separators in a three-clause `for init; condition; post` header. Semicolons inside SQL strings are allowed. Put header initializers on the preceding line: write `err := f()` then `if err != nil {`, never `if err := f(); err != nil {`.

Use `err` for error variables. Do not invent per-call names such as `parseErr` when immediately checking and returning the error. Keep assignment and error check adjacent. Narrow `err` shadowing is acceptable when checked and returned immediately; flag shadowing only when it can hide a bug.

Where the module already uses `github.com/vanclief/ez`, create errors at the origin with `ez.New(code, message, cause)` and propagate them with `ez.Wrap(err)`. Classify ez errors with `ez.ErrorCode(err)`, not by matching error messages.

Before finishing Go changes, scan the code you wrote for forbidden semicolons in Go syntax and rewrite it to comply.
<!-- END golang -->
