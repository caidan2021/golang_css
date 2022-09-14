/*
 * @Date: 2022-09-13 15:34:12
 */
package util

import (
	"os"
)

func DelStoreFile(storagePath string, fileName string) (bool, error) {
	err := os.Remove(storagePath + fileName)
	if err != nil {
		return false, err
	}
	return true, nil

}
