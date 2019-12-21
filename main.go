package main

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	gorm.Model
	Text   string
	Status string
}

type User struct {
	gorm.Model
	UserID   string
	Password string
}

//UserDB初期化
func dbInitUser() {
	db, err := gorm.Open("sqlite3", "user.sqlite3")
	if err != nil {
		panic("データベース開けず！（dbInit）")
	}
	db.AutoMigrate(&User{})
	defer db.Close()
}

//DB初期化
func dbInit() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！（dbInit）")
	}
	db.AutoMigrate(&Todo{})
	defer db.Close()
}

//DBUSER追加
func dbUserInsert(userid string, password string) {
	db, err := gorm.Open("sqlite3", "user.sqlite3")
	if err != nil {
		panic("データベース開けず！（dbInsert)")
	}
	db.Create(&User{UserID: userid, Password: password})
	defer db.Close()
}

//DB追加
func dbInsert(text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！（dbInsert)")
	}
	db.Create(&Todo{Text: text, Status: status})
	defer db.Close()
}

//DB更新
func dbUpdate(id int, text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！（dbUpdate)")
	}
	var todo Todo
	db.First(&todo, id)
	todo.Text = text
	todo.Status = status
	db.Save(&todo)
	db.Close()
}

//DB削除
func dbDelete(id int) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！（dbDelete)")
	}
	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
	db.Close()
}

//DB全取得
func dbGetAll() []Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！(dbGetAll())")
	}
	var todos []Todo
	db.Order("created_at desc").Find(&todos)
	db.Close()
	return todos
}

//DB一つ取得
func dbGetOne(id int) Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！(dbGetOne())")
	}
	var todo Todo
	db.First(&todo, id)
	db.Close()
	return todo
}

//DBUSER一つ取得
func dbGetOneUser(userid string) User {
	db, err := gorm.Open("sqlite3", "user.sqlite3")
	if err != nil {
		panic("データベース開けず！(dbGetOne())")
	}
	var user User
	db.First(&user, userid)
	db.Close()
	return user
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	dbInit()
	dbInitUser()

	//login
	router.GET("/login", func(ctx *gin.Context){
		ctx.HTML(200,"login.html", gin.H{})
	})
	router.POST("/login", func(ctx *gin.Context){
		userid := ctx.PostForm("userid")
		password := ctx.PostForm("password")
		user := dbGetOneUser(userid)
		result := ""
		if user.Password != password {
			result = "IDもしくはPasswordが間違っています"
			ctx.HTML(200,"login.html", gin.H{"abc":result})
		}else {
			ctx.Redirect(302, "/")
		}
	})

	//register
	router.GET("/register", func(ctx *gin.Context){
		ctx.HTML(200,"register.html", gin.H{})
	})
	router.POST("/register", func(ctx *gin.Context){
		userid := ctx.PostForm("userid")
		password := ctx.PostForm("password")
		dbUserInsert(userid, password)
		ctx.Redirect(302,"/login")
	})

	//index
	router.GET("/", func(ctx *gin.Context){
		todos := dbGetAll()
		ctx.HTML(200, "index.html", gin.H{"todos":todos})
	})

	//create
	router.POST("/new",func(ctx *gin.Context){
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		dbInsert(text, status)
		ctx.Redirect(302,"/")
	})

	//detail
	router.GET("/detail/:id",func(ctx *gin.Context){
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := dbGetOne(id)
		ctx.HTML(200, "detail.html", gin.H{"todo":todo})
	})

	// update
	router.POST("/update/:id", func(ctx *gin.Context){
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		dbUpdate(id, text, status)
		ctx.Redirect(302, "/")
	})

	//delete確認
	router.GET("/delete_check/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		todo := dbGetOne(id)
		ctx.HTML(200, "delete.html", gin.H{"todo": todo})
	})

	//Delete
	router.POST("/delete/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		dbDelete(id)
		ctx.Redirect(302, "/")		
	})
	
	router.Run()
}
