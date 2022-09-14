/*
 * @Date: 2022-09-13 15:51:38
 */
package weitu

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"gin/util"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	coderCheckFlagOfWait         = 1
	coderCheckFlagOfOk           = 2
	coderCheckFlagOfNeedStr      = 3
	coderCheckFlagOfNotSupported = 4
)

type LocationInfo struct {
	Name         string
	Longitude    string // 经度
	Latitude     string // 纬度
	NewLongitude string // 经度
	NewLatitude  string // 纬度
}

type WeiTu struct {
}

func NewWeiTu() *WeiTu {
	return &WeiTu{}
}

// 按照行读取text文件
func (WeiTu) ReadTextFileByLine(textFile string) ([]*string, error) {

	var textInfo []*string

	file, err := os.Open(textFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var coderCheckFlag = coderCheckFlagOfWait
	for scanner.Scan() {
		line := scanner.Text()
		if coderCheckFlag == coderCheckFlagOfWait {
			if util.IsUtf8([]byte(line)) {
				coderCheckFlag = coderCheckFlagOfOk
			} else if util.IsGBK([]byte(line)) {
				coderCheckFlag = coderCheckFlagOfNeedStr
			} else {
				coderCheckFlag = coderCheckFlagOfNotSupported
			}
		}

		if coderCheckFlag == coderCheckFlagOfNotSupported {
			return nil, fmt.Errorf("not supported: %v", line)
		}
		if coderCheckFlag == coderCheckFlagOfNeedStr {
			newLine, err := util.GbkToUtf8([]byte(line))
			if err != nil {
				return nil, err
			}
			line = string(newLine)
		}
		textInfo = append(textInfo, &line)
	}

	if err := scanner.Err(); err != nil {
		return textInfo, err
	}

	return textInfo, nil
}

func (w WeiTu) ToCsv(data []LocationInfo, storagePath string) (string, error) {

	strTime := time.Now().Format("20220608181234")
	fileName := fmt.Sprintf("location_info %s.xlsx", strTime)

	xlsFile, err := os.OpenFile(storagePath+fileName, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return "", err
	}

	defer xlsFile.Close()

	xlsFile.WriteString("\xEF\xBB\xBF")
	wStr := csv.NewWriter(xlsFile)
	wStr.Write([]string{"名称", "经度", "纬度", "经度（新）", "纬度（新）"})

	for _, v := range data {
		wStr.Write([]string{v.Name, v.Longitude, v.Latitude, v.NewLongitude, v.NewLatitude})
	}
	wStr.Flush()

	return fileName, nil

}

func (WeiTu) TransformData(textInfo []*string) ([]LocationInfo, error) {

	var data []LocationInfo
	keySortArr := strings.Split(*textInfo[0], ",")

	value := []string{"名称", "经度", "纬度"}
	var key []int

	for _, v := range value {
		k := util.GetKeyFromArr(keySortArr, v, true)
		if k == -1 {
			return data, fmt.Errorf("%s 没有找到", v)
		}
		key = append(key, k)
	}
	fmt.Println(key)

	for kItem, item := range textInfo {
		if kItem == 0 {
			continue
		}

		var perData LocationInfo
		strArr := strings.Split(*item, ",")

		for k, v := range key {
			if k != 0 {
				if len(strArr) <= 3 {
					continue
				}

				newStrArr := strings.Split(strings.Replace(strArr[v], "″", "", -1), "′")
				newLtF, err := strconv.ParseFloat(newStrArr[1], 64)

				if err != nil {
					return data, fmt.Errorf("%s 数据转换失败， e：%s", newStrArr[1], err.Error())
				}

				newStr := fmt.Sprintf("%s′%.3f″", newStrArr[0], newLtF)
				if k == 1 {
					perData.Longitude = strArr[k]
					perData.NewLongitude = newStr

				} else {
					perData.Latitude = strArr[k]
					perData.NewLatitude = newStr
				}
				continue
			}
			perData.Name = strArr[k]
		}
		data = append(data, perData)
	}
	return data, nil
}
