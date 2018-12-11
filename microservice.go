package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)


func main(){
	db,err := sql.Open("mysql","root:password@tcp(localhost:3306)/lists")
	if err != nil{
		fmt.Println(err.Error())
	}
	defer db.Close()

	err = db.Ping() 
	if err != nil{
		fmt.Println(err.Error())
	}

	type List struct {
		Id int `form:"Id" json:"id"`
		Name string `form:"Name" json:"name"`
		Board string `form:"Board" json:"board"`
		Archived bool `form:"Archived" json:"archived"`
	}

	router := gin.Default()
	
	router.GET("/lists-ms/resources/lists/:id", func(c * gin.Context){
		var (
			object List
			result gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("Select * from list where id = ?;", id)
		err = row.Scan(&object.Id,&object.Name,&object.Board,&object.Archived)
		if err != nil {
			result = gin.H {
			}
		}else{
			result = gin.H {
				"id": object.Id,
				"name": object.Name,
				"board": object.Board,
				"archived": object.Archived,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	router.GET("/lists-ms/resources/lists/", func(c * gin.Context){
		var (
			object List
			objects gin.H
			result []gin.H
		)
		rows,err := db.Query("Select * from list;")
		if err != nil {
			fmt.Println(err.Error())
		}  
		for rows.Next(){
			err := rows.Scan(&object.Id,&object.Name,&object.Board,&object.Archived)
			objects = gin.H {
				"id": object.Id,
				"name": object.Name,
				"board": object.Board,
				"archived": object.Archived,
			}
			result = append(result,objects)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, result)
	})

	router.GET("/lists-ms/resources/listsFromBoard/:board", func(c * gin.Context){
		var (
			object List
			objects gin.H
			result []gin.H
		)
		board := c.Param("board")
		rows,err := db.Query("Select * from list where board = ?;", board)
		if err != nil {
			fmt.Println(err.Error())
		}  
		for rows.Next(){
			err := rows.Scan(&object.Id,&object.Name,&object.Board,&object.Archived)
			objects = gin.H {
				"id": object.Id,
				"name": object.Name,
				"board": object.Board,
				"archived": object.Archived,
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
		var input List
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		Name := input.Name
		Board := input.Board
		Archived := false
		stmt, err := db.Prepare("insert into list (Name, Board, Archived) values(?,?,?);")
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = stmt.Exec(Name, Board,Archived)

		if err != nil {
			fmt.Println(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"Mensaje": fmt.Sprintf("se ha creado la lista exitosamente"),
		})
	})
	router.PUT("/lists-ms/resources/lists/:id", func(c * gin.Context){
		var input List
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id := c.Param("id")
		name := input.Name
		archived := input.Archived
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

	router.DELETE("/lists-ms/resources/lists/:id", func(c * gin.Context){
		id := c.Param("id")
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


