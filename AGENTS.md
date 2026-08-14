# AGENTS.md

## Go style rules (strict)

### No explicit semicolons outside `for`

- Never write an explicit semicolon `;` in Go code.
- Exception (only): semicolons are allowed _only_ inside a 3-clause `for` header:
  `for init; condition; post { ... }`
- All other explicit-semicolon forms are forbidden, including:
  - multiple statements on one line: `a(); b()`
  - `if init; condition { ... }` (e.g. `if err := f(); err != nil { ... }`)
  - `switch init; { ... }` and type-switch init forms

### Required rewrite style

- One statement per line.
- If you need a header initializer that would require a semicolon (`if init; cond` or `switch init; expr` / `switch init; {}`), move it to its own line above the statement.

### Error handling (idiomatic)

- Use `err` as the error variable name.
- Do not create per-call error names (`parseErr`, `insertErr`, `keyErr`, etc.) when the error is immediately checked and returned.
- Keep assignment and error check adjacent (no unrelated code between them).
- `err` shadowing in a narrow inner scope is acceptable when the error is immediately checked and returned.

## Review guidelines

- Flag any explicit semicolon `;` usage in Go outside `for init; condition; post { ... }` as an issue unless its part of a SQL query.
- Flag any `if init; cond {}` / `switch init; {}` / `switch init; expr {}` usage as an issue (requires rewrite to multi-line).
- Flag per-call error renaming (`parseErr`, `insertErr`, `keyErr`, etc.) when `err` is sufficient.
- Do not flag `err` shadowing by default.
- Flag `err` shadowing only when it can change behavior or hide bugs (for example, when the outer `err` is relied on later in the same function or when assignment and check are not adjacent).
