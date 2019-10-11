package managers

import (
	"github.com/doge-soft/dogego_module_tasks/models"
	"log"
	"reflect"
	"runtime"
)

type TaskManager struct {
	Tasks models.TaskList
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		Tasks: make(models.TaskList, 0),
	}
}

func (manager *TaskManager) RegisterCronTask(
	cron_trigger string,
	ufn func(s string) models.TaskInput,
	job func(data models.TaskInput) error) {
	manager.Tasks = append(manager.Tasks, models.Task{
		TaskName:    runtime.FuncForPC(reflect.ValueOf(job).Pointer()).Name(),
		Trigger:     models.TimeTrigger,
		CronTrigger: cron_trigger,
		UnmarshalFN: ufn,
		Job:         job,
	})
}

func (manager *TaskManager) RegisterAsyncTask(
	ufn func(s string) models.TaskInput,
	job func(data models.TaskInput) error) {
	manager.Tasks = append(manager.Tasks, models.Task{
		TaskName:    runtime.FuncForPC(reflect.ValueOf(job).Pointer()).Name(),
		Trigger:     models.AsyncTrigger,
		CronTrigger: "",
		UnmarshalFN: ufn,
		Job:         job,
	})
}

func (manager *TaskManager) Trigger(task_name string, input_string string) {
	for _, task := range manager.Tasks {
		if task.TaskName == task_name {
			log.Printf("Job %s be Trigged. Started Execute it.", task.TaskName)

			result := task.UnmarshalFN(input_string)
			err := task.Job(result)
			if err != nil {
				log.Printf("Job %s Execute error %s", task.TaskName, err.Error())
			}

			log.Printf("Job %s Execute finish", task.TaskName)
		}
	}
}
