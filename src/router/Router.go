package router

import (
	"SHUCTES/src/controller"
	"SHUCTES/src/middleware"

	"github.com/gin-gonic/gin"
)

//signin注册， login登录, logout登出
func RegisterRouter(app *gin.Engine) {
	//jwt中间件
	authMiddleware := middleware.GetGinJWTHandler()
	app.Use(middleware.CroHandler())			//跨域
	//404中间件
	app.NoRoute(middleware.GlobalNoRouteHandler())
	//记录请求中间件
	app.Use(middleware.LoggerHandler())
	//服务重启中间件
	app.Use(gin.Recovery())

	r := app.Group("/")
	r.POST("/signin", controller.SigninHandler())
	r.POST("/login", authMiddleware.LoginHandler)

	user := r.Group("/user", authMiddleware.MiddlewareFunc())
	user.GET("/refresh", authMiddleware.RefreshHandler) //刷新token
	user.POST("/testAuth", authMiddleware.MiddlewareFunc(), controller.TestAuthMethod()) //测试
	user.GET("/GetUserInfo", controller.GetUserInfoHandler()) //获取username对应user信息

	course := r.Group("/course")
	course.GET("/GetCourseInfo", controller.GetCourseInfoHandler()) 				//根据course_uid获取课程信息和评星、标签、评论
	course.GET("/SearchWithFulltext", controller.SearchWithFulltextHandler()) 	//TODO 全文搜索
	course.GET("/SearchWithCourseName", controller.SearchWithCourseNameHandler())//根据课程名搜索课程信息和平均评分
	course.GET("/SearchWithCID", controller.SearchWithCIDHandler()) 				//根据课程号搜索课程信息和平均评分
	course.GET("/SearchWithTeachname", controller.SearchWithTeachnameHandler()) 	//根据教师姓名搜索课程信息和平均评分
	course.GET("/SearchWithTeachid", controller.SearchWithTeachidHandler()) 		//根据教师工号搜索课程信息和平均评分
	course.GET("/SearchWithTimeText", controller.SearchWithTimeTextHandler()) 	//根据上课时间搜索课程信息和平均评分
	course.GET("/SearchWithSchool", controller.SearchWithSchoolHandler()) 		//根据开课学院搜索课程信息和平均评分

	star := r.Group("/star")
	star.GET("/GetUserStarredCourse", authMiddleware.MiddlewareFunc(), controller.GetUserStarredCourseHandler()) 	//根据username查打星过的课程和星级
	star.POST("/AddCourseStar", authMiddleware.MiddlewareFunc(), controller.AddCourseStarHandler()) 					//根据username, course_uid增加课程评星
	star.PUT("/UpdateCourseStar", authMiddleware.MiddlewareFunc(), controller.UpdateCourseStarHandler()) 			//根据username, course_uid更改课程评星
	star.DELETE("/DeleteCourseStar", authMiddleware.MiddlewareFunc(), controller.DeleteCourseStarHandler()) 			//根据username, course_uid删除课程评星

	tag := r.Group("/tag")
	tag.GET("/GetUserTaggedCourse", authMiddleware.MiddlewareFunc(), controller.GetUserTaggedCourseHandler()) //根据username查标签过的课程和标签
	tag.POST("/AddCourseTag", authMiddleware.MiddlewareFunc(), controller.AddCourseTagHandler())              //根据username, course_uid增加课程标签
	tag.DELETE("/DeleteCourseTag", authMiddleware.MiddlewareFunc(), controller.DeleteCourseTagHandler())      //删除username课程标签

	comment := r.Group("/comment")
	comment.GET("/GetUserCommentedCourse", authMiddleware.MiddlewareFunc(), controller.GetUserCommentedCourseHandler()) //根据username查评论过的课程和评论
	comment.POST("/AddCourseComment", authMiddleware.MiddlewareFunc(), controller.AddCourseCommentHandler())            //根据username, course_uid设置课程评论
	comment.PUT("/UpdateCourseComment", authMiddleware.MiddlewareFunc(), controller.UpdateCourseCommentHandler())       //根据username, course_uid更改课程评论
	comment.DELETE("/DeleteCourseComment", authMiddleware.MiddlewareFunc(), controller.DeleteCourseCommentHandler())    //删除username课程评论

	test := r.Group("/test")
	test.GET("/get", controller.TestGetMethod())
	test.PUT("/put", controller.TestPutMethod())
	test.POST("/post", controller.TestPostMethod())
	test.DELETE("/delete", controller.TestDeleteMethod())

	test.POST("/signin", controller.SigninHandler())
}
