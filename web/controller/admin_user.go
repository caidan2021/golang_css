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
		UserName   string `json:"username" binding:"required"`
		Password   string `json:"password" binding:"required"`
		RememberMe bool   `json:"rememberMe" default:"true"`
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

	if err := adminUser.UpdateRememberMe(r.RememberMe); err != nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage(err.Error()))
		return
	}
	ctx.SetCookie("css-token", adminUser.RememberToken, 86400*14, "/", util.Config.Demon, false, false)

	rt := models.AdminRender{
		ID:    adminUser.ID,
		Name:  adminUser.Name,
		Email: adminUser.Email,
		Token: adminUser.RememberToken,
	}

	ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{"currentUser": rt}))
	return
}

func Current(ctx *gin.Context) {
	admin, ok := ctx.Get(models.AdminUserKey)
	if !ok {
		ctx.JSON(http.StatusOK, util.SuccessRespPackage(true))
	}

	ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{"currentUser": admin}))

}
