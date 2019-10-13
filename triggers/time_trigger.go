package triggers

import (
	"fmt"
	"github.com/doge-soft/dogego_module_mq"
	"github.com/doge-soft/dogego_module_mutex"
	"github.com/doge-soft/dogego_module_tasks/managers"
	"github.com/doge-soft/dogego_module_tasks/models"
	"github.com/doge-soft/dogego_module_tasks/utils"
	"github.com/robfig/cron"
	"log"
	"time"
)

type TimeTrigger struct {
	MQ      *dogego_module_mq.RedisMQ
	Cron    *cron.Cron
	Manager *managers.TaskManager
	Mutex   *dogego_module_mutex.RedisMutex
}

func NewTimeTrigger(mq *dogego_module_mq.RedisMQ,
	manager *managers.TaskManager,
	mutex *dogego_module_mutex.RedisMutex) *TimeTrigger {
	return &TimeTrigger{
		MQ:      mq,
		Cron:    cron.New(),
		Manager: manager,
		Mutex:   mutex,
	}
}

func (trigger *TimeTrigger) registerTimeTasks() {
	for _, t := range trigger.Manager.Tasks {
		if t.Trigger == models.TimeTrigger {
			func(t models.Task) {
				trigger.Cron.AddFunc(t.CronTrigger, func() {
					utils.PublishTask(trigger.MQ, t.Job, []byte(""))
				})
			}(t)
		}
	}
}

func (trigger *TimeTrigger) ContinueLife() {
	err := trigger.Mutex.RedisClient.Set(
		fmt.Sprintf("lock:%s", "master"), "true", time.Minute*2).Err()

	if err != nil {
		log.Println(err)
		return
	}
}

func (trigger *TimeTrigger) SeizeMaster() {
	trigger.StartCronTrigger(true)
}

func (trigger *TimeTrigger) StartCronTrigger(unlocked bool) {
	if unlocked {
		trigger.Cron.Stop()
		trigger.Cron = cron.New()
	}

	if !trigger.Mutex.Lock("master", time.Minute*2) {
		log.Println("Register master error. In masters running by one.")
		trigger.Cron.AddFunc("@every 2m", trigger.SeizeMaster)
		trigger.Cron.Start()
		return
	}

	trigger.registerTimeTasks()

	trigger.Cron.AddFunc("@every 1m", trigger.ContinueLife)

	trigger.Cron.Start()

	log.Println("Tasks list:")
	for _, i := range trigger.Cron.Entries() {
		log.Println(i)
	}
}
