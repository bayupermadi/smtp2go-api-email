package pkg

import (
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

var tpl = template.Must(template.ParseFiles("static/index.html"))

// UploadPageHandler renders the upload page
func UploadPageHandler(c echo.Context) error {
	log.Println("Rendering upload page...")
	return tpl.Execute(c.Response().Writer, nil)
}

// UploadHandler processes the uploaded CSV file
func UploadHandler(c echo.Context) error {
	log.Println("Starting UploadHandler...")

	// Parse the form data
	from := c.FormValue("from")
	templateID := c.FormValue("template_id")
	log.Printf("Received form data - From: %s, TemplateID: %s\n", from, templateID)

	// Parse the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		log.Println("Error: File upload failed.")
		return tpl.Execute(c.Response().Writer, map[string]string{
			"Message": "Error: File upload failed.",
		})
	}
	log.Printf("Uploaded file: %s\n", file.Filename)

	src, err := file.Open()
	if err != nil {
		log.Println("Error: Cannot read uploaded file.")
		return tpl.Execute(c.Response().Writer, map[string]string{
			"Message": "Error: Cannot read uploaded file.",
		})
	}
	defer src.Close()

	// Save the file to ./upload directory
	filePath := filepath.Join("./upload", file.Filename)
	log.Printf("Saving file to: %s\n", filePath)

	dst, err := os.Create(filePath)
	if err != nil {
		log.Println("Error: Cannot save file.")
		return tpl.Execute(c.Response().Writer, map[string]string{
			"Message": "Error: Cannot save file.",
		})
	}
	defer dst.Close()

	_, err = dst.ReadFrom(src)
	if err != nil {
		log.Println("Error: Failed to copy file.")
		return tpl.Execute(c.Response().Writer, map[string]string{
			"Message": "Error: Failed to copy file.",
		})
	}
	log.Println("File saved successfully.")

	// Process CSV file with workers
	log.Println("Processing CSV file...")
	result := ProcessCSVWithWorkers(filePath, from, templateID)
	log.Println("CSV processing completed.")

	// Delete the uploaded file after processing
	log.Println("Deleting uploaded file...")
	os.Remove(filePath)

	// Pass success message or errors back to the template
	log.Println("Rendering result to template...")
	return tpl.Execute(c.Response().Writer, result)
}
