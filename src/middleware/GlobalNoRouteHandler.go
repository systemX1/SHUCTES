package middleware

import (
	. "SHUCTES/src/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GlobalNoRouteHandler()  gin.HandlerFunc {
	return func(c *gin.Context) {
		Logger.WithField("url", c.Request.RequestURI).Infof("Page Not Found")
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "Page not found",
		})
	}
}

