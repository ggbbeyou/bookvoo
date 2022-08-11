package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/bookvoo/common"
	"github.com/yzimhao/bookvoo/user"
)

func user_query(c *gin.Context) {
	uinfo, _ := c.Get("user")
	common.Success(c, uinfo.(*user.User))
}
