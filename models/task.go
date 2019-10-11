package models

const (
	TimeTrigger = iota
	AsyncTrigger
)

const Sep = "%:#:#:#:%"

type TaskInput interface {
}

type Task struct {
	TaskName    string
	Trigger     int
	CronTrigger string
	UnmarshalFN func(s string) TaskInput
	Job         func(data TaskInput) error
}

type TaskList []Task
