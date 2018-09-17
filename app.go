package main

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

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

func getConnection() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./db/memoapp.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// Memo struct
type Memo struct {
	ID          int64 `gorm:"primary_key"`
	Subject     string
	Description string
	CreatedAt   string
}

func getMemoList(c *gin.Context) {
	db := getConnection()
	defer db.Close()

	rows, err := db.Raw("select id, subject, description, created_at from memo order by created_at", nil).Rows()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	memoList := []Memo{}
	for rows.Next() {
		var memo Memo
		err := db.ScanRows(rows, &memo)
		if err != nil {
			log.Fatal(err)
		}
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

	db.Exec("insert into memo (subject, description, created_at) values (?, ?, datetime('now', 'localtime'))", subject, description)

	c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteMemo(c *gin.Context) {
	id := c.PostForm("id")
	db := getConnection()
	defer db.Close()

	db.Exec("delete from memo where id = ?", id)

	c.Redirect(http.StatusMovedPermanently, "/")
}
