/*
 * @Date: 2022-08-19 16:44:11
 */
package web

import (
	"gin/util"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	uploadFileKey = "file"
	showImgKey    = "imgName"
)

func Upload(ctx *gin.Context) {

	file, err := ctx.FormFile(uploadFileKey)
	if err != nil {
		errorMsg := util.ValidatorError(err)
		ctx.JSON(http.StatusOK, util.FailedRespPackage(errorMsg))
	}

	newFileName := util.Md5Str(strconv.FormatInt(time.Now().Unix(), 10)) + "_" + file.Filename
	dst := path.Join(util.Config.FileStoragePath, newFileName)
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		errorMsg := util.ValidatorError(err)
		ctx.JSON(http.StatusOK, util.FailedRespPackage(errorMsg))
		return
	}

	rt := "http://" + util.Config.FileAddr + "/img/" + newFileName

	ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{"path": rt}))
	return
}

func ShowImage(c *gin.Context) {
	imageName := c.Param(showImgKey)
	fullPath := path.Join(util.Config.FileStoragePath, imageName)
	c.File(fullPath)
}
