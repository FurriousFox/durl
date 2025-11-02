package test

const (
	// Active  = -2
	// Running = -2
	Pending = 0
	Success = 1
	Failed  = 2
	// Skipped = 3
)

type Status struct {
	State int
	Msg   string
}
