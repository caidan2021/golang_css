/*
 * @Date: 2022-08-18 11:35:01
 */
package util

import (
	"crypto/md5"
	"fmt"
)

func IsLocalEnv() bool {
	return Config.Env == "local"
}

func Md5Str(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
