package controller

import (
	. "SHUCTES/src/database"
	. "SHUCTES/src/log"
	"SHUCTES/src/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Signin struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func SigninHandler() gin.HandlerFunc  {
	return func(c *gin.Context) {
		var user model.User

		if err := c.BindJSON(&user); err != nil {
			Logger.Errorf("While binding json: %s", err)
		} else {
			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				Logger.Errorf("While GenerateFromPassword: %s", err)
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
				return
			}

			SQL := `
			INSERT INTO ctes.user (name, password, enrol_year, school, permission) VALUES
	        (?,	?, ?, ?, 1);`

			stmt, err := DB.Prepare(SQL)
			if err != nil {
				Logger.Errorf("While preparing SQL: " + err.Error())
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
				return
			}
			res, err := stmt.Exec(user.Username, hash, user.EnrolYear, user.School)
			if err != nil {
				Logger.Errorf("While executing SQL: " + err.Error())
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
				return
			}
			lastId, err := res.LastInsertId()
			if err != nil {
				Logger.Errorf("While getting SQL result: " + err.Error())
			}
			rowCnt, err := res.RowsAffected()
			if err != nil {
				Logger.Errorf("While getting SQL rowsAffected: " + err.Error())
			}
			Logger.Infof("ID = %d, affected = %d\n", lastId, rowCnt)

			c.JSON(http.StatusOK, gin.H{
				"username":	user.Username,
				"msg":		"Sign in successful",
			})
		}
	}
}

func GetUserInfoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if username, ok := c.GetQuery("username"); !ok {
			c.JSON(http.StatusBadRequest,"Query wrong")
			return
		} else {
			querySQL := `
			SELECT name, enrol_year, school, permission, create_time, update_time
			FROM ctes.user
			WHERE ctes.user.name = ?;`

			var user model.User
			row := DB.QueryRow(querySQL, username)
			if err := row.Scan(&user.Username, &user.EnrolYear, &user.School, &user.Permission, &user.CrTime, &user.UpTime); err != nil{
				Logger.Errorf("While scanning from database %s", err)
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
				return
			} else {
				c.JSON(http.StatusOK, gin.H{
					"username":			user.Username,
					"enrol_year":		user.EnrolYear,
					"school":			user.School,
					"permission":		user.Permission,
					"create_time":		user.CrTime,
					"update_time":		user.UpTime,
					"msg":     "Request successful",
				})

				Logger.Infof("user %s GetUserInfo successful", username)
			}
		}
	}
}

func TestAuthMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, ok := c.GetQuery("username")
		if !ok {
			c.JSON(http.StatusBadRequest,gin.H{"msg":     "Query wrong"})
			return
		}

		c.JSON(200, gin.H{
			"username":	username,
			"msg":     "auth successful",
		})
	}
}
