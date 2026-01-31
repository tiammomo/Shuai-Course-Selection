package router

import (
	"github.com/gin-gonic/gin"

	"course_select/internal/interface/api/handler"
	"course_select/internal/interface/api/middleware"
)

// Router 路由配置
type Router struct {
	authHandler     *handler.AuthHandler
	memberHandler   *handler.MemberHandler
	courseHandler   *handler.CourseHandler
	authMiddleware  *middleware.AuthMiddleware
	limiterMiddleware *middleware.LimiterMiddleware
}

// NewRouter 创建路由
func NewRouter(
	authHandler *handler.AuthHandler,
	memberHandler *handler.MemberHandler,
	courseHandler *handler.CourseHandler,
	authMiddleware *middleware.AuthMiddleware,
	limiterMiddleware *middleware.LimiterMiddleware,
) *Router {
	return &Router{
		authHandler:      authHandler,
		memberHandler:    memberHandler,
		courseHandler:    courseHandler,
		authMiddleware:   authMiddleware,
		limiterMiddleware: limiterMiddleware,
	}
}

// setupAuthRoutes 注册认证路由
func (r *Router) setupAuthRoutes(v1 *gin.RouterGroup) {
	auth := v1.Group("/auth")
	auth.POST("/login", r.authHandler.Login)
	auth.POST("/logout", r.authHandler.Logout)
	auth.GET("/whoami", r.authMiddleware.RequireAuth(), r.authHandler.WhoAmI)
}

// setupMemberRoutes 注册成员管理路由
func (r *Router) setupMemberRoutes(v1 *gin.RouterGroup) {
	member := v1.Group("/member")
	member.GET("", r.memberHandler.GetMember)
	member.GET("/list", r.memberHandler.GetMemberList)
	member.POST("/create", r.authMiddleware.RequireAuth(), r.authMiddleware.RequireAdmin(), r.memberHandler.CreateMember)
	member.POST("/update", r.authMiddleware.RequireAuth(), r.authMiddleware.RequireAdmin(), r.memberHandler.UpdateMember)
	member.POST("/delete", r.authMiddleware.RequireAuth(), r.authMiddleware.RequireAdmin(), r.memberHandler.DeleteMember)
}

// setupCourseRoutes 注册课程管理路由
func (r *Router) setupCourseRoutes(v1 *gin.RouterGroup) {
	course := v1.Group("/course")
	course.GET("/get", r.courseHandler.GetCourse)
	course.POST("/create", r.authMiddleware.RequireAuth(), r.authMiddleware.RequireAdmin(), r.courseHandler.CreateCourse)
	course.POST("/schedule", r.authMiddleware.RequireAuth(), r.authMiddleware.RequireAdmin(), r.courseHandler.ScheduleCourse)
}

// setupTeacherRoutes 注册教师管理路由
func (r *Router) setupTeacherRoutes(v1 *gin.RouterGroup) {
	teacher := v1.Group("/teacher")
	teacher.GET("/get_course", r.courseHandler.GetTeacherCourses)
	teacher.POST("/bind_course", r.authMiddleware.RequireAuth(), r.authMiddleware.RequireAdmin(), r.courseHandler.BindCourse)
	teacher.POST("/unbind_course", r.authMiddleware.RequireAuth(), r.authMiddleware.RequireAdmin(), r.courseHandler.UnbindCourse)
}

// setupStudentRoutes 注册学生选课路由
func (r *Router) setupStudentRoutes(v1 *gin.RouterGroup) {
	student := v1.Group("/student")
	student.POST("/book_course", r.courseHandler.BookCourse)
	student.GET("/course", r.courseHandler.GetStudentCourses)
}

// RegisterRoutes 注册路由
func (r *Router) RegisterRoutes(engine *gin.Engine) {
	engine.Use(r.limiterMiddleware.GlobalLimit(4000))
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	v1 := engine.Group("/api/v1")
	r.setupAuthRoutes(v1)
	r.setupMemberRoutes(v1)
	r.setupCourseRoutes(v1)
	r.setupTeacherRoutes(v1)
	r.setupStudentRoutes(v1)

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
