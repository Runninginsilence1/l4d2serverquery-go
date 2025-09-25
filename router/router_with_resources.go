//go:build embed_ui

package router

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"l4d2serverquery-go/resources"
)

func setUI(r *gin.Engine) {

	// 前端是单页面不存在路由跳转
	fs, err := static.EmbedFolder(resources.UiFs, "dist")
	if err != nil {
		panic(err)
	}
	r.Use(static.Serve("/", fs))
}
