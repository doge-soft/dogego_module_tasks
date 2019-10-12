package executers

import (
	"github.com/doge-soft/dogego_module_mutex"
	"github.com/doge-soft/dogego_module_tasks/global"
	"github.com/doge-soft/dogego_module_tasks/managers"
	"github.com/doge-soft/dogego_module_tasks/models"
	"strings"
)

func TaskExecuter(manager *managers.TaskManager, mutex *dogego_module_mutex.RedisMutex) func(message string) error {
	return func(message string) error {
		splits := strings.Split(message, models.Sep)

		rs, err := mutex.RedisClient.Get(global.DataKey(splits[1])).Result()

		if err != nil {
			return err
		}

		go manager.Invoke(splits[0], rs, mutex)

		return nil
	}
}
