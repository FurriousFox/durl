package test

const (
	Pending = -2
	Active  = -1
	Failed  = false
	Success = true
	Skipped = 2
)

type Status struct {
	State any
	Msg   string
}
