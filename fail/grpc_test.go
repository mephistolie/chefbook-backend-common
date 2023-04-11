package fail

import (
	"github.com/mephistolie/chefbook-backend-auth/pkg/logger"
	"net/http"
	"testing"
)

const (
	testErrorType    = "test_error_type"
	testErrorMessage = "test error message"
)

var testResponse = Response{
	ErrorType: testErrorType,
	Message:   testErrorMessage,
}

func TestGrpcBadRequest(t *testing.T) {
	expected := testResponse
	expected.Code = http.StatusBadRequest
	testGrpcErrorTransfer(expected, t)
}

func TestGrpcUnauthorized(t *testing.T) {
	expected := testResponse
	expected.Code = http.StatusUnauthorized
	testGrpcErrorTransfer(expected, t)
}

func TestGrpcAccessDenied(t *testing.T) {
	expected := testResponse
	expected.Code = http.StatusNotFound
	testGrpcErrorTransfer(expected, t)
}

func TestGrpcNotFound(t *testing.T) {
	expected := testResponse
	expected.Code = http.StatusNotFound
	testGrpcErrorTransfer(expected, t)
}

func TestGrpcServer(t *testing.T) {
	expected := testResponse
	expected.Code = http.StatusInternalServerError
	testGrpcErrorTransfer(expected, t)
}

func testGrpcErrorTransfer(
	expected Response,
	t *testing.T,
) {
	output := ParseGrpc(createGrpcByResponse(expected))

	if expected.ErrorType != output.ErrorType || expected.Code != output.Code ||
		expected.Message != output.Message {
		logger.Error(output)
		t.Fail()
	}
}
