package middleware

import (
	. "SHUCTES/src/database"
	. "SHUCTES/src/log"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type User struct {
	Username string `form:"username" json:"username"`
}

var identityKey = "username"

//Authentication 认证 Authorization 授权
func GetGinJWTHandler() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       	"gin-jwt zone",
		Key:         	[]byte("secret key"),
		Timeout:     	6 * time.Hour,
		MaxRefresh:  	6 * time.Hour,
		IdentityKey: 	identityKey,
		Authenticator:	authenticator,
		Authorizator:	authorizator,
		Unauthorized:	unauthorized,
		PayloadFunc:    payload,
		IdentityHandler:identityHandler,
	})

	if err != nil {
		Logger.Fatal("JWT Error:" + err.Error())
	}
	return authMiddleware
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &User{
		Username: claims[identityKey].(string),
	}
}

func payload(data interface{}) jwt.MapClaims {
	if v, ok := data.(*User); ok {
		return jwt.MapClaims{
			identityKey: v.Username,
		}
	}
	return jwt.MapClaims{}
}

func authenticator(c *gin.Context) (interface{}, error) {
	var loginVal Login

	if err := c.BindJSON(&loginVal); err != nil {
		Logger.Errorf("While binding json: %s", err)
	} else {
		//从user表读username对应密码
		querySQL := `
			SELECT ctes.user.password
			FROM ctes.user
			WHERE ctes.user.name = ?;`

		row := DB.QueryRow(querySQL, loginVal.Username)

		var hashedPW string
		if err := row.Scan(&hashedPW); err != nil{
			Logger.Infof("While reading password from database %s", err)
		} else {
			//验证数据库中hashedPW与表单中Password
			if err = bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(loginVal.Password)); err != nil {
				Logger.Infof("user %s log in failed", loginVal.Username)
			} else {
				Logger.Infof("user %s log in successful", loginVal.Username)
				return &User{Username: loginVal.Username}, nil
			}
		}
	}

	return nil, jwt.ErrFailedAuthentication
}

func authorizator(data interface{}, c *gin.Context) bool {
	query, _ := c.GetQuery("username")
	if v, ok := data.(*User); ok && v.Username == query {
		return true
	}
	return false
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    	code,
		"msg": 		message,
	})
}







