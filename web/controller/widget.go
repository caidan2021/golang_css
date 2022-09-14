/*
 * @Date: 2022-09-13 15:18:52
 */
package controller

import (
	"fmt"
	"gin/util"
	"net/http"

	"gin/component/weitu"

	"github.com/gin-gonic/gin"
)

func WeiTuFmt(ctx *gin.Context) {

	// 先获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage(err.Error()))
		return
	}
	// 暂存上传的文件
	storagePath := "./storage/"
	ctx.SaveUploadedFile(file, storagePath+file.Filename)

	// defer删除暂存的文件
	defer util.DelStoreFile(storagePath, file.Filename)

	weituHandle := weitu.NewWeiTu()

	// 读取暂存的文件内容
	textInfo, err := weituHandle.ReadTextFileByLine(storagePath + file.Filename)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage(err.Error()))
		return
	}

	// 转换数据
	data, err := weituHandle.TransformData(textInfo)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage(err.Error()))
		return
	}

	// 写入excel中
	fileName, err := weituHandle.ToCsv(data, storagePath)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage(err.Error()))
		return
	}

	// defer删除暂存的excel文件
	defer util.DelStoreFile(storagePath, fileName)

	// 返回给前端下载excel
	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.File(storagePath + fileName)

}
