package result

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	CommonOk    = "API-COMMON-SUCCESS"
	CommonError = "API-COMMON-ERROR"
	ParamError  = "API-PARAM-ERROR"
)

type ApiResult[T any] struct {
	Content       T      `json:"content"`
	Count         int    `json:"count"`
	StatusCode    string `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
}

func (result *ApiResult[T]) IsSuccess() bool {
	return result.StatusCode == CommonOk
}

func SuccessResult[T any](content T) (int, ApiResult[T]) {
	return http.StatusOK, ApiResult[T]{
		Content:       content,
		Count:         0,
		StatusCode:    CommonOk,
		StatusMessage: CommonOk,
	}
}

func DefaultErrorResult() (int, ApiResult[any]) {
	return http.StatusInternalServerError, ApiResult[any]{
		Count:         0,
		StatusCode:    CommonError,
		StatusMessage: CommonError,
	}
}

func ErrorResultWithCode(statusCode string) (int, ApiResult[any]) {
	return http.StatusInternalServerError, ApiResult[any]{
		Count:         0,
		StatusCode:    statusCode,
		StatusMessage: CommonError,
	}
}

func ErrorResultWithMsg(statusMsg string) (int, ApiResult[any]) {
	return http.StatusInternalServerError, ApiResult[any]{
		Count:         0,
		StatusCode:    CommonError,
		StatusMessage: statusMsg,
	}
}

func ErrorResultWithCodeAndMsg(statusCode string, statusMsg string) (int, ApiResult[any]) {
	return http.StatusInternalServerError, ApiResult[any]{
		Count:         0,
		StatusCode:    statusCode,
		StatusMessage: statusMsg,
	}
}

func ParamErrorResult() (int, ApiResult[any]) {
	return http.StatusUnprocessableEntity, ApiResult[any]{
		Count:         0,
		StatusCode:    ParamError,
		StatusMessage: ParamError,
	}
}

func ReturnResult(context *gin.Context) {
	//context.JSON()
}
