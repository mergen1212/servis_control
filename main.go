package main

import (
	"fmt"
	"net/http"

	"strconv"
	"watchman/pkg/db"

	"github.com/gin-gonic/gin"
)

func getHostPort() (string, int) {
	return "127.0.0.1", 8000
}

func sayAliveView(c *gin.Context, db db.Database) {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	_, err = db.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
	}

	projectHash := c.Param("project_hash")
	project, err := db.GetProjectByHash(projectHash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
	}

	if project.UserID != userId {
		c.JSON(http.StatusNotFound, gin.H{})
	}

	c.JSON(http.StatusOK, gin.H{})
}

// TODO: how to do it without this?
func buildHandler(fn func(c *gin.Context, db db.Database), db db.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(c, db)
	}
}

// TODO: why does router function return database instance?
func getRouter() (*gin.Engine, db.Database) {
	db, err := db.GetDB()
	if err != nil {
		panic(err)
	}

	err = db.PrepareDB()
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.GET("/:user_id/:project_hash", buildHandler(sayAliveView, db))
	return router, db
}

func main() {
	router, db := getRouter()
	defer db.Close()

	host, port := getHostPort()
	router.Run(fmt.Sprintf("%s:%d", host, port))
}
