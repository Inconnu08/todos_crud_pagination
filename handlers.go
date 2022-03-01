package main

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"strconv"
)

func getPagingParams(c *gin.Context) (page, pageSize int) {
	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 5
	}

	page, err = strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	return
}

func GetAllTodos(c *gin.Context) {
	page, pageSize := getPagingParams(c)

	todos, totalTodoCount := FetchTodos(page, pageSize)

	c.JSON(http.StatusOK, CreateTodoPagedResponse(c.Request, todos, page, pageSize, totalTodoCount))
}

func GetAllPendingTodos(c *gin.Context) {
	page, pageSize := getPagingParams(c)
	todos, totalTodoCount := FetchPendingTodos(page, pageSize, false)

	c.JSON(http.StatusOK, CreateTodoPagedResponse(c.Request, todos, page, pageSize, totalTodoCount))
}

func GetAllCompletedTodos(c *gin.Context) {
	page, pageSize := getPagingParams(c)
	todos, totalTodoCount := FetchPendingTodos(page, pageSize, true)

	c.JSON(http.StatusOK, CreateTodoPagedResponse(c.Request, todos, page, pageSize, totalTodoCount))
}

func GetTodoById(c *gin.Context) {
	id := c.Param("id")
	if id == "completed" {
		GetAllCompletedTodos(c)
		return
	} else if id == "pending" {
		GetAllPendingTodos(c)
		return
	}
	id64, _ := strconv.ParseUint(id, 10, 32)
	todo, err := FetchById(uint(id64))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, CreateErrorDtoWithMessage("Could not find Todo"))
		return
	}

	c.JSON(http.StatusOK, GetSuccessTodoDto(&todo))
}

func CreateTodo(c *gin.Context) {
	var json CreateTodoResponse
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, CreateBadRequestErrorDto(err))
		return
	}
	todo, err := CreateTodoServices(json.Title, json.Description, json.Completed)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, CreateErrorDtoWithMessage(err.Error()))
	}

	c.JSON(http.StatusOK, CreateTodoCreatedDto(&todo))
}

func UpdateTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, CreateErrorDtoWithMessage("You must set an ID"))
		return
	}

	var json CreateTodoResponse
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, CreateBadRequestErrorDto(err))
		return
	}
	todo, err := UpdateTodoServices(uint(id), json.Title, json.Description, json.Completed)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, CreateErrorDtoWithMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, CreateTodoUpdatedDto(&todo))

}

func DeleteTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, CreateErrorDtoWithMessage("You must set an ID"))
		return
	}
	todo, err := FetchById(uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, CreateErrorDtoWithMessage("todo not found"))
		return
	}

	err = DeleteTodoServices(&todo)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, CreateErrorDtoWithMessage("Could not delete Todo"))
		return
	}

	c.JSON(http.StatusOK, CreateSuccessWithMessageDto("todo deleted successfully"))
}

func DeleteAllTodos(c *gin.Context) {
	DeleteAllTodosServices()
	c.JSON(http.StatusOK, CreateErrorDtoWithMessage("All todos deleted successfully"))
}
