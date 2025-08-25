package types

type Job struct {
	ID       string
	URL      string
	CrosExpr string // "*/10 * * * * *"
}
