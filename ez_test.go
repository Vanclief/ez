package ez

import (
	"errors"
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
	err0 := errors.New("A plain error message")
	err1 := New("Depth 1", EINTERNAL, "", err0)
	err2 := New("Depth 2", EINVALID, "Not so original error", err1)
	err3 := Wrap("Depth 3", err2)

	ErrorStacktrace(err3)
}

func TestAddDataToNilData(t *testing.T) {
	err := New("TestOp", ECONFLICT, "message", nil)
	err.AddData("key1", "value1")

	expected := map[string]interface{}{"key1": "value1"}
	assert.Equal(t, expected, err.Data)
}

func TestAddDataToExistingData(t *testing.T) {
	err := &Error{Data: map[string]interface{}{"existing": "data"}}
	err.AddData("key2", 42)

	expected := map[string]interface{}{
		"existing": "data",
		"key2":     42,
	}
	assert.Equal(t, expected, err.Data)
}

func TestAddDataOverrideExisting(t *testing.T) {
	err := &Error{Data: map[string]interface{}{"key": "old"}}
	err.AddData("key", "new")

	expected := map[string]interface{}{"key": "new"}
	assert.Equal(t, expected, err.Data)
}
