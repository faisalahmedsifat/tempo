package types

type Job struct {
	ID       string
	URL      string
	CronExpr string // "*/10 * * * * *"

	Method  string            // "GET", "POST", "PUT", "DELETE"
	Body    string            // "{\"key\": \"value\"}"
	Headers map[string]string // "{\"Content-Type\": \"application/json\"}"
}
