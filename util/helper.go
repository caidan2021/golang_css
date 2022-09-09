/*
 * @Date: 2022-08-18 11:35:01
 */
package util

import (
	"crypto/md5"
	"fmt"
)

type TwoDimensionFmt struct {
	Key   interface{} `json:"key"`
	Value interface{} `json:"value"`
}

func IsLocalEnv() bool {
	return Config.Env == "local"
}

func Md5Str(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func RenderTwoDimensionForIntStr(arr map[int64]string) []TwoDimensionFmt {
	rt := []TwoDimensionFmt{}
	for key, value := range arr {
		rt = append(rt, TwoDimensionFmt{Key: key, Value: value})
	}
	return rt
}
