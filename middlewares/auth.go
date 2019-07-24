package middlewares

import (
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type User struct {
	UserName  string
	FirstName string
	LastName  string
}

const identityKey = "id"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := jwt.New(&jwt.GinJWTMiddleware{
			Realm:       "test zone",
			Key:         []byte("secret key"),
			Timeout:     time.Second * 300,
			MaxRefresh:  time.Second * 300,
			IdentityKey: identityKey,
			PayloadFunc: func(data interface{}) jwt.MapClaims {
				if v, ok := data.(*User); ok {
					return jwt.MapClaims{
						identityKey: v.UserName,
					}
				}
				return jwt.MapClaims{}
			},
			IdentityHandler: func(c *gin.Context) interface{} {
				claims := jwt.ExtractClaims(c)
				return &User{
					UserName: claims[identityKey].(string),
				}
			},
			Authenticator: func(c *gin.Context) (interface{}, error) {
				var loginVals Login
				if err := c.ShouldBind(&loginVals); err != nil {
					return "", jwt.ErrMissingLoginValues
				}
				userID := loginVals.Username
				password := loginVals.Password

				if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
					return &User{
						UserName:  userID,
						LastName:  "Nikita",
						FirstName: "Leskin",
					}, nil
				}

				return nil, jwt.ErrFailedAuthentication
			},
			Authorizator: func(data interface{}, c *gin.Context) bool {
				if v, ok := data.(*User); ok && v.UserName == "admin" {
					return true
				}

				return false
			},
			Unauthorized: func(c *gin.Context, code int, message string) {
				c.JSON(code, gin.H{
					"code":    code,
					"message": message,
				})
			},
			TokenLookup:   "header:Authorization",
			TokenHeadName: "Bearer",
			TimeFunc:      time.Now,
		})

		if err != nil {
			log.Fatal("JWT Error:" + err.Error())
		}

		c.Next()
	}
}
