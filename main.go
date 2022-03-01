package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func migrate(db *gorm.DB) {
	db.AutoMigrate(&Todo{})
}

func main() {
	if env := godotenv.Load(); env != nil {
		panic(env)
	}

	database := OpenDbConnection()
	defer database.Close()
	migrate(database)
	Seed(database)

	goGonicEngine := gin.Default()
	goGonicEngine.Use(cors.Default())

	apiGroup := goGonicEngine.Group("/api")
	apiGroup.GET("/todos", GetAllTodos)
	apiGroup.GET("/todos/completed", GetAllPendingTodos)
	apiGroup.GET("/todos/pending", GetAllCompletedTodos)
	apiGroup.GET("/todos/:id", GetTodoById)
	apiGroup.POST("/todos", CreateTodo)
	apiGroup.PUT("/todos/:id", UpdateTodo)
	apiGroup.PATCH("/todos/:id", CreateTodo)
	apiGroup.DELETE("/todos", DeleteAllTodos)
	apiGroup.DELETE("/todos/:id", DeleteTodo)

	_ = goGonicEngine.Run(":8080")
}
