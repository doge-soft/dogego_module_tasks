package triggers

import (
	"encoding/json"
	"github.com/doge-soft/dogego_module_mq"
	"github.com/doge-soft/dogego_module_tasks/models"
	"github.com/doge-soft/dogego_module_tasks/utils"
)

type AsyncTrigger struct {
	MQ *dogego_module_mq.RedisMQ
}

func NewAsyncTrigger(mq *dogego_module_mq.RedisMQ) *AsyncTrigger {
	return &AsyncTrigger{
		MQ: mq,
	}
}

func (trigger *AsyncTrigger) TriggerAsyncTask(
	job func(data models.TaskInput) error, data models.TaskInput) error {
	// Serialize body
	result, err := json.Marshal(data)

	if err != nil {
		return err
	}

	// Publish Message
	err = utils.PublishTask(trigger.MQ, job, result)

	if err != nil {
		return err
	}

	return nil
}
