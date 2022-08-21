/*
 * @Date: 2022-08-19 10:42:26
 */
package controller

import (
	"gin/drivers"
	"gin/models"
	"gin/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddAdminUser(ctx *gin.Context) {
	type req struct {
		ID            int64  `json:"id"`
		Name          string `json:"name" binding:"required"`
		Password      string `json:"password" binding:"required"`
		Email         string `json:"email"`
		RememberToken string `json:"rememberToken"`
	}

	r := req{}
	if err := ctx.ShouldBindJSON(&r); err != nil {
		errorMsg := util.ValidatorError(err)
		ctx.JSON(http.StatusOK, util.FailedRespPackage(errorMsg))
		return
	}

	adminUser := new(models.AdminUser)

	if adminUser.NameExists(r.Name, r.ID) {
		ctx.JSON(http.StatusOK, util.FailedRespPackage("用户名已存在"))
		return
	}

	adminUser.Name = r.Name
	adminUser.Password = adminUser.EncryptionPw(r.Password)
	adminUser.Email = r.Email
	adminUser.RememberToken = r.RememberToken

	if err := drivers.Mysql().Create(adminUser).Error; err != nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage("保存用户信息错误"))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessRespPackage(true))
	return
}

func LoginAdmin(ctx *gin.Context) {
	type req struct {
		UserName      string `json:"username" binding:"required"`
		Password      string `json:"password" binding:"required"`
		RememberToken string `json:"rememberToken"`
	}

	r := req{}
	if err := ctx.ShouldBindJSON(&r); err != nil {
		errorMsg := util.ValidatorError(err)
		ctx.JSON(http.StatusOK, util.FailedRespPackage(errorMsg))
		return
	}

	adminUser := new(models.AdminUser)
	if err := drivers.Mysql().Model(&models.AdminUser{}).Where("name = ?", r.UserName).First(adminUser).Error; err != nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage("用户名不存在"))
		return
	}
	if adminUser.Password != adminUser.EncryptionPw(r.Password) {
		ctx.JSON(http.StatusOK, util.FailedRespPackage("密码错误"))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessRespPackage(adminUser))
	return
}

func Test(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{"test": "test"}))
	return
}
