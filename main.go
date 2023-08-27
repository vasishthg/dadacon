package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	store := cookie.NewStore([]byte("systumisriyal"))
	router.Use(sessions.Sessions("dadacon", store))
	router.Static("/static", "./static")
	db, err := sql.Open("mysql", "root:toor@tcp(localhost:3306)/dadacon")
	if err != nil {
		log.Fatal("err connecting to db:", err)
	}
	defer db.Close()
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	router.Run(":5500")
}
