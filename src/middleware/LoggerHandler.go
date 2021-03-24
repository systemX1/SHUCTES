package middleware

import (
	. "SHUCTES/src/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func LoggerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//计算执行时间
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		Logger.WithFields(logrus.Fields{
			"statusCode": c.Writer.Status(),
			"latencyTime": latencyTime,
			"reqUri": c.Request.RequestURI,
			"clientIP": c.ClientIP(),
			"reqMethod": c.Request.Method,
		}).Infof("Handled")
	}
}
