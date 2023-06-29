package main

import (
	"github.com/Cedar-81/swype/controllers"
	"github.com/Cedar-81/swype/initializers"
	"github.com/Cedar-81/swype/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDatabase()
	initializers.CreateSuperAdmin()
}

func main() {

	r := gin.Default()

	//handle post routes
	r.GET("/posts", middleware.RequireAuth, controllers.GetPosts)
	r.GET("/post/:id", middleware.RequireAuth, controllers.GetPost)
	r.POST("/post", middleware.RequireAuth, controllers.CreatePost)
	r.PATCH("/post/:id", middleware.RequireAuth, controllers.UpdatePost)
	r.DELETE("/post/:id", controllers.DeletePost)
	//handle post routes end

	//handle comment
	r.GET("/comments", middleware.RequireAuth, middleware.CheckAdmin, controllers.GetComments)
	r.GET("/comments/:postid", middleware.RequireAuth, controllers.GetCommentsByPostId)
	r.GET("/comment/:id", middleware.RequireAuth, controllers.GetComment)
	r.POST("/comment/:postid", middleware.RequireAuth, controllers.CreateComment)
	r.PATCH("/comment/:id", middleware.RequireAuth, controllers.UpdateComment)
	r.DELETE("/comment/:id", middleware.RequireAuth, controllers.DeleteComment)
	//handle comment routes end

	//handle auth routes
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.LogIn)
	r.GET("/logout", controllers.LogOut)
	// Super Admin create Admin
	r.POST("/createadmin", middleware.RequireAuth, middleware.CheckSuperAdmin, controllers.CreateAdmin)
	// check login or not
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	// Chưa laàm ádnl ándkl
	r.GET("/resetpassword", controllers.ResetPassword)
	r.PATCH("/updatepassword", middleware.RequireAuth, controllers.UpdatePassword)
	//handle auth routes end

	//handle users routes
	// only Admin or Super Admin can get all users
	r.GET("/users", middleware.RequireAuth, middleware.CheckAdmin, controllers.GetUsers)
	r.GET("/user", middleware.RequireAuth, controllers.GetLoggedUser)
	r.PATCH("/user", middleware.RequireAuth, controllers.UpdateUser)
	r.GET("/user/:username", middleware.RequireAuth, controllers.GetUser)
	r.DELETE("/user", middleware.RequireAuth, controllers.DeleteUser)
	//handle users end

	//handle special routes
	r.PATCH("/follow/:tofollow", middleware.RequireAuth, controllers.UpdateUser)
	r.GET("/feed", middleware.RequireAuth, controllers.GetFeedPosts)
	r.GET("/search/:searchparam", middleware.RequireAuth, controllers.Search)
	r.GET("/notifications", middleware.RequireAuth, controllers.GetNotifications)
	r.GET("/like/post/:postid", middleware.RequireAuth, controllers.LikePost)
	r.GET("/like/comment/:commentid", middleware.RequireAuth, controllers.LikeComment)
	//end handle special routes

	r.Run() // listen and serve on 0.0.0.0:3000
}
