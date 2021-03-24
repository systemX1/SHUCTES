package controller

import (
	"SHUCTES/src/config"
	. "SHUCTES/src/database"
	. "SHUCTES/src/log"
	"SHUCTES/src/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCourseInfoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		courseUid, ok := c.GetQuery("course_uid")
		if !ok {
			c.JSON(http.StatusBadRequest,gin.H{"msg":     "Query wrong"})
			return
		}
		Logger.Infof("Handling Requset, courseUid: %s", courseUid)

		//课程信息和star数
		sqlA := `
        SELECT course.uid, course.name, credit, cid, teachno, teachname, teachid, timetext, room, cap, peo_n, school, AVG(IF(star_n IS NULL, 0, star_n)) star_n
FROM ctes.course LEFT JOIN ctes.star 
ON course.uid = star.course_uid
WHERE course.uid = ?
GROUP BY course.uid, course.name, credit, cid, teachno, teachname, teachid, timetext, room, cap, peo_n, school
ORDER BY star_n DESC;  `

		rows, err := DB.Query(sqlA, courseUid)
		if err != nil {
			Logger.Errorf("While querying: " + err.Error())
			c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
			return
		}

		var jsonSentA model.Course
		corSlice := make([]model.Course, 0) //课程信息和star数

		for rows.Next() {
			if err := rows.Scan(&jsonSentA.Uid, &jsonSentA.Name, &jsonSentA.Credit, &jsonSentA.Cid, &jsonSentA.Teachno, &jsonSentA.Teachname, &jsonSentA.Teachid, &jsonSentA.Timetext, &jsonSentA.Room, &jsonSentA.Cap, &jsonSentA.PeoN, &jsonSentA.School, &jsonSentA.AvgStar); err != nil {
				Logger.Errorf("While scanning rows: " + err.Error())
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
				return
			} else {
				corSlice = append(corSlice, jsonSentA)
			}
		}
		rows.Close()

		//课程的tag_idx和tag数量
		sqlB := `
		SELECT tag_idx, COUNT(tag_idx) FROM ctes.tag
		WHERE course_uid = ?
		GROUP BY tag_idx
		ORDER BY COUNT(tag_idx) DESC;`

		rows, err = DB.Query(sqlB, courseUid)
		if err != nil {
			Logger.Errorf("While querying: " + err.Error())
			c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
			return
		}

		type tag struct {
			Idx int    `json:"idx"`
			Tag string `json:"tag"`
			Num int    `json:"num"`
		}
		var jsonSentB tag//课程的tag_idx,tag和tag数量
		tagSlice :=  make([]tag, 0)

		for rows.Next() {
			if err := rows.Scan(&jsonSentB.Idx, &jsonSentB.Num); err != nil {
				Logger.Errorf("While scanning rows: " + err.Error())
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
				return
			} else {
				tagSlice = append(tagSlice, jsonSentB)
			}
		}
		rows.Close()

		//为tag添加name string
		for i, t := range tagSlice {
			tagSlice[i].Tag = config.Conf.TagCon[t.Idx - 1]
		}


		//课程评论
		sqlC := `
		SELECT content FROM ctes.comment
		WHERE course_uid = ?;`

		rows, err = DB.Query(sqlC, courseUid)
		if err != nil {
			Logger.Errorf("While querying: " + err.Error())
			c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
			return
		}

		type comment struct {
			Content string `json:"content"`
		}
		var jsonSentC comment//课程的tag_idx,tag和tag数量
		commentSlice :=  make([]comment, 0)

		for rows.Next() {
			if err := rows.Scan(&jsonSentC.Content); err != nil {
				Logger.Errorf("While scanning rows: " + err.Error())
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
				return
			} else {
				commentSlice = append(commentSlice, jsonSentC)
			}
		}
		rows.Close()


		c.JSON(http.StatusOK, gin.H{
			"course": 	corSlice,
			"tag":  	tagSlice,
			"comment":	commentSlice,
			"msg":     "Request successful",
		})
	}
}

func SearchWithFulltextHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}


func SearchWithCourseNameHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if course, ok := c.GetQuery("course"); !ok {
			c.JSON(http.StatusBadRequest,gin.H{
				"msg":     "Query wrong",
			})
			return
		} else {
			Logger.Infof("Handling Requset, CourseName: %s", course)
			sql := `
		 	SELECT course.uid, course.name, credit, cid, teachno, teachname, teachid, timetext, room, cap, peo_n, school, AVG(IF(star_n IS NULL, 0, star_n)) star_n
FROM ctes.course LEFT JOIN ctes.star 
ON course.uid = star.course_uid
WHERE course.name LIKE ?
GROUP BY course.uid, course.name, credit, cid, teachno, teachname, teachid, timetext, room, cap, peo_n, school
ORDER BY star_n DESC;  `

			if rows, err := DB.Query(sql,  "%" + course + "%"); err != nil {
				Logger.Errorf("While querying: " + err.Error())
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
				return
			} else {
				var jsonSent model.Course
				ret := make([]model.Course, 0)
//teachno, teachname, teachid, timetext, room, cap, peo_n, school
				defer rows.Close()
				for rows.Next() {
					if err := rows.Scan(&jsonSent.Uid, &jsonSent.Name, &jsonSent.Credit, &jsonSent.Cid, &jsonSent.Teachno, &jsonSent.Teachname, &jsonSent.Teachid, &jsonSent.Timetext, &jsonSent.Room, &jsonSent.Cap, &jsonSent.PeoN, &jsonSent.School, &jsonSent.AvgStar); err != nil {
						Logger.Errorf("While scanning rows: " + err.Error())
						c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
						return
					} else {
						ret = append(ret, jsonSent)
					}
				}

				c.JSON(http.StatusOK, gin.H{
					"content":	ret,
					"length": 	len(ret),
					"msg":     "Request successful",
				})
			}
		}
	}
}

func SearchWithCIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cid, ok := c.GetQuery("cid"); !ok {
			c.JSON(http.StatusBadRequest,gin.H{
				"msg":     "Query wrong",
			})
			return
		} else {
			Logger.Infof("Handling Requset, cid: %s", cid)
			sql := `
SELECT course.uid, course.name, credit, cid, teachno, teachname, teachid, timetext, room, cap, peo_n, school, AVG(IF(star_n IS NULL, 0, star_n)) star_n
FROM ctes.course LEFT JOIN ctes.star 
ON course.uid = star.course_uid
WHERE course.cid LIKE ?
GROUP BY course.uid, course.name, credit, cid, teachno, teachname, teachid, timetext, room, cap, peo_n, school
ORDER BY star_n DESC;`

			if rows, err := DB.Query(sql,  "%" + cid + "%"); err != nil {
				Logger.Errorf("While querying: " + err.Error())
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
				return
			} else {
				var jsonSent model.Course
				ret := make([]model.Course, 0)
				//teachno, teachname, teachid, timetext, room, cap, peo_n, school
				defer rows.Close()
				for rows.Next() {
					if err := rows.Scan(&jsonSent.Uid, &jsonSent.Name, &jsonSent.Credit, &jsonSent.Cid, &jsonSent.Teachno, &jsonSent.Teachname, &jsonSent.Teachid, &jsonSent.Timetext, &jsonSent.Room, &jsonSent.Cap, &jsonSent.PeoN, &jsonSent.School, &jsonSent.AvgStar); err != nil {
						Logger.Errorf("While scanning rows: " + err.Error())
						c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
						return
					} else {
						ret = append(ret, jsonSent)
					}
				}

				c.JSON(http.StatusOK, gin.H{
					"content":	ret,
					"length": 	len(ret),
					"msg":     "Request successful",
				})
			}
		}
	}
}

func SearchWithTeachnameHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if teachname, ok := c.GetQuery("teachname"); !ok {
			c.JSON(http.StatusBadRequest,gin.H{
				"msg":     "Query wrong",
			})
			return
		} else {
			Logger.Infof("Handling Requset, teachname: %s", teachname)
			sql := `
SELECT course.uid, course.name, credit, cid, teachno, teachname, teachid, timetext, room, cap, peo_n, school, AVG(IF(star_n IS NULL, 0, star_n)) star_n
FROM ctes.course LEFT JOIN ctes.star 
ON course.uid = star.course_uid
WHERE course.teachname LIKE ?
GROUP BY course.uid, course.name, credit, cid, teachno, teachname, teachid, timetext, room, cap, peo_n, school
ORDER BY star_n DESC;`

			if rows, err := DB.Query(sql,  "%" + teachname + "%"); err != nil {
				Logger.Errorf("While querying: " + err.Error())
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
				return
			} else {
				var jsonSent model.Course
				ret := make([]model.Course, 0)
				//teachno, teachname, teachid, timetext, room, cap, peo_n, school
				defer rows.Close()
				for rows.Next() {
					if err := rows.Scan(&jsonSent.Uid, &jsonSent.Name, &jsonSent.Credit, &jsonSent.Cid, &jsonSent.Teachno, &jsonSent.Teachname, &jsonSent.Teachid, &jsonSent.Timetext, &jsonSent.Room, &jsonSent.Cap, &jsonSent.PeoN, &jsonSent.School, &jsonSent.AvgStar); err != nil {
						Logger.Errorf("While scanning rows: " + err.Error())
						c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
						return
					} else {
						ret = append(ret, jsonSent)
					}
				}

				c.JSON(http.StatusOK, gin.H{
					"content":	ret,
					"length": 	len(ret),
					"msg":     "Request successful",
				})
			}
		}
	}
}

func SearchWithTeachidHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if teachid, ok := c.GetQuery("teachid"); !ok {
			c.JSON(http.StatusBadRequest,gin.H{
				"msg":     "Query wrong",
			})
			return
		} else {
			Logger.Infof("Handling Requset, teachid: %s", teachid)
			sql := `
		 	SELECT course.uid, course.name, credit, cid, teachno, teachname, teachid, timetext, room, cap, peo_n, school, AVG(IF(star_n IS NULL, 0, star_n)) star_n
FROM ctes.course LEFT JOIN ctes.star 
ON course.uid = star.course_uid
WHERE course.teachid LIKE ?
GROUP BY course.uid, course.name, credit, cid, teachno, teachname, teachid, timetext, room, cap, peo_n, school
ORDER BY star_n DESC;`

			if rows, err := DB.Query(sql,  "%" + teachid + "%"); err != nil {
				Logger.Errorf("While querying: " + err.Error())
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
				return
			} else {
				var jsonSent model.Course
				ret := make([]model.Course, 0)
				//teachno, teachname, teachid, timetext, room, cap, peo_n, school
				defer rows.Close()
				for rows.Next() {
					if err := rows.Scan(&jsonSent.Uid, &jsonSent.Name, &jsonSent.Credit, &jsonSent.Cid, &jsonSent.Teachno, &jsonSent.Teachname, &jsonSent.Teachid, &jsonSent.Timetext, &jsonSent.Room, &jsonSent.Cap, &jsonSent.PeoN, &jsonSent.School, &jsonSent.AvgStar); err != nil {
						Logger.Errorf("While scanning rows: " + err.Error())
						c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
						return
					} else {
						ret = append(ret, jsonSent)
					}
				}

				c.JSON(http.StatusOK, gin.H{
					"content":	ret,
					"length": 	len(ret),
					"msg":     "Request successful",
				})
			}
		}
	}
}

func SearchWithTimeTextHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if timetext, ok := c.GetQuery("timetext"); !ok {
			c.JSON(http.StatusBadRequest,gin.H{
				"msg":     "Query wrong",
			})
			return
		} else {
			Logger.Infof("Handling Requset, timetext: %s", timetext)
			sql := `
		 	SELECT course.uid, course.name, credit, cid, teachno, teachname, teachid, timetext, room, cap, peo_n, school, AVG(IF(star_n IS NULL, 0, star_n)) star_n
FROM ctes.course LEFT JOIN ctes.star 
ON course.uid = star.course_uid
WHERE course.timetext LIKE ?
GROUP BY course.uid, course.name, credit, cid, teachno, teachname, teachid, timetext, room, cap, peo_n, school
ORDER BY star_n DESC;`

			if rows, err := DB.Query(sql,  "%" + timetext + "%"); err != nil {
				Logger.Errorf("While querying: " + err.Error())
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
				return
			} else {
				var jsonSent model.Course
				ret := make([]model.Course, 0)
				//teachno, teachname, teachid, timetext, room, cap, peo_n, school
				defer rows.Close()
				for rows.Next() {
					if err := rows.Scan(&jsonSent.Uid, &jsonSent.Name, &jsonSent.Credit, &jsonSent.Cid, &jsonSent.Teachno, &jsonSent.Teachname, &jsonSent.Teachid, &jsonSent.Timetext, &jsonSent.Room, &jsonSent.Cap, &jsonSent.PeoN, &jsonSent.School, &jsonSent.AvgStar); err != nil {
						Logger.Errorf("While scanning rows: " + err.Error())
						c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
						return
					} else {
						ret = append(ret, jsonSent)
					}
				}

				c.JSON(http.StatusOK, gin.H{
					"content":	ret,
					"length": 	len(ret),
					"msg":     "Request successful",
				})
			}
		}
	}
}

func SearchWithSchoolHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if school, ok := c.GetQuery("school"); !ok {
			c.JSON(http.StatusBadRequest,gin.H{
				"msg":     "Query wrong",
			})
			return
		} else {
			Logger.Infof("Handling Requset, school: %s", school)
			sql := `
		 	SELECT course.uid, course.name, credit, cid, teachno, teachname, teachid, timetext, room, cap, peo_n, school, AVG(IF(star_n IS NULL, 0, star_n)) star_n
FROM ctes.course LEFT JOIN ctes.star 
ON course.uid = star.course_uid
WHERE course.school LIKE ?
GROUP BY course.uid, course.name, credit, cid, teachno, teachname, teachid, timetext, room, cap, peo_n, school
ORDER BY star_n DESC;`

			if rows, err := DB.Query(sql,  "%" + school + "%"); err != nil {
				Logger.Errorf("While querying: " + err.Error())
				c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
				return
			} else {
				var jsonSent model.Course
				ret := make([]model.Course, 0)
				//teachno, teachname, teachid, timetext, room, cap, peo_n, school
				defer rows.Close()
				for rows.Next() {
					if err := rows.Scan(&jsonSent.Uid, &jsonSent.Name, &jsonSent.Credit, &jsonSent.Cid, &jsonSent.Teachno, &jsonSent.Teachname, &jsonSent.Teachid, &jsonSent.Timetext, &jsonSent.Room, &jsonSent.Cap, &jsonSent.PeoN, &jsonSent.School, &jsonSent.AvgStar); err != nil {
						Logger.Errorf("While scanning rows: " + err.Error())
						c.JSON(http.StatusInternalServerError,  gin.H{"msg": "server error"})
						return
					} else {
						ret = append(ret, jsonSent)
					}
				}

				c.JSON(http.StatusOK, gin.H{
					"content":	ret,
					"length": 	len(ret),
					"msg":     "Request successful",
				})
			}
		}
	}
}



















































