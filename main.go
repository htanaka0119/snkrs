package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	//使ってないけどimportしないとエラーが出るのでimportする
	//その場合は頭に _ をつける
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	gorm.Model
	Text   string
	Status string
}

//DBマイグレート
func dbInit() {
	db, error := gorm.Open("sqlite3", "test.sqlite3")
	if error != nil {
		panic("データベース開けず!(dbInit)")
	}
	db.AutoMigrate(&Todo{})
	defer db.Close()
}

//DB追加
func dbInsert(text string, status string) {
	db, error := gorm.Open("sqlite3", "test.sqlite3")
	if error != nil {
		panic("データベース開けず!(dbInsert)")
	}
	db.Create(&Todo{Text: text, Status: status})
	defer db.Close()
}

func dbUpdate(id int, text string, status string) {
	db, error := gorm.Open("sqlite3", "test.sqlite3")
	if error != nil {
		panic("データベース開けず(dbUpdate)")
	}
	var todo Todo
	db.First(&todo, id)
	todo.Text = text
	todo.Status = status
	db.Save(&todo)
	db.Close()
}

func dbDelete(id int) {
	db, error := gorm.Open("sqlite3", "test.sqlite3")
	if error != nil {
		panic("データベース開けず(dbDelete)")
	}
	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
	db.Close()
}

func dbGetAll() []Todo {
	db, error := gorm.Open("sqlite3", "test.sqlite3")
	if error != nil {
		panic("データベース開けず(doGetAll)")
	}
	var todos []Todo
	db.Order("created_at desc").Find(&todos)
	db.Close()
	return todos
}

func dbGetOne(id int) Todo {
	db, error := gorm.Open("sqlite3", "test.sqlite3")
	if error != nil {
		panic("データベース開けず(dbGetOne)")
	}
	var todo Todo
	db.First(&todo, id)
	db.Close()
	return todo
}

func todaysKicks() []Todo {
	db, error := gorm.Open("sqlite3", "test.sqlite3")
	if error != nil {
		panic("データベース開けず(doGetAll)")
	}
	var todos Todo
	db.Order("created_at desc").Find(&todos)
	db.Close()
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(len(todos))
	return todos[num]
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	dbInit()

	router.GET("/", func(ctx *gin.Context) {
		todos := dbGetAll()
		rand.Seed(time.Now().UnixNano())
		snkrs := todaysKicks()
		fmt.Printf("今日の気分は%v!!\n", snkrs)
		ctx.HTML(200, "index.html", gin.H{"todos": todos})
		ctx.HTML(200, "index.html", gin.H{"kicks": snkrs})
	})

	router.POST("/new", func(ctx *gin.Context) {
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		dbInsert(text, status)
		ctx.Redirect(302, "/")
	})

	router.GET("/detail/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, error := strconv.Atoi(n)
		if error != nil {
			panic(error)
		}
		todo := dbGetOne(id)
		ctx.HTML(200, "detail.html", gin.H{"todo": todo})
	})

	router.POST("/update/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, error := strconv.Atoi(n)
		if error != nil {
			panic(error)
		}
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		dbUpdate(id, text, status)
		ctx.Redirect(302, "/")
	})

	router.GET("/delete_check/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, error := strconv.Atoi(n)
		if error != nil {
			panic(error)
		}
		todo := dbGetOne(id)
		ctx.HTML(200, "delete.html", gin.H{"todo": todo})
	})

	router.POST("/delete/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, error := strconv.Atoi(n)
		if error != nil {
			panic(error)
		}
		dbDelete(id)
		ctx.Redirect(302, "/")
	})
	router.Run()
}
