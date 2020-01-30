package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// Todo is the item.
type Todo struct {
	gorm.Model
	Text   string
	Status string
	Person string
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	dbInit()

	router.GET("/", func(ctx *gin.Context) {
		todos := dbGetAll()
		ctx.HTML(200, "index.html", gin.H{
			"todos": todos,
		})
	})

	router.POST("/new", func(ctx *gin.Context) {
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		person := ctx.PostForm("person")
		dbInsert(text, status, person)
		ctx.Redirect(302, "/")
	})

	router.GET("/details/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := dbGetOne(id)
		ctx.HTML(200, "detail.html", gin.H{"todo": todo})
	})

	router.POST("/update/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}

		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		person := ctx.PostForm("person")
		dbUpdate(id, text, status, person)
		ctx.Redirect(302, "/")
	})

	router.GET("/delete_check/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := dbGetOne(id)
		ctx.HTML(200, "delete.html", gin.H{"todo": todo})
	})

	router.POST("/delete/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		dbDelete(id)
		ctx.Redirect(302, "/")
	})

	router.Run()
}

func dbInit() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("Error.Open DB has failed.")
	}

	db.AutoMigrate(&Todo{})
	defer db.Close()
}

func dbInsert(text string, status string, person string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("Error,Insert DB has failed.")
	}
	db.Create(&Todo{Text: text, Status: status, Person: person})
	defer db.Close()
}

func dbUpdate(id int, text string, status string, person string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("Error,Update DB has failed.")
	}

	var todo Todo
	db.First(&todo, id)
	todo.Text = text
	todo.Status = status
	todo.Person = person
	db.Save(&todo)
	db.Close()

}

func dbDelete(id int) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("Error,Delete DB has failed.")
	}

	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
	db.Close()
}

func dbGetAll() []Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("Error,Delete all DB has failed.")
	}

	var todos []Todo
	db.Order("created_at desc").Find(&todos)
	db.Close()
	return todos
}

func dbGetOne(id int) Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("Error,Could not get DB.")
	}

	var todo Todo
	db.First(&todo, id)
	db.Close()
	return todo
}
