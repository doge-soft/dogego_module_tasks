package dogego_module_tasks

import (
	"github.com/doge-soft/dogego_module_mq"
	"github.com/doge-soft/dogego_module_mutex"
	"github.com/doge-soft/dogego_module_tasks/executers"
	"github.com/doge-soft/dogego_module_tasks/managers"
	"github.com/doge-soft/dogego_module_tasks/triggers"
	"github.com/doge-soft/dogego_module_tasks/utils"
	"log"
)

func NewDogeGoTaskModule(mq *dogego_module_mq.RedisMQ,
	mutex *dogego_module_mutex.RedisMutex) (*managers.TaskManager, *triggers.AsyncTrigger, *triggers.TimeTrigger) {
	manager := managers.NewTaskManager()

	err := mq.Custome(utils.XGetenv("TASKS_QUEUE", "tasks"),
		executers.TaskExecuter(manager, mutex))

	if err != nil {
		log.Println(err)
	}

	async_trigger := triggers.NewAsyncTrigger(mq)
	time_trigger := triggers.NewTimeTrigger(mq, manager, mutex)

	return manager, async_trigger, time_trigger
}
