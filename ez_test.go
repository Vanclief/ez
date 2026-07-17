package ez

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func TestNew(t *testing.T) {
	err := New(ECONFLICT, "An error message", nil)

	assert.NotNil(t, err)
	assert.Equal(t, "conflict", err.Code)
	assert.Equal(t, "An error message", err.Message)
	assert.Equal(t, "ez.TestNew", err.Op)
	assert.Equal(t, nil, err.Err)
}

func TestRoot(t *testing.T) {
	err := Root(ENOTFOUND, "A not found error")

	assert.NotNil(t, err)
	assert.Equal(t, "not_found", err.Code)
	assert.Equal(t, "A not found error", err.Message)
	assert.Equal(t, "ez.TestRoot", err.Op)
	assert.Equal(t, nil, err.Err)
}

func TestWrap(t *testing.T) {
	wrappedErr := New(ECONFLICT, "An error message", nil)

	err := Wrap(wrappedErr)

	assert.NotNil(t, err)
	assert.Equal(t, "conflict", err.Code)
	assert.Equal(t, "An error message", err.Message)
	assert.Equal(t, "ez.TestWrap", err.Op)
	assert.Equal(t, wrappedErr, err.Err)
}

type testService struct{}

func (s *testService) fail() *Error {
	return Root(EINTERNAL, "Something failed")
}

func TestMethodOp(t *testing.T) {
	svc := &testService{}
	err := svc.fail()

	assert.Equal(t, "ez.testService.fail", err.Op)
}

func TestClosureOp(t *testing.T) {
	fn := func() *Error {
		return Root(EINTERNAL, "Something failed")
	}
	err := fn()

	assert.Equal(t, "ez.TestClosureOp.func1", err.Op)
}

func TestError(t *testing.T) {
	err := New(EINTERNAL, "An internal error", nil)

	msg := err.Error()

	assert.NotNil(t, err)
	assert.Equal(t, "ez.TestError: <internal> An internal error", msg)
}

func TestErrorCode(t *testing.T) {
	err := New(EINVALID, "An invalid error", nil)

	code := ErrorCode(err)

	assert.NotNil(t, err)
	assert.Equal(t, "invalid", code)
}

func TestErrorMessage(t *testing.T) {
	err := New(ENOTFOUND, "A not found error", nil)

	msg := ErrorMessage(err)

	assert.NotNil(t, err)
	assert.Equal(t, "A not found error", msg)
}

func TestErrorStacktrace(t *testing.T) {
	err0 := errors.New("A plain error message")
	err1 := New(EINTERNAL, "", err0)
	err2 := New(EINVALID, "Not so original error", err1)
	err3 := Wrap(err2)

	ErrorStacktrace(err3)
}

func TestOpFromFuncName(t *testing.T) {
	cases := map[string]string{
		"github.com/vanclief/project/pkg.(*API).Create": "pkg.API.Create",
		"github.com/vanclief/project/pkg.API.Create":    "pkg.API.Create",
		"github.com/vanclief/project/pkg.monthBounds":   "pkg.monthBounds",
		"main.main":                   "main.main",
		"pkg.(*API).Create.func1":     "pkg.API.Create.func1",
		"pkg.Load[...]":               "pkg.Load",
		"pkg.(*Repository[...]).Find": "pkg.Repository.Find",
	}

	for input, expected := range cases {
		assert.Equal(t, expected, opFromFuncName(input))
	}
}

func TestAddDataToNilData(t *testing.T) {
	err := New(ECONFLICT, "message", nil)
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

func TestWrapPreservesData(t *testing.T) {
	rootErr := Root(ENOTFOUND, "User not found").AddData("user_id", "123")

	err := Wrap(rootErr)

	expected := map[string]interface{}{"user_id": "123"}
	assert.Equal(t, expected, err.Data)
}

func TestNewFromGRPC(t *testing.T) {
	grpcErr := status.Error(codes.NotFound, "User not found")

	err := NewFromGRPC(grpcErr)

	assert.Equal(t, "not_found", err.Code)
	assert.Equal(t, "User not found", err.Message)
	assert.Equal(t, "ez.TestNewFromGRPC", err.Op)
	assert.Equal(t, grpcErr, err.Err)
}
