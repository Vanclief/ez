package ez

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
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

	expected := `ez.TestErrorStacktrace <invalid> "Not so original error"
ez.TestErrorStacktrace <invalid> "Not so original error"
ez.TestErrorStacktrace <internal> ""
A plain error message`

	assert.Equal(t, expected, ErrorStacktrace(err3))
	assert.Equal(t, "", ErrorStacktrace(nil))
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

type timeoutError struct{}

func (timeoutError) Error() string { return "i/o timeout" }
func (timeoutError) Timeout() bool { return true }

func TestErrorCodeDetectsTimeoutsAndCancellations(t *testing.T) {
	assert.Equal(t, ECANCELED, ErrorCode(fmt.Errorf("request: %w", context.Canceled)))
	assert.Equal(t, ETIMEOUT, ErrorCode(fmt.Errorf("request: %w", context.DeadlineExceeded)))
	assert.Equal(t, ETIMEOUT, ErrorCode(fmt.Errorf("read: %w", os.ErrDeadlineExceeded)))
	assert.Equal(t, EINTERNAL, ErrorCode(errors.New("a plain error")))

	// Wrap picks up the detected code for non-ez errors
	assert.Equal(t, ETIMEOUT, Wrap(fmt.Errorf("request: %w", context.DeadlineExceeded)).Code)
}

func TestMixedStdlibAndEzChains(t *testing.T) {
	// Cancellation and deadline sentinels are still found through standard
	// wrappers, even when an ez error is also in the chain.
	mixed := fmt.Errorf("query: %w", Wrap(fmt.Errorf("request: %w", context.DeadlineExceeded)))

	assert.Equal(t, ETIMEOUT, ErrorCode(mixed))

	inner := Root(ENOTFOUND, "User not found").AddData("user_id", "123")
	wrapped := fmt.Errorf("lookup: %w", inner)

	// ez metadata hidden behind a foreign wrapper is not searched for.
	assert.Equal(t, EINTERNAL, ErrorCode(wrapped))
	assert.Equal(t, "An internal error has occurred. Please contact technical support.", ErrorMessage(wrapped))
	assert.Nil(t, ErrorData(wrapped))
	assert.Equal(t, EINTERNAL, ErrorCode(errors.Join(inner, errors.New("other"))))

	wrappedForeign := Wrap(wrapped)
	assert.Equal(t, EINTERNAL, wrappedForeign.Code)
	assert.Equal(t, "", wrappedForeign.Message)
	assert.Nil(t, wrappedForeign.Data)

	// Wrap preserves code, message and data in direct ez chains.
	rewrapped := Wrap(inner)
	assert.Equal(t, ENOTFOUND, rewrapped.Code)
	assert.Equal(t, "User not found", rewrapped.Message)
	assert.Equal(t, map[string]interface{}{"user_id": "123"}, rewrapped.Data)

	// AddData on the wrapper does not mutate the wrapped error
	rewrapped.AddData("extra", true)
	assert.Equal(t, map[string]interface{}{"user_id": "123"}, inner.Data)

	// The copy is shallow: nested mutable values stay shared
	nested := Root(EINVALID, "nested").AddData("tags", []string{"a"})
	Wrap(nested).Data["tags"].([]string)[0] = "changed"
	assert.Equal(t, "changed", nested.Data["tags"].([]string)[0])

	// Wrap derives empty code/message from deeper in the chain instead of
	// propagating empty fields
	codeless := &Error{Op: "handbuilt", Err: Root(ENOTFOUND, "User not found")}
	derived := Wrap(codeless)
	assert.Equal(t, ENOTFOUND, derived.Code)
	assert.Equal(t, "User not found", derived.Message)

	// With nothing in the chain, Code gets a fallback but Message stays
	// empty — ErrorMessage resolves fallback text at read time, never
	// frozen into the struct
	bare := Wrap(&Error{Op: "handbuilt"})
	assert.Equal(t, EINTERNAL, bare.Code)
	assert.Equal(t, "", bare.Message)
	assert.Equal(t, "An internal error has occurred. Please contact technical support.", ErrorMessage(bare))

	// A coded-but-messageless error never gets contradictory text baked in
	notFound := Wrap(&Error{Code: ENOTFOUND})
	assert.Equal(t, ENOTFOUND, notFound.Code)
	assert.Equal(t, "", notFound.Message)

	// ErrorStacktrace formats ez frames until the first foreign error,
	// then appends that error's complete text once and stops
	outer := New(EINVALID, "outer msg", fmt.Errorf("mid: %w", inner))
	expected := `ez.TestMixedStdlibAndEzChains <invalid> "outer msg"
mid: ez.TestMixedStdlibAndEzChains: <not_found> User not found`
	assert.Equal(t, expected, ErrorStacktrace(outer))
	assert.Equal(t, wrapped.Error(), ErrorStacktrace(wrapped))

	// Unwrap lets errors.Is see through ez layers
	sentinel := errors.New("sentinel")
	assert.True(t, errors.Is(New(EINTERNAL, "wrapped", sentinel), sentinel))
}

