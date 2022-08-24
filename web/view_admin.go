/*
 * @Date: 2022-08-20 17:25:36
 */
package web

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Handler(ctx *gin.Context) {
	pageName := ctx.Param("pageName")
	if !strings.HasPrefix(pageName, "admin_") && !strings.HasPrefix(pageName, "tmpl") {
		pageName = "admin_" + pageName
	}
	if !strings.HasSuffix(pageName, ".html") {
		pageName += ".html"
	}
	ctx.HTML(http.StatusOK, pageName, gin.H{})

}
