package ez

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	const op = "TestNew"
	err := New(op, ECONFLICT, "An error message", nil)

	assert.NotNil(t, err)
	assert.Equal(t, err.Code, "conflict")
	assert.Equal(t, err.Message, "An error message")
	assert.Equal(t, err.Op, op)
	assert.Equal(t, err.Err, nil)
}

func TestWrap(t *testing.T) {
	const op = "TestNew"
	wrappedErr := New(op, ECONFLICT, "An error message", nil)

	const newOp = "TestWrap"
	err := Wrap(newOp, wrappedErr)

	assert.NotNil(t, err)
	assert.Equal(t, err.Code, "conflict")
	assert.Equal(t, err.Message, "An error message")
	assert.Equal(t, err.Op, newOp)
	assert.Equal(t, err.Err, wrappedErr)
}

func TestError(t *testing.T) {
	const op = "TestError"
	err := New(op, EINTERNAL, "An internal error", nil)

	msg := err.Error()

	assert.NotNil(t, err)
	assert.Equal(t, msg, "TestError: <internal> An internal error")
}

func TestErrorCode(t *testing.T) {
	const op = "TestErrorCode"
	err := New(op, EINVALID, "An invalid error", nil)

	code := ErrorCode(err)

	assert.NotNil(t, err)
	assert.Equal(t, code, "invalid")
}

func TestErrorMessage(t *testing.T) {
	const op = "TestErrorMessage"
	err := New(op, ENOTFOUND, "A not found error", nil)

	msg := ErrorMessage(err)

	assert.NotNil(t, err)
	assert.Equal(t, msg, "A not found error")
}

func TestErrorStacktrace(t *testing.T) {
	err1 := New("Op 1", EINTERNAL, "Original error", nil)
	err2 := New("Op 2", EINVALID, "Not so original error", err1)
	err3 := Wrap("Op 3", err2)

	ErrorStacktrace(err3)
}
