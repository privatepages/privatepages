package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"package/main/internal/config"
)

var conf *config.Config

func init() {
	conf = config.Cfg
}

func Upload(c *gin.Context) {

	// get the vars off req body
	var body struct {
		Token        string
		Artifactname string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// auth

	if body.Token == conf.API_SECRET { // try bcrypt on browser
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect token",
		})
		return
	}

	// validate artifactname

	// // body.Artifactname
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Incorrect artifactname",
	// 	})
	// 	return
	// }

	// validate filename from form-data

	// validate file size

	// validate tar or gzip

	// unarchive

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}
