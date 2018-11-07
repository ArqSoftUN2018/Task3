package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
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

	router := gin.Default()

	router.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "GET, PUT, POST, DELETE",
		RequestHeaders: "",
		ExposedHeaders: "",
		MaxAge: 50 * time.Second,
		Credentials: true,
		ValidateHeaders: false,
	}))
	
	router.GET("/lists-ms/resources/lists/:id", func(c * gin.Context){
		var (
			object list
			result gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("Select * from list where id = ?;", id)
		err = row.Scan(&object.id,&object.name,&object.board,&object.archived)
		if err != nil {
			result = gin.H {
			}
		}else{
			result = gin.H {
				"id": object.id,
				"name": object.name,
				"board": object.board,
				"archived": object.archived,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	router.GET("/lists-ms/resources/lists/", func(c * gin.Context){
		var (
			object list
			objects gin.H
			result []gin.H
		)
		rows,err := db.Query("Select id,name,board,archived from list;")
		if err != nil {
			fmt.Println(err.Error())
		}  
		for rows.Next(){
			err := rows.Scan(&object.id,&object.name,&object.board,&object.archived)
			objects = gin.H {
				"id": object.id,
				"name": object.name,
				"board": object.board,
				"archived": object.archived,
			}
			result = append(result,objects)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, result)
	})

	router.GET("/lists-ms/resources/lists-board/:id", func(c * gin.Context){
		var (
			object list
			objects gin.H
			result []gin.H
		)
		id := c.Param("id")
		rows,err := db.Query("Select id,name,board,archived from list where board = ?;", id)
		if err != nil {
			fmt.Println(err.Error())
		}  
		for rows.Next(){
			err := rows.Scan(&object.id,&object.name,&object.board,&object.archived)
			objects = gin.H {
				"id": object.id,
				"name": object.name,
				"board": object.board,
				"archived": object.archived,
			}
			result = append(result,objects)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, result)
	})
	
	router.POST("/lists-ms/resources/lists/", func(c * gin.Context){
		name := c.PostForm("name")
		board := c.PostForm("board")
		archived := false
		stmt, err := db.Prepare("insert into list (name, board, archived) values(?,?,?);")
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = stmt.Exec(name, board,archived)

		if err != nil {
			fmt.Println(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"Mensaje": fmt.Sprintf("se ha creado la lista exitosamente"),
		})
	})
	router.PUT("/lists-ms/resources/lists/:id", func(c * gin.Context){
		id := c.Param("id")
		name := c.PostForm("name")
		archived := c.PostForm("archived")
		stmt, err := db.Prepare("update list set name = ?, archived = ? where id = ?;")
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = stmt.Exec(name, archived, id)

		if err != nil {
			fmt.Println(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"Mensaje": fmt.Sprintf("se ha actualizado la lista exitosamente"),
		})
	})

	router.DELETE("/lists-ms/resources/lists/", func(c * gin.Context){
		id := c.PostForm("id")
		stmt, err := db.Prepare("delete from list where id = ?;")
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = stmt.Exec(id)

		if err != nil {
			fmt.Println(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"Mensaje": fmt.Sprintf("se ha borrado la lista exitosamente"),
		})
	})
	router.Run(":3002")
}
