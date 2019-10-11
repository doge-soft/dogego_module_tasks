package executers

import (
	"github.com/doge-soft/dogego_module_tasks/global"
	"github.com/doge-soft/dogego_module_tasks/managers"
	"github.com/doge-soft/dogego_module_tasks/models"
	"github.com/go-redis/redis"
	"strings"
)

func TaskExecuter(manager *managers.TaskManager, redis_client *redis.Client) func(message string) error {
	return func(message string) error {
		splits := strings.Split(message, models.Sep)

		rs, err := redis_client.Get(global.DataKey(splits[1])).Result()

		if err != nil {
			return err
		}

		manager.Trigger(splits[0], rs)

		return nil
	}
}