type falseTimeoutWrapper struct{ err error }

func (w falseTimeoutWrapper) Error() string { return "wrapped: " + w.err.Error() }
func (w falseTimeoutWrapper) Timeout() bool { return false }
func (w falseTimeoutWrapper) Unwrap() error { return w.err }

func TestTimeoutDetectionScope(t *testing.T) {
	// Sentinels are found anywhere errors.Is reaches, including joins
	assert.Equal(t, ETIMEOUT, ErrorCode(errors.Join(errors.New("other"), context.DeadlineExceeded)))
	assert.Equal(t, ECANCELED, ErrorCode(errors.Join(errors.New("other"), context.Canceled)))
	assert.Equal(t, ECANCELED, ErrorCode(errors.Join(&Error{Err: errors.New("other")}, context.Canceled)))

	// The Timeout() interface is consulted on the top-level error only —
	// wrapped or joined interface-only timeouts are out of scope
	assert.Equal(t, ETIMEOUT, ErrorCode(timeoutError{}))
	assert.Equal(t, EINTERNAL, ErrorCode(fmt.Errorf("dial: %w", timeoutError{})))
	assert.Equal(t, EINTERNAL, ErrorCode(errors.Join(timeoutError{})))
	assert.Equal(t, EINTERNAL, ErrorCode(falseTimeoutWrapper{err: timeoutError{}}))
}

func TestTypedNilErrorDoesNotPanic(t *testing.T) {
	var nilErr *Error
	wrapped := fmt.Errorf("call failed: %w", nilErr)

	assert.Equal(t, EINTERNAL, ErrorCode(wrapped))
	assert.Equal(t, "An internal error has occurred. Please contact technical support.", ErrorMessage(wrapped))
	assert.Nil(t, ErrorData(wrapped))
	assert.NotPanics(t, func() { ErrorStacktrace(wrapped) })
	assert.Equal(t, EINTERNAL, Wrap(wrapped).Code)

	// Wrapping a direct typed nil yields a usable error too
	direct := Wrap(nilErr)
	assert.Equal(t, EINTERNAL, direct.Code)
	assert.NotPanics(t, func() { _ = direct.Error() })
	assert.NotPanics(t, func() { _ = ErrorStacktrace(direct) })
	assert.NotPanics(t, func() { _ = nilErr.Error() })
	assert.NotPanics(t, func() { _ = nilErr.String() })
}

func TestErrorMessageMatchesDetectedCode(t *testing.T) {
	timedOut := Wrap(fmt.Errorf("request: %w", context.DeadlineExceeded))
	assert.Equal(t, ETIMEOUT, timedOut.Code)
	assert.Equal(t, "", timedOut.Message)
	assert.Equal(t, "The operation timed out.", ErrorMessage(timedOut))

	canceled := Wrap(fmt.Errorf("request: %w", context.Canceled))
	assert.Equal(t, ECANCELED, canceled.Code)
	assert.Equal(t, "", canceled.Message)
	assert.Equal(t, "The operation was canceled.", ErrorMessage(canceled))

	// A coded error with an empty message still resolves the code-specific
	// fallback instead of the generic internal one
	emptyMsg := NewFromGRPC(status.Error(codes.DeadlineExceeded, ""))
	assert.Equal(t, ETIMEOUT, emptyMsg.Code)
	assert.Equal(t, "The operation timed out.", ErrorMessage(emptyMsg))
}

func TestHTTPStatusToError(t *testing.T) {
	cases := map[int]string{
		http.StatusBadRequest:          EINVALID,
		http.StatusUnprocessableEntity: EINVALID,
		http.StatusMethodNotAllowed:    EINVALID, // unmapped 4xx
		http.StatusNotFound:            ENOTFOUND,
		http.StatusGone:                ENOTFOUND,
		http.StatusConflict:            ECONFLICT,
		http.StatusPreconditionFailed:  ECONFLICT,
		http.StatusUnauthorized:        ENOTAUTHENTICATED,
		http.StatusForbidden:           ENOTAUTHORIZED,
		http.StatusTooManyRequests:     ERESOURCEEXHAUSTED,
		http.StatusRequestTimeout:      ETIMEOUT,
		http.StatusGatewayTimeout:      ETIMEOUT,
		StatusClientClosedRequest:      ECANCELED,
		http.StatusNotImplemented:      ENOTIMPLEMENTED,
		http.StatusInternalServerError: EINTERNAL,
		http.StatusBadGateway:          EUNAVAILABLE,
		http.StatusServiceUnavailable:  EUNAVAILABLE,
		http.StatusInsufficientStorage: EINTERNAL, // unmapped 5xx
		http.StatusMovedPermanently:    EINTERNAL, // not an error status
	}

	for status, expected := range cases {
		assert.Equal(t, expected, HTTPStatusToError(status), "status %d", status)
	}
}

