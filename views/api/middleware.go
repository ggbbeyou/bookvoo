package api

import "github.com/gin-gonic/gin"

func login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//todo 登陆中间件
		ctx.Set("user_id", USERID)
	}
}
