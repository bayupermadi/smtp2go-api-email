package pkg

import (
	"fmt"
	"log"
	"smtp2go-api-email/models"

	"github.com/smtp2go-oss/smtp2go-go"
)

// sendEmail sends an email using SMTP2GO with dynamic config
func sendEmail(name, email, from, templateID string) error {
	log.Println("executed", email)
	e := smtp2go.Email{
		From:       from,
		To:         []string{email},
		TemplateID: templateID,
		TemplateData: map[string]string{
			"name": name,
		},
	}

	_, err := smtp2go.Send(&e)
	return err
}

// EmailWorker processes email sending
func EmailWorker(jobs <-chan []string, results chan<- models.EmailResult, workerID int, from, templateID string) {
	for job := range jobs {
		name, email := job[0], job[1]
		err := sendEmail(name, email, from, templateID)
		if err != nil {
			results <- models.EmailResult{Success: false, Error: fmt.Sprintf("Worker %d: Failed to send email to %s", workerID, email)}
		} else {
			results <- models.EmailResult{Success: true}
		}
	}
}
