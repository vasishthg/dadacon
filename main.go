package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xyproto/randomstring"
)

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
type missions struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Progress    string `json:"progress"`
	Threat      string `json:"threat"`
	Remaining   string `json:"remaining"`
	S1          string `json:"s1"`
	S2          string `json:"s2"`
	S3          string `json:"s3"`
	Agents      []byte `json:"agents"`
	Location    string `json:"location"`
}
type notif struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Readby []byte `json:"readby"`
	Unread bool   `json:"unread"`
}

type resource struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	File string `json:"file"`
	Code string `json:"code"`
}

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
		Agent := agent{}
		err := db.QueryRow("SELECT id, name, alias, pfp, files, missions, code, level FROM agent WHERE id = ?", session.Get("ID")).Scan(&Agent.ID, &Agent.Name, &Agent.Alias, &Agent.Pfp, &Agent.Files, &Agent.Missions, &Agent.Code, &Agent.Level)
		if err != nil {
			fmt.Println(err)
		}
		var fileNames []string
		if err := json.Unmarshal(Agent.Files, &fileNames); err != nil {
			fmt.Println(err)
		}
		var missionIDs []int
		if err := json.Unmarshal(Agent.Missions, &missionIDs); err != nil {
			fmt.Println(err)
		}
		var systum []missions
		for _, missionID := range missionIDs {
			row := db.QueryRow("SELECT id, title, description, progress, threat, remaining, s1, s2, s3, agents, location FROM missions WHERE id=?", missionID)

			var Mission missions
			err := row.Scan(&Mission.ID, &Mission.Title, &Mission.Description, &Mission.Progress, &Mission.Threat, &Mission.Remaining, &Mission.S1, &Mission.S2, &Mission.S3, &Mission.Agents, &Mission.Location)
			if err != nil {
				log.Println(err)
				continue
			}

			systum = append(systum, Mission)
		}
		rows, err := db.Query("SELECT id, title, description, progress, threat, remaining, s1, s2, s3, agents, location FROM missions")
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()
		Mission := missions{}
		var allMissions []missions
		for rows.Next() {
			err := rows.Scan(&Mission.ID, &Mission.Title, &Mission.Description, &Mission.Progress, &Mission.Threat, &Mission.Remaining, &Mission.S1, &Mission.S2, &Mission.S3, &Mission.Agents, &Mission.Location)
			if err != nil {
				log.Println(err)
			}
			allMissions = append(allMissions, Mission)

		}
		var availableMissions []missions
		for _, allMission := range allMissions {
			if !slices.Contains(missionIDs, allMission.ID) {
				availableMissions = append(availableMissions, allMission)
			}
		}
		rows, err = db.Query("SELECT id, title, `desc`, readby FROM notifications")
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()

		var notifs []notif
		for rows.Next() {
			var Notif notif
			err := rows.Scan(&Notif.ID, &Notif.Title, &Notif.Desc, &Notif.Readby)
			if err != nil {
				log.Println(err)
			}
			notifs = append(notifs, Notif)
		}
		fmt.Println(notifs)

		var unreadNotifs []notif
		for _, n := range notifs {
			if n.Readby != nil { // Assuming Readby is a []byte field
				var readbyIDs []int
				err := json.Unmarshal(n.Readby, &readbyIDs)
				if err != nil {
					log.Println(err)
					continue
				}

				if !slices.Contains(readbyIDs, Agent.ID) {
					n.Unread = true
					unreadNotifs = append(unreadNotifs, n)
				}
			}
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"agent":             Agent,
			"fileNames":         fileNames,
			"agentMissions":     systum,
			"missions":          allMissions,
			"availableMissions": availableMissions,
			"notifs":            notifs,
			"unreadNotifs":      unreadNotifs,
		})
	})

	router.POST("/notification/read", func(c *gin.Context) {
		session := sessions.Default(c)
		fmt.Println(session.Get("loggedin"))
		if session.Get("loggedin") != true {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}
		Agent := agent{}
		err := db.QueryRow("SELECT id, name, alias, pfp, files, missions, code, level FROM agent WHERE id = ?", session.Get("ID")).Scan(&Agent.ID, &Agent.Name, &Agent.Alias, &Agent.Pfp, &Agent.Files, &Agent.Missions, &Agent.Code, &Agent.Level)
		if err != nil {
			fmt.Println(err)
		}
		rows, err := db.Query("SELECT id, readby FROM notifications")
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var readby []byte
			err := rows.Scan(&id, &readby)
			if err != nil {
				fmt.Println(err)
				continue
			}

			var readbyIDs []int
			err = json.Unmarshal(readby, &readbyIDs)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if !slices.Contains(readbyIDs, Agent.ID) {
				readbyIDs = append(readbyIDs, Agent.ID)
				updatedReadbyJSON, err := json.Marshal(readbyIDs)
				if err != nil {
					fmt.Println(err)
					continue
				}

				_, err = db.Exec("UPDATE notifications SET readby = ? WHERE id = ?", updatedReadbyJSON, id)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	})

	router.POST("/notification/mayday", func(c *gin.Context) {

		session := sessions.Default(c)
		fmt.Println(session.Get("loggedin"))
		if session.Get("loggedin") != true {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}
		// https://www.google.com/maps/dir/?api=1&destination=34.059808,-118.368152

		Agent := agent{}
		err := db.QueryRow("SELECT id, name, alias, pfp, files, missions, code, level FROM agent WHERE id = ?", session.Get("ID")).Scan(&Agent.ID, &Agent.Name, &Agent.Alias, &Agent.Pfp, &Agent.Files, &Agent.Missions, &Agent.Code, &Agent.Level)
		if err != nil {
			fmt.Println(err)
		}
		notification := "Mayday! Mayday! Mayday! Agent " + "`" + Agent.Alias + "`" + " is in danger! Location: " + "https://www.google.com/maps/dir/?api=1&destination=" + c.PostForm("lat") + "," + c.PostForm("lng")
		fmt.Print("INSERT INTO notifications (title, desc) VALUES (?, ?, ?)", "Mayday!", notification)
		_, err = db.Exec("INSERT INTO `notifications` (title, `desc`, readby) VALUES(?, ?, ?)", "Mayday!", notification, []byte("[]"))
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec("UPDATE agent SET code = ? WHERE id = ?", randomstring.CookieFriendlyString(30), Agent.ID)
		if err != nil {
			fmt.Println(err)
		}
		session.Clear()
		session.Save()
		c.JSON(http.StatusOK, gin.H{"notification": notification})
	})

	router.POST("/request/mission/:id", func(c *gin.Context) {
		session := sessions.Default(c)
		fmt.Println(session.Get("loggedin"))
		id := c.Param("id")
		agentID := session.Get("ID")
		Agent := agent{}
		err := db.QueryRow("SELECT missions FROM agent WHERE id = ?", agentID).Scan(&Agent.Missions)
		if err != nil && err != sql.ErrNoRows {
			c.String(http.StatusInternalServerError, "Database error")
			return
		}

		var missionIDs []int
		if err := json.Unmarshal(Agent.Missions, &missionIDs); err != nil {
			fmt.Println(err)
		}

		fmt.Println(missionIDs)
		newMissionID, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println(err)
		}
		if slices.Contains(missionIDs, newMissionID) {
			c.JSON(http.StatusConflict, gin.H{"error": "Mission already added"})
			return
		} else {
			missionIDs = append(missionIDs, newMissionID)
		}
		updatedMissions, err := json.Marshal(missionIDs)
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec("UPDATE agent SET missions = ? WHERE id = ?", updatedMissions, agentID)
		if err != nil {
			fmt.Println(err)
		}

		c.JSON(http.StatusOK, gin.H{"missionIDs": missionIDs})
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
		// type Message struct {
		// 	ID        int    `json:"id"`
		// 	Sender    string `json:"sender"`
		// 	Receiver  string `json:"receiver"`
		// 	Message   string `json:"message"`
		// 	Timestamp time.Time
		// }
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
		rows, err := db.Query("SELECT id, name, file, file FROM resources")
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()

		var resources []resource
		for rows.Next() {
			var Resource resource
			err := rows.Scan(&Resource.ID, &Resource.Name, &Resource.File, &Resource.Code)
			if err != nil {
				log.Println(err)
			}
			resources = append(resources, Resource)
		}
		fmt.Println(resources)

		c.HTML(http.StatusOK, "resources.html", gin.H{
			"resources": resources,
		})
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

	router.POST("/resource/checkcode", func(c *gin.Context) {
		id := c.PostForm("id")
		code := c.PostForm("code")
		fmt.Println(id, code)
		Resource := resource{}
		err := db.QueryRow("SELECT id, name, file, code FROM resources WHERE id = ?", id).Scan(&Resource.ID, &Resource.Name, &Resource.File, &Resource.Code)
		if err != nil {
			fmt.Println(err)
		}
		if Resource.Code == code {
			c.JSON(http.StatusOK, gin.H{"url": Resource.File, "data": "true"})
		} else {
			c.JSON(http.StatusOK, gin.H{"data": "false"})
			return
		}

	})

	router.Run(":5500")
}
