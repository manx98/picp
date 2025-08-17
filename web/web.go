package web

import (
	"embed"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"net/http"
	"picp/logger"
)

//go:embed dist
var sourceFs embed.FS

func Init(router *gin.Engine, isLogin func(*gin.Context) bool, logout func(*gin.Context)) {
	router.GET("/logout", func(ctx *gin.Context) {
		logout(ctx)
		ctx.Redirect(302, "/login")
	})
	distFs, err := fs.Sub(sourceFs, "dist")
	if err != nil {
		logger.Fatal("read embed file occur error", zap.Error(err))
	}
	httpFs := http.FS(distFs)
	err = fs.WalkDir(distFs, ".", func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			if path == "index.html" {
				file, err := distFs.Open(path)
				if err != nil {
					return err
				}
				indexContent, err := io.ReadAll(file)
				if err != nil {
					return err
				}
				indexHtml := string(indexContent)
				warp := func(ctx *gin.Context) {
					ctx.Header("content-type", "text/html;charset=utf-8")
					ctx.String(http.StatusOK, indexHtml)
				}
				index := func(ctx *gin.Context) {
					if isLogin(ctx) {
						warp(ctx)
					} else {
						ctx.Redirect(302, "/login")
					}
				}
				router.GET("", index)
				router.GET("/wifi", index)
				router.GET("/setting", index)
				router.GET("/login", func(ctx *gin.Context) {
					if isLogin(ctx) {
						ctx.Redirect(302, "")
					} else {
						warp(ctx)
					}
				})
			} else {
				router.GET(path, func(ctx *gin.Context) {
					ctx.Header("Cache-Control", "public, max-age=86400")
					ctx.FileFromFS(path, httpFs)
				})
			}
		}
		return err
	})
	if err != nil {
		logger.Fatal("init web resource failed", zap.Error(err))
	}
}
