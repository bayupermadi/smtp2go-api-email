package main

import (
	"os"
	"smtp2go-api-email/pkg"

	"github.com/labstack/echo/v4"
)

func main() {
	// Create upload directory if it doesn't exist
	os.MkdirAll("./upload", os.ModePerm)

	// Initialize Echo
	e := echo.New()

	// Routes
	e.GET("/", pkg.UploadPageHandler)
	e.POST("/upload", pkg.UploadHandler)

	// Serve static files
	e.Static("/static", "static")

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