func TestErrorToHTTPStatus(t *testing.T) {
	assert.Equal(t, http.StatusConflict, ErrorToHTTPStatus(Root(ECONFLICT, "")))
	assert.Equal(t, http.StatusGatewayTimeout, ErrorToHTTPStatus(Root(ETIMEOUT, "")))
	assert.Equal(t, StatusClientClosedRequest, ErrorToHTTPStatus(Root(ECANCELED, "")))
	assert.Equal(t, http.StatusInternalServerError, ErrorToHTTPStatus(Root(EINTERNAL, "")))
}

func TestGRPCCodeMappings(t *testing.T) {
	assert.Equal(t, codes.FailedPrecondition, ErrorToGRPCCode(Root(ECONFLICT, "")))
	assert.Equal(t, codes.DeadlineExceeded, ErrorToGRPCCode(Root(ETIMEOUT, "")))
	assert.Equal(t, codes.Canceled, ErrorToGRPCCode(Root(ECANCELED, "")))

	assert.Equal(t, ECONFLICT, GRPCCodeToError(codes.Aborted))
	assert.Equal(t, ECONFLICT, GRPCCodeToError(codes.AlreadyExists))
	assert.Equal(t, ECONFLICT, GRPCCodeToError(codes.FailedPrecondition))
	assert.Equal(t, ECONFLICT, GRPCCodeToError(codes.OutOfRange))
	assert.Equal(t, ETIMEOUT, GRPCCodeToError(codes.DeadlineExceeded))
	assert.Equal(t, ECANCELED, GRPCCodeToError(codes.Canceled))
	assert.Equal(t, EINTERNAL, GRPCCodeToError(codes.Internal))
	assert.Equal(t, EINTERNAL, GRPCCodeToError(codes.Unknown))
	assert.Equal(t, EINTERNAL, GRPCCodeToError(codes.DataLoss))
}

func TestNewFromGRPCNil(t *testing.T) {
	err := NewFromGRPC(nil)

	assert.NotNil(t, err)
	assert.Equal(t, EINTERNAL, err.Code)
	assert.Empty(t, err.Message)
	assert.Equal(t, "NewFromGRPC called with nil error", err.Data["reason"])
	assert.Equal(t, "An internal error has occurred. Please contact technical support.", ErrorMessage(err))
	assert.NotPanics(t, func() { err.Op = "method" })

	var asInterface error = NewFromGRPC(nil)
	assert.NotNil(t, asInterface)
}

func TestNewFromGRPCDoesNotExposeStatusMessages(t *testing.T) {
	const upstreamMessage = "database password must not reach users"
	const genericMessage = "An internal error has occurred. Please contact technical support."

	cases := []struct {
		name        string
		grpcCode    codes.Code
		code        string
		userMessage string
	}{
		{name: "canceled", grpcCode: codes.Canceled, code: ECANCELED, userMessage: "The operation was canceled."},
		{name: "unknown", grpcCode: codes.Unknown, code: EINTERNAL, userMessage: genericMessage},
		{name: "invalid argument", grpcCode: codes.InvalidArgument, code: EINVALID, userMessage: genericMessage},
		{name: "deadline exceeded", grpcCode: codes.DeadlineExceeded, code: ETIMEOUT, userMessage: "The operation timed out."},
		{name: "not found", grpcCode: codes.NotFound, code: ENOTFOUND, userMessage: genericMessage},
		{name: "already exists", grpcCode: codes.AlreadyExists, code: ECONFLICT, userMessage: genericMessage},
		{name: "permission denied", grpcCode: codes.PermissionDenied, code: ENOTAUTHORIZED, userMessage: genericMessage},
		{name: "resource exhausted", grpcCode: codes.ResourceExhausted, code: ERESOURCEEXHAUSTED, userMessage: genericMessage},
		{name: "failed precondition", grpcCode: codes.FailedPrecondition, code: ECONFLICT, userMessage: genericMessage},
		{name: "aborted", grpcCode: codes.Aborted, code: ECONFLICT, userMessage: genericMessage},
		{name: "out of range", grpcCode: codes.OutOfRange, code: ECONFLICT, userMessage: genericMessage},
		{name: "unimplemented", grpcCode: codes.Unimplemented, code: ENOTIMPLEMENTED, userMessage: genericMessage},
		{name: "internal", grpcCode: codes.Internal, code: EINTERNAL, userMessage: genericMessage},
		{name: "unavailable", grpcCode: codes.Unavailable, code: EUNAVAILABLE, userMessage: genericMessage},
		{name: "data loss", grpcCode: codes.DataLoss, code: EINTERNAL, userMessage: genericMessage},
		{name: "unauthenticated", grpcCode: codes.Unauthenticated, code: ENOTAUTHENTICATED, userMessage: genericMessage},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			source := status.Error(tc.grpcCode, upstreamMessage)
			err := NewFromGRPC(source)

			assert.Equal(t, tc.code, err.Code)
			assert.Empty(t, err.Message)
			assert.Equal(t, source, err.Err)
			assert.Equal(t, tc.userMessage, ErrorMessage(err))
		})
	}
}

