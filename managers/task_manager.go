package managers

import (
	"github.com/doge-soft/dogego_module_mutex"
	"github.com/doge-soft/dogego_module_tasks/models"
	"log"
	"reflect"
	"runtime"
	"time"
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

func (manager *TaskManager) Invoke(task_name string, input_string string, mutex *dogego_module_mutex.RedisMutex, key string) {
	for _, task := range manager.Tasks {
		if task.TaskName == task_name {
			defer func() {
				if err := recover(); err != nil {
					mutex.UnLock(task.TaskName)
					log.Println(err)
				}
			}()

			result := task.UnmarshalFN(input_string)

			from := time.Now().UnixNano()
			err := task.Job(result)
			to := time.Now().UnixNano()

			if err != nil {
				log.Printf("%s Task Execute Error: %dms\n", task.TaskName, (to-from)/int64(time.Millisecond))
			} else {
				log.Printf("%s Task Execute Success: %dms\n", task.TaskName, (to-from)/int64(time.Millisecond))
			}

			err = mutex.RedisClient.Del(key).Err()

			if err != nil {
				log.Println(err)
			}
		}
	}
}
