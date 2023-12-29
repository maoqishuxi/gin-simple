package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

//go:embed build/*
var app embed.FS

func EmbedReact(urlPrefix, buildDirectory string, em embed.FS) gin.HandlerFunc {
	dir := static.LocalFile(buildDirectory, true)
	embedDir, _ := fs.Sub(em, buildDirectory)
	fileserver := http.FileServer(http.FS(embedDir))

	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}

	return func(c *gin.Context) {
		if !dir.Exists(urlPrefix, c.Request.URL.Path) {
			c.Request.URL.Path = "/"
		}
		fileserver.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

func main() {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.Use(CORSMiddleware())
	router.Use(EmbedReact("/", "build", app))

	port := ":8000"
	log.Printf("Listening and serving HTTP on %s\n", port)

	router.Run(port)
}