type nilGRPCStatusError struct {
	err error
}

func (e nilGRPCStatusError) Error() string            { return "invalid gRPC status" }
func (e nilGRPCStatusError) Unwrap() error            { return e.err }
func (nilGRPCStatusError) GRPCStatus() *status.Status { return nil }

func TestNewFromGRPCFallbacks(t *testing.T) {
	const genericMessage = "An internal error has occurred. Please contact technical support."

	cases := []struct {
		name        string
		err         error
		code        string
		userMessage string
	}{
		{name: "local error", err: errors.New("database password must not reach users"), code: EINTERNAL, userMessage: genericMessage},
		{name: "canceled", err: context.Canceled, code: ECANCELED, userMessage: "The operation was canceled."},
		{name: "wrapped canceled", err: fmt.Errorf("dial: %w", context.Canceled), code: ECANCELED, userMessage: "The operation was canceled."},
		{name: "deadline", err: context.DeadlineExceeded, code: ETIMEOUT, userMessage: "The operation timed out."},
		{name: "wrapped deadline", err: fmt.Errorf("dial: %w", context.DeadlineExceeded), code: ETIMEOUT, userMessage: "The operation timed out."},
		{name: "nil status", err: nilGRPCStatusError{}, code: EINTERNAL, userMessage: genericMessage},
		{name: "nil status with deadline", err: nilGRPCStatusError{err: context.DeadlineExceeded}, code: ETIMEOUT, userMessage: "The operation timed out."},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := NewFromGRPC(tc.err)

			assert.Equal(t, tc.code, err.Code)
			assert.Empty(t, err.Message)
			assert.Equal(t, tc.err, err.Err)
			assert.Equal(t, tc.userMessage, ErrorMessage(err))
		})
	}
}

func TestNewFromGRPCStatusTakesPrecedenceOverContext(t *testing.T) {
	const genericMessage = "An internal error has occurred. Please contact technical support."
	source := status.Error(codes.NotFound, "User not found")
	cases := []struct {
		name string
		err  error
	}{
		{name: "status before deadline", err: errors.Join(source, context.DeadlineExceeded)},
		{name: "deadline before status", err: errors.Join(context.DeadlineExceeded, source)},
		{name: "status before cancellation", err: errors.Join(source, context.Canceled)},
		{name: "cancellation before status", err: errors.Join(context.Canceled, source)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := NewFromGRPC(tc.err)

			assert.Equal(t, ENOTFOUND, err.Code)
			assert.Empty(t, err.Message)
			assert.Equal(t, tc.err, err.Err)
			assert.Equal(t, genericMessage, ErrorMessage(err))
		})
	}
}

func TestNewFromGRPC(t *testing.T) {
	source := status.Error(codes.NotFound, "User not found")

	err := NewFromGRPC(source)

	assert.Equal(t, "not_found", err.Code)
	assert.Empty(t, err.Message)
	assert.Equal(t, "ez.TestNewFromGRPC", err.Op)
	assert.Equal(t, source, err.Err)
}

func TestNewFromGRPCWrappedStatusDoesNotExposeMessages(t *testing.T) {
	const genericMessage = "An internal error has occurred. Please contact technical support."
	source := status.Error(codes.NotFound, "User not found")
	wrapped := fmt.Errorf("database password must not reach users: %w", source)

	err := NewFromGRPC(wrapped)

	assert.Equal(t, ENOTFOUND, err.Code)
	assert.Empty(t, err.Message)
	assert.Equal(t, genericMessage, ErrorMessage(err))
	assert.Equal(t, wrapped, err.Err)
}
