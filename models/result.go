package models

// ProcessResult holds the summary of email processing
type ProcessResult struct {
	TotalEmails   int
	SuccessEmails int
	FailedEmails  int
	Errors        []string
}
