package controller

import (
	. "SHUCTES/src/database"
	. "SHUCTES/src/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserStarredCourseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, ok := c.GetQuery("username")
		if !ok {
			c.JSON(http.StatusBadRequest,gin.H{"msg":     "Query wrong"})
			return
		}
		Logger.Infof("Handling Requset, Username: %s", username)

		sql := `
SELECT course.uid, course.cid, course.name, course.teachno, course.teachname, star.star_n
FROM ctes.star, ctes.course
WHERE star.course_uid = course.uid
AND username = ?;`

		rows, err := DB.Query(sql, username)
		if err != nil {
			Logger.Errorf("While querying: " + err.Error())
			c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
			return
		}

		type AutoGenerated struct {
			CourseUid  string `json:"course_uid"`
			CourseCid  string `json:"course_cid"`
			CourseName string `json:"course"`
			TeachNo    string `json:"teachno"`
			TeachName  string `json:"teachname"`
			StarN      int    `json:"star_n"`
		}
		var jsonSent AutoGenerated
		retSlice :=  make([]AutoGenerated, 0)

		for rows.Next() {
			if err := rows.Scan(&jsonSent.CourseUid, &jsonSent.CourseCid, &jsonSent.CourseName, &jsonSent.TeachNo, &jsonSent.TeachName, &jsonSent.StarN); err != nil {
				Logger.Errorf("While scanning rows: " + err.Error())
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
				return
			} else {
				retSlice = append(retSlice, jsonSent)
			}
		}
		rows.Close()

		c.JSON(http.StatusOK, gin.H{
			"content":	retSlice,
			"length": 	len(retSlice),
			"msg":     "Request successful",
		})
	}
}

func AddCourseStarHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, ok := c.GetQuery("username")
		if !ok {
			c.JSON(http.StatusBadRequest,gin.H{"msg":     "Query wrong"})
			return
		}

		type AutoGenerated struct {
			CourseUid	string	`json:"course_uid"`
			StarN 		int 	`json:"star_n"`
		}
		var jsonReci AutoGenerated

		if err := c.BindJSON(&jsonReci); err != nil {
			Logger.Errorf("While binding json to struct: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Query wrong"})
			return
		} else {
			StarNValidation(&jsonReci.StarN)
			Logger.Infof("Handling Requset, Username: %s, JsonReci = %#v", username, jsonReci)

			sql := `
			INSERT INTO ctes.star(star.username, star.course_uid, star.star_n) 
			SELECT ?, ?, ?;`

			stmt, err := DB.Prepare(sql)
			if err != nil {
				Logger.Errorf("While preparing sql: " + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "server error"})
				return
			}
			res, err := stmt.Exec(username, jsonReci.CourseUid, jsonReci.StarN)
			if err != nil {
				Logger.Errorf("While executing sql: " + err.Error())
				c.JSON(http.StatusForbidden, gin.H{"msg": "Duplicate key or other error while querying"})
				return
			}
			lastId, err := res.LastInsertId()
			if err != nil {
				Logger.Errorf("While getting sql result: " + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "server error"})
				return
			}
			rowCnt, err := res.RowsAffected()
			if err != nil {
				Logger.Errorf("While getting sql rowsAffected: " + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "server error"})
				return
			}
			Logger.Infof("Username: %s, JsonReci = %#v, ID = %d, affected = %d\n", username, jsonReci ,lastId, rowCnt)
			c.JSON(http.StatusOK, gin.H{"msg": "request successful"})
		}
	}
}

func UpdateCourseStarHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, ok := c.GetQuery("username")
		if !ok {
			c.JSON(http.StatusBadRequest,gin.H{"msg":     "Query wrong"})
			return
		}

		type AutoGenerated struct {
			CourseUid	string	`json:"course_uid"`
			StarN 		int 	`json:"star_n"`
		}
		var jsonReci AutoGenerated

		if err := c.BindJSON(&jsonReci); err != nil {
			Logger.Errorf("While binding json to struct: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Query wrong"})
			return
		} else {
			StarNValidation(&jsonReci.StarN)
			Logger.Infof("Handling Requset, Username: %s, JsonReci = %#v", username, jsonReci)

			sql := `
			UPDATE ctes.star
			SET star.star_n = ? 
			WHERE star.course_uid = ?
			AND star.username = ?;`
			result, err := DB.Exec(sql, jsonReci.StarN, jsonReci.CourseUid, username)

			if err != nil{
				Logger.Errorf("While executing sql: " + err.Error())
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
				return
			} else {
				Logger.Infof("Update success, Username: %s, JsonReci = %#v", username, jsonReci)
			}

			rowaffected, err := result.RowsAffected()
			if err != nil {
				Logger.Errorf("Get RowsAffected failed: " + err.Error())
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
			} else {
				Logger.Infof("Affected rows: %d", rowaffected)
				c.JSON(http.StatusOK,  gin.H{"msg": "request successful"})
			}
		}
	}
}

func DeleteCourseStarHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, ok := c.GetQuery("username")
		if !ok {
			c.JSON(http.StatusBadRequest,gin.H{"msg":     "Query wrong"})
			return
		}

		type AutoGenerated struct {
			CourseUid	string	`json:"course_uid"`
		}
		var jsonReci AutoGenerated

		if err := c.BindJSON(&jsonReci); err != nil {
			Logger.Errorf("While binding json to struct: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Query wrong"})
			return
		} else {
			Logger.Infof("Handling Requset, Username: %s, JsonReci = %#v", username, jsonReci)

			sql := `
			DELETE FROM ctes.star
			WHERE star.username = ?
			AND star.course_uid = ?;`
			result, err := DB.Exec(sql, username, jsonReci.CourseUid)

			if err != nil{
				Logger.Errorf("While executing sql: " + err.Error())
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
				return
			} else {
				Logger.Infof("Update success, Username: %s, JsonReci = %#v", username, jsonReci)
			}

			rowaffected, err := result.RowsAffected()
			if err != nil {
				Logger.Errorf("Get RowsAffected failed: " + err.Error())
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
			} else {
				Logger.Infof("Affected rows: %d", rowaffected)
				c.JSON(http.StatusOK,  gin.H{"msg": "request successful"})
			}
		}
	}
}

func StarNValidation(starN *int) {
	if *starN < 1 {
		*starN = 1
	} else if *starN > 5{
		*starN = 5
	}
}
