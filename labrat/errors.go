// labrat.Error is a simple custom error, designed to contain
package labrat

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type Error interface {
	Error() string
	Code() ErrorCode
	WithMeta(key, value string) Error
	Meta(key string) string
	MetaMap()map[string]string
}

type wraperr struct {
	wrapper Error
	cause error
}

type laberr struct {
	code ErrorCode
	meta map[string]string
}

func NewError(code ErrorCode, message string) Error {
	return &wraperr{
		wrapper: &laberr{
			code: code,
			meta: map[string]string{},
		},
		cause: errors.New(message),
	}
	//return ErrorWith(errors.New(message))
}

type ErrorCode string

func (c ErrorCode) String() string { return string(c) }

func (c ErrorCode) HTTP() int {
	switch c {
	case NotFound:
		return http.StatusNotFound
	case NoError:
		return http.StatusOK
	case InternalError:
		return http.StatusInternalServerError
	default:
		return 0
	}
}

func ErrorWith(err error) Error {

	msg := fmt.Sprintf("%+v",err)

	laberr := NewError(InternalError, msg)

	return &wraperr{
		wrapper: laberr,
		cause: err,
	}

}
func (l laberr) MetaMap()map[string]string { return l.meta }

func (e wraperr) MetaMap()map[string]string { return e.wrapper.MetaMap() }

func (l *laberr) Code() ErrorCode {
	if l != nil {
		return l.code
	}
	return NoError
}

func (l laberr) Error() string { return fmt.Sprintf("code: %d, with meta: %+v", l.Code().HTTP(), l.MetaMap()) }

func (l laberr) Meta(key string) string {
	if l.meta != nil {
		return l.meta[key]
	}
	return ""
}

func (e wraperr) Error() string          { return fmt.Sprintf("%+v\n%+v", e.wrapper.Error(), e.cause) }

func (e *wraperr) Code() ErrorCode        {
	if e != nil {
		return e.wrapper.Code()
	}
	return NoError
}

func (e wraperr) Meta(key string) string { return e.wrapper.Meta(key) }

func (l laberr) WithMeta(key, value string) Error {
	err :=  &laberr{
		meta: make(map[string]string, len(l.meta)),
		code: l.code,
	}

	for k, v := range l.meta {
		err.meta[k] = v
	}

	err.meta[key] = value

	return err

}

func (e wraperr) WithMeta(key, value string) Error {
	return &wraperr{
		cause: e.cause,
		wrapper: e.wrapper.WithMeta(key, value),
	}
}

func (e *wraperr) Cause() error { return e.cause }

func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}


const (
	NotFound         ErrorCode = "not_found"
	DeadlineExceeded ErrorCode = "deadline_exceeded"
	InternalError ErrorCode = "internal_server_error"
	// ...
	NoError          ErrorCode = ""
)

