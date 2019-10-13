package utils

import (
	"fmt"
	"github.com/doge-soft/dogego_module_mq"
	"github.com/doge-soft/dogego_module_tasks/global"
	"github.com/doge-soft/dogego_module_tasks/models"
	"github.com/google/uuid"
	"reflect"
	"runtime"
)

func PublishTask(
	mq *dogego_module_mq.RedisMQ,
	job func(data models.TaskInput) error, data []byte) error {
	// Save token to Redis
	// Generage Token
	key, err := uuid.NewUUID()

	if err != nil {
		return err
	}

	// Save to Redis
	err = mq.RedisClient.Set(global.DataKey(key.String()), string(data), 0).Err()

	if err != nil {
		return err
	}

	// Publish Message
	err = mq.Publish(XGetenv("TASKS_QUEUE", "tasks"),
		fmt.Sprintf(
			"%s#:#:#%s",
			runtime.FuncForPC(reflect.ValueOf(job).Pointer()).Name(),
			key.String()))

	if err != nil {
		return err
	}

	return nil
}
