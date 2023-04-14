package fail

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

var (
	GrpcUnknown = CreateGrpcClient(TypeUnknown, "unknown error")

	GrpcInvalidBody         = CreateGrpcClient(TypeInvalidBody, "invalid request body")
	GrpcBigFile             = CreateGrpcClient(TypeInvalidBody, "file too big")
	GrpcUnsupportedFileType = CreateGrpcClient(TypeInvalidBody, "unsupported file type")

	GrpcAccessDenied = CreateGrpcAccessDenied(TypeAccessDenied, "access denied")

	GrpcNotFound = CreateGrpcClient(TypeNotFound, "not found")
)

func CreateGrpcClient(errorType string, errorMessage string) error {
	return createGrpc(errorType, errorMessage, codes.InvalidArgument)
}

func CreateGrpcUnauthorized(errorType string, errorMessage string) error {
	return createGrpc(errorType, errorMessage, codes.Unauthenticated)
}

func CreateGrpcAccessDenied(errorType string, errorMessage string) error {
	return createGrpc(errorType, errorMessage, codes.PermissionDenied)
}

func CreateGrpcNotFound(errorType string, errorMessage string) error {
	return createGrpc(errorType, errorMessage, codes.NotFound)
}

func CreateGrpcServer(errorType string, errorMessage string) error {
	return createGrpc(errorType, errorMessage, codes.Internal)
}

func createGrpcByResponse(
	response Response,
) error {
	return createGrpcWithHttpCode(response.ErrorType, response.Message, response.Code)
}

func createGrpcWithHttpCode(
	errorType string,
	errorMessage string,
	code int,
) error {
	grpcCode := codes.Unknown

	switch code {
	case http.StatusBadRequest:
		grpcCode = codes.InvalidArgument
	case http.StatusUnauthorized:
		grpcCode = codes.Unauthenticated
	case http.StatusForbidden:
		grpcCode = codes.PermissionDenied
	case http.StatusNotFound:
		grpcCode = codes.NotFound
	}

	return createGrpc(errorType, errorMessage, grpcCode)
}

func createGrpc(
	errorType string,
	errorMessage string,
	code codes.Code,
) error {
	details := errdetails.ErrorInfo{Reason: errorType}
	errStatus, err := status.New(code, errorMessage).WithDetails(&details)
	if err != nil {
		return err
	}
	return errStatus.Err()
}

func ParseGrpc(err error) Response {
	response := Response{Code: http.StatusInternalServerError, ErrorType: TypeUnknown}

	errStatus, ok := status.FromError(err)
	if !ok {
		return response
	}

	switch errStatus.Code() {
	case codes.InvalidArgument:
		response.Code, response.ErrorType = http.StatusBadRequest, TypeInvalidBody
	case codes.Unauthenticated:
		response.Code, response.ErrorType = http.StatusUnauthorized, TypeUnauthorized
	case codes.PermissionDenied:
		response.Code, response.ErrorType = http.StatusForbidden, TypeAccessDenied
	case codes.NotFound:
		response.Code, response.ErrorType = http.StatusNotFound, TypeNotFound
	}

	response.Message = errStatus.Message()

	if len(errStatus.Details()) > 0 {
		detail := errStatus.Details()[0]
		if metadata, ok := detail.(*errdetails.ErrorInfo); ok {
			response.ErrorType = metadata.Reason
		} else {
			return Response{}
		}
	}

	return response
}