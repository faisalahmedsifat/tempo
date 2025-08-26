package types

type Job struct {
	ID       string
	URL      string
	CronExpr string // "*/10 * * * * *"
}
