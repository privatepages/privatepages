package controllers

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"

	"package/main/internal/config"
	"package/main/internal/logger"
)

var conf *config.Config
var log = logger.Log

var artifactNamePattern = regexp.MustCompile("^([A-z0-9-_]{2,64})$")

func init() {
	conf = config.Cfg
}

// Upload method for uploading files
func Upload(c *gin.Context) {

	// Get the vars off req body
	var body struct {
		File         *multipart.FileHeader `form:"file" binding:"required"`
		Token        string                `form:"token" binding:"required"`
		Artifactname string                `form:"artifactname" binding:"required"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Auth
	if body.Token != conf.APISecret { // try bcrypt on browser
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Incorrect token",
		})
		return
	}

	// Validate artifactname
	if !artifactNamePattern.MatchString(body.Artifactname) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect artifactname. ^([A-z0-9-_]{2,64})$",
		})
		return
	}

	// Open file for reading
	file, err := body.File.Open()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Failed to read file",
		})
		log.Debug(err)
		return
	}
	defer file.Close()

	// Make temp dir
	tmpDir, err := ioutil.TempDir(conf.ArtifactStoragePath, "tmp")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal problem, see logs",
		})
		log.Error(err)
		return
	}
	defer os.RemoveAll(tmpDir)

	// Validate tar
	log.Debug(body.File.Header) // textproto.MIMEHeader

	// Unarchive in tmp dir
	if err = Untar(tmpDir, file); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Incorrect tar archive",
		})
		log.Debug(err)
		return
	}

	// Check for index.html in root
	if _, err = os.Stat(tmpDir + "/index.html"); conf.CheckIndexPage && err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "No index.html file in archive",
		})
		log.Error(err)
		return
	}

	// Remove old artifact dir
	if err = os.RemoveAll(conf.ArtifactStoragePath + "/" + body.Artifactname); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal problem, see logs",
		})
		log.Error(err)
		return
	}

	// Replace target artifact dir
	if err = os.Rename(tmpDir, conf.ArtifactStoragePath+"/"+body.Artifactname); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal problem, see logs",
		})
		log.Error(err)
		return
	}
	if err := os.Chmod(conf.ArtifactStoragePath+"/"+body.Artifactname, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal problem, see logs",
		})
		log.Error(err)
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"status": "Artifact has been uploaded",
		"name":   body.Artifactname,
		"url":    "",
	})
}

// Remove method for removing artifacts
func Remove(c *gin.Context) {

	// Get the vars off req body
	var body struct {
		Token        string `form:"token" binding:"required"`
		Artifactname string `form:"artifactname" binding:"required"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Auth
	if body.Token != conf.APISecret { // try bcrypt on browser
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Incorrect token",
		})
		return
	}

	// Validate artifactname
	if !artifactNamePattern.MatchString(body.Artifactname) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect artifactname. ^([A-z0-9-_]{2,64})$",
		})
		return
	}

	// if no dir - return 404

	// Remove artifact dir
	err := os.RemoveAll(conf.ArtifactStoragePath + "/" + body.Artifactname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal problem, see logs",
		})
		log.Error(err)
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"status": "Artifact has been removed",
		"name":   body.Artifactname,
		"url":    "",
	})
}
