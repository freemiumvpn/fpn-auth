package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Bind attaches engine to handlers
func Bind(engine *gin.Engine) {
	engine.Use(Errors)
	engine.Use(Recovery)
	logrus.Trace("[ Client ] Binding Middlewares")
}
