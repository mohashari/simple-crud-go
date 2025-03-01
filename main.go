package main

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	var err error
		db, err = gorm.Open("mysql", "root:welcome1@tcp(127.0.0.1:3306)/golang?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	//Migrate the schema
	db.AutoMigrate(&todoModel{})
}

func main() {

	router := gin.Default()
		v1 := router.Group("/api/v1/todos")
	{
		v1.POST("", createTodo)
		v1.GET("", fetchAllTodo)
		v1.GET("/:id", fetchSingleTodo)
		v1.PUT("/:id", updateTodo)
		v1.DELETE("/:id", deleteTodo)
	}
	router.Run(":8081")
}

type (
	todoModel struct {
		gorm.Model
		Title     string `json:"title"`
		Completed int    `json:"completed"`
	}

	transformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)

func createTodo(c *gin.Context) {
	//log.Info("REQUEST FROM GIN: ",c)
	//completed, _:= strconv.Atoi(c.PostForm("completed"))
	//todo := todoModel{Title: c.PostForm("title"), Completed: completed}
	//log.Info(todo)
	//db.Create(&todo)
	//c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated,
	//	"message": "Todo Created Item Successfully", "resourceId": todo.ID})

	var todo todoModel
	c.Bind(&todo)
	if todo.Title != ""{
		db.Create(&todo)
		c.JSON(http.StatusOK,gin.H{"status": http.StatusOK,"result":todo})
	}else {
		c.JSON(http.StatusOK,gin.H{"status":http.StatusOK,"result":"Failure"})
	}


}

func fetchAllTodo(c *gin.Context) {
	var todos [] todoModel
	var _todos []transformedTodo

	db.Find(&todos)

	if len(todos) <= 0 {

		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo Found"})
		return
	}

	for _ , item := range todos {

		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}

		_todos = append(_todos, transformedTodo{ID: item.ID, Title: item.Title, Completed: completed})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todos})
}

func fetchSingleTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo Found!"})
		return
	}

	completed := false
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}

	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": _todo})
}

func updateTodo(c*gin.Context)  {
	var todo todoModel
	todoID:=c.Param("id")
	db.First(&todo,todoID)

	if todo.ID == 0{
		c.JSON(http.StatusNotFound,gin.H{"status":http.StatusNotFound,"message":"No Todo Not Found"})
		return
	}
	db.Model(&todo).Update("title",c.PostForm("title"))
	completed,  _:= strconv.Atoi(c.PostForm("completed"))
	db.Model(&todo).Update("completed",completed)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!"})

}

func deleteTodo(c*gin.Context)  {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo,todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound,gin.H{"status":http.StatusNotFound,"message":"todo Not Found"})
	}

	db.Delete(&todo)
	c.JSON(http.StatusOK,gin.H{"status":http.StatusOK,"message":"Todo deleted sucessfully"})
}