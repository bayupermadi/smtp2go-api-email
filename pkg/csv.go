package pkg

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"smtp2go-api-email/models"
	"strings"
	"sync"
)

// ProcessResult holds the summary of CSV processing and email sending
type ProcessResult struct {
	Message       string
	Errors        []string
	SuccessEmails int
	FailedEmails  int
}

// ProcessCSVWithWorkers parses the CSV file, validates it, and sends emails using workers
func ProcessCSVWithWorkers(filePath, from, templateID string) *ProcessResult {
	log.Println("Starting CSV processing...")

	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening CSV file:", err)
		return &ProcessResult{
			Message: "Error: Unable to open the uploaded CSV file.",
			Errors:  []string{err.Error()},
		}
	}
	defer file.Close()

	// Read all rows
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Println("Error reading CSV file:", err)
		return &ProcessResult{
			Message: "Error: Unable to read the uploaded CSV file.",
			Errors:  []string{err.Error()},
		}
	}

	// Validate headers
	if len(rows) < 1 {
		return &ProcessResult{
			Message: "Error: CSV file is empty or missing headers.",
			Errors:  []string{"Missing headers"},
		}
	}

	headers := rows[0]
	expectedHeaders := []string{"name", "email"}
	if !validateHeaders(headers, expectedHeaders) {
		return &ProcessResult{
			Message: "Error: CSV headers are incorrect.",
			Errors:  []string{fmt.Sprintf("Expected headers: %v", expectedHeaders)},
		}
	}

	// Channels and worker pool setup
	jobs := make(chan []string, len(rows)-1)
	results := make(chan models.EmailResult, len(rows)-1)
	var wg sync.WaitGroup

	// Start workers
	workerCount := 5
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			EmailWorker(jobs, results, workerID, from, templateID)
		}(i + 1)
	}

	// Send jobs to workers
	for i, row := range rows[1:] {
		if len(row) < 2 {
			log.Printf("Row %d: Missing name or email field\n", i+2)
			continue
		}

		name := strings.TrimSpace(row[0])
		email := strings.TrimSpace(row[1])

		if name == "" || !isValidEmail(email) {
			log.Printf("Row %d: Invalid data - Name: %s, Email: %s\n", i+2, name, email)
			continue
		}

		jobs <- []string{name, email}
	}
	close(jobs)

	// Wait for workers to complete and close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var errors []string
	successCount, failureCount := 0, 0
	for result := range results {
		if result.Success {
			successCount++
		} else {
			failureCount++
			errors = append(errors, result.Error)
		}
	}

	// Return the processing result
	return &ProcessResult{
		Message:       "CSV file processed and emails sent.",
		Errors:        errors,
		SuccessEmails: successCount,
		FailedEmails:  failureCount,
	}
}

// validateHeaders checks if the CSV headers match the expected headers
func validateHeaders(headers, expected []string) bool {
	if len(headers) != len(expected) {
		return false
	}
	for i, header := range headers {
		if strings.TrimSpace(strings.ToLower(header)) != strings.ToLower(expected[i]) {
			return false
		}
	}
	return true
}

// isValidEmail validates the email format using regex
func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(regex).MatchString(email)
}
