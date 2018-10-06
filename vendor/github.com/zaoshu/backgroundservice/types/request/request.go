package request

// TaskRequest message task
type TaskRequest struct {
	Handler string `json:"h"`
	// TaskExpire if TaskExpire == -1, message will never expire until queue delete it.
	Expire        int64  `json:"e"`
	ArgumentsJSON string `json:"a"`
}

// MailArgument simplify mail arguments access
type MailArgument struct {
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	Receivers   []string          `json:"receivers"`
	ReplyTo     []string          `json:"replyTo"`
	Bucket      string            `json:"bucket"`
	Attachments map[string]string `json:"attachments"` // filename, attachment file 3s url
}

// MailSignupEmailArgument signup email argument
type MailSignupEmailArgument struct {
	URL      string `json:"url"`
	Receiver string `json:"receiver"`
}

// MailResetPasswordEmailArgument signup email argument
type MailResetPasswordEmailArgument MailSignupEmailArgument

// MailTaskNoticeEmailArgument wrap sending email task arguments
type MailTaskNoticeEmailArgument struct {
	Status     string   `json:"status"`
	FinishTime string   `json:"time"`
	TaskName   string   `json:"name"`
	TaskURL    string   `json:"url"`
	Receivers  []string `json:"receivers"`
}
