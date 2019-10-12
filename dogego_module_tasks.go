package main

import (
	"encoding/json"
	"github.com/doge-soft/dogego_module_mq"
	"github.com/doge-soft/dogego_module_mutex"
	"github.com/doge-soft/dogego_module_tasks/executers"
	"github.com/doge-soft/dogego_module_tasks/managers"
	"github.com/doge-soft/dogego_module_tasks/models"
	"github.com/doge-soft/dogego_module_tasks/triggers"
	"github.com/go-redis/redis"
	"log"
	"time"
)

type A struct {
	AB int
}

func T1(d models.TaskInput) error {
	log.Println(d.(A))
	panic("yqrx")
	return nil
}

func main() {
	client := redis.NewClient(&redis.Options{})
	mq := dogego_module_mq.NewRedisMQ(client)
	manager := managers.NewTaskManager()
	mutex := dogego_module_mutex.NewRedisMutex(client)

	mq.Custome("tasks", executers.TaskExecuter(manager, mutex))

	manager.RegisterAsyncTask(func(s string) models.TaskInput {
		var r A
		json.Unmarshal([]byte(s), &r)

		return r
	}, T1)

	triggers.NewAsyncTrigger(mq).TriggerAsyncTask(T1, A{AB: 1})

	time.Sleep(time.Second * 10)
}
