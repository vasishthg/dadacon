package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
		session := sessions.Default(c)
		fmt.Println(session.Get("loggedin"))
		if session.Get("loggedin") != true {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}
		type agent struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Alias    string `json:"alias"`
			Pfp      string `json:"pfp"`
			Files    []byte `json:"files"`
			Missions []byte `json:"missions"`
			Code     string `json:"code"`
			Level    int    `json:"level"`
		}
		Agent := agent{}
		err := db.QueryRow("SELECT id, name, alias, pfp, files, missions, code, level FROM agent WHERE id = ?", session.Get("ID")).Scan(&Agent.ID, &Agent.Name, &Agent.Alias, &Agent.Pfp, &Agent.Files, &Agent.Missions, &Agent.Code, &Agent.Level)
		if err != nil {
			fmt.Println(err)
		}
		var fileNames []string
		if err := json.Unmarshal(Agent.Files, &fileNames); err != nil {
			fmt.Println(err)
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"agent":     Agent,
			"fileNames": fileNames,
		})
	})
	router.DELETE("/delete", func(c *gin.Context) {
		filename := c.Query("filename")
		if filename == "" {
			c.String(http.StatusBadRequest, "Filename not provided")
			return
		}
		filepath := "static/uploads/files/" + filename
		err := os.Remove(filepath)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error deleting file")
			return
		}
		session := sessions.Default(c)
		agentID := session.Get("ID")
		var existingFilesJSON []byte
		err = db.QueryRow("SELECT files FROM agent WHERE id = ?", agentID).Scan(&existingFilesJSON)
		if err != nil && err != sql.ErrNoRows {
			c.String(http.StatusInternalServerError, "Database error")
			return
		}
		var existingFiles []string
		if err := json.Unmarshal(existingFilesJSON, &existingFiles); err != nil {
			c.String(http.StatusInternalServerError, "Error decoding JSON")
			return
		}
		var updatedFiles []string
		for _, f := range existingFiles {
			if f != filename {
				updatedFiles = append(updatedFiles, f)
			}
		}
		updatedFilesJSON, err := json.Marshal(updatedFiles)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error creating JSON")
			return
		}

		_, err = db.Exec("UPDATE agent SET files = ? WHERE id = ?", updatedFilesJSON, agentID)
		if err != nil {
			c.String(http.StatusInternalServerError, "Database error")
			return
		}

		c.String(http.StatusOK, "File deleted successfully")
		c.String(http.StatusOK, "'%s' deleted", filename)
	})

	router.GET("/chat", func(c *gin.Context) {
		type Message struct {
			ID        int    `json:"id"`
			Sender    string `json:"sender"`
			Receiver  string `json:"receiver"`
			Message   string `json:"message"`
			Timestamp time.Time
		}
		session := sessions.Default(c)
		fmt.Println(session.Get("loggedin"))
		if session.Get("loggedin") != true {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}
		c.HTML(http.StatusOK, "chat.html", gin.H{})
	})

	router.GET("/map", func(c *gin.Context) {
		session := sessions.Default(c)
		fmt.Println(session.Get("loggedin"))
		if session.Get("loggedin") != true {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}
		c.HTML(http.StatusOK, "map.html", gin.H{})
	})

	router.GET("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		fmt.Println(session.Get("loggedin"))
		// login := struct {
		// 	Username string `json:"username"`
		// 	Password string `json:"password"`
		// 	Code     string `json:"code"`
		// }{}

		if session.Get("loggedin") != true {
			c.HTML(http.StatusOK, "login.html", gin.H{})
		} else {
			c.Redirect(http.StatusTemporaryRedirect, "/")

			return
		}
	})
	router.GET("/resources", func(c *gin.Context) {
		session := sessions.Default(c)
		fmt.Println(session.Get("loggedin"))
		if session.Get("loggedin") != true {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}
		c.HTML(http.StatusOK, "resources.html", gin.H{})
	})
	router.GET("/request", func(c *gin.Context) {
		session := sessions.Default(c)
		fmt.Println(session.Get("loggedin"))
		if session.Get("loggedin") != true {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}
		c.HTML(http.StatusOK, "request.html", gin.H{})
	})

	router.POST("/login/verify", func(c *gin.Context) {
		session := sessions.Default(c)
		username := c.PostForm("username")
		password := c.PostForm("password")

		type agent struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Alias    string `json:"alias"`
			Password string `json:"password"`
			Pfp      string `json:"pfp"`
			Files    []byte `json:"files"`
			Missions []byte `json:"missions"`
			Code     string `json:"code"`
			Level    int    `json:"level"`
		}

		Agent := agent{}

		code := session.Get("code")

		err := db.QueryRow("SELECT id, name, alias, password, pfp, files, missions, code, level FROM agent WHERE code = ?", code).Scan(
			&Agent.ID, &Agent.Name, &Agent.Alias, &Agent.Password, &Agent.Pfp, &Agent.Files, &Agent.Missions, &Agent.Code, &Agent.Level)
		if err != nil {
			fmt.Println(err)
		}
		AgentStr := strconv.Itoa(Agent.ID)
		if username != AgentStr {
			c.JSON(http.StatusConflict, gin.H{"error": "Agent cipher different from ID"})
		} else {
			err := db.QueryRow("SELECT id, name, alias, password, pfp, files, missions, code, level FROM agent WHERE id = ?", username).Scan(
				&Agent.ID, &Agent.Name, &Agent.Alias, &Agent.Password, &Agent.Pfp, &Agent.Files, &Agent.Missions, &Agent.Code, &Agent.Level)
			if err != nil {
				fmt.Println(err)
			}
			if Agent.Password == password {
				session.Set("loggedin", true)
				session.Set("name", Agent.Name)
				session.Set("alias", Agent.Alias)
				session.Set("pfp", Agent.Pfp)
				session.Set("files", Agent.Files)
				session.Set("missions", Agent.Missions)
				session.Set("code", Agent.Code)
				session.Set("level", Agent.Level)
				session.Save()
				c.Redirect(http.StatusFound, "/")
			} else {
				c.JSON(http.StatusExpectationFailed, gin.H{"error": "Invalid password"})
			}
		}
	})
	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "Bad request")
			return
		}
		session := sessions.Default(c)
		filepath := "static/uploads/files/" + file.Filename
		err = c.SaveUploadedFile(file, filepath)
		if err != nil {
			fmt.Println("Error saving file:", err)
			c.String(http.StatusInternalServerError, "Error saving file")
			return
		}
		var existingFilesJSON []byte
		err = db.QueryRow("SELECT files FROM agent WHERE id = ?", session.Get("ID")).Scan(&existingFilesJSON)
		if err != nil && err != sql.ErrNoRows {
			c.String(http.StatusInternalServerError, "Database error")
			return
		}
		var existingFiles []string
		if existingFilesJSON != nil {
			err = json.Unmarshal(existingFilesJSON, &existingFiles)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error decoding JSON")
				return
			}
		}
		existingFiles = append(existingFiles, file.Filename)
		updatedFilesJSON, err := json.Marshal(existingFiles)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error creating JSON")
			return
		}
		_, err = db.Exec("UPDATE agent SET files = ? WHERE id = ?", updatedFilesJSON, session.Get("ID"))
		if err != nil {
			c.String(http.StatusInternalServerError, "Database error")
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("Uploaded '%s'", file.Filename))
	})

	router.POST("/check/cipher", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}
		if !strings.HasSuffix(file.Filename, ".txt") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only .txt files are allowed"})
			return
		}
		uploadedFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening file"})
			return
		}
		defer uploadedFile.Close()
		content, err := ioutil.ReadAll(uploadedFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file content"})
			return
		}
		processedContent := string(content)
		type agent struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Alias    string `json:"alias"`
			Password string `json:"password"`
			Pfp      string `json:"pfp"`
			Files    string `json:"files"`
			Missions string `json:"missions"`
			Code     string `json:"code"`
			Level    int    `json:"level"`
		}
		Agent := agent{}
		err = db.QueryRow("SELECT id, name, alias, password, pfp, files, missions, code, level FROM agent WHERE code = ?", processedContent).Scan(&Agent.ID, &Agent.Name, &Agent.Alias, &Agent.Password, &Agent.Pfp, &Agent.Files, &Agent.Missions, &Agent.Code, &Agent.Level)
		if err != nil {
			fmt.Println(err)
		}
		if Agent.Code == processedContent {
			session := sessions.Default(c)
			session.Set("ID", Agent.ID)
			session.Set("code", processedContent)
			session.Save()
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "Invalid code"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": processedContent})
	})

	router.Run(":5500")
}
