package main

import (
	"os"
	"smtp2go-api-email/pkg"

	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/labstack/echo/v4"
)

func main() {
	// Create upload directory if it doesn't exist
	os.MkdirAll("./upload", os.ModePerm)

	// Initialize Echo
	e := echo.New()

	// Embed static files using go.rice
	staticBox := rice.MustFindBox("static")
	assetHandler := http.FileServer(staticBox.HTTPBox())

	// Routes
	e.GET("/", pkg.UploadPageHandler)
	e.POST("/upload", pkg.UploadHandler)

	// Serve embedded static files (CSS, JS, HTML, and CSV)
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
