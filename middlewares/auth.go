package middlewares

import (
	"errors"
	"persona/forms"
	"persona/models"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const (
	identityKey = "id"
	realm       = "Persona App"
	// TODO: move key to a config
	secretKey                  = "]@jvVX,hAiwT|%B`\\d@PYsGb=)74#H|BNb[F^Iu<u0OUn:]AX/BP%EYy:?EJ|Hq"
	tokenExpireTimeoutSeconds  = 300
	tokenRefreshTimeoutSeconds = 300
)

func AuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	middleware := &jwt.GinJWTMiddleware{
		Realm:            realm,
		SigningAlgorithm: "HS256",
		Key:              []byte(secretKey),
		Timeout:          tokenExpireTimeoutSeconds * time.Second,
		MaxRefresh:       tokenRefreshTimeoutSeconds * time.Second,
		Authenticator:    authenticate,
		Authorizator:     authorize,
		PayloadFunc:      payload,
		Unauthorized:     unauthorized,
		IdentityHandler:  identityHandler,
		TokenLookup:      "header:Authorization",
		TokenHeadName:    "Bearer",
		TimeFunc:         time.Now,
	}

	if err := middleware.MiddlewareInit(); err != nil {
		return nil, err
	}

	return middleware, nil
}

func authenticate(c *gin.Context) (interface{}, error) {
	var loginValues forms.UserLogin
	if err := c.ShouldBindJSON(&loginValues); err != nil {
		return nil, errors.New("Incorrect email or password")
	}

	email := loginValues.Email
	password := loginValues.Password

	if email == "admin" && password == "admin" {
		return &models.User{}, nil
	}

	return nil, jwt.ErrFailedAuthentication
}

func authorize(data interface{}, c *gin.Context) bool {
	if _, ok := data.(*models.User); ok {
		return true
	}

	return true
}

func payload(data interface{}) jwt.MapClaims {
	// TODO: add data to jwt token
	// if v, ok := data.(*User); ok {
	// 	return jwt.MapClaims{
	// 		identityKey: v.UserName,
	// 	}
	// }
	return jwt.MapClaims{}
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func identityHandler(*gin.Context) interface{} {
	// claims := jwt.ExtractClaims(c)
	// return &User{
	// 	UserName: claims[identityKey].(string),
	// }
	return nil
}
