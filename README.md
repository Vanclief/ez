# ez

Package ez provides an opinionated, simple way to manage errors in Golang. 
Based on [this awesome post](https://middlemost.com/failure-is-your-domain/).

## Features

* Extremely simple to use.

* Allows nesting errors, so you can create a stacktrace.

* Can receive your own Application error codes.


## Usage

Import the library:

`import "github.com/vanclief/ez"`

Creating an error:

```
const op = "TestNew"
err := New(op, ECONFLICT, "An error message", nil)
```

Creating nested errors:

```
const op = "TestNew"
err1 := New(op, EINTERNAL, "Nested error", nil)
err2 := New(op, ECONFLICT, "Another error", err1)
```

Return the string interpretation of an error:
```
const op = "TestError"
err := New(op, EINTERNAL, "An internal error", nil)

err.Error() >> "TestError: <internal> An internal error"
```

Return the code of the root error:
```
const op = "TestErrorCode"
err := New(op, EINVALID, "An invalid error", nil)

ErrorCode(err) >> "invalid"
```

Return the human readable code:
```
const op = "TestErrorMessage"
err := New(op, ENOTFOUND, "A not found error", nil)

ErrorMessage(err) >> "A not found error"
```