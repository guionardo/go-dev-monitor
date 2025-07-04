package static

import (
	_ "embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed script.js
var script_js []byte

//go:embed index.html
var index_html []byte

//go:embed styles.css
var styles_css []byte

func ApplyStatic(router *gin.Engine) {
	router.GET("/", handler(index_html, "text/html"))
	router.GET("/script.js", handler(script_js, "application/javascript"))
	router.GET("/styles.css", handler(styles_css, "text/css"))
}

func handler(data []byte, contentType string) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, contentType, data)
	}
}
