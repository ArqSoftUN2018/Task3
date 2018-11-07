package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"github.com/itsjamie/gin-cors"
	_ "github.com/go-sql-driver/mysql"
)

func main(){
	db,err := sql.Open("mysql","root:password@tcp(lists-db:3306)/lists")
	if err != nil{
		fmt.Println(err.Error())
	}
	defer db.Close()

	err = db.Ping() 
	if err != nil{
		fmt.Println(err.Error())
	}

	type list struct {
		id int `form:"id" json:"id" binding:"required"`
		name string `form:"name" json:"name" binding:"required"`
		board string `form:"board" json:"board" binding:"required"`
		archived bool `form:"archived" json:"archived" binding:"required"`
	}

	router := gin.New()
	router.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "GET, PUT, POST, DELETE",
		RequestHeaders: "Origin, Authorization, Content-Type",
		ExposedHeaders: "",
		MaxAge: 50 * time.Second,
		Credentials: true,
		ValidateHeaders: false,
	}))
}
