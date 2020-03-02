package filter

import (
	"github.com/gin-gonic/gin"
)

//拦截器
func HTTPInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		api := c.Request.URL.Path
		identity := ApiClassification(api)
		var (
			code   int
			msg    string
			claims *CustomClaims
		)
		switch identity {
		case "admin":
			code, msg, claims = AdminFilter(c, api)
		case "user":
			code, msg, claims = UserFilter(c, api)
		}

		if code != 0 {
			c.JSON(code, msg)
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", claims)
	}
}
