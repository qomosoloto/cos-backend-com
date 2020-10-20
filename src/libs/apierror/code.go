package apierror

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/lib/pq"
	"github.com/qiniu/x/rpc"
	"github.com/wujiu2020/strip/utils/apires"
	"gopkg.in/go-playground/validator.v9"
)

var (
	errors = map[int]*apires.ResError{} // code -> Error, not allowed dup code
)

func New(status, code int, message string) *apires.ResError {
	if _, ok := errors[code]; ok {
		panic(fmt.Sprintf("apierror code:%d, message:`%s` has been taken", code, message))
	}
	e := apires.NewResError(status, code, message)
	errors[code] = e
	return e
}

// 通用错误
var (
	ErrNotFound        = New(404, 404, "not found")
	ErrBadRequest      = New(400, 400, "bad request")
	ErrNeedLogin       = New(401, 401, "need login")
	ErrNoPermission    = New(403, 403, "forbidden")
	ErrTooLargeRequest = New(413, 413, "request entity too large")
	ErrClientClosed    = New(499, 499, "client closed request")
	ErrRequestTimeout  = New(500, 408, "internal request timeout")
	ErrServerError     = New(500, 500, "internal server error")
	ErrNotImplemented  = New(501, 501, "not implemented")
	_                  = 0
	ErrUnkown          = New(599, 90000, "unkown")
	ErrDBDup           = New(400, 2000, "dup err")
	ErrDBInternalError = New(500, 2001, "db internal error")
	ErrDBDeadLockError = New(500, 2002, "db dead lock error")
	ErrDBNullViolation = New(400, 2003, "db not null violation")

	ErrIOError = New(500, 3001, "io error")
)

// 业务错误
var (
	ErrApiAuthInvalidHeader = New(401, 20000, "authorization not in request headers")
	ErrApiAuthInvalidToken  = New(401, 20001, "invalid token")
	ErrApiAuthInvalidType   = New(401, 20002, "invalid auth type")
	ErrApiAuthInvalidInfo   = New(401, 20003, "invalid auth info")
	ErrApiAuthExpiredToken  = New(401, 20004, "expired token")
	ErrApiAuthCreateFailed  = New(403, 20008, "token create failed")
	ErrInvalidSignature     = New(401, 20009, "invalid signature")
)

func HandleError(err error) *apires.ResError {
	switch err {
	case sql.ErrNoRows:
		return ErrNotFound
	case context.Canceled:
		return ErrClientClosed
	}

	switch er := err.(type) {
	case *apires.ResError:
		return er

	case *rpc.ErrorInfo:
		v := apires.NewResError(er.HttpCode(), er.Errno, er.Err)
		v.Data = er.Err
		return v

	case *url.Error:
		return HandleError(er.Err)

	case *pq.Error:
		switch er.Code {
		case "23505":
			return ErrDBDup
		case "23502":
			return ErrDBNullViolation
		}

	case validator.ValidationErrors:
		return ErrBadRequest.WithData(validationErrorConvert(err), err.Error())

	case *FieldError:
		return ResponseFormFieldError(er)

	case *FormError:
		return ErrBadRequest.WithData(er, er.Error())

	case nil:
		panic("err cannot be nil")
	}

	return ErrServerError
}
