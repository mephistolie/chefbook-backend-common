package fail

var (
	HttpUnknown = Response{Code: 400, ErrorType: TypeUnknown, Message: "unknown error"}

	HttpInvalidBody         = Response{Code: 400, ErrorType: TypeInvalidBody, Message: "invalid request body"}
	HttpBigFile             = Response{Code: 400, ErrorType: TypeInvalidBody, Message: "file too big"}
	HtppUnsupportedFileType = Response{Code: 400, ErrorType: TypeInvalidBody, Message: "unsupported file type"}

	HttpAccessDenied = Response{Code: 403, ErrorType: TypeAccessDenied, Message: "access denied"}

	HttpNotFound = Response{Code: 404, ErrorType: TypeNotFound, Message: "not found"}
)
