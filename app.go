package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/", getMemoList)
	r.POST("/", addMemo)
	r.POST("/delete", deleteMemo)

	r.Run()
}

func getConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./db/memoapp.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// Memo struct
type Memo struct {
	ID          int64
	Subject     string
	Description string
	CreatedAt   string
}

func getMemoList(c *gin.Context) {
	db := getConnection()
	defer db.Close()

	rows, err := db.Query("select id, subject, description, created_at from memo order by created_at")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	memoList := []Memo{}
	for rows.Next() {
		var id int64
		var subject string
		var description string
		var createdAt string
		err := rows.Scan(&id, &subject, &description, &createdAt)
		if err != nil {
			log.Fatal(err)
		}

		memo := Memo{}
		memo.ID = id
		memo.Subject = subject
		memo.Description = description
		memo.CreatedAt = createdAt
		memoList = append(memoList, memo)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"memoList": memoList,
	})
}

func addMemo(c *gin.Context) {
	subject := c.PostForm("subject")
	description := c.PostForm("description")

	db := getConnection()
	defer db.Close()

	stmt, err := db.Prepare("insert into memo (subject, description, created_at) values (?, ?, datetime('now', 'localtime'))")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(subject, description)
	if err != nil {
		log.Fatal(err)
	}

	c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteMemo(c *gin.Context) {
	id := c.PostForm("id")
	db := getConnection()
	defer db.Close()

	stmt, err := db.Prepare("delete from memo where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}

	c.Redirect(http.StatusMovedPermanently, "/")
}
