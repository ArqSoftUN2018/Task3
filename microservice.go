package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main(){
	db,err := sql.Open("mysql","root:password@tcp(127.0.0.1:3306/lists)")
	if err != nil{
		fmt.Println(err.Error())
	}
	defer db.Close()

	err = db.Ping() 
	if err != nil{
		fmt.Println(err.Error())
	}

	type listas struct {
		id int
		nombre string
		tablero string 
		archivado bool
	}

	router := gin.Default()
	router.GET("/lista/:id", func(c * gin.Context){
		var (
			objeto listas
			resultado gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("Select * from listas where id = ?;", id)
		err := row.Scan(&objeto.id,&objeto.nombre,&objeto.tablero,&objeto.archivado)
		if err != nil {
			resultado = gin.H {
				"resultado": nil,
				"cantidad":0,
			}
		}else{
			resultado = gin.H {
				"resultado": objeto,
				"cantidad": 1,
			}
			c.JSON(http.StatusOK, resultado)
		}
	})

	router.GET("/listas", func(c * gin.Context){
		var (
			objeto listas
			objetos []listas
		)
		rows,err := db.Query("Select * from listas;")
		if err != nil {
			fmt.Println(err.Error())
		}  
		for rows.Next(){
			err := rows.Scan(&objeto.id,&objeto.nombre,&objeto.tablero,&objeto.archivado)
			objetos = append(objetos,objeto)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"resultado": objetos,
			"cantidad": len(objetos),
		})
	})
	router.POST("/lista", func(c * gin.Context){
		nombre := c.PostForm("nombre")
		tablero := c.PostForm("tablero")
		archivado := false
		stmt, err := db.Prepare("insert into listas (nombre, tablero, archivado) values(?,?);")
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = stmt.Exec(nombre, tablero,archivado)

		if err != nil {
			fmt.Println(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"Mensaje": fmt.Sprintf("se ha creado la lista exitosamente"),
		})
	})
	router.PUT("/lista", func(c * gin.Context){
		id := c.Query("id")
		nombre := c.PostForm("nombre")
		archivado := c.PostForm("tablero")
		stmt, err := db.Prepare("update listas set nombre = ?, archivado = ? where id = ?;")
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = stmt.Exec(nombre, archivado, id)

		if err != nil {
			fmt.Println(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"Mensaje": fmt.Sprintf("se ha actualizado la lista exitosamente"),
		})
	})

	router.DELETE("/lista", func(c * gin.Context){
		id := c.PostForm("id")
		stmt, err := db.Prepare("delete from listas where id = ?;")
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
	router.Run(":8080")
}