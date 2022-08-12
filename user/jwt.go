package user

import (
	"fmt"
	"net/http"
	"time"

	"math/rand"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yzimhao/bookvoo/common"
	"github.com/yzimhao/bookvoo/user/assets"
)

var AuthMiddleware *jwt.GinJWTMiddleware

type req_login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func InitJwt() {
	viper.SetDefault("main.jwt_key", "e94dae72b1e14876")

	auth, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte(viper.GetString("main.jwt_key")),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,

		IdentityKey: "user",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					"user_id":  v.UserId,
					"username": v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			uid := int64(claims["user_id"].(float64))
			return &User{
				UserId:   uid,
				UserName: claims["username"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			// var loginVals login
			// if err := c.ShouldBind(&loginVals); err != nil {
			// 	return "", jwt.ErrMissingLoginValues
			// }
			// userID := "admin"
			// password := "admin"

			// if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
			// 	return &User{
			// 		UserName:  userID,
			// 		LastName:  "Bo-Yi",
			// 		FirstName: "Wu",
			// 	}, nil
			// }

			if viper.GetString("main.mode") == "demo" {
				rand.Seed(time.Now().Unix())
				n := rand.Intn(99)
				demoUid := DemoUserStart + int64(n)

				assets.InitAssetsForDemo(demoUid, DemoUsdSymbol, "100000", "R001")
				assets.InitAssetsForDemo(demoUid, DemoEthSymbol, "100000", "R001")

				return &User{
					UserId: demoUid,
					UserName: func() string {
						return fmt.Sprintf("user%d", demoUid)
					}(),
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		//权限认证
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// if v, ok := data.(*User); ok && v.UserName == "admin" {
			// 	return true
			// }

			// return false
			return true
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(http.StatusOK, gin.H{
				"ok":     0,
				"reason": message,
			})
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {

			exp := expire.Unix() - time.Now().Unix()
			c.SetCookie("jwt", token, int(exp), "/", "*", false, false)

			c.JSON(http.StatusOK, gin.H{
				"ok": 1,
				"data": map[string]string{
					"token":  token,
					"expire": expire.Format(time.RFC3339),
				},
			})
		},

		LogoutResponse: func(c *gin.Context, code int) {
			c.SetCookie("jwt", "", -1, "/", "*", false, false)
			common.Success(c, nil)
		},

		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		logrus.Error("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := auth.MiddlewareInit()
	if errInit != nil {
		logrus.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}
	AuthMiddleware = auth
}
