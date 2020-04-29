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
err := ez.New(op, ECONFLICT, "An error message", nil)
```

Creating nested errors (used when the err is not an ez.Error):
```
const op = "TestNew"
err1 := errors.New("emit macho dwarf: elf header corrupted")
err2 := ez.New(op, ECONFLICT, "Elf header should not be corrupted", err1)
```

Creating wrapped errors:
```
const op = "OriginalOp"
err1 := ez.New(op, EINTERNAL, "Nested error", nil)

const newOp = "NewOp"
err2 := ez.Wrap(newOp, err1)
```

Return the string interpretation of an error:
```
const op = "TestError"
err := ez.New(op, EINTERNAL, "An internal error", nil)

err.Error() >> "TestError: <internal> An internal error"
```

Return the code of the root error:
```
const op = "TestErrorCode"
err := ez.New(op, EINVALID, "An invalid error", nil)

ez.ErrorCode(err) >> "invalid"
```

Return the human readable code:
```
const op = "TestErrorMessage"
err := ez.New(op, ENOTFOUND, "A not found error", nil)

ez.ErrorMessage(err) >> "A not found error"
```
