/*
 * @Date: 2022-08-18 11:45:45
 */
package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type BaseRespPackage struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	Success = 0
	Error   = 1

	DefaultPage     = 0
	DefaultPageSize = 15

	// 非法字符
	IllegalStr = "\n\r\t*/\\"
)

// page & pagesize change
func Paging(page string, total string) (start int, end int) {
	var i, v int
	i, _ = strconv.Atoi(page)
	v, _ = strconv.Atoi(total)
	start = i*v - v
	end = v
	return
}

func ValidatorError(err error) string {
	validatorError, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}
	msg := make([]string, len(validatorError))
	for _, v := range validatorError {
		msg = append(msg, fmt.Sprintf("提交的参数[%s]值不正确 ; ", v.Field()))
	}
	return strings.Join(msg, "")
}

func SuccessRespPackage(data interface{}) *BaseRespPackage {
	return &BaseRespPackage{Code: Success, Data: data, Msg: "success"}
}

func FailedRespPackage(msg string) *BaseRespPackage {
	return &BaseRespPackage{Code: Error, Data: nil, Msg: msg}
}
