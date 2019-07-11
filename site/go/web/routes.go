package web

import (
	"github.com/gin-gonic/gin"
)

func RegisterPages(router* gin.Engine) {
	router.LoadHTMLGlob("./templates/*")

	router.Static("/css", "./css")
	router.Static("/img", "./img")
	router.Static("/ts", "./ts")
	router.StaticFile("/Bundle.js", "./Bundle.js")

	router.GET("/", IndexHandler)
}